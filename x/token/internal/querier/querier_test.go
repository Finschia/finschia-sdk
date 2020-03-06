package querier

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/iam/exported"
	"github.com/line/link/x/token/internal/keeper"
	"github.com/line/link/x/token/internal/types"
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
	require.NoError(t, tkeeper.IssueToken(ctx, types.NewToken(contractID, tokenName, tokenSymbol, tokenImageURL, sdk.NewInt(1), true), sdk.NewInt(tokenAmount), addr1, addr1))
	require.NoError(t, tkeeper.BurnToken(ctx, contractID, sdk.NewInt(tokenBurned), addr1))
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
	querier := NewQuerier(tkeeper)
	res, err := querier(ctx, path, req)
	require.NoError(t, err)
	if len(res) > 0 {
		require.NoError(t, tkeeper.UnmarshalJSON(res, result))
	}
}

func TestNewQuerier_queryAccountPermission(t *testing.T) {
	prepare(t)

	params := types.QueryAccAddressParams{
		Addr: addr1,
	}
	var perms []exported.PermissionI
	query(t, params, types.QueryPerms, &perms)
	require.Equal(t, len(perms), 3)
	require.Equal(t, perms[0].GetResource(), contractID)
	require.Equal(t, perms[0].GetAction(), "modify")
	require.Equal(t, perms[1].GetResource(), contractID)
	require.Equal(t, perms[1].GetAction(), "mint")
	require.Equal(t, perms[2].GetResource(), contractID)
	require.Equal(t, perms[2].GetAction(), "burn")
}

func TestNewQuerier_queryTokens_one(t *testing.T) {
	prepare(t)

	params := types.QueryContractIDParams{
		ContractID: contractID,
	}
	var token types.Token
	query(t, params, types.QueryTokens, &token)
	require.Equal(t, token.GetContractID(), contractID)
	require.Equal(t, token.GetName(), tokenName)
	require.Equal(t, token.GetSymbol(), tokenSymbol)
	require.Equal(t, token.GetImageURI(), tokenImageURL)
}

func TestNewQuerier_queryTokens_all(t *testing.T) {
	prepare(t)

	var tokens types.Tokens
	query(t, nil, types.QueryTokens, &tokens)
	require.Equal(t, len(tokens), 1)
	require.Equal(t, tokens[0].GetContractID(), contractID)
	require.Equal(t, tokens[0].GetName(), tokenName)
	require.Equal(t, tokens[0].GetSymbol(), tokenSymbol)
	require.Equal(t, tokens[0].GetImageURI(), tokenImageURL)
}

func TestNewQuerier_queryBalance(t *testing.T) {
	prepare(t)

	params := types.QueryAccAddressContractIDParams{
		Addr:       addr1,
		ContractID: contractID,
	}
	var balance sdk.Int
	query(t, params, types.QueryBalance, &balance)
	require.Equal(t, balance.Int64(), int64(tokenAmount-tokenBurned))
}

func TestNewQuerier_queryTotalSupply(t *testing.T) {
	prepare(t)

	params := types.QueryContractIDParams{
		ContractID: contractID,
	}
	var supply sdk.Int
	query(t, params, types.QuerySupply, &supply)
	require.Equal(t, supply.Int64(), int64(tokenAmount-tokenBurned))
}

func TestNewQuerier_queryTotalMint(t *testing.T) {
	prepare(t)

	params := types.QueryContractIDParams{
		ContractID: contractID,
	}
	var supply sdk.Int
	query(t, params, types.QueryMint, &supply)
	require.Equal(t, supply.Int64(), int64(tokenAmount))
}

func TestNewQuerier_queryTotalBurn(t *testing.T) {
	prepare(t)

	params := types.QueryContractIDParams{
		ContractID: contractID,
	}
	var supply sdk.Int
	query(t, params, types.QueryBurn, &supply)
	require.Equal(t, supply.Int64(), int64(tokenBurned))
}

func TestNewQuerier_invalid(t *testing.T) {
	prepare(t)
	params := types.QueryAccAddressContractIDParams{
		ContractID: contractID,
		Addr:       addr1,
	}
	querier := NewQuerier(tkeeper)
	path := []string{"noquery"}
	req := abci.RequestQuery{
		Path: "",
		Data: []byte(string(codec.MustMarshalJSONIndent(types.ModuleCdc, params))),
	}
	_, err := querier(ctx, path, req)
	require.EqualError(t, err, sdk.ErrUnknownRequest("unknown token query endpoint").Error())
}
