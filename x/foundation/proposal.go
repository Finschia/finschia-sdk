package foundation

import (
	yaml "gopkg.in/yaml.v2"

	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

const (
	// ProposalTypeUpdateFoundationParams updates parameters of foundation.
	ProposalTypeUpdateFoundationParams = "UpdateFoundationParams"
)

func NewUpdateFoundationParamsProposal(title, description string, params Params) govtypes.Content {
	return &UpdateFoundationParamsProposal{title, description, params}
}

// Assert proposals implements govtypes.Content at compile-time
var _ govtypes.Content = (*UpdateFoundationParamsProposal)(nil)

func init() {
	govtypes.RegisterProposalType(ProposalTypeUpdateFoundationParams)
}

func (p *UpdateFoundationParamsProposal) GetTitle() string       { return p.Title }
func (p *UpdateFoundationParamsProposal) GetDescription() string { return p.Description }
func (p *UpdateFoundationParamsProposal) ProposalRoute() string  { return RouterKey }
func (p *UpdateFoundationParamsProposal) ProposalType() string {
	return ProposalTypeUpdateFoundationParams
}
func (p *UpdateFoundationParamsProposal) ValidateBasic() error {
	if err := p.Params.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

func (p UpdateFoundationParamsProposal) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
