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
			desc: "empty oldCoin in Swapped",
			genState: &types.GenesisState{
				Swapped: types.Swapped{
					NewCoinAmount: sdk.NewInt64Coin("atom", 1),
				},
			},
			valid: false,
		},
		{
			desc: "empty newCoin in Swapped",
			genState: &types.GenesisState{
				Swapped: types.Swapped{
					OldCoinAmount: sdk.NewInt64Coin("cony", 1),
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
