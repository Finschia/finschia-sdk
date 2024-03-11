package ante

import (
	"github.com/Finschia/finschia-sdk/crypto/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	authante "github.com/Finschia/finschia-sdk/x/auth/ante"
	authsigning "github.com/Finschia/finschia-sdk/x/auth/signing"
	zkauthtypes "github.com/Finschia/finschia-sdk/x/zkauth/types"
)

type ZKAuthMsgDecorator struct {
	zk              zkauthtypes.ZKAuthKeeper
	ak              authante.AccountKeeper
	signModeHandler authsigning.SignModeHandler
}

func NewZKAuthMsgDecorator(zk zkauthtypes.ZKAuthKeeper, ak authante.AccountKeeper, signModeHandler authsigning.SignModeHandler) ZKAuthMsgDecorator {
	return ZKAuthMsgDecorator{zk: zk, ak: ak, signModeHandler: signModeHandler}
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

	isZKAuthTx, zkMsgs, pubKeys, err := isZKAuthTx(tx)
	if err != nil {
		return ctx, err
	}

	if !isZKAuthTx {
		svd := authante.NewSigVerificationDecorator(zka.ak, zka.signModeHandler)
		return svd.AnteHandle(ctx, tx, simulate, next)
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
	zk zkauthtypes.ZKAuthKeeper
	ak authante.AccountKeeper
}

func NewZKAuthSetPubKeyDecorator(zk zkauthtypes.ZKAuthKeeper, ak authante.AccountKeeper) ZKAuthSetPubKeyDecorator {
	return ZKAuthSetPubKeyDecorator{zk: zk, ak: ak}
}

func (zsp ZKAuthSetPubKeyDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	isZKAuthTx, zkMsgs, _, err := isZKAuthTx(tx)
	if err != nil {
		return ctx, err
	}

	if !isZKAuthTx {
		spk := authante.NewSetPubKeyDecorator(zsp.ak)
		return spk.AnteHandle(ctx, tx, simulate, next)
	}

	for _, zkMsg := range zkMsgs {
		msgs, err := zkMsg.GetMessages()
		if err != nil {
			return ctx, err
		}
		for _, msg := range msgs {
			for _, signer := range msg.GetSigners() {
				zsp.ak.SetAccount(ctx, zsp.ak.NewAccountWithAddress(ctx, signer))
			}
		}
	}

	return next(ctx, tx, simulate)
}

type ZKAuthIncrementSequenceDecorator struct {
	ak authante.AccountKeeper
}

func NewIncrementSequenceDecorator(ak authante.AccountKeeper) ZKAuthIncrementSequenceDecorator {
	return ZKAuthIncrementSequenceDecorator{
		ak: ak,
	}
}

func (zkisd ZKAuthIncrementSequenceDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	isZKAuthTx, zkMsgs, _, err := isZKAuthTx(tx)
	if err != nil {
		return ctx, err
	}

	if !isZKAuthTx {
		isd := authante.NewIncrementSequenceDecorator(zkisd.ak)
		return isd.AnteHandle(ctx, tx, simulate, next)
	}

	for _, zkMsg := range zkMsgs {
		msgs, err := zkMsg.GetMessages()
		if err != nil {
			return ctx, err
		}
		for _, msg := range msgs {
			for _, signer := range msg.GetSigners() {
				acc := zkisd.ak.GetAccount(ctx, signer)
				if err := acc.SetSequence(acc.GetSequence() + 1); err != nil {
					panic(err)
				}

				zkisd.ak.SetAccount(ctx, acc)
			}
		}
	}

	return next(ctx, tx, simulate)
}

func isZKAuthTx(tx sdk.Tx) (bool, []*zkauthtypes.MsgExecution, []types.PubKey, error) {
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
		return false, nil, pubKeys, nil
	}

	return true, zkMsgs, pubKeys, nil
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
