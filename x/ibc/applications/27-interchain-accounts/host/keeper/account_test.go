package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"

	icatypes "github.com/line/lbm-sdk/x/ibc/applications/27-interchain-accounts/types"
	ibctesting "github.com/line/lbm-sdk/x/ibc/testing"
)

func (suite *KeeperTestSuite) TestRegisterInterchainAccount() {
	suite.SetupTest()

	path := NewICAPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupConnections(path)

	//RegisterInterchainAccount
	err := SetupICAPath(path, TestOwnerAddress)
	suite.Require().NoError(err)

	portID, err := icatypes.NewControllerPortID(TestOwnerAddress)
	suite.Require().NoError(err)

	// Get the address of the interchain account stored in state during handshake step
	storedAddr, found := suite.chainB.GetSimApp().ICAHostKeeper.GetInterchainAccountAddress(suite.chainB.GetContext(), ibctesting.FirstConnectionID, portID)
	suite.Require().True(found)

	icaAddr, err := sdk.AccAddressFromBech32(storedAddr)
	suite.Require().NoError(err)

	// Check if account is created
	interchainAccount := suite.chainB.GetSimApp().AccountKeeper.GetAccount(suite.chainB.GetContext(), icaAddr)
	suite.Require().Equal(interchainAccount.GetAddress().String(), storedAddr)
}
