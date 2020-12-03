package querier

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/collection/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/line/link-modules/x/wasm"
)

func NewQueryEncoder(collectionQuerier sdk.Querier) wasm.EncodeQuerier {
	return func(ctx sdk.Context, jsonQuerier json.RawMessage) ([]byte, error) {
		var customQuerier types.WasmCustomQuerier
		err := json.Unmarshal(jsonQuerier, &customQuerier)
		if err != nil {
			return nil, err
		}
		switch customQuerier.Route {
		case types.QueryCollections:
			return handleQueryCollections(ctx, collectionQuerier, []string{customQuerier.Route}, customQuerier.Data)
		case types.QueryBalance:
			return handleQueryBalances(ctx, collectionQuerier, []string{customQuerier.Route}, customQuerier.Data)
		case types.QueryBalances:
			return handleQueryBalances(ctx, collectionQuerier, []string{customQuerier.Route}, customQuerier.Data)
		case types.QueryTokenTypes:
			return handleQueryTokenTypes(ctx, collectionQuerier, []string{customQuerier.Route}, customQuerier.Data)
		case types.QueryTokens:
			return handleQueryTokens(ctx, collectionQuerier, []string{customQuerier.Route}, customQuerier.Data)
		case types.QueryTokensWithTokenType:
			return handleQueryTokensWithTokenType(ctx, collectionQuerier, []string{customQuerier.Route}, customQuerier.Data)
		case types.QueryNFTCount:
			return handleQueryNFTCount(ctx, collectionQuerier, []string{customQuerier.Route}, customQuerier.Data)
		case types.QueryNFTMint:
			return handleQueryNFTCount(ctx, collectionQuerier, []string{customQuerier.Route}, customQuerier.Data)
		case types.QueryNFTBurn:
			return handleQueryNFTCount(ctx, collectionQuerier, []string{customQuerier.Route}, customQuerier.Data)
		case types.QuerySupply:
			return handleQueryTotal(ctx, collectionQuerier, customQuerier.Data)
		case types.QueryParent:
			return handleQueryRootOrParentOrChildren(ctx, collectionQuerier, []string{customQuerier.Route}, customQuerier.Data)
		case types.QueryRoot:
			return handleQueryRootOrParentOrChildren(ctx, collectionQuerier, []string{customQuerier.Route}, customQuerier.Data)
		case types.QueryChildren:
			return handleQueryRootOrParentOrChildren(ctx, collectionQuerier, []string{customQuerier.Route}, customQuerier.Data)
		case types.QueryPerms:
			return handleQueryPerms(ctx, collectionQuerier, []string{customQuerier.Route}, customQuerier.Data)
		case types.QueryApprovers:
			return handleQueryApprovers(ctx, collectionQuerier, []string{customQuerier.Route}, customQuerier.Data)
		case types.QueryIsApproved:
			return handleQueryApproved(ctx, collectionQuerier, []string{customQuerier.Route}, customQuerier.Data)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized Msg route: %T", customQuerier.Route)
		}
	}
}

func handleQueryCollections(ctx sdk.Context, collectionQuerier sdk.Querier, path []string, msgData json.RawMessage) ([]byte, error) {
	var wrapper types.QueryCollectionWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, err
	}
	req := makeRequestQuery(nil)

	contractID := wrapper.QueryCollectionParam.ContractID
	if contractID != "" {
		path = append(path, contractID)
	}
	return collectionQuerier(ctx, path, req)
}

func handleQueryBalances(ctx sdk.Context, collectionQuerier sdk.Querier, path []string, msgData json.RawMessage) ([]byte, error) {
	var wrapper types.QueryBalanceWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, err
	}
	param := types.NewQueryTokenIDAccAddressParams(wrapper.QueryBalanceParam.TokenID, wrapper.QueryBalanceParam.Addr)
	req := makeRequestQuery(param)

	contractID := wrapper.QueryBalanceParam.ContractID
	if contractID != "" {
		path = append(path, contractID)
	}
	return collectionQuerier(ctx, path, req)
}

func handleQueryTokenTypes(ctx sdk.Context, collectionQuerier sdk.Querier, path []string, msgData json.RawMessage) ([]byte, error) {
	var wrapper types.QueryTokenTypesWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, err
	}
	param := types.NewQueryTokenIDParams(wrapper.QueryTokenTypesParam.TokenID)
	req := makeRequestQuery(param)

	contractID := wrapper.QueryTokenTypesParam.ContractID
	if contractID != "" {
		path = append(path, contractID)
	}
	return collectionQuerier(ctx, path, req)
}

func handleQueryTokens(ctx sdk.Context, collectionQuerier sdk.Querier, path []string, msgData json.RawMessage) ([]byte, error) {
	var wrapper types.QueryTokensWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, err
	}
	param := types.NewQueryTokenIDParams(wrapper.QueryTokensParam.TokenID)
	req := makeRequestQuery(param)

	contractID := wrapper.QueryTokensParam.ContractID
	if contractID != "" {
		path = append(path, contractID)
	}
	return collectionQuerier(ctx, path, req)
}

func handleQueryTokensWithTokenType(ctx sdk.Context, collectionQuerier sdk.Querier, path []string, msgData json.RawMessage) ([]byte, error) {
	var wrapper types.QueryTokenTypeWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, err
	}
	param := types.NewQueryTokenTypeParams(wrapper.QueryTokenTypeParam.TokenType)
	req := makeRequestQuery(param)

	contractID := wrapper.QueryTokenTypeParam.ContractID
	if contractID != "" {
		path = append(path, contractID)
	}
	return collectionQuerier(ctx, path, req)
}

func handleQueryNFTCount(ctx sdk.Context, collectionQuerier sdk.Querier, path []string, msgData json.RawMessage) ([]byte, error) {
	var wrapper types.QueryNFTCountWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, err
	}
	param := types.NewQueryTokenIDParams(wrapper.QueryTokensParam.TokenID)
	req := makeRequestQuery(param)

	contractID := wrapper.QueryTokensParam.ContractID
	if contractID != "" {
		path = append(path, contractID)
	}
	return collectionQuerier(ctx, path, req)
}

func handleQueryTotal(ctx sdk.Context, collectionQuerier sdk.Querier, msgData json.RawMessage) ([]byte, error) {
	var wrapper types.QueryTotalWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, err
	}
	params := types.NewQueryTokenIDParams(wrapper.QueryTotalParam.TokenID)
	req := makeRequestQuery(params)

	path := []string{wrapper.QueryTotalParam.Target}
	contractID := wrapper.QueryTotalParam.ContractID
	if contractID != "" {
		path = append(path, contractID)
	}

	return collectionQuerier(ctx, path, req)
}

func handleQueryRootOrParentOrChildren(ctx sdk.Context, collectionQuerier sdk.Querier, path []string, msgData json.RawMessage) ([]byte, error) {
	var wrapper types.QueryTokensWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, err
	}
	param := types.NewQueryTokenIDParams(wrapper.QueryTokensParam.TokenID)
	req := makeRequestQuery(param)

	contractID := wrapper.QueryTokensParam.ContractID
	if contractID != "" {
		path = append(path, contractID)
	}
	return collectionQuerier(ctx, path, req)
}

func handleQueryPerms(ctx sdk.Context, collectionQuerier sdk.Querier, path []string, msgData json.RawMessage) ([]byte, error) {
	var wrapper types.QueryPermsWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, err
	}
	param := types.NewQueryAccAddressParams(wrapper.QueryPermParam.Address)
	req := makeRequestQuery(param)

	contractID := wrapper.QueryPermParam.ContractID
	if contractID != "" {
		path = append(path, contractID)
	}
	return collectionQuerier(ctx, path, req)
}

func handleQueryApproved(ctx sdk.Context, collectionQuerier sdk.Querier, path []string, msgData json.RawMessage) ([]byte, error) {
	var wrapper types.QueryApprovedWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, err
	}
	param := types.NewQueryIsApprovedParams(wrapper.QueryApprovedParam.Proxy, wrapper.QueryApprovedParam.Approver)
	req := makeRequestQuery(param)

	contractID := wrapper.QueryApprovedParam.ContractID
	if contractID != "" {
		path = append(path, contractID)
	}
	return collectionQuerier(ctx, path, req)
}

func handleQueryApprovers(ctx sdk.Context, collectionQuerier sdk.Querier, path []string, msgData json.RawMessage) ([]byte, error) {
	var wrapper types.QueryApproversWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, err
	}
	param := types.NewQueryApproverParams(wrapper.QueryApproversParam.Proxy)
	req := makeRequestQuery(param)

	contractID := wrapper.QueryApproversParam.ContractID
	if contractID != "" {
		path = append(path, contractID)
	}
	return collectionQuerier(ctx, path, req)
}

func makeRequestQuery(params interface{}) abci.RequestQuery {
	req := abci.RequestQuery{
		Path: "",
		Data: []byte(string(codec.MustMarshalJSONIndent(types.ModuleCdc, params))),
	}
	return req
}
