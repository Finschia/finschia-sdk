package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/contract"
)

var _ contract.Msg = (*MsgModify)(nil)

type MsgModify struct {
	Owner      sdk.AccAddress `json:"owner"`
	ContractID string         `json:"contract_id"`
	Changes    Changes        `json:"changes"`
}

func NewMsgModify(owner sdk.AccAddress, contractID string, changes Changes) MsgModify {
	return MsgModify{
		Owner:      owner,
		ContractID: contractID,
		Changes:    changes,
	}
}

func (msg MsgModify) Route() string         { return RouterKey }
func (msg MsgModify) Type() string          { return "modify_token" }
func (msg MsgModify) GetContractID() string { return msg.ContractID }
func (msg MsgModify) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgModify) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Owner} }

func (msg MsgModify) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}

	validator := NewChangesValidator()
	if err := validator.Validate(msg.Changes); err != nil {
		return err
	}

	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner address cannot be empty")
	}

	return nil
}
