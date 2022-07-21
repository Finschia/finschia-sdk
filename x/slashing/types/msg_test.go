package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/line/lbm-sdk/types"
)

func TestMsgUnjailGetSignBytes(t *testing.T) {
	addr := sdk.ValAddress("abcd")
	msg := NewMsgUnjail(addr)
	bytes := msg.GetSignBytes()
	require.Equal(
		t,
		`{"type":"cosmos-sdk/MsgUnjail","value":{"address":"linkvaloper1v93xxeqn4h65f"}}`,
		string(bytes),
	)
}
