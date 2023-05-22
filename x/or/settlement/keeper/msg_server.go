package keeper

import (
	"context"

	"github.com/Finschia/finschia-sdk/x/or/settlement/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (s msgServer) InitiateChallenge(ctx context.Context, req *types.MsgInitiateChallenge) (*types.MsgInitiateChallengeResponse, error) {
	panic("implement me")
}

func (s msgServer) ProposeState(ctx context.Context, req *types.MsgProposeState) (*types.MsgProposeStateResponse, error) {
	panic("implement me")
}

func (s msgServer) RespondState(ctx context.Context, req *types.MsgRespondState) (*types.MsgRespondStateResponse, error) {
	panic("implement me")
}

func (s msgServer) ConfirmStateTransition(ctx context.Context, req *types.MsgConfirmStateTransition) (*types.MsgConfirmStateTransitionResponse, error) {
	panic("implement me")
}

func (s msgServer) DenyStateTransition(ctx context.Context, req *types.MsgDenyStateTransition) (*types.MsgDenyStateTransitionResponse, error) {
	panic("implement me")
}

func (s msgServer) AddTrieNode(ctx context.Context, req *types.MsgAddTrieNode) (*types.MsgAddTrieNodeResponse, error) {
	panic("implement me")
}
