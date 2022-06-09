package keeper_test

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/collection"
)

func (s *KeeperTestSuite) TestMsgSend() {
	testCases := map[string]struct {
		contractID string
		amount  sdk.Int
		valid   bool
	}{
		"valid request": {
			contractID: s.contractID,
			amount:  s.balance,
			valid:   true,
		},
		"insufficient funds": {
			contractID: "fee1dead",
			amount:  s.balance,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgSend{
				ContractId: tc.contractID,
				From:    s.vendor.String(),
				To:      s.customer.String(),
				Amount:  []collection.Coin{{
					TokenId: s.ftClassID + fmt.Sprintf("%08x", 0),
					Amount: tc.amount,
				}},
			}
			res, err := s.msgServer.Send(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgOperatorSend() {
	testCases := map[string]struct {
		operator sdk.AccAddress
		from  sdk.AccAddress
		amount sdk.Int
		valid bool
	}{
		"valid request": {
			operator: s.operator,
			from:  s.customer,
			amount: s.balance,
			valid: true,
		},
		"not approved": {
			operator: s.vendor,
			from:  s.customer,
			amount: s.balance,
		},
		"insufficient funds": {
			operator: s.operator,
			from:  s.customer,
			amount:  s.balance.Add(sdk.OneInt()),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgOperatorSend{
				ContractId: s.contractID,
				Operator:   tc.operator.String(),
				From:    tc.from.String(),
				To:      s.vendor.String(),
				Amount:  []collection.Coin{{
					TokenId: s.ftClassID + fmt.Sprintf("%08x", 0),
					Amount: tc.amount,
				}},
			}
			res, err := s.msgServer.OperatorSend(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgTransferFT() {
	testCases := map[string]struct {
		contractID string
		amount  sdk.Int
		valid   bool
	}{
		"valid request": {
			contractID: s.contractID,
			amount:  s.balance,
			valid:   true,
		},
		"insufficient funds": {
			contractID: "fee1dead",
			amount:  s.balance,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgTransferFT{
				ContractId: tc.contractID,
				From:    s.vendor.String(),
				To:      s.customer.String(),
				Amount:  []collection.Coin{{
					TokenId: s.ftClassID + fmt.Sprintf("%08x", 0),
					Amount: tc.amount,
				}},
			}
			res, err := s.msgServer.TransferFT(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgTransferFTFrom() {
	testCases := map[string]struct {
		proxy sdk.AccAddress
		from  sdk.AccAddress
		amount sdk.Int
		valid bool
	}{
		"valid request": {
			proxy: s.operator,
			from:  s.customer,
			amount: s.balance,
			valid: true,
		},
		"not approved": {
			proxy: s.vendor,
			from:  s.customer,
			amount: s.balance,
		},
		"insufficient funds": {
			proxy: s.operator,
			from:  s.customer,
			amount:  s.balance.Add(sdk.OneInt()),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgTransferFTFrom{
				ContractId: s.contractID,
				Proxy:   tc.proxy.String(),
				From:    tc.from.String(),
				To:      s.vendor.String(),
				Amount:  []collection.Coin{{
					TokenId: s.ftClassID + fmt.Sprintf("%08x", 0),
					Amount: tc.amount,
				}},
			}
			res, err := s.msgServer.TransferFTFrom(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgTransferNFT() {
	testCases := map[string]struct {
		tokenID string
		valid   bool
	}{
		"valid request": {
			tokenID: s.nftClassID + fmt.Sprintf("%08x", 1),
			valid:   true,
		},
		"insufficient funds": {
			tokenID: s.nftClassID + fmt.Sprintf("%08x", 2),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgTransferNFT{
				ContractId: s.contractID,
				From:    s.customer.String(),
				To:      s.vendor.String(),
				TokenIds:  []string{tc.tokenID},
			}
			res, err := s.msgServer.TransferNFT(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgTransferNFTFrom() {
	testCases := map[string]struct {
		proxy sdk.AccAddress
		from  sdk.AccAddress
		tokenID string
		valid bool
	}{
		"valid request": {
			proxy: s.operator,
			from:  s.customer,
			tokenID: s.nftClassID + fmt.Sprintf("%08x", 1),
			valid: true,
		},
		"not approved": {
			proxy: s.vendor,
			from:  s.customer,
			tokenID: s.nftClassID + fmt.Sprintf("%08x", 1),
		},
		"insufficient funds": {
			proxy: s.operator,
			from:  s.customer,
			tokenID: s.nftClassID + fmt.Sprintf("%08x", 2),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgTransferNFTFrom{
				ContractId: s.contractID,
				Proxy:   tc.proxy.String(),
				From:    tc.from.String(),
				To:      s.vendor.String(),
				TokenIds:  []string{tc.tokenID},
			}
			res, err := s.msgServer.TransferNFTFrom(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgAuthorizeOperator() {
	testCases := map[string]struct {
		holder sdk.AccAddress
		operator    sdk.AccAddress
		valid    bool
	}{
		"valid request": {
			holder: s.customer,
			operator:    s.vendor,
			valid:    true,
		},
		"already approved": {
			holder: s.customer,
			operator:    s.operator,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgAuthorizeOperator{
				ContractId:  s.contractID,
				Holder: tc.holder.String(),
				Operator:    tc.operator.String(),
			}
			res, err := s.msgServer.AuthorizeOperator(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgRevokeOperator() {
	testCases := map[string]struct {
		holder sdk.AccAddress
		operator    sdk.AccAddress
		valid    bool
	}{
		"valid request": {
			holder:   s.customer,
			operator: s.operator,
			valid:    true,
		},
		"no authorization": {
			holder:   s.customer,
			operator: s.vendor,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgRevokeOperator{
				ContractId: s.contractID,
				Holder:     tc.holder.String(),
				Operator:   tc.operator.String(),
			}
			res, err := s.msgServer.RevokeOperator(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgApprove() {
	testCases := map[string]struct {
		approver sdk.AccAddress
		proxy    sdk.AccAddress
		valid    bool
	}{
		"valid request": {
			approver: s.customer,
			proxy:    s.vendor,
			valid:    true,
		},
		"already approved": {
			approver: s.customer,
			proxy:    s.operator,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgApprove{
				ContractId:  s.contractID,
				Approver: tc.approver.String(),
				Proxy:    tc.proxy.String(),
			}
			res, err := s.msgServer.Approve(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgDisapprove() {
	testCases := map[string]struct {
		approver sdk.AccAddress
		proxy    sdk.AccAddress
		valid    bool
	}{
		"valid request": {
			approver: s.customer,
			proxy:    s.operator,
			valid:    true,
		},
		"no authorization": {
			approver: s.customer,
			proxy:    s.vendor,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgDisapprove{
				ContractId:  s.contractID,
				Approver: tc.approver.String(),
				Proxy:    tc.proxy.String(),
			}
			res, err := s.msgServer.Disapprove(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}
