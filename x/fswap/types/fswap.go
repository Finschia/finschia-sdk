package types

import (
	"gopkg.in/yaml.v2"

	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

// ValidateBasic validates the set of Swap
func (s *Swap) ValidateBasic() error {
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
	if s.SwapRate.IsZero() {
		return sdkerrors.ErrInvalidRequest.Wrap("swap rate cannot be zero")
	}
	return nil
}

func (s *Swap) String() string {
	out, _ := yaml.Marshal(s)
	return string(out)
}

func (s *SwapStats) ValidateBasic() error {
	if s.SwapCount < 0 {
		return ErrInvalidState.Wrap("swap count cannot be negative")
	}
	return nil
}

func (s *SwapStats) String() string {
	out, _ := yaml.Marshal(s)
	return string(out)
}

// ValidateBasic validates the set of Swapped
func (s *Swapped) ValidateBasic() error {
	if err := validateCoinAmount(s.FromCoinAmount); err != nil {
		return err
	}
	if err := validateCoinAmount(s.ToCoinAmount); err != nil {
		return err
	}
	return nil
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

func (s *Swapped) String() string {
	out, _ := yaml.Marshal(s)
	return string(out)
}
