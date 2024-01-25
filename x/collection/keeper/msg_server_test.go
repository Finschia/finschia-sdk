package keeper_test

import (
	"encoding/json"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/golang/mock/gomock"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/Finschia/finschia-sdk/x/collection"
	"github.com/Finschia/finschia-sdk/x/collection-token/class"
)

func mustJSONMarshal(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(b)
}

// w wraps input with double quotes if it is a string or fmt.Stringer.
func w(input any) string {
	switch input.(type) {
	case string, fmt.Stringer:
		return fmt.Sprintf("\"%s\"", input)
	default:
		panic("unsupported type")
	}
}

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
						{Key: "amount", Value: mustJSONMarshal(collection.NewCoins(collection.NewFTCoin(s.ftClassID, s.balance))), Index: false},
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "from", Value: w(s.vendor.String()), Index: false},
						{Key: "operator", Value: w(s.vendor.String()), Index: false},
						{Key: "to", Value: w(s.customer.String()), Index: false},
					},
				},
			},
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
				Amount:     collection.NewCoins(collection.NewFTCoin(s.ftClassID, s.balance.Add(math.OneInt()))),
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
			res, err := s.msgServer.SendFT(ctx, tc.req)
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
						{Key: "amount", Value: mustJSONMarshal(collection.NewCoins(collection.NewFTCoin(s.ftClassID, s.balance))), Index: false},
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "from", Value: w(s.customer.String()), Index: false},
						{Key: "operator", Value: w(s.operator.String()), Index: false},
						{Key: "to", Value: w(s.vendor.String()), Index: false},
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
				Amount:     collection.NewCoins(collection.NewFTCoin(s.ftClassID, s.balance.Add(math.OneInt()))),
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
			res, err := s.msgServer.OperatorSendFT(ctx, tc.req)
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
	rootNFTID := collection.NewNFTID(s.nftClassID, 1)
	issuedTokenIDs := s.extractChainedNFTIDs(rootNFTID)

	testCases := map[string]struct {
		contractID string
		tokenID    string
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			tokenID:    rootNFTID,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventOwnerChanged",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "from", Value: w(s.customer.String()), Index: false},
						{Key: "to", Value: w(s.vendor.String()), Index: false},
						{Key: "token_id", Value: w(issuedTokenIDs[1]), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventOwnerChanged",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "from", Value: w(s.customer.String()), Index: false},
						{Key: "to", Value: w(s.vendor.String()), Index: false},
						{Key: "token_id", Value: w(issuedTokenIDs[2]), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventOwnerChanged",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "from", Value: w(s.customer.String()), Index: false},
						{Key: "to", Value: w(s.vendor.String()), Index: false},
						{Key: "token_id", Value: w(issuedTokenIDs[3]), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventSent",
					Attributes: []abci.EventAttribute{
						{Key: "amount", Value: mustJSONMarshal(collection.NewCoins(collection.Coin{TokenId: issuedTokenIDs[0], Amount: math.OneInt()})), Index: false},
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "from", Value: w(s.customer.String()), Index: false},
						{Key: "operator", Value: w(s.customer.String()), Index: false},
						{Key: "to", Value: w(s.vendor.String()), Index: false},
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
			res, err := s.msgServer.SendNFT(ctx, req)
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
	rootNFTID := collection.NewNFTID(s.nftClassID, 1)
	issuedTokenIDs := s.extractChainedNFTIDs(rootNFTID)

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
			tokenID:    rootNFTID,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventOwnerChanged",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "from", Value: w(s.customer.String()), Index: false},
						{Key: "to", Value: w(s.vendor.String()), Index: false},
						{Key: "token_id", Value: w(issuedTokenIDs[1]), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventOwnerChanged",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "from", Value: w(s.customer.String()), Index: false},
						{Key: "to", Value: w(s.vendor.String()), Index: false},
						{Key: "token_id", Value: w(issuedTokenIDs[2]), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventOwnerChanged",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "from", Value: w(s.customer.String()), Index: false},
						{Key: "to", Value: w(s.vendor.String()), Index: false},
						{Key: "token_id", Value: w(issuedTokenIDs[3]), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventSent",
					Attributes: []abci.EventAttribute{
						{Key: "amount", Value: mustJSONMarshal(collection.NewCoins(collection.Coin{TokenId: issuedTokenIDs[0], Amount: math.OneInt()})), Index: false},
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "from", Value: w(s.customer.String()), Index: false},
						{Key: "operator", Value: w(s.operator.String()), Index: false},
						{Key: "to", Value: w(s.vendor.String()), Index: false},
					},
				},
			},
		},
		"contract not found": {
			contractID: "deadbeef",
			operator:   s.operator,
			from:       s.customer,
			tokenID:    rootNFTID,
			err:        class.ErrContractNotExist,
		},
		"not approved": {
			contractID: s.contractID,
			operator:   s.vendor,
			from:       s.customer,
			tokenID:    rootNFTID,
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
			res, err := s.msgServer.OperatorSendNFT(ctx, req)
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
					{Key: "contract_id", Value: w(s.contractID), Index: false},
					{Key: "holder", Value: w(s.customer.String()), Index: false},
					{Key: "operator", Value: w(s.vendor.String()), Index: false},
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
			res, err := s.msgServer.AuthorizeOperator(ctx, tc.req)
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
					{Key: "contract_id", Value: w(s.contractID), Index: false},
					{Key: "holder", Value: w(s.customer.String()), Index: false},
					{Key: "operator", Value: w(s.operator.String()), Index: false},
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
			res, err := s.msgServer.RevokeOperator(ctx, tc.req)
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
	expectedNewContractID := "3336b76f"
	s.mockClassKeeper.EXPECT().NewID(gomock.Any()).Return(expectedNewContractID).Times(1)
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
						{Key: "contract_id", Value: w(expectedNewContractID), Index: false},
						{Key: "creator", Value: w(s.vendor.String()), Index: false},
						{Key: "meta", Value: w(""), Index: false},
						{Key: "name", Value: w(""), Index: false},
						{Key: "uri", Value: w(""), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(expectedNewContractID), Index: false},
						{Key: "grantee", Value: w(s.vendor.String()), Index: false},
						{Key: "granter", Value: w(""), Index: false},
						{Key: "permission", Value: w(collection.Permission(collection.LegacyPermissionIssue).String()), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(expectedNewContractID), Index: false},
						{Key: "grantee", Value: w(s.vendor.String()), Index: false},
						{Key: "granter", Value: w(""), Index: false},
						{Key: "permission", Value: w(collection.Permission(collection.LegacyPermissionModify).String()), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(expectedNewContractID), Index: false},
						{Key: "grantee", Value: w(s.vendor.String()), Index: false},
						{Key: "granter", Value: w(""), Index: false},
						{Key: "permission", Value: w(collection.Permission(collection.LegacyPermissionMint).String()), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(expectedNewContractID), Index: false},
						{Key: "grantee", Value: w(s.vendor.String()), Index: false},
						{Key: "granter", Value: w(""), Index: false},
						{Key: "permission", Value: w(collection.Permission(collection.LegacyPermissionBurn).String()), Index: false},
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgCreateContract{
				Owner: tc.owner.String(),
			}
			res, err := s.msgServer.CreateContract(ctx, req)
			s.Require().Equal(expectedNewContractID, res.ContractId)
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
		amount     math.Int
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			owner:      s.vendor,
			amount:     math.ZeroInt(),
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventCreatedFTClass",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "decimals", Value: "0", Index: false},
						{Key: "meta", Value: w(""), Index: false},
						{Key: "mintable", Value: "false", Index: false},
						{Key: "name", Value: w(""), Index: false},
						{Key: "operator", Value: w(s.vendor.String()), Index: false},
						{Key: "token_id", Value: w(expectedTokenID), Index: false},
					},
				},
			},
		},
		"valid request with supply": {
			contractID: s.contractID,
			owner:      s.vendor,
			amount:     math.OneInt(),
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventCreatedFTClass",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "decimals", Value: "0", Index: false},
						{Key: "meta", Value: w(""), Index: false},
						{Key: "mintable", Value: "false", Index: false},
						{Key: "name", Value: w(""), Index: false},
						{Key: "operator", Value: w(s.vendor.String()), Index: false},
						{Key: "token_id", Value: w(expectedTokenID), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventMintedFT",
					Attributes: []abci.EventAttribute{
						{Key: "amount", Value: mustJSONMarshal(collection.NewCoins(collection.Coin{TokenId: expectedTokenID, Amount: math.OneInt()})), Index: false},
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "operator", Value: w(s.vendor.String()), Index: false},
						{Key: "to", Value: w(s.customer.String()), Index: false},
					},
				},
			},
		},
		"contract not found": {
			contractID: "deadbeef",
			owner:      s.vendor,
			amount:     math.ZeroInt(),
			err:        class.ErrContractNotExist,
		},
		"no permission": {
			contractID: s.contractID,
			owner:      s.customer,
			amount:     math.ZeroInt(),
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
			res, err := s.msgServer.IssueFT(ctx, req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
			s.Require().Equal(tc.events, ctx.EventManager().Events())

			// check balance and tokenID
			tokenID := collection.NewFTID(res.TokenId)
			bal, err := s.queryServer.Balance(ctx, &collection.QueryBalanceRequest{
				ContractId: s.contractID,
				Address:    s.customer.String(),
				TokenId:    tokenID,
			})
			s.Require().NoError(err)
			expectedCoin := collection.Coin{
				TokenId: tokenID,
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
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "meta", Value: w(""), Index: false},
						{Key: "name", Value: w(""), Index: false},
						{Key: "operator", Value: w(s.vendor.String()), Index: false},
						{Key: "token_type", Value: w(expectedTokenType), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "grantee", Value: w(s.vendor.String()), Index: false},
						{Key: "granter", Value: w(""), Index: false},
						{Key: "permission", Value: w(collection.Permission(collection.LegacyPermissionMint).String()), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "grantee", Value: w(s.vendor.String()), Index: false},
						{Key: "granter", Value: w(""), Index: false},
						{Key: "permission", Value: w(collection.Permission(collection.LegacyPermissionBurn).String()), Index: false},
					},
				},
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
			res, err := s.msgServer.IssueNFT(ctx, req)
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
		collection.NewFTCoin(s.ftClassID, math.NewInt(100000)),
	)
	amounts := collection.NewCoins(
		collection.NewFTCoin(s.ftClassID, math.NewInt(100000)),
		collection.NewFTCoin(*mintableFTClassID, math.NewInt(200000)),
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
						{Key: "amount", Value: mustJSONMarshal(amount), Index: false},
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "operator", Value: w(s.vendor.String()), Index: false},
						{Key: "to", Value: w(s.customer.String()), Index: false},
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
						{Key: "amount", Value: mustJSONMarshal(amounts), Index: false},
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "operator", Value: w(s.vendor.String()), Index: false},
						{Key: "to", Value: w(s.customer.String()), Index: false},
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
						{Key: "amount", Value: "[]", Index: false},
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "operator", Value: w(s.vendor.String()), Index: false},
						{Key: "to", Value: w(s.customer.String()), Index: false},
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
				collection.NewFTCoin("00bab10c", math.OneInt()),
			),
			err: collection.ErrTokenNotExist,
		},
		"include invalid tokenId among 2 tokens": {
			contractID: s.contractID,
			from:       s.vendor,
			amount: collection.NewCoins(
				collection.NewFTCoin(s.ftClassID, math.OneInt()),
				collection.NewFTCoin("00bab10b", math.OneInt()), // no exist tokenId
			),
			err: collection.ErrTokenNotExist,
		},
		"mintable false tokenId": {
			contractID: s.contractID,
			from:       s.vendor,
			amount:     collection.NewCoins(collection.NewFTCoin(*nonmintableFTClassID, math.OneInt())),
			err:        collection.ErrTokenNotMintable,
		},
		"include mintable false among 2 tokens": {
			contractID: s.contractID,
			from:       s.vendor,
			amount: collection.NewCoins(
				collection.NewFTCoin(*mintableFTClassID, math.OneInt()),
				collection.NewFTCoin(*nonmintableFTClassID, math.OneInt()),
			),
			err: collection.ErrTokenNotMintable,
		},
	}

	// query the values to be effected by MintFT
	queryValuesEffectedByMintFT := func(ctx sdk.Context, coins collection.Coins, contractID string) (balances collection.Coins, supply, minted []math.Int) {
		for _, am := range coins {
			// save balance
			bal, err := s.queryServer.Balance(ctx, &collection.QueryBalanceRequest{
				ContractId: contractID,
				Address:    s.customer.String(),
				TokenId:    am.TokenId,
			})
			s.Require().NoError(err)
			balances = append(balances, bal.Balance)

			// save supply
			res, err := s.queryServer.FTSupply(ctx, &collection.QueryFTSupplyRequest{
				ContractId: contractID,
				TokenId:    am.TokenId,
			})
			s.Require().NoError(err)
			supply = append(supply, res.Supply)

			// save minted
			m, err := s.queryServer.FTMinted(ctx, &collection.QueryFTMintedRequest{
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
				res, err := s.msgServer.MintFT(ctx, req)
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
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "operator", Value: w(s.vendor.String()), Index: false},
						{Key: "to", Value: w(s.customer.String()), Index: false},
						{Key: "tokens", Value: mustJSONMarshal(expectedTokens), Index: false},
					},
				},
			},
		},
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
			res, err := s.msgServer.MintNFT(ctx, req)
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
	singleAmount := collection.NewCoins(
		collection.NewFTCoin(s.ftClassID, math.NewInt(50000)),
	)

	// create a fungible token class
	mintableFTClassID, err := s.keeper.CreateTokenClass(s.ctx, s.contractID, &collection.FTClass{
		Name:     "tibetian fox2",
		Mintable: true,
	})
	s.Require().NoError(err)
	multiAmounts := collection.NewCoins(
		collection.NewFTCoin(s.ftClassID, math.NewInt(50000)),
		collection.NewFTCoin(*mintableFTClassID, math.NewInt(60000)),
	)

	// mintft
	mintedCoin := collection.NewFTCoin(*mintableFTClassID, math.NewInt(1000000))
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
			amount:     singleAmount,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventBurned",
					Attributes: []abci.EventAttribute{
						{Key: "amount", Value: mustJSONMarshal(collection.NewCoins(
							collection.NewFTCoin(s.ftClassID, math.NewInt(50000)),
						)), Index: false},
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "from", Value: w(s.vendor.String()), Index: false},
						{Key: "operator", Value: w(s.vendor.String()), Index: false},
					},
				},
			},
		},
		"valid multi amount burn": {
			contractID: s.contractID,
			from:       s.vendor,
			amount:     multiAmounts,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventBurned",
					Attributes: []abci.EventAttribute{
						{Key: "amount", Value: mustJSONMarshal(multiAmounts), Index: false},
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "from", Value: w(s.vendor.String()), Index: false},
						{Key: "operator", Value: w(s.vendor.String()), Index: false},
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
						{Key: "amount", Value: "[]", Index: false},
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "from", Value: w(s.vendor.String()), Index: false},
						{Key: "operator", Value: w(s.vendor.String()), Index: false},
					},
				},
			},
		},
		"contract not found": {
			contractID: "deadbeef",
			from:       s.vendor,
			amount:     singleAmount,
			err:        class.ErrContractNotExist,
		},
		"no permission": {
			contractID: s.contractID,
			from:       s.customer,
			amount:     singleAmount,
			err:        collection.ErrTokenNoPermission,
		},
		"insufficient funds": {
			contractID: s.contractID,
			from:       s.vendor,
			amount: collection.NewCoins(
				collection.NewFTCoin("00bab10c", math.OneInt()),
			),
			err: collection.ErrInsufficientToken,
		},
		"include insufficient funds amount 2 amounts": {
			contractID: s.contractID,
			from:       s.vendor,
			amount: collection.NewCoins(
				collection.NewFTCoin(s.ftClassID, s.balance),
				collection.NewFTCoin("00bab10c", math.OneInt()),
			),
			err: collection.ErrInsufficientToken,
		},
	}

	// query the values to be effected by BurnFT
	queryValuesAffectedByBurnFT := func(ctx sdk.Context, coins collection.Coins, contractID, from string) (balances collection.Coins, supply, burnt []math.Int) {
		for _, am := range coins {
			// save balance
			bal, err := s.queryServer.Balance(ctx, &collection.QueryBalanceRequest{
				ContractId: contractID,
				Address:    from,
				TokenId:    am.TokenId,
			})
			s.Require().NoError(err)
			balances = append(balances, bal.Balance)

			// save supply
			res, err := s.queryServer.FTSupply(ctx, &collection.QueryFTSupplyRequest{
				ContractId: contractID,
				TokenId:    am.TokenId,
			})
			s.Require().NoError(err)
			supply = append(supply, res.Supply)

			// save minted
			b, err := s.queryServer.FTBurnt(ctx, &collection.QueryFTBurntRequest{
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
				res, err := s.msgServer.BurnFT(ctx, req)
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
	singleAmount := collection.NewCoins(
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
			amount:     singleAmount,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventBurned",
					Attributes: []abci.EventAttribute{
						{Key: "amount", Value: mustJSONMarshal(singleAmount), Index: false},
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "from", Value: w(s.customer.String()), Index: false},
						{Key: "operator", Value: w(s.operator.String()), Index: false},
					},
				},
			},
		},
		"contract not found": {
			contractID: "deadbeef",
			operator:   s.operator,
			from:       s.customer,
			amount:     singleAmount,
			err:        class.ErrContractNotExist,
		},
		"no authorization": {
			contractID: s.contractID,
			operator:   s.vendor,
			from:       s.customer,
			amount:     singleAmount,
			err:        collection.ErrCollectionNotApproved,
		},
		"no permission": {
			contractID: s.contractID,
			operator:   s.stranger,
			from:       s.customer,
			amount:     singleAmount,
			err:        collection.ErrTokenNoPermission,
		},
		"insufficient funds - exist token": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			amount: collection.NewCoins(
				collection.NewFTCoin(s.ftClassID, s.balance.Add(math.OneInt())),
			),
			err: collection.ErrInsufficientToken,
		},
		"insufficient funds - non-exist token": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			amount: collection.NewCoins(
				collection.NewFTCoin("00bab10c", math.OneInt()),
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
			res, err := s.msgServer.OperatorBurnFT(ctx, req)
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
	rootNFTID := collection.NewNFTID(s.nftClassID, s.numNFTs*2+1)
	issuedTokenIDs := s.extractChainedNFTIDs(rootNFTID)
	coins := make([]collection.Coin, 0)
	for _, id := range issuedTokenIDs {
		coins = append(coins, collection.NewCoin(id, math.NewInt(1)))
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
			tokenIDs:   []string{rootNFTID},
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventBurned",
					Attributes: []abci.EventAttribute{
						{Key: "amount", Value: mustJSONMarshal(coins), Index: false},
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "from", Value: w(s.vendor.String()), Index: false},
						{Key: "operator", Value: w(s.vendor.String()), Index: false},
					},
				},
			},
		},
		"contract not found": {
			contractID: "deadbeef",
			from:       s.vendor,
			tokenIDs:   []string{rootNFTID},
			err:        class.ErrContractNotExist,
		},
		"no permission": {
			contractID: s.contractID,
			from:       s.customer,
			tokenIDs:   []string{rootNFTID},
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
			res, err := s.msgServer.BurnNFT(ctx, req)
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
	rootNFTID := collection.NewNFTID(s.nftClassID, 1)
	issuedTokenIDs := s.extractChainedNFTIDs(rootNFTID)
	coins := make([]collection.Coin, 0)
	for _, id := range issuedTokenIDs {
		coins = append(coins, collection.NewCoin(id, math.NewInt(1)))
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
			tokenIDs:   []string{rootNFTID},
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventBurned",
					Attributes: []abci.EventAttribute{
						{Key: "amount", Value: mustJSONMarshal(coins), Index: false},
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "from", Value: w(s.customer.String()), Index: false},
						{Key: "operator", Value: w(s.operator.String()), Index: false},
					},
				},
			},
		},
		"contract not found": {
			contractID: "deadbeef",
			operator:   s.operator,
			from:       s.customer,
			tokenIDs:   []string{rootNFTID},
			err:        class.ErrContractNotExist,
		},
		"no authorization": {
			contractID: s.contractID,
			operator:   s.vendor,
			from:       s.customer,
			tokenIDs:   []string{rootNFTID},
			err:        collection.ErrCollectionNotApproved,
		},
		"no permission": {
			contractID: s.contractID,
			operator:   s.stranger,
			from:       s.customer,
			tokenIDs:   []string{rootNFTID},
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
			res, err := s.msgServer.OperatorBurnNFT(ctx, req)
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
	expectedTokenIndex := collection.NewNFTID(s.nftClassID, 1)[8:]
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
						{Key: "changes", Value: mustJSONMarshal(changes), Index: false},
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "operator", Value: w(s.vendor.String()), Index: false},
					},
				},
			},
		},
		"contract not found": {
			contractID: "deadbeef",
			operator:   s.vendor,
			err:        class.ErrContractNotExist,
		},
		"no permission": {
			contractID: s.contractID,
			operator:   s.customer,
			tokenType:  s.nftClassID,
			tokenIndex: expectedTokenIndex,
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
			res, err := s.msgServer.Modify(ctx, req)
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
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "grantee", Value: w(s.operator.String()), Index: false},
						{Key: "granter", Value: w(s.vendor.String()), Index: false},
						{Key: "permission", Value: w(collection.Permission(collection.LegacyPermissionModify).String()), Index: false},
					},
				},
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
			res, err := s.msgServer.GrantPermission(ctx, req)
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
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "grantee", Value: w(s.operator.String()), Index: false},
						{Key: "permission", Value: w(collection.Permission(collection.LegacyPermissionMint).String()), Index: false},
					},
				},
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
			res, err := s.msgServer.RevokePermission(ctx, req)
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
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "holder", Value: w(s.customer.String()), Index: false},
						{Key: "operator", Value: w(s.customer.String()), Index: false},
						{Key: "subject", Value: w(collection.NewNFTID(s.nftClassID, s.depthLimit+1)), Index: false},
						{Key: "target", Value: w(collection.NewNFTID(s.nftClassID, 1)), Index: false},
					},
				},
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
			res, err := s.msgServer.Attach(ctx, req)
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
	issuedNfts := make([]string, s.depthLimit)
	for i := 1; i <= s.depthLimit; i++ {
		issuedNfts[i-1] = collection.NewNFTID(s.nftClassID, i)
	}

	testCases := map[string]struct {
		contractID string
		subjectID  string
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			subjectID:  issuedNfts[1],
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventDetached",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "holder", Value: w(s.customer.String()), Index: false},
						{Key: "operator", Value: w(s.customer.String()), Index: false},
						{Key: "previous_parent", Value: w(issuedNfts[0]), Index: false},
						{Key: "subject", Value: w(issuedNfts[1]), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventRootChanged",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "from", Value: w(issuedNfts[0]), Index: false},
						{Key: "to", Value: w(issuedNfts[1]), Index: false},
						{Key: "token_id", Value: w(issuedNfts[2]), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventRootChanged",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "from", Value: w(issuedNfts[0]), Index: false},
						{Key: "to", Value: w(issuedNfts[1]), Index: false},
						{Key: "token_id", Value: w(issuedNfts[3]), Index: false},
					},
				},
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
			res, err := s.msgServer.Detach(ctx, req)
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
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "holder", Value: w(s.customer.String()), Index: false},
						{Key: "operator", Value: w(s.operator.String()), Index: false},
						{Key: "subject", Value: w(collection.NewNFTID(s.nftClassID, s.depthLimit+1)), Index: false},
						{Key: "target", Value: w(collection.NewNFTID(s.nftClassID, 1)), Index: false},
					},
				},
			},
		},
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
			res, err := s.msgServer.OperatorAttach(ctx, req)
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
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "holder", Value: w(s.customer.String()), Index: false},
						{Key: "operator", Value: w(s.operator.String()), Index: false},
						{Key: "previous_parent", Value: w(nfts[0]), Index: false},
						{Key: "subject", Value: w(nfts[1]), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventRootChanged",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "from", Value: w(nfts[0]), Index: false},
						{Key: "to", Value: w(nfts[1]), Index: false},
						{Key: "token_id", Value: w(nfts[2]), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventRootChanged",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "from", Value: w(nfts[0]), Index: false},
						{Key: "to", Value: w(nfts[1]), Index: false},
						{Key: "token_id", Value: w(nfts[3]), Index: false},
					},
				},
			},
		},
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
			res, err := s.msgServer.OperatorDetach(ctx, req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
			s.Require().Equal(tc.events, ctx.EventManager().Events())
		})
	}
}

func (s *KeeperTestSuite) extractChainedNFTIDs(root string) []string {
	allTokenIDs := make([]string, 0)
	allTokenIDs = append(allTokenIDs, root)
	cursor := allTokenIDs[0]
	for {
		ctx, _ := s.ctx.CacheContext()
		res, err := s.queryServer.Children(ctx, &collection.QueryChildrenRequest{
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
	return allTokenIDs
}
