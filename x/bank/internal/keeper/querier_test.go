package keeper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cbank "github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/line/link/x/bank/internal/keeper/mocks"
	"github.com/line/link/x/bank/internal/types"
)

type testInput struct {
	cdc *codec.Codec
	ctx sdk.Context
	k   cbank.Keeper
	ak  auth.AccountKeeper
	pk  params.Keeper
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()

	cdc := codec.New()
	auth.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	authCapKey := sdk.NewKVStoreKey("authCapKey")
	keyParams := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	_ = ms.LoadLatestVersion()

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[sdk.AccAddress([]byte("moduleAcc")).String()] = true

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)

	ak := auth.NewAccountKeeper(
		cdc, authCapKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount,
	)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	ak.SetParams(ctx, auth.DefaultParams())

	bankKeeper := cbank.NewBaseKeeper(
		ak, pk.Subspace(cbank.DefaultParamspace), cbank.DefaultCodespace, blacklistedAddrs)
	bankKeeper.SetSendEnabled(ctx, true)

	return testInput{cdc: cdc, ctx: ctx, k: bankKeeper, ak: ak, pk: pk}
}

func TestBalanceOf(t *testing.T) {
	input := setupTestInput()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/bank/%s", QueryBalanceOf),
		Data: []byte{},
	}

	mock := mocks.Queryable{}
	fallbackQuerier := mock.Querier

	querier := NewQuerier(input.k, fallbackQuerier)

	res, err := querier(input.ctx, []string{"balance_of"}, req)
	require.NotNil(t, err)
	require.Nil(t, res)

	_, _, addr := authtypes.KeyTestPubAddr()
	req.Data = input.cdc.MustMarshalJSON(types.NewQueryBalanceOfParams(addr, "foo"))
	res, err = querier(input.ctx, []string{"balance_of"}, req)
	require.Nil(t, err) // the account does not exist, no error returned anyway
	require.NotNil(t, res)

	var amount sdk.Int
	require.NoError(t, input.cdc.UnmarshalJSON(res, &amount))
	require.True(t, amount.IsZero())

	acc := input.ak.NewAccountWithAddress(input.ctx, addr)
	_ = acc.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("foo", 10)))
	input.ak.SetAccount(input.ctx, acc)
	res, err = querier(input.ctx, []string{"balance_of"}, req)
	require.Nil(t, err)
	require.NotNil(t, res)
	require.NoError(t, input.cdc.UnmarshalJSON(res, &amount))
	require.True(t, amount.Equal(sdk.NewInt(10)))
}

func TestQuerierRouteNotFoundThenCalledFallbackQuerier(t *testing.T) {
	input := setupTestInput()
	req := abci.RequestQuery{
		Path: "custom/bank/notfound",
		Data: []byte{},
	}
	path := []string{"notfound"}

	mock := mocks.Queryable{}
	mock.On("Querier", input.ctx, path, req).Return(nil, nil).Once()
	fallbackQuerier := mock.Querier

	querier := NewQuerier(input.k, fallbackQuerier)
	_, _ = querier(input.ctx, path, req)

	mock.AssertExpectations(t)
}
