package signing_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lfb-sdk/crypto/keys/secp256k1"
	sdk "github.com/line/lfb-sdk/types"
	signingtypes "github.com/line/lfb-sdk/types/tx/signing"
	"github.com/line/lfb-sdk/x/auth/legacy/legacytx"
	"github.com/line/lfb-sdk/x/auth/signing"
	banktypes "github.com/line/lfb-sdk/x/bank/types"
)

func MakeTestHandlerMap() signing.SignModeHandler {
	return signing.NewSignModeHandlerMap(
		signingtypes.SignMode_SIGN_MODE_LEGACY_AMINO_JSON,
		[]signing.SignModeHandler{
			legacytx.NewStdTxSignModeHandler(),
		},
	)
}

func TestHandlerMap_GetSignBytes(t *testing.T) {
	priv1 := secp256k1.GenPrivKey()
	addr1 := sdk.AccAddress(priv1.PubKey().Address())
	priv2 := secp256k1.GenPrivKey()
	addr2 := sdk.AccAddress(priv2.PubKey().Address())

	coins := sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}

	fee := legacytx.StdFee{
		Amount: coins,
		Gas:    10000,
	}
	memo := "foo"
	msgs := []sdk.Msg{
		&banktypes.MsgSend{
			FromAddress: addr1.String(),
			ToAddress:   addr2.String(),
			Amount:      coins,
		},
	}

	var (
		chainId        = "test-chain"
		sbh     uint64 = 7
		seqNum  uint64 = 7
	)

	tx := legacytx.StdTx{
		Msgs:           msgs,
		Fee:            fee,
		Signatures:     nil,
		SigBlockHeight: sbh,
		Memo:           memo,
	}

	handler := MakeTestHandlerMap()
	aminoJSONHandler := legacytx.NewStdTxSignModeHandler()

	signingData := signing.SignerData{
		ChainID:        chainId,
		Sequence:       seqNum,
	}
	signBz, err := handler.GetSignBytes(signingtypes.SignMode_SIGN_MODE_LEGACY_AMINO_JSON, signingData, tx)
	require.NoError(t, err)

	expectedSignBz, err := aminoJSONHandler.GetSignBytes(signingtypes.SignMode_SIGN_MODE_LEGACY_AMINO_JSON, signingData, tx)
	require.NoError(t, err)

	require.Equal(t, expectedSignBz, signBz)

	// expect error with wrong sign mode
	_, err = aminoJSONHandler.GetSignBytes(signingtypes.SignMode_SIGN_MODE_DIRECT, signingData, tx)
	require.Error(t, err)
}

func TestHandlerMap_DefaultMode(t *testing.T) {
	handler := MakeTestHandlerMap()
	require.Equal(t, signingtypes.SignMode_SIGN_MODE_LEGACY_AMINO_JSON, handler.DefaultMode())
}

func TestHandlerMap_Modes(t *testing.T) {
	handler := MakeTestHandlerMap()
	modes := handler.Modes()
	require.Contains(t, modes, signingtypes.SignMode_SIGN_MODE_LEGACY_AMINO_JSON)
	require.Len(t, modes, 1)
}
