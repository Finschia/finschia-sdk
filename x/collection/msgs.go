package collection

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	lengthClassID = 8

	nameLengthLimit = 20
	uriLengthLimit  = 1000
	metaLengthLimit = 1000
	changesLimit    = 100
)

var (
	patternAll  = fmt.Sprintf(`[0-9a-f]{%d}`, lengthClassID)
	patternZero = fmt.Sprintf(`0{%d}`, lengthClassID)

	patternClassID          = patternAll
	patternLegacyNFTClassID = fmt.Sprintf(`[1-9a-f][0-9a-f]{%d}`, lengthClassID-1)

	// regexps for class ids
	reClassID          = regexp.MustCompile(fmt.Sprintf(`^%s$`, patternClassID))
	reLegacyNFTClassID = regexp.MustCompile(fmt.Sprintf(`^%s$`, patternLegacyNFTClassID))

	// regexps for token ids
	reTokenID        = regexp.MustCompile(fmt.Sprintf(`^%s%s$`, patternClassID, patternAll))
	reLegacyNFTID    = regexp.MustCompile(fmt.Sprintf(`^%s%s$`, patternLegacyNFTClassID, patternAll))
	reLegacyIdxNFTID = regexp.MustCompile(fmt.Sprintf(`^%s%s$`, patternLegacyNFTClassID, patternZero))

	// regexps for contract ids
	reContractIDString = `[0-9a-f]{8,8}`
	reContractID       = regexp.MustCompile(fmt.Sprintf(`^%s$`, reContractIDString))
)

func NewNFTID(classID string, number int) string {
	return newTokenID(classID, math.NewUint(uint64(number)))
}

func newTokenID(classID string, number math.Uint) string {
	numberFormat := "%0" + fmt.Sprintf("%d", lengthClassID) + "x"
	return classID + fmt.Sprintf(numberFormat, number.Uint64())
}

func SplitTokenID(tokenID string) (classID string) {
	return tokenID[:lengthClassID]
}

func ValidateContractID(id string) error {
	if !reContractID.MatchString(id) {
		return ErrInvalidContractID.Wrapf("invalid contract id: %s", id)
	}
	return nil
}

func ValidateClassID(id string) error {
	return validateID(id, reClassID)
}

// Deprecated: do not use (no successor).
func ValidateLegacyNFTClassID(id string) error {
	// daphne emits ErrInvalidTokenID here, but it's against to the spec.
	if err := validateID(id, reLegacyNFTClassID); err != nil {
		return ErrInvalidTokenType.Wrap(err.Error())
	}

	return nil
}

func ValidateTokenID(id string) error {
	if err := validateID(id, reTokenID); err != nil {
		return ErrInvalidTokenID.Wrap(err.Error())
	}

	return nil
}

func ValidateNFTID(id string) error {
	if err := ValidateTokenID(id); err != nil {
		return err
	}
	return nil
}

// Deprecated: do not use (no successor).
func ValidateLegacyNFTID(id string) error {
	if err := validateID(id, reLegacyNFTID); err != nil {
		return ErrInvalidTokenID.Wrap(err.Error())
	}

	return nil
}

func ValidateLegacyIdxNFT(id string) error {
	if err := validateID(id, reLegacyIdxNFTID); err != nil {
		return ErrInvalidTokenID.Wrap(err.Error())
	}

	return nil
}

func validateID(id string, reg *regexp.Regexp) error {
	if !reg.MatchString(id) {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid id: %s", id)
	}
	return nil
}

func validateName(name string) error {
	if err := validateStringSize(name, nameLengthLimit, "name"); err != nil {
		return ErrInvalidNameLength.Wrap(err.Error())
	}

	return nil
}

func validateURI(uri string) error {
	if err := validateStringSize(uri, uriLengthLimit, "uri"); err != nil {
		return ErrInvalidBaseImgURILength.Wrap(err.Error())
	}

	return nil
}

func validateMeta(meta string) error {
	if err := validateStringSize(meta, metaLengthLimit, "meta"); err != nil {
		return ErrInvalidMetaLength.Wrap(err.Error())
	}

	return nil
}

func validateStringSize(str string, limit int, name string) error {
	if length := utf8.RuneCountInString(str); length > limit {
		return sdkerrors.ErrInvalidRequest.Wrapf("%s cannot exceed %d in length: current %d", name, limit, length)
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
		AttributeKeyBaseImgURI.String(): validateURI,
		AttributeKeyMeta.String():       validateMeta,
		AttributeKeyURI.String():        validateURI,
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
		return ErrInvalidChangesField.Wrapf("invalid field: %s", change.Key)
	}
	return validator(change.Value)
}

// ValidateBasic implements Msg.
func (m MsgSendNFT) ValidateBasic() error {
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
		return ErrEmptyField.Wrap("token ids cannot be empty")
	}
	for _, id := range m.TokenIds {
		if err := ValidateTokenID(id); err != nil {
			return err
		}
	}

	return nil
}

// ValidateBasic implements Msg.
func (m MsgOperatorSendNFT) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.Operator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", m.Operator)
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}
	if _, err := sdk.AccAddressFromBech32(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", m.To)
	}

	if len(m.TokenIds) == 0 {
		return ErrEmptyField.Wrap("token ids cannot be empty")
	}
	for _, id := range m.TokenIds {
		if err := ValidateTokenID(id); err != nil {
			return err
		}
	}

	return nil
}

// ValidateBasic implements Msg.
func (m MsgAuthorizeOperator) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.Holder); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid holder address: %s", m.Holder)
	}
	if _, err := sdk.AccAddressFromBech32(m.Operator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", m.Operator)
	}

	if m.Operator == m.Holder {
		return ErrApproverProxySame
	}

	return nil
}

// ValidateBasic implements Msg.
func (m MsgRevokeOperator) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.Holder); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid holder address: %s", m.Holder)
	}
	if _, err := sdk.AccAddressFromBech32(m.Operator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", m.Operator)
	}

	if m.Operator == m.Holder {
		return ErrApproverProxySame
	}

	return nil
}

// ValidateBasic implements Msg.
func (m MsgCreateContract) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", m.Owner)
	}

	if err := validateName(m.Name); err != nil {
		return err
	}

	if err := validateURI(m.Uri); err != nil {
		return err
	}

	if err := validateMeta(m.Meta); err != nil {
		return err
	}

	return nil
}

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
		return ErrEmptyField.Wrap("mint params cannot be empty")
	}
	for _, param := range m.Params {
		classID := param.TokenType
		if err := ValidateLegacyNFTClassID(classID); err != nil {
			return err
		}

		if len(param.Name) == 0 {
			return ErrInvalidTokenName
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

// ValidateBasic implements Msg.
func (m MsgBurnNFT) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}

	if len(m.TokenIds) == 0 {
		return ErrEmptyField.Wrap("token ids cannot be empty")
	}
	for _, id := range m.TokenIds {
		if err := ValidateLegacyNFTID(id); err != nil {
			return err
		}
	}

	return nil
}

// ValidateBasic implements Msg.
func (m MsgOperatorBurnNFT) ValidateBasic() error {
	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.Operator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", m.Operator)
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}

	if len(m.TokenIds) == 0 {
		return ErrEmptyField.Wrap("token ids cannot be empty")
	}
	for _, id := range m.TokenIds {
		if err := ValidateLegacyNFTID(id); err != nil {
			return err
		}
	}

	return nil
}

func UpdateMsgModify(msg *MsgModify) {
	for i, change := range msg.Changes {
		key := change.Key
		converted := CollectionAttrCanonicalKey(key)
		if converted != key {
			msg.Changes[i].Key = converted
		}
	}
}

// ValidateBasic implements Msg.
func (m MsgModify) ValidateBasic() error {
	UpdateMsgModify(&m)

	if err := ValidateContractID(m.ContractId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", m.Owner)
	}

	if len(m.TokenType) != 0 {
		classID := m.TokenType
		if err := ValidateClassID(classID); err != nil {
			return ErrInvalidTokenType.Wrap(err.Error())
		}
	}

	if len(m.TokenIndex) != 0 {
		tokenID := m.TokenType + m.TokenIndex
		if err := ValidateTokenID(tokenID); err != nil {
			return ErrInvalidTokenIndex.Wrap(err.Error())
		}
		// reject modifying nft class with token index filled (daphne compat.)
		if ValidateLegacyIdxNFT(tokenID) == nil {
			return ErrInvalidTokenIndex.Wrap("cannot modify nft class with index filled")
		}
	}

	validator := validateTokenClassChange
	if len(m.TokenType) == 0 {
		if len(m.TokenIndex) == 0 {
			validator = validateContractChange
		} else {
			return ErrTokenIndexWithoutType.Wrap("token index without type")
		}
	}
	if len(m.Changes) == 0 {
		return ErrEmptyChanges.Wrap("empty changes")
	}
	if len(m.Changes) > changesLimit {
		return ErrInvalidChangesFieldCount.Wrapf("the number of changes exceeds the limit: %d > %d", len(m.Changes), changesLimit)
	}
	seenKeys := map[string]bool{}
	for _, change := range m.Changes {
		key := change.Key
		if seenKeys[key] {
			return ErrDuplicateChangesField.Wrapf("duplicate keys: %s", change.Key)
		}
		seenKeys[key] = true

		attribute := Attribute{
			Key:   change.Key,
			Value: change.Value,
		}
		if err := validator(attribute); err != nil {
			return err
		}
	}

	return nil
}

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
