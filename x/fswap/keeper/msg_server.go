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

func (s MsgServer) Swap(ctx context.Context, req *types.MsgSwap) (*types.MsgSwapResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)

	from, err := sdk.AccAddressFromBech32(req.FromAddress)
	if err != nil {
		return nil, err
	}

	if err := s.keeper.IsSendEnabledCoins(c, req.GetFromCoinAmount()); err != nil {
		return &types.MsgSwapResponse{}, err
	}

	if err := s.keeper.Swap(c, from, req.GetFromCoinAmount(), req.GetToDenom()); err != nil {
		return nil, err
	}

	return &types.MsgSwapResponse{}, nil
}

func (s MsgServer) SwapAll(ctx context.Context, req *types.MsgSwapAll) (*types.MsgSwapAllResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)

	from, err := sdk.AccAddressFromBech32(req.FromAddress)
	if err != nil {
		return nil, err
	}

	balance := s.keeper.GetBalance(c, from, req.FromDenom)
	if balance.IsZero() {
		return nil, sdkerrors.ErrInsufficientFunds
	}

	if err := s.keeper.IsSendEnabledCoins(c, balance); err != nil {
		return nil, err
	}

	if err := s.keeper.Swap(c, from, balance, req.GetToDenom()); err != nil {
		return nil, err
	}

	return &types.MsgSwapAllResponse{}, nil
}

func (s MsgServer) SetSwap(ctx context.Context, req *types.MsgSetSwap) (*types.MsgSetSwapResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)

	if err := s.keeper.validateAuthority(req.Authority); err != nil {
		return nil, err
	}

	if err := s.keeper.MakeSwap(c, req.GetSwap(), req.GetToDenomMetadata()); err != nil {
		return nil, err
	}

	return &types.MsgSetSwapResponse{}, nil
}
