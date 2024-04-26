package types

import (
	"gopkg.in/yaml.v2"

	sdk "github.com/Finschia/finschia-sdk/types"
)

// Validate validates the set of swapped
func (fi FswapInit) ValidateBasic() error {
	// todo vadalidate
	return nil
}

// String implements the Stringer interface.
func (fi FswapInit) String() string {
	out, _ := yaml.Marshal(fi)
	return string(out)
}

// NewFswapInit creates a new FswapInit instance
func NewFswapInit(
	fromDenom string,
	toDenom string,
	amountLimit sdk.Int,
	swapRate sdk.Int,
) FswapInit {
	return FswapInit{
		FromDenom:   fromDenom,
		ToDenom:     toDenom,
		AmountLimit: amountLimit,
		SwapRate:    swapRate,
	}
}

// DefaultFswapInit returns an initial FswapInit object
func DefaultFswapInit() FswapInit {
	return NewFswapInit("", "", sdk.ZeroInt(), sdk.ZeroInt())
}
