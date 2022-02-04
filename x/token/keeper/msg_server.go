package keeper

import (
	"context"

	// sdk "github.com/line/lbm-sdk/types"
	// sdkerrors "github.com/line/lbm-sdk/types/errors"
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
func (s msgServer) Transfer(context.Context, *token.MsgTransfer) (*token.MsgTransferResponse, error) {
	panic("Not implemented")
}

// TransferFrom defines a method to transfer tokens from one account to another account by the proxy
func (s msgServer) TransferFrom(context.Context, *token.MsgTransferFrom) (*token.MsgTransferFromResponse, error) {
	panic("Not implemented")
}

// Approve allows one to transfer tokens on behalf of the approver
func (s msgServer) Approve(context.Context, *token.MsgApprove) (*token.MsgApproveResponse, error) {
	panic("Not implemented")
}

// Issue defines a method to issue a token
func (s msgServer) Issue(context.Context, *token.MsgIssue) (*token.MsgIssueResponse, error) {
	panic("Not implemented")
}

// Grant allows one to mint or burn tokens or modify a token metadata
func (s msgServer) Grant(context.Context, *token.MsgGrant) (*token.MsgGrantResponse, error) {
	panic("Not implemented")
}

// Revoke revokes the grant
func (s msgServer) Revoke(context.Context, *token.MsgRevoke) (*token.MsgRevokeResponse, error) {
	panic("Not implemented")
}

// Mint defines a method to mint tokens
func (s msgServer) Mint(context.Context, *token.MsgMint) (*token.MsgMintResponse, error) {
	panic("Not implemented")
}

// Burn defines a method to burn tokens
func (s msgServer) Burn(context.Context, *token.MsgBurn) (*token.MsgBurnResponse, error) {
	panic("Not implemented")
}

// BurnFrom defines a method to burn tokens
func (s msgServer) BurnFrom(context.Context, *token.MsgBurnFrom) (*token.MsgBurnFromResponse, error) {
	panic("Not implemented")
}

// Modify defines a method to modify a token metadata
func (s msgServer) Modify(context.Context, *token.MsgModify) (*token.MsgModifyResponse, error) {
	panic("Not implemented")
}
