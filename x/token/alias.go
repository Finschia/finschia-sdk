package token

import (
	"github.com/link-chain/link/x/token/types"
)

const (
	ModuleName = types.ModuleName
	StoreKey   = types.StoreKey
	RouterKey  = types.RouterKey

	DefaultCodespace = types.DefaultCodespace
)

type (
	Token = types.Token

	MsgPublishToken = types.MsgPublishToken
	MsgMint         = types.MsgMint
	MsgBurn         = types.MsgBurn

	QueryTokenParams = types.QueryTokenParams
)

var (
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec

	TokenSymbolKey       = types.TokenSymbolKey
	TokenSymbolKeyPrefix = types.TokenSymbolKeyPrefix

	ErrTokenExist          = types.ErrTokenExist
	ErrTokenNotExist       = types.ErrTokenNotExist
	ErrTokenNotMintable    = types.ErrTokenNotMintable
	ErrTokenPermissionMint = types.ErrTokenPermissionMint
	ErrTokenPermissionBurn = types.ErrTokenPermissionBurn

	EventTypePublishToken  = types.EventTypePublishToken
	EventTypeTransferToken = types.EventTypeTransferToken
	EventTypeMintToken     = types.EventTypeMintToken
	EventTypeBurnToken     = types.EventTypeBurnToken
	AttributeKeyName       = types.AttributeKeyName
	AttributeKeySymbol     = types.AttributeKeySymbol
	AttributeKeyOwner      = types.AttributeKeyOwner
	AttributeKeyAmount     = types.AttributeKeyAmount
	AttributeKeyMintable   = types.AttributeKeyMintable
	AttributeKeyFrom       = types.AttributeKeyFrom
	AttributeKeyTo         = types.AttributeKeyTo
)
