package keeper

import (
	"context"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/token"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServer returns an implementation of the token MsgServer interface
// for the provided Keeper.
func NewMsgServer(keeper Keeper) token.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

var _ token.MsgServer = msgServer{}

// Send defines a method to send tokens from one account to another account
func (s msgServer) Send(c context.Context, req *token.MsgSend) (*token.MsgSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := s.keeper.Send(ctx, req.ContractId, sdk.AccAddress(req.From), sdk.AccAddress(req.To), req.Amount); err != nil {
		return nil, err
	}

	event := token.EventSent{
		ContractId: req.ContractId,
		Operator:   req.From,
		From:       req.From,
		To:         req.To,
		Amount:     req.Amount,
	}
	ctx.EventManager().EmitEvent(token.NewEventTransfer(event))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &token.MsgSendResponse{}, nil
}

// OperatorSend defines a method to send tokens from one account to another account by the operator
func (s msgServer) OperatorSend(c context.Context, req *token.MsgOperatorSend) (*token.MsgOperatorSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if _, err := s.keeper.GetAuthorization(ctx, req.ContractId, sdk.AccAddress(req.From), sdk.AccAddress(req.Operator)); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	if err := s.keeper.Send(ctx, req.ContractId, sdk.AccAddress(req.From), sdk.AccAddress(req.To), req.Amount); err != nil {
		return nil, err
	}

	event := token.EventSent{
		ContractId: req.ContractId,
		Operator:   req.Operator,
		From:       req.From,
		To:         req.To,
		Amount:     req.Amount,
	}
	ctx.EventManager().EmitEvent(token.NewEventTransferFrom(event))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &token.MsgOperatorSendResponse{}, nil
}

// TransferFrom defines a method to send tokens from one account to another account by the proxy
func (s msgServer) TransferFrom(c context.Context, req *token.MsgTransferFrom) (*token.MsgTransferFromResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if _, err := s.keeper.GetAuthorization(ctx, req.ContractId, sdk.AccAddress(req.From), sdk.AccAddress(req.Proxy)); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	if err := s.keeper.Send(ctx, req.ContractId, sdk.AccAddress(req.From), sdk.AccAddress(req.To), req.Amount); err != nil {
		return nil, err
	}

	event := token.EventSent{
		ContractId: req.ContractId,
		Operator:   req.Proxy,
		From:       req.From,
		To:         req.To,
		Amount:     req.Amount,
	}
	ctx.EventManager().EmitEvent(token.NewEventTransferFrom(event))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &token.MsgTransferFromResponse{}, nil
}

// AuthorizeOperator allows one to send tokens on behalf of the token holder
func (s msgServer) AuthorizeOperator(c context.Context, req *token.MsgAuthorizeOperator) (*token.MsgAuthorizeOperatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.keeper.AuthorizeOperator(ctx, req.ContractId, sdk.AccAddress(req.Holder), sdk.AccAddress(req.Operator)); err != nil {
		return nil, err
	}

	event := token.EventAuthorizedOperator{
		ContractId: req.ContractId,
		Holder:     req.Holder,
		Operator:   req.Operator,
	}
	ctx.EventManager().EmitEvent(token.NewEventApproveToken(event))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &token.MsgAuthorizeOperatorResponse{}, nil
}

// RevokeOperator revokes one to send tokens on behalf of the token holder
func (s msgServer) RevokeOperator(c context.Context, req *token.MsgRevokeOperator) (*token.MsgRevokeOperatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.keeper.RevokeOperator(ctx, req.ContractId, sdk.AccAddress(req.Holder), sdk.AccAddress(req.Operator)); err != nil {
		return nil, err
	}

	if err := ctx.EventManager().EmitTypedEvent(&token.EventRevokedOperator{
		ContractId: req.ContractId,
		Holder:     req.Holder,
		Operator:   req.Operator,
	}); err != nil {
		panic(err)
	}

	return &token.MsgRevokeOperatorResponse{}, nil
}

// Approve allows one to send tokens on behalf of the approver
func (s msgServer) Approve(c context.Context, req *token.MsgApprove) (*token.MsgApproveResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.keeper.AuthorizeOperator(ctx, req.ContractId, sdk.AccAddress(req.Approver), sdk.AccAddress(req.Proxy)); err != nil {
		return nil, err
	}

	event := token.EventAuthorizedOperator{
		ContractId: req.ContractId,
		Holder:     req.Approver,
		Operator:   req.Proxy,
	}
	ctx.EventManager().EmitEvent(token.NewEventApproveToken(event))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &token.MsgApproveResponse{}, nil
}

// Issue defines a method to issue a token
func (s msgServer) Issue(c context.Context, req *token.MsgIssue) (*token.MsgIssueResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	contractID := s.keeper.classKeeper.NewID(ctx)
	class := token.TokenClass{
		ContractId: contractID,
		Name:       req.Name,
		Symbol:     req.Symbol,
		ImageUri:   req.ImageUri,
		Meta:       req.Meta,
		Decimals:   req.Decimals,
		Mintable:   req.Mintable,
	}

	s.keeper.Issue(ctx, class, sdk.AccAddress(req.Owner), sdk.AccAddress(req.To), req.Amount)

	return &token.MsgIssueResponse{Id: contractID}, nil
}

// Grant allows one to mint or burn tokens or modify a token metadata
func (s msgServer) Grant(c context.Context, req *token.MsgGrant) (*token.MsgGrantResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	granter := sdk.AccAddress(req.Granter)
	grantee := sdk.AccAddress(req.Grantee)
	permission := token.Permission(token.Permission_value[req.Permission])
	if _, err := s.keeper.GetGrant(ctx, req.ContractId, granter, permission); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}
	if _, err := s.keeper.GetGrant(ctx, req.ContractId, grantee, permission); err == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("%s is already granted for %s", grantee, permission.String())
	}

	s.keeper.Grant(ctx, req.ContractId, granter, grantee, permission)

	return &token.MsgGrantResponse{}, nil
}

// Abandon abandons the permission
func (s msgServer) Abandon(c context.Context, req *token.MsgAbandon) (*token.MsgAbandonResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	grantee := sdk.AccAddress(req.Grantee)
	permission := token.Permission(token.Permission_value[req.Permission])
	if _, err := s.keeper.GetGrant(ctx, req.ContractId, grantee, permission); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	s.keeper.Abandon(ctx, req.ContractId, grantee, permission)

	return &token.MsgAbandonResponse{}, nil
}

// GrantPermission allows one to mint or burn tokens or modify a token metadata
func (s msgServer) GrantPermission(c context.Context, req *token.MsgGrantPermission) (*token.MsgGrantPermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	granter := sdk.AccAddress(req.From)
	grantee := sdk.AccAddress(req.To)
	permission := token.Permission(token.Permission_value[req.Permission])
	if _, err := s.keeper.GetGrant(ctx, req.ContractId, granter, permission); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}
	if _, err := s.keeper.GetGrant(ctx, req.ContractId, grantee, permission); err == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("%s is already granted for %s", grantee, permission.String())
	}

	s.keeper.Grant(ctx, req.ContractId, granter, grantee, permission)

	return &token.MsgGrantPermissionResponse{}, nil
}

// RevokePermission abandons the permission
func (s msgServer) RevokePermission(c context.Context, req *token.MsgRevokePermission) (*token.MsgRevokePermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	grantee := sdk.AccAddress(req.From)
	permission := token.Permission(token.Permission_value[req.Permission])
	if _, err := s.keeper.GetGrant(ctx, req.ContractId, grantee, permission); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	s.keeper.Abandon(ctx, req.ContractId, grantee, permission)

	return &token.MsgRevokePermissionResponse{}, nil
}

// Mint defines a method to mint tokens
func (s msgServer) Mint(c context.Context, req *token.MsgMint) (*token.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.keeper.Mint(ctx, req.ContractId, sdk.AccAddress(req.From), sdk.AccAddress(req.To), req.Amount); err != nil {
		return nil, err
	}

	return &token.MsgMintResponse{}, nil
}

// Burn defines a method to burn tokens
func (s msgServer) Burn(c context.Context, req *token.MsgBurn) (*token.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.keeper.Burn(ctx, req.ContractId, sdk.AccAddress(req.From), req.Amount); err != nil {
		return nil, err
	}

	return &token.MsgBurnResponse{}, nil
}

// OperatorBurn defines a method for the operator to burn tokens on the behalf of the holder.
func (s msgServer) OperatorBurn(c context.Context, req *token.MsgOperatorBurn) (*token.MsgOperatorBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.keeper.OperatorBurn(ctx, req.ContractId, sdk.AccAddress(req.Operator), sdk.AccAddress(req.From), req.Amount); err != nil {
		return nil, err
	}

	return &token.MsgOperatorBurnResponse{}, nil
}

// BurnFrom defines a method for the proxy to burn tokens on the behalf of the holder.
func (s msgServer) BurnFrom(c context.Context, req *token.MsgBurnFrom) (*token.MsgBurnFromResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.keeper.OperatorBurn(ctx, req.ContractId, sdk.AccAddress(req.Proxy), sdk.AccAddress(req.From), req.Amount); err != nil {
		return nil, err
	}

	return &token.MsgBurnFromResponse{}, nil
}

// Modify defines a method to modify a token metadata
func (s msgServer) Modify(c context.Context, req *token.MsgModify) (*token.MsgModifyResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	grantee := sdk.AccAddress(req.Owner)
	if _, err := s.keeper.GetGrant(ctx, req.ContractId, grantee, token.Permission_Modify); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	if err := s.keeper.Modify(ctx, req.ContractId, grantee, req.Changes); err != nil {
		return nil, err
	}

	return &token.MsgModifyResponse{}, nil
}
