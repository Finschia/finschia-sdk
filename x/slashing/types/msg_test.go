package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/line/lfb-sdk/types"
)

func TestMsgUnjailGetSignBytes(t *testing.T) {
	addr := sdk.BytesToValAddress([]byte("abcd"))
	msg := NewMsgUnjail(addr)
	bytes := msg.GetSignBytes()
	require.Equal(
		t,
		`{"type":"lfb-sdk/MsgUnjail","value":{"address":"linkvaloper1v93xxeqn4h65f"}}`,
		string(bytes),
	)
}
