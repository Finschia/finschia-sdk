package querier

import (
	"encoding/json"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/token/internal/types"
	"github.com/line/link-modules/x/wasm"
	"github.com/stretchr/testify/require"
)

var (
	tokenQueryEncoder wasm.EncodeQuerier
)

func setupQueryEncoder() {
	tokenQuerier := NewQuerier(tkeeper)

	tokenQueryEncoder = NewQueryEncoder(tokenQuerier)
}

func encodeQuery(t *testing.T, jsonQuerier json.RawMessage, result interface{}) error {
	res, err := tokenQueryEncoder(ctx, jsonQuerier)
	if len(res) > 0 {
		require.NoError(t, tkeeper.UnmarshalJSON(res, result))
	}
	return err
}

func TestNewQuerier_encodeQueryTokens(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"tokens","data":{"query_token_param":{"contract_id":"%s"}}}`, contractID)

	var token types.Token
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &token)
	require.NoError(t, err)
	require.Equal(t, token.GetContractID(), contractID)
	require.Equal(t, token.GetName(), tokenName)
	require.Equal(t, token.GetSymbol(), tokenSymbol)
	require.Equal(t, token.GetImageURI(), tokenImageURL)
}

func TestNewQuerier_encodeQueryAccountPermission(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"perms","data":{"query_perm_param":{"contract_id":"%s","address":"%s"}}}`, contractID, addr1)

	var perms types.Permissions
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &perms)
	require.NoError(t, err)
	require.Equal(t, len(perms), 3)
	require.Equal(t, perms[0].String(), "modify")
	require.Equal(t, perms[1].String(), "mint")
	require.Equal(t, perms[2].String(), "burn")
}

func TestNewQuerier_encodeQueryBalance(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"balance","data":{"query_balance_param":{"contract_id":"%s","address":"%s"}}}`, contractID, addr1)

	var balance sdk.Int
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &balance)
	require.NoError(t, err)
	require.Equal(t, balance.Int64(), int64(tokenAmount-tokenBurned))
}

func TestNewQuerier_encodeQueryTotalSupply(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"supply","data":{"query_total_param":{"contract_id":"%s","target":"%s"}}}`, contractID, types.QuerySupply)

	var supply sdk.Int
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &supply)
	require.NoError(t, err)
	require.Equal(t, supply.Int64(), int64(tokenAmount-tokenBurned))
}

func TestNewQuerier_encodeQueryTotalMint(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"supply","data":{"query_total_param":{"contract_id":"%s","target":"%s"}}}`, contractID, types.QueryMint)

	var supply sdk.Int
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &supply)
	require.NoError(t, err)
	require.Equal(t, supply.Int64(), int64(tokenAmount))
}

func TestNewQuerier_encodeQueryTotalBurn(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintf(`{"route":"supply","data":{"query_total_param":{"contract_id":"%s","target":"%s"}}}`, contractID, types.QueryBurn)

	var supply sdk.Int
	err := encodeQuery(t, json.RawMessage(jsonQuerier), &supply)
	require.NoError(t, err)
	require.Equal(t, supply.Int64(), int64(tokenBurned))
}

func TestNewQuerier_invalidEncode(t *testing.T) {
	prepare(t)
	setupQueryEncoder()
	jsonQuerier := fmt.Sprintln(`{"route":"noquery","data":{"query_invalid_param":""}}`)

	err := encodeQuery(t, json.RawMessage(jsonQuerier), nil)
	require.EqualError(t, err, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized Msg route: %T", "noquery").Error())
}
