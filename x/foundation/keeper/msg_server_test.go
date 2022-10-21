package keeper_test

import (
	"time"

	"github.com/line/lbm-sdk/testutil/testdata"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
	"github.com/line/lbm-sdk/x/stakingplus"
)

func (s *KeeperTestSuite) TestMsgUpdateParams() {
	testCases := map[string]struct {
		authority sdk.AccAddress
		params    foundation.Params
		valid     bool
	}{
		"valid request": {
			authority: s.authority,
			params:    foundation.DefaultParams(),
			valid:     true,
		},
		"invalid authority": {
			authority: s.stranger,
			params:    foundation.DefaultParams(),
		},
		"enabling feature": {
			authority: s.authority,
			params: foundation.Params{
				FoundationTax: sdk.ZeroDec(),
				CensoredMsgTypeUrls: []string{
					sdk.MsgTypeURL((*testdata.TestMsg)(nil)),
				},
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &foundation.MsgUpdateParams{
				Authority: tc.authority.String(),
				Params:    tc.params,
			}
			res, err := s.msgServer.UpdateParams(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgFundTreasury() {
	testCases := map[string]struct {
		amount sdk.Int
		valid  bool
	}{
		"valid request": {
			amount: s.balance,
			valid:  true,
		},
		"insufficient funds": {
			amount: s.balance.Add(sdk.OneInt()),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &foundation.MsgFundTreasury{
				From:   s.stranger.String(),
				Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, tc.amount)),
			}
			res, err := s.msgServer.FundTreasury(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgWithdrawFromTreasury() {
	testCases := map[string]struct {
		authority sdk.AccAddress
		to        sdk.AccAddress
		amount    sdk.Int
		valid     bool
	}{
		"valid request": {
			authority: s.authority,
			to:        s.stranger,
			amount:    s.balance,
			valid:     true,
		},
		"authority not authorized": {
			authority: s.stranger,
			to:        s.stranger,
			amount:    s.balance,
		},
		"receiver not authorized": {
			authority: s.authority,
			to:        s.members[0],
			amount:    s.balance,
		},
		"insufficient funds": {
			authority: s.authority,
			to:        s.stranger,
			amount:    s.balance.Add(sdk.OneInt()),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &foundation.MsgWithdrawFromTreasury{
				Authority: tc.authority.String(),
				To:        tc.to.String(),
				Amount:    sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, tc.amount)),
			}
			res, err := s.msgServer.WithdrawFromTreasury(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgUpdateDecisionPolicy() {
	testCases := map[string]struct {
		authority sdk.AccAddress
		policy    foundation.DecisionPolicy
		valid     bool
	}{
		"valid request": {
			authority: s.authority,
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: sdk.OneDec(),
				Windows:   &foundation.DecisionPolicyWindows{},
			},
			valid: true,
		},
		"not authorized": {
			authority: s.stranger,
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: sdk.OneDec(),
				Windows:   &foundation.DecisionPolicyWindows{},
			},
		},
		"invalid policy": {
			authority: s.authority,
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: sdk.OneDec(),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod:       time.Hour,
					MinExecutionPeriod: foundation.DefaultConfig().MaxExecutionPeriod + time.Hour,
				},
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &foundation.MsgUpdateDecisionPolicy{
				Authority: tc.authority.String(),
			}
			err := req.SetDecisionPolicy(tc.policy)
			s.Require().NoError(err)

			res, err := s.msgServer.UpdateDecisionPolicy(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgUpdateMembers() {
	testCases := map[string]struct {
		authority sdk.AccAddress
		member    foundation.MemberRequest
		valid     bool
	}{
		"valid request": {
			authority: s.authority,
			member: foundation.MemberRequest{
				Address: s.members[0].String(),
			},
			valid: true,
		},
		"not authorized": {
			authority: s.stranger,
			member: foundation.MemberRequest{
				Address: s.members[0].String(),
			},
		},
		"remove a non-member": {
			authority: s.authority,
			member: foundation.MemberRequest{
				Address: s.stranger.String(),
				Remove:  true,
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &foundation.MsgUpdateMembers{
				Authority:     tc.authority.String(),
				MemberUpdates: []foundation.MemberRequest{tc.member},
			}
			res, err := s.msgServer.UpdateMembers(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgSubmitProposal() {
	members := make([]string, len(s.members))
	for i, member := range s.members {
		members[i] = member.String()
	}

	testCases := map[string]struct {
		malleate  func(ctx sdk.Context)
		proposers []string
		metadata  string
		msg       sdk.Msg
		exec      foundation.Exec
		valid     bool
	}{
		"valid request (submit)": {
			proposers: members,
			msg:       testdata.NewTestMsg(s.authority),
			valid:     true,
		},
		"valid request (submit & execute)": {
			proposers: members,
			msg:       testdata.NewTestMsg(s.authority),
			exec:      foundation.Exec_EXEC_TRY,
			valid:     true,
		},
		"valid request (submit & unable to reach quorum)": {
			proposers: []string{members[0]},
			msg:       testdata.NewTestMsg(s.authority),
			exec:      foundation.Exec_EXEC_TRY,
			valid:     true,
		},
		"not a member": {
			proposers: []string{s.stranger.String()},
			msg:       testdata.NewTestMsg(s.authority),
		},
		"unauthorized msg": {
			proposers: []string{members[0]},
			msg:       testdata.NewTestMsg(s.stranger),
		},
		"exec fails": {
			malleate: func(ctx sdk.Context) {
				// try exec will fail because of a non-zero MinExecutionPeriod.
				err := s.keeper.UpdateDecisionPolicy(ctx, &foundation.ThresholdDecisionPolicy{
					Threshold: sdk.OneDec(),
					Windows: &foundation.DecisionPolicyWindows{
						VotingPeriod:       time.Hour,
						MinExecutionPeriod: time.Second,
					},
				})
				s.Require().NoError(err)
			},
			proposers: members,
			msg:       testdata.NewTestMsg(s.authority),
			exec:      foundation.Exec_EXEC_TRY,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			if tc.malleate != nil {
				tc.malleate(ctx)
			}

			req := &foundation.MsgSubmitProposal{
				Proposers: tc.proposers,
				Metadata:  tc.metadata,
				Exec:      tc.exec,
			}
			err := req.SetMsgs([]sdk.Msg{tc.msg})
			s.Require().NoError(err)

			res, err := s.msgServer.SubmitProposal(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgWithdrawProposal() {
	testCases := map[string]struct {
		proposalID uint64
		address    sdk.AccAddress
		valid      bool
	}{
		"valid request (proposer)": {
			proposalID: s.activeProposal,
			address:    s.members[0],
			valid:      true,
		},
		"valid request (authority)": {
			proposalID: s.activeProposal,
			address:    s.authority,
			valid:      true,
		},
		"not authorized": {
			proposalID: s.activeProposal,
			address:    s.stranger,
		},
		"inactive proposal": {
			proposalID: s.withdrawnProposal,
			address:    s.members[0],
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &foundation.MsgWithdrawProposal{
				ProposalId: tc.proposalID,
				Address:    tc.address.String(),
			}
			res, err := s.msgServer.WithdrawProposal(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgVote() {
	testCases := map[string]struct {
		malleate   func(ctx sdk.Context)
		proposalID uint64
		voter      sdk.AccAddress
		msg        sdk.Msg
		exec       foundation.Exec
		valid      bool
	}{
		"valid request (vote)": {
			proposalID: s.activeProposal,
			voter:      s.members[0],
			valid:      true,
		},
		"valid request (vote & execute)": {
			proposalID: s.activeProposal,
			voter:      s.members[0],
			exec:       foundation.Exec_EXEC_TRY,
			valid:      true,
		},
		"not authorized": {
			proposalID: s.activeProposal,
			voter:      s.stranger,
		},
		"already voted": {
			proposalID: s.votedProposal,
			voter:      s.members[0],
		},
		"exec fails": {
			malleate: func(ctx sdk.Context) {
				// try exec will fail because of a non-zero MinExecutionPeriod.
				err := s.keeper.UpdateDecisionPolicy(ctx, &foundation.ThresholdDecisionPolicy{
					Threshold: sdk.OneDec(),
					Windows: &foundation.DecisionPolicyWindows{
						VotingPeriod:       time.Hour,
						MinExecutionPeriod: time.Second,
					},
				})
				s.Require().NoError(err)

				// submit a proposal
				proposers := make([]string, len(s.members))
				for i, member := range s.members {
					proposers[i] = member.String()
				}
				req := &foundation.MsgSubmitProposal{
					Proposers: proposers,
				}
				err = req.SetMsgs([]sdk.Msg{testdata.NewTestMsg(s.authority)})
				s.Require().NoError(err)

				res, err := s.msgServer.SubmitProposal(sdk.WrapSDKContext(ctx), req)
				s.Require().NoError(err)
				s.Require().NotNil(res)

				s.Require().Equal(s.nextProposal, res.ProposalId)
			},
			proposalID: s.nextProposal,
			voter:      s.members[0],
			exec:       foundation.Exec_EXEC_TRY,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			if tc.malleate != nil {
				tc.malleate(ctx)
			}

			req := &foundation.MsgVote{
				ProposalId: tc.proposalID,
				Voter:      tc.voter.String(),
				Option:     foundation.VOTE_OPTION_YES,
				Exec:       tc.exec,
			}
			res, err := s.msgServer.Vote(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgExec() {
	testCases := map[string]struct {
		malleate   func(ctx sdk.Context)
		proposalID uint64
		signer     sdk.AccAddress
		valid      bool
	}{
		"valid request (execute)": {
			proposalID: s.votedProposal,
			signer:     s.members[0],
			valid:      true,
		},
		"valid request (not finalized)": {
			proposalID: s.activeProposal,
			signer:     s.members[0],
			valid:      true,
		},
		"not authorized": {
			proposalID: s.votedProposal,
			signer:     s.stranger,
		},
		"no such a proposal": {
			proposalID: s.nextProposal,
			signer:     s.members[0],
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &foundation.MsgExec{
				ProposalId: tc.proposalID,
				Signer:     tc.signer.String(),
			}
			res, err := s.msgServer.Exec(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgLeaveFoundation() {
	testCases := map[string]struct {
		malleate func(ctx sdk.Context)
		address  sdk.AccAddress
		valid    bool
	}{
		"valid request": {
			address: s.members[0],
			valid:   true,
		},
		"not authorized": {
			address: s.stranger,
		},
		"policy violation": {
			malleate: func(ctx sdk.Context) {
				// remove all members but the first one
				// preset policy is a threshold policy with its threshold 1
				requests := make([]foundation.MemberRequest, len(s.members)-1)
				for i, member := range s.members[1:] {
					requests[i] = foundation.MemberRequest{
						Address: member.String(),
						Remove:  true,
					}
				}
				err := s.keeper.UpdateMembers(ctx, requests)
				s.Require().NoError(err)
			},
			address: s.members[0],
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			if tc.malleate != nil {
				tc.malleate(ctx)
			}

			req := &foundation.MsgLeaveFoundation{
				Address: tc.address.String(),
			}
			res, err := s.msgServer.LeaveFoundation(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgGrant() {
	testCases := map[string]struct {
		authority     sdk.AccAddress
		grantee       sdk.AccAddress
		authorization foundation.Authorization
		valid         bool
	}{
		"valid request": {
			authority:     s.authority,
			grantee:       s.members[0],
			authorization: &foundation.ReceiveFromTreasuryAuthorization{},
			valid:         true,
		},
		"not authorized": {
			authority:     s.stranger,
			grantee:       s.members[0],
			authorization: &foundation.ReceiveFromTreasuryAuthorization{},
		},
		"already granted": {
			authority:     s.authority,
			grantee:       s.stranger,
			authorization: &foundation.ReceiveFromTreasuryAuthorization{},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &foundation.MsgGrant{
				Authority: tc.authority.String(),
				Grantee:   tc.grantee.String(),
			}
			err := req.SetAuthorization(tc.authorization)
			s.Require().NoError(err)

			res, err := s.msgServer.Grant(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgRevoke() {
	testCases := map[string]struct {
		authority  sdk.AccAddress
		grantee    sdk.AccAddress
		msgTypeURL string
		valid      bool
	}{
		"valid request": {
			authority:  s.authority,
			grantee:    s.stranger,
			msgTypeURL: foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
			valid:      true,
		},
		"no grant": {
			authority:  s.authority,
			grantee:    s.members[0],
			msgTypeURL: foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
		},
		"not authorized": {
			authority:  s.stranger,
			grantee:    s.stranger,
			msgTypeURL: foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
		},
		"wrong granter": {
			authority:  s.authority,
			grantee:    s.stranger,
			msgTypeURL: stakingplus.CreateValidatorAuthorization{}.MsgTypeURL(),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &foundation.MsgRevoke{
				Authority:  tc.authority.String(),
				Grantee:    tc.grantee.String(),
				MsgTypeUrl: tc.msgTypeURL,
			}
			res, err := s.msgServer.Revoke(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgGovMint() {
	testCases := map[string]struct {
		authority      sdk.AccAddress
		amount         sdk.Coins
		emptyCountTest bool
		valid          bool
	}{
		"valid request": {
			authority: s.authority,
			amount:    sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))),
			valid:     true,
		},
		"empty count": {
			authority:      s.authority,
			amount:         sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))),
			emptyCountTest: true,
		},
		"not authorized": {
			authority: s.stranger,
			amount:    sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			if tc.emptyCountTest {
				s.keeper.SetGovMintLeftCount(ctx, 0)
			}

			req := &foundation.MsgGovMint{
				Authority: tc.authority.String(),
				Amount:    tc.amount,
			}
			res, err := s.msgServer.GovMint(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}
