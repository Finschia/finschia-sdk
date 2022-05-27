package class

import (
	"fmt"
	"regexp"

	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

var (
	// reClassIDString must be a hex string of 8 characters long
	reClassIDString = `[0-9a-f]{8,8}`
	reClassID       = regexp.MustCompile(fmt.Sprintf(`^%s$`, reClassIDString))
)

// ValidateID returns whether the class id is valid
func ValidateID(id string) error {
	if !reClassID.MatchString(id) {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid class id: %s", id)
	}
	return nil
}
