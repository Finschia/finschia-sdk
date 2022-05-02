package keeper_test

import (
	"time"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

func (s *KeeperTestSuite) TestUpdateDecisionPolicy() {
	config := foundation.DefaultConfig()
	testCases := map[string]struct{
		policy foundation.DecisionPolicy
		valid bool
	}{
		"threshold policy (valid)": {
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: config.MinThreshold,
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
			valid: true,
		},
		"threshold policy (low threshold)": {
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: sdk.OneDec(),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
		},
		"threshold policy (invalid min execution period)": {
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: config.MinThreshold,
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
					MinExecutionPeriod: time.Hour + config.MaxExecutionPeriod,
				},
			},
		},
		"percentage policy (valid)": {
			policy: &foundation.PercentageDecisionPolicy{
				Percentage: config.MinPercentage,
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
			valid: true,
		},
		"percentage policy (low percentage)": {
			policy: &foundation.PercentageDecisionPolicy{
				Percentage: sdk.MustNewDecFromStr("0.1"),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
		},
		"percentage policy (invalid min execution period)": {
			policy: &foundation.PercentageDecisionPolicy{
				Percentage: config.MinPercentage,
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
					MinExecutionPeriod: time.Hour + config.MaxExecutionPeriod,
				},
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			err := s.keeper.UpdateDecisionPolicy(s.ctx, tc.policy)
			if tc.valid {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestUpdateMembers() {
	testCases := map[string]struct{
		updates []foundation.Member
		valid bool
	}{
		"add a new member": {
			updates: []foundation.Member{
				{
					Address: s.comingMember.String(),
					Weight: sdk.OneDec(),
				},
			},
			valid: true,
		},
		"remove a member": {
			updates: []foundation.Member{
				{
					Address: s.badMember.String(),
					Weight: sdk.ZeroDec(),
				},
			},
			valid: true,
		},
		"remove a non-member": {
			updates: []foundation.Member{
				{
					Address: s.stranger.String(),
					Weight: sdk.ZeroDec(),
				},
			},
		},
		"long metadata": {
			updates: []foundation.Member{
				{
					Address: s.comingMember.String(),
					Weight: sdk.OneDec(),
					Metadata: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
				},
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			err := s.keeper.UpdateMembers(s.ctx, tc.updates)
			if tc.valid {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestUpdateOperator() {
	testCases := []struct{
		name string
		operator sdk.AccAddress
		valid bool
	}{
		{
			name: "already the operator",
			operator: s.operator,
		},
		{
			name: "valid operator",
			operator: s.stranger,
			valid: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := s.keeper.UpdateOperator(s.ctx, tc.operator)
			if tc.valid {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}
