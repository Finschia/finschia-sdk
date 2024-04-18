package types

import (
	"errors"
	fmt "fmt"
	"strings"

	sdk "github.com/Finschia/finschia-sdk/types"
)

// NewSwapped creates a new Swapped instance
func NewSwapped(
	oldCoin sdk.Coin,
	newCoin sdk.Coin,
) Swapped {
	return Swapped{
		OldCoinAmount: oldCoin,
		NewCoinAmount: newCoin,
	}
}

// DefaultSwapped returns an initial Swapped object
func DefaultSwapped() Swapped {
	return NewSwapped(
		sdk.NewCoin("cony", sdk.NewInt(0)),
		sdk.NewCoin(DefaultNewCoinDenom, sdk.NewInt(0)),
	)
}

func validateCoin(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid coin amount: %T", i)
	}

	if strings.TrimSpace(v.Denom) == "" {
		return errors.New("denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v.Denom); err != nil {
		return err
	}

	return nil
}

// Validate validates the set of swapped
func (s Swapped) Validate() error {
	if err := validateCoin(s.OldCoinAmount); err != nil {
		return err
	}
	if err := validateCoin(s.NewCoinAmount); err != nil {
		return err
	}
	return nil
}
