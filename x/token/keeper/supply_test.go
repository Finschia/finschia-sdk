package keeper_test

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
)

func (s *KeeperTestSuite) TestIssue() {
	// create a not mintable class
	class := token.Token{
		Id:       "fee1dead",
		Name:     "NOT Mintable",
		Symbol:   "NO",
		Mintable: false,
	}
	err := s.keeper.Issue(s.ctx, class, s.vendor, s.vendor, sdk.OneInt())
	s.Require().NoError(err)

	mintActions := []string{
		token.ActionMint,
		token.ActionBurn,
	}
	for _, action := range mintActions {
		s.Require().False(s.keeper.GetGrant(s.ctx, s.vendor, class.Id, action))
	}
	s.Require().True(s.keeper.GetGrant(s.ctx, s.vendor, class.Id, token.ActionModify))

	// override fails
	class.Id = s.classID
	err = s.keeper.Issue(s.ctx, class, s.vendor, s.vendor, sdk.OneInt())
	s.Require().Error(err)
}

func (s *KeeperTestSuite) TestMint() {
	users := []sdk.AccAddress{s.vendor, s.operator, s.customer}
	to := s.vendor
	amount := token.FT{ClassId: s.classID, Amount: sdk.OneInt()}
	for _, grantee := range users {
		name := fmt.Sprintf("Grantee: %s", grantee)
		s.Run(name, func() {
			granted := s.keeper.GetGrant(s.ctx, grantee, amount.ClassId, token.ActionMint)
			err := s.keeper.Mint(s.ctx, grantee, to, []token.FT{amount})
			if granted {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestBurn() {
	users := []sdk.AccAddress{s.vendor, s.operator, s.customer}
	amount := token.FT{ClassId: s.classID, Amount: sdk.OneInt()}
	for _, from := range users {
		name := fmt.Sprintf("From: %s", from)
		s.Run(name, func() {
			granted := s.keeper.GetGrant(s.ctx, from, amount.ClassId, token.ActionBurn)
			err := s.keeper.Burn(s.ctx, from, []token.FT{amount})
			if granted {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestBurnFrom() {
	users := []sdk.AccAddress{s.vendor, s.operator, s.customer}
	amount := token.FT{ClassId: s.classID, Amount: sdk.OneInt()}
	for _, grantee := range users {
		for _, from := range users {
			name := fmt.Sprintf("Grantee: %s, From: %s", grantee, from)
			s.Run(name, func() {
				granted := s.keeper.GetGrant(s.ctx, grantee, amount.ClassId, token.ActionBurn)
				approved := s.keeper.GetApprove(s.ctx, from, grantee, amount.ClassId)
				err := s.keeper.BurnFrom(s.ctx, grantee, from, []token.FT{amount})
				if granted && approved {
					s.Require().NoError(err)
				} else {
					s.Require().Error(err)
				}
			})
		}
	}
}

func (s *KeeperTestSuite) TestModify() {
	users := []sdk.AccAddress{s.vendor, s.operator, s.customer}
	changes := []token.Pair{
		{Key: token.AttributeKeyName, Value: "new name"},
		{Key: token.AttributeKeyImageURI, Value: "new uri"},
		{Key: token.AttributeKeyMeta, Value: "new meta"},
	}
	for _, grantee := range users {
		name := fmt.Sprintf("Grantee: %s", grantee)
		s.Run(name, func() {
			granted := s.keeper.GetGrant(s.ctx, grantee, s.classID, token.ActionModify)
			err := s.keeper.Modify(s.ctx, s.classID, grantee, changes)
			if granted {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGrant() {
	users := []sdk.AccAddress{s.vendor, s.operator, s.customer}
	actions := []string{token.ActionMint, token.ActionBurn, token.ActionModify}
	for _, granter := range users {
		for _, grantee := range users {
			for _, action := range actions {
				name := fmt.Sprintf("Granter: %s, Grantee: %s", granter, grantee)
				s.Run(name, func() {
					granterGranted := s.keeper.GetGrant(s.ctx, granter, s.classID, action)
					granteeGranted := s.keeper.GetGrant(s.ctx, grantee, s.classID, action)
					err := s.keeper.Grant(s.ctx, granter, grantee, s.classID, action)
					if granterGranted && !granteeGranted {
						s.Require().NoError(err)
						s.Require().True(s.keeper.GetGrant(s.ctx, grantee, s.classID, action))
					} else {
						s.Require().Error(err)
					}
				})
			}
		}
	}
}

func (s *KeeperTestSuite) TestRevoke() {
	users := []sdk.AccAddress{s.vendor, s.operator, s.customer}
	actions := []string{token.ActionMint, token.ActionBurn, token.ActionModify}
	for _, grantee := range users {
		for _, action := range actions {
			name := fmt.Sprintf("Grantee: %s", grantee)
			s.Run(name, func() {
				granted := s.keeper.GetGrant(s.ctx, grantee, s.classID, action)
				err := s.keeper.Revoke(s.ctx, grantee, s.classID, action)
				if granted {
					s.Require().NoError(err)
					s.Require().False(s.keeper.GetGrant(s.ctx, grantee, s.classID, action))
				} else {
					s.Require().Error(err)
				}
			})
		}
	}
}
