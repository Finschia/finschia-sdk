package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

func (s *KeeperTestSuite) TestUpdateParams() {
	ctx, _ := s.ctx.CacheContext()

	// add a dummy url to censorship list
	msgTypeURL := sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil))
	dummyURL := sdk.MsgTypeURL((*foundation.MsgFundTreasury)(nil))
	s.keeper.SetParams(ctx, foundation.Params{
		FoundationTax: sdk.OneDec(),
		CensoredMsgTypeUrls: []string{
			msgTypeURL,
			dummyURL,
		},
	})

	// check preconditions
	s.Require().True(s.keeper.IsCensoredMessage(ctx, msgTypeURL))
	_, err := s.keeper.GetAuthorization(ctx, s.stranger, msgTypeURL)
	s.Require().NoError(err)

	// test update params
	// 1. remove the features
	//   a. foundation tax
	//   b. message censorship
	removingParams := foundation.Params{
		FoundationTax: sdk.ZeroDec(),
		CensoredMsgTypeUrls: []string{
			dummyURL,
		},
	}
	s.Require().NoError(removingParams.ValidateBasic())
	err = s.keeper.UpdateParams(ctx, removingParams)
	s.Require().NoError(err)

	// check params
	s.Require().Equal(removingParams, s.keeper.GetParams(ctx))
	s.Require().Equal(removingParams.FoundationTax, s.keeper.GetFoundationTax(ctx))
	s.Require().False(s.keeper.IsCensoredMessage(ctx, msgTypeURL))

	// check authorizations
	_, err = s.keeper.GetAuthorization(ctx, s.stranger, msgTypeURL)
	s.Require().Error(err)

	// 2. re-enable the features, which must fail
	//   a. foundation tax
	taxParams := foundation.Params{
		FoundationTax: sdk.OneDec(),
	}
	s.Require().NoError(taxParams.ValidateBasic())
	err = s.keeper.UpdateParams(ctx, taxParams)
	s.Require().Error(err)

	//   b. message censorship
	msgParams := foundation.Params{
		FoundationTax: sdk.ZeroDec(),
		CensoredMsgTypeUrls: []string{
			msgTypeURL,
		},
	}
	s.Require().NoError(msgParams.ValidateBasic())
	err = s.keeper.UpdateParams(ctx, msgParams)
	s.Require().Error(err)
}
