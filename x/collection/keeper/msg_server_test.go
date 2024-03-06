package keeper_test

import (
	"encoding/json"
	"fmt"
	"strings"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/gogoproto/proto"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"

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
	case sdk.AccAddress:
		ac := authcodec.NewBech32Codec("link")
		s, err := ac.BytesToString(input.([]byte))
		if err != nil {
			panic(err)
		}
		return fmt.Sprintf("\"%s\"", s)
	case sdk.ValAddress:
		ac := authcodec.NewBech32Codec("linkvaloper")
		s, err := ac.BytesToString(input.([]byte))
		if err != nil {
			panic(err)
		}
		return fmt.Sprintf("\"%s\"", s)
	case sdk.ConsAddress:
		ac := authcodec.NewBech32Codec("linkvalcons")
		s, err := ac.BytesToString(input.([]byte))
		if err != nil {
			panic(err)
		}
		return fmt.Sprintf("\"%s\"", s)
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
		from       sdk.AccAddress
		to         sdk.AccAddress
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			tokenID:    rootNFTID,
			from:       s.customer,
			to:         s.vendor,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventSent",
					Attributes: []abci.EventAttribute{
						{Key: "amount", Value: mustJSONMarshal(collection.NewCoins(collection.Coin{TokenId: rootNFTID, Amount: math.OneInt()})), Index: false},
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "from", Value: w(s.bytesToString(s.customer)), Index: false},
						{Key: "operator", Value: w(s.bytesToString(s.customer)), Index: false},
						{Key: "to", Value: w(s.bytesToString(s.vendor)), Index: false},
					},
				},
			},
		},
		"contract not found": {
			contractID: "deadbeef",
			tokenID:    collection.NewNFTID(s.nftClassID, 1),
			from:       s.customer,
			to:         s.vendor,
			err:        collection.ErrContractNotExist,
		},
		"NFT not found": {
			contractID: s.contractID,
			tokenID:    collection.NewNFTID("deadbeef", 1),
			from:       s.customer,
			to:         s.vendor,
			err:        collection.ErrTokenNotExist,
		},
		"not owned by": {
			contractID: s.contractID,
			tokenID:    collection.NewNFTID(s.nftClassID, s.numNFTs+1),
			from:       s.customer,
			to:         s.vendor,
			err:        collection.ErrTokenNotOwnedBy,
		},
		"invalid from": {
			contractID: s.contractID,
			tokenID:    rootNFTID,
			to:         s.vendor,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid contract id": {
			tokenID: rootNFTID,
			from:    s.customer,
			to:      s.vendor,
			err:     collection.ErrInvalidContractID,
		},
		"invalid to": {
			contractID: s.contractID,
			tokenID:    rootNFTID,
			from:       s.customer,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"empty token ids": {
			contractID: s.contractID,
			from:       s.customer,
			to:         s.vendor,
			err:        collection.ErrEmptyField,
		},
		"invalid token ids": {
			contractID: s.contractID,
			tokenID:    "null",
			from:       s.customer,
			to:         s.vendor,
			err:        collection.ErrInvalidTokenID,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			var req *collection.MsgSendNFT
			if tc.tokenID == "" {
				req = &collection.MsgSendNFT{
					ContractId: tc.contractID,
					From:       s.bytesToString(tc.from),
					To:         s.bytesToString(tc.to),
				}
			} else {
				req = &collection.MsgSendNFT{
					ContractId: tc.contractID,
					From:       s.bytesToString(tc.from),
					To:         s.bytesToString(tc.to),
					TokenIds:   []string{tc.tokenID},
				}
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
		to         sdk.AccAddress
		tokenID    string
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			to:         s.vendor,
			tokenID:    rootNFTID,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventSent",
					Attributes: []abci.EventAttribute{
						{Key: "amount", Value: mustJSONMarshal(collection.NewCoins(collection.Coin{TokenId: rootNFTID, Amount: math.OneInt()})), Index: false},
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "from", Value: w(s.bytesToString(s.customer)), Index: false},
						{Key: "operator", Value: w(s.bytesToString(s.operator)), Index: false},
						{Key: "to", Value: w(s.bytesToString(s.vendor)), Index: false},
					},
				},
			},
		},
		"contract not found": {
			contractID: "deadbeef",
			operator:   s.operator,
			from:       s.customer,
			to:         s.vendor,
			tokenID:    rootNFTID,
			err:        collection.ErrContractNotExist,
		},
		"not approved": {
			contractID: s.contractID,
			operator:   s.vendor,
			from:       s.customer,
			to:         s.vendor,
			tokenID:    rootNFTID,
			err:        collection.ErrCollectionNotApproved,
		},
		"NFT not found": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			to:         s.vendor,
			tokenID:    collection.NewNFTID("deadbeef", 1),
			err:        collection.ErrTokenNotExist,
		},
		"not owned by": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			to:         s.vendor,
			tokenID:    collection.NewNFTID(s.nftClassID, s.numNFTs+1),
			err:        collection.ErrTokenNotOwnedBy,
		},
		"invalid operator": {
			contractID: s.contractID,
			from:       s.customer,
			to:         s.vendor,
			tokenID:    rootNFTID,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid contract id": {
			operator: s.operator,
			from:     s.customer,
			to:       s.vendor,
			tokenID:  rootNFTID,
			err:      collection.ErrInvalidContractID,
		},
		"invalid from": {
			contractID: s.contractID,
			operator:   s.operator,
			to:         s.vendor,
			tokenID:    rootNFTID,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid to": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			tokenID:    rootNFTID,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"empty ids": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			to:         s.vendor,
			err:        collection.ErrEmptyField,
		},
		"invalid id": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			to:         s.vendor,
			tokenID:    "null",
			err:        collection.ErrInvalidTokenID,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			var req *collection.MsgOperatorSendNFT
			if tc.tokenID == "" {
				req = &collection.MsgOperatorSendNFT{
					ContractId: tc.contractID,
					Operator:   s.bytesToString(tc.operator),
					From:       s.bytesToString(tc.from),
					To:         s.bytesToString(tc.to),
				}
			} else {
				req = &collection.MsgOperatorSendNFT{
					ContractId: tc.contractID,
					Operator:   s.bytesToString(tc.operator),
					From:       s.bytesToString(tc.from),
					To:         s.bytesToString(tc.to),
					TokenIds:   []string{tc.tokenID},
				}
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
		contractID string
		holder     sdk.AccAddress
		operator   sdk.AccAddress
		events     sdk.Events
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			holder:     s.customer,
			operator:   s.vendor,
			events: sdk.Events{sdk.Event{
				Type: "lbm.collection.v1.EventAuthorizedOperator",
				Attributes: []abci.EventAttribute{
					{Key: "contract_id", Value: w(s.contractID), Index: false},
					{Key: "holder", Value: w(s.bytesToString(s.customer)), Index: false},
					{Key: "operator", Value: w(s.bytesToString(s.vendor)), Index: false},
				},
			}},
		},
		"contract not found": {
			contractID: "deadbeef",
			holder:     s.customer,
			operator:   s.vendor,
			err:        collection.ErrContractNotExist,
		},
		"already approved": {
			contractID: s.contractID,
			holder:     s.customer,
			operator:   s.operator,
			err:        collection.ErrCollectionAlreadyApproved,
		},
		"invalid contract id": {
			holder:   s.customer,
			operator: s.operator,
			err:      collection.ErrInvalidContractID,
		},
		"invalid holder": {
			contractID: s.contractID,
			operator:   s.operator,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"empty operator": {
			contractID: s.contractID,
			holder:     s.customer,
			err:        sdkerrors.ErrInvalidAddress,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			// Arrange
			ctx, _ := s.ctx.CacheContext()
			prevAuth, _ := s.keeper.GetAuthorization(ctx, tc.contractID, tc.holder, tc.operator)

			// Act
			req := &collection.MsgAuthorizeOperator{
				ContractId: tc.contractID,
				Holder:     s.bytesToString(tc.holder),
				Operator:   s.bytesToString(tc.operator),
			}
			res, err := s.msgServer.AuthorizeOperator(ctx, req)
			if tc.err != nil {
				s.Require().ErrorIs(err, tc.err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
			s.Require().Equal(tc.events, ctx.EventManager().Events())
			curAuth, err := s.keeper.GetAuthorization(ctx, tc.contractID, tc.holder, tc.operator)
			s.Require().NoError(err)
			s.Require().Nil(prevAuth)
			h, err := s.addressCodec.BytesToString(tc.holder.Bytes())
			s.Require().NoError(err)
			s.Require().Equal(h, curAuth.Holder)
			o, err := s.addressCodec.BytesToString(tc.operator.Bytes())
			s.Require().NoError(err)
			s.Require().Equal(o, curAuth.Operator)
		})
	}
}

func (s *KeeperTestSuite) TestMsgRevokeOperator() {
	testCases := map[string]struct {
		contractID string
		holder     sdk.AccAddress
		operator   sdk.AccAddress
		events     sdk.Events
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			holder:     s.customer,
			operator:   s.operator,
			events: sdk.Events{sdk.Event{
				Type: "lbm.collection.v1.EventRevokedOperator",
				Attributes: []abci.EventAttribute{
					{Key: "contract_id", Value: w(s.contractID), Index: false},
					{Key: "holder", Value: w(s.bytesToString(s.customer)), Index: false},
					{Key: "operator", Value: w(s.bytesToString(s.operator)), Index: false},
				},
			}},
		},
		"contract not found": {
			contractID: "deadbeef",
			holder:     s.customer,
			operator:   s.vendor,
			err:        collection.ErrContractNotExist,
		},
		"no authorization": {
			contractID: s.contractID,
			holder:     s.customer,
			operator:   s.vendor,
			err:        collection.ErrCollectionNotApproved,
		},
		"invalid contract id": {
			holder:   s.customer,
			operator: s.operator,
			err:      collection.ErrInvalidContractID,
		},
		"invalid holder": {
			contractID: s.contractID,
			operator:   s.operator,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"empty operator": {
			contractID: s.contractID,
			holder:     s.customer,
			err:        sdkerrors.ErrInvalidAddress,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			// Arrange
			ctx, _ := s.ctx.CacheContext()
			prevAuth, _ := s.keeper.GetAuthorization(ctx, tc.contractID, tc.holder, tc.operator)

			// Act
			req := &collection.MsgRevokeOperator{
				ContractId: tc.contractID,
				Holder:     s.bytesToString(tc.holder),
				Operator:   s.bytesToString(tc.operator),
			}
			res, err := s.msgServer.RevokeOperator(ctx, req)
			if tc.err != nil {
				s.Require().ErrorIs(err, tc.err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)

			s.Require().Equal(tc.events, ctx.EventManager().Events())
			s.Require().NotNil(prevAuth)
			s.Require().Equal(s.bytesToString(tc.holder), prevAuth.Holder)
			s.Require().Equal(s.bytesToString(tc.operator), prevAuth.Operator)
			curAuth, err := s.keeper.GetAuthorization(ctx, tc.contractID, tc.holder, tc.operator)
			s.Require().ErrorIs(err, collection.ErrCollectionNotApproved)
			s.Require().Nil(curAuth)
		})
	}
}

func (s *KeeperTestSuite) TestMsgCreateContract() {
	expectedNewContractID := "3336b76f"
	testCases := map[string]struct {
		owner  sdk.AccAddress
		name   string
		uri    string
		meta   string
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
						{Key: "creator", Value: w(s.bytesToString(s.vendor)), Index: false},
						{Key: "meta", Value: w(""), Index: false},
						{Key: "name", Value: w(""), Index: false},
						{Key: "uri", Value: w(""), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(expectedNewContractID), Index: false},
						{Key: "grantee", Value: w(s.bytesToString(s.vendor)), Index: false},
						{Key: "granter", Value: w(""), Index: false},
						{Key: "permission", Value: w(collection.Permission(collection.LegacyPermissionIssue).String()), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(expectedNewContractID), Index: false},
						{Key: "grantee", Value: w(s.bytesToString(s.vendor)), Index: false},
						{Key: "granter", Value: w(""), Index: false},
						{Key: "permission", Value: w(collection.Permission(collection.LegacyPermissionModify).String()), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(expectedNewContractID), Index: false},
						{Key: "grantee", Value: w(s.bytesToString(s.vendor)), Index: false},
						{Key: "granter", Value: w(""), Index: false},
						{Key: "permission", Value: w(collection.Permission(collection.LegacyPermissionMint).String()), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(expectedNewContractID), Index: false},
						{Key: "grantee", Value: w(s.bytesToString(s.vendor)), Index: false},
						{Key: "granter", Value: w(""), Index: false},
						{Key: "permission", Value: w(collection.Permission(collection.LegacyPermissionBurn).String()), Index: false},
					},
				},
			},
		},
		"invalid owner": {
			name: "tibetian fox",
			uri:  "file:///tibetian_fox.png",
			meta: "Tibetian fox",
			err:  sdkerrors.ErrInvalidAddress,
		},
		"long name": {
			owner: s.vendor,
			name:  string(make([]rune, 21)),
			uri:   "file:///tibetian_fox.png",
			meta:  "Tibetian fox",
			err:   collection.ErrInvalidNameLength,
		},
		"invalid base image uri": {
			owner: s.vendor,
			name:  "tibetian fox",
			uri:   string(make([]rune, 1001)),
			meta:  "Tibetian fox",
			err:   collection.ErrInvalidBaseImgURILength,
		},
		"invalid meta": {
			owner: s.vendor,
			name:  "tibetian fox",
			uri:   "file:///tibetian_fox.png",
			meta:  string(make([]rune, 1001)),
			err:   collection.ErrInvalidMetaLength,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgCreateContract{
				Owner: s.bytesToString(tc.owner),
				Name:  tc.name,
				Uri:   tc.uri,
				Meta:  tc.meta,
			}
			res, err := s.msgServer.CreateContract(ctx, req)
			if tc.err != nil {
				s.Require().ErrorIs(err, tc.err)
				return
			}
			s.Require().Equal(expectedNewContractID, res.ContractId)

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
		name       string
		meta       string
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
						{Key: "operator", Value: w(s.bytesToString(s.vendor)), Index: false},
						{Key: "token_type", Value: w(expectedTokenType), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "grantee", Value: w(s.bytesToString(s.vendor)), Index: false},
						{Key: "granter", Value: w(""), Index: false},
						{Key: "permission", Value: w(collection.Permission(collection.LegacyPermissionMint).String()), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventGranted",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "grantee", Value: w(s.bytesToString(s.vendor)), Index: false},
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
		"invalid contract id": {
			owner: s.vendor,
			name:  "tibetian fox",
			meta:  "Tibetian Fox",
			err:   collection.ErrInvalidContractID,
		},
		"invalid operator": {
			contractID: s.contractID,
			name:       "tibetian fox",
			meta:       "Tibetian Fox",
			err:        sdkerrors.ErrInvalidAddress,
		},
		"long name": {
			contractID: s.contractID,
			owner:      s.vendor,
			meta:       "Tibetian Fox",
			name:       string(make([]rune, 21)),
			err:        collection.ErrInvalidNameLength,
		},
		"invalid meta": {
			contractID: s.contractID,
			owner:      s.vendor,
			name:       "tibetian fox",
			meta:       string(make([]rune, 1001)),
			err:        collection.ErrInvalidMetaLength,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgIssueNFT{
				ContractId: tc.contractID,
				Owner:      s.bytesToString(tc.owner),
				Name:       tc.name,
				Meta:       tc.meta,
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
		to         sdk.AccAddress
		params     []collection.MintNFTParam
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			from:       s.vendor,
			to:         s.customer,
			params:     params,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventMintedNFT",
					Attributes: []abci.EventAttribute{
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "operator", Value: w(s.bytesToString(s.vendor)), Index: false},
						{Key: "to", Value: w(s.bytesToString(s.customer)), Index: false},
						{Key: "tokens", Value: mustJSONMarshal(expectedTokens), Index: false},
					},
				},
			},
		},
		"contract not found": {
			contractID: "deadbeef",
			from:       s.vendor,
			to:         s.customer,
			params:     params,
			err:        collection.ErrContractNotExist,
		},
		"no permission": {
			contractID: s.contractID,
			from:       s.customer,
			to:         s.customer,
			params:     params,
			err:        collection.ErrTokenNoPermission,
		},
		"no class of the token": {
			contractID: s.contractID,
			from:       s.vendor,
			to:         s.customer,
			params: []collection.MintNFTParam{{
				Name:      "tibetian fox",
				TokenType: "deadbeef",
			}},
			err: collection.ErrTokenTypeNotExist,
		},
		"invalid contract id": {
			from:   s.vendor,
			to:     s.customer,
			params: params,
			err:    collection.ErrInvalidContractID,
		},
		"invalid operator": {
			contractID: s.contractID,
			to:         s.customer,
			params:     params,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid to": {
			contractID: s.contractID,
			from:       s.vendor,
			params:     params,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"empty params": {
			contractID: s.contractID,
			from:       s.vendor,
			to:         s.customer,
			err:        collection.ErrEmptyField,
		},
		"param of invalid token type": {
			contractID: s.contractID,
			from:       s.vendor,
			to:         s.customer,
			params: []collection.MintNFTParam{{
				Name: "tibetian fox",
			}},
			err: collection.ErrInvalidTokenType,
		},
		"param of empty name": {
			contractID: s.contractID,
			from:       s.vendor,
			to:         s.customer,
			params: []collection.MintNFTParam{{
				TokenType: s.nftClassID,
			}},
			err: collection.ErrInvalidTokenName,
		},
		"param of too long name": {
			contractID: s.contractID,
			from:       s.vendor,
			to:         s.customer,
			params: []collection.MintNFTParam{{
				TokenType: s.nftClassID,
				Name:      string(make([]rune, 21)),
			}},
			err: collection.ErrInvalidNameLength,
		},
		"param of invalid meta": {
			contractID: s.contractID,
			from:       s.vendor,
			to:         s.customer,
			params: []collection.MintNFTParam{{
				TokenType: s.nftClassID,
				Name:      "tibetian fox",
				Meta:      string(make([]rune, 1001)),
			}},
			err: collection.ErrInvalidMetaLength,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgMintNFT{
				ContractId: tc.contractID,
				From:       s.bytesToString(tc.from),
				To:         s.bytesToString(tc.to),
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
						{Key: "from", Value: w(s.bytesToString(s.vendor)), Index: false},
						{Key: "operator", Value: w(s.bytesToString(s.vendor)), Index: false},
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
		"NFT not found": {
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
		"invalid contract id": {
			from:     s.vendor,
			tokenIDs: []string{rootNFTID},
			err:      collection.ErrInvalidContractID,
		},
		"invalid from": {
			contractID: s.contractID,
			tokenIDs:   []string{rootNFTID},
			err:        sdkerrors.ErrInvalidAddress,
		},
		"empty ids": {
			contractID: s.contractID,
			from:       s.vendor,
			err:        collection.ErrEmptyField,
		},
		"invalid id": {
			contractID: s.contractID,
			from:       s.vendor,
			tokenIDs:   []string{""},
			err:        collection.ErrInvalidTokenID,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgBurnNFT{
				ContractId: tc.contractID,
				From:       s.bytesToString(tc.from),
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
						{Key: "from", Value: w(s.bytesToString(s.customer)), Index: false},
						{Key: "operator", Value: w(s.bytesToString(s.operator)), Index: false},
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
		"NFT not found": {
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
		"invalid contract id": {
			operator: s.operator,
			from:     s.customer,
			tokenIDs: []string{rootNFTID},
			err:      collection.ErrInvalidContractID,
		},
		"invalid grantee": {
			contractID: s.contractID,
			from:       s.customer,
			tokenIDs:   []string{rootNFTID},
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid from": {
			contractID: s.contractID,
			operator:   s.operator,
			tokenIDs:   []string{rootNFTID},
			err:        sdkerrors.ErrInvalidAddress,
		},
		"empty ids": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			err:        collection.ErrEmptyField,
		},
		"invalid id": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			tokenIDs:   []string{""},
			err:        collection.ErrInvalidTokenID,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgOperatorBurnNFT{
				ContractId: tc.contractID,
				Operator:   s.bytesToString(tc.operator),
				From:       s.bytesToString(tc.from),
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
	validChanges := []collection.Attribute{{
		Key:   collection.AttributeKeyName.String(),
		Value: "test",
	}}

	testCases := map[string]struct {
		operator   sdk.AccAddress
		contractID string
		tokenType  string
		tokenIndex string
		changes    []collection.Attribute
		err        error
		events     sdk.Events
	}{
		"valid request - Contract": {
			contractID: s.contractID,
			operator:   s.vendor,
			changes:    validChanges,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventModifiedContract",
					Attributes: []abci.EventAttribute{
						{Key: "changes", Value: mustJSONMarshal(validChanges), Index: false},
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "operator", Value: w(s.bytesToString(s.vendor)), Index: false},
					},
				},
			},
		},
		"valid request - TokenClass": {
			contractID: s.contractID,
			operator:   s.vendor,
			tokenType:  s.nftClassID,
			changes:    validChanges,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventModifiedTokenClass",
					Attributes: []abci.EventAttribute{
						{Key: "changes", Value: mustJSONMarshal(validChanges), Index: false},
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "operator", Value: w(s.bytesToString(s.vendor)), Index: false},
						{Key: "token_type", Value: w(s.nftClassID), Index: false},
						{Key: "type_name", Value: w(proto.MessageName(&collection.NFTClass{})), Index: false},
					},
				},
			},
		},
		"valid request - NFT": {
			contractID: s.contractID,
			operator:   s.vendor,
			tokenType:  s.nftClassID,
			tokenIndex: s.issuedNFTs[s.bytesToString(s.vendor)][0].TokenId[8:],
			changes:    validChanges,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventModifiedNFT",
					Attributes: []abci.EventAttribute{
						{Key: "changes", Value: mustJSONMarshal(validChanges), Index: false},
						{Key: "contract_id", Value: w(s.contractID), Index: false},
						{Key: "operator", Value: w(s.bytesToString(s.vendor)), Index: false},
						{Key: "token_id", Value: w(s.issuedNFTs[s.bytesToString(s.vendor)][0].TokenId), Index: false},
					},
				},
			},
		},
		"contract not found": {
			contractID: "deadbeef",
			operator:   s.vendor,
			changes:    validChanges,
			err:        collection.ErrContractNotExist,
		},
		"no permission": {
			contractID: s.contractID,
			operator:   s.customer,
			tokenType:  s.nftClassID,
			tokenIndex: expectedTokenIndex,
			changes:    validChanges,
			err:        collection.ErrTokenNoPermission,
		},
		"nft not found": {
			contractID: s.contractID,
			operator:   s.vendor,
			tokenType:  s.nftClassID,
			tokenIndex: collection.NewNFTID(s.nftClassID, s.numNFTs*3+1)[8:],
			changes:    validChanges,
			err:        collection.ErrTokenNotExist,
		},
		"nft class not found": {
			contractID: s.contractID,
			operator:   s.vendor,
			tokenType:  "deadbeef",
			changes:    validChanges,
			err:        collection.ErrTokenTypeNotExist,
		},
		"invalid contract id": {
			operator: s.vendor,
			changes:  validChanges,
			err:      collection.ErrInvalidContractID,
		},
		"invalid owner": {
			contractID: s.contractID,
			changes:    validChanges,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid key of change": {
			contractID: s.contractID,
			operator:   s.vendor,
			changes:    []collection.Attribute{{Key: strings.ToUpper(collection.AttributeKeyName.String()), Value: "tt"}},
			err:        collection.ErrInvalidChangesField,
		},
		"invalid value of change": {
			contractID: s.contractID,
			operator:   s.vendor,
			changes:    []collection.Attribute{{Key: collection.AttributeKeyName.String(), Value: string(make([]rune, 21))}},
			err:        collection.ErrInvalidNameLength,
		},
		"empty changes": {
			contractID: s.contractID,
			operator:   s.vendor,
			err:        collection.ErrEmptyChanges,
		},
		"too many changes": {
			contractID: s.contractID,
			operator:   s.vendor,
			changes:    make([]collection.Attribute, 101),
			err:        collection.ErrInvalidChangesFieldCount,
		},
		"duplicated changes": {
			contractID: s.contractID,
			operator:   s.vendor,
			changes: []collection.Attribute{
				{Key: collection.AttributeKeyBaseImgURI.String(), Value: "hello"},
				{Key: collection.AttributeKeyURI.String(), Value: "world"},
			},
			err: collection.ErrDuplicateChangesField,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			req := &collection.MsgModify{
				ContractId: tc.contractID,
				Owner:      s.bytesToString(tc.operator),
				TokenType:  tc.tokenType,
				TokenIndex: tc.tokenIndex,
				Changes:    tc.changes,
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
						{Key: "grantee", Value: w(s.bytesToString(s.operator)), Index: false},
						{Key: "granter", Value: w(s.bytesToString(s.vendor)), Index: false},
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
		"invalid contract id": {
			granter:    s.vendor,
			grantee:    s.operator,
			permission: collection.LegacyPermissionMint.String(),
			err:        collection.ErrInvalidContractID,
		},
		"invalid from": {
			contractID: s.contractID,
			grantee:    s.operator,
			permission: collection.LegacyPermissionMint.String(),
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid to": {
			contractID: s.contractID,
			granter:    s.vendor,
			permission: collection.LegacyPermissionMint.String(),
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid permission": {
			contractID: s.contractID,
			granter:    s.vendor,
			grantee:    s.operator,
			err:        collection.ErrInvalidPermission,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgGrantPermission{
				ContractId: tc.contractID,
				From:       s.bytesToString(tc.granter),
				To:         s.bytesToString(tc.grantee),
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
						{Key: "grantee", Value: w(s.bytesToString(s.operator)), Index: false},
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
		"invalid contract id": {
			from:       s.operator,
			permission: collection.LegacyPermissionMint.String(),
			err:        collection.ErrInvalidContractID,
		},
		"invalid from": {
			contractID: s.contractID,
			permission: collection.LegacyPermissionMint.String(),
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid permission": {
			contractID: s.contractID,
			from:       s.operator,
			err:        collection.ErrInvalidPermission,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &collection.MsgRevokePermission{
				ContractId: tc.contractID,
				From:       s.bytesToString(tc.from),
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
