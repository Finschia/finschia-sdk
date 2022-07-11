package keeper

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/line/lbm-sdk/store/prefix"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/query"
	"github.com/line/lbm-sdk/x/staking/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

// Validators queries all validators that match the given status
func (k Querier) Validators(c context.Context, req *types.QueryValidatorsRequest) (*types.QueryValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// validate the provided status, return all the validators if the status is empty
	if req.Status != "" && !(req.Status == types.Bonded.String() || req.Status == types.Unbonded.String() || req.Status == types.Unbonding.String()) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid validator status %s", req.Status)
	}

	var validators types.Validators
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	valStore := prefix.NewStore(store, types.ValidatorsKey)

	pageRes, err := query.FilteredPaginate(valStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		val, err := types.UnmarshalValidator(k.cdc, value)
		if err != nil {
			return false, err
		}

		if req.Status != "" && !strings.EqualFold(val.GetStatus().String(), req.Status) {
			return false, nil
		}

		if accumulate {
			validators = append(validators, val)
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryValidatorsResponse{Validators: validators, Pagination: pageRes}, nil
}

// Validator queries validator info for given validator address
func (k Querier) Validator(c context.Context, req *types.QueryValidatorRequest) (*types.QueryValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.ValidatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}
	if err := sdk.ValidateValAddress(req.ValidatorAddr); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	valAddr := sdk.ValAddress(req.ValidatorAddr)

	ctx := sdk.UnwrapSDKContext(c)
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return nil, status.Errorf(codes.NotFound, "validator %s not found", req.ValidatorAddr)
	}

	return &types.QueryValidatorResponse{Validator: validator}, nil
}

// ValidatorDelegations queries delegate info for given validator
func (k Querier) ValidatorDelegations(c context.Context, req *types.QueryValidatorDelegationsRequest) (*types.QueryValidatorDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.ValidatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}
	if err := sdk.ValidateValAddress(req.ValidatorAddr); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var delegations []types.Delegation
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	valStore := prefix.NewStore(store, types.DelegationKey)
	pageRes, err := query.FilteredPaginate(valStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		delegation, err := types.UnmarshalDelegation(k.cdc, value)
		if err != nil {
			return false, err
		}

		valAddr := sdk.ValAddress(req.ValidatorAddr)

		if !delegation.GetValidatorAddr().Equals(valAddr) {
			return false, nil
		}

		if accumulate {
			delegations = append(delegations, delegation)
		}
		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	delResponses, err := DelegationsToDelegationResponses(ctx, k.Keeper, delegations)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryValidatorDelegationsResponse{
		DelegationResponses: delResponses, Pagination: pageRes}, nil
}

// ValidatorUnbondingDelegations queries unbonding delegations of a validator
func (k Querier) ValidatorUnbondingDelegations(c context.Context, req *types.QueryValidatorUnbondingDelegationsRequest) (*types.QueryValidatorUnbondingDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.ValidatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}
	if err := sdk.ValidateValAddress(req.ValidatorAddr); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	var ubds types.UnbondingDelegations
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)

	valAddr := sdk.ValAddress(req.ValidatorAddr)

	srcValPrefix := types.GetUBDsByValIndexKey(valAddr)
	ubdStore := prefix.NewStore(store, srcValPrefix)
	pageRes, err := query.Paginate(ubdStore, req.Pagination, func(key []byte, value []byte) error {
		storeKey := types.GetUBDKeyFromValIndexKey(append(srcValPrefix, key...))
		storeValue := store.Get(storeKey)

		ubd, err := types.UnmarshalUBD(k.cdc, storeValue)
		if err != nil {
			return err
		}
		ubds = append(ubds, ubd)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryValidatorUnbondingDelegationsResponse{
		UnbondingResponses: ubds,
		Pagination:         pageRes,
	}, nil
}

// Delegation queries delegate info for given validator delegator pair
func (k Querier) Delegation(c context.Context, req *types.QueryDelegationRequest) (*types.QueryDelegationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if err := sdk.ValidateAccAddress(req.DelegatorAddr); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if req.ValidatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}
	if err := sdk.ValidateValAddress(req.ValidatorAddr); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)
	delAddr := sdk.AccAddress(req.DelegatorAddr)

	valAddr := sdk.ValAddress(req.ValidatorAddr)

	delegation, found := k.GetDelegation(ctx, delAddr, valAddr)
	if !found {
		return nil, status.Errorf(
			codes.NotFound,
			"delegation with delegator %s not found for validator %s",
			req.DelegatorAddr, req.ValidatorAddr)
	}

	delResponse, err := DelegationToDelegationResponse(ctx, k.Keeper, delegation)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegationResponse{DelegationResponse: &delResponse}, nil
}

// UnbondingDelegation queries unbonding info for give validator delegator pair
func (k Querier) UnbondingDelegation(c context.Context, req *types.QueryUnbondingDelegationRequest) (*types.QueryUnbondingDelegationResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddr == "" {
		return nil, status.Errorf(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if err := sdk.ValidateAccAddress(req.DelegatorAddr); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if req.ValidatorAddr == "" {
		return nil, status.Errorf(codes.InvalidArgument, "validator address cannot be empty")
	}
	if err := sdk.ValidateValAddress(req.ValidatorAddr); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)

	delAddr := sdk.AccAddress(req.DelegatorAddr)

	valAddr := sdk.ValAddress(req.ValidatorAddr)

	unbond, found := k.GetUnbondingDelegation(ctx, delAddr, valAddr)
	if !found {
		return nil, status.Errorf(
			codes.NotFound,
			"unbonding delegation with delegator %s not found for validator %s",
			req.DelegatorAddr, req.ValidatorAddr)
	}

	return &types.QueryUnbondingDelegationResponse{Unbond: unbond}, nil
}

// DelegatorDelegations queries all delegations of a give delegator address
func (k Querier) DelegatorDelegations(c context.Context, req *types.QueryDelegatorDelegationsRequest) (*types.QueryDelegatorDelegationsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if err := sdk.ValidateAccAddress(req.DelegatorAddr); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var delegations types.Delegations
	ctx := sdk.UnwrapSDKContext(c)

	delAddr := sdk.AccAddress(req.DelegatorAddr)

	store := ctx.KVStore(k.storeKey)
	delStore := prefix.NewStore(store, types.GetDelegationsKey(delAddr))
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		delegation, err := types.UnmarshalDelegation(k.cdc, value)
		if err != nil {
			return err
		}
		delegations = append(delegations, delegation)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	delegationResps, err := DelegationsToDelegationResponses(ctx, k.Keeper, delegations)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorDelegationsResponse{DelegationResponses: delegationResps, Pagination: pageRes}, nil

}

// DelegatorValidator queries validator info for given delegator validator pair
func (k Querier) DelegatorValidator(c context.Context, req *types.QueryDelegatorValidatorRequest) (*types.QueryDelegatorValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if err := sdk.ValidateAccAddress(req.DelegatorAddr); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if req.ValidatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}
	if err := sdk.ValidateValAddress(req.ValidatorAddr); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)
	delAddr := sdk.AccAddress(req.DelegatorAddr)

	valAddr := sdk.ValAddress(req.ValidatorAddr)

	validator, err := k.GetDelegatorValidator(ctx, delAddr, valAddr)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorValidatorResponse{Validator: validator}, nil
}

// DelegatorUnbondingDelegations queries all unbonding delegations of a given delegator address
func (k Querier) DelegatorUnbondingDelegations(c context.Context, req *types.QueryDelegatorUnbondingDelegationsRequest) (*types.QueryDelegatorUnbondingDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if err := sdk.ValidateAccAddress(req.DelegatorAddr); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var unbondingDelegations types.UnbondingDelegations
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	delAddr := sdk.AccAddress(req.DelegatorAddr)

	unbStore := prefix.NewStore(store, types.GetUBDsKey(delAddr))
	pageRes, err := query.Paginate(unbStore, req.Pagination, func(key []byte, value []byte) error {
		unbond, err := types.UnmarshalUBD(k.cdc, value)
		if err != nil {
			return err
		}
		unbondingDelegations = append(unbondingDelegations, unbond)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorUnbondingDelegationsResponse{
		UnbondingResponses: unbondingDelegations, Pagination: pageRes}, nil
}

// HistoricalInfo queries the historical info for given height
func (k Querier) HistoricalInfo(c context.Context, req *types.QueryHistoricalInfoRequest) (*types.QueryHistoricalInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Height < 0 {
		return nil, status.Error(codes.InvalidArgument, "height cannot be negative")
	}
	ctx := sdk.UnwrapSDKContext(c)
	hi, found := k.GetHistoricalInfo(ctx, req.Height)
	if !found {
		return nil, status.Errorf(codes.NotFound, "historical info for height %d not found", req.Height)
	}

	return &types.QueryHistoricalInfoResponse{Hist: &hi}, nil
}

// Redelegations queries redelegations of given address
func (k Querier) Redelegations(c context.Context, req *types.QueryRedelegationsRequest) (*types.QueryRedelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	var redels types.Redelegations
	var pageRes *query.PageResponse
	var err error

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	switch {
	case req.DelegatorAddr != "" && req.SrcValidatorAddr != "" && req.DstValidatorAddr != "":
		if err := sdk.ValidateAccAddress(req.DelegatorAddr); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if err := sdk.ValidateValAddress(req.SrcValidatorAddr); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if err := sdk.ValidateValAddress(req.DstValidatorAddr); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		redels, err = queryRedelegation(ctx, k, req)
	case req.DelegatorAddr == "" && req.SrcValidatorAddr != "" && req.DstValidatorAddr == "":
		if err := sdk.ValidateValAddress(req.SrcValidatorAddr); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		redels, pageRes, err = queryRedelegationsFromSrcValidator(store, k, req)
	default:
		if err := sdk.ValidateAccAddress(req.DelegatorAddr); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		redels, pageRes, err = queryAllRedelegations(store, k, req)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	redelResponses, err := RedelegationsToRedelegationResponses(ctx, k.Keeper, redels)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRedelegationsResponse{RedelegationResponses: redelResponses, Pagination: pageRes}, nil
}

// DelegatorValidators queries all validators info for given delegator address
func (k Querier) DelegatorValidators(c context.Context, req *types.QueryDelegatorValidatorsRequest) (*types.QueryDelegatorValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if err := sdk.ValidateAccAddress(req.DelegatorAddr); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var validators types.Validators
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	delAddr := sdk.AccAddress(req.DelegatorAddr)

	delStore := prefix.NewStore(store, types.GetDelegationsKey(delAddr))
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		delegation, err := types.UnmarshalDelegation(k.cdc, value)
		if err != nil {
			return err
		}

		validator, found := k.GetValidator(ctx, delegation.GetValidatorAddr())
		if !found {
			return types.ErrNoValidatorFound
		}

		validators = append(validators, validator)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorValidatorsResponse{Validators: validators, Pagination: pageRes}, nil
}

// Pool queries the pool info
func (k Querier) Pool(c context.Context, _ *types.QueryPoolRequest) (*types.QueryPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	bondDenom := k.BondDenom(ctx)
	bondedPool := k.GetBondedPool(ctx)
	notBondedPool := k.GetNotBondedPool(ctx)

	pool := types.NewPool(
		k.bankKeeper.GetBalance(ctx, notBondedPool.GetAddress(), bondDenom).Amount,
		k.bankKeeper.GetBalance(ctx, bondedPool.GetAddress(), bondDenom).Amount,
	)

	return &types.QueryPoolResponse{Pool: pool}, nil
}

// Params queries the staking parameters
func (k Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func queryRedelegation(ctx sdk.Context, k Querier, req *types.QueryRedelegationsRequest) (redels types.Redelegations, err error) {

	delAddr := sdk.AccAddress(req.DelegatorAddr)

	srcValAddr := sdk.ValAddress(req.SrcValidatorAddr)

	dstValAddr := sdk.ValAddress(req.DstValidatorAddr)

	redel, found := k.GetRedelegation(ctx, delAddr, srcValAddr, dstValAddr)
	if !found {
		return nil, status.Errorf(
			codes.NotFound,
			"redelegation not found for delegator address %s from validator address %s",
			req.DelegatorAddr, req.SrcValidatorAddr)
	}
	redels = []types.Redelegation{redel}

	return redels, err
}

func queryRedelegationsFromSrcValidator(store sdk.KVStore, k Querier, req *types.QueryRedelegationsRequest) (redels types.Redelegations, res *query.PageResponse, err error) {
	valAddr := sdk.ValAddress(req.SrcValidatorAddr)

	srcValPrefix := types.GetREDsFromValSrcIndexKey(valAddr)
	redStore := prefix.NewStore(store, srcValPrefix)
	res, err = query.Paginate(redStore, req.Pagination, func(key []byte, value []byte) error {
		storeKey := types.GetREDKeyFromValSrcIndexKey(append(srcValPrefix, key...))
		storeValue := store.Get(storeKey)
		red, err := types.UnmarshalRED(k.cdc, storeValue)
		if err != nil {
			return err
		}
		redels = append(redels, red)
		return nil
	})

	return redels, res, err
}

func queryAllRedelegations(store sdk.KVStore, k Querier, req *types.QueryRedelegationsRequest) (redels types.Redelegations, res *query.PageResponse, err error) {
	delAddr := sdk.AccAddress(req.DelegatorAddr)

	redStore := prefix.NewStore(store, types.GetREDsKey(delAddr))
	res, err = query.Paginate(redStore, req.Pagination, func(key []byte, value []byte) error {
		redelegation, err := types.UnmarshalRED(k.cdc, value)
		if err != nil {
			return err
		}
		redels = append(redels, redelegation)
		return nil
	})

	return redels, res, err
}
