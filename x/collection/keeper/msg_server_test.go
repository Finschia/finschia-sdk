package keeper_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/Finschia/finschia-sdk/types"
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
				helperBuildEventSent(s.contractID, s.vendor, s.customer, s.vendor, collection.NewCoins(collection.NewFTCoin(s.ftClassID, s.balance))),
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
				helperBuildEventSent(s.contractID, s.customer, s.vendor, s.operator, collection.NewCoins(collection.NewFTCoin(s.ftClassID, s.balance))),
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

func helperBuildEventSent(contractID string, from, to, operator sdk.AccAddress, amount collection.Coins) sdk.Event {
	return sdk.Event{
		Type: "lbm.collection.v1.EventSent",
		Attributes: []abci.EventAttribute{
			{Key: []byte("amount"), Value: []byte(asJsonStr(amount)), Index: false},
			{Key: []byte("contract_id"), Value: []byte(wrapQuot(contractID)), Index: false},
			{Key: []byte("from"), Value: []byte(wrapQuot(from.String())), Index: false},
			{Key: []byte("operator"), Value: []byte(wrapQuot(operator.String())), Index: false},
			{Key: []byte("to"), Value: []byte(wrapQuot(to.String())), Index: false},
		},
	}
}

func TestHelperBuildEventSent(t *testing.T) {
	testCases := map[string]struct {
		expectedEvent sdk.Event
		contractID    string
		from          string
		to            string
		operator      string
		amount        collection.Coins
	}{
		"helper function should keep event compatibility": {
			expectedEvent: sdk.Event{Type: "lbm.collection.v1.EventSent", Attributes: []abci.EventAttribute{
				{Key: []byte("amount"), Value: []byte(`[{"token_id":"0000000100000000","amount":"1000000"}]`), Index: false},
				{Key: []byte("contract_id"), Value: []byte(`"9be17165"`), Index: false},
				{Key: []byte("from"), Value: []byte(`"link1v9jxgun9wdenqa2xzfx"`), Index: false},
				{Key: []byte("operator"), Value: []byte(`"link1v9jxgun9wdenqa2xzfx"`), Index: false},
				{Key: []byte("to"), Value: []byte(`"link1v9jxgun9wdenyjqyyxu"`), Index: false},
			}},
			contractID: "9be17165",
			from:       "link1v9jxgun9wdenqa2xzfx",
			to:         "link1v9jxgun9wdenyjqyyxu",
			operator:   "link1v9jxgun9wdenqa2xzfx",
			amount:     collection.NewCoins(collection.NewFTCoin("00000001", sdk.NewInt(1000000))),
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Act
			from, err := sdk.AccAddressFromBech32(tc.from)
			assert.NoError(t, err)
			to, err := sdk.AccAddressFromBech32(tc.to)
			assert.NoError(t, err)
			event := helperBuildEventSent(tc.contractID, from, to, from, tc.amount)

			// Assert
			assert.Equal(t, tc.expectedEvent, event)
		})
	}
}

func (s *KeeperTestSuite) TestMsgSendNFT() {
	testCases := map[string]struct {
		contractID string
		tokenID    string
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			tokenID:    collection.NewNFTID(s.nftClassID, 1),
			events:     sdk.Events{sdk.Event{Type: "lbm.collection.v1.EventOwnerChanged", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x66, 0x72, 0x6f, 0x6d}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0x22}, Index: false}}}, sdk.Event{Type: "lbm.collection.v1.EventOwnerChanged", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x66, 0x72, 0x6f, 0x6d}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x33, 0x22}, Index: false}}}, sdk.Event{Type: "lbm.collection.v1.EventOwnerChanged", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x66, 0x72, 0x6f, 0x6d}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x34, 0x22}, Index: false}}}, sdk.Event{Type: "lbm.collection.v1.EventSent", Attributes: []abci.EventAttribute{{Key: []uint8{0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74}, Value: []uint8{0x5b, 0x7b, 0x22, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x22, 0x3a, 0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x22, 0x2c, 0x22, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x3a, 0x22, 0x31, 0x22, 0x7d, 0x5d}, Index: false}, {Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x66, 0x72, 0x6f, 0x6d}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}, {Key: []uint8{0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}}}},
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

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgOperatorSendNFT() {
	tokenID := collection.NewNFTID(s.nftClassID, 1)
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
			events:     sdk.Events{sdk.Event{Type: "lbm.collection.v1.EventOwnerChanged", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x66, 0x72, 0x6f, 0x6d}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0x22}, Index: false}}}, sdk.Event{Type: "lbm.collection.v1.EventOwnerChanged", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x66, 0x72, 0x6f, 0x6d}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x33, 0x22}, Index: false}}}, sdk.Event{Type: "lbm.collection.v1.EventOwnerChanged", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x66, 0x72, 0x6f, 0x6d}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x34, 0x22}, Index: false}}}, sdk.Event{Type: "lbm.collection.v1.EventSent", Attributes: []abci.EventAttribute{{Key: []uint8{0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74}, Value: []uint8{0x5b, 0x7b, 0x22, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x22, 0x3a, 0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x22, 0x2c, 0x22, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x3a, 0x22, 0x31, 0x22, 0x7d, 0x5d}, Index: false}, {Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x66, 0x72, 0x6f, 0x6d}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}, {Key: []uint8{0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x7a, 0x77, 0x30, 0x38, 0x70, 0x36, 0x74, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}}}},
		},
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

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgAuthorizeOperator() {
	const authorizedOperatorEventType = "lbm.collection.v1.EventAuthorizedOperator"
	testCases := map[string]struct {
		isNegativeCase bool
		req            *collection.MsgAuthorizeOperator
		expectedEvents sdk.Events
		expectedError  error
	}{
		"valid request": {
			req: &collection.MsgAuthorizeOperator{
				ContractId: s.contractID,
				Holder:     s.customer.String(),
				Operator:   s.vendor.String(),
			},
			expectedEvents: sdk.Events{
				helperBuildEventAuthRelated(authorizedOperatorEventType, s.contractID, s.customer, s.vendor),
			},
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

			// Assert
			events := ctx.EventManager().Events()
			s.Require().Equal(tc.expectedEvents, events)
			curAuth, err := s.keeper.GetAuthorization(ctx, tc.req.ContractId, holder, operator)
			s.Require().NoError(err)
			s.Require().Nil(prevAuth)
			s.Require().Equal(tc.req.Holder, curAuth.Holder)
			s.Require().Equal(tc.req.Operator, curAuth.Operator)
		})
	}
}

func (s *KeeperTestSuite) TestMsgRevokeOperator() {
	const revokedOperatorEventType = "lbm.collection.v1.EventRevokedOperator"
	testCases := map[string]struct {
		isNegativeCase bool
		req            *collection.MsgRevokeOperator
		expectedEvents sdk.Events
		expectedError  error
	}{
		"valid request": {
			req: &collection.MsgRevokeOperator{
				ContractId: s.contractID,
				Holder:     s.customer.String(),
				Operator:   s.operator.String(),
			},
			expectedEvents: sdk.Events{
				helperBuildEventAuthRelated(revokedOperatorEventType, s.contractID, s.customer, s.operator),
			},
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

			// Assert
			events := ctx.EventManager().Events()
			s.Require().Equal(tc.expectedEvents, events)
			s.Require().NotNil(prevAuth)
			s.Require().Equal(tc.req.Holder, prevAuth.Holder)
			s.Require().Equal(tc.req.Operator, prevAuth.Operator)
			curAuth, err := s.keeper.GetAuthorization(ctx, tc.req.ContractId, holder, operator)
			s.Require().ErrorIs(err, collection.ErrCollectionNotApproved)
			s.Require().Nil(curAuth)
		})
	}
}

func helperBuildEventAuthRelated(evtType, contractID string, holder, operator sdk.AccAddress) sdk.Event {
	return sdk.Event{
		Type: evtType,
		Attributes: []abci.EventAttribute{
			{Key: []byte("contract_id"), Value: []byte(wrapQuot(contractID)), Index: false},
			{Key: []byte("holder"), Value: []byte(wrapQuot(holder.String())), Index: false},
			{Key: []byte("operator"), Value: []byte(wrapQuot(operator.String())), Index: false},
		},
	}
}

func TestHelperBuildAuthRelatedEvent(t *testing.T) {
	testCases := map[string]struct {
		expectedEvent sdk.Event
		contractID    string
		holder        string
		operator      string
	}{
		"helper function should keep event compatibility for approve": {
			expectedEvent: sdk.Event{Type: "lbm.collection.v1.EventSent", Attributes: []abci.EventAttribute{
				{Key: []byte("contract_id"), Value: []byte(`"9be17165"`), Index: false},
				{Key: []byte("holder"), Value: []byte(`"link1v9jxgun9wdenyjqyyxu"`), Index: false},
				{Key: []byte("operator"), Value: []byte(`"link1v9jxgun9wdenqa2xzfx"`), Index: false},
			}},
			contractID: "9be17165",
			holder:     "link1v9jxgun9wdenyjqyyxu",
			operator:   "link1v9jxgun9wdenqa2xzfx",
		},
		"helper function should keep event compatibility for revoke": {
			expectedEvent: sdk.Event{Type: "lbm.collection.v1.EventSent", Attributes: []abci.EventAttribute{
				{Key: []byte("contract_id"), Value: []byte(`"9be17165"`), Index: false},
				{Key: []byte("holder"), Value: []byte(`"link1v9jxgun9wdenyjqyyxu"`), Index: false},
				{Key: []byte("operator"), Value: []byte(`"link1v9jxgun9wdenzw08p6t"`), Index: false},
			}},
			contractID: "9be17165",
			holder:     "link1v9jxgun9wdenyjqyyxu",
			operator:   "link1v9jxgun9wdenzw08p6t",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Act
			holder, err := sdk.AccAddressFromBech32(tc.holder)
			assert.NoError(t, err)
			operator, err := sdk.AccAddressFromBech32(tc.operator)
			assert.NoError(t, err)
			event := helperBuildEventAuthRelated(tc.expectedEvent.Type, tc.contractID, holder, operator)

			// Assert
			assert.Equal(t, tc.expectedEvent, event)
		})
	}
}

func (s *KeeperTestSuite) TestMsgCreateContract() {
	testCases := map[string]struct {
		owner  sdk.AccAddress
		err    error
		events sdk.Events
	}{
		"valid request": {
			owner:  s.vendor,
			events: sdk.Events{sdk.Event{Type: "lbm.collection.v1.EventCreatedContract", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x33, 0x33, 0x33, 0x36, 0x62, 0x37, 0x36, 0x66, 0x22}, Index: false}, {Key: []uint8{0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}, {Key: []uint8{0x6d, 0x65, 0x74, 0x61}, Value: []uint8{0x22, 0x22}, Index: false}, {Key: []uint8{0x6e, 0x61, 0x6d, 0x65}, Value: []uint8{0x22, 0x22}, Index: false}, {Key: []uint8{0x75, 0x72, 0x69}, Value: []uint8{0x22, 0x22}, Index: false}}}, sdk.Event{Type: "lbm.collection.v1.EventGranted", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x33, 0x33, 0x33, 0x36, 0x62, 0x37, 0x36, 0x66, 0x22}, Index: false}, {Key: []uint8{0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x65}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}, {Key: []uint8{0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x72}, Value: []uint8{0x22, 0x22}, Index: false}, {Key: []uint8{0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e}, Value: []uint8{0x22, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x49, 0x53, 0x53, 0x55, 0x45, 0x22}, Index: false}}}, sdk.Event{Type: "lbm.collection.v1.EventGranted", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x33, 0x33, 0x33, 0x36, 0x62, 0x37, 0x36, 0x66, 0x22}, Index: false}, {Key: []uint8{0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x65}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}, {Key: []uint8{0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x72}, Value: []uint8{0x22, 0x22}, Index: false}, {Key: []uint8{0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e}, Value: []uint8{0x22, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x4d, 0x4f, 0x44, 0x49, 0x46, 0x59, 0x22}, Index: false}}}, sdk.Event{Type: "lbm.collection.v1.EventGranted", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x33, 0x33, 0x33, 0x36, 0x62, 0x37, 0x36, 0x66, 0x22}, Index: false}, {Key: []uint8{0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x65}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}, {Key: []uint8{0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x72}, Value: []uint8{0x22, 0x22}, Index: false}, {Key: []uint8{0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e}, Value: []uint8{0x22, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x4d, 0x49, 0x4e, 0x54, 0x22}, Index: false}}}, sdk.Event{Type: "lbm.collection.v1.EventGranted", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x33, 0x33, 0x33, 0x36, 0x62, 0x37, 0x36, 0x66, 0x22}, Index: false}, {Key: []uint8{0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x65}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}, {Key: []uint8{0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x72}, Value: []uint8{0x22, 0x22}, Index: false}, {Key: []uint8{0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e}, Value: []uint8{0x22, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x42, 0x55, 0x52, 0x4e, 0x22}, Index: false}}}},
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

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgIssueFT() {
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
						{Key: []uint8("contract_id"), Value: []uint8("\"9be17165\""), Index: false},
						{Key: []uint8("decimals"), Value: []uint8("0"), Index: false},
						{Key: []uint8("meta"), Value: []uint8("\"\""), Index: false},
						{Key: []uint8("mintable"), Value: []uint8("false"), Index: false},
						{Key: []uint8("name"), Value: []uint8("\"\""), Index: false},
						{Key: []uint8("operator"), Value: []uint8(fmt.Sprintf("\"%s\"", s.vendor)), Index: false},
						{Key: []uint8("token_id"), Value: []uint8("\"0000000200000000\""), Index: false},
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
						{Key: []uint8("contract_id"), Value: []uint8("\"9be17165\""), Index: false},
						{Key: []uint8("decimals"), Value: []uint8("0"), Index: false},
						{Key: []uint8("meta"), Value: []uint8("\"\""), Index: false},
						{Key: []uint8("mintable"), Value: []uint8("false"), Index: false},
						{Key: []uint8("name"), Value: []uint8("\"\""), Index: false},
						{Key: []uint8("operator"), Value: []uint8(fmt.Sprintf("\"%s\"", s.vendor)), Index: false},
						{Key: []uint8("token_id"), Value: []uint8("\"0000000200000000\""), Index: false},
					},
				},
				sdk.Event{
					Type: "lbm.collection.v1.EventMintedFT",
					Attributes: []abci.EventAttribute{
						{Key: []uint8("amount"), Value: []uint8("[{\"token_id\":\"0000000200000000\",\"amount\":\"1\"}]"), Index: false},
						{Key: []uint8("contract_id"), Value: []uint8("\"9be17165\""), Index: false},
						{Key: []uint8("operator"), Value: []uint8(fmt.Sprintf("\"%s\"", s.vendor)), Index: false},
						{Key: []uint8("to"), Value: []uint8(fmt.Sprintf("\"%s\"", s.customer)), Index: false},
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
	testCases := map[string]struct {
		contractID string
		owner      sdk.AccAddress
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			owner:      s.vendor,
			events:     sdk.Events{sdk.Event{Type: "lbm.collection.v1.EventCreatedNFTClass", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x6d, 0x65, 0x74, 0x61}, Value: []uint8{0x22, 0x22}, Index: false}, {Key: []uint8{0x6e, 0x61, 0x6d, 0x65}, Value: []uint8{0x22, 0x22}, Index: false}, {Key: []uint8{0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0x22}, Index: false}}}, sdk.Event{Type: "lbm.collection.v1.EventGranted", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x65}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}, {Key: []uint8{0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x72}, Value: []uint8{0x22, 0x22}, Index: false}, {Key: []uint8{0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e}, Value: []uint8{0x22, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x4d, 0x49, 0x4e, 0x54, 0x22}, Index: false}}}, sdk.Event{Type: "lbm.collection.v1.EventGranted", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x65}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}, {Key: []uint8{0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x72}, Value: []uint8{0x22, 0x22}, Index: false}, {Key: []uint8{0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e}, Value: []uint8{0x22, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x42, 0x55, 0x52, 0x4e, 0x22}, Index: false}}}},
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

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
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
					Type: "lbm.collection.v1.EventMintedFT",
					Attributes: []abci.EventAttribute{
						{Key: []uint8("amount"), Value: []uint8("[{\"token_id\":\"0000000100000000\",\"amount\":\"100000\"}]"), Index: false},
						{Key: []uint8("contract_id"), Value: []uint8("\"9be17165\""), Index: false},
						{Key: []uint8("operator"), Value: []uint8(fmt.Sprintf("\"%s\"", s.vendor)), Index: false},
						{Key: []uint8("to"), Value: []uint8(fmt.Sprintf("\"%s\"", s.customer)), Index: false},
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
		"valid request - empty amount": {
			contractID: s.contractID,
			from:       s.vendor,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventMintedFT",
					Attributes: []abci.EventAttribute{
						{Key: []uint8("amount"), Value: []uint8("[]"), Index: false},
						{Key: []uint8("contract_id"), Value: []uint8("\"9be17165\""), Index: false},
						{Key: []uint8("operator"), Value: []uint8(fmt.Sprintf("\"%s\"", s.vendor)), Index: false},
						{Key: []uint8("to"), Value: []uint8(fmt.Sprintf("\"%s\"", s.customer)), Index: false},
					},
				},
			},
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
		"valid request - multi tokens": {
			contractID: s.contractID,
			from:       s.vendor,
			amount: collection.NewCoins(
				collection.NewFTCoin(s.ftClassID, sdk.NewInt(100000)),
				collection.NewFTCoin(*mintableFTClassID, sdk.NewInt(200000)),
			),
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventMintedFT",
					Attributes: []abci.EventAttribute{
						{Key: []uint8("amount"), Value: []uint8("[{\"token_id\":\"0000000100000000\",\"amount\":\"100000\"},{\"token_id\":\"0000000200000000\",\"amount\":\"200000\"}]"), Index: false},
						{Key: []uint8("contract_id"), Value: []uint8("\"9be17165\""), Index: false},
						{Key: []uint8("operator"), Value: []uint8(fmt.Sprintf("\"%s\"", s.vendor)), Index: false},
						{Key: []uint8("to"), Value: []uint8(fmt.Sprintf("\"%s\"", s.customer)), Index: false},
					},
				},
			},
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
	}}
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
			events:     sdk.Events{sdk.Event{Type: "lbm.collection.v1.EventMintedNFT", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x73}, Value: []uint8{0x5b, 0x7b, 0x22, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x22, 0x3a, 0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x36, 0x22, 0x2c, 0x22, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x3a, 0x22, 0x22, 0x2c, 0x22, 0x6d, 0x65, 0x74, 0x61, 0x22, 0x3a, 0x22, 0x22, 0x7d, 0x5d}, Index: false}}}},
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
			res, err := s.msgServer.MintNFT(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgBurnFT() {
	// prepare mutli token burn test
	amount := collection.NewCoins(
		collection.NewFTCoin(s.ftClassID, s.balance),
	)

	// create a fungible token class
	mintableFTClassID, err := s.keeper.CreateTokenClass(s.ctx, s.contractID, &collection.FTClass{
		Name:     "tibetian fox2",
		Mintable: true,
	})
	s.Require().NoError(err)
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
			amount: collection.NewCoins(
				collection.NewFTCoin(s.ftClassID, sdk.NewInt(50000)),
			),
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventBurned",
					Attributes: []abci.EventAttribute{
						{Key: []uint8("amount"), Value: []uint8("[{\"token_id\":\"0000000100000000\",\"amount\":\"50000\"}]"), Index: false},
						{Key: []uint8("contract_id"), Value: []uint8("\"9be17165\""), Index: false},
						{Key: []uint8("from"), Value: []uint8(fmt.Sprintf("\"%s\"", s.vendor)), Index: false},
						{Key: []uint8("operator"), Value: []uint8(fmt.Sprintf("\"%s\"", s.vendor)), Index: false},
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
		"no amount - valid": {
			contractID: s.contractID,
			from:       s.vendor,
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventBurned",
					Attributes: []abci.EventAttribute{
						{Key: []uint8("amount"), Value: []uint8("[]"), Index: false},
						{Key: []uint8("contract_id"), Value: []uint8("\"9be17165\""), Index: false},
						{Key: []uint8("from"), Value: []uint8(fmt.Sprintf("\"%s\"", s.vendor)), Index: false},
						{Key: []uint8("operator"), Value: []uint8(fmt.Sprintf("\"%s\"", s.vendor)), Index: false},
					},
				},
			},
		},
		"valid multi amount burn": {
			contractID: s.contractID,
			from:       s.vendor,
			amount: collection.NewCoins(
				collection.NewFTCoin(s.ftClassID, sdk.NewInt(50000)),
				collection.NewFTCoin(*mintableFTClassID, sdk.NewInt(60000)),
			),
			events: sdk.Events{
				sdk.Event{
					Type: "lbm.collection.v1.EventBurned",
					Attributes: []abci.EventAttribute{
						{Key: []uint8("amount"), Value: []uint8("[{\"token_id\":\"0000000100000000\",\"amount\":\"50000\"},{\"token_id\":\"0000000200000000\",\"amount\":\"60000\"}]"), Index: false},
						{Key: []uint8("contract_id"), Value: []uint8("\"9be17165\""), Index: false},
						{Key: []uint8("from"), Value: []uint8(fmt.Sprintf("\"%s\"", s.vendor)), Index: false},
						{Key: []uint8("operator"), Value: []uint8(fmt.Sprintf("\"%s\"", s.vendor)), Index: false},
					},
				},
			},
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
			events:     sdk.Events{sdk.Event{Type: "lbm.collection.v1.EventBurned", Attributes: []abci.EventAttribute{{Key: []uint8{0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74}, Value: []uint8{0x5b, 0x7b, 0x22, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x22, 0x3a, 0x22, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x22, 0x2c, 0x22, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x3a, 0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x22, 0x7d, 0x5d}, Index: false}, {Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x66, 0x72, 0x6f, 0x6d}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}, {Key: []uint8{0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x7a, 0x77, 0x30, 0x38, 0x70, 0x36, 0x74, 0x22}, Index: false}}}},
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
		"insufficient funds": {
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

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
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
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			from:       s.vendor,
			tokenIDs:   tokenIDs,
			events:     sdk.Events{sdk.Event{Type: "lbm.collection.v1.EventBurned", Attributes: []abci.EventAttribute{{Key: []uint8{0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74}, Value: []uint8{0x5b, 0x7b, 0x22, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x22, 0x3a, 0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x66, 0x22, 0x2c, 0x22, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x3a, 0x22, 0x31, 0x22, 0x7d, 0x2c, 0x7b, 0x22, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x22, 0x3a, 0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x22, 0x2c, 0x22, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x3a, 0x22, 0x31, 0x22, 0x7d, 0x2c, 0x7b, 0x22, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x22, 0x3a, 0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x31, 0x22, 0x2c, 0x22, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x3a, 0x22, 0x31, 0x22, 0x7d, 0x2c, 0x7b, 0x22, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x22, 0x3a, 0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x32, 0x22, 0x2c, 0x22, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x3a, 0x22, 0x31, 0x22, 0x7d, 0x5d}, Index: false}, {Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x66, 0x72, 0x6f, 0x6d}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}, {Key: []uint8{0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}}}},
		},
		"contract not found": {
			contractID: "deadbeef",
			from:       s.vendor,
			tokenIDs:   tokenIDs,
			err:        class.ErrContractNotExist,
		},
		"no permission": {
			contractID: s.contractID,
			from:       s.customer,
			tokenIDs:   tokenIDs,
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

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgOperatorBurnNFT() {
	tokenIDs := []string{
		collection.NewNFTID(s.nftClassID, 1),
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
			tokenIDs:   tokenIDs,
			events:     sdk.Events{sdk.Event{Type: "lbm.collection.v1.EventBurned", Attributes: []abci.EventAttribute{{Key: []uint8{0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74}, Value: []uint8{0x5b, 0x7b, 0x22, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x22, 0x3a, 0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x22, 0x2c, 0x22, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x3a, 0x22, 0x31, 0x22, 0x7d, 0x2c, 0x7b, 0x22, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x22, 0x3a, 0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0x22, 0x2c, 0x22, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x3a, 0x22, 0x31, 0x22, 0x7d, 0x2c, 0x7b, 0x22, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x22, 0x3a, 0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x33, 0x22, 0x2c, 0x22, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x3a, 0x22, 0x31, 0x22, 0x7d, 0x2c, 0x7b, 0x22, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x22, 0x3a, 0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x34, 0x22, 0x2c, 0x22, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x3a, 0x22, 0x31, 0x22, 0x7d, 0x5d}, Index: false}, {Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x66, 0x72, 0x6f, 0x6d}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}, {Key: []uint8{0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x7a, 0x77, 0x30, 0x38, 0x70, 0x36, 0x74, 0x22}, Index: false}}}},
		},
		"contract not found": {
			contractID: "deadbeef",
			operator:   s.operator,
			from:       s.customer,
			tokenIDs:   tokenIDs,
			err:        class.ErrContractNotExist,
		},
		"no authorization": {
			contractID: s.contractID,
			operator:   s.vendor,
			from:       s.customer,
			tokenIDs:   tokenIDs,
			err:        collection.ErrCollectionNotApproved,
		},
		"no permission": {
			contractID: s.contractID,
			operator:   s.stranger,
			from:       s.customer,
			tokenIDs:   tokenIDs,
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

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
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
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.vendor,
			events:     sdk.Events{sdk.Event{Type: "lbm.collection.v1.EventModifiedContract", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73}, Value: []uint8{0x5b, 0x7b, 0x22, 0x6b, 0x65, 0x79, 0x22, 0x3a, 0x22, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x2c, 0x22, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x3a, 0x22, 0x74, 0x65, 0x73, 0x74, 0x22, 0x7d, 0x5d}, Index: false}, {Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}}}},
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

			changes := []collection.Attribute{{
				Key:   collection.AttributeKeyName.String(),
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

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
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
			events:     sdk.Events{sdk.Event{Type: "lbm.collection.v1.EventGranted", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x65}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x7a, 0x77, 0x30, 0x38, 0x70, 0x36, 0x74, 0x22}, Index: false}, {Key: []uint8{0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x71, 0x61, 0x32, 0x78, 0x7a, 0x66, 0x78, 0x22}, Index: false}, {Key: []uint8{0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e}, Value: []uint8{0x22, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x4d, 0x4f, 0x44, 0x49, 0x46, 0x59, 0x22}, Index: false}}}},
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

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
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
			events:     sdk.Events{sdk.Event{Type: "lbm.collection.v1.EventRenounced", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x65}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x7a, 0x77, 0x30, 0x38, 0x70, 0x36, 0x74, 0x22}, Index: false}, {Key: []uint8{0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e}, Value: []uint8{0x22, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x4d, 0x49, 0x4e, 0x54, 0x22}, Index: false}}}},
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

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
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
			events:     sdk.Events{sdk.Event{Type: "lbm.collection.v1.EventAttached", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x68, 0x6f, 0x6c, 0x64, 0x65, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}, {Key: []uint8{0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}, {Key: []uint8{0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x35, 0x22}, Index: false}, {Key: []uint8{0x74, 0x61, 0x72, 0x67, 0x65, 0x74}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x22}, Index: false}}}},
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

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgDetach() {
	testCases := map[string]struct {
		contractID string
		subjectID  string
		err        error
		events     sdk.Events
	}{
		"valid request": {
			contractID: s.contractID,
			subjectID:  collection.NewNFTID(s.nftClassID, 2),
			events:     sdk.Events{sdk.Event{Type: "lbm.collection.v1.EventDetached", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x68, 0x6f, 0x6c, 0x64, 0x65, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}, {Key: []uint8{0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}, {Key: []uint8{0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x22}, Index: false}, {Key: []uint8{0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0x22}, Index: false}}}, sdk.Event{Type: "lbm.collection.v1.EventRootChanged", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x66, 0x72, 0x6f, 0x6d}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x33, 0x22}, Index: false}}}, sdk.Event{Type: "lbm.collection.v1.EventRootChanged", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x66, 0x72, 0x6f, 0x6d}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x34, 0x22}, Index: false}}}},
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

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
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
			events:     sdk.Events{sdk.Event{Type: "lbm.collection.v1.EventAttached", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x68, 0x6f, 0x6c, 0x64, 0x65, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}, {Key: []uint8{0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x7a, 0x77, 0x30, 0x38, 0x70, 0x36, 0x74, 0x22}, Index: false}, {Key: []uint8{0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x35, 0x22}, Index: false}, {Key: []uint8{0x74, 0x61, 0x72, 0x67, 0x65, 0x74}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x22}, Index: false}}}},
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
			res, err := s.msgServer.OperatorAttach(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgOperatorDetach() {
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
			events:     sdk.Events{sdk.Event{Type: "lbm.collection.v1.EventDetached", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x68, 0x6f, 0x6c, 0x64, 0x65, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x79, 0x6a, 0x71, 0x79, 0x79, 0x78, 0x75, 0x22}, Index: false}, {Key: []uint8{0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72}, Value: []uint8{0x22, 0x6c, 0x69, 0x6e, 0x6b, 0x31, 0x76, 0x39, 0x6a, 0x78, 0x67, 0x75, 0x6e, 0x39, 0x77, 0x64, 0x65, 0x6e, 0x7a, 0x77, 0x30, 0x38, 0x70, 0x36, 0x74, 0x22}, Index: false}, {Key: []uint8{0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x22}, Index: false}, {Key: []uint8{0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0x22}, Index: false}}}, sdk.Event{Type: "lbm.collection.v1.EventRootChanged", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x66, 0x72, 0x6f, 0x6d}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x33, 0x22}, Index: false}}}, sdk.Event{Type: "lbm.collection.v1.EventRootChanged", Attributes: []abci.EventAttribute{{Key: []uint8{0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x39, 0x62, 0x65, 0x31, 0x37, 0x31, 0x36, 0x35, 0x22}, Index: false}, {Key: []uint8{0x66, 0x72, 0x6f, 0x6d}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0x22}, Index: false}, {Key: []uint8{0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64}, Value: []uint8{0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x34, 0x22}, Index: false}}}},
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
			res, err := s.msgServer.OperatorDetach(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)

			if s.deterministic {
				s.Require().Equal(tc.events, ctx.EventManager().Events())
			}
		})
	}
}

func asJsonStr(attrs collection.Coins) string {
	var buf strings.Builder
	enc := json.NewEncoder(&buf)
	enc.Encode(attrs)
	return strings.TrimSpace(buf.String())
}

// wrapQuot ("text") -> `"text"`
func wrapQuot(s string) string {
	return `"` + strings.TrimSpace(s) + `"`
}
