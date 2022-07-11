package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
	"github.com/line/lbm-sdk/x/stakingplus"
)

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
		operator sdk.AccAddress
		to       sdk.AccAddress
		amount   sdk.Int
		valid    bool
	}{
		"valid request": {
			operator: s.operator,
			to:       s.stranger,
			amount:   s.balance,
			valid:    true,
		},
		"operator not authorized": {
			operator: s.stranger,
			to:       s.stranger,
			amount:   s.balance,
		},
		"receiver not authorized": {
			operator: s.operator,
			to:       s.members[0],
			amount:   s.balance,
		},
		"insufficient funds": {
			operator: s.operator,
			to:       s.stranger,
			amount:   s.balance.Add(sdk.OneInt()),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &foundation.MsgWithdrawFromTreasury{
				Operator: tc.operator.String(),
				To:       tc.to.String(),
				Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, tc.amount)),
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
		operator sdk.AccAddress
		policy   foundation.DecisionPolicy
		valid    bool
	}{
		"valid request": {
			operator: s.operator,
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: foundation.DefaultConfig().MinThreshold,
				Windows:   &foundation.DecisionPolicyWindows{},
			},
			valid: true,
		},
		"not authorized": {
			operator: s.stranger,
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: foundation.DefaultConfig().MinThreshold,
				Windows:   &foundation.DecisionPolicyWindows{},
			},
		},
		"low threshold": {
			operator: s.operator,
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: sdk.OneDec(),
				Windows:   &foundation.DecisionPolicyWindows{},
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &foundation.MsgUpdateDecisionPolicy{
				Operator: tc.operator.String(),
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
		operator sdk.AccAddress
		member   foundation.Member
		valid    bool
	}{
		"valid request": {
			operator: s.operator,
			member: foundation.Member{
				Address: s.members[0].String(),
			},
			valid: true,
		},
		"not authorized": {
			operator: s.stranger,
			member: foundation.Member{
				Address: s.members[0].String(),
			},
		},
		"remove a non-member": {
			operator: s.operator,
			member: foundation.Member{
				Address: s.stranger.String(),
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &foundation.MsgUpdateMembers{
				Operator:      tc.operator.String(),
				MemberUpdates: []foundation.Member{tc.member},
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
		proposers []string
		metadata  string
		msg       sdk.Msg
		exec      foundation.Exec
		valid     bool
	}{
		"valid request (submit)": {
			proposers: members,
			msg: &foundation.MsgWithdrawFromTreasury{
				Operator: s.operator.String(),
				To:       s.stranger.String(),
				Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance)),
			},
			valid: true,
		},
		"valid request (submit & execute)": {
			proposers: members,
			msg: &foundation.MsgWithdrawFromTreasury{
				Operator: s.operator.String(),
				To:       s.stranger.String(),
				Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance)),
			},
			exec:  foundation.Exec_EXEC_TRY,
			valid: true,
		},
		"valid request (submit & unable to reach quorum)": {
			proposers: []string{members[0]},
			msg: &foundation.MsgWithdrawFromTreasury{
				Operator: s.operator.String(),
				To:       s.stranger.String(),
				Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance)),
			},
			exec:  foundation.Exec_EXEC_TRY,
			valid: true,
		},
		"not a member": {
			proposers: []string{s.stranger.String()},
			msg: &foundation.MsgWithdrawFromTreasury{
				Operator: s.operator.String(),
				To:       s.stranger.String(),
				Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance)),
			},
		},
		"unauthorized msg": {
			proposers: []string{members[0]},
			msg: &foundation.MsgWithdrawFromTreasury{
				Operator: s.stranger.String(),
				To:       s.stranger.String(),
				Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance)),
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

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
		"valid request (operator)": {
			proposalID: s.activeProposal,
			address:    s.operator,
			valid:      true,
		},
		"not authorized": {
			address: s.stranger,
		},
		"inactive proposal": {
			proposalID: s.abortedProposal,
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
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

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
		address sdk.AccAddress
		valid   bool
	}{
		"valid request": {
			address: s.members[0],
			valid:   true,
		},
		"not authorized": {
			address: s.stranger,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

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
		operator      sdk.AccAddress
		authorization foundation.Authorization
		valid         bool
	}{
		"valid request": {
			operator:      s.operator,
			authorization: &foundation.ReceiveFromTreasuryAuthorization{},
			valid:         true,
		},
		"not authorized": {
			operator:      s.stranger,
			authorization: &foundation.ReceiveFromTreasuryAuthorization{},
		},
		"wrong granter": {
			operator:      s.operator,
			authorization: &stakingplus.CreateValidatorAuthorization{},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &foundation.MsgGrant{
				Operator: tc.operator.String(),
				Grantee:  s.operator.String(),
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
		operator   sdk.AccAddress
		grantee    sdk.AccAddress
		msgTypeURL string
		valid      bool
	}{
		"valid request": {
			operator:   s.operator,
			grantee:    s.stranger,
			msgTypeURL: foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
			valid:      true,
		},
		"no grant": {
			operator:   s.operator,
			grantee:    s.members[0],
			msgTypeURL: foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
		},
		"not authorized": {
			operator:   s.stranger,
			grantee:    s.stranger,
			msgTypeURL: foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
		},
		"wrong granter": {
			operator:   s.operator,
			grantee:    s.stranger,
			msgTypeURL: stakingplus.CreateValidatorAuthorization{}.MsgTypeURL(),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &foundation.MsgRevoke{
				Operator:   tc.operator.String(),
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
