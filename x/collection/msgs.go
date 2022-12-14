package collection

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/token/class"
)

const (
	lengthClassID = 8

	nameLengthLimit       = 20
	baseImgURILengthLimit = 1000
	metaLengthLimit       = 1000
	changesLimit          = 100
)

var (
	patternAll  = fmt.Sprintf(`[0-9a-f]{%d}`, lengthClassID)
	patternZero = fmt.Sprintf(`0{%d}`, lengthClassID)

	patternClassID          = patternAll
	patternLegacyFTClassID  = fmt.Sprintf(`0[0-9a-f]{%d}`, lengthClassID-1)
	patternLegacyNFTClassID = fmt.Sprintf(`[1-9a-f][0-9a-f]{%d}`, lengthClassID-1)

	// regexps for class ids
	reClassID          = regexp.MustCompile(fmt.Sprintf(`^%s$`, patternClassID))
	reLegacyFTClassID  = regexp.MustCompile(fmt.Sprintf(`^%s$`, patternLegacyFTClassID))
	reLegacyNFTClassID = regexp.MustCompile(fmt.Sprintf(`^%s$`, patternLegacyNFTClassID))

	// regexps for token ids
	reTokenID     = regexp.MustCompile(fmt.Sprintf(`^%s%s$`, patternClassID, patternAll))
	reFTID        = regexp.MustCompile(fmt.Sprintf(`^%s%s$`, patternClassID, patternZero))
	reLegacyNFTID = regexp.MustCompile(fmt.Sprintf(`^%s%s$`, patternLegacyNFTClassID, patternAll))
)

func validateAmount(amount sdk.Int) error {
	if !amount.IsPositive() {
		return sdkerrors.ErrInvalidRequest.Wrapf("amount must be positive: %s", amount)
	}
	return nil
}

// deprecated
func validateCoins(amount []Coin) error {
	if err := validateCoinsWithIDValidator(amount, ValidateTokenID); err != nil {
		return ErrInvalidCoins.Wrap(err.Error())
	}

	return nil
}

// deprecated
func validateCoinsWithIDValidator(amount []Coin, validator func(string) error) error {
	for _, amt := range amount {
		if err := validator(amt.TokenId); err != nil {
			return err
		}
		if err := validateAmount(amt.Amount); err != nil {
			return err
		}
	}
	return nil
}

func NewFTID(classID string) string {
	return newTokenID(classID, sdk.ZeroUint())
}

func NewNFTID(classID string, number int) string {
	return newTokenID(classID, sdk.NewUint(uint64(number)))
}

func newTokenID(classID string, number sdk.Uint) string {
	numberFormat := "%0" + fmt.Sprintf("%d", lengthClassID) + "x"
	return classID + fmt.Sprintf(numberFormat, number.Uint64())
}

func SplitTokenID(tokenID string) (classID string) {
	return tokenID[:lengthClassID]
}

func ValidateContractID(id string) error {
	if err := class.ValidateID(id); err != nil {
		return ErrInvalidContractID.Wrap(id)
	}

	return nil
}

func ValidateClassID(id string) error {
	if err := validateID(id, reClassID); err != nil {
		return ErrInvalidClassID.Wrap(id)
	}

	return nil
}

// Deprecated: do not use (no successor).
func ValidateLegacyFTClassID(id string) error {
	return validateID(id, reLegacyFTClassID)
}

// Deprecated: do not use (no successor).
func ValidateLegacyNFTClassID(id string) error {
	if err := validateID(id, reLegacyNFTClassID); err != nil {
		return ErrInvalidClassID.Wrapf("%s not nft class", id)
	}

	return nil
}

func ValidateTokenID(id string) error {
	if err := validateID(id, reTokenID); err != nil {
		return ErrInvalidTokenID.Wrap(id)
	}

	return nil
}

func ValidateFTID(id string) error {
	if err := validateID(id, reFTID); err != nil {
		return ErrInvalidTokenID.Wrapf("%s not ft", id)
	}

	return nil
}

func ValidateNFTID(id string) error {
	if err := ValidateTokenID(id); err != nil {
		return err
	}
	if err := ValidateFTID(id); err == nil {
		return ErrInvalidTokenID.Wrapf("%s not nft", id)
	}
	return nil
}

// Deprecated: do not use (no successor).
func ValidateLegacyNFTID(id string) error {
	if err := validateID(id, reLegacyNFTID); err != nil {
		return ErrInvalidTokenID.Wrapf("%s not nft", id)
	}

	return nil
}

func validateID(id string, reg *regexp.Regexp) error {
	if !reg.MatchString(id) {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid id; %s", id)
	}
	return nil
}

func validateName(name string) error {
	if err := validateStringSize(name, nameLengthLimit); err != nil {
		return ErrInvalidName.Wrap(err.Error())
	}

	return nil
}

func validateBaseImgURI(baseImgURI string) error {
	if err := validateStringSize(baseImgURI, baseImgURILengthLimit); err != nil {
		return ErrInvalidBaseImgURI.Wrap(err.Error())
	}

	return nil
}

func validateMeta(meta string) error {
	if err := validateStringSize(meta, metaLengthLimit); err != nil {
		return ErrInvalidMeta.Wrap(err.Error())
	}

	return nil
}

func validateStringSize(str string, limit int) error {
	if length := utf8.RuneCountInString(str); length > limit {
		return sdkerrors.ErrInvalidRequest.Wrapf("%d exceeds its limit %d in length", length, limit)
	}
	return nil
}

func validateDecimals(decimals int32) error {
	if decimals < 0 || decimals > 18 {
		return ErrInvalidDecimals.Wrapf("must be >=0 and <=18, got; %d", decimals)
	}
	return nil
}

func validateLegacyPermission(permission string) error {
	return ValidatePermission(Permission(LegacyPermissionFromString(permission)))
}

func ValidatePermission(permission Permission) error {
	if p := Permission_value[Permission_name[int32(permission)]]; p == 0 {
		return ErrInvalidPermission.Wrap(permission.String())
	}
	return nil
}

func validateContractChange(change Attribute) error {
	validators := map[string]func(string) error{
		AttributeKeyName.String():       validateName,
		AttributeKeyBaseImgURI.String(): validateBaseImgURI,
		AttributeKeyMeta.String():       validateMeta,
	}

	return validateChange(change, validators)
}

func validateTokenClassChange(change Attribute) error {
	validators := map[string]func(string) error{
		AttributeKeyName.String(): validateName,
		AttributeKeyMeta.String(): validateMeta,
	}

	return validateChange(change, validators)
}

func validateChange(change Attribute, validators map[string]func(string) error) error {
	validator, ok := validators[change.Key]
	if !ok {
		return ErrInvalidChanges.Wrapf("invalid key: %s", change.Key)
	}
	return validator(change.Value)
}

var _ sdk.Msg = (*MsgTransferFT)(nil)

// ValidateBasic implements Msg.
func (m MsgTransferFT) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}
	if _, err := sdk.AccAddressFromBech32(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", m.To)
	}

	if err := validateCoins(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgTransferFT) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.From)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgTransferFT) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgTransferFT) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgTransferFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = (*MsgTransferFTFrom)(nil)

// ValidateBasic implements Msg.
func (m MsgTransferFTFrom) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.Proxy); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid proxy address: %s", m.Proxy)
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}
	if _, err := sdk.AccAddressFromBech32(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", m.To)
	}

	if err := validateCoins(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgTransferFTFrom) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Proxy)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgTransferFTFrom) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgTransferFTFrom) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgTransferFTFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = (*MsgTransferNFT)(nil)

// ValidateBasic implements Msg.
func (m MsgTransferNFT) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}
	if _, err := sdk.AccAddressFromBech32(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", m.To)
	}

	if len(m.TokenIds) == 0 {
		return ErrEmptyTokenIDs
	}
	for _, id := range m.TokenIds {
		if err := ValidateTokenID(id); err != nil {
			return err
		}
	}

	return nil
}

// GetSigners implements Msg
func (m MsgTransferNFT) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.From)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgTransferNFT) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgTransferNFT) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgTransferNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = (*MsgTransferNFTFrom)(nil)

// ValidateBasic implements Msg.
func (m MsgTransferNFTFrom) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.Proxy); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid proxy address: %s", m.Proxy)
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}
	if _, err := sdk.AccAddressFromBech32(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", m.To)
	}

	if len(m.TokenIds) == 0 {
		return ErrEmptyTokenIDs
	}
	for _, id := range m.TokenIds {
		if err := ValidateTokenID(id); err != nil {
			return err
		}
	}

	return nil
}

// GetSigners implements Msg
func (m MsgTransferNFTFrom) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Proxy)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgTransferNFTFrom) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgTransferNFTFrom) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgTransferNFTFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = (*MsgApprove)(nil)

// ValidateBasic implements Msg.
func (m MsgApprove) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.Approver); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid approver address: %s", m.Approver)
	}
	if _, err := sdk.AccAddressFromBech32(m.Proxy); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid proxy address: %s", m.Proxy)
	}

	if m.Proxy == m.Approver {
		return ErrOperatorIsHolder
	}

	return nil
}

// GetSigners implements Msg
func (m MsgApprove) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Approver)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgApprove) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgApprove) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgApprove) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = (*MsgDisapprove)(nil)

// ValidateBasic implements Msg.
func (m MsgDisapprove) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.Approver); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid approver address: %s", m.Approver)
	}
	if _, err := sdk.AccAddressFromBech32(m.Proxy); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid proxy address: %s", m.Proxy)
	}

	if m.Proxy == m.Approver {
		return ErrOperatorIsHolder
	}

	return nil
}

// GetSigners implements Msg
func (m MsgDisapprove) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Approver)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgDisapprove) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgDisapprove) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgDisapprove) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = (*MsgCreateContract)(nil)

// ValidateBasic implements Msg.
func (m MsgCreateContract) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", m.Owner)
	}

	if err := validateName(m.Name); err != nil {
		return err
	}

	if err := validateBaseImgURI(m.BaseImgUri); err != nil {
		return err
	}

	if err := validateMeta(m.Meta); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgCreateContract) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Owner)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgCreateContract) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgCreateContract) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgCreateContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = (*MsgIssueFT)(nil)

// ValidateBasic implements Msg.
func (m MsgIssueFT) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", m.Owner)
	}

	if len(m.Name) == 0 {
		return ErrInvalidName.Wrap("empty")
	}
	if err := validateName(m.Name); err != nil {
		return err
	}

	if err := validateMeta(m.Meta); err != nil {
		return err
	}

	if err := validateDecimals(m.Decimals); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", m.To)
	}

	// daphne compat.
	if m.Amount.Equal(sdk.OneInt()) && m.Decimals == 0 && !m.Mintable {
		return ErrBadUseCase.Wrap("condition (amount == 0 & decimals == 0 & mintable == false) is invalid")
	}

	return nil
}

// GetSigners implements Msg
func (m MsgIssueFT) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Owner)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgIssueFT) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgIssueFT) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgIssueFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = (*MsgIssueNFT)(nil)

// ValidateBasic implements Msg.
func (m MsgIssueNFT) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if err := validateName(m.Name); err != nil {
		return err
	}

	if err := validateMeta(m.Meta); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", m.Owner)
	}

	return nil
}

// GetSigners implements Msg
func (m MsgIssueNFT) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Owner)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgIssueNFT) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgIssueNFT) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgIssueNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = (*MsgMintFT)(nil)

// ValidateBasic implements Msg.
func (m MsgMintFT) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}
	if _, err := sdk.AccAddressFromBech32(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", m.To)
	}

	if err := validateCoins(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgMintFT) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.From)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgMintFT) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgMintFT) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgMintFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = (*MsgMintNFT)(nil)

// ValidateBasic implements Msg.
func (m MsgMintNFT) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}
	if _, err := sdk.AccAddressFromBech32(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", m.To)
	}

	if len(m.Params) == 0 {
		return ErrInvalidMintNFTParams.Wrap("empty")
	}
	for _, param := range m.Params {
		classID := param.TokenType
		if err := ValidateLegacyNFTClassID(classID); err != nil {
			return err
		}

		if err := validateName(param.Name); err != nil {
			return err
		}

		if err := validateMeta(param.Meta); err != nil {
			return err
		}
	}

	return nil
}

// GetSigners implements Msg
func (m MsgMintNFT) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.From)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgMintNFT) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgMintNFT) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgMintNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = (*MsgBurnFT)(nil)

// ValidateBasic implements Msg.
func (m MsgBurnFT) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}

	if err := validateCoins(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgBurnFT) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.From)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgBurnFT) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgBurnFT) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgBurnFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = (*MsgBurnFTFrom)(nil)

// ValidateBasic implements Msg.
func (m MsgBurnFTFrom) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.Proxy); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid proxy address: %s", m.Proxy)
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}

	if err := validateCoins(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgBurnFTFrom) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Proxy)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgBurnFTFrom) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgBurnFTFrom) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgBurnFTFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = (*MsgBurnNFT)(nil)

// ValidateBasic implements Msg.
func (m MsgBurnNFT) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}

	if len(m.TokenIds) == 0 {
		return ErrEmptyTokenIDs
	}
	for _, id := range m.TokenIds {
		if err := ValidateLegacyNFTID(id); err != nil {
			return err
		}
	}

	return nil
}

// GetSigners implements Msg
func (m MsgBurnNFT) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.From)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgBurnNFT) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgBurnNFT) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgBurnNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = (*MsgBurnNFTFrom)(nil)

// ValidateBasic implements Msg.
func (m MsgBurnNFTFrom) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.Proxy); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid proxy address: %s", m.Proxy)
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}

	if len(m.TokenIds) == 0 {
		return ErrEmptyTokenIDs
	}
	for _, id := range m.TokenIds {
		if err := ValidateLegacyNFTID(id); err != nil {
			return err
		}
	}

	return nil
}

// GetSigners implements Msg
func (m MsgBurnNFTFrom) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Proxy)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgBurnNFTFrom) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgBurnNFTFrom) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgBurnNFTFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = (*MsgModify)(nil)

// ValidateBasic implements Msg.
func (m MsgModify) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", m.Owner)
	}

	if len(m.TokenType) != 0 {
		classID := m.TokenType
		if err := ValidateClassID(classID); err != nil {
			return err
		}
		if err := ValidateLegacyFTClassID(classID); err == nil && len(m.TokenIndex) == 0 {
			return ErrInvalidModificationTarget.Wrap("fungible token type without index")
		}
	}

	if len(m.TokenIndex) != 0 {
		tokenID := m.TokenType + m.TokenIndex
		if err := ValidateTokenID(tokenID); err != nil {
			return err
		}
	}

	validator := validateTokenClassChange
	if len(m.TokenType) == 0 {
		if len(m.TokenIndex) == 0 {
			validator = validateContractChange
		} else {
			return ErrInvalidModificationTarget.Wrap("token index without type")
		}
	}
	if len(m.Changes) == 0 {
		return ErrInvalidChanges.Wrap("empty")
	}
	if len(m.Changes) > changesLimit {
		return ErrInvalidChanges.Wrapf("number of changes exceeds its limit: %d > %d", len(m.Changes), changesLimit)
	}
	seenKeys := map[string]bool{}
	for _, change := range m.Changes {
		if seenKeys[change.Field] {
			return ErrInvalidChanges.Wrapf("duplicate keys: %s", change.Field)
		}
		seenKeys[change.Field] = true

		attribute := Attribute{
			Key:   change.Field,
			Value: change.Value,
		}
		if err := validator(attribute); err != nil {
			return err
		}
	}

	return nil
}

// GetSigners implements Msg
func (m MsgModify) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Owner)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgModify) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgModify) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgModify) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = (*MsgGrantPermission)(nil)

// ValidateBasic implements Msg.
func (m MsgGrantPermission) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}
	if _, err := sdk.AccAddressFromBech32(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", m.To)
	}

	if err := validateLegacyPermission(m.Permission); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgGrantPermission) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.From)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgGrantPermission) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgGrantPermission) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgGrantPermission) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = (*MsgRevokePermission)(nil)

// ValidateBasic implements Msg.
func (m MsgRevokePermission) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}

	if err := validateLegacyPermission(m.Permission); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgRevokePermission) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.From)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgRevokePermission) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgRevokePermission) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgRevokePermission) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = (*MsgAttach)(nil)

// ValidateBasic implements Msg.
func (m MsgAttach) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}

	if err := ValidateTokenID(m.TokenId); err != nil {
		return err
	}
	if err := ValidateTokenID(m.ToTokenId); err != nil {
		return err
	}

	if m.TokenId == m.ToTokenId {
		return ErrInvalidComposition.Wrap("target and subject should be different")
	}

	return nil
}

// GetSigners implements Msg
func (m MsgAttach) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.From)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgAttach) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgAttach) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgAttach) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = (*MsgDetach)(nil)

// ValidateBasic implements Msg.
func (m MsgDetach) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}

	if err := ValidateTokenID(m.TokenId); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgDetach) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.From)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgDetach) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgDetach) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgDetach) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = (*MsgAttachFrom)(nil)

// ValidateBasic implements Msg.
func (m MsgAttachFrom) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.Proxy); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid proxy address: %s", m.Proxy)
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}

	if err := ValidateTokenID(m.TokenId); err != nil {
		return err
	}
	if err := ValidateTokenID(m.ToTokenId); err != nil {
		return err
	}

	if m.TokenId == m.ToTokenId {
		return ErrInvalidComposition.Wrap("target and subject should be different")
	}

	return nil
}

// GetSigners implements Msg
func (m MsgAttachFrom) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Proxy)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgAttachFrom) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgAttachFrom) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgAttachFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = (*MsgDetachFrom)(nil)

// ValidateBasic implements Msg.
func (m MsgDetachFrom) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.Proxy); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid proxy address: %s", m.Proxy)
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}

	if err := ValidateTokenID(m.TokenId); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgDetachFrom) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Proxy)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgDetachFrom) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgDetachFrom) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgDetachFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
