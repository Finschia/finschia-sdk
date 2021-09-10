package store

import (
	tmdb "github.com/line/tm-db/v2"

	"github.com/line/lbm-sdk/store/cache"
	"github.com/line/lbm-sdk/store/rootmulti"
	"github.com/line/lbm-sdk/store/types"
)

func NewCommitMultiStore(db tmdb.DB) types.CommitMultiStore {
	return rootmulti.NewStore(db)
}

func NewCommitKVStoreCacheManager(cacheSize int, metricsProvider cache.MetricsProvider) types.MultiStorePersistentCache {
	return cache.NewCommitKVStoreCacheManager(cacheSize, metricsProvider)
}
