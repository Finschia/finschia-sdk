package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
)

func (s *KeeperTestSuite) TestIssue() {
	ctx, _ := s.ctx.CacheContext()

	// create a not mintable class
	class := token.Contract{
		Name:     "NOT Mintable",
		Symbol:   "NO",
		Mintable: false,
	}
	contractID := s.keeper.Issue(ctx, class, s.vendor, s.vendor, sdk.OneInt())

	mintPermissions := []token.Permission{
		token.PermissionMint,
		token.PermissionBurn,
	}
	for _, permission := range mintPermissions {
		s.Require().Nil(s.keeper.GetGrant(ctx, contractID, s.vendor, permission))
	}
	s.Require().NotNil(s.keeper.GetGrant(ctx, contractID, s.vendor, token.PermissionModify))
}

func (s *KeeperTestSuite) TestMint() {
	testCases := map[string]struct {
		grantee sdk.AccAddress
		err     error
	}{
		"valid request": {
			grantee: s.operator,
		},
		"no permission": {
			grantee: s.customer,
			err:     token.ErrTokenNoPermission,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			err := s.keeper.Mint(ctx, s.contractID, tc.grantee, s.stranger, sdk.OneInt())
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}
		})
	}
}

func (s *KeeperTestSuite) TestBurn() {
	testCases := map[string]struct {
		from   sdk.AccAddress
		amount sdk.Int
		err    error
	}{
		"valid request": {
			from:   s.vendor,
			amount: s.balance,
		},
		"no permission": {
			from:   s.customer,
			amount: s.balance,
			err:    token.ErrTokenNoPermission,
		},
		"insufficient tokens": {
			from:   s.vendor,
			amount: s.balance.Add(sdk.OneInt()),
			err:    token.ErrInsufficientBalance,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			err := s.keeper.Burn(ctx, s.contractID, tc.from, tc.amount)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}
		})
	}
}

func (s *KeeperTestSuite) TestOperatorBurn() {
	testCases := map[string]struct {
		operator sdk.AccAddress
		from     sdk.AccAddress
		amount   sdk.Int
		err      error
	}{
		"valid request": {
			operator: s.operator,
			from:     s.customer,
			amount:   s.balance,
		},
		"not authorized": {
			operator: s.vendor,
			from:     s.stranger,
			amount:   s.balance,
			err:      token.ErrTokenNotApproved,
		},
		"no permission": {
			operator: s.stranger,
			from:     s.customer,
			amount:   s.balance,
			err:      token.ErrTokenNoPermission,
		},
		"insufficient tokens": {
			operator: s.operator,
			from:     s.customer,
			amount:   s.balance.Add(sdk.OneInt()),
			err:      token.ErrInsufficientBalance,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			err := s.keeper.OperatorBurn(ctx, s.contractID, tc.operator, tc.from, tc.amount)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}
		})
	}
}

func (s *KeeperTestSuite) TestModify() {
	changes := []token.Attribute{
		{Key: token.AttributeKeyName.String(), Value: "new name"},
		{Key: token.AttributeKeyImageURI.String(), Value: "new uri"},
		{Key: token.AttributeKeyMeta.String(), Value: "new meta"},
	}

	ctx, _ := s.ctx.CacheContext()

	err := s.keeper.Modify(ctx, s.contractID, s.vendor, changes)
	s.Require().NoError(err)

	class, err := s.keeper.GetClass(ctx, s.contractID)
	s.Require().NoError(err)

	s.Require().Equal(changes[0].Value, class.Name)
	s.Require().Equal(changes[1].Value, class.Uri)
	s.Require().Equal(changes[2].Value, class.Meta)
}
