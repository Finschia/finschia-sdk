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

func (s msgServer) StartChallenge(c context.Context, req *types.MsgStartChallenge) (*types.MsgStartChallengeResponse, error) {
	panic("implement me")
}

func (s msgServer) NsectChallenge(c context.Context, req *types.MsgNsectChallenge) (*types.MsgNsectChallengeResponse, error) {
	panic("implement me")
}

func (s msgServer) FinishChallenge(c context.Context, req *types.MsgFinishChallenge) (*types.MsgFinishChallengeResponse, error) {
	panic("implement me")
}
