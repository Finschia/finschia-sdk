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
		ContractId: "fee1dead",
		Name:       "NOT Mintable",
		Symbol:     "NO",
		Mintable:   false,
	}
	s.keeper.Issue(ctx, class, s.vendor, s.vendor, sdk.OneInt())

	mintPermissions := []token.Permission{
		token.PermissionMint,
		token.PermissionBurn,
	}
	for _, permission := range mintPermissions {
		s.Require().Nil(s.keeper.GetGrant(ctx, class.ContractId, s.vendor, permission))
	}
	s.Require().NotNil(s.keeper.GetGrant(ctx, class.ContractId, s.vendor, token.PermissionModify))

	// override fails
	class.ContractId = s.contractID
	s.Require().Panics(func() {
		s.keeper.Issue(ctx, class, s.vendor, s.vendor, sdk.OneInt())
	})
}

func (s *KeeperTestSuite) TestMint() {
	userDescriptions := map[sdk.AccAddress]string{
		s.vendor:   "vendor",
		s.operator: "operator",
		s.customer: "customer",
	}
	to := s.vendor
	amount := sdk.OneInt()
	for grantee, desc := range userDescriptions {
		name := fmt.Sprintf("Grantee: %s", desc)
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			_, grantErr := s.keeper.GetGrant(ctx, s.contractID, grantee, token.PermissionMint)
			err := s.keeper.Mint(ctx, s.contractID, grantee, to, amount)
			if grantErr == nil {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestBurn() {
	userDescriptions := map[sdk.AccAddress]string{
		s.vendor:   "vendor",
		s.operator: "operator",
		s.customer: "customer",
	}
	amountDescriptions := map[sdk.Int]string{
		s.balance:                   "limit",
		s.balance.Add(sdk.OneInt()): "excess",
	}
	for from, fromDesc := range userDescriptions {
		for amount, amountDesc := range amountDescriptions {
			name := fmt.Sprintf("From: %s, Amount: %s", fromDesc, amountDesc)
			s.Run(name, func() {
				ctx, _ := s.ctx.CacheContext()

				_, grantErr := s.keeper.GetGrant(ctx, s.contractID, from, token.PermissionBurn)
				err := s.keeper.Burn(ctx, s.contractID, from, amount)
				if grantErr == nil && amount.LTE(s.balance) {
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
		s.vendor:   "vendor",
		s.operator: "operator",
		s.customer: "customer",
	}
	amountDescriptions := map[sdk.Int]string{
		s.balance:                   "limit",
		s.balance.Add(sdk.OneInt()): "excess",
	}
	for operator, operatorDesc := range userDescriptions {
		for from, fromDesc := range userDescriptions {
			for amount, amountDesc := range amountDescriptions {
				name := fmt.Sprintf("Operator: %s, From: %s, Amount: %s", operatorDesc, fromDesc, amountDesc)
				s.Run(name, func() {
					ctx, _ := s.ctx.CacheContext()

					_, grantErr := s.keeper.GetGrant(ctx, s.contractID, operator, token.PermissionBurn)
					_, authErr := s.keeper.GetAuthorization(ctx, s.contractID, from, operator)
					err := s.keeper.OperatorBurn(ctx, s.contractID, operator, from, amount)
					if grantErr == nil && authErr == nil && amount.LTE(s.balance) {
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
		"fee1dead":   "not-exist",
	}
	userDescriptions := map[sdk.AccAddress]string{
		s.vendor:   "vendor",
		s.operator: "operator",
		s.customer: "customer",
	}
	changes := []token.Pair{
		{Field: token.AttributeKeyName.String(), Value: "new name"},
		{Field: token.AttributeKeyImageURI.String(), Value: "new uri"},
		{Field: token.AttributeKeyMeta.String(), Value: "new meta"},
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
