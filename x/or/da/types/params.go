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
	DefaultMinQueueTxGas           uint64 = DefaultEnqueueL2GasPrepaid * 2
	DefaultL2GasDiscountDivisor    uint64 = 100
	DefaultEnqueueL2GasPrepaid     uint64 = 30000
	DefaultQueueTxExpirationWindow uint64 = 600
	DefaultFraudProofWindow        uint64 = 7 * 24 * 60 * 60 / 6
	DefaultSequencerPublishWindow  uint64 = 600
)

// NewParams creates a new Params instance
func NewParams(
	CCBatchMaxBytes uint64,
	MaxQueueTxSize uint64,
	MinQueueTxGas uint64,
	L2GasDiscountDivisor uint64,
	EnqueueL2GasPrepaid uint64,
	QueueTxExpirationWindow uint64,
	SCCBatchMaxBytes uint64,
	FraudProofWindow uint64,
	SequencerPublishWindow uint64,
) Params {
	return Params{
		CCBatchMaxBytes:         CCBatchMaxBytes,
		MaxQueueTxSize:          MaxQueueTxSize,
		MinQueueTxGas:           MinQueueTxGas,
		L2GasDiscountDivisor:    L2GasDiscountDivisor,
		EnqueueL2GasPrepaid:     EnqueueL2GasPrepaid,
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
		DefaultL2GasDiscountDivisor,
		DefaultEnqueueL2GasPrepaid,
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
