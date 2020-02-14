package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestKeys(t *testing.T) {
	addr1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	addr2 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	symbol1 := "symbol1"
	symbol2 := "symbol2"

	require.NotEqual(t, CollectionApprovedKey(addr1, addr2, symbol1), CollectionApprovedKey(addr2, addr1, symbol1))
	require.NotEqual(t, CollectionApprovedKey(addr1, addr2, symbol1), CollectionApprovedKey(addr1, addr2, symbol2))
}
