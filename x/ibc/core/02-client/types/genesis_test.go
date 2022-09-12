package types_test

import (
	"time"

	octypes "github.com/line/ostracon/types"

	client "github.com/line/lbm-sdk/x/ibc/core/02-client"
	"github.com/line/lbm-sdk/x/ibc/core/02-client/types"
	commitmenttypes "github.com/line/lbm-sdk/x/ibc/core/23-commitment/types"
	"github.com/line/lbm-sdk/x/ibc/core/exported"
	localhosttypes "github.com/line/lbm-sdk/x/ibc/light-clients/09-localhost/types"
	ibcoctypes "github.com/line/lbm-sdk/x/ibc/light-clients/99-ostracon/types"
	ibctesting "github.com/line/lbm-sdk/x/ibc/testing"
	ibctestingmock "github.com/line/lbm-sdk/x/ibc/testing/mock"
)

const (
	chainID         = "chainID"
	tmClientID0     = "99-ostracon-0"
	tmClientID1     = "99-ostracon-1"
	invalidClientID = "myclient-0"
	clientID        = tmClientID0

	height = 10
)

var clientHeight = types.NewHeight(0, 10)

func (suite *TypesTestSuite) TestMarshalGenesisState() {
	cdc := suite.chainA.App.AppCodec()
	path := ibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(path)
	err := path.EndpointA.UpdateClient()
	suite.Require().NoError(err)

	genesis := client.ExportGenesis(suite.chainA.GetContext(), suite.chainA.App.GetIBCKeeper().ClientKeeper)

	bz, err := cdc.MarshalJSON(&genesis)
	suite.Require().NoError(err)
	suite.Require().NotNil(bz)

	var gs types.GenesisState
	err = cdc.UnmarshalJSON(bz, &gs)
	suite.Require().NoError(err)
}

func (suite *TypesTestSuite) TestValidateGenesis() {
	privVal := ibctestingmock.NewPV()
	pubKey, err := privVal.GetPubKey()
	suite.Require().NoError(err)

	now := time.Now().UTC()

	val := ibctesting.NewTestValidator(pubKey, 10)
	valSet := octypes.NewValidatorSet([]*octypes.Validator{val})
	voterSet := octypes.WrapValidatorsToVoterSet(valSet.Validators)

	heightMinus1 := types.NewHeight(0, height-1)
	header := suite.chainA.CreateTMClientHeader(chainID, int64(clientHeight.RevisionHeight), heightMinus1, now, valSet, valSet, voterSet, voterSet, []octypes.PrivValidator{privVal})

	testCases := []struct {
		name     string
		genState types.GenesisState
		expPass  bool
	}{
		{
			name:     "default",
			genState: types.DefaultGenesisState(),
			expPass:  true,
		},
		{
			name: "valid custom genesis",
			genState: types.NewGenesisState(
				[]types.IdentifiedClientState{
					types.NewIdentifiedClientState(
						tmClientID0, ibcoctypes.NewClientState(chainID, ibctesting.DefaultTrustLevel, ibctesting.TrustingPeriod, ibctesting.UnbondingPeriod, ibctesting.MaxClockDrift, clientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
					),
					types.NewIdentifiedClientState(
						exported.Localhost+"-1", localhosttypes.NewClientState("chainID", clientHeight),
					),
				},
				[]types.ClientConsensusStates{
					types.NewClientConsensusStates(
						tmClientID0,
						[]types.ConsensusStateWithHeight{
							types.NewConsensusStateWithHeight(
								header.GetHeight().(types.Height),
								ibcoctypes.NewConsensusState(
									header.GetTime(), commitmenttypes.NewMerkleRoot(header.Header.GetAppHash()), header.Header.NextValidatorsHash,
								),
							),
						},
					),
				},
				[]types.IdentifiedGenesisMetadata{
					types.NewIdentifiedGenesisMetadata(
						clientID,
						[]types.GenesisMetadata{
							types.NewGenesisMetadata([]byte("key1"), []byte("val1")),
							types.NewGenesisMetadata([]byte("key2"), []byte("val2")),
						},
					),
				},
				types.NewParams(exported.Ostracon, exported.Localhost),
				false,
				2,
			),
			expPass: true,
		},
		{
			name: "invalid clientid",
			genState: types.NewGenesisState(
				[]types.IdentifiedClientState{
					types.NewIdentifiedClientState(
						invalidClientID, ibcoctypes.NewClientState(chainID, ibcoctypes.DefaultTrustLevel, ibctesting.TrustingPeriod, ibctesting.UnbondingPeriod, ibctesting.MaxClockDrift, clientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
					),
					types.NewIdentifiedClientState(
						exported.Localhost, localhosttypes.NewClientState("chainID", clientHeight),
					),
				},
				[]types.ClientConsensusStates{
					types.NewClientConsensusStates(
						invalidClientID,
						[]types.ConsensusStateWithHeight{
							types.NewConsensusStateWithHeight(
								header.GetHeight().(types.Height),
								ibcoctypes.NewConsensusState(
									header.GetTime(), commitmenttypes.NewMerkleRoot(header.Header.GetAppHash()), header.Header.NextValidatorsHash,
								),
							),
						},
					),
				},
				nil,
				types.NewParams(exported.Ostracon),
				false,
				0,
			),
			expPass: false,
		},
		{
			name: "invalid client",
			genState: types.NewGenesisState(
				[]types.IdentifiedClientState{
					types.NewIdentifiedClientState(
						tmClientID0, ibcoctypes.NewClientState(chainID, ibcoctypes.DefaultTrustLevel, ibctesting.TrustingPeriod, ibctesting.UnbondingPeriod, ibctesting.MaxClockDrift, clientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
					),
					types.NewIdentifiedClientState(exported.Localhost, localhosttypes.NewClientState("chaindID", types.ZeroHeight())),
				},
				nil,
				nil,
				types.NewParams(exported.Ostracon),
				false,
				0,
			),
			expPass: false,
		},
		{
			name: "consensus state client id does not match client id in genesis clients",
			genState: types.NewGenesisState(
				[]types.IdentifiedClientState{
					types.NewIdentifiedClientState(
						tmClientID0, ibcoctypes.NewClientState(chainID, ibcoctypes.DefaultTrustLevel, ibctesting.TrustingPeriod, ibctesting.UnbondingPeriod, ibctesting.MaxClockDrift, clientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
					),
					types.NewIdentifiedClientState(
						exported.Localhost, localhosttypes.NewClientState("chaindID", clientHeight),
					),
				},
				[]types.ClientConsensusStates{
					types.NewClientConsensusStates(
						tmClientID1,
						[]types.ConsensusStateWithHeight{
							types.NewConsensusStateWithHeight(
								types.NewHeight(0, 1),
								ibcoctypes.NewConsensusState(
									header.GetTime(), commitmenttypes.NewMerkleRoot(header.Header.GetAppHash()), header.Header.NextValidatorsHash,
								),
							),
						},
					),
				},
				nil,
				types.NewParams(exported.Ostracon),
				false,
				0,
			),
			expPass: false,
		},
		{
			name: "invalid consensus state height",
			genState: types.NewGenesisState(
				[]types.IdentifiedClientState{
					types.NewIdentifiedClientState(
						tmClientID0, ibcoctypes.NewClientState(chainID, ibcoctypes.DefaultTrustLevel, ibctesting.TrustingPeriod, ibctesting.UnbondingPeriod, ibctesting.MaxClockDrift, clientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
					),
					types.NewIdentifiedClientState(
						exported.Localhost, localhosttypes.NewClientState("chaindID", clientHeight),
					),
				},
				[]types.ClientConsensusStates{
					types.NewClientConsensusStates(
						tmClientID0,
						[]types.ConsensusStateWithHeight{
							types.NewConsensusStateWithHeight(
								types.ZeroHeight(),
								ibcoctypes.NewConsensusState(
									header.GetTime(), commitmenttypes.NewMerkleRoot(header.Header.GetAppHash()), header.Header.NextValidatorsHash,
								),
							),
						},
					),
				},
				nil,
				types.NewParams(exported.Ostracon),
				false,
				0,
			),
			expPass: false,
		},
		{
			name: "invalid consensus state",
			genState: types.NewGenesisState(
				[]types.IdentifiedClientState{
					types.NewIdentifiedClientState(
						tmClientID0, ibcoctypes.NewClientState(chainID, ibcoctypes.DefaultTrustLevel, ibctesting.TrustingPeriod, ibctesting.UnbondingPeriod, ibctesting.MaxClockDrift, clientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
					),
					types.NewIdentifiedClientState(
						exported.Localhost, localhosttypes.NewClientState("chaindID", clientHeight),
					),
				},
				[]types.ClientConsensusStates{
					types.NewClientConsensusStates(
						tmClientID0,
						[]types.ConsensusStateWithHeight{
							types.NewConsensusStateWithHeight(
								types.NewHeight(0, 1),
								ibcoctypes.NewConsensusState(
									time.Time{}, commitmenttypes.NewMerkleRoot(header.Header.GetAppHash()), header.Header.NextValidatorsHash,
								),
							),
						},
					),
				},
				nil,
				types.NewParams(exported.Ostracon),
				false,
				0,
			),
			expPass: false,
		},
		{
			name: "client in genesis clients is disallowed by params",
			genState: types.NewGenesisState(
				[]types.IdentifiedClientState{
					types.NewIdentifiedClientState(
						tmClientID0, ibcoctypes.NewClientState(chainID, ibctesting.DefaultTrustLevel, ibctesting.TrustingPeriod, ibctesting.UnbondingPeriod, ibctesting.MaxClockDrift, clientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
					),
					types.NewIdentifiedClientState(
						exported.Localhost, localhosttypes.NewClientState("chainID", clientHeight),
					),
				},
				[]types.ClientConsensusStates{
					types.NewClientConsensusStates(
						tmClientID0,
						[]types.ConsensusStateWithHeight{
							types.NewConsensusStateWithHeight(
								header.GetHeight().(types.Height),
								ibcoctypes.NewConsensusState(
									header.GetTime(), commitmenttypes.NewMerkleRoot(header.Header.GetAppHash()), header.Header.NextValidatorsHash,
								),
							),
						},
					),
				},
				nil,
				types.NewParams(exported.Solomachine),
				false,
				0,
			),
			expPass: false,
		},
		{
			name: "metadata client-id does not match a genesis client",
			genState: types.NewGenesisState(
				[]types.IdentifiedClientState{
					types.NewIdentifiedClientState(
						clientID, ibcoctypes.NewClientState(chainID, ibctesting.DefaultTrustLevel, ibctesting.TrustingPeriod, ibctesting.UnbondingPeriod, ibctesting.MaxClockDrift, clientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
					),
					types.NewIdentifiedClientState(
						exported.Localhost, localhosttypes.NewClientState("chainID", clientHeight),
					),
				},
				[]types.ClientConsensusStates{
					types.NewClientConsensusStates(
						clientID,
						[]types.ConsensusStateWithHeight{
							types.NewConsensusStateWithHeight(
								header.GetHeight().(types.Height),
								ibcoctypes.NewConsensusState(
									header.GetTime(), commitmenttypes.NewMerkleRoot(header.Header.GetAppHash()), header.Header.NextValidatorsHash,
								),
							),
						},
					),
				},
				[]types.IdentifiedGenesisMetadata{
					types.NewIdentifiedGenesisMetadata(
						"wrongclientid",
						[]types.GenesisMetadata{
							types.NewGenesisMetadata([]byte("key1"), []byte("val1")),
							types.NewGenesisMetadata([]byte("key2"), []byte("val2")),
						},
					),
				},
				types.NewParams(exported.Ostracon, exported.Localhost),
				false,
				0,
			),
			expPass: false,
		},
		{
			name: "invalid metadata",
			genState: types.NewGenesisState(
				[]types.IdentifiedClientState{
					types.NewIdentifiedClientState(
						clientID, ibcoctypes.NewClientState(chainID, ibcoctypes.DefaultTrustLevel, ibctesting.TrustingPeriod, ibctesting.UnbondingPeriod, ibctesting.MaxClockDrift, clientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
					),
				},
				[]types.ClientConsensusStates{
					types.NewClientConsensusStates(
						clientID,
						[]types.ConsensusStateWithHeight{
							types.NewConsensusStateWithHeight(
								header.GetHeight().(types.Height),
								ibcoctypes.NewConsensusState(
									header.GetTime(), commitmenttypes.NewMerkleRoot(header.Header.GetAppHash()), header.Header.NextValidatorsHash,
								),
							),
						},
					),
				},
				[]types.IdentifiedGenesisMetadata{
					types.NewIdentifiedGenesisMetadata(
						clientID,
						[]types.GenesisMetadata{
							types.NewGenesisMetadata([]byte(""), []byte("val1")),
							types.NewGenesisMetadata([]byte("key2"), []byte("val2")),
						},
					),
				},
				types.NewParams(exported.Ostracon),
				false,
				0,
			),
		},
		{
			name: "invalid params",
			genState: types.NewGenesisState(
				[]types.IdentifiedClientState{
					types.NewIdentifiedClientState(
						tmClientID0, ibcoctypes.NewClientState(chainID, ibctesting.DefaultTrustLevel, ibctesting.TrustingPeriod, ibctesting.UnbondingPeriod, ibctesting.MaxClockDrift, clientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
					),
					types.NewIdentifiedClientState(
						exported.Localhost, localhosttypes.NewClientState("chainID", clientHeight),
					),
				},
				[]types.ClientConsensusStates{
					types.NewClientConsensusStates(
						tmClientID0,
						[]types.ConsensusStateWithHeight{
							types.NewConsensusStateWithHeight(
								header.GetHeight().(types.Height),
								ibcoctypes.NewConsensusState(
									header.GetTime(), commitmenttypes.NewMerkleRoot(header.Header.GetAppHash()), header.Header.NextValidatorsHash,
								),
							),
						},
					),
				},
				nil,
				types.NewParams(" "),
				false,
				0,
			),
			expPass: false,
		},
		{
			name: "invalid param",
			genState: types.NewGenesisState(
				[]types.IdentifiedClientState{
					types.NewIdentifiedClientState(
						tmClientID0, ibcoctypes.NewClientState(chainID, ibctesting.DefaultTrustLevel, ibctesting.TrustingPeriod, ibctesting.UnbondingPeriod, ibctesting.MaxClockDrift, clientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
					),
					types.NewIdentifiedClientState(
						exported.Localhost, localhosttypes.NewClientState("chainID", clientHeight),
					),
				},
				[]types.ClientConsensusStates{
					types.NewClientConsensusStates(
						tmClientID0,
						[]types.ConsensusStateWithHeight{
							types.NewConsensusStateWithHeight(
								header.GetHeight().(types.Height),
								ibcoctypes.NewConsensusState(
									header.GetTime(), commitmenttypes.NewMerkleRoot(header.Header.GetAppHash()), header.Header.NextValidatorsHash,
								),
							),
						},
					),
				},
				nil,
				types.NewParams(" "),
				true,
				0,
			),
			expPass: false,
		},
		{
			name: "localhost client not registered on allowlist",
			genState: types.NewGenesisState(
				[]types.IdentifiedClientState{
					types.NewIdentifiedClientState(
						tmClientID1, ibcoctypes.NewClientState(chainID, ibctesting.DefaultTrustLevel, ibctesting.TrustingPeriod, ibctesting.UnbondingPeriod, ibctesting.MaxClockDrift, clientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
					),
					types.NewIdentifiedClientState(
						exported.Localhost+"-0", localhosttypes.NewClientState("chainID", clientHeight),
					),
				},
				[]types.ClientConsensusStates{
					types.NewClientConsensusStates(
						tmClientID1,
						[]types.ConsensusStateWithHeight{
							types.NewConsensusStateWithHeight(
								header.GetHeight().(types.Height),
								ibcoctypes.NewConsensusState(
									header.GetTime(), commitmenttypes.NewMerkleRoot(header.Header.GetAppHash()), header.Header.NextValidatorsHash,
								),
							),
						},
					),
				},
				nil,
				types.NewParams(exported.Ostracon),
				true,
				2,
			),
			expPass: false,
		},
		{
			name: "next sequence too small",
			genState: types.NewGenesisState(
				[]types.IdentifiedClientState{
					types.NewIdentifiedClientState(
						tmClientID0, ibcoctypes.NewClientState(chainID, ibctesting.DefaultTrustLevel, ibctesting.TrustingPeriod, ibctesting.UnbondingPeriod, ibctesting.MaxClockDrift, clientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
					),
					types.NewIdentifiedClientState(
						exported.Localhost+"-1", localhosttypes.NewClientState("chainID", clientHeight),
					),
				},
				[]types.ClientConsensusStates{
					types.NewClientConsensusStates(
						tmClientID0,
						[]types.ConsensusStateWithHeight{
							types.NewConsensusStateWithHeight(
								header.GetHeight().(types.Height),
								ibcoctypes.NewConsensusState(
									header.GetTime(), commitmenttypes.NewMerkleRoot(header.Header.GetAppHash()), header.Header.NextValidatorsHash,
								),
							),
						},
					),
				},
				nil,
				types.NewParams(exported.Ostracon, exported.Localhost),
				false,
				0,
			),
			expPass: false,
		},
		{
			name: "failed to parse client identifier in client state loop",
			genState: types.NewGenesisState(
				[]types.IdentifiedClientState{
					types.NewIdentifiedClientState(
						"my-client", ibcoctypes.NewClientState(chainID, ibctesting.DefaultTrustLevel, ibctesting.TrustingPeriod, ibctesting.UnbondingPeriod, ibctesting.MaxClockDrift, clientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
					),
					types.NewIdentifiedClientState(
						exported.Localhost+"-1", localhosttypes.NewClientState("chainID", clientHeight),
					),
				},
				[]types.ClientConsensusStates{
					types.NewClientConsensusStates(
						tmClientID0,
						[]types.ConsensusStateWithHeight{
							types.NewConsensusStateWithHeight(
								header.GetHeight().(types.Height),
								ibcoctypes.NewConsensusState(
									header.GetTime(), commitmenttypes.NewMerkleRoot(header.Header.GetAppHash()), header.Header.NextValidatorsHash,
								),
							),
						},
					),
				},
				nil,
				types.NewParams(exported.Ostracon, exported.Localhost),
				false,
				5,
			),
			expPass: false,
		},
		{
			name: "consensus state different than client state type",
			genState: types.NewGenesisState(
				[]types.IdentifiedClientState{
					types.NewIdentifiedClientState(
						exported.Localhost+"-1", localhosttypes.NewClientState("chainID", clientHeight),
					),
				},
				[]types.ClientConsensusStates{
					types.NewClientConsensusStates(
						exported.Localhost+"-1",
						[]types.ConsensusStateWithHeight{
							types.NewConsensusStateWithHeight(
								header.GetHeight().(types.Height),
								ibcoctypes.NewConsensusState(
									header.GetTime(), commitmenttypes.NewMerkleRoot(header.Header.GetAppHash()), header.Header.NextValidatorsHash,
								),
							),
						},
					),
				},
				nil,
				types.NewParams(exported.Ostracon, exported.Localhost),
				false,
				5,
			),
			expPass: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		err := tc.genState.Validate()
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}
