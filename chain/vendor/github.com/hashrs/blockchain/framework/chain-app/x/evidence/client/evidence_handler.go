package client

import (
	"github.com/spf13/cobra"

	"github.com/hashrs/blockchain/framework/chain-app/client/context"
	"github.com/hashrs/blockchain/framework/chain-app/codec"
	"github.com/hashrs/blockchain/framework/chain-app/x/evidence/client/rest"
)

type (
	// RESTHandlerFn defines a REST service handler for evidence submission
	RESTHandlerFn func(context.CLIContext) rest.EvidenceRESTHandler

	// CLIHandlerFn defines a CLI command handler for evidence submission
	CLIHandlerFn func(*codec.Codec) *cobra.Command

	// EvidenceHandler defines a type that exposes REST and CLI client handlers for
	// evidence submission.
	EvidenceHandler struct {
		CLIHandler  CLIHandlerFn
		RESTHandler RESTHandlerFn
	}
)

func NewEvidenceHandler(cliHandler CLIHandlerFn, restHandler RESTHandlerFn) EvidenceHandler {
	return EvidenceHandler{
		CLIHandler:  cliHandler,
		RESTHandler: restHandler,
	}
}
