package token

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/token/class"
)

var _ sdk.Msg = (*MsgSend)(nil)

// ValidateBasic implements Msg.
func (m MsgSend) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid from address: %s", m.From)
	}
	if err := sdk.ValidateAccAddress(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid to address: %s", m.To)
	}

	if err := validateAmount(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgSend) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.From)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgOperatorSend)(nil)

// ValidateBasic implements Msg.
func (m MsgOperatorSend) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.Operator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid operator address: %s", m.Operator)
	}
	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid from address: %s", m.From)
	}
	if err := sdk.ValidateAccAddress(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid to address: %s", m.To)
	}

	if err := validateAmount(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgOperatorSend) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Operator)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgTransferFrom)(nil)

// ValidateBasic implements Msg.
func (m MsgTransferFrom) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.Proxy); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid proxy address: %s", m.Proxy)
	}
	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid from address: %s", m.From)
	}
	if err := sdk.ValidateAccAddress(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid to address: %s", m.To)
	}

	if err := validateAmount(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgTransferFrom) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Proxy)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgAuthorizeOperator)(nil)

// ValidateBasic implements Msg.
func (m MsgAuthorizeOperator) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.Holder); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid holder address: %s", m.Holder)
	}
	if err := sdk.ValidateAccAddress(m.Operator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid operator address: %s", m.Operator)
	}

	return nil
}

// GetSigners implements Msg.
func (m MsgAuthorizeOperator) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Holder)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgRevokeOperator)(nil)

// ValidateBasic implements Msg.
func (m MsgRevokeOperator) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.Holder); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid holder address: %s", m.Holder)
	}
	if err := sdk.ValidateAccAddress(m.Operator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid operator address: %s", m.Operator)
	}

	return nil
}

// GetSigners implements Msg.
func (m MsgRevokeOperator) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Holder)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgApprove)(nil)

// ValidateBasic implements Msg.
func (m MsgApprove) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.Approver); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid approver address: %s", m.Approver)
	}
	if err := sdk.ValidateAccAddress(m.Proxy); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid proxy address: %s", m.Proxy)
	}

	return nil
}

// GetSigners implements Msg.
func (m MsgApprove) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Approver)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgIssue)(nil)

// ValidateBasic implements Msg.
func (m MsgIssue) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(m.Owner); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid owner address: %s", m.Owner)
	}

	if err := sdk.ValidateAccAddress(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid to address: %s", m.To)
	}

	if err := validateName(m.Name); err != nil {
		return err
	}

	if err := validateSymbol(m.Symbol); err != nil {
		return err
	}

	if err := validateImageURI(m.ImageUri); err != nil {
		return err
	}

	if err := validateMeta(m.Meta); err != nil {
		return err
	}

	if err := validateDecimals(m.Decimals); err != nil {
		return err
	}

	if err := validateAmount(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg.
func (m MsgIssue) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Owner)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgGrant)(nil)

// ValidateBasic implements Msg.
func (m MsgGrant) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.Granter); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid granter address: %s", m.Granter)
	}
	if err := sdk.ValidateAccAddress(m.Grantee); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid grantee address: %s", m.Grantee)
	}

	if err := validatePermission(m.Permission); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgGrant) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Granter)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgAbandon)(nil)

// ValidateBasic implements Msg.
func (m MsgAbandon) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.Grantee); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid grantee address: %s", m.Grantee)
	}

	if err := validatePermission(m.Permission); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgAbandon) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Grantee)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgGrantPermission)(nil)

// ValidateBasic implements Msg.
func (m MsgGrantPermission) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid granter address: %s", m.From)
	}
	if err := sdk.ValidateAccAddress(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid grantee address: %s", m.To)
	}

	if err := validatePermission(m.Permission); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgGrantPermission) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.From)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgRevokePermission)(nil)

// ValidateBasic implements Msg.
func (m MsgRevokePermission) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid from address: %s", m.From)
	}

	if err := validatePermission(m.Permission); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgRevokePermission) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.From)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgMint)(nil)

// ValidateBasic implements Msg.
func (m MsgMint) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}
	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid grantee address: %s", m.From)
	}

	if err := sdk.ValidateAccAddress(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid to address: %s", m.To)
	}

	if err := validateAmount(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgMint) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.From)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgBurn)(nil)

// ValidateBasic implements Msg.
func (m MsgBurn) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}
	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid from address: %s", m.From)
	}

	if err := validateAmount(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgBurn) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.From)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgOperatorBurn)(nil)

// ValidateBasic implements Msg.
func (m MsgOperatorBurn) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.Operator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid operator address: %s", m.Operator)
	}
	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid from address: %s", m.From)
	}

	if err := validateAmount(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgOperatorBurn) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Operator)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgBurnFrom)(nil)

// ValidateBasic implements Msg.
func (m MsgBurnFrom) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.Proxy); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid proxy address: %s", m.Proxy)
	}
	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid from address: %s", m.From)
	}

	if err := validateAmount(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgBurnFrom) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Proxy)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgModify)(nil)

// ValidateBasic implements Msg.
func (m MsgModify) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}
	if err := sdk.ValidateAccAddress(m.Owner); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid grantee address: %s", m.Owner)
	}

	checkedFields := map[string]bool{}
	for _, change := range m.Changes {
		if checkedFields[change.Field] {
			return sdkerrors.ErrInvalidRequest.Wrapf("Duplicated field: %s", change.Field)
		}
		checkedFields[change.Field] = true

		if err := validateChange(change); err != nil {
			return err
		}
	}
	if len(checkedFields) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrapf("No field provided")
	}

	return nil
}

// GetSigners implements Msg
func (m MsgModify) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Owner)
	return []sdk.AccAddress{signer}
}
