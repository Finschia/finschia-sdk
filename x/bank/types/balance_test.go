package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/line/lbm-sdk/types"
	bank "github.com/line/lbm-sdk/x/bank/types"
)

func TestBalanceValidate(t *testing.T) {
	testCases := []struct {
		name    string
		balance bank.Balance
		expErr  bool
	}{
		{
			"valid balance",
			bank.Balance{
				Address: "link1yq8lgssgxlx9smjhes6ryjasmqmd3ts2p6925r",
				Coins:   sdk.Coins{sdk.NewInt64Coin("uatom", 1)},
			},
			false,
		},
		{"empty balance", bank.Balance{}, true},
		{
			"nil balance coins",
			bank.Balance{
				Address: "link1yq8lgssgxlx9smjhes6ryjasmqmd3ts2p6925r",
			},
			false,
		},
		{
			"dup coins",
			bank.Balance{
				Address: "link1yq8lgssgxlx9smjhes6ryjasmqmd3ts2p6925r",
				Coins: sdk.Coins{
					sdk.NewInt64Coin("uatom", 1),
					sdk.NewInt64Coin("uatom", 1),
				},
			},
			true,
		},
		{
			"invalid coin denom",
			bank.Balance{
				Address: "link1yq8lgssgxlx9smjhes6ryjasmqmd3ts2p6925r",
				Coins: sdk.Coins{
					sdk.Coin{Denom: "", Amount: sdk.OneInt()},
				},
			},
			true,
		},
		{
			"negative coin",
			bank.Balance{
				Address: "link1yq8lgssgxlx9smjhes6ryjasmqmd3ts2p6925r",
				Coins: sdk.Coins{
					sdk.Coin{Denom: "uatom", Amount: sdk.NewInt(-1)},
				},
			},
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			err := tc.balance.Validate()

			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBalance_GetAddress(t *testing.T) {
	tests := []struct {
		name    string
		Address string
		expErr  bool
	}{
		{"empty address", "", true},
		{"malformed address", "invalid", true},
		{"valid address", "link1vy0ga0klndqy92ceqehfkvgmn4t94ete4mhemy", false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			b := bank.Balance{Address: tt.Address}
			err := sdk.ValidateAccAddress(b.GetAddress().String())
			if tt.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
