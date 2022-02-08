package keeper

import (
	"context"

	sdk "github.com/line/lbm-sdk/types"
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

// Transfer defines a method to transfer tokens from one account to another account
func (s msgServer) Transfer(c context.Context, req *token.MsgTransfer) (*token.MsgTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	amount := token.FT{ClassId: req.ClassId, Amount: req.Amount}
	if err := s.keeper.transfer(ctx, sdk.AccAddress(req.From), sdk.AccAddress(req.To), []token.FT{amount}); err != nil {
		return nil, err
	}

	return &token.MsgTransferResponse{}, nil
}

// TransferFrom defines a method to transfer tokens from one account to another account by the proxy
func (s msgServer) TransferFrom(c context.Context, req *token.MsgTransferFrom) (*token.MsgTransferFromResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	amount := token.FT{ClassId: req.ClassId, Amount: req.Amount}
	if err := s.keeper.transferFrom(ctx, sdk.AccAddress(req.Proxy), sdk.AccAddress(req.From), sdk.AccAddress(req.To), []token.FT{amount}); err != nil {
		return nil, err
	}

	return &token.MsgTransferFromResponse{}, nil
}

// Approve allows one to transfer tokens on behalf of the approver
func (s msgServer) Approve(c context.Context, req *token.MsgApprove) (*token.MsgApproveResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.keeper.approve(ctx, sdk.AccAddress(req.Approver), sdk.AccAddress(req.Proxy), req.ClassId); err != nil {
		return nil, err
	}

	return &token.MsgApproveResponse{}, nil
}

// Issue defines a method to issue a token
func (s msgServer) Issue(c context.Context, req *token.MsgIssue) (*token.MsgIssueResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	classId := s.keeper.classKeeper.NewId(ctx)
	class := token.Token{
		Id: classId,
		Name: req.Name,
		Symbol: req.Symbol,
		ImageUri: req.ImageUri,
		Meta: req.Meta,
		Decimals: req.Decimals,
		Mintable: req.Mintable,
	}
	if err := s.keeper.issue(ctx, class, sdk.AccAddress(req.Owner), sdk.AccAddress(req.To), req.Amount); err != nil {
		return nil, err
	}

	return &token.MsgIssueResponse{}, nil
}

// Grant allows one to mint or burn tokens or modify a token metadata
func (s msgServer) Grant(c context.Context, req *token.MsgGrant) (*token.MsgGrantResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.keeper.grant(ctx, sdk.AccAddress(req.Granter), sdk.AccAddress(req.Grantee), req.ClassId, req.Action); err != nil {
		return nil, err
	}

	return &token.MsgGrantResponse{}, nil
}

// Revoke revokes the grant
func (s msgServer) Revoke(c context.Context, req *token.MsgRevoke) (*token.MsgRevokeResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.keeper.revoke(ctx, sdk.AccAddress(req.Grantee), req.ClassId, req.Action); err != nil {
		return nil, err
	}

	return &token.MsgRevokeResponse{}, nil
}

// Mint defines a method to mint tokens
func (s msgServer) Mint(c context.Context, req *token.MsgMint) (*token.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	amount := token.FT{ClassId: req.ClassId, Amount: req.Amount}
	if err := s.keeper.mint(ctx, sdk.AccAddress(req.Grantee), sdk.AccAddress(req.To), []token.FT{amount}); err != nil {
		return nil, err
	}

	return &token.MsgMintResponse{}, nil
}

// Burn defines a method to burn tokens
func (s msgServer) Burn(c context.Context, req *token.MsgBurn) (*token.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	amount := token.FT{ClassId: req.ClassId, Amount: req.Amount}
	if err := s.keeper.burn(ctx, sdk.AccAddress(req.From), []token.FT{amount}); err != nil {
		return nil, err
	}

	return &token.MsgBurnResponse{}, nil
}

// BurnFrom defines a method to burn tokens
func (s msgServer) BurnFrom(c context.Context, req *token.MsgBurnFrom) (*token.MsgBurnFromResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	amount := token.FT{ClassId: req.ClassId, Amount: req.Amount}
	if err := s.keeper.burnFrom(ctx, sdk.AccAddress(req.Grantee), sdk.AccAddress(req.From), []token.FT{amount}); err != nil {
		return nil, err
	}

	return &token.MsgBurnFromResponse{}, nil
}

// Modify defines a method to modify a token metadata
func (s msgServer) Modify(c context.Context, req *token.MsgModify) (*token.MsgModifyResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.keeper.modify(ctx, req.ClassId, sdk.AccAddress(req.Grantee), req.Changes); err != nil {
		return nil, err
	}

	return &token.MsgModifyResponse{}, nil
}
