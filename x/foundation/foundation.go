package foundation

import (
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"

	"github.com/Finschia/finschia-rdk/codec"
	codectypes "github.com/Finschia/finschia-rdk/codec/types"
	sdk "github.com/Finschia/finschia-rdk/types"
	sdkerrors "github.com/Finschia/finschia-rdk/types/errors"
)

func validateProposers(proposers []string) error {
	if len(proposers) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("no proposers")
	}

	addrs := map[string]bool{}
	for _, proposer := range proposers {
		if addrs[proposer] {
			return sdkerrors.ErrInvalidRequest.Wrapf("duplicated proposer: %s", proposer)
		}
		addrs[proposer] = true

		if _, err := sdk.AccAddressFromBech32(proposer); err != nil {
			return sdkerrors.ErrInvalidAddress.Wrapf("invalid proposer address: %s", proposer)
		}
	}

	return nil
}

func validateMsgs(msgs []sdk.Msg) error {
	if len(msgs) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("no msgs")
	}

	for i, msg := range msgs {
		if err := msg.ValidateBasic(); err != nil {
			return sdkerrors.Wrapf(err, "msg %d", i)
		}
	}

	return nil
}

func validateVoteOption(option VoteOption) error {
	if option == VOTE_OPTION_UNSPECIFIED {
		return sdkerrors.ErrInvalidRequest.Wrap("empty vote option")
	}
	if _, ok := VoteOption_name[int32(option)]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid vote option")
	}

	return nil
}

func (c Censorship) ValidateBasic() error {
	url := c.MsgTypeUrl
	if len(url) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty msg type url")
	}

	authority := c.Authority
	if _, found := CensorshipAuthority_name[int32(authority)]; !found {
		return sdkerrors.ErrInvalidRequest.Wrapf("censorship authority %s over %s", authority, url)
	}

	return nil
}

func (p Params) ValidateBasic() error {
	if err := validateRatio(p.FoundationTax, "tax rate"); err != nil {
		return err
	}

	return nil
}

func (m Member) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid member address: %s", m.Address)
	}

	return nil
}

// ValidateBasic performs stateless validation on a member.
func (m MemberRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid member address: %s", m.Address)
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
	Validate(info FoundationInfo, config Config) error
}

// DefaultTallyResult returns a TallyResult with all counts set to 0.
func DefaultTallyResult() TallyResult {
	return NewTallyResult(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
}

func NewTallyResult(yes, abstain, no, veto sdk.Dec) TallyResult {
	return TallyResult{
		YesCount:        yes,
		AbstainCount:    abstain,
		NoCount:         no,
		NoWithVetoCount: veto,
	}
}

func (t *TallyResult) Add(option VoteOption) error {
	weight := sdk.OneDec()

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
		return sdkerrors.ErrInvalidRequest.Wrapf("unknown vote option %s", option)
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

func (p Proposal) ValidateBasic() error {
	if p.Id == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("id must be > 0")
	}
	if p.FoundationVersion == 0 {
		return sdkerrors.ErrInvalidVersion.Wrap("foundation version must be > 0")
	}
	if err := validateProposers(p.Proposers); err != nil {
		return err
	}
	if err := validateMsgs(p.GetMsgs()); err != nil {
		return err
	}
	return nil
}

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

// for the tests
func (p Proposal) WithMsgs(msgs []sdk.Msg) *Proposal {
	proposal := p
	if err := proposal.SetMsgs(msgs); err != nil {
		return nil
	}
	return &proposal
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

func validateDecisionPolicyWindows(windows DecisionPolicyWindows, config Config) error {
	if windows.MinExecutionPeriod >= windows.VotingPeriod+config.MaxExecutionPeriod {
		return sdkerrors.ErrInvalidRequest.Wrap("min_execution_period should be smaller than voting_period + max_execution_period")
	}
	return nil
}

func validateDecisionPolicyWindowsBasic(windows *DecisionPolicyWindows) error {
	if windows == nil || windows.VotingPeriod == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("voting period cannot be zero")
	}

	return nil
}

var _ DecisionPolicy = (*ThresholdDecisionPolicy)(nil)

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
		return &DecisionPolicyResult{Final: true}, nil
	}

	return &DecisionPolicyResult{}, nil
}

func (p ThresholdDecisionPolicy) GetVotingPeriod() time.Duration {
	return p.Windows.VotingPeriod
}

func (p ThresholdDecisionPolicy) ValidateBasic() error {
	if p.Threshold.IsNil() || !p.Threshold.IsPositive() {
		return sdkerrors.ErrInvalidRequest.Wrap("threshold must be a positive number")
	}

	if err := validateDecisionPolicyWindowsBasic(p.Windows); err != nil {
		return err
	}
	return nil
}

func (p ThresholdDecisionPolicy) Validate(info FoundationInfo, config Config) error {
	if !info.TotalWeight.IsPositive() {
		return sdkerrors.ErrInvalidRequest.Wrapf("total weight must be positive")
	}

	if err := validateDecisionPolicyWindows(*p.Windows, config); err != nil {
		return err
	}

	return nil
}

var _ DecisionPolicy = (*PercentageDecisionPolicy)(nil)

func (p PercentageDecisionPolicy) Allow(result TallyResult, totalWeight sdk.Dec, sinceSubmission time.Duration) (*DecisionPolicyResult, error) {
	if sinceSubmission < p.Windows.MinExecutionPeriod {
		return nil, sdkerrors.ErrUnauthorized.Wrapf("must wait %s after submission before execution, currently at %s", p.Windows.MinExecutionPeriod, sinceSubmission)
	}

	notAbstaining := totalWeight.Sub(result.AbstainCount)
	// If no one votes (everyone abstains), proposal fails
	if notAbstaining.IsZero() {
		return &DecisionPolicyResult{Final: true}, nil
	}

	yesPercentage := result.YesCount.Quo(notAbstaining)
	if yesPercentage.GTE(p.Percentage) {
		return &DecisionPolicyResult{Allow: true, Final: true}, nil
	}

	totalCounts := result.TotalCounts()
	undecided := totalWeight.Sub(totalCounts)
	maxYesCount := result.YesCount.Add(undecided)
	maxYesPercentage := maxYesCount.Quo(notAbstaining)
	if maxYesPercentage.LT(p.Percentage) {
		return &DecisionPolicyResult{Final: true}, nil
	}

	return &DecisionPolicyResult{}, nil
}

func (p PercentageDecisionPolicy) GetVotingPeriod() time.Duration {
	return p.Windows.VotingPeriod
}

func (p PercentageDecisionPolicy) ValidateBasic() error {
	if err := validateDecisionPolicyWindowsBasic(p.Windows); err != nil {
		return err
	}

	if err := validateRatio(p.Percentage, "percentage"); err != nil {
		return err
	}

	return nil
}

func (p PercentageDecisionPolicy) Validate(info FoundationInfo, config Config) error {
	// total weight must be positive, because the admin is group policy
	// (in x/group words)
	if !info.TotalWeight.IsPositive() {
		return sdkerrors.ErrInvalidRequest.Wrapf("total weight must be positive")
	}

	if err := validateDecisionPolicyWindows(*p.Windows, config); err != nil {
		return err
	}

	return nil
}

func validateRatio(ratio sdk.Dec, name string) error {
	if ratio.IsNil() {
		return sdkerrors.ErrInvalidRequest.Wrapf("%s is nil", name)
	}

	if ratio.GT(sdk.OneDec()) || ratio.IsNegative() {
		return sdkerrors.ErrInvalidRequest.Wrapf("%s must be >= 0 and <= 1", name)
	}
	return nil
}

var _ DecisionPolicy = (*OutsourcingDecisionPolicy)(nil)

func (p OutsourcingDecisionPolicy) Allow(result TallyResult, totalWeight sdk.Dec, sinceSubmission time.Duration) (*DecisionPolicyResult, error) {
	return nil, sdkerrors.ErrInvalidRequest.Wrap(p.Description)
}

func (p OutsourcingDecisionPolicy) GetVotingPeriod() time.Duration {
	return 0
}

func (p OutsourcingDecisionPolicy) ValidateBasic() error {
	return nil
}

func (p OutsourcingDecisionPolicy) Validate(info FoundationInfo, config Config) error {
	return sdkerrors.ErrInvalidRequest.Wrap(p.Description)
}

var _ codectypes.UnpackInterfacesMessage = (*FoundationInfo)(nil)

func (i FoundationInfo) GetDecisionPolicy() DecisionPolicy {
	if i.DecisionPolicy == nil {
		return nil
	}

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

// for the tests
func (i FoundationInfo) WithDecisionPolicy(policy DecisionPolicy) *FoundationInfo {
	info := i
	if err := info.SetDecisionPolicy(policy); err != nil {
		return nil
	}
	return &info
}

func (i *FoundationInfo) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var policy DecisionPolicy
	return unpacker.UnpackAny(i.DecisionPolicy, &policy)
}

func GetAuthorization(any *codectypes.Any, name string) (Authorization, error) {
	cached := any.GetCachedValue()
	if cached == nil {
		return nil, fmt.Errorf("any cached value is nil, %s authorization must be correctly packed any values", name)
	}

	a, ok := cached.(Authorization)
	if !ok {
		return nil, sdkerrors.ErrInvalidType.Wrapf("can't proto unmarshal %T", a)
	}
	return a, nil
}

func SetAuthorization(a Authorization) (*codectypes.Any, error) {
	msg, ok := a.(proto.Message)
	if !ok {
		return nil, sdkerrors.ErrInvalidType.Wrapf("can't proto marshal %T", msg)
	}

	any, err := codectypes.NewAnyWithValue(a)
	if err != nil {
		return nil, err
	}
	return any, nil
}

func (p Pool) ValidateBasic() error {
	if err := p.Treasury.Validate(); err != nil {
		return err
	}

	return nil
}

// Members defines a repeated slice of Member objects.
type Members struct {
	Members []Member
}

// ValidateBasic performs stateless validation on an array of members. On top
// of validating each member individually, it also makes sure there are no
// duplicate addresses.
func (ms Members) ValidateBasic() error {
	index := make(map[string]struct{}, len(ms.Members))
	for i := range ms.Members {
		member := ms.Members[i]
		if err := member.ValidateBasic(); err != nil {
			return err
		}
		addr := member.Address
		if _, exists := index[addr]; exists {
			return sdkerrors.ErrInvalidRequest.Wrapf("duplicated address: %s", member.Address)
		}
		index[addr] = struct{}{}
	}
	return nil
}

// MemberRequests defines a repeated slice of MemberRequest objects.
type MemberRequests struct {
	Members []MemberRequest
}

// ValidateBasic performs stateless validation on an array of members. On top
// of validating each member individually, it also makes sure there are no
// duplicate addresses.
func (ms MemberRequests) ValidateBasic() error {
	index := make(map[string]struct{}, len(ms.Members))
	for i := range ms.Members {
		member := ms.Members[i]
		if err := member.ValidateBasic(); err != nil {
			return err
		}
		addr := member.Address
		if _, exists := index[addr]; exists {
			return sdkerrors.ErrInvalidRequest.Wrapf("duplicated address: %s", member.Address)
		}
		index[addr] = struct{}{}
	}
	return nil
}
