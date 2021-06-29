package types

import (
	"encoding/json"
	"fmt"

	"github.com/gogo/protobuf/jsonpb"
	sdk "github.com/line/lfb-sdk/types"
	sdkerrors "github.com/line/lfb-sdk/types/errors"
	paramtypes "github.com/line/lfb-sdk/x/params/types"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	// DefaultParamspace for params keeper
	DefaultParamspace = ModuleName
	// DefaultMaxWasmCodeSize limit max bytes read to prevent gzip bombs
	DefaultMaxWasmCodeSize = 600 * 1024
	// GasMultiplier is how many cosmwasm gas points = 1 sdk gas point
	// SDK reference costs can be found here: https://github.com/cosmos/cosmos-sdk/blob/02c6c9fafd58da88550ab4d7d494724a477c8a68/store/types/gas.go#L153-L164
	// A write at ~3000 gas and ~200us = 10 gas per us (microsecond) cpu/io
	// Rough timing have 88k gas at 90us, which is equal to 1k sdk gas... (one read)
	//
	// Please not that all gas prices returned to the wasmer engine should have this multiplied
	DefaultGasMultiplier uint64 = 100
	// InstanceCost is how much SDK gas we charge each time we load a WASM instance.
	// Creating a new instance is costly, and this helps put a recursion limit to contracts calling contracts.
	DefaultInstanceCost = 40_000
	// CompileCost is how much SDK gas we charge *per byte* for compiling WASM code.
	DefaultCompileCost = 2
)

var ParamStoreKeyUploadAccess = []byte("uploadAccess")
var ParamStoreKeyInstantiateAccess = []byte("instantiateAccess")
var ParamStoreKeyContractStatusAccess = []byte("contractStatusAccess")
var ParamStoreKeyMaxWasmCodeSize = []byte("maxWasmCodeSize")
var ParamStoreKeyGasMultiplier = []byte("gasMultiplier")
var ParamStoreKeyInstanceCost = []byte("instanceCost")
var ParamStoreKeyCompileCost = []byte("compileCost")

var AllAccessTypes = []AccessType{
	AccessTypeNobody,
	AccessTypeOnlyAddress,
	AccessTypeEverybody,
}

func (a AccessType) With(addr sdk.AccAddress) AccessConfig {
	switch a {
	case AccessTypeNobody:
		return AllowNobody
	case AccessTypeOnlyAddress:
		if err := sdk.VerifyAddressFormat(addr); err != nil {
			panic(err)
		}
		return AccessConfig{Permission: AccessTypeOnlyAddress, Address: addr.String()}
	case AccessTypeEverybody:
		return AllowEverybody
	}
	panic("unsupported access type")
}

func (a AccessType) String() string {
	switch a {
	case AccessTypeNobody:
		return "Nobody"
	case AccessTypeOnlyAddress:
		return "OnlyAddress"
	case AccessTypeEverybody:
		return "Everybody"
	}
	return "Unspecified"
}

func (a *AccessType) UnmarshalText(text []byte) error {
	for _, v := range AllAccessTypes {
		if v.String() == string(text) {
			*a = v
			return nil
		}
	}
	*a = AccessTypeUnspecified
	return nil
}
func (a AccessType) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

func (a *AccessType) MarshalJSONPB(_ *jsonpb.Marshaler) ([]byte, error) {
	return json.Marshal(a)
}

func (a *AccessType) UnmarshalJSONPB(_ *jsonpb.Unmarshaler, data []byte) error {
	return json.Unmarshal(data, a)
}

func (a AccessConfig) Equals(o AccessConfig) bool {
	return a.Permission == o.Permission && a.Address == o.Address
}

var (
	DefaultUploadAccess         = AllowEverybody
	DefaultContractStatusAccess = AllowNobody
	AllowEverybody              = AccessConfig{Permission: AccessTypeEverybody}
	AllowNobody                 = AccessConfig{Permission: AccessTypeNobody}
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default wasm parameters
func DefaultParams() Params {
	return Params{
		CodeUploadAccess:             AllowEverybody,
		InstantiateDefaultPermission: AccessTypeEverybody,
		MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
		ContractStatusAccess:         DefaultContractStatusAccess,
		GasMultiplier:                DefaultGasMultiplier,
		InstanceCost:                 DefaultInstanceCost,
		CompileCost:                  DefaultCompileCost,
	}
}

func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyUploadAccess, &p.CodeUploadAccess, validateAccessConfig),
		paramtypes.NewParamSetPair(ParamStoreKeyInstantiateAccess, &p.InstantiateDefaultPermission, validateAccessType),
		paramtypes.NewParamSetPair(ParamStoreKeyContractStatusAccess, &p.ContractStatusAccess, validateAccessConfig),
		paramtypes.NewParamSetPair(ParamStoreKeyMaxWasmCodeSize, &p.MaxWasmCodeSize, validateMaxWasmCodeSize),
		paramtypes.NewParamSetPair(ParamStoreKeyGasMultiplier, &p.GasMultiplier, validateGasMultiplier),
		paramtypes.NewParamSetPair(ParamStoreKeyInstanceCost, &p.InstanceCost, validateInstanceCost),
		paramtypes.NewParamSetPair(ParamStoreKeyCompileCost, &p.CompileCost, validateCompileCost),
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
	if err := validateAccessConfig(p.ContractStatusAccess); err != nil {
		return errors.Wrap(err, "contract status access")
	}
	if err := validateMaxWasmCodeSize(p.MaxWasmCodeSize); err != nil {
		return errors.Wrap(err, "max wasm code size")
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
	v, ok := i.(AccessConfig)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return v.ValidateBasic()
}

func validateAccessType(i interface{}) error {
	a, ok := i.(AccessType)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if a == AccessTypeUnspecified {
		return sdkerrors.Wrap(ErrEmpty, "type")
	}
	for _, v := range AllAccessTypes {
		if v == a {
			return nil
		}
	}
	return sdkerrors.Wrapf(ErrInvalid, "unknown type: %q", a)
}

func validateMaxWasmCodeSize(i interface{}) error {
	a, ok := i.(uint64)
	if !ok {
		return sdkerrors.Wrapf(ErrInvalid, "type: %T", i)
	}
	if a == 0 {
		return sdkerrors.Wrap(ErrInvalid, "must be greater 0")
	}
	return nil
}

func validateGasMultiplier(i interface{}) error {
	a, ok := i.(uint64)
	if !ok {
		return sdkerrors.Wrapf(ErrInvalid, "type: %T", i)
	}
	if a == 0 {
		return sdkerrors.Wrap(ErrInvalid, "must be greater 0")
	}
	return nil
}

func validateInstanceCost(i interface{}) error {
	a, ok := i.(uint64)
	if !ok {
		return sdkerrors.Wrapf(ErrInvalid, "type: %T", i)
	}
	if a == 0 {
		return sdkerrors.Wrap(ErrInvalid, "must be greater 0")
	}
	return nil
}

func validateCompileCost(i interface{}) error {
	a, ok := i.(uint64)
	if !ok {
		return sdkerrors.Wrapf(ErrInvalid, "type: %T", i)
	}
	if a == 0 {
		return sdkerrors.Wrap(ErrInvalid, "must be greater 0")
	}
	return nil
}

func (a AccessConfig) ValidateBasic() error {
	switch a.Permission {
	case AccessTypeUnspecified:
		return sdkerrors.Wrap(ErrEmpty, "type")
	case AccessTypeNobody, AccessTypeEverybody:
		if len(a.Address) != 0 {
			return sdkerrors.Wrap(ErrInvalid, "address not allowed for this type")
		}
		return nil
	case AccessTypeOnlyAddress:
		_, err := sdk.AccAddressFromBech32(a.Address)
		return err
	}
	return sdkerrors.Wrapf(ErrInvalid, "unknown type: %q", a.Permission)
}

func (a AccessConfig) Allowed(actor sdk.AccAddress) bool {
	switch a.Permission {
	case AccessTypeNobody:
		return false
	case AccessTypeEverybody:
		return true
	case AccessTypeOnlyAddress:
		return a.Address == actor.String()
	default:
		panic("unknown type")
	}
}
