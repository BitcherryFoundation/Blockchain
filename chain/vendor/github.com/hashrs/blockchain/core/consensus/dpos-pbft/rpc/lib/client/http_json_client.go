package rpcclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"

	amino "github.com/hashrs/blockchain/libs/amino"
	"github.com/pkg/errors"

	types "github.com/hashrs/blockchain/core/consensus/dpos-pbft/rpc/lib/types"
)

const (
	protoHTTP  = "http"
	protoHTTPS = "https"
	protoWSS   = "wss"
	protoWS    = "ws"
	protoTCP   = "tcp"
)

//-------------------------------------------------------------

// Parsed URL structure
type parsedURL struct {
	url.URL
}

// Parse URL and set defaults
func newParsedURL(remoteAddr string) (*parsedURL, error) {
	u, err := url.Parse(remoteAddr)
	if err != nil {
		return nil, err
	}

	// default to tcp if nothing specified
	if u.Scheme == "" {
		u.Scheme = protoTCP
	}

	return &parsedURL{*u}, nil
}

// Change protocol to HTTP for unknown protocols and TCP protocol - useful for RPC connections
func (u *parsedURL) SetDefaultSchemeHTTP() {
	// protocol to use for http operations, to support both http and https
	switch u.Scheme {
	case protoHTTP, protoHTTPS, protoWS, protoWSS:
		// known protocols not changed
	default:
		// default to http for unknown protocols (ex. tcp)
		u.Scheme = protoHTTP
	}
}

// Get full address without the protocol - useful for Dialer connections
func (u parsedURL) GetHostWithPath() string {
	// Remove protocol, userinfo and # fragment, assume opaque is empty
	return u.Host + u.EscapedPath()
}

// Get a trimmed address - useful for WS connections
func (u parsedURL) GetTrimmedHostWithPath() string {
	// replace / with . for http requests (kvstore domain)
	return strings.Replace(u.GetHostWithPath(), "/", ".", -1)
}

// Get a trimmed address with protocol - useful as address in RPC connections
func (u parsedURL) GetTrimmedURL() string {
	return u.Scheme + "://" + u.GetTrimmedHostWithPath()
}

//-------------------------------------------------------------

// HTTPClient is a common interface for JSON-RPC HTTP clients.
type HTTPClient interface {
	// Call calls the given method with the params and returns a result.
	Call(method string, params map[string]interface{}, result interface{}) (interface{}, error)
	// Codec returns an amino codec used.
	Codec() *amino.Codec
	// SetCodec sets an amino codec.
	SetCodec(*amino.Codec)
}

// JSONRPCCaller implementers can facilitate calling the JSON-RPC endpoint.
type JSONRPCCaller interface {
	Call(method string, params map[string]interface{}, result interface{}) (interface{}, error)
}

//-------------------------------------------------------------

// JSONRPCClient is a JSON-RPC client, which sends POST HTTP requests to the
// remote server.
//
// Request values are amino encoded. Response is expected to be amino encoded.
// New amino codec is used if no other codec was set using SetCodec.
//
// JSONRPCClient is safe for concurrent use by multiple goroutines.
type JSONRPCClient struct {
	address  string
	username string
	password string

	client *http.Client
	cdc    *amino.Codec

	mtx       sync.Mutex
	nextReqID int
}

var _ HTTPClient = (*JSONRPCClient)(nil)

// Both JSONRPCClient and JSONRPCRequestBatch can facilitate calls to the JSON
// RPC endpoint.
var _ JSONRPCCaller = (*JSONRPCClient)(nil)
var _ JSONRPCCaller = (*JSONRPCRequestBatch)(nil)

// NewJSONRPCClient returns a JSONRPCClient pointed at the given address.
// An error is returned on invalid remote. The function panics when remote is nil.
func NewJSONRPCClient(remote string) (*JSONRPCClient, error) {
	httpClient, err := DefaultHTTPClient(remote)
	if err != nil {
		return nil, err
	}
	return NewJSONRPCClientWithHTTPClient(remote, httpClient)
}

// NewJSONRPCClientWithHTTPClient returns a JSONRPCClient pointed at the given
// address using a custom http client. An error is returned on invalid remote.
// The function panics when remote is nil.
func NewJSONRPCClientWithHTTPClient(remote string, client *http.Client) (*JSONRPCClient, error) {
	if client == nil {
		panic("nil http.Client provided")
	}

	parsedURL, err := newParsedURL(remote)
	if err != nil {
		return nil, fmt.Errorf("invalid remote %s: %s", remote, err)
	}

	parsedURL.SetDefaultSchemeHTTP()

	address := parsedURL.GetTrimmedURL()
	username := parsedURL.User.Username()
	password, _ := parsedURL.User.Password()

	rpcClient := &JSONRPCClient{
		address:  address,
		username: username,
		password: password,
		client:   client,
		cdc:      amino.NewCodec(),
	}

	return rpcClient, nil
}

// Call issues a POST HTTP request. Requests are JSON encoded. Content-Type:
// text/json.
func (c *JSONRPCClient) Call(method string, params map[string]interface{}, result interface{}) (interface{}, error) {
	id := c.nextRequestID()

	request, err := types.MapToRequest(c.cdc, id, method, params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode params")
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request")
	}

	requestBuf := bytes.NewBuffer(requestBytes)
	httpRequest, err := http.NewRequest(http.MethodPost, c.address, requestBuf)
	if err != nil {
		return nil, errors.Wrap(err, "Request failed")
	}
	httpRequest.Header.Set("Content-Type", "text/json")
	if c.username != "" || c.password != "" {
		httpRequest.SetBasicAuth(c.username, c.password)
	}
	httpResponse, err := c.client.Do(httpRequest)
	if err != nil {
		return nil, errors.Wrap(err, "Post failed")
	}
	defer httpResponse.Body.Close() // nolint: errcheck

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	return unmarshalResponseBytes(c.cdc, responseBytes, id, result)
}

func (c *JSONRPCClient) Codec() *amino.Codec       { return c.cdc }
func (c *JSONRPCClient) SetCodec(cdc *amino.Codec) { c.cdc = cdc }

// NewRequestBatch starts a batch of requests for this client.
func (c *JSONRPCClient) NewRequestBatch() *JSONRPCRequestBatch {
	return &JSONRPCRequestBatch{
		requests: make([]*jsonRPCBufferedRequest, 0),
		client:   c,
	}
}

func (c *JSONRPCClient) sendBatch(requests []*jsonRPCBufferedRequest) ([]interface{}, error) {
	reqs := make([]types.RPCRequest, 0, len(requests))
	results := make([]interface{}, 0, len(requests))
	for _, req := range requests {
		reqs = append(reqs, req.request)
		results = append(results, req.result)
	}

	// serialize the array of requests into a single JSON object
	requestBytes, err := json.Marshal(reqs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal requests")
	}

	httpRequest, err := http.NewRequest(http.MethodPost, c.address, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, errors.Wrap(err, "Request failed")
	}
	httpRequest.Header.Set("Content-Type", "text/json")
	if c.username != "" || c.password != "" {
		httpRequest.SetBasicAuth(c.username, c.password)
	}
	httpResponse, err := c.client.Do(httpRequest)
	if err != nil {
		return nil, errors.Wrap(err, "Post failed")
	}
	defer httpResponse.Body.Close() // nolint: errcheck

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	// collect ids to check responses IDs in unmarshalResponseBytesArray
	ids := make([]types.JSONRPCIntID, len(requests))
	for i, req := range requests {
		ids[i] = req.request.ID.(types.JSONRPCIntID)
	}

	return unmarshalResponseBytesArray(c.cdc, responseBytes, ids, results)
}

func (c *JSONRPCClient) nextRequestID() types.JSONRPCIntID {
	c.mtx.Lock()
	id := c.nextReqID
	c.nextReqID++
	c.mtx.Unlock()
	return types.JSONRPCIntID(id)
}

//------------------------------------------------------------------------------------

// jsonRPCBufferedRequest encapsulates a single buffered request, as well as its
// anticipated response structure.
type jsonRPCBufferedRequest struct {
	request types.RPCRequest
	result  interface{} // The result will be deserialized into this object.
}

// JSONRPCRequestBatch allows us to buffer multiple request/response structures
// into a single batch request. Note that this batch acts like a FIFO queue, and
// is thread-safe.
type JSONRPCRequestBatch struct {
	client *JSONRPCClient

	mtx      sync.Mutex
	requests []*jsonRPCBufferedRequest
}

// Count returns the number of enqueued requests waiting to be sent.
func (b *JSONRPCRequestBatch) Count() int {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	return len(b.requests)
}

func (b *JSONRPCRequestBatch) enqueue(req *jsonRPCBufferedRequest) {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	b.requests = append(b.requests, req)
}

// Clear empties out the request batch.
func (b *JSONRPCRequestBatch) Clear() int {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	return b.clear()
}

func (b *JSONRPCRequestBatch) clear() int {
	count := len(b.requests)
	b.requests = make([]*jsonRPCBufferedRequest, 0)
	return count
}

// Send will attempt to send the current batch of enqueued requests, and then
// will clear out the requests once done. On success, this returns the
// deserialized list of results from each of the enqueued requests.
func (b *JSONRPCRequestBatch) Send() ([]interface{}, error) {
	b.mtx.Lock()
	defer func() {
		b.clear()
		b.mtx.Unlock()
	}()
	return b.client.sendBatch(b.requests)
}

// Call enqueues a request to call the given RPC method with the specified
// parameters, in the same way that the `JSONRPCClient.Call` function would.
func (b *JSONRPCRequestBatch) Call(
	method string,
	params map[string]interface{},
	result interface{},
) (interface{}, error) {
	id := b.client.nextRequestID()
	request, err := types.MapToRequest(b.client.cdc, id, method, params)
	if err != nil {
		return nil, err
	}
	b.enqueue(&jsonRPCBufferedRequest{request: request, result: result})
	return result, nil
}

//-------------------------------------------------------------

func makeHTTPDialer(remoteAddr string) (func(string, string) (net.Conn, error), error) {
	u, err := newParsedURL(remoteAddr)
	if err != nil {
		return nil, err
	}

	protocol := u.Scheme

	// accept http(s) as an alias for tcp
	switch protocol {
	case protoHTTP, protoHTTPS:
		protocol = protoTCP
	}

	dialFn := func(proto, addr string) (net.Conn, error) {
		return net.Dial(protocol, u.GetHostWithPath())
	}

	return dialFn, nil
}

// DefaultHTTPClient is used to create an http client with some default parameters.
// We overwrite the http.Client.Dial so we can do http over tcp or unix.
// remoteAddr should be fully featured (eg. with tcp:// or unix://).
// An error will be returned in case of invalid remoteAddr.
func DefaultHTTPClient(remoteAddr string) (*http.Client, error) {
	dialFn, err := makeHTTPDialer(remoteAddr)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			// Set to true to prevent GZIP-bomb DoS attacks
			DisableCompression: true,
			Dial:               dialFn,
		},
	}

	return client, nil
}
