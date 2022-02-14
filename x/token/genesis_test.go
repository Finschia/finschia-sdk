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
	addr := sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
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
		"invalid address in a balance": {
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
		"invalid amount in a balance": {
			&token.GenesisState{
				Balances: []token.Balance{
					{
						Address: addr.String(),
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
		"empty tokens in a balance": {
			&token.GenesisState{
				Balances: []token.Balance{
					{
						Address: addr.String(),
					},
				},
			},
			false,
		},
		"invalid name of class": {
			&token.GenesisState{
				Classes: []token.Token{{
					Name: string(make([]rune, 21)),
				}},
			},
			false,
		},
		"invalid symbol of class": {
			&token.GenesisState{
				Classes: []token.Token{{
					Symbol: "t",
				}},
			},
			false,
		},
		"invalid image uri of class": {
			&token.GenesisState{
				Classes: []token.Token{{
					ImageUri: string(make([]rune, 1001)),
				}},
			},
			false,
		},
		"invalid meta of class": {
			&token.GenesisState{
				Classes: []token.Token{{
					Meta: string(make([]rune, 1001)),
				}},
			},
			false,
		},
		"invalid decimals of class": {
			&token.GenesisState{
				Classes: []token.Token{{
					Decimals: -1,
				}},
			},
			false,
		},
		"invalid grantee of grant": {
			&token.GenesisState{
				Grants: []token.Grant{{
					Grantee: "invalid",
					Action: token.ActionMint,
				}},
			},
			false,
		},
		"invalid action of grant": {
			&token.GenesisState{
				Grants: []token.Grant{{
					Grantee: addr.String(),
					Action: "invalid",
				}},
			},
			false,
		},
		"invalid approver of approval": {
			&token.GenesisState{
				Approves: []token.Approve{{
					Approver: "invalid",
					Proxy: addr.String(),
				}},
			},
			false,
		},
		"invalid proxy of approval": {
			&token.GenesisState{
				Approves: []token.Approve{{
					Approver: addr.String(),
					Proxy: "invalid",
				}},
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
