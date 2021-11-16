package ante

import (
	sdk "github.com/line/lbm-sdk/types"
	keeper2 "github.com/line/lbm-sdk/x/auth/keeper"
	"github.com/line/lbm-sdk/x/auth/signing"
	"github.com/line/lbm-sdk/x/auth/types"
	"github.com/line/lbm-sdk/x/feegrant/keeper"
	types2 "github.com/line/lbm-sdk/x/feegrant/types"
)

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & sig block height, and deducts fees from the first
// signer.
func NewAnteHandler(
	ak AccountKeeper, bankKeeper types.BankKeeper, feegrantKeeper keeper.Keeper,
	sigGasConsumer SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler,
) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		NewRejectExtensionOptionsDecorator(),
		NewMempoolFeeDecorator(),
		NewValidateBasicDecorator(),
		NewTxSigBlockHeightDecorator(ak),
		TxTimeoutHeightDecorator{},
		NewValidateMemoDecorator(ak),
		NewConsumeGasForTxSizeDecorator(ak),
		NewDeductGrantedFeeDecorator(ak.(keeper2.AccountKeeper), bankKeeper.(types2.BankKeeper), feegrantKeeper),
		// The above handlers should not call `GetAccount` or `GetSignerAcc` for signer
		NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		// The handlers below may call `GetAccount` or `GetSignerAcc` for signer
		NewValidateSigCountDecorator(ak),
		NewSigGasConsumeDecorator(ak, sigGasConsumer),
		NewSigVerificationDecorator(ak, signModeHandler),
		NewIncrementSequenceDecorator(ak),
	)
}
