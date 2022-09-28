package foundation_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	"github.com/line/lbm-sdk/testutil/testdata"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

func TestDefaultGenesisState(t *testing.T) {
	gs := foundation.DefaultGenesisState()
	require.Equal(t, sdk.ZeroDec(), gs.Params.FoundationTax)
	require.Empty(t, gs.Params.CensoredMsgTypeUrls)
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
			data: foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
			},
			valid: true,
		},
		"members": {
			data: foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
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
				Params: foundation.DefaultParams(),
				Foundation: *foundation.FoundationInfo{
					Operator: createAddress().String(),
					Version:  1,
				}.WithDecisionPolicy(&foundation.OutsourcingDecisionPolicy{}),
			},
			valid: true,
		},
		"authorizations": {
			data: foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
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
				Params: foundation.Params{
					FoundationTax: sdk.NewDec(2),
				},
				Foundation: foundation.DefaultFoundation(),
			},
		},
		"invalid members": {
			data: foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
				Members:    []foundation.Member{{}},
			},
		},
		"invalid operator address": {
			data: foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: *foundation.FoundationInfo{
					Operator: "invalid-address",
					Version:  1,
				}.WithDecisionPolicy(foundation.DefaultDecisionPolicy()),
			},
		},
		"invalid foundation version": {
			data: foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: *foundation.FoundationInfo{}.WithDecisionPolicy(foundation.DefaultDecisionPolicy()),
			},
		},
		"no decision policy": {
			data: foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: foundation.FoundationInfo{
					Version: 1,
				},
			},
		},
		"invalid decision policy": {
			data: foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: *foundation.FoundationInfo{
					Version: 1,
				}.WithDecisionPolicy(&foundation.ThresholdDecisionPolicy{
					Windows: &foundation.DecisionPolicyWindows{},
				}),
			},
		},
		"non empty members with outsourcing decision policy": {
			data: foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: *foundation.FoundationInfo{
					Operator: createAddress().String(),
					Version:  1,
				}.WithDecisionPolicy(&foundation.OutsourcingDecisionPolicy{
					Description: "using x/group",
				}),
				Members: []foundation.Member{
					{
						Address: createAddress().String(),
					},
				},
			},
		},
		"non empty proposals with outsourcing decision policy": {
			data: foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: *foundation.FoundationInfo{
					Operator: createAddress().String(),
					Version:  1,
				}.WithDecisionPolicy(&foundation.OutsourcingDecisionPolicy{
					Description: "using x/group",
				}),
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id:                1,
						Proposers:         []string{createAddress().String()},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{testdata.NewTestMsg()}),
				},
			},
		},
		"invalid proposal": {
			data: foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
				Proposals:  []foundation.Proposal{{}},
			},
		},
		"duplicate proposals": {
			data: foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id:                1,
						Proposers:         []string{createAddress().String()},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{testdata.NewTestMsg()}),
					*foundation.Proposal{
						Id:                1,
						Proposers:         []string{createAddress().String()},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{testdata.NewTestMsg()}),
				},
			},
		},
		"no proposal for the vote": {
			data: foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
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
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id:                1,
						Proposers:         []string{createAddress().String()},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{testdata.NewTestMsg()}),
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
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id:                1,
						Proposers:         []string{createAddress().String()},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{testdata.NewTestMsg()}),
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
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
				Authorizations: []foundation.GrantAuthorization{{
					Grantee: createAddress().String(),
				}},
			},
		},
		"invalid grantee": {
			data: foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
				Authorizations: []foundation.GrantAuthorization{
					*foundation.GrantAuthorization{}.WithAuthorization(&foundation.ReceiveFromTreasuryAuthorization{}),
				},
			},
		},
		"invalid pool": {
			data: foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
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
