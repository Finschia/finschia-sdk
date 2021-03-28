package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/lbm-sdk/v2/errors"
	"github.com/line/lbm-sdk/v2/x/contract"
)

var _ sdk.Msg = (*MsgGrantPermission)(nil)

func NewMsgGrantPermission(from sdk.AccAddress, contractID string, to sdk.AccAddress, perm Permission) MsgGrantPermission {
	return MsgGrantPermission{
		From:       from,
		ContractID: contractID,
		To:         to,
		Permission: perm,
	}
}

type MsgGrantPermission struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	To         sdk.AccAddress `json:"to"`
	Permission Permission     `json:"permission"`
}

func (MsgGrantPermission) Route() string                    { return RouterKey }
func (MsgGrantPermission) Type() string                     { return "grant_perm" }
func (msg MsgGrantPermission) GetContractID() string        { return msg.ContractID }
func (msg MsgGrantPermission) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg MsgGrantPermission) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgGrantPermission) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}
	if msg.From.Empty() || msg.To.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "addresses cannot be empty")
	}

	if msg.From.Equals(msg.To) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "from, to address can not be the same")
	}

	return validateAction(msg.Permission.String(), MintAction, BurnAction, IssueAction, ModifyAction)
}

var _ sdk.Msg = (*MsgRevokePermission)(nil)

func NewMsgRevokePermission(from sdk.AccAddress, contractID string, perm Permission) MsgRevokePermission {
	return MsgRevokePermission{
		From:       from,
		ContractID: contractID,
		Permission: perm,
	}
}

type MsgRevokePermission struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	Permission Permission     `json:"permission"`
}

func (MsgRevokePermission) Route() string                    { return RouterKey }
func (MsgRevokePermission) Type() string                     { return "revoke_perm" }
func (msg MsgRevokePermission) GetContractID() string        { return msg.ContractID }
func (msg MsgRevokePermission) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg MsgRevokePermission) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgRevokePermission) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "addresses cannot be empty")
	}

	return validateAction(msg.Permission.String(), MintAction, BurnAction, IssueAction, ModifyAction)
}
func validateAction(action string, actions ...string) error {
	for _, a := range actions {
		if action == a {
			return nil
		}
	}
	return sdkerrors.Wrap(errors.ErrInvalidPermission, fmt.Sprintf("permission should be one of [%s]", strings.Join(actions, ",")))
}
