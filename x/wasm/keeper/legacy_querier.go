package keeper

import (
	"encoding/json"
	"reflect"
	"strconv"

	abci "github.com/line/ostracon/abci/types"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/wasm/types"
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

// NewLegacyQuerier creates a new querier
func NewLegacyQuerier(keeper types.ViewKeeper, gasLimit sdk.Gas) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		var (
			rsp interface{}
			err error
		)
		switch path[0] {
		case QueryGetContract:
			err := sdk.ValidateAccAddress(path[1])
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
			}
			//nolint:staticcheck
			rsp, err = queryContractInfo(ctx, sdk.AccAddress(path[1]), keeper)
		case QueryListContractByCode:
			codeID, err := strconv.ParseUint(path[1], 10, 64)
			if err != nil {
				return nil, sdkerrors.Wrapf(types.ErrInvalid, "code id: %s", err.Error())
			}
			//nolint:staticcheck
			rsp = queryContractListByCode(ctx, codeID, keeper)
		case QueryGetContractState:
			err := sdk.ValidateAccAddress(path[1])
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
			}
			if len(path) < 3 {
				return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown data query endpoint")
			}
			return queryContractState(ctx, path[1], path[2], req.Data, gasLimit, keeper)
		case QueryGetCode:
			codeID, err := strconv.ParseUint(path[1], 10, 64)
			if err != nil {
				return nil, sdkerrors.Wrapf(types.ErrInvalid, "code id: %s", err.Error())
			}
			//nolint:staticcheck
			rsp, err = queryCode(ctx, codeID, keeper)
		case QueryListCode:
			//nolint:staticcheck
			rsp, err = queryCodeList(ctx, keeper)
		case QueryContractHistory:
			err := sdk.ValidateAccAddress(path[1])
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
			}
			//nolint:staticcheck
			rsp, err = queryContractHistory(ctx, sdk.AccAddress(path[1]), keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown data query endpoint")
		}
		if err != nil {
			return nil, err
		}
		if rsp == nil || reflect.ValueOf(rsp).IsNil() {
			return nil, nil
		}
		bz, err := json.MarshalIndent(rsp, "", "  ")
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return bz, nil
	}
}

func queryContractState(ctx sdk.Context, bech, queryMethod string, data []byte, gasLimit sdk.Gas, keeper types.ViewKeeper) (json.RawMessage, error) {
	err := sdk.ValidateAccAddress(bech)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, bech)
	}
	contractAddr := sdk.AccAddress(bech)
	switch queryMethod {
	case QueryMethodContractStateAll:
		resultData := make([]types.Model, 0)
		// this returns a serialized json object (which internally encoded binary fields properly)
		for iter := keeper.GetContractState(ctx, contractAddr); iter.Valid(); iter.Next() {
			resultData = append(resultData, types.Model{
				Key:   iter.Key(),
				Value: iter.Value(),
			})
		}
		bz, err := json.Marshal(resultData)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return bz, nil
	case QueryMethodContractStateRaw:
		// this returns the raw data from the state, base64-encoded
		return keeper.QueryRaw(ctx, contractAddr, data), nil
	case QueryMethodContractStateSmart:
		// we enforce a subjective gas limit on all queries to avoid infinite loops
		ctx = ctx.WithGasMeter(sdk.NewGasMeter(gasLimit))
		msg := types.RawContractMessage(data)
		if err := msg.ValidateBasic(); err != nil {
			return nil, sdkerrors.Wrap(err, "json msg")
		}
		// this returns raw bytes (must be base64-encoded)
		return keeper.QuerySmart(ctx, contractAddr, msg)
	default:
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, queryMethod)
	}
}

//nolint:unparam
func queryCodeList(ctx sdk.Context, keeper types.ViewKeeper) ([]types.CodeInfoResponse, error) {
	var info []types.CodeInfoResponse
	keeper.IterateCodeInfos(ctx, func(i uint64, res types.CodeInfo) bool {
		info = append(info, types.CodeInfoResponse{
			CodeID:                i,
			Creator:               res.Creator,
			DataHash:              res.CodeHash,
			InstantiatePermission: res.InstantiateConfig,
		})
		return false
	})
	return info, nil
}

//nolint:unparam
func queryContractHistory(ctx sdk.Context, contractAddr sdk.AccAddress, keeper types.ViewKeeper) ([]types.ContractCodeHistoryEntry, error) {
	history := keeper.GetContractHistory(ctx, contractAddr)
	// redact response
	for i := range history {
		history[i].Updated = nil
	}
	return history, nil
}

//nolint:unparam
func queryContractListByCode(ctx sdk.Context, codeID uint64, keeper types.ViewKeeper) []string {
	var contracts []string
	keeper.IterateContractsByCode(ctx, codeID, func(addr sdk.AccAddress) bool {
		contracts = append(contracts, addr.String())
		return false
	})
	return contracts
}
