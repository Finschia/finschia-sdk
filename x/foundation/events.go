package foundation

import (
	"github.com/gogo/protobuf/proto"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (m EventUpdateDecisionPolicy) GetDecisionPolicy() DecisionPolicy {
	if m.DecisionPolicy == nil {
		return nil
	}

	policy, ok := m.DecisionPolicy.GetCachedValue().(DecisionPolicy)
	if !ok {
		return nil
	}
	return policy
}

func (m *EventUpdateDecisionPolicy) SetDecisionPolicy(policy DecisionPolicy) error {
	event, ok := policy.(proto.Message)
	if !ok {
		return sdkerrors.ErrInvalidType.Wrapf("can't proto marshal %T", event)
	}

	any, err := codectypes.NewAnyWithValue(event)
	if err != nil {
		return err
	}
	m.DecisionPolicy = any

	return nil
}

func (m EventUpdateDecisionPolicy) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var policy DecisionPolicy
	return unpacker.UnpackAny(m.DecisionPolicy, &policy)
}
