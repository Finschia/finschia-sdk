package handler

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/contract"
	"github.com/line/link-modules/x/token/internal/keeper"
	"github.com/line/link-modules/x/token/internal/types"
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
		if _, ok := msg.(types.MsgIssue); ok {
			contractID := keeper.NewContractID(ctx)
			ctx = ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, contractID))
		}
		if ctx.Context().Value(contract.CtxKey{}) == nil {
			panic("contract id does not set")
		}
		switch msg := msg.(type) {
		case types.MsgIssue:
			return handleMsgIssue(ctx, keeper, msg)
		case types.MsgMint:
			return handleMsgMint(ctx, keeper, msg)
		case types.MsgBurn:
			return handleMsgBurn(ctx, keeper, msg)
		case types.MsgBurnFrom:
			return handleMsgBurnFrom(ctx, keeper, msg)
		case types.MsgTransfer:
			return handleMsgTransfer(ctx, keeper, msg)
		case types.MsgGrantPermission:
			return handleMsgGrant(ctx, keeper, msg)
		case types.MsgRevokePermission:
			return handleMsgRevoke(ctx, keeper, msg)
		case types.MsgModify:
			return handleMsgModify(ctx, keeper, msg)
		case types.MsgTransferFrom:
			return handleMsgTransferFrom(ctx, keeper, msg)
		case types.MsgApprove:
			return handleMsgApprove(ctx, keeper, msg)
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
