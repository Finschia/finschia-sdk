package keeper

import (
	"context"
	"encoding/hex"

	sdk "github.com/Finschia/finschia-sdk/types"
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
	ctx := sdk.UnwrapSDKContext(c)

	// TODO:
	// - get start state from rollup module
	// - validate if rollup name exist
	// - validate if challenger exist in rollup
	// - validate if defender exist in rollup
	// - validate if block height exist in rollup blocks
	startState, _ := hex.DecodeString("E335AC012EC8DFD2B339094F06B32A65BB350BA05ACC259402385F823BE055AD") // start state of miniapp (demo application), `jq -r .pre x/or/demo/out/proof-0.json`
	challenge := types.Challenge{
		L:                   0,
		R:                   req.StepCount,
		AssertedStateHashes: map[uint64][]byte{0: startState[:]},
		DefendedStateHashes: map[uint64][]byte{0: startState[:]},
		Challenger:          req.From,
		Defender:            req.To,
		BlockHeight:         req.BlockHeight,
		RollupName:          req.RollupName,
	}

	if s.Keeper.HasChallenge(ctx, challenge.ID()) {
		return nil, types.ErrChallengeExists
	}
	s.Keeper.SetChallenge(ctx, challenge.ID(), challenge)

	return &types.MsgStartChallengeResponse{
		ChallengeId: challenge.ID(),
	}, nil
}

func (s msgServer) NsectChallenge(c context.Context, req *types.MsgNsectChallenge) (*types.MsgNsectChallengeResponse, error) {
	panic("implement me")
}

func (s msgServer) FinishChallenge(c context.Context, req *types.MsgFinishChallenge) (*types.MsgFinishChallengeResponse, error) {
	panic("implement me")
}
