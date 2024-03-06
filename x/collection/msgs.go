package collection

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	"cosmossdk.io/math"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	lengthClassID = 8

	nameLengthLimit = 20
	uriLengthLimit  = 1000
	metaLengthLimit = 1000
	ChangesLimit    = 100
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

func ValidateName(name string) error {
	if err := validateStringSize(name, nameLengthLimit, "name"); err != nil {
		return ErrInvalidNameLength.Wrap(err.Error())
	}

	return nil
}

func ValidateURI(uri string) error {
	if err := validateStringSize(uri, uriLengthLimit, "uri"); err != nil {
		return ErrInvalidBaseImgURILength.Wrap(err.Error())
	}

	return nil
}

func ValidateMeta(meta string) error {
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

func ValidateLegacyPermission(permission string) error {
	return ValidatePermission(Permission(LegacyPermissionFromString(permission)))
}

func ValidatePermission(permission Permission) error {
	if p := Permission_value[Permission_name[int32(permission)]]; p == 0 {
		return ErrInvalidPermission.Wrap(permission.String())
	}
	return nil
}

func ValidateContractChange(change Attribute) error {
	validators := map[string]func(string) error{
		AttributeKeyName.String():       ValidateName,
		AttributeKeyBaseImgURI.String(): ValidateURI,
		AttributeKeyMeta.String():       ValidateMeta,
		AttributeKeyURI.String():        ValidateURI,
	}

	return validateChange(change, validators)
}

func ValidateTokenClassChange(change Attribute) error {
	validators := map[string]func(string) error{
		AttributeKeyName.String(): ValidateName,
		AttributeKeyMeta.String(): ValidateMeta,
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
