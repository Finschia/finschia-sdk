package keeper

import (
	sdk "github.com/line/lbm-sdk/types"

	icatypes "github.com/line/lbm-sdk/x/ibc/applications/27-interchain-accounts/types"
	"github.com/line/lbm-sdk/x/ibc/core/exported"
)

// EmitWriteErrorAcknowledgementEvent emits an event signalling an error acknowledgement and including the error details
func EmitWriteErrorAcknowledgementEvent(ctx sdk.Context, packet exported.PacketI, err error) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			icatypes.EventTypePacket,
			sdk.NewAttribute(sdk.AttributeKeyModule, icatypes.ModuleName),
			sdk.NewAttribute(icatypes.AttributeKeyAckError, err.Error()),
			sdk.NewAttribute(icatypes.AttributeKeyHostChannelID, packet.GetDestChannel()),
		),
	)
}
