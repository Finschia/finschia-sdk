package foundation

import (
	"fmt"

	codectypes "github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

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
