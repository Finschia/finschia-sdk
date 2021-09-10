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
	// | Type | Name | Prefix | Length | Notes |
	// | ---- | ---- | ------ | ----- | ------ |
	// | PrivKeyLedgerSecp256k1 | ostracon/PrivKeyLedgerSecp256k1 | 0x10CAB393 | variable |  |
	// | PubKeyBLS12 | ostracon/PubKeyBLS12 | 0xD68FFBC1 | 0x30 | |
	// | PubKey | ostracon/PubKeyEd25519 | 0x1624DE64 | variable |  |
	// | PubKey | ostracon/PubKeySr25519 | 0x0DFB1005 | variable |  |
	// | PubKey | ostracon/PubKeySecp256k1 | 0xEB5AE987 | variable |  |
	// | PubKeyMultisigThreshold | ostracon/PubKeyMultisigThreshold | 0x22C1F7E2 | variable |  |
	// | PubKeyComposite | ostracon/PubKeyComposite | 0x01886E34 | variable |  |
	// | PrivKeyBLS12 | ostracon/PrivKeyBLS12 | 0xEAECF03F | 0x20 |  |
	// | PrivKey | ostracon/PrivKeyEd25519 | 0xA3288910 | variable |  |
	// | PrivKey | ostracon/PrivKeySr25519 | 0x2F82D78B | variable |  |
	// | PrivKey | ostracon/PrivKeySecp256k1 | 0xE1B0F79B | variable |  |
	// | PrivKeyComposite | ostracon/PrivKeyComposite | 0x9F3EE8F0 | variable | |
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
