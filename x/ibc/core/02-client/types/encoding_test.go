package types_test

import (
	"github.com/line/lbm-sdk/x/ibc/core/02-client/types"
	ibcoctypes "github.com/line/lbm-sdk/x/ibc/light-clients/99-ostracon/types"
)

func (suite *TypesTestSuite) TestMarshalHeader() {

	cdc := suite.chainA.App.AppCodec()
	h := &ibcoctypes.Header{
		TrustedHeight: types.NewHeight(4, 100),
	}

	// marshal header
	bz, err := types.MarshalHeader(cdc, h)
	suite.Require().NoError(err)

	// unmarshal header
	newHeader, err := types.UnmarshalHeader(cdc, bz)
	suite.Require().NoError(err)

	suite.Require().Equal(h, newHeader)

	// use invalid bytes
	invalidHeader, err := types.UnmarshalHeader(cdc, []byte("invalid bytes"))
	suite.Require().Error(err)
	suite.Require().Nil(invalidHeader)

}
