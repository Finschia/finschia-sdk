package token_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
)

func TestEventTypeStringer(t *testing.T) {
	for _, name := range token.EventType_name {
		t.Run(name, func(t *testing.T) {
			value := token.EventType(token.EventType_value[name])
			customName := value.String()
			require.EqualValues(t, value, token.EventTypeFromString(customName))
		})
	}
}

func TestAttributeKeyStringer(t *testing.T) {
	for _, name := range token.AttributeKey_name {
		t.Run(name, func(t *testing.T) {
			value := token.AttributeKey(token.AttributeKey_value[name])
			customName := value.String()
			require.EqualValues(t, value, token.AttributeKeyFromString(customName))
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

func TestNewEventIssueToken(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := token.EventIssued{
		Creator:    str(),
		ContractId: str(),
		Name:       str(),
		Symbol:     str(),
		Uri:        str(),
		Meta:       str(),
		Decimals:   0,
		Mintable:   true,
	}
	to := sdk.AccAddress(str())
	amount := sdk.OneInt()
	legacy := token.NewEventIssueToken(event, to, amount)

	require.Equal(t, token.EventTypeIssueToken.String(), legacy.Type)

	attributes := map[token.AttributeKey]string{
		token.AttributeKeyContractID: event.ContractId,
		token.AttributeKeyName:       event.Name,
		token.AttributeKeySymbol:     event.Symbol,
		token.AttributeKeyImageURI:   event.Uri,
		token.AttributeKeyMeta:       event.Meta,
		token.AttributeKeyMintable:   fmt.Sprintf("%v", event.Mintable),
		token.AttributeKeyDecimals:   fmt.Sprintf("%d", event.Decimals),
		token.AttributeKeyAmount:     amount.String(),
		token.AttributeKeyOwner:      event.Creator,
		token.AttributeKeyTo:         to.String(),
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventMintToken(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := token.EventMinted{
		ContractId: str(),
		Operator:   str(),
		To:         str(),
		Amount:     sdk.OneInt(),
	}
	legacy := token.NewEventMintToken(event)

	require.Equal(t, token.EventTypeMintToken.String(), legacy.Type)

	attributes := map[token.AttributeKey]string{
		token.AttributeKeyContractID: event.ContractId,
		token.AttributeKeyFrom:       event.Operator,
		token.AttributeKeyTo:         event.To,
		token.AttributeKeyAmount:     event.Amount.String(),
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventBurnToken(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	from := str()
	event := token.EventBurned{
		ContractId: str(),
		Operator:   from,
		From:       from,
		Amount:     sdk.OneInt(),
	}
	legacy := token.NewEventBurnToken(event)
	require.NotNil(t, legacy)

	require.Equal(t, token.EventTypeBurnToken.String(), legacy.Type)

	attributes := map[token.AttributeKey]string{
		token.AttributeKeyContractID: event.ContractId,
		token.AttributeKeyFrom:       event.From,
		token.AttributeKeyAmount:     event.Amount.String(),
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventBurnTokenFrom(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := token.EventBurned{
		ContractId: str(),
		Operator:   str(),
		From:       str(),
		Amount:     sdk.OneInt(),
	}
	legacy := token.NewEventBurnTokenFrom(event)
	require.NotNil(t, legacy)

	require.Equal(t, token.EventTypeBurnTokenFrom.String(), legacy.Type)

	attributes := map[token.AttributeKey]string{
		token.AttributeKeyContractID: event.ContractId,
		token.AttributeKeyProxy:      event.Operator,
		token.AttributeKeyFrom:       event.From,
		token.AttributeKeyAmount:     event.Amount.String(),
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventModifyToken(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := token.EventModified{
		ContractId: str(),
		Operator:   str(),
		Changes: []token.Pair{{
			Field: token.AttributeKeyName.String(),
			Value: str(),
		}},
	}
	legacies := token.NewEventModifyToken(event)
	require.Greater(t, len(legacies), 1)

	require.Equal(t, token.EventTypeModifyToken.String(), legacies[0].Type)
	attributes := map[token.AttributeKey]string{
		token.AttributeKeyContractID: event.ContractId,
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacies[0], key.String(), value), key)
	}

	for i, legacy := range legacies[1:] {
		require.Equal(t, token.EventTypeModifyToken.String(), legacy.Type)

		attributes := map[string]string{
			event.Changes[i].Field: event.Changes[i].Value,
		}
		for key, value := range attributes {
			require.True(t, assertAttribute(legacy, key, value), key)
		}
	}
}

func TestNewEventTransfer(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	from := str()
	event := token.EventSent{
		ContractId: str(),
		Operator:   from,
		From:       from,
		Amount:     sdk.OneInt(),
	}
	legacy := token.NewEventTransfer(event)
	require.NotNil(t, legacy)

	require.Equal(t, token.EventTypeTransfer.String(), legacy.Type)

	attributes := map[token.AttributeKey]string{
		token.AttributeKeyContractID: event.ContractId,
		token.AttributeKeyFrom:       event.From,
		token.AttributeKeyTo:         event.To,
		token.AttributeKeyAmount:     event.Amount.String(),
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventTransferFrom(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := token.EventSent{
		ContractId: str(),
		Operator:   str(),
		From:       str(),
		Amount:     sdk.OneInt(),
	}
	legacy := token.NewEventTransferFrom(event)
	require.NotNil(t, legacy)

	require.Equal(t, token.EventTypeTransferFrom.String(), legacy.Type)

	attributes := map[token.AttributeKey]string{
		token.AttributeKeyContractID: event.ContractId,
		token.AttributeKeyProxy:      event.Operator,
		token.AttributeKeyFrom:       event.From,
		token.AttributeKeyTo:         event.To,
		token.AttributeKeyAmount:     event.Amount.String(),
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventGrantPermToken(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }
	permission := func() token.Permission {
		n := len(token.Permission_value) - 1
		return token.Permission(1 + rand.Intn(n))
	}

	event := token.EventGranted{
		ContractId: str(),
		Granter:    str(),
		Grantee:    str(),
		Permission: permission(),
	}
	legacy := token.NewEventGrantPermToken(event)

	require.Equal(t, token.EventTypeGrantPermToken.String(), legacy.Type)

	attributes := map[token.AttributeKey]string{
		token.AttributeKeyContractID: event.ContractId,
		token.AttributeKeyFrom:       event.Granter,
		token.AttributeKeyTo:         event.Grantee,
		token.AttributeKeyPerm:       token.LegacyPermission(event.Permission).String(),
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventGrantPermTokenHead(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }
	permission := func() token.Permission {
		n := len(token.Permission_value) - 1
		return token.Permission(1 + rand.Intn(n))
	}

	event := token.EventGranted{
		ContractId: str(),
		Granter:    str(),
		Grantee:    str(),
		Permission: permission(),
	}
	legacy := token.NewEventGrantPermTokenHead(event)

	require.Equal(t, token.EventTypeGrantPermToken.String(), legacy.Type)

	attributes := map[token.AttributeKey]string{
		token.AttributeKeyContractID: event.ContractId,
		token.AttributeKeyTo:         event.Grantee,
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventGrantPermTokenBody(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }
	permission := func() token.Permission {
		n := len(token.Permission_value) - 1
		return token.Permission(1 + rand.Intn(n))
	}

	event := token.EventGranted{
		ContractId: str(),
		Granter:    str(),
		Grantee:    str(),
		Permission: permission(),
	}
	legacy := token.NewEventGrantPermTokenBody(event)

	require.Equal(t, token.EventTypeGrantPermToken.String(), legacy.Type)

	attributes := map[token.AttributeKey]string{
		token.AttributeKeyPerm: token.LegacyPermission(event.Permission).String(),
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventRevokePermToken(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }
	permission := func() token.Permission {
		n := len(token.Permission_value) - 1
		return token.Permission(1 + rand.Intn(n))
	}

	event := token.EventRenounced{
		ContractId: str(),
		Grantee:    str(),
		Permission: permission(),
	}
	legacy := token.NewEventRevokePermToken(event)

	require.Equal(t, token.EventTypeRevokePermToken.String(), legacy.Type)

	attributes := map[token.AttributeKey]string{
		token.AttributeKeyContractID: event.ContractId,
		token.AttributeKeyFrom:       event.Grantee,
		token.AttributeKeyPerm:       token.LegacyPermission(event.Permission).String(),
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}

func TestNewEventApproveToken(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	str := func() string { return randomString(8) }

	event := token.EventAuthorizedOperator{
		ContractId: str(),
		Holder:     str(),
		Operator:   str(),
	}
	legacy := token.NewEventApproveToken(event)

	require.Equal(t, token.EventTypeApproveToken.String(), legacy.Type)

	attributes := map[token.AttributeKey]string{
		token.AttributeKeyContractID: event.ContractId,
		token.AttributeKeyApprover:   event.Holder,
		token.AttributeKeyProxy:      event.Operator,
	}
	for key, value := range attributes {
		require.True(t, assertAttribute(legacy, key.String(), value), key)
	}
}
