package keeper_test

import (
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	"github.com/line/lbm-sdk/simapp"
	"github.com/line/lbm-sdk/testutil/testdata"
	sdk "github.com/line/lbm-sdk/types"

	"github.com/line/lbm-sdk/x/foundation"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
	"github.com/line/lbm-sdk/x/stakingplus"
)

func TestImportExportGenesis(t *testing.T) {
	checkTx := false
	app := simapp.Setup(checkTx)
	testdata.RegisterInterfaces(app.InterfaceRegistry())

	ctx := app.BaseApp.NewContext(checkTx, ocproto.Header{})
	keeper := app.FoundationKeeper

	createAddress := func() sdk.AccAddress {
		return sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	operator := foundation.DefaultOperator()
	member := createAddress()
	stranger := createAddress()

	validator := createAddress()
	app.StakingKeeper.SetValidator(ctx, stakingtypes.Validator{
		OperatorAddress: sdk.ValAddress(validator).String(),
	})

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
				Params: foundation.DefaultParams(),
				Foundation: *foundation.FoundationInfo{
					Operator:    operator.String(),
					Version:     1,
					TotalWeight: sdk.ZeroDec(),
				}.WithDecisionPolicy(foundation.DefaultDecisionPolicy()),
			},
		},
		"members": {
			init: &foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
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
					Operator:    operator.String(),
					Version:     1,
					TotalWeight: sdk.OneDec(),
				}.WithDecisionPolicy(foundation.DefaultDecisionPolicy()),
				Members: []foundation.Member{
					{
						Address: member.String(),
					},
				},
			},
		},
		"proposals": {
			init: &foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
				Proposals: []foundation.Proposal{
					*foundation.Proposal{
						Id:                1,
						Proposers:         []string{member.String()},
						FoundationVersion: 1,
					}.WithMsgs([]sdk.Msg{testdata.NewTestMsg(operator)}),
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
					Operator:    operator.String(),
					Version:     1,
					TotalWeight: sdk.ZeroDec(),
				}.WithDecisionPolicy(foundation.DefaultDecisionPolicy()),
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
					}.WithMsgs([]sdk.Msg{testdata.NewTestMsg(operator)}),
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
				Authorizations: []foundation.GrantAuthorization{
					*foundation.GrantAuthorization{
						Grantee: stranger.String(),
					}.WithAuthorization(&foundation.ReceiveFromTreasuryAuthorization{}),
				},
			},
			valid: true,
			export: &foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: *foundation.FoundationInfo{
					Operator:    operator.String(),
					Version:     1,
					TotalWeight: sdk.ZeroDec(),
				}.WithDecisionPolicy(foundation.DefaultDecisionPolicy()),
				Authorizations: []foundation.GrantAuthorization{
					*foundation.GrantAuthorization{
						Grantee: stranger.String(),
					}.WithAuthorization(&foundation.ReceiveFromTreasuryAuthorization{}),
				},
			},
		},
		"genesis members from the grantees": {
			init: &foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
				Authorizations: []foundation.GrantAuthorization{
					*foundation.GrantAuthorization{
						Grantee: stranger.String(),
					}.WithAuthorization(&stakingplus.CreateValidatorAuthorization{
						ValidatorAddress: sdk.ValAddress(stranger).String(),
					}),
				},
			},
			valid: true,
			export: &foundation.GenesisState{
				Params: foundation.DefaultParams(),
				Foundation: *foundation.FoundationInfo{
					Operator:    operator.String(),
					Version:     1,
					TotalWeight: sdk.OneDec(),
				}.WithDecisionPolicy(foundation.DefaultDecisionPolicy()),
				Authorizations: []foundation.GrantAuthorization{
					*foundation.GrantAuthorization{
						Grantee: stranger.String(),
					}.WithAuthorization(&stakingplus.CreateValidatorAuthorization{
						ValidatorAddress: sdk.ValAddress(stranger).String(),
					}),
				},
				Members: []foundation.Member{{
					Address:  stranger.String(),
					Metadata: "genesis member",
				}},
			},
		},
		"grantees from the validators": {
			init: &foundation.GenesisState{
				Params: foundation.Params{
					FoundationTax: sdk.ZeroDec(),
					CensoredMsgTypeUrls: []string{
						sdk.MsgTypeURL((*stakingtypes.MsgCreateValidator)(nil)),
					},
				},
				Foundation: *foundation.FoundationInfo{
					Operator: operator.String(),
					Version:  1,
				}.WithDecisionPolicy(&foundation.OutsourcingDecisionPolicy{
					Description: "using x/group",
				}),
			},
			valid: true,
			export: &foundation.GenesisState{
				Params: foundation.Params{
					FoundationTax: sdk.ZeroDec(),
					CensoredMsgTypeUrls: []string{
						sdk.MsgTypeURL((*stakingtypes.MsgCreateValidator)(nil)),
					},
				},
				Foundation: *foundation.FoundationInfo{
					Operator:    operator.String(),
					Version:     1,
					TotalWeight: sdk.ZeroDec(),
				}.WithDecisionPolicy(&foundation.OutsourcingDecisionPolicy{
					Description: "using x/group",
				}),
				Authorizations: []foundation.GrantAuthorization{
					*foundation.GrantAuthorization{
						Grantee: validator.String(),
					}.WithAuthorization(&stakingplus.CreateValidatorAuthorization{
						ValidatorAddress: sdk.ValAddress(validator).String(),
					}),
				},
			},
		},
		"member of long metadata": {
			init: &foundation.GenesisState{
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
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
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
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
				Params:     foundation.DefaultParams(),
				Foundation: foundation.DefaultFoundation(),
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

		err = keeper.InitGenesis(ctx, app.StakingKeeper, tc.init)
		if !tc.valid {
			require.Error(t, err, name)
			continue
		}
		require.NoError(t, err, name)

		actual := keeper.ExportGenesis(ctx)
		require.Equal(t, tc.export, actual, name)
	}
}
