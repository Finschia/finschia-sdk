package types

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/params"
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

	DefaultHumanizeCost = 5 * DefaultGasMultiplier

	DefaultCanonicalCost = 4 * DefaultGasMultiplier
)

var ParamStoreKeyUploadAccess = []byte("uploadAccess")
var ParamStoreKeyInstantiateAccess = []byte("instantiateAccess")
var ParamStoreKeyMaxWasmCodeSize = []byte("maxWasmCodeSize")
var ParamStoreKeyGasMultiplier = []byte("gasMultiplier")
var ParamStoreKeyMaxGas = []byte("maxGas")
var ParamStoreKeyInstanceCost = []byte("instanceCost")
var ParamStoreKeyCompileCost = []byte("compileCost")
var ParamStoreKeyHumanizeCost = []byte("humanizeCost")
var ParamStoreKeyCanonicalCost = []byte("canonicalizeCost")

type AccessType string

const (
	Undefined   AccessType = "Undefined"
	Nobody      AccessType = "Nobody"
	OnlyAddress AccessType = "OnlyAddress"
	Everybody   AccessType = "Everybody"
)

var AllAccessTypes = map[AccessType]struct{}{
	Nobody:      {},
	OnlyAddress: {},
	Everybody:   {},
}

func (a AccessType) With(addr sdk.AccAddress) AccessConfig {
	switch a {
	case Nobody:
		return AllowNobody
	case OnlyAddress:
		if err := sdk.VerifyAddressFormat(addr); err != nil {
			panic(err)
		}
		return AccessConfig{Type: OnlyAddress, Address: addr}
	case Everybody:
		return AllowEverybody
	}
	panic("unsupported access type")
}

func (a *AccessType) UnmarshalText(text []byte) error {
	s := AccessType(text)
	if _, ok := AllAccessTypes[s]; ok {
		*a = s
		return nil
	}
	*a = Undefined
	return nil
}

func (a AccessType) MarshalText() ([]byte, error) {
	if _, ok := AllAccessTypes[a]; ok {
		return []byte(a), nil
	}
	return []byte(Undefined), nil
}

type AccessConfig struct {
	Type    AccessType     `json:"permission" yaml:"permission"`
	Address sdk.AccAddress `json:"address,omitempty" yaml:"address"`
}

func (v AccessConfig) Equals(o AccessConfig) bool {
	return v.Type == o.Type && v.Address.Equals(o.Address)
}

var (
	DefaultUploadAccess = AllowEverybody
	AllowEverybody      = AccessConfig{Type: Everybody}
	AllowNobody         = AccessConfig{Type: Nobody}
)

// Params defines the set of wasm parameters.
type Params struct {
	UploadAccess                 AccessConfig `json:"code_upload_access" yaml:"code_upload_access"`
	DefaultInstantiatePermission AccessType   `json:"instantiate_default_permission" yaml:"instantiate_default_permission"`
	MaxWasmCodeSize              uint64       `json:"max_wasm_code_size" yaml:"max_wasm_code_size"`
	GasMultiplier                uint64       `json:"gas_multiplier" yaml:"gas_multiplier"`
	InstanceCost                 uint64       `json:"instance_cost" yaml:"instance_cost"`
	CompileCost                  uint64       `json:"compile_cost" yaml:"compile_cost"`
	HumanizeCost                 uint64       `json:"humanize_cost" yaml:"humanize_cost"`
	CanonicalizeCost             uint64       `json:"canonicalize_cost" yaml:"canonicalize_cost"`
}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default wasm parameters
func DefaultParams() Params {
	return Params{
		UploadAccess:                 AllowEverybody,
		DefaultInstantiatePermission: Everybody,
		MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
		GasMultiplier:                DefaultGasMultiplier,
		InstanceCost:                 DefaultInstanceCost,
		CompileCost:                  DefaultCompileCost,
		HumanizeCost:                 DefaultHumanizeCost,
		CanonicalizeCost:             DefaultCanonicalCost,
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
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(ParamStoreKeyUploadAccess, &p.UploadAccess, validateAccessConfig),
		params.NewParamSetPair(ParamStoreKeyInstantiateAccess, &p.DefaultInstantiatePermission, validateAccessType),
		params.NewParamSetPair(ParamStoreKeyMaxWasmCodeSize, &p.MaxWasmCodeSize, validateMaxWasmCodeSize),
		params.NewParamSetPair(ParamStoreKeyGasMultiplier, &p.GasMultiplier, validateGasMultiplier),
		params.NewParamSetPair(ParamStoreKeyInstanceCost, &p.InstanceCost, validateInstanceCost),
		params.NewParamSetPair(ParamStoreKeyCompileCost, &p.CompileCost, validateCompileCost),
		params.NewParamSetPair(ParamStoreKeyHumanizeCost, &p.HumanizeCost, validateHumanizeCost),
		params.NewParamSetPair(ParamStoreKeyCanonicalCost, &p.CanonicalizeCost, validateCanonicalCost),
	}
}

// ValidateBasic performs basic validation on wasm parameters
func (p Params) ValidateBasic() error {
	if err := validateAccessType(p.DefaultInstantiatePermission); err != nil {
		return errors.Wrap(err, "instantiate default permission")
	}
	if err := validateAccessConfig(p.UploadAccess); err != nil {
		return errors.Wrap(err, "upload access")
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
	if err := validateHumanizeCost(p.HumanizeCost); err != nil {
		return errors.Wrap(err, "humanize cost")
	}
	if err := validateCanonicalCost(p.HumanizeCost); err != nil {
		return errors.Wrap(err, "canonical cost")
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
	v, ok := i.(AccessType)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == Undefined {
		return sdkerrors.Wrap(ErrEmpty, "type")
	}
	if _, ok := AllAccessTypes[v]; !ok {
		return sdkerrors.Wrapf(ErrInvalid, "unknown type: %q", v)
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

func validateHumanizeCost(i interface{}) error {
	a, ok := i.(uint64)
	if !ok {
		return sdkerrors.Wrapf(ErrInvalid, "type: %T", i)
	}
	if a == 0 {
		return sdkerrors.Wrap(ErrInvalid, "must be greater 0")
	}
	return nil
}

func validateCanonicalCost(i interface{}) error {
	a, ok := i.(uint64)
	if !ok {
		return sdkerrors.Wrapf(ErrInvalid, "type: %T", i)
	}
	if a == 0 {
		return sdkerrors.Wrap(ErrInvalid, "must be greater 0")
	}
	return nil
}

func (v AccessConfig) ValidateBasic() error {
	switch v.Type {
	case Undefined, "":
		return sdkerrors.Wrap(ErrEmpty, "type")
	case Nobody, Everybody:
		if len(v.Address) != 0 {
			return sdkerrors.Wrap(ErrInvalid, "address not allowed for this type")
		}
		return nil
	case OnlyAddress:
		return sdk.VerifyAddressFormat(v.Address)
	}
	return sdkerrors.Wrapf(ErrInvalid, "unknown type: %q", v.Type)
}

func (v AccessConfig) Allowed(actor sdk.AccAddress) bool {
	switch v.Type {
	case Nobody:
		return false
	case Everybody:
		return true
	case OnlyAddress:
		return v.Address.Equals(actor)
	default:
		panic("unknown type")
	}
}
