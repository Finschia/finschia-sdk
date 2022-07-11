package foundation_test

import (
	"testing"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
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
		authorization := &foundation.ReceiveFromTreasuryAuthorization{}

		resp, err := authorization.Accept(sdk.Context{}, tc.msg)
		if !tc.valid {
			require.Error(t, err, name)
			continue
		}
		require.NoError(t, err, name)

		require.Equal(t, tc.accept, resp.Accept)
	}
}
