package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

// todo: add tests
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
