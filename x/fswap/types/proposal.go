package types

import (
	"gopkg.in/yaml.v2"

	gov "github.com/Finschia/finschia-sdk/x/gov/types"
)

const (
	ProposalTypeSwapInit string = "SwapInit"
)

// NewSwapInitProposal creates a new SwapInitProposal instance.
// Deprecated: this proposal is considered legacy and is deprecated in favor of
// Msg-based gov proposals. See MsgSwapInit.
func NewSwapInitProposal(title, description string, swapInit SwapInit) *SwapInitProposal {
	return &SwapInitProposal{title, description, swapInit}
}

// Implements Proposal Interface
var _ gov.Content = &SwapInitProposal{}

func init() {
	gov.RegisterProposalType(ProposalTypeSwapInit)
}

// ProposalRoute gets the proposal's router key
func (m *SwapInitProposal) ProposalRoute() string { return RouterKey }

// ProposalType is "SwapInit"
func (m *SwapInitProposal) ProposalType() string { return ProposalTypeSwapInit }

// String implements the Stringer interface.
func (m *SwapInitProposal) String() string {
	out, _ := yaml.Marshal(m)
	return string(out)
}

// ValidateBasic validates the proposal
func (m *SwapInitProposal) ValidateBasic() error {
	if err := m.SwapInit.ValidateBasic(); err != nil {
		return err
	}
	return gov.ValidateAbstract(m)
}
