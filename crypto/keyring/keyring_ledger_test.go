//go:build ledger || test_ledger_mock
// +build ledger test_ledger_mock

package keyring

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-sdk/crypto/hd"
	"github.com/Finschia/finschia-sdk/crypto/keys/secp256k1"
	"github.com/Finschia/finschia-sdk/crypto/ledger"
	cryptotypes "github.com/Finschia/finschia-sdk/crypto/types"
	sdk "github.com/Finschia/finschia-sdk/types"
)

func TestInMemoryCreateLedger(t *testing.T) {
	kb := NewInMemory()

	ledger, err := kb.SaveLedgerKey("some_account", hd.Secp256k1, "link", 438, 3, 1)
	if err != nil {
		require.Error(t, err)
		require.Equal(t, "ledger nano S: support for ledger devices is not available in this executable", err.Error())
		require.Nil(t, ledger)
		t.Skip("ledger nano S: support for ledger devices is not available in this executable")
		return
	}

	// The mock is available, check that the address is correct
	pubKey := ledger.GetPubKey()
	expectedPkStr := "PubKeySecp256k1{038F17B38DF1EFC0714D1D3BA0AC1388C32E8B38AD87FD769BAC0B4A11DCE0EBE1}"
	require.Equal(t, expectedPkStr, pubKey.String())

	// Check that restoring the key gets the same results
	restoredKey, err := kb.Key("some_account")
	require.NoError(t, err)
	require.NotNil(t, restoredKey)
	require.Equal(t, "some_account", restoredKey.GetName())
	require.Equal(t, TypeLedger, restoredKey.GetType())
	pubKey = restoredKey.GetPubKey()
	require.Equal(t, expectedPkStr, pubKey.String())

	path, err := restoredKey.GetPath()
	require.NoError(t, err)
	require.Equal(t, "m/44'/438'/3'/0/1", path.String())
}

// TestSignVerify does some detailed checks on how we sign and validate
// signatures
func TestSignVerifyKeyRingWithLedger(t *testing.T) {
	dir := t.TempDir()

	kb, err := New("keybasename", "test", dir, nil)
	require.NoError(t, err)

	i1, err := kb.SaveLedgerKey("key", hd.Secp256k1, "link", 438, 0, 0)
	if err != nil {
		require.Equal(t, "ledger nano S: support for ledger devices is not available in this executable", err.Error())
		t.Skip("ledger nano S: support for ledger devices is not available in this executable")
		return
	}
	require.Equal(t, "key", i1.GetName())

	d1 := []byte("my first message")
	s1, pub1, err := kb.Sign("key", d1)
	require.NoError(t, err)

	s2, pub2, err := SignWithLedger(i1, d1)
	require.NoError(t, err)

	require.True(t, pub1.Equals(pub2))
	require.True(t, bytes.Equal(s1, s2))

	require.Equal(t, i1.GetPubKey(), pub1)
	require.Equal(t, i1.GetPubKey(), pub2)
	require.True(t, pub1.VerifySignature(d1, s1))
	require.True(t, i1.GetPubKey().VerifySignature(d1, s1))
	require.True(t, bytes.Equal(s1, s2))

	localInfo, _, err := kb.NewMnemonic("test", English, sdk.FullFundraiserPath, DefaultBIP39Passphrase, hd.Secp256k1)
	require.NoError(t, err)
	_, _, err = SignWithLedger(localInfo, d1)
	require.Error(t, err)
	require.Equal(t, "not a ledger object", err.Error())
}

func TestAltKeyring_SaveLedgerKey(t *testing.T) {
	dir := t.TempDir()

	keyring, err := New(t.Name(), BackendTest, dir, nil)
	require.NoError(t, err)

	// Test unsupported Algo
	_, err = keyring.SaveLedgerKey("key", notSupportedAlgo{}, "link", 438, 0, 0)
	require.Error(t, err)
	require.Contains(t, err.Error(), ErrUnsupportedSigningAlgo.Error())

	ledger, err := keyring.SaveLedgerKey("some_account", hd.Secp256k1, "link", 438, 3, 1)
	if err != nil {
		require.Equal(t, "ledger nano S: support for ledger devices is not available in this executable", err.Error())
		t.Skip("ledger nano S: support for ledger devices is not available in this executable")
		return
	}

	// The mock is available, check that the address is correct
	require.Equal(t, "some_account", ledger.GetName())
	pubKey := ledger.GetPubKey()
	expectedPkStr := "PubKeySecp256k1{038F17B38DF1EFC0714D1D3BA0AC1388C32E8B38AD87FD769BAC0B4A11DCE0EBE1}"
	require.Equal(t, expectedPkStr, pubKey.String())

	// Check that restoring the key gets the same results
	restoredKey, err := keyring.Key("some_account")
	require.NoError(t, err)
	require.NotNil(t, restoredKey)
	require.Equal(t, "some_account", restoredKey.GetName())
	require.Equal(t, TypeLedger, restoredKey.GetType())
	pubKey = restoredKey.GetPubKey()
	require.Equal(t, expectedPkStr, pubKey.String())

	path, err := restoredKey.GetPath()
	require.NoError(t, err)
	require.Equal(t, "m/44'/438'/3'/0/1", path.String())
}

func TestSignWithLedger(t *testing.T) {
	// Create two distinct Ledger records: infoA and infoB.
	// InfoA is added to the Ledger but infoB is not added.
	pathA := hd.NewFundraiserParams(0, sdk.CoinType, 0)
	privA, _, err := ledger.NewPrivKeySecp256k1(*pathA, "cosmos")
	require.NoError(t, err)
	infoA := newLedgerInfo("ledgerA", privA.PubKey(), *pathA, hd.Secp256k1Type)
	pubA := infoA.GetPubKey()

	pathB := hd.NewFundraiserParams(0, sdk.CoinType, 1)
	// privB won't be added to the Ledger because it doesn't use ledger.NewPrivKeySecp256k1
	privB := secp256k1.GenPrivKey()
	infoB := newLedgerInfo("ledgerB", privB.PubKey(), *pathB, hd.Secp256k1Type)
	require.NoError(t, err)
	pubB := infoB.GetPubKey()

	require.NotEqual(t, pubA, pubB)
	type testCase struct {
		name            string
		info            Info
		msg             []byte
		wantSig         []byte
		wantPub         cryptotypes.PubKey
		wantErr         bool
		wantErrContains string
	}
	testCases := []testCase{
		{
			name:    "ordinary ledger tx",
			info:    infoA,
			msg:     []byte("msg"),
			wantSig: []byte{0xf6, 0xa4, 0x8d, 0x2b, 0x57, 0xd, 0x24, 0x8e, 0x37, 0x66, 0xd9, 0x1b, 0xe9, 0x2b, 0x5, 0x12, 0x48, 0x62, 0x72, 0x8f, 0x87, 0x26, 0x16, 0x36, 0x46, 0x15, 0xff, 0xa4, 0xa7, 0x7, 0xdf, 0x74, 0x20, 0xcd, 0x92, 0x9b, 0x5e, 0xab, 0x7b, 0x63, 0xfb, 0x13, 0x75, 0xa6, 0x74, 0xb4, 0x3d, 0x61, 0x34, 0x33, 0x6, 0x2a, 0x7, 0x59, 0x9b, 0x45, 0x62, 0xc6, 0xb6, 0x59, 0xcc, 0x38, 0xad, 0xbd},
			wantPub: pubA,
			wantErr: false,
		},
		{
			name:            "want error when the public key the user attempted to sign with doesn't match the public key on the ledger",
			info:            infoB,
			msg:             []byte("msg"),
			wantSig:         []byte(nil),
			wantPub:         nil,
			wantErr:         true,
			wantErrContains: "the public key that the user attempted to sign with does not match the public key on the ledger device",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sig, pub, err := SignWithLedger(tc.info, tc.msg)
			assert.Equal(t, tc.wantSig, sig)
			assert.Equal(t, tc.wantPub, pub)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErrContains)
			}
		})
	}
}
