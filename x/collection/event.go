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
func NewEventCreateCollection(event EventCreatedContract, creator sdk.AccAddress) sdk.Event {
	eventType := EventTypeCreateCollection.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyName:       event.Name,
		AttributeKeyMeta:       event.Meta,
		AttributeKeyBaseImgURI: event.BaseImgUri,

		AttributeKeyOwner: creator.String(),
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return res
}

// Deprecated: use EventCreatedFTClass.
func NewEventIssueFT(event EventCreatedFTClass, operator, to sdk.AccAddress, amount sdk.Int) sdk.Event {
	eventType := EventTypeIssueFT.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyTokenID:    NewFTID(event.ClassId),
		AttributeKeyName:       event.Name,
		AttributeKeyMeta:       event.Meta,
		AttributeKeyDecimals:   fmt.Sprintf("%d", event.Decimals),
		AttributeKeyMintable:   fmt.Sprintf("%t", event.Mintable),

		AttributeKeyOwner:  operator.String(),
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

// Deprecated: use EventCreatedNFTClass.
func NewEventIssueNFT(event EventCreatedNFTClass) sdk.Event {
	eventType := EventTypeIssueNFT.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyTokenType:  event.ClassId,
		AttributeKeyName:       event.Name,
		AttributeKeyMeta:       event.Meta,
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return res
}

// Deprecated: use EventMintedFT.
func NewEventMintFT(event EventMintedFT) sdk.Event {
	eventType := EventTypeMintFT.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyFrom:       event.Operator,
		AttributeKeyTo:         event.To,
		AttributeKeyAmount:     Coins(event.Amount).String(),
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return res
}

// Deprecated: use EventMintedNFT.
func NewEventMintNFT(event EventMintedNFT) sdk.Events {
	eventType := EventTypeMintNFT.String()

	res := make(sdk.Events, len(event.Tokens))
	for i, token := range event.Tokens {
		e := sdk.NewEvent(eventType)
		attributes := map[AttributeKey]string{
			AttributeKeyContractID: event.ContractId,
			AttributeKeyFrom:       event.Operator,
			AttributeKeyTo:         event.To,

			AttributeKeyTokenID: token.Id,
			AttributeKeyName:    token.Name,
			AttributeKeyMeta:    token.Meta,
		}
		for key, value := range attributes {
			attribute := sdk.NewAttribute(key.String(), value)
			e = e.AppendAttributes(attribute)
		}

		res[i] = e
	}

	return res
}

// Deprecated: use EventBurned.
func NewEventBurnFT(event EventBurned) *sdk.Event {
	eventType := EventTypeBurnFT.String()

	amount := []Coin{}
	for _, coin := range event.Amount {
		if err := ValidateFTID(coin.TokenId); err == nil {
			amount = append(amount, coin)
		}
	}
	if len(amount) == 0 {
		return nil
	}

	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyFrom:       event.From,
		AttributeKeyAmount:     Coins(amount).String(),
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return &res
}

// Deprecated: use EventBurned.
func NewEventBurnNFT(event EventBurned) sdk.Events {
	eventType := EventTypeBurnNFT.String()

	amount := []Coin{}
	for _, coin := range event.Amount {
		if err := ValidateNFTID(coin.TokenId); err == nil {
			amount = append(amount, coin)
		}
	}
	if len(amount) == 0 {
		return nil
	}

	res := make(sdk.Events, 0, len(amount)+1)

	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyFrom:       event.From,
	}
	head := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		head = head.AppendAttributes(attribute)
	}
	res = append(res, head)

	for _, coin := range amount {
		attribute := sdk.NewAttribute(AttributeKeyTokenID.String(), coin.TokenId)
		event := sdk.NewEvent(eventType, attribute)
		res = append(res, event)
	}

	return res
}

// Deprecated: use EventBurned.
func NewEventBurnFTFrom(event EventBurned) *sdk.Event {
	eventType := EventTypeBurnFTFrom.String()

	amount := []Coin{}
	for _, coin := range event.Amount {
		if err := ValidateFTID(coin.TokenId); err == nil {
			amount = append(amount, coin)
		}
	}
	if len(amount) == 0 {
		return nil
	}

	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyProxy:      event.Operator,
		AttributeKeyFrom:       event.From,
		AttributeKeyAmount:     Coins(amount).String(),
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return &res
}

// Deprecated: use EventBurned.
func NewEventBurnNFTFrom(event EventBurned) sdk.Events {
	eventType := EventTypeBurnNFTFrom.String()

	amount := []Coin{}
	for _, coin := range event.Amount {
		if err := ValidateNFTID(coin.TokenId); err == nil {
			amount = append(amount, coin)
		}
	}
	if len(amount) == 0 {
		return nil
	}

	res := make(sdk.Events, 0, len(amount)+1)

	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyProxy:      event.Operator,
		AttributeKeyFrom:       event.From,
	}
	head := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		head = head.AppendAttributes(attribute)
	}
	res = append(res, head)

	for _, coin := range amount {
		attribute := sdk.NewAttribute(AttributeKeyTokenID.String(), coin.TokenId)
		event := sdk.NewEvent(eventType, attribute)
		res = append(res, event)
	}

	return res
}

// Deprecated: use EventModifiedContract
func NewEventModifyCollection(event EventModifiedContract) sdk.Events {
	eventType := EventTypeModifyCollection.String()
	res := make(sdk.Events, 0, 1+len(event.Changes))

	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
	}
	head := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		head = head.AppendAttributes(attribute)
	}
	res = append(res, head)

	for _, pair := range event.Changes {
		attribute := sdk.NewAttribute(pair.Key, pair.Value)
		event := sdk.NewEvent(eventType, attribute)
		res = append(res, event)
	}
	return res
}

// Deprecated: use EventModifiedTokenClass
func NewEventModifyTokenType(event EventModifiedTokenClass) sdk.Events {
	eventType := EventTypeModifyTokenType.String()
	res := make(sdk.Events, 0, 1+len(event.Changes))

	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyTokenType:  event.ClassId,
	}
	head := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		head = head.AppendAttributes(attribute)
	}
	res = append(res, head)

	for _, pair := range event.Changes {
		attribute := sdk.NewAttribute(pair.Key, pair.Value)
		event := sdk.NewEvent(eventType, attribute)
		res = append(res, event)
	}
	return res
}

// Deprecated: use EventModifiedTokenClass
func NewEventModifyTokenOfFTClass(event EventModifiedTokenClass) sdk.Events {
	eventType := EventTypeModifyToken.String()
	res := make(sdk.Events, 0, 1+len(event.Changes))

	tokenID := NewFTID(event.ClassId)
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyTokenID:    tokenID,
	}
	head := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		head = head.AppendAttributes(attribute)
	}
	res = append(res, head)

	for _, pair := range event.Changes {
		attribute := sdk.NewAttribute(pair.Key, pair.Value)
		event := sdk.NewEvent(eventType, attribute)
		res = append(res, event)
	}
	return res
}

// Deprecated: use EventModifiedNFT
func NewEventModifyTokenOfNFT(event EventModifiedNFT) sdk.Events {
	eventType := EventTypeModifyToken.String()
	res := make(sdk.Events, 0, 1+len(event.Changes))

	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyTokenID:    event.TokenId,
	}
	head := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		head = head.AppendAttributes(attribute)
	}
	res = append(res, head)

	for _, pair := range event.Changes {
		attribute := sdk.NewAttribute(pair.Key, pair.Value)
		event := sdk.NewEvent(eventType, attribute)
		res = append(res, event)
	}
	return res
}

// Deprecated: use EventSent.
func NewEventTransferFT(event EventSent) *sdk.Event {
	eventType := EventTypeTransferFT.String()

	amount := []Coin{}
	for _, coin := range event.Amount {
		if err := ValidateFTID(coin.TokenId); err == nil {
			amount = append(amount, coin)
		}
	}
	if len(amount) == 0 {
		return nil
	}

	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyFrom:       event.From,
		AttributeKeyTo:         event.To,
		AttributeKeyAmount:     Coins(amount).String(),
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return &res
}

// Deprecated: use EventSent.
func NewEventTransferNFT(event EventSent) sdk.Events {
	eventType := EventTypeTransferNFT.String()

	amount := []Coin{}
	for _, coin := range event.Amount {
		if err := ValidateNFTID(coin.TokenId); err == nil {
			amount = append(amount, coin)
		}
	}
	if len(amount) == 0 {
		return nil
	}

	res := make(sdk.Events, 0, 1+len(amount))

	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyFrom:       event.From,
		AttributeKeyTo:         event.To,
	}
	head := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		head = head.AppendAttributes(attribute)
	}
	res = append(res, head)

	for _, coin := range amount {
		attribute := sdk.NewAttribute(AttributeKeyTokenID.String(), coin.TokenId)
		event := sdk.NewEvent(eventType, attribute)
		res = append(res, event)
	}
	return res
}

// Deprecated: use EventSent.
func NewEventTransferFTFrom(event EventSent) *sdk.Event {
	eventType := EventTypeTransferFTFrom.String()

	amount := []Coin{}
	for _, coin := range event.Amount {
		if err := ValidateFTID(coin.TokenId); err == nil {
			amount = append(amount, coin)
		}
	}
	if len(amount) == 0 {
		return nil
	}

	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyProxy:      event.Operator,
		AttributeKeyFrom:       event.From,
		AttributeKeyTo:         event.To,
		AttributeKeyAmount:     Coins(amount).String(),
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return &res
}

// Deprecated: use EventSent.
func NewEventTransferNFTFrom(event EventSent) sdk.Events {
	eventType := EventTypeTransferNFTFrom.String()

	amount := []Coin{}
	for _, coin := range event.Amount {
		if err := ValidateNFTID(coin.TokenId); err == nil {
			amount = append(amount, coin)
		}
	}
	if len(amount) == 0 {
		return nil
	}

	res := make(sdk.Events, 0, 1+len(amount))

	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyProxy:      event.Operator,
		AttributeKeyFrom:       event.From,
		AttributeKeyTo:         event.To,
	}
	head := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		head = head.AppendAttributes(attribute)
	}
	res = append(res, head)

	for _, coin := range amount {
		attribute := sdk.NewAttribute(AttributeKeyTokenID.String(), coin.TokenId)
		event := sdk.NewEvent(eventType, attribute)
		res = append(res, event)
	}
	return res
}

// Deprecated: use EventGrant.
func NewEventGrantPermToken(event EventGranted) sdk.Event {
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

// Deprecated: use EventGrant.
func NewEventGrantPermTokenHead(event EventGranted) sdk.Event {
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

// Deprecated: use EventGrant.
func NewEventGrantPermTokenBody(event EventGranted) sdk.Event {
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

// Deprecated: use EventAbandon.
func NewEventRevokePermToken(event EventRenounced) sdk.Event {
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

// Deprecated: use EventAuthorizedOperator.
func NewEventApproveCollection(event EventAuthorizedOperator) sdk.Event {
	eventType := EventTypeApproveCollection.String()
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

// Deprecated: use EventRevokedOperator.
func NewEventDisapproveCollection(event EventRevokedOperator) sdk.Event {
	eventType := EventTypeDisapproveCollection.String()
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

// Deprecated: use EventAttached.
func NewEventAttachToken(event EventAttached, newRoot string) sdk.Event {
	eventType := EventTypeAttachToken.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyFrom:       event.Holder,
		AttributeKeyTokenID:    event.Subject,
		AttributeKeyToTokenID:  event.Target,

		AttributeKeyOldRoot: event.Subject,
		AttributeKeyNewRoot: newRoot,
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return res
}

// Deprecated: use EventDetached.
func NewEventDetachToken(event EventDetached, oldRoot string) sdk.Event {
	eventType := EventTypeDetachToken.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID:  event.ContractId,
		AttributeKeyFrom:        event.Holder,
		AttributeKeyFromTokenID: event.Subject,

		AttributeKeyOldRoot: oldRoot,
		AttributeKeyNewRoot: event.Subject,
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return res
}

// Deprecated: use EventAttached.
func NewEventAttachFrom(event EventAttached, newRoot string) sdk.Event {
	eventType := EventTypeAttachFrom.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyProxy:      event.Operator,
		AttributeKeyFrom:       event.Holder,
		AttributeKeyTokenID:    event.Subject,
		AttributeKeyToTokenID:  event.Target,

		AttributeKeyOldRoot: event.Subject,
		AttributeKeyNewRoot: newRoot,
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return res
}

// Deprecated: use EventDetached.
func NewEventDetachFrom(event EventDetached, oldRoot string) sdk.Event {
	eventType := EventTypeDetachFrom.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID:  event.ContractId,
		AttributeKeyProxy:       event.Operator,
		AttributeKeyFrom:        event.Holder,
		AttributeKeyFromTokenID: event.Subject,

		AttributeKeyOldRoot: oldRoot,
		AttributeKeyNewRoot: event.Subject,
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return res
}

// Deprecated: do not use.
func NewEventOperationTransferNFT(event EventOwnerChanged) sdk.Event {
	eventType := EventTypeOperationTransferNFT.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyTokenID:    event.TokenId,
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return res
}

// Deprecated: use EventBurned.
func NewEventOperationBurnNFT(contractID string, tokenID string) sdk.Event {
	eventType := EventTypeOperationBurnNFT.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: contractID,
		AttributeKeyTokenID:    tokenID,
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return res
}

// Deprecated: do not use.
func NewEventOperationRootChanged(event EventRootChanged) sdk.Event {
	eventType := EventTypeOperationRootChanged.String()
	attributes := map[AttributeKey]string{
		AttributeKeyContractID: event.ContractId,
		AttributeKeyTokenID:    event.TokenId,
	}

	res := sdk.NewEvent(eventType)
	for key, value := range attributes {
		attribute := sdk.NewAttribute(key.String(), value)
		res = res.AppendAttributes(attribute)
	}
	return res
}
