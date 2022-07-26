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
	return validateCoinsWithIDValidator(amount, ValidateTokenID)
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
	return class.ValidateID(id)
}

func ValidateClassID(id string) error {
	return validateID(id, reClassID)
}

// Deprecated: do not use (no successor).
func ValidateLegacyFTClassID(id string) error {
	return validateID(id, reLegacyFTClassID)
}

// Deprecated: do not use (no successor).
func ValidateLegacyNFTClassID(id string) error {
	return validateID(id, reLegacyNFTClassID)
}

func ValidateTokenID(id string) error {
	return validateID(id, reTokenID)
}

func ValidateFTID(id string) error {
	return validateID(id, reFTID)
}

func ValidateNFTID(id string) error {
	if err := ValidateTokenID(id); err != nil {
		return err
	}
	if err := ValidateFTID(id); err == nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid id: %s", id)
	}
	return nil
}

// Deprecated: do not use (no successor).
func ValidateLegacyNFTID(id string) error {
	return validateID(id, reLegacyNFTID)
}

func validateID(id string, reg *regexp.Regexp) error {
	if !reg.MatchString(id) {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid id: %s", id)
	}
	return nil
}

func validateName(name string) error {
	return validateStringSize(name, nameLengthLimit, "name")
}

func validateBaseImgURI(baseImgURI string) error {
	return validateStringSize(baseImgURI, baseImgURILengthLimit, "base_img_uri")
}

func validateMeta(meta string) error {
	return validateStringSize(meta, metaLengthLimit, "meta")
}

func validateStringSize(str string, limit int, name string) error {
	if length := utf8.RuneCountInString(str); length > limit {
		return sdkerrors.ErrInvalidRequest.Wrapf("%s cannot exceed %d in length: current %d", name, limit, length)
	}
	return nil
}

func validateDecimals(decimals int32) error {
	if decimals < 0 || decimals > 18 {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid decimals: %d", decimals)
	}
	return nil
}

func validateLegacyPermission(permission string) error {
	return ValidatePermission(Permission(LegacyPermissionFromString(permission)))
}

func ValidatePermission(permission Permission) error {
	if p := Permission_value[Permission_name[int32(permission)]]; p == 0 {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid permission: %s", permission)
	}
	return nil
}

func validateContractChange(change Attribute) error {
	validators := map[AttributeKey]func(string) error{
		AttributeKeyName:       validateName,
		AttributeKeyBaseImgURI: validateBaseImgURI,
		AttributeKeyMeta:       validateMeta,
	}

	return validateChange(change, validators)
}

func validateTokenClassChange(change Attribute) error {
	validators := map[AttributeKey]func(string) error{
		AttributeKeyName: validateName,
		AttributeKeyMeta: validateMeta,
	}

	return validateChange(change, validators)
}

func validateChange(change Attribute, validators map[AttributeKey]func(string) error) error {
	validator, ok := validators[AttributeKeyFromString(change.Key)]
	if !ok {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid field: %s", change.Key)
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
		return sdkerrors.ErrInvalidRequest.Wrap("token ids cannot be empty")
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
		return sdkerrors.ErrInvalidRequest.Wrap("token ids cannot be empty")
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

	return nil
}

// GetSigners implements Msg
func (m MsgApprove) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Approver)
	return []sdk.AccAddress{signer}
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

	return nil
}

// GetSigners implements Msg
func (m MsgDisapprove) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Approver)
	return []sdk.AccAddress{signer}
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
		return sdkerrors.ErrInvalidRequest.Wrapf("empty name")
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
		return sdkerrors.ErrInvalidRequest.Wrap("invalid issue of ft")
	}

	return nil
}

// GetSigners implements Msg
func (m MsgIssueFT) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Owner)
	return []sdk.AccAddress{signer}
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
		return sdkerrors.ErrInvalidRequest.Wrap("mint params cannot be empty")
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
		return sdkerrors.ErrInvalidRequest.Wrap("token ids cannot be empty")
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
		return sdkerrors.ErrInvalidRequest.Wrap("token ids cannot be empty")
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
			return sdkerrors.ErrInvalidRequest.Wrap("fungible token type without index")
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
			return sdkerrors.ErrInvalidRequest.Wrap("token index without type")
		}
	}
	if len(m.Changes) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty changes")
	}
	if len(m.Changes) > changesLimit {
		return sdkerrors.ErrInvalidRequest.Wrapf("the number of changes exceeds the limit: %d > %d", len(m.Changes), changesLimit)
	}
	seenKeys := map[string]bool{}
	for _, change := range m.Changes {
		if seenKeys[change.Field] {
			return sdkerrors.ErrInvalidRequest.Wrapf("duplicate keys: %s", change.Field)
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
		return sdkerrors.ErrInvalidRequest.Wrap("cannot attach token to itself")
	}

	return nil
}

// GetSigners implements Msg
func (m MsgAttach) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.From)
	return []sdk.AccAddress{signer}
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
		return sdkerrors.ErrInvalidRequest.Wrap("cannot attach token to itself")
	}

	return nil
}

// GetSigners implements Msg
func (m MsgAttachFrom) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Proxy)
	return []sdk.AccAddress{signer}
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
