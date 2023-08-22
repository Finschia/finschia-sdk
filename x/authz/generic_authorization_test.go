package authz_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-rdk/x/authz"
	banktypes "github.com/Finschia/finschia-rdk/x/bank/types"
)

func TestGenericAuthorization(t *testing.T) {
	t.Log("verify ValidateBasic returns nil for service msg")
	a := authz.NewGenericAuthorization(banktypes.SendAuthorization{}.MsgTypeURL())
	require.NoError(t, a.ValidateBasic())
	require.Equal(t, banktypes.SendAuthorization{}.MsgTypeURL(), a.Msg)
}
