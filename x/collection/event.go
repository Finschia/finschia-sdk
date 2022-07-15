package collection

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

// Deprecated: use typed events.
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

// Deprecated: use EventCreatedContract.
func NewEventCreateCollection(e EventCreatedContract, creator sdk.AccAddress) sdk.Event {
	eventType := EventTypeCreateCollection.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyName:       e.Name,
		AttributeKeyMeta:       e.Meta,
		AttributeKeyBaseImgURI: e.BaseImgUri,

		AttributeKeyOwner: creator.String(),
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new = new.AppendAttributes(attribute)
	}
	return new
}

// Deprecated: use EventCreatedFTClass.
func NewEventIssueFT(e EventCreatedFTClass, operator, to sdk.AccAddress, amount sdk.Int) sdk.Event {
	eventType := EventTypeIssueFT.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyTokenID:    NewFTID(e.ClassId),
		AttributeKeyName:       e.Name,
		AttributeKeyMeta:       e.Meta,
		AttributeKeyDecimals:   fmt.Sprintf("%d", e.Decimals),
		AttributeKeyMintable:   fmt.Sprintf("%t", e.Mintable),

		AttributeKeyOwner:  operator.String(),
		AttributeKeyTo:     to.String(),
		AttributeKeyAmount: amount.String(),
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new = new.AppendAttributes(attribute)
	}
	return new
}

// Deprecated: use EventCreatedNFTClass.
func NewEventIssueNFT(e EventCreatedNFTClass) sdk.Event {
	eventType := EventTypeIssueNFT.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyTokenType:  e.ClassId,
		AttributeKeyName:       e.Name,
		AttributeKeyMeta:       e.Meta,
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new = new.AppendAttributes(attribute)
	}
	return new
}

// Deprecated: use EventMintedFT.
func NewEventMintFT(e EventMintedFT) sdk.Event {
	eventType := EventTypeMintFT.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyFrom:       e.Operator,
		AttributeKeyTo:         e.To,
		AttributeKeyAmount:     Coins(e.Amount).String(),
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new = new.AppendAttributes(attribute)
	}
	return new
}

// Deprecated: use EventMintedNFT.
func NewEventMintNFT(e EventMintedNFT) sdk.Events {
	eventType := EventTypeMintNFT.String()

	new := make(sdk.Events, len(e.Tokens))
	for i, token := range e.Tokens {
		event := sdk.NewEvent(eventType)
		attributes := map[AttributeKey]string{
			AttributeKeyContractID: e.ContractId,
			AttributeKeyFrom:       e.Operator,
			AttributeKeyTo:         e.To,

			AttributeKeyTokenID: token.Id,
			AttributeKeyName:    token.Name,
			AttributeKeyMeta:    token.Meta,
		}
		for key, value := range attributes {
			attribute := sdk.NewAttribute(key.String(), value)
			event = event.AppendAttributes(attribute)
		}

		new[i] = event
	}

	return new
}

// Deprecated: use EventBurned.
func NewEventBurnFT(e EventBurned) *sdk.Event {
	eventType := EventTypeBurnFT.String()

	amount := []Coin{}
	for _, coin := range e.Amount {
		if err := ValidateFTID(coin.TokenId); err == nil {
			amount = append(amount, coin)
		}
	}
	if len(amount) == 0 {
		return nil
	}

	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyFrom:       e.From,
		AttributeKeyAmount:     Coins(amount).String(),
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new = new.AppendAttributes(attribute)
	}
	return &new
}

// Deprecated: use EventBurned.
func NewEventBurnNFT(e EventBurned) sdk.Events {
	eventType := EventTypeBurnNFT.String()

	amount := []Coin{}
	for _, coin := range e.Amount {
		if err := ValidateNFTID(coin.TokenId); err == nil {
			amount = append(amount, coin)
		}
	}
	if len(amount) == 0 {
		return nil
	}

	new := make(sdk.Events, 0, len(amount)+1)

	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyFrom:       e.From,
	}
	head := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		head = head.AppendAttributes(attribute)
	}
	new = append(new, head)

	for _, coin := range amount {
		attribute := sdk.NewAttribute(AttributeKeyTokenID.String(), coin.TokenId)
		event := sdk.NewEvent(eventType, attribute)
		new = append(new, event)
	}

	return new
}

// Deprecated: use EventBurned.
func NewEventBurnFTFrom(e EventBurned) *sdk.Event {
	eventType := EventTypeBurnFTFrom.String()

	amount := []Coin{}
	for _, coin := range e.Amount {
		if err := ValidateFTID(coin.TokenId); err == nil {
			amount = append(amount, coin)
		}
	}
	if len(amount) == 0 {
		return nil
	}

	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyProxy:      e.Operator,
		AttributeKeyFrom:       e.From,
		AttributeKeyAmount:     Coins(amount).String(),
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new = new.AppendAttributes(attribute)
	}
	return &new
}

// Deprecated: use EventBurned.
func NewEventBurnNFTFrom(e EventBurned) sdk.Events {
	eventType := EventTypeBurnNFTFrom.String()

	amount := []Coin{}
	for _, coin := range e.Amount {
		if err := ValidateNFTID(coin.TokenId); err == nil {
			amount = append(amount, coin)
		}
	}
	if len(amount) == 0 {
		return nil
	}

	new := make(sdk.Events, 0, len(amount)+1)

	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyProxy:      e.Operator,
		AttributeKeyFrom:       e.From,
	}
	head := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		head = head.AppendAttributes(attribute)
	}
	new = append(new, head)

	for _, coin := range amount {
		attribute := sdk.NewAttribute(AttributeKeyTokenID.String(), coin.TokenId)
		event := sdk.NewEvent(eventType, attribute)
		new = append(new, event)
	}

	return new
}

// Deprecated: use EventModifiedContract
func NewEventModifyCollection(e EventModifiedContract) sdk.Events {
	eventType := EventTypeModifyCollection.String()
	new := make(sdk.Events, 0, 1+len(e.Changes))

	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
	}
	head := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		head = head.AppendAttributes(attribute)
	}
	new = append(new, head)

	for _, pair := range e.Changes {
		attribute := sdk.NewAttribute(pair.Key, pair.Value)
		event := sdk.NewEvent(eventType, attribute)
		new = append(new, event)
	}
	return new
}

// Deprecated: use EventModifiedTokenClass
func NewEventModifyTokenType(e EventModifiedTokenClass) sdk.Events {
	eventType := EventTypeModifyTokenType.String()
	new := make(sdk.Events, 0, 1+len(e.Changes))

	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyTokenType:  e.ClassId,
	}
	head := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		head = head.AppendAttributes(attribute)
	}
	new = append(new, head)

	for _, pair := range e.Changes {
		attribute := sdk.NewAttribute(pair.Key, pair.Value)
		event := sdk.NewEvent(eventType, attribute)
		new = append(new, event)
	}
	return new
}

// Deprecated: use EventModifiedTokenClass
func NewEventModifyTokenOfFTClass(e EventModifiedTokenClass) sdk.Events {
	eventType := EventTypeModifyToken.String()
	new := make(sdk.Events, 0, 1+len(e.Changes))

	tokenID := NewFTID(e.ClassId)
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyTokenID:    tokenID,
	}
	head := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		head = head.AppendAttributes(attribute)
	}
	new = append(new, head)

	for _, pair := range e.Changes {
		attribute := sdk.NewAttribute(pair.Key, pair.Value)
		event := sdk.NewEvent(eventType, attribute)
		new = append(new, event)
	}
	return new
}

// Deprecated: use EventModifiedNFT
func NewEventModifyTokenOfNFT(e EventModifiedNFT) sdk.Events {
	eventType := EventTypeModifyToken.String()
	new := make(sdk.Events, 0, 1+len(e.Changes))

	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyTokenID:    e.TokenId,
	}
	head := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		head = head.AppendAttributes(attribute)
	}
	new = append(new, head)

	for _, pair := range e.Changes {
		attribute := sdk.NewAttribute(pair.Key, pair.Value)
		event := sdk.NewEvent(eventType, attribute)
		new = append(new, event)
	}
	return new
}

// Deprecated: use EventSent.
func NewEventTransferFT(e EventSent) *sdk.Event {
	eventType := EventTypeTransferFT.String()

	amount := []Coin{}
	for _, coin := range e.Amount {
		if err := ValidateFTID(coin.TokenId); err == nil {
			amount = append(amount, coin)
		}
	}
	if len(amount) == 0 {
		return nil
	}

	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyFrom:       e.From,
		AttributeKeyTo:         e.To,
		AttributeKeyAmount:     Coins(amount).String(),
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new = new.AppendAttributes(attribute)
	}
	return &new
}

// Deprecated: use EventSent.
func NewEventTransferNFT(e EventSent) sdk.Events {
	eventType := EventTypeTransferNFT.String()

	amount := []Coin{}
	for _, coin := range e.Amount {
		if err := ValidateNFTID(coin.TokenId); err == nil {
			amount = append(amount, coin)
		}
	}
	if len(amount) == 0 {
		return nil
	}

	new := make(sdk.Events, 0, 1+len(amount))

	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyFrom:       e.From,
		AttributeKeyTo:         e.To,
	}
	head := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		head = head.AppendAttributes(attribute)
	}
	new = append(new, head)

	for _, coin := range amount {
		attribute := sdk.NewAttribute(AttributeKeyTokenID.String(), coin.TokenId)
		event := sdk.NewEvent(eventType, attribute)
		new = append(new, event)
	}
	return new
}

// Deprecated: use EventSent.
func NewEventTransferFTFrom(e EventSent) *sdk.Event {
	eventType := EventTypeTransferFTFrom.String()

	amount := []Coin{}
	for _, coin := range e.Amount {
		if err := ValidateFTID(coin.TokenId); err == nil {
			amount = append(amount, coin)
		}
	}
	if len(amount) == 0 {
		return nil
	}

	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyProxy:      e.Operator,
		AttributeKeyFrom:       e.From,
		AttributeKeyTo:         e.To,
		AttributeKeyAmount:     Coins(amount).String(),
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new = new.AppendAttributes(attribute)
	}
	return &new
}

// Deprecated: use EventSent.
func NewEventTransferNFTFrom(e EventSent) sdk.Events {
	eventType := EventTypeTransferNFTFrom.String()

	amount := []Coin{}
	for _, coin := range e.Amount {
		if err := ValidateNFTID(coin.TokenId); err == nil {
			amount = append(amount, coin)
		}
	}
	if len(amount) == 0 {
		return nil
	}

	new := make(sdk.Events, 0, 1+len(amount))

	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyProxy:      e.Operator,
		AttributeKeyFrom:       e.From,
		AttributeKeyTo:         e.To,
	}
	head := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		head = head.AppendAttributes(attribute)
	}
	new = append(new, head)

	for _, coin := range amount {
		attribute := sdk.NewAttribute(AttributeKeyTokenID.String(), coin.TokenId)
		event := sdk.NewEvent(eventType, attribute)
		new = append(new, event)
	}
	return new
}

// Deprecated: use EventGrant.
func NewEventGrantPermToken(e EventGrant) sdk.Event {
	eventType := EventTypeGrantPermToken.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyTo:         e.Grantee,
		AttributeKeyPerm:       LegacyPermission(e.Permission).String(),
	}
	if len(e.Granter) != 0 {
		attributes[AttributeKeyFrom] = e.Granter
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new = new.AppendAttributes(attribute)
	}
	return new
}

// Deprecated: use EventGrant.
func NewEventGrantPermTokenHead(e EventGrant) sdk.Event {
	eventType := EventTypeGrantPermToken.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyTo:         e.Grantee,
	}
	if len(e.Granter) != 0 {
		attributes[AttributeKeyFrom] = e.Granter
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new = new.AppendAttributes(attribute)
	}
	return new
}

// Deprecated: use EventGrant.
func NewEventGrantPermTokenBody(e EventGrant) sdk.Event {
	eventType := EventTypeGrantPermToken.String()
	attributes := map[AttributeKey]string{
		AttributeKeyPerm: LegacyPermission(e.Permission).String(),
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new = new.AppendAttributes(attribute)
	}
	return new
}

// Deprecated: use EventAbandon.
func NewEventRevokePermToken(e EventAbandon) sdk.Event {
	eventType := EventTypeRevokePermToken.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyFrom:       e.Grantee,
		AttributeKeyPerm:       LegacyPermission(e.Permission).String(),
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new = new.AppendAttributes(attribute)
	}
	return new
}

// Deprecated: use EventAuthorizedOperator.
func NewEventApproveCollection(e EventAuthorizedOperator) sdk.Event {
	eventType := EventTypeApproveCollection.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyApprover:   e.Holder,
		AttributeKeyProxy:      e.Operator,
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new = new.AppendAttributes(attribute)
	}
	return new
}

// Deprecated: use EventRevokedOperator.
func NewEventDisapproveCollection(e EventRevokedOperator) sdk.Event {
	eventType := EventTypeDisapproveCollection.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyApprover:   e.Holder,
		AttributeKeyProxy:      e.Operator,
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new = new.AppendAttributes(attribute)
	}
	return new
}

// Deprecated: use EventAttached.
func NewEventAttachToken(e EventAttached, newRoot string) sdk.Event {
	eventType := EventTypeAttachToken.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyFrom:       e.Holder,
		AttributeKeyTokenID:    e.Subject,
		AttributeKeyToTokenID:  e.Target,

		AttributeKeyOldRoot: e.Subject,
		AttributeKeyNewRoot: newRoot,
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new = new.AppendAttributes(attribute)
	}
	return new
}

// Deprecated: use EventDetached.
func NewEventDetachToken(e EventDetached, oldRoot string) sdk.Event {
	eventType := EventTypeDetachToken.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID:  e.ContractId,
		AttributeKeyFrom:        e.Holder,
		AttributeKeyFromTokenID: e.Subject,

		AttributeKeyOldRoot: oldRoot,
		AttributeKeyNewRoot: e.Subject,
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new = new.AppendAttributes(attribute)
	}
	return new
}

// Deprecated: use EventAttached.
func NewEventAttachFrom(e EventAttached, newRoot string) sdk.Event {
	eventType := EventTypeAttachFrom.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: e.ContractId,
		AttributeKeyProxy:      e.Operator,
		AttributeKeyFrom:       e.Holder,
		AttributeKeyTokenID:    e.Subject,
		AttributeKeyToTokenID:  e.Target,

		AttributeKeyOldRoot: e.Subject,
		AttributeKeyNewRoot: newRoot,
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new = new.AppendAttributes(attribute)
	}
	return new
}

// Deprecated: use EventDetached.
func NewEventDetachFrom(e EventDetached, oldRoot string) sdk.Event {
	eventType := EventTypeDetachFrom.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID:  e.ContractId,
		AttributeKeyProxy:       e.Operator,
		AttributeKeyFrom:        e.Holder,
		AttributeKeyFromTokenID: e.Subject,

		AttributeKeyOldRoot: oldRoot,
		AttributeKeyNewRoot: e.Subject,
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new = new.AppendAttributes(attribute)
	}
	return new
}

// Deprecated: do not use.
func NewEventOperationTransferNFT(contractID string, tokenID string) sdk.Event {
	eventType := EventTypeOperationTransferNFT.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: contractID,
		AttributeKeyTokenID:    tokenID,
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new = new.AppendAttributes(attribute)
	}
	return new
}

// Deprecated: use EventBurned.
func NewEventOperationBurnNFT(contractID string, tokenID string) sdk.Event {
	eventType := EventTypeOperationBurnNFT.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: contractID,
		AttributeKeyTokenID:    tokenID,
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new = new.AppendAttributes(attribute)
	}
	return new
}

// Deprecated: do not use.
func NewEventOperationRootChanged(contractID string, tokenID string) sdk.Event {
	eventType := EventTypeOperationRootChanged.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: contractID,
		AttributeKeyTokenID:    tokenID,
	}

	new := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		new = new.AppendAttributes(attribute)
	}
	return new
}
