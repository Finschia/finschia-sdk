package keeper_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/line/lbm-sdk/baseapp"
	"github.com/line/lbm-sdk/codec"
	cryptocodec "github.com/line/lbm-sdk/crypto/codec"
	sdk "github.com/line/lbm-sdk/types"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
	ocbytes "github.com/line/ostracon/libs/bytes"
	ocproto "github.com/line/ostracon/proto/ostracon/types"
	octypes "github.com/line/ostracon/types"
	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/simapp"
	"github.com/line/lbm-sdk/x/ibc/core/02-client/keeper"
	"github.com/line/lbm-sdk/x/ibc/core/02-client/types"
	commitmenttypes "github.com/line/lbm-sdk/x/ibc/core/23-commitment/types"
	"github.com/line/lbm-sdk/x/ibc/core/exported"
	localhosttypes "github.com/line/lbm-sdk/x/ibc/light-clients/09-localhost/types"
	ibcoctypes "github.com/line/lbm-sdk/x/ibc/light-clients/99-ostracon/types"
	ibctesting "github.com/line/lbm-sdk/x/ibc/testing"
	ibctestingmock "github.com/line/lbm-sdk/x/ibc/testing/mock"
)

const (
	testChainID          = "gaiahub-0"
	testChainIDRevision1 = "gaiahub-1"

	testClientID  = "ostracon-0"
	testClientID2 = "ostracon-1"
	testClientID3 = "ostracon-2"

	height = 5

	trustingPeriod time.Duration = time.Hour * 24 * 7 * 2
	ubdPeriod      time.Duration = time.Hour * 24 * 7 * 3
	maxClockDrift  time.Duration = time.Second * 10
)

var (
	testClientHeight          = types.NewHeight(0, 5)
	testClientHeightRevision1 = types.NewHeight(1, 5)
	newClientHeight           = types.NewHeight(1, 1)
)

type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain

	cdc            codec.Codec
	ctx            sdk.Context
	keeper         *keeper.Keeper
	consensusState *ibcoctypes.ConsensusState
	header         *ibcoctypes.Header
	valSet         *octypes.ValidatorSet
	valSetHash     ocbytes.HexBytes
	privVal        octypes.PrivValidator
	now            time.Time
	past           time.Time

	// TODO: deprecate
	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)

	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))

	isCheckTx := false
	suite.now = time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	suite.past = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	now2 := suite.now.Add(time.Hour)
	app := simapp.Setup(isCheckTx)

	suite.cdc = app.AppCodec()
	suite.ctx = app.BaseApp.NewContext(isCheckTx, ocproto.Header{Height: height, ChainID: testClientID, Time: now2})
	suite.keeper = &app.IBCKeeper.ClientKeeper
	suite.privVal = ibctestingmock.NewPV()

	pubKey, err := suite.privVal.GetPubKey()
	suite.Require().NoError(err)

	testClientHeightMinus1 := types.NewHeight(0, height-1)

	validator := ibctesting.NewTestValidator(pubKey, 1)
	suite.valSet = octypes.NewValidatorSet([]*octypes.Validator{validator})
	suite.valSetHash = suite.valSet.Hash()
	voterSet := octypes.WrapValidatorsToVoterSet(suite.valSet.Validators)
	suite.header = suite.chainA.CreateTMClientHeader(testChainID, int64(testClientHeight.RevisionHeight), testClientHeightMinus1, now2, suite.valSet, suite.valSet, voterSet, voterSet, []octypes.PrivValidator{suite.privVal})
	suite.consensusState = ibcoctypes.NewConsensusState(suite.now, commitmenttypes.NewMerkleRoot([]byte("hash")), suite.valSetHash)

	var validators stakingtypes.Validators
	for i := 1; i < 11; i++ {
		privVal := ibctestingmock.NewPV()
		tmPk, err := privVal.GetPubKey()
		suite.Require().NoError(err)
		pk, err := cryptocodec.FromOcPubKeyInterface(tmPk)
		suite.Require().NoError(err)
		val, err := stakingtypes.NewValidator(sdk.ValAddress(pk.Address()), pk, stakingtypes.Description{})
		suite.Require().NoError(err)

		val.Status = stakingtypes.Bonded
		val.Tokens = sdk.NewInt(rand.Int63())
		validators = append(validators, val)

		hi := stakingtypes.NewHistoricalInfo(suite.ctx.BlockHeader(), validators, sdk.DefaultPowerReduction)
		app.StakingKeeper.SetHistoricalInfo(suite.ctx, int64(i), &hi)
	}

	// add localhost client
	revision := types.ParseChainID(suite.chainA.ChainID)
	localHostClient := localhosttypes.NewClientState(
		suite.chainA.ChainID, types.NewHeight(revision, uint64(suite.chainA.GetContext().BlockHeight())),
	)
	suite.chainA.App.GetIBCKeeper().ClientKeeper.SetClientState(suite.chainA.GetContext(), exported.Localhost, localHostClient)

	// TODO: deprecate
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.IBCKeeper.ClientKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestSetClientState() {
	clientState := ibcoctypes.NewClientState(testChainID, ibcoctypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false)
	suite.keeper.SetClientState(suite.ctx, testClientID, clientState)

	retrievedState, found := suite.keeper.GetClientState(suite.ctx, testClientID)
	suite.Require().True(found, "GetClientState failed")
	suite.Require().Equal(clientState, retrievedState, "Client states are not equal")
}

func (suite *KeeperTestSuite) TestSetClientConsensusState() {
	suite.keeper.SetClientConsensusState(suite.ctx, testClientID, testClientHeight, suite.consensusState)

	retrievedConsState, found := suite.keeper.GetClientConsensusState(suite.ctx, testClientID, testClientHeight)
	suite.Require().True(found, "GetConsensusState failed")

	tmConsState, ok := retrievedConsState.(*ibcoctypes.ConsensusState)
	suite.Require().True(ok)
	suite.Require().Equal(suite.consensusState, tmConsState, "ConsensusState not stored correctly")
}

func (suite *KeeperTestSuite) TestValidateSelfClient() {
	testClientHeight := types.NewHeight(0, uint64(suite.chainA.GetContext().BlockHeight()-1))

	testCases := []struct {
		name        string
		clientState exported.ClientState
		expPass     bool
	}{
		{
			"success",
			ibcoctypes.NewClientState(suite.chainA.ChainID, ibcoctypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
			true,
		},
		{
			"success with nil UpgradePath",
			ibcoctypes.NewClientState(suite.chainA.ChainID, ibcoctypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), nil, false, false),
			true,
		},
		{
			"invalid client type",
			localhosttypes.NewClientState(suite.chainA.ChainID, testClientHeight),
			false,
		},
		{
			"frozen client",
			&ibcoctypes.ClientState{suite.chainA.ChainID, ibcoctypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, testClientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false},
			false,
		},
		{
			"incorrect chainID",
			ibcoctypes.NewClientState("gaiatestnet", ibcoctypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
			false,
		},
		{
			"invalid client height",
			ibcoctypes.NewClientState(suite.chainA.ChainID, ibcoctypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.NewHeight(0, uint64(suite.chainA.GetContext().BlockHeight())), commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
			false,
		},
		{
			"invalid client revision",
			ibcoctypes.NewClientState(suite.chainA.ChainID, ibcoctypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeightRevision1, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
			false,
		},
		{
			"invalid proof specs",
			ibcoctypes.NewClientState(suite.chainA.ChainID, ibcoctypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, nil, ibctesting.UpgradePath, false, false),
			false,
		},
		{
			"invalid trust level",
			ibcoctypes.NewClientState(suite.chainA.ChainID, ibcoctypes.Fraction{0, 1}, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
			false,
		},
		{
			"invalid unbonding period",
			ibcoctypes.NewClientState(suite.chainA.ChainID, ibcoctypes.DefaultTrustLevel, trustingPeriod, ubdPeriod+10, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
			false,
		},
		{
			"invalid trusting period",
			ibcoctypes.NewClientState(suite.chainA.ChainID, ibcoctypes.DefaultTrustLevel, ubdPeriod+10, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
			false,
		},
		{
			"invalid upgrade path",
			ibcoctypes.NewClientState(suite.chainA.ChainID, ibcoctypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), []string{"bad", "upgrade", "path"}, false, false),
			false,
		},
	}

	for _, tc := range testCases {
		err := suite.chainA.App.GetIBCKeeper().ClientKeeper.ValidateSelfClient(suite.chainA.GetContext(), tc.clientState)
		if tc.expPass {
			suite.Require().NoError(err, "expected valid client for case: %s", tc.name)
		} else {
			suite.Require().Error(err, "expected invalid client for case: %s", tc.name)
		}
	}
}

func (suite KeeperTestSuite) TestGetAllGenesisClients() {
	clientIDs := []string{
		testClientID2, testClientID3, testClientID,
	}
	expClients := []exported.ClientState{
		ibcoctypes.NewClientState(testChainID, ibcoctypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
		ibcoctypes.NewClientState(testChainID, ibcoctypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
		ibcoctypes.NewClientState(testChainID, ibcoctypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
	}

	expGenClients := make(types.IdentifiedClientStates, len(expClients))

	for i := range expClients {
		suite.chainA.App.GetIBCKeeper().ClientKeeper.SetClientState(suite.chainA.GetContext(), clientIDs[i], expClients[i])
		expGenClients[i] = types.NewIdentifiedClientState(clientIDs[i], expClients[i])
	}

	// add localhost client
	localHostClient, found := suite.chainA.App.GetIBCKeeper().ClientKeeper.GetClientState(suite.chainA.GetContext(), exported.Localhost)
	suite.Require().True(found)
	expGenClients = append(expGenClients, types.NewIdentifiedClientState(exported.Localhost, localHostClient))

	genClients := suite.chainA.App.GetIBCKeeper().ClientKeeper.GetAllGenesisClients(suite.chainA.GetContext())

	suite.Require().Equal(expGenClients.Sort(), genClients)
}

func (suite KeeperTestSuite) TestGetAllGenesisMetadata() {
	expectedGenMetadata := []types.IdentifiedGenesisMetadata{
		types.NewIdentifiedGenesisMetadata(
			"99-ostracon-1",
			[]types.GenesisMetadata{
				types.NewGenesisMetadata(ibcoctypes.ProcessedTimeKey(types.NewHeight(0, 1)), []byte("foo")),
				types.NewGenesisMetadata(ibcoctypes.ProcessedTimeKey(types.NewHeight(0, 2)), []byte("bar")),
				types.NewGenesisMetadata(ibcoctypes.ProcessedTimeKey(types.NewHeight(0, 3)), []byte("baz")),
			},
		),
		types.NewIdentifiedGenesisMetadata(
			"clientB",
			[]types.GenesisMetadata{
				types.NewGenesisMetadata(ibcoctypes.ProcessedTimeKey(types.NewHeight(1, 100)), []byte("val1")),
				types.NewGenesisMetadata(ibcoctypes.ProcessedTimeKey(types.NewHeight(2, 300)), []byte("val2")),
			},
		),
	}

	genClients := []types.IdentifiedClientState{
		types.NewIdentifiedClientState("99-ostracon-1", &ibcoctypes.ClientState{}), types.NewIdentifiedClientState("clientB", &ibcoctypes.ClientState{}),
		types.NewIdentifiedClientState("clientC", &ibcoctypes.ClientState{}), types.NewIdentifiedClientState("clientD", &localhosttypes.ClientState{}),
	}

	suite.chainA.App.GetIBCKeeper().ClientKeeper.SetAllClientMetadata(suite.chainA.GetContext(), expectedGenMetadata)

	actualGenMetadata, err := suite.chainA.App.GetIBCKeeper().ClientKeeper.GetAllClientMetadata(suite.chainA.GetContext(), genClients)
	suite.Require().NoError(err, "get client metadata returned error unexpectedly")
	suite.Require().Equal(expectedGenMetadata, actualGenMetadata, "retrieved metadata is unexpected")
}

func (suite KeeperTestSuite) TestGetConsensusState() {
	suite.ctx = suite.ctx.WithBlockHeight(10)
	cases := []struct {
		name    string
		height  types.Height
		expPass bool
	}{
		{"zero height", types.ZeroHeight(), false},
		{"height > latest height", types.NewHeight(0, uint64(suite.ctx.BlockHeight())+1), false},
		{"latest height - 1", types.NewHeight(0, uint64(suite.ctx.BlockHeight())-1), true},
		{"latest height", types.GetSelfHeight(suite.ctx), true},
	}

	for i, tc := range cases {
		tc := tc
		cs, err := suite.keeper.GetSelfConsensusState(suite.ctx, tc.height)
		if tc.expPass {
			suite.Require().NoError(err, "Case %d should have passed: %s", i, tc.name)
			suite.Require().NotNil(cs, "Case %d should have passed: %s", i, tc.name)
		} else {
			suite.Require().Error(err, "Case %d should have failed: %s", i, tc.name)
			suite.Require().Nil(cs, "Case %d should have failed: %s", i, tc.name)
		}
	}
}

func (suite KeeperTestSuite) TestConsensusStateHelpers() {
	// initial setup
	clientState := ibcoctypes.NewClientState(testChainID, ibcoctypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false)

	suite.keeper.SetClientState(suite.ctx, testClientID, clientState)
	suite.keeper.SetClientConsensusState(suite.ctx, testClientID, testClientHeight, suite.consensusState)

	nextState := ibcoctypes.NewConsensusState(suite.now, commitmenttypes.NewMerkleRoot([]byte("next")), suite.valSetHash)

	testClientHeightPlus5 := types.NewHeight(0, height+5)
	voterSet := octypes.WrapValidatorsToVoterSet(suite.valSet.Validators)

	header := suite.chainA.CreateTMClientHeader(testClientID, int64(testClientHeightPlus5.RevisionHeight), testClientHeight, suite.header.Header.Time.Add(time.Minute),
		suite.valSet, suite.valSet, voterSet, voterSet, []octypes.PrivValidator{suite.privVal})

	// mock update functionality
	clientState.LatestHeight = header.GetHeight().(types.Height)
	suite.keeper.SetClientConsensusState(suite.ctx, testClientID, header.GetHeight(), nextState)
	suite.keeper.SetClientState(suite.ctx, testClientID, clientState)

	latest, ok := suite.keeper.GetLatestClientConsensusState(suite.ctx, testClientID)
	suite.Require().True(ok)
	suite.Require().Equal(nextState, latest, "Latest client not returned correctly")
}

// 2 clients in total are created on chainA. The first client is updated so it contains an initial consensus state
// and a consensus state at the update height.
func (suite KeeperTestSuite) TestGetAllConsensusStates() {
	path := ibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)

	clientState := path.EndpointA.GetClientState()
	expConsensusHeight0 := clientState.GetLatestHeight()
	consensusState0, ok := suite.chainA.GetConsensusState(path.EndpointA.ClientID, expConsensusHeight0)
	suite.Require().True(ok)

	// update client to create a second consensus state
	err := path.EndpointA.UpdateClient()
	suite.Require().NoError(err)

	clientState = path.EndpointA.GetClientState()
	expConsensusHeight1 := clientState.GetLatestHeight()
	suite.Require().True(expConsensusHeight1.GT(expConsensusHeight0))
	consensusState1, ok := suite.chainA.GetConsensusState(path.EndpointA.ClientID, expConsensusHeight1)
	suite.Require().True(ok)

	expConsensus := []exported.ConsensusState{
		consensusState0,
		consensusState1,
	}

	// create second client on chainA
	path2 := ibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path2)
	clientState = path2.EndpointA.GetClientState()

	expConsensusHeight2 := clientState.GetLatestHeight()
	consensusState2, ok := suite.chainA.GetConsensusState(path2.EndpointA.ClientID, expConsensusHeight2)
	suite.Require().True(ok)

	expConsensus2 := []exported.ConsensusState{consensusState2}

	expConsensusStates := types.ClientsConsensusStates{
		types.NewClientConsensusStates(path.EndpointA.ClientID, []types.ConsensusStateWithHeight{
			types.NewConsensusStateWithHeight(expConsensusHeight0.(types.Height), expConsensus[0]),
			types.NewConsensusStateWithHeight(expConsensusHeight1.(types.Height), expConsensus[1]),
		}),
		types.NewClientConsensusStates(path2.EndpointA.ClientID, []types.ConsensusStateWithHeight{
			types.NewConsensusStateWithHeight(expConsensusHeight2.(types.Height), expConsensus2[0]),
		}),
	}.Sort()

	consStates := suite.chainA.App.GetIBCKeeper().ClientKeeper.GetAllConsensusStates(suite.chainA.GetContext())
	suite.Require().Equal(expConsensusStates, consStates, "%s \n\n%s", expConsensusStates, consStates)
}
