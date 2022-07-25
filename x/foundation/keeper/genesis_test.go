package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"

	"github.com/line/lbm-sdk/x/foundation"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
	"github.com/line/lbm-sdk/x/stakingplus"
)

func (s *KeeperTestSuite) TestImportExportGenesis() {
	testCases := map[string]struct {
		init   *foundation.GenesisState
		valid  bool
		export *foundation.GenesisState
	}{
		"minimal": {
			init:  &foundation.GenesisState{},
			valid: true,
			export: &foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: foundation.FoundationInfo{
					Operator:    s.keeper.GetAdmin(s.ctx).String(),
					Version:     1,
					TotalWeight: sdk.ZeroDec(),
				}.WithDecisionPolicy(foundation.DefaultDecisionPolicy(foundation.DefaultConfig())),
			},
		},
		"enabled with no create validator grantees": {
			init: &foundation.GenesisState{
				Params: &foundation.Params{
					Enabled:       true,
					FoundationTax: sdk.ZeroDec(),
				},
			},
			valid: true,
			export: &foundation.GenesisState{
				Params: &foundation.Params{
					Enabled:       true,
					FoundationTax: sdk.ZeroDec(),
				},
				Foundation: foundation.FoundationInfo{
					Operator:    s.keeper.GetAdmin(s.ctx).String(),
					Version:     1,
					TotalWeight: sdk.ZeroDec(),
				}.WithDecisionPolicy(foundation.DefaultDecisionPolicy(foundation.DefaultConfig())),
			},
		},
		"members": {
			init: &foundation.GenesisState{
				Members: []foundation.Member{
					{
						Address:       s.members[0].String(),
						Participating: true,
					},
				},
			},
			valid: true,
			export: &foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: foundation.FoundationInfo{
					Operator:    s.keeper.GetAdmin(s.ctx).String(),
					Version:     1,
					TotalWeight: sdk.OneDec(),
				}.WithDecisionPolicy(foundation.DefaultDecisionPolicy(foundation.DefaultConfig())),
				Members: []foundation.Member{
					{
						Address:       s.members[0].String(),
						Participating: true,
					},
				},
			},
		},
		"proposals": {
			init: &foundation.GenesisState{
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id:                1,
						Proposers:         []string{s.members[0].String()},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{&foundation.MsgWithdrawFromTreasury{
						Operator: s.operator.String(),
						To:       s.stranger.String(),
						Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance)),
					}}),
				},
				Votes: []foundation.Vote{
					{
						ProposalId: 1,
						Voter:      s.members[0].String(),
						Option:     foundation.VOTE_OPTION_YES,
					},
				},
			},
			valid: true,
			export: &foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: foundation.FoundationInfo{
					Operator:    s.keeper.GetAdmin(s.ctx).String(),
					Version:     1,
					TotalWeight: sdk.ZeroDec(),
				}.WithDecisionPolicy(foundation.DefaultDecisionPolicy(foundation.DefaultConfig())),
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id:                1,
						Proposers:         []string{s.members[0].String()},
						FoundationVersion: 1,
						FinalTallyResult: foundation.TallyResult{
							YesCount:        sdk.ZeroDec(),
							NoCount:         sdk.ZeroDec(),
							AbstainCount:    sdk.ZeroDec(),
							NoWithVetoCount: sdk.ZeroDec(),
						},
					}.WithMsgs([]sdk.Msg{&foundation.MsgWithdrawFromTreasury{
						Operator: s.operator.String(),
						To:       s.stranger.String(),
						Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance)),
					}}),
				},
				Votes: []foundation.Vote{
					{
						ProposalId: 1,
						Voter:      s.members[0].String(),
						Option:     foundation.VOTE_OPTION_YES,
					},
				},
			},
		},
		"authorizations": {
			init: &foundation.GenesisState{
				Authorizations: []foundation.GrantAuthorization{
					*foundation.GrantAuthorization{
						Granter: foundation.ModuleName,
						Grantee: s.stranger.String(),
					}.WithAuthorization(&foundation.ReceiveFromTreasuryAuthorization{}),
				},
			},
			valid: true,
			export: &foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: foundation.FoundationInfo{
					Operator:    s.keeper.GetAdmin(s.ctx).String(),
					Version:     1,
					TotalWeight: sdk.ZeroDec(),
				}.WithDecisionPolicy(foundation.DefaultDecisionPolicy(foundation.DefaultConfig())),
				Authorizations: []foundation.GrantAuthorization{
					*foundation.GrantAuthorization{
						Granter: foundation.ModuleName,
						Grantee: s.stranger.String(),
					}.WithAuthorization(&foundation.ReceiveFromTreasuryAuthorization{}),
				},
			},
		},
		"create validator authorizations": {
			init: &foundation.GenesisState{
				Authorizations: []foundation.GrantAuthorization{
					*foundation.GrantAuthorization{
						Granter: govtypes.ModuleName,
						Grantee: s.stranger.String(),
					}.WithAuthorization(&stakingplus.CreateValidatorAuthorization{
						ValidatorAddress: sdk.ValAddress(s.stranger).String(),
					}),
				},
			},
			valid: true,
			export: &foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: foundation.FoundationInfo{
					Operator:    s.keeper.GetAdmin(s.ctx).String(),
					Version:     1,
					TotalWeight: sdk.OneDec(),
				}.WithDecisionPolicy(foundation.DefaultDecisionPolicy(foundation.DefaultConfig())),
				Authorizations: []foundation.GrantAuthorization{
					*foundation.GrantAuthorization{
						Granter: govtypes.ModuleName,
						Grantee: s.stranger.String(),
					}.WithAuthorization(&stakingplus.CreateValidatorAuthorization{
						ValidatorAddress: sdk.ValAddress(s.stranger).String(),
					}),
				},
				Members: []foundation.Member{{
					Address:       s.stranger.String(),
					Participating: true,
					Metadata:      "genesis member",
				}},
			},
		},
		"member of long metadata": {
			init: &foundation.GenesisState{
				Members: []foundation.Member{
					{
						Address:       s.members[0].String(),
						Participating: true,
						Metadata:      string(make([]rune, 256)),
					},
				},
			},
		},
		"proposal of long metadata": {
			init: &foundation.GenesisState{
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id:                1,
						Metadata:          string(make([]rune, 256)),
						Proposers:         []string{s.members[0].String()},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{&foundation.MsgWithdrawFromTreasury{
						Operator: s.operator.String(),
						To:       s.stranger.String(),
						Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance)),
					}}),
				},
			},
		},
		"vote of long metadata": {
			init: &foundation.GenesisState{
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id:                1,
						Proposers:         []string{s.members[0].String()},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{&foundation.MsgWithdrawFromTreasury{
						Operator: s.operator.String(),
						To:       s.stranger.String(),
						Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance)),
					}}),
				},
				Votes: []foundation.Vote{
					{
						ProposalId: 1,
						Voter:      s.members[0].String(),
						Option:     foundation.VOTE_OPTION_YES,
						Metadata:   string(make([]rune, 256)),
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			s.keeper.ResetState(ctx)

			err := foundation.ValidateGenesis(*tc.init)
			s.Require().NoError(err)

			err = s.keeper.InitGenesis(ctx, s.app.StakingKeeper, tc.init)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			actual := s.keeper.ExportGenesis(ctx)
			s.Require().Equal(tc.export, actual)
		})
	}
}
