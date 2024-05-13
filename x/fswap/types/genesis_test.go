package types_test

import (
	"testing"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisStateValidate(t *testing.T) {
	testSwapRate, _ := sdk.NewDecFromStr("1234567890")
	exampleGenesis := &types.GenesisState{
		Swaps: []types.Swap{
			{
				FromDenom:           "aaa",
				ToDenom:             "bbb",
				AmountCapForToDenom: sdk.NewInt(1234567890000),
				SwapRate:            testSwapRate,
			},
		},
		SwapStats: types.SwapStats{
			SwapCount: 1,
		},
		Swappeds: []types.Swapped{
			{
				FromCoinAmount: sdk.Coin{
					Denom:  "aaa",
					Amount: sdk.ZeroInt(),
				},
				ToCoinAmount: sdk.Coin{
					Denom:  "bbb",
					Amount: sdk.ZeroInt(),
				},
			},
		},
	}
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		modify   func(*types.GenesisState)
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc:     "example is valid",
			genState: exampleGenesis,
			valid:    true,
		},
		{
			desc:     "SwapCount is nagative in SwapStats is invalid",
			genState: exampleGenesis,
			modify: func(gs *types.GenesisState) {
				gs.SwapStats.SwapCount = -1
			},
			valid: false,
		},
		{
			desc:     "number of swaps does not match number of Swappeds",
			genState: exampleGenesis,
			modify: func(gs *types.GenesisState) {
				gs.Swaps = append(gs.Swaps, types.Swap{})
			},
			valid: false,
		},
		{
			desc:     "number of swaps does not match number of Swappeds",
			genState: exampleGenesis,
			modify: func(gs *types.GenesisState) {
				gs.Swaps = append(gs.Swaps, types.Swap{})
				gs.Swappeds = append(gs.Swappeds, types.Swapped{})
			},
			valid: false,
		},
		{
			desc:     "fromDenom=toDenom in Swap is invalid",
			genState: exampleGenesis,
			modify: func(gs *types.GenesisState) {
				gs.Swaps[0].ToDenom = "aaa"
			},
			valid: false,
		},
		{
			desc:     "AmountCapForToDenom=0 in Swap is invalid",
			genState: exampleGenesis,
			modify: func(gs *types.GenesisState) {
				gs.Swaps[0].AmountCapForToDenom = sdk.ZeroInt()
			},
			valid: false,
		},
		{
			desc:     "SwapRate=0 in Swap is invalid",
			genState: exampleGenesis,
			modify: func(gs *types.GenesisState) {
				gs.Swaps[0].SwapRate = sdk.ZeroDec()
			},
			valid: false,
		},
		{
			desc:     "FromCoinAmount is nagative in Swappeds is invalid",
			genState: exampleGenesis,
			modify: func(gs *types.GenesisState) {
				gs.Swappeds[0].FromCoinAmount.Amount = sdk.NewInt(-1)
			},
			valid: false,
		},
		{
			desc:     "ToCoinAmount is nagative in Swappeds is invalid",
			genState: exampleGenesis,
			modify: func(gs *types.GenesisState) {
				gs.Swappeds[0].ToCoinAmount.Amount = sdk.NewInt(-1)
			},
			valid: false,
		},
		{
			desc:     "FromCoin in swap and swapped do not correspond",
			genState: exampleGenesis,
			modify: func(gs *types.GenesisState) {
				gs.Swappeds[0].FromCoinAmount.Denom = "ccc"
			},
			valid: false,
		},
		{
			desc:     "ToCoin in swap and swapped do not correspond",
			genState: exampleGenesis,
			modify: func(gs *types.GenesisState) {
				gs.Swappeds[0].ToCoinAmount.Denom = "ccc"
			},
			valid: false,
		},
		{
			desc:     "AmountCapForToDenom has been exceeded is invalid",
			genState: exampleGenesis,
			modify: func(gs *types.GenesisState) {
				gs.Swappeds[0].ToCoinAmount.Amount = sdk.NewInt(12345678900000000)
			},
			valid: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			if tc.modify != nil {
				tc.modify(tc.genState)
			}
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
