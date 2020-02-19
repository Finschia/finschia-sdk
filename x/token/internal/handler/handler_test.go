package handler

import (
	"os"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	testCommon "github.com/line/link/x/token/internal/keeper"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

var (
	ms  store.CommitMultiStore
	ctx sdk.Context
	k   testCommon.Keeper
)

func setup() {
	println("setup")
	ctx, ms, k = testCommon.TestKeeper()
}

func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	os.Exit(ret)
}

func cacheKeeper() (sdk.Context, sdk.Handler) {
	msCache := ms.CacheMultiStore()
	return ctx.WithMultiStore(msCache), NewHandler(k)
}

func verifyEventFunc(t *testing.T, expected sdk.Events, actual sdk.Events) {
	require.Equal(t, sdk.StringifyEvents(expected.ToABCIEvents()).String(), sdk.StringifyEvents(actual.ToABCIEvents()).String())
}

const (
	symbol     = "token001"
	coinSymbol = "link"
	name       = "description"
	tokenuri   = "tokenuri"
)

var (
	amount   = sdk.NewInt(1000)
	decimals = sdk.NewInt(6)
	addr1    = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	addr2    = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
)

func TestHandlerUnrecognized(t *testing.T) {
	ctx, h := cacheKeeper()

	res := h(ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Unrecognized  Msg type"))
	require.False(t, res.Code.IsOK())
}
