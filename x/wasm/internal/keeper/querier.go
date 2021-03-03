package keeper

import (
	"sort"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/line/lbm-sdk/x/wasm/internal/types"
)

const (
	QueryListContractByCode = "list-contracts-by-code"
	QueryGetContract        = "contract-info"
	QueryGetContractState   = "contract-state"
	QueryGetCode            = "code"
	QueryListCode           = "list-code"
	QueryContractHistory    = "contract-history"
)

const (
	QueryMethodContractStateSmart = "smart"
	QueryMethodContractStateAll   = "all"
	QueryMethodContractStateRaw   = "raw"
)

// NewQuerier creates a new querier
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryGetContract:
			return queryContractInfo(ctx, path[1], keeper)
		case QueryListContractByCode:
			return queryContractListByCode(ctx, path[1], keeper)
		case QueryGetContractState:
			if len(path) < 3 {
				return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown data query endpoint")
			}
			return queryContractState(ctx, path[1], path[2], req, keeper)
		case QueryGetCode:
			return queryCode(ctx, path[1], keeper)
		case QueryListCode:
			return queryCodeList(ctx, keeper)
		case QueryContractHistory:
			return queryContractHistory(ctx, path[1], keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown data query endpoint")
		}
	}
}

func queryContractInfo(ctx sdk.Context, bech string, keeper Keeper) ([]byte, error) {
	addr, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	info := keeper.GetContractInfo(ctx, addr)
	if info == nil {
		return []byte("null"), nil
	}
	redact(info)
	infoWithAddress := types.NewContractInfoResponse(*info, addr)
	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, infoWithAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

// redact clears all fields not in the public api
func redact(info *types.ContractInfo) {
	info.Created = nil
}

func queryContractListByCode(ctx sdk.Context, codeIDstr string, keeper Keeper) ([]byte, error) {
	codeID, err := strconv.ParseUint(codeIDstr, 10, 64)
	if err != nil {
		return nil, err
	}

	var contracts []types.ContractInfoResponse
	keeper.IterateContractInfo(ctx, func(addr sdk.AccAddress, info types.ContractInfo) bool {
		if info.CodeID == codeID {
			// and add the address
			infoWithAddress := types.NewContractInfoResponse(info, addr)
			contracts = append(contracts, infoWithAddress)
		}
		return false
	})

	// now we sort them by AbsoluteTxPosition
	sort.Slice(contracts, func(i, j int) bool {
		return contracts[i].LessThan(contracts[j])
	})

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, contracts)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryContractState(ctx sdk.Context, bech, queryMethod string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	contractAddr, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, bech)
	}

	var resultData []types.Model
	switch queryMethod {
	case QueryMethodContractStateAll:
		// this returns a serialized json object (which internally encoded binary fields properly)
		for iter := keeper.GetContractState(ctx, contractAddr); iter.Valid(); iter.Next() {
			resultData = append(resultData, types.Model{
				Key:   iter.Key(),
				Value: iter.Value(),
			})
		}
		if resultData == nil {
			resultData = make([]types.Model, 0)
		}
	case QueryMethodContractStateRaw:
		// this returns the raw data from the state, base64-encoded
		return keeper.QueryRaw(ctx, contractAddr, req.Data), nil
	case QueryMethodContractStateSmart:
		// we enforce a subjective gas limit on all queries to avoid infinite loops
		ctx = ctx.WithGasMeter(sdk.NewGasMeter(keeper.queryGasLimit))
		// this returns raw bytes (must be base64-encoded)
		return keeper.QuerySmart(ctx, contractAddr, req.Data)
	default:
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, queryMethod)
	}
	bz, err := types.ModuleCdc.MarshalJSON(resultData)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryCode(ctx sdk.Context, codeIDstr string, keeper Keeper) ([]byte, error) {
	codeID, err := strconv.ParseUint(codeIDstr, 10, 64)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "invalid codeID: "+err.Error())
	}

	res := keeper.GetCodeInfo(ctx, codeID)
	if res == nil {
		// nil, nil leads to 404 in rest handler
		return nil, nil
	}

	code, err := keeper.GetByteCode(ctx, codeID)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "loading wasm code")
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, types.NewCodeInfoResponse(codeID, *res, code))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryCodeList(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var info []types.CodeInfoResponse
	keeper.IterateCodeInfos(ctx, func(i uint64, res types.CodeInfo) bool {
		info = append(info, types.NewCodeInfoResponse(i, res, nil))
		return false
	})

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, info)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryContractHistory(ctx sdk.Context, bech string, keeper Keeper) ([]byte, error) {
	contractAddr, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	entries := keeper.GetContractHistory(ctx, contractAddr)
	if entries == nil {
		// nil, nil leads to 404 in rest handler
		return nil, nil
	}

	histories := make([]types.ContractHistoryResponse, len(entries))
	for i, entry := range entries {
		histories[i] = types.NewContractHistoryResponse(entry)
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, histories)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}
