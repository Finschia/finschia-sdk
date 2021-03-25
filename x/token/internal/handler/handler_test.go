package handler

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link-modules/x/contract"
	testCommon "github.com/line/link-modules/x/token/internal/keeper"
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
	ctx = ctx.WithMultiStore(msCache)
	ctx = ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, defaultContractID))
	return ctx, NewHandler(k)
}

func verifyEventFunc(t *testing.T, expected sdk.Events, actual sdk.Events) {
	require.Equal(t, sdk.StringifyEvents(expected.ToABCIEvents()).String(), sdk.StringifyEvents(actual.ToABCIEvents()).String())
}

const (
	defaultName       = "name"
	defaultContractID = "9be17165"
	defaultSymbol     = "BTC"
	defaultImageURI   = "image-uri"
	defaultMeta       = "{}"
	defaultDecimals   = 6
	defaultAmount     = 1000
	defaultCoin       = "link"
)

var (
	addr1 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	addr2 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
)

func TestHandlerUnrecognized(t *testing.T) {
	ctx, h := cacheKeeper()

	_, err := h(ctx, sdk.NewTestMsg())
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "unrecognized Msg type"))
}
