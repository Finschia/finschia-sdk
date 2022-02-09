package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/token"
)

// NewHandler creates an sdk.Handler for all the token type messages
func NewHandler(k Keeper) sdk.Handler {
	msgServer := NewMsgServer(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *token.MsgTransfer:
			res, err := msgServer.Transfer(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *token.MsgTransferFrom:
			res, err := msgServer.TransferFrom(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *token.MsgApprove:
			res, err := msgServer.Approve(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *token.MsgIssue:
			res, err := msgServer.Issue(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *token.MsgGrant:
			res, err := msgServer.Grant(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *token.MsgRevoke:
			res, err := msgServer.Revoke(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *token.MsgMint:
			res, err := msgServer.Mint(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *token.MsgBurn:
			res, err := msgServer.Burn(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *token.MsgBurnFrom:
			res, err := msgServer.BurnFrom(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *token.MsgModify:
			res, err := msgServer.Modify(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, errors.Wrapf(errors.ErrUnknownRequest, "unrecognized %s message type: %T", token.ModuleName, msg)
		}
	}
}
