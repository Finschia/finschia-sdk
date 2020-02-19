package handler

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/keeper"
	"github.com/line/link/x/collection/internal/types"
)

func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgCreateCollection:
			return handleMsgCreateCollection(ctx, keeper, msg)
		case types.MsgIssueCFT:
			return handleMsgIssueCFT(ctx, keeper, msg)
		case types.MsgMintCNFT:
			return handleMsgMintCNFT(ctx, keeper, msg)
		case types.MsgBurnCNFT:
			return handleMsgBurnCNFT(ctx, keeper, msg)
		case types.MsgBurnCNFTFrom:
			return handleMsgBurnCNFTFrom(ctx, keeper, msg)
		case types.MsgIssueCNFT:
			return handleMsgIssueCNFT(ctx, keeper, msg)
		case types.MsgMintCFT:
			return handleMsgMintCFT(ctx, keeper, msg)
		case types.MsgBurnCFT:
			return handleMsgBurnCFT(ctx, keeper, msg)
		case types.MsgBurnCFTFrom:
			return handleMsgBurnCFTFrom(ctx, keeper, msg)
		case types.MsgGrantPermission:
			return handleMsgGrant(ctx, keeper, msg)
		case types.MsgRevokePermission:
			return handleMsgRevoke(ctx, keeper, msg)
		case types.MsgModifyTokenURI:
			return handleMsgModifyTokenURI(ctx, keeper, msg)
		case types.MsgTransferCFT:
			return handleMsgTransferCFT(ctx, keeper, msg)
		case types.MsgTransferCNFT:
			return handleMsgTransferCNFT(ctx, keeper, msg)
		case types.MsgTransferCFTFrom:
			return handleMsgTransferCFTFrom(ctx, keeper, msg)
		case types.MsgTransferCNFTFrom:
			return handleMsgTransferCNFTFrom(ctx, keeper, msg)
		case types.MsgAttach:
			return handleMsgAttach(ctx, keeper, msg)
		case types.MsgDetach:
			return handleMsgDetach(ctx, keeper, msg)
		case types.MsgAttachFrom:
			return handleMsgAttachFrom(ctx, keeper, msg)
		case types.MsgDetachFrom:
			return handleMsgDetachFrom(ctx, keeper, msg)
		case types.MsgApprove:
			return handleMsgApprove(ctx, keeper, msg)
		case types.MsgDisapprove:
			return handleMsgDisapprove(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized  Msg type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}
