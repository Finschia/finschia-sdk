package types

import (
	fmt "fmt"

	"gopkg.in/yaml.v2"

	sdk "github.com/Finschia/finschia-sdk/types"
)

// NewSwapped creates a new Swapped instance
func NewSwapped(
	oldCoinAmount sdk.Int,
	newCoinAmount sdk.Int,
) Swapped {
	return Swapped{
		OldCoinAmount: oldCoinAmount,
		NewCoinAmount: newCoinAmount,
	}
}

// DefaultSwapped returns an initial Swapped object
func DefaultSwapped() Swapped {
	return NewSwapped(sdk.ZeroInt(), sdk.ZeroInt())
}

func validateCoinAmount(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid coin amount: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("coin amount must be not nil")
	}

	if v.LT(sdk.ZeroInt()) {
		return fmt.Errorf("coin amount cannot be lower than 0")
	}

	return nil
}

// Validate validates the set of swapped
func (s Swapped) Validate() error {
	if err := validateCoinAmount(s.OldCoinAmount); err != nil {
		return err
	}
	if err := validateCoinAmount(s.NewCoinAmount); err != nil {
		return err
	}
	return nil
}

// String implements the Stringer interface.
func (s Swapped) String() string {
	out, _ := yaml.Marshal(s)
	return string(out)
}
