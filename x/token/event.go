package token

import (
	"fmt"
	"strings"

	sdk "github.com/line/lbm-sdk/types"
)

const (
	prefixEventType    = "EVENT_TYPE_"
	prefixAttributeKey = "ATTRIBUTE_KEY_"
)

func (x EventType) String() string {
	lenPrefix := len(prefixEventType)
	return strings.ToLower(EventType_name[int32(x)][lenPrefix:])
}

func EventTypeFromString(name string) EventType {
	eventTypeName := prefixEventType + strings.ToUpper(name)
	return EventType(EventType_value[eventTypeName])
}

func (x AttributeKey) String() string {
	lenPrefix := len(prefixAttributeKey)
	return strings.ToLower(AttributeKey_name[int32(x)][lenPrefix:])
}

func AttributeKeyFromString(name string) AttributeKey {
	attributeKeyName := prefixAttributeKey + strings.ToUpper(name)
	return AttributeKey(AttributeKey_value[attributeKeyName])
}

func NewEventIssueToken(event EventIssue, to sdk.AccAddress, amount sdk.Int) sdk.Event {
	eventType := EventTypeIssueToken.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyName:       event.Name,
		AttributeKeySymbol:     event.Symbol,
		AttributeKeyImageURI:   event.Uri,
		AttributeKeyMeta:       event.Meta,
		AttributeKeyDecimals:   fmt.Sprintf("%d", event.Decimals),
		AttributeKeyMintable:   fmt.Sprintf("%t", event.Mintable),

		AttributeKeyOwner:  event.Creator,
		AttributeKeyTo:     to.String(),
		AttributeKeyAmount: amount.String(),
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return res
}

func NewEventMintToken(event EventMinted) sdk.Event {
	eventType := EventTypeMintToken.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyFrom:       event.Operator,
		AttributeKeyTo:         event.To,
		AttributeKeyAmount:     event.Amount.String(),
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return res
}

func NewEventBurnToken(event EventBurned) sdk.Event {
	eventType := EventTypeBurnToken.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyFrom:       event.From,
		AttributeKeyAmount:     event.Amount.String(),
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return res
}

func NewEventBurnTokenFrom(event EventBurned) sdk.Event {
	eventType := EventTypeBurnTokenFrom.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyProxy:      event.Operator,
		AttributeKeyFrom:       event.From,
		AttributeKeyAmount:     event.Amount.String(),
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return res
}

func NewEventModifyToken(event EventModified) []sdk.Event {
	eventType := EventTypeModifyToken.String()
	res := []sdk.Event{
		sdk.NewEvent(eventType,
			sdk.NewAttribute(AttributeKeyContractID.String(), event.ContractId),
		),
	}

	for _, pair := range event.Changes {
		attribute := sdk.NewAttribute(pair.Field, pair.Value)
		event := sdk.NewEvent(eventType, attribute)
		res = append(res, event)
	}
	return res
}

func NewEventTransfer(event EventSent) sdk.Event {
	eventType := EventTypeTransfer.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyFrom:       event.From,
		AttributeKeyTo:         event.To,
		AttributeKeyAmount:     event.Amount.String(),
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return res
}

func NewEventTransferFrom(event EventSent) sdk.Event {
	eventType := EventTypeTransferFrom.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyProxy:      event.Operator,
		AttributeKeyFrom:       event.From,
		AttributeKeyTo:         event.To,
		AttributeKeyAmount:     event.Amount.String(),
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return res
}

func NewEventGrantPermToken(event EventGrant) sdk.Event {
	eventType := EventTypeGrantPermToken.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyTo:         event.Grantee,
		AttributeKeyPerm:       LegacyPermission(event.Permission).String(),
	}
	if len(event.Granter) != 0 {
		attributes[AttributeKeyFrom] = event.Granter
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return res
}

func NewEventGrantPermTokenHead(event EventGrant) sdk.Event {
	eventType := EventTypeGrantPermToken.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyTo:         event.Grantee,
	}
	if len(event.Granter) != 0 {
		attributes[AttributeKeyFrom] = event.Granter
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return res
}

func NewEventGrantPermTokenBody(event EventGrant) sdk.Event {
	eventType := EventTypeGrantPermToken.String()
	attributes := map[AttributeKey]string{
		AttributeKeyPerm: LegacyPermission(event.Permission).String(),
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return res
}

func NewEventRevokePermToken(event EventAbandon) sdk.Event {
	eventType := EventTypeRevokePermToken.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyFrom:       event.Grantee,
		AttributeKeyPerm:       LegacyPermission(event.Permission).String(),
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return res
}

func NewEventApproveToken(event EventAuthorizedOperator) sdk.Event {
	eventType := EventTypeApproveToken.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyApprover:   event.Holder,
		AttributeKeyProxy:      event.Operator,
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return res
}
