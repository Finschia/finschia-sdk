package types

import (
	"fmt"
	"math"

	"gopkg.in/yaml.v2"
)

const (
	DefaultCCBatchMaxBytes         uint64 = 1000_000
	DefaultSCCBatchMaxBytes        uint64 = 1000_000
	DefaultMaxQueueTxSize          uint64 = math.MaxUint16
	DefaultMinQueueTxGas           uint64 = 300000
	DefaultQueueTxExpirationWindow uint64 = 600
	DefaultFraudProofWindow        int64  = 7 * 24 * 60 * 60 / 6
	DefaultSequencerPublishWindow  int64  = 600
)

// NewParams creates a new Params instance
func NewParams(
	CCBatchMaxBytes uint64,
	MaxQueueTxSize uint64,
	MinQueueTxGas uint64,
	QueueTxExpirationWindow uint64,
	SCCBatchMaxBytes uint64,
	FraudProofWindow int64,
	SequencerPublishWindow int64,
) Params {
	return Params{
		CCBatchMaxBytes:         CCBatchMaxBytes,
		MaxQueueTxSize:          MaxQueueTxSize,
		MinQueueTxGas:           MinQueueTxGas,
		QueueTxExpirationWindow: QueueTxExpirationWindow,
		SCCBatchMaxBytes:        SCCBatchMaxBytes,
		FraudProofWindow:        FraudProofWindow,
		SequencerPublishWindow:  SequencerPublishWindow,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultCCBatchMaxBytes,
		DefaultMaxQueueTxSize,
		DefaultMinQueueTxGas,
		DefaultQueueTxExpirationWindow,
		DefaultSCCBatchMaxBytes,
		DefaultFraudProofWindow,
		DefaultSequencerPublishWindow,
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
		return fmt.Errorf("cc batch max bytes must be positive: %d", v)
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
