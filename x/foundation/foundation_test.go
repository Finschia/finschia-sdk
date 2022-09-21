package foundation_test

import (
	"testing"
	"time"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
	"github.com/stretchr/testify/require"
)

func TestDecisionPolicy(t *testing.T) {
	config := foundation.DefaultConfig()
	policy := foundation.DefaultDecisionPolicy()

	require.NoError(t, policy.ValidateBasic())
	require.NoError(t, policy.Validate(config))
}

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
		validBasic         bool
		valid              bool
	}{
		"valid policy": {
			threshold:          sdk.OneDec(),
			votingPeriod:       time.Hour,
			minExecutionPeriod: config.MaxExecutionPeriod + time.Hour - time.Nanosecond,
			validBasic:         true,
			valid:              true,
		},
		"invalid threshold": {
			votingPeriod:       time.Hour,
			minExecutionPeriod: config.MaxExecutionPeriod + time.Hour - time.Nanosecond,
		},
		"invalid voting period": {
			threshold:          sdk.OneDec(),
			minExecutionPeriod: config.MaxExecutionPeriod - time.Nanosecond,
		},
		"invalid min execution period": {
			threshold:          sdk.OneDec(),
			votingPeriod:       time.Hour,
			minExecutionPeriod: config.MaxExecutionPeriod + time.Hour,
			validBasic:         true,
		},
	}

	for name, tc := range testCases {
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
			require.Error(t, err, name)
			continue
		}
		require.NoError(t, err, name)

		err = policy.Validate(config)
		if !tc.valid {
			require.Error(t, err, name)
			continue
		}
		require.NoError(t, err, name)
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
	require.NoError(t, policy.Validate(config))
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
		result, err := policy.Allow(tc.tally, tc.totalWeight, tc.sinceSubmission)
		if !tc.valid {
			require.Error(t, err, name)
			continue
		}
		require.NoError(t, err, name)

		require.Equal(t, tc.final, result.Final, name)
		if tc.final {
			require.Equal(t, tc.allow, result.Allow, name)
		}
	}
}

func TestPercentageDecisionPolicy(t *testing.T) {
	config := foundation.DefaultConfig()

	testCases := map[string]struct {
		percentage         sdk.Dec
		votingPeriod       time.Duration
		minExecutionPeriod time.Duration
		validBasic         bool
		valid              bool
	}{
		"valid policy": {
			percentage:         sdk.OneDec(),
			votingPeriod:       time.Hour,
			minExecutionPeriod: config.MaxExecutionPeriod + time.Hour - time.Nanosecond,
			validBasic:         true,
			valid:              true,
		},
		"invalid percentage": {
			votingPeriod:       time.Hour,
			minExecutionPeriod: config.MaxExecutionPeriod + time.Hour - time.Nanosecond,
		},
		"invalid voting period": {
			percentage:         sdk.OneDec(),
			minExecutionPeriod: config.MaxExecutionPeriod - time.Nanosecond,
		},
		"invalid min execution period": {
			percentage:         sdk.OneDec(),
			votingPeriod:       time.Hour,
			minExecutionPeriod: config.MaxExecutionPeriod + time.Hour,
			validBasic:         true,
		},
	}

	for name, tc := range testCases {
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
			require.Error(t, err, name)
			continue
		}
		require.NoError(t, err, name)

		err = policy.Validate(config)
		if !tc.valid {
			require.Error(t, err, name)
			continue
		}
		require.NoError(t, err, name)
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
	require.NoError(t, policy.Validate(config))
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
		result, err := policy.Allow(tc.tally, totalWeight, tc.sinceSubmission)
		if !tc.valid {
			require.Error(t, err, name)
			continue
		}
		require.NoError(t, err, name)

		require.Equal(t, tc.final, result.Final, name)
		if tc.final {
			require.Equal(t, tc.allow, result.Allow, name)
		}
	}
}
