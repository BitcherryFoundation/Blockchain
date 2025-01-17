package starter

import (
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hashrs/blockchain/core/consensus/dpos-pbft/libs/cli"
	amino "github.com/hashrs/blockchain/libs/amino"

	"github.com/hashrs/blockchain/framework/chain-app/client"
	"github.com/hashrs/blockchain/framework/chain-app/client/flags"
	"github.com/hashrs/blockchain/framework/chain-app/client/keys"
	"github.com/hashrs/blockchain/framework/chain-app/client/lcd"
	"github.com/hashrs/blockchain/framework/chain-app/client/rpc"
	"github.com/hashrs/blockchain/framework/chain-app/server"
	sdk "github.com/hashrs/blockchain/framework/chain-app/types"
	"github.com/hashrs/blockchain/framework/chain-app/x/auth"
	authcmd "github.com/hashrs/blockchain/framework/chain-app/x/auth/client/cli"
	bankcmd "github.com/hashrs/blockchain/framework/chain-app/x/bank/client/cli"
	genutilcli "github.com/hashrs/blockchain/framework/chain-app/x/genutil/client/cli"
	"github.com/hashrs/blockchain/framework/chain-app/x/staking"
)

// NewCLICommand returns a basic root CLI cmd to interact with a running SDK chain.
func NewCLICommand(chainName string) *cobra.Command {

	cobra.EnableCommandSorting = false

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	rootCmd := &cobra.Command{
		Use:   strings.ToLower(chainName) + "-cli",
		Short: chainName + " Chain CLI",
	}

	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of node")

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(DefaultCLIHome),
		flags.LineBreak,
		lcd.ServeCommand(Cdc, registerRoutes),
		flags.LineBreak,
		keys.Commands(),
		flags.LineBreak,
	)
	return rootCmd

}

func registerRoutes(rs *lcd.RestServer) {
	client.RegisterRoutes(rs.CliCtx, rs.Mux)
	ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
}

// QueryCmd builds a basic collection of query commands for your SDK CLI tool.
func QueryCmd(cdc *amino.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		flags.LineBreak,
		authcmd.GetAccountCmd(cdc),
		authcmd.QueryTxCmd(cdc),
		authcmd.QueryTxsByEventsCmd(cdc),
	)

	return queryCmd
}

// TxCmd builds a basic collection of transaction commands.
func TxCmd(cdc *amino.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		bankcmd.SendTxCmd(cdc),
		flags.LineBreak,
		authcmd.GetSignCommand(cdc),
		authcmd.GetMultiSignCommand(cdc),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(cdc),
		authcmd.GetEncodeCommand(cdc),
		flags.LineBreak,
	)

	return txCmd
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	if err := viper.BindPFlag(flags.FlagChainID, cmd.PersistentFlags().Lookup(flags.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}

// ServerCommandParams described the params needed to build a basic server CLI command.
type ServerCommandParams struct {
	CmdName     string             // name of the CLI command
	CmdDesc     string             // short description of its function
	AppCreator  server.AppCreator  // method for constructing an ABCI application
	AppExporter server.AppExporter // method for exporting the chain state of an ABCI application
}

// NewServerCommandParams collects the params for a server command
func NewServerCommandParams(name string, desc string, creator server.AppCreator,
	exporter server.AppExporter) ServerCommandParams {
	return ServerCommandParams{name, desc, creator, exporter}
}

// NewServerCommand creates a new ServerCommandParams object
func NewServerCommand(params ServerCommandParams) *cobra.Command {

	cobra.EnableCommandSorting = false

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	ctx := server.NewDefaultContext()

	cdc := MakeCodec()

	rootCmd := &cobra.Command{
		Use:               params.CmdName,
		Short:             params.CmdDesc,
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(
		genutilcli.InitCmd(ctx, cdc, ModuleBasics, DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(ctx, cdc, auth.GenesisAccountIterator{}, DefaultNodeHome),
		genutilcli.GenTxCmd(ctx, cdc, ModuleBasics, staking.AppModuleBasic{},
			auth.GenesisAccountIterator{}, DefaultNodeHome, DefaultCLIHome),
		genutilcli.ValidateGenesisCmd(ctx, cdc, ModuleBasics),
		AddGenesisAccountCmd(ctx, cdc, DefaultNodeHome, DefaultCLIHome),
	)

	server.AddCommands(ctx, cdc, rootCmd, params.AppCreator, params.AppExporter)
	return rootCmd
}
