package querier

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/token/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/line/link-modules/x/wasm"
)

func NewQueryEncoder(tokenQuerier sdk.Querier) wasm.EncodeQuerier {
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
			return handleQueryTotal(ctx, tokenQuerier, []string{customQuerier.Route}, customQuerier.Data)
		case types.QueryMint:
			return handleQueryTotal(ctx, tokenQuerier, []string{customQuerier.Route}, customQuerier.Data)
		case types.QueryBurn:
			return handleQueryTotal(ctx, tokenQuerier, []string{customQuerier.Route}, customQuerier.Data)
		case types.QueryPerms:
			return handleQueryPerms(ctx, tokenQuerier, []string{customQuerier.Route}, customQuerier.Data)
		case types.QueryIsApproved:
			return handleQueryIsApproved(ctx, tokenQuerier, []string{customQuerier.Route}, customQuerier.Data)
		case types.QueryApprovers:
			return handleQueryApprovers(ctx, tokenQuerier, []string{customQuerier.Route}, customQuerier.Data)
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

	contractID := wrapper.TokenParam.ContractID
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
		Addr: wrapper.BalanceParam.Address,
	})

	contractID := wrapper.BalanceParam.ContractID
	if contractID != "" {
		path = append(path, contractID)
	}
	return tokenQuerier(ctx, path, req)
}

func handleQueryTotal(ctx sdk.Context, tokenQuerier sdk.Querier, path []string, msgData json.RawMessage) ([]byte, error) {
	var wrapper types.QueryTotalWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, err
	}
	req := makeRequestQuery(nil)

	contractID := wrapper.TotalParam.ContractID
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
		Addr: wrapper.PermParam.Address,
	})

	contractID := wrapper.PermParam.ContractID
	if contractID != "" {
		path = append(path, contractID)
	}
	return tokenQuerier(ctx, path, req)
}

func handleQueryIsApproved(ctx sdk.Context, tokenQuerier sdk.Querier, path []string, msgData json.RawMessage) ([]byte, error) {
	var wrapper types.QueryIsApprovedWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, err
	}

	req := makeRequestQuery(types.QueryIsApprovedParams{
		Proxy:    wrapper.IsApprovedParam.Proxy,
		Approver: wrapper.IsApprovedParam.Approver,
	})

	contractID := wrapper.IsApprovedParam.ContractID
	if contractID != "" {
		path = append(path, contractID)
	}
	return tokenQuerier(ctx, path, req)
}

func handleQueryApprovers(ctx sdk.Context, tokenQuerier sdk.Querier, path []string, msgData json.RawMessage) ([]byte, error) {
	var wrapper types.QueryApproversWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, err
	}

	req := makeRequestQuery(types.QueryProxyParams{
		Proxy: wrapper.ApproversParam.Proxy,
	})

	contractID := wrapper.ApproversParam.ContractID
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
