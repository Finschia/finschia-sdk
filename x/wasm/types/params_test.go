package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/codec"
	codectypes "github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
)

func TestValidateParams(t *testing.T) {
	var (
		anyAddress     = sdk.BytesToAccAddress(make([]byte, ContractAddrLen))
		invalidAddress = "invalid address"
	)

	specs := map[string]struct {
		src    Params
		expErr bool
	}{
		"all good with defaults": {
			src: DefaultParams(),
		},
		"all good with nobody": {
			src: Params{
				CodeUploadAccess:             AllowNobody,
				InstantiateDefaultPermission: AccessTypeNobody,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				GasMultiplier:                DefaultGasMultiplier,
				InstanceCost:                 DefaultInstanceCost,
				CompileCost:                  DefaultCompileCost,
			},
		},
		"all good with everybody": {
			src: Params{
				CodeUploadAccess:             AllowEverybody,
				InstantiateDefaultPermission: AccessTypeEverybody,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				GasMultiplier:                DefaultGasMultiplier,
				InstanceCost:                 DefaultInstanceCost,
				CompileCost:                  DefaultCompileCost,
			},
		},
		"all good with only address": {
			src: Params{
				CodeUploadAccess:             AccessTypeOnlyAddress.With(anyAddress),
				InstantiateDefaultPermission: AccessTypeOnlyAddress,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				GasMultiplier:                DefaultGasMultiplier,
				InstanceCost:                 DefaultInstanceCost,
				CompileCost:                  DefaultCompileCost,
			},
		},
		"reject empty type in instantiate permission": {
			src: Params{
				CodeUploadAccess: AllowNobody,
				MaxWasmCodeSize:  DefaultMaxWasmCodeSize,
				GasMultiplier:    DefaultGasMultiplier,
				InstanceCost:     DefaultInstanceCost,
				CompileCost:      DefaultCompileCost,
			},
			expErr: true,
		},
		"reject unknown type in instantiate": {
			src: Params{
				CodeUploadAccess:             AllowNobody,
				InstantiateDefaultPermission: 1111,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				GasMultiplier:                DefaultGasMultiplier,
				InstanceCost:                 DefaultInstanceCost,
				CompileCost:                  DefaultCompileCost,
			},
			expErr: true,
		},
		"reject CodeUploadAccess invalid address in only address": {
			src: Params{
				CodeUploadAccess:             AccessConfig{Permission: AccessTypeOnlyAddress, Address: invalidAddress},
				InstantiateDefaultPermission: AccessTypeOnlyAddress,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				GasMultiplier:                DefaultGasMultiplier,
				InstanceCost:                 DefaultInstanceCost,
				CompileCost:                  DefaultCompileCost,
			},
			expErr: true,
		},
		"reject CodeUploadAccess Everybody with obsolete address": {
			src: Params{
				CodeUploadAccess:             AccessConfig{Permission: AccessTypeEverybody, Address: anyAddress.String()},
				InstantiateDefaultPermission: AccessTypeOnlyAddress,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				GasMultiplier:                DefaultGasMultiplier,
				InstanceCost:                 DefaultInstanceCost,
				CompileCost:                  DefaultCompileCost,
			},
			expErr: true,
		},
		"reject CodeUploadAccess Nobody with obsolete address": {
			src: Params{
				CodeUploadAccess:             AccessConfig{Permission: AccessTypeNobody, Address: anyAddress.String()},
				InstantiateDefaultPermission: AccessTypeOnlyAddress,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				GasMultiplier:                DefaultGasMultiplier,
				InstanceCost:                 DefaultInstanceCost,
				CompileCost:                  DefaultCompileCost,
			},
			expErr: true,
		},
		"reject empty CodeUploadAccess": {
			src: Params{
				InstantiateDefaultPermission: AccessTypeOnlyAddress,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				GasMultiplier:                DefaultGasMultiplier,
				InstanceCost:                 DefaultInstanceCost,
				CompileCost:                  DefaultCompileCost,
			},
			expErr: true,
		},
		"reject undefined permission in CodeUploadAccess": {
			src: Params{
				CodeUploadAccess:             AccessConfig{Permission: AccessTypeUnspecified},
				InstantiateDefaultPermission: AccessTypeOnlyAddress,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				GasMultiplier:                DefaultGasMultiplier,
				InstanceCost:                 DefaultInstanceCost,
				CompileCost:                  DefaultCompileCost,
			},
			expErr: true,
		},
		"reject empty max wasm code size": {
			src: Params{
				CodeUploadAccess:             AllowNobody,
				InstantiateDefaultPermission: AccessTypeNobody,
				GasMultiplier:                DefaultGasMultiplier,
				InstanceCost:                 DefaultInstanceCost,
				CompileCost:                  DefaultCompileCost,
			},
			expErr: true,
		},
		"reject empty gas multiplier": {
			src: Params{
				CodeUploadAccess:             AllowNobody,
				InstantiateDefaultPermission: AccessTypeNobody,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				InstanceCost:                 DefaultInstanceCost,
				CompileCost:                  DefaultCompileCost,
			},
			expErr: true,
		},
		"reject empty instance cost": {
			src: Params{
				CodeUploadAccess:             AllowNobody,
				InstantiateDefaultPermission: AccessTypeNobody,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				GasMultiplier:                DefaultGasMultiplier,
				CompileCost:                  DefaultCompileCost,
			},
			expErr: true,
		},
		"reject empty compile cost": {
			src: Params{
				CodeUploadAccess:             AllowNobody,
				InstantiateDefaultPermission: AccessTypeNobody,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				GasMultiplier:                DefaultGasMultiplier,
				InstanceCost:                 DefaultInstanceCost,
			},
			expErr: true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			err := spec.src.ValidateBasic()
			if spec.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAccessTypeMarshalJson(t *testing.T) {
	specs := map[string]struct {
		src AccessType
		exp string
	}{
		"Unspecified": {src: AccessTypeUnspecified, exp: `"Unspecified"`},
		"Nobody":      {src: AccessTypeNobody, exp: `"Nobody"`},
		"OnlyAddress": {src: AccessTypeOnlyAddress, exp: `"OnlyAddress"`},
		"Everybody":   {src: AccessTypeEverybody, exp: `"Everybody"`},
		"unknown":     {src: 999, exp: `"Unspecified"`},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			got, err := json.Marshal(spec.src)
			require.NoError(t, err)
			assert.Equal(t, []byte(spec.exp), got)
		})
	}
}

func TestAccessTypeUnmarshalJson(t *testing.T) {
	specs := map[string]struct {
		src string
		exp AccessType
	}{
		"Unspecified": {src: `"Unspecified"`, exp: AccessTypeUnspecified},
		"Nobody":      {src: `"Nobody"`, exp: AccessTypeNobody},
		"OnlyAddress": {src: `"OnlyAddress"`, exp: AccessTypeOnlyAddress},
		"Everybody":   {src: `"Everybody"`, exp: AccessTypeEverybody},
		"unknown":     {src: `""`, exp: AccessTypeUnspecified},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			var got AccessType
			err := json.Unmarshal([]byte(spec.src), &got)
			require.NoError(t, err)
			assert.Equal(t, spec.exp, got)
		})
	}
}
func TestParamsUnmarshalJson(t *testing.T) {
	specs := map[string]struct {
		src string
		exp Params
	}{

		"defaults": {
			src: `{"code_upload_access": {"permission": "Everybody"},
				"instantiate_default_permission": "Everybody",
				"max_wasm_code_size": 1228800,
				"gas_multiplier": 140000000,
				"instance_cost": 60000,
				"compile_cost": 3}`,
			exp: DefaultParams(),
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			var val Params
			interfaceRegistry := codectypes.NewInterfaceRegistry()
			marshaler := codec.NewProtoCodec(interfaceRegistry)

			err := marshaler.UnmarshalJSON([]byte(spec.src), &val)
			require.NoError(t, err)
			assert.Equal(t, spec.exp, val)
		})
	}
}
