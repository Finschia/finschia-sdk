package internal_test

import (
	"time"

	abci "github.com/cometbft/cometbft/abci/types"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

func (s *KeeperTestSuite) TestMsgFundTreasury() {
	testCases := map[string]struct {
		from   sdk.AccAddress
		amount math.Int
		valid  bool
		events sdk.Events
	}{
		"valid request": {
			from:   s.stranger,
			amount: s.balance,
			valid:  true,
			events: sdk.Events{{Type: "lbm.foundation.v1.EventFundTreasury", Attributes: []abci.EventAttribute{{Key: "amount", Value: "[{\"denom\":\"stake\",\"amount\":\"987654321\"}]", Index: false}, {Key: "from", Value: "\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgsjpha7m\"", Index: false}}}},
		},
		"empty from": {
			amount: s.balance,
		},
		"zero amount": {
			from:   s.stranger,
			amount: math.ZeroInt(),
		},
		"insufficient funds": {
			from:   s.stranger,
			amount: s.balance.Add(math.OneInt()),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &foundation.MsgFundTreasury{
				From:   s.bytesToString(tc.from),
				Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, tc.amount)),
			}
			res, err := s.msgServer.FundTreasury(ctx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			s.Require().NotNil(res)

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgWithdrawFromTreasury() {
	testCases := map[string]struct {
		authority sdk.AccAddress
		to        sdk.AccAddress
		amount    math.Int
		valid     bool
		events    sdk.Events
	}{
		"valid request": {
			authority: s.authority,
			to:        s.stranger,
			amount:    s.balance,
			valid:     true,
			events:    sdk.Events{{Type: "lbm.foundation.v1.EventWithdrawFromTreasury", Attributes: []abci.EventAttribute{{Key: "amount", Value: "[{\"denom\":\"stake\",\"amount\":\"987654321\"}]", Index: false}, {Key: "to", Value: "\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgsjpha7m\"", Index: false}}}},
		},
		"empty authority": {
			to:     s.stranger,
			amount: math.OneInt(),
		},
		"empty to": {
			authority: s.authority,
			amount:    math.OneInt(),
		},
		"zero amount": {
			authority: s.authority,
			to:        s.stranger,
			amount:    math.ZeroInt(),
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
			amount:    s.balance.Add(math.OneInt()),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &foundation.MsgWithdrawFromTreasury{
				Authority: s.bytesToString(tc.authority),
				To:        s.bytesToString(tc.to),
				Amount:    sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, tc.amount)),
			}
			res, err := s.msgServer.WithdrawFromTreasury(ctx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			s.Require().NotNil(res)

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgUpdateDecisionPolicy() {
	testCases := map[string]struct {
		authority sdk.AccAddress
		policy    foundation.DecisionPolicy
		valid     bool
		events    sdk.Events
	}{
		"valid threshold policy": {
			authority: s.authority,
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: math.LegacyOneDec(),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
			valid:  true,
			events: sdk.Events{{Type: "lbm.foundation.v1.EventUpdateDecisionPolicy", Attributes: []abci.EventAttribute{{Key: "decision_policy", Value: "{\"@type\":\"/lbm.foundation.v1.ThresholdDecisionPolicy\",\"threshold\":\"1.000000000000000000\",\"windows\":{\"voting_period\":\"3600s\",\"min_execution_period\":\"0s\"}}", Index: false}}}},
		},
		"valid percentage policy": {
			authority: s.authority,
			policy: &foundation.PercentageDecisionPolicy{
				Percentage: math.LegacyOneDec(),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
			valid:  true,
			events: sdk.Events{{Type: "lbm.foundation.v1.EventUpdateDecisionPolicy", Attributes: []abci.EventAttribute{{Key: "decision_policy", Value: "{\"@type\":\"/lbm.foundation.v1.PercentageDecisionPolicy\",\"percentage\":\"1.000000000000000000\",\"windows\":{\"voting_period\":\"3600s\",\"min_execution_period\":\"0s\"}}", Index: false}}}},
		},
		"empty authority": {
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: math.LegacyOneDec(),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
		},
		"empty policy": {
			authority: s.authority,
		},
		"zero threshold": {
			authority: s.authority,
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: math.LegacyZeroDec(),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
		},
		"zero voting period": {
			authority: s.authority,
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: math.LegacyOneDec(),
				Windows:   &foundation.DecisionPolicyWindows{},
			},
		},
		"invalid percentage": {
			authority: s.authority,
			policy: &foundation.PercentageDecisionPolicy{
				Percentage: math.LegacyNewDec(2),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
		},
		"not authorized": {
			authority: s.stranger,
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: math.LegacyOneDec(),
				Windows:   &foundation.DecisionPolicyWindows{},
			},
		},
		"invalid policy": {
			authority: s.authority,
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: math.LegacyOneDec(),
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
				Authority: s.bytesToString(tc.authority),
			}
			if tc.policy != nil {
				err := req.SetDecisionPolicy(tc.policy)
				s.Require().NoError(err)
			}

			res, err := s.msgServer.UpdateDecisionPolicy(ctx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			s.Require().NotNil(res)

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgUpdateMembers() {
	testCases := map[string]struct {
		authority sdk.AccAddress
		members   []foundation.MemberRequest
		valid     bool
		events    sdk.Events
	}{
		"valid request": {
			authority: s.authority,
			members: []foundation.MemberRequest{{
				Address: s.bytesToString(s.members[0]),
			}},
			valid:  true,
			events: sdk.Events{{Type: "lbm.foundation.v1.EventUpdateMembers", Attributes: []abci.EventAttribute{{Key: "member_updates", Value: "[{\"address\":\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq8ql9tg\",\"remove\":false,\"metadata\":\"\"}]", Index: false}}}},
		},
		"empty authority": {
			members: []foundation.MemberRequest{{
				Address: s.bytesToString(s.members[0]),
			}},
		},
		"empty requests": {
			authority: s.authority,
		},
		"invalid requests": {
			authority: s.authority,
			members:   []foundation.MemberRequest{{}},
		},
		"not authorized": {
			authority: s.stranger,
			members: []foundation.MemberRequest{{
				Address: s.bytesToString(s.members[0]),
			}},
		},
		"remove a non-member": {
			authority: s.authority,
			members: []foundation.MemberRequest{{
				Address: s.bytesToString(s.stranger),
				Remove:  true,
			}},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &foundation.MsgUpdateMembers{
				Authority:     s.bytesToString(tc.authority),
				MemberUpdates: tc.members,
			}
			res, err := s.msgServer.UpdateMembers(ctx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			s.Require().NotNil(res)

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgSubmitProposal() {
	members := make([]string, len(s.members))
	for i, member := range s.members {
		members[i] = s.bytesToString(member)
	}

	testCases := map[string]struct {
		malleate  func(ctx sdk.Context)
		proposers []string
		metadata  string
		msgs      []sdk.Msg
		exec      foundation.Exec
		valid     bool
		events    sdk.Events
	}{
		"valid request (submit)": {
			proposers: members,
			msgs:      []sdk.Msg{s.newTestMsg(s.authority)},
			valid:     true,
			events:    sdk.Events{{Type: "lbm.foundation.v1.EventSubmitProposal", Attributes: []abci.EventAttribute{{Key: "proposal", Value: "{\"id\":\"6\",\"metadata\":\"\",\"proposers\":[\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq8ql9tg\",\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgp6ktsk6\",\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgz597xc9\",\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgrfn2n9h\",\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgyg2aryj\",\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cg94ufkeq\",\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgxm0uqhl\",\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cg8xeg42d\",\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgge5mf44\",\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgfyz0ug8\"],\"submit_time\":\"2023-11-07T19:32:00Z\",\"foundation_version\":\"1\",\"status\":\"PROPOSAL_STATUS_SUBMITTED\",\"final_tally_result\":{\"yes_count\":\"0.000000000000000000\",\"abstain_count\":\"0.000000000000000000\",\"no_count\":\"0.000000000000000000\",\"no_with_veto_count\":\"0.000000000000000000\"},\"voting_period_end\":\"2023-11-14T19:32:00Z\",\"executor_result\":\"PROPOSAL_EXECUTOR_RESULT_NOT_RUN\",\"messages\":[{\"@type\":\"/testpb.TestMsg\",\"signers\":[\"link190vt0vxc8c8vj24a7mm3fjsenfu8f5yxxj76cp\"]}]}", Index: false}}}},
		},
		"valid request (submit & execute)": {
			proposers: members,
			msgs:      []sdk.Msg{s.newTestMsg(s.authority)},
			exec:      foundation.Exec_EXEC_TRY,
			valid:     true,
			events:    sdk.Events{{Type: "lbm.foundation.v1.EventSubmitProposal", Attributes: []abci.EventAttribute{{Key: "proposal", Value: "{\"id\":\"6\",\"metadata\":\"\",\"proposers\":[\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq8ql9tg\",\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgp6ktsk6\",\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgz597xc9\",\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgrfn2n9h\",\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgyg2aryj\",\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cg94ufkeq\",\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgxm0uqhl\",\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cg8xeg42d\",\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgge5mf44\",\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgfyz0ug8\"],\"submit_time\":\"2023-11-07T19:32:00Z\",\"foundation_version\":\"1\",\"status\":\"PROPOSAL_STATUS_SUBMITTED\",\"final_tally_result\":{\"yes_count\":\"0.000000000000000000\",\"abstain_count\":\"0.000000000000000000\",\"no_count\":\"0.000000000000000000\",\"no_with_veto_count\":\"0.000000000000000000\"},\"voting_period_end\":\"2023-11-14T19:32:00Z\",\"executor_result\":\"PROPOSAL_EXECUTOR_RESULT_NOT_RUN\",\"messages\":[{\"@type\":\"/testpb.TestMsg\",\"signers\":[\"link190vt0vxc8c8vj24a7mm3fjsenfu8f5yxxj76cp\"]}]}", Index: false}}}, {Type: "lbm.foundation.v1.EventVote", Attributes: []abci.EventAttribute{{Key: "vote", Value: "{\"proposal_id\":\"6\",\"voter\":\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq8ql9tg\",\"option\":\"VOTE_OPTION_YES\",\"metadata\":\"\",\"submit_time\":\"2023-11-07T19:32:00Z\"}", Index: false}}}, {Type: "lbm.foundation.v1.EventVote", Attributes: []abci.EventAttribute{{Key: "vote", Value: "{\"proposal_id\":\"6\",\"voter\":\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgp6ktsk6\",\"option\":\"VOTE_OPTION_YES\",\"metadata\":\"\",\"submit_time\":\"2023-11-07T19:32:00Z\"}", Index: false}}}, {Type: "lbm.foundation.v1.EventVote", Attributes: []abci.EventAttribute{{Key: "vote", Value: "{\"proposal_id\":\"6\",\"voter\":\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgz597xc9\",\"option\":\"VOTE_OPTION_YES\",\"metadata\":\"\",\"submit_time\":\"2023-11-07T19:32:00Z\"}", Index: false}}}, {Type: "lbm.foundation.v1.EventVote", Attributes: []abci.EventAttribute{{Key: "vote", Value: "{\"proposal_id\":\"6\",\"voter\":\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgrfn2n9h\",\"option\":\"VOTE_OPTION_YES\",\"metadata\":\"\",\"submit_time\":\"2023-11-07T19:32:00Z\"}", Index: false}}}, {Type: "lbm.foundation.v1.EventVote", Attributes: []abci.EventAttribute{{Key: "vote", Value: "{\"proposal_id\":\"6\",\"voter\":\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgyg2aryj\",\"option\":\"VOTE_OPTION_YES\",\"metadata\":\"\",\"submit_time\":\"2023-11-07T19:32:00Z\"}", Index: false}}}, {Type: "lbm.foundation.v1.EventVote", Attributes: []abci.EventAttribute{{Key: "vote", Value: "{\"proposal_id\":\"6\",\"voter\":\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cg94ufkeq\",\"option\":\"VOTE_OPTION_YES\",\"metadata\":\"\",\"submit_time\":\"2023-11-07T19:32:00Z\"}", Index: false}}}, {Type: "lbm.foundation.v1.EventVote", Attributes: []abci.EventAttribute{{Key: "vote", Value: "{\"proposal_id\":\"6\",\"voter\":\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgxm0uqhl\",\"option\":\"VOTE_OPTION_YES\",\"metadata\":\"\",\"submit_time\":\"2023-11-07T19:32:00Z\"}", Index: false}}}, {Type: "lbm.foundation.v1.EventVote", Attributes: []abci.EventAttribute{{Key: "vote", Value: "{\"proposal_id\":\"6\",\"voter\":\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cg8xeg42d\",\"option\":\"VOTE_OPTION_YES\",\"metadata\":\"\",\"submit_time\":\"2023-11-07T19:32:00Z\"}", Index: false}}}, {Type: "lbm.foundation.v1.EventVote", Attributes: []abci.EventAttribute{{Key: "vote", Value: "{\"proposal_id\":\"6\",\"voter\":\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgge5mf44\",\"option\":\"VOTE_OPTION_YES\",\"metadata\":\"\",\"submit_time\":\"2023-11-07T19:32:00Z\"}", Index: false}}}, {Type: "lbm.foundation.v1.EventVote", Attributes: []abci.EventAttribute{{Key: "vote", Value: "{\"proposal_id\":\"6\",\"voter\":\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgfyz0ug8\",\"option\":\"VOTE_OPTION_YES\",\"metadata\":\"\",\"submit_time\":\"2023-11-07T19:32:00Z\"}", Index: false}}}, {Type: "lbm.foundation.v1.EventExec", Attributes: []abci.EventAttribute{{Key: "logs", Value: "\"proposal execution failed on proposal 6, because of error no message handler found for \\\"/testpb.TestMsg\\\": unknown request\"", Index: false}, {Key: "proposal_id", Value: "\"6\"", Index: false}, {Key: "result", Value: "\"PROPOSAL_EXECUTOR_RESULT_FAILURE\"", Index: false}}}},
		},
		"valid request (submit & unable to reach quorum)": {
			proposers: []string{members[0]},
			msgs:      []sdk.Msg{s.newTestMsg(s.authority)},
			exec:      foundation.Exec_EXEC_TRY,
			valid:     true,
			events:    sdk.Events{{Type: "lbm.foundation.v1.EventSubmitProposal", Attributes: []abci.EventAttribute{{Key: "proposal", Value: "{\"id\":\"6\",\"metadata\":\"\",\"proposers\":[\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq8ql9tg\"],\"submit_time\":\"2023-11-07T19:32:00Z\",\"foundation_version\":\"1\",\"status\":\"PROPOSAL_STATUS_SUBMITTED\",\"final_tally_result\":{\"yes_count\":\"0.000000000000000000\",\"abstain_count\":\"0.000000000000000000\",\"no_count\":\"0.000000000000000000\",\"no_with_veto_count\":\"0.000000000000000000\"},\"voting_period_end\":\"2023-11-14T19:32:00Z\",\"executor_result\":\"PROPOSAL_EXECUTOR_RESULT_NOT_RUN\",\"messages\":[{\"@type\":\"/testpb.TestMsg\",\"signers\":[\"link190vt0vxc8c8vj24a7mm3fjsenfu8f5yxxj76cp\"]}]}", Index: false}}}, {Type: "lbm.foundation.v1.EventVote", Attributes: []abci.EventAttribute{{Key: "vote", Value: "{\"proposal_id\":\"6\",\"voter\":\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq8ql9tg\",\"option\":\"VOTE_OPTION_YES\",\"metadata\":\"\",\"submit_time\":\"2023-11-07T19:32:00Z\"}", Index: false}}}, {Type: "lbm.foundation.v1.EventExec", Attributes: []abci.EventAttribute{{Key: "logs", Value: "\"proposal execution failed on proposal 6, because of error no message handler found for \\\"/testpb.TestMsg\\\": unknown request\"", Index: false}, {Key: "proposal_id", Value: "\"6\"", Index: false}, {Key: "result", Value: "\"PROPOSAL_EXECUTOR_RESULT_FAILURE\"", Index: false}}}},
		},
		"empty proposers": {
			msgs: []sdk.Msg{s.newTestMsg()},
		},
		"invalid proposer": {
			proposers: []string{},
			msgs:      []sdk.Msg{s.newTestMsg()},
		},
		"duplicate proposers": {
			proposers: []string{members[0], members[0]},
			msgs:      []sdk.Msg{s.newTestMsg()},
		},
		"empty msgs": {
			proposers: []string{members[0]},
		},
		"invalid msg": {
			proposers: []string{members[0]},
			msgs:      []sdk.Msg{&foundation.MsgWithdrawFromTreasury{}},
		},
		"invalid exec": {
			proposers: []string{members[0]},
			msgs:      []sdk.Msg{s.newTestMsg()},
			exec:      -1,
		},
		"not a member": {
			proposers: []string{s.bytesToString(s.stranger)},
			msgs:      []sdk.Msg{s.newTestMsg(s.authority)},
		},
		"unauthorized msg": {
			proposers: []string{members[0]},
			msgs:      []sdk.Msg{s.newTestMsg(s.stranger)},
		},
		"exec fails": {
			malleate: func(ctx sdk.Context) {
				// try exec will fail because of a non-zero MinExecutionPeriod.
				err := s.impl.UpdateDecisionPolicy(ctx, &foundation.ThresholdDecisionPolicy{
					Threshold: math.LegacyOneDec(),
					Windows: &foundation.DecisionPolicyWindows{
						VotingPeriod:       time.Hour,
						MinExecutionPeriod: time.Second,
					},
				})
				s.Require().NoError(err)
			},
			proposers: members,
			msgs:      []sdk.Msg{s.newTestMsg(s.authority)},
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
			err := req.SetMsgs(tc.msgs)
			s.Require().NoError(err)

			res, err := s.msgServer.SubmitProposal(ctx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			s.Require().NotNil(res)

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgWithdrawProposal() {
	testCases := map[string]struct {
		proposalID uint64
		address    sdk.AccAddress
		valid      bool
		events     sdk.Events
	}{
		"valid request (proposer)": {
			proposalID: s.activeProposal,
			address:    s.members[0],
			valid:      true,
			events:     sdk.Events{{Type: "lbm.foundation.v1.EventWithdrawProposal", Attributes: []abci.EventAttribute{{Key: "proposal_id", Value: "\"1\"", Index: false}}}},
		},
		"valid request (authority)": {
			proposalID: s.activeProposal,
			address:    s.authority,
			valid:      true,
			events:     sdk.Events{{Type: "lbm.foundation.v1.EventWithdrawProposal", Attributes: []abci.EventAttribute{{Key: "proposal_id", Value: "\"1\"", Index: false}}}},
		},
		"empty proposal id": {
			address: s.members[0],
		},
		"empty address": {
			proposalID: 1,
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
				Address:    s.bytesToString(tc.address),
			}
			res, err := s.msgServer.WithdrawProposal(ctx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			s.Require().NotNil(res)

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgVote() {
	testCases := map[string]struct {
		malleate   func(ctx sdk.Context)
		proposalID uint64
		voter      sdk.AccAddress
		option     foundation.VoteOption
		exec       foundation.Exec
		valid      bool
		events     sdk.Events
	}{
		"valid request (vote)": {
			proposalID: s.activeProposal,
			voter:      s.members[0],
			option:     foundation.VOTE_OPTION_YES,
			valid:      true,
			events:     sdk.Events{{Type: "lbm.foundation.v1.EventVote", Attributes: []abci.EventAttribute{{Key: "vote", Value: "{\"proposal_id\":\"1\",\"voter\":\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq8ql9tg\",\"option\":\"VOTE_OPTION_YES\",\"metadata\":\"\",\"submit_time\":\"2023-11-07T19:32:00Z\"}", Index: false}}}},
		},
		"valid request (vote & execute)": {
			proposalID: s.activeProposal,
			voter:      s.members[0],
			option:     foundation.VOTE_OPTION_YES,
			exec:       foundation.Exec_EXEC_TRY,
			valid:      true,
			events:     sdk.Events{{Type: "lbm.foundation.v1.EventVote", Attributes: []abci.EventAttribute{{Key: "vote", Value: "{\"proposal_id\":\"1\",\"voter\":\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq8ql9tg\",\"option\":\"VOTE_OPTION_YES\",\"metadata\":\"\",\"submit_time\":\"2023-11-07T19:32:00Z\"}", Index: false}}}, {Type: "lbm.foundation.v1.EventWithdrawFromTreasury", Attributes: []abci.EventAttribute{{Key: "amount", Value: "[{\"denom\":\"stake\",\"amount\":\"987654321\"}]", Index: false}, {Key: "to", Value: "\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgsjpha7m\"", Index: false}}}, {Type: "lbm.foundation.v1.EventExec", Attributes: []abci.EventAttribute{{Key: "logs", Value: "\"\"", Index: false}, {Key: "proposal_id", Value: "\"1\"", Index: false}, {Key: "result", Value: "\"PROPOSAL_EXECUTOR_RESULT_SUCCESS\"", Index: false}}}},
		},
		"empty proposal id": {
			voter:  s.members[0],
			option: foundation.VOTE_OPTION_YES,
		},
		"empty voter": {
			proposalID: 1,
			option:     foundation.VOTE_OPTION_YES,
		},
		"empty option": {
			proposalID: 1,
			voter:      s.members[0],
		},
		"invalid option": {
			proposalID: 1,
			voter:      s.members[0],
			option:     -1,
		},
		"invalid exec": {
			proposalID: 1,
			voter:      s.members[0],
			option:     foundation.VOTE_OPTION_YES,
			exec:       -1,
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
				err := s.impl.UpdateDecisionPolicy(ctx, &foundation.ThresholdDecisionPolicy{
					Threshold: math.LegacyOneDec(),
					Windows: &foundation.DecisionPolicyWindows{
						VotingPeriod:       time.Hour,
						MinExecutionPeriod: time.Second,
					},
				})
				s.Require().NoError(err)

				// submit a proposal
				proposers := make([]string, len(s.members))
				for i, member := range s.members {
					proposers[i] = s.bytesToString(member)
				}
				req := &foundation.MsgSubmitProposal{
					Proposers: proposers,
				}
				err = req.SetMsgs([]sdk.Msg{s.newTestMsg(s.authority)})
				s.Require().NoError(err)

				res, err := s.msgServer.SubmitProposal(ctx, req)
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
				Voter:      s.bytesToString(tc.voter),
				Option:     tc.option,
				Exec:       tc.exec,
			}
			res, err := s.msgServer.Vote(ctx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			s.Require().NotNil(res)

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgExec() {
	testCases := map[string]struct {
		malleate   func(ctx sdk.Context)
		proposalID uint64
		signer     sdk.AccAddress
		valid      bool
		events     sdk.Events
	}{
		"valid request (execute)": {
			proposalID: s.activeProposal,
			signer:     s.members[0],
			valid:      true,
			events:     sdk.Events{{Type: "lbm.foundation.v1.EventWithdrawFromTreasury", Attributes: []abci.EventAttribute{{Key: "amount", Value: "[{\"denom\":\"stake\",\"amount\":\"987654321\"}]", Index: false}, {Key: "to", Value: "\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgsjpha7m\"", Index: false}}}, {Type: "lbm.foundation.v1.EventExec", Attributes: []abci.EventAttribute{{Key: "logs", Value: "\"\"", Index: false}, {Key: "proposal_id", Value: "\"1\"", Index: false}, {Key: "result", Value: "\"PROPOSAL_EXECUTOR_RESULT_SUCCESS\"", Index: false}}}},
		},
		"valid request (not finalized)": {
			proposalID: s.votedProposal,
			signer:     s.members[0],
			valid:      true,
			events:     sdk.Events{{Type: "lbm.foundation.v1.EventExec", Attributes: []abci.EventAttribute{{Key: "logs", Value: "\"\"", Index: false}, {Key: "proposal_id", Value: "\"2\"", Index: false}, {Key: "result", Value: "\"PROPOSAL_EXECUTOR_RESULT_NOT_RUN\"", Index: false}}}},
		},
		"empty proposal id": {
			signer: s.members[0],
		},
		"empty signer": {
			proposalID: 1,
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
				Signer:     s.bytesToString(tc.signer),
			}
			res, err := s.msgServer.Exec(ctx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			s.Require().NotNil(res)

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgLeaveFoundation() {
	testCases := map[string]struct {
		malleate func(ctx sdk.Context)
		address  sdk.AccAddress
		valid    bool
		events   sdk.Events
	}{
		"valid request": {
			address: s.members[0],
			valid:   true,
			events:  sdk.Events{{Type: "lbm.foundation.v1.EventLeaveFoundation", Attributes: []abci.EventAttribute{{Key: "address", Value: "\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq8ql9tg\"", Index: false}}}},
		},
		"empty address": {},
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
						Address: s.bytesToString(member),
						Remove:  true,
					}
				}
				err := s.impl.UpdateMembers(ctx, requests)
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
				Address: s.bytesToString(tc.address),
			}
			res, err := s.msgServer.LeaveFoundation(ctx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			s.Require().NotNil(res)

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgUpdateCensorship() {
	testCases := map[string]struct {
		authority  sdk.AccAddress
		censorship foundation.Censorship
		valid      bool
		events     sdk.Events
	}{
		"valid request": {
			authority: s.authority,
			censorship: foundation.Censorship{
				MsgTypeUrl: sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
				Authority:  foundation.CensorshipAuthorityGovernance,
			},
			valid:  true,
			events: sdk.Events{{Type: "lbm.foundation.v1.EventUpdateCensorship", Attributes: []abci.EventAttribute{{Key: "censorship", Value: "{\"msg_type_url\":\"/lbm.foundation.v1.MsgWithdrawFromTreasury\",\"authority\":\"CENSORSHIP_AUTHORITY_GOVERNANCE\"}", Index: false}}}},
		},
		"invalid authority": {
			authority: s.stranger,
			censorship: foundation.Censorship{
				MsgTypeUrl: sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
				Authority:  foundation.CensorshipAuthorityGovernance,
			},
		},
		"enabling feature": {
			authority: s.authority,
			censorship: foundation.Censorship{
				MsgTypeUrl: sdk.MsgTypeURL((*testdata.TestMsg)(nil)),
				Authority:  foundation.CensorshipAuthorityFoundation,
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &foundation.MsgUpdateCensorship{
				Authority:  s.bytesToString(tc.authority),
				Censorship: tc.censorship,
			}
			res, err := s.msgServer.UpdateCensorship(ctx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			s.Require().NotNil(res)

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgGrant() {
	testCases := map[string]struct {
		authority     sdk.AccAddress
		grantee       sdk.AccAddress
		authorization foundation.Authorization
		valid         bool
		events        sdk.Events
	}{
		"valid request": {
			authority:     s.authority,
			grantee:       s.members[0],
			authorization: &foundation.ReceiveFromTreasuryAuthorization{},
			valid:         true,
			events:        sdk.Events{{Type: "lbm.foundation.v1.EventGrant", Attributes: []abci.EventAttribute{{Key: "authorization", Value: "{\"@type\":\"/lbm.foundation.v1.ReceiveFromTreasuryAuthorization\"}", Index: false}, {Key: "grantee", Value: "\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq8ql9tg\"", Index: false}}}},
		},
		"empty authority": {
			grantee:       s.members[0],
			authorization: &foundation.ReceiveFromTreasuryAuthorization{},
		},
		"empty grantee": {
			authority:     s.authority,
			authorization: &foundation.ReceiveFromTreasuryAuthorization{},
		},
		"empty authorization": {
			authority: s.authority,
			grantee:   s.members[0],
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
				Authority: s.bytesToString(tc.authority),
				Grantee:   s.bytesToString(tc.grantee),
			}
			if tc.authorization != nil {
				err := req.SetAuthorization(tc.authorization)
				s.Require().NoError(err)
			}

			res, err := s.msgServer.Grant(ctx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			s.Require().NotNil(res)

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgRevoke() {
	testCases := map[string]struct {
		authority  sdk.AccAddress
		grantee    sdk.AccAddress
		msgTypeURL string
		valid      bool
		events     sdk.Events
	}{
		"valid request": {
			authority:  s.authority,
			grantee:    s.stranger,
			msgTypeURL: foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
			valid:      true,
			events:     sdk.Events{{Type: "lbm.foundation.v1.EventRevoke", Attributes: []abci.EventAttribute{{Key: "grantee", Value: "\"link15ky9du8a2wlstz6fpx3p4mqpjyrm5cgsjpha7m\"", Index: false}, {Key: "msg_type_url", Value: "\"/lbm.foundation.v1.MsgWithdrawFromTreasury\"", Index: false}}}},
		},
		"empty authority": {
			grantee:    s.stranger,
			msgTypeURL: foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
		},
		"empty grantee": {
			authority:  s.authority,
			msgTypeURL: foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
		},
		"empty url": {
			authority: s.authority,
			grantee:   s.stranger,
		},
		"not authorized": {
			authority:  s.stranger,
			grantee:    s.stranger,
			msgTypeURL: foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
		},
		"no grant": {
			authority:  s.authority,
			grantee:    s.members[0],
			msgTypeURL: foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &foundation.MsgRevoke{
				Authority:  s.bytesToString(tc.authority),
				Grantee:    s.bytesToString(tc.grantee),
				MsgTypeUrl: tc.msgTypeURL,
			}
			res, err := s.msgServer.Revoke(ctx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			s.Require().NotNil(res)

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}
