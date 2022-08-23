package keeper

import (
	"testing"
	"time"

	"github.com/line/ostracon/libs/log"
	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"
	dbm "github.com/tendermint/tm-db"

	"github.com/line/lbm-sdk/baseapp"
	"github.com/line/lbm-sdk/store"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/module"
	authkeeper "github.com/line/lbm-sdk/x/auth/keeper"
	authtypes "github.com/line/lbm-sdk/x/auth/types"
	authzkeeper "github.com/line/lbm-sdk/x/authz/keeper"
	"github.com/line/lbm-sdk/x/bank"
	banktypes "github.com/line/lbm-sdk/x/bank/types"
	"github.com/line/lbm-sdk/x/bankplus"
	bankpluskeeper "github.com/line/lbm-sdk/x/bankplus/keeper"
	capabilitykeeper "github.com/line/lbm-sdk/x/capability/keeper"
	capabilitytypes "github.com/line/lbm-sdk/x/capability/types"
	crisistypes "github.com/line/lbm-sdk/x/crisis/types"
	"github.com/line/lbm-sdk/x/distribution"
	distributionkeeper "github.com/line/lbm-sdk/x/distribution/keeper"
	distributiontypes "github.com/line/lbm-sdk/x/distribution/types"
	evidencetypes "github.com/line/lbm-sdk/x/evidence/types"
	"github.com/line/lbm-sdk/x/feegrant"
	govkeeper "github.com/line/lbm-sdk/x/gov/keeper"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
	ibctransfertypes "github.com/line/lbm-sdk/x/ibc/applications/transfer/types"
	ibchost "github.com/line/lbm-sdk/x/ibc/core/24-host"
	ibckeeper "github.com/line/lbm-sdk/x/ibc/core/keeper"
	minttypes "github.com/line/lbm-sdk/x/mint/types"
	"github.com/line/lbm-sdk/x/params"
	paramskeeper "github.com/line/lbm-sdk/x/params/keeper"
	paramstypes "github.com/line/lbm-sdk/x/params/types"
	paramproposal "github.com/line/lbm-sdk/x/params/types/proposal"
	slashingtypes "github.com/line/lbm-sdk/x/slashing/types"
	"github.com/line/lbm-sdk/x/staking"
	stakingkeeper "github.com/line/lbm-sdk/x/staking/keeper"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
	upgradekeeper "github.com/line/lbm-sdk/x/upgrade/keeper"
	upgradetypes "github.com/line/lbm-sdk/x/upgrade/types"
	wasmkeeper "github.com/line/lbm-sdk/x/wasm/keeper"
	"github.com/line/lbm-sdk/x/wasm/keeper/wasmtesting"
	lbmwasmtypes "github.com/line/lbm-sdk/x/wasm/lbm/types"
	wasmtypes "github.com/line/lbm-sdk/x/wasm/types"
)

type TestKeepers struct {
	wasmkeeper.TestKeepers

	ContractKeeper lbmwasmtypes.ContractOpsKeeper
	WasmKeeper     *Keeper
}

func CreateTestInput(t testing.TB, isCheckTx bool, supportedFeatures string, encoders *wasmkeeper.MessageEncoders, queriers *wasmkeeper.QueryPlugins, opts ...wasmkeeper.Option) (sdk.Context, TestKeepers) {
	// Load default wasm config
	return createTestInput(t, isCheckTx, supportedFeatures, encoders, queriers, wasmtypes.DefaultWasmConfig(), dbm.NewMemDB(), opts...)
}

func createTestInput(
	t testing.TB,
	isCheckTx bool,
	supportedFeatures string,
	encoders *wasmkeeper.MessageEncoders,
	queriers *wasmkeeper.QueryPlugins,
	wasmConfig wasmtypes.WasmConfig,
	db dbm.DB,
	opts ...wasmkeeper.Option,
) (sdk.Context, TestKeepers) {
	tempDir := t.TempDir()

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distributiontypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, ibchost.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, ibctransfertypes.StoreKey,
		capabilitytypes.StoreKey, feegrant.StoreKey, authzkeeper.StoreKey,
		wasmtypes.StoreKey,
	)
	ms := store.NewCommitMultiStore(db)
	for _, v := range keys {
		ms.MountStoreWithDB(v, sdk.StoreTypeIAVL, db)
	}
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	for _, v := range tkeys {
		ms.MountStoreWithDB(v, sdk.StoreTypeTransient, db)
	}

	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
	for _, v := range memKeys {
		ms.MountStoreWithDB(v, sdk.StoreTypeMemory, db)
	}

	require.NoError(t, ms.LoadLatestVersion())

	ctx := sdk.NewContext(ms, ocproto.Header{
		Height: 1234567,
		Time:   time.Date(2020, time.April, 22, 12, 0, 0, 0, time.UTC),
	}, isCheckTx, log.NewNopLogger())
	ctx = wasmtypes.WithTXCounter(ctx, 0)

	encodingConfig := wasmkeeper.MakeEncodingConfig(t)
	appCodec, legacyAmino := encodingConfig.Marshaler, encodingConfig.Amino

	paramsKeeper := paramskeeper.NewKeeper(
		appCodec,
		legacyAmino,
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.TStoreKey],
	)
	for _, m := range []string{
		authtypes.ModuleName,
		banktypes.ModuleName,
		stakingtypes.ModuleName,
		minttypes.ModuleName,
		distributiontypes.ModuleName,
		slashingtypes.ModuleName,
		crisistypes.ModuleName,
		ibctransfertypes.ModuleName,
		capabilitytypes.ModuleName,
		ibchost.ModuleName,
		govtypes.ModuleName,
		wasmtypes.ModuleName,
	} {
		paramsKeeper.Subspace(m)
	}
	subspace := func(m string) paramstypes.Subspace {
		r, ok := paramsKeeper.GetSubspace(m)
		require.True(t, ok)
		return r
	}
	maccPerms := map[string][]string{ // module account permissions
		authtypes.FeeCollectorName:     nil,
		distributiontypes.ModuleName:   nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		wasmtypes.ModuleName:           {authtypes.Burner},
	}
	accountKeeper := authkeeper.NewAccountKeeper(
		appCodec,
		keys[authtypes.StoreKey], // target store
		subspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount, // prototype
		maccPerms,
	)
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	bankKeeper := bankpluskeeper.NewBaseKeeper(
		appCodec,
		keys[banktypes.StoreKey],
		accountKeeper,
		subspace(banktypes.ModuleName),
		blockedAddrs,
	)
	bankKeeper.SetParams(ctx, banktypes.DefaultParams())

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		keys[stakingtypes.StoreKey],
		accountKeeper,
		bankKeeper,
		subspace(stakingtypes.ModuleName),
	)
	stakingKeeper.SetParams(ctx, wasmkeeper.TestingStakeParams)

	distKeeper := distributionkeeper.NewKeeper(
		appCodec,
		keys[distributiontypes.StoreKey],
		subspace(distributiontypes.ModuleName),
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		authtypes.FeeCollectorName,
		nil,
	)
	distKeeper.SetParams(ctx, distributiontypes.DefaultParams())
	stakingKeeper.SetHooks(distKeeper.Hooks())

	// set genesis items required for distribution
	distKeeper.SetFeePool(ctx, distributiontypes.InitialFeePool())

	upgradeKeeper := upgradekeeper.NewKeeper(
		map[int64]bool{},
		keys[upgradetypes.StoreKey],
		appCodec,
		tempDir,
		nil,
	)

	faucet := wasmkeeper.NewTestFaucet(t, ctx, bankKeeper, minttypes.ModuleName, sdk.NewCoin("stake", sdk.NewInt(100_000_000_000)))

	// set some funds ot pay out validatores, based on code from:
	// https://github.com/line/lbm-sdk/blob/95b22d3a685f7eb531198e0023ef06873835e632/x/distribution/keeper/keeper_test.go#L49-L56
	distrAcc := distKeeper.GetDistributionAccount(ctx)
	faucet.Fund(ctx, distrAcc.GetAddress(), sdk.NewCoin("stake", sdk.NewInt(2000000)))
	accountKeeper.SetModuleAccount(ctx, distrAcc)

	capabilityKeeper := capabilitykeeper.NewKeeper(
		appCodec,
		keys[capabilitytypes.StoreKey],
		memKeys[capabilitytypes.MemStoreKey],
	)
	scopedIBCKeeper := capabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedWasmKeeper := capabilityKeeper.ScopeToModule(wasmtypes.ModuleName)

	ibcKeeper := ibckeeper.NewKeeper(
		appCodec,
		keys[ibchost.StoreKey],
		subspace(ibchost.ModuleName),
		stakingKeeper,
		upgradeKeeper,
		scopedIBCKeeper,
	)

	router := baseapp.NewRouter()
	bh := bank.NewHandler(bankKeeper)
	router.AddRoute(sdk.NewRoute(banktypes.RouterKey, bh))
	sh := staking.NewHandler(stakingKeeper)
	router.AddRoute(sdk.NewRoute(stakingtypes.RouterKey, sh))
	dh := distribution.NewHandler(distKeeper)
	router.AddRoute(sdk.NewRoute(distributiontypes.RouterKey, dh))

	querier := baseapp.NewGRPCQueryRouter()
	querier.SetInterfaceRegistry(encodingConfig.InterfaceRegistry)
	msgRouter := baseapp.NewMsgServiceRouter()
	msgRouter.SetInterfaceRegistry(encodingConfig.InterfaceRegistry)

	cfg := sdk.GetConfig()
	cfg.SetAddressVerifier(wasmtypes.VerifyAddressLen())

	// set default params of lbm wasm
	paramSpace := subspace(wasmtypes.ModuleName).WithKeyTable(lbmwasmtypes.ParamKeyTable())
	defaultParams := lbmwasmtypes.DefaultParams()
	paramSpace.SetParamSet(ctx, &defaultParams)

	keeper := NewKeeper(
		appCodec,
		keys[wasmtypes.StoreKey],
		paramSpace,
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		distKeeper,
		ibcKeeper.ChannelKeeper,
		&ibcKeeper.PortKeeper,
		scopedWasmKeeper,
		wasmtesting.MockIBCTransferKeeper{},
		msgRouter,
		querier,
		tempDir,
		wasmConfig,
		supportedFeatures,
		encoders,
		queriers,
		opts...,
	)
	// add wasm handler, so we can loop-back (contracts calling contracts)
	contractKeeper := wasmkeeper.NewDefaultPermissionKeeper(&keeper)
	router.AddRoute(sdk.NewRoute(wasmtypes.RouterKey, wasmkeeper.TestHandler(contractKeeper)))

	am := module.NewManager( // minimal module set that we use for message/ query tests
		bankplus.NewAppModule(appCodec, bankKeeper, accountKeeper),
		staking.NewAppModule(appCodec, stakingKeeper, accountKeeper, bankKeeper),
		distribution.NewAppModule(appCodec, distKeeper, accountKeeper, bankKeeper, stakingKeeper),
	)
	am.RegisterServices(module.NewConfigurator(appCodec, msgRouter, querier))
	wasmtypes.RegisterMsgServer(msgRouter, wasmkeeper.NewMsgServerImpl(wasmkeeper.NewDefaultPermissionKeeper(keeper)))
	wasmtypes.RegisterQueryServer(querier, wasmkeeper.NewGrpcQuerier(appCodec, keys[wasmtypes.ModuleName], keeper, wasmConfig.SmartQueryGasLimit))

	govRouter := govtypes.NewRouter().
		AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(paramsKeeper)).
		AddRoute(distributiontypes.RouterKey, distribution.NewCommunityPoolSpendProposalHandler(distKeeper)).
		AddRoute(wasmtypes.RouterKey, NewWasmProposalHandler(keeper, lbmwasmtypes.EnableAllProposals))

	govKeeper := govkeeper.NewKeeper(
		appCodec,
		keys[govtypes.StoreKey],
		subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable()),
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		govRouter,
	)

	govKeeper.SetProposalID(ctx, govtypes.DefaultStartingProposalID)
	govKeeper.SetDepositParams(ctx, govtypes.DefaultDepositParams())
	govKeeper.SetVotingParams(ctx, govtypes.DefaultVotingParams())
	govKeeper.SetTallyParams(ctx, govtypes.DefaultTallyParams())

	wasmKeepers := wasmkeeper.TestKeepers{
		AccountKeeper:  accountKeeper,
		StakingKeeper:  stakingKeeper,
		DistKeeper:     distKeeper,
		ContractKeeper: contractKeeper,
		WasmKeeper:     &keeper.Keeper,
		BankKeeper:     bankKeeper,
		GovKeeper:      govKeeper,
		IBCKeeper:      ibcKeeper,
		Router:         router,
		EncodingConfig: encodingConfig,
		Faucet:         faucet,
		MultiStore:     ms,
	}

	lbmWasmContractKeeper := NewDefaultPermissionKeeper(&keeper)

	return ctx, TestKeepers{
		TestKeepers:    wasmKeepers,
		ContractKeeper: lbmWasmContractKeeper,
		WasmKeeper:     &keeper,
	}
}
