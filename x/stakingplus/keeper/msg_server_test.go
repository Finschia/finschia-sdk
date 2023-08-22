package keeper_test

import (
	"github.com/Finschia/finschia-rdk/simapp"
	sdk "github.com/Finschia/finschia-rdk/types"
	stakingtypes "github.com/Finschia/finschia-rdk/x/staking/types"
)

func (s *KeeperTestSuite) TestMsgCreateValidator() {
	testCases := map[string]struct {
		delegator sdk.AccAddress
		valid     bool
	}{
		"valid request": {
			delegator: s.grantee,
			valid:     true,
		},
		"no grant found": {
			delegator: s.stranger,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			pk := simapp.CreateTestPubKeys(1)[0]
			delegation := sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())
			req, err := stakingtypes.NewMsgCreateValidator(
				sdk.ValAddress(tc.delegator),
				pk,
				delegation,
				stakingtypes.Description{},
				stakingtypes.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
				delegation.Amount,
			)
			s.Require().NoError(err)

			res, err := s.msgServer.CreateValidator(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}
