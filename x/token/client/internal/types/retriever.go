package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link-modules/x/token/internal/types"
)

type Retriever struct {
	querier types.NodeQuerier
}

func NewRetriever(querier types.NodeQuerier) Retriever {
	return Retriever{querier: querier}
}

func (r Retriever) query(path, contractID string, data []byte) ([]byte, int64, error) {
	return r.querier.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, path, contractID), data)
}

func (r Retriever) GetAccountPermission(ctx context.CLIContext, contractID string, addr sdk.AccAddress) (types.Permissions, int64, error) {
	var pms types.Permissions
	bs, err := ctx.Codec.MarshalJSON(types.NewQueryContractIDAccAddressParams(addr))
	if err != nil {
		return pms, 0, err
	}

	res, height, err := r.query(types.QueryPerms, contractID, bs)
	if err != nil {
		return pms, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &pms); err != nil {
		return pms, height, err
	}

	return pms, height, nil
}

func (r Retriever) GetAccountBalance(ctx context.CLIContext, contractID string, addr sdk.AccAddress) (sdk.Int, int64, error) {
	var supply sdk.Int
	bs, err := ctx.Codec.MarshalJSON(types.NewQueryContractIDAccAddressParams(addr))
	if err != nil {
		return supply, 0, err
	}

	res, height, err := r.query(types.QueryBalance, contractID, bs)
	if err != nil {
		return supply, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &supply); err != nil {
		return supply, height, err
	}

	return supply, height, nil
}
func (r Retriever) GetTotal(ctx context.CLIContext, contractID string, target string) (sdk.Int, int64, error) {
	var total sdk.Int

	res, height, err := r.query(target, contractID, nil)
	if err != nil {
		return total, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &total); err != nil {
		return total, height, err
	}

	return total, height, nil
}

func (r Retriever) GetToken(ctx context.CLIContext, contractID string) (types.Token, int64, error) {
	var token types.Token

	res, height, err := r.query(types.QueryTokens, contractID, nil)
	if err != nil {
		return token, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &token); err != nil {
		return token, height, err
	}
	return token, height, nil
}

func (r Retriever) IsApproved(ctx context.CLIContext, contractID string, proxy sdk.AccAddress, approver sdk.AccAddress) (approved bool, height int64, err error) {
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQueryIsApprovedParams(proxy, approver))
	if err != nil {
		return false, 0, err
	}

	res, height, err := r.query(types.QueryIsApproved, contractID, bs)
	if err != nil {
		return false, 0, err
	}

	err = ctx.Codec.UnmarshalJSON(res, &approved)
	if err != nil {
		return false, 0, err
	}

	return approved, height, nil
}

func (r Retriever) GetApprovers(ctx context.CLIContext, contractID string, proxy sdk.AccAddress) (accAdds []sdk.AccAddress, height int64, err error) {
	bs, err := ctx.Codec.MarshalJSON(types.NewQueryApproverParams(proxy))
	if err != nil {
		return accAdds, 0, err
	}

	res, height, err := r.query(types.QueryApprovers, contractID, bs)
	if err != nil {
		return accAdds, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &accAdds); err != nil {
		return accAdds, height, err
	}

	return accAdds, height, nil
}
