package querier

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/token/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQueryEncoder(tokenQuerier sdk.Querier) types.EncodeQuerier {
	return func(ctx sdk.Context, jsonQuerier json.RawMessage) ([]byte, error) {
		var customQuerier types.WasmCustomQuerier
		err := json.Unmarshal(jsonQuerier, &customQuerier)
		if err != nil {
			return nil, err
		}
		switch customQuerier.Route {
		case types.QueryTokens:
			return handleQueryToken(ctx, tokenQuerier, []string{customQuerier.Route}, customQuerier.Data)
		case types.QueryBalance:
			return handleQueryBalance(ctx, tokenQuerier, []string{customQuerier.Route}, customQuerier.Data)
		case types.QuerySupply:
			return handleQueryTotal(ctx, tokenQuerier, customQuerier.Data)
		case types.QueryPerms:
			return handleQueryPerms(ctx, tokenQuerier, []string{customQuerier.Route}, customQuerier.Data)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized Msg route: %T", customQuerier.Route)
		}
	}
}

func handleQueryToken(ctx sdk.Context, tokenQuerier sdk.Querier, path []string, msgData json.RawMessage) ([]byte, error) {
	var wrapper types.QueryTokenWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, err
	}
	req := makeRequestQuery(nil)

	contractID := wrapper.QueryTokenParam.ContractID
	if contractID != "" {
		path = append(path, contractID)
	}
	return tokenQuerier(ctx, path, req)
}

func handleQueryBalance(ctx sdk.Context, tokenQuerier sdk.Querier, path []string, msgData json.RawMessage) ([]byte, error) {
	var wrapper types.QueryBalanceWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, err
	}

	req := makeRequestQuery(types.QueryContractIDAccAddressParams{
		Addr: wrapper.QueryBalanceParam.Address,
	})

	contractID := wrapper.QueryBalanceParam.ContractID
	if contractID != "" {
		path = append(path, contractID)
	}
	return tokenQuerier(ctx, path, req)
}

func handleQueryTotal(ctx sdk.Context, tokenQuerier sdk.Querier, msgData json.RawMessage) ([]byte, error) {
	var wrapper types.QueryTotalWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, err
	}
	req := makeRequestQuery(nil)

	path := []string{wrapper.QueryTotalParam.Target}
	contractID := wrapper.QueryTotalParam.ContractID
	if contractID != "" {
		path = append(path, contractID)
	}
	return tokenQuerier(ctx, path, req)
}

func handleQueryPerms(ctx sdk.Context, tokenQuerier sdk.Querier, path []string, msgData json.RawMessage) ([]byte, error) {
	var wrapper types.QueryPermWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, err
	}

	req := makeRequestQuery(types.QueryContractIDAccAddressParams{
		Addr: wrapper.QueryPermParam.Address,
	})

	contractID := wrapper.QueryPermParam.ContractID
	if contractID != "" {
		path = append(path, contractID)
	}
	return tokenQuerier(ctx, path, req)
}

func makeRequestQuery(params interface{}) abci.RequestQuery {
	req := abci.RequestQuery{
		Path: "",
		Data: []byte(string(codec.MustMarshalJSONIndent(types.ModuleCdc, params))),
	}
	return req
}
