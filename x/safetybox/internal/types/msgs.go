package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewMsgSafetyBoxCreate(safetyBoxID string, safetyBoxOwner sdk.AccAddress, contractID string) MsgSafetyBoxCreate {
	return MsgSafetyBoxCreate{safetyBoxID, safetyBoxOwner, contractID}
}

type MsgSafetyBoxCreate struct {
	SafetyBoxID    string         `json:"safety_box_id"`
	SafetyBoxOwner sdk.AccAddress `json:"safety_box_owner"`
	ContractID     string         `json:"contract_id"`
}

func (msgSbCreate MsgSafetyBoxCreate) Route() string { return RouterKey }

func (msgSbCreate MsgSafetyBoxCreate) Type() string { return MsgTypeSafetyBoxCreate }

func (msgSbCreate MsgSafetyBoxCreate) ValidateBasic() error {
	if len(msgSbCreate.SafetyBoxID) == 0 {
		return ErrSafetyBoxIDRequired
	}

	if msgSbCreate.SafetyBoxOwner.Empty() {
		return ErrSafetyBoxOwnerRequired
	}

	if len(msgSbCreate.ContractID) == 0 {
		return ErrSafetyBoxContractIDRequired
	}

	return nil
}

func (msgSbCreate MsgSafetyBoxCreate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbCreate))
}

func (msgSbCreate MsgSafetyBoxCreate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbCreate.SafetyBoxOwner}
}

func NewMsgSafetyBoxAllocateToken(safetyBoxID string, allocatorAddress sdk.AccAddress, contractID string, amount sdk.Int) MsgSafetyBoxAllocateToken {
	return MsgSafetyBoxAllocateToken{safetyBoxID, allocatorAddress, contractID, amount}
}

type MsgSafetyBoxAllocateToken struct {
	SafetyBoxID      string         `json:"safety_box_id"`
	AllocatorAddress sdk.AccAddress `json:"allocator_address"`
	ContractID       string         `json:"contract_id"`
	Amount           sdk.Int        `json:"amount"`
}

func (msgSbSendToken MsgSafetyBoxAllocateToken) Route() string { return RouterKey }

func (msgSbSendToken MsgSafetyBoxAllocateToken) Type() string { return MsgTypeSafetyBoxAllocateToken }

func (msgSbSendToken MsgSafetyBoxAllocateToken) ValidateBasic() error {
	if len(msgSbSendToken.SafetyBoxID) == 0 {
		return ErrSafetyBoxIDRequired
	}

	if msgSbSendToken.AllocatorAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Allocator address is required")
	}

	if len(msgSbSendToken.ContractID) == 0 {
		return ErrSafetyBoxContractIDRequired
	}

	if !msgSbSendToken.Amount.IsPositive() {
		return sdkerrors.Wrap(ErrSafetyBoxInvalidAmount, msgSbSendToken.Amount.String())
	}

	return nil
}

func (msgSbSendToken MsgSafetyBoxAllocateToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbSendToken))
}

func (msgSbSendToken MsgSafetyBoxAllocateToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbSendToken.AllocatorAddress}
}

func NewMsgSafetyBoxRecallToken(safetyBoxID string, allocatorAddress sdk.AccAddress, contractID string, amount sdk.Int) MsgSafetyBoxRecallToken {
	return MsgSafetyBoxRecallToken{safetyBoxID, allocatorAddress, contractID, amount}
}

type MsgSafetyBoxRecallToken struct {
	SafetyBoxID      string         `json:"safety_box_id"`
	AllocatorAddress sdk.AccAddress `json:"allocator_address"`
	ContractID       string         `json:"contract_id"`
	Amount           sdk.Int        `json:"amount"`
}

func (msgSbSendToken MsgSafetyBoxRecallToken) Route() string { return RouterKey }

func (msgSbSendToken MsgSafetyBoxRecallToken) Type() string { return MsgTypeSafetyBoxRecallToken }

func (msgSbSendToken MsgSafetyBoxRecallToken) ValidateBasic() error {
	if len(msgSbSendToken.SafetyBoxID) == 0 {
		return ErrSafetyBoxIDRequired
	}

	if msgSbSendToken.AllocatorAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Allocator address is required")
	}

	if len(msgSbSendToken.ContractID) == 0 {
		return ErrSafetyBoxContractIDRequired
	}

	if !msgSbSendToken.Amount.IsPositive() {
		return sdkerrors.Wrap(ErrSafetyBoxInvalidAmount, msgSbSendToken.Amount.String())
	}

	return nil
}

func (msgSbSendToken MsgSafetyBoxRecallToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbSendToken))
}

func (msgSbSendToken MsgSafetyBoxRecallToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbSendToken.AllocatorAddress}
}

func NewMsgSafetyBoxIssueToken(safetyBoxID string, fromAddress, toAddress sdk.AccAddress, contractID string, amount sdk.Int) MsgSafetyBoxIssueToken {
	return MsgSafetyBoxIssueToken{safetyBoxID, fromAddress, toAddress, contractID, amount}
}

type MsgSafetyBoxIssueToken struct {
	SafetyBoxID string         `json:"safety_box_id"`
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	ContractID  string         `json:"contract_id"`
	Amount      sdk.Int        `json:"amount"`
}

func (msgSbSendToken MsgSafetyBoxIssueToken) Route() string { return RouterKey }

func (msgSbSendToken MsgSafetyBoxIssueToken) Type() string { return MsgTypeSafetyBoxIssueToken }

func (msgSbSendToken MsgSafetyBoxIssueToken) ValidateBasic() error {
	if len(msgSbSendToken.SafetyBoxID) == 0 {
		return ErrSafetyBoxIDRequired
	}

	if msgSbSendToken.FromAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "From address is required")
	}

	if msgSbSendToken.ToAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "To address is required")
	}

	if len(msgSbSendToken.ContractID) == 0 {
		return ErrSafetyBoxContractIDRequired
	}

	if !msgSbSendToken.Amount.IsPositive() {
		return sdkerrors.Wrap(ErrSafetyBoxInvalidAmount, msgSbSendToken.Amount.String())
	}

	return nil
}

func (msgSbSendToken MsgSafetyBoxIssueToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbSendToken))
}

func (msgSbSendToken MsgSafetyBoxIssueToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbSendToken.FromAddress}
}

func NewMsgSafetyBoxReturnToken(safetyBoxID string, returnerAddress sdk.AccAddress, contractID string, amount sdk.Int) MsgSafetyBoxReturnToken {
	return MsgSafetyBoxReturnToken{safetyBoxID, returnerAddress, contractID, amount}
}

type MsgSafetyBoxReturnToken struct {
	SafetyBoxID     string         `json:"safety_box_id"`
	ReturnerAddress sdk.AccAddress `json:"returner_address"`
	ContractID      string         `json:"contract_id"`
	Amount          sdk.Int        `json:"amount"`
}

func (msgSbSendToken MsgSafetyBoxReturnToken) Route() string { return RouterKey }

func (msgSbSendToken MsgSafetyBoxReturnToken) Type() string { return MsgTypeSafetyBoxReturnToken }

func (msgSbSendToken MsgSafetyBoxReturnToken) ValidateBasic() error {
	if len(msgSbSendToken.SafetyBoxID) == 0 {
		return ErrSafetyBoxIDRequired
	}

	if msgSbSendToken.ReturnerAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Returner address is required")
	}

	if len(msgSbSendToken.ContractID) == 0 {
		return ErrSafetyBoxContractIDRequired
	}

	if !msgSbSendToken.Amount.IsPositive() {
		return sdkerrors.Wrap(ErrSafetyBoxInvalidAmount, msgSbSendToken.Amount.String())
	}

	return nil
}

func (msgSbSendToken MsgSafetyBoxReturnToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbSendToken))
}

func (msgSbSendToken MsgSafetyBoxReturnToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbSendToken.ReturnerAddress}
}

type MsgSafetyBoxRegisterIssuer struct {
	SafetyBoxID string         `json:"safety_box_id"`
	Operator    sdk.AccAddress `json:"operator"`
	Address     sdk.AccAddress `json:"address"`
}

func (MsgSafetyBoxRegisterIssuer) Route() string { return RouterKey }

func (MsgSafetyBoxRegisterIssuer) Type() string { return MsgTypeSafetyBoxGrantIssuerPermission }

func (msgSbPermission MsgSafetyBoxRegisterIssuer) ValidateBasic() error {
	return validateBasic(msgSbPermission.SafetyBoxID, msgSbPermission.Operator, msgSbPermission.Address)
}

func (msgSbPermission MsgSafetyBoxRegisterIssuer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbPermission))
}

func (msgSbPermission MsgSafetyBoxRegisterIssuer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbPermission.Operator}
}

type MsgSafetyBoxRegisterReturner struct {
	SafetyBoxID string         `json:"safety_box_id"`
	Operator    sdk.AccAddress `json:"operator"`
	Address     sdk.AccAddress `json:"address"`
}

func (MsgSafetyBoxRegisterReturner) Route() string { return RouterKey }

func (MsgSafetyBoxRegisterReturner) Type() string { return MsgTypeSafetyBoxGrantReturnerPermission }

func (msgSbPermission MsgSafetyBoxRegisterReturner) ValidateBasic() error {
	return validateBasic(msgSbPermission.SafetyBoxID, msgSbPermission.Operator, msgSbPermission.Address)
}

func (msgSbPermission MsgSafetyBoxRegisterReturner) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbPermission))
}

func (msgSbPermission MsgSafetyBoxRegisterReturner) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbPermission.Operator}
}

type MsgSafetyBoxRegisterAllocator struct {
	SafetyBoxID string         `json:"safety_box_id"`
	Operator    sdk.AccAddress `json:"operator"`
	Address     sdk.AccAddress `json:"address"`
}

func (MsgSafetyBoxRegisterAllocator) Route() string { return RouterKey }

func (MsgSafetyBoxRegisterAllocator) Type() string { return MsgTypeSafetyBoxGrantAllocatorPermission }

func (msgSbPermission MsgSafetyBoxRegisterAllocator) ValidateBasic() error {
	return validateBasic(msgSbPermission.SafetyBoxID, msgSbPermission.Operator, msgSbPermission.Address)
}

func (msgSbPermission MsgSafetyBoxRegisterAllocator) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbPermission))
}

func (msgSbPermission MsgSafetyBoxRegisterAllocator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbPermission.Operator}
}

type MsgSafetyBoxRegisterOperator struct {
	SafetyBoxID    string         `json:"safety_box_id"`
	SafetyBoxOwner sdk.AccAddress `json:"safety_box_owner"`
	Address        sdk.AccAddress `json:"address"`
}

func (MsgSafetyBoxRegisterOperator) Route() string { return RouterKey }

func (MsgSafetyBoxRegisterOperator) Type() string { return MsgTypeSafetyBoxGrantOperatorPermission }

func (msgSbPermission MsgSafetyBoxRegisterOperator) ValidateBasic() error {
	return validateBasic(msgSbPermission.SafetyBoxID, msgSbPermission.SafetyBoxOwner, msgSbPermission.Address)
}

func (msgSbPermission MsgSafetyBoxRegisterOperator) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbPermission))
}

func (msgSbPermission MsgSafetyBoxRegisterOperator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbPermission.SafetyBoxOwner}
}

func validateBasic(sbID string, operator, address sdk.AccAddress) error {
	if len(sbID) == 0 {
		return ErrSafetyBoxIDRequired
	}

	if operator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Operator/SafetyBoxOwner is required")
	}

	if address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Address is required")
	}

	return nil
}

type MsgSafetyBoxDeregisterIssuer struct {
	SafetyBoxID string         `json:"safety_box_id"`
	Operator    sdk.AccAddress `json:"operator"`
	Address     sdk.AccAddress `json:"address"`
}

func (MsgSafetyBoxDeregisterIssuer) Route() string { return RouterKey }

func (MsgSafetyBoxDeregisterIssuer) Type() string { return MsgTypeSafetyBoxRevokeIssuerPermission }

func (msgSbPermission MsgSafetyBoxDeregisterIssuer) ValidateBasic() error {
	return validateBasic(msgSbPermission.SafetyBoxID, msgSbPermission.Operator, msgSbPermission.Address)
}

func (msgSbPermission MsgSafetyBoxDeregisterIssuer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbPermission))
}

func (msgSbPermission MsgSafetyBoxDeregisterIssuer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbPermission.Operator}
}

type MsgSafetyBoxDeregisterReturner struct {
	SafetyBoxID string         `json:"safety_box_id"`
	Operator    sdk.AccAddress `json:"operator"`
	Address     sdk.AccAddress `json:"address"`
}

func (MsgSafetyBoxDeregisterReturner) Route() string { return RouterKey }

func (MsgSafetyBoxDeregisterReturner) Type() string { return MsgTypeSafetyBoxRevokeReturnerPermission }

func (msgSbPermission MsgSafetyBoxDeregisterReturner) ValidateBasic() error {
	return validateBasic(msgSbPermission.SafetyBoxID, msgSbPermission.Operator, msgSbPermission.Address)
}

func (msgSbPermission MsgSafetyBoxDeregisterReturner) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbPermission))
}

func (msgSbPermission MsgSafetyBoxDeregisterReturner) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbPermission.Operator}
}

type MsgSafetyBoxDeregisterAllocator struct {
	SafetyBoxID string         `json:"safety_box_id"`
	Operator    sdk.AccAddress `json:"operator"`
	Address     sdk.AccAddress `json:"address"`
}

func (MsgSafetyBoxDeregisterAllocator) Route() string { return RouterKey }

func (MsgSafetyBoxDeregisterAllocator) Type() string { return MsgTypeSafetyBoxRevokeAllocatorPermission }

func (msgSbPermission MsgSafetyBoxDeregisterAllocator) ValidateBasic() error {
	return validateBasic(msgSbPermission.SafetyBoxID, msgSbPermission.Operator, msgSbPermission.Address)
}

func (msgSbPermission MsgSafetyBoxDeregisterAllocator) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbPermission))
}

func (msgSbPermission MsgSafetyBoxDeregisterAllocator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbPermission.Operator}
}

type MsgSafetyBoxDeregisterOperator struct {
	SafetyBoxID    string         `json:"safety_box_id"`
	SafetyBoxOwner sdk.AccAddress `json:"safety_box_owner"`
	Address        sdk.AccAddress `json:"address"`
}

func (MsgSafetyBoxDeregisterOperator) Route() string { return RouterKey }

func (MsgSafetyBoxDeregisterOperator) Type() string { return MsgTypeSafetyBoxRevokeOperatorPermission }

func (msgSbPermission MsgSafetyBoxDeregisterOperator) ValidateBasic() error {
	return validateBasic(msgSbPermission.SafetyBoxID, msgSbPermission.SafetyBoxOwner, msgSbPermission.Address)
}

func (msgSbPermission MsgSafetyBoxDeregisterOperator) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbPermission))
}

func (msgSbPermission MsgSafetyBoxDeregisterOperator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbPermission.SafetyBoxOwner}
}
