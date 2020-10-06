package token

import (
	"github.com/line/link-modules/x/token/internal/keeper"
	"github.com/line/link-modules/x/token/internal/types"
)

const (
	ModuleName = types.ModuleName
	StoreKey   = types.StoreKey
	RouterKey  = types.RouterKey
)

type (
	Account     = types.Account
	Token       = types.Token
	Permissions = types.Permissions
	Keeper      = keeper.Keeper
	Permission  = types.Permission
)

var (
	NewMsgIssue            = types.NewMsgIssue
	NewMsgMint             = types.NewMsgMint
	NewMsgBurn             = types.NewMsgBurn
	NewMsgBurnFrom         = types.NewMsgBurnFrom
	NewMsgTransfer         = types.NewMsgTransfer
	NewMsgApprove          = types.NewMsgApprove
	NewMsgTransferFrom     = types.NewMsgTransferFrom
	NewMsgModify           = types.NewMsgModify
	NewChangesWithMap      = types.NewChangesWithMap
	NewMsgGrantPermission  = types.NewMsgGrantPermission
	NewMsgRevokePermission = types.NewMsgRevokePermission
	ModuleCdc              = types.ModuleCdc
	RegisterCodec          = types.RegisterCodec
	NewToken               = types.NewToken
	NewKeeper              = keeper.NewKeeper

	ErrTokenNotExist       = types.ErrTokenNotExist
	ErrInsufficientBalance = types.ErrInsufficientBalance
)
