package keeper_test

import (
	"github.com/line/ostracon/libs/log"
	ostproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/line/tm-db/v2/memdb"

	"github.com/line/lfb-sdk/simapp"

	"github.com/line/lfb-sdk/codec"
	"github.com/line/lfb-sdk/store"
	sdk "github.com/line/lfb-sdk/types"
	paramskeeper "github.com/line/lfb-sdk/x/params/keeper"
)

func testComponents() (*codec.LegacyAmino, sdk.Context, sdk.StoreKey, paramskeeper.Keeper) {
	marshaler := simapp.MakeTestEncodingConfig().Marshaler
	legacyAmino := createTestCodec()
	mkey := sdk.NewKVStoreKey("test")
	ctx := defaultContext(mkey)
	keeper := paramskeeper.NewKeeper(marshaler, legacyAmino, mkey)

	return legacyAmino, ctx, mkey, keeper
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

func defaultContext(key sdk.StoreKey) sdk.Context {
	db := memdb.NewDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	err := cms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}
	ctx := sdk.NewContext(cms, ostproto.Header{}, false, log.NewNopLogger())
	return ctx
}
