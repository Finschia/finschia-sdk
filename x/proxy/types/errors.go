package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeProxyInvalidMsgType              sdk.CodeType = 1
	CodeProxyAddressRequired             sdk.CodeType = 2
	CodeProxyOnBehalfOfAddressRequired   sdk.CodeType = 3
	CodeProxyDenomRequired               sdk.CodeType = 4
	CodeProxyAmountMustBePositiveInteger sdk.CodeType = 5
	CodeProxyNotExist                    sdk.CodeType = 6
	CodeProxyDenomDoesNotMatch           sdk.CodeType = 7
	CodeProxyNotEnoughApprovedCoins      sdk.CodeType = 8
	CodeProxyToAddressRequired           sdk.CodeType = 9
)

func ErrProxyInvalidMsgType(codespace sdk.CodespaceType, msgType string) sdk.Error {
	return sdk.NewError(codespace, CodeProxyInvalidMsgType, "Invalid msg type: %s", msgType)
}

func ErrProxyAddressRequired(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeProxyAddressRequired, "Proxy Address is required")
}

func ErrProxyOnBehalfOfAddressRequired(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeProxyOnBehalfOfAddressRequired, "OnBehalfOf Address is required")
}

func ErrProxyDenomRequired(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeProxyDenomRequired, "Denom is required")
}

func ErrProxyAmountMustBePositiveInteger(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeProxyAmountMustBePositiveInteger, "Amount must be a positive integer")
}

func ErrProxyNotExist(codespace sdk.CodespaceType, proxy, onBehalfOf string) sdk.Error {
	return sdk.NewError(codespace, CodeProxyNotExist, "Proxy (%s) is not approved for the account (%s)", proxy, onBehalfOf)
}

func ErrProxyDenomDoesNotMatch(codespace sdk.CodespaceType, pa1, pa2 ProxyAllowance) sdk.Error {
	return sdk.NewError(codespace, CodeProxyDenomDoesNotMatch, "Two proxies must have the same Proxy, OnBehalOf and Denom to add or subtract (%v compared to %v)", pa1.ProxyDenom, pa2.ProxyDenom)
}

func ErrProxyNotEnoughApprovedCoins(codespace sdk.CodespaceType, approvedAmount, requestedAmount sdk.Int) sdk.Error {
	return sdk.NewError(codespace, CodeProxyNotEnoughApprovedCoins, "Not enough approved coins (Approved: %v, Requested: %v)", approvedAmount, requestedAmount)
}

func ErrProxyToAddressRequired(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeProxyToAddressRequired, "ToAddress is required to send coins")
}
