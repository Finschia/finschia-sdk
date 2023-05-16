package types

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

const (
	DefaultCCBatchMaxBytes  uint64 = 1000000
	DefaultSCCBatchMaxBytes uint64 = 1000000
)

// NewParams creates a new Params instance
func NewParams(
	CCBatchMaxBytes uint64,
	SCCBatchMaxBytes uint64,
) Params {
	return Params{
		CCBatchMaxBytes:  CCBatchMaxBytes,
		SCCBatchMaxBytes: SCCBatchMaxBytes,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultCCBatchMaxBytes,
		DefaultSCCBatchMaxBytes,
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateCCBatchMaxBytes(p.CCBatchMaxBytes); err != nil {
		return err
	}
	if err := validateSCCBatchMaxBytes(p.SCCBatchMaxBytes); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateCCBatchMaxBytes(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("ctc batch max bytes must be positive: %d", v)
	}

	return nil
}

func validateSCCBatchMaxBytes(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("scc batch max bytes must be positive: %d", v)
	}

	return nil
}
