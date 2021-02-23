package types

import (
	"encoding/json"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/contract"
)

var _ contract.Msg = (*MsgTransfer)(nil)

type MsgTransfer struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	To         sdk.AccAddress `json:"to"`
	Amount     sdk.Int        `json:"amount"`
}

func NewMsgTransfer(from sdk.AccAddress, to sdk.AccAddress, contractID string, amount sdk.Int) MsgTransfer {
	return MsgTransfer{From: from, To: to, ContractID: contractID, Amount: amount}
}

func (msg MsgTransfer) Route() string { return RouterKey }

func (msg MsgTransfer) Type() string { return "transfer_ft" }

func (msg MsgTransfer) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}

	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "from cannot be empty")
	}

	if msg.To.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "to cannot be empty")
	}

	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "send amount must be positive")
	}
	return nil
}

func (msg MsgTransfer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgTransfer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

func (msg MsgTransfer) GetContractID() string {
	return msg.ContractID
}

var _ contract.Msg = (*MsgTransferFrom)(nil)

type MsgTransferFrom struct {
	Proxy      sdk.AccAddress `json:"proxy"`
	ContractID string         `json:"contract_id"`
	From       sdk.AccAddress `json:"from"`
	To         sdk.AccAddress `json:"to"`
	Amount     sdk.Int        `json:"amount"`
}

func NewMsgTransferFrom(proxy sdk.AccAddress, contractID string, from sdk.AccAddress, to sdk.AccAddress, amount sdk.Int) MsgTransferFrom {
	return MsgTransferFrom{
		Proxy:      proxy,
		ContractID: contractID,
		From:       from,
		To:         to,
		Amount:     amount,
	}
}
func (msg MsgTransferFrom) MarshalJSON() ([]byte, error) {
	type msgAlias MsgTransferFrom
	return json.Marshal(msgAlias(msg))
}

func (msg *MsgTransferFrom) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgTransferFrom
	return json.Unmarshal(data, msgAlias(msg))
}

func (MsgTransferFrom) Route() string { return RouterKey }

func (MsgTransferFrom) Type() string { return "transfer_from" }

func (msg MsgTransferFrom) GetContractID() string { return msg.ContractID }

func (msg MsgTransferFrom) ValidateBasic() error {
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
	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(ErrInvalidAmount, "invalid amount")
	}
	return nil
}

func (msg MsgTransferFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgTransferFrom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proxy}
}
