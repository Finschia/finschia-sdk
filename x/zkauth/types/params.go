package types

import (
	"gopkg.in/yaml.v2"
)

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{}
}

// Validate validates the set of params
func (p Params) Validate() error {

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
