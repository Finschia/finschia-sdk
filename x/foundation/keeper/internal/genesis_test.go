package internal_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

func workingPolicy() foundation.DecisionPolicy {
	return &foundation.ThresholdDecisionPolicy{
		Threshold: math.LegacyOneDec(),
		Windows: &foundation.DecisionPolicyWindows{
			VotingPeriod: 7 * 24 * time.Hour, // one week
		},
	}
}

func TestImportExportGenesis(t *testing.T) {
	_, keeper, _, _, _, addressCodec, ctx := setupFoundationKeeper(t, nil, nil)

	bytesToString := func(addr sdk.AccAddress) string {
		str, err := addressCodec.BytesToString(addr)
		require.NoError(t, err)
		return str
	}

	authority, err := addressCodec.StringToBytes(keeper.GetAuthority())
	require.NoError(t, err)

	createAddress := func() sdk.AccAddress {
		return sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	member := createAddress()
	stranger := createAddress()

	testCases := map[string]struct {
		init   *foundation.GenesisState
		valid  bool
		export *foundation.GenesisState
	}{
		"minimal": {
			init: &foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
			},
			valid: true,
			export: &foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
			},
		},
		"members": {
			init: &foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: *foundation.FoundationInfo{
					Version:     1,
					TotalWeight: math.LegacyOneDec(),
				}.WithDecisionPolicy(workingPolicy()),
				Members: []foundation.Member{
					{
						Address: bytesToString(member),
					},
				},
			},
			valid: true,
			export: &foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: *foundation.FoundationInfo{
					Version:     1,
					TotalWeight: math.LegacyOneDec(),
				}.WithDecisionPolicy(workingPolicy()),
				Members: []foundation.Member{
					{
						Address: bytesToString(member),
					},
				},
			},
		},
		"proposals": {
			init: &foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: *foundation.FoundationInfo{
					Version:     1,
					TotalWeight: math.LegacyOneDec(),
				}.WithDecisionPolicy(workingPolicy()),
				Members: []foundation.Member{
					{
						Address: bytesToString(member),
					},
				},
				PreviousProposalId: 1,
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id:                1,
						Proposers:         []string{bytesToString(member)},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{&testdata.TestMsg{Signers: []string{bytesToString(authority)}}}),
				},
				Votes: []foundation.Vote{
					{
						ProposalId: 1,
						Voter:      bytesToString(member),
						Option:     foundation.VOTE_OPTION_YES,
					},
				},
			},
			valid: true,
			export: &foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: *foundation.FoundationInfo{
					Version:     1,
					TotalWeight: math.LegacyOneDec(),
				}.WithDecisionPolicy(workingPolicy()),
				Members: []foundation.Member{
					{
						Address: bytesToString(member),
					},
				},
				PreviousProposalId: 1,
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id:                1,
						Proposers:         []string{bytesToString(member)},
						FoundationVersion: 1,
						FinalTallyResult: foundation.TallyResult{
							YesCount:        math.LegacyZeroDec(),
							NoCount:         math.LegacyZeroDec(),
							AbstainCount:    math.LegacyZeroDec(),
							NoWithVetoCount: math.LegacyZeroDec(),
						},
					}.WithMsgs([]sdk.Msg{&testdata.TestMsg{Signers: []string{bytesToString(authority)}}}),
				},
				Votes: []foundation.Vote{
					{
						ProposalId: 1,
						Voter:      bytesToString(member),
						Option:     foundation.VOTE_OPTION_YES,
					},
				},
			},
		},
		"authorizations": {
			init: &foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
				Censorships: []foundation.Censorship{
					{
						MsgTypeUrl: sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
						Authority:  foundation.CensorshipAuthorityFoundation,
					},
				},
				Authorizations: []foundation.GrantAuthorization{
					*foundation.GrantAuthorization{
						Grantee: bytesToString(stranger),
					}.WithAuthorization(&foundation.ReceiveFromTreasuryAuthorization{}),
				},
			},
			valid: true,
			export: &foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
				Censorships: []foundation.Censorship{
					{
						MsgTypeUrl: sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
						Authority:  foundation.CensorshipAuthorityFoundation,
					},
				},
				Authorizations: []foundation.GrantAuthorization{
					*foundation.GrantAuthorization{
						Grantee: bytesToString(stranger),
					}.WithAuthorization(&foundation.ReceiveFromTreasuryAuthorization{}),
				},
			},
		},
		"pool": {
			init: &foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
				Pool: foundation.Pool{
					Treasury: sdk.NewDecCoins(sdk.NewDecCoin(sdk.DefaultBondDenom, math.OneInt())),
				},
			},
			valid: true,
			export: &foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
				Pool: foundation.Pool{
					Treasury: sdk.NewDecCoins(sdk.NewDecCoin(sdk.DefaultBondDenom, math.OneInt())),
				},
			},
		},
		"member of long metadata": {
			init: &foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: *foundation.FoundationInfo{
					Version:     1,
					TotalWeight: math.LegacyOneDec(),
				}.WithDecisionPolicy(workingPolicy()),
				Members: []foundation.Member{
					{
						Address:  bytesToString(member),
						Metadata: string(make([]rune, 256)),
					},
				},
			},
		},
		"proposal of long metadata": {
			init: &foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: *foundation.FoundationInfo{
					Version:     1,
					TotalWeight: math.LegacyOneDec(),
				}.WithDecisionPolicy(workingPolicy()),
				Members: []foundation.Member{
					{
						Address: bytesToString(member),
					},
				},
				PreviousProposalId: 1,
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id:                1,
						Metadata:          string(make([]rune, 256)),
						Proposers:         []string{bytesToString(member)},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{testdata.NewTestMsg()}),
				},
			},
		},
		"vote of long metadata": {
			init: &foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: *foundation.FoundationInfo{
					Version:     1,
					TotalWeight: math.LegacyOneDec(),
				}.WithDecisionPolicy(workingPolicy()),
				Members: []foundation.Member{
					{
						Address: bytesToString(member),
					},
				},
				PreviousProposalId: 1,
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id:                1,
						Proposers:         []string{bytesToString(member)},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{testdata.NewTestMsg()}),
				},
				Votes: []foundation.Vote{
					{
						ProposalId: 1,
						Voter:      bytesToString(member),
						Option:     foundation.VOTE_OPTION_YES,
						Metadata:   string(make([]rune, 256)),
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		ctx, _ := ctx.CacheContext()

		err := foundation.ValidateGenesis(*tc.init, addressCodec)
		require.NoError(t, err, name)

		err = keeper.InitGenesis(ctx, tc.init)
		if !tc.valid {
			require.Error(t, err, name)
			continue
		}
		require.NoError(t, err, name)

		actual := keeper.ExportGenesis(ctx)
		require.Equal(t, tc.export, actual, name)
	}
}
