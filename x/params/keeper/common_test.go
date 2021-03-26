package keeper_test

import (
	"github.com/line/ostracon/libs/log"
	ostproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/line/tm-db/v2/memdb"

	"github.com/line/lbm-sdk/v2/simapp"

	"github.com/line/lbm-sdk/v2/codec"
	"github.com/line/lbm-sdk/v2/store"
	sdk "github.com/line/lbm-sdk/v2/types"
	paramskeeper "github.com/line/lbm-sdk/v2/x/params/keeper"
)

func testComponents() (*codec.LegacyAmino, sdk.Context, sdk.StoreKey, sdk.StoreKey, paramskeeper.Keeper) {
	marshaler := simapp.MakeTestEncodingConfig().Marshaler
	legacyAmino := createTestCodec()
	mkey := sdk.NewKVStoreKey("test")
	tkey := sdk.NewTransientStoreKey("transient_test")
	ctx := defaultContext(mkey, tkey)
	keeper := paramskeeper.NewKeeper(marshaler, legacyAmino, mkey, tkey)

	return legacyAmino, ctx, mkey, tkey, keeper
}

type invalid struct{}

type s struct {
	I int
}

func createTestCodec() *codec.LegacyAmino {
	cdc := codec.NewLegacyAmino()
	sdk.RegisterLegacyAminoCodec(cdc)
	cdc.RegisterConcrete(s{}, "test/s", nil)
	cdc.RegisterConcrete(invalid{}, "test/invalid", nil)
	return cdc
}

func defaultContext(key sdk.StoreKey, tkey sdk.StoreKey) sdk.Context {
	db := memdb.NewDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, db)
	err := cms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}
	ctx := sdk.NewContext(cms, ostproto.Header{}, false, log.NewNopLogger())
	return ctx
}
