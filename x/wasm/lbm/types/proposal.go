package types

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
	wasmtypes "github.com/line/lbm-sdk/x/wasm/types"
)

const (
	ProposalTypeDeactivateContract wasmtypes.ProposalType = "DeactivateContract"
	ProposalTypeActivateContract   wasmtypes.ProposalType = "ActivateContract"
)

var EnableAllProposals = append([]wasmtypes.ProposalType{
	ProposalTypeDeactivateContract,
	ProposalTypeActivateContract,
}, wasmtypes.EnableAllProposals...)

func init() {
	govtypes.RegisterProposalType(string(ProposalTypeDeactivateContract))
	govtypes.RegisterProposalType(string(ProposalTypeActivateContract))
	govtypes.RegisterProposalTypeCodec(&DeactivateContract{}, "wasm/DeactivateContractProposal")
	govtypes.RegisterProposalTypeCodec(&ActivateContract{}, "wasm/ActivateContractProposal")
}

func (p DeactivateContract) GetTitle() string { return p.Title }

func (p DeactivateContract) GetDescription() string { return p.Description }

func (p DeactivateContract) ProposalRoute() string { return wasmtypes.RouterKey }

func (p DeactivateContract) ProposalType() string { return string(ProposalTypeDeactivateContract) }

func (p DeactivateContract) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(p.Contract); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "contract")
	}

	return nil
}

func (p DeactivateContract) String() string {
	return fmt.Sprintf(`Deactivate Contract Proposal:
  Title: 	   %s
  Description: %s
  Contract:    %s
`, p.Title, p.Description, p.Contract)
}

func (p ActivateContract) GetTitle() string { return p.Title }

func (p ActivateContract) GetDescription() string { return p.Description }

func (p ActivateContract) ProposalRoute() string { return wasmtypes.RouterKey }

func (p ActivateContract) ProposalType() string { return string(ProposalTypeActivateContract) }

func (p ActivateContract) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(p.Contract); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "contract")
	}

	return nil
}

func (p ActivateContract) String() string {
	return fmt.Sprintf(`Activate Contract Proposal:
  Title: 	   %s
  Description: %s
  Contract:    %s
`, p.Title, p.Description, p.Contract)
}
