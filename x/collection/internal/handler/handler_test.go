package handler

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/line/lbm-sdk/store"
	sdk "github.com/line/lbm-sdk/types"
	testCommon "github.com/line/lbm-sdk/x/collection/internal/keeper"
	"github.com/line/lbm-sdk/x/collection/internal/types"
	"github.com/line/lbm-sdk/x/contract"
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

var verifyEventFunc = func(t *testing.T, expected sdk.Events, actual sdk.Events) {
	require.Equal(t, sdk.StringifyEvents(expected.ToABCIEvents()).String(), sdk.StringifyEvents(actual.ToABCIEvents()).String())
}

const (
	defaultContractID = "9be17165"
	defaultName       = "name"
	defaultMeta       = "{}"
	defaultImgURI     = "img-uri"
	defaultDecimals   = 6
	defaultAmount     = 1000
	defaultTokenType  = "10000001"
	defaultTokenType2 = "10000002"
	defaultTokenType3 = "10000003"
	defaultTokenIndex = "00000001"
	defaultTokenID1   = defaultTokenType + defaultTokenIndex
	defaultTokenID2   = defaultTokenType + "00000002"
	defaultTokenID3   = defaultTokenType + "00000003"
	defaultTokenIDFT  = "0000000100000000"
)

var (
	addr1 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	addr2 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
)

func GetMadeContractID(events sdk.Events) string {
	for _, event := range events.ToABCIEvents() {
		for _, attr := range event.Attributes {
			if string(attr.Key) == types.AttributeKeyContractID {
				return string(attr.Value)
			}
		}
	}
	return ""
}

func TestHandlerUnrecognized(t *testing.T) {
	ctx, h := cacheKeeper()
	_, err := h(ctx, sdk.NewTestMsg())
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "unrecognized Msg type"))
}
