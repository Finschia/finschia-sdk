package types

import (
	"gopkg.in/yaml.v2"
)

const (
	DefaultFetchIntervals uint64 = 3600
)

// NewParams creates a new Params instance
func NewParams(
	FetchIntervals uint64,
) Params {
	return Params{
		FetchIntervals: FetchIntervals,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(DefaultFetchIntervals)
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
