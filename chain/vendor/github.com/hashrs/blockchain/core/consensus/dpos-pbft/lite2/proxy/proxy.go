package proxy

import (
	"context"
	"net"
	"net/http"

	"github.com/pkg/errors"

	amino "github.com/hashrs/blockchain/libs/amino"

	"github.com/hashrs/blockchain/core/consensus/dpos-pbft/libs/log"
	tmpubsub "github.com/hashrs/blockchain/core/consensus/dpos-pbft/libs/pubsub"
	lrpc "github.com/hashrs/blockchain/core/consensus/dpos-pbft/lite2/rpc"
	ctypes "github.com/hashrs/blockchain/core/consensus/dpos-pbft/rpc/core/types"
	rpcserver "github.com/hashrs/blockchain/core/consensus/dpos-pbft/rpc/lib/server"
)

// A Proxy defines parameters for running an HTTP server proxy.
type Proxy struct {
	Addr     string // TCP address to listen on, ":http" if empty
	Config   *rpcserver.Config
	Codec    *amino.Codec
	Client   *lrpc.Client
	Logger   log.Logger
	Listener net.Listener
}

// ListenAndServe configures the rpcserver.WebsocketManager, sets up the RPC
// routes to proxy via Client, and starts up an HTTP server on the TCP network
// address p.Addr.
// See http#Server#ListenAndServe.
func (p *Proxy) ListenAndServe() error {
	listener, mux, err := p.listen()
	if err != nil {
		return err
	}
	p.Listener = listener

	return rpcserver.StartHTTPServer(
		listener,
		mux,
		p.Logger,
		p.Config,
	)
}

// ListenAndServeTLS acts identically to ListenAndServe, except that it expects
// HTTPS connections.
// See http#Server#ListenAndServeTLS.
func (p *Proxy) ListenAndServeTLS(certFile, keyFile string) error {
	listener, mux, err := p.listen()
	if err != nil {
		return err
	}
	p.Listener = listener

	return rpcserver.StartHTTPAndTLSServer(
		listener,
		mux,
		certFile,
		keyFile,
		p.Logger,
		p.Config,
	)
}

func (p *Proxy) listen() (net.Listener, *http.ServeMux, error) {
	ctypes.RegisterAmino(p.Codec)

	mux := http.NewServeMux()

	// 1) Register regular routes.
	r := RPCRoutes(p.Client)
	rpcserver.RegisterRPCFuncs(mux, r, p.Codec, p.Logger)

	// 2) Allow websocket connections.
	wmLogger := p.Logger.With("protocol", "websocket")
	wm := rpcserver.NewWebsocketManager(r, p.Codec,
		rpcserver.OnDisconnect(func(remoteAddr string) {
			err := p.Client.UnsubscribeAll(context.Background(), remoteAddr)
			if err != nil && err != tmpubsub.ErrSubscriptionNotFound {
				wmLogger.Error("Failed to unsubscribe addr from events", "addr", remoteAddr, "err", err)
			}
		}),
		rpcserver.ReadLimit(p.Config.MaxBodyBytes),
	)
	wm.SetLogger(wmLogger)
	mux.HandleFunc("/websocket", wm.WebsocketHandler)

	// 3) Start a client.
	if !p.Client.IsRunning() {
		if err := p.Client.Start(); err != nil {
			return nil, mux, errors.Wrap(err, "Client#Start")
		}
	}

	// 4) Start listening for new connections.
	listener, err := rpcserver.Listen(p.Addr, p.Config)
	if err != nil {
		return nil, mux, err
	}

	return listener, mux, nil
}
