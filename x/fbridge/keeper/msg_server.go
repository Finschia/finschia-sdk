package keeper

import (
	"context"

	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

func NewMsgServer(k Keeper) types.MsgServer {
	return &msgServer{k}
}

func (m msgServer) Transfer(goCtx context.Context, msg *types.MsgTransfer) (*types.MsgTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if m.IsBridgeHalted(ctx) {
		return nil, types.ErrInactiveBridge
	}

	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if err := IsValidEthereumAddress(msg.Receiver); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid receiver address")
	}

	seq, err := m.handleBridgeTransfer(ctx, from, msg.Amount)
	if err != nil {
		return nil, err
	}

	if err := ctx.EventManager().EmitTypedEvent(&types.EventTransfer{
		Sender:   msg.Sender,
		Receiver: msg.Receiver,
		Amount:   msg.Amount.String(),
		Seq:      seq,
	}); err != nil {
		panic(err)
	}

	return &types.MsgTransferResponse{}, nil
}

func (m msgServer) Provision(ctx context.Context, msg *types.MsgProvision) (*types.MsgProvisionResponse, error) {
	panic("implement me")
}

func (m msgServer) HoldTransfer(ctx context.Context, msg *types.MsgHoldTransfer) (*types.MsgHoldTransferResponse, error) {
	panic("implement me")
}

func (m msgServer) ReleaseTransfer(ctx context.Context, msg *types.MsgReleaseTransfer) (*types.MsgReleaseTransferResponse, error) {
	panic("implement me")
}

func (m msgServer) RemoveProvision(ctx context.Context, msg *types.MsgRemoveProvision) (*types.MsgRemoveProvisionResponse, error) {
	panic("implement me")
}

func (m msgServer) ClaimBatch(ctx context.Context, msg *types.MsgClaimBatch) (*types.MsgClaimBatchResponse, error) {
	panic("implement me")
}

func (m msgServer) Claim(ctx context.Context, msg *types.MsgClaim) (*types.MsgClaimResponse, error) {
	panic("implement me")
}

func (m msgServer) SuggestRole(goCtx context.Context, msg *types.MsgSuggestRole) (*types.MsgSuggestRoleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	proposer, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid proposer address (%s)", err)
	}

	target, err := sdk.AccAddressFromBech32(msg.Target)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid target address (%s)", err)
	}

	if err := types.IsValidRole(msg.Role); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	proposal, err := m.RegisterRoleProposal(ctx, proposer, target, msg.Role)
	if err != nil {
		return nil, err
	}

	if err := ctx.EventManager().EmitTypedEvent(&types.EventSuggestRole{
		Proposal: proposal,
	}); err != nil {
		panic(err)
	}

	return &types.MsgSuggestRoleResponse{}, nil
}

func (m msgServer) AddVoteForRole(goCtx context.Context, msg *types.MsgAddVoteForRole) (*types.MsgAddVoteForRoleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	voter, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid proposer address (%s)", err)
	}

	if err := types.IsValidVoteOption(msg.Option); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if err := m.addVote(ctx, msg.ProposalId, voter, msg.Option); err != nil {
		return nil, err
	}

	if err := ctx.EventManager().EmitTypedEvent(&types.EventAddVoteForRole{
		Voter:      msg.From,
		ProposalId: msg.ProposalId,
		Option:     msg.Option,
	}); err != nil {
		panic(err)
	}

	return &types.MsgAddVoteForRoleResponse{}, nil
}

func (m msgServer) Halt(goCtx context.Context, msg *types.MsgHalt) (*types.MsgHaltResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, err := sdk.AccAddressFromBech32(msg.Guardian)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid guardian address (%s)", err)
	}

	if err := m.setBridgeSwitch(ctx, addr, types.StatusInactive); err != nil {
		return nil, err
	}

	return &types.MsgHaltResponse{}, nil
}

func (m msgServer) Resume(goCtx context.Context, msg *types.MsgResume) (*types.MsgResumeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, err := sdk.AccAddressFromBech32(msg.Guardian)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid guardian address (%s)", err)
	}

	if err := m.setBridgeSwitch(ctx, addr, types.StatusActive); err != nil {
		return nil, err
	}

	return &types.MsgResumeResponse{}, nil
}
