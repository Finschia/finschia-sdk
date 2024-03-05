package internal_test

import (
	"cosmossdk.io/math"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

func (s *KeeperTestSuite) TestUpdateParams() {
	ctx, _ := s.ctx.CacheContext()

	s.impl.SetParams(ctx, foundation.Params{
		FoundationTax: math.LegacyOneDec(),
	})

	// test update params
	// 1. remove the features
	//   a. foundation tax
	removingParams := foundation.Params{
		FoundationTax: math.LegacyZeroDec(),
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
		FoundationTax: math.LegacyOneDec(),
	}
	s.Require().NoError(taxParams.ValidateBasic())
	err = s.impl.UpdateParams(ctx, taxParams)
	s.Require().Error(err)
}
