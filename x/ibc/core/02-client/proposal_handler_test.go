package client_test

import (
	sdk "github.com/line/lbm-sdk/types"
	distributiontypes "github.com/line/lbm-sdk/x/distribution/types"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
	client "github.com/line/lbm-sdk/x/ibc/core/02-client"
	clienttypes "github.com/line/lbm-sdk/x/ibc/core/02-client/types"
	"github.com/line/lbm-sdk/x/ibc/core/exported"
	ibctmtypes "github.com/line/lbm-sdk/x/ibc/light-clients/99-ostracon/types"
	ibctesting "github.com/line/lbm-sdk/x/ibc/testing"
)

func (suite *ClientTestSuite) TestNewClientUpdateProposalHandler() {
	var (
		content govtypes.Content
		err     error
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"valid update client proposal", func() {
				subject, _ := suite.coordinator.SetupClients(suite.chainA, suite.chainB, exported.Ostracon)
				subjectClientState := suite.chainA.GetClientState(subject)
				substitute, _ := suite.coordinator.SetupClients(suite.chainA, suite.chainB, exported.Ostracon)
				initialHeight := clienttypes.NewHeight(subjectClientState.GetLatestHeight().GetRevisionNumber(), subjectClientState.GetLatestHeight().GetRevisionHeight()+1)

				// update substitute twice
				suite.coordinator.UpdateClient(suite.chainA, suite.chainB, substitute, exported.Ostracon)
				suite.coordinator.UpdateClient(suite.chainA, suite.chainB, substitute, exported.Ostracon)
				substituteClientState := suite.chainA.GetClientState(substitute)

				tmClientState, ok := subjectClientState.(*ibctmtypes.ClientState)
				suite.Require().True(ok)
				tmClientState.AllowUpdateAfterMisbehaviour = true
				tmClientState.FrozenHeight = tmClientState.LatestHeight
				suite.chainA.App.IBCKeeper.ClientKeeper.SetClientState(suite.chainA.GetContext(), subject, tmClientState)

				// replicate changes to substitute (they must match)
				tmClientState, ok = substituteClientState.(*ibctmtypes.ClientState)
				suite.Require().True(ok)
				tmClientState.AllowUpdateAfterMisbehaviour = true
				suite.chainA.App.IBCKeeper.ClientKeeper.SetClientState(suite.chainA.GetContext(), substitute, tmClientState)

				content = clienttypes.NewClientUpdateProposal(ibctesting.Title, ibctesting.Description, subject, substitute, initialHeight)
			}, true,
		},
		{
			"nil proposal", func() {
				content = nil
			}, false,
		},
		{
			"unsupported proposal type", func() {
				content = distributiontypes.NewCommunityPoolSpendProposal(ibctesting.Title, ibctesting.Description, suite.chainA.SenderAccount.GetAddress(), sdk.NewCoins(sdk.NewCoin("communityfunds", sdk.NewInt(10))))
			}, false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest() // reset

			tc.malleate()

			proposalHandler := client.NewClientProposalHandler(suite.chainA.App.IBCKeeper.ClientKeeper)

			err = proposalHandler(suite.chainA.GetContext(), content)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}

}
