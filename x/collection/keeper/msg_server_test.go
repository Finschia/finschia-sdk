package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/collection"
)

func (s *KeeperTestSuite) TestMsgSend() {
	testCases := map[string]struct {
		contractID string
		amount     sdk.Int
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			amount:     s.balance,
			valid:      true,
		},
		"insufficient funds": {
			contractID: "deadbeef",
			amount:     s.balance,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgSend{
				ContractId: tc.contractID,
				From:       s.vendor.String(),
				To:         s.customer.String(),
				Amount: collection.NewCoins(
					collection.NewFTCoin(s.ftClassID, tc.amount),
				),
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
		from     sdk.AccAddress
		amount   sdk.Int
		valid    bool
	}{
		"valid request": {
			operator: s.operator,
			from:     s.customer,
			amount:   s.balance,
			valid:    true,
		},
		"not approved": {
			operator: s.vendor,
			from:     s.customer,
			amount:   s.balance,
		},
		"insufficient funds": {
			operator: s.operator,
			from:     s.customer,
			amount:   s.balance.Add(sdk.OneInt()),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgOperatorSend{
				ContractId: s.contractID,
				Operator:   tc.operator.String(),
				From:       tc.from.String(),
				To:         s.vendor.String(),
				Amount: collection.NewCoins(
					collection.NewFTCoin(s.ftClassID, tc.amount),
				),
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
		amount     sdk.Int
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			amount:     s.balance,
			valid:      true,
		},
		"insufficient funds": {
			contractID: "deadbeef",
			amount:     s.balance,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgTransferFT{
				ContractId: tc.contractID,
				From:       s.vendor.String(),
				To:         s.customer.String(),
				Amount: collection.NewCoins(
					collection.NewFTCoin(s.ftClassID, tc.amount),
				),
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
		proxy  sdk.AccAddress
		from   sdk.AccAddress
		amount sdk.Int
		valid  bool
	}{
		"valid request": {
			proxy:  s.operator,
			from:   s.customer,
			amount: s.balance,
			valid:  true,
		},
		"not approved": {
			proxy:  s.vendor,
			from:   s.customer,
			amount: s.balance,
		},
		"insufficient funds": {
			proxy:  s.operator,
			from:   s.customer,
			amount: s.balance.Add(sdk.OneInt()),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgTransferFTFrom{
				ContractId: s.contractID,
				Proxy:      tc.proxy.String(),
				From:       tc.from.String(),
				To:         s.vendor.String(),
				Amount: collection.NewCoins(
					collection.NewFTCoin(s.ftClassID, tc.amount),
				),
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
			tokenID: collection.NewNFTID(s.nftClassID, 1),
			valid:   true,
		},
		"insufficient funds": {
			tokenID: collection.NewNFTID("deadbeef", 1),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgTransferNFT{
				ContractId: s.contractID,
				From:       s.customer.String(),
				To:         s.vendor.String(),
				TokenIds:   []string{tc.tokenID},
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
	tokenID := collection.NewNFTID(s.nftClassID, 1)
	testCases := map[string]struct {
		proxy   sdk.AccAddress
		from    sdk.AccAddress
		tokenID string
		valid   bool
	}{
		"valid request": {
			proxy:   s.operator,
			from:    s.customer,
			tokenID: tokenID,
			valid:   true,
		},
		"not approved": {
			proxy:   s.vendor,
			from:    s.customer,
			tokenID: tokenID,
		},
		"insufficient funds": {
			proxy:   s.operator,
			from:    s.customer,
			tokenID: collection.NewNFTID("deadbeef", 1),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgTransferNFTFrom{
				ContractId: s.contractID,
				Proxy:      tc.proxy.String(),
				From:       tc.from.String(),
				To:         s.vendor.String(),
				TokenIds:   []string{tc.tokenID},
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
		holder   sdk.AccAddress
		operator sdk.AccAddress
		valid    bool
	}{
		"valid request": {
			holder:   s.customer,
			operator: s.vendor,
			valid:    true,
		},
		"already approved": {
			holder:   s.customer,
			operator: s.operator,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgAuthorizeOperator{
				ContractId: s.contractID,
				Holder:     tc.holder.String(),
				Operator:   tc.operator.String(),
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
		holder   sdk.AccAddress
		operator sdk.AccAddress
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
				ContractId: s.contractID,
				Approver:   tc.approver.String(),
				Proxy:      tc.proxy.String(),
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
				ContractId: s.contractID,
				Approver:   tc.approver.String(),
				Proxy:      tc.proxy.String(),
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
		valid bool
	}{
		"valid request": {
			owner: s.vendor,
			valid: true,
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

func (s *KeeperTestSuite) TestMsgCreateFTClass() {
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		supply     sdk.Int
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.vendor,
			supply:     sdk.ZeroInt(),
			valid:      true,
		},
		"valid request with supply": {
			contractID: s.contractID,
			operator:   s.vendor,
			supply:     sdk.OneInt(),
			valid:      true,
		},
		"no permission": {
			contractID: s.contractID,
			operator:   s.customer,
			supply:     sdk.ZeroInt(),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgCreateFTClass{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				To:         s.customer.String(),
				Supply:     tc.supply,
			}
			res, err := s.msgServer.CreateFTClass(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgCreateNFTClass() {
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.vendor,
			valid:      true,
		},
		"no permission": {
			contractID: s.contractID,
			operator:   s.customer,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgCreateNFTClass{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
			}
			res, err := s.msgServer.CreateNFTClass(sdk.WrapSDKContext(ctx), req)
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
		owner      sdk.AccAddress
		amount     sdk.Int
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			owner:      s.vendor,
			amount:     sdk.ZeroInt(),
			valid:      true,
		},
		"valid request with supply": {
			contractID: s.contractID,
			owner:      s.vendor,
			amount:     sdk.OneInt(),
			valid:      true,
		},
		"no permission": {
			contractID: s.contractID,
			owner:      s.customer,
			amount:     sdk.ZeroInt(),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgIssueFT{
				ContractId: tc.contractID,
				Owner:      tc.owner.String(),
				Amount:     tc.amount,
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
		owner      sdk.AccAddress
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			owner:      s.vendor,
			valid:      true,
		},
		"no permission": {
			contractID: s.contractID,
			owner:      s.customer,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgIssueNFT{
				ContractId: tc.contractID,
				Owner:      tc.owner.String(),
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
	amount := collection.NewCoins(
		collection.NewFTCoin(s.ftClassID, sdk.OneInt()),
	)
	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		amount     []collection.Coin
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			from:       s.vendor,
			amount:     amount,
			valid:      true,
		},
		"no permission": {
			contractID: s.contractID,
			from:       s.customer,
			amount:     amount,
		},
		"no class of the token": {
			contractID: s.contractID,
			from:       s.vendor,
			amount: collection.NewCoins(
				collection.NewFTCoin("00bab10c", sdk.OneInt()),
			),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgMintFT{
				ContractId: tc.contractID,
				From:       tc.from.String(),
				To:         s.customer.String(),
				Amount:     tc.amount,
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
	params := []collection.MintNFTParam{{
		TokenType: s.nftClassID,
	}}
	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		params     []collection.MintNFTParam
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			from:       s.vendor,
			params:     params,
			valid:      true,
		},
		"no permission": {
			contractID: s.contractID,
			from:       s.customer,
			params:     params,
		},
		"no class of the token": {
			contractID: s.contractID,
			from:       s.vendor,
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
				From:       tc.from.String(),
				To:         s.customer.String(),
				Params:     tc.params,
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
	amount := collection.NewCoins(
		collection.NewFTCoin(s.ftClassID, s.balance),
	)
	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		amount     []collection.Coin
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			from:       s.vendor,
			amount:     amount,
			valid:      true,
		},
		"no permission": {
			contractID: s.contractID,
			from:       s.customer,
			amount:     amount,
		},
		"insufficient funds": {
			contractID: s.contractID,
			from:       s.vendor,
			amount: collection.NewCoins(
				collection.NewFTCoin("00bab10c", sdk.OneInt()),
			),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgBurn{
				ContractId: tc.contractID,
				From:       tc.from.String(),
				Amount:     tc.amount,
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
	amount := collection.NewCoins(
		collection.NewFTCoin(s.ftClassID, s.balance),
	)
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		from       sdk.AccAddress
		amount     []collection.Coin
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			amount:     amount,
			valid:      true,
		},
		"no authorization": {
			contractID: s.contractID,
			operator:   s.vendor,
			from:       s.customer,
			amount:     amount,
		},
		"no permission": {
			contractID: s.contractID,
			operator:   s.stranger,
			from:       s.customer,
			amount:     amount,
		},
		"insufficient funds": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			amount: collection.NewCoins(
				collection.NewFTCoin("00bab10c", sdk.OneInt()),
			),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgOperatorBurn{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				From:       tc.from.String(),
				Amount:     tc.amount,
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
	amount := collection.NewCoins(
		collection.NewFTCoin(s.ftClassID, s.balance),
	)
	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		amount     []collection.Coin
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			from:       s.vendor,
			amount:     amount,
			valid:      true,
		},
		"no permission": {
			contractID: s.contractID,
			from:       s.customer,
			amount:     amount,
		},
		"insufficient funds": {
			contractID: s.contractID,
			from:       s.vendor,
			amount: collection.NewCoins(
				collection.NewFTCoin("00bab10c", sdk.OneInt()),
			),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgBurnFT{
				ContractId: tc.contractID,
				From:       tc.from.String(),
				Amount:     tc.amount,
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
	amount := collection.NewCoins(
		collection.NewFTCoin(s.ftClassID, s.balance),
	)
	testCases := map[string]struct {
		contractID string
		proxy      sdk.AccAddress
		from       sdk.AccAddress
		amount     []collection.Coin
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			proxy:      s.operator,
			from:       s.customer,
			amount:     amount,
			valid:      true,
		},
		"no authorization": {
			contractID: s.contractID,
			proxy:      s.vendor,
			from:       s.customer,
			amount:     amount,
		},
		"no permission": {
			contractID: s.contractID,
			proxy:      s.stranger,
			from:       s.customer,
			amount:     amount,
		},
		"insufficient funds": {
			contractID: s.contractID,
			proxy:      s.operator,
			from:       s.customer,
			amount: collection.NewCoins(
				collection.NewFTCoin("00bab10c", sdk.OneInt()),
			),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgBurnFTFrom{
				ContractId: tc.contractID,
				Proxy:      tc.proxy.String(),
				From:       tc.from.String(),
				Amount:     tc.amount,
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
	tokenIDs := []string{
		collection.NewNFTID(s.nftClassID, s.numNFTs*2+1),
	}
	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		tokenIDs   []string
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			from:       s.vendor,
			tokenIDs:   tokenIDs,
			valid:      true,
		},
		"no permission": {
			contractID: s.contractID,
			from:       s.customer,
			tokenIDs:   tokenIDs,
		},
		"insufficient funds": {
			contractID: s.contractID,
			from:       s.vendor,
			tokenIDs: []string{
				collection.NewNFTID("deadbeef", 1),
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgBurnNFT{
				ContractId: tc.contractID,
				From:       tc.from.String(),
				TokenIds:   tc.tokenIDs,
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
	tokenIDs := []string{
		collection.NewNFTID(s.nftClassID, 1),
	}
	testCases := map[string]struct {
		contractID string
		proxy      sdk.AccAddress
		from       sdk.AccAddress
		tokenIDs   []string
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			proxy:      s.operator,
			from:       s.customer,
			tokenIDs:   tokenIDs,
			valid:      true,
		},
		"no authorization": {
			contractID: s.contractID,
			proxy:      s.vendor,
			from:       s.customer,
			tokenIDs:   tokenIDs,
		},
		"no permission": {
			contractID: s.contractID,
			proxy:      s.stranger,
			from:       s.customer,
			tokenIDs:   tokenIDs,
		},
		"insufficient funds": {
			contractID: s.contractID,
			proxy:      s.operator,
			from:       s.customer,
			tokenIDs: []string{
				collection.NewNFTID("deadbeef", 1),
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgBurnNFTFrom{
				ContractId: tc.contractID,
				Proxy:      tc.proxy.String(),
				From:       tc.from.String(),
				TokenIds:   tc.tokenIDs,
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

func (s *KeeperTestSuite) TestMsgModifyContract() {
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.vendor,
			valid:      true,
		},
		"no permission": {
			contractID: s.contractID,
			operator:   s.customer,
		},
		"contract not found": {
			contractID: "deadbeef",
			operator:   s.vendor,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			changes := []collection.Attribute{{
				Key:   collection.AttributeKeyName.String(),
				Value: "fox",
			}}
			req := &collection.MsgModifyContract{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				Changes:    changes,
			}
			res, err := s.msgServer.ModifyContract(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgModifyTokenClass() {
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		classID    string
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.vendor,
			classID:    s.nftClassID,
			valid:      true,
		},
		"no permission": {
			contractID: s.contractID,
			operator:   s.customer,
			classID:    s.nftClassID,
		},
		"token class not found": {
			contractID: s.contractID,
			operator:   s.vendor,
			classID:    "deadbeef",
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			changes := []collection.Attribute{{
				Key:   collection.AttributeKeyName.String(),
				Value: "arctic fox",
			}}
			req := &collection.MsgModifyTokenClass{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				ClassId:    tc.classID,
				Changes:    changes,
			}
			res, err := s.msgServer.ModifyTokenClass(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgModifyNFT() {
	tokenID := collection.NewNFTID(s.nftClassID, 1)
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		tokenID    string
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.vendor,
			tokenID:    tokenID,
			valid:      true,
		},
		"no permission": {
			contractID: s.contractID,
			operator:   s.customer,
			tokenID:    tokenID,
		},
		"token not found": {
			contractID: s.contractID,
			operator:   s.vendor,
			tokenID:    collection.NewNFTID("deadbeef", 1),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			changes := []collection.Attribute{{
				Key:   collection.AttributeKeyName.String(),
				Value: "fennec fox 1",
			}}
			req := &collection.MsgModifyNFT{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				TokenId:    tc.tokenID,
				Changes:    changes,
			}
			res, err := s.msgServer.ModifyNFT(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgModify() {
	tokenIndex := collection.NewNFTID(s.nftClassID, 1)[8:]
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		tokenType  string
		tokenIndex string
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.vendor,
			valid:      true,
		},
		"no permission": {
			contractID: s.contractID,
			operator:   s.customer,
			tokenType:  s.nftClassID,
			tokenIndex: tokenIndex,
		},
		"nft not found": {
			contractID: s.contractID,
			operator:   s.vendor,
			tokenType:  s.nftClassID,
			tokenIndex: collection.NewNFTID(s.nftClassID, s.numNFTs*3+1)[8:],
		},
		"ft class not found": {
			contractID: s.contractID,
			operator:   s.vendor,
			tokenType:  "00bab10c",
			tokenIndex: collection.NewFTID("00bab10c")[8:],
		},
		"nft class not found": {
			contractID: s.contractID,
			operator:   s.vendor,
			tokenType:  "deadbeef",
		},
		"token index without type": {
			contractID: s.contractID,
			operator:   s.vendor,
			tokenIndex: "deadbeef",
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			changes := []collection.Change{{
				Field: collection.AttributeKeyName.String(),
				Value: "test",
			}}
			req := &collection.MsgModify{
				ContractId: tc.contractID,
				Owner:      tc.operator.String(),
				TokenType:  tc.tokenType,
				TokenIndex: tc.tokenIndex,
				Changes:    changes,
			}
			res, err := s.msgServer.Modify(sdk.WrapSDKContext(ctx), req)
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
		granter    sdk.AccAddress
		grantee    sdk.AccAddress
		permission collection.Permission
		valid      bool
	}{
		"valid request": {
			granter:    s.vendor,
			grantee:    s.operator,
			permission: collection.PermissionModify,
			valid:      true,
		},
		"already granted": {
			granter:    s.vendor,
			grantee:    s.operator,
			permission: collection.PermissionMint,
		},
		"granter has no permission": {
			granter:    s.customer,
			grantee:    s.operator,
			permission: collection.PermissionModify,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgGrant{
				ContractId: s.contractID,
				Granter:    tc.granter.String(),
				Grantee:    tc.grantee.String(),
				Permission: tc.permission,
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
		grantee    sdk.AccAddress
		permission collection.Permission
		valid      bool
	}{
		"valid request": {
			grantee:    s.operator,
			permission: collection.PermissionMint,
			valid:      true,
		},
		"not granted yet": {
			grantee:    s.operator,
			permission: collection.PermissionModify,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgAbandon{
				ContractId: s.contractID,
				Grantee:    tc.grantee.String(),
				Permission: tc.permission,
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
		granter    sdk.AccAddress
		grantee    sdk.AccAddress
		permission string
		valid      bool
	}{
		"valid request": {
			granter:    s.vendor,
			grantee:    s.operator,
			permission: collection.LegacyPermissionModify.String(),
			valid:      true,
		},
		"already granted": {
			granter:    s.vendor,
			grantee:    s.operator,
			permission: collection.LegacyPermissionMint.String(),
		},
		"granter has no permission": {
			granter:    s.customer,
			grantee:    s.operator,
			permission: collection.LegacyPermissionModify.String(),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgGrantPermission{
				ContractId: s.contractID,
				From:       tc.granter.String(),
				To:         tc.grantee.String(),
				Permission: tc.permission,
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
		from       sdk.AccAddress
		permission string
		valid      bool
	}{
		"valid request": {
			from:       s.operator,
			permission: collection.LegacyPermissionMint.String(),
			valid:      true,
		},
		"not granted yet": {
			from:       s.operator,
			permission: collection.LegacyPermissionModify.String(),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgRevokePermission{
				ContractId: s.contractID,
				From:       tc.from.String(),
				Permission: tc.permission,
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

func (s *KeeperTestSuite) TestMsgAttach() {
	testCases := map[string]struct {
		contractID string
		subjectID  string
		targetID   string
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			subjectID:  collection.NewNFTID(s.nftClassID, collection.DefaultDepthLimit+1),
			targetID:   collection.NewNFTID(s.nftClassID, 1),
			valid:      true,
		},
		"not owner of the token": {
			contractID: s.contractID,
			subjectID:  collection.NewNFTID(s.nftClassID, s.numNFTs+1),
			targetID:   collection.NewNFTID(s.nftClassID, 1),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgAttach{
				ContractId: tc.contractID,
				From:       s.customer.String(),
				TokenId:    tc.subjectID,
				ToTokenId:  tc.targetID,
			}
			res, err := s.msgServer.Attach(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgDetach() {
	testCases := map[string]struct {
		contractID string
		subjectID  string
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			subjectID:  collection.NewNFTID(s.nftClassID, 2),
			valid:      true,
		},
		"not owner of the token": {
			contractID: s.contractID,
			subjectID:  collection.NewNFTID(s.nftClassID, s.numNFTs+2),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgDetach{
				ContractId: tc.contractID,
				From:       s.customer.String(),
				TokenId:    tc.subjectID,
			}
			res, err := s.msgServer.Detach(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgOperatorAttach() {
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		subjectID  string
		targetID   string
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			subjectID:  collection.NewNFTID(s.nftClassID, collection.DefaultDepthLimit+1),
			targetID:   collection.NewNFTID(s.nftClassID, 1),
			valid:      true,
		},
		"not authorized": {
			contractID: s.contractID,
			operator:   s.vendor,
			subjectID:  collection.NewNFTID(s.nftClassID, collection.DefaultDepthLimit+1),
			targetID:   collection.NewNFTID(s.nftClassID, 1),
		},
		"not owner of the token": {
			contractID: s.contractID,
			operator:   s.operator,
			subjectID:  collection.NewNFTID(s.nftClassID, s.numNFTs+1),
			targetID:   collection.NewNFTID(s.nftClassID, 1),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgOperatorAttach{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				Owner:      s.customer.String(),
				Subject:    tc.subjectID,
				Target:     tc.targetID,
			}
			res, err := s.msgServer.OperatorAttach(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgOperatorDetach() {
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		subjectID  string
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			subjectID:  collection.NewNFTID(s.nftClassID, 2),
			valid:      true,
		},
		"not authorized": {
			contractID: s.contractID,
			operator:   s.vendor,
			subjectID:  collection.NewNFTID(s.nftClassID, 2),
		},
		"not owner of the token": {
			contractID: s.contractID,
			operator:   s.operator,
			subjectID:  collection.NewNFTID(s.nftClassID, s.numNFTs+2),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgOperatorDetach{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				Owner:      s.customer.String(),
				Subject:    tc.subjectID,
			}
			res, err := s.msgServer.OperatorDetach(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgAttachFrom() {
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		subjectID  string
		targetID   string
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			subjectID:  collection.NewNFTID(s.nftClassID, collection.DefaultDepthLimit+1),
			targetID:   collection.NewNFTID(s.nftClassID, 1),
			valid:      true,
		},
		"not authorized": {
			contractID: s.contractID,
			operator:   s.vendor,
			subjectID:  collection.NewNFTID(s.nftClassID, collection.DefaultDepthLimit+1),
			targetID:   collection.NewNFTID(s.nftClassID, 1),
		},
		"not owner of the token": {
			contractID: s.contractID,
			operator:   s.operator,
			subjectID:  collection.NewNFTID(s.nftClassID, s.numNFTs+1),
			targetID:   collection.NewNFTID(s.nftClassID, 1),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgAttachFrom{
				ContractId: tc.contractID,
				Proxy:      tc.operator.String(),
				From:       s.customer.String(),
				TokenId:    tc.subjectID,
				ToTokenId:  tc.targetID,
			}
			res, err := s.msgServer.AttachFrom(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgDetachFrom() {
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		subjectID  string
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			subjectID:  collection.NewNFTID(s.nftClassID, 2),
			valid:      true,
		},
		"not authorized": {
			contractID: s.contractID,
			operator:   s.vendor,
			subjectID:  collection.NewNFTID(s.nftClassID, 2),
		},
		"not owner of the token": {
			contractID: s.contractID,
			operator:   s.operator,
			subjectID:  collection.NewNFTID(s.nftClassID, s.numNFTs+2),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgDetachFrom{
				ContractId: tc.contractID,
				Proxy:      tc.operator.String(),
				From:       s.customer.String(),
				TokenId:    tc.subjectID,
			}
			res, err := s.msgServer.DetachFrom(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}
