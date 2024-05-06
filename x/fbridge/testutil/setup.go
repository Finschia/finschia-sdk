package testutil

import (
	"github.com/Finschia/finschia-sdk/crypto/keys/secp256k1"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/codec"
	codectypes "github.com/Finschia/finschia-sdk/codec/types"
	"github.com/Finschia/finschia-sdk/std"
	"github.com/Finschia/finschia-sdk/store"
	storetypes "github.com/Finschia/finschia-sdk/store/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/auth/tx"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func DefaultContextWithDB(tb testing.TB, key, mkey, tkey storetypes.StoreKey) sdk.Context {
	tb.Helper()
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(mkey, storetypes.StoreTypeMemory, db)
	cms.MountStoreWithDB(tkey, storetypes.StoreTypeTransient, db)
	err := cms.LoadLatestVersion()
	assert.NoError(tb, err)

	ctx := sdk.NewContext(cms, tmproto.Header{}, false, log.NewNopLogger())

	return ctx
}

type TestEncodingConfig struct {
	InterfaceRegistry codectypes.InterfaceRegistry
	Codec             codec.Codec
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}

func MakeTestEncodingConfig() TestEncodingConfig {
	cdc := codec.NewLegacyAmino()
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	encCfg := TestEncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             marshaler,
		TxConfig:          tx.NewTxConfig(marshaler, tx.DefaultSignModes),
		Amino:             cdc,
	}

	std.RegisterLegacyAminoCodec(encCfg.Amino)
	std.RegisterInterfaces(encCfg.InterfaceRegistry)
	types.RegisterLegacyAminoCodec(encCfg.Amino)
	types.RegisterInterfaces(encCfg.InterfaceRegistry)
	return encCfg
}

func PrepareFbridgeTest(tb testing.TB, n int) (*sdk.KVStoreKey, *sdk.MemoryStoreKey, sdk.Context, TestEncodingConfig, *MockAccountKeeper, *MockBankKeeper, []sdk.AccAddress) {
	tb.Helper()

	ctrl := gomock.NewController(tb)
	key := storetypes.NewKVStoreKey(types.StoreKey)
	memKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	ctx := DefaultContextWithDB(tb, key, memKey, sdk.NewTransientStoreKey("transient_test"))
	encCfg := MakeTestEncodingConfig()

	authKeeper := NewMockAccountKeeper(ctrl)
	authKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(authtypes.NewEmptyModuleAccount("fbridge").GetAddress()).AnyTimes()
	bankKeeper := NewMockBankKeeper(ctrl)

	addrs := make([]sdk.AccAddress, n)
	for i := 0; i < n; i++ {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	return key, memKey, ctx, encCfg, authKeeper, bankKeeper, addrs
}
