package keyring

import (
	"testing"

	"github.com/stretchr/testify/require"

	kmultisig "github.com/line/lfb-sdk/crypto/keys/multisig"
	"github.com/line/lfb-sdk/crypto/keys/secp256k1"
	"github.com/line/lfb-sdk/crypto/types"
	sdk "github.com/line/lfb-sdk/types"
)

func TestBech32KeysOutput(t *testing.T) {
	tmpKey := secp256k1.GenPrivKey().PubKey()
	bechTmpKey := sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, tmpKey)
	tmpAddr := sdk.BytesToAccAddress(tmpKey.Address())

	multisigPks := kmultisig.NewLegacyAminoPubKey(1, []types.PubKey{tmpKey})
	multiInfo := NewMultiInfo("multisig", multisigPks)
	accAddr := sdk.BytesToAccAddress(multiInfo.GetPubKey().Address())
	bechPubKey := sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, multiInfo.GetPubKey())

	expectedOutput := NewKeyOutput(multiInfo.GetName(), multiInfo.GetType().String(), accAddr.String(), bechPubKey)
	expectedOutput.Threshold = 1
	expectedOutput.PubKeys = []multisigPubKeyOutput{{tmpAddr.String(), bechTmpKey, 1}}

	outputs, err := Bech32KeysOutput([]Info{multiInfo})
	require.NoError(t, err)
	require.Equal(t, expectedOutput, outputs[0])
}
