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

func NewEventIssueToken(e EventIssue, grantee, to sdk.AccAddress, amount sdk.Int) sdk.Event {
	eventType := EventTypeIssueToken.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyName:       e.Name,
		AttributeKeySymbol:     e.Symbol,
		AttributeKeyImageURI:   e.Uri,
		AttributeKeyMeta:       e.Meta,
		AttributeKeyDecimals:   fmt.Sprintf("%d", e.Decimals),
		AttributeKeyMintable:   fmt.Sprintf("%t", e.Mintable),

		AttributeKeyOwner:  grantee.String(),
		AttributeKeyTo:     to.String(),
		AttributeKeyAmount: amount.String(),
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new.AppendAttributes(attribute)
	}
	return new
}

func NewEventMintToken(e EventMinted) sdk.Event {
	eventType := EventTypeMintToken.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyFrom:       e.Operator,
		AttributeKeyTo:         e.To,
		AttributeKeyAmount:     e.Amount.String(),
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new.AppendAttributes(attribute)
	}
	return new
}

func NewEventBurnToken(e EventBurned) sdk.Event {
	eventType := EventTypeBurnToken.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyFrom:       e.From,
		AttributeKeyAmount:     e.Amount.String(),
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new.AppendAttributes(attribute)
	}
	return new
}

func NewEventBurnTokenFrom(e EventBurned) sdk.Event {
	eventType := EventTypeBurnTokenFrom.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyProxy:      e.Operator,
		AttributeKeyFrom:       e.From,
		AttributeKeyAmount:     e.Amount.String(),
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new.AppendAttributes(attribute)
	}
	return new
}

func NewEventModifyToken(e EventModified) []sdk.Event {
	eventType := EventTypeModifyToken.String()
	new := []sdk.Event{
		sdk.NewEvent(eventType,
			sdk.NewAttribute(AttributeKeyContractID.String(), e.ContractId),
		),
	}

	for _, pair := range e.Changes {
		attribute := sdk.NewAttribute(pair.Field, pair.Value)
		event := sdk.NewEvent(eventType, attribute)
		new = append(new, event)
	}
	return new
}

func NewEventTransfer(e EventSent) sdk.Event {
	eventType := EventTypeTransfer.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyFrom:       e.From,
		AttributeKeyTo:         e.To,
		AttributeKeyAmount:     e.Amount.String(),
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new.AppendAttributes(attribute)
	}
	return new
}

func NewEventTransferFrom(e EventSent) sdk.Event {
	eventType := EventTypeTransferFrom.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyProxy:      e.Operator,
		AttributeKeyFrom:       e.From,
		AttributeKeyTo:         e.To,
		AttributeKeyAmount:     e.Amount.String(),
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new.AppendAttributes(attribute)
	}
	return new
}

func NewEventGrantPermToken(e EventGrant) sdk.Event {
	eventType := EventTypeGrantPermToken.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyTo:         e.Grantee,
		AttributeKeyPerm:       LegacyPermission(e.Permission).String(),
	}
	if e.Granter != e.Grantee {
		attributes[AttributeKeyFrom] = e.Granter
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new.AppendAttributes(attribute)
	}
	return new
}

func NewEventRevokePermToken(e EventAbandon) sdk.Event {
	eventType := EventTypeRevokePermToken.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyFrom:       e.Grantee,
		AttributeKeyPerm:       e.Permission,
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new.AppendAttributes(attribute)
	}
	return new
}

func NewEventApproveToken(e EventAuthorizedOperator) sdk.Event {
	eventType := EventTypeApproveToken.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyApprover:   e.Holder,
		AttributeKeyProxy:      e.Operator,
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new.AppendAttributes(attribute)
	}
	return new
}
