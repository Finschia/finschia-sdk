package keeper_test

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/collection"
)

func (s *KeeperTestSuite) TestSendCoins() {
	testCases := map[string]struct {
		amount collection.Coin
		valid  bool
	}{
		"valid send (fungible token)": {
			amount: collection.NewFTCoin(s.ftClassID, s.balance),
			valid:  true,
		},
		"valid send (non-fungible token)": {
			amount: collection.NewNFTCoin(s.nftClassID, 1),
			valid:  true,
		},
		"insufficient tokens": {
			amount: collection.NewFTCoin(s.ftClassID, s.balance.Add(sdk.OneInt())),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			tokenID := tc.amount.TokenId
			customerBalance := s.keeper.GetBalance(ctx, s.contractID, s.customer, tokenID)
			operatorBalance := s.keeper.GetBalance(ctx, s.contractID, s.operator, tokenID)

			err := s.keeper.SendCoins(ctx, s.contractID, s.customer, s.operator, collection.NewCoins(tc.amount))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			newCustomerBalance := s.keeper.GetBalance(ctx, s.contractID, s.customer, tokenID)
			newOperatorBalance := s.keeper.GetBalance(ctx, s.contractID, s.operator, tokenID)
			s.Require().True(customerBalance.Sub(tc.amount.Amount).Equal(newCustomerBalance))
			s.Require().True(operatorBalance.Add(tc.amount.Amount).Equal(newOperatorBalance))
		})
	}
}

func (s *KeeperTestSuite) TestAuthorizeOperator() {
	// make sure the dummy contract does not exist
	dummyContractID := "deadbeef"
	_, err := s.keeper.GetContract(s.ctx, dummyContractID)
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
		for operator, operatorDesc := range userDescriptions {
			for from, fromDesc := range userDescriptions {
				name := fmt.Sprintf("ContractID: %s, Operator: %s, From: %s", idDesc, operatorDesc, fromDesc)
				s.Run(name, func() {
					ctx, _ := s.ctx.CacheContext()

					_, idErr := s.keeper.GetContract(ctx, id)

					fromAddr, err := sdk.AccAddressFromBech32(from)
					s.Require().NoError(err)
					operatorAddr, err := sdk.AccAddressFromBech32(operator)
					s.Require().NoError(err)

					_, authErr := s.keeper.GetAuthorization(ctx, id, fromAddr, operatorAddr)
					err = s.keeper.AuthorizeOperator(ctx, id, fromAddr, operatorAddr)
					if idErr == nil && authErr != nil {
						s.Require().NoError(err)
						_, authErr = s.keeper.GetAuthorization(ctx, id, fromAddr, operatorAddr)
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
	// make sure the dummy contract does not exist
	dummyContractID := "deadbeef"
	_, err := s.keeper.GetContract(s.ctx, dummyContractID)
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
		for operator, operatorDesc := range userDescriptions {
			for from, fromDesc := range userDescriptions {
				name := fmt.Sprintf("ContractID: %s, Operator: %s, From: %s", idDesc, operatorDesc, fromDesc)
				s.Run(name, func() {
					ctx, _ := s.ctx.CacheContext()

					_, idErr := s.keeper.GetContract(ctx, id)

					fromAddr, err := sdk.AccAddressFromBech32(from)
					s.Require().NoError(err)
					operatorAddr, err := sdk.AccAddressFromBech32(operator)
					s.Require().NoError(err)

					_, authErr := s.keeper.GetAuthorization(ctx, id, fromAddr, operatorAddr)
					err = s.keeper.RevokeOperator(ctx, id, fromAddr, operatorAddr)
					if idErr == nil && authErr == nil {
						s.Require().NoError(err)
						_, authErr = s.keeper.GetAuthorization(ctx, id, fromAddr, operatorAddr)
						s.Require().Error(authErr)
					} else {
						s.Require().Error(err)
					}
				})
			}
		}
	}
}
