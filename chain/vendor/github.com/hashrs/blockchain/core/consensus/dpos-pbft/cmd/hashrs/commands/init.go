package commands

import (
	"fmt"

	cfg "github.com/hashrs/blockchain/core/consensus/dpos-pbft/config"
	tmos "github.com/hashrs/blockchain/core/consensus/dpos-pbft/libs/os"
	tmrand "github.com/hashrs/blockchain/core/consensus/dpos-pbft/libs/rand"
	"github.com/hashrs/blockchain/core/consensus/dpos-pbft/p2p"
	"github.com/hashrs/blockchain/core/consensus/dpos-pbft/privval"
	"github.com/hashrs/blockchain/core/consensus/dpos-pbft/types"
	tmtime "github.com/hashrs/blockchain/core/consensus/dpos-pbft/types/time"
	"github.com/spf13/cobra"
)

// InitFilesCmd initialises a fresh HashRs Core instance.
var InitFilesCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize HashRs",
	RunE:  initFiles,
}

func initFiles(cmd *cobra.Command, args []string) error {
	return initFilesWithConfig(config)
}

func initFilesWithConfig(config *cfg.Config) error {
	// private validator
	privValKeyFile := config.PrivValidatorKeyFile()
	privValStateFile := config.PrivValidatorStateFile()
	var pv *privval.FilePV
	if tmos.FileExists(privValKeyFile) {
		pv = privval.LoadFilePV(privValKeyFile, privValStateFile)
		logger.Info("Found private validator", "keyFile", privValKeyFile,
			"stateFile", privValStateFile)
	} else {
		pv = privval.GenFilePV(privValKeyFile, privValStateFile)
		pv.Save()
		logger.Info("Generated private validator", "keyFile", privValKeyFile,
			"stateFile", privValStateFile)
	}

	nodeKeyFile := config.NodeKeyFile()
	if tmos.FileExists(nodeKeyFile) {
		logger.Info("Found node key", "path", nodeKeyFile)
	} else {
		if _, err := p2p.LoadOrGenNodeKey(nodeKeyFile); err != nil {
			return err
		}
		logger.Info("Generated node key", "path", nodeKeyFile)
	}

	// genesis file
	genFile := config.GenesisFile()
	if tmos.FileExists(genFile) {
		logger.Info("Found genesis file", "path", genFile)
	} else {
		genDoc := types.GenesisDoc{
			ChainID:         fmt.Sprintf("test-chain-%v", tmrand.Str(6)),
			GenesisTime:     tmtime.Now(),
			ConsensusParams: types.DefaultConsensusParams(),
		}
		key := pv.GetPubKey()
		genDoc.Validators = []types.GenesisValidator{{
			Address: key.Address(),
			PubKey:  key,
			Power:   10,
		}}

		if err := genDoc.SaveAs(genFile); err != nil {
			return err
		}
		logger.Info("Generated genesis file", "path", genFile)
	}

	return nil
}