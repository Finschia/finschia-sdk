package keeper

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/line/lbm-sdk/v2/x/contract"
	"github.com/line/lbm-sdk/v2/x/token/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeeper() (sdk.Context, store.CommitMultiStore, Keeper) {
	keyAuth := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyToken := sdk.NewKVStoreKey(types.StoreKey)
	keyContract := sdk.NewKVStoreKey(contract.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAuth, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyToken, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyContract, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	if err := ms.LoadLatestVersion(); err != nil {
		panic(err)
	}

	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	cdc.Seal()

	paramsKeeper := params.NewKeeper(cdc, keyParams, tkeyParams)
	authSubspace := paramsKeeper.Subspace(auth.DefaultParamspace)

	// add keepers
	accountKeeper := auth.NewAccountKeeper(cdc, keyAuth, authSubspace, auth.ProtoBaseAccount)
	keeper := NewKeeper(cdc, accountKeeper, contract.NewContractKeeper(cdc, keyContract), keyToken)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	return ctx, ms, keeper
}
