package token

import (
	"unicode/utf8"
	"regexp"

	"github.com/line/lbm-sdk/codec/legacy"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

const (
	TypeMsgTransfer = "transfer"
	TypeMsgTransferFrom = "transfer-from"
	TypeMsgApprove = "approve"
	TypeMsgIssue = "issue"
	TypeMsgGrant = "grant"
	TypeMsgRevoke = "revoke"
	TypeMsgMint = "mint"
	TypeMsgBurn = "burn"
	TypeMsgBurnFrom = "burn-from"
	TypeMsgModify = "modify"

	ActionMint = "mint"
	ActionBurn = "burn"
	ActionModify = "modify"

	AttributeKeyName = "name"
	AttributeKeyImageUri = "image_uri"
	AttributeKeyMeta = "meta"

	maxName = 20
	maxImageUri = 1000
	maxMeta = 1000
)

func validateAmount(amount sdk.Int) error {
	if !amount.IsPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Amount must be positive: %s", amount)
	}
	return nil
}

var _ sdk.Msg = (*MsgTransfer)(nil)

// Route implements Msg.
func (m MsgTransfer) Route() string { return RouterKey }

// Type implements Msg.
func (m MsgTransfer) Type() string { return sdk.MsgTypeURL(&m) }

// ValidateBasic implements Msg.
func (m MsgTransfer) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid from address: %s", m.From)
	}

	if err := sdk.ValidateAccAddress(m.To); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid to address: %s", m.To)
	}

	if err := validateAmount(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSignBytes implements Msg.
func (m MsgTransfer) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// GetSigners implements Msg
func (m MsgTransfer) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.From)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgTransferFrom)(nil)

// Route implements Msg.
func (m MsgTransferFrom) Route() string { return RouterKey }

// Type implements Msg.
func (m MsgTransferFrom) Type() string { return sdk.MsgTypeURL(&m) }

// ValidateBasic implements Msg.
func (m MsgTransferFrom) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(m.Proxy); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid proxy address: %s", m.Proxy)
	}

	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid from address: %s", m.From)
	}

	if err := sdk.ValidateAccAddress(m.To); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid to address: %s", m.To)
	}

	if err := validateAmount(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSignBytes implements Msg.
func (m MsgTransferFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// GetSigners implements Msg
func (m MsgTransferFrom) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Proxy)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgApprove)(nil)

// Route implements Msg.
func (m MsgApprove) Route() string { return RouterKey }

// Type implements Msg.
func (m MsgApprove) Type() string { return sdk.MsgTypeURL(&m) }

// ValidateBasic implements Msg.
func (m MsgApprove) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(m.Approver); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid approver address: %s", m.Approver)
	}

	if err := sdk.ValidateAccAddress(m.Proxy); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid proxy address: %s", m.Proxy)
	}

	return nil
}

// GetSignBytes implements Msg.
func (m MsgApprove) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// GetSigners implements Msg.
func (m MsgApprove) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Approver)
	return []sdk.AccAddress{signer}
}

func stringInSize(str string, size int) bool {
	return utf8.RuneCountInString(str) <= size
}

func validateName(name string) error {
	if len(name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Name cannot be empty")
	} else if !stringInSize(name, maxName) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Name cannot be longer than %d", maxName)
	}
	return nil
}

func validateSymbol(symbol string) error {
	reg := regexp.MustCompile(`[A-Z][A-Z0-9]{1,4}`)
	if !reg.MatchString(symbol) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Invalid symbol: %s, valid expression is: %s", symbol, reg.String())
	}
	return nil
}

func validateImageUri(uri string) error {
	if !stringInSize(uri, maxImageUri) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "ImageUri cannot be longer than %d", maxImageUri)
	}
	return nil
}

func validateMeta(meta string) error {
	if !stringInSize(meta, maxMeta) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Meta cannot be longer than %d", maxMeta)
	}
	return nil
}

func validateDecimals(decimals int32) error {
	if decimals < 0 || decimals > 18 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Invalid decimals: %d", decimals)
	}
	return nil
}

var _ sdk.Msg = (*MsgIssue)(nil)

// Route implements Msg.
func (m MsgIssue) Route() string { return RouterKey }

// Type implements Msg.
func (m MsgIssue) Type() string { return sdk.MsgTypeURL(&m) }

// ValidateBasic implements Msg.
func (m MsgIssue) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(m.Owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid owner address: %s", m.Owner)
	}

	if err := sdk.ValidateAccAddress(m.To); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid to address: %s", m.To)
	}

	if err := validateName(m.Name); err != nil {
		return err
	}

	if err := validateSymbol(m.Symbol); err != nil {
		return err
	}

	if err := validateImageUri(m.ImageUri); err != nil {
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

// GetSignBytes implements Msg.
func (m MsgIssue) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// GetSigners implements Msg.
func (m MsgIssue) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Owner)
	return []sdk.AccAddress{signer}
}

func validateAction(action string) error {
	actions := []string{
		ActionMint,
		ActionBurn,
		ActionModify,
	}
	for _, a := range actions {
		if action == a {
			return nil
		}
	}
	return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Invalid action: %s", action)
}

var _ sdk.Msg = (*MsgGrant)(nil)

// Route implements Msg.
func (m MsgGrant) Route() string { return RouterKey }

// Type implements Msg.
func (m MsgGrant) Type() string { return sdk.MsgTypeURL(&m) }

// ValidateBasic implements Msg.
func (m MsgGrant) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(m.Granter); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid granter address: %s", m.Granter)
	}

	if err := sdk.ValidateAccAddress(m.Grantee); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid grantee address: %s", m.Grantee)
	}

	if err := validateAction(m.Action); err != nil {
		return err
	}

	return nil
}

// GetSignBytes implements Msg.
func (m MsgGrant) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// GetSigners implements Msg
func (m MsgGrant) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Granter)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgRevoke)(nil)

// Route implements Msg.
func (m MsgRevoke) Route() string { return RouterKey }

// Type implements Msg.
func (m MsgRevoke) Type() string { return sdk.MsgTypeURL(&m) }

// ValidateBasic implements Msg.
func (m MsgRevoke) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(m.Grantee); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid grantee address: %s", m.Grantee)
	}

	if err := validateAction(m.Action); err != nil {
		return err
	}

	return nil
}

// GetSignBytes implements Msg.
func (m MsgRevoke) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// GetSigners implements Msg
func (m MsgRevoke) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Grantee)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgMint)(nil)

// Route implements Msg.
func (m MsgMint) Route() string { return RouterKey }

// Type implements Msg.
func (m MsgMint) Type() string { return sdk.MsgTypeURL(&m) }

// ValidateBasic implements Msg.
func (m MsgMint) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(m.Grantee); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid grantee address: %s", m.Grantee)
	}

	if err := sdk.ValidateAccAddress(m.To); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid to address: %s", m.To)
	}

	if err := validateAmount(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSignBytes implements Msg.
func (m MsgMint) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// GetSigners implements Msg
func (m MsgMint) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Grantee)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgBurn)(nil)

// Route implements Msg.
func (m MsgBurn) Route() string { return RouterKey }

// Type implements Msg.
func (m MsgBurn) Type() string { return sdk.MsgTypeURL(&m) }

// ValidateBasic implements Msg.
func (m MsgBurn) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid from address: %s", m.From)
	}

	if err := validateAmount(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSignBytes implements Msg.
func (m MsgBurn) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// GetSigners implements Msg
func (m MsgBurn) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.From)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgBurnFrom)(nil)

// Route implements Msg.
func (m MsgBurnFrom) Route() string { return RouterKey }

// Type implements Msg.
func (m MsgBurnFrom) Type() string { return sdk.MsgTypeURL(&m) }

// ValidateBasic implements Msg.
func (m MsgBurnFrom) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(m.Grantee); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid grantee address: %s", m.Grantee)
	}

	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid from address: %s", m.From)
	}

	if err := validateAmount(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSignBytes implements Msg.
func (m MsgBurnFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// GetSigners implements Msg
func (m MsgBurnFrom) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Grantee)
	return []sdk.AccAddress{signer}
}

func validateChange(change Pair) error {
	validators := map[string]func(string)error{
		AttributeKeyName: validateName,
		AttributeKeyImageUri: validateImageUri,
		AttributeKeyMeta: validateMeta,
	}
	validator, ok := validators[change.Key]
	if !ok {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Invalid field: %s", change.Key)
	}
	return validator(change.Value)
}
	
var _ sdk.Msg = (*MsgModify)(nil)

// Route implements Msg.
func (m MsgModify) Route() string { return RouterKey }

// Type implements Msg.
func (m MsgModify) Type() string { return sdk.MsgTypeURL(&m) }

// ValidateBasic implements Msg.
func (m MsgModify) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(m.Grantee); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid grantee address: %s", m.Grantee)
	}

	checkedFields := map[string]bool{}
	for _, change := range m.Changes {
		if checkedFields[change.Key] {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Duplicated field: %s", change.Key)
		}
		checkedFields[change.Key] = true

		if err := validateChange(change); err != nil {
			return err
		}
	}
	if len(checkedFields) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "No field provided")
	}

	return nil
}

// GetSignBytes implements Msg.
func (m MsgModify) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// GetSigners implements Msg
func (m MsgModify) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Grantee)
	return []sdk.AccAddress{signer}
}

