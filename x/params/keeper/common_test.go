package keeper_test

import (
	"github.com/Finschia/finschia-sdk/codec"
	"github.com/Finschia/finschia-sdk/simapp"
	"github.com/Finschia/finschia-sdk/testutil"
	sdk "github.com/Finschia/finschia-sdk/types"
	paramskeeper "github.com/Finschia/finschia-sdk/x/params/keeper"
)

func testComponents() (*codec.LegacyAmino, sdk.Context, sdk.StoreKey, paramskeeper.Keeper) {
	marshaler := simapp.MakeTestEncodingConfig().Marshaler
	legacyAmino := createTestCodec()
	mkey := sdk.NewKVStoreKey("test")
	tkey := sdk.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(mkey, tkey)
	keeper := paramskeeper.NewKeeper(marshaler, legacyAmino, mkey, tkey)

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
