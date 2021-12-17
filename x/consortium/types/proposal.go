package types

import (
	"fmt"

	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

const (
	// ProposalTypeUpdateConsortiumParams updates parameters of consortium.
	ProposalTypeUpdateConsortiumParams = "UpdateConsortiumParams"
	// ProposalTypeUpdateValidatorAuths updates validator authorizations.
	ProposalTypeUpdateValidatorAuths = "UpdateValidatorAuths"
)

func NewUpdateConsortiumParamsProposal(title, description string, params *Params) govtypes.Content {
	return &UpdateConsortiumParamsProposal{title, description, params}
}

// Assert proposals implements govtypes.Content at compile-time
var _ govtypes.Content = &UpdateConsortiumParamsProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeUpdateConsortiumParams)
	govtypes.RegisterProposalTypeCodec(&UpdateConsortiumParamsProposal{}, "lbm-sdk/UpdateConsortiumParamsProposal")
	govtypes.RegisterProposalType(ProposalTypeUpdateValidatorAuths)
	govtypes.RegisterProposalTypeCodec(&UpdateValidatorAuthsProposal{}, "lbm-sdk/UpdateValidatorAuths")
}

func (p *UpdateConsortiumParamsProposal) GetTitle() string       { return p.Title }
func (p *UpdateConsortiumParamsProposal) GetDescription() string { return p.Description }
func (p *UpdateConsortiumParamsProposal) ProposalRoute() string  { return RouterKey }
func (p *UpdateConsortiumParamsProposal) ProposalType() string {
	return ProposalTypeUpdateConsortiumParams
}
func (p *UpdateConsortiumParamsProposal) ValidateBasic() error {
	params := p.Params
	if params.Enabled {
		return ErrInvalidParams
	}

	return nil
}

func (p UpdateConsortiumParamsProposal) String() string {
	return fmt.Sprintf(`Update Consortium Params Proposal:
  Title:       %s
  Description: %s
  Enabled:     %t
`, p.Title, p.Description, p.Params.Enabled)
}

func NewUpdateValidatorAuthsProposal(title, description string,
	auths []*ValidatorAuth) govtypes.Content {
	return &UpdateValidatorAuthsProposal{title, description, auths}
}

var _ govtypes.Content = &UpdateValidatorAuthsProposal{}

func (p *UpdateValidatorAuthsProposal) GetTitle() string       { return p.Title }
func (p *UpdateValidatorAuthsProposal) GetDescription() string { return p.Description }
func (p *UpdateValidatorAuthsProposal) ProposalRoute() string  { return RouterKey }
func (p *UpdateValidatorAuthsProposal) ProposalType() string   { return ProposalTypeUpdateValidatorAuths }

func (p *UpdateValidatorAuthsProposal) ValidateBasic() error {
	if len(p.Auths) == 0 {
		return ErrInvalidProposalValidator
	}

	usedAddrs := map[string]bool{}
	for _, auth := range p.Auths {
		addr := auth.OperatorAddress
		if usedAddrs[addr] {
			return ErrInvalidProposalValidator
		}
		usedAddrs[addr] = true
	}

	return nil
}

func (p UpdateValidatorAuthsProposal) String() string {
	authsStr := "Auths:"
	for _, auth := range p.Auths {
		authsStr += fmt.Sprintf(`
    - OperatorAddress: %s
    - CreationAllowed: %t`,
			auth.OperatorAddress, auth.CreationAllowed)
	}

	return fmt.Sprintf(`Edit Allowed Validators Proposal:
  Title:       %s
  Description: %s
%s
`, p.Title, p.Description, authsStr)
}
