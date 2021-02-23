package types

import (
	"testing"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestKeys(t *testing.T) {
	addr1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	addr2 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	contractID1 := "abcdef012"
	contractID2 := "abcdef013"

	require.NotEqual(t, CollectionApprovedKey(contractID1, addr1, addr2), CollectionApprovedKey(contractID1, addr2, addr1))
	require.NotEqual(t, CollectionApprovedKey(contractID1, addr1, addr2), CollectionApprovedKey(contractID2, addr1, addr2))
}
