package querier

import (
	"context"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/contract"
	"github.com/line/link-modules/x/token/internal/keeper"
	"github.com/line/link-modules/x/token/internal/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	contractID    = "9be17165"
	tokenName     = "linko token"
	tokenSymbol   = "LINKO"
	tokenImageURL = "url"
	tokenAmount   = 1000
	tokenBurned   = 10
	tokenMeta     = "{}"
)

var (
	ms      store.CommitMultiStore
	ctx     sdk.Context
	tkeeper keeper.Keeper
	addr1   sdk.AccAddress
	addr2   sdk.AccAddress
	addr3   sdk.AccAddress
)

func prepare(t *testing.T) {
	ctx, ms, tkeeper = keeper.TestKeeper()
	msCache := ms.CacheMultiStore()
	ctx = ctx.WithMultiStore(msCache)

	addr1 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	addr2 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	addr3 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	// prepare token
	ctx2 := ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, contractID))
	require.NoError(t, tkeeper.IssueToken(ctx2, types.NewToken(contractID, tokenName, tokenSymbol, tokenMeta, tokenImageURL, sdk.NewInt(1), true), sdk.NewInt(tokenAmount), addr1, addr1))
	require.NoError(t, tkeeper.BurnToken(ctx2, sdk.NewInt(tokenBurned), addr1))

	require.NoError(t, tkeeper.GrantPermission(ctx2, addr1, addr2, types.NewBurnPermission()))
	require.NoError(t, tkeeper.SetApproved(ctx2, addr1, addr2))

	// prepare one more approver for test proxy
	require.NoError(t, tkeeper.SetApproved(ctx2, addr1, addr3))
}

func query(t *testing.T, params interface{}, query string, result interface{}) {
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
	querier := NewQuerier(tkeeper)
	res, err := querier(ctx, path, req)
	require.NoError(t, err)
	if len(res) > 0 {
		require.NoError(t, tkeeper.UnmarshalJSON(res, result))
	}
}

func TestNewQuerier_queryAccountPermission(t *testing.T) {
	prepare(t)

	params := types.NewQueryContractIDAccAddressParams(addr1)
	var perms types.Permissions
	query(t, params, types.QueryPerms, &perms)
	require.Equal(t, len(perms), 3)
	require.Equal(t, perms[0].String(), "modify")
	require.Equal(t, perms[1].String(), "mint")
	require.Equal(t, perms[2].String(), "burn")
}

func TestNewQuerier_queryTokens_one(t *testing.T) {
	prepare(t)

	var token types.Token
	query(t, nil, types.QueryTokens, &token)
	require.Equal(t, token.GetContractID(), contractID)
	require.Equal(t, token.GetName(), tokenName)
	require.Equal(t, token.GetSymbol(), tokenSymbol)
	require.Equal(t, token.GetImageURI(), tokenImageURL)
}

func TestNewQuerier_queryBalance(t *testing.T) {
	prepare(t)

	params := types.QueryContractIDAccAddressParams{
		Addr: addr1,
	}
	var balance sdk.Int
	query(t, params, types.QueryBalance, &balance)
	require.Equal(t, balance.Int64(), int64(tokenAmount-tokenBurned))
}

func TestNewQuerier_queryTotalSupply(t *testing.T) {
	prepare(t)

	var supply sdk.Int
	query(t, nil, types.QuerySupply, &supply)
	require.Equal(t, supply.Int64(), int64(tokenAmount-tokenBurned))
}

func TestNewQuerier_queryTotalMint(t *testing.T) {
	prepare(t)

	var supply sdk.Int
	query(t, nil, types.QueryMint, &supply)
	require.Equal(t, supply.Int64(), int64(tokenAmount))
}

func TestNewQuerier_queryTotalBurn(t *testing.T) {
	prepare(t)

	var supply sdk.Int
	query(t, nil, types.QueryBurn, &supply)
	require.Equal(t, supply.Int64(), int64(tokenBurned))
}

func TestNewQuerier_invalid(t *testing.T) {
	prepare(t)
	params := types.QueryContractIDAccAddressParams{
		Addr: addr1,
	}
	querier := NewQuerier(tkeeper)
	path := []string{"noquery", contractID}
	req := abci.RequestQuery{
		Path: "",
		Data: []byte(string(codec.MustMarshalJSONIndent(types.ModuleCdc, params))),
	}
	_, err := querier(ctx, path, req)
	require.EqualError(t, err, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown token query endpoint").Error())
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

func TestNewQuerier_queryApprovers(t *testing.T) {
	prepare(t)
	params := types.QueryProxyParams{
		Proxy: addr1,
	}
	var approvers []sdk.AccAddress
	query(t, params, types.QueryApprovers, &approvers)
	require.Equal(t, 2, len(approvers))
	require.True(t, types.IsAddressContains(approvers, addr3))
	require.True(t, types.IsAddressContains(approvers, addr2))

	var acAdEmpty []sdk.AccAddress
	paramsEmpty := types.QueryProxyParams{
		Proxy: addr2,
	}
	query(t, paramsEmpty, types.QueryApprovers, &acAdEmpty)
	require.Empty(t, acAdEmpty)
}
