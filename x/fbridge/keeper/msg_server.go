package keeper

import (
	"context"

	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

func NewMsgServer(k Keeper) types.MsgServer {
	return &msgServer{k}
}

func (m msgServer) Transfer(ctx context.Context, transfer *types.MsgTransfer) (*types.MsgTransferResponse, error) {
	panic("implement me")
}

func (m msgServer) Provision(ctx context.Context, provision *types.MsgProvision) (*types.MsgProvisionResponse, error) {
	panic("implement me")
}

func (m msgServer) HoldTransfer(ctx context.Context, transfer *types.MsgHoldTransfer) (*types.MsgHoldTransferResponse, error) {
	panic("implement me")
}

func (m msgServer) ReleaseTransfer(ctx context.Context, transfer *types.MsgReleaseTransfer) (*types.MsgReleaseTransferResponse, error) {
	panic("implement me")
}

func (m msgServer) RemoveProvision(ctx context.Context, provision *types.MsgRemoveProvision) (*types.MsgRemoveProvisionResponse, error) {
	panic("implement me")
}

func (m msgServer) ClaimBatch(ctx context.Context, batch *types.MsgClaimBatch) (*types.MsgClaimBatchResponse, error) {
	panic("implement me")
}

func (m msgServer) Claim(ctx context.Context, claim *types.MsgClaim) (*types.MsgClaimResponse, error) {
	panic("implement me")
}

func (m msgServer) UpdateRole(ctx context.Context, role *types.MsgUpdateRole) (*types.MsgUpdateRoleResponse, error) {
	panic("implement me")
}

func (m msgServer) Halt(ctx context.Context, halt *types.MsgHalt) (*types.MsgHaltResponse, error) {
	panic("implement me")
}

func (m msgServer) Resume(ctx context.Context, resume *types.MsgResume) (*types.MsgResumeResponse, error) {
	panic("implement me")
}
