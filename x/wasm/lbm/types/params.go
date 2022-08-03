package types

import (
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	sdkerrors "github.com/line/lbm-sdk/types/errors"
	paramtypes "github.com/line/lbm-sdk/x/params/types"
	wasmtypes "github.com/line/lbm-sdk/x/wasm/types"
)

//const (
//	// DefaultGasMultiplier is how many CosmWasm gas points = 1 Cosmos SDK gas point.
//	//
//	// CosmWasm gas strategy is documented in https://https://github.com/line/cosmwasm/blob/v1.0.0-0.6.0/docs/GAS.md.
//	// LBM SDK reference costs can be found here: https://github.com/line/lbm-sdk/blob/main/store/types/gas.go#L199-L209.
//	//
//	// The original multiplier of 100 up to CosmWasm 0.16 was based on
//	//     "A write at ~3000 gas and ~200us = 10 gas per us (microsecond) cpu/io
//	//     Rough timing have 88k gas at 90us, which is equal to 1k sdk gas... (one read)"
//	// as well as manual Wasmer benchmarks from 2019. This was then multiplied by 150_000
//	// in the 0.16 -> 1.0 upgrade (https://github.com/CosmWasm/cosmwasm/pull/1120).
//	//
//	// The multiplier deserves more reproducible benchmarking and a strategy that allows easy adjustments.
//	// This is tracked in https://github.com/CosmWasm/wasmd/issues/566 and https://github.com/CosmWasm/wasmd/issues/631.
//	// Gas adjustments are consensus breaking but may happen in any release marked as consensus breaking.
//	// Do not make assumptions on how much gas an operation will consume in places that are hard to adjust,
//	// such as hardcoding them in contracts.
//	//
//	// Please note that all gas prices returned to wasmvm should have this multiplied.
//	DefaultGasMultiplier uint64 = 140_000_000
//	// InstanceCost is how much SDK gas we charge each time we load a WASM instance.
//	// Creating a new instance is costly, and this helps put a recursion limit to contracts calling contracts.
//	DefaultInstanceCost = 60_000
//	// CompileCost is how much SDK gas we charge *per byte* for compiling WASM code.
//	DefaultCompileCost = 3
//)

//var ParamStoreKeyUploadAccess = []byte("uploadAccess")
//var ParamStoreKeyInstantiateAccess = []byte("instantiateAccess")
//var ParamStoreKeyGasMultiplier = []byte("gasMultiplier")
//var ParamStoreKeyInstanceCost = []byte("instanceCost")
//var ParamStoreKeyCompileCost = []byte("compileCost")

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default wasm parameters
func DefaultParams() Params {
	return Params{
		CodeUploadAccess:             wasmtypes.AllowEverybody,
		InstantiateDefaultPermission: wasmtypes.AccessTypeEverybody,
		GasMultiplier:                wasmtypes.DefaultGasMultiplier,
		InstanceCost:                 wasmtypes.DefaultInstanceCost,
		CompileCost:                  wasmtypes.DefaultCompileCost,
	}
}

func (p Params) String() string {
	out, err := yaml.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(wasmtypes.ParamStoreKeyUploadAccess, &p.CodeUploadAccess, validateAccessConfig),
		paramtypes.NewParamSetPair(wasmtypes.ParamStoreKeyInstantiateAccess, &p.InstantiateDefaultPermission, validateAccessType),
		paramtypes.NewParamSetPair(wasmtypes.ParamStoreKeyGasMultiplier, &p.GasMultiplier, validateGasMultiplier),
		paramtypes.NewParamSetPair(wasmtypes.ParamStoreKeyInstanceCost, &p.InstanceCost, validateInstanceCost),
		paramtypes.NewParamSetPair(wasmtypes.ParamStoreKeyCompileCost, &p.CompileCost, validateCompileCost),
	}
}

// ValidateBasic performs basic validation on wasm parameters
func (p Params) ValidateBasic() error {
	if err := validateAccessType(p.InstantiateDefaultPermission); err != nil {
		return errors.Wrap(err, "instantiate default permission")
	}
	if err := validateAccessConfig(p.CodeUploadAccess); err != nil {
		return errors.Wrap(err, "upload access")
	}
	if err := validateGasMultiplier(p.GasMultiplier); err != nil {
		return errors.Wrap(err, "gas multiplier")
	}
	if err := validateInstanceCost(p.InstanceCost); err != nil {
		return errors.Wrap(err, "instance cost")
	}
	if err := validateCompileCost(p.CompileCost); err != nil {
		return errors.Wrap(err, "compile cost")
	}
	return nil
}

func validateAccessConfig(i interface{}) error {
	v, ok := i.(wasmtypes.AccessConfig)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return v.ValidateBasic()
}

func validateAccessType(i interface{}) error {
	a, ok := i.(wasmtypes.AccessType)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if a == wasmtypes.AccessTypeUnspecified {
		return sdkerrors.Wrap(wasmtypes.ErrEmpty, "type")
	}
	for _, v := range wasmtypes.AllAccessTypes {
		if v == a {
			return nil
		}
	}
	return sdkerrors.Wrapf(wasmtypes.ErrInvalid, "unknown type: %q", a)
}

func validateGasMultiplier(i interface{}) error {
	a, ok := i.(uint64)
	if !ok {
		return sdkerrors.Wrapf(wasmtypes.ErrInvalid, "type: %T", i)
	}
	if a == 0 {
		return sdkerrors.Wrap(wasmtypes.ErrInvalid, "must be greater 0")
	}
	return nil
}

func validateInstanceCost(i interface{}) error {
	a, ok := i.(uint64)
	if !ok {
		return sdkerrors.Wrapf(wasmtypes.ErrInvalid, "type: %T", i)
	}
	if a == 0 {
		return sdkerrors.Wrap(wasmtypes.ErrInvalid, "must be greater 0")
	}
	return nil
}

func validateCompileCost(i interface{}) error {
	a, ok := i.(uint64)
	if !ok {
		return sdkerrors.Wrapf(wasmtypes.ErrInvalid, "type: %T", i)
	}
	if a == 0 {
		return sdkerrors.Wrap(wasmtypes.ErrInvalid, "must be greater 0")
	}
	return nil
}
