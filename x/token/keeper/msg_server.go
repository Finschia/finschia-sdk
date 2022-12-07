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

	from, err := sdk.AccAddressFromBech32(req.From)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", req.From)
	}
	to, err := sdk.AccAddressFromBech32(req.To)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", req.To)
	}
	if err := s.keeper.Send(ctx, req.ContractId, from, to, req.Amount); err != nil {
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

// TransferFrom defines a method to send tokens from one account to another account by the proxy
func (s msgServer) TransferFrom(c context.Context, req *token.MsgTransferFrom) (*token.MsgTransferFromResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	from, err := sdk.AccAddressFromBech32(req.From)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", req.From)
	}
	proxy, err := sdk.AccAddressFromBech32(req.Proxy)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid proxy address: %s", req.Proxy)
	}
	to, err := sdk.AccAddressFromBech32(req.To)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", req.To)
	}

	if _, err := s.keeper.GetAuthorization(ctx, req.ContractId, from, proxy); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	if err := s.keeper.Send(ctx, req.ContractId, from, to, req.Amount); err != nil {
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

// RevokeOperator revokes one to send tokens on behalf of the token holder
func (s msgServer) RevokeOperator(c context.Context, req *token.MsgRevokeOperator) (*token.MsgRevokeOperatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	holder, err := sdk.AccAddressFromBech32(req.Holder)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid holder address: %s", req.Holder)
	}
	operator, err := sdk.AccAddressFromBech32(req.Operator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", req.Operator)
	}
	if err := s.keeper.RevokeOperator(ctx, req.ContractId, holder, operator); err != nil {
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
	approver, err := sdk.AccAddressFromBech32(req.Approver)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid approver address: %s", req.Approver)
	}
	proxy, err := sdk.AccAddressFromBech32(req.Proxy)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid proxy address: %s", req.Proxy)
	}
	if err := s.keeper.AuthorizeOperator(ctx, req.ContractId, approver, proxy); err != nil {
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

	owner, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", req.Owner)
	}
	to, err := sdk.AccAddressFromBech32(req.To)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", req.To)
	}
	s.keeper.Issue(ctx, class, owner, to, req.Amount)

	return &token.MsgIssueResponse{Id: contractID}, nil
}

// GrantPermission allows one to mint or burn tokens or modify a token metadata
func (s msgServer) GrantPermission(c context.Context, req *token.MsgGrantPermission) (*token.MsgGrantPermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	granter, err := sdk.AccAddressFromBech32(req.From)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid granter address: %s", req.From)
	}
	grantee, err := sdk.AccAddressFromBech32(req.To)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid grantee address: %s", req.To)
	}
	permission := token.Permission(token.LegacyPermissionFromString(req.Permission))
	if _, err := s.keeper.GetGrant(ctx, req.ContractId, granter, permission); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}
	if _, err := s.keeper.GetGrant(ctx, req.ContractId, grantee, permission); err == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("%s is already granted for %s", grantee, permission)
	}

	s.keeper.Grant(ctx, req.ContractId, granter, grantee, permission)

	event := token.EventGranted{
		ContractId: req.ContractId,
		Granter:    req.From,
		Grantee:    req.To,
		Permission: permission,
	}
	ctx.EventManager().EmitEvent(token.NewEventGrantPermToken(event))

	return &token.MsgGrantPermissionResponse{}, nil
}

// RevokePermission abandons the permission
func (s msgServer) RevokePermission(c context.Context, req *token.MsgRevokePermission) (*token.MsgRevokePermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	grantee, err := sdk.AccAddressFromBech32(req.From)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid grantee address: %s", req.From)
	}
	permission := token.Permission(token.LegacyPermissionFromString(req.Permission))
	if _, err := s.keeper.GetGrant(ctx, req.ContractId, grantee, permission); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	s.keeper.Abandon(ctx, req.ContractId, grantee, permission)

	event := token.EventRenounced{
		ContractId: req.ContractId,
		Grantee:    req.From,
		Permission: permission,
	}
	ctx.EventManager().EmitEvent(token.NewEventRevokePermToken(event))

	return &token.MsgRevokePermissionResponse{}, nil
}

// Mint defines a method to mint tokens
func (s msgServer) Mint(c context.Context, req *token.MsgMint) (*token.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	from, err := sdk.AccAddressFromBech32(req.From)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", req.From)
	}
	to, err := sdk.AccAddressFromBech32(req.To)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", req.To)
	}
	if err := s.keeper.Mint(ctx, req.ContractId, from, to, req.Amount); err != nil {
		return nil, err
	}

	return &token.MsgMintResponse{}, nil
}

// Burn defines a method to burn tokens
func (s msgServer) Burn(c context.Context, req *token.MsgBurn) (*token.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	from, err := sdk.AccAddressFromBech32(req.From)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", req.From)
	}
	if err := s.keeper.Burn(ctx, req.ContractId, from, req.Amount); err != nil {
		return nil, err
	}

	return &token.MsgBurnResponse{}, nil
}

// BurnFrom defines a method for the proxy to burn tokens on the behalf of the holder.
func (s msgServer) BurnFrom(c context.Context, req *token.MsgBurnFrom) (*token.MsgBurnFromResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	proxy, err := sdk.AccAddressFromBech32(req.Proxy)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid proxy address: %s", req.Proxy)
	}
	from, err := sdk.AccAddressFromBech32(req.From)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", req.From)
	}
	if err := s.keeper.OperatorBurn(ctx, req.ContractId, proxy, from, req.Amount); err != nil {
		return nil, err
	}

	return &token.MsgBurnFromResponse{}, nil
}

// Modify defines a method to modify a token metadata
func (s msgServer) Modify(c context.Context, req *token.MsgModify) (*token.MsgModifyResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	grantee, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid grantee address: %s", req.Owner)
	}
	if _, err := s.keeper.GetGrant(ctx, req.ContractId, grantee, token.PermissionModify); err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	if err := s.keeper.Modify(ctx, req.ContractId, grantee, req.Changes); err != nil {
		return nil, err
	}

	return &token.MsgModifyResponse{}, nil
}
