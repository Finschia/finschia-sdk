package keeper

import (
	"context"
	"os"
	"testing"

	"github.com/line/lbm-sdk/store"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/collection/internal/types"
	"github.com/line/lbm-sdk/x/contract"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	defaultName        = "name"
	defaultMeta        = "{}"
	defaultContractID  = "abcdef01"
	defaultContractID2 = "abcdef02"
	wrongContractID    = "abcd1234"
	defaultImgURI      = "img-uri"
	defaultDecimals    = 6
	defaultAmount      = 1000
	defaultTokenType   = "10000001"
	defaultTokenType2  = "10000002"
	defaultTokenType3  = "10000003"
	defaultTokenType4  = "10000004"
	defaultTokenIndex  = "00000001"
	defaultTokenID1    = defaultTokenType + defaultTokenIndex
	defaultTokenID2    = defaultTokenType + "00000002"
	defaultTokenID3    = defaultTokenType + "00000003"
	defaultTokenID4    = defaultTokenType + "00000004"
	defaultTokenID5    = defaultTokenType + "00000005"
	defaultTokenID6    = defaultTokenType + "00000006"
	defaultTokenID7    = defaultTokenType + "00000007"
	defaultTokenID8    = defaultTokenType + "00000008"
	defaultTokenID9    = defaultTokenType + "00000009"
	wrongTokenID       = defaultTokenType2 + "00000001"
	defaultTokenIDFT   = "0000000100000000"
	defaultTokenIDFT2  = "0000000200000000"
	defaultTokenIDFT3  = "0000000300000000"
	defaultTokenIDFT4  = "0000000400000000"
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
	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultContractID, "name", "{}",
		defaultImgURI), addr1))

	ctx2 := ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, defaultContractID2))
	require.NoError(t, keeper.CreateCollection(ctx2, types.NewCollection(defaultContractID2, "name", "{}",
		defaultImgURI), addr1))

	// issue 6 tokens
	// token1 = contract1id1 by addr1
	// token2 = contract1id2 by addr1
	// token3 = contract1id3 by addr1
	// token4 = contract1id4 by addr1
	// token5 = contract2id5 by addr1
	// token6 = contract1id6 by addr2
	// token7 = contract1 by addr1
	require.NoError(t, keeper.IssueNFT(ctx, types.NewBaseTokenType(defaultContractID, defaultTokenType, defaultName, defaultMeta), addr1))
	require.NoError(t, keeper.IssueNFT(ctx2, types.NewBaseTokenType(defaultContractID2, defaultTokenType, defaultName, defaultMeta), addr1))
	require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(defaultContractID, defaultTokenID1, defaultName, defaultMeta, addr1)))
	require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(defaultContractID, defaultTokenID2, defaultName, defaultMeta, addr1)))
	require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(defaultContractID, defaultTokenID3, defaultName, defaultMeta, addr1)))
	require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(defaultContractID, defaultTokenID4, defaultName, defaultMeta, addr1)))
	require.NoError(t, keeper.MintNFT(ctx2, addr1, types.NewNFT(defaultContractID2, defaultTokenID1, defaultName, defaultMeta, addr1)))
	require.NoError(t, keeper.GrantPermission(ctx, addr1, addr2, types.NewMintPermission()))
	require.NoError(t, keeper.MintNFT(ctx, addr2, types.NewNFT(defaultContractID, defaultTokenID5, defaultName, defaultMeta, addr2)))
	require.NoError(t, keeper.IssueFT(ctx, addr1, addr1, types.NewFT(defaultContractID, defaultTokenIDFT, defaultName, defaultMeta, sdk.NewInt(1), true), sdk.NewInt(defaultAmount)))
}

func prepareProxy(ctx sdk.Context, t *testing.T) {
	require.NoError(t, keeper.SetApproved(ctx, addr1, addr2))
	require.NoError(t, keeper.SetApproved(ctx, addr2, addr1))
	require.NoError(t, keeper.TransferFT(ctx, addr1, addr2, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount))))
	require.NoError(t, keeper.TransferNFT(ctx, addr1, addr2, defaultTokenID1))
	require.NoError(t, keeper.TransferNFT(ctx, addr1, addr2, defaultTokenID2))
	require.NoError(t, keeper.TransferNFT(ctx, addr1, addr2, defaultTokenID3))
	require.NoError(t, keeper.TransferNFT(ctx, addr1, addr2, defaultTokenID4))
}

func verifyTokenFunc(t *testing.T, expected types.Token, actual types.Token) {
	switch e := expected.(type) {
	case types.FT:
		a, ok := actual.(types.FT)
		require.True(t, ok)
		require.Equal(t, e.GetContractID(), a.GetContractID())
		require.Equal(t, e.GetName(), a.GetName())
		require.Equal(t, e.GetTokenID(), a.GetTokenID())
		require.Equal(t, e.GetTokenType(), a.GetTokenType())
		require.Equal(t, e.GetTokenIndex(), a.GetTokenIndex())
		require.Equal(t, e.GetDecimals(), a.GetDecimals())
		require.Equal(t, e.GetMintable(), a.GetMintable())
	case types.NFT:
		a, ok := actual.(types.NFT)
		require.True(t, ok)
		require.Equal(t, e.GetContractID(), a.GetContractID())
		require.Equal(t, e.GetName(), a.GetName())
		require.Equal(t, e.GetTokenID(), a.GetTokenID())
		require.Equal(t, e.GetTokenType(), a.GetTokenType())
		require.Equal(t, e.GetTokenIndex(), a.GetTokenIndex())
		require.Equal(t, e.GetOwner(), a.GetOwner())
	default:
		panic("never happen")
	}
}

func verifyTokenTypeFunc(t *testing.T, expected types.TokenType, actual types.TokenType) {
	require.Equal(t, expected.GetName(), actual.GetName())
	require.Equal(t, expected.GetTokenType(), actual.GetTokenType())
}
