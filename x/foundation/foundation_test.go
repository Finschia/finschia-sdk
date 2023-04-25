package foundation_test

import (
	"testing"
	"time"

	"github.com/Finschia/finschia-sdk/crypto/keys/secp256k1"
	"github.com/Finschia/finschia-sdk/testutil/testdata"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/foundation"
	"github.com/stretchr/testify/require"
)

func TestTallyResult(t *testing.T) {
	result := foundation.DefaultTallyResult()

	err := result.Add(foundation.VOTE_OPTION_UNSPECIFIED)
	require.Error(t, err)

	err = result.Add(foundation.VOTE_OPTION_YES)
	require.NoError(t, err)
	require.Equal(t, sdk.OneDec(), result.YesCount)

	result.Add(foundation.VOTE_OPTION_ABSTAIN)
	require.NoError(t, err)
	require.Equal(t, sdk.OneDec(), result.AbstainCount)

	result.Add(foundation.VOTE_OPTION_NO)
	require.NoError(t, err)
	require.Equal(t, sdk.OneDec(), result.NoCount)

	result.Add(foundation.VOTE_OPTION_NO_WITH_VETO)
	require.NoError(t, err)
	require.Equal(t, sdk.OneDec(), result.NoWithVetoCount)

	require.Equal(t, sdk.NewDec(4), result.TotalCounts())
}

func TestThresholdDecisionPolicy(t *testing.T) {
	config := foundation.DefaultConfig()

	testCases := map[string]struct {
		threshold          sdk.Dec
		votingPeriod       time.Duration
		minExecutionPeriod time.Duration
		totalWeight        sdk.Dec
		validBasic         bool
		valid              bool
	}{
		"valid policy": {
			threshold:          sdk.OneDec(),
			votingPeriod:       time.Hour,
			minExecutionPeriod: config.MaxExecutionPeriod + time.Hour - time.Nanosecond,
			totalWeight:        sdk.OneDec(),
			validBasic:         true,
			valid:              true,
		},
		"invalid threshold": {
			votingPeriod:       time.Hour,
			minExecutionPeriod: config.MaxExecutionPeriod + time.Hour - time.Nanosecond,
			totalWeight:        sdk.OneDec(),
		},
		"invalid voting period": {
			threshold:          sdk.OneDec(),
			minExecutionPeriod: config.MaxExecutionPeriod - time.Nanosecond,
			totalWeight:        sdk.OneDec(),
		},
		"invalid min execution period": {
			threshold:          sdk.OneDec(),
			votingPeriod:       time.Hour,
			minExecutionPeriod: config.MaxExecutionPeriod + time.Hour,
			totalWeight:        sdk.OneDec(),
			validBasic:         true,
		},
		"invalid total weight": {
			threshold:          sdk.OneDec(),
			votingPeriod:       time.Hour,
			minExecutionPeriod: config.MaxExecutionPeriod + time.Hour - time.Nanosecond,
			totalWeight:        sdk.ZeroDec(),
			validBasic:         true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			policy := foundation.ThresholdDecisionPolicy{
				Threshold: tc.threshold,
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod:       tc.votingPeriod,
					MinExecutionPeriod: tc.minExecutionPeriod,
				},
			}
			require.Equal(t, tc.votingPeriod, policy.GetVotingPeriod())

			err := policy.ValidateBasic()
			if !tc.validBasic {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			info := foundation.FoundationInfo{
				TotalWeight: tc.totalWeight,
			}
			err = policy.Validate(info, config)
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestThresholdDecisionPolicyAllow(t *testing.T) {
	config := foundation.DefaultConfig()
	policy := foundation.ThresholdDecisionPolicy{
		Threshold: sdk.NewDec(10),
		Windows: &foundation.DecisionPolicyWindows{
			VotingPeriod: time.Hour,
		},
	}
	require.NoError(t, policy.ValidateBasic())

	info := foundation.FoundationInfo{
		TotalWeight: sdk.OneDec(),
	}
	require.NoError(t, policy.Validate(info, config))
	require.Equal(t, time.Hour, policy.GetVotingPeriod())

	testCases := map[string]struct {
		sinceSubmission time.Duration
		totalWeight     sdk.Dec
		tally           foundation.TallyResult
		valid           bool
		final           bool
		allow           bool
	}{
		"allow": {
			sinceSubmission: policy.Windows.MinExecutionPeriod,
			totalWeight:     policy.Threshold,
			tally:           foundation.NewTallyResult(policy.Threshold, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			valid:           true,
			final:           true,
			allow:           true,
		},
		"allow (member size < threshold)": {
			sinceSubmission: policy.Windows.MinExecutionPeriod,
			totalWeight:     sdk.OneDec(),
			tally:           foundation.NewTallyResult(sdk.OneDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			valid:           true,
			final:           true,
			allow:           true,
		},
		"not final": {
			sinceSubmission: policy.Windows.MinExecutionPeriod,
			totalWeight:     policy.Threshold,
			tally:           foundation.NewTallyResult(policy.Threshold.Sub(sdk.OneDec()), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			valid:           true,
		},
		"deny": {
			sinceSubmission: policy.Windows.MinExecutionPeriod,
			totalWeight:     policy.Threshold.Add(sdk.OneDec()),
			tally:           foundation.NewTallyResult(sdk.ZeroDec(), sdk.OneDec(), sdk.OneDec(), sdk.ZeroDec()),
			valid:           true,
			final:           true,
		},
		"too early": {
			sinceSubmission: policy.Windows.MinExecutionPeriod - time.Nanosecond,
			totalWeight:     policy.Threshold,
			tally:           foundation.NewTallyResult(policy.Threshold, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result, err := policy.Allow(tc.tally, tc.totalWeight, tc.sinceSubmission)
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, tc.final, result.Final)
			if tc.final {
				require.Equal(t, tc.allow, result.Allow)
			}
		})
	}
}

func TestPercentageDecisionPolicy(t *testing.T) {
	config := foundation.DefaultConfig()

	testCases := map[string]struct {
		percentage         sdk.Dec
		votingPeriod       time.Duration
		minExecutionPeriod time.Duration
		totalWeight        sdk.Dec
		validBasic         bool
		valid              bool
	}{
		"valid policy": {
			percentage:         sdk.OneDec(),
			votingPeriod:       time.Hour,
			minExecutionPeriod: config.MaxExecutionPeriod + time.Hour - time.Nanosecond,
			totalWeight:        sdk.OneDec(),
			validBasic:         true,
			valid:              true,
		},
		"invalid percentage": {
			votingPeriod:       time.Hour,
			minExecutionPeriod: config.MaxExecutionPeriod + time.Hour - time.Nanosecond,
			totalWeight:        sdk.OneDec(),
		},
		"invalid voting period": {
			percentage:         sdk.OneDec(),
			minExecutionPeriod: config.MaxExecutionPeriod - time.Nanosecond,
			totalWeight:        sdk.OneDec(),
		},
		"invalid min execution period": {
			percentage:         sdk.OneDec(),
			votingPeriod:       time.Hour,
			minExecutionPeriod: config.MaxExecutionPeriod + time.Hour,
			totalWeight:        sdk.OneDec(),
			validBasic:         true,
		},
		"invalid total weight": {
			percentage:         sdk.OneDec(),
			votingPeriod:       time.Hour,
			minExecutionPeriod: config.MaxExecutionPeriod + time.Hour - time.Nanosecond,
			totalWeight:        sdk.ZeroDec(),
			validBasic:         true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			policy := foundation.PercentageDecisionPolicy{
				Percentage: tc.percentage,
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod:       tc.votingPeriod,
					MinExecutionPeriod: tc.minExecutionPeriod,
				},
			}
			require.Equal(t, tc.votingPeriod, policy.GetVotingPeriod())

			err := policy.ValidateBasic()
			if !tc.validBasic {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			info := foundation.FoundationInfo{
				TotalWeight: tc.totalWeight,
			}
			err = policy.Validate(info, config)
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestPercentageDecisionPolicyAllow(t *testing.T) {
	config := foundation.DefaultConfig()
	policy := foundation.PercentageDecisionPolicy{
		Percentage: sdk.MustNewDecFromStr("0.8"),
		Windows: &foundation.DecisionPolicyWindows{
			VotingPeriod: time.Hour,
		},
	}
	require.NoError(t, policy.ValidateBasic())

	info := foundation.FoundationInfo{
		TotalWeight: sdk.OneDec(),
	}
	require.NoError(t, policy.Validate(info, config))
	require.Equal(t, time.Hour, policy.GetVotingPeriod())

	totalWeight := sdk.NewDec(10)
	testCases := map[string]struct {
		sinceSubmission time.Duration
		tally           foundation.TallyResult
		valid           bool
		final           bool
		allow           bool
	}{
		"allow": {
			sinceSubmission: policy.Windows.MinExecutionPeriod,
			tally:           foundation.NewTallyResult(sdk.NewDec(8), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			valid:           true,
			final:           true,
			allow:           true,
		},
		"allow (abstain)": {
			sinceSubmission: policy.Windows.MinExecutionPeriod,
			tally:           foundation.NewTallyResult(sdk.NewDec(4), sdk.NewDec(5), sdk.ZeroDec(), sdk.ZeroDec()),
			valid:           true,
			final:           true,
			allow:           true,
		},
		"not final": {
			sinceSubmission: policy.Windows.MinExecutionPeriod,
			tally:           foundation.NewTallyResult(sdk.ZeroDec(), sdk.NewDec(5), sdk.NewDec(1), sdk.ZeroDec()),
			valid:           true,
		},
		"deny": {
			sinceSubmission: policy.Windows.MinExecutionPeriod,
			tally:           foundation.NewTallyResult(sdk.ZeroDec(), sdk.ZeroDec(), sdk.NewDec(3), sdk.ZeroDec()),
			valid:           true,
			final:           true,
		},
		"deny (all abstain)": {
			sinceSubmission: policy.Windows.MinExecutionPeriod,
			tally:           foundation.NewTallyResult(sdk.ZeroDec(), sdk.NewDec(10), sdk.ZeroDec(), sdk.ZeroDec()),
			valid:           true,
			final:           true,
		},
		"too early": {
			sinceSubmission: policy.Windows.MinExecutionPeriod - time.Nanosecond,
			tally:           foundation.NewTallyResult(sdk.NewDec(8), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result, err := policy.Allow(tc.tally, totalWeight, tc.sinceSubmission)
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, tc.final, result.Final)
			if tc.final {
				require.Equal(t, tc.allow, result.Allow)
			}
		})
	}
}

func TestMembers(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		members []foundation.Member
		valid   bool
	}{
		"valid updates": {
			members: []foundation.Member{
				{
					Address: addrs[0].String(),
				},
				{
					Address: addrs[1].String(),
				},
			},
			valid: true,
		},
		"invalid member": {
			members: []foundation.Member{{}},
		},
		"duplicate members": {
			members: []foundation.Member{
				{
					Address: addrs[0].String(),
				},
				{
					Address: addrs[0].String(),
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			members := foundation.Members{tc.members}
			err := members.ValidateBasic()
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMemberRequests(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		members []foundation.MemberRequest
		valid   bool
	}{
		"valid requests": {
			members: []foundation.MemberRequest{
				{
					Address: addrs[0].String(),
				},
				{
					Address: addrs[1].String(),
					Remove:  true,
				},
			},
			valid: true,
		},
		"invalid member": {
			members: []foundation.MemberRequest{{}},
		},
		"duplicate requests": {
			members: []foundation.MemberRequest{
				{
					Address: addrs[0].String(),
				},
				{
					Address: addrs[0].String(),
					Remove:  true,
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			requests := foundation.MemberRequests{tc.members}
			err := requests.ValidateBasic()
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestProposal(t *testing.T) {
	addrs := make([]sdk.AccAddress, 4)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		id        uint64
		proposers []string
		version   uint64
		msgs      []sdk.Msg
		valid     bool
	}{
		"valid proposal": {
			id: 1,
			proposers: []string{
				addrs[0].String(),
				addrs[1].String(),
			},
			version: 1,
			msgs: []sdk.Msg{
				testdata.NewTestMsg(),
			},
			valid: true,
		},
		"invalid id": {
			proposers: []string{
				addrs[0].String(),
				addrs[1].String(),
			},
			version: 1,
			msgs: []sdk.Msg{
				testdata.NewTestMsg(),
			},
		},
		"empty proposers": {
			id:      1,
			version: 1,
			msgs: []sdk.Msg{
				testdata.NewTestMsg(),
			},
		},
		"invalid proposer": {
			id:        1,
			proposers: []string{""},
			version:   1,
			msgs: []sdk.Msg{
				testdata.NewTestMsg(),
			},
		},
		"duplicate proposers": {
			id: 1,
			proposers: []string{
				addrs[0].String(),
				addrs[0].String(),
			},
			version: 1,
			msgs: []sdk.Msg{
				testdata.NewTestMsg(),
			},
		},
		"invalid version": {
			id: 1,
			proposers: []string{
				addrs[0].String(),
				addrs[1].String(),
			},
			msgs: []sdk.Msg{
				testdata.NewTestMsg(),
			},
		},
		"empty msgs": {
			id: 1,
			proposers: []string{
				addrs[0].String(),
				addrs[1].String(),
			},
			version: 1,
		},
		"invalid msg": {
			id: 1,
			proposers: []string{
				addrs[0].String(),
				addrs[1].String(),
			},
			version: 1,
			msgs: []sdk.Msg{
				&foundation.MsgWithdrawFromTreasury{},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			proposal := foundation.Proposal{
				Id:                tc.id,
				Proposers:         tc.proposers,
				FoundationVersion: tc.version,
			}.WithMsgs(tc.msgs)
			require.NotNil(t, proposal)

			err := proposal.ValidateBasic()
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestOutsourcingDecisionPolicy(t *testing.T) {
	config := foundation.DefaultConfig()

	testCases := map[string]struct {
		totalWeight sdk.Dec
		validBasic  bool
		valid       bool
	}{
		"invalid policy": {
			totalWeight: sdk.OneDec(),
			validBasic:  true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			policy := foundation.OutsourcingDecisionPolicy{}
			require.Zero(t, policy.GetVotingPeriod())

			err := policy.ValidateBasic()
			if !tc.validBasic {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			info := foundation.FoundationInfo{
				TotalWeight: tc.totalWeight,
			}
			err = policy.Validate(info, config)
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestOutsourcingDecisionPolicyAllow(t *testing.T) {
	config := foundation.DefaultConfig()
	policy := foundation.OutsourcingDecisionPolicy{}
	require.NoError(t, policy.ValidateBasic())

	info := foundation.FoundationInfo{
		TotalWeight: sdk.OneDec(),
	}
	require.Error(t, policy.Validate(info, config))
	require.Zero(t, policy.GetVotingPeriod())

	testCases := map[string]struct {
		sinceSubmission time.Duration
		totalWeight     sdk.Dec
		tally           foundation.TallyResult
		valid           bool
		final           bool
		allow           bool
	}{
		"deny": {
			sinceSubmission: 0,
			totalWeight:     sdk.OneDec(),
			tally:           foundation.NewTallyResult(sdk.OneDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result, err := policy.Allow(tc.tally, tc.totalWeight, tc.sinceSubmission)
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, tc.final, result.Final)
			if tc.final {
				require.Equal(t, tc.allow, result.Allow)
			}
		})
	}
}
