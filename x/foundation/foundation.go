package foundation

import (
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"

	"github.com/line/lbm-sdk/codec"
	codectypes "github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

func validateMembers(members []Member) error {
	addrs := map[string]bool{}
	for _, member := range members {
		if err := member.ValidateBasic(); err != nil {
			return err
		}
		if addrs[member.Address] {
			return sdkerrors.ErrInvalidRequest.Wrapf("duplicated address: %s", member.Address)
		}
		addrs[member.Address] = true
	}

	return nil
}

func (m Member) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(m.Address); err != nil {
		return err
	}

	if !m.Weight.Equal(sdk.OneDec()) && !m.Weight.IsZero() {
		return sdkerrors.ErrInvalidRequest.Wrapf("expected a zero or one, got %s", m.Weight)
	}

	return nil
}

type DecisionPolicyResult struct {
	Allow bool
	Final bool
}

type DecisionPolicy interface {
	codec.ProtoMarshaler

	// GetVotingPeriod returns the duration after proposal submission where
	// votes are accepted.
	GetVotingPeriod() time.Duration
	// Allow defines policy-specific logic to allow a proposal to pass or not,
	// based on its tally result, the foundation's total power and the time since
	// the proposal was submitted.
	Allow(tallyResult TallyResult, totalPower sdk.Dec, sinceSubmission time.Duration) (*DecisionPolicyResult, error)

	ValidateBasic() error
	Validate(config Config) error
}

func (t *TallyResult) Add(option VoteOption, weight sdk.Dec) error {
	switch option {
	case VOTE_OPTION_YES:
		t.YesCount = t.YesCount.Add(weight)
	case VOTE_OPTION_NO:
		t.NoCount = t.NoCount.Add(weight)
	case VOTE_OPTION_ABSTAIN:
		t.AbstainCount = t.AbstainCount.Add(weight)
	case VOTE_OPTION_NO_WITH_VETO:
		t.NoWithVetoCount = t.NoWithVetoCount.Add(weight)
	default:
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "unknown vote option %s", option.String())
	}

	return nil
}

// TotalCounts is the sum of all weights.
func (t TallyResult) TotalCounts() sdk.Dec {
	totalCounts := sdk.ZeroDec()

	totalCounts = totalCounts.Add(t.YesCount)
	totalCounts = totalCounts.Add(t.NoCount)
	totalCounts = totalCounts.Add(t.AbstainCount)
	totalCounts = totalCounts.Add(t.NoWithVetoCount)

	return totalCounts
}

var _ codectypes.UnpackInterfacesMessage = (*Proposal)(nil)

func (p *Proposal) GetMsgs() []sdk.Msg {
	msgs, err := GetMsgs(p.Messages, "proposal")
	if err != nil {
		panic(err)
	}
	return msgs
}

func (p *Proposal) SetMsgs(msgs []sdk.Msg) error {
	anys, err := SetMsgs(msgs)
	if err != nil {
		return err
	}
	p.Messages = anys
	return nil
}

func (p Proposal) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return UnpackInterfaces(unpacker, p.Messages)
}

// UnpackInterfaces unpacks Any's to sdk.Msg's.
func UnpackInterfaces(unpacker codectypes.AnyUnpacker, anys []*codectypes.Any) error {
	for _, any := range anys {
		var msg sdk.Msg
		err := unpacker.UnpackAny(any, &msg)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetMsgs takes a slice of Any's and turn them into sdk.Msg's.
func GetMsgs(anys []*codectypes.Any, name string) ([]sdk.Msg, error) {
	msgs := make([]sdk.Msg, len(anys))
	for i, any := range anys {
		cached := any.GetCachedValue()
		if cached == nil {
			return nil, fmt.Errorf("any cached value is nil, %s messages must be correctly packed any values", name)
		}
		msgs[i] = cached.(sdk.Msg)
	}
	return msgs, nil
}

// SetMsgs takes a slice of sdk.Msg's and turn them into Any's.
func SetMsgs(msgs []sdk.Msg) ([]*codectypes.Any, error) {
	anys := make([]*codectypes.Any, len(msgs))
	for i, msg := range msgs {
		var err error
		anys[i], err = codectypes.NewAnyWithValue(msg)
		if err != nil {
			return nil, err
		}
	}
	return anys, nil
}

var _ codectypes.UnpackInterfacesMessage = (*FoundationInfo)(nil)

func (i FoundationInfo) GetDecisionPolicy() DecisionPolicy {
	policy, ok := i.DecisionPolicy.GetCachedValue().(DecisionPolicy)
	if !ok {
		return nil
	}
	return policy
}

func (i *FoundationInfo) SetDecisionPolicy(policy DecisionPolicy) error {
	msg, ok := policy.(proto.Message)
	if !ok {
		return sdkerrors.ErrInvalidType.Wrapf("can't proto marshal %T", msg)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return err
	}
	i.DecisionPolicy = any

	return nil
}

func (i *FoundationInfo) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var policy DecisionPolicy
	return unpacker.UnpackAny(i.DecisionPolicy, &policy)
}

var _ DecisionPolicy = (*ThresholdDecisionPolicy)(nil)

func (p ThresholdDecisionPolicy) Validate(config Config) error {
	if p.Threshold.LT(config.MinThreshold) {
		return sdkerrors.ErrInvalidRequest.Wrap("threshold must be greater than or equal to min_threshold")
	}

	if p.Windows.MinExecutionPeriod > p.Windows.VotingPeriod+config.MaxExecutionPeriod {
		return sdkerrors.ErrInvalidRequest.Wrap("min_execution_period should be smaller than voting_period + max_execution_period")
	}

	return nil
}

func (p ThresholdDecisionPolicy) Allow(result TallyResult, totalWeight sdk.Dec, sinceSubmission time.Duration) (*DecisionPolicyResult, error) {
	if sinceSubmission < p.Windows.MinExecutionPeriod {
		return nil, sdkerrors.ErrUnauthorized.Wrapf("must wait %s after submission before execution, currently at %s", p.Windows.MinExecutionPeriod, sinceSubmission)
	}

	// the real threshold of the policy is `min(threshold,total_weight)`. If
	// the foundation member weights changes (member leaving, member weight update)
	// and the threshold doesn't, we can end up with threshold > total_weight.
	// In this case, as long as everyone votes yes (in which case
	// `yesCount`==`realThreshold`), then the proposal still passes.
	realThreshold := sdk.MinDec(p.Threshold, totalWeight)
	if result.YesCount.GTE(realThreshold) {
		return &DecisionPolicyResult{Allow: true, Final: true}, nil
	}

	totalCounts := result.TotalCounts()
	undecided := totalWeight.Sub(totalCounts)

	// maxYesCount is the max potential number of yes count, i.e the current yes count
	// plus all undecided count (supposing they all vote yes).
	maxYesCount := result.YesCount.Add(undecided)
	if maxYesCount.LT(realThreshold) {
		return &DecisionPolicyResult{Allow: false, Final: true}, nil
	}

	return &DecisionPolicyResult{Final: false}, nil
}

func (p ThresholdDecisionPolicy) GetVotingPeriod() time.Duration {
	return p.Windows.VotingPeriod
}

func (p ThresholdDecisionPolicy) ValidateBasic() error {
	if !p.Threshold.IsPositive() {
		return sdkerrors.ErrInvalidRequest.Wrap("threshold must be a positive number")
	}

	if p.Windows == nil || p.Windows.VotingPeriod == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("voting period cannot be zero")
	}

	return nil
}

var _ DecisionPolicy = (*PercentageDecisionPolicy)(nil)

func (p PercentageDecisionPolicy) Validate(config Config) error {
	if p.Percentage.LT(config.MinPercentage) {
		return sdkerrors.ErrInvalidRequest.Wrap("percentage must be greater than or equal to min_percentage")
	}

	if p.Windows.MinExecutionPeriod > p.Windows.VotingPeriod+config.MaxExecutionPeriod {
		return sdkerrors.ErrInvalidRequest.Wrap("min_execution_period should be smaller than voting_period + max_execution_period")
	}

	return nil
}

func (p PercentageDecisionPolicy) Allow(result TallyResult, totalWeight sdk.Dec, sinceSubmission time.Duration) (*DecisionPolicyResult, error) {
	if sinceSubmission < p.Windows.MinExecutionPeriod {
		return nil, sdkerrors.ErrUnauthorized.Wrapf("must wait %s after submission before execution, currently at %s", p.Windows.MinExecutionPeriod, sinceSubmission)
	}

	yesPercentage := result.YesCount.Quo(totalWeight)
	if yesPercentage.GTE(p.Percentage) {
		return &DecisionPolicyResult{Allow: true, Final: true}, nil
	}

	totalCounts := result.TotalCounts()
	undecided := totalWeight.Sub(totalCounts)
	maxYesCount := result.YesCount.Add(undecided)
	maxYesPercentage := maxYesCount.Quo(totalWeight)
	if maxYesPercentage.LT(p.Percentage) {
		return &DecisionPolicyResult{Allow: false, Final: true}, nil
	}

	return &DecisionPolicyResult{Allow: false, Final: false}, nil
}

func (p PercentageDecisionPolicy) GetVotingPeriod() time.Duration {
	return p.Windows.VotingPeriod
}

func (p PercentageDecisionPolicy) ValidateBasic() error {
	if p.Windows == nil || p.Windows.VotingPeriod == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("voting period cannot be zero")
	}

	if p.Percentage.GT(sdk.OneDec()) || p.Percentage.LTE(sdk.ZeroDec()) {
		return sdkerrors.ErrInvalidRequest.Wrap("percentage must be > 0 and <= 1")
	}

	return nil
}
