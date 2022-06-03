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
	//nolint:gosec
	suffixTokenID = `[0-9a-f]{8}`
	suffixClassID = `0{8}`

	//nolint:gosec
	prefixTokenID = `[0-9a-f]{8}`
	prefixFTID    = `0[0-9a-f]{7}`
	prefixNFTID   = `[1-9a-f][0-9a-f]{7}`

	nameLengthLimit       = 20
	baseImgURILengthLimit = 1000
	metaLengthLimit       = 1000
	changesLimit          = 100
)

var (
	// regexps for class ids
	reClassID    = regexp.MustCompile(fmt.Sprintf(`^%s%s$`, prefixTokenID, suffixClassID))
	reFTClassID  = regexp.MustCompile(fmt.Sprintf(`^%s%s$`, prefixFTID, suffixClassID))
	reNFTClassID = regexp.MustCompile(fmt.Sprintf(`^%s%s$`, prefixNFTID, suffixClassID))

	// regexps for token ids
	reTokenID = regexp.MustCompile(fmt.Sprintf(`^%s%s$`, prefixTokenID, suffixTokenID))
	reFTID    = regexp.MustCompile(fmt.Sprintf(`^%s%s$`, prefixFTID, suffixTokenID))
	reNFTID   = regexp.MustCompile(fmt.Sprintf(`^%s%s$`, prefixNFTID, suffixTokenID))
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

func ValidateClassID(id string) error {
	return validateID(id, reClassID)
}

func ValidateFTClassID(id string) error {
	return validateID(id, reFTClassID)
}

func ValidateNFTClassID(id string) error {
	return validateID(id, reNFTClassID)
}

func ValidateTokenID(id string) error {
	return validateID(id, reTokenID)
}

func ValidateFTID(id string) error {
	return validateID(id, reFTID)
}

func ValidateNFTID(id string) error {
	return validateID(id, reNFTID)
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

func validatePermission(permission string) error {
	if value := Permission_value[permission]; value == 0 {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid permission: %s", permission)
	}
	return nil
}

func validateContractChange(change Pair) error {
	validators := map[AttributeKey]func(string) error{
		AttributeKey_Name:       validateName,
		AttributeKey_BaseImgURI: validateBaseImgURI,
		AttributeKey_Meta:       validateMeta,
	}

	validator, ok := validators[AttributeKey(AttributeKey_value[change.Field])]
	if !ok {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid field: %s", change.Field)
	}
	return validator(change.Value)
}

func validateTokenClassChange(change Pair) error {
	validators := map[AttributeKey]func(string) error{
		AttributeKey_Name: validateName,
		AttributeKey_Meta: validateMeta,
	}

	validator, ok := validators[AttributeKey(AttributeKey_value[change.Field])]
	if !ok {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid field: %s", change.Field)
	}
	return validator(change.Value)
}

var _ sdk.Msg = (*MsgSend)(nil)

// ValidateBasic implements Msg.
func (m MsgSend) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}
	if err := sdk.ValidateAccAddress(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", m.To)
	}

	if err := m.Amount.ValidateBasic(); err != nil {
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
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", m.Operator)
	}
	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}
	if err := sdk.ValidateAccAddress(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", m.To)
	}

	if err := m.Amount.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgOperatorSend) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Operator)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgTransferFT)(nil)

// ValidateBasic implements Msg.
func (m MsgTransferFT) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}
	if err := sdk.ValidateAccAddress(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", m.To)
	}

	if err := validateCoins(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgTransferFT) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.From)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgTransferFTFrom)(nil)

// ValidateBasic implements Msg.
func (m MsgTransferFTFrom) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.Proxy); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid proxy address: %s", m.Proxy)
	}
	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}
	if err := sdk.ValidateAccAddress(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", m.To)
	}

	if err := validateCoins(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgTransferFTFrom) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Proxy)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgTransferNFT)(nil)

// ValidateBasic implements Msg.
func (m MsgTransferNFT) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}
	if err := sdk.ValidateAccAddress(m.To); err != nil {
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
	signer := sdk.AccAddress(m.From)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgTransferNFTFrom)(nil)

// ValidateBasic implements Msg.
func (m MsgTransferNFTFrom) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.Proxy); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid proxy address: %s", m.Proxy)
	}
	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}
	if err := sdk.ValidateAccAddress(m.To); err != nil {
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
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid holder address: %s", m.Holder)
	}
	if err := sdk.ValidateAccAddress(m.Operator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", m.Operator)
	}

	return nil
}

// GetSigners implements Msg
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
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid holder address: %s", m.Holder)
	}
	if err := sdk.ValidateAccAddress(m.Operator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", m.Operator)
	}

	return nil
}

// GetSigners implements Msg
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
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid approver address: %s", m.Approver)
	}
	if err := sdk.ValidateAccAddress(m.Proxy); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid proxy address: %s", m.Proxy)
	}

	return nil
}

// GetSigners implements Msg
func (m MsgApprove) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Approver)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgDisapprove)(nil)

// ValidateBasic implements Msg.
func (m MsgDisapprove) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.Approver); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid approver address: %s", m.Approver)
	}
	if err := sdk.ValidateAccAddress(m.Proxy); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid proxy address: %s", m.Proxy)
	}

	return nil
}

// GetSigners implements Msg
func (m MsgDisapprove) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Approver)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgCreateContract)(nil)

// ValidateBasic implements Msg.
func (m MsgCreateContract) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(m.Owner); err != nil {
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
	signer := sdk.AccAddress(m.Owner)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgIssueFT)(nil)

// ValidateBasic implements Msg.
func (m MsgIssueFT) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.Owner); err != nil {
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

	if err := sdk.ValidateAccAddress(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", m.To)
	}

	if err := validateAmount(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgIssueFT) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Owner)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgIssueNFT)(nil)

// ValidateBasic implements Msg.
func (m MsgIssueNFT) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := validateName(m.Name); err != nil {
		return err
	}

	if err := validateMeta(m.Meta); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.Owner); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", m.Owner)
	}

	return nil
}

// GetSigners implements Msg
func (m MsgIssueNFT) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Owner)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgMintFT)(nil)

// ValidateBasic implements Msg.
func (m MsgMintFT) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}
	if err := sdk.ValidateAccAddress(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", m.To)
	}

	if err := validateCoins(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgMintFT) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.From)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgMintNFT)(nil)

// ValidateBasic implements Msg.
func (m MsgMintNFT) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}
	if err := sdk.ValidateAccAddress(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", m.To)
	}

	if len(m.Params) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("mint params cannot be empty")
	}
	for _, param := range m.Params {
		classID := param.TokenType + fmt.Sprintf("%08x", 0)
		if err := ValidateNFTClassID(classID); err != nil {
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
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}

	if err := m.Amount.ValidateBasic(); err != nil {
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
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", m.Operator)
	}
	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}

	if err := m.Amount.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgOperatorBurn) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Operator)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgBurnFT)(nil)

// ValidateBasic implements Msg.
func (m MsgBurnFT) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}

	if err := validateCoins(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgBurnFT) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.From)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgBurnFTFrom)(nil)

// ValidateBasic implements Msg.
func (m MsgBurnFTFrom) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.Proxy); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid proxy address: %s", m.Proxy)
	}
	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}

	if err := validateCoins(m.Amount); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgBurnFTFrom) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Proxy)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgBurnNFT)(nil)

// ValidateBasic implements Msg.
func (m MsgBurnNFT) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}

	if len(m.TokenIds) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("token ids cannot be empty")
	}
	for _, id := range m.TokenIds {
		if err := ValidateNFTID(id); err != nil {
			return err
		}
	}

	return nil
}

// GetSigners implements Msg
func (m MsgBurnNFT) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.From)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgBurnNFTFrom)(nil)

// ValidateBasic implements Msg.
func (m MsgBurnNFTFrom) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.Proxy); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid proxy address: %s", m.Proxy)
	}
	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}

	if len(m.TokenIds) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("token ids cannot be empty")
	}
	for _, id := range m.TokenIds {
		if err := ValidateNFTID(id); err != nil {
			return err
		}
	}

	return nil
}

// GetSigners implements Msg
func (m MsgBurnNFTFrom) GetSigners() []sdk.AccAddress {
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
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", m.Owner)
	}

	if len(m.TokenType) != 0 {
		classID := m.TokenType + fmt.Sprintf("%08x", 0)
		if err := ValidateClassID(classID); err != nil {
			return err
		}
		if err := ValidateFTClassID(classID); err == nil && len(m.TokenIndex) == 0 {
			// smells
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

		if err := validator(change); err != nil {
			return err
		}
	}

	return nil
}

// GetSigners implements Msg
func (m MsgModify) GetSigners() []sdk.AccAddress {
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
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid granter address: %s", m.Granter)
	}
	if err := sdk.ValidateAccAddress(m.Grantee); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid grantee address: %s", m.Grantee)
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
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid grantee address: %s", m.Grantee)
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
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}
	if err := sdk.ValidateAccAddress(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", m.To)
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
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
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

var _ sdk.Msg = (*MsgAttach)(nil)

// ValidateBasic implements Msg.
func (m MsgAttach) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.From); err != nil {
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
	signer := sdk.AccAddress(m.From)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgDetach)(nil)

// ValidateBasic implements Msg.
func (m MsgDetach) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}

	if err := ValidateTokenID(m.TokenId); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgDetach) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.From)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgOperatorAttach)(nil)

// ValidateBasic implements Msg.
func (m MsgOperatorAttach) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.Operator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", m.Operator)
	}
	if err := sdk.ValidateAccAddress(m.Owner); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", m.Owner)
	}

	if err := ValidateTokenID(m.Id); err != nil {
		return err
	}
	if err := ValidateTokenID(m.To); err != nil {
		return err
	}

	if m.Id == m.To {
		return sdkerrors.ErrInvalidRequest.Wrap("cannot attach token to itself")
	}

	return nil
}

// GetSigners implements Msg
func (m MsgOperatorAttach) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Operator)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgOperatorDetach)(nil)

// ValidateBasic implements Msg.
func (m MsgOperatorDetach) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.Operator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", m.Operator)
	}
	if err := sdk.ValidateAccAddress(m.Owner); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", m.Owner)
	}

	if err := ValidateTokenID(m.Id); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgOperatorDetach) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Operator)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgAttachFrom)(nil)

// ValidateBasic implements Msg.
func (m MsgAttachFrom) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.Proxy); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid proxy address: %s", m.Proxy)
	}
	if err := sdk.ValidateAccAddress(m.From); err != nil {
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
	signer := sdk.AccAddress(m.Proxy)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgDetachFrom)(nil)

// ValidateBasic implements Msg.
func (m MsgDetachFrom) ValidateBasic() error {
	if err := class.ValidateID(m.ContractId); err != nil {
		return err
	}

	if err := sdk.ValidateAccAddress(m.Proxy); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid proxy address: %s", m.Proxy)
	}
	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}

	if err := ValidateTokenID(m.TokenId); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgDetachFrom) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Proxy)
	return []sdk.AccAddress{signer}
}
