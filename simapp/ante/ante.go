package ante

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/auth/ante"
	zkauthante "github.com/Finschia/finschia-sdk/x/zkauth/ante"
	zkauthtypes "github.com/Finschia/finschia-sdk/x/zkauth/types"
)

type HandlerOptions struct {
	ante.HandlerOptions

	ZKAuthKeeper zkauthtypes.ZKAuthKeeper
}

func NewAnteHandler(opts HandlerOptions) (sdk.AnteHandler, error) {
	if opts.AccountKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "account keeper is required for ante builder")
	}

	if opts.BankKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "bank keeper is required for ante builder")
	}

	if opts.SignModeHandler == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}

	if opts.ZKAuthKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "zkauth keeper is required for ante builder")
	}

	sigGasConsumer := opts.SigGasConsumer
	if sigGasConsumer == nil {
		sigGasConsumer = ante.DefaultSigVerificationGasConsumer
	}

	anteDecorators := []sdk.AnteDecorator{
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		ante.NewRejectExtensionOptionsDecorator(),
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(opts.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(opts.AccountKeeper),
		zkauthante.NewZKAuthDeductFeeDecorator(opts.AccountKeeper, opts.BankKeeper, opts.FeegrantKeeper),
		zkauthante.NewZKAuthSetPubKeyDecorator(opts.ZKAuthKeeper, opts.AccountKeeper), // replaces NewSetPubKeyDecorator(opts.AccountKeeper)
		ante.NewValidateSigCountDecorator(opts.AccountKeeper),
		zkauthante.NewZKAuthSigGasConsumeDecorator(opts.AccountKeeper, sigGasConsumer),
		zkauthante.NewZKAuthMsgDecorator(opts.ZKAuthKeeper, opts.AccountKeeper, opts.SignModeHandler), // replaces NewSigVerificationDecorator
		zkauthante.NewIncrementSequenceDecorator(opts.AccountKeeper),                                  // replaces NewIncrementSequenceDecorator(opts.AccountKeeper)
	}

	return sdk.ChainAnteDecorators(anteDecorators...), nil
}
