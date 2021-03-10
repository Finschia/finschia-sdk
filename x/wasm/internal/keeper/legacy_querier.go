package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/lbm-sdk/x/wasm/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
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
func NewLegacyQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryContractHistory:
			return queryContractHistory(ctx, path[1], keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown data query endpoint")
		}
	}
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
