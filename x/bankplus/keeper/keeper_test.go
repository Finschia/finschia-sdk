package keeper

import (
	"testing"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	coreaddress "cosmossdk.io/core/address"
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
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

const (
	initialPower = int64(100)
	holder       = "holder"
	blocker      = "blocker"
)

var (
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
	cut          BaseKeeper
	authKeeper   *banktestutil.MockAccountKeeper
	baseAcc      *authtypes.BaseAccount
	cdc          codec.Codec
	storeService store.KVStoreService
	addrCdc      coreaddress.Codec
	holderAcc    *authtypes.ModuleAccount
	blockedAcc   *authtypes.ModuleAccount
	mintAcc      *authtypes.ModuleAccount
	burnerAcc    *authtypes.ModuleAccount
}

func (s *IntegrationTestSuite) SetupSuite() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("link", "linkpub")
	config.SetBech32PrefixForValidator("linkvaloper", "linkvaloperpub")
	config.SetBech32PrefixForConsensusNode("linkvalcons", "linkvalconspub")
	config.SetPurpose(44)
	config.SetCoinType(438)
	config.Seal()
	s.holderAcc = authtypes.NewEmptyModuleAccount(holder)
	s.blockedAcc = authtypes.NewEmptyModuleAccount(blocker)
	s.mintAcc = authtypes.NewEmptyModuleAccount(minttypes.ModuleName, authtypes.Minter)
	s.burnerAcc = authtypes.NewEmptyModuleAccount(authtypes.Burner, authtypes.Burner)
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
	s.cdc = encCfg.Codec

	storeService := runtime.NewKVStoreService(key)
	s.storeService = storeService

	ctrl := gomock.NewController(s.T())
	s.mockCtrl = ctrl
	authKeeper := banktestutil.NewMockAccountKeeper(ctrl)
	s.authKeeper = authKeeper
	newAcc := sdk.AccAddress("valid")
	s.baseAcc = authtypes.NewBaseAccountWithAddress(newAcc)
	s.authKeeper.EXPECT().GetAccount(gomock.Any(), newAcc).Return(nil).AnyTimes()
	s.authKeeper.EXPECT().AddressCodec().Return(address.NewBech32Codec("link")).AnyTimes()
	s.authKeeper.EXPECT().NewAccountWithAddress(gomock.Any(), newAcc).Return(s.baseAcc).AnyTimes()
	s.authKeeper.EXPECT().GetModuleAddress("").Return(nil).AnyTimes()
	s.authKeeper.EXPECT().GetModuleAccount(gomock.Any(), "").Return(nil).AnyTimes()
	s.authKeeper.EXPECT().HasAccount(s.ctx, s.baseAcc.GetAddress()).Return(true).AnyTimes()
	s.mockAuthKeeperFor(s.holderAcc)
	s.mockAuthKeeperFor(s.blockedAcc)
	s.mockAuthKeeperFor(s.mintAcc)
	s.mockAuthKeeperFor(s.burnerAcc)
	addrCdc := encCfg.Codec.InterfaceRegistry().SigningContext().AddressCodec()
	authorityString, err := addrCdc.BytesToString(authtypes.NewModuleAddress(govtypes.ModuleName))
	s.Require().NoError(err)
	s.addrCdc = addrCdc

	s.cut = NewBaseKeeper(
		encCfg.Codec,
		storeService,
		s.authKeeper,
		map[string]bool{},
		true,
		authorityString,
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
	s.Require().NoError(s.cut.MintCoins(s.ctx, minttypes.ModuleName, initCoins))
	s.Require().NoError(s.cut.SendCoinsFromModuleToAccount(s.ctx, minttypes.ModuleName, s.holderAcc.GetAddress(), initCoins))
	s.Require().NoError(s.cut.MintCoins(s.ctx, minttypes.ModuleName, initCoins))
	s.Require().Equal(initCoins, s.cut.GetAllBalances(s.ctx, s.holderAcc.GetAddress()))
}

func (s *IntegrationTestSuite) verifySendCoinsWithInvalidAddr() {
	s.Require().Panics(func() {
		s.cut.SendCoinsFromModuleToModule(s.ctx, "", s.holderAcc.GetName(), initCoins) // nolint:errcheck // No need
	})
	s.Require().Panics(func() {
		s.cut.SendCoinsFromModuleToModule(s.ctx, authtypes.Burner, "", initCoins) // nolint:errcheck // No need
	})
	s.Require().Panics(func() {
		s.cut.SendCoinsFromModuleToAccount(s.ctx, "", s.baseAcc.GetAddress(), initCoins) // nolint:errcheck // No need
	})
}

func (s *IntegrationTestSuite) errorForInsufficientBalance() {
	s.Require().Error(s.cut.SendCoinsFromModuleToAccount(s.ctx, s.holderAcc.GetName(), s.baseAcc.GetAddress(), initCoins.Add(initCoins...)))
}

func (s *IntegrationTestSuite) verifySendModuleToModule() {
	s.Require().NoError(s.cut.SendCoinsFromModuleToModule(s.ctx, s.holderAcc.GetName(), authtypes.Burner, initCoins))
	s.Require().Equal(sdk.NewCoins().String(), getCoinsByName(s.ctx, s.cut, s.authKeeper, s.holderAcc.GetName()).String())
	s.Require().Equal(initCoins, getCoinsByName(s.ctx, s.cut, s.authKeeper, authtypes.Burner))
}

func (s *IntegrationTestSuite) verifySendModuleToAccount() {
	s.Require().NoError(s.cut.SendCoinsFromModuleToAccount(s.ctx, authtypes.Burner, s.baseAcc.GetAddress(), initCoins))
	s.Require().Equal(sdk.NewCoins().String(), getCoinsByName(s.ctx, s.cut, s.authKeeper, authtypes.Burner).String())
	s.Require().Equal(initCoins, s.cut.GetAllBalances(s.ctx, s.baseAcc.GetAddress()))
}

func (s *IntegrationTestSuite) verifySendAccountToModule() {
	s.Require().NoError(s.cut.SendCoinsFromAccountToModule(s.ctx, s.baseAcc.GetAddress(), authtypes.Burner, initCoins))
	s.Require().Equal(sdk.NewCoins().String(), s.cut.GetAllBalances(s.ctx, s.baseAcc.GetAddress()).String())
	s.Require().Equal(initCoins, getCoinsByName(s.ctx, s.cut, s.authKeeper, authtypes.Burner))
}

func (s *IntegrationTestSuite) TestInactiveAddrOfSendCoins() {
	s.mintInitialBalances()
	s.verifySendToBlockedAddr()
	s.mintInitialBalances()
	s.verifySendToBlockedAddrAfterRemoveIt()
}

func (s *IntegrationTestSuite) verifySendToBlockedAddr() {
	s.Require().False(s.cut.IsInactiveAddr(s.blockedAcc.GetAddress()))
	s.cut.AddToInactiveAddr(s.ctx, s.blockedAcc.GetAddress())
	s.Require().True(s.cut.IsInactiveAddr(s.blockedAcc.GetAddress()))

	err := s.cut.SendCoins(s.ctx, s.holderAcc.GetAddress(), s.blockedAcc.GetAddress(), initCoins)
	s.Require().Contains(err.Error(), "is not allowed to receive funds")
}

func (s *IntegrationTestSuite) verifySendToBlockedAddrAfterRemoveIt() {
	s.cut.DeleteFromInactiveAddr(s.ctx, s.blockedAcc.GetAddress())
	s.Require().False(s.cut.IsInactiveAddr(s.blockedAcc.GetAddress()))
	s.Require().NoError(s.cut.SendCoins(s.ctx, s.holderAcc.GetAddress(), s.blockedAcc.GetAddress(), initCoins))
	s.Require().Equal(sdk.NewCoins().String(), s.cut.GetAllBalances(s.ctx, s.holderAcc.GetAddress()).String())
}

func (s *IntegrationTestSuite) TestInitializeBankPlus() {
	authorityString, err := s.addrCdc.BytesToString(authtypes.NewModuleAddress(govtypes.ModuleName))
	s.Require().NoError(err)
	newKeeper := NewBaseKeeper(
		s.cdc,
		s.storeService,
		s.authKeeper,
		map[string]bool{},
		true,
		authorityString,
		log.NewNopLogger(),
	)

	newKeeper.AddToInactiveAddr(s.ctx, s.blockedAcc.GetAddress())
	s.Require().True(newKeeper.IsInactiveAddr(s.blockedAcc.GetAddress()))

	anotherNewKeeper := NewBaseKeeper(
		s.cdc,
		s.storeService,
		s.authKeeper,
		map[string]bool{},
		true,
		authorityString,
		log.NewNopLogger(),
	)

	anotherNewKeeper.InitializeBankPlus(s.ctx)
	s.Require().True(anotherNewKeeper.IsInactiveAddr(s.blockedAcc.GetAddress()))
}

func (s *IntegrationTestSuite) TestSendCoinsFromModuleToAccount_Blacklist() {
	addr1 := sdk.AccAddress("addr1_______________")
	addr1String, err := s.addrCdc.BytesToString(addr1)
	s.Require().NoError(err)
	authorityString, err := s.addrCdc.BytesToString(authtypes.NewModuleAddress(govtypes.ModuleName))
	s.Require().NoError(err)
	newKeeper := NewBaseKeeper(
		s.cdc,
		s.storeService,
		s.authKeeper,
		map[string]bool{addr1String: true},
		true,
		authorityString,
		log.NewNopLogger(),
	)

	s.Require().NoError(newKeeper.MintCoins(s.ctx, minttypes.ModuleName, initCoins))
	s.Require().Error(newKeeper.SendCoinsFromModuleToAccount(s.ctx, minttypes.ModuleName, addr1, initCoins))
}

func (s *IntegrationTestSuite) TestInputOutputCoins() {
	// set initial balances
	s.Require().NoError(s.cut.MintCoins(s.ctx, minttypes.ModuleName, initCoins))
	s.Require().NoError(s.cut.SendCoinsFromModuleToAccount(s.ctx, minttypes.ModuleName, s.baseAcc.GetAddress(), initCoins))

	input := types.NewInput(s.baseAcc.GetAddress(), initCoins)
	outputs := []types.Output{types.NewOutput(s.burnerAcc.GetAddress(), initCoins)}

	authorityString, err := s.addrCdc.BytesToString(authtypes.NewModuleAddress(govtypes.ModuleName))
	s.Require().NoError(err)

	targetKeeper := func(isDeact bool) BaseKeeper {
		return NewBaseKeeper(
			s.cdc,
			s.storeService,
			s.authKeeper,
			make(map[string]bool),
			isDeact,
			authorityString,
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
