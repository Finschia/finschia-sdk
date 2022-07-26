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
			valid:  true,
		},
		"insufficient tokens": {
			amount: s.balance.Add(sdk.OneInt()),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			err := s.keeper.Send(ctx, s.contractID, s.vendor, s.operator, tc.amount)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			vendorBalance := s.keeper.GetBalance(ctx, s.contractID, s.vendor)
			s.Require().True(s.balance.Sub(tc.amount).Equal(vendorBalance))

			operatorBalance := s.keeper.GetBalance(ctx, s.contractID, s.operator)
			s.Require().True(s.balance.Add(tc.amount).Equal(operatorBalance))
		})
	}
}

func (s *KeeperTestSuite) TestAuthorizeOperator() {
	dummyContractID := "fee1dead"
	_, err := s.keeper.GetClass(s.ctx, dummyContractID)
	s.Require().Error(err)

	contractDescriptions := map[string]string{
		s.contractID:    "valid",
		dummyContractID: "not-exists",
	}
	userDescriptions := map[string]string{
		s.vendor.String():   "vendor",
		s.operator.String(): "operator",
		s.customer.String(): "customer",
		s.stranger.String(): "stranger",
	}
	for id, idDesc := range contractDescriptions {
		for operatorStr, operatorDesc := range userDescriptions {
			operator, err := sdk.AccAddressFromBech32(operatorStr)
			s.Require().NoError(err)
			for fromStr, fromDesc := range userDescriptions {
				from, err := sdk.AccAddressFromBech32(fromStr)
				s.Require().NoError(err)
				name := fmt.Sprintf("ContractID: %s, Operator: %s, From: %s", idDesc, operatorDesc, fromDesc)
				s.Run(name, func() {
					ctx, _ := s.ctx.CacheContext()

					_, idErr := s.keeper.GetClass(ctx, id)
					_, authErr := s.keeper.GetAuthorization(ctx, id, from, operator)
					err := s.keeper.AuthorizeOperator(ctx, id, from, operator)
					if idErr == nil && authErr != nil {
						s.Require().NoError(err)
						_, authErr = s.keeper.GetAuthorization(ctx, id, from, operator)
						s.Require().NoError(authErr)
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
		s.contractID:    "valid",
		dummyContractID: "not-exists",
	}
	userDescriptions := map[string]string{
		s.vendor.String():   "vendor",
		s.operator.String(): "operator",
		s.customer.String(): "customer",
	}
	for id, idDesc := range contractDescriptions {
		for operatorStr, operatorDesc := range userDescriptions {
			operator, err := sdk.AccAddressFromBech32(operatorStr)
			s.Require().NoError(err)
			for fromStr, fromDesc := range userDescriptions {
				from, err := sdk.AccAddressFromBech32(fromStr)
				s.Require().NoError(err)
				name := fmt.Sprintf("ContractID: %s, Operator: %s, From: %s", idDesc, operatorDesc, fromDesc)
				s.Run(name, func() {
					ctx, _ := s.ctx.CacheContext()

					_, idErr := s.keeper.GetClass(ctx, id)
					_, authErr := s.keeper.GetAuthorization(ctx, id, from, operator)
					err := s.keeper.RevokeOperator(ctx, id, from, operator)
					if idErr == nil && authErr == nil {
						s.Require().NoError(err)
						_, authErr = s.keeper.GetAuthorization(ctx, id, from, operator)
						s.Require().Error(authErr)
					} else {
						s.Require().Error(err)
					}
				})
			}
		}
	}
}
