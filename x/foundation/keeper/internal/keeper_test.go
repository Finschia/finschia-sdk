package internal_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"cosmossdk.io/core/address"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	codectestutil "github.com/cosmos/cosmos-sdk/codec/testutil"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/Finschia/finschia-sdk/x/foundation"
	"github.com/Finschia/finschia-sdk/x/foundation/keeper"
	"github.com/Finschia/finschia-sdk/x/foundation/keeper/internal"
	"github.com/Finschia/finschia-sdk/x/foundation/module"
	foundationtestutil "github.com/Finschia/finschia-sdk/x/foundation/testutil"
)

type KeeperTestSuite struct {
	suite.Suite

	deterministic bool

	ctx sdk.Context

	bankKeeper *foundationtestutil.MockBankKeeper
	keeper     keeper.Keeper
	impl       internal.Keeper

	addressCodec address.Codec

	queryServer     foundation.QueryServer
	msgServer       foundation.MsgServer
	proposalHandler govtypes.Handler

	authority sdk.AccAddress
	members   []sdk.AccAddress
	stranger  sdk.AccAddress

	activeProposal    uint64
	votedProposal     uint64
	withdrawnProposal uint64
	invalidProposal   uint64
	noHandlerProposal uint64
	nextProposal      uint64

	balance math.Int
}

func newMsgCreateDog(name string) sdk.Msg {
	return &testdata.MsgCreateDog{
		Dog: &testdata.Dog{
			Name: name,
		},
	}
}

func (s *KeeperTestSuite) createAddresses(accNum int) []sdk.AccAddress {
	if s.deterministic {
		return simtestutil.CreateIncrementalAccounts(accNum)
	} else {
		return simtestutil.CreateRandomAccounts(accNum)
	}
}

func (s *KeeperTestSuite) bytesToString(addr sdk.AccAddress) string {
	str, err := s.addressCodec.BytesToString(addr)
	s.Require().NoError(err)
	return str
}

func (s *KeeperTestSuite) newTestMsg(addrs ...sdk.AccAddress) *testdata.TestMsg {
	accAddresses := make([]string, len(addrs))

	for i, addr := range addrs {
		accAddresses[i] = s.bytesToString(addr)
	}

	return &testdata.TestMsg{
		Signers: accAddresses,
	}
}

func setupFoundationKeeper(t *testing.T, balance *math.Int, addrs []sdk.AccAddress) (
	internal.Keeper,
	keeper.Keeper,
	*foundationtestutil.MockAuthKeeper,
	*foundationtestutil.MockBankKeeper,
	moduletestutil.TestEncodingConfig,
	address.Codec,
	sdk.Context,
) {
	key := storetypes.NewKVStoreKey(foundation.StoreKey)
	tkey := storetypes.NewTransientStoreKey("transient_test")
	testCtx := testutil.DefaultContextWithDB(t, key, tkey)
	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModuleBasic{})

	ir := codectestutil.CodecOptions{
		AccAddressPrefix: "link",
		ValAddressPrefix: "linkvaloper",
	}.NewInterfaceRegistry()
	encCfg.InterfaceRegistry = ir
	encCfg.Codec = codec.NewProtoCodec(ir)

	foundation.RegisterInterfaces(ir)
	testdata.RegisterInterfaces(ir)

	addressCodec := ir.SigningContext().AddressCodec()

	bapp := baseapp.NewBaseApp(
		"foundation",
		log.NewNopLogger(),
		testCtx.DB,
		encCfg.TxConfig.TxDecoder(),
	)
	bapp.SetInterfaceRegistry(ir)

	ctrl := gomock.NewController(t)
	authKeeper := foundationtestutil.NewMockAuthKeeper(ctrl)
	bankKeeper := foundationtestutil.NewMockBankKeeper(ctrl)
	subspace := paramstypes.NewSubspace(encCfg.Codec, encCfg.Amino, key, tkey, "params")

	authority, err := addressCodec.BytesToString(foundation.DefaultAuthority())
	require.NoError(t, err)

	config := foundation.DefaultConfig()
	feeCollector := authtypes.FeeCollectorName
	k := keeper.NewKeeper(encCfg.Codec, addressCodec, runtime.NewKVStoreService(key), bapp.MsgServiceRouter(), authKeeper, bankKeeper, feeCollector, config, authority, subspace)

	impl := internal.NewKeeper(
		encCfg.Codec,
		addressCodec,
		runtime.NewKVStoreService(key),
		bapp.MsgServiceRouter(),
		authKeeper,
		bankKeeper,
		feeCollector,
		config,
		authority,
		subspace,
	)

	msgServer := keeper.NewMsgServer(k)
	queryServer := keeper.NewQueryServer(k)

	foundation.RegisterMsgServer(bapp.MsgServiceRouter(), msgServer)
	foundation.RegisterQueryServer(bapp.GRPCQueryRouter(), queryServer)

	testdata.RegisterMsgServer(bapp.MsgServiceRouter(), testdata.MsgServerImpl{})

	// mock bank keeper
	prefix := []byte{0xff}
	getBalance := func(ctx context.Context, addr sdk.AccAddress) sdk.Coin {
		store := runtime.NewKVStoreService(key).OpenKVStore(ctx)

		bz, err := store.Get(append(prefix, addr...))
		require.NoError(t, err)

		if bz == nil {
			return sdk.NewCoin(sdk.DefaultBondDenom, math.ZeroInt())
		}

		var amt math.Int
		err = amt.Unmarshal(bz)
		require.NoError(t, err)

		return sdk.NewCoin(sdk.DefaultBondDenom, amt)
	}
	setBalance := func(ctx context.Context, addr sdk.AccAddress, amt sdk.Coin) {
		store := runtime.NewKVStoreService(key).OpenKVStore(ctx)

		bz, err := amt.Amount.Marshal()
		require.NoError(t, err)

		err = store.Set(append(prefix, addr...), bz)
		require.NoError(t, err)
	}
	send := func(ctx context.Context, sender, recipient sdk.AccAddress, amt sdk.Coins) error {
		require.LessOrEqual(t, len(amt), 1)

		if len(amt) == 0 {
			return nil
		}

		src := getBalance(ctx, sender)
		src, err := src.SafeSub(amt[0])
		if err != nil {
			return err
		}
		setBalance(ctx, sender, src)

		dst := getBalance(ctx, recipient).Add(amt[0])
		setBalance(ctx, recipient, dst)

		return nil
	}

	bankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, sender sdk.AccAddress, name string, amt sdk.Coins) error {
		recipient := authtypes.NewModuleAddress(name)
		return send(ctx, sender, recipient, amt)
	}).AnyTimes()
	bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, name string, recipient sdk.AccAddress, amt sdk.Coins) error {
		sender := authtypes.NewModuleAddress(name)
		return send(ctx, sender, recipient, amt)
	}).AnyTimes()

	bankKeeper.EXPECT().GetAllBalances(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, addr sdk.AccAddress) sdk.Coins {
		return sdk.NewCoins(getBalance(ctx, addr))
	}).AnyTimes()

	ctx := testCtx.Ctx

	// set balance
	for _, addr := range addrs{
		setBalance(ctx, addr, sdk.NewCoin(sdk.DefaultBondDenom, *balance))
	}

	return impl, k, authKeeper, bankKeeper, encCfg, addressCodec, ctx
}

func (s *KeeperTestSuite) SetupTest() {
	numMembers := 10
	addresses := s.createAddresses(numMembers + 1)
	s.members = addresses[:numMembers]
	s.stranger = addresses[len(addresses)-1]

	s.balance = math.NewInt(987654321)
	coinHolders := []sdk.AccAddress{
		s.stranger,
		authtypes.NewModuleAddress(foundation.TreasuryName),
		authtypes.NewModuleAddress(authtypes.FeeCollectorName),
	}

	var authKeeper *foundationtestutil.MockAuthKeeper
	s.impl, s.keeper, authKeeper, s.bankKeeper, _, s.addressCodec, s.ctx = setupFoundationKeeper(s.T(), &s.balance, coinHolders)

	if s.deterministic {
		s.ctx = s.ctx.WithBlockTime(time.Date(2023, 11, 7, 19, 32, 0, 0, time.UTC))
	}

	s.authority = foundation.DefaultAuthority()

	s.queryServer = keeper.NewQueryServer(s.keeper)
	s.msgServer = keeper.NewMsgServer(s.keeper)

	s.proposalHandler = keeper.NewFoundationProposalsHandler(s.keeper)

	// genesis
	gs := &foundation.GenesisState{}

	params := foundation.DefaultParams()
	params.FoundationTax = math.LegacyOneDec()
	gs.Params = params
	
	members := make([]foundation.Member, len(s.members))
	for i := range s.members {
		members[i] = foundation.Member{
			Address: s.bytesToString(s.members[i]),
		}
	}
	gs.Members = members

	info := foundation.DefaultFoundation()
	info.TotalWeight = math.LegacyNewDec(int64(len(s.members)))
	err := info.SetDecisionPolicy(workingPolicy())
	s.Require().NoError(err)
	gs.Foundation = info

	gs.Censorships = []foundation.Censorship{{
		MsgTypeUrl: sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
		Authority:  foundation.CensorshipAuthorityFoundation,
	}}

	gs.Pool = foundation.Pool{
		Treasury: sdk.NewDecCoinsFromCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance)),
	}

	err = s.keeper.InitGenesis(s.ctx, gs)
	s.Require().NoError(err)

	for _, name := range []string{
		foundation.TreasuryName,
		authtypes.FeeCollectorName,
	} {
		addr := authtypes.NewModuleAddress(name)
		account := &authtypes.ModuleAccount{
			BaseAccount: &authtypes.BaseAccount{
				Address: s.bytesToString(addr),
			},
			Name: name,
		}
		authKeeper.EXPECT().GetModuleAccount(gomock.Any(), name).Return(account).AnyTimes()
	}

	// create an active proposal, voted yes by all members except the first member
	activeProposal, err := s.impl.SubmitProposal(s.ctx, []string{s.bytesToString(s.members[0])}, "", []sdk.Msg{
		&foundation.MsgWithdrawFromTreasury{
			Authority: s.bytesToString(s.authority),
			To:        s.bytesToString(s.stranger),
			Amount:    sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance)),
		},
	})
	s.Require().NoError(err)
	s.activeProposal = *activeProposal

	for _, member := range s.members[1:] {
		err := s.impl.Vote(s.ctx, foundation.Vote{
			ProposalId: s.activeProposal,
			Voter:      s.bytesToString(member),
			Option:     foundation.VOTE_OPTION_YES,
		})
		s.Require().NoError(err)
	}

	// create a proposal voted no by all members
	votedProposal, err := s.impl.SubmitProposal(s.ctx, []string{s.bytesToString(s.members[0])}, "", []sdk.Msg{
		&foundation.MsgWithdrawFromTreasury{
			Authority: s.bytesToString(s.authority),
			To:        s.bytesToString(s.stranger),
			Amount:    sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance))},
	})
	s.Require().NoError(err)
	s.votedProposal = *votedProposal

	for _, member := range s.members {
		err := s.impl.Vote(s.ctx, foundation.Vote{
			ProposalId: s.votedProposal,
			Voter:      s.bytesToString(member),
			Option:     foundation.VOTE_OPTION_NO,
		})
		s.Require().NoError(err)
	}

	// create an withdrawn proposal
	withdrawnProposal, err := s.impl.SubmitProposal(s.ctx, []string{s.bytesToString(s.members[0])}, "", []sdk.Msg{
		&foundation.MsgWithdrawFromTreasury{
			Authority: s.bytesToString(s.authority),
			To:        s.bytesToString(s.stranger),
			Amount:    sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance))},
	})

	s.Require().NoError(err)
	s.withdrawnProposal = *withdrawnProposal

	err = s.impl.WithdrawProposal(s.ctx, s.withdrawnProposal)
	s.Require().NoError(err)

	// create an invalid proposal which contains invalid message
	invalidProposal, err := s.impl.SubmitProposal(s.ctx, []string{s.bytesToString(s.members[0])}, "", []sdk.Msg{
		&foundation.MsgWithdrawFromTreasury{
			Authority: s.bytesToString(s.authority),
			To:        s.bytesToString(s.stranger),
			Amount:    sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance.Add(math.OneInt()))),
		},
	})
	s.Require().NoError(err)
	s.invalidProposal = *invalidProposal

	for _, member := range s.members {
		err := s.impl.Vote(s.ctx, foundation.Vote{
			ProposalId: s.invalidProposal,
			Voter:      s.bytesToString(member),
			Option:     foundation.VOTE_OPTION_YES,
		})
		s.Require().NoError(err)
	}

	// create an invalid proposal which contains invalid message
	noHandlerProposal, err := s.impl.SubmitProposal(s.ctx, []string{s.bytesToString(s.members[0])}, "", []sdk.Msg{s.newTestMsg(s.authority)})
	s.Require().NoError(err)
	s.noHandlerProposal = *noHandlerProposal

	for _, member := range s.members {
		err := s.impl.Vote(s.ctx, foundation.Vote{
			ProposalId: s.noHandlerProposal,
			Voter:      s.bytesToString(member),
			Option:     foundation.VOTE_OPTION_YES,
		})
		s.Require().NoError(err)
	}

	// next proposal is the proposal id for the upcoming proposal
	s.nextProposal = s.noHandlerProposal + 1

	// grant stranger to receive foundation treasury
	err = s.impl.Grant(s.ctx, s.stranger, &foundation.ReceiveFromTreasuryAuthorization{})
	s.Require().NoError(err)
}

func TestKeeperTestSuite(t *testing.T) {
	for _, deterministic := range []bool{
		false,
		true,
	} {
		suite.Run(t, &KeeperTestSuite{deterministic: deterministic})
	}
}

func TestNewKeeper(t *testing.T) {
	createAddress := func() sdk.AccAddress {
		return sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}
	authority := foundation.DefaultAuthority()

	testCases := map[string]struct {
		authority sdk.AccAddress
		panics    bool
	}{
		"default authority": {
			authority: authority,
		},
		"invalid account": {
			panics: true,
		},
		"not the default authority": {
			authority: createAddress(),
			panics:    true,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			newKeeper := func() keeper.Keeper {
				key := storetypes.NewKVStoreKey(foundation.StoreKey)
				tkey := storetypes.NewTransientStoreKey("transient_test")
				testCtx := testutil.DefaultContextWithDB(t, key, tkey)
				encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModuleBasic{})

				ir := codectestutil.CodecOptions{
					AccAddressPrefix: "link",
					ValAddressPrefix: "linkvaloper",
				}.NewInterfaceRegistry()
				encCfg.InterfaceRegistry = ir
				encCfg.Codec = codec.NewProtoCodec(ir)

				foundation.RegisterInterfaces(ir)
				testdata.RegisterInterfaces(ir)

				addressCodec := ir.SigningContext().AddressCodec()

				bapp := baseapp.NewBaseApp(
					"foundation",
					log.NewNopLogger(),
					testCtx.DB,
					encCfg.TxConfig.TxDecoder(),
				)
				bapp.SetInterfaceRegistry(ir)

				ctrl := gomock.NewController(t)
				authKeeper := foundationtestutil.NewMockAuthKeeper(ctrl)
				bankKeeper := foundationtestutil.NewMockBankKeeper(ctrl)
				subspace := paramstypes.NewSubspace(encCfg.Codec, encCfg.Amino, key, tkey, "params")

				authority, err := addressCodec.BytesToString(tc.authority)
				require.NoError(t, err)

				config := foundation.DefaultConfig()
				feeCollector := authtypes.FeeCollectorName
				return keeper.NewKeeper(encCfg.Codec, addressCodec, runtime.NewKVStoreService(key), bapp.MsgServiceRouter(), authKeeper, bankKeeper, feeCollector, config, authority, subspace)
			}

			if tc.panics {
				require.Panics(t, func() { newKeeper() })
				return
			}
			require.NotPanics(t, func() { newKeeper() })

			k := newKeeper()
			addressCodec := addresscodec.NewBech32Codec("link")
			bytesToString := func(addr sdk.AccAddress) string {
				str, err := addressCodec.BytesToString(addr)
				require.NoError(t, err)
				return str
			}

			require.Equal(t,  bytesToString(authority), k.GetAuthority())
		})
	}
}
