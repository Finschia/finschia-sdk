package types

import (
	"fmt"
	"strings"

	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

const (
	// ProposalTypeDisableConsortium disables consortium module
	ProposalTypeDisableConsortium = "DisableConsortium"
	// ProposalTypeEditAllowedValidators edits allowed validators
	ProposalTypeEditAllowedValidators = "EditAllowedValidators"
)

func NewDisableConsortiumProposal(title, description string) govtypes.Content {
	return &DisableConsortiumProposal{title, description}
}

// Assert proposals implements govtypes.Content at compile-time
var _ govtypes.Content = &DisableConsortiumProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeDisableConsortium)
	govtypes.RegisterProposalTypeCodec(&DisableConsortiumProposal{}, "lbm-sdk/DisableConsortiumProposal")
	govtypes.RegisterProposalType(ProposalTypeEditAllowedValidators)
	govtypes.RegisterProposalTypeCodec(&EditAllowedValidatorsProposal{}, "lbm-sdk/EditAllowedValidators")
}

func (p *DisableConsortiumProposal) GetTitle() string       { return p.Title }
func (p *DisableConsortiumProposal) GetDescription() string { return p.Description }
func (p *DisableConsortiumProposal) ProposalRoute() string  { return RouterKey }
func (p *DisableConsortiumProposal) ProposalType() string   { return ProposalTypeDisableConsortium }
func (p *DisableConsortiumProposal) ValidateBasic() error   { return nil }

func (p DisableConsortiumProposal) String() string {
	return fmt.Sprintf(`Disable Consortium Proposal:
  Title:       %s
  Description: %s
`, p.Title, p.Description)
}

func NewEditAllowedValidatorsProposal(title, description string,
	adding_validators, removing_validators []string) govtypes.Content {
	return &EditAllowedValidatorsProposal{title, description, adding_validators, removing_validators}
}

var _ govtypes.Content = &EditAllowedValidatorsProposal{}

func (p *EditAllowedValidatorsProposal) GetTitle() string       { return p.Title }
func (p *EditAllowedValidatorsProposal) GetDescription() string { return p.Description }
func (p *EditAllowedValidatorsProposal) ProposalRoute() string  { return RouterKey }
func (p *EditAllowedValidatorsProposal) ProposalType() string   { return ProposalTypeEditAllowedValidators }

func (p *EditAllowedValidatorsProposal) ValidateBasic() error {
	addingEmpty := len(p.AddingAddresses) == 0
	removingEmpty := len(p.RemovingAddresses) == 0
	if addingEmpty && removingEmpty {
		return ErrInvalidProposalValidator
	} else if addingEmpty != removingEmpty {
		return nil
	} else {
		addings := map[string]bool{}
		for _, adding := range p.AddingAddresses {
			addings[adding] = true
		}
		for _, removing := range p.RemovingAddresses {
			if overlapped := addings[removing]; overlapped {
				return ErrInvalidProposalValidator
			}
		}
		return nil
	}
}

func (p EditAllowedValidatorsProposal) String() string {
	var addingStr string
	if len(p.AddingAddresses) != 0 {
		addingStr = fmt.Sprintf(`Adding validators:
    %s
`, strings.Join(p.AddingAddresses, "\n    "))
	}

	var removingStr string
	if len(p.RemovingAddresses) != 0 {
		removingStr = fmt.Sprintf(`Removing validators:
    %s
`, strings.Join(p.RemovingAddresses, "\n    "))
	}

	return fmt.Sprintf(`Edit Allowed Validators Proposal:
  Title:       %s
  Description: %s
%s%s`, p.Title, p.Description, addingStr, removingStr)
}
