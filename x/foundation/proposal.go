package foundation

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

const (
	// ProposalTypeUpdateFoundationParams updates parameters of foundation.
	ProposalTypeUpdateFoundationParams = "UpdateFoundationParams"
	// ProposalTypeUpdateValidatorAuths updates validator authorizations.
	ProposalTypeUpdateValidatorAuths = "UpdateValidatorAuths"
)

func NewUpdateFoundationParamsProposal(title, description string, params *Params) govtypes.Content {
	return &UpdateFoundationParamsProposal{title, description, params}
}

// Assert proposals implements govtypes.Content at compile-time
var _ govtypes.Content = (*UpdateFoundationParamsProposal)(nil)

func init() {
	govtypes.RegisterProposalType(ProposalTypeUpdateFoundationParams)
	govtypes.RegisterProposalType(ProposalTypeUpdateValidatorAuths)
}

func (p *UpdateFoundationParamsProposal) GetTitle() string       { return p.Title }
func (p *UpdateFoundationParamsProposal) GetDescription() string { return p.Description }
func (p *UpdateFoundationParamsProposal) ProposalRoute() string  { return RouterKey }
func (p *UpdateFoundationParamsProposal) ProposalType() string {
	return ProposalTypeUpdateFoundationParams
}
func (p *UpdateFoundationParamsProposal) ValidateBasic() error {
	params := p.Params
	if params.Enabled {
		return sdkerrors.ErrInvalidRequest.Wrap("cannot enable foundation")
	}

	if err := validateRatio(params.FoundationTax, "tax rate"); err != nil {
		return err
	}

	return nil
}

func (p UpdateFoundationParamsProposal) String() string {
	return fmt.Sprintf(`Update Foundation Params Proposal:
  Title:       %s
  Description: %s
  Enabled:     %t
`, p.Title, p.Description, p.Params.Enabled)
}

func NewUpdateValidatorAuthsProposal(title, description string,
	auths []ValidatorAuth) govtypes.Content {
	return &UpdateValidatorAuthsProposal{title, description, auths}
}

var _ govtypes.Content = (*UpdateValidatorAuthsProposal)(nil)

func (p *UpdateValidatorAuthsProposal) GetTitle() string       { return p.Title }
func (p *UpdateValidatorAuthsProposal) GetDescription() string { return p.Description }
func (p *UpdateValidatorAuthsProposal) ProposalRoute() string  { return RouterKey }
func (p *UpdateValidatorAuthsProposal) ProposalType() string   { return ProposalTypeUpdateValidatorAuths }

func (p *UpdateValidatorAuthsProposal) ValidateBasic() error {
	if len(p.Auths) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty auths")
	}

	usedAddrs := map[string]bool{}
	for _, auth := range p.Auths {
		addr := auth.OperatorAddress
		if _, err := sdk.ValAddressFromBech32(addr); err != nil {
			return err
		}
		if usedAddrs[addr] {
			return sdkerrors.ErrInvalidRequest.Wrapf("multiple auths for same validator: %s", addr)
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

	return fmt.Sprintf(`Update Validator Auths Proposal:
  Title:       %s
  Description: %s
%s
`, p.Title, p.Description, authsStr)
}
