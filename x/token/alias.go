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

	FT            = types.FT
	CollectiveNFT = types.CollectiveNFT

	Collection  = types.BaseCollection
	Collections = types.Collections

	MsgIssue = types.MsgIssue
	MsgMint  = types.MsgMint
	MsgBurn  = types.MsgBurn

	MsgCreateCollection = types.MsgCreateCollection
	MsgIssueCFT         = types.MsgIssueCFT
	MsgIssueCNFT        = types.MsgIssueCNFT
	MsgMintCNFT         = types.MsgMintCNFT
	MsgMintCFT          = types.MsgMintCFT
	MsgBurnCFT          = types.MsgBurnCFT

	MsgGrantPermission  = types.MsgGrantPermission
	MsgRevokePermission = types.MsgRevokePermission

	MsgModifyTokenURI = types.MsgModifyTokenURI

	MsgTransferFT   = types.MsgTransferFT
	MsgTransferCFT  = types.MsgTransferCFT
	MsgTransferCNFT = types.MsgTransferCNFT

	MsgAttach = types.MsgAttach
	MsgDetach = types.MsgDetach

	Permissions = types.Permissions

	Keeper = keeper.Keeper
)

var (
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec

	NewKeeper = keeper.NewKeeper
)
