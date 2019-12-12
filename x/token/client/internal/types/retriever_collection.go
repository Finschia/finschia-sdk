package types

import (
	"fmt"
	"github.com/line/link/x/token/internal/types"
)

type CollectionRetriever struct {
	querier types.NodeQuerier
}

func NewCollectionRetriever(querier types.NodeQuerier) CollectionRetriever {
	return CollectionRetriever{querier: querier}
}

func (ar CollectionRetriever) GetCollection(symbol string) (types.CollectionWithTokens, error) {
	collection, _, err := ar.GetCollectionWithHeight(symbol)
	return collection, err
}

func (ar CollectionRetriever) GetAllCollections() (types.CollectionsWithTokens, error) {
	collections, _, err := ar.GetAllCollectionsWithHeight()
	return collections, err
}

func (ar CollectionRetriever) GetCollectionWithHeight(symbol string) (types.CollectionWithTokens, int64, error) {
	bs, err := types.ModuleCdc.MarshalJSON(types.NewQueryCollectionParams(symbol))
	if err != nil {
		return types.CollectionWithTokens{}, 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCollections), bs)
	if err != nil {
		return types.CollectionWithTokens{}, height, err
	}

	var collection types.CollectionWithTokens
	if err := types.ModuleCdc.UnmarshalJSON(res, &collection); err != nil {
		return types.CollectionWithTokens{}, height, err
	}

	return collection, height, nil
}

func (ar CollectionRetriever) GetAllCollectionsWithHeight() (types.CollectionsWithTokens, int64, error) {
	var collections types.CollectionsWithTokens

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCollections), nil)
	if err != nil {
		return collections, 0, err
	}

	err = types.ModuleCdc.UnmarshalJSON(res, &collections)
	if err != nil {
		return collections, 0, err
	}

	return collections, height, nil
}

func (ar CollectionRetriever) EnsureExists(symbol string) error {
	if _, err := ar.GetCollection(symbol); err != nil {
		return err
	}
	return nil
}
