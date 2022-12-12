package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/collection"
)

func (s *KeeperTestSuite) TestMsgTransferFT() {
	testCases := map[string]struct {
		contractID string
		amount     sdk.Int
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			amount:     s.balance,
		},
		"contract not found": {
			contractID: "deadbeef",
			amount:     s.balance,
			err:        collection.ErrNotFound,
		},
		"insufficient funds": {
			contractID: s.contractID,
			amount:     s.balance.Add(sdk.OneInt()),
			err:        collection.ErrInsufficientTokens,
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
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgTransferFTFrom() {
	testCases := map[string]struct {
		contractID string
		proxy      sdk.AccAddress
		from       sdk.AccAddress
		amount     sdk.Int
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			proxy:      s.operator,
			from:       s.customer,
			amount:     s.balance,
		},
		"contract not found": {
			contractID: "deadbeef",
			proxy:      s.operator,
			from:       s.customer,
			amount:     s.balance,
			err:        collection.ErrNotFound,
		},
		"not approved": {
			contractID: s.contractID,
			proxy:      s.vendor,
			from:       s.customer,
			amount:     s.balance,
			err:        sdkerrors.ErrUnauthorized,
		},
		"insufficient funds": {
			contractID: s.contractID,
			proxy:      s.operator,
			from:       s.customer,
			amount:     s.balance.Add(sdk.OneInt()),
			err:        collection.ErrInsufficientTokens,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgTransferFTFrom{
				ContractId: tc.contractID,
				Proxy:      tc.proxy.String(),
				From:       tc.from.String(),
				To:         s.vendor.String(),
				Amount: collection.NewCoins(
					collection.NewFTCoin(s.ftClassID, tc.amount),
				),
			}
			res, err := s.msgServer.TransferFTFrom(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgTransferNFT() {
	testCases := map[string]struct {
		contractID string
		tokenID    string
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			tokenID:    collection.NewNFTID(s.nftClassID, 1),
		},
		"contract not found": {
			contractID: "deadbeef",
			tokenID:    collection.NewNFTID(s.nftClassID, 1),
			err:        collection.ErrNotFound,
		},
		"insufficient funds": {
			contractID: s.contractID,
			tokenID:    collection.NewNFTID("deadbeef", 1),
			err:        collection.ErrInsufficientTokens,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgTransferNFT{
				ContractId: tc.contractID,
				From:       s.customer.String(),
				To:         s.vendor.String(),
				TokenIds:   []string{tc.tokenID},
			}
			res, err := s.msgServer.TransferNFT(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgTransferNFTFrom() {
	tokenID := collection.NewNFTID(s.nftClassID, 1)
	testCases := map[string]struct {
		contractID string
		proxy      sdk.AccAddress
		from       sdk.AccAddress
		tokenID    string
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			proxy:      s.operator,
			from:       s.customer,
			tokenID:    tokenID,
		},
		"contract not found": {
			contractID: "deadbeef",
			proxy:      s.operator,
			from:       s.customer,
			tokenID:    tokenID,
			err:        collection.ErrNotFound,
		},
		"not approved": {
			contractID: s.contractID,
			proxy:      s.vendor,
			from:       s.customer,
			tokenID:    tokenID,
			err:        sdkerrors.ErrUnauthorized,
		},
		"insufficient funds": {
			contractID: s.contractID,
			proxy:      s.operator,
			from:       s.customer,
			tokenID:    collection.NewNFTID("deadbeef", 1),
			err:        collection.ErrInsufficientTokens,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgTransferNFTFrom{
				ContractId: tc.contractID,
				Proxy:      tc.proxy.String(),
				From:       tc.from.String(),
				To:         s.vendor.String(),
				TokenIds:   []string{tc.tokenID},
			}
			res, err := s.msgServer.TransferNFTFrom(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgApprove() {
	testCases := map[string]struct {
		contractID string
		approver   sdk.AccAddress
		proxy      sdk.AccAddress
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			approver:   s.customer,
			proxy:      s.vendor,
		},
		"contract not found": {
			contractID: "deadbeef",
			approver:   s.customer,
			proxy:      s.vendor,
			err:        collection.ErrNotFound,
		},
		"already approved": {
			contractID: s.contractID,
			approver:   s.customer,
			proxy:      s.operator,
			err:        collection.ErrAlreadyExists,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgApprove{
				ContractId: tc.contractID,
				Approver:   tc.approver.String(),
				Proxy:      tc.proxy.String(),
			}
			res, err := s.msgServer.Approve(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgDisapprove() {
	testCases := map[string]struct {
		contractID string
		approver   sdk.AccAddress
		proxy      sdk.AccAddress
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			approver:   s.customer,
			proxy:      s.operator,
		},
		"contract not found": {
			contractID: "deadbeef",
			approver:   s.customer,
			proxy:      s.operator,
			err:        collection.ErrNotFound,
		},
		"no authorization": {
			contractID: s.contractID,
			approver:   s.customer,
			proxy:      s.vendor,
			err:        collection.ErrNotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgDisapprove{
				ContractId: tc.contractID,
				Approver:   tc.approver.String(),
				Proxy:      tc.proxy.String(),
			}
			res, err := s.msgServer.Disapprove(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgCreateContract() {
	testCases := map[string]struct {
		owner sdk.AccAddress
		err   error
	}{
		"valid request": {
			owner: s.vendor,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgCreateContract{
				Owner: tc.owner.String(),
			}
			res, err := s.msgServer.CreateContract(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgIssueFT() {
	testCases := map[string]struct {
		contractID string
		owner      sdk.AccAddress
		amount     sdk.Int
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			owner:      s.vendor,
			amount:     sdk.ZeroInt(),
		},
		"valid request with supply": {
			contractID: s.contractID,
			owner:      s.vendor,
			amount:     sdk.OneInt(),
		},
		"contract not found": {
			contractID: "deadbeef",
			owner:      s.vendor,
			amount:     sdk.ZeroInt(),
			err:        collection.ErrNotFound,
		},
		"no permission": {
			contractID: s.contractID,
			owner:      s.customer,
			amount:     sdk.ZeroInt(),
			err:        sdkerrors.ErrUnauthorized,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgIssueFT{
				ContractId: tc.contractID,
				Owner:      tc.owner.String(),
				To:         s.customer.String(),
				Amount:     tc.amount,
			}
			res, err := s.msgServer.IssueFT(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgIssueNFT() {
	testCases := map[string]struct {
		contractID string
		owner      sdk.AccAddress
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			owner:      s.vendor,
		},
		"contract not found": {
			contractID: "deadbeef",
			owner:      s.vendor,
			err:        collection.ErrNotFound,
		},
		"no permission": {
			contractID: s.contractID,
			owner:      s.customer,
			err:        sdkerrors.ErrUnauthorized,
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
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

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
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			from:       s.vendor,
			amount:     amount,
		},
		"contract not found": {
			contractID: "deadbeef",
			from:       s.vendor,
			amount:     amount,
			err:        collection.ErrNotFound,
		},
		"no permission": {
			contractID: s.contractID,
			from:       s.customer,
			amount:     amount,
			err:        sdkerrors.ErrUnauthorized,
		},
		"no class of the token": {
			contractID: s.contractID,
			from:       s.vendor,
			amount: collection.NewCoins(
				collection.NewFTCoin("00bab10c", sdk.OneInt()),
			),
			err: collection.ErrNotFound,
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
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

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
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			from:       s.vendor,
			params:     params,
		},
		"contract not found": {
			contractID: "deadbeef",
			from:       s.vendor,
			params:     params,
			err:        collection.ErrNotFound,
		},
		"no permission": {
			contractID: s.contractID,
			from:       s.customer,
			params:     params,
			err:        sdkerrors.ErrUnauthorized,
		},
		"no class of the token": {
			contractID: s.contractID,
			from:       s.vendor,
			params: []collection.MintNFTParam{{
				TokenType: "deadbeef",
			}},
			err: collection.ErrNotFound,
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
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

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
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			from:       s.vendor,
			amount:     amount,
		},
		"contract not found": {
			contractID: "deadbeef",
			from:       s.vendor,
			amount:     amount,
			err:        collection.ErrNotFound,
		},
		"no permission": {
			contractID: s.contractID,
			from:       s.customer,
			amount:     amount,
			err:        sdkerrors.ErrUnauthorized,
		},
		"insufficient funds": {
			contractID: s.contractID,
			from:       s.vendor,
			amount: collection.NewCoins(
				collection.NewFTCoin("00bab10c", sdk.OneInt()),
			),
			err: collection.ErrInsufficientTokens,
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
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

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
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			proxy:      s.operator,
			from:       s.customer,
			amount:     amount,
		},
		"contract not found": {
			contractID: "deadbeef",
			proxy:      s.operator,
			from:       s.customer,
			amount:     amount,
			err:        collection.ErrNotFound,
		},
		"no authorization": {
			contractID: s.contractID,
			proxy:      s.vendor,
			from:       s.customer,
			amount:     amount,
			err:        sdkerrors.ErrUnauthorized,
		},
		"no permission": {
			contractID: s.contractID,
			proxy:      s.stranger,
			from:       s.customer,
			amount:     amount,
			err:        sdkerrors.ErrUnauthorized,
		},
		"insufficient funds": {
			contractID: s.contractID,
			proxy:      s.operator,
			from:       s.customer,
			amount: collection.NewCoins(
				collection.NewFTCoin("00bab10c", sdk.OneInt()),
			),
			err: collection.ErrInsufficientTokens,
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
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

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
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			from:       s.vendor,
			tokenIDs:   tokenIDs,
		},
		"contract not found": {
			contractID: "deadbeef",
			from:       s.vendor,
			tokenIDs:   tokenIDs,
			err:        collection.ErrNotFound,
		},
		"no permission": {
			contractID: s.contractID,
			from:       s.customer,
			tokenIDs:   tokenIDs,
			err:        sdkerrors.ErrUnauthorized,
		},
		"insufficient funds": {
			contractID: s.contractID,
			from:       s.vendor,
			tokenIDs: []string{
				collection.NewNFTID("deadbeef", 1),
			},
			err: collection.ErrInsufficientTokens,
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
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

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
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			proxy:      s.operator,
			from:       s.customer,
			tokenIDs:   tokenIDs,
		},
		"contract not found": {
			contractID: "deadbeef",
			proxy:      s.operator,
			from:       s.customer,
			tokenIDs:   tokenIDs,
			err:        collection.ErrNotFound,
		},
		"no authorization": {
			contractID: s.contractID,
			proxy:      s.vendor,
			from:       s.customer,
			tokenIDs:   tokenIDs,
			err:        sdkerrors.ErrUnauthorized,
		},
		"no permission": {
			contractID: s.contractID,
			proxy:      s.stranger,
			from:       s.customer,
			tokenIDs:   tokenIDs,
			err:        sdkerrors.ErrUnauthorized,
		},
		"insufficient funds": {
			contractID: s.contractID,
			proxy:      s.operator,
			from:       s.customer,
			tokenIDs: []string{
				collection.NewNFTID("deadbeef", 1),
			},
			err: collection.ErrInsufficientTokens,
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
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

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
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.vendor,
		},
		"contract not found": {
			contractID: "deadbeef",
			operator:   s.vendor,
			err:        collection.ErrNotFound,
		},
		"no permission": {
			contractID: s.contractID,
			operator:   s.customer,
			tokenType:  s.nftClassID,
			tokenIndex: tokenIndex,
			err:        sdkerrors.ErrUnauthorized,
		},
		"nft not found": {
			contractID: s.contractID,
			operator:   s.vendor,
			tokenType:  s.nftClassID,
			tokenIndex: collection.NewNFTID(s.nftClassID, s.numNFTs*3+1)[8:],
			err:        collection.ErrNotFound,
		},
		"ft class not found": {
			contractID: s.contractID,
			operator:   s.vendor,
			tokenType:  "00bab10c",
			tokenIndex: collection.NewFTID("00bab10c")[8:],
			err:        collection.ErrNotFound,
		},
		"nft class not found": {
			contractID: s.contractID,
			operator:   s.vendor,
			tokenType:  "deadbeef",
			err:        collection.ErrNotFound,
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
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgGrantPermission() {
	testCases := map[string]struct {
		contractID string
		granter    sdk.AccAddress
		grantee    sdk.AccAddress
		permission string
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			granter:    s.vendor,
			grantee:    s.operator,
			permission: collection.LegacyPermissionModify.String(),
		},
		"contract not found": {
			contractID: "deadbeef",
			granter:    s.vendor,
			grantee:    s.operator,
			permission: collection.LegacyPermissionModify.String(),
			err:        collection.ErrNotFound,
		},
		"granter has no permission": {
			contractID: s.contractID,
			granter:    s.customer,
			grantee:    s.operator,
			permission: collection.LegacyPermissionModify.String(),
			err:        sdkerrors.ErrUnauthorized,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgGrantPermission{
				ContractId: tc.contractID,
				From:       tc.granter.String(),
				To:         tc.grantee.String(),
				Permission: tc.permission,
			}
			res, err := s.msgServer.GrantPermission(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgRevokePermission() {
	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		permission string
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			from:       s.operator,
			permission: collection.LegacyPermissionMint.String(),
		},
		"contract not found": {
			contractID: "deadbeef",
			from:       s.operator,
			permission: collection.LegacyPermissionMint.String(),
			err:        collection.ErrNotFound,
		},
		"not granted yet": {
			contractID: s.contractID,
			from:       s.operator,
			permission: collection.LegacyPermissionModify.String(),
			err:        collection.ErrNotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgRevokePermission{
				ContractId: tc.contractID,
				From:       tc.from.String(),
				Permission: tc.permission,
			}
			res, err := s.msgServer.RevokePermission(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgAttach() {
	testCases := map[string]struct {
		contractID string
		subjectID  string
		targetID   string
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			subjectID:  collection.NewNFTID(s.nftClassID, collection.DefaultDepthLimit+1),
			targetID:   collection.NewNFTID(s.nftClassID, 1),
		},
		"contract not found": {
			contractID: "deadbeef",
			subjectID:  collection.NewNFTID(s.nftClassID, collection.DefaultDepthLimit+1),
			targetID:   collection.NewNFTID(s.nftClassID, 1),
			err:        collection.ErrNotFound,
		},
		"not owner of the token": {
			contractID: s.contractID,
			subjectID:  collection.NewNFTID(s.nftClassID, s.numNFTs+1),
			targetID:   collection.NewNFTID(s.nftClassID, 1),
			err:        collection.ErrInsufficientTokens,
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
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgDetach() {
	testCases := map[string]struct {
		contractID string
		subjectID  string
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			subjectID:  collection.NewNFTID(s.nftClassID, 2),
		},
		"contract not found": {
			contractID: "deadbeef",
			subjectID:  collection.NewNFTID(s.nftClassID, 2),
			err:        collection.ErrNotFound,
		},
		"not owner of the token": {
			contractID: s.contractID,
			subjectID:  collection.NewNFTID(s.nftClassID, s.numNFTs+2),
			err:        collection.ErrInsufficientTokens,
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
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

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
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			subjectID:  collection.NewNFTID(s.nftClassID, collection.DefaultDepthLimit+1),
			targetID:   collection.NewNFTID(s.nftClassID, 1),
		},
		"contract not found": {
			contractID: "deadbeef",
			operator:   s.operator,
			subjectID:  collection.NewNFTID(s.nftClassID, collection.DefaultDepthLimit+1),
			targetID:   collection.NewNFTID(s.nftClassID, 1),
			err:        collection.ErrNotFound,
		},
		"not authorized": {
			contractID: s.contractID,
			operator:   s.vendor,
			subjectID:  collection.NewNFTID(s.nftClassID, collection.DefaultDepthLimit+1),
			targetID:   collection.NewNFTID(s.nftClassID, 1),
			err:        sdkerrors.ErrUnauthorized,
		},
		"not owner of the token": {
			contractID: s.contractID,
			operator:   s.operator,
			subjectID:  collection.NewNFTID(s.nftClassID, s.numNFTs+1),
			targetID:   collection.NewNFTID(s.nftClassID, 1),
			err:        collection.ErrInsufficientTokens,
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
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgDetachFrom() {
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		subjectID  string
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			subjectID:  collection.NewNFTID(s.nftClassID, 2),
		},
		"contract not found": {
			contractID: "deadbeef",
			operator:   s.operator,
			subjectID:  collection.NewNFTID(s.nftClassID, 2),
			err:        collection.ErrNotFound,
		},
		"not authorized": {
			contractID: s.contractID,
			operator:   s.vendor,
			subjectID:  collection.NewNFTID(s.nftClassID, 2),
			err:        sdkerrors.ErrUnauthorized,
		},
		"not owner of the token": {
			contractID: s.contractID,
			operator:   s.operator,
			subjectID:  collection.NewNFTID(s.nftClassID, s.numNFTs+2),
			err:        collection.ErrInsufficientTokens,
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
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}
