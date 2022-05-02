package ledger

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	cryptotypes "github.com/line/lbm-sdk/crypto/types"
)

type byter interface {
	Bytes() []byte
}

func checkAminoJSON(t *testing.T, src interface{}, dst interface{}, isNil bool) {
	// Marshal to JSON bytes.
	js, err := cdc.MarshalJSON(src)
	require.Nil(t, err, "%+v", err)
	if isNil {
		require.Equal(t, string(js), `null`)
	} else {
		require.Contains(t, string(js), `"type":`)
		require.Contains(t, string(js), `"value":`)
	}
	// Unmarshal.
	err = cdc.UnmarshalJSON(js, dst)
	require.Nil(t, err, "%+v", err)
}

// nolint: govet
func ExamplePrintRegisteredTypes() {
	cdc.PrintTypes(os.Stdout)
	//| Type | Name | Prefix | Length | Notes |
	//| ---- | ---- | ------ | ----- | ------ |
	//| PrivKeyLedgerSecp256k1 | ostracon/PrivKeyLedgerSecp256k1 | 0x5421414C | variable |  |
	//| PubKey | ostracon/PubKeySr25519 | 0x09EF29BD | variable |  |
	//| PubKey | ostracon/PubKeyEd25519 | 0xCA5F2BB0 | variable |  |
	//| PubKey | ostracon/PubKeySecp256k1 | 0xC03701B7 | variable |  |
	//| LegacyAminoPubKey | ostracon/PubKeyMultisigThreshold | 0x77A72198 | variable |  |
	//| PrivKey | ostracon/PrivKeySr25519 | 0x2C3D3053 | variable |  |
	//| PrivKey | ostracon/PrivKeyEd25519 | 0xF53C89CD | variable |  |
	//| PrivKey | ostracon/PrivKeySecp256k1 | 0x423EB2BA | variable |  |
}

func TestNilEncodings(t *testing.T) {

	// Check nil Signature.
	var a, b []byte
	checkAminoJSON(t, &a, &b, true)
	require.EqualValues(t, a, b)

	// Check nil PubKey.
	var c, d cryptotypes.PubKey
	checkAminoJSON(t, &c, &d, true)
	require.EqualValues(t, c, d)

	// Check nil PrivKey.
	var e, f cryptotypes.PrivKey
	checkAminoJSON(t, &e, &f, true)
	require.EqualValues(t, e, f)

}
