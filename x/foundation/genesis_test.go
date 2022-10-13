package foundation_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

func TestDefaultGenesisState(t *testing.T) {
	gs := foundation.DefaultGenesisState()
	require.Equal(t, false, gs.Params.Enabled)
	require.Equal(t, sdk.ZeroDec(), gs.Params.FoundationTax)
}

func TestValidateGenesis(t *testing.T) {
	createAddress := func() sdk.AccAddress {
		return sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		data  foundation.GenesisState
		valid bool
	}{
		"minimal": {
			data:  foundation.GenesisState{},
			valid: true,
		},
		"members": {
			data: foundation.GenesisState{
				Members: []foundation.Member{
					{
						Address: createAddress().String(),
					},
				},
			},
			valid: true,
		},
		"foundation info": {
			data: foundation.GenesisState{
				Foundation: &foundation.FoundationInfo{
					Operator: createAddress().String(),
					Version:  1,
				},
			},
			valid: true,
		},
		"authorizations": {
			data: foundation.GenesisState{
				Authorizations: []foundation.GrantAuthorization{
					*foundation.GrantAuthorization{
						Grantee: createAddress().String(),
					}.WithAuthorization(&foundation.ReceiveFromTreasuryAuthorization{}),
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
		"invalid members": {
			data: foundation.GenesisState{
				Members: []foundation.Member{{}},
			},
		},
		"invalid operator address": {
			data: foundation.GenesisState{
				Foundation: &foundation.FoundationInfo{
					Operator: "invalid-address",
					Version:  1,
				},
			},
		},
		"invalid foundation version": {
			data: foundation.GenesisState{
				Foundation: &foundation.FoundationInfo{},
			},
		},
		"invalid decision policy": {
			data: foundation.GenesisState{
				Foundation: foundation.FoundationInfo{
					Operator: createAddress().String(),
					Version:  1,
				}.WithDecisionPolicy(&foundation.ThresholdDecisionPolicy{
					Windows: &foundation.DecisionPolicyWindows{},
				}),
			},
		},
		"invalid proposals": {
			data: foundation.GenesisState{
				Proposals: []foundation.Proposal{{}},
			},
		},
		"duplicate proposals": {
			data: foundation.GenesisState{
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id:                1,
						Proposers:         []string{createAddress().String()},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{&foundation.MsgWithdrawFromTreasury{
						Operator: createAddress().String(),
						To:       createAddress().String(),
						Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
					}}),
					*foundation.Proposal{
						Id:                1,
						Proposers:         []string{createAddress().String()},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{&foundation.MsgWithdrawFromTreasury{
						Operator: createAddress().String(),
						To:       createAddress().String(),
						Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
					}}),
				},
			},
		},
		"no proposal for the vote": {
			data: foundation.GenesisState{
				Votes: []foundation.Vote{
					{
						ProposalId: 1,
						Voter:      createAddress().String(),
						Option:     foundation.VOTE_OPTION_YES,
					},
				},
			},
		},
		"invalid voter": {
			data: foundation.GenesisState{
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id:                1,
						Proposers:         []string{createAddress().String()},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{&foundation.MsgWithdrawFromTreasury{
						Operator: createAddress().String(),
						To:       createAddress().String(),
						Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
					}}),
				},
				Votes: []foundation.Vote{
					{
						ProposalId: 1,
						Voter:      "invalid-address",
						Option:     foundation.VOTE_OPTION_YES,
					},
				},
			},
		},
		"invalid vote option": {
			data: foundation.GenesisState{
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id:                1,
						Proposers:         []string{createAddress().String()},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{&foundation.MsgWithdrawFromTreasury{
						Operator: createAddress().String(),
						To:       createAddress().String(),
						Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
					}}),
				},
				Votes: []foundation.Vote{
					{
						ProposalId: 1,
						Voter:      createAddress().String(),
					},
				},
			},
		},
		"invalid authorization": {
			data: foundation.GenesisState{
				Authorizations: []foundation.GrantAuthorization{{
					Grantee: createAddress().String(),
				}},
			},
		},
		"invalid grantee": {
			data: foundation.GenesisState{
				Authorizations: []foundation.GrantAuthorization{
					*foundation.GrantAuthorization{}.WithAuthorization(&foundation.ReceiveFromTreasuryAuthorization{}),
				},
			},
		},
		"invalid pool": {
			data: foundation.GenesisState{
				Pool: foundation.Pool{
					Treasury: sdk.DecCoins{
						{
							Denom:  sdk.DefaultBondDenom,
							Amount: sdk.ZeroDec(),
						},
					},
				},
			},
		},
		"invalid one-time-mint left count": {
			data: foundation.GenesisState{
				OneTimeMintLeftCount: 2,
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
