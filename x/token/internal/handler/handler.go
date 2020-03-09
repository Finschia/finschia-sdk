package handler

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/contract"
	"github.com/line/link/x/token/internal/keeper"
	"github.com/line/link/x/token/internal/types"
)

func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		if msg, ok := msg.(contract.Msg); ok {
			err := handleMsgContract(ctx, keeper, msg)
			if err != nil {
				return err.Result()
			}
		}
		switch msg := msg.(type) {
		case types.MsgIssue:
			return handleMsgIssue(ctx, keeper, msg)
		case types.MsgMint:
			return handleMsgMint(ctx, keeper, msg)
		case types.MsgBurn:
			return handleMsgBurn(ctx, keeper, msg)
		case types.MsgTransfer:
			return handleMsgTransfer(ctx, keeper, msg)
		case types.MsgGrantPermission:
			return handleMsgGrant(ctx, keeper, msg)
		case types.MsgRevokePermission:
			return handleMsgRevoke(ctx, keeper, msg)
		case types.MsgModify:
			return handleMsgModify(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized  Msg type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}
func handleMsgContract(ctx sdk.Context, keeper keeper.Keeper, msg contract.Msg) sdk.Error {
	if !keeper.HasContractID(ctx, msg.GetContractID()) {
		return contract.ErrContractNotExist(contract.ContractCodeSpace, msg.GetContractID())
	}
	return nil
}
