package types

import (
	sdk "github.com/line/lbm-sdk/types"
	wasmTypes "github.com/line/lbm-sdk/x/wasm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStoreCodeAndInstantiateContractValidation(t *testing.T) {
	bad, err := sdk.AccAddressFromHex("012345")
	require.NoError(t, err)
	badAddress := bad.String()
	require.NoError(t, err)
	// proper address size
	goodAddress := sdk.AccAddress(make([]byte, wasmTypes.ContractAddrLen)).String()
	sdk.GetConfig().SetAddressVerifier(wasmTypes.VerifyAddressLen())

	cases := map[string]struct {
		msg   MsgStoreCodeAndInstantiateContract
		valid bool
	}{
		"empty": {
			msg:   MsgStoreCodeAndInstantiateContract{},
			valid: false,
		},
		"correct minimal": {
			msg: MsgStoreCodeAndInstantiateContract{
				Sender:       goodAddress,
				WASMByteCode: []byte("foo"),
				Label:        "foo",
				Msg:          []byte("{}"),
			},
			valid: true,
		},
		"missing code": {
			msg: MsgStoreCodeAndInstantiateContract{
				Sender: goodAddress,
				Label:  "foo",
				Msg:    []byte("{}"),
			},
			valid: false,
		},
		"missing label": {
			msg: MsgStoreCodeAndInstantiateContract{
				Sender:       goodAddress,
				WASMByteCode: []byte("foo"),
				Msg:          []byte("{}"),
			},
			valid: false,
		},
		"missing init message": {
			msg: MsgStoreCodeAndInstantiateContract{
				Sender:       goodAddress,
				WASMByteCode: []byte("foo"),
				Label:        "foo",
			},
			valid: false,
		},
		"correct maximal": {
			msg: MsgStoreCodeAndInstantiateContract{
				Sender:       goodAddress,
				WASMByteCode: []byte("foo"),
				Label:        "foo",
				Msg:          []byte(`{"some": "data"}`),
				Funds:        sdk.Coins{sdk.Coin{Denom: "foobar", Amount: sdk.NewInt(200)}},
			},
			valid: true,
		},
		"invalid InstantiatePermission": {
			msg: MsgStoreCodeAndInstantiateContract{
				Sender:                goodAddress,
				WASMByteCode:          []byte("foo"),
				InstantiatePermission: &wasmTypes.AccessConfig{Permission: wasmTypes.AccessTypeOnlyAddress, Address: badAddress},
				Label:                 "foo",
				Msg:                   []byte(`{"some": "data"}`),
				Funds:                 sdk.Coins{sdk.Coin{Denom: "foobar", Amount: sdk.NewInt(200)}},
			},
			valid: false,
		},
		"negative funds": {
			msg: MsgStoreCodeAndInstantiateContract{
				Sender:       goodAddress,
				WASMByteCode: []byte("foo"),
				Msg:          []byte(`{"some": "data"}`),
				// we cannot use sdk.NewCoin() constructors as they panic on creating invalid data (before we can test)
				Funds: sdk.Coins{sdk.Coin{Denom: "foobar", Amount: sdk.NewInt(-200)}},
			},
			valid: false,
		},
		"non json init msg": {
			msg: MsgStoreCodeAndInstantiateContract{
				Sender:       goodAddress,
				WASMByteCode: []byte("foo"),
				Label:        "foo",
				Msg:          []byte("invalid-json"),
			},
			valid: false,
		},
		"bad sender minimal": {
			msg: MsgStoreCodeAndInstantiateContract{
				Sender:       badAddress,
				WASMByteCode: []byte("foo"),
				Label:        "foo",
				Msg:          []byte("{}"),
			},
			valid: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
