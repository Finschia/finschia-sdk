package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/line/lfb-sdk/types"
	"github.com/line/lfb-sdk/x/bank/types"
)

func TestGenesisStateValidate(t *testing.T) {

	testCases := []struct {
		name         string
		genesisState types.GenesisState
		expErr       bool
	}{
		{
			"valid genesisState",
			types.GenesisState{
				Params: types.DefaultParams(),
				Balances: []types.Balance{
					{
						Address: "link1mejkku76a2ec35262rdqddggzwrgtrh52t3t0c",
						Coins:   sdk.Coins{sdk.NewInt64Coin("uatom", 1)},
					},
				},
				Supply: sdk.Coins{sdk.NewInt64Coin("uatom", 1)},
				DenomMetadata: []types.Metadata{
					{
						Description: "The native staking token of the Cosmos Hub.",
						DenomUnits: []*types.DenomUnit{
							{"uatom", uint32(0), []string{"microatom"}},
							{"matom", uint32(3), []string{"milliatom"}},
							{"atom", uint32(6), nil},
						},
						Base:    "uatom",
						Display: "atom",
					},
				},
			},
			false,
		},
		{"empty genesisState", types.GenesisState{}, false},
		{
			"invalid params ",
			types.GenesisState{
				Params: types.Params{
					SendEnabled: []*types.SendEnabled{
						{"", true},
					},
				},
			},
			true,
		},
		{
			"dup balances",
			types.GenesisState{
				Balances: []types.Balance{
					{
						Address: "link1mejkku76a2ec35262rdqddggzwrgtrh52t3t0c",
						Coins:   sdk.Coins{sdk.NewInt64Coin("uatom", 1)},
					},
					{
						Address: "link1mejkku76a2ec35262rdqddggzwrgtrh52t3t0c",
						Coins:   sdk.Coins{sdk.NewInt64Coin("uatom", 1)},
					},
				},
			},
			true,
		},
		{
			"0  balance",
			types.GenesisState{
				Balances: []types.Balance{
					{
						Address: "link1mejkku76a2ec35262rdqddggzwrgtrh52t3t0c",
					},
				},
			},
			false,
		},
		{
			"dup Metadata",
			types.GenesisState{
				DenomMetadata: []types.Metadata{
					{
						Description: "The native staking token of the Cosmos Hub.",
						DenomUnits: []*types.DenomUnit{
							{"uatom", uint32(0), []string{"microatom"}},
							{"matom", uint32(3), []string{"milliatom"}},
							{"atom", uint32(6), nil},
						},
						Base:    "uatom",
						Display: "atom",
					},
					{
						Description: "The native staking token of the Cosmos Hub.",
						DenomUnits: []*types.DenomUnit{
							{"uatom", uint32(0), []string{"microatom"}},
							{"matom", uint32(3), []string{"milliatom"}},
							{"atom", uint32(6), nil},
						},
						Base:    "uatom",
						Display: "atom",
					},
				},
			},
			true,
		},
		{
			"invalid Metadata",
			types.GenesisState{
				DenomMetadata: []types.Metadata{
					{
						Description: "The native staking token of the Cosmos Hub.",
						DenomUnits: []*types.DenomUnit{
							{"uatom", uint32(0), []string{"microatom"}},
							{"matom", uint32(3), []string{"milliatom"}},
							{"atom", uint32(6), nil},
						},
						Base:    "",
						Display: "",
					},
				},
			},
			true,
		},
		{
			"invalid supply",
			types.GenesisState{
				Supply: sdk.Coins{sdk.Coin{Denom: "", Amount: sdk.OneInt()}},
			},
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			err := tc.genesisState.Validate()

			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
