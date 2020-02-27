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
		case types.MsgIssueFT:
			return handleMsgIssueFT(ctx, keeper, msg)
		case types.MsgMintNFT:
			return handleMsgMintNFT(ctx, keeper, msg)
		case types.MsgBurnNFT:
			return handleMsgBurnNFT(ctx, keeper, msg)
		case types.MsgBurnNFTFrom:
			return handleMsgBurnNFTFrom(ctx, keeper, msg)
		case types.MsgIssueNFT:
			return handleMsgIssueNFT(ctx, keeper, msg)
		case types.MsgMintFT:
			return handleMsgMintFT(ctx, keeper, msg)
		case types.MsgBurnFT:
			return handleMsgBurnFT(ctx, keeper, msg)
		case types.MsgBurnFTFrom:
			return handleMsgBurnFTFrom(ctx, keeper, msg)
		case types.MsgGrantPermission:
			return handleMsgGrant(ctx, keeper, msg)
		case types.MsgRevokePermission:
			return handleMsgRevoke(ctx, keeper, msg)
		case types.MsgModify:
			return handleMsgModify(ctx, keeper, msg)
		case types.MsgTransferFT:
			return handleMsgTransferFT(ctx, keeper, msg)
		case types.MsgTransferNFT:
			return handleMsgTransferNFT(ctx, keeper, msg)
		case types.MsgTransferFTFrom:
			return handleMsgTransferFTFrom(ctx, keeper, msg)
		case types.MsgTransferNFTFrom:
			return handleMsgTransferNFTFrom(ctx, keeper, msg)
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
