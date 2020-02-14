package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/line/link/x/iam"
	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"

	cbank "github.com/cosmos/cosmos-sdk/x/bank"
)

type TestInput struct {
	Cdc    *codec.Codec
	Ctx    sdk.Context
	Keeper Keeper
	Ak     auth.AccountKeeper
	Bk     cbank.BaseKeeper
	Iam    iam.Keeper
}

func newTestCodec() *codec.Codec {
	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	cbank.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	iam.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

func SetupTestInput(t *testing.T) *TestInput {
	keyAuth := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyIam := sdk.NewKVStoreKey(iam.StoreKey)

	keyLrc := sdk.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAuth, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyLrc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyIam, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.NoError(t, err)

	cdc := newTestCodec()

	// init params keeper and subspaces
	paramsKeeper := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	authSubspace := paramsKeeper.Subspace(auth.DefaultParamspace)
	cbankSubspace := paramsKeeper.Subspace(cbank.DefaultParamspace)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[sdk.AccAddress([]byte("moduleAcc")).String()] = true

	// add keepers
	accountKeeper := auth.NewAccountKeeper(cdc, keyAuth, authSubspace, auth.ProtoBaseAccount)
	cbankKeeper := cbank.NewBaseKeeper(accountKeeper, cbankSubspace, cbank.DefaultCodespace, blacklistedAddrs)
	iamKeeper := iam.NewKeeper(cdc, keyIam)

	// module account permissions
	maccPerms := map[string][]string{
		types.ModuleName: {supply.Burner, supply.Minter},
	}

	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, cbankKeeper, maccPerms)
	keeper := NewKeeper(cdc, supplyKeeper, iamKeeper.WithPrefix(types.ModuleName), accountKeeper, cbankKeeper, keyLrc)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	supplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.NewCoins()))

	return &TestInput{Cdc: cdc, Ctx: ctx, Keeper: keeper, Ak: accountKeeper, Bk: cbankKeeper}
}
