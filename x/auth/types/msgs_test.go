package types

import (
	"fmt"
	"testing"

	sdk "github.com/line/lfb-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgSendRoute(t *testing.T) {
	msg := NewMsgEmpty(sdk.AccAddress("from"))

	require.Equal(t, ModuleName, msg.Route())
	require.Equal(t, "empty", msg.Type())

	srvMsg := NewServiceMsgEmpty(sdk.AccAddress("from"))
	require.Equal(t, "/lfb.auth.v1beta1.Msg/Empty", srvMsg.Route())
	require.Equal(t, "/lfb.auth.v1beta1.Msg/Empty", srvMsg.Type())
}

func TestMsgSendValidation(t *testing.T) {
	addr1 := sdk.AccAddress("from________________")
	addrEmpty := sdk.AccAddress("")
	addrTooLong := sdk.AccAddress("Accidentally used 33 bytes pubkey")

	cases := []struct {
		expectedErr string // empty means no error expected
		msg         *MsgEmpty
	}{
		{"", NewMsgEmpty(addr1)}, // valid
		{"Invalid sender address (empty address string is not allowed): invalid address", NewMsgEmpty(addrEmpty)},
		{"Invalid sender address (incorrect address length (expected: 20, actual: 33)): invalid address", NewMsgEmpty(addrTooLong)},
	}

	for _, tc := range cases {
		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestMsgSendGetSignBytes(t *testing.T) {
	res := NewMsgEmpty(sdk.AccAddress("input")).GetSignBytes()

	expected := `{"type":"lfb-sdk/MsgEmpty","value":{"from_address":"link1d9h8qat5fnwd3e"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgSendGetSigners(t *testing.T) {
	res := NewMsgEmpty(sdk.AccAddress("input111111111111111")).GetSigners()
	require.Equal(t, fmt.Sprintf("%v", res), "[696E707574313131313131313131313131313131]")

	require.Panics(t, func() {
		NewMsgEmpty(sdk.AccAddress("input")).GetSigners()
	})
}
