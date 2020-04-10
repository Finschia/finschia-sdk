package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	antetypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(ak keeper.AccountKeeper, supplyKeeper antetypes.SupplyKeeper, sigGasConsumer sdkante.SignatureVerificationGasConsumer) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		sdkante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		sdkante.NewMempoolFeeDecorator(),
		sdkante.NewValidateBasicDecorator(),
		sdkante.NewValidateMemoDecorator(ak),
		sdkante.NewConsumeGasForTxSizeDecorator(ak),
		sdkante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		sdkante.NewValidateSigCountDecorator(ak),
		sdkante.NewDeductFeeDecorator(ak, supplyKeeper),
		sdkante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		sdkante.NewSigVerificationDecorator(ak),
		NewIncrementSequenceDecorator(ak), // innermost AnteDecorator
	)
}
