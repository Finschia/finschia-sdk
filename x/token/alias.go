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

	FT              = types.FT
	NFT             = types.NFT
	CollectiveToken = types.CollectiveToken

	Collection  = types.BaseCollection
	Collections = types.Collections

	MsgIssue    = types.MsgIssue
	MsgIssueNFT = types.MsgIssueNFT
	MsgMint     = types.MsgMint
	MsgBurn     = types.MsgBurn

	MsgCreateCollection    = types.MsgCreateCollection
	MsgIssueCollection     = types.MsgIssueCollection
	MsgIssueNFTCollection  = types.MsgIssueNFTCollection
	MsgCollectionTokenMint = types.MsgMintCollection
	MsgCollectionTokenBurn = types.MsgBurnCollection

	MsgGrantPermission  = types.MsgGrantPermission
	MsgRevokePermission = types.MsgRevokePermission

	MsgModifyTokenURI = types.MsgModifyTokenURI

	MsgTransferFT   = types.MsgTransferFT
	MsgTransferCFT  = types.MsgTransferCFT
	MsgTransferNFT  = types.MsgTransferNFT
	MsgTransferCNFT = types.MsgTransferCNFT

	MsgAttach = types.MsgAttach
	MsgDetach = types.MsgDetach

	PermissionI = types.PermissionI
	Permissions = types.Permissions

	Keeper = keeper.Keeper
)

var (
	NewFT  = types.NewFT
	NewNFT = types.NewNFT

	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec

	ErrTokenExist           = types.ErrTokenExist
	ErrCollectionTokenExist = types.ErrCollectionTokenExist

	NewIssuePermission          = types.NewIssuePermission
	NewModifyTokenURIPermission = types.NewModifyTokenURIPermission

	NewKeeper = keeper.NewKeeper
)
