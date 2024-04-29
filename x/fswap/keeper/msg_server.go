package keeper

import (
	"context"

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
	//c := sdk.UnwrapSDKContext(ctx)
	//swapped, err := s.keeper.GetSwapped(c)
	//if err != nil {
	//	return nil, err
	//}
	//if req.GetAmount().Denom != swapped.OldCoinAmount.Denom {
	//	return nil, sdkerrors.ErrInvalidCoins
	//}
	//from, err := sdk.AccAddressFromBech32(req.FromAddress)
	//if err != nil {
	//	return nil, err
	//}
	//if err := s.keeper.Swap(c, from, req.GetAmount()); err != nil {
	//	return nil, err
	//}
	return &types.MsgSwapResponse{}, nil
}

func (s MsgServer) SwapAll(ctx context.Context, req *types.MsgSwapAllRequest) (*types.MsgSwapAllResponse, error) {
	//c := sdk.UnwrapSDKContext(ctx)
	//from, err := sdk.AccAddressFromBech32(req.FromAddress)
	//if err != nil {
	//	return nil, err
	//}
	//if err := s.keeper.SwapAll(c, from); err != nil {
	//	return nil, err
	//}
	return &types.MsgSwapAllResponse{}, nil
}
