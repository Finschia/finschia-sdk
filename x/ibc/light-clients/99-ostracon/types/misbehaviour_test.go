package types_test

import (
	"time"

	"github.com/line/ostracon/crypto/tmhash"
	tmproto "github.com/line/ostracon/proto/ostracon/types"
	tmtypes "github.com/line/ostracon/types"

	clienttypes "github.com/line/lbm-sdk/x/ibc/core/02-client/types"
	"github.com/line/lbm-sdk/x/ibc/core/exported"
	"github.com/line/lbm-sdk/x/ibc/light-clients/99-ostracon/types"
	ibctesting "github.com/line/lbm-sdk/x/ibc/testing"
	ibctestingmock "github.com/line/lbm-sdk/x/ibc/testing/mock"
)

func (suite *OstraconTestSuite) TestMisbehaviour() {
	signers := []tmtypes.PrivValidator{suite.privVal}
	heightMinus1 := clienttypes.NewHeight(0, height.RevisionHeight-1)

	voterSet := tmtypes.WrapValidatorsToVoterSet(suite.valSet.Validators)
	misbehaviour := &types.Misbehaviour{
		Header1:  suite.header,
		Header2:  suite.chainA.CreateTMClientHeader(chainID, int64(height.RevisionHeight), heightMinus1, suite.now, suite.valSet, suite.valSet, voterSet, voterSet, signers),
		ClientId: clientID,
	}

	suite.Require().Equal(exported.Ostracon, misbehaviour.ClientType())
	suite.Require().Equal(clientID, misbehaviour.GetClientID())
}

func (suite *OstraconTestSuite) TestMisbehaviourValidateBasic() {
	altPrivVal := ibctestingmock.NewPV()
	altPubKey, err := altPrivVal.GetPubKey()
	suite.Require().NoError(err)

	revisionHeight := int64(height.RevisionHeight)

	altVal := ibctesting.NewTestValidator(altPubKey, revisionHeight)

	// Create bothValSet with both suite validator and altVal
	bothValSet := tmtypes.NewValidatorSet(append(suite.valSet.Validators, altVal))
	// Create alternative validator set with only altVal
	altValSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{altVal})

	signers := []tmtypes.PrivValidator{suite.privVal}

	// Create signer array and ensure it is in same order as bothValSet
	_, suiteVal := suite.valSet.GetByIndex(0)
	bothSigners := ibctesting.CreateSortedSignerArray(altPrivVal, suite.privVal, altVal, suiteVal)

	altSigners := []tmtypes.PrivValidator{altPrivVal}

	heightMinus1 := clienttypes.NewHeight(0, height.RevisionHeight-1)

	voterSet := tmtypes.WrapValidatorsToVoterSet(suite.valSet.Validators)
	bothVoterSet := tmtypes.WrapValidatorsToVoterSet(bothValSet.Validators)

	testCases := []struct {
		name                 string
		misbehaviour         *types.Misbehaviour
		malleateMisbehaviour func(misbehaviour *types.Misbehaviour) error
		expPass              bool
	}{
		{
			"valid fork misbehaviour, two headers at same height have different time",
			&types.Misbehaviour{
				Header1:  suite.header,
				Header2:  suite.chainA.CreateTMClientHeader(chainID, int64(height.RevisionHeight), heightMinus1, suite.now.Add(time.Minute), suite.valSet, suite.valSet, voterSet, voterSet, signers),
				ClientId: clientID,
			},
			func(misbehaviour *types.Misbehaviour) error { return nil },
			true,
		},
		{
			"valid time misbehaviour, both headers at different heights are at same time",
			&types.Misbehaviour{
				Header1:  suite.chainA.CreateTMClientHeader(chainID, int64(height.RevisionHeight), heightMinus1, suite.now.Add(time.Minute), suite.valSet, suite.valSet, voterSet, voterSet, signers),
				Header2:  suite.header,
				ClientId: clientID,
			},
			func(misbehaviour *types.Misbehaviour) error { return nil },
			true,
		},
		{
			"misbehaviour Header1 is nil",
			types.NewMisbehaviour(clientID, nil, suite.header),
			func(m *types.Misbehaviour) error { return nil },
			false,
		},
		{
			"misbehaviour Header2 is nil",
			types.NewMisbehaviour(clientID, suite.header, nil),
			func(m *types.Misbehaviour) error { return nil },
			false,
		},
		{
			"valid misbehaviour with different trusted headers",
			&types.Misbehaviour{
				Header1:  suite.header,
				Header2:  suite.chainA.CreateTMClientHeader(chainID, int64(height.RevisionHeight), clienttypes.NewHeight(0, height.RevisionHeight-3), suite.now.Add(time.Minute), suite.valSet, bothValSet, voterSet, bothVoterSet, signers),
				ClientId: clientID,
			},
			func(misbehaviour *types.Misbehaviour) error { return nil },
			true,
		},
		{
			"trusted height is 0 in Header1",
			&types.Misbehaviour{
				Header1:  suite.chainA.CreateTMClientHeader(chainID, int64(height.RevisionHeight), clienttypes.ZeroHeight(), suite.now.Add(time.Minute), suite.valSet, suite.valSet, voterSet, voterSet, signers),
				Header2:  suite.header,
				ClientId: clientID,
			},
			func(misbehaviour *types.Misbehaviour) error { return nil },
			false,
		},
		{
			"trusted height is 0 in Header2",
			&types.Misbehaviour{
				Header1:  suite.header,
				Header2:  suite.chainA.CreateTMClientHeader(chainID, int64(height.RevisionHeight), clienttypes.ZeroHeight(), suite.now.Add(time.Minute), suite.valSet, suite.valSet, voterSet, voterSet, signers),
				ClientId: clientID,
			},
			func(misbehaviour *types.Misbehaviour) error { return nil },
			false,
		},
		{
			"trusted valset is nil in Header1",
			&types.Misbehaviour{
				Header1:  suite.chainA.CreateTMClientHeader(chainID, int64(height.RevisionHeight), heightMinus1, suite.now.Add(time.Minute), suite.valSet, nil, voterSet, nil, signers),
				Header2:  suite.header,
				ClientId: clientID,
			},
			func(misbehaviour *types.Misbehaviour) error { return nil },
			false,
		},
		{
			"trusted valset is nil in Header2",
			&types.Misbehaviour{
				Header1:  suite.header,
				Header2:  suite.chainA.CreateTMClientHeader(chainID, int64(height.RevisionHeight), heightMinus1, suite.now.Add(time.Minute), suite.valSet, nil, voterSet, nil, signers),
				ClientId: clientID,
			},
			func(misbehaviour *types.Misbehaviour) error { return nil },
			false,
		},
		{
			"invalid client ID ",
			&types.Misbehaviour{
				Header1:  suite.header,
				Header2:  suite.chainA.CreateTMClientHeader(chainID, int64(height.RevisionHeight), heightMinus1, suite.now, suite.valSet, suite.valSet, voterSet, voterSet, signers),
				ClientId: "GAIA",
			},
			func(misbehaviour *types.Misbehaviour) error { return nil },
			false,
		},
		{
			"chainIDs do not match",
			&types.Misbehaviour{
				Header1:  suite.header,
				Header2:  suite.chainA.CreateTMClientHeader("ethermint", int64(height.RevisionHeight), heightMinus1, suite.now, suite.valSet, suite.valSet, voterSet, voterSet, signers),
				ClientId: clientID,
			},
			func(misbehaviour *types.Misbehaviour) error { return nil },
			false,
		},
		{
			"header2 height is greater",
			&types.Misbehaviour{
				Header1:  suite.header,
				Header2:  suite.chainA.CreateTMClientHeader(chainID, 6, clienttypes.NewHeight(0, height.RevisionHeight+1), suite.now, suite.valSet, suite.valSet, voterSet, voterSet, signers),
				ClientId: clientID,
			},
			func(misbehaviour *types.Misbehaviour) error { return nil },
			false,
		},
		{
			"header 1 doesn't have 2/3 majority",
			&types.Misbehaviour{
				Header1:  suite.chainA.CreateTMClientHeader(chainID, int64(height.RevisionHeight), heightMinus1, suite.now, bothValSet, suite.valSet, bothVoterSet, voterSet, bothSigners),
				Header2:  suite.header,
				ClientId: clientID,
			},
			func(misbehaviour *types.Misbehaviour) error {
				// voteSet contains only altVal which is less than 2/3 of total power (height/1height)
				wrongVoteSet := tmtypes.NewVoteSet(chainID, int64(misbehaviour.Header1.GetHeight().GetRevisionHeight()), 1, tmproto.PrecommitType, tmtypes.WrapValidatorsToVoterSet(altValSet.Validators))
				blockID, err := tmtypes.BlockIDFromProto(&misbehaviour.Header1.Commit.BlockID)
				if err != nil {
					return err
				}

				tmCommit, err := tmtypes.MakeCommit(*blockID, int64(misbehaviour.Header2.GetHeight().GetRevisionHeight()), misbehaviour.Header1.Commit.Round, wrongVoteSet, altSigners, suite.now)
				misbehaviour.Header1.Commit = tmCommit.ToProto()
				return err
			},
			false,
		},
		{
			"header 2 doesn't have 2/3 majority",
			&types.Misbehaviour{
				Header1:  suite.header,
				Header2:  suite.chainA.CreateTMClientHeader(chainID, int64(height.RevisionHeight), heightMinus1, suite.now, bothValSet, suite.valSet, bothVoterSet, voterSet, bothSigners),
				ClientId: clientID,
			},
			func(misbehaviour *types.Misbehaviour) error {
				// voteSet contains only altVal which is less than 2/3 of total power (height/1height)
				wrongVoteSet := tmtypes.NewVoteSet(chainID, int64(misbehaviour.Header2.GetHeight().GetRevisionHeight()), 1, tmproto.PrecommitType, tmtypes.WrapValidatorsToVoterSet(altValSet.Validators))
				blockID, err := tmtypes.BlockIDFromProto(&misbehaviour.Header2.Commit.BlockID)
				if err != nil {
					return err
				}

				tmCommit, err := tmtypes.MakeCommit(*blockID, int64(misbehaviour.Header2.GetHeight().GetRevisionHeight()), misbehaviour.Header2.Commit.Round, wrongVoteSet, altSigners, suite.now)
				misbehaviour.Header2.Commit = tmCommit.ToProto()
				return err
			},
			false,
		},
		{
			"validators sign off on wrong commit",
			&types.Misbehaviour{
				Header1:  suite.header,
				Header2:  suite.chainA.CreateTMClientHeader(chainID, int64(height.RevisionHeight), heightMinus1, suite.now, bothValSet, suite.valSet, bothVoterSet, voterSet, bothSigners),
				ClientId: clientID,
			},
			func(misbehaviour *types.Misbehaviour) error {
				tmBlockID := ibctesting.MakeBlockID(tmhash.Sum([]byte("other_hash")), 3, tmhash.Sum([]byte("other_partset")))
				misbehaviour.Header2.Commit.BlockID = tmBlockID.ToProto()
				return nil
			},
			false,
		},
	}

	for i, tc := range testCases {
		tc := tc

		err := tc.malleateMisbehaviour(tc.misbehaviour)
		suite.Require().NoError(err)

		if tc.expPass {
			suite.Require().NoError(tc.misbehaviour.ValidateBasic(), "valid test case %d failed: %s", i, tc.name)
		} else {
			suite.Require().Error(tc.misbehaviour.ValidateBasic(), "invalid test case %d passed: %s", i, tc.name)
		}
	}
}
