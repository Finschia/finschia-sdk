package foundation_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

func TestCreateValidatorAuthorization(t *testing.T) {
	validVal := "valid_val"
	invalidVal := "invalid_val"

	testCases := map[string]struct {
		msg    sdk.Msg
		valid  bool
		accept bool
	}{
		"valid validator": {
			msg: &stakingtypes.MsgCreateValidator{
				ValidatorAddress: validVal,
			},
			valid:  true,
			accept: true,
		},
		"invalid validator": {
			msg: &stakingtypes.MsgCreateValidator{
				ValidatorAddress: invalidVal,
			},
			valid: false,
		},
		"msg mismatch": {
			msg: &foundation.MsgVote{},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			authorization := &foundation.CreateValidatorAuthorization{
				ValidatorAddress: validVal,
			}

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
