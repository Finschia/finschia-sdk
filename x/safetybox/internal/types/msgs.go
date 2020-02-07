package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMsgSafetyBoxCreate(safetyBoxID string, safetyBoxOwner sdk.AccAddress, safetyBoxDenoms []string) MsgSafetyBoxCreate {
	return MsgSafetyBoxCreate{safetyBoxID, safetyBoxOwner, safetyBoxDenoms}
}

type MsgSafetyBoxCreate struct {
	SafetyBoxID     string         `json:"safety_box_id"`
	SafetyBoxOwner  sdk.AccAddress `json:"safety_box_owner"`
	SafetyBoxDenoms []string       `json:"safety_box_denoms"`
}

func (msgSbCreate MsgSafetyBoxCreate) Route() string { return RouterKey }

func (msgSbCreate MsgSafetyBoxCreate) Type() string { return MsgTypeSafetyBoxCreate }

func (msgSbCreate MsgSafetyBoxCreate) ValidateBasic() sdk.Error {
	if len(msgSbCreate.SafetyBoxID) == 0 {
		return ErrSafetyBoxIDRequired(DefaultCodespace)
	}

	if msgSbCreate.SafetyBoxOwner.Empty() {
		return ErrSafetyBoxOwnerRequired(DefaultCodespace)
	}

	if len(msgSbCreate.SafetyBoxDenoms) == 0 {
		return ErrSafetyBoxDenomRequired(DefaultCodespace)
	}

	return nil
}

func (msgSbCreate MsgSafetyBoxCreate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbCreate))
}

func (msgSbCreate MsgSafetyBoxCreate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbCreate.SafetyBoxOwner}
}

func NewMsgSafetyBoxAllocateCoins(safetyBoxID string, allocatorAddress sdk.AccAddress, coins sdk.Coins) MsgSafetyBoxAllocateCoins {
	return MsgSafetyBoxAllocateCoins{safetyBoxID, allocatorAddress, coins}
}

type MsgSafetyBoxAllocateCoins struct {
	SafetyBoxID      string         `json:"safety_box_id"`
	AllocatorAddress sdk.AccAddress `json:"allocator_address"`
	Coins            sdk.Coins      `json:"coins"`
}

func (msgSbSendCoin MsgSafetyBoxAllocateCoins) Route() string { return RouterKey }

func (msgSbSendCoin MsgSafetyBoxAllocateCoins) Type() string { return MsgTypeSafetyBoxAllocateCoin }

func (msgSbSendCoin MsgSafetyBoxAllocateCoins) ValidateBasic() sdk.Error {
	if len(msgSbSendCoin.SafetyBoxID) == 0 {
		return ErrSafetyBoxIDRequired(DefaultCodespace)
	}

	if msgSbSendCoin.AllocatorAddress.Empty() {
		return sdk.ErrUnknownRequest("Allocator address is required")
	}

	if msgSbSendCoin.Coins.Empty() {
		return ErrSafetyBoxCoinsRequired(DefaultCodespace)
	}

	return nil
}

func (msgSbSendCoin MsgSafetyBoxAllocateCoins) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbSendCoin))
}

func (msgSbSendCoin MsgSafetyBoxAllocateCoins) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbSendCoin.AllocatorAddress}
}

func NewMsgSafetyBoxRecallCoins(safetyBoxID string, allocatorAddress sdk.AccAddress, coins sdk.Coins) MsgSafetyBoxRecallCoins {
	return MsgSafetyBoxRecallCoins{safetyBoxID, allocatorAddress, coins}
}

type MsgSafetyBoxRecallCoins struct {
	SafetyBoxID      string         `json:"safety_box_id"`
	AllocatorAddress sdk.AccAddress `json:"allocator_address"`
	Coins            sdk.Coins      `json:"coins"`
}

func (msgSbSendCoin MsgSafetyBoxRecallCoins) Route() string { return RouterKey }

func (msgSbSendCoin MsgSafetyBoxRecallCoins) Type() string { return MsgTypeSafetyBoxRecallCoin }

func (msgSbSendCoin MsgSafetyBoxRecallCoins) ValidateBasic() sdk.Error {
	if len(msgSbSendCoin.SafetyBoxID) == 0 {
		return ErrSafetyBoxIDRequired(DefaultCodespace)
	}

	if msgSbSendCoin.AllocatorAddress.Empty() {
		return sdk.ErrUnknownRequest("Allocator address is required")
	}

	if msgSbSendCoin.Coins.Empty() {
		return ErrSafetyBoxCoinsRequired(DefaultCodespace)
	}

	return nil
}

func (msgSbSendCoin MsgSafetyBoxRecallCoins) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbSendCoin))
}

func (msgSbSendCoin MsgSafetyBoxRecallCoins) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbSendCoin.AllocatorAddress}
}

func NewMsgSafetyBoxIssueCoins(safetyBoxID string, fromAddress, toAddress sdk.AccAddress, coins sdk.Coins) MsgSafetyBoxIssueCoins {
	return MsgSafetyBoxIssueCoins{safetyBoxID, fromAddress, toAddress, coins}
}

type MsgSafetyBoxIssueCoins struct {
	SafetyBoxID string         `json:"safety_box_id"`
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	Coins       sdk.Coins      `json:"coins"`
}

func (msgSbSendCoin MsgSafetyBoxIssueCoins) Route() string { return RouterKey }

func (msgSbSendCoin MsgSafetyBoxIssueCoins) Type() string { return MsgTypeSafetyBoxIssueCoin }

func (msgSbSendCoin MsgSafetyBoxIssueCoins) ValidateBasic() sdk.Error {
	if len(msgSbSendCoin.SafetyBoxID) == 0 {
		return ErrSafetyBoxIDRequired(DefaultCodespace)
	}

	if msgSbSendCoin.FromAddress.Empty() {
		return sdk.ErrUnknownRequest("From address is required")
	}

	if msgSbSendCoin.ToAddress.Empty() {
		return sdk.ErrUnknownRequest("To address is required")
	}

	if msgSbSendCoin.Coins.Empty() {
		return ErrSafetyBoxCoinsRequired(DefaultCodespace)
	}

	return nil
}

func (msgSbSendCoin MsgSafetyBoxIssueCoins) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbSendCoin))
}

func (msgSbSendCoin MsgSafetyBoxIssueCoins) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbSendCoin.FromAddress}
}

func NewMsgSafetyBoxReturnCoins(safetyBoxID string, returnerAddress sdk.AccAddress, coins sdk.Coins) MsgSafetyBoxReturnCoins {
	return MsgSafetyBoxReturnCoins{safetyBoxID, returnerAddress, coins}
}

type MsgSafetyBoxReturnCoins struct {
	SafetyBoxID     string         `json:"safety_box_id"`
	ReturnerAddress sdk.AccAddress `json:"returner_address"`
	Coins           sdk.Coins      `json:"coins"`
}

func (msgSbSendCoin MsgSafetyBoxReturnCoins) Route() string { return RouterKey }

func (msgSbSendCoin MsgSafetyBoxReturnCoins) Type() string { return MsgTypeSafetyBoxReturnCoin }

func (msgSbSendCoin MsgSafetyBoxReturnCoins) ValidateBasic() sdk.Error {
	if len(msgSbSendCoin.SafetyBoxID) == 0 {
		return ErrSafetyBoxIDRequired(DefaultCodespace)
	}

	if msgSbSendCoin.ReturnerAddress.Empty() {
		return sdk.ErrUnknownRequest("Returner address is required")
	}

	if msgSbSendCoin.Coins.Empty() {
		return ErrSafetyBoxCoinsRequired(DefaultCodespace)
	}

	return nil
}

func (msgSbSendCoin MsgSafetyBoxReturnCoins) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbSendCoin))
}

func (msgSbSendCoin MsgSafetyBoxReturnCoins) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbSendCoin.ReturnerAddress}
}

type MsgSafetyBoxRegisterIssuer struct {
	SafetyBoxID string         `json:"safety_box_id"`
	Operator    sdk.AccAddress `json:"operator"`
	Address     sdk.AccAddress `json:"address"`
}

func (MsgSafetyBoxRegisterIssuer) Route() string { return RouterKey }

func (MsgSafetyBoxRegisterIssuer) Type() string { return MsgTypeSafetyBoxGrantIssuerPermission }

func (msgSbPermission MsgSafetyBoxRegisterIssuer) ValidateBasic() sdk.Error {
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

func (msgSbPermission MsgSafetyBoxRegisterReturner) ValidateBasic() sdk.Error {
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

func (msgSbPermission MsgSafetyBoxRegisterAllocator) ValidateBasic() sdk.Error {
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

func (msgSbPermission MsgSafetyBoxRegisterOperator) ValidateBasic() sdk.Error {
	return validateBasic(msgSbPermission.SafetyBoxID, msgSbPermission.SafetyBoxOwner, msgSbPermission.Address)
}

func (msgSbPermission MsgSafetyBoxRegisterOperator) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbPermission))
}

func (msgSbPermission MsgSafetyBoxRegisterOperator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbPermission.SafetyBoxOwner}
}

func validateBasic(sbID string, operator, address sdk.AccAddress) sdk.Error {
	if len(sbID) == 0 {
		return ErrSafetyBoxIDRequired(DefaultCodespace)
	}

	if operator.Empty() {
		return sdk.ErrInvalidAddress("Operator/SafetyBoxOwner is required")
	}

	if address.Empty() {
		return sdk.ErrInvalidAddress("Address is required")
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

func (msgSbPermission MsgSafetyBoxDeregisterIssuer) ValidateBasic() sdk.Error {
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

func (msgSbPermission MsgSafetyBoxDeregisterReturner) ValidateBasic() sdk.Error {
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

func (msgSbPermission MsgSafetyBoxDeregisterAllocator) ValidateBasic() sdk.Error {
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

func (msgSbPermission MsgSafetyBoxDeregisterOperator) ValidateBasic() sdk.Error {
	return validateBasic(msgSbPermission.SafetyBoxID, msgSbPermission.SafetyBoxOwner, msgSbPermission.Address)
}

func (msgSbPermission MsgSafetyBoxDeregisterOperator) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbPermission))
}

func (msgSbPermission MsgSafetyBoxDeregisterOperator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbPermission.SafetyBoxOwner}
}
