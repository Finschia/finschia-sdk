package keeper_test

import (
	"time"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

func (s *KeeperTestSuite) TestUpdateDecisionPolicy() {
	config := foundation.DefaultConfig()
	testCases := map[string]struct {
		policy foundation.DecisionPolicy
		valid  bool
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
					VotingPeriod:       time.Hour,
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
					VotingPeriod:       time.Hour,
					MinExecutionPeriod: time.Hour + config.MaxExecutionPeriod,
				},
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			err := s.keeper.UpdateDecisionPolicy(ctx, tc.policy)
			if tc.valid {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestUpdateMembers() {
	testCases := map[string]struct {
		updates []foundation.Member
		valid   bool
	}{
		"add a new member": {
			updates: []foundation.Member{
				{
					Address:       s.stranger.String(),
					Participating: true,
				},
			},
			valid: true,
		},
		"remove a member": {
			updates: []foundation.Member{
				{
					Address: s.members[0].String(),
				},
			},
			valid: true,
		},
		"remove a non-member": {
			updates: []foundation.Member{
				{
					Address: s.stranger.String(),
				},
			},
		},
		"long metadata": {
			updates: []foundation.Member{
				{
					Address:       s.stranger.String(),
					Participating: true,
					Metadata:      string(make([]rune, 256)),
				},
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			err := s.keeper.UpdateMembers(ctx, tc.updates)
			if tc.valid {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestUpdateOperator() {
	testCases := map[string]struct {
		operator sdk.AccAddress
		valid    bool
	}{
		"valid new operator": {
			operator: s.stranger,
			valid:    true,
		},
		"already the operator": {
			operator: s.operator,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			err := s.keeper.UpdateOperator(ctx, tc.operator)
			if tc.valid {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}
