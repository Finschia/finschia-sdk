package ante

import (
	sdk "github.com/line/lfb-sdk/types"
	"github.com/line/lfb-sdk/x/auth/signing"
	"github.com/line/lfb-sdk/x/auth/types"
)

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(
	ak AccountKeeper, bankKeeper types.BankKeeper,
	sigGasConsumer SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler,
) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		NewRejectExtensionOptionsDecorator(),
		NewMempoolFeeDecorator(),
		NewValidateBasicDecorator(),
		TxTimeoutHeightDecorator{},
		NewValidateMemoDecorator(ak),
		NewConsumeGasForTxSizeDecorator(ak),
		NewRejectFeeGranterDecorator(),
		NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		NewValidateSigCountDecorator(ak),
		NewDeductFeeDecorator(ak, bankKeeper),
		NewSigGasConsumeDecorator(ak, sigGasConsumer),
		NewSigVerificationDecorator(ak, signModeHandler),
		NewIncrementSequenceDecorator(ak),
	)
}
