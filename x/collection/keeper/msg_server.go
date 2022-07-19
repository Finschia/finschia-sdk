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

	// emit legacy events.
	event := collection.EventSent{
		ContractId: req.ContractId,
		Operator:   req.From,
		From:       req.From,
		To:         req.To,
		Amount:     req.Amount,
	}
	if legacyEvent := collection.NewEventTransferFT(event); legacyEvent != nil {
		ctx.EventManager().EmitEvent(*legacyEvent)
	}
	ctx.EventManager().EmitEvents(collection.NewEventTransferNFT(event))

	fromAddr, _ := sdk.AccAddressFromBech32(req.From)

	toAddr, _ := sdk.AccAddressFromBech32(req.To)

	if err := s.keeper.SendCoins(ctx, req.ContractId, fromAddr, toAddr, req.Amount); err != nil {
		return nil, err
	}

	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgSendResponse{}, nil
}

func (s msgServer) OperatorSend(c context.Context, req *collection.MsgOperatorSend) (*collection.MsgOperatorSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	fromAddr, _ := sdk.AccAddressFromBech32(req.From)

	operatorAddr, _ := sdk.AccAddressFromBech32(req.Operator)

	if _, err := s.keeper.GetAuthorization(ctx, req.ContractId, fromAddr, operatorAddr); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	// emit legacy events.
	event := collection.EventSent{
		ContractId: req.ContractId,
		Operator:   req.Operator,
		From:       req.From,
		To:         req.To,
		Amount:     req.Amount,
	}
	if legacyEvent := collection.NewEventTransferFTFrom(event); legacyEvent != nil {
		ctx.EventManager().EmitEvent(*legacyEvent)
	}
	ctx.EventManager().EmitEvents(collection.NewEventTransferNFTFrom(event))

	toAddr, _ := sdk.AccAddressFromBech32(req.To)

	if err := s.keeper.SendCoins(ctx, req.ContractId, fromAddr, toAddr, req.Amount); err != nil {
		return nil, err
	}

	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgOperatorSendResponse{}, nil
}

func (s msgServer) TransferFT(c context.Context, req *collection.MsgTransferFT) (*collection.MsgTransferFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	fromAddr, _ := sdk.AccAddressFromBech32(req.From)

	toAddr, _ := sdk.AccAddressFromBech32(req.To)

	if err := s.keeper.SendCoins(ctx, req.ContractId, fromAddr, toAddr, req.Amount); err != nil {
		return nil, err
	}

	event := collection.EventSent{
		ContractId: req.ContractId,
		Operator:   req.From,
		From:       req.From,
		To:         req.To,
		Amount:     req.Amount,
	}
	ctx.EventManager().EmitEvent(*collection.NewEventTransferFT(event))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgTransferFTResponse{}, nil
}

func (s msgServer) TransferFTFrom(c context.Context, req *collection.MsgTransferFTFrom) (*collection.MsgTransferFTFromResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	fromAddr, _ := sdk.AccAddressFromBech32(req.From)

	proxyAddr, _ := sdk.AccAddressFromBech32(req.Proxy)

	if _, err := s.keeper.GetAuthorization(ctx, req.ContractId, fromAddr, proxyAddr); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	toAddr, _ := sdk.AccAddressFromBech32(req.To)

	if err := s.keeper.SendCoins(ctx, req.ContractId, fromAddr, toAddr, req.Amount); err != nil {
		return nil, err
	}

	event := collection.EventSent{
		ContractId: req.ContractId,
		Operator:   req.Proxy,
		From:       req.From,
		To:         req.To,
		Amount:     req.Amount,
	}
	ctx.EventManager().EmitEvent(*collection.NewEventTransferFTFrom(event))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgTransferFTFromResponse{}, nil
}

func (s msgServer) TransferNFT(c context.Context, req *collection.MsgTransferNFT) (*collection.MsgTransferNFTResponse, error) {
	amount := make([]collection.Coin, len(req.TokenIds))
	for i, id := range req.TokenIds {
		amount[i] = collection.Coin{TokenId: id, Amount: sdk.OneInt()}
	}

	ctx := sdk.UnwrapSDKContext(c)

	// emit legacy events
	event := collection.EventSent{
		ContractId: req.ContractId,
		Operator:   req.From,
		From:       req.From,
		To:         req.To,
		Amount:     amount,
	}
	ctx.EventManager().EmitEvents(collection.NewEventTransferNFT(event))

	fromAddr, _ := sdk.AccAddressFromBech32(req.From)

	toAddr, _ := sdk.AccAddressFromBech32(req.To)

	if err := s.keeper.SendCoins(ctx, req.ContractId, fromAddr, toAddr, amount); err != nil {
		return nil, err
	}

	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgTransferNFTResponse{}, nil
}

func (s msgServer) TransferNFTFrom(c context.Context, req *collection.MsgTransferNFTFrom) (*collection.MsgTransferNFTFromResponse, error) {
	amount := make([]collection.Coin, len(req.TokenIds))
	for i, id := range req.TokenIds {
		amount[i] = collection.Coin{TokenId: id, Amount: sdk.OneInt()}
	}

	ctx := sdk.UnwrapSDKContext(c)

	fromAddr, _ := sdk.AccAddressFromBech32(req.From)

	proxyAddr, _ := sdk.AccAddressFromBech32(req.Proxy)

	if _, err := s.keeper.GetAuthorization(ctx, req.ContractId, fromAddr, proxyAddr); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	// emit legacy events
	event := collection.EventSent{
		ContractId: req.ContractId,
		Operator:   req.Proxy,
		From:       req.From,
		To:         req.To,
		Amount:     amount,
	}
	ctx.EventManager().EmitEvents(collection.NewEventTransferNFTFrom(event))
	toAddr, _ := sdk.AccAddressFromBech32(req.To)

	if err := s.keeper.SendCoins(ctx, req.ContractId, fromAddr, toAddr, amount); err != nil {
		return nil, err
	}

	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgTransferNFTFromResponse{}, nil
}

func (s msgServer) AuthorizeOperator(c context.Context, req *collection.MsgAuthorizeOperator) (*collection.MsgAuthorizeOperatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	holderAddr, _ := sdk.AccAddressFromBech32(req.Holder)

	operatorAddr, _ := sdk.AccAddressFromBech32(req.Operator)

	if err := s.keeper.AuthorizeOperator(ctx, req.ContractId, holderAddr, operatorAddr); err != nil {
		return nil, err
	}

	event := collection.EventAuthorizedOperator{
		ContractId: req.ContractId,
		Holder:     req.Holder,
		Operator:   req.Operator,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgAuthorizeOperatorResponse{}, nil
}

func (s msgServer) RevokeOperator(c context.Context, req *collection.MsgRevokeOperator) (*collection.MsgRevokeOperatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	holderAddr, _ := sdk.AccAddressFromBech32(req.Holder)

	operatorAddr, _ := sdk.AccAddressFromBech32(req.Operator)

	if err := s.keeper.RevokeOperator(ctx, req.ContractId, holderAddr, operatorAddr); err != nil {
		return nil, err
	}

	event := collection.EventRevokedOperator{
		ContractId: req.ContractId,
		Holder:     req.Holder,
		Operator:   req.Operator,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgRevokeOperatorResponse{}, nil
}

func (s msgServer) Approve(c context.Context, req *collection.MsgApprove) (*collection.MsgApproveResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	approverAddr, _ := sdk.AccAddressFromBech32(req.Approver)

	proxyAddr, _ := sdk.AccAddressFromBech32(req.Proxy)

	if err := s.keeper.AuthorizeOperator(ctx, req.ContractId, approverAddr, proxyAddr); err != nil {
		return nil, err
	}

	event := collection.EventAuthorizedOperator{
		ContractId: req.ContractId,
		Holder:     req.Approver,
		Operator:   req.Proxy,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgApproveResponse{}, nil
}

func (s msgServer) Disapprove(c context.Context, req *collection.MsgDisapprove) (*collection.MsgDisapproveResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	approverAddr, _ := sdk.AccAddressFromBech32(req.Approver)

	proxyAddr, _ := sdk.AccAddressFromBech32(req.Proxy)

	if err := s.keeper.RevokeOperator(ctx, req.ContractId, approverAddr, proxyAddr); err != nil {
		return nil, err
	}

	event := collection.EventRevokedOperator{
		ContractId: req.ContractId,
		Holder:     req.Approver,
		Operator:   req.Proxy,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgDisapproveResponse{}, nil
}

func (s msgServer) CreateContract(c context.Context, req *collection.MsgCreateContract) (*collection.MsgCreateContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	contract := collection.Contract{
		Name:       req.Name,
		BaseImgUri: req.BaseImgUri,
		Meta:       req.Meta,
	}

	ownerAddr, _ := sdk.AccAddressFromBech32(req.Owner)

	id := s.keeper.CreateContract(ctx, ownerAddr, contract)

	return &collection.MsgCreateContractResponse{Id: id}, nil
}

func (s msgServer) CreateFTClass(c context.Context, req *collection.MsgCreateFTClass) (*collection.MsgCreateFTClassResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	operatorAddr, _ := sdk.AccAddressFromBech32(req.Operator)

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, operatorAddr, collection.PermissionIssue); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	mintable := !req.Supply.IsPositive()
	class := &collection.FTClass{
		Name:     req.Name,
		Meta:     req.Meta,
		Decimals: req.Decimals,
		Mintable: mintable,
	}
	id, err := s.keeper.CreateTokenClass(ctx, req.ContractId, class)
	if err != nil {
		return nil, err
	}

	event := collection.EventCreatedFTClass{
		ContractId: req.ContractId,
		ClassId:    *id,
		Name:       class.Name,
		Meta:       class.Meta,
		Decimals:   class.Decimals,
		Mintable:   class.Mintable,
	}

	toAddr, _ := sdk.AccAddressFromBech32(req.To)

	ctx.EventManager().EmitEvent(collection.NewEventIssueFT(event, operatorAddr, toAddr, req.Supply))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	// supply tokens
	if req.Supply.IsPositive() {
		s.keeper.mintFT(ctx, req.ContractId, toAddr, *id, req.Supply)

		event := collection.EventMintedFT{
			ContractId: req.ContractId,
			Operator:   req.Operator,
			To:         req.To,
			Amount:     collection.NewCoins(collection.NewFTCoin(*id, req.Supply)),
		}
		if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
			panic(err)
		}
	}

	return &collection.MsgCreateFTClassResponse{Id: *id}, nil
}

func (s msgServer) CreateNFTClass(c context.Context, req *collection.MsgCreateNFTClass) (*collection.MsgCreateNFTClassResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	operatorAddr, _ := sdk.AccAddressFromBech32(req.Operator)

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, operatorAddr, collection.PermissionIssue); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	class := &collection.NFTClass{
		Name: req.Name,
		Meta: req.Meta,
	}
	id, err := s.keeper.CreateTokenClass(ctx, req.ContractId, class)
	if err != nil {
		return nil, err
	}

	event := collection.EventCreatedNFTClass{
		ContractId: req.ContractId,
		ClassId:    *id,
		Name:       class.Name,
		Meta:       class.Meta,
	}
	ctx.EventManager().EmitEvent(collection.NewEventIssueNFT(event))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgCreateNFTClassResponse{Id: *id}, nil
}

func (s msgServer) IssueFT(c context.Context, req *collection.MsgIssueFT) (*collection.MsgIssueFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	ownerAddr, _ := sdk.AccAddressFromBech32(req.Owner)

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, ownerAddr, collection.PermissionIssue); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	class := &collection.FTClass{
		Name:     req.Name,
		Meta:     req.Meta,
		Decimals: req.Decimals,
		Mintable: req.Mintable,
	}
	id, err := s.keeper.CreateTokenClass(ctx, req.ContractId, class)
	if err != nil {
		return nil, err
	}

	event := collection.EventCreatedFTClass{
		ContractId: req.ContractId,
		ClassId:    *id,
		Name:       class.Name,
		Meta:       class.Meta,
		Decimals:   class.Decimals,
		Mintable:   class.Mintable,
	}

	toAddr, _ := sdk.AccAddressFromBech32(req.To)

	ctx.EventManager().EmitEvent(collection.NewEventIssueFT(event, ownerAddr, toAddr, req.Amount))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	// supply tokens
	if req.Amount.IsPositive() {
		s.keeper.mintFT(ctx, req.ContractId, toAddr, *id, req.Amount)

		event := collection.EventMintedFT{
			ContractId: req.ContractId,
			Operator:   req.Owner,
			To:         req.To,
			Amount:     collection.NewCoins(collection.NewFTCoin(*id, req.Amount)),
		}
		if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
			panic(err)
		}
	}

	return &collection.MsgIssueFTResponse{Id: *id}, nil
}

func (s msgServer) IssueNFT(c context.Context, req *collection.MsgIssueNFT) (*collection.MsgIssueNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	ownerAddr, _ := sdk.AccAddressFromBech32(req.Owner)

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, ownerAddr, collection.PermissionIssue); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	class := &collection.NFTClass{
		Name: req.Name,
		Meta: req.Meta,
	}
	id, err := s.keeper.CreateTokenClass(ctx, req.ContractId, class)
	if err != nil {
		return nil, err
	}

	event := collection.EventCreatedNFTClass{
		ContractId: req.ContractId,
		ClassId:    *id,
		Name:       class.Name,
		Meta:       class.Meta,
	}
	ctx.EventManager().EmitEvent(collection.NewEventIssueNFT(event))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	for _, permission := range []collection.Permission{
		collection.PermissionMint,
		collection.PermissionBurn,
	} {
		s.keeper.Grant(ctx, req.ContractId, []byte{}, ownerAddr, permission)
	}

	return &collection.MsgIssueNFTResponse{Id: *id}, nil
}

func (s msgServer) MintFT(c context.Context, req *collection.MsgMintFT) (*collection.MsgMintFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	fromAddr, _ := sdk.AccAddressFromBech32(req.From)

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, fromAddr, collection.PermissionMint); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	toAddr, _ := sdk.AccAddressFromBech32(req.To)

	if err := s.keeper.MintFT(ctx, req.ContractId, toAddr, req.Amount); err != nil {
		return nil, err
	}

	event := collection.EventMintedFT{
		ContractId: req.ContractId,
		Operator:   req.From,
		To:         req.To,
		Amount:     req.Amount,
	}
	ctx.EventManager().EmitEvent(collection.NewEventMintFT(event))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgMintFTResponse{}, nil
}

func (s msgServer) MintNFT(c context.Context, req *collection.MsgMintNFT) (*collection.MsgMintNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	fromAddr, _ := sdk.AccAddressFromBech32(req.From)

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, fromAddr, collection.PermissionMint); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	toAddr, _ := sdk.AccAddressFromBech32(req.To)

	tokens, err := s.keeper.MintNFT(ctx, req.ContractId, toAddr, req.Params)
	if err != nil {
		return nil, err
	}

	event := collection.EventMintedNFT{
		ContractId: req.ContractId,
		Operator:   req.From,
		To:         req.To,
		Tokens:     tokens,
	}
	ctx.EventManager().EmitEvents(collection.NewEventMintNFT(event))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	tokenIDs := make([]string, 0, len(tokens))
	for _, token := range tokens {
		tokenIDs = append(tokenIDs, token.Id)
	}
	return &collection.MsgMintNFTResponse{Ids: tokenIDs}, nil
}

func (s msgServer) Burn(c context.Context, req *collection.MsgBurn) (*collection.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	fromAddr, _ := sdk.AccAddressFromBech32(req.From)

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, fromAddr, collection.PermissionBurn); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	// legacy: emit events against the original request.
	event := collection.EventBurned{
		ContractId: req.ContractId,
		Operator:   req.From,
		From:       req.From,
		Amount:     req.Amount,
	}
	if e := collection.NewEventBurnFT(event); e != nil {
		ctx.EventManager().EmitEvent(*e)
	}
	ctx.EventManager().EmitEvents(collection.NewEventBurnNFT(event))

	burnt, err := s.keeper.BurnCoins(ctx, req.ContractId, fromAddr, req.Amount)
	if err != nil {
		return nil, err
	}

	// emit events against all burnt tokens.
	event.Amount = burnt
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgBurnResponse{}, nil
}

func (s msgServer) OperatorBurn(c context.Context, req *collection.MsgOperatorBurn) (*collection.MsgOperatorBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	fromAddr, _ := sdk.AccAddressFromBech32(req.From)

	operatorAddr, _ := sdk.AccAddressFromBech32(req.Operator)

	if _, err := s.keeper.GetAuthorization(ctx, req.ContractId, fromAddr, operatorAddr); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, operatorAddr, collection.PermissionBurn); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	// legacy: emit events against the original request.
	event := collection.EventBurned{
		ContractId: req.ContractId,
		Operator:   req.Operator,
		From:       req.From,
		Amount:     req.Amount,
	}
	if e := collection.NewEventBurnFTFrom(event); e != nil {
		ctx.EventManager().EmitEvent(*e)
	}
	ctx.EventManager().EmitEvents(collection.NewEventBurnNFTFrom(event))

	burnt, err := s.keeper.BurnCoins(ctx, req.ContractId, fromAddr, req.Amount)
	if err != nil {
		return nil, err
	}

	// emit events against all burnt tokens.
	event.Amount = burnt
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgOperatorBurnResponse{}, nil
}

func (s msgServer) BurnFT(c context.Context, req *collection.MsgBurnFT) (*collection.MsgBurnFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	fromAddr, _ := sdk.AccAddressFromBech32(req.From)

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, fromAddr, collection.PermissionBurn); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	burnt, err := s.keeper.BurnCoins(ctx, req.ContractId, fromAddr, req.Amount)
	if err != nil {
		return nil, err
	}

	event := collection.EventBurned{
		ContractId: req.ContractId,
		Operator:   req.From,
		From:       req.From,
		Amount:     burnt,
	}
	if e := collection.NewEventBurnFT(event); e != nil {
		ctx.EventManager().EmitEvent(*e)
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgBurnFTResponse{}, nil
}

func (s msgServer) BurnFTFrom(c context.Context, req *collection.MsgBurnFTFrom) (*collection.MsgBurnFTFromResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	fromAddr, _ := sdk.AccAddressFromBech32(req.From)

	proxyAddr, _ := sdk.AccAddressFromBech32(req.Proxy)

	if _, err := s.keeper.GetAuthorization(ctx, req.ContractId, fromAddr, proxyAddr); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, proxyAddr, collection.PermissionBurn); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	burnt, err := s.keeper.BurnCoins(ctx, req.ContractId, fromAddr, req.Amount)
	if err != nil {
		return nil, err
	}

	event := collection.EventBurned{
		ContractId: req.ContractId,
		Operator:   req.Proxy,
		From:       req.From,
		Amount:     burnt,
	}
	if e := collection.NewEventBurnFTFrom(event); e != nil {
		ctx.EventManager().EmitEvent(*e)
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgBurnFTFromResponse{}, nil
}

func (s msgServer) BurnNFT(c context.Context, req *collection.MsgBurnNFT) (*collection.MsgBurnNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	fromAddr, _ := sdk.AccAddressFromBech32(req.From)

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, fromAddr, collection.PermissionBurn); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	coins := make([]collection.Coin, 0, len(req.TokenIds))
	for _, id := range req.TokenIds {
		coins = append(coins, collection.NewCoin(id, sdk.OneInt()))
	}

	// legacy: emit events against the original request.
	event := collection.EventBurned{
		ContractId: req.ContractId,
		Operator:   req.From,
		From:       req.From,
		Amount:     coins,
	}
	ctx.EventManager().EmitEvents(collection.NewEventBurnNFT(event))

	burnt, err := s.keeper.BurnCoins(ctx, req.ContractId, fromAddr, coins)
	if err != nil {
		return nil, err
	}

	// emit events against all burnt tokens.
	event.Amount = burnt
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgBurnNFTResponse{}, nil
}

func (s msgServer) BurnNFTFrom(c context.Context, req *collection.MsgBurnNFTFrom) (*collection.MsgBurnNFTFromResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	fromAddr, _ := sdk.AccAddressFromBech32(req.From)

	proxyAddr, _ := sdk.AccAddressFromBech32(req.Proxy)

	if _, err := s.keeper.GetAuthorization(ctx, req.ContractId, fromAddr, proxyAddr); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, proxyAddr, collection.PermissionBurn); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	coins := make([]collection.Coin, 0, len(req.TokenIds))
	for _, id := range req.TokenIds {
		coins = append(coins, collection.NewCoin(id, sdk.OneInt()))
	}

	// legacy: emit events against the original request.
	event := collection.EventBurned{
		ContractId: req.ContractId,
		Operator:   req.Proxy,
		From:       req.From,
		Amount:     coins,
	}
	ctx.EventManager().EmitEvents(collection.NewEventBurnNFTFrom(event))

	burnt, err := s.keeper.BurnCoins(ctx, req.ContractId, fromAddr, coins)
	if err != nil {
		return nil, err
	}

	// emit events against all burnt tokens.
	event.Amount = burnt
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgBurnNFTFromResponse{}, nil
}

func (s msgServer) ModifyContract(c context.Context, req *collection.MsgModifyContract) (*collection.MsgModifyContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	operator, _ := sdk.AccAddressFromBech32(req.Operator)

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, operator, collection.PermissionModify); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	if err := s.keeper.ModifyContract(ctx, req.ContractId, operator, req.Changes); err != nil {
		return nil, err
	}

	event := collection.EventModifiedContract{
		ContractId: req.ContractId,
		Operator:   req.Operator,
		Changes:    req.Changes,
	}
	ctx.EventManager().EmitEvents(collection.NewEventModifyCollection(event))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgModifyContractResponse{}, nil
}

func (s msgServer) ModifyTokenClass(c context.Context, req *collection.MsgModifyTokenClass) (*collection.MsgModifyTokenClassResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	operator, _ := sdk.AccAddressFromBech32(req.Operator)

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, operator, collection.PermissionModify); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	if err := s.keeper.ModifyTokenClass(ctx, req.ContractId, req.ClassId, operator, req.Changes); err != nil {
		return nil, err
	}

	event := collection.EventModifiedTokenClass{
		ContractId: req.ContractId,
		Operator:   req.Operator,
		ClassId:    req.ClassId,
		Changes:    req.Changes,
	}
	class, err := s.keeper.GetTokenClass(ctx, req.ContractId, req.ClassId)
	if err != nil {
		panic(err)
	}
	if _, ok := class.(*collection.FTClass); ok {
		ctx.EventManager().EmitEvents(collection.NewEventModifyTokenOfFTClass(event))
	}
	if _, ok := class.(*collection.NFTClass); ok {
		ctx.EventManager().EmitEvents(collection.NewEventModifyTokenType(event))
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgModifyTokenClassResponse{}, nil
}

func (s msgServer) ModifyNFT(c context.Context, req *collection.MsgModifyNFT) (*collection.MsgModifyNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	operator, _ := sdk.AccAddressFromBech32(req.Operator)

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, operator, collection.PermissionModify); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	if err := s.keeper.ModifyNFT(ctx, req.ContractId, req.TokenId, operator, req.Changes); err != nil {
		return nil, err
	}

	event := collection.EventModifiedNFT{
		ContractId: req.ContractId,
		Operator:   req.Operator,
		TokenId:    req.TokenId,
		Changes:    req.Changes,
	}
	ctx.EventManager().EmitEvents(collection.NewEventModifyTokenOfNFT(event))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgModifyNFTResponse{}, nil
}

func (s msgServer) Modify(c context.Context, req *collection.MsgModify) (*collection.MsgModifyResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	operator, _ := sdk.AccAddressFromBech32(req.Owner)

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, operator, collection.PermissionModify); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	// copied from daphne
	modify := func(tokenType, tokenIndex string) error {
		changes := make([]collection.Attribute, len(req.Changes))
		for i, change := range req.Changes {
			changes[i] = collection.Attribute{
				Key:   change.Field,
				Value: change.Field,
			}
		}

		classID := tokenType
		tokenID := classID + tokenIndex
		if tokenType != "" {
			if tokenIndex != "" {
				if collection.ValidateNFTID(tokenID) == nil {
					event := collection.EventModifiedNFT{
						ContractId: req.ContractId,
						Operator:   operator.String(),
						TokenId:    tokenID,
						Changes:    changes,
					}
					ctx.EventManager().EmitEvents(collection.NewEventModifyTokenOfNFT(event))
					if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
						panic(err)
					}

					return s.keeper.ModifyNFT(ctx, req.ContractId, tokenID, operator, changes)
				}

				event := collection.EventModifiedTokenClass{
					ContractId: req.ContractId,
					Operator:   operator.String(),
					ClassId:    classID,
					Changes:    changes,
				}

				ctx.EventManager().EmitEvents(collection.NewEventModifyTokenOfFTClass(event))
				if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
					panic(err)
				}

				return s.keeper.ModifyTokenClass(ctx, req.ContractId, classID, operator, changes)
			}

			event := collection.EventModifiedTokenClass{
				ContractId: req.ContractId,
				Operator:   operator.String(),
				ClassId:    classID,
				Changes:    changes,
			}
			ctx.EventManager().EmitEvents(collection.NewEventModifyTokenType(event))
			if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
				panic(err)
			}

			return s.keeper.ModifyTokenClass(ctx, req.ContractId, classID, operator, changes)
		}
		if req.TokenIndex == "" {
			event := collection.EventModifiedContract{
				ContractId: req.ContractId,
				Operator:   operator.String(),
				Changes:    changes,
			}
			ctx.EventManager().EmitEvents(collection.NewEventModifyCollection(event))
			if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
				panic(err)
			}

			return s.keeper.ModifyContract(ctx, req.ContractId, operator, changes)
		}

		return sdkerrors.ErrInvalidRequest.Wrap("token index without type")
	}

	if err := modify(req.TokenType, req.TokenIndex); err != nil {
		return nil, err
	}

	return &collection.MsgModifyResponse{}, nil
}

func (s msgServer) Grant(c context.Context, req *collection.MsgGrant) (*collection.MsgGrantResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	granter, _ := sdk.AccAddressFromBech32(req.Granter)

	grantee, _ := sdk.AccAddressFromBech32(req.Grantee)

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, granter, req.Permission); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrapf("%s is not authorized for %s", granter, req.Permission)
	}
	if _, err := s.keeper.GetGrant(ctx, req.ContractId, grantee, req.Permission); err == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("%s is already granted for %s", grantee, req.Permission)
	}

	s.keeper.Grant(ctx, req.ContractId, granter, grantee, req.Permission)

	event := collection.EventGrant{
		ContractId: req.ContractId,
		Granter:    granter.String(),
		Grantee:    grantee.String(),
		Permission: req.Permission,
	}
	ctx.EventManager().EmitEvent(collection.NewEventGrantPermToken(event))
	// it emits typed event inside s.keeper.Grant()

	return &collection.MsgGrantResponse{}, nil
}

func (s msgServer) Abandon(c context.Context, req *collection.MsgAbandon) (*collection.MsgAbandonResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	grantee, _ := sdk.AccAddressFromBech32(req.Grantee)

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, grantee, req.Permission); err != nil {
		return nil, sdkerrors.ErrNotFound.Wrapf("%s is not authorized for %s", grantee, req.Permission)
	}

	s.keeper.Abandon(ctx, req.ContractId, grantee, req.Permission)

	event := collection.EventAbandon{
		ContractId: req.ContractId,
		Grantee:    grantee.String(),
		Permission: req.Permission,
	}
	ctx.EventManager().EmitEvent(collection.NewEventRevokePermToken(event))
	// it emits typed event inside s.keeper.Abandon()

	return &collection.MsgAbandonResponse{}, nil
}

func (s msgServer) GrantPermission(c context.Context, req *collection.MsgGrantPermission) (*collection.MsgGrantPermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	granter, _ := sdk.AccAddressFromBech32(req.From)

	grantee, _ := sdk.AccAddressFromBech32(req.To)

	permission := collection.Permission(collection.LegacyPermissionFromString(req.Permission))
	if _, err := s.keeper.GetGrant(ctx, req.ContractId, granter, permission); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrapf("%s is not authorized for %s", granter, permission)
	}
	if _, err := s.keeper.GetGrant(ctx, req.ContractId, grantee, permission); err == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("%s is already granted for %s", grantee, permission)
	}

	s.keeper.Grant(ctx, req.ContractId, granter, grantee, permission)

	event := collection.EventGrant{
		ContractId: req.ContractId,
		Granter:    granter.String(),
		Grantee:    grantee.String(),
		Permission: permission,
	}
	ctx.EventManager().EmitEvent(collection.NewEventGrantPermToken(event))
	// it emits typed event inside s.keeper.Grant()

	return &collection.MsgGrantPermissionResponse{}, nil
}

func (s msgServer) RevokePermission(c context.Context, req *collection.MsgRevokePermission) (*collection.MsgRevokePermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	grantee, _ := sdk.AccAddressFromBech32(req.From)

	permission := collection.Permission(collection.LegacyPermissionFromString(req.Permission))
	if _, err := s.keeper.GetGrant(ctx, req.ContractId, grantee, permission); err != nil {
		return nil, sdkerrors.ErrNotFound.Wrapf("%s is not authorized for %s", grantee, permission)
	}

	s.keeper.Abandon(ctx, req.ContractId, grantee, permission)

	event := collection.EventAbandon{
		ContractId: req.ContractId,
		Grantee:    grantee.String(),
		Permission: permission,
	}
	ctx.EventManager().EmitEvent(collection.NewEventRevokePermToken(event))
	// it emits typed event inside s.keeper.Abandon()

	return &collection.MsgRevokePermissionResponse{}, nil
}

func (s msgServer) Attach(c context.Context, req *collection.MsgAttach) (*collection.MsgAttachResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	event := collection.EventAttached{
		ContractId: req.ContractId,
		Operator:   req.From,
		Holder:     req.From,
		Subject:    req.TokenId,
		Target:     req.ToTokenId,
	}
	newRoot := s.keeper.GetRoot(ctx, req.ContractId, req.ToTokenId)
	ctx.EventManager().EmitEvent(collection.NewEventAttachToken(event, newRoot))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	fromAddr, _ := sdk.AccAddressFromBech32(req.From)

	if err := s.keeper.Attach(ctx, req.ContractId, fromAddr, req.TokenId, req.ToTokenId); err != nil {
		return nil, err
	}

	return &collection.MsgAttachResponse{}, nil
}

func (s msgServer) Detach(c context.Context, req *collection.MsgDetach) (*collection.MsgDetachResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// legacy
	if err := s.keeper.hasNFT(ctx, req.ContractId, req.TokenId); err != nil {
		return nil, err
	}
	oldRoot := s.keeper.GetRoot(ctx, req.ContractId, req.TokenId)

	event := collection.EventDetached{
		ContractId: req.ContractId,
		Operator:   req.From,
		Holder:     req.From,
		Subject:    req.TokenId,
	}
	ctx.EventManager().EmitEvent(collection.NewEventDetachToken(event, oldRoot))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	fromAddr, _ := sdk.AccAddressFromBech32(req.From)

	if err := s.keeper.Detach(ctx, req.ContractId, fromAddr, req.TokenId); err != nil {
		return nil, err
	}

	return &collection.MsgDetachResponse{}, nil
}

func (s msgServer) OperatorAttach(c context.Context, req *collection.MsgOperatorAttach) (*collection.MsgOperatorAttachResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	ownerAddr, _ := sdk.AccAddressFromBech32(req.Owner)

	operatorAddr, _ := sdk.AccAddressFromBech32(req.Operator)

	if _, err := s.keeper.GetAuthorization(ctx, req.ContractId, ownerAddr, operatorAddr); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	event := collection.EventAttached{
		ContractId: req.ContractId,
		Operator:   req.Operator,
		Holder:     req.Owner,
		Subject:    req.Subject,
		Target:     req.Target,
	}
	newRoot := s.keeper.GetRoot(ctx, req.ContractId, req.Target)
	ctx.EventManager().EmitEvent(collection.NewEventAttachFrom(event, newRoot))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	if err := s.keeper.Attach(ctx, req.ContractId, ownerAddr, req.Subject, req.Target); err != nil {
		return nil, err
	}

	return &collection.MsgOperatorAttachResponse{}, nil
}

func (s msgServer) OperatorDetach(c context.Context, req *collection.MsgOperatorDetach) (*collection.MsgOperatorDetachResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	ownerAddr, _ := sdk.AccAddressFromBech32(req.Owner)

	operatorAddr, _ := sdk.AccAddressFromBech32(req.Operator)

	if _, err := s.keeper.GetAuthorization(ctx, req.ContractId, ownerAddr, operatorAddr); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	// legacy
	if err := s.keeper.hasNFT(ctx, req.ContractId, req.Subject); err != nil {
		return nil, err
	}
	oldRoot := s.keeper.GetRoot(ctx, req.ContractId, req.Subject)

	event := collection.EventDetached{
		ContractId: req.ContractId,
		Operator:   req.Operator,
		Holder:     req.Owner,
		Subject:    req.Subject,
	}
	ctx.EventManager().EmitEvent(collection.NewEventDetachFrom(event, oldRoot))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	if err := s.keeper.Detach(ctx, req.ContractId, ownerAddr, req.Subject); err != nil {
		return nil, err
	}

	return &collection.MsgOperatorDetachResponse{}, nil
}

func (s msgServer) AttachFrom(c context.Context, req *collection.MsgAttachFrom) (*collection.MsgAttachFromResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	fromAddr, _ := sdk.AccAddressFromBech32(req.From)

	proxyAddr, _ := sdk.AccAddressFromBech32(req.Proxy)

	if _, err := s.keeper.GetAuthorization(ctx, req.ContractId, fromAddr, proxyAddr); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	event := collection.EventAttached{
		ContractId: req.ContractId,
		Operator:   req.Proxy,
		Holder:     req.From,
		Subject:    req.TokenId,
		Target:     req.ToTokenId,
	}
	newRoot := s.keeper.GetRoot(ctx, req.ContractId, req.ToTokenId)
	ctx.EventManager().EmitEvent(collection.NewEventAttachFrom(event, newRoot))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	if err := s.keeper.Attach(ctx, req.ContractId, fromAddr, req.TokenId, req.ToTokenId); err != nil {
		return nil, err
	}

	return &collection.MsgAttachFromResponse{}, nil
}

func (s msgServer) DetachFrom(c context.Context, req *collection.MsgDetachFrom) (*collection.MsgDetachFromResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	fromAddr, _ := sdk.AccAddressFromBech32(req.From)

	proxyAddr, _ := sdk.AccAddressFromBech32(req.Proxy)

	if _, err := s.keeper.GetAuthorization(ctx, req.ContractId, fromAddr, proxyAddr); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	// legacy
	if err := s.keeper.hasNFT(ctx, req.ContractId, req.TokenId); err != nil {
		return nil, err
	}
	oldRoot := s.keeper.GetRoot(ctx, req.ContractId, req.TokenId)

	event := collection.EventDetached{
		ContractId: req.ContractId,
		Operator:   req.Proxy,
		Holder:     req.From,
		Subject:    req.TokenId,
	}
	ctx.EventManager().EmitEvent(collection.NewEventDetachFrom(event, oldRoot))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	if err := s.keeper.Detach(ctx, req.ContractId, fromAddr, req.TokenId); err != nil {
		return nil, err
	}

	return &collection.MsgDetachFromResponse{}, nil
}
