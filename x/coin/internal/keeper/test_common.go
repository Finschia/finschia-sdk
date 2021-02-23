package keeper

// DONTCOVER

import (
	"github.com/line/lbm-sdk/x/coin/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/line/lbm-sdk/codec"
	"github.com/line/lbm-sdk/store"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/auth"
	"github.com/line/lbm-sdk/x/bank"
	"github.com/line/lbm-sdk/x/params"
)

type TestInput struct {
	Cdc *codec.Codec
	Ctx sdk.Context
	K   Keeper
	Ak  auth.AccountKeeper
	Pk  params.Keeper
}

func SetupTestInput() TestInput {
	db := dbm.NewMemDB()

	cdc := codec.New()
	auth.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	authCapKey := sdk.NewKVStoreKey("authCapKey")
	keyParams := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")
	keyBank := sdk.NewKVStoreKey(types.StoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyBank, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)

	if err := ms.LoadLatestVersion(); err != nil {
		panic(err)
	}

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[sdk.AccAddress([]byte("moduleAcc")).String()] = true

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)

	ak := auth.NewAccountKeeper(
		cdc, authCapKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount,
	)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	ak.SetParams(ctx, auth.DefaultParams())

	bankKeeper := bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), blacklistedAddrs)
	bankKeeper.SetSendEnabled(ctx, true)

	keeper := NewKeeper(bankKeeper, keyBank)

	return TestInput{Cdc: cdc, Ctx: ctx, K: keeper, Ak: ak, Pk: pk}
}
