package zkauth_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-sdk/x/zkauth"
	datest "github.com/Finschia/finschia-sdk/x/zkauth/testutil"
	"github.com/Finschia/finschia-sdk/x/zkauth/types"
)

func TestGenesis(t *testing.T) {
	testCases := map[string]struct {
		genesisState types.GenesisState
		valid        bool
	}{
		"default genesis": {
			types.GenesisState{
				Params: types.DefaultParams(),
			},
			true,
		},
	}

	k, ctx := datest.ZkAuthKeeper(t)
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.valid {
				zkauth.InitGenesis(ctx, *k, tc.genesisState)
				got := zkauth.ExportGenesis(ctx, *k)
				require.NotNil(t, got)
			} else {
				require.Panics(t, func() { zkauth.InitGenesis(ctx, *k, tc.genesisState) })
			}
		})
	}

	// TODO: compare got & genesisState for each field
}
