package keeper

import (
	"context"

	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

var _ types.MsgServer = MsgServer{}

type MsgServer struct {
	keeper Keeper
}

func NewMsgServer(keeper Keeper) *MsgServer {
	return &MsgServer{keeper}
}

func (s MsgServer) Swap(ctx context.Context, req *types.MsgSwapRequest) (*types.MsgSwapResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)

	swapInit, err := s.keeper.getSwapInit(c)
	if err != nil {
		return &types.MsgSwapResponse{}, err
	}

	if req.GetFromCoinAmount().Denom != swapInit.GetFromDenom() {
		return nil, sdkerrors.ErrInvalidCoins
	}

	from, err := sdk.AccAddressFromBech32(req.FromAddress)
	if err != nil {
		return nil, err
	}

	if err := s.keeper.Swap(c, from, req.GetFromCoinAmount()); err != nil {
		return nil, err
	}

	return &types.MsgSwapResponse{}, nil
}

func (s MsgServer) SwapAll(ctx context.Context, req *types.MsgSwapAllRequest) (*types.MsgSwapAllResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)

	if !s.keeper.hasBeenInitialized(c) {
		return &types.MsgSwapAllResponse{}, types.ErrSwapNotInitialized
	}

	from, err := sdk.AccAddressFromBech32(req.FromAddress)
	if err != nil {
		return nil, err
	}

	if err := s.keeper.SwapAll(c, from); err != nil {
		return nil, err
	}

	return &types.MsgSwapAllResponse{}, nil
}
