package types

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/line/lbm-sdk/types"
)

func TestMsgSendRoute(t *testing.T) {
	msg := NewMsgEmpty(sdk.AccAddress("from"))

	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "empty", msg.Type())
}

func TestMsgSendValidation(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from________________"))
	addrEmpty := sdk.AccAddress("")
	addrLong := sdk.AccAddress([]byte("Accidentally used 33 bytes pubkey"))

	cases := []struct {
		expectedErr string // empty means no error expected
		msg         *MsgEmpty
	}{
		{"", NewMsgEmpty(addr1)}, // valid
		{"Invalid sender address (empty address string is not allowed): invalid address", NewMsgEmpty(addrEmpty)},
		{"", NewMsgEmpty(addrLong)},
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
	res := NewMsgEmpty(sdk.AccAddress([]byte("input"))).GetSignBytes()

	expected := `{"type":"cosmos-sdk/MsgEmpty","value":{"from_address":"link1d9h8qat5fnwd3e"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgSendGetSigners(t *testing.T) {
	res := NewMsgEmpty(sdk.AccAddress([]byte("input111111111111111"))).GetSigners()
	bytes, _ := sdk.AccAddressFromBech32(res[0].String())
	require.Equal(t, "696e707574313131313131313131313131313131", fmt.Sprintf("%v", hex.EncodeToString(bytes)))
}
