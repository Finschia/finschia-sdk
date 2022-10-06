package foundation

import (
	"fmt"

	sdkerrors "github.com/line/lbm-sdk/types/errors"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

const (
	// ProposalTypeUpdateFoundationParams updates parameters of foundation.
	ProposalTypeUpdateFoundationParams = "UpdateFoundationParams"
)

func NewUpdateFoundationParamsProposal(title, description string, params *Params) govtypes.Content {
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
