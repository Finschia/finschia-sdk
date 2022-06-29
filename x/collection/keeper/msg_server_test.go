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
			tokenID: "deadbeef" + fmt.Sprintf("%08x", 1),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgTransferNFT{
				ContractId: s.contractID,
				From:    s.vendor.String(),
				To:      s.customer.String(),
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
	tokenID := s.nftClassID + fmt.Sprintf("%08x", 3)
	testCases := map[string]struct {
		proxy sdk.AccAddress
		from  sdk.AccAddress
		tokenID string
		valid bool
	}{
		"valid request": {
			proxy: s.operator,
			from:  s.customer,
			tokenID: tokenID,
			valid: true,
		},
		"not approved": {
			proxy: s.vendor,
			from:  s.customer,
			tokenID: tokenID,
		},
		"insufficient funds": {
			proxy: s.operator,
			from:  s.customer,
			tokenID: "deadbeef" + fmt.Sprintf("%08x", 1),
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

func (s *KeeperTestSuite) TestMsgCreateContract() {
	testCases := map[string]struct {
		owner sdk.AccAddress
		valid    bool
	}{
		"valid request": {
			owner: s.vendor,
			valid:    true,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgCreateContract{
				Owner: tc.owner.String(),
			}
			res, err := s.msgServer.CreateContract(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgIssueFT() {
	testCases := map[string]struct {
		contractID string
		owner sdk.AccAddress
		valid  bool
	}{
		"valid request": {
			contractID: s.contractID,
			owner: s.vendor,
			valid:  true,
		},
		"no permission": {
			contractID: s.contractID,
			owner: s.customer,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgIssueFT{
				ContractId: tc.contractID,
				Owner:  tc.owner.String(),
			}
			res, err := s.msgServer.IssueFT(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgIssueNFT() {
	testCases := map[string]struct {
		contractID string
		owner sdk.AccAddress
		valid  bool
	}{
		"valid request": {
			contractID: s.contractID,
			owner: s.vendor,
			valid:  true,
		},
		"no permission": {
			contractID: s.contractID,
			owner: s.customer,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgIssueNFT{
				ContractId: tc.contractID,
				Owner:  tc.owner.String(),
			}
			res, err := s.msgServer.IssueNFT(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgMintFT() {
	amount := collection.NewCoin(s.ftClassID + fmt.Sprintf("%08x", 0), sdk.OneInt())
	testCases := map[string]struct {
		contractID string
		from sdk.AccAddress
		amount []collection.Coin
		valid  bool
	}{
		"valid request": {
			contractID: s.contractID,
			from: s.vendor,
			amount: []collection.Coin{amount},
			valid:  true,
		},
		"no permission": {
			contractID: s.contractID,
			from: s.customer,
			amount: []collection.Coin{amount},
		},
		"no class of the token": {
			contractID: s.contractID,
			from: s.vendor,
			amount: []collection.Coin{collection.NewCoin("deadbeef" + fmt.Sprintf("%08x", 0), sdk.OneInt())},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgMintFT{
				ContractId: tc.contractID,
				From:  tc.from.String(),
				To: s.customer.String(),
				Amount: tc.amount,
			}
			res, err := s.msgServer.MintFT(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgMintNFT() {
	param := collection.MintNFTParam{
		TokenType: s.nftClassID,
	}
	testCases := map[string]struct {
		contractID string
		from sdk.AccAddress
		params []collection.MintNFTParam
		valid  bool
	}{
		"valid request": {
			contractID: s.contractID,
			from: s.vendor,
			params: []collection.MintNFTParam{param},
			valid:  true,
		},
		"no permission": {
			contractID: s.contractID,
			from: s.customer,
			params: []collection.MintNFTParam{param},
		},
		"no class of the token": {
			contractID: s.contractID,
			from: s.vendor,
			params: []collection.MintNFTParam{{
				TokenType: "deadbeef",
			}},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgMintNFT{
				ContractId: tc.contractID,
				From:  tc.from.String(),
				To: s.customer.String(),
				Params: tc.params,
			}
			res, err := s.msgServer.MintNFT(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgBurn() {
	amount := collection.NewCoin(s.ftClassID + fmt.Sprintf("%08x", 0), s.balance)
	testCases := map[string]struct {
		contractID string
		from sdk.AccAddress
		amount []collection.Coin
		valid  bool
	}{
		"valid request": {
			contractID: s.contractID,
			from: s.vendor,
			amount: []collection.Coin{amount},
			valid:  true,
		},
		"no permission": {
			contractID: s.contractID,
			from: s.customer,
			amount: []collection.Coin{amount},
		},
		"insufficient funds": {
			contractID: s.contractID,
			from: s.vendor,
			amount: []collection.Coin{collection.NewCoin("deadbeef" + fmt.Sprintf("%08x", 0), sdk.OneInt())},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgBurn{
				ContractId: tc.contractID,
				From:  tc.from.String(),
				Amount: tc.amount,
			}
			res, err := s.msgServer.Burn(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgOperatorBurn() {
	amount := collection.NewCoin(s.ftClassID + fmt.Sprintf("%08x", 0), s.balance)
	testCases := map[string]struct {
		contractID string
		operator sdk.AccAddress
		from sdk.AccAddress
		amount []collection.Coin
		valid  bool
	}{
		"valid request": {
			contractID: s.contractID,
			operator: s.operator,
			from: s.customer,
			amount: []collection.Coin{amount},
			valid:  true,
		},
		"no authorization": {
			contractID: s.contractID,
			operator: s.vendor,
			from: s.customer,
			amount: []collection.Coin{amount},
		},
		"no permission": {
			contractID: s.contractID,
			operator: s.stranger,
			from: s.customer,
			amount: []collection.Coin{amount},
		},
		"insufficient funds": {
			contractID: s.contractID,
			operator: s.operator,
			from: s.customer,
			amount: []collection.Coin{collection.NewCoin("deadbeef" + fmt.Sprintf("%08x", 0), sdk.OneInt())},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgOperatorBurn{
				ContractId: tc.contractID,
				Operator: tc.operator.String(),
				From:  tc.from.String(),
				Amount: tc.amount,
			}
			res, err := s.msgServer.OperatorBurn(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgBurnFT() {
	amount := collection.NewCoin(s.ftClassID + fmt.Sprintf("%08x", 0), s.balance)
	testCases := map[string]struct {
		contractID string
		from sdk.AccAddress
		amount []collection.Coin
		valid  bool
	}{
		"valid request": {
			contractID: s.contractID,
			from: s.vendor,
			amount: []collection.Coin{amount},
			valid:  true,
		},
		"no permission": {
			contractID: s.contractID,
			from: s.customer,
			amount: []collection.Coin{amount},
		},
		"insufficient funds": {
			contractID: s.contractID,
			from: s.vendor,
			amount: []collection.Coin{collection.NewCoin("deadbeef" + fmt.Sprintf("%08x", 0), sdk.OneInt())},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgBurnFT{
				ContractId: tc.contractID,
				From:  tc.from.String(),
				Amount: tc.amount,
			}
			res, err := s.msgServer.BurnFT(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgBurnFTFrom() {
	amount := collection.NewCoin(s.ftClassID + fmt.Sprintf("%08x", 0), s.balance)
	testCases := map[string]struct {
		contractID string
		proxy sdk.AccAddress
		from sdk.AccAddress
		amount []collection.Coin
		valid  bool
	}{
		"valid request": {
			contractID: s.contractID,
			proxy: s.operator,
			from: s.customer,
			amount: []collection.Coin{amount},
			valid:  true,
		},
		"no authorization": {
			contractID: s.contractID,
			proxy: s.vendor,
			from: s.customer,
			amount: []collection.Coin{amount},
		},
		"no permission": {
			contractID: s.contractID,
			proxy: s.stranger,
			from: s.customer,
			amount: []collection.Coin{amount},
		},
		"insufficient funds": {
			contractID: s.contractID,
			proxy: s.operator,
			from: s.customer,
			amount: []collection.Coin{collection.NewCoin("deadbeef" + fmt.Sprintf("%08x", 0), sdk.OneInt())},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgBurnFTFrom{
				ContractId: tc.contractID,
				Proxy: tc.proxy.String(),
				From:  tc.from.String(),
				Amount: tc.amount,
			}
			res, err := s.msgServer.BurnFTFrom(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgBurnNFT() {
	param := s.nftClassID + fmt.Sprintf("%08x", 1)
	testCases := map[string]struct {
		contractID string
		from sdk.AccAddress
		tokenIDs []string
		valid  bool
	}{
		"valid request": {
			contractID: s.contractID,
			from: s.vendor,
			tokenIDs: []string{param},
			valid:  true,
		},
		"no permission": {
			contractID: s.contractID,
			from: s.customer,
			tokenIDs: []string{param},
		},
		"insufficient funds": {
			contractID: s.contractID,
			from: s.vendor,
			tokenIDs: []string{"deadbeef" + fmt.Sprintf("%08x", 1)},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgBurnNFT{
				ContractId: tc.contractID,
				From:  tc.from.String(),
				TokenIds: tc.tokenIDs,
			}
			res, err := s.msgServer.BurnNFT(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgBurnNFTFrom() {
	param := s.nftClassID + fmt.Sprintf("%08x", 3)
	testCases := map[string]struct {
		contractID string
		proxy sdk.AccAddress
		from sdk.AccAddress
		tokenIDs []string
		valid  bool
	}{
		"valid request": {
			contractID: s.contractID,
			proxy: s.operator,
			from: s.customer,
			tokenIDs: []string{param},
			valid:  true,
		},
		"no authorization": {
			contractID: s.contractID,
			proxy: s.vendor,
			from: s.customer,
			tokenIDs: []string{param},
		},
		"no permission": {
			contractID: s.contractID,
			proxy: s.stranger,
			from: s.customer,
			tokenIDs: []string{param},
		},
		"insufficient funds": {
			contractID: s.contractID,
			proxy: s.operator,
			from: s.customer,
			tokenIDs: []string{"deadbeef" + fmt.Sprintf("%08x", 1)},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgBurnNFTFrom{
				ContractId: tc.contractID,
				Proxy: tc.proxy.String(),
				From:  tc.from.String(),
				TokenIds: tc.tokenIDs,
			}
			res, err := s.msgServer.BurnNFTFrom(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgGrant() {
	testCases := map[string]struct {
		granter sdk.AccAddress
		grantee sdk.AccAddress
		permission  string
		valid   bool
	}{
		"valid request": {
			granter: s.vendor,
			grantee: s.operator,
			permission:  collection.Permission_Modify.String(),
			valid:   true,
		},
		"already granted": {
			granter: s.vendor,
			grantee: s.operator,
			permission:  collection.Permission_Mint.String(),
		},
		"granter has no permission": {
			granter: s.customer,
			grantee: s.operator,
			permission:  collection.Permission_Modify.String(),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgGrant{
				ContractId: s.contractID,
				Granter: tc.granter.String(),
				Grantee: tc.grantee.String(),
				Permission:  tc.permission,
			}
			res, err := s.msgServer.Grant(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgAbandon() {
	testCases := map[string]struct {
		grantee sdk.AccAddress
		permission  string
		valid   bool
	}{
		"valid request": {
			grantee: s.operator,
			permission:  collection.Permission_Mint.String(),
			valid:   true,
		},
		"not granted yet": {
			grantee: s.operator,
			permission:  collection.Permission_Modify.String(),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgAbandon{
				ContractId: s.contractID,
				Grantee: tc.grantee.String(),
				Permission:  tc.permission,
			}
			res, err := s.msgServer.Abandon(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgGrantPermission() {
	testCases := map[string]struct {
		granter sdk.AccAddress
		grantee sdk.AccAddress
		permission  string
		valid   bool
	}{
		"valid request": {
			granter: s.vendor,
			grantee: s.operator,
			permission:  collection.Permission_Modify.String(),
			valid:   true,
		},
		"already granted": {
			granter: s.vendor,
			grantee: s.operator,
			permission:  collection.Permission_Mint.String(),
		},
		"granter has no permission": {
			granter: s.customer,
			grantee: s.operator,
			permission:  collection.Permission_Modify.String(),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgGrantPermission{
				ContractId: s.contractID,
				From: tc.granter.String(),
				To: tc.grantee.String(),
				Permission:  tc.permission,
			}
			res, err := s.msgServer.GrantPermission(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgRevokePermission() {
	testCases := map[string]struct {
		from sdk.AccAddress
		permission  string
		valid   bool
	}{
		"valid request": {
			from: s.operator,
			permission:  collection.Permission_Mint.String(),
			valid:   true,
		},
		"not granted yet": {
			from: s.operator,
			permission:  collection.Permission_Modify.String(),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgRevokePermission{
				ContractId: s.contractID,
				From: tc.from.String(),
				Permission:  tc.permission,
			}
			res, err := s.msgServer.RevokePermission(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}
