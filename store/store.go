package store

import (
	"github.com/Finschia/ostracon/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/Finschia/finschia-rdk/store/cache"
	"github.com/Finschia/finschia-rdk/store/rootmulti"
	"github.com/Finschia/finschia-rdk/store/types"
)

func NewCommitMultiStore(db dbm.DB) types.CommitMultiStore {
	return rootmulti.NewStore(db, log.NewNopLogger())
}

func NewCommitKVStoreCacheManager(cacheSize int, metricsProvider cache.MetricsProvider) types.MultiStorePersistentCache {
	return cache.NewCommitKVStoreCacheManager(cacheSize, metricsProvider)
}
