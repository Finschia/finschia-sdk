package types

import (
	"gopkg.in/yaml.v2"

	bank "github.com/Finschia/finschia-sdk/x/bank/types"
	gov "github.com/Finschia/finschia-sdk/x/gov/types"
)

const (
	ProposalTypeSwap string = "Swap"
)

// NewMakeSwapProposal creates a new SwapProposal instance.
// Deprecated: this proposal is considered legacy and is deprecated in favor of
// Msg-based gov proposals. See MsgSwap.
func NewMakeSwapProposal(title, description string, swap Swap, toDenomMetadata bank.Metadata) *MakeSwapProposal {
	return &MakeSwapProposal{title, description, swap, toDenomMetadata}
}

// Implements Proposal Interface
var _ gov.Content = &MakeSwapProposal{}

func init() {
	gov.RegisterProposalType(ProposalTypeSwap)
}

// ProposalRoute gets the proposal's router key
func (m *MakeSwapProposal) ProposalRoute() string { return RouterKey }

// ProposalType is "Swap"
func (m *MakeSwapProposal) ProposalType() string { return ProposalTypeSwap }

// String implements the Stringer interface.
func (m *MakeSwapProposal) String() string {
	out, _ := yaml.Marshal(m)
	return string(out)
}

// ValidateBasic validates the proposal
func (m *MakeSwapProposal) ValidateBasic() error {
	if err := m.Swap.ValidateBasic(); err != nil {
		return err
	}
	return gov.ValidateAbstract(m)
}
