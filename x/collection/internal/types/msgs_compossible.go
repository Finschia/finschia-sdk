package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/contract"
)

var _ contract.Msg = (*MsgAttach)(nil)
var _ contract.Msg = (*MsgDetach)(nil)
var _ contract.Msg = (*MsgAttachFrom)(nil)
var _ contract.Msg = (*MsgDetachFrom)(nil)

type MsgAttach struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	ToTokenID  string         `json:"to_token_id"`
	TokenID    string         `json:"token_id"`
}

func NewMsgAttach(from sdk.AccAddress, contractID string, toTokenID string, tokenID string) MsgAttach {
	return MsgAttach{
		From:       from,
		ContractID: contractID,
		ToTokenID:  toTokenID,
		TokenID:    tokenID,
	}
}

func (msg MsgAttach) MarshalJSON() ([]byte, error) {
	type msgAlias MsgAttach
	return json.Marshal(msgAlias(msg))
}

func (msg *MsgAttach) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgAttach
	return json.Unmarshal(data, msgAlias(msg))
}

func (MsgAttach) Route() string { return RouterKey }

func (MsgAttach) Type() string { return "attach" }

func (msg MsgAttach) GetContractID() string { return msg.ContractID }

func (msg MsgAttach) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}

	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "From cannot be empty")
	}

	if err := ValidateTokenID(msg.ToTokenID); err != nil {
		return sdkerrors.Wrap(ErrInvalidTokenID, msg.ToTokenID)
	}

	if err := ValidateTokenID(msg.TokenID); err != nil {
		return sdkerrors.Wrap(ErrInvalidTokenID, msg.TokenID)
	}

	if msg.ToTokenID == msg.TokenID {
		return sdkerrors.Wrapf(ErrCannotAttachToItself, "TokenID: %s", msg.TokenID)
	}

	return nil
}

func (msg MsgAttach) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgAttach) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

type MsgDetach struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	TokenID    string         `json:"token_id"`
}

func NewMsgDetach(from sdk.AccAddress, contractID string, tokenID string) MsgDetach {
	return MsgDetach{
		From:       from,
		ContractID: contractID,
		TokenID:    tokenID,
	}
}

func (msg MsgDetach) MarshalJSON() ([]byte, error) {
	type msgAlias MsgDetach
	return json.Marshal(msgAlias(msg))
}

func (msg *MsgDetach) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgDetach
	return json.Unmarshal(data, msgAlias(msg))
}

func (MsgDetach) Route() string { return RouterKey }

func (MsgDetach) Type() string { return "detach" }

func (msg MsgDetach) GetContractID() string { return msg.ContractID }

func (msg MsgDetach) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}

	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "From cannot be empty")
	}

	if err := ValidateTokenID(msg.TokenID); err != nil {
		return sdkerrors.Wrap(ErrInvalidTokenID, msg.TokenID)
	}

	return nil
}

func (msg MsgDetach) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgDetach) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

type MsgAttachFrom struct {
	Proxy      sdk.AccAddress `json:"proxy"`
	ContractID string         `json:"contract_id"`
	From       sdk.AccAddress `json:"from"`
	ToTokenID  string         `json:"to_token_id"`
	TokenID    string         `json:"token_id"`
}

func NewMsgAttachFrom(proxy sdk.AccAddress, contractID string, from sdk.AccAddress, toTokenID string, tokenID string) MsgAttachFrom {
	return MsgAttachFrom{
		Proxy:      proxy,
		ContractID: contractID,
		From:       from,
		ToTokenID:  toTokenID,
		TokenID:    tokenID,
	}
}

func (msg MsgAttachFrom) MarshalJSON() ([]byte, error) {
	type msgAlias MsgAttachFrom
	return json.Marshal(msgAlias(msg))
}

func (msg *MsgAttachFrom) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgAttachFrom
	return json.Unmarshal(data, msgAlias(msg))
}

func (MsgAttachFrom) Route() string { return RouterKey }

func (MsgAttachFrom) Type() string { return "attach_from" }

func (msg MsgAttachFrom) GetContractID() string { return msg.ContractID }

func (msg MsgAttachFrom) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}
	if msg.Proxy.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Proxy cannot be empty")
	}
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "From cannot be empty")
	}
	if err := ValidateTokenID(msg.ToTokenID); err != nil {
		return sdkerrors.Wrap(ErrInvalidTokenID, msg.ToTokenID)
	}
	if err := ValidateTokenID(msg.TokenID); err != nil {
		return sdkerrors.Wrap(ErrInvalidTokenID, msg.TokenID)
	}

	if msg.ToTokenID == msg.TokenID {
		return sdkerrors.Wrapf(ErrCannotAttachToItself, "TokenID: %s", msg.TokenID)
	}

	return nil
}

func (msg MsgAttachFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgAttachFrom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proxy}
}

type MsgDetachFrom struct {
	Proxy      sdk.AccAddress `json:"proxy"`
	ContractID string         `json:"contract_id"`
	From       sdk.AccAddress `json:"from"`
	TokenID    string         `json:"token_id"`
}

func NewMsgDetachFrom(proxy sdk.AccAddress, contractID string, from sdk.AccAddress, tokenID string) MsgDetachFrom {
	return MsgDetachFrom{
		Proxy:      proxy,
		ContractID: contractID,
		From:       from,
		TokenID:    tokenID,
	}
}

func (msg MsgDetachFrom) MarshalJSON() ([]byte, error) {
	type msgAlias MsgDetachFrom
	return json.Marshal(msgAlias(msg))
}

func (msg *MsgDetachFrom) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgDetachFrom
	return json.Unmarshal(data, msgAlias(msg))
}

func (MsgDetachFrom) Route() string { return RouterKey }

func (MsgDetachFrom) Type() string { return "detach_from" }

func (msg MsgDetachFrom) GetContractID() string { return msg.ContractID }

func (msg MsgDetachFrom) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}
	if msg.Proxy.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Proxy cannot be empty")
	}
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "From cannot be empty")
	}
	if err := ValidateTokenID(msg.TokenID); err != nil {
		return sdkerrors.Wrap(ErrInvalidTokenID, msg.TokenID)
	}

	return nil
}

func (msg MsgDetachFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgDetachFrom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proxy}
}
