package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrProxyInvalidMsgType              = sdkerrors.Register(ModuleName, 1, "Invalid msg type")
	ErrProxyAddressRequired             = sdkerrors.Register(ModuleName, 2, "Proxy Address is required")
	ErrProxyOnBehalfOfAddressRequired   = sdkerrors.Register(ModuleName, 3, "OnBehalfOf Address is required")
	ErrProxyDenomRequired               = sdkerrors.Register(ModuleName, 4, "Denom is required")
	ErrProxyAmountMustBePositiveInteger = sdkerrors.Register(ModuleName, 5, "Amount must be a positive integer")
	ErrProxyNotExist                    = sdkerrors.Register(ModuleName, 6, "Proxy is not approved for the account")
	ErrProxyDenomDoesNotMatch           = sdkerrors.Register(ModuleName, 7, "Two proxies must have the same Proxy, OnBehalOf and Denom to add or subtract")
	ErrProxyNotEnoughApprovedCoins      = sdkerrors.Register(ModuleName, 8, "Not enough approved coins")
	ErrProxyToAddressRequired           = sdkerrors.Register(ModuleName, 9, "ToAddress is required to send coins")
)
