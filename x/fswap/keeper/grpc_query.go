package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Finschia/finschia-sdk/store/prefix"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/types/query"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

var _ types.QueryServer = QueryServer{}

type QueryServer struct {
	Keeper
}

func NewQueryServer(keeper Keeper) *QueryServer {
	return &QueryServer{
		keeper,
	}
}

func (s QueryServer) Swapped(ctx context.Context, req *types.QuerySwappedRequest) (*types.QuerySwappedResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	if err := sdk.ValidateDenom(req.GetFromDenom()); err != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap(err.Error())
	}
	if err := sdk.ValidateDenom(req.GetToDenom()); err != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap(err.Error())
	}

	swapped, err := s.Keeper.getSwapped(c, req.GetFromDenom(), req.GetToDenom())
	if err != nil {
		return nil, err
	}

	return &types.QuerySwappedResponse{
		FromCoinAmount: swapped.GetFromCoinAmount(),
		ToCoinAmount:   swapped.GetToCoinAmount(),
	}, nil
}

func (s QueryServer) TotalSwappableToCoinAmount(ctx context.Context, req *types.QueryTotalSwappableToCoinAmountRequest) (*types.QueryTotalSwappableToCoinAmountResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	if err := sdk.ValidateDenom(req.GetFromDenom()); err != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap(err.Error())
	}
	if err := sdk.ValidateDenom(req.GetToDenom()); err != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap(err.Error())
	}

	amount, err := s.Keeper.getSwappableNewCoinAmount(c, req.GetFromDenom(), req.GetToDenom())
	if err != nil {
		return &types.QueryTotalSwappableToCoinAmountResponse{}, err
	}

	return &types.QueryTotalSwappableToCoinAmountResponse{SwappableAmount: amount}, nil
}

func (s QueryServer) Swap(ctx context.Context, req *types.QuerySwapRequest) (*types.QuerySwapResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	c := sdk.UnwrapSDKContext(ctx)
	swap, err := s.Keeper.getSwap(c, req.GetFromDenom(), req.GetToDenom())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &types.QuerySwapResponse{swap}, nil
}

func (s QueryServer) Swaps(ctx context.Context, req *types.QuerySwapsRequest) (*types.QuerySwapsResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)

	swaps := []types.Swap{}
	store := c.KVStore(s.storeKey)
	swapStore := prefix.NewStore(store, swapPrefix)
	pageResponse, err := query.Paginate(swapStore, req.Pagination, func(key, value []byte) error {
		swap := types.Swap{}
		if err := s.Keeper.cdc.Unmarshal(value, &swap); err != nil {
			return err
		}
		swaps = append(swaps, swap)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &types.QuerySwapsResponse{
		Swaps:      swaps,
		Pagination: pageResponse,
	}, nil
}
