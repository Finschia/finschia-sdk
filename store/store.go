package store

import (
	dbm "github.com/tendermint/tm-db"

	"github.com/line/lbm-sdk/store/cache"
	"github.com/line/lbm-sdk/store/rootmulti"
	"github.com/line/lbm-sdk/store/types"
)

func NewCommitMultiStore(db dbm.DB) types.CommitMultiStore {
	return rootmulti.NewStore(db)
}

func NewCommitKVStoreCacheManager(cacheSize int, metricsProvider cache.MetricsProvider) types.MultiStorePersistentCache {
	return cache.NewCommitKVStoreCacheManager(cacheSize, metricsProvider)
}
