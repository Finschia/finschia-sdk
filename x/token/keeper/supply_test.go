package keeper_test

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
)

func (s *KeeperTestSuite) TestIssue() {
	ctx, _ := s.ctx.CacheContext()

	// create a not mintable class
	class := token.Contract{
		Id:       "fee1dead",
		Name:     "NOT Mintable",
		Symbol:   "NO",
		Mintable: false,
	}
	s.keeper.Issue(ctx, class, s.vendor, s.vendor, sdk.OneInt())

	mintPermissions := []token.Permission{
		token.PermissionMint,
		token.PermissionBurn,
	}
	for _, permission := range mintPermissions {
		s.Require().Nil(s.keeper.GetGrant(ctx, class.Id, s.vendor, permission))
	}
	s.Require().NotNil(s.keeper.GetGrant(ctx, class.Id, s.vendor, token.PermissionModify))

	// override fails
	class.Id = s.contractID
	s.Require().Panics(func() {
		s.keeper.Issue(ctx, class, s.vendor, s.vendor, sdk.OneInt())
	})
}

func (s *KeeperTestSuite) TestMint() {
	userDescriptions := map[string]string{
		s.vendor.String():   "vendor",
		s.operator.String(): "operator",
		s.customer.String(): "customer",
	}
	to := s.vendor
	amount := sdk.OneInt()
	for granteeStr, desc := range userDescriptions {
		grantee, err := sdk.AccAddressFromBech32(granteeStr)
		s.Require().NoError(err)
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
	userDescriptions := map[string]string{
		s.vendor.String():   "vendor",
		s.operator.String(): "operator",
		s.customer.String(): "customer",
	}
	amountDescriptions := map[sdk.Int]string{
		s.balance:                   "limit",
		s.balance.Add(sdk.OneInt()): "excess",
	}
	for fromStr, fromDesc := range userDescriptions {
		from, err := sdk.AccAddressFromBech32(fromStr)
		s.Require().NoError(err)
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
	userDescriptions := map[string]string{
		s.vendor.String():   "vendor",
		s.operator.String(): "operator",
		s.customer.String(): "customer",
	}
	amountDescriptions := map[sdk.Int]string{
		s.balance:                   "limit",
		s.balance.Add(sdk.OneInt()): "excess",
	}
	for operatorStr, operatorDesc := range userDescriptions {
		operator, err := sdk.AccAddressFromBech32(operatorStr)
		s.Require().NoError(err)
		for fromStr, fromDesc := range userDescriptions {
			from, err := sdk.AccAddressFromBech32(fromStr)
			s.Require().NoError(err)
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
	userDescriptions := map[string]string{
		s.vendor.String():   "vendor",
		s.operator.String(): "operator",
		s.customer.String(): "customer",
	}
	changes := []token.Attribute{
		{Key: token.AttributeKeyName.String(), Value: "new name"},
		{Key: token.AttributeKeyImageURI.String(), Value: "new uri"},
		{Key: token.AttributeKeyMeta.String(), Value: "new meta"},
	}

	for contractID, contractDesc := range contractDescriptions {
		for granteeStr, granteeDesc := range userDescriptions {
			grantee, err := sdk.AccAddressFromBech32(granteeStr)
			s.Require().NoError(err)
			name := fmt.Sprintf("Grantee: %s, Contract: %s", granteeDesc, contractDesc)
			s.Run(name, func() {
				ctx, _ := s.ctx.CacheContext()

				err := s.keeper.Modify(ctx, contractID, grantee, changes)
				if contractID != s.contractID {
					s.Require().Error(err)
					return
				}
				s.Require().NoError(err)

				class, err := s.keeper.GetClass(ctx, contractID)
				s.Require().NoError(err)

				s.Require().Equal(changes[0].Value, class.Name)
				s.Require().Equal(changes[1].Value, class.Uri)
				s.Require().Equal(changes[2].Value, class.Meta)
			})
		}
	}
}
