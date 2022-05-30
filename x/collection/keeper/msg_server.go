package keeper

import (
	"context"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/collection"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServer returns an implementation of the collection MsgServer interface
// for the provided Keeper.
func NewMsgServer(keeper Keeper) collection.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

var _ collection.MsgServer = (*msgServer)(nil)

func (s msgServer) Send(c context.Context, req *collection.MsgSend) (*collection.MsgSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.keeper.SendCoins(ctx, req.ContractId, sdk.AccAddress(req.From), sdk.AccAddress(req.To), req.Amount); err != nil {
		return nil, err
	}

	event := collection.EventSent{
		ContractId: req.ContractId,
		Operator:   req.From,
		From:       req.From,
		To:         req.To,
		Amount:     req.Amount,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		return nil, err
	}
	return &collection.MsgSendResponse{}, nil
}

func (s msgServer) OperatorSend(c context.Context, req *collection.MsgOperatorSend) (*collection.MsgOperatorSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	// TODO: check authorization
	if req.Operator != req.From {
		return nil, sdkerrors.ErrNotSupported
	}
	if err := s.keeper.SendCoins(ctx, req.ContractId, sdk.AccAddress(req.From), sdk.AccAddress(req.To), req.Amount); err != nil {
		return nil, err
	}

	event := collection.EventSent{
		ContractId: req.ContractId,
		Operator:   req.Operator,
		From:       req.From,
		To:         req.To,
		Amount:     req.Amount,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		return nil, err
	}
	return &collection.MsgOperatorSendResponse{}, nil
}

func (s msgServer) TransferFT(c context.Context, req *collection.MsgTransferFT) (*collection.MsgTransferFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.keeper.SendCoins(ctx, req.ContractId, sdk.AccAddress(req.From), sdk.AccAddress(req.To), req.Amount); err != nil {
		return nil, err
	}

	event := collection.EventSent{
		ContractId: req.ContractId,
		Operator:   req.From,
		From:       req.From,
		To:         req.To,
		Amount:     req.Amount,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		return nil, err
	}
	return &collection.MsgTransferFTResponse{}, nil
}

func (s msgServer) TransferFTFrom(c context.Context, req *collection.MsgTransferFTFrom) (*collection.MsgTransferFTFromResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	// TODO: check authorization
	if req.Proxy != req.From {
		return nil, sdkerrors.ErrNotSupported
	}
	if err := s.keeper.SendCoins(ctx, req.ContractId, sdk.AccAddress(req.From), sdk.AccAddress(req.To), req.Amount); err != nil {
		return nil, err
	}

	event := collection.EventSent{
		ContractId: req.ContractId,
		Operator:   req.Proxy,
		From:       req.From,
		To:         req.To,
		Amount:     req.Amount,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		return nil, err
	}
	return &collection.MsgTransferFTFromResponse{}, nil
}

func (s msgServer) TransferNFT(c context.Context, req *collection.MsgTransferNFT) (*collection.MsgTransferNFTResponse, error) {
	amount := make([]collection.Coin, len(req.TokenIds))
	for i, id := range req.TokenIds {
		amount[i] = collection.Coin{TokenId: id, Amount: sdk.OneInt()}
	}

	ctx := sdk.UnwrapSDKContext(c)
	if err := s.keeper.SendCoins(ctx, req.ContractId, sdk.AccAddress(req.From), sdk.AccAddress(req.To), amount); err != nil {
		return nil, err
	}

	event := collection.EventSent{
		ContractId: req.ContractId,
		Operator:   req.From,
		From:       req.From,
		To:         req.To,
		Amount:     amount,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		return nil, err
	}
	return &collection.MsgTransferNFTResponse{}, nil
}

func (s msgServer) TransferNFTFrom(c context.Context, req *collection.MsgTransferNFTFrom) (*collection.MsgTransferNFTFromResponse, error) {
	amount := make([]collection.Coin, len(req.TokenIds))
	for i, id := range req.TokenIds {
		amount[i] = collection.Coin{TokenId: id, Amount: sdk.OneInt()}
	}

	ctx := sdk.UnwrapSDKContext(c)
	// TODO: check authorization
	if req.Proxy != req.From {
		return nil, sdkerrors.ErrNotSupported
	}
	if err := s.keeper.SendCoins(ctx, req.ContractId, sdk.AccAddress(req.From), sdk.AccAddress(req.To), amount); err != nil {
		return nil, err
	}

	event := collection.EventSent{
		ContractId: req.ContractId,
		Operator:   req.Proxy,
		From:       req.From,
		To:         req.To,
		Amount:     amount,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		return nil, err
	}
	return &collection.MsgTransferNFTFromResponse{}, nil
}

func (s msgServer) AuthorizeOperator(c context.Context, req *collection.MsgAuthorizeOperator) (*collection.MsgAuthorizeOperatorResponse, error) {
	return nil, sdkerrors.ErrNotSupported
}

func (s msgServer) RevokeOperator(c context.Context, req *collection.MsgRevokeOperator) (*collection.MsgRevokeOperatorResponse, error) {
	return nil, sdkerrors.ErrNotSupported
}

func (s msgServer) CreateContract(c context.Context, req *collection.MsgCreateContract) (*collection.MsgCreateContractResponse, error) {
	return nil, sdkerrors.ErrNotSupported
}

func (s msgServer) IssueFT(c context.Context, req *collection.MsgIssueFT) (*collection.MsgIssueFTResponse, error) {
	return nil, sdkerrors.ErrNotSupported
}

func (s msgServer) IssueNFT(c context.Context, req *collection.MsgIssueNFT) (*collection.MsgIssueNFTResponse, error) {
	return nil, sdkerrors.ErrNotSupported
}

func (s msgServer) MintFT(c context.Context, req *collection.MsgMintFT) (*collection.MsgMintFTResponse, error) {
	return nil, sdkerrors.ErrNotSupported
}

func (s msgServer) MintNFT(c context.Context, req *collection.MsgMintNFT) (*collection.MsgMintNFTResponse, error) {
	return nil, sdkerrors.ErrNotSupported
}

func (s msgServer) BurnFT(c context.Context, req *collection.MsgBurnFT) (*collection.MsgBurnFTResponse, error) {
	return nil, sdkerrors.ErrNotSupported
}

func (s msgServer) BurnFTFrom(c context.Context, req *collection.MsgBurnFTFrom) (*collection.MsgBurnFTFromResponse, error) {
	return nil, sdkerrors.ErrNotSupported
}

func (s msgServer) BurnNFT(c context.Context, req *collection.MsgBurnNFT) (*collection.MsgBurnNFTResponse, error) {
	return nil, sdkerrors.ErrNotSupported
}

func (s msgServer) BurnNFTFrom(c context.Context, req *collection.MsgBurnNFTFrom) (*collection.MsgBurnNFTFromResponse, error) {
	return nil, sdkerrors.ErrNotSupported
}

func (s msgServer) Burn(c context.Context, req *collection.MsgBurn) (*collection.MsgBurnResponse, error) {
	return nil, sdkerrors.ErrNotSupported
}

func (s msgServer) OperatorBurn(c context.Context, req *collection.MsgOperatorBurn) (*collection.MsgOperatorBurnResponse, error) {
	return nil, sdkerrors.ErrNotSupported
}

func (s msgServer) Modify(c context.Context, req *collection.MsgModify) (*collection.MsgModifyResponse, error) {
	return nil, sdkerrors.ErrNotSupported
}

func (s msgServer) Grant(c context.Context, req *collection.MsgGrant) (*collection.MsgGrantResponse, error) {
	return nil, sdkerrors.ErrNotSupported
}

func (s msgServer) Abandon(c context.Context, req *collection.MsgAbandon) (*collection.MsgAbandonResponse, error) {
	return nil, sdkerrors.ErrNotSupported
}

func (s msgServer) Attach(c context.Context, req *collection.MsgAttach) (*collection.MsgAttachResponse, error) {
	return nil, sdkerrors.ErrNotSupported
}

func (s msgServer) Detach(c context.Context, req *collection.MsgDetach) (*collection.MsgDetachResponse, error) {
	return nil, sdkerrors.ErrNotSupported
}

func (s msgServer) OperatorAttach(c context.Context, req *collection.MsgOperatorAttach) (*collection.MsgOperatorAttachResponse, error) {
	return nil, sdkerrors.ErrNotSupported
}

func (s msgServer) OperatorDetach(c context.Context, req *collection.MsgOperatorDetach) (*collection.MsgOperatorDetachResponse, error) {
	return nil, sdkerrors.ErrNotSupported
}
