package foundation_test

import (
	"testing"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/foundation"
	"github.com/stretchr/testify/require"
)

func TestReceiveFromTreasuryAuthorization(t *testing.T) {
	testCases := map[string]struct {
		msg    sdk.Msg
		valid  bool
		accept bool
	}{
		"valid": {
			msg:    &foundation.MsgWithdrawFromTreasury{},
			valid:  true,
			accept: true,
		},
		"msg mismatch": {
			msg: &foundation.MsgVote{},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			authorization := &foundation.ReceiveFromTreasuryAuthorization{}

			resp, err := authorization.Accept(sdk.Context{}, tc.msg)
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, tc.accept, resp.Accept)
		})
	}
}
