package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

func (s *KeeperTestSuite) TestMsgFundTreasury() {
	testCases := map[string]struct {
		amount sdk.Int
		valid bool
	}{
		"valid request": {
			amount: sdk.OneInt(),
			valid: true,
		},
		"insufficient funds": {
			amount: s.balance.Add(s.balance),
			valid: false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &foundation.MsgFundTreasury{
				From:    s.stranger.String(),
				Amount:  sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, tc.amount)),
			}

			res, err := s.msgServer.FundTreasury(s.goCtx, req)
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
		amount sdk.Int
		valid bool
	}{
		"valid request": {
			operator: s.operator,
			amount: sdk.OneInt(),
			valid: true,
		},
		"not authorized": {
			operator: s.stranger,
			amount: sdk.OneInt(),
			valid: false,
		},
		"insufficient funds": {
			operator: s.operator,
			amount: s.balance.Add(s.balance),
			valid: false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &foundation.MsgWithdrawFromTreasury{
				Operator: tc.operator.String(),
				To: s.stranger.String(),
				Amount:  sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, tc.amount)),
			}

			res, err := s.msgServer.WithdrawFromTreasury(s.goCtx, req)
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
		member foundation.Member
		valid bool
	}{
		"valid request (add member)": {
			operator: s.operator,
			member: foundation.Member{
				Address: s.comingMember.String(),
				Weight: sdk.OneDec(),
			},
			valid: true,
		},
		"valid request (remove member)": {
			operator: s.operator,
			member: foundation.Member{
				Address: s.badMember.String(),
				Weight: sdk.ZeroDec(),
			},
			valid: true,
		},
		"not authorized": {
			operator: s.stranger,
			member: foundation.Member{
				Address: s.member.String(),
				Weight: sdk.ZeroDec(),
			},
			valid: false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &foundation.MsgUpdateMembers{
				Operator: tc.operator.String(),
				MemberUpdates: []foundation.Member{tc.member},
			}

			res, err := s.msgServer.UpdateMembers(s.goCtx, req)
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
		policy foundation.DecisionPolicy
		valid bool
	}{
		"valid request (threshold)": {
			operator: s.operator,
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: sdk.NewDec(10),
				Windows: &foundation.DecisionPolicyWindows{},
			},
			valid: true,
		},
		"valid request (percentage)": {
			operator: s.operator,
			policy: &foundation.PercentageDecisionPolicy{
				Percentage: sdk.NewDec(1),
				Windows: &foundation.DecisionPolicyWindows{},
			},
			valid: true,
		},
		"not authorized": {
			operator: s.stranger,
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: sdk.NewDec(10),
				Windows: &foundation.DecisionPolicyWindows{},
			},
			valid: false,
		},
		"low threshold": {
			operator: s.operator,
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: sdk.NewDec(2),
				Windows: &foundation.DecisionPolicyWindows{},
			},
			valid: false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &foundation.MsgUpdateDecisionPolicy{
				Operator: tc.operator.String(),
			}

			err := req.SetDecisionPolicy(tc.policy)
			s.Require().NoError(err)
			res, err := s.msgServer.UpdateDecisionPolicy(s.goCtx, req)
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
	getMembers := func(ctx sdk.Context) []sdk.AccAddress {
		var members []sdk.AccAddress
		for _, member := range s.keeper.GetMembers(ctx) {
			members = append(members, sdk.AccAddress(member.Address))
		}
		return members
	}

	testCases := map[string]struct {
		proposers []sdk.AccAddress
		msg sdk.Msg
		exec foundation.Exec
		valid bool
	}{
		"valid request (submit)": {
			proposers: getMembers(s.ctx),
			msg: &foundation.MsgWithdrawFromTreasury{
				Operator: s.operator.String(),
				To: s.stranger.String(),
				Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			},
			valid: true,
		},
		"valid request (submit & execute)": {
			proposers: getMembers(s.ctx),
			msg: &foundation.MsgWithdrawFromTreasury{
				Operator: s.operator.String(),
				To: s.stranger.String(),
				Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			},
			exec: foundation.Exec_EXEC_TRY,
			valid: true,
		},
		"valid request (submit & execute fail)": {
			proposers: []sdk.AccAddress{s.member},
			msg: &foundation.MsgWithdrawFromTreasury{
				Operator: s.operator.String(),
				To: s.stranger.String(),
				Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			},
			exec: foundation.Exec_EXEC_TRY,
			valid: true,
		},
		"not authorized": {
			proposers: []sdk.AccAddress{s.stranger},
			msg: &foundation.MsgWithdrawFromTreasury{
				Operator: s.operator.String(),
				To: s.stranger.String(),
				Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			},
			valid: false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			var proposers []string
			for _, proposer := range tc.proposers {
				proposers = append(proposers, proposer.String())
			}
			req := &foundation.MsgSubmitProposal{
				Proposers: proposers,
				Exec: tc.exec,
			}
			err := req.SetMsgs([]sdk.Msg{tc.msg})
			s.Require().NoError(err)

			res, err := s.msgServer.SubmitProposal(s.goCtx, req)
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
		address sdk.AccAddress
		valid bool
	}{
		"valid request (proposer)": {
			address: s.member,
			valid: true,
		},
		"valid request (operator)": {
			address: s.operator,
			valid: true,
		},
		"not authorized": {
			address: s.stranger,
			valid: false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			// submit a proposal first
			proposal := &foundation.MsgSubmitProposal{Proposers: []string{s.member.String()}}
			err := proposal.SetMsgs([]sdk.Msg{&foundation.MsgWithdrawFromTreasury{
				Operator: s.operator.String(),
				To: s.stranger.String(),
				Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			}})
			s.Require().NoError(err)

			proposalRes, err := s.msgServer.SubmitProposal(s.goCtx, proposal)
			s.Require().NoError(err)

			proposalId := proposalRes.ProposalId

			// withdraw the proposal
			req := &foundation.MsgWithdrawProposal{
				ProposalId: proposalId,
				Address: tc.address.String(),
			}

			res, err := s.msgServer.WithdrawProposal(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)

			// double withdraw which fails
			// it feeds "already invalidated"
			_, err = s.msgServer.WithdrawProposal(s.goCtx, req)
			s.Require().Error(err)
		})
	}
}

func (s *KeeperTestSuite) TestMsgVote() {
	testCases := map[string]struct {
		voter sdk.AccAddress
		msg sdk.Msg
		exec foundation.Exec
		valid bool
	}{
		"valid request (vote)": {
			voter: s.member,
			valid: true,
		},
		"valid request (vote & execute)": {
			voter: s.member,
			exec: foundation.Exec_EXEC_TRY,
			valid: true,
		},
		"not authorized": {
			voter: s.stranger,
			valid: false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			// submit a proposal first
			proposal := &foundation.MsgSubmitProposal{Proposers: []string{s.member.String()}}
			err := proposal.SetMsgs([]sdk.Msg{&foundation.MsgWithdrawFromTreasury{
				Operator: s.operator.String(),
				To: s.stranger.String(),
				Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			}})
			s.Require().NoError(err)

			proposalRes, err := s.msgServer.SubmitProposal(s.goCtx, proposal)
			s.Require().NoError(err)

			proposalId := proposalRes.ProposalId

			// members except the voter vote first
			voters := map[string]bool{}
			for _, voter := range s.keeper.GetMembers(s.ctx) {
				voters[voter.Address] = true
			}
			delete(voters, tc.voter.String())
			for voter := range voters {
				s.msgServer.Vote(s.goCtx, &foundation.MsgVote{
					ProposalId: proposalId,
					Voter: voter,
					Option: foundation.VOTE_OPTION_YES,
				})
			}

			// do the test on the subject
			req := &foundation.MsgVote{
				ProposalId: proposalId,
				Voter: tc.voter.String(),
				Option: foundation.VOTE_OPTION_YES,
				Exec: tc.exec,
			}

			res, err := s.msgServer.Vote(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)

			// double vote which fails
			// it feeds "no proposal" and "already voted"
			req.Option = foundation.VOTE_OPTION_YES
			_, err = s.msgServer.Vote(s.goCtx, req)
			s.Require().Error(err)
		})
	}
}

func (s *KeeperTestSuite) TestMsgExec() {
	testCases := map[string]struct {
		signer sdk.AccAddress
		voteOption foundation.VoteOption
		valid bool
	}{
		"valid request (execute)": {
			signer: s.member,
			voteOption: foundation.VOTE_OPTION_YES,
			valid: true,
		},
		"valid request (not finalized)": {
			signer: s.member,
			valid: true,
		},
		"valid request (rejected)": {
			signer: s.member,
			voteOption: foundation.VOTE_OPTION_NO,
			valid: true,
		},
		"not authorized": {
			signer: s.stranger,
			voteOption: foundation.VOTE_OPTION_YES,
			valid: false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			// submit a proposal first
			proposal := &foundation.MsgSubmitProposal{Proposers: []string{s.member.String()}}
			err := proposal.SetMsgs([]sdk.Msg{&foundation.MsgWithdrawFromTreasury{
				Operator: s.operator.String(),
				To: s.stranger.String(),
				Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			}})
			s.Require().NoError(err)

			proposalRes, err := s.msgServer.SubmitProposal(s.goCtx, proposal)
			s.Require().NoError(err)

			proposalId := proposalRes.ProposalId

			// all members vote first
			for _, voter := range s.keeper.GetMembers(s.ctx) {
				s.msgServer.Vote(s.goCtx, &foundation.MsgVote{
					ProposalId: proposalId,
					Voter: voter.String(),
					Option: tc.voteOption,
				})
			}

			// do the test on the subject
			req := &foundation.MsgExec{
				ProposalId: proposalId,
				Signer: tc.signer.String(),
			}

			res, err := s.msgServer.Exec(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)

			// double exec which fails
			// it feeds "no proposal"
			_, err = s.msgServer.Exec(s.goCtx, req)
			s.Require().Error(err)
		})
	}
}

func (s *KeeperTestSuite) TestMsgLeaveFoundation() {
	testCases := map[string]struct {
		address sdk.AccAddress
		valid bool
	}{
		"valid request": {
			address: s.leavingMember,
			valid: true,
		},
		"not authorized": {
			address: s.stranger,
			valid: false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &foundation.MsgLeaveFoundation{
				Address: tc.address.String(),
			}

			res, err := s.msgServer.LeaveFoundation(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}
