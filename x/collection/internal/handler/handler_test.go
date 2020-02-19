package handler

import (
	"os"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	testCommon "github.com/line/link/x/collection/internal/keeper"
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

var verifyEventFunc = func(t *testing.T, expected sdk.Events, actual sdk.Events) {
	require.Equal(t, sdk.StringifyEvents(expected.ToABCIEvents()).String(), sdk.StringifyEvents(actual.ToABCIEvents()).String())
}

const (
	defaultName      = "name"
	defaultSymbol    = "token001"
	defaultTokenURI  = "token-uri"
	defaultDecimals  = 6
	defaultAmount    = 1000
	defaultTokenType = "1001"
	defaultTokenID1  = defaultTokenType + "0001"
	defaultTokenID2  = defaultTokenType + "0002"
	defaultTokenIDFT = "00010000"
)

var (
	addr1 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	addr2 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
)

func TestHandlerUnrecognized(t *testing.T) {
	ctx, h := cacheKeeper()
	res := h(ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Unrecognized  Msg type"))
	require.False(t, res.Code.IsOK())
}
