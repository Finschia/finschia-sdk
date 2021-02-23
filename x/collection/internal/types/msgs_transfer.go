package types

import (
	"encoding/json"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/contract"
)

var _ contract.Msg = (*MsgTransferFT)(nil)
var _ contract.Msg = (*MsgTransferNFT)(nil)
var _ contract.Msg = (*MsgTransferFTFrom)(nil)
var _ contract.Msg = (*MsgTransferNFTFrom)(nil)

var _ json.Marshaler = (*MsgTransferFT)(nil)
var _ json.Unmarshaler = (*MsgTransferFT)(nil)
var _ json.Marshaler = (*MsgTransferNFT)(nil)
var _ json.Unmarshaler = (*MsgTransferNFT)(nil)
var _ json.Marshaler = (*MsgTransferFTFrom)(nil)
var _ json.Unmarshaler = (*MsgTransferFTFrom)(nil)
var _ json.Marshaler = (*MsgTransferNFTFrom)(nil)
var _ json.Unmarshaler = (*MsgTransferNFTFrom)(nil)

type MsgTransferFT struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	To         sdk.AccAddress `json:"to"`
	Amount     Coins          `json:"amount"`
}

func NewMsgTransferFT(from sdk.AccAddress, contractID string, to sdk.AccAddress, amount ...Coin) MsgTransferFT {
	return MsgTransferFT{
		From:       from,
		ContractID: contractID,
		To:         to,
		Amount:     amount,
	}
}

func (msg MsgTransferFT) MarshalJSON() ([]byte, error) {
	type msgAlias MsgTransferFT
	return json.Marshal(msgAlias(msg))
}

func (msg *MsgTransferFT) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgTransferFT
	return json.Unmarshal(data, msgAlias(msg))
}

func (MsgTransferFT) Route() string { return RouterKey }

func (MsgTransferFT) Type() string { return "transfer_ft" }

func (msg MsgTransferFT) GetContractID() string { return msg.ContractID }

func (msg MsgTransferFT) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}

	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "From cannot be empty")
	}

	if msg.To.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "To cannot be empty")
	}

	for _, tokenID := range msg.Amount {
		if err := ValidateDenom(tokenID.Denom); err != nil {
			return sdkerrors.Wrap(ErrInvalidTokenID, "invalid token id")
		}
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(ErrInvalidAmount, "invalid amount")
	}
	return nil
}

func (msg MsgTransferFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgTransferFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

type MsgTransferNFT struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	To         sdk.AccAddress `json:"to"`
	TokenIDs   []string       `json:"token_ids"`
}

func NewMsgTransferNFT(from sdk.AccAddress, contractID string, to sdk.AccAddress, tokenIDs ...string) MsgTransferNFT {
	return MsgTransferNFT{
		From:       from,
		ContractID: contractID,
		To:         to,
		TokenIDs:   tokenIDs,
	}
}

func (msg MsgTransferNFT) MarshalJSON() ([]byte, error) {
	type msgAlias MsgTransferNFT
	return json.Marshal(msgAlias(msg))
}

func (msg *MsgTransferNFT) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgTransferNFT
	return json.Unmarshal(data, msgAlias(msg))
}

func (MsgTransferNFT) Route() string { return RouterKey }

func (MsgTransferNFT) Type() string { return "transfer_nft" }

func (msg MsgTransferNFT) GetContractID() string { return msg.ContractID }

func (msg MsgTransferNFT) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}

	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "From cannot be empty")
	}

	if msg.To.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "To cannot be empty")
	}

	if len(msg.TokenIDs) == 0 {
		return sdkerrors.Wrap(ErrEmptyField, "token_ids cannot be empty")
	}
	for _, tokenID := range msg.TokenIDs {
		if err := ValidateTokenID(tokenID); err != nil {
			return sdkerrors.Wrap(ErrInvalidTokenID, err.Error())
		}
	}

	return nil
}

func (msg MsgTransferNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgTransferNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

type MsgTransferFTFrom struct {
	Proxy      sdk.AccAddress `json:"proxy"`
	ContractID string         `json:"contract_id"`
	From       sdk.AccAddress `json:"from"`
	To         sdk.AccAddress `json:"to"`
	Amount     Coins          `json:"amount"`
}

func NewMsgTransferFTFrom(proxy sdk.AccAddress, contractID string, from sdk.AccAddress, to sdk.AccAddress, amount ...Coin) MsgTransferFTFrom {
	return MsgTransferFTFrom{
		Proxy:      proxy,
		ContractID: contractID,
		From:       from,
		To:         to,
		Amount:     amount,
	}
}

func (msg MsgTransferFTFrom) MarshalJSON() ([]byte, error) {
	type msgAlias MsgTransferFTFrom
	return json.Marshal(msgAlias(msg))
}

func (msg *MsgTransferFTFrom) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgTransferFTFrom
	return json.Unmarshal(data, msgAlias(msg))
}

func (MsgTransferFTFrom) Route() string { return RouterKey }

func (MsgTransferFTFrom) Type() string { return "transfer_ft_from" }

func (msg MsgTransferFTFrom) GetContractID() string { return msg.ContractID }

func (msg MsgTransferFTFrom) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}
	if msg.Proxy.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Proxy cannot be empty")
	}
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "From cannot be empty")
	}
	if msg.To.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "To cannot be empty")
	}
	if msg.From.Equals(msg.Proxy) {
		return sdkerrors.Wrapf(ErrApproverProxySame, "Approver: %s", msg.From.String())
	}
	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(ErrInvalidAmount, "invalid amount")
	}
	return nil
}

func (msg MsgTransferFTFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgTransferFTFrom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proxy}
}

type MsgTransferNFTFrom struct {
	Proxy      sdk.AccAddress `json:"proxy"`
	ContractID string         `json:"contract_id"`
	From       sdk.AccAddress `json:"from"`
	To         sdk.AccAddress `json:"to"`
	TokenIDs   []string       `json:"token_ids"`
}

func NewMsgTransferNFTFrom(proxy sdk.AccAddress, contractID string, from sdk.AccAddress, to sdk.AccAddress, tokenIDs ...string) MsgTransferNFTFrom {
	return MsgTransferNFTFrom{
		Proxy:      proxy,
		ContractID: contractID,
		From:       from,
		To:         to,
		TokenIDs:   tokenIDs,
	}
}

func (msg MsgTransferNFTFrom) MarshalJSON() ([]byte, error) {
	type msgAlias MsgTransferNFTFrom
	return json.Marshal(msgAlias(msg))
}

func (msg *MsgTransferNFTFrom) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgTransferNFTFrom
	return json.Unmarshal(data, msgAlias(msg))
}

func (MsgTransferNFTFrom) Route() string { return RouterKey }

func (MsgTransferNFTFrom) Type() string { return "transfer_nft_from" }

func (msg MsgTransferNFTFrom) GetContractID() string { return msg.ContractID }

func (msg MsgTransferNFTFrom) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}
	if msg.Proxy.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Proxy cannot be empty")
	}
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "From cannot be empty")
	}
	if msg.To.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "To cannot be empty")
	}
	if msg.From.Equals(msg.Proxy) {
		return sdkerrors.Wrapf(ErrApproverProxySame, "Approver: %s", msg.From.String())
	}

	if len(msg.TokenIDs) == 0 {
		return sdkerrors.Wrap(ErrEmptyField, "token_ids cannot be empty")
	}
	for _, tokenID := range msg.TokenIDs {
		if err := ValidateTokenID(tokenID); err != nil {
			return sdkerrors.Wrap(ErrInvalidTokenID, err.Error())
		}
	}
	return nil
}

func (msg MsgTransferNFTFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgTransferNFTFrom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proxy}
}
