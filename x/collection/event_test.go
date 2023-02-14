package collection_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/collection"
)

func TestEventTypeStringer(t *testing.T) {
	for _, name := range collection.EventType_name {
		t.Run(name, func(t *testing.T) {
			value := collection.EventType(collection.EventType_value[name])
			customName := value.String()
			require.EqualValues(t, value, collection.EventTypeFromString(customName))
		})
	}
}

func TestAttributeKeyStringer(t *testing.T) {
	for _, name := range collection.AttributeKey_name {
		t.Run(name, func(t *testing.T) {
			value := collection.AttributeKey(collection.AttributeKey_value[name])
			customName := value.String()
			require.EqualValues(t, value, collection.AttributeKeyFromString(customName))
		})
	}
}

func randomString(length int) string {
	letters := []rune("0123456789abcdef")
	res := make([]rune, length)
	for i := range res {
		res[i] = letters[rand.Intn(len(letters))]
	}
	return string(res)
}

func assertAttribute(e sdk.Event, key, value string) bool {
	for _, attr := range e.Attributes {
		if string(attr.Key) == key {
			return string(attr.Value) == value
		}
	}
	return false
}

func TestNewEventCreateCollection(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := collection.EventCreatedContract{
		Creator:    str(),
		ContractId: str(),
		Name:       str(),
		Meta:       str(),
		Uri:        str(),
	}
	legacy := collection.NewEventCreateCollection(event)

	require.Equal(t, collection.EventTypeCreateCollection.String(), legacy.Type)

	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyName:       event.Name,
		collection.AttributeKeyMeta:       event.Meta,
		collection.AttributeKeyBaseImgURI: event.Uri,
		collection.AttributeKeyOwner:      event.Creator,
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventIssueFT(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := collection.EventCreatedFTClass{
		ContractId: str(),
		Operator:   str(),
		TokenId:    str(),
		Name:       str(),
		Meta:       str(),
		Decimals:   0,
		Mintable:   true,
	}
	to := sdk.AccAddress(str())
	amount := sdk.OneInt()
	legacy := collection.NewEventIssueFT(event, to, amount)

	require.Equal(t, collection.EventTypeIssueFT.String(), legacy.Type)

	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyTokenID:    event.TokenId,
		collection.AttributeKeyName:       event.Name,
		collection.AttributeKeyMeta:       event.Meta,
		collection.AttributeKeyMintable:   fmt.Sprintf("%v", event.Mintable),
		collection.AttributeKeyDecimals:   fmt.Sprintf("%d", event.Decimals),
		collection.AttributeKeyAmount:     amount.String(),
		collection.AttributeKeyOwner:      event.Operator,
		collection.AttributeKeyTo:         to.String(),
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventIssueNFT(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := collection.EventCreatedNFTClass{
		ContractId: str(),
		TokenType:  str(),
		Name:       str(),
		Meta:       str(),
	}
	legacy := collection.NewEventIssueNFT(event)

	require.Equal(t, collection.EventTypeIssueNFT.String(), legacy.Type)

	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyTokenType:  event.TokenType,
		collection.AttributeKeyName:       event.Name,
		collection.AttributeKeyMeta:       event.Meta,
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventMintFT(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := collection.EventMintedFT{
		ContractId: str(),
		Operator:   str(),
		To:         str(),
		Amount:     collection.NewCoins(collection.NewFTCoin(str(), sdk.OneInt())),
	}
	legacy := collection.NewEventMintFT(event)

	require.Equal(t, collection.EventTypeMintFT.String(), legacy.Type)

	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyFrom:       event.Operator,
		collection.AttributeKeyTo:         event.To,
		collection.AttributeKeyAmount:     collection.Coins(event.Amount).String(),
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventMintNFT(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := collection.EventMintedNFT{
		ContractId: str(),
		Operator:   str(),
		To:         str(),
		Tokens: []collection.NFT{{
			TokenId: str(),
			Name:    str(),
			Meta:    str(),
		}},
	}
	legacies := collection.NewEventMintNFT(event)

	for i, legacy := range legacies {
		require.Equal(t, collection.EventTypeMintNFT.String(), legacy.Type)

		attributes := map[collection.AttributeKey]string{
			collection.AttributeKeyContractID: event.ContractId,
			collection.AttributeKeyFrom:       event.Operator,
			collection.AttributeKeyTo:         event.To,
			collection.AttributeKeyTokenID:    event.Tokens[i].TokenId,
			collection.AttributeKeyName:       event.Tokens[i].Name,
			collection.AttributeKeyMeta:       event.Tokens[i].Meta,
		}
		for key, value := range attributes {
			require.True(t, assertAttribute(legacy, key.String(), value), key)
		}
	}
}

func TestNewEventBurnFT(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	from := str()
	event := collection.EventBurned{
		ContractId: str(),
		Operator:   from,
		From:       from,
		Amount: collection.NewCoins(
			collection.NewFTCoin(str(), sdk.OneInt()),
		),
	}
	legacy := collection.NewEventBurnFT(event)
	require.NotNil(t, legacy)

	require.Equal(t, collection.EventTypeBurnFT.String(), legacy.Type)

	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyFrom:       event.From,
		collection.AttributeKeyAmount:     collection.Coins(event.Amount).String(),
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(*legacy, key.String(), value), key)
	}

	empty := collection.NewEventBurnNFT(event)
	require.Empty(t, empty)
}

func TestNewEventBurnNFT(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	from := str()
	event := collection.EventBurned{
		ContractId: str(),
		Operator:   from,
		From:       from,
		Amount: collection.NewCoins(
			collection.NewNFTCoin(str(), 1),
		),
	}
	legacies := collection.NewEventBurnNFT(event)
	require.Greater(t, len(legacies), 1)

	require.Equal(t, collection.EventTypeBurnNFT.String(), legacies[0].Type)
	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyFrom:       event.From,
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacies[0], key.String(), value), key)
	}

	for i, legacy := range legacies[1:] {
		require.Equal(t, collection.EventTypeBurnNFT.String(), legacy.Type)

		attributes := map[collection.AttributeKey]string{
			collection.AttributeKeyTokenID: event.Amount[i].TokenId,
		}
		for key, value := range attributes {
			require.True(t, assertAttribute(legacy, key.String(), value), key)
		}
	}

	empty := collection.NewEventBurnFT(event)
	require.Nil(t, empty)
}

func TestNewEventBurnFTFrom(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := collection.EventBurned{
		ContractId: str(),
		Operator:   str(),
		From:       str(),
		Amount: collection.NewCoins(
			collection.NewFTCoin(str(), sdk.OneInt()),
		),
	}
	legacy := collection.NewEventBurnFTFrom(event)
	require.NotNil(t, legacy)

	require.Equal(t, collection.EventTypeBurnFTFrom.String(), legacy.Type)

	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyProxy:      event.Operator,
		collection.AttributeKeyFrom:       event.From,
		collection.AttributeKeyAmount:     collection.Coins(event.Amount).String(),
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(*legacy, key.String(), value), key)
	}

	empty := collection.NewEventBurnNFTFrom(event)
	require.Empty(t, empty)
}

func TestNewEventBurnNFTFrom(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := collection.EventBurned{
		ContractId: str(),
		Operator:   str(),
		From:       str(),
		Amount: collection.NewCoins(
			collection.NewNFTCoin(str(), 1),
		),
	}
	legacies := collection.NewEventBurnNFTFrom(event)
	require.Greater(t, len(legacies), 1)

	require.Equal(t, collection.EventTypeBurnNFTFrom.String(), legacies[0].Type)
	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyProxy:      event.Operator,
		collection.AttributeKeyFrom:       event.From,
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacies[0], key.String(), value), key)
	}

	for i, legacy := range legacies[1:] {
		require.Equal(t, collection.EventTypeBurnNFTFrom.String(), legacy.Type)

		attributes := map[collection.AttributeKey]string{
			collection.AttributeKeyTokenID: event.Amount[i].TokenId,
		}
		for key, value := range attributes {
			require.True(t, assertAttribute(legacy, key.String(), value), key)
		}
	}

	empty := collection.NewEventBurnFTFrom(event)
	require.Nil(t, empty)
}

func TestNewEventModifyCollection(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := collection.EventModifiedContract{
		ContractId: str(),
		Operator:   str(),
		Changes: []collection.Attribute{
			{
				Key:   collection.AttributeKeyName.String(),
				Value: str(),
			},
			{
				Key:   collection.AttributeKeyBaseImgURI.String(),
				Value: str(),
			},
		},
	}
	collection.UpdateEventModifiedContract(&event)

	legacies := collection.NewEventModifyCollection(event)
	require.Greater(t, len(legacies), 1)

	require.Equal(t, collection.EventTypeModifyCollection.String(), legacies[0].Type)
	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacies[0], key.String(), value), key)
	}

	for i, legacy := range legacies[1:] {
		require.Equal(t, collection.EventTypeModifyCollection.String(), legacy.Type)

		attributes := map[string]string{
			event.Changes[i].Key: event.Changes[i].Value,
		}
		for key, value := range attributes {
			require.True(t, assertAttribute(legacy, key, value), key)
		}
	}
}

func TestNewEventModifyTokenType(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := collection.EventModifiedTokenClass{
		ContractId: str(),
		Operator:   str(),
		TokenType:  str(),
		Changes: []collection.Attribute{{
			Key:   collection.AttributeKeyName.String(),
			Value: str(),
		}},
	}
	legacies := collection.NewEventModifyTokenType(event)
	require.Greater(t, len(legacies), 1)

	require.Equal(t, collection.EventTypeModifyTokenType.String(), legacies[0].Type)
	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyTokenType:  event.TokenType,
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacies[0], key.String(), value), key)
	}

	for i, legacy := range legacies[1:] {
		require.Equal(t, collection.EventTypeModifyTokenType.String(), legacy.Type)

		attributes := map[string]string{
			event.Changes[i].Key: event.Changes[i].Value,
		}
		for key, value := range attributes {
			require.True(t, assertAttribute(legacy, key, value), key)
		}
	}
}

func TestNewEventModifyTokenOfFTClass(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := collection.EventModifiedTokenClass{
		ContractId: str(),
		Operator:   str(),
		TokenType:  str(),
		Changes: []collection.Attribute{{
			Key:   collection.AttributeKeyName.String(),
			Value: str(),
		}},
	}
	legacies := collection.NewEventModifyTokenOfFTClass(event)
	require.Greater(t, len(legacies), 1)

	require.Equal(t, collection.EventTypeModifyToken.String(), legacies[0].Type)
	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyTokenID:    collection.NewFTID(event.TokenType),
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacies[0], key.String(), value), key)
	}

	for i, legacy := range legacies[1:] {
		require.Equal(t, collection.EventTypeModifyToken.String(), legacy.Type)

		attributes := map[string]string{
			event.Changes[i].Key: event.Changes[i].Value,
		}
		for key, value := range attributes {
			require.True(t, assertAttribute(legacy, key, value), key)
		}
	}
}

func TestNewEventModifyTokenOfNFT(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := collection.EventModifiedNFT{
		ContractId: str(),
		Operator:   str(),
		TokenId:    str(),
		Changes: []collection.Attribute{{
			Key:   collection.AttributeKeyName.String(),
			Value: str(),
		}},
	}
	legacies := collection.NewEventModifyTokenOfNFT(event)
	require.Greater(t, len(legacies), 1)

	require.Equal(t, collection.EventTypeModifyToken.String(), legacies[0].Type)
	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyTokenID:    event.TokenId,
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacies[0], key.String(), value), key)
	}

	for i, legacy := range legacies[1:] {
		require.Equal(t, collection.EventTypeModifyToken.String(), legacy.Type)

		attributes := map[string]string{
			event.Changes[i].Key: event.Changes[i].Value,
		}
		for key, value := range attributes {
			require.True(t, assertAttribute(legacy, key, value), key)
		}
	}
}

func TestNewEventTransferFT(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	from := str()
	event := collection.EventSent{
		ContractId: str(),
		Operator:   from,
		From:       from,
		Amount: collection.NewCoins(
			collection.NewFTCoin(str(), sdk.OneInt()),
		),
	}
	legacy := collection.NewEventTransferFT(event)
	require.NotNil(t, legacy)

	require.Equal(t, collection.EventTypeTransferFT.String(), legacy.Type)

	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyFrom:       event.From,
		collection.AttributeKeyTo:         event.To,
		collection.AttributeKeyAmount:     collection.Coins(event.Amount).String(),
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(*legacy, key.String(), value), key)
	}

	empty := collection.NewEventTransferNFT(event)
	require.Empty(t, empty)
}

func TestNewEventTransferNFT(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	from := str()
	event := collection.EventSent{
		ContractId: str(),
		Operator:   from,
		From:       from,
		Amount: collection.NewCoins(
			collection.NewNFTCoin(str(), 1),
		),
	}
	legacies := collection.NewEventTransferNFT(event)
	require.Greater(t, len(legacies), 1)

	require.Equal(t, collection.EventTypeTransferNFT.String(), legacies[0].Type)
	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyFrom:       event.From,
		collection.AttributeKeyTo:         event.To,
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacies[0], key.String(), value), key)
	}

	for i, legacy := range legacies[1:] {
		require.Equal(t, collection.EventTypeTransferNFT.String(), legacy.Type)

		attributes := map[collection.AttributeKey]string{
			collection.AttributeKeyTokenID: event.Amount[i].TokenId,
		}
		for key, value := range attributes {
			require.True(t, assertAttribute(legacy, key.String(), value), key)
		}
	}

	empty := collection.NewEventTransferFT(event)
	require.Nil(t, empty)
}

func TestNewEventTransferFTFrom(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := collection.EventSent{
		ContractId: str(),
		Operator:   str(),
		From:       str(),
		Amount: collection.NewCoins(
			collection.NewFTCoin(str(), sdk.OneInt()),
		),
	}
	legacy := collection.NewEventTransferFTFrom(event)
	require.NotNil(t, legacy)

	require.Equal(t, collection.EventTypeTransferFTFrom.String(), legacy.Type)

	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyProxy:      event.Operator,
		collection.AttributeKeyFrom:       event.From,
		collection.AttributeKeyTo:         event.To,
		collection.AttributeKeyAmount:     collection.Coins(event.Amount).String(),
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(*legacy, key.String(), value), key)
	}

	empty := collection.NewEventTransferNFTFrom(event)
	require.Empty(t, empty)
}

func TestNewEventTransferNFTFrom(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := collection.EventSent{
		ContractId: str(),
		Operator:   str(),
		From:       str(),
		Amount: collection.NewCoins(
			collection.NewNFTCoin(str(), 1),
		),
	}
	legacies := collection.NewEventTransferNFTFrom(event)
	require.Greater(t, len(legacies), 1)

	require.Equal(t, collection.EventTypeTransferNFTFrom.String(), legacies[0].Type)
	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyProxy:      event.Operator,
		collection.AttributeKeyFrom:       event.From,
		collection.AttributeKeyTo:         event.To,
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacies[0], key.String(), value), key)
	}

	for i, legacy := range legacies[1:] {
		require.Equal(t, collection.EventTypeTransferNFTFrom.String(), legacy.Type)

		attributes := map[collection.AttributeKey]string{
			collection.AttributeKeyTokenID: event.Amount[i].TokenId,
		}
		for key, value := range attributes {
			require.True(t, assertAttribute(legacy, key.String(), value), key)
		}
	}

	empty := collection.NewEventTransferFTFrom(event)
	require.Nil(t, empty)
}

func TestNewEventGrantPermToken(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }
	permission := func() collection.Permission {
		n := len(collection.Permission_value) - 1
		return collection.Permission(1 + rand.Intn(n))
	}

	event := collection.EventGranted{
		ContractId: str(),
		Granter:    str(),
		Grantee:    str(),
		Permission: permission(),
	}
	legacy := collection.NewEventGrantPermToken(event)

	require.Equal(t, collection.EventTypeGrantPermToken.String(), legacy.Type)

	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyFrom:       event.Granter,
		collection.AttributeKeyTo:         event.Grantee,
		collection.AttributeKeyPerm:       collection.LegacyPermission(event.Permission).String(),
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventGrantPermTokenHead(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }
	permission := func() collection.Permission {
		n := len(collection.Permission_value) - 1
		return collection.Permission(1 + rand.Intn(n))
	}

	event := collection.EventGranted{
		ContractId: str(),
		Granter:    str(),
		Grantee:    str(),
		Permission: permission(),
	}
	legacy := collection.NewEventGrantPermTokenHead(event)

	require.Equal(t, collection.EventTypeGrantPermToken.String(), legacy.Type)

	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyTo:         event.Grantee,
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventGrantPermTokenBody(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }
	permission := func() collection.Permission {
		n := len(collection.Permission_value) - 1
		return collection.Permission(1 + rand.Intn(n))
	}

	event := collection.EventGranted{
		ContractId: str(),
		Granter:    str(),
		Grantee:    str(),
		Permission: permission(),
	}
	legacy := collection.NewEventGrantPermTokenBody(event)

	require.Equal(t, collection.EventTypeGrantPermToken.String(), legacy.Type)

	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyPerm: collection.LegacyPermission(event.Permission).String(),
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventRevokePermToken(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }
	permission := func() collection.Permission {
		n := len(collection.Permission_value) - 1
		return collection.Permission(1 + rand.Intn(n))
	}

	event := collection.EventRenounced{
		ContractId: str(),
		Grantee:    str(),
		Permission: permission(),
	}
	legacy := collection.NewEventRevokePermToken(event)

	require.Equal(t, collection.EventTypeRevokePermToken.String(), legacy.Type)

	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyFrom:       event.Grantee,
		collection.AttributeKeyPerm:       collection.LegacyPermission(event.Permission).String(),
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventApproveCollection(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := collection.EventAuthorizedOperator{
		ContractId: str(),
		Holder:     str(),
		Operator:   str(),
	}
	legacy := collection.NewEventApproveCollection(event)

	require.Equal(t, collection.EventTypeApproveCollection.String(), legacy.Type)

	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyApprover:   event.Holder,
		collection.AttributeKeyProxy:      event.Operator,
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventDisapproveCollection(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := collection.EventRevokedOperator{
		ContractId: str(),
		Holder:     str(),
		Operator:   str(),
	}
	legacy := collection.NewEventDisapproveCollection(event)

	require.Equal(t, collection.EventTypeDisapproveCollection.String(), legacy.Type)

	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyApprover:   event.Holder,
		collection.AttributeKeyProxy:      event.Operator,
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventAttachToken(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := collection.EventAttached{
		ContractId: str(),
		Operator:   str(),
		Holder:     str(),
		Subject:    collection.NewNFTID(str(), 1),
		Target:     collection.NewNFTID(str(), 2),
	}
	newRoot := collection.NewNFTID(str(), 3)
	legacy := collection.NewEventAttachToken(event, newRoot)

	require.Equal(t, collection.EventTypeAttachToken.String(), legacy.Type)

	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyFrom:       event.Holder,
		collection.AttributeKeyTokenID:    event.Subject,
		collection.AttributeKeyToTokenID:  event.Target,
		collection.AttributeKeyOldRoot:    event.Subject,
		collection.AttributeKeyNewRoot:    newRoot,
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventDetachToken(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := collection.EventDetached{
		ContractId: str(),
		Operator:   str(),
		Holder:     str(),
		Subject:    collection.NewNFTID(str(), 1),
	}
	oldRoot := collection.NewNFTID(str(), 3)
	legacy := collection.NewEventDetachToken(event, oldRoot)

	require.Equal(t, collection.EventTypeDetachToken.String(), legacy.Type)

	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID:  event.ContractId,
		collection.AttributeKeyFrom:        event.Holder,
		collection.AttributeKeyFromTokenID: event.Subject,
		collection.AttributeKeyOldRoot:     oldRoot,
		collection.AttributeKeyNewRoot:     event.Subject,
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventAttachFrom(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := collection.EventAttached{
		ContractId: str(),
		Operator:   str(),
		Holder:     str(),
		Subject:    collection.NewNFTID(str(), 1),
		Target:     collection.NewNFTID(str(), 2),
	}
	newRoot := collection.NewNFTID(str(), 3)
	legacy := collection.NewEventAttachFrom(event, newRoot)

	require.Equal(t, collection.EventTypeAttachFrom.String(), legacy.Type)

	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyFrom:       event.Holder,
		collection.AttributeKeyTokenID:    event.Subject,
		collection.AttributeKeyToTokenID:  event.Target,
		collection.AttributeKeyOldRoot:    event.Subject,
		collection.AttributeKeyNewRoot:    newRoot,
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventDetachFrom(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := collection.EventDetached{
		ContractId: str(),
		Operator:   str(),
		Holder:     str(),
		Subject:    collection.NewNFTID(str(), 1),
	}
	oldRoot := collection.NewNFTID(str(), 3)
	legacy := collection.NewEventDetachFrom(event, oldRoot)

	require.Equal(t, collection.EventTypeDetachFrom.String(), legacy.Type)

	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID:  event.ContractId,
		collection.AttributeKeyFrom:        event.Holder,
		collection.AttributeKeyFromTokenID: event.Subject,
		collection.AttributeKeyOldRoot:     oldRoot,
		collection.AttributeKeyNewRoot:     event.Subject,
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventOperationTransferNFT(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := collection.EventOwnerChanged{
		ContractId: str(),
		TokenId:    str(),
	}
	legacy := collection.NewEventOperationTransferNFT(event)

	require.Equal(t, collection.EventTypeOperationTransferNFT.String(), legacy.Type)

	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyTokenID:    event.TokenId,
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventOperationBurnNFT(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	contractID := str()
	tokenID := str()
	legacy := collection.NewEventOperationBurnNFT(contractID, tokenID)

	require.Equal(t, collection.EventTypeOperationBurnNFT.String(), legacy.Type)

	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: contractID,
		collection.AttributeKeyTokenID:    tokenID,
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventOperationRootChanged(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := collection.EventRootChanged{
		ContractId: str(),
		TokenId:    str(),
	}
	legacy := collection.NewEventOperationRootChanged(event)

	require.Equal(t, collection.EventTypeOperationRootChanged.String(), legacy.Type)

	attributes := map[collection.AttributeKey]string{
		collection.AttributeKeyContractID: event.ContractId,
		collection.AttributeKeyTokenID:    event.TokenId,
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestUpdateModifiedContract(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := collection.EventModifiedContract{
		ContractId: str(),
		Operator:   str(),
		Changes: []collection.Attribute{
			{
				Key:   collection.AttributeKeyName.String(),
				Value: str(),
			},
			{
				Key:   collection.AttributeKeyBaseImgURI.String(),
				Value: str(),
			},
		},
	}
	collection.UpdateEventModifiedContract(&event)

	newChange := event.Changes[len(event.Changes) - 1]
	require.Equal(t, collection.AttributeKeyURI.String(), newChange.Key)
}
