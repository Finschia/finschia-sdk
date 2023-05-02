package internal_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Finschia/finschia-sdk/crypto/keys/secp256k1"
	"github.com/Finschia/finschia-sdk/simapp"
	"github.com/Finschia/finschia-sdk/testutil/testdata"
	sdk "github.com/Finschia/finschia-sdk/types"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

func workingPolicy() foundation.DecisionPolicy {
	return &foundation.ThresholdDecisionPolicy{
		Threshold: sdk.OneDec(),
		Windows: &foundation.DecisionPolicyWindows{
			VotingPeriod: 7 * 24 * time.Hour, // one week
		},
	}
}

func TestImportExportGenesis(t *testing.T) {
	checkTx := false
	app := simapp.Setup(checkTx)
	testdata.RegisterInterfaces(app.InterfaceRegistry())

	ctx := app.BaseApp.NewContext(checkTx, tmproto.Header{})
	keeper := app.FoundationKeeper

	createAddress := func() sdk.AccAddress {
		return sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	authority := foundation.DefaultAuthority()
	existingAccount := createAddress()
	app.AccountKeeper.SetAccount(ctx, app.AccountKeeper.NewAccountWithAddress(ctx, existingAccount))

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
					TotalWeight: sdk.OneDec(),
				}.WithDecisionPolicy(workingPolicy()),
				Members: []foundation.Member{
					{
						Address: member.String(),
					},
				},
			},
			valid: true,
			export: &foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: *foundation.FoundationInfo{
					Version:     1,
					TotalWeight: sdk.OneDec(),
				}.WithDecisionPolicy(workingPolicy()),
				Members: []foundation.Member{
					{
						Address: member.String(),
					},
				},
			},
		},
		"proposals": {
			init: &foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: *foundation.FoundationInfo{
					Version:     1,
					TotalWeight: sdk.OneDec(),
				}.WithDecisionPolicy(workingPolicy()),
				Members: []foundation.Member{
					{
						Address: member.String(),
					},
				},
				PreviousProposalId: 1,
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id:                1,
						Proposers:         []string{member.String()},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{testdata.NewTestMsg(authority)}),
				},
				Votes: []foundation.Vote{
					{
						ProposalId: 1,
						Voter:      member.String(),
						Option:     foundation.VOTE_OPTION_YES,
					},
				},
			},
			valid: true,
			export: &foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: *foundation.FoundationInfo{
					Version:     1,
					TotalWeight: sdk.OneDec(),
				}.WithDecisionPolicy(workingPolicy()),
				Members: []foundation.Member{
					{
						Address: member.String(),
					},
				},
				PreviousProposalId: 1,
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id:                1,
						Proposers:         []string{member.String()},
						FoundationVersion: 1,
						FinalTallyResult: foundation.TallyResult{
							YesCount:        sdk.ZeroDec(),
							NoCount:         sdk.ZeroDec(),
							AbstainCount:    sdk.ZeroDec(),
							NoWithVetoCount: sdk.ZeroDec(),
						},
					}.WithMsgs([]sdk.Msg{testdata.NewTestMsg(authority)}),
				},
				Votes: []foundation.Vote{
					{
						ProposalId: 1,
						Voter:      member.String(),
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
						Grantee: stranger.String(),
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
						Grantee: stranger.String(),
					}.WithAuthorization(&foundation.ReceiveFromTreasuryAuthorization{}),
				},
			},
		},
		"pool": {
			init: &foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
				Pool: foundation.Pool{
					Treasury: sdk.NewDecCoins(sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.OneInt())),
				},
			},
			valid: true,
			export: &foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
				Pool: foundation.Pool{
					Treasury: sdk.NewDecCoins(sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.OneInt())),
				},
			},
		},
		"member of long metadata": {
			init: &foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: *foundation.FoundationInfo{
					Version:     1,
					TotalWeight: sdk.OneDec(),
				}.WithDecisionPolicy(workingPolicy()),
				Members: []foundation.Member{
					{
						Address:  member.String(),
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
					TotalWeight: sdk.OneDec(),
				}.WithDecisionPolicy(workingPolicy()),
				Members: []foundation.Member{
					{
						Address: member.String(),
					},
				},
				PreviousProposalId: 1,
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id:                1,
						Metadata:          string(make([]rune, 256)),
						Proposers:         []string{member.String()},
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
					TotalWeight: sdk.OneDec(),
				}.WithDecisionPolicy(workingPolicy()),
				Members: []foundation.Member{
					{
						Address: member.String(),
					},
				},
				PreviousProposalId: 1,
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id:                1,
						Proposers:         []string{member.String()},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{testdata.NewTestMsg()}),
				},
				Votes: []foundation.Vote{
					{
						ProposalId: 1,
						Voter:      member.String(),
						Option:     foundation.VOTE_OPTION_YES,
						Metadata:   string(make([]rune, 256)),
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		ctx, _ := ctx.CacheContext()

		err := foundation.ValidateGenesis(*tc.init)
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
