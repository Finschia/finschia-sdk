package ante

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	authsigning "github.com/Finschia/finschia-sdk/x/auth/signing"
	zkauthtypes "github.com/Finschia/finschia-sdk/x/zkauth/types"
)

type ZKAuthMsgDecorator struct {
	zk zkauthtypes.ZKAuthKeeper
}

func NewZKAuthMsgDecorator(zk zkauthtypes.ZKAuthKeeper) ZKAuthMsgDecorator {
	return ZKAuthMsgDecorator{zk: zk}
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

	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	msgs := sigTx.GetMsgs()
	pubKeys, err := sigTx.GetPubKeys()
	if err != nil {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey, "invalid public key, %s", err)
	}

	for i, msg := range msgs {
		if zkMsg, ok := msg.(*zkauthtypes.MsgExecution); ok {
			// verify ZKAuth signature
			if err := zkauthtypes.VerifyZKAuthSignature(ctx, zka.zk, pubKeys[i].Bytes(), zkMsg); err != nil {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "invalid zkauth signature")
			}
		}
	}

	return next(ctx, tx, simulate)
}
