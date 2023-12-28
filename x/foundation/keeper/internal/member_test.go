package internal_test

import (
	"time"

	"cosmossdk.io/math"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

func (s *KeeperTestSuite) TestUpdateDecisionPolicy() {
	config := foundation.DefaultConfig()
	testCases := map[string]struct {
		policy foundation.DecisionPolicy
		valid  bool
	}{
		"valid policy": {
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: math.LegacyOneDec(),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
			valid: true,
		},
		"invalid policy (invalid min execution period)": {
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: math.LegacyOneDec(),
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

			err := s.impl.UpdateDecisionPolicy(ctx, tc.policy)
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
		updates []foundation.MemberRequest
		valid   bool
	}{
		"add a new member": {
			updates: []foundation.MemberRequest{
				{
					Address: s.bytesToString(s.stranger),
				},
			},
			valid: true,
		},
		"remove a member": {
			updates: []foundation.MemberRequest{
				{
					Address: s.bytesToString(s.members[0]),
					Remove:  true,
				},
			},
			valid: true,
		},
		"remove a non-member": {
			updates: []foundation.MemberRequest{
				{
					Address: s.bytesToString(s.stranger),
					Remove:  true,
				},
			},
		},
		"long metadata": {
			updates: []foundation.MemberRequest{
				{
					Address:  s.bytesToString(s.stranger),
					Metadata: string(make([]rune, 256)),
				},
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			err := s.impl.UpdateMembers(ctx, tc.updates)
			if tc.valid {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}
