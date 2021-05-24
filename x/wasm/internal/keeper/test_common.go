package keeper

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/line/lbm-sdk/v2/baseapp"
	"github.com/line/lbm-sdk/v2/codec"
	codectypes "github.com/line/lbm-sdk/v2/codec/types"
	params2 "github.com/line/lbm-sdk/v2/simapp/params"
	"github.com/line/lbm-sdk/v2/std"
	"github.com/line/lbm-sdk/v2/store"
	storetypes "github.com/line/lbm-sdk/v2/store/types"
	sdk "github.com/line/lbm-sdk/v2/types"
	sdkerrors "github.com/line/lbm-sdk/v2/types/errors"
	"github.com/line/lbm-sdk/v2/types/module"
	"github.com/line/lbm-sdk/v2/x/auth"
	authkeeper "github.com/line/lbm-sdk/v2/x/auth/keeper"
	"github.com/line/lbm-sdk/v2/x/auth/tx"
	authtypes "github.com/line/lbm-sdk/v2/x/auth/types"
	"github.com/line/lbm-sdk/v2/x/bank"
	bankkeeper "github.com/line/lbm-sdk/v2/x/bank/keeper"
	banktypes "github.com/line/lbm-sdk/v2/x/bank/types"
	"github.com/line/lbm-sdk/v2/x/capability"
	capabilitykeeper "github.com/line/lbm-sdk/v2/x/capability/keeper"
	capabilitytypes "github.com/line/lbm-sdk/v2/x/capability/types"
	"github.com/line/lbm-sdk/v2/x/crisis"
	crisistypes "github.com/line/lbm-sdk/v2/x/crisis/types"
	"github.com/line/lbm-sdk/v2/x/distribution"
	distrclient "github.com/line/lbm-sdk/v2/x/distribution/client"
	distributionkeeper "github.com/line/lbm-sdk/v2/x/distribution/keeper"
	distributiontypes "github.com/line/lbm-sdk/v2/x/distribution/types"
	"github.com/line/lbm-sdk/v2/x/evidence"
	"github.com/line/lbm-sdk/v2/x/gov"
	govkeeper "github.com/line/lbm-sdk/v2/x/gov/keeper"
	govtypes "github.com/line/lbm-sdk/v2/x/gov/types"
	"github.com/line/lbm-sdk/v2/x/ibc/applications/transfer"
	ibctransfertypes "github.com/line/lbm-sdk/v2/x/ibc/applications/transfer/types"
	ibc "github.com/line/lbm-sdk/v2/x/ibc/core"
	ibchost "github.com/line/lbm-sdk/v2/x/ibc/core/24-host"
	ibckeeper "github.com/line/lbm-sdk/v2/x/ibc/core/keeper"
	"github.com/line/lbm-sdk/v2/x/mint"
	minttypes "github.com/line/lbm-sdk/v2/x/mint/types"
	"github.com/line/lbm-sdk/v2/x/params"
	paramsclient "github.com/line/lbm-sdk/v2/x/params/client"
	paramskeeper "github.com/line/lbm-sdk/v2/x/params/keeper"
	paramstypes "github.com/line/lbm-sdk/v2/x/params/types"
	paramproposal "github.com/line/lbm-sdk/v2/x/params/types/proposal"
	"github.com/line/lbm-sdk/v2/x/slashing"
	slashingtypes "github.com/line/lbm-sdk/v2/x/slashing/types"
	"github.com/line/lbm-sdk/v2/x/staking"
	stakingkeeper "github.com/line/lbm-sdk/v2/x/staking/keeper"
	stakingtypes "github.com/line/lbm-sdk/v2/x/staking/types"
	"github.com/line/lbm-sdk/v2/x/upgrade"
	upgradeclient "github.com/line/lbm-sdk/v2/x/upgrade/client"
	"github.com/line/lbm-sdk/v2/x/wasm/internal/types"
	"github.com/line/ostracon/crypto"
	"github.com/line/ostracon/crypto/ed25519"
	"github.com/line/ostracon/libs/log"
	"github.com/line/ostracon/libs/rand"
	tmproto "github.com/line/ostracon/proto/ostracon/types"
	dbm "github.com/line/tm-db/v2"
	"github.com/line/tm-db/v2/memdb"
	"github.com/stretchr/testify/require"
)

type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
	TempDir() string
	Helper()
}

var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distribution.AppModuleBasic{},
	gov.NewAppModuleBasic(
		paramsclient.ProposalHandler, distrclient.ProposalHandler, upgradeclient.ProposalHandler,
	),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	ibc.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	transfer.AppModuleBasic{},
)

func MakeTestCodec(t TestingT) codec.Marshaler {
	return MakeEncodingConfig(t).Marshaler
}

func MakeEncodingConfig(_ TestingT) params2.EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(marshaler, tx.DefaultSignModes)

	std.RegisterInterfaces(interfaceRegistry)
	std.RegisterLegacyAminoCodec(amino)

	ModuleBasics.RegisterLegacyAminoCodec(amino)
	ModuleBasics.RegisterInterfaces(interfaceRegistry)
	types.RegisterInterfaces(interfaceRegistry)
	types.RegisterLegacyAminoCodec(amino)

	return params2.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
}

var TestingStakeParams = stakingtypes.Params{
	UnbondingTime:     100,
	MaxValidators:     10,
	MaxEntries:        10,
	HistoricalEntries: 10,
	BondDenom:         "stake",
}

type TestKeepers struct {
	AccountKeeper authkeeper.AccountKeeper
	StakingKeeper stakingkeeper.Keeper
	DistKeeper    distributionkeeper.Keeper
	BankKeeper    bankkeeper.Keeper
	GovKeeper     govkeeper.Keeper
	WasmKeeper    *Keeper
	IBCKeeper     *ibckeeper.Keeper
}

// CreateDefaultTestInput common settings for CreateTestInput
func CreateDefaultTestInput(t TestingT) (sdk.Context, TestKeepers) {
	return CreateTestInput(t, false, "staking", nil, nil)
}

// encoders can be nil to accept the defaults, or set it to override some of the message handlers (like default)
func CreateTestInput(t TestingT, isCheckTx bool, supportedFeatures string, encoders *MessageEncoders, queriers *QueryPlugins) (sdk.Context, TestKeepers) {
	// Load default wasm config
	wasmConfig := types.DefaultWasmConfig()
	db := memdb.NewDB()
	return createTestInput(t, isCheckTx, supportedFeatures, encoders, queriers, wasmConfig, db)
}

// encoders can be nil to accept the defaults, or set it to override some of the message handlers (like default)
func createTestInput(
	t TestingT,
	isCheckTx bool,
	supportedFeatures string,
	encoders *MessageEncoders,
	queriers *QueryPlugins,
	wasmConfig types.WasmConfig,
	db dbm.DB,
) (sdk.Context, TestKeepers) {
	tempDir := t.TempDir()

	keyWasm := sdk.NewKVStoreKey(types.StoreKey)
	keyAcc := sdk.NewKVStoreKey(authtypes.StoreKey)
	keyBank := sdk.NewKVStoreKey(banktypes.StoreKey)
	keyStaking := sdk.NewKVStoreKey(stakingtypes.StoreKey)
	keyDistro := sdk.NewKVStoreKey(distributiontypes.StoreKey)
	keyParams := sdk.NewKVStoreKey(paramstypes.StoreKey)
	keyGov := sdk.NewKVStoreKey(govtypes.StoreKey)
	keyIBC := sdk.NewKVStoreKey(ibchost.StoreKey)
	keyCapability := sdk.NewKVStoreKey(capabilitytypes.StoreKey)
	keyCapabilityTransient := storetypes.NewMemoryStoreKey(capabilitytypes.MemStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyWasm, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyBank, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDistro, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyGov, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyIBC, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyCapability, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyCapabilityTransient, sdk.StoreTypeMemory, db)
	require.NoError(t, ms.LoadLatestVersion())

	ctx := sdk.NewContext(ms, tmproto.Header{
		Height: 1234567,
		Time:   time.Date(2020, time.April, 22, 12, 0, 0, 0, time.UTC),
	}, isCheckTx, log.NewNopLogger())
	encodingConfig := MakeEncodingConfig(t)
	appCodec, legacyAmino := encodingConfig.Marshaler, encodingConfig.Amino

	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, keyParams)
	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distributiontypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(capabilitytypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)

	maccPerms := map[string][]string{ // module account permissions
		authtypes.FeeCollectorName:     nil,
		distributiontypes.ModuleName:   nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
	}
	authSubsp, _ := paramsKeeper.GetSubspace(authtypes.ModuleName)
	authKeeper := authkeeper.NewAccountKeeper(
		appCodec,
		keyAcc, // target store
		authSubsp,
		authtypes.ProtoBaseAccount, // prototype
		maccPerms,
	)
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		allowReceivingFunds := acc != distributiontypes.ModuleName
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = allowReceivingFunds
	}

	bankSubsp, _ := paramsKeeper.GetSubspace(banktypes.ModuleName)
	bankKeeper := bankkeeper.NewBaseKeeper(
		appCodec,
		keyBank,
		authKeeper,
		bankSubsp,
		blockedAddrs,
	)
	bankParams := banktypes.DefaultParams()
	bankParams = bankParams.SetSendEnabledParam("stake", true)
	bankKeeper.SetParams(ctx, bankParams)

	stakingSubsp, _ := paramsKeeper.GetSubspace(stakingtypes.ModuleName)
	stakingKeeper := stakingkeeper.NewKeeper(appCodec, keyStaking, authKeeper, bankKeeper, stakingSubsp)
	stakingKeeper.SetParams(ctx, TestingStakeParams)

	distSubsp, _ := paramsKeeper.GetSubspace(distributiontypes.ModuleName)
	distKeeper := distributionkeeper.NewKeeper(appCodec, keyDistro, distSubsp, authKeeper, bankKeeper, stakingKeeper, authtypes.FeeCollectorName, nil)
	distKeeper.SetParams(ctx, distributiontypes.DefaultParams())
	stakingKeeper.SetHooks(distKeeper.Hooks())

	// set genesis items required for distribution
	distKeeper.SetFeePool(ctx, distributiontypes.InitialFeePool())

	// set some funds ot pay out validatores, based on code from:
	// https://github.com/cosmos/cosmos-sdk/blob/fea231556aee4d549d7551a6190389c4328194eb/x/distribution/keeper/keeper_test.go#L50-L57
	distrAcc := distKeeper.GetDistributionAccount(ctx)
	err := bankKeeper.SetBalances(ctx, distrAcc.GetAddress(), sdk.NewCoins(
		sdk.NewCoin("stake", sdk.NewInt(2000000)),
	))
	require.NoError(t, err)
	authKeeper.SetModuleAccount(ctx, distrAcc)
	capabilityKeeper := capabilitykeeper.NewKeeper(appCodec, keyCapability, keyCapabilityTransient)
	scopedIBCKeeper := capabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedWasmKeeper := capabilityKeeper.ScopeToModule(types.ModuleName)

	ibcSubsp, _ := paramsKeeper.GetSubspace(ibchost.ModuleName)

	ibcKeeper := ibckeeper.NewKeeper(
		appCodec, keyIBC, ibcSubsp, stakingKeeper, scopedIBCKeeper,
	)

	router := baseapp.NewRouter()
	bh := bank.NewHandler(bankKeeper)
	router.AddRoute(sdk.NewRoute(banktypes.RouterKey, bh))
	sh := staking.NewHandler(stakingKeeper)
	router.AddRoute(sdk.NewRoute(stakingtypes.RouterKey, sh))
	dh := distribution.NewHandler(distKeeper)
	router.AddRoute(sdk.NewRoute(distributiontypes.RouterKey, dh))

	querier := baseapp.NewGRPCQueryRouter()
	banktypes.RegisterQueryServer(querier, bankKeeper)
	stakingtypes.RegisterQueryServer(querier, stakingkeeper.Querier{Keeper: stakingKeeper})
	distributiontypes.RegisterQueryServer(querier, distKeeper)

	keeper := NewKeeper(
		appCodec,
		keyWasm,
		paramsKeeper.Subspace(types.DefaultParamspace),
		authKeeper,
		bankKeeper,
		stakingKeeper,
		distKeeper,
		ibcKeeper.ChannelKeeper,
		&ibcKeeper.PortKeeper,
		scopedWasmKeeper,
		router,
		nil,
		querier,
		tempDir,
		wasmConfig,
		supportedFeatures,
		encoders,
		queriers,
	)
	keeper.setParams(ctx, types.DefaultParams())
	// add wasm handler so we can loop-back (contracts calling contracts)
	router.AddRoute(sdk.NewRoute(types.RouterKey, TestHandler(&keeper)))

	govRouter := govtypes.NewRouter().
		AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(paramsKeeper)).
		AddRoute(distributiontypes.RouterKey, distribution.NewCommunityPoolSpendProposalHandler(distKeeper)).
		AddRoute(types.RouterKey, NewWasmProposalHandler(&keeper, types.EnableAllProposals))

	govKeeper := govkeeper.NewKeeper(
		appCodec, keyGov, paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable()), authKeeper, bankKeeper, stakingKeeper, govRouter,
	)

	govKeeper.SetProposalID(ctx, govtypes.DefaultStartingProposalID)
	govKeeper.SetDepositParams(ctx, govtypes.DefaultDepositParams())
	govKeeper.SetVotingParams(ctx, govtypes.DefaultVotingParams())
	govKeeper.SetTallyParams(ctx, govtypes.DefaultTallyParams())

	keepers := TestKeepers{
		AccountKeeper: authKeeper,
		StakingKeeper: stakingKeeper,
		DistKeeper:    distKeeper,
		WasmKeeper:    &keeper,
		BankKeeper:    bankKeeper,
		GovKeeper:     govKeeper,
		IBCKeeper:     ibcKeeper,
	}
	return ctx, keepers
}

// TestHandler returns a wasm handler for tests (to avoid circular imports)
func TestHandler(k *Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgStoreCode:
			return handleStoreCode(ctx, k, msg)
		case *types.MsgInstantiateContract:
			return handleInstantiate(ctx, k, msg)
		case *types.MsgExecuteContract:
			return handleExecute(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized wasm message type: %T", msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleStoreCode(ctx sdk.Context, k *Keeper, msg *types.MsgStoreCode) (*sdk.Result, error) {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	codeID, err := k.Create(ctx, senderAddr, msg.WASMByteCode, msg.Source, msg.Builder, msg.InstantiatePermission)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{
		Data:   []byte(fmt.Sprintf("%d", codeID)),
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}

func handleInstantiate(ctx sdk.Context, k *Keeper, msg *types.MsgInstantiateContract) (*sdk.Result, error) {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	var adminAddr sdk.AccAddress
	if msg.Admin != "" {
		if adminAddr, err = sdk.AccAddressFromBech32(msg.Admin); err != nil {
			return nil, sdkerrors.Wrap(err, "admin")
		}
	}

	contractAddr, _, err := k.Instantiate(ctx, msg.CodeID, senderAddr, adminAddr, msg.InitMsg, msg.Label, msg.Funds)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{
		Data:   contractAddr,
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func handleExecute(ctx sdk.Context, k *Keeper, msg *types.MsgExecuteContract) (*sdk.Result, error) {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	contractAddr, err := sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "admin")
	}
	res, err := k.Execute(ctx, contractAddr, senderAddr, msg.Msg, msg.Funds)
	if err != nil {
		return nil, err
	}

	res.Events = ctx.EventManager().Events().ToABCIEvents()
	return res, nil
}

func RandomAccountAddress(_ TestingT) sdk.AccAddress {
	_, _, addr := keyPubAddr()
	return addr
}

func RandomBech32AccountAddress(t TestingT) string {
	return RandomAccountAddress(t).String()
}

type ExampleContract struct {
	InitialAmount sdk.Coins
	Creator       crypto.PrivKey
	CreatorAddr   sdk.AccAddress
	CodeID        uint64
}

func StoreHackatomExampleContract(t TestingT, ctx sdk.Context, keepers TestKeepers) ExampleContract {
	return StoreExampleContract(t, ctx, keepers, "./testdata/hackatom.wasm")
}

func StoreBurnerExampleContract(t TestingT, ctx sdk.Context, keepers TestKeepers) ExampleContract {
	return StoreExampleContract(t, ctx, keepers, "./testdata/burner.wasm")
}

func StoreIBCReflectContract(t TestingT, ctx sdk.Context, keepers TestKeepers) ExampleContract {
	return StoreExampleContract(t, ctx, keepers, "./testdata/ibc_reflect.wasm")
}

func StoreReflectContract(t TestingT, ctx sdk.Context, keepers TestKeepers) uint64 {
	wasmCode, err := ioutil.ReadFile("./testdata/reflect.wasm")
	require.NoError(t, err)

	_, _, creatorAddr := keyPubAddr()
	codeID, err := keepers.WasmKeeper.Create(ctx, creatorAddr, wasmCode, "", "", nil)
	require.NoError(t, err)
	return codeID
}

func StoreExampleContract(t TestingT, ctx sdk.Context, keepers TestKeepers, wasmFile string) ExampleContract {
	anyAmount := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000))
	creator, _, creatorAddr := keyPubAddr()
	fundAccounts(t, ctx, keepers.AccountKeeper, keepers.BankKeeper, creatorAddr, anyAmount)

	wasmCode, err := ioutil.ReadFile(wasmFile)
	require.NoError(t, err)

	codeID, err := keepers.WasmKeeper.Create(ctx, creatorAddr, wasmCode, "", "", nil)
	require.NoError(t, err)
	return ExampleContract{anyAmount, creator, creatorAddr, codeID}
}

var wasmIdent = []byte("\x00\x61\x73\x6D")

type ExampleContractInstance struct {
	ExampleContract
	Contract sdk.AccAddress
}

// SeedNewContractInstance sets the mock wasmerEngine in keeper and calls store + instantiate to init the contract's metadata
func SeedNewContractInstance(t TestingT, ctx sdk.Context, keepers TestKeepers, mock types.WasmerEngine) ExampleContractInstance {
	t.Helper()
	exampleContract := StoreRandomContract(t, ctx, keepers, mock)
	contractAddr, _, err := keepers.WasmKeeper.Instantiate(ctx, exampleContract.CodeID, exampleContract.CreatorAddr, exampleContract.CreatorAddr, []byte(`{}`), "", nil)
	require.NoError(t, err)
	return ExampleContractInstance{
		ExampleContract: exampleContract,
		Contract:        contractAddr,
	}
}

// StoreRandomContract sets the mock wasmerEngine in keeper and calls store
func StoreRandomContract(t TestingT, ctx sdk.Context, keepers TestKeepers, mock types.WasmerEngine) ExampleContract {
	t.Helper()
	anyAmount := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000))
	creator, _, creatorAddr := keyPubAddr()
	fundAccounts(t, ctx, keepers.AccountKeeper, keepers.BankKeeper, creatorAddr, anyAmount)
	keepers.WasmKeeper.wasmer = mock
	wasmCode := append(wasmIdent, rand.Bytes(10)...)
	codeID, err := keepers.WasmKeeper.Create(ctx, creatorAddr, wasmCode, "", "", nil)
	require.NoError(t, err)
	exampleContract := ExampleContract{InitialAmount: anyAmount, Creator: creator, CreatorAddr: creatorAddr, CodeID: codeID}
	return exampleContract
}

type HackatomExampleInstance struct {
	ExampleContract
	Contract        sdk.AccAddress
	Verifier        crypto.PrivKey
	VerifierAddr    sdk.AccAddress
	Beneficiary     crypto.PrivKey
	BeneficiaryAddr sdk.AccAddress
}

// InstantiateHackatomExampleContract load and instantiate the "./testdata/hackatom.wasm" contract
func InstantiateHackatomExampleContract(t TestingT, ctx sdk.Context, keepers TestKeepers) HackatomExampleInstance {
	contract := StoreHackatomExampleContract(t, ctx, keepers)

	verifier, _, verifierAddr := keyPubAddr()
	fundAccounts(t, ctx, keepers.AccountKeeper, keepers.BankKeeper, verifierAddr, contract.InitialAmount)

	beneficiary, _, beneficiaryAddr := keyPubAddr()
	initMsgBz := HackatomExampleInitMsg{
		Verifier:    verifierAddr,
		Beneficiary: beneficiaryAddr,
	}.GetBytes(t)
	initialAmount := sdk.NewCoins(sdk.NewInt64Coin("denom", 100))

	adminAddr := contract.CreatorAddr
	contractAddr, _, err := keepers.WasmKeeper.Instantiate(ctx, contract.CodeID, contract.CreatorAddr, adminAddr, initMsgBz, "demo contract to query", initialAmount)
	require.NoError(t, err)
	return HackatomExampleInstance{
		ExampleContract: contract,
		Contract:        contractAddr,
		Verifier:        verifier,
		VerifierAddr:    verifierAddr,
		Beneficiary:     beneficiary,
		BeneficiaryAddr: beneficiaryAddr,
	}
}

type HackatomExampleInitMsg struct {
	Verifier    sdk.AccAddress `json:"verifier"`
	Beneficiary sdk.AccAddress `json:"beneficiary"`
}

func (m HackatomExampleInitMsg) GetBytes(t TestingT) []byte {
	initMsgBz, err := json.Marshal(m)
	require.NoError(t, err)
	return initMsgBz
}

type IBCReflectExampleInstance struct {
	Contract      sdk.AccAddress
	Admin         sdk.AccAddress
	CodeID        uint64
	ReflectCodeID uint64
}

// InstantiateIBCReflectContract load and instantiate the "./testdata/ibc_reflect.wasm" contract
func InstantiateIBCReflectContract(t TestingT, ctx sdk.Context, keepers TestKeepers) IBCReflectExampleInstance {
	reflectID := StoreReflectContract(t, ctx, keepers)
	ibcReflectID := StoreIBCReflectContract(t, ctx, keepers).CodeID

	initMsgBz := IBCReflectInitMsg{
		ReflectCodeID: reflectID,
	}.GetBytes(t)
	adminAddr := RandomAccountAddress(t)

	contractAddr, _, err := keepers.WasmKeeper.Instantiate(ctx, ibcReflectID, adminAddr, adminAddr, initMsgBz, "ibc-reflect-factory", nil)
	require.NoError(t, err)
	return IBCReflectExampleInstance{
		Admin:         adminAddr,
		Contract:      contractAddr,
		CodeID:        ibcReflectID,
		ReflectCodeID: reflectID,
	}
}

type IBCReflectInitMsg struct {
	ReflectCodeID uint64 `json:"reflect_code_id"`
}

func (m IBCReflectInitMsg) GetBytes(t TestingT) []byte {
	initMsgBz, err := json.Marshal(m)
	require.NoError(t, err)
	return initMsgBz
}

type BurnerExampleInitMsg struct {
	Payout sdk.AccAddress `json:"payout"`
}

func (m BurnerExampleInitMsg) GetBytes(t TestingT) []byte {
	initMsgBz, err := json.Marshal(m)
	require.NoError(t, err)
	return initMsgBz
}

//nolint:deadcode,unused
func createFakeFundedAccount(t TestingT, ctx sdk.Context, am authkeeper.AccountKeeper, bank bankkeeper.Keeper, coins sdk.Coins) sdk.AccAddress {
	_, _, addr := keyPubAddr()
	fundAccounts(t, ctx, am, bank, addr, coins)
	return addr
}

func fundAccounts(t TestingT, ctx sdk.Context, am authkeeper.AccountKeeper, bank bankkeeper.Keeper, addr sdk.AccAddress, coins sdk.Coins) {
	acc := am.NewAccountWithAddress(ctx, addr)
	am.SetAccount(ctx, acc)
	require.NoError(t, bank.SetBalances(ctx, addr, coins))
}

var keyCounter uint64 = 0

// we need to make this deterministic (same every test run), as encoded address size and thus gas cost,
// depends on the actual bytes (due to ugly CanonicalAddress encoding)
//nolint:unparam
func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	keyCounter++
	seed := make([]byte, 8)
	binary.BigEndian.PutUint64(seed, keyCounter)

	key := ed25519.GenPrivKeyFromSecret(seed)
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}
