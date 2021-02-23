package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestTokenApproveKey(t *testing.T) {
	addr1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	addr2 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	contractID1 := "abcdef012"
	contractID2 := "abcdef013"

	require.NotEqual(t, TokenApprovedKey(contractID1, addr1, addr2), TokenApprovedKey(contractID1, addr2, addr1))
	require.NotEqual(t, TokenApprovedKey(contractID1, addr1, addr2), TokenApprovedKey(contractID2, addr1, addr2))
}

func TestTokenApproversKey(t *testing.T) {
	addr1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	addr2 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	contractID := "abcdef012"

	tokenApproversKey := TokenApproversKey(contractID, addr1)

	require.NotEqual(t, tokenApproversKey, TokenApproversKey(contractID, addr2))
	require.Contains(t, string(tokenApproversKey), contractID)
	require.Contains(t, string(tokenApproversKey), string(addr1.Bytes()))
}
