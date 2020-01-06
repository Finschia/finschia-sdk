package token

import (
	"github.com/line/link/x/token/internal/keeper"
	"github.com/line/link/x/token/internal/types"
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

	Collection            = types.Collection
	Collections           = types.Collections
	CollectionWithTokens  = types.CollectionWithTokens
	CollectionsWithTokens = types.CollectionsWithTokens

	MsgIssue              = types.MsgIssue
	MsgIssueCollection    = types.MsgIssueCollection
	MsgIssueNFT           = types.MsgIssueNFT
	MsgIssueNFTCollection = types.MsgIssueNFTCollection
	MsgMint               = types.MsgMint
	MsgBurn               = types.MsgBurn
	MsgGrantPermission    = types.MsgGrantPermission
	MsgRevokePermission   = types.MsgRevokePermission

	PermissionI = types.PermissionI
	Permissions = types.Permissions

	Keeper = keeper.Keeper
)

var (
	NewFT  = types.NewFT
	NewNFT = types.NewNFT

	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec

	ErrTokenExist = types.ErrTokenExist

	NewIssuePermission = types.NewIssuePermission

	NewKeeper = keeper.NewKeeper
)
