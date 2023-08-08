package keeper_test

import (
	"fmt"

	sdk "github.com/Finschia/finschia-rdk/types"
	"github.com/Finschia/finschia-rdk/x/token"
)

func (s *KeeperTestSuite) TestSend() {
	testCases := map[string]struct {
		amount sdk.Int
		err    error
	}{
		"valid send": {
			amount: sdk.OneInt(),
		},
		"insufficient tokens": {
			amount: s.balance.Add(sdk.OneInt()),
			err:    token.ErrInsufficientBalance,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			err := s.keeper.Send(ctx, s.contractID, s.vendor, s.operator, tc.amount)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			vendorBalance := s.keeper.GetBalance(ctx, s.contractID, s.vendor)
			s.Require().True(s.balance.Sub(tc.amount).Equal(vendorBalance))

			operatorBalance := s.keeper.GetBalance(ctx, s.contractID, s.operator)
			s.Require().True(s.balance.Add(tc.amount).Equal(operatorBalance))
		})
	}
}

func (s *KeeperTestSuite) TestAuthorizeOperator() {
	userDescriptions := map[string]string{
		s.vendor.String():   "vendor",
		s.operator.String(): "operator",
		s.customer.String(): "customer",
		s.stranger.String(): "stranger",
	}
	for operatorStr, operatorDesc := range userDescriptions {
		operator, err := sdk.AccAddressFromBech32(operatorStr)
		s.Require().NoError(err)
		for fromStr, fromDesc := range userDescriptions {
			from, err := sdk.AccAddressFromBech32(fromStr)
			s.Require().NoError(err)
			name := fmt.Sprintf("Operator: %s, From: %s", operatorDesc, fromDesc)
			s.Run(name, func() {
				ctx, _ := s.ctx.CacheContext()

				_, queryErr := s.keeper.GetAuthorization(ctx, s.contractID, from, operator)
				err := s.keeper.AuthorizeOperator(ctx, s.contractID, from, operator)
				if queryErr == nil { // authorize must fail
					s.Require().ErrorIs(err, token.ErrTokenAlreadyApproved)
				} else {
					s.Require().ErrorIs(queryErr, token.ErrTokenNotApproved)
					s.Require().NoError(err)
					_, queryErr := s.keeper.GetAuthorization(ctx, s.contractID, from, operator)
					s.Require().NoError(queryErr)
				}
			})
		}
	}
}

func (s *KeeperTestSuite) TestRevokeOperator() {
	userDescriptions := map[string]string{
		s.vendor.String():   "vendor",
		s.operator.String(): "operator",
		s.customer.String(): "customer",
	}
	for operatorStr, operatorDesc := range userDescriptions {
		operator, err := sdk.AccAddressFromBech32(operatorStr)
		s.Require().NoError(err)
		for fromStr, fromDesc := range userDescriptions {
			from, err := sdk.AccAddressFromBech32(fromStr)
			s.Require().NoError(err)
			name := fmt.Sprintf("Operator: %s, From: %s", operatorDesc, fromDesc)
			s.Run(name, func() {
				ctx, _ := s.ctx.CacheContext()

				_, queryErr := s.keeper.GetAuthorization(ctx, s.contractID, from, operator)
				err := s.keeper.RevokeOperator(ctx, s.contractID, from, operator)
				if queryErr != nil { // revoke must fail
					s.Require().ErrorIs(queryErr, token.ErrTokenNotApproved)
					s.Require().ErrorIs(err, token.ErrTokenNotApproved)
				} else {
					s.Require().NoError(err)
					_, queryErr := s.keeper.GetAuthorization(ctx, s.contractID, from, operator)
					s.Require().ErrorIs(queryErr, token.ErrTokenNotApproved)
				}
			})
		}
	}
}
