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
	stranger := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	grantee := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	// approve Msg/CreateValidator to grantee
	s.foundationKeeper.
		EXPECT().
		Accept(gomock.Any(), grantee, gomock.Any()).
		Return(nil)
	s.foundationKeeper.
		EXPECT().
		Accept(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sdkerrors.ErrUnauthorized)

	// approve any msg
	s.stakingMsgServer.
		EXPECT().
		CreateValidator(gomock.Any(), gomock.Any()).
		Return(&stakingtypes.MsgCreateValidatorResponse{}, nil)

	testCases := map[string]struct {
		validator sdk.AccAddress
		valid     bool
	}{
		"valid request": {
			validator: grantee,
			valid:     true,
		},
		"no grant found": {
			validator: stranger,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			pk := simtestutil.CreateTestPubKeys(1)[0]
			delegation := sdk.NewCoin(sdk.DefaultBondDenom, math.OneInt())
			val, err := s.valCodec.BytesToString(tc.validator)
			s.Require().NoError(err)
			req, err := stakingtypes.NewMsgCreateValidator(
				val,
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
