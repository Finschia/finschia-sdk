package keeper

import (
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	defaultName       = "name"
	defaultSymbol     = "token001"
	defaultSymbol2    = "token002"
	defaultTokenURI   = "token-uri"
	defaultDecimals   = 6
	defaultAmount     = 1000
	defaultTokenType  = "1001"
	defaultTokenType2 = "1002"
	defaultTokenID1   = defaultTokenType + "0001"
	defaultTokenID2   = defaultTokenType + "0002"
	defaultTokenID3   = defaultTokenType + "0003"
	defaultTokenID4   = defaultTokenType + "0004"
	defaultTokenID5   = defaultTokenType + "0005"
	defaultTokenID8   = defaultTokenType + "0008"
	defaultTokenIDFT  = "00010000"
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
	return ctx.WithMultiStore(msCache)
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

func prepareCollectionTokens(ctx sdk.Context, t *testing.T) {
	// prepare collection
	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultSymbol, "name"), addr1))

	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultSymbol2, "name"), addr1))

	// issue 6 tokens
	// token1 = symbol1id1 by addr1
	// token2 = symbol1id2 by addr1
	// token3 = symbol1id3 by addr1
	// token4 = symbol1id4 by addr1
	// token5 = symbol2id5 by addr1
	// token6 = symbol1id6 by addr2
	// token7 = symbol1 by addr1
	require.NoError(t, keeper.IssueNFT(ctx, addr1, defaultSymbol))
	require.NoError(t, keeper.IssueNFT(ctx, addr1, defaultSymbol2))
	collection, err := keeper.GetCollection(ctx, defaultSymbol)
	require.NoError(t, err)
	require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(collection, defaultName, defaultTokenType, defaultTokenURI, addr1)))
	collection, err = keeper.GetCollection(ctx, defaultSymbol)
	require.NoError(t, err)
	require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(collection, defaultName, defaultTokenType, defaultTokenURI, addr1)))
	collection, err = keeper.GetCollection(ctx, defaultSymbol)
	require.NoError(t, err)
	require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(collection, defaultName, defaultTokenType, defaultTokenURI, addr1)))
	collection, err = keeper.GetCollection(ctx, defaultSymbol)
	require.NoError(t, err)
	require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(collection, defaultName, defaultTokenType, defaultTokenURI, addr1)))
	collection2, err := keeper.GetCollection(ctx, defaultSymbol2)
	require.NoError(t, err)
	require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(collection2, defaultName, defaultTokenType, defaultTokenURI, addr1)))
	collection, err = keeper.GetCollection(ctx, defaultSymbol)
	require.NoError(t, err)
	require.NoError(t, keeper.GrantPermission(ctx, addr1, addr2, types.NewMintPermission(defaultSymbol+defaultTokenType)))
	require.NoError(t, keeper.MintNFT(ctx, addr2, types.NewNFT(collection, defaultName, defaultTokenType, defaultTokenURI, addr2)))
	require.NoError(t, keeper.IssueFT(ctx, addr1, types.NewFT(collection, defaultName, defaultTokenURI, sdk.NewInt(1), true), sdk.NewInt(defaultAmount)))
}
