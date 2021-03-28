package keeper

import (
	"context"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/lbm-sdk/v2/x/contract"
	"github.com/line/lbm-sdk/v2/x/token/internal/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	defaultName       = "name"
	defaultSymbol     = "BTC"
	defaultContractID = "9be17165"
	anotherContractID = "56171eb9"
	defaultMeta       = "{}"
	defaultImageURI   = "image-uri"
	defaultDecimals   = 6
	defaultAmount     = 1000
)

var (
	addr1 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	addr2 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	addr3 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
)

var (
	ms     store.CommitMultiStore
	ctx    sdk.Context
	keeper Keeper
)

func setup() {
	println("setup")
	ctx, ms, keeper = TestKeeper()
}

func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	os.Exit(ret)
}

func cacheKeeper() sdk.Context {
	msCache := ms.CacheMultiStore()
	ctx = ctx.WithMultiStore(msCache)
	ctx = ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, defaultContractID))
	return ctx
}

func verifyTokenFunc(t *testing.T, expected types.Token, actual types.Token) {
	require.Equal(t, expected.GetContractID(), actual.GetContractID())
	require.Equal(t, expected.GetName(), actual.GetName())
	require.Equal(t, expected.GetImageURI(), actual.GetImageURI())
	require.Equal(t, expected.GetDecimals(), actual.GetDecimals())
	require.Equal(t, expected.GetMintable(), actual.GetMintable())
}

func TestKeeper_MarshalJSONLogger(t *testing.T) {
	ctx := cacheKeeper()
	dummy := struct {
		Key   string
		Value string
	}{
		Key:   "key",
		Value: "value",
	}
	bz, err := keeper.MarshalJSON(dummy)
	require.NoError(t, err)

	dummy2 := struct {
		Key   string
		Value string
	}{}

	err = keeper.UnmarshalJSON(bz, &dummy2)
	require.NoError(t, err)
	require.Equal(t, dummy.Key, dummy2.Key)
	require.Equal(t, dummy.Value, dummy2.Value)
	logger := keeper.Logger(ctx)
	logger.Info("test", dummy, dummy2)
}
