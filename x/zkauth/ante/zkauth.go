package ante

import (
	"github.com/Finschia/finschia-sdk/crypto/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	authante "github.com/Finschia/finschia-sdk/x/auth/ante"
	authsigning "github.com/Finschia/finschia-sdk/x/auth/signing"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	zkauthtypes "github.com/Finschia/finschia-sdk/x/zkauth/types"
)

type ZKAuthMsgDecorator struct {
	zk  zkauthtypes.ZKAuthKeeper
	ak  authante.AccountKeeper
	svd *authante.SigVerificationDecorator
}

func NewZKAuthMsgDecorator(zk zkauthtypes.ZKAuthKeeper, ak authante.AccountKeeper, signModeHandler authsigning.SignModeHandler) ZKAuthMsgDecorator {
	return ZKAuthMsgDecorator{
		zk:  zk,
		ak:  ak,
		svd: authante.NewSigVerificationDecorator(ak, signModeHandler),
	}
}

func (zka ZKAuthMsgDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	/*
		todo:
		If there are multiple msg, the order of the pubKey of the signer and signature must be the same.
		- If msg is zkauth, use zkauth signature verification.
		- If msg is a general tx, it is verified by general signature verification.
		(In this implementation, it is assumed that there is only zkauth msg.)

		If the number of msg and the number of pubKey are not the same, how should matching be done?
		Basically, in the case of zkauth msg, ephPubKey must be idempotent for each msg.
	*/

	isZKAuthTx, zkMsgs, pubKeys, err := getZKAuthInfoFromTx(tx)
	if err != nil {
		return ctx, err
	}

	if !isZKAuthTx {
		return zka.svd.AnteHandle(ctx, tx, simulate, next)
	}

	for i, zkMsg := range zkMsgs {
		// verify ZKAuth signature
		if err := zkauthtypes.VerifyZKAuthSignature(ctx, zka.zk, pubKeys[i].Bytes(), zkMsg); err != nil {
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "invalid zkauth signature, %s", err)
		}
	}

	return next(ctx, tx, simulate)
}

type ZKAuthSetPubKeyDecorator struct {
	zk  zkauthtypes.ZKAuthKeeper
	ak  authante.AccountKeeper
	spk authante.SetPubKeyDecorator
}

func NewZKAuthSetPubKeyDecorator(zk zkauthtypes.ZKAuthKeeper, ak authante.AccountKeeper) ZKAuthSetPubKeyDecorator {
	return ZKAuthSetPubKeyDecorator{
		zk:  zk,
		ak:  ak,
		spk: authante.NewSetPubKeyDecorator(ak),
	}
}

func (zsp ZKAuthSetPubKeyDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	isZKAuthTx, zkMsgs, _, err := getZKAuthInfoFromTx(tx)
	if err != nil {
		return ctx, err
	}

	if !isZKAuthTx {
		return zsp.spk.AnteHandle(ctx, tx, simulate, next)
	}

	for _, zkMsg := range zkMsgs {
		for _, signer := range zkMsg.GetSigners() {
			accExists := zsp.ak.HasAccount(ctx, signer)
			if !accExists {
				zsp.ak.SetAccount(ctx, zsp.ak.NewAccountWithAddress(ctx, signer))
			}
		}
	}

	return next(ctx, tx, simulate)
}

type ZKAuthIncrementSequenceDecorator struct {
	ak  authante.AccountKeeper
	isd authante.IncrementSequenceDecorator
}

func NewIncrementSequenceDecorator(ak authante.AccountKeeper) ZKAuthIncrementSequenceDecorator {
	return ZKAuthIncrementSequenceDecorator{
		ak:  ak,
		isd: authante.NewIncrementSequenceDecorator(ak),
	}
}

func (zkisd ZKAuthIncrementSequenceDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	isZKAuthTx, zkMsgs, _, err := getZKAuthInfoFromTx(tx)
	if err != nil {
		return ctx, err
	}

	if !isZKAuthTx {
		return zkisd.isd.AnteHandle(ctx, tx, simulate, next)
	}

	for _, zkMsg := range zkMsgs {
		for _, signer := range zkMsg.GetSigners() {
			acc := zkisd.ak.GetAccount(ctx, signer)
			if err := acc.SetSequence(acc.GetSequence() + 1); err != nil {
				panic(err)
			}
			zkisd.ak.SetAccount(ctx, acc)
		}
	}

	return next(ctx, tx, simulate)
}

type ZKAuthDeductFeeDecorator struct {
	ak             authante.AccountKeeper
	bankKeeper     authtypes.BankKeeper
	feegrantKeeper authante.FeegrantKeeper
	dfd            authante.DeductFeeDecorator
}

func NewZKAuthDeductFeeDecorator(ak authante.AccountKeeper, bankKeeper authtypes.BankKeeper, feegrantKeeper authante.FeegrantKeeper) ZKAuthDeductFeeDecorator {
	return ZKAuthDeductFeeDecorator{
		ak:             ak,
		bankKeeper:     bankKeeper,
		feegrantKeeper: feegrantKeeper,
		dfd:            authante.NewDeductFeeDecorator(ak, bankKeeper, feegrantKeeper),
	}
}

func (zdf ZKAuthDeductFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	isZKAuthTx, _, _, err := getZKAuthInfoFromTx(tx)
	if err != nil {
		return ctx, err
	}

	if !isZKAuthTx {
		return zdf.dfd.AnteHandle(ctx, tx, simulate, next)
	}

	// Case of zkauth msg, does nothing in this case
	return next(ctx, tx, simulate)
}

type ZKAuthSigGasConsumeDecorator struct {
	ak             authante.AccountKeeper
	sigGasConsumer authante.SignatureVerificationGasConsumer
	sgc            authante.SigGasConsumeDecorator
}

func NewZKAuthSigGasConsumeDecorator(ak authante.AccountKeeper, sigGasConsumer authante.SignatureVerificationGasConsumer) ZKAuthSigGasConsumeDecorator {
	return ZKAuthSigGasConsumeDecorator{
		ak:             ak,
		sigGasConsumer: sigGasConsumer,
		sgc:            authante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
	}
}

func (zsg ZKAuthSigGasConsumeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	isZKAuthTx, _, _, err := getZKAuthInfoFromTx(tx)
	if err != nil {
		return ctx, err
	}

	if !isZKAuthTx {
		return zsg.sgc.AnteHandle(ctx, tx, simulate, next)
	}

	// Case of zkauth msg, does nothing in this case
	// TODO: We need an algorithm to deduct fees from zkauth addresses.
	return next(ctx, tx, simulate)
}

func getZKAuthInfoFromTx(tx sdk.Tx) (bool, []*zkauthtypes.MsgExecution, []types.PubKey, error) {
	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return false, nil, nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	pubKeys, err := getPubkeysFromTx(sigTx)
	if err != nil {
		return false, nil, nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey, "invalid public key, %s", err)
	}

	isOnlyMsgExecution, zkMsgs := getMsgExecutionFromTx(sigTx.GetMsgs())
	if !isOnlyMsgExecution {
		return isOnlyMsgExecution, nil, pubKeys, nil
	}

	return isOnlyMsgExecution, zkMsgs, pubKeys, nil
}

func getPubkeysFromTx(sigTx authsigning.SigVerifiableTx) ([]types.PubKey, error) {
	pubKeys, err := sigTx.GetPubKeys()
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey, "invalid public key, %s", err)
	}

	return pubKeys, nil
}

func getMsgExecutionFromTx(msgs []sdk.Msg) (bool, []*zkauthtypes.MsgExecution) {
	// In this implementation, it is assumed that there is only zkauth msg.
	zkMsgs := make([]*zkauthtypes.MsgExecution, 0, len(msgs))
	for _, msg := range msgs {
		zkMsg, ok := msg.(*zkauthtypes.MsgExecution)
		if !ok {
			return false, nil
		}
		zkMsgs = append(zkMsgs, zkMsg)
	}

	return true, zkMsgs
}
