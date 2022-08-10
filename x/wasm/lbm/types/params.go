package types

import (
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	sdkerrors "github.com/line/lbm-sdk/types/errors"
	paramtypes "github.com/line/lbm-sdk/x/params/types"
	wasmtypes "github.com/line/lbm-sdk/x/wasm/types"
)

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
