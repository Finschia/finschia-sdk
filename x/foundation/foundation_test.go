package foundation_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"

	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

func TestTallyResult(t *testing.T) {
	result := foundation.DefaultTallyResult()

	err := result.Add(foundation.VOTE_OPTION_UNSPECIFIED)
	require.Error(t, err)

	err = result.Add(foundation.VOTE_OPTION_YES)
	require.NoError(t, err)
	require.Equal(t, math.LegacyOneDec(), result.YesCount)

	err = result.Add(foundation.VOTE_OPTION_ABSTAIN)
	require.NoError(t, err)
	require.Equal(t, math.LegacyOneDec(), result.AbstainCount)

	err = result.Add(foundation.VOTE_OPTION_NO)
	require.NoError(t, err)
	require.Equal(t, math.LegacyOneDec(), result.NoCount)

	err = result.Add(foundation.VOTE_OPTION_NO_WITH_VETO)
	require.NoError(t, err)
	require.Equal(t, math.LegacyOneDec(), result.NoWithVetoCount)

	require.Equal(t, math.LegacyNewDec(4), result.TotalCounts())
}

func TestThresholdDecisionPolicy(t *testing.T) {
	config := foundation.DefaultConfig()

	testCases := map[string]struct {
		threshold          math.LegacyDec
		votingPeriod       time.Duration
		minExecutionPeriod time.Duration
		totalWeight        math.LegacyDec
		validBasic         bool
		valid              bool
	}{
		"valid policy": {
			threshold:          math.LegacyOneDec(),
			votingPeriod:       time.Hour,
			minExecutionPeriod: config.MaxExecutionPeriod + time.Hour - time.Nanosecond,
			totalWeight:        math.LegacyOneDec(),
			validBasic:         true,
			valid:              true,
		},
		"invalid threshold": {
			votingPeriod:       time.Hour,
			minExecutionPeriod: config.MaxExecutionPeriod + time.Hour - time.Nanosecond,
			totalWeight:        math.LegacyOneDec(),
		},
		"invalid voting period": {
			threshold:          math.LegacyOneDec(),
			minExecutionPeriod: config.MaxExecutionPeriod - time.Nanosecond,
			totalWeight:        math.LegacyOneDec(),
		},
		"invalid min execution period": {
			threshold:          math.LegacyOneDec(),
			votingPeriod:       time.Hour,
			minExecutionPeriod: config.MaxExecutionPeriod + time.Hour,
			totalWeight:        math.LegacyOneDec(),
			validBasic:         true,
		},
		"invalid total weight": {
			threshold:          math.LegacyOneDec(),
			votingPeriod:       time.Hour,
			minExecutionPeriod: config.MaxExecutionPeriod + time.Hour - time.Nanosecond,
			totalWeight:        math.LegacyZeroDec(),
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
		Threshold: math.LegacyNewDec(10),
		Windows: &foundation.DecisionPolicyWindows{
			VotingPeriod: time.Hour,
		},
	}
	require.NoError(t, policy.ValidateBasic())

	info := foundation.FoundationInfo{
		TotalWeight: math.LegacyOneDec(),
	}
	require.NoError(t, policy.Validate(info, config))
	require.Equal(t, time.Hour, policy.GetVotingPeriod())

	testCases := map[string]struct {
		sinceSubmission time.Duration
		totalWeight     math.LegacyDec
		tally           foundation.TallyResult
		valid           bool
		final           bool
		allow           bool
	}{
		"allow": {
			sinceSubmission: policy.Windows.MinExecutionPeriod,
			totalWeight:     policy.Threshold,
			tally:           foundation.NewTallyResult(policy.Threshold, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
			valid:           true,
			final:           true,
			allow:           true,
		},
		"allow (member size < threshold)": {
			sinceSubmission: policy.Windows.MinExecutionPeriod,
			totalWeight:     math.LegacyOneDec(),
			tally:           foundation.NewTallyResult(math.LegacyOneDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
			valid:           true,
			final:           true,
			allow:           true,
		},
		"not final": {
			sinceSubmission: policy.Windows.MinExecutionPeriod,
			totalWeight:     policy.Threshold,
			tally:           foundation.NewTallyResult(policy.Threshold.Sub(math.LegacyOneDec()), math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
			valid:           true,
		},
		"deny": {
			sinceSubmission: policy.Windows.MinExecutionPeriod,
			totalWeight:     policy.Threshold.Add(math.LegacyOneDec()),
			tally:           foundation.NewTallyResult(math.LegacyZeroDec(), math.LegacyOneDec(), math.LegacyOneDec(), math.LegacyZeroDec()),
			valid:           true,
			final:           true,
		},
		"too early": {
			sinceSubmission: policy.Windows.MinExecutionPeriod - time.Nanosecond,
			totalWeight:     policy.Threshold,
			tally:           foundation.NewTallyResult(policy.Threshold, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
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
		percentage         math.LegacyDec
		votingPeriod       time.Duration
		minExecutionPeriod time.Duration
		totalWeight        math.LegacyDec
		validBasic         bool
		valid              bool
	}{
		"valid policy": {
			percentage:         math.LegacyOneDec(),
			votingPeriod:       time.Hour,
			minExecutionPeriod: config.MaxExecutionPeriod + time.Hour - time.Nanosecond,
			totalWeight:        math.LegacyOneDec(),
			validBasic:         true,
			valid:              true,
		},
		"invalid percentage": {
			votingPeriod:       time.Hour,
			minExecutionPeriod: config.MaxExecutionPeriod + time.Hour - time.Nanosecond,
			totalWeight:        math.LegacyOneDec(),
		},
		"invalid voting period": {
			percentage:         math.LegacyOneDec(),
			minExecutionPeriod: config.MaxExecutionPeriod - time.Nanosecond,
			totalWeight:        math.LegacyOneDec(),
		},
		"invalid min execution period": {
			percentage:         math.LegacyOneDec(),
			votingPeriod:       time.Hour,
			minExecutionPeriod: config.MaxExecutionPeriod + time.Hour,
			totalWeight:        math.LegacyOneDec(),
			validBasic:         true,
		},
		"invalid total weight": {
			percentage:         math.LegacyOneDec(),
			votingPeriod:       time.Hour,
			minExecutionPeriod: config.MaxExecutionPeriod + time.Hour - time.Nanosecond,
			totalWeight:        math.LegacyZeroDec(),
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
		Percentage: math.LegacyMustNewDecFromStr("0.8"),
		Windows: &foundation.DecisionPolicyWindows{
			VotingPeriod: time.Hour,
		},
	}
	require.NoError(t, policy.ValidateBasic())

	info := foundation.FoundationInfo{
		TotalWeight: math.LegacyOneDec(),
	}
	require.NoError(t, policy.Validate(info, config))
	require.Equal(t, time.Hour, policy.GetVotingPeriod())

	totalWeight := math.LegacyNewDec(10)
	testCases := map[string]struct {
		sinceSubmission time.Duration
		tally           foundation.TallyResult
		valid           bool
		final           bool
		allow           bool
	}{
		"allow": {
			sinceSubmission: policy.Windows.MinExecutionPeriod,
			tally:           foundation.NewTallyResult(math.LegacyNewDec(8), math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
			valid:           true,
			final:           true,
			allow:           true,
		},
		"allow (abstain)": {
			sinceSubmission: policy.Windows.MinExecutionPeriod,
			tally:           foundation.NewTallyResult(math.LegacyNewDec(4), math.LegacyNewDec(5), math.LegacyZeroDec(), math.LegacyZeroDec()),
			valid:           true,
			final:           true,
			allow:           true,
		},
		"not final": {
			sinceSubmission: policy.Windows.MinExecutionPeriod,
			tally:           foundation.NewTallyResult(math.LegacyZeroDec(), math.LegacyNewDec(5), math.LegacyNewDec(1), math.LegacyZeroDec()),
			valid:           true,
		},
		"deny": {
			sinceSubmission: policy.Windows.MinExecutionPeriod,
			tally:           foundation.NewTallyResult(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyNewDec(3), math.LegacyZeroDec()),
			valid:           true,
			final:           true,
		},
		"deny (all abstain)": {
			sinceSubmission: policy.Windows.MinExecutionPeriod,
			tally:           foundation.NewTallyResult(math.LegacyZeroDec(), math.LegacyNewDec(10), math.LegacyZeroDec(), math.LegacyZeroDec()),
			valid:           true,
			final:           true,
		},
		"too early": {
			sinceSubmission: policy.Windows.MinExecutionPeriod - time.Nanosecond,
			tally:           foundation.NewTallyResult(math.LegacyNewDec(8), math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
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
	addressCodec := addresscodec.NewBech32Codec("link")
	bytesToString := func(addr sdk.AccAddress) string {
		str, err := addressCodec.BytesToString(addr)
		require.NoError(t, err)
		return str
	}

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
					Address: bytesToString(addrs[0]),
				},
				{
					Address: bytesToString(addrs[1]),
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
					Address: bytesToString(addrs[0]),
				},
				{
					Address: bytesToString(addrs[0]),
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			members := foundation.Members{tc.members}
			err := members.ValidateBasic(addressCodec)
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMemberRequests(t *testing.T) {
	addressCodec := addresscodec.NewBech32Codec("link")
	bytesToString := func(addr sdk.AccAddress) string {
		str, err := addressCodec.BytesToString(addr)
		require.NoError(t, err)
		return str
	}

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
					Address: bytesToString(addrs[0]),
				},
				{
					Address: bytesToString(addrs[1]),
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
					Address: bytesToString(addrs[0]),
				},
				{
					Address: bytesToString(addrs[0]),
					Remove:  true,
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			requests := foundation.MemberRequests{tc.members}
			err := requests.ValidateBasic(addressCodec)
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestProposal(t *testing.T) {
	addressCodec := addresscodec.NewBech32Codec("link")
	bytesToString := func(addr sdk.AccAddress) string {
		str, err := addressCodec.BytesToString(addr)
		require.NoError(t, err)
		return str
	}

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
				bytesToString(addrs[0]),
				bytesToString(addrs[1]),
			},
			version: 1,
			msgs: []sdk.Msg{
				testdata.NewTestMsg(),
			},
			valid: true,
		},
		"invalid id": {
			proposers: []string{
				bytesToString(addrs[0]),
				bytesToString(addrs[1]),
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
				bytesToString(addrs[0]),
				bytesToString(addrs[0]),
			},
			version: 1,
			msgs: []sdk.Msg{
				testdata.NewTestMsg(),
			},
		},
		"invalid version": {
			id: 1,
			proposers: []string{
				bytesToString(addrs[0]),
				bytesToString(addrs[1]),
			},
			msgs: []sdk.Msg{
				testdata.NewTestMsg(),
			},
		},
		"empty msgs": {
			id: 1,
			proposers: []string{
				bytesToString(addrs[0]),
				bytesToString(addrs[1]),
			},
			version: 1,
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

			err := proposal.ValidateBasic(addressCodec)
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
		totalWeight math.LegacyDec
		validBasic  bool
		valid       bool
	}{
		"invalid policy": {
			totalWeight: math.LegacyOneDec(),
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
		TotalWeight: math.LegacyOneDec(),
	}
	require.Error(t, policy.Validate(info, config))
	require.Zero(t, policy.GetVotingPeriod())

	testCases := map[string]struct {
		sinceSubmission time.Duration
		totalWeight     math.LegacyDec
		tally           foundation.TallyResult
		valid           bool
		final           bool
		allow           bool
	}{
		"deny": {
			sinceSubmission: 0,
			totalWeight:     math.LegacyOneDec(),
			tally:           foundation.NewTallyResult(math.LegacyOneDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
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
