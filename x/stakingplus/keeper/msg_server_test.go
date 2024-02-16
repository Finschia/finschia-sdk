package keeper_test

import (
	"github.com/cometbft/cometbft/crypto/secp256k1"
	"github.com/golang/mock/gomock"

	"cosmossdk.io/math"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (s *KeeperTestSuite) TestMsgCreateValidator() {
	// approve any msg
	s.stakingMsgServer.
		EXPECT().
		CreateValidator(gomock.Any(), gomock.Any()).
		Return(&stakingtypes.MsgCreateValidatorResponse{}, nil).
		AnyTimes()

	testCases := map[string]struct {
		valid     bool
		acceptRet error
	}{
		"valid request": {
			valid:     true,
			acceptRet: nil,
		},
		"no grant found": {
			valid:     false,
			acceptRet: sdkerrors.ErrUnauthorized,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			s.foundationKeeper.
				EXPECT().
				Accept(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(tc.acceptRet).
				Times(1)

			ctx, _ := s.ctx.CacheContext()

			pk := simtestutil.CreateTestPubKeys(1)[0]
			delegation := sdk.NewCoin(sdk.DefaultBondDenom, math.OneInt())
			val := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())
			valStr, err := s.valCodec.BytesToString(val)
			s.Require().NoError(err)
			req, err := stakingtypes.NewMsgCreateValidator(
				valStr,
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
