package collection

import (
	"github.com/line/link/x/collection/internal/keeper"
	"github.com/line/link/x/collection/internal/types"
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

	FT  = types.FT
	NFT = types.NFT

	Collection  = types.BaseCollection
	Collections = types.Collections

	MsgCreateCollection = types.MsgCreateCollection
	MsgIssueCFT         = types.MsgIssueCFT
	MsgIssueCNFT        = types.MsgIssueCNFT
	MsgMintCNFT         = types.MsgMintCNFT
	MsgBurnCNFT         = types.MsgBurnCNFT
	MsgBurnCNFTFrom     = types.MsgBurnCNFTFrom
	MsgMintCFT          = types.MsgMintCFT
	MsgBurnCFT          = types.MsgBurnCFT
	MsgBurnCFTFrom      = types.MsgBurnCFTFrom

	MsgGrantPermission  = types.MsgGrantPermission
	MsgRevokePermission = types.MsgRevokePermission

	MsgModifyTokenURI = types.MsgModifyTokenURI

	MsgTransferCFT  = types.MsgTransferCFT
	MsgTransferCNFT = types.MsgTransferCNFT

	MsgTransferCFTFrom  = types.MsgTransferCFTFrom
	MsgTransferCNFTFrom = types.MsgTransferCNFTFrom

	MsgAttach = types.MsgAttach
	MsgDetach = types.MsgDetach

	MsgAttachFrom = types.MsgAttachFrom
	MsgDetachFrom = types.MsgDetachFrom

	MsgApproveCollection    = types.MsgApprove
	MsgDisapproveCollection = types.MsgDisapprove

	Permissions = types.Permissions

	Keeper = keeper.Keeper
)

var (
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec

	NewKeeper = keeper.NewKeeper
)
