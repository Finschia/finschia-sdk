package testutil

import (
	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/codec"
	"github.com/Finschia/finschia-sdk/codec/types"
	"github.com/Finschia/finschia-sdk/std"
	"github.com/Finschia/finschia-sdk/types/module"
	"github.com/Finschia/finschia-sdk/x/auth/tx"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/Finschia/finschia-sdk/store"
	storetypes "github.com/Finschia/finschia-sdk/store/types"
	sdk "github.com/Finschia/finschia-sdk/types"
)

func DefaultContextWithDB(tb testing.TB, key, tkey storetypes.StoreKey) sdk.Context {
	tb.Helper()
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(tkey, storetypes.StoreTypeTransient, db)
	err := cms.LoadLatestVersion()
	assert.NoError(tb, err)

	ctx := sdk.NewContext(cms, tmproto.Header{}, false, log.NewNopLogger())

	return ctx
}

type TestEncodingConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Codec             codec.Codec
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}

func MakeTestEncodingConfig(modules ...module.AppModuleBasic) TestEncodingConfig {
	cdc := codec.NewLegacyAmino()
	interfaceRegistry := types.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	encCfg := TestEncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             marshaler,
		TxConfig:          tx.NewTxConfig(marshaler, tx.DefaultSignModes),
		Amino:             cdc,
	}

	std.RegisterLegacyAminoCodec(encCfg.Amino)
	std.RegisterInterfaces(encCfg.InterfaceRegistry)
	mb := module.NewBasicManager(modules...)
	mb.RegisterLegacyAminoCodec(encCfg.Amino)
	mb.RegisterInterfaces(encCfg.InterfaceRegistry)
	return encCfg
}
