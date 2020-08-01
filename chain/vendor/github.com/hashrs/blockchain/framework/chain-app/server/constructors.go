package server

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	abci "github.com/hashrs/blockchain/core/consensus/dpos-pbft/abci/types"
	"github.com/hashrs/blockchain/core/consensus/dpos-pbft/libs/log"
	tmtypes "github.com/hashrs/blockchain/core/consensus/dpos-pbft/types"
	dbm "github.com/hashrs/blockchain/libs/state-db/tm-db"

	sdk "github.com/hashrs/blockchain/framework/chain-app/types"
)

type (
	// AppCreator is a function that allows us to lazily initialize an
	// application using various configurations.
	AppCreator func(log.Logger, dbm.DB, io.Writer) abci.Application

	// AppExporter is a function that dumps all app state to
	// JSON-serializable structure and returns the current validator set.
	AppExporter func(log.Logger, dbm.DB, io.Writer, int64, bool, []string) (json.RawMessage, []tmtypes.GenesisValidator, error)
)

func openDB(rootDir string) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	db, err := sdk.NewLevelDB("application", dataDir)
	return db, err
}

func openTraceWriter(traceWriterFile string) (w io.Writer, err error) {
	if traceWriterFile != "" {
		w, err = os.OpenFile(
			traceWriterFile,
			os.O_WRONLY|os.O_APPEND|os.O_CREATE,
			0666,
		)
		return
	}
	return
}