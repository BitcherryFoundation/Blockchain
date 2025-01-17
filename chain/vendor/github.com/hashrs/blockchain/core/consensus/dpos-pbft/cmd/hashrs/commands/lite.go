package commands

import (
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	amino "github.com/hashrs/blockchain/libs/amino"
	dbm "github.com/hashrs/blockchain/libs/state-db/tm-db"

	tmos "github.com/hashrs/blockchain/core/consensus/dpos-pbft/libs/os"
	lite "github.com/hashrs/blockchain/core/consensus/dpos-pbft/lite2"
	httpp "github.com/hashrs/blockchain/core/consensus/dpos-pbft/lite2/provider/http"
	lproxy "github.com/hashrs/blockchain/core/consensus/dpos-pbft/lite2/proxy"
	lrpc "github.com/hashrs/blockchain/core/consensus/dpos-pbft/lite2/rpc"
	dbs "github.com/hashrs/blockchain/core/consensus/dpos-pbft/lite2/store/db"
	rpcclient "github.com/hashrs/blockchain/core/consensus/dpos-pbft/rpc/client"
	rpcserver "github.com/hashrs/blockchain/core/consensus/dpos-pbft/rpc/lib/server"
)

// LiteCmd represents the base command when called without any subcommands
var LiteCmd = &cobra.Command{
	Use:   "lite",
	Short: "Run lite-client proxy server, verifying hashrs rpc",
	Long: `This node will run a secure proxy to a hashrs rpc server.

All calls that can be tracked back to a block header by a proof
will be verified before passing them back to the caller. Other that
that it will present the same interface as a full hashrs node,
just with added trust and running locally.`,
	RunE:         runProxy,
	SilenceUsage: true,
}

var (
	listenAddr         string
	nodeAddr           string
	chainID            string
	home               string
	maxOpenConnections int

	trustingPeriod time.Duration
	trustedHeight  int64
	trustedHash    []byte
)

func init() {
	LiteCmd.Flags().StringVar(&listenAddr, "laddr", "tcp://localhost:8888", "Serve the proxy on the given address")
	LiteCmd.Flags().StringVar(&nodeAddr, "node", "tcp://localhost:26657", "Connect to a HashRs node at this address")
	LiteCmd.Flags().StringVar(&chainID, "chain-id", "hashrs", "Specify the HashRs chain ID")
	LiteCmd.Flags().StringVar(&home, "home-dir", ".hashrs-lite", "Specify the home directory")
	LiteCmd.Flags().IntVar(
		&maxOpenConnections,
		"max-open-connections",
		900,
		"Maximum number of simultaneous connections (including WebSocket).")

	LiteCmd.Flags().DurationVar(&trustingPeriod, "trusting-period", 168*time.Hour, "Trusting period. Should be significantly less than the unbonding period")
	LiteCmd.Flags().Int64Var(&trustedHeight, "trusted-height", 1, "Trusted header's height")
	LiteCmd.Flags().BytesHexVar(&trustedHash, "trusted-hash", []byte{}, "Trusted header's hash")
}

func runProxy(cmd *cobra.Command, args []string) error {
	liteLogger := logger.With("module", "lite")

	logger.Info("Connecting to HashRs node...")
	// First, connect a client
	node, err := rpcclient.NewHTTP(nodeAddr, "/websocket")
	if err != nil {
		return errors.Wrap(err, "new HTTP client")
	}

	logger.Info("Creating client...")
	db, err := dbm.NewGoLevelDB("lite-client-db", home)
	if err != nil {
		return err
	}
	c, err := lite.NewClient(
		chainID,
		lite.TrustOptions{
			Period: trustingPeriod,
			Height: trustedHeight,
			Hash:   trustedHash,
		},
		httpp.NewWithClient(chainID, node),
		dbs.New(db, chainID),
	)
	if err != nil {
		return err
	}
	c.SetLogger(liteLogger)

	p := lproxy.Proxy{
		Addr:   listenAddr,
		Config: &rpcserver.Config{MaxOpenConnections: maxOpenConnections},
		Codec:  amino.NewCodec(),
		Client: lrpc.NewClient(node, c),
		Logger: liteLogger,
	}
	// Stop upon receiving SIGTERM or CTRL-C.
	tmos.TrapSignal(liteLogger, func() {
		p.Listener.Close()
	})

	logger.Info("Starting proxy...")
	if err := p.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		logger.Error("proxy ListenAndServe", "err", err)
	}

	return nil
}
