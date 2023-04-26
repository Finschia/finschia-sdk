package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	paramtypes "github.com/Finschia/finschia-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyPlaceholder = []byte("Placeholder")
	// TODO: Determine the default value
	DefaultPlaceholder string = "placeholder"
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	placeholder string,
) Params {
	return Params{
		Placeholder: placeholder,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultPlaceholder,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyPlaceholder, &p.Placeholder, validatePlaceholder),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validatePlaceholder(p.Placeholder); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// validatePlaceholder validates the Placeholder param
func validatePlaceholder(v interface{}) error {
	placeholder, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	// TODO implement validation
	_ = placeholder

	return nil
}
