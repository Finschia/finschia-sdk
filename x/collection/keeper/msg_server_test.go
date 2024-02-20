package keeper_test

import (
	"encoding/json"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Finschia/finschia-sdk/x/collection"
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

func (s *KeeperTestSuite) TestMsgSendNFT() {
	rootNFTID := collection.NewNFTID(s.nftClassID, 1)

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
					Type: "lbm.collection.v1.EventSent",
					Attributes: []abci.EventAttribute{
						{Key: "amount", Value: mustJSONMarshal(collection.NewCoins(collection.Coin{TokenId: rootNFTID, Amount: math.OneInt()})), Index: false},
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
			err:        collection.ErrContractNotExist,
		},
		"not found": {
			contractID: s.contractID,
			tokenID:    collection.NewNFTID("deadbeef", 1),
			err:        collection.ErrTokenNotExist,
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
					Type: "lbm.collection.v1.EventSent",
					Attributes: []abci.EventAttribute{
						{Key: "amount", Value: mustJSONMarshal(collection.NewCoins(collection.Coin{TokenId: rootNFTID, Amount: math.OneInt()})), Index: false},
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
			err:        collection.ErrContractNotExist,
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
			expectedError: collection.ErrContractNotExist,
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
			expectedError: collection.ErrContractNotExist,
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
			err:        collection.ErrContractNotExist,
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

func (s *KeeperTestSuite) TestMsgMintNFT() {
	params := []collection.MintNFTParam{{
		TokenType: s.nftClassID,
		Name:      "tester",
		Meta:      "Mint NFT",
	}}
	expectedTokens := []collection.NFT{
		{
			TokenId: "100000010000000d",
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
			err:        collection.ErrContractNotExist,
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

func (s *KeeperTestSuite) TestMsgBurnNFT() {
	rootNFTID := collection.NewNFTID(s.nftClassID, s.numNFTs*2+1)
	coins := []collection.Coin{collection.NewCoin(rootNFTID, math.NewInt(1))}

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
			err:        collection.ErrContractNotExist,
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
	coins := []collection.Coin{collection.NewCoin(rootNFTID, math.NewInt(1))}

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
			err:        collection.ErrContractNotExist,
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
			err:        collection.ErrContractNotExist,
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
			err:        collection.ErrContractNotExist,
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
			err:        collection.ErrContractNotExist,
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
