package internal_test

import (
	"github.com/line/lbm-sdk/testutil/testdata"
	sdk "github.com/line/lbm-sdk/types"
	authtypes "github.com/line/lbm-sdk/x/auth/types"
	"github.com/line/lbm-sdk/x/foundation"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

func (s *KeeperTestSuite) TestProposalHandler() {
	testCases := map[string]struct {
		malleate func(ctx sdk.Context)
		msg      sdk.Msg
		valid    bool
		require  func(ctx sdk.Context)
	}{
		"valid": {
			malleate: func(ctx sdk.Context) {
				s.impl.SetCensorship(ctx, foundation.Censorship{
					MsgTypeUrl: sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
					Authority:  foundation.CensorshipAuthorityGovernance,
				})
			},
			msg: &foundation.MsgUpdateCensorship{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Censorship: foundation.Censorship{
					MsgTypeUrl: sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
					Authority:  foundation.CensorshipAuthorityUnspecified,
				},
			},
			valid: true,
			require: func(ctx sdk.Context) {
				s.Require().False(s.impl.IsCensoredMessage(ctx, sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil))))
			},
		},
		"bad signer": {
			malleate: func(ctx sdk.Context) {
				s.impl.SetCensorship(ctx, foundation.Censorship{
					MsgTypeUrl: sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
					Authority:  foundation.CensorshipAuthorityGovernance,
				})
			},
			msg: &foundation.MsgUpdateCensorship{
				Authority: s.impl.GetAuthority(),
				Censorship: foundation.Censorship{
					MsgTypeUrl: sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
					Authority:  foundation.CensorshipAuthorityUnspecified,
				},
			},
		},
		"message type not allowed": {
			malleate: func(ctx sdk.Context) {
				s.impl.SetCensorship(ctx, foundation.Censorship{
					MsgTypeUrl: sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
					Authority:  foundation.CensorshipAuthorityGovernance,
				})
			},
			msg: newMsgCreateDog("doge"),
		},
		"no handler found": {
			malleate: func(ctx sdk.Context) {
				s.impl.SetCensorship(ctx, foundation.Censorship{
					MsgTypeUrl: sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
					Authority:  foundation.CensorshipAuthorityGovernance,
				})
			},
			msg: testdata.NewTestMsg(authtypes.NewModuleAddress(govtypes.ModuleName)),
		},
		"message execution failed": {
			malleate: func(ctx sdk.Context) {
				s.impl.SetCensorship(ctx, foundation.Censorship{
					MsgTypeUrl: sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
					Authority:  foundation.CensorshipAuthorityGovernance,
				})
			},
			msg: &foundation.MsgUpdateCensorship{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Censorship: foundation.Censorship{
					MsgTypeUrl: sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
					Authority:  foundation.CensorshipAuthorityFoundation,
				},
			},
		},
		"authority is not x/gov yet": {
			msg: &foundation.MsgUpdateCensorship{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Censorship: foundation.Censorship{
					MsgTypeUrl: sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
					Authority:  foundation.CensorshipAuthorityUnspecified,
				},
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			if tc.malleate != nil {
				tc.malleate(ctx)
			}

			proposal := &foundation.FoundationExecProposal{}
			proposal.SetMessages([]sdk.Msg{tc.msg})

			err := s.proposalHandler(ctx, proposal)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			tc.require(ctx)
		})
	}
}
