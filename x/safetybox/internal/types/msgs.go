package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMsgSafetyBoxCreate(safetyBoxId string, safetyBoxOwner sdk.AccAddress) MsgSafetyBoxCreate {
	return MsgSafetyBoxCreate{safetyBoxId, safetyBoxOwner}
}

type MsgSafetyBoxCreate struct {
	SafetyBoxId    string         `json:"safety_box_id"`
	SafetyBoxOwner sdk.AccAddress `json:"safety_box_owner"`
}

func (msgSbCreate MsgSafetyBoxCreate) Route() string { return RouterKey }

func (msgSbCreate MsgSafetyBoxCreate) Type() string { return MsgTypeSafetyBoxCreate }

func (msgSbCreate MsgSafetyBoxCreate) ValidateBasic() sdk.Error {
	if len(msgSbCreate.SafetyBoxId) == 0 {
		return ErrSafetyBoxIdRequired(DefaultCodespace)
	}

	if msgSbCreate.SafetyBoxOwner.Empty() {
		return ErrSafetyBoxOwnerRequired(DefaultCodespace)
	}

	return nil
}

func (msgSbCreate MsgSafetyBoxCreate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbCreate))
}

func (msgSbCreate MsgSafetyBoxCreate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbCreate.SafetyBoxOwner}
}

func NewMsgSafetyBoxAllocateCoins(safetyBoxId string, allocatorAddress sdk.AccAddress, coins sdk.Coins) MsgSafetyBoxAllocateCoins {
	return MsgSafetyBoxAllocateCoins{safetyBoxId, allocatorAddress, coins}
}

type MsgSafetyBoxAllocateCoins struct {
	SafetyBoxId      string         `json:"safety_box_id"`
	AllocatorAddress sdk.AccAddress `json:"allocator_address"`
	Coins            sdk.Coins      `json:"coins"`
}

func (msgSbSendCoin MsgSafetyBoxAllocateCoins) Route() string { return RouterKey }

func (msgSbSendCoin MsgSafetyBoxAllocateCoins) Type() string { return MsgTypeSafetyBoxAllocateCoin }

func (msgSbSendCoin MsgSafetyBoxAllocateCoins) ValidateBasic() sdk.Error {
	if len(msgSbSendCoin.SafetyBoxId) == 0 {
		return ErrSafetyBoxIdRequired(DefaultCodespace)
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

func NewMsgSafetyBoxRecallCoins(safetyBoxId string, allocatorAddress sdk.AccAddress, coins sdk.Coins) MsgSafetyBoxRecallCoins {
	return MsgSafetyBoxRecallCoins{safetyBoxId, allocatorAddress, coins}
}

type MsgSafetyBoxRecallCoins struct {
	SafetyBoxId      string         `json:"safety_box_id"`
	AllocatorAddress sdk.AccAddress `json:"allocator_address"`
	Coins            sdk.Coins      `json:"coins"`
}

func (msgSbSendCoin MsgSafetyBoxRecallCoins) Route() string { return RouterKey }

func (msgSbSendCoin MsgSafetyBoxRecallCoins) Type() string { return MsgTypeSafetyBoxRecallCoin }

func (msgSbSendCoin MsgSafetyBoxRecallCoins) ValidateBasic() sdk.Error {
	if len(msgSbSendCoin.SafetyBoxId) == 0 {
		return ErrSafetyBoxIdRequired(DefaultCodespace)
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

func NewMsgSafetyBoxIssueCoins(safetyBoxId string, fromAddress, toAddress sdk.AccAddress, coins sdk.Coins) MsgSafetyBoxIssueCoins {
	return MsgSafetyBoxIssueCoins{safetyBoxId, fromAddress, toAddress, coins}
}

type MsgSafetyBoxIssueCoins struct {
	SafetyBoxId string         `json:"safety_box_id"`
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	Coins       sdk.Coins      `json:"coins"`
}

func (msgSbSendCoin MsgSafetyBoxIssueCoins) Route() string { return RouterKey }

func (msgSbSendCoin MsgSafetyBoxIssueCoins) Type() string { return MsgTypeSafetyBoxIssueCoin }

func (msgSbSendCoin MsgSafetyBoxIssueCoins) ValidateBasic() sdk.Error {
	if len(msgSbSendCoin.SafetyBoxId) == 0 {
		return ErrSafetyBoxIdRequired(DefaultCodespace)
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

func NewMsgSafetyBoxReturnCoins(safetyBoxId string, returnerAddress sdk.AccAddress, coins sdk.Coins) MsgSafetyBoxReturnCoins {
	return MsgSafetyBoxReturnCoins{safetyBoxId, returnerAddress, coins}
}

type MsgSafetyBoxReturnCoins struct {
	SafetyBoxId     string         `json:"safety_box_id"`
	ReturnerAddress sdk.AccAddress `json:"returner_address"`
	Coins           sdk.Coins      `json:"coins"`
}

func (msgSbSendCoin MsgSafetyBoxReturnCoins) Route() string { return RouterKey }

func (msgSbSendCoin MsgSafetyBoxReturnCoins) Type() string { return MsgTypeSafetyBoxReturnCoin }

func (msgSbSendCoin MsgSafetyBoxReturnCoins) ValidateBasic() sdk.Error {
	if len(msgSbSendCoin.SafetyBoxId) == 0 {
		return ErrSafetyBoxIdRequired(DefaultCodespace)
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

func NewMsgSafetyBoxRegisterIssuer(safetyBoxId string, operator, address sdk.AccAddress) MsgSafetyBoxRegisterIssuer {
	return MsgSafetyBoxRegisterIssuer{safetyBoxId, operator, address}
}

type MsgSafetyBoxRegisterIssuer struct {
	SafetyBoxId string         `json:"safety_box_id"`
	Operator    sdk.AccAddress `json:"operator"`
	Address     sdk.AccAddress `json:"address"`
}

func (MsgSafetyBoxRegisterIssuer) Route() string { return RouterKey }

func (MsgSafetyBoxRegisterIssuer) Type() string { return MsgTypeSafetyBoxGrantIssuerPermission }

func (msgSbPermission MsgSafetyBoxRegisterIssuer) ValidateBasic() sdk.Error {
	return validateBasic(msgSbPermission.SafetyBoxId, msgSbPermission.Operator, msgSbPermission.Address)
}

func (msgSbPermission MsgSafetyBoxRegisterIssuer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbPermission))
}

func (msgSbPermission MsgSafetyBoxRegisterIssuer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbPermission.Operator}
}

func NewMsgSafetyBoxRegisterReturner(safetyBoxId string, operator, address sdk.AccAddress) MsgSafetyBoxRegisterReturner {
	return MsgSafetyBoxRegisterReturner{safetyBoxId, operator, address}
}

type MsgSafetyBoxRegisterReturner struct {
	SafetyBoxId string         `json:"safety_box_id"`
	Operator    sdk.AccAddress `json:"operator"`
	Address     sdk.AccAddress `json:"address"`
}

func (MsgSafetyBoxRegisterReturner) Route() string { return RouterKey }

func (MsgSafetyBoxRegisterReturner) Type() string { return MsgTypeSafetyBoxGrantReturnerPermission }

func (msgSbPermission MsgSafetyBoxRegisterReturner) ValidateBasic() sdk.Error {
	return validateBasic(msgSbPermission.SafetyBoxId, msgSbPermission.Operator, msgSbPermission.Address)
}

func (msgSbPermission MsgSafetyBoxRegisterReturner) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbPermission))
}

func (msgSbPermission MsgSafetyBoxRegisterReturner) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbPermission.Operator}
}

func NewMsgSafetyBoxRegisterAllocator(safetyBoxId string, operator, address sdk.AccAddress) MsgSafetyBoxRegisterAllocator {
	return MsgSafetyBoxRegisterAllocator{safetyBoxId, operator, address}
}

type MsgSafetyBoxRegisterAllocator struct {
	SafetyBoxId string         `json:"safety_box_id"`
	Operator    sdk.AccAddress `json:"operator"`
	Address     sdk.AccAddress `json:"address"`
}

func (MsgSafetyBoxRegisterAllocator) Route() string { return RouterKey }

func (MsgSafetyBoxRegisterAllocator) Type() string { return MsgTypeSafetyBoxGrantAllocatorPermission }

func (msgSbPermission MsgSafetyBoxRegisterAllocator) ValidateBasic() sdk.Error {
	return validateBasic(msgSbPermission.SafetyBoxId, msgSbPermission.Operator, msgSbPermission.Address)
}

func (msgSbPermission MsgSafetyBoxRegisterAllocator) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbPermission))
}

func (msgSbPermission MsgSafetyBoxRegisterAllocator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbPermission.Operator}
}

func NewMsgSafetyBoxRegisterOperator(safetyBoxId string, safetyBoxOwner, address sdk.AccAddress) MsgSafetyBoxRegisterOperator {
	return MsgSafetyBoxRegisterOperator{safetyBoxId, safetyBoxOwner, address}
}

type MsgSafetyBoxRegisterOperator struct {
	SafetyBoxId    string         `json:"safety_box_id"`
	SafetyBoxOwner sdk.AccAddress `json:"safety_box_owner"`
	Address        sdk.AccAddress `json:"address"`
}

func (MsgSafetyBoxRegisterOperator) Route() string { return RouterKey }

func (MsgSafetyBoxRegisterOperator) Type() string { return MsgTypeSafetyBoxGrantOperatorPermission }

func (msgSbPermission MsgSafetyBoxRegisterOperator) ValidateBasic() sdk.Error {
	return validateBasic(msgSbPermission.SafetyBoxId, msgSbPermission.SafetyBoxOwner, msgSbPermission.Address)
}

func (msgSbPermission MsgSafetyBoxRegisterOperator) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbPermission))
}

func (msgSbPermission MsgSafetyBoxRegisterOperator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbPermission.SafetyBoxOwner}
}

func validateBasic(sbId string, operator, address sdk.AccAddress) sdk.Error {
	if len(sbId) == 0 {
		return ErrSafetyBoxIdRequired(DefaultCodespace)
	}

	if operator.Empty() {
		return sdk.ErrInvalidAddress("Operator/SafetyBoxOwner is required")
	}

	if address.Empty() {
		return sdk.ErrInvalidAddress("Address is required")
	}

	return nil
}

func NewMsgSafetyBoxDeregisterIssuer(safetyBoxId string, operator, address sdk.AccAddress) MsgSafetyBoxDeregisterIssuer {
	return MsgSafetyBoxDeregisterIssuer{safetyBoxId, operator, address}
}

type MsgSafetyBoxDeregisterIssuer struct {
	SafetyBoxId string         `json:"safety_box_id"`
	Operator    sdk.AccAddress `json:"operator"`
	Address     sdk.AccAddress `json:"address"`
}

func (MsgSafetyBoxDeregisterIssuer) Route() string { return RouterKey }

func (MsgSafetyBoxDeregisterIssuer) Type() string { return MsgTypeSafetyBoxRevokeIssuerPermission }

func (msgSbPermission MsgSafetyBoxDeregisterIssuer) ValidateBasic() sdk.Error {
	return validateBasic(msgSbPermission.SafetyBoxId, msgSbPermission.Operator, msgSbPermission.Address)
}

func (msgSbPermission MsgSafetyBoxDeregisterIssuer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbPermission))
}

func (msgSbPermission MsgSafetyBoxDeregisterIssuer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbPermission.Operator}
}

func NewMsgSafetyBoxDeregisterReturner(safetyBoxId string, operator, address sdk.AccAddress) MsgSafetyBoxDeregisterReturner {
	return MsgSafetyBoxDeregisterReturner{safetyBoxId, operator, address}
}

type MsgSafetyBoxDeregisterReturner struct {
	SafetyBoxId string         `json:"safety_box_id"`
	Operator    sdk.AccAddress `json:"operator"`
	Address     sdk.AccAddress `json:"address"`
}

func (MsgSafetyBoxDeregisterReturner) Route() string { return RouterKey }

func (MsgSafetyBoxDeregisterReturner) Type() string { return MsgTypeSafetyBoxRevokeReturnerPermission }

func (msgSbPermission MsgSafetyBoxDeregisterReturner) ValidateBasic() sdk.Error {
	return validateBasic(msgSbPermission.SafetyBoxId, msgSbPermission.Operator, msgSbPermission.Address)
}

func (msgSbPermission MsgSafetyBoxDeregisterReturner) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbPermission))
}

func (msgSbPermission MsgSafetyBoxDeregisterReturner) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbPermission.Operator}
}

func NewMsgSafetyBoxDeregisterAllocator(safetyBoxId string, operator, address sdk.AccAddress) MsgSafetyBoxDeregisterAllocator {
	return MsgSafetyBoxDeregisterAllocator{safetyBoxId, operator, address}
}

type MsgSafetyBoxDeregisterAllocator struct {
	SafetyBoxId string         `json:"safety_box_id"`
	Operator    sdk.AccAddress `json:"operator"`
	Address     sdk.AccAddress `json:"address"`
}

func (MsgSafetyBoxDeregisterAllocator) Route() string { return RouterKey }

func (MsgSafetyBoxDeregisterAllocator) Type() string { return MsgTypeSafetyBoxRevokeAllocatorPermission }

func (msgSbPermission MsgSafetyBoxDeregisterAllocator) ValidateBasic() sdk.Error {
	return validateBasic(msgSbPermission.SafetyBoxId, msgSbPermission.Operator, msgSbPermission.Address)
}

func (msgSbPermission MsgSafetyBoxDeregisterAllocator) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbPermission))
}

func (msgSbPermission MsgSafetyBoxDeregisterAllocator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbPermission.Operator}
}

func NewMsgSafetyBoxDeregisterOperator(safetyBoxId string, operator, address sdk.AccAddress) MsgSafetyBoxDeregisterOperator {
	return MsgSafetyBoxDeregisterOperator{safetyBoxId, operator, address}
}

type MsgSafetyBoxDeregisterOperator struct {
	SafetyBoxId    string         `json:"safety_box_id"`
	SafetyBoxOwner sdk.AccAddress `json:"safety_box_owner"`
	Address        sdk.AccAddress `json:"address"`
}

func (MsgSafetyBoxDeregisterOperator) Route() string { return RouterKey }

func (MsgSafetyBoxDeregisterOperator) Type() string { return MsgTypeSafetyBoxRevokeOperatorPermission }

func (msgSbPermission MsgSafetyBoxDeregisterOperator) ValidateBasic() sdk.Error {
	return validateBasic(msgSbPermission.SafetyBoxId, msgSbPermission.SafetyBoxOwner, msgSbPermission.Address)
}

func (msgSbPermission MsgSafetyBoxDeregisterOperator) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSbPermission))
}

func (msgSbPermission MsgSafetyBoxDeregisterOperator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgSbPermission.SafetyBoxOwner}
}
