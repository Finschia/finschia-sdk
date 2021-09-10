package types

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

type ProposalType string

const (
	ProposalTypeStoreCode            ProposalType = "StoreCode"
	ProposalTypeInstantiateContract  ProposalType = "InstantiateContract"
	ProposalTypeMigrateContract      ProposalType = "MigrateContract"
	ProposalTypeUpdateAdmin          ProposalType = "UpdateAdmin"
	ProposalTypeClearAdmin           ProposalType = "ClearAdmin"
	ProposalTypePinCodes             ProposalType = "PinCodes"
	ProposalTypeUnpinCodes           ProposalType = "UnpinCodes"
	ProposalTypeUpdateContractStatus ProposalType = "UpdateContractStatus"
)

// DisableAllProposals contains no wasm gov types.
var DisableAllProposals []ProposalType

// EnableAllProposals contains all wasm gov types as keys.
var EnableAllProposals = []ProposalType{
	ProposalTypeStoreCode,
	ProposalTypeInstantiateContract,
	ProposalTypeMigrateContract,
	ProposalTypeUpdateAdmin,
	ProposalTypeClearAdmin,
	ProposalTypePinCodes,
	ProposalTypeUnpinCodes,
	ProposalTypeUpdateContractStatus,
}

// ConvertToProposals maps each key to a ProposalType and returns a typed list.
// If any string is not a valid type (in this file), then return an error
func ConvertToProposals(keys []string) ([]ProposalType, error) {
	valid := make(map[string]bool, len(EnableAllProposals))
	for _, key := range EnableAllProposals {
		valid[string(key)] = true
	}

	proposals := make([]ProposalType, len(keys))
	for i, key := range keys {
		if _, ok := valid[key]; !ok {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "'%s' is not a valid ProposalType", key)
		}
		proposals[i] = ProposalType(key)
	}
	return proposals, nil
}

func init() { // register new content types with the sdk
	govtypes.RegisterProposalType(string(ProposalTypeStoreCode))
	govtypes.RegisterProposalType(string(ProposalTypeInstantiateContract))
	govtypes.RegisterProposalType(string(ProposalTypeMigrateContract))
	govtypes.RegisterProposalType(string(ProposalTypeUpdateAdmin))
	govtypes.RegisterProposalType(string(ProposalTypeClearAdmin))
	govtypes.RegisterProposalType(string(ProposalTypePinCodes))
	govtypes.RegisterProposalType(string(ProposalTypeUnpinCodes))
	govtypes.RegisterProposalTypeCodec(&StoreCodeProposal{}, "wasm/StoreCodeProposal")
	govtypes.RegisterProposalTypeCodec(&InstantiateContractProposal{}, "wasm/InstantiateContractProposal")
	govtypes.RegisterProposalTypeCodec(&MigrateContractProposal{}, "wasm/MigrateContractProposal")
	govtypes.RegisterProposalTypeCodec(&UpdateAdminProposal{}, "wasm/UpdateAdminProposal")
	govtypes.RegisterProposalTypeCodec(&ClearAdminProposal{}, "wasm/ClearAdminProposal")
	govtypes.RegisterProposalTypeCodec(&PinCodesProposal{}, "wasm/PinCodesProposal")
	govtypes.RegisterProposalTypeCodec(&UnpinCodesProposal{}, "wasm/UnpinCodesProposal")
	govtypes.RegisterProposalTypeCodec(UpdateContractStatusProposal{}, "wasm/UpdateContractStatusProposal")
}

// ProposalRoute returns the routing key of a parameter change proposal.
func (p StoreCodeProposal) ProposalRoute() string { return RouterKey }

// GetTitle returns the title of the proposal
func (p *StoreCodeProposal) GetTitle() string { return p.Title }

// GetDescription returns the human readable description of the proposal
func (p StoreCodeProposal) GetDescription() string { return p.Description }

// ProposalType returns the type
func (p StoreCodeProposal) ProposalType() string { return string(ProposalTypeStoreCode) }

// ValidateBasic validates the proposal
func (p StoreCodeProposal) ValidateBasic() error {
	if err := validateProposalCommons(p.Title, p.Description); err != nil {
		return err
	}
	if err := sdk.ValidateAccAddress(p.RunAs); err != nil {
		return sdkerrors.Wrap(err, "run as")
	}

	if err := validateWasmCode(p.WASMByteCode); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "code bytes %s", err.Error())
	}

	if err := validateSourceURL(p.Source); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "source %s", err.Error())
	}

	if err := validateBuilder(p.Builder); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "builder %s", err.Error())
	}
	if p.InstantiatePermission != nil {
		if err := p.InstantiatePermission.ValidateBasic(); err != nil {
			return sdkerrors.Wrap(err, "instantiate permission")
		}
	}
	return nil
}

// String implements the Stringer interface.
func (p StoreCodeProposal) String() string {
	return fmt.Sprintf(`Store Code Proposal:
  Title:       %s
  Description: %s
  Run as:      %s
  WasmCode:    %X
  Source:      %s
  Builder:     %s
`, p.Title, p.Description, p.RunAs, p.WASMByteCode, p.Source, p.Builder)
}

// MarshalYAML pretty prints the wasm byte code
func (p StoreCodeProposal) MarshalYAML() (interface{}, error) {
	return struct {
		Title                 string        `yaml:"title"`
		Description           string        `yaml:"description"`
		RunAs                 string        `yaml:"run_as"`
		WASMByteCode          string        `yaml:"wasm_byte_code"`
		Source                string        `yaml:"source"`
		Builder               string        `yaml:"builder"`
		InstantiatePermission *AccessConfig `yaml:"instantiate_permission"`
	}{
		Title:                 p.Title,
		Description:           p.Description,
		RunAs:                 p.RunAs,
		WASMByteCode:          base64.StdEncoding.EncodeToString(p.WASMByteCode),
		Source:                p.Source,
		Builder:               p.Builder,
		InstantiatePermission: p.InstantiatePermission,
	}, nil
}

// ProposalRoute returns the routing key of a parameter change proposal.
func (p InstantiateContractProposal) ProposalRoute() string { return RouterKey }

// GetTitle returns the title of the proposal
func (p *InstantiateContractProposal) GetTitle() string { return p.Title }

// GetDescription returns the human readable description of the proposal
func (p InstantiateContractProposal) GetDescription() string { return p.Description }

// ProposalType returns the type
func (p InstantiateContractProposal) ProposalType() string {
	return string(ProposalTypeInstantiateContract)
}

// ValidateBasic validates the proposal
func (p InstantiateContractProposal) ValidateBasic() error {
	if err := validateProposalCommons(p.Title, p.Description); err != nil {
		return err
	}
	if err := sdk.ValidateAccAddress(p.RunAs); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "run as")
	}

	if p.CodeID == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "code id is required")
	}

	if err := validateLabel(p.Label); err != nil {
		return err
	}

	if !p.Funds.IsValid() {
		return sdkerrors.ErrInvalidCoins
	}

	if len(p.Admin) != 0 {
		if err := sdk.ValidateAccAddress(p.Admin); err != nil {
			return err
		}
	}
	if !json.Valid(p.InitMsg) {
		return sdkerrors.Wrap(ErrInvalid, "init msg json")
	}

	return nil
}

// String implements the Stringer interface.
func (p InstantiateContractProposal) String() string {
	return fmt.Sprintf(`Instantiate Code Proposal:
  Title:       %s
  Description: %s
  Run as:      %s
  Admin:       %s
  Code id:     %d
  Label:       %s
  InitMsg:     %q
  Funds:       %s
`, p.Title, p.Description, p.RunAs, p.Admin, p.CodeID, p.Label, p.InitMsg, p.Funds)
}

// MarshalYAML pretty prints the init message
func (p InstantiateContractProposal) MarshalYAML() (interface{}, error) {
	return struct {
		Title       string    `yaml:"title"`
		Description string    `yaml:"description"`
		RunAs       string    `yaml:"run_as"`
		Admin       string    `yaml:"admin"`
		CodeID      uint64    `yaml:"code_id"`
		Label       string    `yaml:"label"`
		InitMsg     string    `yaml:"init_msg"`
		Funds       sdk.Coins `yaml:"funds"`
	}{
		Title:       p.Title,
		Description: p.Description,
		RunAs:       p.RunAs,
		Admin:       p.Admin,
		CodeID:      p.CodeID,
		Label:       p.Label,
		InitMsg:     string(p.InitMsg),
		Funds:       p.Funds,
	}, nil
}

// ProposalRoute returns the routing key of a parameter change proposal.
func (p MigrateContractProposal) ProposalRoute() string { return RouterKey }

// GetTitle returns the title of the proposal
func (p *MigrateContractProposal) GetTitle() string { return p.Title }

// GetDescription returns the human readable description of the proposal
func (p MigrateContractProposal) GetDescription() string { return p.Description }

// ProposalType returns the type
func (p MigrateContractProposal) ProposalType() string { return string(ProposalTypeMigrateContract) }

// ValidateBasic validates the proposal
func (p MigrateContractProposal) ValidateBasic() error {
	if err := validateProposalCommons(p.Title, p.Description); err != nil {
		return err
	}
	if p.CodeID == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "code_id is required")
	}
	if err := sdk.ValidateAccAddress(p.Contract); err != nil {
		return sdkerrors.Wrap(err, "contract")
	}
	if err := sdk.ValidateAccAddress(p.RunAs); err != nil {
		return sdkerrors.Wrap(err, "run as")
	}
	if !json.Valid(p.MigrateMsg) {
		return sdkerrors.Wrap(ErrInvalid, "migrate msg json")
	}
	return nil
}

// String implements the Stringer interface.
func (p MigrateContractProposal) String() string {
	return fmt.Sprintf(`Migrate Contract Proposal:
  Title:       %s
  Description: %s
  Contract:    %s
  Code id:     %d
  Run as:      %s
  MigrateMsg   %q
`, p.Title, p.Description, p.Contract, p.CodeID, p.RunAs, p.MigrateMsg)
}

// MarshalYAML pretty prints the migrate message
func (p MigrateContractProposal) MarshalYAML() (interface{}, error) {
	return struct {
		Title       string `yaml:"title"`
		Description string `yaml:"description"`
		Contract    string `yaml:"contract"`
		CodeID      uint64 `yaml:"code_id"`
		MigrateMsg  string `yaml:"msg"`
		RunAs       string `yaml:"run_as"`
	}{
		Title:       p.Title,
		Description: p.Description,
		Contract:    p.Contract,
		CodeID:      p.CodeID,
		MigrateMsg:  string(p.MigrateMsg),
		RunAs:       p.RunAs,
	}, nil
}

// ProposalRoute returns the routing key of a parameter change proposal.
func (p UpdateAdminProposal) ProposalRoute() string { return RouterKey }

// GetTitle returns the title of the proposal
func (p *UpdateAdminProposal) GetTitle() string { return p.Title }

// GetDescription returns the human readable description of the proposal
func (p UpdateAdminProposal) GetDescription() string { return p.Description }

// ProposalType returns the type
func (p UpdateAdminProposal) ProposalType() string { return string(ProposalTypeUpdateAdmin) }

// ValidateBasic validates the proposal
func (p UpdateAdminProposal) ValidateBasic() error {
	if err := validateProposalCommons(p.Title, p.Description); err != nil {
		return err
	}
	if err := sdk.ValidateAccAddress(p.Contract); err != nil {
		return sdkerrors.Wrap(err, "contract")
	}
	if err := sdk.ValidateAccAddress(p.NewAdmin); err != nil {
		return sdkerrors.Wrap(err, "new admin")
	}
	return nil
}

// String implements the Stringer interface.
func (p UpdateAdminProposal) String() string {
	return fmt.Sprintf(`Update Contract Admin Proposal:
  Title:       %s
  Description: %s
  Contract:    %s
  New Admin:   %s
`, p.Title, p.Description, p.Contract, p.NewAdmin)
}

// ProposalRoute returns the routing key of a parameter change proposal.
func (p ClearAdminProposal) ProposalRoute() string { return RouterKey }

// GetTitle returns the title of the proposal
func (p *ClearAdminProposal) GetTitle() string { return p.Title }

// GetDescription returns the human readable description of the proposal
func (p ClearAdminProposal) GetDescription() string { return p.Description }

// ProposalType returns the type
func (p ClearAdminProposal) ProposalType() string { return string(ProposalTypeClearAdmin) }

// ValidateBasic validates the proposal
func (p ClearAdminProposal) ValidateBasic() error {
	if err := validateProposalCommons(p.Title, p.Description); err != nil {
		return err
	}
	if err := sdk.ValidateAccAddress(p.Contract); err != nil {
		return sdkerrors.Wrap(err, "contract")
	}
	return nil
}

// String implements the Stringer interface.
func (p ClearAdminProposal) String() string {
	return fmt.Sprintf(`Clear Contract Admin Proposal:
  Title:       %s
  Description: %s
  Contract:    %s
`, p.Title, p.Description, p.Contract)
}

// ProposalRoute returns the routing key of a parameter change proposal.
func (p PinCodesProposal) ProposalRoute() string { return RouterKey }

// GetTitle returns the title of the proposal
func (p *PinCodesProposal) GetTitle() string { return p.Title }

// GetDescription returns the human readable description of the proposal
func (p PinCodesProposal) GetDescription() string { return p.Description }

// ProposalType returns the type
func (p PinCodesProposal) ProposalType() string { return string(ProposalTypePinCodes) }

// ValidateBasic validates the proposal
func (p PinCodesProposal) ValidateBasic() error {
	if err := validateProposalCommons(p.Title, p.Description); err != nil {
		return err
	}
	if len(p.CodeIDs) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "code ids")
	}
	return nil
}

// String implements the Stringer interface.
func (p PinCodesProposal) String() string {
	return fmt.Sprintf(`Pin Wasm Codes Proposal:
  Title:       %s
  Description: %s
  Codes:       %v
`, p.Title, p.Description, p.CodeIDs)
}

// ProposalRoute returns the routing key of a parameter change proposal.
func (p UnpinCodesProposal) ProposalRoute() string { return RouterKey }

// GetTitle returns the title of the proposal
func (p *UnpinCodesProposal) GetTitle() string { return p.Title }

// GetDescription returns the human readable description of the proposal
func (p UnpinCodesProposal) GetDescription() string { return p.Description }

// ProposalType returns the type
func (p UnpinCodesProposal) ProposalType() string { return string(ProposalTypeUnpinCodes) }

// ValidateBasic validates the proposal
func (p UnpinCodesProposal) ValidateBasic() error {
	if err := validateProposalCommons(p.Title, p.Description); err != nil {
		return err
	}
	if len(p.CodeIDs) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "code ids")
	}
	return nil
}

// String implements the Stringer interface.
func (p UnpinCodesProposal) String() string {
	return fmt.Sprintf(`Unpin Wasm Codes Proposal:
  Title:       %s
  Description: %s
  Codes:       %v
`, p.Title, p.Description, p.CodeIDs)
}

// ProposalRoute returns the routing key of a parameter change proposal.
func (p UpdateContractStatusProposal) ProposalRoute() string { return RouterKey }

// GetTitle returns the title of the proposal
func (p *UpdateContractStatusProposal) GetTitle() string { return p.Title }

// GetDescription returns the human readable description of the proposal
func (p UpdateContractStatusProposal) GetDescription() string { return p.Description }

// ProposalType returns the type
func (p UpdateContractStatusProposal) ProposalType() string { return string(ProposalTypeUpdateAdmin) }

// ValidateBasic validates the proposal
func (p UpdateContractStatusProposal) ValidateBasic() error {
	if err := validateProposalCommons(p.Title, p.Description); err != nil {
		return err
	}
	if err := sdk.ValidateAccAddress(p.Contract); err != nil {
		return sdkerrors.Wrap(err, "contract")
	}
	found := false
	for _, v := range AllContractStatus {
		if p.Status == v {
			found = true
			break
		}
	}
	if !found || p.Status == ContractStatusUnspecified {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "invalid status")
	}
	return nil
}

// String implements the Stringer interface.
func (p UpdateContractStatusProposal) String() string {
	return fmt.Sprintf(`Update Contract Status Proposal:
  Title:       %s
  Description: %s
  Contract:    %s
  Status:   %s
`, p.Title, p.Description, p.Contract, p.Status.String())
}

func validateProposalCommons(title, description string) error {
	if strings.TrimSpace(title) != title {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal title must not start/end with white spaces")
	}
	if len(title) == 0 {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal title cannot be blank")
	}
	if len(title) > govtypes.MaxTitleLength {
		return sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "proposal title is longer than max length of %d", govtypes.MaxTitleLength)
	}
	if strings.TrimSpace(description) != description {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal description must not start/end with white spaces")
	}
	if len(description) == 0 {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal description cannot be blank")
	}
	if len(description) > govtypes.MaxDescriptionLength {
		return sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "proposal description is longer than max length of %d", govtypes.MaxDescriptionLength)
	}
	return nil
}
