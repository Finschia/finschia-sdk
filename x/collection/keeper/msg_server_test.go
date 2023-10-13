package keeper_test

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Finschia/finschia-sdk/testutil"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/query"
	"github.com/Finschia/finschia-sdk/x/collection"
	"github.com/Finschia/finschia-sdk/x/token/class"
)

func (s *KeeperTestSuite) TestMsgSendFT() {
	testCases := map[string]struct {
		isNegativeCase bool
		req            *collection.MsgSendFT
		ftID           string
		expectedEvents sdk.Events
		expectedError  error
	}{
		"valid request": {
			req: &collection.MsgSendFT{
				ContractId: s.contractID,
				From:       s.vendor.String(),
				To:         s.customer.String(),
				Amount:     collection.NewCoins(collection.NewFTCoin(s.ftClassID, s.balance)),
			},
			ftID: collection.NewFTID(s.ftClassID),
			expectedEvents: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventSent",
					Attributes: []abci.EventAttribute{
						{Key: []byte("amount"), Value: testutil.MustJSONMarshal(collection.NewCoins(collection.NewFTCoin(s.ftClassID, s.balance))), Index: false},
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("to"), Value: testutil.W(s.customer.String()), Index: false},
					},
				}},
		},
		"contract not found": {
			isNegativeCase: true,
			req: &collection.MsgSendFT{
				ContractId: "deadbeef",
				From:       s.vendor.String(),
				To:         s.customer.String(),
				Amount:     collection.NewCoins(collection.NewFTCoin(s.ftClassID, s.balance)),
			},
			ftID:          collection.NewFTID(s.ftClassID),
			expectedError: class.ErrContractNotExist,
		},
		"insufficient funds": {
			isNegativeCase: true,
			req: &collection.MsgSendFT{
				ContractId: s.contractID,
				From:       s.vendor.String(),
				To:         s.customer.String(),
				Amount:     collection.NewCoins(collection.NewFTCoin(s.ftClassID, s.balance.Add(sdk.OneInt()))),
			},
			ftID:          collection.NewFTID(s.ftClassID),
			expectedError: collection.ErrInsufficientToken,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			// Arrange
			s.Require().NoError(tc.req.ValidateBasic())
			from, err := sdk.AccAddressFromBech32(tc.req.From)
			s.Require().NoError(err)
			to, err := sdk.AccAddressFromBech32(tc.req.To)
			s.Require().NoError(err)
			ctx, _ := s.ctx.CacheContext()
			prevFromBalance := s.keeper.GetBalance(ctx, tc.req.ContractId, from, tc.ftID)
			prevToBalance := s.keeper.GetBalance(ctx, tc.req.ContractId, to, tc.ftID)

			// Act
			res, err := s.msgServer.SendFT(sdk.WrapSDKContext(ctx), tc.req)
			if tc.isNegativeCase {
				s.Require().ErrorIs(err, tc.expectedError)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)

			// Assert
			events := ctx.EventManager().Events()
			s.Require().Equal(tc.expectedEvents, events)
			curFromBalance := s.keeper.GetBalance(ctx, tc.req.ContractId, from, tc.ftID)
			curToBalance := s.keeper.GetBalance(ctx, tc.req.ContractId, to, tc.ftID)
			s.Require().Equal(prevFromBalance.Sub(tc.req.Amount[0].Amount).Abs(), curFromBalance.Abs())
			s.Require().Equal(prevToBalance.Add(tc.req.Amount[0].Amount), curToBalance)
		})
	}
}

func (s *KeeperTestSuite) TestMsgOperatorSendFT() {
	testCases := map[string]struct {
		isNegativeCase bool
		req            *collection.MsgOperatorSendFT
		ftID           string
		expectedEvents sdk.Events
		expectedError  error
	}{
		"valid request": {
			req: &collection.MsgOperatorSendFT{
				ContractId: s.contractID,
				Operator:   s.operator.String(),
				From:       s.customer.String(),
				To:         s.vendor.String(),
				Amount:     collection.NewCoins(collection.NewFTCoin(s.ftClassID, s.balance)),
			},
			ftID: collection.NewFTID(s.ftClassID),
			expectedEvents: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventSent",
					Attributes: []abci.EventAttribute{
						{Key: []byte("amount"), Value: testutil.MustJSONMarshal(collection.NewCoins(collection.NewFTCoin(s.ftClassID, s.balance))), Index: false},
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(s.customer.String()), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.operator.String()), Index: false},
						{Key: []byte("to"), Value: testutil.W(s.vendor.String()), Index: false},
					},
				},
			},
		},
		"contract not found": {
			isNegativeCase: true,
			req: &collection.MsgOperatorSendFT{
				ContractId: "deadbeef",
				Operator:   s.operator.String(),
				From:       s.customer.String(),
				To:         s.vendor.String(),
				Amount:     collection.NewCoins(collection.NewFTCoin(s.ftClassID, s.balance)),
			},
			expectedError: class.ErrContractNotExist,
		},
		"not approved": {
			isNegativeCase: true,
			req: &collection.MsgOperatorSendFT{
				ContractId: s.contractID,
				Operator:   s.vendor.String(),
				From:       s.customer.String(),
				To:         s.vendor.String(),
				Amount:     collection.NewCoins(collection.NewFTCoin(s.ftClassID, s.balance)),
			},
			expectedError: collection.ErrCollectionNotApproved,
		},
		"insufficient funds": {
			isNegativeCase: true,
			req: &collection.MsgOperatorSendFT{
				ContractId: s.contractID,
				Operator:   s.operator.String(),
				From:       s.customer.String(),
				To:         s.vendor.String(),
				Amount:     collection.NewCoins(collection.NewFTCoin(s.ftClassID, s.balance.Add(sdk.OneInt()))),
			},
			expectedError: collection.ErrInsufficientToken,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			// Arrange
			s.Require().NoError(tc.req.ValidateBasic())
			from, err := sdk.AccAddressFromBech32(tc.req.From)
			s.Require().NoError(err)
			to, err := sdk.AccAddressFromBech32(tc.req.To)
			s.Require().NoError(err)
			operator, err := sdk.AccAddressFromBech32(tc.req.Operator)
			s.Require().NoError(err)
			ctx, _ := s.ctx.CacheContext()
			prevFromBalance := s.keeper.GetBalance(ctx, tc.req.ContractId, from, tc.ftID)
			prevToBalance := s.keeper.GetBalance(ctx, tc.req.ContractId, to, tc.ftID)
			prevOperatorBalance := s.keeper.GetBalance(ctx, tc.req.ContractId, operator, tc.ftID)

			// Act
			res, err := s.msgServer.OperatorSendFT(sdk.WrapSDKContext(ctx), tc.req)
			if tc.isNegativeCase {
				s.Require().ErrorIs(err, tc.expectedError)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)

			// Assert
			events := ctx.EventManager().Events()
			s.Require().Equal(tc.expectedEvents, events)
			curFromBalance := s.keeper.GetBalance(ctx, tc.req.ContractId, from, tc.ftID)
			curToBalance := s.keeper.GetBalance(ctx, tc.req.ContractId, to, tc.ftID)
			curOperatorBalance := s.keeper.GetBalance(ctx, tc.req.ContractId, operator, tc.ftID)
			s.Require().Equal(prevFromBalance.Sub(tc.req.Amount[0].Amount).Abs(), curFromBalance.Abs())
			s.Require().Equal(prevToBalance.Add(tc.req.Amount[0].Amount), curToBalance)
			s.Require().Equal(prevOperatorBalance, curOperatorBalance)
		})
	}
}

func (s *KeeperTestSuite) TestMsgSendNFT() {
	allTokenIDs := make([]string, 0)
	allTokenIDs = append(allTokenIDs, collection.NewNFTID(s.nftClassID, 1))
	cursor := allTokenIDs[0]
	for {
		ctx, _ := s.ctx.CacheContext()
		res, err := s.queryServer.Children(sdk.WrapSDKContext(ctx), &collection.QueryChildrenRequest{
			s.contractID,
			cursor,
			&query.PageRequest{},
		})
		s.Require().NoError(err)
		if res.Children == nil {
			break
		}
		allTokenIDs = append(allTokenIDs, res.Children[0].TokenId)
		cursor = allTokenIDs[len(allTokenIDs)-1]
	}

	testCases := map[string]struct {
		contractID string
		tokenID    string
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			tokenID:    collection.NewNFTID(s.nftClassID, 1),
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventOwnerChanged",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(s.customer.String()), Index: false},
						{Key: []byte("to"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("token_id"), Value: testutil.W(allTokenIDs[1]), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventOwnerChanged",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(s.customer.String()), Index: false},
						{Key: []byte("to"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("token_id"), Value: testutil.W(allTokenIDs[2]), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventOwnerChanged",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(s.customer.String()), Index: false},
						{Key: []byte("to"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("token_id"), Value: testutil.W(allTokenIDs[3]), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventSent",
					Attributes: []abci.EventAttribute{
						{Key: []byte("amount"), Value: testutil.MustJSONMarshal(collection.NewCoins(collection.Coin{TokenId: allTokenIDs[0], Amount: sdk.OneInt()})), Index: false},
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(s.customer.String()), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.customer.String()), Index: false},
						{Key: []byte("to"), Value: testutil.W(s.vendor.String()), Index: false},
					},
				},
			},
		},
		"contract not found": {
			contractID: "deadbeef",
			tokenID:    collection.NewNFTID(s.nftClassID, 1),
			err:        class.ErrContractNotExist,
		},
		"not found": {
			contractID: s.contractID,
			tokenID:    collection.NewNFTID("deadbeef", 1),
			err:        collection.ErrTokenNotExist,
		},
		"child": {
			contractID: s.contractID,
			tokenID:    collection.NewNFTID(s.nftClassID, 2),
			err:        collection.ErrTokenCannotTransferChildToken,
		},
		"not owned by": {
			contractID: s.contractID,
			tokenID:    collection.NewNFTID(s.nftClassID, s.numNFTs+1),
			err:        collection.ErrTokenNotOwnedBy,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgSendNFT{
				ContractId: tc.contractID,
				From:       s.customer.String(),
				To:         s.vendor.String(),
				TokenIds:   []string{tc.tokenID},
			}
			res, err := s.msgServer.SendNFT(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
			s.Require().Equal(tc.events, ctx.EventManager().Events())
		})
	}
}

func (s *KeeperTestSuite) TestMsgOperatorSendNFT() {
	tokenID := collection.NewNFTID(s.nftClassID, 1)
	allTokenIDs := make([]string, 0)
	allTokenIDs = append(allTokenIDs, tokenID)
	cursor := allTokenIDs[0]
	for {
		ctx, _ := s.ctx.CacheContext()
		res, err := s.queryServer.Children(sdk.WrapSDKContext(ctx), &collection.QueryChildrenRequest{
			ContractId: s.contractID,
			TokenId:    cursor,
			Pagination: &query.PageRequest{},
		})
		s.Require().NoError(err)
		if res.Children == nil {
			break
		}
		allTokenIDs = append(allTokenIDs, res.Children[0].TokenId)
		cursor = allTokenIDs[len(allTokenIDs)-1]
	}

	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		from       sdk.AccAddress
		tokenID    string
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			tokenID:    tokenID,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventOwnerChanged",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(s.customer.String()), Index: false},
						{Key: []byte("to"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("token_id"), Value: testutil.W(allTokenIDs[1]), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventOwnerChanged",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(s.customer.String()), Index: false},
						{Key: []byte("to"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("token_id"), Value: testutil.W(allTokenIDs[2]), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventOwnerChanged",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(s.customer.String()), Index: false},
						{Key: []byte("to"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("token_id"), Value: testutil.W(allTokenIDs[3]), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventSent",
					Attributes: []abci.EventAttribute{
						{Key: []byte("amount"), Value: testutil.MustJSONMarshal(collection.NewCoins(collection.Coin{TokenId: allTokenIDs[0], Amount: sdk.OneInt()})), Index: false},
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(s.customer.String()), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.operator.String()), Index: false},
						{Key: []byte("to"), Value: testutil.W(s.vendor.String()), Index: false},
					},
				},
			}},
		"contract not found": {
			contractID: "deadbeef",
			operator:   s.operator,
			from:       s.customer,
			tokenID:    tokenID,
			err:        class.ErrContractNotExist,
		},
		"not approved": {
			contractID: s.contractID,
			operator:   s.vendor,
			from:       s.customer,
			tokenID:    tokenID,
			err:        collection.ErrCollectionNotApproved,
		},
		"not found": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			tokenID:    collection.NewNFTID("deadbeef", 1),
			err:        collection.ErrTokenNotExist,
		},
		"child": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			tokenID:    collection.NewNFTID(s.nftClassID, 2),
			err:        collection.ErrTokenCannotTransferChildToken,
		},
		"not owned by": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			tokenID:    collection.NewNFTID(s.nftClassID, s.numNFTs+1),
			err:        collection.ErrTokenNotOwnedBy,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgOperatorSendNFT{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				From:       tc.from.String(),
				To:         s.vendor.String(),
				TokenIds:   []string{tc.tokenID},
			}
			res, err := s.msgServer.OperatorSendNFT(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
			s.Require().Equal(tc.events, ctx.EventManager().Events())
		})
	}
}

func (s *KeeperTestSuite) TestMsgAuthorizeOperator() {
	testCases := map[string]struct {
		isNegativeCase bool
		req            *collection.MsgAuthorizeOperator
		events         sdk.Events
		expectedError  error
	}{
		"valid request": {
			req: &collection.MsgAuthorizeOperator{
				ContractId: s.contractID,
				Holder:     s.customer.String(),
				Operator:   s.vendor.String(),
			},
			events: sdk.Events{sdk.Event{
				Type: "lbm.collection.v1.EventAuthorizedOperator",
				Attributes: []abci.EventAttribute{
					{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
					{Key: []byte("holder"), Value: testutil.W(s.customer.String()), Index: false},
					{Key: []byte("operator"), Value: testutil.W(s.vendor.String()), Index: false},
				},
			}},
		},
		"contract not found": {
			isNegativeCase: true,
			req: &collection.MsgAuthorizeOperator{
				ContractId: "deadbeef",
				Holder:     s.customer.String(),
				Operator:   s.vendor.String(),
			},
			expectedError: class.ErrContractNotExist,
		},
		"already approved": {
			isNegativeCase: true,
			req: &collection.MsgAuthorizeOperator{
				ContractId: s.contractID,
				Holder:     s.customer.String(),
				Operator:   s.operator.String(),
			},
			expectedError: collection.ErrCollectionAlreadyApproved,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			// Arrange
			s.Require().NoError(tc.req.ValidateBasic())
			holder, err := sdk.AccAddressFromBech32(tc.req.Holder)
			s.Require().NoError(err)
			operator, err := sdk.AccAddressFromBech32(tc.req.Operator)
			s.Require().NoError(err)
			ctx, _ := s.ctx.CacheContext()
			prevAuth, _ := s.keeper.GetAuthorization(ctx, tc.req.ContractId, holder, operator)

			// Act
			res, err := s.msgServer.AuthorizeOperator(sdk.WrapSDKContext(ctx), tc.req)
			if tc.isNegativeCase {
				s.Require().ErrorIs(err, tc.expectedError)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
			s.Require().Equal(tc.events, ctx.EventManager().Events())
			curAuth, err := s.keeper.GetAuthorization(ctx, tc.req.ContractId, holder, operator)
			s.Require().NoError(err)
			s.Require().Nil(prevAuth)
			s.Require().Equal(tc.req.Holder, curAuth.Holder)
			s.Require().Equal(tc.req.Operator, curAuth.Operator)
		})
	}
}

func (s *KeeperTestSuite) TestMsgRevokeOperator() {
	testCases := map[string]struct {
		isNegativeCase bool
		req            *collection.MsgRevokeOperator
		events         sdk.Events
		expectedError  error
	}{
		"valid request": {
			req: &collection.MsgRevokeOperator{
				ContractId: s.contractID,
				Holder:     s.customer.String(),
				Operator:   s.operator.String(),
			},
			events: sdk.Events{sdk.Event{
				Type: "lbm.collection.v1.EventRevokedOperator",
				Attributes: []abci.EventAttribute{
					{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
					{Key: []byte("holder"), Value: testutil.W(s.customer.String()), Index: false},
					{Key: []byte("operator"), Value: testutil.W(s.operator.String()), Index: false},
				},
			}},
		},
		"contract not found": {
			isNegativeCase: true,
			req: &collection.MsgRevokeOperator{
				ContractId: "deadbeef",
				Holder:     s.customer.String(),
				Operator:   s.operator.String(),
			},
			expectedError: class.ErrContractNotExist,
		},
		"no authorization": {
			isNegativeCase: true,
			req: &collection.MsgRevokeOperator{
				ContractId: s.contractID,
				Holder:     s.customer.String(),
				Operator:   s.vendor.String(),
			},
			expectedError: collection.ErrCollectionNotApproved,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			// Arrange
			s.Require().NoError(tc.req.ValidateBasic())
			holder, err := sdk.AccAddressFromBech32(tc.req.Holder)
			s.Require().NoError(err)
			operator, err := sdk.AccAddressFromBech32(tc.req.Operator)
			s.Require().NoError(err)
			ctx, _ := s.ctx.CacheContext()
			prevAuth, _ := s.keeper.GetAuthorization(ctx, tc.req.ContractId, holder, operator)

			// Act
			res, err := s.msgServer.RevokeOperator(sdk.WrapSDKContext(ctx), tc.req)
			if tc.isNegativeCase {
				s.Require().ErrorIs(err, tc.expectedError)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)

			s.Require().Equal(tc.events, ctx.EventManager().Events())
			s.Require().NotNil(prevAuth)
			s.Require().Equal(tc.req.Holder, prevAuth.Holder)
			s.Require().Equal(tc.req.Operator, prevAuth.Operator)
			curAuth, err := s.keeper.GetAuthorization(ctx, tc.req.ContractId, holder, operator)
			s.Require().ErrorIs(err, collection.ErrCollectionNotApproved)
			s.Require().Nil(curAuth)
		})
	}
}

func (s *KeeperTestSuite) TestMsgCreateContract() {
	newContractID := "3336b76f"
	testCases := map[string]struct {
		owner  sdk.AccAddress
		err    error
		events sdk.Events
	}{
		"valid request": {
			owner: s.vendor,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventCreatedContract",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(newContractID), Index: false},
						{Key: []byte("creator"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("meta"), Value: testutil.W(""), Index: false},
						{Key: []byte("name"), Value: testutil.W(""), Index: false},
						{Key: []byte("uri"), Value: testutil.W(""), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(newContractID), Index: false},
						{Key: []byte("grantee"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("granter"), Value: testutil.W(""), Index: false},
						{Key: []byte("permission"), Value: testutil.W(collection.Permission(collection.LegacyPermissionIssue).String()), Index: false},
					}},
				sdk.Event{
					Type: "lbm.collection.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(newContractID), Index: false},
						{Key: []byte("grantee"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("granter"), Value: testutil.W(""), Index: false},
						{Key: []byte("permission"), Value: testutil.W(collection.Permission(collection.LegacyPermissionModify).String()), Index: false},
					}},
				sdk.Event{
					Type: "lbm.collection.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(newContractID), Index: false},
						{Key: []byte("grantee"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("granter"), Value: testutil.W(""), Index: false},
						{Key: []byte("permission"), Value: testutil.W(collection.Permission(collection.LegacyPermissionMint).String()), Index: false},
					}},
				sdk.Event{
					Type: "lbm.collection.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(newContractID), Index: false},
						{Key: []byte("grantee"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("granter"), Value: testutil.W(""), Index: false},
						{Key: []byte("permission"), Value: testutil.W(collection.Permission(collection.LegacyPermissionBurn).String()), Index: false},
					}},
			}},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgCreateContract{
				Owner: tc.owner.String(),
			}
			res, err := s.msgServer.CreateContract(sdk.WrapSDKContext(ctx), req)
			s.Require().Equal(newContractID, res.ContractId)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
			s.Require().Equal(tc.events, ctx.EventManager().Events())

		})
	}
}

func (s *KeeperTestSuite) TestMsgIssueFT() {
	expectedClassID := "00000002"
	expectedTokenID := collection.NewFTID(expectedClassID)

	testCases := map[string]struct {
		contractID string
		owner      sdk.AccAddress
		amount     sdk.Int
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			owner:      s.vendor,
			amount:     sdk.ZeroInt(),
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventCreatedFTClass",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("decimals"), Value: []byte("0"), Index: false},
						{Key: []byte("meta"), Value: testutil.W(""), Index: false},
						{Key: []byte("mintable"), Value: []byte("false"), Index: false},
						{Key: []byte("name"), Value: testutil.W(""), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("token_id"), Value: testutil.W(expectedTokenID), Index: false},
					},
				},
			},
		},
		"valid request with supply": {
			contractID: s.contractID,
			owner:      s.vendor,
			amount:     sdk.OneInt(),
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventCreatedFTClass",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("decimals"), Value: []byte("0"), Index: false},
						{Key: []byte("meta"), Value: testutil.W(""), Index: false},
						{Key: []byte("mintable"), Value: []byte("false"), Index: false},
						{Key: []byte("name"), Value: testutil.W(""), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("token_id"), Value: testutil.W(expectedTokenID), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventMintedFT",
					Attributes: []abci.EventAttribute{
						{Key: []byte("amount"), Value: testutil.MustJSONMarshal(collection.NewCoins(collection.Coin{TokenId: expectedTokenID, Amount: sdk.OneInt()})), Index: false},
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("to"), Value: testutil.W(s.customer.String()), Index: false},
					},
				},
			},
		},
		"contract not found": {
			contractID: "deadbeef",
			owner:      s.vendor,
			amount:     sdk.ZeroInt(),
			err:        class.ErrContractNotExist,
		},
		"no permission": {
			contractID: s.contractID,
			owner:      s.customer,
			amount:     sdk.ZeroInt(),
			err:        collection.ErrTokenNoPermission,
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
			s.Require().Equal(tc.events, ctx.EventManager().Events())

			// check balance and tokenId
			tokenId := collection.NewFTID(res.TokenId)
			bal, err := s.queryServer.Balance(sdk.WrapSDKContext(ctx), &collection.QueryBalanceRequest{
				ContractId: s.contractID,
				Address:    s.customer.String(),
				TokenId:    tokenId,
			})
			s.Require().NoError(err)
			expectedCoin := collection.Coin{
				TokenId: tokenId,
				Amount:  tc.amount,
			}
			s.Require().Equal(expectedCoin, bal.Balance)
		})
	}
}

func (s *KeeperTestSuite) TestMsgIssueNFT() {
	expectedTokenType := "10000002"

	testCases := map[string]struct {
		contractID string
		owner      sdk.AccAddress
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			owner:      s.vendor,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventCreatedNFTClass",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("meta"), Value: testutil.W(""), Index: false},
						{Key: []byte("name"), Value: testutil.W(""), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("token_type"), Value: testutil.W(expectedTokenType), Index: false},
					}},
				sdk.Event{
					Type: "lbm.collection.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("grantee"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("granter"), Value: testutil.W(""), Index: false},
						{Key: []byte("permission"), Value: testutil.W(collection.Permission(collection.LegacyPermissionMint).String()), Index: false},
					}},
				sdk.Event{
					Type: "lbm.collection.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("grantee"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("granter"), Value: testutil.W(""), Index: false},
						{Key: []byte("permission"), Value: testutil.W(collection.Permission(collection.LegacyPermissionBurn).String()), Index: false},
					}},
			},
		},
		"contract not found": {
			contractID: "deadbeef",
			owner:      s.vendor,
			err:        class.ErrContractNotExist,
		},
		"no permission": {
			contractID: s.contractID,
			owner:      s.customer,
			err:        collection.ErrTokenNoPermission,
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
			s.Require().Equal(tc.events, ctx.EventManager().Events())

		})
	}
}

func (s *KeeperTestSuite) TestMsgMintFT() {
	// prepare multi tokens for test
	// create a fungible token class (mintable true)
	mintableFTClassID, err := s.keeper.CreateTokenClass(s.ctx, s.contractID, &collection.FTClass{
		Name:     "tibetian fox2",
		Mintable: true,
	})
	s.Require().NoError(err)

	// create a fungible token class (mintable false)
	nonmintableFTClassID, err := s.keeper.CreateTokenClass(s.ctx, s.contractID, &collection.FTClass{
		Name:     "tibetian fox3",
		Mintable: false,
	})
	s.Require().NoError(err)

	amount := collection.NewCoins(
		collection.NewFTCoin(s.ftClassID, sdk.NewInt(100000)),
	)
	amounts := collection.NewCoins(
		collection.NewFTCoin(s.ftClassID, sdk.NewInt(100000)),
		collection.NewFTCoin(*mintableFTClassID, sdk.NewInt(200000)),
	)

	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		amount     []collection.Coin
		err        error
		events     sdk.Events
	}{
		"valid request - single token": {
			contractID: s.contractID,
			from:       s.vendor,
			amount:     amount,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventMintedFT",
					Attributes: []abci.EventAttribute{
						{Key: []byte("amount"), Value: testutil.MustJSONMarshal(amount), Index: false},
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("to"), Value: testutil.W(s.customer.String()), Index: false},
					},
				},
			},
		},
		"valid request - multi tokens": {
			contractID: s.contractID,
			from:       s.vendor,
			amount:     amounts,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventMintedFT",
					Attributes: []abci.EventAttribute{
						{Key: []byte("amount"), Value: testutil.MustJSONMarshal(amounts), Index: false},
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("to"), Value: testutil.W(s.customer.String()), Index: false},
					},
				},
			},
		},
		"valid request - empty amount": {
			contractID: s.contractID,
			from:       s.vendor,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventMintedFT",
					Attributes: []abci.EventAttribute{
						{Key: []byte("amount"), Value: []byte("[]"), Index: false},
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("to"), Value: testutil.W(s.customer.String()), Index: false},
					},
				},
			},
		},
		"contract not found": {
			contractID: "deadbeef",
			from:       s.vendor,
			amount:     amount,
			err:        class.ErrContractNotExist,
		},
		"no permission": {
			contractID: s.contractID,
			from:       s.customer,
			amount:     amount,
			err:        collection.ErrTokenNoPermission,
		},
		"no class of the token": {
			contractID: s.contractID,
			from:       s.vendor,
			amount: collection.NewCoins(
				collection.NewFTCoin("00bab10c", sdk.OneInt()),
			),
			err: collection.ErrTokenNotExist,
		},
		"include invalid tokenId among 2 tokens": {
			contractID: s.contractID,
			from:       s.vendor,
			amount: collection.NewCoins(
				collection.NewFTCoin(s.ftClassID, sdk.OneInt()),
				collection.NewFTCoin("00bab10b", sdk.OneInt()), // no exist tokenId
			),
			err: collection.ErrTokenNotExist,
		},
		"mintable false tokenId": {
			contractID: s.contractID,
			from:       s.vendor,
			amount:     collection.NewCoins(collection.NewFTCoin(*nonmintableFTClassID, sdk.OneInt())),
			err:        collection.ErrTokenNotMintable,
		},
		"include mintable false among 2 tokens": {
			contractID: s.contractID,
			from:       s.vendor,
			amount: collection.NewCoins(
				collection.NewFTCoin(*mintableFTClassID, sdk.OneInt()),
				collection.NewFTCoin(*nonmintableFTClassID, sdk.OneInt()),
			),
			err: collection.ErrTokenNotMintable,
		},
	}

	// query the values to be effected by MintFT
	queryValuesEffectedByMintFT := func(ctx sdk.Context, coins collection.Coins, contractID string) (balances collection.Coins, supply []sdk.Int, minted []sdk.Int) {
		for _, am := range coins {
			// save balance
			bal, err := s.queryServer.Balance(sdk.WrapSDKContext(ctx), &collection.QueryBalanceRequest{
				ContractId: contractID,
				Address:    s.customer.String(),
				TokenId:    am.TokenId,
			})
			s.Require().NoError(err)
			balances = append(balances, bal.Balance)

			// save supply
			res, err := s.queryServer.FTSupply(sdk.WrapSDKContext(ctx), &collection.QueryFTSupplyRequest{
				ContractId: contractID,
				TokenId:    am.TokenId,
			})
			s.Require().NoError(err)
			supply = append(supply, res.Supply)

			// save minted
			m, err := s.queryServer.FTMinted(sdk.WrapSDKContext(ctx), &collection.QueryFTMintedRequest{
				ContractId: contractID,
				TokenId:    am.TokenId,
			})
			s.Require().NoError(err)
			minted = append(minted, m.Minted)
		}
		return
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			// test multiple times
			ctx := s.ctx
			for t := 0; t < 3; t++ {
				ctx, _ = ctx.CacheContext()

				prevAmount, prevSupply, prevMinted := queryValuesEffectedByMintFT(ctx, tc.amount, tc.contractID)

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
				s.Require().Equal(tc.events, ctx.EventManager().Events())

				// check results
				afterAmount, afterSupply, afterMinted := queryValuesEffectedByMintFT(ctx, tc.amount, tc.contractID)
				for i, am := range tc.amount {
					expectedBalance := collection.Coin{
						TokenId: am.TokenId,
						Amount:  prevAmount[i].Amount.Add(am.Amount),
					}
					s.Require().Equal(expectedBalance, afterAmount[i])

					expectedSupply := prevSupply[i].Add(am.Amount)
					s.Require().True(expectedSupply.Equal(afterSupply[i]))

					expectedMinted := prevMinted[i].Add(am.Amount)
					s.Require().True(expectedMinted.Equal(afterMinted[i]))
				}
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgMintNFT() {
	params := []collection.MintNFTParam{{
		TokenType: s.nftClassID,
		Name:      "tester",
		Meta:      "Mint NFT",
	}}
	expectedTokens := []collection.NFT{
		{
			TokenId: "1000000100000016",
			Name:    params[0].Name,
			Meta:    params[0].Meta,
		},
	}

	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		params     []collection.MintNFTParam
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			from:       s.vendor,
			params:     params,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventMintedNFT",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("to"), Value: testutil.W(s.customer.String()), Index: false},
						{Key: []byte("tokens"), Value: testutil.MustJSONMarshal(expectedTokens), Index: false},
					},
				},
			}},
		"contract not found": {
			contractID: "deadbeef",
			from:       s.vendor,
			params:     params,
			err:        class.ErrContractNotExist,
		},
		"no permission": {
			contractID: s.contractID,
			from:       s.customer,
			params:     params,
			err:        collection.ErrTokenNoPermission,
		},
		"no class of the token": {
			contractID: s.contractID,
			from:       s.vendor,
			params: []collection.MintNFTParam{{
				TokenType: "deadbeef",
			}},
			err: collection.ErrTokenTypeNotExist,
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
			s.Require().Equal(tc.events, ctx.EventManager().Events())
		})
	}
}

func (s *KeeperTestSuite) TestMsgBurnFT() {
	// prepare mutli token burn test
	amount := collection.NewCoins(
		collection.NewFTCoin(s.ftClassID, sdk.NewInt(50000)),
	)

	// create a fungible token class
	mintableFTClassID, err := s.keeper.CreateTokenClass(s.ctx, s.contractID, &collection.FTClass{
		Name:     "tibetian fox2",
		Mintable: true,
	})
	s.Require().NoError(err)
	amounts := collection.NewCoins(
		collection.NewFTCoin(s.ftClassID, sdk.NewInt(50000)),
		collection.NewFTCoin(*mintableFTClassID, sdk.NewInt(60000)),
	)

	// mintft
	mintedCoin := collection.NewFTCoin(*mintableFTClassID, sdk.NewInt(1000000))
	err = s.keeper.MintFT(s.ctx, s.contractID, s.vendor, []collection.Coin{mintedCoin})
	s.Require().NoError(err)

	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		amount     []collection.Coin
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			from:       s.vendor,
			amount:     amount,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventBurned",
					Attributes: []abci.EventAttribute{
						{Key: []byte("amount"), Value: testutil.MustJSONMarshal(collection.NewCoins(
							collection.NewFTCoin(s.ftClassID, sdk.NewInt(50000)),
						)), Index: false},
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.vendor.String()), Index: false},
					},
				},
			},
		},
		"valid multi amount burn": {
			contractID: s.contractID,
			from:       s.vendor,
			amount:     amounts,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventBurned",
					Attributes: []abci.EventAttribute{
						{Key: []byte("amount"), Value: testutil.MustJSONMarshal(amounts), Index: false},
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.vendor.String()), Index: false},
					},
				},
			},
		},
		"no amount - valid": {
			contractID: s.contractID,
			from:       s.vendor,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventBurned",
					Attributes: []abci.EventAttribute{
						{Key: []byte("amount"), Value: []byte("[]"), Index: false},
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.vendor.String()), Index: false},
					},
				},
			},
		},
		"contract not found": {
			contractID: "deadbeef",
			from:       s.vendor,
			amount:     amount,
			err:        class.ErrContractNotExist,
		},
		"no permission": {
			contractID: s.contractID,
			from:       s.customer,
			amount:     amount,
			err:        collection.ErrTokenNoPermission,
		},
		"insufficient funds": {
			contractID: s.contractID,
			from:       s.vendor,
			amount: collection.NewCoins(
				collection.NewFTCoin("00bab10c", sdk.OneInt()),
			),
			err: collection.ErrInsufficientToken,
		},
		"include insufficient funds amount 2 amounts": {
			contractID: s.contractID,
			from:       s.vendor,
			amount: collection.NewCoins(
				collection.NewFTCoin(s.ftClassID, s.balance),
				collection.NewFTCoin("00bab10c", sdk.OneInt()),
			),
			err: collection.ErrInsufficientToken,
		},
	}

	// query the values to be effected by BurnFT
	queryValuesAffectedByBurnFT := func(ctx sdk.Context, coins collection.Coins, contractID, from string) (balances collection.Coins, supply []sdk.Int, burnt []sdk.Int) {
		for _, am := range coins {
			// save balance
			bal, err := s.queryServer.Balance(sdk.WrapSDKContext(ctx), &collection.QueryBalanceRequest{
				ContractId: contractID,
				Address:    from,
				TokenId:    am.TokenId,
			})
			s.Require().NoError(err)
			balances = append(balances, bal.Balance)

			// save supply
			res, err := s.queryServer.FTSupply(sdk.WrapSDKContext(ctx), &collection.QueryFTSupplyRequest{
				ContractId: contractID,
				TokenId:    am.TokenId,
			})
			s.Require().NoError(err)
			supply = append(supply, res.Supply)

			// save minted
			b, err := s.queryServer.FTBurnt(sdk.WrapSDKContext(ctx), &collection.QueryFTBurntRequest{
				ContractId: contractID,
				TokenId:    am.TokenId,
			})
			s.Require().NoError(err)
			burnt = append(burnt, b.Burnt)
		}
		return
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			// test multiple times
			ctx := s.ctx
			for t := 0; t < 3; t++ {
				ctx, _ = ctx.CacheContext()
				prevAmount, prevSupply, prevBurnt := queryValuesAffectedByBurnFT(ctx, tc.amount, tc.contractID, tc.from.String())

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
				s.Require().Equal(tc.events, ctx.EventManager().Events())

				// check changed amount
				afterAmount, afterSupply, afterBurnt := queryValuesAffectedByBurnFT(ctx, tc.amount, tc.contractID, tc.from.String())
				for i, am := range tc.amount {
					expectedBalance := prevAmount[i].Amount.Sub(am.Amount)
					s.Require().Equal(am.TokenId, afterAmount[i].TokenId)
					s.Require().True(expectedBalance.Equal(afterAmount[i].Amount))

					expectedSupply := prevSupply[i].Sub(am.Amount)
					s.Require().True(expectedSupply.Equal(afterSupply[i]))

					expectedBurnt := prevBurnt[i].Add(am.Amount)
					s.Require().True(expectedBurnt.Equal(afterBurnt[i]))
				}
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgOperatorBurnFT() {
	amount := collection.NewCoins(
		collection.NewFTCoin(s.ftClassID, s.balance),
	)

	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		from       sdk.AccAddress
		amount     []collection.Coin
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			amount:     amount,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventBurned",
					Attributes: []abci.EventAttribute{
						{Key: []byte("amount"), Value: testutil.MustJSONMarshal(amount), Index: false},
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(s.customer.String()), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.operator.String()), Index: false},
					}}},
		},
		"contract not found": {
			contractID: "deadbeef",
			operator:   s.operator,
			from:       s.customer,
			amount:     amount,
			err:        class.ErrContractNotExist,
		},
		"no authorization": {
			contractID: s.contractID,
			operator:   s.vendor,
			from:       s.customer,
			amount:     amount,
			err:        collection.ErrCollectionNotApproved,
		},
		"no permission": {
			contractID: s.contractID,
			operator:   s.stranger,
			from:       s.customer,
			amount:     amount,
			err:        collection.ErrTokenNoPermission,
		},
		"insufficient funds - exist token": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			amount: collection.NewCoins(
				collection.NewFTCoin(s.ftClassID, s.balance.Add(sdk.OneInt())),
			),
			err: collection.ErrInsufficientToken,
		},
		"insufficient funds - non-exist token": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			amount: collection.NewCoins(
				collection.NewFTCoin("00bab10c", sdk.OneInt()),
			),
			err: collection.ErrInsufficientToken,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgOperatorBurnFT{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				From:       tc.from.String(),
				Amount:     tc.amount,
			}
			res, err := s.msgServer.OperatorBurnFT(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
			s.Require().Equal(tc.events, ctx.EventManager().Events())
		})
	}
}

func (s *KeeperTestSuite) TestMsgBurnNFT() {
	tokenID := collection.NewNFTID(s.nftClassID, s.numNFTs*2+1)
	target := []string{tokenID}
	tokenIDs := make([]string, 0)
	tokenIDs = append(tokenIDs, tokenID)
	cursor := tokenIDs[0]
	for {
		ctx, _ := s.ctx.CacheContext()
		res, err := s.queryServer.Children(sdk.WrapSDKContext(ctx), &collection.QueryChildrenRequest{
			ContractId: s.contractID,
			TokenId:    cursor,
			Pagination: &query.PageRequest{},
		})
		s.Require().NoError(err)
		if res.Children == nil {
			break
		}
		tokenIDs = append(tokenIDs, res.Children[0].TokenId)
		cursor = tokenIDs[len(tokenIDs)-1]
	}
	coins := make([]collection.Coin, 0)
	for _, id := range tokenIDs {
		coins = append(coins, collection.NewCoin(id, sdk.NewInt(1)))
	}

	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		tokenIDs   []string
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			from:       s.vendor,
			tokenIDs:   target,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventBurned",
					Attributes: []abci.EventAttribute{
						{Key: []byte("amount"), Value: testutil.MustJSONMarshal(coins), Index: false},
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.vendor.String()), Index: false},
					}}},
		},
		"contract not found": {
			contractID: "deadbeef",
			from:       s.vendor,
			tokenIDs:   target,
			err:        class.ErrContractNotExist,
		},
		"no permission": {
			contractID: s.contractID,
			from:       s.customer,
			tokenIDs:   target,
			err:        collection.ErrTokenNoPermission,
		},
		"not found": {
			contractID: s.contractID,
			from:       s.vendor,
			tokenIDs: []string{
				collection.NewNFTID("deadbeef", 1),
			},
			err: collection.ErrTokenNotExist,
		},
		"child": {
			contractID: s.contractID,
			from:       s.vendor,
			tokenIDs: []string{
				collection.NewNFTID(s.nftClassID, 2),
			},
			err: collection.ErrBurnNonRootNFT,
		},
		"not owned by": {
			contractID: s.contractID,
			from:       s.vendor,
			tokenIDs: []string{
				collection.NewNFTID(s.nftClassID, s.numNFTs+1),
			},
			err: collection.ErrTokenNotOwnedBy,
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
			s.Require().Equal(tc.events, ctx.EventManager().Events())
		})
	}
}

func (s *KeeperTestSuite) TestMsgOperatorBurnNFT() {
	tokenID := collection.NewNFTID(s.nftClassID, 1)
	target := []string{tokenID}
	tokenIDs := make([]string, 0)
	tokenIDs = append(tokenIDs, tokenID)
	cursor := tokenIDs[0]
	for {
		ctx, _ := s.ctx.CacheContext()
		res, err := s.queryServer.Children(sdk.WrapSDKContext(ctx), &collection.QueryChildrenRequest{
			ContractId: s.contractID,
			TokenId:    cursor,
			Pagination: &query.PageRequest{},
		})
		s.Require().NoError(err)
		if res.Children == nil {
			break
		}
		tokenIDs = append(tokenIDs, res.Children[0].TokenId)
		cursor = tokenIDs[len(tokenIDs)-1]
	}
	coins := make([]collection.Coin, 0)
	for _, id := range tokenIDs {
		coins = append(coins, collection.NewCoin(id, sdk.NewInt(1)))
	}

	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		from       sdk.AccAddress
		tokenIDs   []string
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			tokenIDs:   target,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventBurned",
					Attributes: []abci.EventAttribute{
						{Key: []byte("amount"), Value: testutil.MustJSONMarshal(coins), Index: false},
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(s.customer.String()), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.operator.String()), Index: false},
					}}},
		},
		"contract not found": {
			contractID: "deadbeef",
			operator:   s.operator,
			from:       s.customer,
			tokenIDs:   target,
			err:        class.ErrContractNotExist,
		},
		"no authorization": {
			contractID: s.contractID,
			operator:   s.vendor,
			from:       s.customer,
			tokenIDs:   target,
			err:        collection.ErrCollectionNotApproved,
		},
		"no permission": {
			contractID: s.contractID,
			operator:   s.stranger,
			from:       s.customer,
			tokenIDs:   target,
			err:        collection.ErrTokenNoPermission,
		},
		"not found": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			tokenIDs: []string{
				collection.NewNFTID("deadbeef", 1),
			},
			err: collection.ErrTokenNotExist,
		},
		"child": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			tokenIDs: []string{
				collection.NewNFTID(s.nftClassID, 2),
			},
			err: collection.ErrBurnNonRootNFT,
		},
		"not owned by": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			tokenIDs: []string{
				collection.NewNFTID(s.nftClassID, s.numNFTs+1),
			},
			err: collection.ErrTokenNotOwnedBy,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgOperatorBurnNFT{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				From:       tc.from.String(),
				TokenIds:   tc.tokenIDs,
			}
			res, err := s.msgServer.OperatorBurnNFT(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
			s.Require().Equal(tc.events, ctx.EventManager().Events())

		})
	}
}

func (s *KeeperTestSuite) TestMsgModify() {
	tokenIndex := collection.NewNFTID(s.nftClassID, 1)[8:]
	changes := []collection.Attribute{{
		Key:   collection.AttributeKeyName.String(),
		Value: "test",
	}}

	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		tokenType  string
		tokenIndex string
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.vendor,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventModifiedContract",
					Attributes: []abci.EventAttribute{
						{Key: []byte("changes"), Value: testutil.MustJSONMarshal(changes), Index: false},
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.vendor.String()), Index: false},
					}}}},
		"contract not found": {
			contractID: "deadbeef",
			operator:   s.vendor,
			err:        class.ErrContractNotExist,
		},
		"no permission": {
			contractID: s.contractID,
			operator:   s.customer,
			tokenType:  s.nftClassID,
			tokenIndex: tokenIndex,
			err:        collection.ErrTokenNoPermission,
		},
		"nft not found": {
			contractID: s.contractID,
			operator:   s.vendor,
			tokenType:  s.nftClassID,
			tokenIndex: collection.NewNFTID(s.nftClassID, s.numNFTs*3+1)[8:],
			err:        collection.ErrTokenNotExist,
		},
		"ft class not found": {
			contractID: s.contractID,
			operator:   s.vendor,
			tokenType:  "00bab10c",
			tokenIndex: collection.NewFTID("00bab10c")[8:],
			err:        collection.ErrTokenNotExist,
		},
		"nft class not found": {
			contractID: s.contractID,
			operator:   s.vendor,
			tokenType:  "deadbeef",
			err:        collection.ErrTokenTypeNotExist,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
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
			s.Require().Equal(tc.events, ctx.EventManager().Events())
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
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			granter:    s.vendor,
			grantee:    s.operator,
			permission: collection.LegacyPermissionModify.String(),
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("grantee"), Value: testutil.W(s.operator.String()), Index: false},
						{Key: []byte("granter"), Value: testutil.W(s.vendor.String()), Index: false},
						{Key: []byte("permission"), Value: testutil.W(collection.Permission(collection.LegacyPermissionModify).String()), Index: false},
					}},
			},
		},
		"contract not found": {
			contractID: "deadbeef",
			granter:    s.vendor,
			grantee:    s.operator,
			permission: collection.LegacyPermissionModify.String(),
			err:        class.ErrContractNotExist,
		},
		"granter has no permission": {
			contractID: s.contractID,
			granter:    s.customer,
			grantee:    s.operator,
			permission: collection.LegacyPermissionModify.String(),
			err:        collection.ErrTokenNoPermission,
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
			s.Require().Equal(tc.events, ctx.EventManager().Events())
		})
	}
}

func (s *KeeperTestSuite) TestMsgRevokePermission() {
	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		permission string
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			from:       s.operator,
			permission: collection.LegacyPermissionMint.String(),
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventRenounced",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("grantee"), Value: testutil.W(s.operator.String()), Index: false},
						{Key: []byte("permission"), Value: testutil.W(collection.Permission(collection.LegacyPermissionMint).String()), Index: false},
					}},
			},
		},
		"contract not found": {
			contractID: "deadbeef",
			from:       s.operator,
			permission: collection.LegacyPermissionMint.String(),
			err:        class.ErrContractNotExist,
		},
		"not granted yet": {
			contractID: s.contractID,
			from:       s.operator,
			permission: collection.LegacyPermissionModify.String(),
			err:        collection.ErrTokenNoPermission,
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
			s.Require().Equal(tc.events, ctx.EventManager().Events())
		})
	}
}

func (s *KeeperTestSuite) TestMsgAttach() {
	testCases := map[string]struct {
		contractID string
		subjectID  string
		targetID   string
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			subjectID:  collection.NewNFTID(s.nftClassID, s.depthLimit+1),
			targetID:   collection.NewNFTID(s.nftClassID, 1),
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventAttached",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("holder"), Value: testutil.W(s.customer.String()), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.customer.String()), Index: false},
						{Key: []byte("subject"), Value: testutil.W(collection.NewNFTID(s.nftClassID, s.depthLimit+1)), Index: false},
						{Key: []byte("target"), Value: testutil.W(collection.NewNFTID(s.nftClassID, 1)), Index: false},
					}},
			},
		},
		"contract not found": {
			contractID: "deadbeef",
			subjectID:  collection.NewNFTID(s.nftClassID, collection.DefaultDepthLimit+1),
			targetID:   collection.NewNFTID(s.nftClassID, 1),
			err:        class.ErrContractNotExist,
		},
		"not owner of the token": {
			contractID: s.contractID,
			subjectID:  collection.NewNFTID(s.nftClassID, s.numNFTs+1),
			targetID:   collection.NewNFTID(s.nftClassID, 1),
			err:        collection.ErrTokenNotOwnedBy,
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
			s.Require().Equal(tc.events, ctx.EventManager().Events())
		})
	}
}

func (s *KeeperTestSuite) TestMsgDetach() {
	nfts := make([]string, s.depthLimit)
	for i := 1; i <= s.depthLimit; i++ {
		nfts[i-1] = collection.NewNFTID(s.nftClassID, i)
	}

	testCases := map[string]struct {
		contractID string
		subjectID  string
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			subjectID:  nfts[1],
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventDetached",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("holder"), Value: testutil.W(s.customer.String()), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.customer.String()), Index: false},
						{Key: []byte("previous_parent"), Value: testutil.W(nfts[0]), Index: false},
						{Key: []byte("subject"), Value: testutil.W(nfts[1]), Index: false},
					}},
				sdk.Event{
					Type: "lbm.collection.v1.EventRootChanged",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(nfts[0]), Index: false},
						{Key: []byte("to"), Value: testutil.W(nfts[1]), Index: false},
						{Key: []byte("token_id"), Value: testutil.W(nfts[2]), Index: false},
					}},
				sdk.Event{
					Type: "lbm.collection.v1.EventRootChanged",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(nfts[0]), Index: false},
						{Key: []byte("to"), Value: testutil.W(nfts[1]), Index: false},
						{Key: []byte("token_id"), Value: testutil.W(nfts[3]), Index: false},
					}},
			},
		},
		"contract not found": {
			contractID: "deadbeef",
			subjectID:  collection.NewNFTID(s.nftClassID, 2),
			err:        class.ErrContractNotExist,
		},
		"not owner of the token": {
			contractID: s.contractID,
			subjectID:  collection.NewNFTID(s.nftClassID, s.numNFTs+2),
			err:        collection.ErrTokenNotOwnedBy,
		},
		"not a child": {
			contractID: s.contractID,
			subjectID:  collection.NewNFTID(s.nftClassID, 1),
			err:        collection.ErrTokenNotAChild,
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
			s.Require().Equal(tc.events, ctx.EventManager().Events())

		})
	}
}

func (s *KeeperTestSuite) TestMsgOperatorAttach() {
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		subjectID  string
		targetID   string
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			subjectID:  collection.NewNFTID(s.nftClassID, s.depthLimit+1),
			targetID:   collection.NewNFTID(s.nftClassID, 1),
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventAttached",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("holder"), Value: testutil.W(s.customer.String()), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.operator.String()), Index: false},
						{Key: []byte("subject"), Value: testutil.W(collection.NewNFTID(s.nftClassID, s.depthLimit+1)), Index: false},
						{Key: []byte("target"), Value: testutil.W(collection.NewNFTID(s.nftClassID, 1)), Index: false},
					}},
			}},
		"contract not found": {
			contractID: "deadbeef",
			operator:   s.operator,
			subjectID:  collection.NewNFTID(s.nftClassID, collection.DefaultDepthLimit+1),
			targetID:   collection.NewNFTID(s.nftClassID, 1),
			err:        class.ErrContractNotExist,
		},
		"not authorized": {
			contractID: s.contractID,
			operator:   s.vendor,
			subjectID:  collection.NewNFTID(s.nftClassID, s.depthLimit+1),
			targetID:   collection.NewNFTID(s.nftClassID, 1),
			err:        collection.ErrCollectionNotApproved,
		},
		"not owner of the token": {
			contractID: s.contractID,
			operator:   s.operator,
			subjectID:  collection.NewNFTID(s.nftClassID, s.numNFTs+1),
			targetID:   collection.NewNFTID(s.nftClassID, 1),
			err:        collection.ErrTokenNotOwnedBy,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgOperatorAttach{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				From:       s.customer.String(),
				TokenId:    tc.subjectID,
				ToTokenId:  tc.targetID,
			}
			res, err := s.msgServer.OperatorAttach(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
			s.Require().Equal(tc.events, ctx.EventManager().Events())
		})
	}
}

func (s *KeeperTestSuite) TestMsgOperatorDetach() {
	nfts := make([]string, s.depthLimit)
	for i := 1; i <= s.depthLimit; i++ {
		nfts[i-1] = collection.NewNFTID(s.nftClassID, i)
	}

	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		subjectID  string
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			subjectID:  collection.NewNFTID(s.nftClassID, 2),
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventDetached",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("holder"), Value: testutil.W(s.customer.String()), Index: false},
						{Key: []byte("operator"), Value: testutil.W(s.operator.String()), Index: false},
						{Key: []byte("previous_parent"), Value: testutil.W(nfts[0]), Index: false},
						{Key: []byte("subject"), Value: testutil.W(nfts[1]), Index: false},
					}},
				sdk.Event{
					Type: "lbm.collection.v1.EventRootChanged",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(nfts[0]), Index: false},
						{Key: []byte("to"), Value: testutil.W(nfts[1]), Index: false},
						{Key: []byte("token_id"), Value: testutil.W(nfts[2]), Index: false},
					}},
				sdk.Event{
					Type: "lbm.collection.v1.EventRootChanged",
					Attributes: []abci.EventAttribute{
						{Key: []byte("contract_id"), Value: testutil.W(s.contractID), Index: false},
						{Key: []byte("from"), Value: testutil.W(nfts[0]), Index: false},
						{Key: []byte("to"), Value: testutil.W(nfts[1]), Index: false},
						{Key: []byte("token_id"), Value: testutil.W(nfts[3]), Index: false},
					}},
			}},
		"contract not found": {
			contractID: "deadbeef",
			operator:   s.operator,
			subjectID:  collection.NewNFTID(s.nftClassID, 2),
			err:        class.ErrContractNotExist,
		},
		"not authorized": {
			contractID: s.contractID,
			operator:   s.vendor,
			subjectID:  collection.NewNFTID(s.nftClassID, 2),
			err:        collection.ErrCollectionNotApproved,
		},
		"not owner of the token": {
			contractID: s.contractID,
			operator:   s.operator,
			subjectID:  collection.NewNFTID(s.nftClassID, s.numNFTs+2),
			err:        collection.ErrTokenNotOwnedBy,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgOperatorDetach{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				From:       s.customer.String(),
				TokenId:    tc.subjectID,
			}
			res, err := s.msgServer.OperatorDetach(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
			s.Require().Equal(tc.events, ctx.EventManager().Events())

		})
	}
}
