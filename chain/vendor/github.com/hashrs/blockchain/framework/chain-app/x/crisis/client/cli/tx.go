// nolint
package cli

import (
	"bufio"

	"github.com/spf13/cobra"

	"github.com/hashrs/blockchain/framework/chain-app/client"
	"github.com/hashrs/blockchain/framework/chain-app/client/context"
	"github.com/hashrs/blockchain/framework/chain-app/client/flags"
	"github.com/hashrs/blockchain/framework/chain-app/codec"
	sdk "github.com/hashrs/blockchain/framework/chain-app/types"
	"github.com/hashrs/blockchain/framework/chain-app/x/auth"
	"github.com/hashrs/blockchain/framework/chain-app/x/auth/client/utils"
	"github.com/hashrs/blockchain/framework/chain-app/x/crisis/internal/types"
)

// command to replace a delegator's withdrawal address
func GetCmdInvariantBroken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invariant-broken [module-name] [invariant-route]",
		Short: "submit proof that an invariant broken to halt the chain",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			senderAddr := cliCtx.GetFromAddress()
			moduleName, route := args[0], args[1]
			msg := types.NewMsgVerifyInvariant(senderAddr, moduleName, route)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Crisis transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(flags.PostCommands(
		GetCmdInvariantBroken(cdc),
	)...)
	return txCmd
}
