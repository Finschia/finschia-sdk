package keeper

import (
	wasmvmtypes "github.com/line/wasmvm/types"
	sdk "github.com/line/lbm-sdk/types"
)

const (
	// DefaultEventAttributeDataCost is how much SDK gas is charged *per byte* for attribute data in events.
	// This is used with len(key) + len(value)
	DefaultEventAttributeDataCost uint64 = 1
	// DefaultContractMessageDataCost is how much SDK gas is charged *per byte* of the message that goes to the contract
	// This is used with len(msg)
	DefaultContractMessageDataCost uint64 = 1
	// DefaultPerAttributeCost is how much SDK gas we charge per attribute count.
	DefaultPerAttributeCost uint64 = 10
	// DefaultEventAttributeDataFreeTier number of bytes of total attribute data we do not charge.
	DefaultEventAttributeDataFreeTier = 100
)

// GasRegister abstract source for gas costs
type GasRegister interface {
	// NewContractInstanceCosts costs to crate a new contract instance from code
	// EventCosts costs to persist an event
	EventCosts(evts []wasmvmtypes.EventAttribute) sdk.Gas
}

// WasmGasRegisterConfig config type
type WasmGasRegisterConfig struct {
	// EventPerAttributeCost is how much SDK gas is charged *per byte* for attribute data in events.
	// This is used with len(key) + len(value)
	EventPerAttributeCost sdk.Gas
	// EventAttributeDataCost is how much SDK gas is charged *per byte* for attribute data in events.
	// This is used with len(key) + len(value)
	EventAttributeDataCost sdk.Gas
	// EventAttributeDataFreeTier number of bytes of total attribute data that is free of charge
	EventAttributeDataFreeTier int
	// ContractMessageDataCost SDK gas charged *per byte* of the message that goes to the contract
	// This is used with len(msg)
	ContractMessageDataCost sdk.Gas
}

// DefaultGasRegisterConfig default values
func DefaultGasRegisterConfig() WasmGasRegisterConfig {
	return WasmGasRegisterConfig{
		EventPerAttributeCost:      DefaultPerAttributeCost,
		EventAttributeDataCost:     DefaultEventAttributeDataCost,
		EventAttributeDataFreeTier: DefaultEventAttributeDataFreeTier,
		ContractMessageDataCost:    DefaultContractMessageDataCost,
	}
}

// WasmGasRegister implements GasRegister interface
type WasmGasRegister struct {
	c WasmGasRegisterConfig
}

// NewDefaultWasmGasRegister creates instance with default values
func NewDefaultWasmGasRegister() WasmGasRegister {
	return NewWasmGasRegister(DefaultGasRegisterConfig())
}

// NewWasmGasRegister constructor
func NewWasmGasRegister(c WasmGasRegisterConfig) WasmGasRegister {
	return WasmGasRegister{
		c: c,
	}
}

// EventCosts costs to persist an event
func (g WasmGasRegister) EventCosts(evts []wasmvmtypes.EventAttribute) sdk.Gas {
	if len(evts) == 0 {
		return 0
	}
	var storedBytes int
	for _, l := range evts {
		storedBytes += len(l.Key) + len(l.Value)
	}
	// apply free tier
	if storedBytes <= g.c.EventAttributeDataFreeTier {
		storedBytes = 0
	} else {
		storedBytes -= g.c.EventAttributeDataFreeTier
	}
	// total Length * costs + attribute count * costs
	r := sdk.NewIntFromUint64(g.c.EventAttributeDataCost).Mul(sdk.NewIntFromUint64(uint64(storedBytes))).
		Add(sdk.NewIntFromUint64(g.c.EventPerAttributeCost).Mul(sdk.NewIntFromUint64(uint64(len(evts)))))
	if !r.IsUint64() {
		panic(sdk.ErrorOutOfGas{Descriptor: "overflow"})
	}
	return r.Uint64()
}
