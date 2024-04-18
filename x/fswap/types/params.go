package types

import (
	"errors"
	fmt "fmt"
	"strings"

	"gopkg.in/yaml.v2"

	sdk "github.com/Finschia/finschia-sdk/types"
)

const (
	DefaultNewCoinDenom string = "PDT"
)

// NewParams creates a new Params instance
func NewParams(
	newCoinDenom string,
) Params {
	return Params{
		NewCoinDenom: newCoinDenom,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultNewCoinDenom,
	)
}

func validateNewCoinDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("new denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateNewCoinDenom(p.NewCoinDenom); err != nil {
		return err
	}
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
