package simapp

import (
	"fmt"
	"testing"

	types3 "github.com/line/lbm-sdk/v2/store/types"
	"github.com/line/ostracon/abci/types"
	types2 "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/v2/codec"
	"github.com/line/lbm-sdk/v2/std"
	sdk "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/types/module"
	authtypes "github.com/line/lbm-sdk/v2/x/auth/types"
)

func makeCodec(bm module.BasicManager) *codec.LegacyAmino {
	cdc := codec.NewLegacyAmino()

	bm.RegisterLegacyAminoCodec(cdc)
	std.RegisterLegacyAminoCodec(cdc)

	return cdc
}

func TestSetup(t *testing.T) {
	app := Setup(false)
	ctx := app.BaseApp.NewContext(false, types2.Header{})

	app.InitChain(
		types.RequestInitChain{
			AppStateBytes: []byte("{}"),
			ChainId:       "test-chain-id",
		},
	)

	acc := app.AccountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(authtypes.FeeCollectorName))
	require.NotNil(t, acc)
}

func TestGetSimulationLog(t *testing.T) {
	cdc := makeCodec(ModuleBasics)

	decoders := make(sdk.StoreDecoderRegistry)
	decoders[authtypes.StoreKey] = sdk.StoreDecoder{
		Marshal:   nil,
		Unmarshal: nil,
		LogPair:   func(kvAs, kvBs types3.KOPair) string { return "10" },
	}

	tests := []struct {
		store       string
		kvPairs     []types3.KOPair
		expectedLog string
	}{
		{
			"Empty",
			[]types3.KOPair{{}},
			"",
		},
		{
			authtypes.StoreKey,
			[]types3.KOPair{{Key: authtypes.GlobalAccountNumberKey, Value: cdc.MustMarshalBinaryBare(uint64(10))}},
			"10",
		},
		{
			"OtherStore",
			[]types3.KOPair{{Key: []byte("key"), Value: []byte("value")}},
			fmt.Sprintf("store A %X => %X\nstore B %X => %X\n", []byte("key"), []byte("value"), []byte("key"), []byte("value")),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.store, func(t *testing.T) {
			require.Equal(t, tt.expectedLog, GetSimulationLog(tt.store, decoders, tt.kvPairs, tt.kvPairs), tt.store)
		})
	}
}
