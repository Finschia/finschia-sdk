package keeper_test

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
)

func (s *KeeperTestSuite) TestTransfer() {
	amount := token.FT{ClassId: s.classID, Amount: s.balance}

	testCases := map[string]struct{
		amount token.FT
		valid bool
	}{
		"valid transfer": {
			amount,
			true,
		},
		"insufficient tokens": {
			token.FT{
				ClassId: amount.ClassId,
				Amount: amount.Amount.Mul(sdk.NewInt(2)),
			},
			false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			err := s.keeper.Transfer(s.ctx, s.vendor, s.operator, []token.FT{tc.amount})
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			vendorBalance := s.keeper.GetBalance(s.ctx, s.vendor, tc.amount.ClassId).Amount
			s.Require().True(s.balance.Sub(tc.amount.Amount).Equal(vendorBalance))

			operatorBalance := s.keeper.GetBalance(s.ctx, s.operator, tc.amount.ClassId).Amount
			s.Require().True(s.balance.Add(tc.amount.Amount).Equal(operatorBalance))
		})
	}
}

func (s *KeeperTestSuite) TestTransferFrom() {
	users := []sdk.AccAddress{s.vendor, s.operator, s.customer}
	to := s.vendor
	amount := token.FT{ClassId: s.classID, Amount: s.balance}
	for _, proxy := range users {
		for _, from := range users {
			name := fmt.Sprintf("Proxy: %s, From: %s", proxy, from)
			s.Run(name, func() {
				approved := s.keeper.GetApprove(s.ctx, from, proxy, amount.ClassId)
				err := s.keeper.TransferFrom(s.ctx, proxy, from, to, []token.FT{amount})
				if approved {
					s.Require().NoError(err)
				} else {
					s.Require().Error(err)
				}
			})
		}
	}
}

func (s *KeeperTestSuite) TestApprove() {
	users := []sdk.AccAddress{s.vendor, s.operator, s.customer}
	for _, proxy := range users {
		for _, from := range users {
			name := fmt.Sprintf("Proxy: %s, From: %s", proxy, from)
			s.Run(name, func() {
				approved := s.keeper.GetApprove(s.ctx, from, proxy, s.classID)
				err := s.keeper.Approve(s.ctx, from, proxy, s.classID)
				if !approved {
					s.Require().NoError(err)
					approved = s.keeper.GetApprove(s.ctx, from, proxy, s.classID)
					s.Require().True(approved)
				} else {
					s.Require().Error(err)
				}
			})
		}
	}

	// no such a class
	err := s.keeper.Approve(s.ctx, s.vendor, s.customer, "00000000")
	s.Require().Error(err)
}
