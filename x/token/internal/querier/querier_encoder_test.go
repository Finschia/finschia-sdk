package querier

// import (
// 	"encoding/json"
// 	"fmt"
// 	"testing"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
// 	"github.com/line/lbm-sdk/v2/x/token/internal/types"
// 	"github.com/line/lbm-sdk/v2/x/wasm"
// 	"github.com/stretchr/testify/require"
// )

// var (
// 	tokenQueryEncoder wasm.EncodeQuerier
// )

// func setupQueryEncoder() {
// 	tokenQuerier := NewQuerier(tkeeper)

// 	tokenQueryEncoder = NewQueryEncoder(tokenQuerier)
// }

// func encodeQuery(t *testing.T, jsonQuerier json.RawMessage, result interface{}) error {
// 	res, err := tokenQueryEncoder(ctx, jsonQuerier)
// 	if len(res) > 0 {
// 		require.NoError(t, tkeeper.UnmarshalJSON(res, result))
// 	}
// 	return err
// }

// func TestNewQuerier_encodeQueryTokens(t *testing.T) {
// 	prepare(t)
// 	setupQueryEncoder()
// 	jsonQuerier := fmt.Sprintf(`{"route":"tokens","data":{"token_param":{"contract_id":"%s"}}}`, contractID)

// 	var token types.Token
// 	err := encodeQuery(t, json.RawMessage(jsonQuerier), &token)
// 	require.NoError(t, err)
// 	require.Equal(t, token.GetContractID(), contractID)
// 	require.Equal(t, token.GetName(), tokenName)
// 	require.Equal(t, token.GetSymbol(), tokenSymbol)
// 	require.Equal(t, token.GetImageURI(), tokenImageURL)
// }

// func TestNewQuerier_encodeQueryAccountPermission(t *testing.T) {
// 	prepare(t)
// 	setupQueryEncoder()
// 	jsonQuerier := fmt.Sprintf(`{"route":"perms","data":{"perm_param":{"contract_id":"%s","address":"%s"}}}`, contractID, addr1)

// 	var perms types.Permissions
// 	err := encodeQuery(t, json.RawMessage(jsonQuerier), &perms)
// 	require.NoError(t, err)
// 	require.Equal(t, len(perms), 3)
// 	require.Equal(t, perms[0].String(), "modify")
// 	require.Equal(t, perms[1].String(), "mint")
// 	require.Equal(t, perms[2].String(), "burn")
// }

// func TestNewQuerier_encodeQueryBalance(t *testing.T) {
// 	prepare(t)
// 	setupQueryEncoder()
// 	jsonQuerier := fmt.Sprintf(`{"route":"balance","data":{"balance_param":{"contract_id":"%s","address":"%s"}}}`, contractID, addr1)

// 	var balance sdk.Int
// 	err := encodeQuery(t, json.RawMessage(jsonQuerier), &balance)
// 	require.NoError(t, err)
// 	require.Equal(t, balance.Int64(), int64(tokenAmount-tokenBurned))
// }

// func TestNewQuerier_encodeQueryTotalSupply(t *testing.T) {
// 	prepare(t)
// 	setupQueryEncoder()
// 	jsonQuerier := fmt.Sprintf(`{"route":"%s","data":{"total_param":{"contract_id":"%s"}}}`, types.QuerySupply, contractID)

// 	var supply sdk.Int
// 	err := encodeQuery(t, json.RawMessage(jsonQuerier), &supply)
// 	require.NoError(t, err)
// 	require.Equal(t, supply.Int64(), int64(tokenAmount-tokenBurned))
// }

// func TestNewQuerier_encodeQueryTotalMint(t *testing.T) {
// 	prepare(t)
// 	setupQueryEncoder()
// 	jsonQuerier := fmt.Sprintf(`{"route":"%s","data":{"total_param":{"contract_id":"%s"}}}`, types.QueryMint, contractID)

// 	var supply sdk.Int
// 	err := encodeQuery(t, json.RawMessage(jsonQuerier), &supply)
// 	require.NoError(t, err)
// 	require.Equal(t, supply.Int64(), int64(tokenAmount))
// }

// func TestNewQuerier_encodeQueryTotalBurn(t *testing.T) {
// 	prepare(t)
// 	setupQueryEncoder()
// 	jsonQuerier := fmt.Sprintf(`{"route":"%s","data":{"total_param":{"contract_id":"%s"}}}`, types.QueryBurn, contractID)

// 	var supply sdk.Int
// 	err := encodeQuery(t, json.RawMessage(jsonQuerier), &supply)
// 	require.NoError(t, err)
// 	require.Equal(t, supply.Int64(), int64(tokenBurned))
// }

// func TestNewQuerier_encodeQueryIsApproved_true(t *testing.T) {
// 	prepare(t)
// 	setupQueryEncoder()
// 	jsonQuerier := fmt.Sprintf(`{"route":"approved","data":{"is_approved_param":{"proxy":"%s", "contract_id":"%s","approver":"%s"}}}`, addr1, contractID, addr2)

// 	var approved bool
// 	err := encodeQuery(t, json.RawMessage(jsonQuerier), &approved)
// 	require.NoError(t, err)
// 	require.True(t, approved)
// }

// func TestNewQuerier_encodeQueryApprovers(t *testing.T) {
// 	prepare(t)
// 	setupQueryEncoder()
// 	jsonQuerier := fmt.Sprintf(`{"route":"approvers","data":{"approvers_param":{"proxy":"%s", "contract_id":"%s"}}}`, addr1, contractID)

// 	var approvers []sdk.AccAddress
// 	err := encodeQuery(t, json.RawMessage(jsonQuerier), &approvers)
// 	require.NoError(t, err)
// 	require.Equal(t, 2, len(approvers))
// 	require.True(t, types.IsAddressContains(approvers, addr3))
// 	require.True(t, types.IsAddressContains(approvers, addr2))

// 	var acAdEmpty []sdk.AccAddress
// 	jsonQuerier = fmt.Sprintf(`{"route":"approvers","data":{"approvers_param":{"proxy":"%s", "contract_id":"%s"}}}`, addr2, contractID)

// 	err = encodeQuery(t, json.RawMessage(jsonQuerier), &acAdEmpty)
// 	require.NoError(t, err)
// 	require.Empty(t, acAdEmpty)
// }

// func TestNewQuerier_invalidEncode(t *testing.T) {
// 	prepare(t)
// 	setupQueryEncoder()
// 	jsonQuerier := fmt.Sprintln(`{"route":"noquery","data":{"query_invalid_param":""}}`)

// 	err := encodeQuery(t, json.RawMessage(jsonQuerier), nil)
// 	require.EqualError(t, err, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized Msg route: %T", "noquery").Error())
// }
