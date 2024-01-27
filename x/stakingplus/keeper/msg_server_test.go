package keeper_test

import (
	"cosmossdk.io/math"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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

			pk := simtestutil.CreateTestPubKeys(1)[0]
			delegation := sdk.NewCoin(sdk.DefaultBondDenom, math.OneInt())
			req, err := stakingtypes.NewMsgCreateValidator(
				sdk.ValAddress(tc.delegator).String(),
				pk,
				delegation,
				stakingtypes.Description{
					Moniker: "Test Validator",
				},
				stakingtypes.NewCommissionRates(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
				delegation.Amount,
			)
			s.Require().NoError(err)

			res, err := s.msgServer.CreateValidator(ctx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}
