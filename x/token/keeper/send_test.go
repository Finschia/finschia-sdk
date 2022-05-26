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
			amount: s.balance.Add(sdk.OneInt()),
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
	dummyContractID := "fee1dead"
	_, err := s.keeper.GetClass(s.ctx, dummyContractID)
	s.Require().Error(err)

	contractDescriptions := map[string]string{
		s.classID: "valid",
		dummyContractID: "not-exists",
	}
	userDescriptions := map[sdk.AccAddress]string{
		s.vendor: "vendor",
		s.operator: "operator",
		s.customer: "customer",
	}
	for id, idDesc := range contractDescriptions {
		for operator, operatorDesc := range userDescriptions {
			for from, fromDesc := range userDescriptions {
				name := fmt.Sprintf("ContractID: %s, Operator: %s, From: %s", idDesc, operatorDesc, fromDesc)
				s.Run(name, func() {
					ctx, _ := s.ctx.CacheContext()

					_, idErr := s.keeper.GetClass(ctx, id)
					authorization := s.keeper.GetAuthorization(ctx, id, from, operator)
					err := s.keeper.AuthorizeOperator(ctx, id, from, operator)
					if idErr == nil && authorization == nil {
						s.Require().NoError(err)
						authorization = s.keeper.GetAuthorization(ctx, id, from, operator)
						s.Require().NotNil(authorization)
					} else {
						s.Require().Error(err)
					}
				})
			}
		}
	}
}

func (s *KeeperTestSuite) TestRevokeOperator() {
	dummyContractID := "fee1dead"
	_, err := s.keeper.GetClass(s.ctx, dummyContractID)
	s.Require().Error(err)

	contractDescriptions := map[string]string{
		s.classID: "valid",
		dummyContractID: "not-exists",
	}
	userDescriptions := map[sdk.AccAddress]string{
		s.vendor: "vendor",
		s.operator: "operator",
		s.customer: "customer",
	}
	for id, idDesc := range contractDescriptions {
		for operator, operatorDesc := range userDescriptions {
			for from, fromDesc := range userDescriptions {
				name := fmt.Sprintf("ContractID: %s, Operator: %s, From: %s", idDesc, operatorDesc, fromDesc)
				s.Run(name, func() {
					ctx, _ := s.ctx.CacheContext()

					_, idErr := s.keeper.GetClass(ctx, id)
					authorization := s.keeper.GetAuthorization(ctx, id, from, operator)
					err := s.keeper.RevokeOperator(ctx, id, from, operator)
					if idErr == nil && authorization != nil {
						s.Require().NoError(err)
						authorization = s.keeper.GetAuthorization(ctx, id, from, operator)
						s.Require().Nil(authorization)
					} else {
						s.Require().Error(err)
					}
				})
			}
		}
	}
}
