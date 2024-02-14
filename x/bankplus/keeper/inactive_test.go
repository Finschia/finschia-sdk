package keeper

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	codectestutil "github.com/cosmos/cosmos-sdk/codec/testutil"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func TestBankPlus(t *testing.T) {
	suite.Run(t, &BankPlusTestSuite{})
}

type BankPlusTestSuite struct {
	suite.Suite
	mockCtrl *gomock.Controller
	cut      BaseKeeper
	ctx      context.Context
}

func (s *BankPlusTestSuite) SetupTest() {
	key := storetypes.NewKVStoreKey(banktypes.StoreKey)
	tkey := storetypes.NewTransientStoreKey("transient_test")
	kvStoreService := runtime.NewKVStoreService(key)
	s.ctx = testutil.DefaultContextWithDB(s.T(), key, tkey).Ctx

	ctrl := gomock.NewController(s.T())
	mockAccKeeper := banktestutil.NewMockAccountKeeper(ctrl)
	mockAccKeeper.EXPECT().AddressCodec().Return(codec.NewBech32Codec("link")).AnyTimes()

	codec := codectestutil.CodecOptions{
		AccAddressPrefix: "link",
		ValAddressPrefix: "linkvaloper",
	}.NewCodec()
	s.mockCtrl = ctrl
	authority, err := codec.InterfaceRegistry().
		SigningContext().
		AddressCodec().
		BytesToString(authtypes.NewModuleAddress(govtypes.ModuleName))
	s.Require().NoError(err)

	s.cut = NewBaseKeeper(
		codec,
		kvStoreService,
		mockAccKeeper,
		map[string]bool{},
		true,
		authority,
		log.NewNopLogger())
}

func (s *BankPlusTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *BankPlusTestSuite) TestInactiveAddress() {
	require.Equal(s.T(), 0, len(s.cut.inactiveAddrs))
	addr := genAddr()
	anotherAddr := genAddr()
	s.addAddrOk(addr)
	s.duplicateAddOk(addr)
	s.deleteAddrOk(addr)
	s.falseForUnknownAddr(anotherAddr)
	s.noErrorWhenDeletionOfUnknownAddr(anotherAddr)
	s.testLoadAllInactiveAddrs(addr, anotherAddr)
}

func (s *BankPlusTestSuite) addAddrOk(addr sdk.AccAddress) {
	s.cut.addToInactiveAddr(s.ctx, addr)
	require.True(s.T(), s.cut.isStoredInactiveAddr(s.ctx, addr))
}

func (s *BankPlusTestSuite) duplicateAddOk(addr sdk.AccAddress) {
	s.cut.addToInactiveAddr(s.ctx, addr)
	require.True(s.T(), s.cut.isStoredInactiveAddr(s.ctx, addr))
}

func (s *BankPlusTestSuite) deleteAddrOk(addr sdk.AccAddress) {
	s.cut.deleteFromInactiveAddr(s.ctx, addr)
	require.False(s.T(), s.cut.isStoredInactiveAddr(s.ctx, addr))
}

func (s *BankPlusTestSuite) falseForUnknownAddr(anotherAddr sdk.AccAddress) {
	require.False(s.T(), s.cut.isStoredInactiveAddr(s.ctx, anotherAddr))
}

func (s *BankPlusTestSuite) noErrorWhenDeletionOfUnknownAddr(anotherAddr sdk.AccAddress) {
	require.NotPanicsf(s.T(), func() {
		s.cut.deleteFromInactiveAddr(s.ctx, anotherAddr)
	}, "no panic expected")
}

func (s *BankPlusTestSuite) testLoadAllInactiveAddrs(addr, anotherAddr sdk.AccAddress) {
	s.cut.addToInactiveAddr(s.ctx, addr)
	s.cut.addToInactiveAddr(s.ctx, anotherAddr)
	require.Equal(s.T(), 0, len(s.cut.inactiveAddrs))
	s.cut.loadAllInactiveAddrs(s.ctx)
	require.Equal(s.T(), 2, len(s.cut.inactiveAddrs))
}

func genAddr() sdk.AccAddress {
	pk := secp256k1.GenPrivKey().PubKey()
	return sdk.AccAddress(pk.Address())
}
