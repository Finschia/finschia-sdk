package foundation

import (
	errorsmod "cosmossdk.io/errors"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeFoundationExec string = "FoundationExec"
)

func NewFoundationExecProposal(title, description string, msgs []sdk.Msg) govtypes.Content {
	proposal := &FoundationExecProposal{
		Title:       title,
		Description: description,
	}
	proposal.SetMessages(msgs)

	return proposal
}

// Implements Proposal Interface
var _ govtypes.Content = &FoundationExecProposal{}

func (p FoundationExecProposal) ProposalRoute() string { return RouterKey }
func (p FoundationExecProposal) ProposalType() string  { return ProposalTypeFoundationExec }
func (p FoundationExecProposal) ValidateBasic() error {
	msgs := GetMessagesFromFoundationExecProposal(p)
	for idx, msg := range msgs {
		m, ok := msg.(sdk.HasValidateBasic)
		if !ok {
			continue
		}

		if err := m.ValidateBasic(); err != nil {
			return errorsmod.Wrapf(err, "msg: %d", idx)
		}
	}

	return govtypes.ValidateAbstract(&p)
}

func GetMessagesFromFoundationExecProposal(p FoundationExecProposal) []sdk.Msg {
	anys := p.Messages
	res := make([]sdk.Msg, len(anys))
	for i, any := range anys {
		cached := any.GetCachedValue()
		if cached == nil {
			panic("Any cached value is nil")
		}
		res[i] = cached.(sdk.Msg)
	}
	return res
}

func (p *FoundationExecProposal) SetMessages(msgs []sdk.Msg) {
	p.Messages = make([]*codectypes.Any, len(msgs))
	for i, msg := range msgs {
		any, err := codectypes.NewAnyWithValue(msg)
		if err != nil {
			panic(err)
		}
		p.Messages[i] = any
	}
}

func (p FoundationExecProposal) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return UnpackInterfaces(unpacker, p.Messages)
}

func init() {
	govtypes.RegisterProposalType(ProposalTypeFoundationExec)
}
