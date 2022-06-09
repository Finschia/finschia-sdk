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
			// TODO: feature not supported
		},
		// "not approved": {
		// 	operator: s.vendor,
		// 	from:  s.customer,
		// 	amount: s.balance,
		// },
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
			// TODO: feature not supported
		},
		// "not approved": {
		// 	proxy: s.vendor,
		// 	from:  s.customer,
		// 	amount: s.balance,
		// },
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
		contractID string
		valid   bool
	}{
		"valid request": {
			contractID: s.contractID,
			valid:   true,
		},
		"insufficient funds": {
			contractID: "fee1dead",
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgTransferNFT{
				ContractId: tc.contractID,
				From:    s.customer.String(),
				To:      s.vendor.String(),
				TokenIds:  []string{
					s.nftClassID + fmt.Sprintf("%08x", 1),
				},
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
		contractID string
		proxy sdk.AccAddress
		from  sdk.AccAddress
		valid bool
	}{
		"valid request": {
			contractID: s.contractID,
			proxy: s.operator,
			from:  s.customer,
			// TODO: feature not supported
		},
		// "not approved": {
		// 	proxy: s.vendor,
		// 	from:  s.customer,
		// },
		"insufficient funds": {
			contractID: "fee1dead",
			proxy: s.operator,
			from:  s.customer,
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
				TokenIds:  []string{
					s.nftClassID + fmt.Sprintf("%08x", 1),
				},
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
