package token_test

import (
	"testing"
	"math"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
)

func TestValidateGenesis(t *testing.T) {
	addr := secp256k1.GenPrivKey().PubKey().Address().String()
	testCases := map[string]struct{
		gs *token.GenesisState
		valid bool
	}{
		"default genesis": {
			token.DefaultGenesisState(),
			true,
		},
		"invalid class nonce": {
			&token.GenesisState{
				ClassState: &token.ClassGenesisState{
					Nonce: sdk.NewUint(math.MaxUint64).Incr(),
				},
			},
			false,
		},
		"invalid balances (invalid address)": {
			&token.GenesisState{
				Balances: []token.Balance{
					{
						Address: "INVALID",
						Tokens: []token.FT{
							{
								ClassId: "deadbeef",
								Amount: sdk.OneInt(),
							},
						},
					},
				},
			},
			false,
		},
		"invalid balances (invalid amount)": {
			&token.GenesisState{
				Balances: []token.Balance{
					{
						Address: addr,
						Tokens: []token.FT{
							{
								ClassId: "deadbeef",
								Amount: sdk.ZeroInt(),
							},
						},
					},
				},
			},
			false,
		},
		"invalid balances (empty tokens)": {
			&token.GenesisState{
				Balances: []token.Balance{
					{
						Address: addr,
					},
				},
			},
			false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := token.ValidateGenesis(*tc.gs)
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
