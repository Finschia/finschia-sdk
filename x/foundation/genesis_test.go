package foundation_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

func TestValidateGenesis(t *testing.T) {
	createAddress := func() sdk.AccAddress {
		return sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct{
		data foundation.GenesisState
		valid bool
	}{
		"minimal": {
			data: foundation.GenesisState{
			},
			valid: true,
		},
		"members": {
			data: foundation.GenesisState{
				Members: []foundation.Member{
					{
						Address: createAddress().String(),
						Participating: true,
					},
				},
			},
			valid: true,
		},
		"foundation info": {
			data: foundation.GenesisState{
				Foundation: &foundation.FoundationInfo{
					Operator: createAddress().String(),
					Version: 1,
				},
			},
			valid: true,
		},
		"invalid foundation tax": {
			data: foundation.GenesisState{
				Params: &foundation.Params{
					FoundationTax: sdk.NewDec(2),
				},
			},
		},
		"member of invalid address": {
			data: foundation.GenesisState{
				Members: []foundation.Member{
					{
						Address: "invalid-address",
						Participating: true,
					},
				},
			},
		},
		"invalid operator address": {
			data: foundation.GenesisState{
				Foundation: &foundation.FoundationInfo{
					Operator: "invalid-address",
					Version: 1,
				},
			},
		},
		"invalid foundation version": {
			data: foundation.GenesisState{
				Foundation: &foundation.FoundationInfo{
				},
			},
		},
		"invalid decision policy": {
			data: foundation.GenesisState{
				Foundation: foundation.FoundationInfo{
					Operator: createAddress().String(),
					Version: 1,
				}.WithDecisionPolicy(&foundation.ThresholdDecisionPolicy{
					Windows: &foundation.DecisionPolicyWindows{},
				}),
			},
		},
		"proposal of no proposers": {
			data: foundation.GenesisState{
				Proposals: []foundation.Proposal{
					{
						Id: 1,
						FoundationVersion: 1,
					},
				},
			},
		},
		"proposal of invalid foundation version": {
			data: foundation.GenesisState{
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id: 1,
						Proposers: []string{createAddress().String()},
					}.WithMsgs([]sdk.Msg{&foundation.MsgWithdrawFromTreasury{
						Operator: createAddress().String(),
						To: createAddress().String(),
						Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
					}}),
				},
			},
		},
		"proposal of empty msgs": {
			data: foundation.GenesisState{
				Proposals: []foundation.Proposal{
					{
						Id: 1,
						Proposers: []string{createAddress().String()},
						FoundationVersion: 1,
					},
				},
			},
		},
		"duplicate proposals": {
			data: foundation.GenesisState{
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id: 1,
						Proposers: []string{createAddress().String()},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{&foundation.MsgWithdrawFromTreasury{
						Operator: createAddress().String(),
						To: createAddress().String(),
						Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
					}}),
					*foundation.Proposal{
						Id: 1,
						Proposers: []string{createAddress().String()},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{&foundation.MsgWithdrawFromTreasury{
						Operator: createAddress().String(),
						To: createAddress().String(),
						Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
					}}),
				},
			},
		},
		"no proposal for the vote": {
			data: foundation.GenesisState{
				Votes: []foundation.Vote{
					{
						ProposalId: 1,
						Voter: createAddress().String(),
						Option: foundation.VOTE_OPTION_YES,
					},
				},
			},
		},
		"invalid voter": {
			data: foundation.GenesisState{
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id: 1,
						Proposers: []string{createAddress().String()},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{&foundation.MsgWithdrawFromTreasury{
						Operator: createAddress().String(),
						To: createAddress().String(),
						Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
					}}),
				},
				Votes: []foundation.Vote{
					{
						ProposalId: 1,
						Voter: "invalid-address",
						Option: foundation.VOTE_OPTION_YES,
					},
				},
			},
		},
		"invalid vote option": {
			data: foundation.GenesisState{
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id: 1,
						Proposers: []string{createAddress().String()},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{&foundation.MsgWithdrawFromTreasury{
						Operator: createAddress().String(),
						To: createAddress().String(),
						Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
					}}),
				},
				Votes: []foundation.Vote{
					{
						ProposalId: 1,
						Voter: createAddress().String(),
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		err := foundation.ValidateGenesis(tc.data)
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}
