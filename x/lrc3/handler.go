package lrc3

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/link-chain/link/x/lrc3/internal/keeper"
	"github.com/link-chain/link/x/lrc3/internal/types"
	nft "github.com/link-chain/link/x/nft"
)

// NewHandler routes the messages to the handlers
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgInit:
			return HandleMsgInit(ctx, msg, k)
		case types.MsgMintNFT:
			return HandleMsgMintNFT(ctx, msg, k)
		case nft.MsgBurnNFT:
			return HandleMsgBurnNFT(ctx, msg, k)
		case types.MsgTransfer:
			return HandleMsgTransferNFT(ctx, msg, k)
		case types.MsgApprove:
			return HandleMsgApprove(ctx, msg, k)
		case types.MsgSetApprovalForAll:
			return HandleMsgSetApprovalForAll(ctx, msg, k)
		default:
			errMsg := fmt.Sprintf("Unrecognized message type: %T", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// HandleMsgInit handler for MsgInit
func HandleMsgInit(ctx sdk.Context, msg types.MsgInit, k keeper.Keeper,
) sdk.Result {
	collection := nft.NewCollection(msg.Denom, nft.NewNFTs())

	err := k.SetLRC3(ctx, msg.Denom, collection)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeInit,
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyOwnerAddress, msg.OwnerAddress.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

// HandleMsgMintNFT handles MsgMintNFT
func HandleMsgMintNFT(ctx sdk.Context, msg types.MsgMintNFT, k keeper.Keeper,
) sdk.Result {
	lrc3, found := k.NFTKeeper.GetCollection(ctx, msg.Denom)
	if !found {
		return types.ErrNotExistLRC3(types.DefaultCodespace).Result()
	}

	len := len(lrc3.NFTs)
	var tokenId uint64
	if len == 0 {
		tokenId = 0
	} else {
		latestNFT := lrc3.NFTs[len-1]
		tokenId, _ = strconv.ParseUint(latestNFT.GetID(), 10, 64)
		tokenId++
	}
	msgMintNFT := nft.NewMsgMintNFT(msg.Sender, msg.Recipient, strconv.FormatUint(tokenId, 10), msg.Denom, msg.TokenURI)

	return nft.HandleMsgMintNFT(ctx, msgMintNFT, k.NFTKeeper)
}

// HandleMsgBurnNFT handles MsgBurnNFT
func HandleMsgBurnNFT(ctx sdk.Context, msg nft.MsgBurnNFT, k keeper.Keeper,
) sdk.Result {
	err := checkApprove(ctx, msg.Denom, msg.ID, msg.Sender, k)
	if err != nil {
		return err.Result()
	}
	return nft.HandleMsgBurnNFT(ctx, msg, k.NFTKeeper)
}

// HandleMsgTransferNFT handler for MsgTransferNFT
func HandleMsgTransferNFT(ctx sdk.Context, msg types.MsgTransfer, k keeper.Keeper,
) sdk.Result {
	err := checkApprove(ctx, msg.MsgTransferNFT.Denom, msg.MsgTransferNFT.ID, msg.MsgTransferNFT.Sender, k)
	if err != nil {
		return err.Result()
	}
	return nft.HandleMsgTransferNFT(ctx, msg.MsgTransferNFT, k.NFTKeeper)
}

// HandleMsgApprove handles MsgTokenApprove
func HandleMsgApprove(ctx sdk.Context, msg types.MsgApprove, k keeper.Keeper) sdk.Result {

	denom := msg.Denom
	sender := msg.Sender
	recipient := msg.Recipient

	nft, err := k.NFTKeeper.GetNFT(ctx, denom, msg.TokenID)
	if err != nil {
		return err.Result()
	}
	owner := nft.GetOwner()

	if owner.Equals(recipient) {
		return sdk.ErrInvalidAddress("Recipient address and token owner are the same.").Result()
	}

	operators, _ := k.GetOperatorApprovals(ctx, denom, owner)
	found := operators.Find(sender)

	if owner.Equals(sender) == false && found == false {
		return types.ErrInvalidPermission(types.DefaultCodespace).Result()
	}

	k.SetApproval(ctx, denom, msg.TokenID, recipient)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeApprove,
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyTokenId, msg.TokenID),
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Recipient.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// HandleMsgSetApprovalForAll handles MsgSetApprovalForAll
func HandleMsgSetApprovalForAll(ctx sdk.Context, msg types.MsgSetApprovalForAll, k keeper.Keeper) sdk.Result {
	denom := msg.Denom
	owner := msg.Owner
	operator := msg.Operator
	approved := msg.Approved

	operators, _ := k.GetOperatorApprovals(ctx, denom, owner)
	if approved == true {
		operators = append(operators, operator)
	} else {
		operators, _ = operators.DeleteOperator(operator)
	}

	k.SetOperatorApprovals(ctx, denom, owner, operators)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSetApprovalForAll,
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
			sdk.NewAttribute(types.AttributeKeyOperator, msg.Operator.String()),
			sdk.NewAttribute(types.AttributeKeyTo, strconv.FormatBool(msg.Approved)),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func checkApprove(ctx sdk.Context, denom, id string, sender sdk.AccAddress, k keeper.Keeper) sdk.Error {
	token, err := k.NFTKeeper.GetNFT(ctx, denom, id)
	if err != nil {
		return err
	}
	approve := k.IsApprovedOrOwner(ctx, sender, denom, token)
	if approve == false {
		return types.ErrNotExistApprove(types.DefaultCodespace)
	}
	return nil
}
