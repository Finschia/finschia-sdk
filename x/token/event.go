package token

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
)

func NewEventIssueToken(e EventIssue, grantee, to sdk.AccAddress, amount sdk.Int) sdk.Event {
	eventType := EventType_IssueToken.String()
	attributes := map[AttributeKey]string{
		AttributeKey_ContractID: e.ContractId,
		AttributeKey_Name:       e.Name,
		AttributeKey_Symbol:     e.Symbol,
		AttributeKey_ImageURI:   e.Uri,
		AttributeKey_Meta:       e.Meta,
		AttributeKey_Decimals:   fmt.Sprintf("%d", e.Decimals),
		AttributeKey_Mintable:   fmt.Sprintf("%t", e.Mintable),

		AttributeKey_Owner:  grantee.String(),
		AttributeKey_To:     to.String(),
		AttributeKey_Amount: amount.String(),
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new.AppendAttributes(attribute)
	}
	return new
}

func NewEventMintToken(e EventMinted) sdk.Event {
	eventType := EventType_MintToken.String()
	attributes := map[AttributeKey]string{
		AttributeKey_ContractID: e.ContractId,
		AttributeKey_From:       e.Operator,
		AttributeKey_To:         e.To,
		AttributeKey_Amount:     e.Amount.String(),
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new.AppendAttributes(attribute)
	}
	return new
}

func NewEventBurnToken(e EventBurned) sdk.Event {
	eventType := EventType_BurnToken.String()
	attributes := map[AttributeKey]string{
		AttributeKey_ContractID: e.ContractId,
		AttributeKey_From:       e.From,
		AttributeKey_Amount:     e.Amount.String(),
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new.AppendAttributes(attribute)
	}
	return new
}

func NewEventBurnTokenFrom(e EventBurned) sdk.Event {
	eventType := EventType_BurnTokenFrom.String()
	attributes := map[AttributeKey]string{
		AttributeKey_ContractID: e.ContractId,
		AttributeKey_Proxy:      e.Operator,
		AttributeKey_From:       e.From,
		AttributeKey_Amount:     e.Amount.String(),
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new.AppendAttributes(attribute)
	}
	return new
}

func NewEventModifyToken(e EventModified) []sdk.Event {
	eventType := EventType_ModifyToken.String()
	new := []sdk.Event{
		sdk.NewEvent(eventType,
			sdk.NewAttribute(AttributeKey_ContractID.String(), e.ContractId),
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
	eventType := EventType_Transfer.String()
	attributes := map[AttributeKey]string{
		AttributeKey_ContractID: e.ContractId,
		AttributeKey_From:       e.From,
		AttributeKey_To:         e.To,
		AttributeKey_Amount:     e.Amount.String(),
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new.AppendAttributes(attribute)
	}
	return new
}

func NewEventTransferFrom(e EventSent) sdk.Event {
	eventType := EventType_TransferFrom.String()
	attributes := map[AttributeKey]string{
		AttributeKey_ContractID: e.ContractId,
		AttributeKey_Proxy:      e.Operator,
		AttributeKey_From:       e.From,
		AttributeKey_To:         e.To,
		AttributeKey_Amount:     e.Amount.String(),
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new.AppendAttributes(attribute)
	}
	return new
}

func NewEventGrantPermToken(e EventGrant) sdk.Event {
	eventType := EventType_GrantPermToken.String()
	attributes := map[AttributeKey]string{
		AttributeKey_ContractID: e.ContractId,
		AttributeKey_To:         e.Grantee,
		AttributeKey_Perm:       e.Permission,
	}
	if e.Granter != e.Grantee {
		attributes[AttributeKey_From] = e.Granter
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new.AppendAttributes(attribute)
	}
	return new
}

func NewEventRevokePermToken(e EventAbandon) sdk.Event {
	eventType := EventType_RevokePermToken.String()
	attributes := map[AttributeKey]string{
		AttributeKey_ContractID: e.ContractId,
		AttributeKey_From:       e.Grantee,
		AttributeKey_Perm:       e.Permission,
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new.AppendAttributes(attribute)
	}
	return new
}

func NewEventApproveToken(e EventAuthorizedOperator) sdk.Event {
	eventType := EventType_ApproveToken.String()
	attributes := map[AttributeKey]string{
		AttributeKey_ContractID: e.ContractId,
		AttributeKey_Approver:   e.Holder,
		AttributeKey_Proxy:      e.Operator,
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new.AppendAttributes(attribute)
	}
	return new
}
