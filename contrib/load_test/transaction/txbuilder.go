package transaction

import (
	crkeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/line/link/contrib/load_test/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

/*
	auth.types.TxBuilder uses a keybase for signing.
	So it is too slow to generate a large number of keys.
	Therefore, defined TxBuilderWithoutKeybase.
*/
type TxBuilderWithoutKeybase struct {
	txBuilder authtypes.TxBuilder
}

func NewTxBuilder(gas uint64) TxBuilderWithoutKeybase {
	return TxBuilderWithoutKeybase{
		txBuilder: authtypes.TxBuilder{},
	}.WithGas(gas)
}

func (bldr TxBuilderWithoutKeybase) WithTxEncoder(txEncoder sdk.TxEncoder) TxBuilderWithoutKeybase {
	bldr.txBuilder = bldr.txBuilder.WithTxEncoder(txEncoder)
	return bldr
}

func (bldr TxBuilderWithoutKeybase) WithChainID(chainID string) TxBuilderWithoutKeybase {
	bldr.txBuilder = bldr.txBuilder.WithChainID(chainID)
	return bldr
}

func (bldr TxBuilderWithoutKeybase) WithGas(gas uint64) TxBuilderWithoutKeybase {
	bldr.txBuilder = bldr.txBuilder.WithGas(gas)
	return bldr
}

func (bldr TxBuilderWithoutKeybase) WithFees(fees string) TxBuilderWithoutKeybase {
	bldr.txBuilder = bldr.txBuilder.WithFees(fees)
	return bldr
}

func (bldr TxBuilderWithoutKeybase) WithGasPrices(gasPrices string) TxBuilderWithoutKeybase {
	bldr.txBuilder = bldr.txBuilder.WithGasPrices(gasPrices)
	return bldr
}

func (bldr TxBuilderWithoutKeybase) WithKeybase(keybase crkeys.Keybase) error {
	return types.InaccessibleFieldError("TxBuilderWithoutKeybase can not access keybase")
}

func (bldr TxBuilderWithoutKeybase) WithSequence(sequence uint64) TxBuilderWithoutKeybase {
	bldr.txBuilder = bldr.txBuilder.WithSequence(sequence)
	return bldr
}

func (bldr TxBuilderWithoutKeybase) WithMemo(memo string) TxBuilderWithoutKeybase {
	bldr.txBuilder = bldr.txBuilder.WithMemo(memo)
	return bldr
}

func (bldr TxBuilderWithoutKeybase) WithAccountNumber(accnum uint64) TxBuilderWithoutKeybase {
	bldr.txBuilder = bldr.txBuilder.WithAccountNumber(accnum)
	return bldr
}

func (bldr TxBuilderWithoutKeybase) BuildAndSign(priv secp256k1.PrivKeySecp256k1, msgs []sdk.Msg) (stdTx authtypes.StdTx, err error) {
	msg, err := bldr.txBuilder.BuildSignMsg(msgs)
	if err != nil {
		return
	}
	stdTx, err = bldr.Sign(priv, msg)
	return
}

func (bldr TxBuilderWithoutKeybase) Sign(priv secp256k1.PrivKeySecp256k1, msg authtypes.StdSignMsg) (stdTx authtypes.StdTx, err error) {
	sig, err := MakeSignature(priv, msg)
	if err != nil {
		return
	}
	return authtypes.NewStdTx(msg.Msgs, msg.Fee, []authtypes.StdSignature{sig}, msg.Memo), nil
}

func MakeSignature(priv secp256k1.PrivKeySecp256k1, msg authtypes.StdSignMsg) (sig authtypes.StdSignature, err error) {
	sigBytes, err := priv.Sign(msg.Bytes())
	if err != nil {
		return
	}
	pubkey := priv.PubKey()

	return auth.StdSignature{
		PubKey:    pubkey,
		Signature: sigBytes,
	}, nil
}
