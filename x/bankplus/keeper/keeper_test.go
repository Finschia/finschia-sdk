package keeper_test

import (
	"testing"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	codectestutil "github.com/cosmos/cosmos-sdk/codec/testutil"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/Finschia/finschia-sdk/x/bankplus/keeper"
)

const (
	initialPower = int64(100)
	holder       = "holder"
	blocker      = "blocker"
)

var (
	holderAcc  = authtypes.NewEmptyModuleAccount(holder)
	blockedAcc = authtypes.NewEmptyModuleAccount(blocker)
	mintAcc    = authtypes.NewEmptyModuleAccount(minttypes.ModuleName, authtypes.Minter)
	burnerAcc  = authtypes.NewEmptyModuleAccount(authtypes.Burner, authtypes.Burner)

	initTokens = sdk.TokensFromConsensusPower(initialPower, sdk.DefaultPowerReduction)
	initCoins  = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens))
)

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, &IntegrationTestSuite{})
}

type IntegrationTestSuite struct {
	suite.Suite
	ctx          sdk.Context
	mockCtrl     *gomock.Controller
	bkPlusKeeper keeper.BaseKeeper
	authKeeper   *banktestutil.MockAccountKeeper
	baseAcc      *authtypes.BaseAccount
	codec        codec.Codec
	storeService store.KVStoreService
}

func (s *IntegrationTestSuite) SetupTest() {
	key := storetypes.NewKVStoreKey(types.StoreKey)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	s.ctx = testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: cmttime.Now()})

	encCfg := moduletestutil.MakeTestEncodingConfig()
	encCfg.Codec = codectestutil.CodecOptions{
		AccAddressPrefix: "link",
		ValAddressPrefix: "linkvaloper",
	}.NewCodec()
	s.codec = encCfg.Codec

	storeService := runtime.NewKVStoreService(key)
	s.storeService = storeService

	ctrl := gomock.NewController(s.T())
	s.mockCtrl = ctrl
	authKeeper := banktestutil.NewMockAccountKeeper(ctrl)
	s.authKeeper = authKeeper
	newAcc := sdk.AccAddress("valid")
	s.baseAcc = authtypes.NewBaseAccountWithAddress(newAcc)
	s.authKeeper.EXPECT().GetAccount(gomock.Any(), newAcc).Return(nil).AnyTimes()
	s.authKeeper.EXPECT().AddressCodec().Return(address.NewBech32Codec("cosmos")).AnyTimes()
	s.authKeeper.EXPECT().NewAccountWithAddress(gomock.Any(), newAcc).Return(s.baseAcc).AnyTimes()
	s.authKeeper.EXPECT().GetModuleAddress("").Return(nil).AnyTimes()
	s.authKeeper.EXPECT().GetModuleAccount(gomock.Any(), "").Return(nil).AnyTimes()
	s.authKeeper.EXPECT().HasAccount(s.ctx, s.baseAcc.GetAddress()).Return(true).AnyTimes()
	s.mockAuthKeeperFor(holderAcc)
	s.mockAuthKeeperFor(blockedAcc)
	s.mockAuthKeeperFor(mintAcc)
	s.mockAuthKeeperFor(burnerAcc)

	s.bkPlusKeeper = keeper.NewBaseKeeper(
		encCfg.Codec,
		storeService,
		s.authKeeper,
		map[string]bool{},
		true,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		log.NewNopLogger(),
	)
}

func (s *IntegrationTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func getCoinsByName(ctx sdk.Context, bk keeper.Keeper, ak types.AccountKeeper, moduleName string) sdk.Coins {
	moduleAddress := ak.GetModuleAddress(moduleName)
	macc := ak.GetAccount(ctx, moduleAddress)
	if macc == nil {
		return sdk.NewCoins()
	}

	return bk.GetAllBalances(ctx, macc.GetAddress())
}

func (s *IntegrationTestSuite) mockAuthKeeperFor(acc *authtypes.ModuleAccount) {
	s.authKeeper.EXPECT().GetModuleAccount(s.ctx, acc.GetName()).Return(acc).AnyTimes()
	s.authKeeper.EXPECT().GetModuleAddress(acc.GetName()).Return(acc.GetAddress()).AnyTimes()
	s.authKeeper.EXPECT().GetAccount(s.ctx, acc.GetAddress()).Return(acc).AnyTimes()
	s.authKeeper.EXPECT().HasAccount(s.ctx, acc.GetAddress()).Return(true).AnyTimes()
}

func (s *IntegrationTestSuite) TestSupply_SendCoins() {
	s.mintInitialBalances()
	s.verifySendCoinsWithInvalidAddr()
	s.errorForInsufficientBalance()
	s.verifySendModuleToModule()
	s.verifySendModuleToAccount()
	s.verifySendAccountToModule()
}

func (s *IntegrationTestSuite) mintInitialBalances() {
	s.Require().NoError(s.bkPlusKeeper.MintCoins(s.ctx, minttypes.ModuleName, initCoins))
	s.Require().NoError(s.bkPlusKeeper.SendCoinsFromModuleToAccount(s.ctx, minttypes.ModuleName, holderAcc.GetAddress(), initCoins))
	s.Require().NoError(s.bkPlusKeeper.MintCoins(s.ctx, minttypes.ModuleName, initCoins))
	s.Require().Equal(initCoins, s.bkPlusKeeper.GetAllBalances(s.ctx, holderAcc.GetAddress()))
}

func (s *IntegrationTestSuite) verifySendCoinsWithInvalidAddr() {
	s.Require().Panics(func() {
		s.bkPlusKeeper.SendCoinsFromModuleToModule(s.ctx, "", holderAcc.GetName(), initCoins) // nolint:errcheck // No need
	})
	s.Require().Panics(func() {
		s.bkPlusKeeper.SendCoinsFromModuleToModule(s.ctx, authtypes.Burner, "", initCoins) // nolint:errcheck // No need
	})
	s.Require().Panics(func() {
		s.bkPlusKeeper.SendCoinsFromModuleToAccount(s.ctx, "", s.baseAcc.GetAddress(), initCoins) // nolint:errcheck // No need
	})
}

func (s *IntegrationTestSuite) errorForInsufficientBalance() {
	s.Require().Error(s.bkPlusKeeper.SendCoinsFromModuleToAccount(s.ctx, holderAcc.GetName(), s.baseAcc.GetAddress(), initCoins.Add(initCoins...)))
}

func (s *IntegrationTestSuite) verifySendModuleToModule() {
	s.Require().NoError(s.bkPlusKeeper.SendCoinsFromModuleToModule(s.ctx, holderAcc.GetName(), authtypes.Burner, initCoins))
	s.Require().Equal(sdk.NewCoins().String(), getCoinsByName(s.ctx, s.bkPlusKeeper, s.authKeeper, holderAcc.GetName()).String())
	s.Require().Equal(initCoins, getCoinsByName(s.ctx, s.bkPlusKeeper, s.authKeeper, authtypes.Burner))
}

func (s *IntegrationTestSuite) verifySendModuleToAccount() {
	s.Require().NoError(s.bkPlusKeeper.SendCoinsFromModuleToAccount(s.ctx, authtypes.Burner, s.baseAcc.GetAddress(), initCoins))
	s.Require().Equal(sdk.NewCoins().String(), getCoinsByName(s.ctx, s.bkPlusKeeper, s.authKeeper, authtypes.Burner).String())
	s.Require().Equal(initCoins, s.bkPlusKeeper.GetAllBalances(s.ctx, s.baseAcc.GetAddress()))
}

func (s *IntegrationTestSuite) verifySendAccountToModule() {
	s.Require().NoError(s.bkPlusKeeper.SendCoinsFromAccountToModule(s.ctx, s.baseAcc.GetAddress(), authtypes.Burner, initCoins))
	s.Require().Equal(sdk.NewCoins().String(), s.bkPlusKeeper.GetAllBalances(s.ctx, s.baseAcc.GetAddress()).String())
	s.Require().Equal(initCoins, getCoinsByName(s.ctx, s.bkPlusKeeper, s.authKeeper, authtypes.Burner))
}

func (s *IntegrationTestSuite) TestInactiveAddrOfSendCoins() {
	s.mintInitialBalances()
	s.verifySendToBlockedAddr()
	s.verifySendToBlockedAddrAfterRemoveIt()
}

func (s *IntegrationTestSuite) verifySendToBlockedAddr() {
	s.Require().False(s.bkPlusKeeper.IsInactiveAddr(blockedAcc.GetAddress()))
	s.bkPlusKeeper.AddToInactiveAddr(s.ctx, blockedAcc.GetAddress())
	s.Require().True(s.bkPlusKeeper.IsInactiveAddr(blockedAcc.GetAddress()))

	err := s.bkPlusKeeper.SendCoins(s.ctx, holderAcc.GetAddress(), blockedAcc.GetAddress(), initCoins)
	s.Require().Contains(err.Error(), "is not allowed to receive funds")
	s.Require().Equal(initCoins, s.bkPlusKeeper.GetAllBalances(s.ctx, holderAcc.GetAddress()))
}

func (s *IntegrationTestSuite) verifySendToBlockedAddrAfterRemoveIt() {
	s.bkPlusKeeper.DeleteFromInactiveAddr(s.ctx, blockedAcc.GetAddress())
	s.Require().False(s.bkPlusKeeper.IsInactiveAddr(blockedAcc.GetAddress()))
	s.Require().NoError(s.bkPlusKeeper.SendCoins(s.ctx, holderAcc.GetAddress(), blockedAcc.GetAddress(), initCoins))
	s.Require().Equal(sdk.NewCoins().String(), s.bkPlusKeeper.GetAllBalances(s.ctx, holderAcc.GetAddress()).String())
}

func (s *IntegrationTestSuite) TestInitializeBankPlus() {
	newKeeper := keeper.NewBaseKeeper(
		s.codec,
		s.storeService,
		s.authKeeper,
		map[string]bool{},
		true,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		log.NewNopLogger(),
	)
	newKeeper.AddToInactiveAddr(s.ctx, blockedAcc.GetAddress())
	s.Require().True(newKeeper.IsInactiveAddr(blockedAcc.GetAddress()))

	anotherNewKeeper := keeper.NewBaseKeeper(
		s.codec,
		s.storeService,
		s.authKeeper,
		map[string]bool{},
		true,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		log.NewNopLogger(),
	)
	anotherNewKeeper.InitializeBankPlus(s.ctx)
	s.Require().True(anotherNewKeeper.IsInactiveAddr(blockedAcc.GetAddress()))
}

func (s *IntegrationTestSuite) TestSendCoinsFromModuleToAccount_Blacklist() {
	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	newKeeper := keeper.NewBaseKeeper(
		s.codec,
		s.storeService,
		s.authKeeper,
		map[string]bool{addr1.String(): true},
		true,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		log.NewNopLogger(),
	)

	s.Require().NoError(newKeeper.MintCoins(s.ctx, minttypes.ModuleName, initCoins))
	s.Require().Error(newKeeper.SendCoinsFromModuleToAccount(s.ctx, minttypes.ModuleName, addr1, initCoins))
}

func (s *IntegrationTestSuite) TestInputOutputCoins() {
	// set initial balances
	s.Require().NoError(s.bkPlusKeeper.MintCoins(s.ctx, minttypes.ModuleName, initCoins))
	s.Require().NoError(s.bkPlusKeeper.SendCoinsFromModuleToAccount(s.ctx, minttypes.ModuleName, s.baseAcc.GetAddress(), initCoins))

	input := types.NewInput(s.baseAcc.GetAddress(), initCoins)
	outputs := []types.Output{types.NewOutput(burnerAcc.GetAddress(), initCoins)}

	targetKeeper := func(isDeact bool) keeper.BaseKeeper {
		return keeper.NewBaseKeeper(
			s.codec,
			s.storeService,
			s.authKeeper,
			make(map[string]bool),
			isDeact,
			authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			log.NewNopLogger(),
		)
	}
	tcs := map[string]struct {
		deactMultiSend bool
		err            error
	}{
		"MultiSend Off": {
			true,
			sdkerrors.ErrNotSupported.Wrap("MultiSend was deactivated"),
		},
		"MultiSend On": {
			false,
			nil,
		},
	}

	for name, tc := range tcs {
		tc := tc
		s.T().Run(name, func(t *testing.T) {
			if tc.err != nil {
				s.EqualError(targetKeeper(tc.deactMultiSend).InputOutputCoins(s.ctx, input, outputs), tc.err.Error())
			} else {
				s.Require().NoError(targetKeeper(tc.deactMultiSend).InputOutputCoins(s.ctx, input, outputs))
			}
		})
	}
}
