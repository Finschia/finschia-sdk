package evidence

import (
	sdk "github.com/line/lbm-sdk/v2/types"
	sdkerrors "github.com/line/lbm-sdk/v2/types/errors"
	"github.com/line/lbm-sdk/v2/x/evidence/keeper"
	"github.com/line/lbm-sdk/v2/x/evidence/types"
)

// NewHandler returns a handler for evidence messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgSubmitEvidence:
			res, err := msgServer.SubmitEvidence(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}
