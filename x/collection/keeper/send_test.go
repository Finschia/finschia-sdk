package keeper_test

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/collection"
)

func (s *KeeperTestSuite) TestSendCoins() {
	testCases := map[string]struct {
		amount collection.Coin
		err    error
	}{
		"valid send (fungible token)": {
			amount: collection.NewFTCoin(s.ftClassID, s.balance),
		},
		"valid send (non-fungible token)": {
			amount: collection.NewNFTCoin(s.nftClassID, 1),
		},
		"insufficient tokens": {
			amount: collection.NewFTCoin(s.ftClassID, s.balance.Add(sdk.OneInt())),
			err:    collection.ErrInsufficientTokens,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			tokenID := tc.amount.TokenId
			customerBalance := s.keeper.GetBalance(ctx, s.contractID, s.customer, tokenID)
			operatorBalance := s.keeper.GetBalance(ctx, s.contractID, s.operator, tokenID)

			err := s.keeper.SendCoins(ctx, s.contractID, s.customer, s.operator, collection.NewCoins(tc.amount))
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

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

	userDescriptions := map[string]string{
		s.vendor.String():   "vendor",
		s.operator.String(): "operator",
		s.customer.String(): "customer",
	}
	for operator, operatorDesc := range userDescriptions {
		for from, fromDesc := range userDescriptions {
			name := fmt.Sprintf("Operator: %s, From: %s", operatorDesc, fromDesc)
			s.Run(name, func() {
				ctx, _ := s.ctx.CacheContext()

				fromAddr, err := sdk.AccAddressFromBech32(from)
				s.Require().NoError(err)
				operatorAddr, err := sdk.AccAddressFromBech32(operator)
				s.Require().NoError(err)

				_, queryErr := s.keeper.GetAuthorization(ctx, s.contractID, fromAddr, operatorAddr)
				err = s.keeper.AuthorizeOperator(ctx, s.contractID, fromAddr, operatorAddr)
				if queryErr == nil { // authorize must fail
					s.Require().ErrorIs(err, collection.ErrAuthorizationAlreadyExists)
				} else {
					s.Require().ErrorIs(queryErr, collection.ErrAuthorizationNotFound)
					s.Require().NoError(err)
					_, queryErr := s.keeper.GetAuthorization(ctx, s.contractID, fromAddr, operatorAddr)
					s.Require().NoError(queryErr)
				}
			})
		}
	}
}

func (s *KeeperTestSuite) TestRevokeOperator() {
	// make sure the dummy contract does not exist
	dummyContractID := "deadbeef"
	_, err := s.keeper.GetContract(s.ctx, dummyContractID)
	s.Require().Error(err)

	userDescriptions := map[string]string{
		s.vendor.String():   "vendor",
		s.operator.String(): "operator",
		s.customer.String(): "customer",
	}
	for operator, operatorDesc := range userDescriptions {
		for from, fromDesc := range userDescriptions {
			name := fmt.Sprintf("Operator: %s, From: %s", operatorDesc, fromDesc)
			s.Run(name, func() {
				ctx, _ := s.ctx.CacheContext()

				fromAddr, err := sdk.AccAddressFromBech32(from)
				s.Require().NoError(err)
				operatorAddr, err := sdk.AccAddressFromBech32(operator)
				s.Require().NoError(err)

				_, queryErr := s.keeper.GetAuthorization(ctx, s.contractID, fromAddr, operatorAddr)
				err = s.keeper.RevokeOperator(ctx, s.contractID, fromAddr, operatorAddr)
				if queryErr != nil { // revoke must fail
					s.Require().ErrorIs(queryErr, collection.ErrAuthorizationNotFound)
					s.Require().ErrorIs(err, collection.ErrAuthorizationNotFound)
				} else {
					s.Require().NoError(err)
					_, queryErr := s.keeper.GetAuthorization(ctx, s.contractID, fromAddr, operatorAddr)
					s.Require().ErrorIs(queryErr, collection.ErrAuthorizationNotFound)
				}
			})
		}
	}
}
