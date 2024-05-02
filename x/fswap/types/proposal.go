package types

import (
	"gopkg.in/yaml.v2"

	gov "github.com/Finschia/finschia-sdk/x/gov/types"
)

const (
	ProposalTypeSwap string = "Swap"
)

// NewSwapProposal creates a new SwapProposal instance.
// Deprecated: this proposal is considered legacy and is deprecated in favor of
// Msg-based gov proposals. See MsgSwap.
func NewSwapProposal(title, description string, swap Swap) *SwapProposal {
	return &SwapProposal{title, description, swap}
}

// Implements Proposal Interface
var _ gov.Content = &SwapProposal{}

func init() {
	gov.RegisterProposalType(ProposalTypeSwap)
}

// ProposalRoute gets the proposal's router key
func (m *SwapProposal) ProposalRoute() string { return RouterKey }

// ProposalType is "Swap"
func (m *SwapProposal) ProposalType() string { return ProposalTypeSwap }

// String implements the Stringer interface.
func (m *SwapProposal) String() string {
	out, _ := yaml.Marshal(m)
	return string(out)
}

// ValidateBasic validates the proposal
func (m *SwapProposal) ValidateBasic() error {
	if err := m.Swap.ValidateBasic(); err != nil {
		return err
	}
	return gov.ValidateAbstract(m)
}
