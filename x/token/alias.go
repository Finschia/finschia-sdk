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
	Token  = types.Token
	Tokens = types.Tokens

	MsgPublishToken     = types.MsgPublishToken
	MsgMint             = types.MsgMint
	MsgBurn             = types.MsgBurn
	MsgGrantPermission  = types.MsgGrantPermission
	MsgRevokePermission = types.MsgRevokePermission

	QueryTokenParams             = types.QueryTokenParams
	QueryAccountPermissionParams = types.QueryAccountPermissionParams

	PermissionI = types.PermissionI
	Permissions = types.Permissions
)

var (
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec

	TokenSymbolKey       = types.TokenSymbolKey
	TokenSymbolKeyPrefix = types.TokenSymbolKeyPrefix

	ErrTokenExist          = types.ErrTokenExist
	ErrTokenNotExist       = types.ErrTokenNotExist
	ErrTokenNotMintable    = types.ErrTokenNotMintable
	ErrTokenPermission     = types.ErrTokenPermission
	ErrTokenPermissionMint = types.ErrTokenPermissionMint
	ErrTokenPermissionBurn = types.ErrTokenPermissionBurn

	EventTypePublishToken    = types.EventTypePublishToken
	EventTypeMintToken       = types.EventTypeMintToken
	EventTypeBurnToken       = types.EventTypeBurnToken
	EventTypeGrantPermToken  = types.EventTypeGrantPermToken
	EventTypeRevokePermToken = types.EventTypeRevokePermToken

	AttributeKeyName     = types.AttributeKeyName
	AttributeKeySymbol   = types.AttributeKeySymbol
	AttributeKeyOwner    = types.AttributeKeyOwner
	AttributeKeyAmount   = types.AttributeKeyAmount
	AttributeKeyMintable = types.AttributeKeyMintable
	AttributeKeyFrom     = types.AttributeKeyFrom
	AttributeKeyTo       = types.AttributeKeyTo
	AttributeKeyResource = types.AttributeKeyResource
	AttributeKeyAction   = types.AttributeKeyAction

	AttributeValueCategory = types.AttributeValueCategory

	NewMintPermission = types.NewMintPermission
	NewBurnPermission = types.NewBurnPermission
)
