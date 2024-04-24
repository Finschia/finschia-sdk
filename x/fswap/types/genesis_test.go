package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

func TestGenesisStateValidate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc:     "empty genesisState",
			genState: &types.GenesisState{},
			valid:    false,
		},
		{
			desc: "empty params",
			genState: &types.GenesisState{
				Swapped: types.DefaultSwapped(),
			},
			valid: false,
		},
		{
			desc: "empty swapped",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
			},
			valid: false,
		},
		{
			desc: "empty swappableNewCoinAmount in params",
			genState: &types.GenesisState{
				Params:  types.Params{},
				Swapped: types.DefaultSwapped(),
			},
			valid: false,
		},
		{
			desc: "empty oldCoin in Swapped",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				Swapped: types.Swapped{
					NewCoinAmount: sdk.ZeroInt(),
				},
			},
			valid: false,
		},
		{
			desc: "empty newCoin in Swapped",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				Swapped: types.Swapped{
					OldCoinAmount: sdk.ZeroInt(),
				},
			},
			valid: false,
		},
		{
			desc: "coinAmount is negative",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				Swapped: types.Swapped{
					NewCoinAmount: sdk.ZeroInt(),
					OldCoinAmount: sdk.NewInt(-1),
				},
			},
			valid: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
