package handler

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/lbm-sdk/v2/x/collection/internal/keeper"
	"github.com/line/lbm-sdk/v2/x/collection/internal/types"
	"github.com/line/lbm-sdk/v2/x/contract"
)

func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		if msg, ok := msg.(contract.Msg); ok {
			ctx = ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, msg.GetContractID()))
			err := handleMsgContract(ctx, keeper, msg)
			if err != nil {
				return nil, err
			}
		}

		if _, ok := msg.(types.MsgCreateCollection); ok {
			contractID := keeper.NewContractID(ctx)
			ctx = ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, contractID))
		}
		if ctx.Context().Value(contract.CtxKey{}) == nil {
			panic("contract id does not set")
		}
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
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized Msg type: %T", msg)
		}
	}
}
func handleMsgContract(ctx sdk.Context, keeper keeper.Keeper, msg contract.Msg) error {
	if !keeper.HasContractID(ctx) {
		return sdkerrors.Wrapf(contract.ErrContractNotExist, "contract id: %s", msg.GetContractID())
	}
	return nil
}
