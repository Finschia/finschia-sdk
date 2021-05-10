package store

import (
	tmdb "github.com/line/tm-db/v2"

	"github.com/line/lbm-sdk/v2/store/cache"
	"github.com/line/lbm-sdk/v2/store/rootmulti"
	"github.com/line/lbm-sdk/v2/store/types"
)

func NewCommitMultiStore(db tmdb.DB) types.CommitMultiStore {
	return rootmulti.NewStore(db)
}

func NewCommitKVStoreCacheManager(cacheSize int) types.MultiStorePersistentCache {
	return cache.NewCommitKVStoreCacheManager(cacheSize)
}
