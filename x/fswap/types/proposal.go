package types

import (
	"gopkg.in/yaml.v2"

	gov "github.com/Finschia/finschia-sdk/x/gov/types"
)

const (
	ProposalTypeFswapInit string = "FswapInit"
)

// NewFswapInitProposal creates a new FswapInitProposal instance.
// Deprecated: this proposal is considered legacy and is deprecated in favor of
// Msg-based gov proposals. See MsgFswapInit.
func NewFswapInitProposal(title, description string, fswapInit FswapInit) *FswapInitProposal {
	return &FswapInitProposal{title, description, fswapInit}
}

// Implements Proposal Interface
var _ gov.Content = &FswapInitProposal{}

func init() {
	gov.RegisterProposalType(ProposalTypeFswapInit)
}

// ProposalRoute gets the proposal's router key
func (m *FswapInitProposal) ProposalRoute() string { return RouterKey }

// ProposalType is "FswapInit"
func (m *FswapInitProposal) ProposalType() string { return ProposalTypeFswapInit }

// String implements the Stringer interface.
func (m *FswapInitProposal) String() string {
	out, _ := yaml.Marshal(m)
	return string(out)
}

// ValidateBasic validates the proposal
func (m *FswapInitProposal) ValidateBasic() error {
	if err := m.FswapInit.ValidateBasic(); err != nil {
		return err
	}
	return gov.ValidateAbstract(m)
}
