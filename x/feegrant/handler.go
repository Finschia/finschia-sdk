package feegrant

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/feegrant/keeper"
	"github.com/line/lbm-sdk/x/feegrant/types"
)

// NewHandler creates an sdk.Handler for all the gov type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgGrantFeeAllowance:
			res, err := msgServer.GrantAllowance(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgRevokeFeeAllowance:
			res, err := msgServer.RevokeAllowance(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, errors.Wrapf(errors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}
