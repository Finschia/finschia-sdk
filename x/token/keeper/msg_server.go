package keeper

import (
	"context"

	sdk "github.com/Finschia/finschia-rdk/types"
	"github.com/Finschia/finschia-rdk/x/token"
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

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	from := sdk.MustAccAddressFromBech32(req.From)
	to := sdk.MustAccAddressFromBech32(req.To)

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
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &token.MsgSendResponse{}, nil
}

// OperatorSend defines a method to send tokens from one account to another account by the operator
func (s msgServer) OperatorSend(c context.Context, req *token.MsgOperatorSend) (*token.MsgOperatorSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	from := sdk.MustAccAddressFromBech32(req.From)
	operator := sdk.MustAccAddressFromBech32(req.Operator)
	to := sdk.MustAccAddressFromBech32(req.To)

	if _, err := s.keeper.GetAuthorization(ctx, req.ContractId, from, operator); err != nil {
		return nil, token.ErrTokenNotApproved.Wrap(err.Error())
	}

	if err := s.keeper.Send(ctx, req.ContractId, from, to, req.Amount); err != nil {
		return nil, err
	}

	event := token.EventSent{
		ContractId: req.ContractId,
		Operator:   req.Operator,
		From:       req.From,
		To:         req.To,
		Amount:     req.Amount,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &token.MsgOperatorSendResponse{}, nil
}

// RevokeOperator revokes one to send tokens on behalf of the token holder
func (s msgServer) RevokeOperator(c context.Context, req *token.MsgRevokeOperator) (*token.MsgRevokeOperatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	holder := sdk.MustAccAddressFromBech32(req.Holder)
	operator := sdk.MustAccAddressFromBech32(req.Operator)

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

// AuthorizeOperator allows one to send tokens on behalf of the holder
func (s msgServer) AuthorizeOperator(c context.Context, req *token.MsgAuthorizeOperator) (*token.MsgAuthorizeOperatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	holder := sdk.MustAccAddressFromBech32(req.Holder)
	operator := sdk.MustAccAddressFromBech32(req.Operator)

	if err := s.keeper.AuthorizeOperator(ctx, req.ContractId, holder, operator); err != nil {
		return nil, err
	}

	event := token.EventAuthorizedOperator{
		ContractId: req.ContractId,
		Holder:     req.Holder,
		Operator:   req.Operator,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &token.MsgAuthorizeOperatorResponse{}, nil
}

// Issue defines a method to issue a token
func (s msgServer) Issue(c context.Context, req *token.MsgIssue) (*token.MsgIssueResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	class := token.Contract{
		Name:     req.Name,
		Symbol:   req.Symbol,
		Uri:      req.Uri,
		Meta:     req.Meta,
		Decimals: req.Decimals,
		Mintable: req.Mintable,
	}

	owner := sdk.MustAccAddressFromBech32(req.Owner)
	to := sdk.MustAccAddressFromBech32(req.To)
	contractID := s.keeper.Issue(ctx, class, owner, to, req.Amount)

	return &token.MsgIssueResponse{ContractId: contractID}, nil
}

// GrantPermission allows one to mint or burn tokens or modify a token metadata
func (s msgServer) GrantPermission(c context.Context, req *token.MsgGrantPermission) (*token.MsgGrantPermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	granter := sdk.MustAccAddressFromBech32(req.From)
	grantee := sdk.MustAccAddressFromBech32(req.To)
	permission := token.Permission(token.LegacyPermissionFromString(req.Permission))
	if _, err := s.keeper.GetGrant(ctx, req.ContractId, granter, permission); err != nil {
		return nil, token.ErrTokenNoPermission.Wrap(err.Error())
	}

	s.keeper.Grant(ctx, req.ContractId, granter, grantee, permission)

	return &token.MsgGrantPermissionResponse{}, nil
}

// RevokePermission abandons the permission
func (s msgServer) RevokePermission(c context.Context, req *token.MsgRevokePermission) (*token.MsgRevokePermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	grantee := sdk.MustAccAddressFromBech32(req.From)
	permission := token.Permission(token.LegacyPermissionFromString(req.Permission))
	if _, err := s.keeper.GetGrant(ctx, req.ContractId, grantee, permission); err != nil {
		return nil, token.ErrTokenNoPermission.Wrap(err.Error())
	}

	s.keeper.Abandon(ctx, req.ContractId, grantee, permission)

	return &token.MsgRevokePermissionResponse{}, nil
}

// Mint defines a method to mint tokens
func (s msgServer) Mint(c context.Context, req *token.MsgMint) (*token.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	from := sdk.MustAccAddressFromBech32(req.From)
	to := sdk.MustAccAddressFromBech32(req.To)

	if err := s.keeper.Mint(ctx, req.ContractId, from, to, req.Amount); err != nil {
		return nil, err
	}

	return &token.MsgMintResponse{}, nil
}

// Burn defines a method to burn tokens
func (s msgServer) Burn(c context.Context, req *token.MsgBurn) (*token.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	from := sdk.MustAccAddressFromBech32(req.From)

	if err := s.keeper.Burn(ctx, req.ContractId, from, req.Amount); err != nil {
		return nil, err
	}

	return &token.MsgBurnResponse{}, nil
}

// OperatorBurn defines a method for the operator to burn tokens on the behalf of the holder.
func (s msgServer) OperatorBurn(c context.Context, req *token.MsgOperatorBurn) (*token.MsgOperatorBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	operator := sdk.MustAccAddressFromBech32(req.Operator)
	from := sdk.MustAccAddressFromBech32(req.From)

	if err := s.keeper.OperatorBurn(ctx, req.ContractId, operator, from, req.Amount); err != nil {
		return nil, err
	}

	return &token.MsgOperatorBurnResponse{}, nil
}

// Modify defines a method to modify a token metadata
func (s msgServer) Modify(c context.Context, req *token.MsgModify) (*token.MsgModifyResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	token.UpdateMsgModify(req)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	grantee := sdk.MustAccAddressFromBech32(req.Owner)

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, grantee, token.PermissionModify); err != nil {
		return nil, token.ErrTokenNoPermission.Wrap(err.Error())
	}

	if err := s.keeper.Modify(ctx, req.ContractId, grantee, req.Changes); err != nil {
		return nil, err
	}

	return &token.MsgModifyResponse{}, nil
}
