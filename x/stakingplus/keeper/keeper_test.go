package keeper_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/runtime"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/Finschia/finschia-sdk/x/stakingplus"
	"github.com/Finschia/finschia-sdk/x/stakingplus/keeper"
	"github.com/Finschia/finschia-sdk/x/stakingplus/testutil"
)

type KeeperTestSuite struct {
	suite.Suite
	ctx sdk.Context

	app           *runtime.App
	accountKeeper authkeeper.AccountKeeper
	bankKeeper    bankkeeper.Keeper
	stakingKeeper *stakingkeeper.Keeper
	msgServer     stakingtypes.MsgServer

	stranger sdk.AccAddress
	grantee  sdk.AccAddress

	balance math.Int
}

func (s *KeeperTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	foundationKeeper := testutil.NewMockFoundationKeeper(ctrl)

	app, err := simtestutil.Setup(
		depinject.Configs(
			testutil.AppConfig,
			depinject.Supply(log.NewNopLogger()),
		),
		&s.accountKeeper,
		&s.bankKeeper,
		&s.stakingKeeper,
	)
	s.Require().NoError(err)

	s.app = app
	s.ctx = s.app.BaseApp.NewContext(false)

	s.msgServer = keeper.NewMsgServerImpl(s.stakingKeeper, foundationKeeper)

	createAddress := func() sdk.AccAddress {
		return sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	s.stranger = createAddress()
	s.grantee = createAddress()

	s.balance = math.NewInt(1000000)
	holders := []sdk.AccAddress{
		s.stranger,
		s.grantee,
	}
	for _, holder := range holders {
		amount := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance))

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
	foundationKeeper.
		EXPECT().
		Accept(gomock.Any(), s.grantee, NewCreateValidatorAuthorizationMatcher(s.grantee)).
		Return(nil)
	foundationKeeper.
		EXPECT().
		Accept(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sdkerrors.ErrUnauthorized)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

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
