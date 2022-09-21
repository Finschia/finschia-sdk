package keeper_test

import (
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
	minttypes "github.com/line/lbm-sdk/x/mint/types"
	stakingkeeper "github.com/line/lbm-sdk/x/staking/keeper"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
	"github.com/line/lbm-sdk/x/stakingplus"
	"github.com/line/lbm-sdk/x/stakingplus/keeper"
)

type KeeperTestSuite struct {
	suite.Suite
	ctx sdk.Context

	app       *simapp.SimApp
	keeper    stakingkeeper.Keeper
	msgServer stakingtypes.MsgServer

	stranger sdk.AccAddress
	grantee  sdk.AccAddress

	balance sdk.Int
}

func (s *KeeperTestSuite) SetupTest() {
	checkTx := false
	s.app = simapp.Setup(checkTx)
	s.ctx = s.app.BaseApp.NewContext(checkTx, ocproto.Header{})
	s.keeper = s.app.StakingKeeper

	s.msgServer = keeper.NewMsgServerImpl(s.keeper, s.app.FoundationKeeper)

	createAddress := func() sdk.AccAddress {
		return sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	s.stranger = createAddress()
	s.grantee = createAddress()

	s.balance = sdk.NewInt(1000000)
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
		err := s.app.BankKeeper.MintCoins(s.ctx, minterName, amount)
		s.Require().NoError(err)

		minter := s.app.AccountKeeper.GetModuleAccount(s.ctx, minterName).GetAddress()
		err = s.app.BankKeeper.SendCoins(s.ctx, minter, holder, amount)
		s.Require().NoError(err)
	}

	// allow Msg/CreateValidator
	s.app.FoundationKeeper.SetParams(s.ctx, &foundation.Params{
		Enabled:       true,
		FoundationTax: sdk.ZeroDec(),
	})
	err := s.app.FoundationKeeper.Grant(s.ctx, s.grantee, &stakingplus.CreateValidatorAuthorization{
		ValidatorAddress: sdk.ValAddress(s.grantee).String(),
	})
	s.Require().NoError(err)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
