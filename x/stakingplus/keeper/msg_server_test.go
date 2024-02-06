package keeper_test

import (
	"fmt"

	"github.com/cometbft/cometbft/crypto/secp256k1"
	"github.com/golang/mock/gomock"

	"cosmossdk.io/math"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/Finschia/finschia-sdk/x/stakingplus"
)

func (s *KeeperTestSuite) TestMsgCreateValidator() {
	stranger := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	grantee := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	balance := math.NewInt(1000000)

	for _, holder := range []sdk.AccAddress{stranger, grantee} {
		amount := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, balance))

		// using minttypes here introduces dependency on x/mint
		// the work around would be registering a new module account on this suite
		// because x/bank already has dependency on x/mint, and we must have dependency
		// on x/bank, it's OK to use x/mint here.
		minterName := minttypes.ModuleName
		err := s.bankKeeper.MintCoins(s.ctx, minterName, amount)
		s.Require().NoError(err)

		minter := s.accountKeeper.GetModuleAccount(s.ctx, minterName).GetAddress()
		err = s.bankKeeper.SendCoins(s.ctx, minter, holder, amount)
		s.Require().NoError(err)
	}

	// approve Msg/CreateValidator to grantee
	s.foundationKeeper.
		EXPECT().
		Accept(gomock.Any(), grantee, NewCreateValidatorAuthorizationMatcher(grantee)).
		Return(nil)
	s.foundationKeeper.
		EXPECT().
		Accept(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sdkerrors.ErrUnauthorized)

	testCases := map[string]struct {
		delegator sdk.AccAddress
		valid     bool
	}{
		"valid request": {
			delegator: grantee,
			valid:     true,
		},
		"no grant found": {
			delegator: stranger,
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

var _ gomock.Matcher = (*CreateValidatorAuthorizationMatcher)(nil)

type CreateValidatorAuthorizationMatcher struct {
	authz stakingplus.CreateValidatorAuthorization
}

func NewCreateValidatorAuthorizationMatcher(grantee sdk.AccAddress) *CreateValidatorAuthorizationMatcher {
	return &CreateValidatorAuthorizationMatcher{
		authz: stakingplus.CreateValidatorAuthorization{
			ValidatorAddress: sdk.ValAddress(grantee).String(),
		},
	}
}

func (c CreateValidatorAuthorizationMatcher) Matches(x interface{}) bool {
	msg, ok := x.(sdk.Msg)
	if !ok {
		return false
	}

	resp, err := c.authz.Accept(sdk.Context{}, msg)
	return resp.Accept && (err == nil)
}

func (c CreateValidatorAuthorizationMatcher) String() string {
	return fmt.Sprintf("grants %s to %s", c.authz.MsgTypeURL(), c.authz.ValidatorAddress)
}
