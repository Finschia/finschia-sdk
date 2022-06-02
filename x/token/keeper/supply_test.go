package keeper_test

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
)

func (s *KeeperTestSuite) TestIssue() {
	ctx, _ := s.ctx.CacheContext()

	// create a not mintable class
	class := token.TokenClass{
		ContractId:       "fee1dead",
		Name:     "NOT Mintable",
		Symbol:   "NO",
		Mintable: false,
	}
	err := s.keeper.Issue(ctx, class, s.vendor, s.vendor, sdk.OneInt())
	s.Require().NoError(err)

	mintPermissions := []token.Permission{
		token.Permission_Mint,
		token.Permission_Burn,
	}
	for _, permission := range mintPermissions {
		s.Require().Nil(s.keeper.GetGrant(ctx, class.ContractId, s.vendor, permission))
	}
	s.Require().NotNil(s.keeper.GetGrant(ctx, class.ContractId, s.vendor, token.Permission_Modify))

	// override fails
	class.ContractId = s.contractID
	err = s.keeper.Issue(ctx, class, s.vendor, s.vendor, sdk.OneInt())
	s.Require().Error(err)
}

func (s *KeeperTestSuite) TestMint() {
	userDescriptions := map[sdk.AccAddress]string{
		s.vendor: "vendor",
		s.operator: "operator",
		s.customer: "customer",
	}
	to := s.vendor
	amount := sdk.OneInt()
	for grantee, desc := range userDescriptions {
		name := fmt.Sprintf("Grantee: %s", desc)
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			grant := s.keeper.GetGrant(ctx, s.contractID, grantee, token.Permission_Mint)
			err := s.keeper.Mint(ctx, s.contractID, grantee, to, amount)
			if grant != nil {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestBurn() {
	userDescriptions := map[sdk.AccAddress]string{
		s.vendor: "vendor",
		s.operator: "operator",
		s.customer: "customer",
	}
	amountDescriptions := map[sdk.Int]string {
		s.balance: "limit",
		s.balance.Add(sdk.OneInt()): "excess",
	}
	for from, fromDesc := range userDescriptions {
		for amount, amountDesc := range amountDescriptions {
			name := fmt.Sprintf("From: %s, Amount: %s", fromDesc, amountDesc)
			s.Run(name, func() {
				ctx, _ := s.ctx.CacheContext()

				grant := s.keeper.GetGrant(ctx, s.contractID, from, token.Permission_Burn)
				err := s.keeper.Burn(ctx, s.contractID, from, amount)
				if grant != nil && amount.LTE(s.balance) {
					s.Require().NoError(err)
				} else {
					s.Require().Error(err)
				}
			})
		}
	}
}

func (s *KeeperTestSuite) TestOperatorBurn() {
	userDescriptions := map[sdk.AccAddress]string{
		s.vendor: "vendor",
		s.operator: "operator",
		s.customer: "customer",
	}
	amountDescriptions := map[sdk.Int]string {
		s.balance: "limit",
		s.balance.Add(sdk.OneInt()): "excess",
	}
	for operator, operatorDesc := range userDescriptions {
		for from, fromDesc := range userDescriptions {
			for amount, amountDesc := range amountDescriptions {
				name := fmt.Sprintf("Operator: %s, From: %s, Amount: %s", operatorDesc, fromDesc, amountDesc)
				s.Run(name, func() {
					ctx, _ := s.ctx.CacheContext()

					grant := s.keeper.GetGrant(ctx, s.contractID, operator, token.Permission_Burn)
					authorization := s.keeper.GetAuthorization(ctx, s.contractID, from, operator)
					err := s.keeper.OperatorBurn(ctx, s.contractID, operator, from, amount)
					if grant != nil && authorization != nil && amount.LTE(s.balance) {
						s.Require().NoError(err)
					} else {
						s.Require().Error(err)
					}
				})
			}
		}
	}
}

func (s *KeeperTestSuite) TestModify() {
	contractDescriptions := map[string]string{
		s.contractID: "valid",
		"fee1dead": "not-exist",
	}
	userDescriptions := map[sdk.AccAddress]string{
		s.vendor: "vendor",
		s.operator: "operator",
		s.customer: "customer",
	}
	changes := []token.Pair{
		{Field: token.AttributeKey_Name.String(), Value: "new name"},
		{Field: token.AttributeKey_ImageURI.String(), Value: "new uri"},
		{Field: token.AttributeKey_Meta.String(), Value: "new meta"},
	}

	for contractID, contractDesc := range contractDescriptions {
		for grantee, granteeDesc := range userDescriptions {
			name := fmt.Sprintf("Grantee: %s, Contract: %s", granteeDesc, contractDesc)
			s.Run(name, func() {
				ctx, _ := s.ctx.CacheContext()

				err := s.keeper.Modify(ctx, contractID, grantee, changes)
				if contractID == s.contractID {
					s.Require().NoError(err)
				} else {
					s.Require().Error(err)
				}
			})
		}
	}
}
