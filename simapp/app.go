package simapp

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	ocabci "github.com/Finschia/ostracon/abci/types"
	"github.com/Finschia/ostracon/libs/log"
	ostos "github.com/Finschia/ostracon/libs/os"

	"github.com/Finschia/finschia-sdk/baseapp"
	"github.com/Finschia/finschia-sdk/client"
	nodeservice "github.com/Finschia/finschia-sdk/client/grpc/node"
	"github.com/Finschia/finschia-sdk/client/grpc/ocservice"
	"github.com/Finschia/finschia-sdk/client/grpc/tmservice"
	"github.com/Finschia/finschia-sdk/codec"
	"github.com/Finschia/finschia-sdk/codec/types"
	"github.com/Finschia/finschia-sdk/server/api"
	"github.com/Finschia/finschia-sdk/server/config"
	servertypes "github.com/Finschia/finschia-sdk/server/types"
	appante "github.com/Finschia/finschia-sdk/simapp/ante"
	simappparams "github.com/Finschia/finschia-sdk/simapp/params"
	"github.com/Finschia/finschia-sdk/store/streaming"
	"github.com/Finschia/finschia-sdk/testutil/testdata"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/module"
	"github.com/Finschia/finschia-sdk/version"
	"github.com/Finschia/finschia-sdk/x/auth"
	"github.com/Finschia/finschia-sdk/x/auth/ante"
	authkeeper "github.com/Finschia/finschia-sdk/x/auth/keeper"
	authsims "github.com/Finschia/finschia-sdk/x/auth/simulation"
	authtx "github.com/Finschia/finschia-sdk/x/auth/tx"
	authtx2 "github.com/Finschia/finschia-sdk/x/auth/tx2"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	"github.com/Finschia/finschia-sdk/x/auth/vesting"
	vestingtypes "github.com/Finschia/finschia-sdk/x/auth/vesting/types"
	"github.com/Finschia/finschia-sdk/x/authz"
	authzkeeper "github.com/Finschia/finschia-sdk/x/authz/keeper"
	authzmodule "github.com/Finschia/finschia-sdk/x/authz/module"
	"github.com/Finschia/finschia-sdk/x/bank"
	bankkeeper "github.com/Finschia/finschia-sdk/x/bank/keeper"
	banktypes "github.com/Finschia/finschia-sdk/x/bank/types"
	"github.com/Finschia/finschia-sdk/x/bankplus"
	bankpluskeeper "github.com/Finschia/finschia-sdk/x/bankplus/keeper"
	"github.com/Finschia/finschia-sdk/x/capability"
	capabilitykeeper "github.com/Finschia/finschia-sdk/x/capability/keeper"
	capabilitytypes "github.com/Finschia/finschia-sdk/x/capability/types"
	"github.com/Finschia/finschia-sdk/x/collection"
	collectionkeeper "github.com/Finschia/finschia-sdk/x/collection/keeper"
	collectionmodule "github.com/Finschia/finschia-sdk/x/collection/module"
	"github.com/Finschia/finschia-sdk/x/crisis"
	crisiskeeper "github.com/Finschia/finschia-sdk/x/crisis/keeper"
	crisistypes "github.com/Finschia/finschia-sdk/x/crisis/types"
	distr "github.com/Finschia/finschia-sdk/x/distribution"
	distrclient "github.com/Finschia/finschia-sdk/x/distribution/client"
	distrkeeper "github.com/Finschia/finschia-sdk/x/distribution/keeper"
	distrtypes "github.com/Finschia/finschia-sdk/x/distribution/types"
	"github.com/Finschia/finschia-sdk/x/evidence"
	evidencekeeper "github.com/Finschia/finschia-sdk/x/evidence/keeper"
	evidencetypes "github.com/Finschia/finschia-sdk/x/evidence/types"
	"github.com/Finschia/finschia-sdk/x/feegrant"
	feegrantkeeper "github.com/Finschia/finschia-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/Finschia/finschia-sdk/x/feegrant/module"
	"github.com/Finschia/finschia-sdk/x/foundation"
	foundationclient "github.com/Finschia/finschia-sdk/x/foundation/client"
	foundationkeeper "github.com/Finschia/finschia-sdk/x/foundation/keeper"
	foundationmodule "github.com/Finschia/finschia-sdk/x/foundation/module"
	"github.com/Finschia/finschia-sdk/x/genutil"
	genutiltypes "github.com/Finschia/finschia-sdk/x/genutil/types"
	"github.com/Finschia/finschia-sdk/x/gov"
	govkeeper "github.com/Finschia/finschia-sdk/x/gov/keeper"
	govtypes "github.com/Finschia/finschia-sdk/x/gov/types"
	"github.com/Finschia/finschia-sdk/x/mint"
	mintkeeper "github.com/Finschia/finschia-sdk/x/mint/keeper"
	minttypes "github.com/Finschia/finschia-sdk/x/mint/types"
	"github.com/Finschia/finschia-sdk/x/params"
	paramsclient "github.com/Finschia/finschia-sdk/x/params/client"
	paramskeeper "github.com/Finschia/finschia-sdk/x/params/keeper"
	paramstypes "github.com/Finschia/finschia-sdk/x/params/types"
	paramproposal "github.com/Finschia/finschia-sdk/x/params/types/proposal"
	"github.com/Finschia/finschia-sdk/x/slashing"
	slashingkeeper "github.com/Finschia/finschia-sdk/x/slashing/keeper"
	slashingtypes "github.com/Finschia/finschia-sdk/x/slashing/types"
	"github.com/Finschia/finschia-sdk/x/staking"
	stakingkeeper "github.com/Finschia/finschia-sdk/x/staking/keeper"
	stakingtypes "github.com/Finschia/finschia-sdk/x/staking/types"
	stakingplusmodule "github.com/Finschia/finschia-sdk/x/stakingplus/module"
	"github.com/Finschia/finschia-sdk/x/token"
	"github.com/Finschia/finschia-sdk/x/token/class"
	classkeeper "github.com/Finschia/finschia-sdk/x/token/class/keeper"
	tokenkeeper "github.com/Finschia/finschia-sdk/x/token/keeper"
	tokenmodule "github.com/Finschia/finschia-sdk/x/token/module"
	"github.com/Finschia/finschia-sdk/x/upgrade"
	upgradeclient "github.com/Finschia/finschia-sdk/x/upgrade/client"
	upgradekeeper "github.com/Finschia/finschia-sdk/x/upgrade/keeper"
	upgradetypes "github.com/Finschia/finschia-sdk/x/upgrade/types"
	"github.com/Finschia/finschia-sdk/x/zkauth"
	zkauthkeeper "github.com/Finschia/finschia-sdk/x/zkauth/keeper"
	zkauthtypes "github.com/Finschia/finschia-sdk/x/zkauth/types"

	// unnamed import of statik for swagger UI support
	_ "github.com/Finschia/finschia-sdk/client/docs/statik"
)

const appName = "SimApp"

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		stakingplusmodule.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		foundationmodule.AppModuleBasic{},
		gov.NewAppModuleBasic(
			paramsclient.ProposalHandler,
			distrclient.ProposalHandler,
			upgradeclient.ProposalHandler,
			upgradeclient.CancelProposalHandler,
			foundationclient.ProposalHandler,
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		vesting.AppModuleBasic{},
		tokenmodule.AppModuleBasic{},
		collectionmodule.AppModuleBasic{},
		zkauth.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		foundation.ModuleName:          nil,
		foundation.TreasuryName:        nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		// govtypes.ModuleName: true, // TODO: uncomment it when authority is ready
	}
)

var (
	_ App                     = (*SimApp)(nil)
	_ servertypes.Application = (*SimApp)(nil)
)

// SimApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type SimApp struct {
	*baseapp.BaseApp
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// keepers
	AccountKeeper    authkeeper.AccountKeeper
	BankKeeper       bankkeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	StakingKeeper    stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	MintKeeper       mintkeeper.Keeper
	DistrKeeper      distrkeeper.Keeper
	FoundationKeeper foundationkeeper.Keeper
	GovKeeper        govkeeper.Keeper
	CrisisKeeper     crisiskeeper.Keeper
	UpgradeKeeper    upgradekeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper
	AuthzKeeper      authzkeeper.Keeper
	EvidenceKeeper   evidencekeeper.Keeper
	FeeGrantKeeper   feegrantkeeper.Keeper
	ClassKeeper      classkeeper.Keeper
	TokenKeeper      tokenkeeper.Keeper
	CollectionKeeper collectionkeeper.Keeper
	ZKAuthKeeper     zkauthkeeper.Keeper

	// the module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager

	// the configurator
	configurator module.Configurator
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".simapp")
}

// NewSimApp returns a reference to an initialized SimApp.
func NewSimApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool,
	homePath string, invCheckPeriod uint, encodingConfig simappparams.EncodingConfig,
	appOpts servertypes.AppOptions, baseAppOptions ...func(*baseapp.BaseApp),
) *SimApp {
	appCodec := encodingConfig.Marshaler
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(appName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		minttypes.StoreKey,
		distrtypes.StoreKey,
		slashingtypes.StoreKey,
		govtypes.StoreKey,
		paramstypes.StoreKey,
		upgradetypes.StoreKey,
		evidencetypes.StoreKey,
		capabilitytypes.StoreKey,
		feegrant.StoreKey,
		foundation.StoreKey,
		class.StoreKey,
		token.StoreKey,
		collection.StoreKey,
		authzkeeper.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	// NOTE: The testingkey is just mounted for testing purposes. Actual applications should
	// not include this key.
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey, "testingkey")

	// configure state listening capabilities using AppOptions
	// we are doing nothing with the returned streamingServices and waitGroup in this case
	if _, _, err := streaming.LoadStreamingServices(bApp, appOpts, appCodec, keys); err != nil {
		ostos.Exit(err.Error())
	}

	app := &SimApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		memKeys:           memKeys,
	}

	app.ParamsKeeper = initParamsKeeper(appCodec, legacyAmino, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

	// set the BaseApp's parameter store
	bApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])

	// Applications that wish to enforce statically created ScopedKeepers should call `Seal` after creating
	// their scoped modules in `NewApp` with `ScopeToModule`
	app.CapabilityKeeper.Seal()

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], app.GetSubspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms,
	)
	app.BankKeeper = bankpluskeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], app.AccountKeeper, app.GetSubspace(banktypes.ModuleName), app.BlockedAddrs(), false)
	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec, keys[stakingtypes.StoreKey], app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName),
	)
	app.MintKeeper = mintkeeper.NewKeeper(
		appCodec, keys[minttypes.StoreKey], app.GetSubspace(minttypes.ModuleName), &stakingKeeper,
		app.AccountKeeper, app.BankKeeper, authtypes.FeeCollectorName,
	)
	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec, keys[distrtypes.StoreKey], app.GetSubspace(distrtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, authtypes.FeeCollectorName, app.ModuleAccountAddrs(),
	)
	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec, keys[slashingtypes.StoreKey], &stakingKeeper, app.GetSubspace(slashingtypes.ModuleName),
	)
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		app.GetSubspace(crisistypes.ModuleName), invCheckPeriod, app.BankKeeper, authtypes.FeeCollectorName,
	)

	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(appCodec, keys[feegrant.StoreKey], app.AccountKeeper)
	app.UpgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, keys[upgradetypes.StoreKey], appCodec, homePath, app.BaseApp)

	foundationConfig := foundation.DefaultConfig()
	app.FoundationKeeper = foundationkeeper.NewKeeper(appCodec, keys[foundation.StoreKey], app.BaseApp.MsgServiceRouter(), app.AccountKeeper, app.BankKeeper, authtypes.FeeCollectorName, foundationConfig, foundation.DefaultAuthority().String(), app.GetSubspace(foundation.ModuleName))

	app.ClassKeeper = classkeeper.NewKeeper(appCodec, keys[class.StoreKey])
	app.TokenKeeper = tokenkeeper.NewKeeper(appCodec, keys[token.StoreKey], app.ClassKeeper)
	app.CollectionKeeper = collectionkeeper.NewKeeper(appCodec, keys[collection.StoreKey], app.ClassKeeper)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(keys[authzkeeper.StoreKey], appCodec, app.BaseApp.MsgServiceRouter())

	// register the proposal types
	govRouter := govtypes.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(foundation.RouterKey, foundationkeeper.NewFoundationProposalsHandler(app.FoundationKeeper))

	govKeeper := govkeeper.NewKeeper(
		appCodec, keys[govtypes.StoreKey], app.GetSubspace(govtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, govRouter,
	)

	app.GovKeeper = *govKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
		// register the governance hooks
		),
	)

	// create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], &app.StakingKeeper, app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	// create zkauth keeper
	jwKsMap := zkauthtypes.NewJWKs()
	// todo: verification key should be loaded from file.
	var verificationKey = []byte("{\n \"protocol\": \"groth16\",\n \"curve\": \"bn128\",\n \"nPublic\": 1,\n \"vk_alpha_1\": [\n  \"20491192805390485299153009773594534940189261866228447918068658471970481763042\",\n  \"9383485363053290200918347156157836566562967994039712273449902621266178545958\",\n  \"1\"\n ],\n \"vk_beta_2\": [\n  [\n   \"6375614351688725206403948262868962793625744043794305715222011528459656738731\",\n   \"4252822878758300859123897981450591353533073413197771768651442665752259397132\"\n  ],\n  [\n   \"10505242626370262277552901082094356697409835680220590971873171140371331206856\",\n   \"21847035105528745403288232691147584728191162732299865338377159692350059136679\"\n  ],\n  [\n   \"1\",\n   \"0\"\n  ]\n ],\n \"vk_gamma_2\": [\n  [\n   \"10857046999023057135944570762232829481370756359578518086990519993285655852781\",\n   \"11559732032986387107991004021392285783925812861821192530917403151452391805634\"\n  ],\n  [\n   \"8495653923123431417604973247489272438418190587263600148770280649306958101930\",\n   \"4082367875863433681332203403145435568316851327593401208105741076214120093531\"\n  ],\n  [\n   \"1\",\n   \"0\"\n  ]\n ],\n \"vk_delta_2\": [\n  [\n   \"21349319915249622662700217004338779716430783387183352766870647565870141979289\",\n   \"8213816744021090866451311756048660670381089332123677295675725952502733471420\"\n  ],\n  [\n   \"4787213629490370557685854255230879988945206163033639129474026644007741911075\",\n   \"20003855859301921415178037270191878217707285640767940877063768682564788786247\"\n  ],\n  [\n   \"1\",\n   \"0\"\n  ]\n ],\n \"vk_alphabeta_12\": [\n  [\n   [\n    \"2029413683389138792403550203267699914886160938906632433982220835551125967885\",\n    \"21072700047562757817161031222997517981543347628379360635925549008442030252106\"\n   ],\n   [\n    \"5940354580057074848093997050200682056184807770593307860589430076672439820312\",\n    \"12156638873931618554171829126792193045421052652279363021382169897324752428276\"\n   ],\n   [\n    \"7898200236362823042373859371574133993780991612861777490112507062703164551277\",\n    \"7074218545237549455313236346927434013100842096812539264420499035217050630853\"\n   ]\n  ],\n  [\n   [\n    \"7077479683546002997211712695946002074877511277312570035766170199895071832130\",\n    \"10093483419865920389913245021038182291233451549023025229112148274109565435465\"\n   ],\n   [\n    \"4595479056700221319381530156280926371456704509942304414423590385166031118820\",\n    \"19831328484489333784475432780421641293929726139240675179672856274388269393268\"\n   ],\n   [\n    \"11934129596455521040620786944827826205713621633706285934057045369193958244500\",\n    \"8037395052364110730298837004334506829870972346962140206007064471173334027475\"\n   ]\n  ]\n ],\n \"IC\": [\n  [\n   \"801233197807402683764630185033839955156034586542543249813920835808534245147\",\n   \"13286420793149616228297035344471157585445615731792629462934831296345279687002\",\n   \"1\"\n  ],\n  [\n   \"17608180544527043978731301492557909061209088433544687588079992534282036547698\",\n   \"11240405619785894451348234456278767489162139374206168239508590931049712428392\",\n   \"1\"\n  ]\n ]\n}")
	zkAuthVerifier := zkauthtypes.NewZKAuthVerifier(verificationKey)
	zkauthKeeper := zkauthkeeper.NewKeeper(appCodec, keys[zkauthtypes.StoreKey], jwKsMap, zkAuthVerifier, app.MsgServiceRouter())
	app.ZKAuthKeeper = *zkauthKeeper

	// Fetch JWK
	var wg sync.WaitGroup
	wg.Add(1)

	ctx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})
	go app.ZKAuthKeeper.FetchJWK(ctx, &wg)

	wg.Wait()

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bankplus.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		foundationmodule.NewAppModule(appCodec, app.FoundationKeeper),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		stakingplusmodule.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.FoundationKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		params.NewAppModule(app.ParamsKeeper),
		tokenmodule.NewAppModule(appCodec, app.TokenKeeper),
		collectionmodule.NewAppModule(appCodec, app.CollectionKeeper),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		zkauth.NewAppModule(appCodec, app.ZKAuthKeeper, app.AccountKeeper, app.BankKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	// NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		minttypes.ModuleName,
		foundation.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		token.ModuleName,
		collection.ModuleName,
		zkauthtypes.ModuleName,
	)
	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		foundation.ModuleName,
		token.ModuleName,
		collection.ModuleName,
		zkauthtypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		foundation.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		token.ModuleName,
		collection.ModuleName,
		zkauthtypes.ModuleName,
	)

	// Uncomment if you want to set a custom migration order here.
	// app.mm.SetOrderMigrations(custom order)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)

	// add test gRPC service for testing gRPC queries in isolation
	testdata.RegisterQueryServer(app.GRPCQueryRouter(), testdata.QueryImpl{})

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	overrideModules := map[string]module.AppModuleSimulation{
		authtypes.ModuleName:    auth.NewAppModule(app.appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts),
		stakingtypes.ModuleName: staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
	}
	app.sm = module.NewSimulationManagerFromAppModules(app.mm.Modules, overrideModules)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)

	// todo: please use this code when integrate zkauth's anteHandler.
	anteHandler, err := appante.NewAnteHandler(
		appante.HandlerOptions{
			HandlerOptions: ante.HandlerOptions{
				AccountKeeper:   app.AccountKeeper,
				BankKeeper:      app.BankKeeper,
				FeegrantKeeper:  app.FeeGrantKeeper,
				SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			},
			ZKAuthKeeper: app.ZKAuthKeeper,
		},
	)
	// anteHandler, err := ante.NewAnteHandler(
	// 	ante.HandlerOptions{
	// 		AccountKeeper:   app.AccountKeeper,
	// 		BankKeeper:      app.BankKeeper,
	// 		SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
	// 		FeegrantKeeper:  app.FeeGrantKeeper,
	// 		SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
	// 	},
	// )
	if err != nil {
		panic(err)
	}

	app.SetAnteHandler(anteHandler)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			ostos.Exit(err.Error())
		}

		ctx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})
		app.BankKeeper.(bankpluskeeper.Keeper).InitializeBankPlus(ctx)
	}

	return app
}

// Name returns the name of the App
func (app *SimApp) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *SimApp) BeginBlocker(ctx sdk.Context, req ocabci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *SimApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *SimApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *SimApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *SimApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlockedAddrs returns all the app's module account addresses that are not
// allowed to receive external tokens.
func (app *SimApp) BlockedAddrs() map[string]bool {
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blockedAddrs
}

// GetModuleManager returns the app module manager
// NOTE: used for testing purposes
func (app *SimApp) GetModuleManager() *module.Manager {
	return app.mm
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *SimApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// TestingApp functions
// GetBaseApp implements the TestingApp interface.
func (app *SimApp) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

// GetStakingKeeper implements the TestingApp interface.
func (app *SimApp) GetStakingKeeper() stakingkeeper.Keeper {
	return app.StakingKeeper
}

// GetTxConfig implements the TestingApp interface.
func (app *SimApp) GetTxConfig() client.TxConfig {
	return MakeTestEncodingConfig().TxConfig
}

// AppCodec returns SimApp's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *SimApp) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns SimApp's InterfaceRegistry
func (app *SimApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *SimApp) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *SimApp) GetMemKey(storeKey string) *sdk.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *SimApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// SimulationManager implements the SimulationApp interface
func (app *SimApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *SimApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx

	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	authtx2.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	ocservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register node gRPC service for grpc-gateway.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *SimApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
	authtx2.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *SimApp) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
	ocservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

func (app *SimApp) RegisterNodeService(clientCtx client.Context) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(ctx client.Context, rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(foundation.ModuleName)

	return paramsKeeper
}
