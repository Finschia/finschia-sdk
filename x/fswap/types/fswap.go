package types

import (
	"errors"
	fmt "fmt"
	"strings"

	sdk "github.com/Finschia/finschia-sdk/types"
)

// NewSwapped creates a new Swapped instance
func NewSwapped(
	config Config,
) Swapped {
	return Swapped{
		OldCoinAmount: sdk.NewCoin(config.OldCoinDenom, sdk.NewInt(0)),
		NewCoinAmount: sdk.NewCoin(config.NewCoinDenom, sdk.NewInt(0)),
	}
}

// DefaultSwapped returns an initial Swapped object
func DefaultSwapped() Swapped {
	return NewSwapped(DefaultConfig())
}

func ValidateCoin(i interface{}) error {
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
	if err := ValidateCoin(s.OldCoinAmount); err != nil {
		return err
	}
	if err := ValidateCoin(s.NewCoinAmount); err != nil {
		return err
	}
	return nil
}
