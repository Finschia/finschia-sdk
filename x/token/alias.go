package token

import (
	"github.com/line/lbm-sdk/v2/x/token/internal/keeper"
	"github.com/line/lbm-sdk/v2/x/token/internal/querier"
	"github.com/line/lbm-sdk/v2/x/token/internal/types"
)

const (
	ModuleName      = types.ModuleName
	StoreKey        = types.StoreKey
	RouterKey       = types.RouterKey
	EncodeRouterKey = types.EncodeRouterKey
)

type (
	MsgIssue    = types.MsgIssue
	MsgTransfer = types.MsgTransfer
	MsgMint     = types.MsgMint
	MsgBurn     = types.MsgBurn
	MsgModify   = types.MsgModify

	Account     = types.Account
	Token       = types.Token
	Permissions = types.Permissions
	Keeper      = keeper.Keeper
	Permission  = types.Permission

	EncodeHandler = types.EncodeHandler
	EncodeQuerier = types.EncodeQuerier
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
	NewQuerier             = querier.NewQuerier

	NewMintPermission   = types.NewMintPermission
	NewBurnPermission   = types.NewBurnPermission
	NewModifyPermission = types.NewModifyPermission

	NewMsgEncodeHandler = keeper.NewMsgEncodeHandler
	NewQueryEncoder     = querier.NewQueryEncoder

	NewChanges = types.NewChanges
	NewChange  = types.NewChange

	ErrTokenNotExist       = types.ErrTokenNotExist
	ErrInsufficientBalance = types.ErrInsufficientBalance
)
