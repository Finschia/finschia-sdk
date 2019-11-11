package lrc3

import (
	"github.com/link-chain/link/x/lrc3/internal/keeper"
	"github.com/link-chain/link/x/lrc3/internal/types"
	nft "github.com/link-chain/link/x/nft"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	NewKeeper               = keeper.NewKeeper
	NewQuerier              = keeper.NewQuerier
	NewMsgInit              = types.NewMsgInit
	NewMsgMintNFT           = types.NewMsgMintNFT
	NewMsgBurnNFT           = nft.NewMsgBurnNFT
	NewMsgTransferNFT       = nft.NewMsgTransferNFT
	NewMsgEditNFTMetadata   = nft.NewMsgEditNFTMetadata
	NewMsgApprove           = types.NewMsgApprove
	NewMsgSetApprovalForAll = types.NewMsgSetApprovalForAll

	EventTypeInit              = types.EventTypeInit
	EventTypeMint              = types.EventTypeMintNFT
	EventTypeApprove           = types.EventTypeApprove
	EventTypeSetApprovalForAll = types.EventTypeSetApprovalForAll

	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	Keeper               = keeper.Keeper
	TokenBalance         = types.TokenBalance
	Approval             = types.Approval
	TokenOwner           = types.TokenOwner
	OperatorApprovals    = types.OperatorApprovals
	MsgInit              = types.MsgInit
	MsgMintNFT           = types.MsgMintNFT
	MsgBurn              = types.MsgBurn
	MsgTransfer          = types.MsgTransfer
	MsgEditMetadata      = types.MsgEditMetadata
	MsgMsgApprove        = types.MsgApprove
	MsgSetApprovalForAll = types.MsgSetApprovalForAll
)
