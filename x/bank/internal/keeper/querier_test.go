package keeper

import (
	"fmt"
	"github.com/line/link/x/bank/internal/types"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func TestBalances(t *testing.T) {
	input := SetupTestInput()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/bank/%s", QueryBalance),
		Data: []byte{},
	}

	querier := NewQuerier(input.K)

	res, err := querier(input.Ctx, []string{"balances"}, req)
	require.NotNil(t, err)
	require.Nil(t, res)

	_, _, addr := authtypes.KeyTestPubAddr()
	req.Data = input.Cdc.MustMarshalJSON(types.NewQueryBalanceParams(addr, ""))
	res, err = querier(input.Ctx, []string{"balances"}, req)
	require.Nil(t, err) // the account does not exist, no error returned anyway
	require.NotNil(t, res)

	var coins sdk.Coins
	require.NoError(t, input.Cdc.UnmarshalJSON(res, &coins))
	require.True(t, coins.IsZero())

	acc := input.Ak.NewAccountWithAddress(input.Ctx, addr)
	require.NoError(t, acc.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("foo", 10))))
	input.Ak.SetAccount(input.Ctx, acc)
	res, err = querier(input.Ctx, []string{"balances"}, req)
	require.Nil(t, err)
	require.NotNil(t, res)
	require.NoError(t, input.Cdc.UnmarshalJSON(res, &coins))
	require.True(t, coins.AmountOf("foo").Equal(sdk.NewInt(10)))

	// Query with denomination
	var amount sdk.Int
	req.Data = input.Cdc.MustMarshalJSON(types.NewQueryBalanceParams(addr, "foo"))
	res, err = querier(input.Ctx, []string{"balances"}, req)
	require.Nil(t, err) // the account does not exist, no error returned anyway
	require.NotNil(t, res)

	require.NoError(t, input.Cdc.UnmarshalJSON(res, &amount))
	require.True(t, amount.Equal(sdk.NewInt(10)))

}

func TestQuerierRouteNotFound(t *testing.T) {
	input := SetupTestInput()
	req := abci.RequestQuery{
		Path: "custom/bank/notfound",
		Data: []byte{},
	}

	querier := NewQuerier(input.K)
	_, err := querier(input.Ctx, []string{"notfound"}, req)
	require.Error(t, err)
}
