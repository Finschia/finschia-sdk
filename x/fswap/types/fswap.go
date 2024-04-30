package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	sdk "github.com/Finschia/finschia-sdk/types"
)

// NewSwapped creates a new Swapped instance
func NewSwapped(
	oldCoinAmount sdk.Coin,
	newCoinAmount sdk.Coin,
) Swapped {
	return Swapped{
		FromCoinAmount: oldCoinAmount,
		ToCoinAmount:   newCoinAmount,
	}
}

// DefaultSwapped returns an initial Swapped object
func DefaultSwapped() Swapped {
	return NewSwapped(sdk.Coin{}, sdk.Coin{})
}

func validateCoinAmount(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid coin amount: %T", i)
	}
	if v.IsNil() {
		return fmt.Errorf("coin amount must be not nil")
	}
	if err := v.Validate(); err != nil {
		return err
	}
	return nil
}

// Validate validates the set of swapped
func (s Swapped) Validate() error {
	if err := validateCoinAmount(s.FromCoinAmount); err != nil {
		return err
	}
	if err := validateCoinAmount(s.ToCoinAmount); err != nil {
		return err
	}
	return nil
}

// String implements the Stringer interface.
func (s Swapped) String() string {
	out, _ := yaml.Marshal(s)
	return string(out)
}

// Validate validates the set of swapped
func (s SwapInit) ValidateBasic() error {
	if s.FromDenom == "" {
		return fmt.Errorf("from denomination cannot be empty")
	}
	if s.ToDenom == "" {
		return fmt.Errorf("to denomination cannot be empty")
	}
	if s.FromDenom == s.ToDenom {
		return fmt.Errorf("from denomination cannot be equal to to denomination")
	}
	if s.AmountCapForToDenom.LT(sdk.ZeroInt()) {
		return fmt.Errorf("amount cannot be less than zero")
	}
	if s.SwapMultiple.LT(sdk.ZeroInt()) {
		return fmt.Errorf("swap multiple cannot be less than zero")
	}
	return nil
}

// String implements the Stringer interface.
func (s SwapInit) String() string {
	out, _ := yaml.Marshal(s)
	return string(out)
}
