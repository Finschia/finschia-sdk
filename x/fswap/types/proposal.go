package types

import (
	gov "github.com/Finschia/finschia-sdk/x/gov/types"
	"gopkg.in/yaml.v2"
)

const (
	ProposalTypeFswapInit string = "FswapInit"
)

// NewFswapInitProposal creates a new FswapInitProposal instance.
// Deprecated: this proposal is considered legacy and is deprecated in favor of
// Msg-based gov proposals. See MsgFswapInit.
func NewFswapInitProposal(title string, description string, fswapInit FswapInit) *FswapInitProposal {
	return &FswapInitProposal{title, description, fswapInit}
}

// Implements Proposal Interface
var _ gov.Content = &FswapInitProposal{}

func init() {
	gov.RegisterProposalType(ProposalTypeFswapInit)
}

// ProposalRoute gets the proposal's router key
func (fip *FswapInitProposal) ProposalRoute() string { return RouterKey }

// ProposalType is "FswapInit"
func (fip *FswapInitProposal) ProposalType() string { return ProposalTypeFswapInit }

// String implements the Stringer interface.
func (fip *FswapInitProposal) String() string {
	out, _ := yaml.Marshal(fip)
	return string(out)
}

// ValidateBasic validates the proposal
func (fip *FswapInitProposal) ValidateBasic() error {
	if err := fip.FswapInit.ValidateBasic(); err != nil {
		return err
	}
	return gov.ValidateAbstract(fip)
}
