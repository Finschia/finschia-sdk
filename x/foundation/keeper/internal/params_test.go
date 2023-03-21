package internal_test

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

func (s *KeeperTestSuite) TestUpdateParams() {
	ctx, _ := s.ctx.CacheContext()

	s.impl.SetParams(ctx, foundation.Params{
		FoundationTax: sdk.OneDec(),
	})

	// test update params
	// 1. remove the features
	//   a. foundation tax
	removingParams := foundation.Params{
		FoundationTax: sdk.ZeroDec(),
	}
	s.Require().NoError(removingParams.ValidateBasic())
	err := s.impl.UpdateParams(ctx, removingParams)
	s.Require().NoError(err)

	// check params
	s.Require().Equal(removingParams, s.impl.GetParams(ctx))
	s.Require().Equal(removingParams.FoundationTax, s.impl.GetFoundationTax(ctx))

	// 2. re-enable the features, which must fail
	//   a. foundation tax
	taxParams := foundation.Params{
		FoundationTax: sdk.OneDec(),
	}
	s.Require().NoError(taxParams.ValidateBasic())
	err = s.impl.UpdateParams(ctx, taxParams)
	s.Require().Error(err)
}
