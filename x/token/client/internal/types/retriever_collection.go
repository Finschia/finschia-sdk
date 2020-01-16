package types

import (
	"fmt"
	context "github.com/line/link/client"
	"github.com/line/link/x/token/internal/types"
)

type CollectionRetriever struct {
	querier types.NodeQuerier
}

func NewCollectionRetriever(querier types.NodeQuerier) CollectionRetriever {
	return CollectionRetriever{querier: querier}
}

func (ar CollectionRetriever) GetCollection(ctx context.CLIContext, symbol string) (types.CollectionWithTokens, error) {
	collection, _, err := ar.GetCollectionWithHeight(ctx, symbol)
	return collection, err
}

func (ar CollectionRetriever) GetAllCollections(ctx context.CLIContext) (types.CollectionsWithTokens, error) {
	collections, _, err := ar.GetAllCollectionsWithHeight(ctx)
	return collections, err
}

func (ar CollectionRetriever) GetCollectionWithHeight(ctx context.CLIContext, symbol string) (types.CollectionWithTokens, int64, error) {
	bs, err := ctx.Codec.MarshalJSON(types.NewQueryCollectionParams(symbol))
	if err != nil {
		return types.CollectionWithTokens{}, 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCollections), bs)
	if err != nil {
		return types.CollectionWithTokens{}, height, err
	}

	var collection types.CollectionWithTokens
	if err := ctx.Codec.UnmarshalJSON(res, &collection); err != nil {
		return types.CollectionWithTokens{}, height, err
	}

	return collection, height, nil
}

func (ar CollectionRetriever) GetAllCollectionsWithHeight(ctx context.CLIContext) (types.CollectionsWithTokens, int64, error) {
	var collections types.CollectionsWithTokens

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCollections), nil)
	if err != nil {
		return collections, 0, err
	}

	err = ctx.Codec.UnmarshalJSON(res, &collections)
	if err != nil {
		return collections, 0, err
	}

	return collections, height, nil
}

func (ar CollectionRetriever) EnsureExists(ctx context.CLIContext, symbol string) error {
	if _, err := ar.GetCollection(ctx, symbol); err != nil {
		return err
	}
	return nil
}
