package simapp

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/auth/ante"
	"github.com/line/lbm-sdk/x/auth/signing"
	feegrantkeeper "github.com/line/lbm-sdk/x/feegrant/keeper"
	"github.com/line/lbm-sdk/x/gov/types"
	channelkeeper "github.com/line/lbm-sdk/x/ibc/core/04-channel/keeper"
	ibcante "github.com/line/lbm-sdk/x/ibc/core/ante"
)

func NewAnteHandler(
	ak ante.AccountKeeper,
	bankKeeper types.BankKeeper, //nolint:interfacer
	feegrantKeeper feegrantkeeper.Keeper,
	sigGasConsumer ante.SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler,
	channelKeeper channelkeeper.Keeper,
) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		ante.NewRejectExtensionOptionsDecorator(),
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.TxTimeoutHeightDecorator{},
		ante.NewValidateMemoDecorator(ak),
		ante.NewConsumeGasForTxSizeDecorator(ak),
		ante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(ak),
		// ante.NewDeductFeeDecorator(ak, bankKeeper),
		ante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		ante.NewSigVerificationDecorator(ak, signModeHandler),
		ante.NewIncrementSequenceDecorator(ak),
		ibcante.NewAnteDecorator(channelKeeper),
	)
}
