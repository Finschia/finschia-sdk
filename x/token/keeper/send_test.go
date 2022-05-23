package keeper_test

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
)

func (s *KeeperTestSuite) TestSend() {
	testCases := map[string]struct {
		amount sdk.Int
		valid  bool
	}{
		"valid send": {
			amount: sdk.OneInt(),
			valid: true,
		},
		"insufficient tokens": {
			amount: s.balance.Mul(sdk.NewInt(2)),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			err := s.keeper.Send(ctx, s.classID, s.vendor, s.operator, tc.amount)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			vendorBalance := s.keeper.GetBalance(ctx, s.classID, s.vendor)
			s.Require().True(s.balance.Sub(tc.amount).Equal(vendorBalance))

			operatorBalance := s.keeper.GetBalance(ctx, s.classID, s.operator)
			s.Require().True(s.balance.Add(tc.amount).Equal(operatorBalance))
		})
	}
}

func (s *KeeperTestSuite) TestAuthorizeOperator() {
	// no such a class
	authorization := s.keeper.GetAuthorization(s.ctx, "deadbeef", s.vendor, s.customer)
	s.Require().Nil(authorization)

	users := []sdk.AccAddress{s.vendor, s.operator, s.customer}
	for _, operator := range users {
		for _, from := range users {
			name := fmt.Sprintf("Operator: %s, From: %s", operator, from)
			s.Run(name, func() {
				ctx, _ := s.ctx.CacheContext()

				authorization := s.keeper.GetAuthorization(ctx, s.classID, from, operator)
				err := s.keeper.AuthorizeOperator(ctx, s.classID, from, operator)
				if authorization == nil {
					s.Require().NoError(err)
					authorization = s.keeper.GetAuthorization(ctx, s.classID, from, operator)
					s.Require().NotNil(authorization)
				} else {
					s.Require().Error(err)
				}
			})
		}
	}
}
