package types_test

import (
	"testing"
	"time"

	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
	ocbytes "github.com/line/ostracon/libs/bytes"
	ocproto "github.com/line/ostracon/proto/ostracon/types"
	octypes "github.com/line/ostracon/types"
	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/simapp"
	clienttypes "github.com/line/lbm-sdk/x/ibc/core/02-client/types"
	ibcoctypes "github.com/line/lbm-sdk/x/ibc/light-clients/99-ostracon/types"
	ibctesting "github.com/line/lbm-sdk/x/ibc/testing"
	ibctestingmock "github.com/line/lbm-sdk/x/ibc/testing/mock"
)

const (
	chainID                        = "lbm"
	chainIDRevision0               = "lbm-revision-0"
	chainIDRevision1               = "lbm-revision-1"
	clientID                       = "lbmmainnet"
	trustingPeriod   time.Duration = time.Hour * 24 * 7 * 2
	ubdPeriod        time.Duration = time.Hour * 24 * 7 * 3
	maxClockDrift    time.Duration = time.Second * 10
)

var (
	height          = clienttypes.NewHeight(0, 4)
	newClientHeight = clienttypes.NewHeight(1, 1)
	upgradePath     = []string{"upgrade", "upgradedIBCState"}
)

type OstraconTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain

	// TODO: deprecate usage in favor of testing package
	ctx        sdk.Context
	cdc        codec.Codec
	privVal    octypes.PrivValidator
	valSet     *octypes.ValidatorSet
	valsHash   ocbytes.HexBytes
	header     *ibcoctypes.Header
	now        time.Time
	headerTime time.Time
	clientTime time.Time
}

func (suite *OstraconTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))
	// commit some blocks so that QueryProof returns valid proof (cannot return valid query if height <= 1)
	suite.coordinator.CommitNBlocks(suite.chainA, 2)
	suite.coordinator.CommitNBlocks(suite.chainB, 2)

	// TODO: deprecate usage in favor of testing package
	checkTx := false
	app := simapp.Setup(checkTx)

	suite.cdc = app.AppCodec()

	// now is the time of the current chain, must be after the updating header
	// mocks ctx.BlockTime()
	suite.now = time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	suite.clientTime = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	// Header time is intended to be time for any new header used for updates
	suite.headerTime = time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)

	suite.privVal = ibctestingmock.NewPV()

	pubKey, err := suite.privVal.GetPubKey()
	suite.Require().NoError(err)

	heightMinus1 := clienttypes.NewHeight(0, height.RevisionHeight-1)

	val := ibctesting.NewTestValidator(pubKey, 10)
	suite.valSet = octypes.NewValidatorSet([]*octypes.Validator{val})
	suite.valsHash = suite.valSet.Hash()
	voterSet := octypes.WrapValidatorsToVoterSet(suite.valSet.Validators)
	suite.header = suite.chainA.CreateOCClientHeader(chainID, int64(height.RevisionHeight), heightMinus1, suite.now, suite.valSet, suite.valSet, voterSet, voterSet, []octypes.PrivValidator{suite.privVal})
	suite.ctx = app.BaseApp.NewContext(checkTx, ocproto.Header{Height: 1, Time: suite.now})
}

func getSuiteSigners(suite *OstraconTestSuite) []octypes.PrivValidator {
	return []octypes.PrivValidator{suite.privVal}
}

func getBothSigners(suite *OstraconTestSuite, altVal *octypes.Validator, altPrivVal octypes.PrivValidator) (*octypes.ValidatorSet, []octypes.PrivValidator) {
	// Create bothValSet with both suite validator and altVal. Would be valid update
	bothValSet := octypes.NewValidatorSet(append(suite.valSet.Validators, altVal))
	// Create signer array and ensure it is in same order as bothValSet
	_, suiteVal := suite.valSet.GetByIndex(0)
	bothSigners := ibctesting.CreateSortedSignerArray(altPrivVal, suite.privVal, altVal, suiteVal)
	return bothValSet, bothSigners
}

func TestOstraconTestSuite(t *testing.T) {
	suite.Run(t, new(OstraconTestSuite))
}
