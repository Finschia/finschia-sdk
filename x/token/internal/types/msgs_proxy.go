package types

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/contract"
)

var _ contract.Msg = (*MsgApprove)(nil)

type MsgApprove struct {
	Approver   sdk.AccAddress `json:"approver"`
	ContractID string         `json:"contract_id"`
	Proxy      sdk.AccAddress `json:"proxy"`
}

func NewMsgApprove(approver sdk.AccAddress, contractID string, proxy sdk.AccAddress) MsgApprove {
	return MsgApprove{
		Approver:   approver,
		ContractID: contractID,
		Proxy:      proxy,
	}
}

func (msg MsgApprove) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return nil
	}
	if msg.Approver.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Approver cannot be empty")
	}
	if msg.Proxy.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Proxy cannot be empty")
	}
	if msg.Approver.Equals(msg.Proxy) {
		return sdkerrors.Wrapf(ErrApproverProxySame, "Approver: %s", msg.Approver.String())
	}
	return nil
}

func (MsgApprove) Route() string             { return RouterKey }
func (MsgApprove) Type() string              { return "approve_token" }
func (msg MsgApprove) GetContractID() string { return msg.ContractID }
func (msg MsgApprove) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Approver}
}
func (msg MsgApprove) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
