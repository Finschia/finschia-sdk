package querier

import (
	"bytes"
	"context"
	"testing"

	"github.com/line/lbm-sdk/codec"
	"github.com/line/lbm-sdk/store"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/collection/internal/keeper"
	"github.com/line/lbm-sdk/x/collection/internal/types"
	"github.com/line/lbm-sdk/x/contract"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	contractID       = "9be17165"
	collectionName   = "mycol"
	imageURL         = "url"
	meta             = "meta"
	tokenFTType      = "00000001"
	tokenFTIndex     = "00000000"
	tokenFTID        = tokenFTType + tokenFTIndex
	tokenFTSupply    = 1000
	tokenNFTType     = "10000001"
	tokenNFTIndex1   = "00000001"
	tokenNFTIndex2   = "00000002"
	tokenNFTIndex3   = "00000003"
	tokenNFTID1      = tokenNFTType + tokenNFTIndex1
	tokenNFTID2      = tokenNFTType + tokenNFTIndex2 /* #nosec */
	tokenNFTID3      = tokenNFTType + tokenNFTIndex3 /* #nosec */
	tokenNFTTypeName = "sword"
	tokenFTName      = "ft_token"
	tokenNFTName1    = "nft_token1" /* #nosec */
	tokenNFTName2    = "nft_token2" /* #nosec */
	tokenNFTName3    = "nft_token3" /* #nosec */
)

var (
	ms      store.CommitMultiStore
	ctx     sdk.Context
	ckeeper keeper.Keeper
	addr1   sdk.AccAddress
	addr2   sdk.AccAddress
	addr3   sdk.AccAddress
)

func prepare(t *testing.T) {
	ctx, ms, ckeeper = keeper.TestKeeper()
	msCache := ms.CacheMultiStore()
	ctx = ctx.WithMultiStore(msCache)

	addr1 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	addr2 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	addr3 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	// prepare contract ID
	newContractID := ckeeper.NewContractID(ctx)
	require.Equal(t, contractID, newContractID)
	ctx2 := ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, contractID))

	// prepare collection
	require.NoError(t, ckeeper.CreateCollection(ctx2, types.NewCollection(contractID, collectionName, meta, imageURL), addr1))
	require.NoError(t, ckeeper.IssueFT(ctx2, addr1, addr1, types.NewFT(contractID, tokenFTID, tokenFTName, meta, sdk.NewInt(1), true), sdk.NewInt(tokenFTSupply)))
	require.NoError(t, ckeeper.IssueNFT(ctx2, types.NewBaseTokenType(contractID, tokenNFTType, tokenNFTTypeName, meta), addr1))
	require.NoError(t, ckeeper.MintNFT(ctx2, addr1, types.NewNFT(contractID, tokenNFTID1, tokenNFTName1, meta, addr1)))
	require.NoError(t, ckeeper.MintNFT(ctx2, addr1, types.NewNFT(contractID, tokenNFTID2, tokenNFTName2, meta, addr1)))
	require.NoError(t, ckeeper.MintNFT(ctx2, addr1, types.NewNFT(contractID, tokenNFTID3, tokenNFTName3, meta, addr1)))

	require.NoError(t, ckeeper.Attach(ctx2, addr1, tokenNFTID1, tokenNFTID2))
	require.NoError(t, ckeeper.Attach(ctx2, addr1, tokenNFTID1, tokenNFTID3))
	require.NoError(t, ckeeper.GrantPermission(ctx2, addr1, addr2, types.NewMintPermission()))
	require.NoError(t, ckeeper.SetApproved(ctx2, addr1, addr2))
	require.NoError(t, ckeeper.SetApproved(ctx2, addr1, addr3))
}

func query(t *testing.T, params interface{}, query string, result interface{}) {
	res, err := queryInternal(params, query, contractID)
	require.NoError(t, err)
	if len(res) > 0 {
		require.NoError(t, ckeeper.UnmarshalJSON(res, result))
	}
}

func queryInternal(params interface{}, query, contractID string) ([]byte, error) {
	req := abci.RequestQuery{
		Path: "",
		Data: []byte(string(codec.MustMarshalJSONIndent(types.ModuleCdc, params))),
	}
	if params == nil {
		req.Data = nil
	}
	path := []string{query}
	if contractID != "" {
		path = append(path, contractID)
	}
	querier := NewQuerier(ckeeper)
	return querier(ctx, path, req)
}

func TestNewQuerier_queryBalance(t *testing.T) {
	prepare(t)
	params := types.QueryTokenIDAccAddressParams{
		TokenID: tokenFTID,
		Addr:    addr1,
	}
	var balance sdk.Int
	query(t, params, types.QueryBalance, &balance)
	require.True(t, balance.Equal(sdk.NewInt(1000)))
}

func TestNewQuerier_queryBalances(t *testing.T) {
	prepare(t)
	params := types.QueryAccAddressParams{
		Addr: addr1,
	}
	var coins types.Coins
	query(t, params, types.QueryBalances, &coins)
	require.Equal(t, tokenFTID, coins[0].Denom)
	require.Equal(t, sdk.NewInt(tokenFTSupply), coins[0].Amount)
	require.Equal(t, tokenNFTID1, coins[1].Denom)
	require.Equal(t, sdk.NewInt(1), coins[1].Amount)
	require.Equal(t, tokenNFTID2, coins[2].Denom)
	require.Equal(t, sdk.NewInt(1), coins[2].Amount)
	require.Equal(t, tokenNFTID3, coins[3].Denom)
	require.Equal(t, sdk.NewInt(1), coins[3].Amount)
	paramsNoExist1 := types.QueryAccAddressParams{
		Addr: addr3,
	}
	var coinsEmpty types.Coins
	query(t, paramsNoExist1, types.QueryBalances, &coinsEmpty)
	require.Empty(t, coinsEmpty)
}

func TestNewQuerier_queryBalanceOwnedNFT(t *testing.T) {
	prepare(t)
	params := types.QueryTokenIDAccAddressParams{
		TokenID: tokenNFTID1,
		Addr:    addr1,
	}
	var balance sdk.Int
	query(t, params, types.QueryBalance, &balance)
	require.True(t, balance.Equal(sdk.NewInt(1)))
}

func TestNewQuerier_queryBalanceNoOwnedNFT(t *testing.T) {
	prepare(t)
	params := types.QueryTokenIDAccAddressParams{
		TokenID: tokenNFTID1,
		Addr:    addr2,
	}
	var balance sdk.Int
	query(t, params, types.QueryBalance, &balance)
	require.True(t, balance.Equal(sdk.NewInt(0)))
}

func TestNewQuerier_queryBalanceNonExistentAccount(t *testing.T) {
	prepare(t)

	params := types.QueryTokenIDAccAddressParams{
		TokenID: tokenFTID,
		Addr:    addr3,
	}
	var balance sdk.Int
	query(t, params, types.QueryBalance, &balance)
	require.True(t, balance.Equal(sdk.NewInt(0)))
}

func TestNewQuerier_queryBalanceNonExistentContractID(t *testing.T) {
	prepare(t)

	contractID := "12345678"
	params := types.QueryTokenIDAccAddressParams{
		TokenID: tokenFTID,
		Addr:    addr1,
	}
	_, err := queryInternal(params, types.QueryBalance, contractID)
	require.Error(t, err, sdkerrors.Wrap(types.ErrCollectionNotExist, contractID))
}

func TestNewQuerier_queryBalanceNonExistentTokenID(t *testing.T) {
	prepare(t)

	tokenID := "00000009" + tokenFTIndex
	params := types.QueryTokenIDAccAddressParams{
		TokenID: tokenID,
		Addr:    addr1,
	}
	_, err := queryInternal(params, types.QueryBalance, contractID)
	require.Error(t, err, sdkerrors.Wrapf(types.ErrCollectionNotExist, "%s %s", contractID, tokenID))
}

func TestNewQuerier_queryAccountPermission(t *testing.T) {
	prepare(t)

	params := types.NewQueryAccAddressParams(addr1)
	var permissions types.Permissions
	query(t, params, types.QueryPerms, &permissions)
	require.Equal(t, len(permissions), 4)
	require.Equal(t, permissions[0].String(), "issue")
	require.Equal(t, permissions[1].String(), "mint")
	require.Equal(t, permissions[2].String(), "burn")
	require.Equal(t, permissions[3].String(), "modify")
}

func TestNewQuerier_queryTokens_FT(t *testing.T) {
	prepare(t)

	params := types.QueryTokenIDParams{
		TokenID: tokenFTID,
	}
	var token types.Token
	query(t, params, types.QueryTokens, &token)
	require.Equal(t, token.GetContractID(), contractID)
	require.Equal(t, token.GetTokenID(), tokenFTID)
	require.Equal(t, token.GetName(), tokenFTName)
	require.Equal(t, token.GetTokenType(), tokenFTType)
	require.Equal(t, token.GetTokenIndex(), tokenFTIndex)
}

func TestNewQuerier_queryTokens_NFT(t *testing.T) {
	prepare(t)

	params := types.QueryTokenIDParams{
		TokenID: tokenNFTID1,
	}
	var token types.Token
	query(t, params, types.QueryTokens, &token)
	require.Equal(t, token.GetContractID(), contractID)
	require.Equal(t, token.GetTokenID(), tokenNFTID1)
	require.Equal(t, token.GetName(), tokenNFTName1)
	require.Equal(t, token.GetTokenType(), tokenNFTType)
	require.Equal(t, token.GetTokenIndex(), tokenNFTIndex1)
}

func TestNewQuerier_queryTokens_all(t *testing.T) {
	prepare(t)

	params := types.QueryTokenIDParams{
		TokenID: "",
	}
	var tokens types.Tokens
	query(t, params, types.QueryTokens, &tokens)
	require.Equal(t, len(tokens), 4)
	require.Equal(t, tokens[0].GetContractID(), contractID)
	require.Equal(t, tokens[0].GetTokenID(), tokenFTID)
	require.Equal(t, tokens[0].GetName(), tokenFTName)
	require.Equal(t, tokens[0].GetTokenType(), tokenFTType)
	require.Equal(t, tokens[0].GetTokenIndex(), tokenFTIndex)
	require.Equal(t, tokens[1].GetContractID(), contractID)
	require.Equal(t, tokens[1].GetTokenID(), tokenNFTID1)
	require.Equal(t, tokens[1].GetName(), tokenNFTName1)
	require.Equal(t, tokens[1].GetTokenType(), tokenNFTType)
	require.Equal(t, tokens[1].GetTokenIndex(), tokenNFTIndex1)
	require.Equal(t, tokens[2].GetContractID(), contractID)
	require.Equal(t, tokens[2].GetTokenID(), tokenNFTID2)
	require.Equal(t, tokens[2].GetName(), tokenNFTName2)
	require.Equal(t, tokens[2].GetTokenType(), tokenNFTType)
	require.Equal(t, tokens[2].GetTokenIndex(), tokenNFTIndex2)
	require.Equal(t, tokens[3].GetContractID(), contractID)
	require.Equal(t, tokens[3].GetTokenID(), tokenNFTID3)
	require.Equal(t, tokens[3].GetName(), tokenNFTName3)
	require.Equal(t, tokens[3].GetTokenType(), tokenNFTType)
	require.Equal(t, tokens[3].GetTokenIndex(), tokenNFTIndex3)
}

func TestNewQuerier_queryTokenTypes_one(t *testing.T) {
	prepare(t)

	params := types.QueryTokenIDParams{
		TokenID: tokenNFTType,
	}
	var tokenType types.TokenType
	query(t, params, types.QueryTokenTypes, &tokenType)
	require.Equal(t, tokenType.GetContractID(), contractID)
	require.Equal(t, tokenType.GetTokenType(), tokenNFTType)
	require.Equal(t, tokenType.GetName(), tokenNFTTypeName)
}

func TestNewQuerier_queryTokenTypes_all(t *testing.T) {
	prepare(t)

	params := types.QueryTokenIDParams{
		TokenID: "",
	}
	var tokenTypes types.TokenTypes
	query(t, params, types.QueryTokenTypes, &tokenTypes)
	require.Equal(t, len(tokenTypes), 1)
	require.Equal(t, tokenTypes[0].GetContractID(), contractID)
	require.Equal(t, tokenTypes[0].GetTokenType(), tokenNFTType)
	require.Equal(t, tokenTypes[0].GetName(), tokenNFTTypeName)
}

func TestNewQuerier_queryTokensWithTokenType(t *testing.T) {
	prepare(t)

	params := types.QueryTokenTypeParams{
		TokenType: tokenNFTType,
	}
	var tokens types.Tokens
	query(t, params, types.QueryTokensWithTokenType, &tokens)
	require.Equal(t, len(tokens), 3)
	require.Equal(t, tokens[0].GetContractID(), contractID)
	require.Equal(t, tokens[0].GetName(), tokenNFTName1)
	require.Equal(t, tokens[0].GetTokenType(), tokenNFTType)
	require.Equal(t, tokens[1].GetContractID(), contractID)
	require.Equal(t, tokens[1].GetName(), tokenNFTName2)
	require.Equal(t, tokens[1].GetTokenType(), tokenNFTType)
	require.Equal(t, tokens[2].GetContractID(), contractID)
	require.Equal(t, tokens[2].GetName(), tokenNFTName3)
	require.Equal(t, tokens[2].GetTokenType(), tokenNFTType)
	params2 := types.QueryTokenTypeParams{
		TokenType: tokenFTType,
	}
	var tokens2 types.Tokens
	query(t, params2, types.QueryTokensWithTokenType, &tokens2)
	require.Equal(t, len(tokens2), 1)
	require.Equal(t, tokens2[0].GetContractID(), contractID)
	require.Equal(t, tokens2[0].GetName(), tokenFTName)
	require.Equal(t, tokens2[0].GetTokenType(), tokenFTType)
	paramsNoExist := types.QueryTokenTypeParams{
		TokenType: "99999999",
	}
	var tokensNoExist types.Tokens
	query(t, paramsNoExist, types.QueryTokensWithTokenType, &tokensNoExist)
	require.Equal(t, len(tokensNoExist), 0)
	require.Empty(t, tokensNoExist)
}

func TestNewQuerier_queryCollections_one(t *testing.T) {
	prepare(t)

	var collection types.Collection
	query(t, nil, types.QueryCollections, &collection)
	require.Equal(t, collection.GetContractID(), contractID)
	require.Equal(t, collection.GetName(), collectionName)
	require.Equal(t, collection.GetBaseImgURI(), imageURL)
}

func TestNewQuerier_queryNFTCount(t *testing.T) {
	prepare(t)

	params := types.QueryTokenIDParams{
		TokenID: tokenNFTType,
	}
	var count sdk.Int
	query(t, params, types.QueryNFTCount, &count)
	require.Equal(t, count, sdk.NewInt(3))
}

func TestNewQuerier_queryTotalSupply_FT(t *testing.T) {
	prepare(t)

	params := types.QueryTokenIDParams{
		TokenID: tokenFTID,
	}
	var supply sdk.Int
	query(t, params, types.QuerySupply, &supply)
	require.Equal(t, supply.Int64(), int64(tokenFTSupply))
}

func TestNewQuerier_queryTotalMint_FT(t *testing.T) {
	prepare(t)

	params := types.QueryTokenIDParams{
		TokenID: tokenFTID,
	}
	var supply sdk.Int
	query(t, params, types.QueryMint, &supply)
	require.Equal(t, supply.Int64(), int64(tokenFTSupply))
}

func TestNewQuerier_queryTotalBurn_FT(t *testing.T) {
	prepare(t)

	params := types.QueryTokenIDParams{
		TokenID: tokenFTID,
	}
	var supply sdk.Int
	query(t, params, types.QueryBurn, &supply)
	require.Equal(t, supply.Int64(), int64(0))
}

func TestNewQuerier_queryTotalSupply_NFT(t *testing.T) {
	prepare(t)

	params := types.QueryTokenIDParams{
		TokenID: tokenNFTType,
	}
	var supply sdk.Int
	query(t, params, types.QueryNFTCount, &supply)
	require.Equal(t, int64(3), supply.Int64())
}

func TestNewQuerier_queryTotalMint_NFT(t *testing.T) {
	prepare(t)

	params := types.QueryTokenIDParams{
		TokenID: tokenNFTType,
	}
	var supply sdk.Int
	query(t, params, types.QueryNFTMint, &supply)
	require.Equal(t, int64(3), supply.Int64())
}

func TestNewQuerier_queryTotalBurn_NFT(t *testing.T) {
	prepare(t)

	params := types.QueryTokenIDParams{
		TokenID: tokenNFTType,
	}
	var supply sdk.Int
	query(t, params, types.QueryNFTBurn, &supply)
	require.Equal(t, supply.Int64(), int64(0))
}

func TestNewQuerier_queryParent(t *testing.T) {
	prepare(t)

	params := types.QueryTokenIDParams{
		TokenID: tokenNFTID2,
	}
	var token types.Token
	query(t, params, types.QueryParent, &token)
	require.Equal(t, token.GetContractID(), contractID)
	require.Equal(t, token.GetTokenID(), tokenNFTID1)
}

func TestNewQuerier_queryParent_nil(t *testing.T) {
	prepare(t)

	params := types.QueryTokenIDParams{
		TokenID: tokenNFTID1,
	}
	var token types.Token
	query(t, params, types.QueryParent, &token)
	require.Equal(t, token, nil)
}

func TestNewQuerier_queryRoot(t *testing.T) {
	prepare(t)

	params := types.QueryTokenIDParams{
		TokenID: tokenNFTID3,
	}
	var token types.Token
	query(t, params, types.QueryRoot, &token)
	require.Equal(t, token.GetContractID(), contractID)
	require.Equal(t, token.GetTokenID(), tokenNFTID1)
}

func TestNewQuerier_queryRoot_self(t *testing.T) {
	prepare(t)

	params := types.QueryTokenIDParams{
		TokenID: tokenNFTID1,
	}
	var token types.Token
	query(t, params, types.QueryRoot, &token)
	require.Equal(t, token.GetContractID(), contractID)
	require.Equal(t, token.GetTokenID(), tokenNFTID1)
}

func TestNewQuerier_queryChildren(t *testing.T) {
	prepare(t)

	params := types.QueryTokenIDParams{
		TokenID: tokenNFTID1,
	}
	var tokens types.Tokens
	query(t, params, types.QueryChildren, &tokens)
	require.Equal(t, len(tokens), 2)
	require.Equal(t, tokens[0].GetContractID(), contractID)
	require.Equal(t, tokens[0].GetTokenID(), tokenNFTID2)
	require.Equal(t, tokens[1].GetContractID(), contractID)
	require.Equal(t, tokens[1].GetTokenID(), tokenNFTID3)
}

func TestNewQuerier_queryChildren_empty(t *testing.T) {
	prepare(t)

	params := types.QueryTokenIDParams{
		TokenID: tokenNFTID2,
	}
	var tokens types.Tokens
	query(t, params, types.QueryChildren, &tokens)
	require.Equal(t, len(tokens), 0)
}

func TestNewQuerier_queryApprovers(t *testing.T) {
	prepare(t)
	params := types.QueryProxyParams{
		Proxy: addr1,
	}
	var acAd1 []sdk.AccAddress
	query(t, params, types.QueryApprovers, &acAd1)
	require.Equal(t, 2, len(acAd1))
	if bytes.Compare(addr2, addr3) <= 0 {
		require.Equal(t, addr2, acAd1[0])
		require.Equal(t, addr3, acAd1[1])
	} else {
		require.Equal(t, addr2, acAd1[1])
		require.Equal(t, addr3, acAd1[0])
	}

	var acAdEmpty []sdk.AccAddress
	paramsEmpty := types.QueryProxyParams{
		Proxy: addr2,
	}
	query(t, paramsEmpty, types.QueryApprovers, &acAdEmpty)
	require.Empty(t, acAdEmpty)
}

func TestNewQuerier_queryIsApproved_true(t *testing.T) {
	prepare(t)
	params := types.QueryIsApprovedParams{
		Proxy:    addr1,
		Approver: addr2,
	}
	var approved bool
	query(t, params, types.QueryIsApproved, &approved)
	require.True(t, approved)
}

func TestNewQuerier_queryIsApproved_false(t *testing.T) {
	prepare(t)

	params := types.QueryIsApprovedParams{
		Proxy:    addr2,
		Approver: addr1,
	}
	var approved bool
	query(t, params, types.QueryIsApproved, &approved)
	require.False(t, approved)
}

func TestNewQuerier_invalid(t *testing.T) {
	prepare(t)
	params := types.QueryTokenIDParams{
		TokenID: tokenNFTID1,
	}
	querier := NewQuerier(ckeeper)
	path := []string{"noquery"}
	req := abci.RequestQuery{
		Path: "",
		Data: []byte(string(codec.MustMarshalJSONIndent(types.ModuleCdc, params))),
	}
	_, err := querier(ctx, path, req)
	require.EqualError(t, err, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown collection query endpoint").Error())
}
