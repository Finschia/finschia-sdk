package types

import (
	"gopkg.in/yaml.v2"

	sdk "github.com/Finschia/finschia-sdk/types"
)

// NewParams creates a new Params instance
func NewParams(
	swappableNewCoinAmount sdk.Int,
) Params {
	return Params{
		SwappableNewCoinAmount: swappableNewCoinAmount,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(sdk.ZeroInt())
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateCoinAmount(p.SwappableNewCoinAmount); err != nil {
		return err
	}
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
