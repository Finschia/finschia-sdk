package collection

import (
	"github.com/line/lbm-sdk/v2/x/collection/internal/keeper"
	"github.com/line/lbm-sdk/v2/x/collection/internal/querier"
	"github.com/line/lbm-sdk/v2/x/collection/internal/types"
)

const (
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
	DefaultParamspace = types.DefaultParamspace
	EncodeRouterKey   = types.EncodeRouterKey
)

type (
	Token  = types.Token
	Tokens = types.Tokens
	Coins  = types.Coins
	FT     = types.FT
	NFT    = types.NFT

	TokenType = types.TokenType

	Collection = types.BaseCollection

	MsgCreateCollection = types.MsgCreateCollection
	MsgIssueFT          = types.MsgIssueFT
	MsgIssueNFT         = types.MsgIssueNFT
	MsgMintNFT          = types.MsgMintNFT
	MsgBurnNFT          = types.MsgBurnNFT
	MsgBurnNFTFrom      = types.MsgBurnNFTFrom
	MsgMintFT           = types.MsgMintFT
	MsgBurnFT           = types.MsgBurnFT
	MsgBurnFTFrom       = types.MsgBurnFTFrom

	MsgGrantPermission  = types.MsgGrantPermission
	MsgRevokePermission = types.MsgRevokePermission

	MsgModify = types.MsgModify

	MsgTransferFT  = types.MsgTransferFT
	MsgTransferNFT = types.MsgTransferNFT

	MsgTransferFTFrom  = types.MsgTransferFTFrom
	MsgTransferNFTFrom = types.MsgTransferNFTFrom

	MsgAttach = types.MsgAttach
	MsgDetach = types.MsgDetach

	MsgAttachFrom = types.MsgAttachFrom
	MsgDetachFrom = types.MsgDetachFrom

	MsgApproveCollection    = types.MsgApprove
	MsgDisapproveCollection = types.MsgDisapprove

	Permissions = types.Permissions
	Permission  = types.Permission

	MintNFTParam = types.MintNFTParam

	Keeper = keeper.Keeper
)

var (
	NewMsgCreateCollection = types.NewMsgCreateCollection
	NewMsgIssueFT          = types.NewMsgIssueFT
	NewMsgIssueNFT         = types.NewMsgIssueNFT
	NewMsgMintNFT          = types.NewMsgMintNFT
	NewMsgBurnNFT          = types.NewMsgBurnNFT
	NewMsgBurnNFTFrom      = types.NewMsgBurnNFTFrom
	NewMsgBurnFTFrom       = types.NewMsgBurnFTFrom
	NewMsgMintFT           = types.NewMsgMintFT
	NewMsgBurnFT           = types.NewMsgBurnFT
	NewMsgGrantPermission  = types.NewMsgGrantPermission
	NewMsgRevokePermission = types.NewMsgRevokePermission
	NewMsgModify           = types.NewMsgModify
	NewChangesWithMap      = types.NewChangesWithMap
	NewMsgTransferFT       = types.NewMsgTransferFT
	NewMsgTransferNFT      = types.NewMsgTransferNFT
	NewMsgTransferFTFrom   = types.NewMsgTransferFTFrom
	NewMsgTransferNFTFrom  = types.NewMsgTransferNFTFrom
	NewMsgAttach           = types.NewMsgAttach
	NewMsgDetach           = types.NewMsgDetach
	NewMsgAttachFrom       = types.NewMsgAttachFrom
	NewMsgDetachFrom       = types.NewMsgDetachFrom
	NewMsgApprove          = types.NewMsgApprove
	NewMsgDisapprove       = types.NewMsgDisapprove
	NewMintNFTParam        = types.NewMintNFTParam
	NewCoin                = types.NewCoin
	NewPermissions         = types.NewPermissions

	NewMintPermission   = types.NewMintPermission
	NewBurnPermission   = types.NewBurnPermission
	NewIssuePermission  = types.NewIssuePermission
	NewModifyPermission = types.NewModifyPermission

	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec

	NewKeeper  = keeper.NewKeeper
	NewQuerier = querier.NewQuerier

	NewMsgEncodeHandler = keeper.NewMsgEncodeHandler
	NewQueryEncoder     = querier.NewQueryEncoder
)
