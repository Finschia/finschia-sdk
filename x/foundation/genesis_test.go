package foundation_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	"github.com/line/lbm-sdk/testutil/testdata"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

func workingFoundation() foundation.FoundationInfo {
	return *foundation.FoundationInfo{
		Version:     1,
		TotalWeight: sdk.OneDec(),
	}.WithDecisionPolicy(workingPolicy())
}

func workingPolicy() foundation.DecisionPolicy {
	return &foundation.ThresholdDecisionPolicy{
		Threshold: sdk.OneDec(),
		Windows: &foundation.DecisionPolicyWindows{
			VotingPeriod: 7 * 24 * time.Hour, // one week
		},
	}
}

func TestDefaultGenesisState(t *testing.T) {
	gs := foundation.DefaultGenesisState()
	require.NoError(t, foundation.ValidateGenesis(*gs))

	require.True(t, gs.Params.FoundationTax.IsZero())
	require.Empty(t, gs.Params.CensoredMsgTypeUrls)

	require.EqualValues(t, 1, gs.Foundation.Version)
	require.True(t, gs.Foundation.TotalWeight.IsZero())

	require.Empty(t, gs.Members)
	require.Zero(t, gs.PreviousProposalId)
	require.Empty(t, gs.Proposals)
	require.Empty(t, gs.Votes)

	require.Empty(t, gs.Authorizations)
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
				Foundation: workingFoundation(),
				Members: []foundation.Member{
					{
						Address: createAddress().String(),
					},
				},
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
				Foundation: workingFoundation(),
				Members:    []foundation.Member{{}},
			},
		},
		"invalid foundation info": {
			data: foundation.GenesisState{
				Params: foundation.DefaultParams(),
			},
		},
		"number of members is different from total weight": {
			data: foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
				Members: []foundation.Member{
					{
						Address: createAddress().String(),
					},
				},
			},
		},
		"non empty proposals with outsourcing decision policy": {
			data: foundation.GenesisState{
				Params:             foundation.DefaultParams(),
				Foundation:         foundation.DefaultFoundation(),
				PreviousProposalId: 1,
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
				Foundation: workingFoundation(),
				Members: []foundation.Member{
					{
						Address: createAddress().String(),
					},
				},
				PreviousProposalId: 1,
				Proposals:          []foundation.Proposal{{}},
			},
		},
		"proposal of too far ahead id": {
			data: foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: workingFoundation(),
				Members: []foundation.Member{
					{
						Address: createAddress().String(),
					},
				},
				PreviousProposalId: 0,
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id:                1,
						Proposers:         []string{createAddress().String()},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{testdata.NewTestMsg()}),
				},
			},
		},
		"proposal of too far ahead version": {
			data: foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: workingFoundation(),
				Members: []foundation.Member{
					{
						Address: createAddress().String(),
					},
				},
				PreviousProposalId: 1,
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id:                1,
						Proposers:         []string{createAddress().String()},
						FoundationVersion: 2,
					}.WithMsgs([]sdk.Msg{testdata.NewTestMsg()}),
				},
			},
		},
		"duplicate proposals": {
			data: foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: workingFoundation(),
				Members: []foundation.Member{
					{
						Address: createAddress().String(),
					},
				},
				PreviousProposalId: 1,
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
				Foundation: workingFoundation(),
				Members: []foundation.Member{
					{
						Address: createAddress().String(),
					},
				},
				PreviousProposalId: 1,
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
						Option:     foundation.VOTE_OPTION_YES,
					},
				},
			},
		},
		"invalid vote option": {
			data: foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: workingFoundation(),
				Members: []foundation.Member{
					{
						Address: createAddress().String(),
					},
				},
				PreviousProposalId: 1,
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
		"invalid gov-mint left count": {
			data: foundation.GenesisState{
				Params:           foundation.DefaultParams(),
				Foundation:       foundation.DefaultFoundation(),
				GovMintLeftCount: foundation.GovMintMaxCount + 1,
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

func TestFoundationInfo(t *testing.T) {
	testCases := map[string]struct {
		version     uint64
		totalWeight sdk.Dec
		policy      foundation.DecisionPolicy
		valid       bool
	}{
		"valid info (default)": {
			version:     1,
			totalWeight: sdk.ZeroDec(),
			policy:      foundation.DefaultDecisionPolicy(),
			valid:       true,
		},
		"valid info (working policy)": {
			version:     1,
			totalWeight: sdk.OneDec(),
			policy:      workingPolicy(),
			valid:       true,
		},
		"invalid version": {
			totalWeight: sdk.ZeroDec(),
			policy:      foundation.DefaultDecisionPolicy(),
		},
		"invalid total weight": {
			version: 1,
			policy:  foundation.DefaultDecisionPolicy(),
		},
		"empty policy": {
			version:     1,
			totalWeight: sdk.ZeroDec(),
		},
		"invalid policy": {
			version:     1,
			totalWeight: sdk.ZeroDec(),
			policy:      &foundation.ThresholdDecisionPolicy{},
		},
		"outsourcing with members": {
			version:     1,
			totalWeight: sdk.OneDec(),
			policy:      foundation.DefaultDecisionPolicy(),
		},
		"working policy without members": {
			version:     1,
			totalWeight: sdk.ZeroDec(),
			policy:      workingPolicy(),
		},
	}

	for name, tc := range testCases {
		info := foundation.FoundationInfo{
			Version:     tc.version,
			TotalWeight: tc.totalWeight,
		}
		if tc.policy != nil {
			err := info.SetDecisionPolicy(tc.policy)
			require.NoError(t, err, name)
		}

		err := info.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			continue
		}
		require.NoError(t, err, name)
	}
}
