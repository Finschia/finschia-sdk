package types

import (
	"encoding/json"
	"testing"

	"github.com/line/lbm-sdk/v2/codec"
	codectypes "github.com/line/lbm-sdk/v2/codec/types"
	sdk "github.com/line/lbm-sdk/v2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateParams(t *testing.T) {
	var (
		anyAddress     sdk.AccAddress = make([]byte, sdk.AddrLen)
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
				CodeUploadAccess:             AllowNobody,
				InstantiateDefaultPermission: AccessTypeNobody,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				ContractStatusAccess:         AllowNobody,
			},
		},
		"all good with everybody": {
			src: Params{
				CodeUploadAccess:             AllowEverybody,
				InstantiateDefaultPermission: AccessTypeEverybody,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				ContractStatusAccess:         AllowEverybody,
			},
		},
		"all good with only address": {
			src: Params{
				CodeUploadAccess:             AccessTypeOnlyAddress.With(anyAddress),
				InstantiateDefaultPermission: AccessTypeOnlyAddress,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				ContractStatusAccess:         AccessTypeOnlyAddress.With(anyAddress),
			},
		},
		"reject empty type in instantiate permission": {
			src: Params{
				CodeUploadAccess:     AllowNobody,
				MaxWasmCodeSize:      DefaultMaxWasmCodeSize,
				ContractStatusAccess: DefaultContractStatusAccess,
			},
			expErr: true,
		},
		"reject unknown type in instantiate": {
			src: Params{
				CodeUploadAccess:             AllowNobody,
				InstantiateDefaultPermission: 1111,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				ContractStatusAccess:         DefaultContractStatusAccess,
			},
			expErr: true,
		},
		"reject CodeUploadAccess invalid address in only address": {
			src: Params{
				CodeUploadAccess:             AccessConfig{Permission: AccessTypeOnlyAddress, Address: invalidAddress},
				InstantiateDefaultPermission: AccessTypeOnlyAddress,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				ContractStatusAccess:         DefaultContractStatusAccess,
			},
			expErr: true,
		},
		"reject CodeUploadAccess Everybody with obsolete address": {
			src: Params{
				CodeUploadAccess:             AccessConfig{Permission: AccessTypeEverybody, Address: anyAddress.String()},
				InstantiateDefaultPermission: AccessTypeOnlyAddress,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				ContractStatusAccess:         DefaultContractStatusAccess,
			},
			expErr: true,
		},
		"reject CodeUploadAccess Nobody with obsolete address": {
			src: Params{
				CodeUploadAccess:             AccessConfig{Permission: AccessTypeNobody, Address: anyAddress.String()},
				InstantiateDefaultPermission: AccessTypeOnlyAddress,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				ContractStatusAccess:         DefaultContractStatusAccess,
			},
			expErr: true,
		},
		"reject empty CodeUploadAccess": {
			src: Params{
				InstantiateDefaultPermission: AccessTypeOnlyAddress,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				ContractStatusAccess:         DefaultContractStatusAccess,
			},
			expErr: true,
		},
		"reject undefined permission in CodeUploadAccess": {
			src: Params{
				CodeUploadAccess:             AccessConfig{Permission: AccessTypeUnspecified},
				InstantiateDefaultPermission: AccessTypeOnlyAddress,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				ContractStatusAccess:         DefaultContractStatusAccess,
			},
			expErr: true,
		},
		"reject empty max wasm code size": {
			src: Params{
				CodeUploadAccess:             AllowNobody,
				InstantiateDefaultPermission: AccessTypeNobody,
				ContractStatusAccess:         DefaultContractStatusAccess,
			},
			expErr: true,
		},
		"reject ContractStatusAccess invalid address in only address": {
			src: Params{
				CodeUploadAccess:             DefaultUploadAccess,
				InstantiateDefaultPermission: AccessTypeOnlyAddress,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				ContractStatusAccess:         AccessConfig{Permission: AccessTypeOnlyAddress, Address: invalidAddress},
			},
			expErr: true,
		},
		"reject ContractStatusAccess Everybody with obsolete address": {
			src: Params{
				CodeUploadAccess:             DefaultUploadAccess,
				InstantiateDefaultPermission: AccessTypeOnlyAddress,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				ContractStatusAccess:         AccessConfig{Permission: AccessTypeEverybody, Address: anyAddress.String()},
			},
			expErr: true,
		},
		"reject ContractStatusAccess Nobody with obsolete address": {
			src: Params{
				CodeUploadAccess:             DefaultUploadAccess,
				InstantiateDefaultPermission: AccessTypeOnlyAddress,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				ContractStatusAccess:         AccessConfig{Permission: AccessTypeNobody, Address: anyAddress.String()},
			},
			expErr: true,
		},
		"reject empty ContractStatusAccess": {
			src: Params{
				CodeUploadAccess:             DefaultUploadAccess,
				InstantiateDefaultPermission: AccessTypeOnlyAddress,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
			},
			expErr: true,
		},
		"reject undefined permission in ContractStatusAccess": {
			src: Params{
				CodeUploadAccess:             DefaultUploadAccess,
				InstantiateDefaultPermission: AccessTypeOnlyAddress,
				MaxWasmCodeSize:              DefaultMaxWasmCodeSize,
				ContractStatusAccess:         AccessConfig{Permission: AccessTypeUnspecified},
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
				"max_wasm_code_size": 614400,
				"contract_status_access": {"permission": "Nobody"}}`,
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
