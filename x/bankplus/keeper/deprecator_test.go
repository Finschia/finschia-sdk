package keeper

import (
	"context"
	"testing"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"
	"github.com/stretchr/testify/suite"

	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	codectestutil "github.com/cosmos/cosmos-sdk/codec/testutil"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/Finschia/finschia-sdk/x/bankplus/types"
)

func TestDeprecateTestSuite(t *testing.T) {
	suite.Run(t, &DeprecationTestSuite{})
}

type DeprecationTestSuite struct {
	suite.Suite
	ctx          sdk.Context
	cdc          codec.Codec
	storeService store.KVStoreService
}

func (s *DeprecationTestSuite) SetupTest() {
	key := storetypes.NewKVStoreKey(banktypes.StoreKey)
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
}

func (s *DeprecationTestSuite) TestDeprecateBankPlus() {
	oldAcc := authtypes.NewBaseAccountWithAddress(sdk.AccAddress("acc1"))
	anotherOldAcc := authtypes.NewBaseAccountWithAddress(sdk.AccAddress("acc2"))
	s.Require().False(isStoredInactiveAddr(s.ctx, s.storeService, oldAcc.GetAddress()))
	s.Require().False(isStoredInactiveAddr(s.ctx, s.storeService, anotherOldAcc.GetAddress()))

	addrCdc := s.cdc.InterfaceRegistry().SigningContext().AddressCodec()
	addToInactiveAddr(s.ctx, s.storeService, s.cdc, addrCdc, oldAcc.GetAddress())
	addToInactiveAddr(s.ctx, s.storeService, s.cdc, addrCdc, anotherOldAcc.GetAddress())
	s.Require().True(isStoredInactiveAddr(s.ctx, s.storeService, oldAcc.GetAddress()))
	s.Require().True(isStoredInactiveAddr(s.ctx, s.storeService, anotherOldAcc.GetAddress()))

	err := DeprecateBankPlus(s.ctx, BaseKeeper{storeService: s.storeService})

	s.Require().NoError(err)
	s.Require().False(isStoredInactiveAddr(s.ctx, s.storeService, oldAcc.GetAddress()))
	s.Require().False(isStoredInactiveAddr(s.ctx, s.storeService, anotherOldAcc.GetAddress()))
}

// isStoredInactiveAddr checks if the address is stored or not as blocked address
func isStoredInactiveAddr(ctx context.Context, storeService store.KVStoreService, address sdk.AccAddress) bool {
	kvStore := storeService.OpenKVStore(ctx)
	bz, _ := kvStore.Get(inactiveAddrKey(address))
	return bz != nil
}

// addToInactiveAddr adds a blocked address to the store.
func addToInactiveAddr(ctx context.Context, storeService store.KVStoreService, cdc codec.Codec, addrCdc address.Codec, address sdk.AccAddress) {
	kvStore := storeService.OpenKVStore(ctx)
	addrString, err := addrCdc.BytesToString(address)
	if err != nil {
		panic(err)
	}

	blockedCAddr := types.InactiveAddr{Address: addrString}
	bz := cdc.MustMarshal(&blockedCAddr)
	if err := kvStore.Set(inactiveAddrKey(address), bz); err != nil {
		panic(err)
	}
}
