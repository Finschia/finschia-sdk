package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/lbm-sdk/v2/x/collection/internal/types"
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

func (r Retriever) GetAccountBalance(ctx context.CLIContext, contractID, tokenID string, addr sdk.AccAddress) (sdk.Int, int64, error) {
	var balance sdk.Int
	bs, err := ctx.Codec.MarshalJSON(types.NewQueryTokenIDAccAddressParams(tokenID, addr))
	if err != nil {
		return balance, 0, err
	}

	res, height, err := r.query(types.QueryBalance, contractID, bs)
	if err != nil {
		return balance, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &balance); err != nil {
		return balance, height, err
	}

	return balance, height, nil
}

func (r Retriever) GetAccountBalances(ctx context.CLIContext, contractID string, addr sdk.AccAddress) (types.Coins, int64, error) {
	var coins types.Coins
	bs, err := ctx.Codec.MarshalJSON(types.NewQueryAccAddressParams(addr))
	if err != nil {
		return coins, 0, err
	}

	res, height, err := r.query(types.QueryBalances, contractID, bs)

	if err != nil {
		return coins, height, err
	}
	if err := ctx.Codec.UnmarshalJSON(res, &coins); err != nil {
		return coins, height, err
	}
	return coins, height, nil
}

func (r Retriever) GetAccountPermission(ctx context.CLIContext, contractID string, addr sdk.AccAddress) (types.Permissions, int64, error) {
	var pms types.Permissions
	bs, err := ctx.Codec.MarshalJSON(types.NewQueryAccAddressParams(addr))
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

func (r Retriever) GetCollection(ctx context.CLIContext, contractID string) (types.BaseCollection, int64, error) {
	var collection types.BaseCollection
	res, height, err := r.query(types.QueryCollections, contractID, nil)
	if err != nil {
		return collection, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &collection); err != nil {
		return collection, height, err
	}

	return collection, height, nil
}

func (r Retriever) GetCollectionNFTCount(ctx context.CLIContext, contractID, tokenID, target string) (sdk.Int, int64, error) {
	var nftcount sdk.Int
	bs, err := ctx.Codec.MarshalJSON(types.NewQueryTokenIDParams(tokenID))
	if err != nil {
		return nftcount, 0, err
	}
	if target != types.QueryNFTCount && target != types.QueryNFTMint && target != types.QueryNFTBurn {
		return nftcount, 0, fmt.Errorf("invalid target : %s", target)
	}

	res, height, err := r.query(target, contractID, bs)
	if err != nil {
		return nftcount, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &nftcount); err != nil {
		return nftcount, height, err
	}

	return nftcount, height, nil
}

func (r Retriever) GetTotal(ctx context.CLIContext, contractID, tokenID, target string) (sdk.Int, int64, error) {
	var supply sdk.Int
	bs, err := ctx.Codec.MarshalJSON(types.NewQueryTokenIDParams(tokenID))
	if err != nil {
		return supply, 0, err
	}

	res, height, err := r.query(target, contractID, bs)
	if err != nil {
		return supply, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &supply); err != nil {
		return supply, height, err
	}

	return supply, height, nil
}

func (r Retriever) GetToken(ctx context.CLIContext, contractID, tokenID string) (types.Token, int64, error) {
	var token types.Token
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQueryTokenIDParams(tokenID))
	if err != nil {
		return token, 0, err
	}

	res, height, err := r.query(types.QueryTokens, contractID, bs)
	if err != nil {
		return token, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &token); err != nil {
		return token, height, err
	}
	return token, height, nil
}

func (r Retriever) GetTokens(ctx context.CLIContext, contractID string) (types.Tokens, int64, error) {
	var tokens types.Tokens
	res, height, err := r.query(types.QueryTokens, contractID, nil)
	if err != nil {
		return tokens, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &tokens); err != nil {
		return tokens, height, err
	}
	return tokens, height, nil
}

func (r Retriever) GetTokensWithTokenType(ctx context.CLIContext, contractID string, tokenType string) (types.Tokens, int64, error) {
	var tokens types.Tokens
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQueryTokenTypeParams(tokenType))

	if err != nil {
		return tokens, 0, err
	}
	if err = types.ValidateTokenType(tokenType); err != nil {
		return tokens, 0, err
	}
	res, height, err := r.query(types.QueryTokensWithTokenType, contractID, bs)
	if err != nil {
		return tokens, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &tokens); err != nil {
		return tokens, height, err
	}
	return tokens, height, nil
}

func (r Retriever) GetTokenType(ctx context.CLIContext, contractID, tokenTypeID string) (types.TokenType, int64, error) {
	var tokenType types.TokenType
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQueryTokenIDParams(tokenTypeID))
	if err != nil {
		return tokenType, 0, err
	}

	res, height, err := r.query(types.QueryTokenTypes, contractID, bs)
	if err != nil {
		return tokenType, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &tokenType); err != nil {
		return tokenType, height, err
	}
	return tokenType, height, nil
}

func (r Retriever) GetTokenTypes(ctx context.CLIContext, contractID string) (types.TokenTypes, int64, error) {
	var tokenTypes types.TokenTypes

	res, height, err := r.query(types.QueryTokenTypes, contractID, nil)
	if err != nil {
		return tokenTypes, height, err
	}

	if err := ctx.Codec.UnmarshalJSON(res, &tokenTypes); err != nil {
		return tokenTypes, height, err
	}
	return tokenTypes, height, nil
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

func (r Retriever) EnsureExists(ctx context.CLIContext, contractID, tokenID string) error {
	if _, _, err := r.GetToken(ctx, contractID, tokenID); err != nil {
		return err
	}
	return nil
}

func (r Retriever) GetParent(ctx context.CLIContext, contractID, tokenID string) (types.Token, int64, error) {
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQueryTokenIDParams(tokenID))
	if err != nil {
		return nil, 0, err
	}

	res, height, err := r.query(types.QueryParent, contractID, bs)
	if res == nil {
		return nil, 0, err
	}

	var token types.Token
	if err := ctx.Codec.UnmarshalJSON(res, &token); err != nil {
		return nil, 0, err
	}

	return token, height, nil
}

func (r Retriever) GetRoot(ctx context.CLIContext, contractID, tokenID string) (types.Token, int64, error) {
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQueryTokenIDParams(tokenID))
	if err != nil {
		return nil, 0, err
	}

	res, height, err := r.query(types.QueryRoot, contractID, bs)
	if res == nil {
		return nil, 0, err
	}

	var token types.Token
	if err := ctx.Codec.UnmarshalJSON(res, &token); err != nil {
		return nil, 0, err
	}

	return token, height, nil
}

func (r Retriever) GetChildren(ctx context.CLIContext, contractID, tokenID string) (types.Tokens, int64, error) {
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQueryTokenIDParams(tokenID))
	if err != nil {
		return nil, 0, err
	}

	res, height, err := r.query(types.QueryChildren, contractID, bs)
	if res == nil {
		return nil, 0, err
	}

	var tokens types.Tokens
	if err := ctx.Codec.UnmarshalJSON(res, &tokens); err != nil {
		return nil, 0, err
	}

	return tokens, height, nil
}
