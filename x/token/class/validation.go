package class

import (
	"fmt"
	"regexp"

	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

var (
	// reContractIDString must be a hex string of 8 characters long
	reContractIDString = `[0-9a-f]{8,8}`
	reContractID       = regexp.MustCompile(fmt.Sprintf(`^%s$`, reContractIDString))
)

// ValidateID returns whether the contract id is valid
func ValidateID(id string) error {
	if !reContractID.MatchString(id) {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid contract id: %s", id)
	}
	return nil
}
