package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
)

func (s *KeeperTestSuite) TestTransfer() {
	class := s.mintableClass
	amount := token.FT{ClassId: class, Amount: s.balance}
	err := s.keeper.Transfer(s.ctx, s.vendor, s.operator, []token.FT{amount})
	s.Require().NoError(err)
	err = s.keeper.Transfer(s.ctx, s.vendor, s.operator, []token.FT{amount})
	s.Require().Error(err)

	s.Require().Equal(
		token.FT{ClassId: class, Amount: s.balance.Add(s.balance)},
		s.keeper.GetBalance(s.ctx, s.operator, class))
	s.Require().Equal(
		token.FT{ClassId: class, Amount: sdk.ZeroInt()},
		s.keeper.GetBalance(s.ctx, s.vendor, class))
}

func (s *KeeperTestSuite) TestTransferFrom() {
	users := []sdk.AccAddress{s.vendor, s.operator, s.customer}
	to := s.vendor
	class := s.mintableClass
	amount := token.FT{ClassId: class, Amount: sdk.OneInt()}
	for _, proxy := range users {
		for _, from := range users {
			approved := s.keeper.GetApprove(s.ctx, from, proxy, class)
			err := s.keeper.TransferFrom(s.ctx, proxy, from, to, []token.FT{amount})
			if approved {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		}
	}
}

func (s *KeeperTestSuite) TestApprove() {
	users := []sdk.AccAddress{s.vendor, s.operator, s.customer}
	class := s.mintableClass
	for _, proxy := range users {
		for _, from := range users {
			approved := s.keeper.GetApprove(s.ctx, from, proxy, class)
			if !approved {
				err := s.keeper.Approve(s.ctx, from, proxy, class)
				s.Require().NoError(err)
				approved = s.keeper.GetApprove(s.ctx, from, proxy, class)
				s.Require().True(approved)
			}
		}
	}
}
