package server

// DONTCOVER

import (
	"fmt"
	"io/ioutil"
	"os"

	tmtypes "github.com/hashrs/blockchain/core/consensus/dpos-pbft/types"
	dbm "github.com/hashrs/blockchain/libs/state-db/tm-db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hashrs/blockchain/framework/chain-app/client/flags"
	"github.com/hashrs/blockchain/framework/chain-app/codec"
	sdk "github.com/hashrs/blockchain/framework/chain-app/types"
)

const (
	flagHeight        = "height"
	flagForZeroHeight = "for-zero-height"
	flagJailWhitelist = "jail-whitelist"
)

// ExportCmd dumps app state to JSON.
func ExportCmd(ctx *Context, cdc *codec.Codec, appExporter AppExporter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export state to JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(flags.FlagHome))

			traceWriterFile := viper.GetString(flagTraceStore)

			db, err := openDB(config.RootDir)
			if err != nil {
				return err
			}

			if isEmptyState(db) || appExporter == nil {
				if _, err := fmt.Fprintln(os.Stderr, "WARNING: State is not initialized. Returning genesis file."); err != nil {
					return err
				}

				genesis, err := ioutil.ReadFile(config.GenesisFile())
				if err != nil {
					return err
				}

				fmt.Println(string(genesis))
				return nil
			}

			traceWriter, err := openTraceWriter(traceWriterFile)
			if err != nil {
				return err
			}

			height := viper.GetInt64(flagHeight)
			forZeroHeight := viper.GetBool(flagForZeroHeight)
			jailWhiteList := viper.GetStringSlice(flagJailWhitelist)

			appState, validators, err := appExporter(ctx.Logger, db, traceWriter, height, forZeroHeight, jailWhiteList)
			if err != nil {
				return fmt.Errorf("error exporting state: %v", err)
			}

			doc, err := tmtypes.GenesisDocFromFile(ctx.Config.GenesisFile())
			if err != nil {
				return err
			}

			doc.AppState = appState
			doc.Validators = validators

			encoded, err := codec.MarshalJSONIndent(cdc, doc)
			if err != nil {
				return err
			}

			fmt.Println(string(sdk.MustSortJSON(encoded)))
			return nil
		},
	}

	cmd.Flags().Int64(flagHeight, -1, "Export state from a particular height (-1 means latest height)")
	cmd.Flags().Bool(flagForZeroHeight, false, "Export state to start at height zero (perform preproccessing)")
	cmd.Flags().StringSlice(flagJailWhitelist, []string{}, "List of validators to not jail state export")
	return cmd
}

func isEmptyState(db dbm.DB) bool {
	return db.Stats()["leveldb.sstables"] == ""
}
