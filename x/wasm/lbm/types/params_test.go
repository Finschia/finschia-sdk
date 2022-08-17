package types

import (
	"encoding/json"
	wasmtypes "github.com/line/lbm-sdk/x/wasm/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/codec"
	codectypes "github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
)

func TestValidateParams(t *testing.T) {
	var (
		anyAddress     sdk.AccAddress = make([]byte, wasmtypes.ContractAddrLen)
		invalidAddress                = "invalid address"
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
				CodeUploadAccess:             wasmtypes.AllowNobody,
				InstantiateDefaultPermission: wasmtypes.AccessTypeNobody,
				GasMultiplier:                wasmtypes.DefaultGasMultiplier,
				InstanceCost:                 wasmtypes.DefaultInstanceCost,
				CompileCost:                  wasmtypes.DefaultCompileCost,
			},
		},
		"all good with everybody": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AllowEverybody,
				InstantiateDefaultPermission: wasmtypes.AccessTypeEverybody,
				GasMultiplier:                wasmtypes.DefaultGasMultiplier,
				InstanceCost:                 wasmtypes.DefaultInstanceCost,
				CompileCost:                  wasmtypes.DefaultCompileCost,
			},
		},
		"all good with only address": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AccessTypeOnlyAddress.With(anyAddress),
				InstantiateDefaultPermission: wasmtypes.AccessTypeOnlyAddress,
				GasMultiplier:                wasmtypes.DefaultGasMultiplier,
				InstanceCost:                 wasmtypes.DefaultInstanceCost,
				CompileCost:                  wasmtypes.DefaultCompileCost,
			},
		},
		"reject empty type in instantiate permission": {
			src: Params{
				CodeUploadAccess: wasmtypes.AllowNobody,
				GasMultiplier:    wasmtypes.DefaultGasMultiplier,
				InstanceCost:     wasmtypes.DefaultInstanceCost,
				CompileCost:      wasmtypes.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject unknown type in instantiate": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AllowNobody,
				InstantiateDefaultPermission: 1111,
				GasMultiplier:                wasmtypes.DefaultGasMultiplier,
				InstanceCost:                 wasmtypes.DefaultInstanceCost,
				CompileCost:                  wasmtypes.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject CodeUploadAccess invalid address in only address": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeOnlyAddress, Address: invalidAddress},
				InstantiateDefaultPermission: wasmtypes.AccessTypeOnlyAddress,
				GasMultiplier:                wasmtypes.DefaultGasMultiplier,
				InstanceCost:                 wasmtypes.DefaultInstanceCost,
				CompileCost:                  wasmtypes.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject CodeUploadAccess Everybody with obsolete address": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeEverybody, Address: anyAddress.String()},
				InstantiateDefaultPermission: wasmtypes.AccessTypeOnlyAddress,
				GasMultiplier:                wasmtypes.DefaultGasMultiplier,
				InstanceCost:                 wasmtypes.DefaultInstanceCost,
				CompileCost:                  wasmtypes.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject CodeUploadAccess Nobody with obsolete address": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeNobody, Address: anyAddress.String()},
				InstantiateDefaultPermission: wasmtypes.AccessTypeOnlyAddress,
				GasMultiplier:                wasmtypes.DefaultGasMultiplier,
				InstanceCost:                 wasmtypes.DefaultInstanceCost,
				CompileCost:                  wasmtypes.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject empty CodeUploadAccess": {
			src: Params{
				InstantiateDefaultPermission: wasmtypes.AccessTypeOnlyAddress,
				GasMultiplier:                wasmtypes.DefaultGasMultiplier,
				InstanceCost:                 wasmtypes.DefaultInstanceCost,
				CompileCost:                  wasmtypes.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject undefined permission in CodeUploadAccess": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeUnspecified},
				InstantiateDefaultPermission: wasmtypes.AccessTypeOnlyAddress,
				GasMultiplier:                wasmtypes.DefaultGasMultiplier,
				InstanceCost:                 wasmtypes.DefaultInstanceCost,
				CompileCost:                  wasmtypes.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject empty gas multiplier": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AllowNobody,
				InstantiateDefaultPermission: wasmtypes.AccessTypeNobody,
				InstanceCost:                 wasmtypes.DefaultInstanceCost,
				CompileCost:                  wasmtypes.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject empty instance cost": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AllowNobody,
				InstantiateDefaultPermission: wasmtypes.AccessTypeNobody,
				GasMultiplier:                wasmtypes.DefaultGasMultiplier,
				CompileCost:                  wasmtypes.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject empty compile cost": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AllowNobody,
				InstantiateDefaultPermission: wasmtypes.AccessTypeNobody,
				GasMultiplier:                wasmtypes.DefaultGasMultiplier,
				InstanceCost:                 wasmtypes.DefaultInstanceCost,
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
		src wasmtypes.AccessType
		exp string
	}{
		"Unspecified": {src: wasmtypes.AccessTypeUnspecified, exp: `"Unspecified"`},
		"Nobody":      {src: wasmtypes.AccessTypeNobody, exp: `"Nobody"`},
		"OnlyAddress": {src: wasmtypes.AccessTypeOnlyAddress, exp: `"OnlyAddress"`},
		"Everybody":   {src: wasmtypes.AccessTypeEverybody, exp: `"Everybody"`},
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
		exp wasmtypes.AccessType
	}{
		"Unspecified": {src: `"Unspecified"`, exp: wasmtypes.AccessTypeUnspecified},
		"Nobody":      {src: `"Nobody"`, exp: wasmtypes.AccessTypeNobody},
		"OnlyAddress": {src: `"OnlyAddress"`, exp: wasmtypes.AccessTypeOnlyAddress},
		"Everybody":   {src: `"Everybody"`, exp: wasmtypes.AccessTypeEverybody},
		"unknown":     {src: `""`, exp: wasmtypes.AccessTypeUnspecified},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			var got wasmtypes.AccessType
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
