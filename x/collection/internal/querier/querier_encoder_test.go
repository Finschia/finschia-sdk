package querier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/lbm-sdk/v2/x/collection/internal/types"
	"github.com/line/lbm-sdk/v2/x/wasm"
	"github.com/stretchr/testify/require"
)

var (
	collectionQueryEncoder wasm.EncodeQuerier
)

func setupQueryEncoder() {
	collectionQuerier := NewQuerier(ckeeper)

	collectionQueryEncoder = NewQueryEncoder(collectionQuerier)
}

func encodeQuery(t *testing.T, jsonQuerier json.RawMessage, result interface{}) error {
	res, err := collectionQueryEncoder(ctx, jsonQuerier)
	if len(res) > 0 {
		require.NoError(t, ckeeper.UnmarshalJSON(res, result))
	}
	return err
}

func TestNewQuerier_encodeQueryBalance(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"balance","data":{"balance_param":{"contract_id":"%s", "token_id":"%s", "addr":"%s"}}}`, contractID, tokenFTID, addr1)

	var balance sdk.Int
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &balance)
	require.NoError(t, err)
	require.True(t, balance.Equal(sdk.NewInt(1000)))
}

func TestNewQuerier_encodeQueryBalances(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"balances","data":{"balance_param":{"contract_id":"%s", "token_id":"%s", "addr":"%s"}}}`, contractID, "", addr1)

	var coins types.Coins
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &coins)
	require.NoError(t, err)
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

func TestNewQuerier_encodeQueryBalanceNonExistentAccount(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"balance","data":{"balance_param":{"contract_id":"%s", "token_id":"%s", "addr":"%s"}}}`, contractID, tokenFTID, addr3)

	var balance sdk.Int
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &balance)
	require.NoError(t, err)
	require.True(t, balance.Equal(sdk.NewInt(0)))
}

func TestNewQuerier_encodeQueryBalanceNonExistentContractID(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	contractID := "12345678"
	jsonQuerier := fmt.Sprintf(`{"route":"balance","data":{"balance_param":{"contract_id":"%s", "token_id":"%s", "addr":"%s"}}}`, contractID, tokenFTID, addr1)

	var balance sdk.Int
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &balance)
	require.Error(t, err, sdkerrors.Wrap(types.ErrCollectionNotExist, contractID))
}

func TestNewQuerier_encodeQueryBalanceNonExistentTokenID(t *testing.T) {
	prepare(t)

	tokenID := "00000009" + tokenFTIndex
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"balance","data":{"balance_param":{"contract_id":"%s", "token_id":"%s", "addr":"%s"}}}`, contractID, tokenID, addr1)

	var balance sdk.Int
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &balance)
	require.Error(t, err, sdkerrors.Wrapf(types.ErrCollectionNotExist, "%s %s", contractID, tokenID))
}

func TestNewQuerier_encodeQueryAccountPermission(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"perms","data":{"perm_param":{"contract_id":"%s", "address":"%s"}}}`, contractID, addr1)

	var permissions types.Permissions
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &permissions)
	require.NoError(t, err)

	require.Equal(t, len(permissions), 4)
	require.Equal(t, permissions[0].String(), "issue")
	require.Equal(t, permissions[1].String(), "mint")
	require.Equal(t, permissions[2].String(), "burn")
	require.Equal(t, permissions[3].String(), "modify")
}

func TestNewQuerier_encodeQueryTokens_FT(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"tokens","data":{"tokens_param":{"contract_id":"%s", "token_id":"%s"}}}`, contractID, tokenFTID)

	var token types.Token
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &token)
	require.NoError(t, err)
	require.Equal(t, token.GetContractID(), contractID)
	require.Equal(t, token.GetTokenID(), tokenFTID)
	require.Equal(t, token.GetName(), tokenFTName)
	require.Equal(t, token.GetTokenType(), tokenFTType)
	require.Equal(t, token.GetTokenIndex(), tokenFTIndex)
}

func TestNewQuerier_encodeQueryTokens_NFT(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"tokens","data":{"tokens_param":{"contract_id":"%s", "token_id":"%s"}}}`, contractID, tokenNFTID1)

	var token types.Token
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &token)
	require.NoError(t, err)
	require.Equal(t, token.GetContractID(), contractID)
	require.Equal(t, token.GetTokenID(), tokenNFTID1)
	require.Equal(t, token.GetName(), tokenNFTName1)
	require.Equal(t, token.GetTokenType(), tokenNFTType)
	require.Equal(t, token.GetTokenIndex(), tokenNFTIndex1)
}

func TestNewQuerier_encodeQueryTokens_all(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"tokens","data":{"tokens_param":{"contract_id":"%s", "token_id":""}}}`, contractID)

	var tokens types.Tokens
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &tokens)
	require.NoError(t, err)
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

func TestNewQuerier_encodeQueryTokenTypes_one(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"tokentypes","data":{"tokentypes_param":{"contract_id":"%s", "token_id":"%s"}}}`, contractID, tokenNFTType)

	var tokenType types.TokenType
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &tokenType)
	require.NoError(t, err)
	require.Equal(t, tokenType.GetContractID(), contractID)
	require.Equal(t, tokenType.GetTokenType(), tokenNFTType)
	require.Equal(t, tokenType.GetName(), tokenNFTTypeName)
}

func TestNewQuerier_encodeQueryTokenTypes_all(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"tokentypes","data":{"tokentypes_param":{"contract_id":"%s", "token_id":""}}}`, contractID)

	var tokenTypes types.TokenTypes
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &tokenTypes)
	require.NoError(t, err)
	require.Equal(t, len(tokenTypes), 1)
	require.Equal(t, tokenTypes[0].GetContractID(), contractID)
	require.Equal(t, tokenTypes[0].GetTokenType(), tokenNFTType)
	require.Equal(t, tokenTypes[0].GetName(), tokenNFTTypeName)
}

func TestNewQuerier_encodeQueryTokensWithTokenType(t *testing.T) {
	prepare(t)
	setupQueryEncoder()

	jsonQuerier := fmt.Sprintf(`{"route":"tokensWithTokenType","data":{"token_type_param":{"contract_id":"%s", "token_type":"%s"}}}`, contractID, tokenNFTType)
	var tokens types.Tokens
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &tokens)
	require.NoError(t, err)
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

	var tokens2 types.Tokens
	jsonQuerier = fmt.Sprintf(`{"route":"tokensWithTokenType","data":{"token_type_param":{"contract_id":"%s", "token_type":"%s"}}}`, contractID, tokenFTType)
	err = encodeQuery(t, json.RawMessage(jsonQuerier), &tokens2)
	require.NoError(t, err)
	require.Equal(t, len(tokens2), 1)
	require.Equal(t, tokens2[0].GetContractID(), contractID)
	require.Equal(t, tokens2[0].GetName(), tokenFTName)
	require.Equal(t, tokens2[0].GetTokenType(), tokenFTType)

	tokenType := "99999999"
	var tokensNoExist types.Tokens
	jsonQuerier = fmt.Sprintf(`{"route":"tokensWithTokenType","data":{"token_type_param":{"contract_id":"%s", "token_type":"%s"}}}`, contractID, tokenType)
	err = encodeQuery(t, json.RawMessage(jsonQuerier), &tokensNoExist)
	require.NoError(t, err)
	require.Equal(t, len(tokensNoExist), 0)
	require.Empty(t, tokensNoExist)
}

func TestNewQuerier_encodeQueryCollections_one(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"collections","data":{"collection_param":{"contract_id":"%s"}}}`, contractID)

	var collection types.Collection
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &collection)
	require.NoError(t, err)
	require.Equal(t, collection.GetContractID(), contractID)
	require.Equal(t, collection.GetName(), collectionName)
	require.Equal(t, collection.GetBaseImgURI(), imageURL)
}

func TestNewQuerier_encodeQueryNFTCount(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"nftcount","data":{"tokens_param":{"contract_id":"%s", "token_id":"%s"}}}`, contractID, tokenNFTType)

	var count sdk.Int
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &count)
	require.NoError(t, err)
	require.Equal(t, count, sdk.NewInt(3))
}

func TestNewQuerier_encodeQueryTotalMint_NFT(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"nftmint","data":{"tokens_param":{"contract_id":"%s", "token_id":"%s"}}}`, contractID, tokenNFTType)

	var supply sdk.Int
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &supply)
	require.NoError(t, err)
	require.Equal(t, int64(3), supply.Int64())
}

func TestNewQuerier_encodeQueryTotalBurn_NFT(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"nftburn","data":{"tokens_param":{"contract_id":"%s", "token_id":"%s"}}}`, contractID, tokenNFTType)

	var supply sdk.Int
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &supply)
	require.NoError(t, err)
	require.Equal(t, supply.Int64(), int64(0))
}

func TestNewQuerier_encodeQueryTotalSupply_FT(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"%s","data":{"total_param":{"contract_id":"%s", "token_id":"%s"}}}`, types.QuerySupply, contractID, tokenFTID)

	var supply sdk.Int
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &supply)
	require.NoError(t, err)
	require.Equal(t, supply.Int64(), int64(tokenFTSupply))
}

func TestNewQuerier_encodeQueryTotalMint_FT(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"%s","data":{"total_param":{"contract_id":"%s", "token_id":"%s"}}}`, types.QueryMint, contractID, tokenFTID)

	var supply sdk.Int
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &supply)
	require.NoError(t, err)
	require.Equal(t, supply.Int64(), int64(tokenFTSupply))
}

func TestNewQuerier_encodeQueryTotalBurn_FT(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"%s","data":{"total_param":{"contract_id":"%s", "token_id":"%s"}}}`, types.QueryBurn, contractID, tokenFTID)

	var supply sdk.Int
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &supply)
	require.NoError(t, err)
	require.Equal(t, supply.Int64(), int64(0))
}

func TestNewQuerier_encodeQueryParent(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"parent","data":{"tokens_param":{"contract_id":"%s", "token_id":"%s"}}}`, contractID, tokenNFTID2)

	var token types.Token
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &token)
	require.NoError(t, err)
	require.Equal(t, token.GetContractID(), contractID)
	require.Equal(t, token.GetTokenID(), tokenNFTID1)
}

func TestNewQuerier_encodeQueryParent_nil(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"parent","data":{"tokens_param":{"contract_id":"%s", "token_id":"%s"}}}`, contractID, tokenNFTID1)

	var token types.Token
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &token)
	require.NoError(t, err)
	require.Equal(t, token, nil)
}

func TestNewQuerier_encodeQueryRoot(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"root","data":{"tokens_param":{"contract_id":"%s", "token_id":"%s"}}}`, contractID, tokenNFTID3)

	var token types.Token
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &token)
	require.NoError(t, err)
	require.Equal(t, token.GetContractID(), contractID)
	require.Equal(t, token.GetTokenID(), tokenNFTID1)
}

func TestNewQuerier_encodeQueryRoot_self(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"root","data":{"tokens_param":{"contract_id":"%s", "token_id":"%s"}}}`, contractID, tokenNFTID1)

	var token types.Token
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &token)
	require.NoError(t, err)
	require.Equal(t, token.GetContractID(), contractID)
	require.Equal(t, token.GetTokenID(), tokenNFTID1)
}

func TestNewQuerier_encodeQueryChildren(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"children","data":{"tokens_param":{"contract_id":"%s", "token_id":"%s"}}}`, contractID, tokenNFTID1)

	var tokens types.Tokens
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &tokens)
	require.NoError(t, err)
	require.Equal(t, len(tokens), 2)
	require.Equal(t, tokens[0].GetContractID(), contractID)
	require.Equal(t, tokens[0].GetTokenID(), tokenNFTID2)
	require.Equal(t, tokens[1].GetContractID(), contractID)
	require.Equal(t, tokens[1].GetTokenID(), tokenNFTID3)
}

func TestNewQuerier_encodeQueryChildren_empty(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"children","data":{"tokens_param":{"contract_id":"%s", "token_id":"%s"}}}`, contractID, tokenNFTID2)

	var tokens types.Tokens
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &tokens)
	require.NoError(t, err)
	require.Equal(t, len(tokens), 0)
}

func TestNewQuerier_encodeQueryApprovers(t *testing.T) {
	prepare(t)

	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"approver","data":{"approvers_param":{"contract_id":"%s", "proxy":"%s"}}}`, contractID, addr1)

	var acAd1 []sdk.AccAddress
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &acAd1)
	require.NoError(t, err)
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

func TestNewQuerier_encodeQueryIsApproved_true(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"approved","data":{"is_approved_param":{"contract_id":"%s", "proxy":"%s", "approver":"%s"}}}`, contractID, addr1, addr2)

	var approved bool
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &approved)
	require.NoError(t, err)
	require.True(t, approved)
}

func TestNewQuerier_encodeQueryIsApproved_false(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"approved","data":{"is_approved_param":{"contract_id":"%s", "proxy":"%s", "approver":"%s"}}}`, contractID, addr2, addr1)

	var approved bool
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &approved)
	require.NoError(t, err)
	require.False(t, approved)
}
