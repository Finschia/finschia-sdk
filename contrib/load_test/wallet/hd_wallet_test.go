// +build !integration

package wallet

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	TestMnemonic    = "fever tell fancy ridge fly glow reflect decline voice coil reflect ski empty forum frost rebuild slide nut invite chase swarm flag dizzy diet"
	InvalidMnemonic = "invalid mnemonic"
)

func TestNewHDWallet(t *testing.T) {
	config := sdk.GetConfig()
	config.SetCoinType(types.CoinType)

	t.Log("test with valid mnemonic")
	{
		hd, err := NewHDWallet(TestMnemonic)
		require.NoError(t, err)
		require.Equal(t, int(hd.coinType), types.CoinType)
		require.Len(t, hd.masterPrivateKey, 32)
		require.Len(t, hd.masterChainCode, 32)
	}
	t.Log("test with invalid mnemonic")
	{
		_, err := NewHDWallet(InvalidMnemonic)
		require.EqualError(t, err, "Invalid mnemonic")
	}
	t.Log("test with empty mnemonic")
	{
		_, err := NewHDWallet("")
		require.EqualError(t, err, "Invalid mnemonic")
	}
}

func TestHDWallet_GeneratePrivateKey(t *testing.T) {
	config := sdk.GetConfig()
	config.SetCoinType(types.CoinType)
	hd, err := NewHDWallet(TestMnemonic)
	require.NoError(t, err)

	t.Log("test with valid account number and index")
	{
		// Given private key with min values
		minAccountNumber := uint32(0)
		minIndex := uint32(0)
		minPriv, err := hd.GeneratePrivateKey(minAccountNumber, minIndex)

		require.NoError(t, err)
		require.Len(t, minPriv, 32)
		require.Equal(t, "PubKeySecp256k1{02DA1EB422D27EA000C6700FE8B516D360A1F291D1515B8630A46B3D7147916CDB}",
			minPriv.PubKey().(secp256k1.PubKeySecp256k1).String())

		// Given private key with max values
		MaxUint32 := 4294967295
		maxAccountNumber := uint32(MaxUint32)
		maxIndex := uint32(MaxUint32)
		maxPriv, err := hd.GeneratePrivateKey(maxAccountNumber, maxIndex)

		require.NoError(t, err)
		require.Len(t, maxPriv, 32)
		require.Equal(t, "PubKeySecp256k1{036A45349779F1D35439D5A007217AFA739B1023D41BE75AA165E0B400A1C0C3D5}",
			maxPriv.PubKey().(secp256k1.PubKeySecp256k1).String())

		// Then
		require.NotEqual(t, minPriv, maxPriv)
	}
}

func TestHDWallet_GetKeyWallet(t *testing.T) {
	// Given Configs
	config := sdk.GetConfig()
	config.SetCoinType(types.CoinType)
	config.SetBech32PrefixForAccount(types.Bech32PrefixAcc(false), types.Bech32PrefixAccPub(false))
	config.SetBech32PrefixForValidator(types.Bech32PrefixValAddr(false), types.Bech32PrefixValPub(false))
	config.SetBech32PrefixForConsensusNode(types.Bech32PrefixConsAddr(false), types.Bech32PrefixConsPub(false))
	// And HD Wallet
	hd, err := NewHDWallet(TestMnemonic)
	require.NoError(t, err)
	// And account number, index
	accountNumber := uint32(0)
	index := uint32(0)

	// When
	keyWallet, err := hd.GetKeyWallet(accountNumber, index)

	// Then
	require.NoError(t, err)
	require.Len(t, keyWallet.privateKey, 32)
	require.Len(t, keyWallet.address, 20)
	require.Equal(t, keyWallet.address.String(), "link1pzv7nmx9zw04d0huc4jm5492nunxpe4pfrppvt")
}
