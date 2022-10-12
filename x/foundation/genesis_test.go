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
			data: foundation.GenesisState{
				OneTimeMintMaxCount: 1,
			},
			valid: true,
		},
		"members": {
			data: foundation.GenesisState{
				Members: []foundation.Member{
					{
						Address: createAddress().String(),
					},
				},
				OneTimeMintMaxCount: 1,
			},
			valid: true,
		},
		"foundation info": {
			data: foundation.GenesisState{
				Foundation: &foundation.FoundationInfo{
					Operator: createAddress().String(),
					Version:  1,
				},
				OneTimeMintMaxCount: 1,
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
				OneTimeMintMaxCount: 1,
			},
			valid: true,
		},
		"invalid foundation tax": {
			data: foundation.GenesisState{
				Params: &foundation.Params{
					FoundationTax: sdk.NewDec(2),
				},
				OneTimeMintMaxCount: 1,
			},
		},
		"invalid members": {
			data: foundation.GenesisState{
				Members:             []foundation.Member{{}},
				OneTimeMintMaxCount: 1,
			},
		},
		"invalid operator address": {
			data: foundation.GenesisState{
				Foundation: &foundation.FoundationInfo{
					Operator: "invalid-address",
					Version:  1,
				},
				OneTimeMintMaxCount: 1,
			},
		},
		"invalid foundation version": {
			data: foundation.GenesisState{
				Foundation:          &foundation.FoundationInfo{},
				OneTimeMintMaxCount: 1,
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
				OneTimeMintMaxCount: 1,
			},
		},
		"invalid proposals": {
			data: foundation.GenesisState{
				Proposals:           []foundation.Proposal{{}},
				OneTimeMintMaxCount: 1,
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
				OneTimeMintMaxCount: 1,
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
				OneTimeMintMaxCount: 1,
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
				OneTimeMintMaxCount: 1,
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
				OneTimeMintMaxCount: 1,
			},
		},
		"invalid authorization": {
			data: foundation.GenesisState{
				Authorizations: []foundation.GrantAuthorization{{
					Grantee: createAddress().String(),
				}},
				OneTimeMintMaxCount: 1,
			},
		},
		"invalid grantee": {
			data: foundation.GenesisState{
				Authorizations: []foundation.GrantAuthorization{
					*foundation.GrantAuthorization{}.WithAuthorization(&foundation.ReceiveFromTreasuryAuthorization{}),
				},
				OneTimeMintMaxCount: 1,
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
				OneTimeMintMaxCount: 1,
			},
		},
		"invalid one-time-mint max count": {
			data: foundation.GenesisState{
				OneTimeMintMaxCount: 0,
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
