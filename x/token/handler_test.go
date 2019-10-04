package token

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	input := setupTestInput(t)
	ctx, keeper := input.ctx, input.keeper

	h := NewHandler(keeper)

	res := h(ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Unrecognized  Msg type"))
	require.False(t, res.Code.IsOK())

	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	msg := MsgPublishToken{
		Name:     "Link Token Description",
		Symbol:   "cony",
		Owner:    addr,
		Amount:   sdk.NewInt(1000),
		Mintable: true,
	}
	res = h(ctx, msg)
	require.True(t, res.Code.IsOK())
}
