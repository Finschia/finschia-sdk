package types

import (
	"gopkg.in/yaml.v2"

	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
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
		return sdkerrors.ErrInvalidCoins.Wrapf("invalid coin amount: %T", i)
	}
	if v.IsNil() {
		return sdkerrors.ErrInvalidCoins.Wrap("coin amount must be not nil")
	}
	if err := v.Validate(); err != nil {
		return sdkerrors.ErrInvalidCoins.Wrap(err.Error())
	}
	return nil
}

// Validate validates the set of Swapped
func (s *Swapped) Validate() error {
	if err := validateCoinAmount(s.FromCoinAmount); err != nil {
		return err
	}
	if err := validateCoinAmount(s.ToCoinAmount); err != nil {
		return err
	}
	return nil
}

// String implements the Stringer interface.
func (s *Swapped) String() string {
	out, _ := yaml.Marshal(s)
	return string(out)
}

// ValidateBasic validates the set of SwapInit
func (s *SwapInit) ValidateBasic() error {
	if s.FromDenom == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("from denomination cannot be empty")
	}
	if s.ToDenom == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("to denomination cannot be empty")
	}
	if s.FromDenom == s.ToDenom {
		return sdkerrors.ErrInvalidRequest.Wrap("from denomination cannot be equal to to denomination")
	}
	if s.AmountCapForToDenom.LT(sdk.OneInt()) {
		return sdkerrors.ErrInvalidRequest.Wrap("amount cannot be less than one")
	}
	if s.SwapMultiple.LT(sdk.OneInt()) {
		return sdkerrors.ErrInvalidRequest.Wrap("swap multiple cannot be less than one")
	}
	return nil
}

// String implements the Stringer interface.
func (s *SwapInit) String() string {
	out, _ := yaml.Marshal(s)
	return string(out)
}

func (s *SwapInit) IsEmpty() bool {
	if s.FromDenom == "" {
		return true
	}
	if s.ToDenom == "" {
		return true
	}
	if s.AmountCapForToDenom.IsZero() {
		return true
	}
	if s.SwapMultiple.IsZero() {
		return true
	}
	return false
}
