// +build !integration

package wallet

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestNewKeyWallet(t *testing.T) {
	privateKey := secp256k1.GenPrivKey()

	keyWallet := NewKeyWallet(privateKey)

	require.Equal(t, privateKey, keyWallet.PrivateKey())
	require.Len(t, keyWallet.Address(), 20)
}
