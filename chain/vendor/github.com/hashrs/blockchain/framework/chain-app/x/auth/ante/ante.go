package ante

import (
	sdk "github.com/hashrs/blockchain/framework/chain-app/types"
	"github.com/hashrs/blockchain/framework/chain-app/x/auth/keeper"
	"github.com/hashrs/blockchain/framework/chain-app/x/auth/types"
)

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(ak keeper.AccountKeeper, supplyKeeper types.SupplyKeeper, sigGasConsumer SignatureVerificationGasConsumer) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		NewMempoolFeeDecorator(),
		NewValidateBasicDecorator(),
		NewValidateMemoDecorator(ak),
		NewConsumeGasForTxSizeDecorator(ak),
		NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		NewValidateSigCountDecorator(ak),
		NewDeductFeeDecorator(ak, supplyKeeper),
		NewSigGasConsumeDecorator(ak, sigGasConsumer),
		NewSigVerificationDecorator(ak),
		NewIncrementSequenceDecorator(ak), // innermost AnteDecorator
	)
}
