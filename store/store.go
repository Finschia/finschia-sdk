package store

import (
	dbm "github.com/tendermint/tm-db"

	"github.com/line/lbm-sdk/v2/store/cache"
	"github.com/line/lbm-sdk/v2/store/rootmulti"
	"github.com/line/lbm-sdk/v2/store/types"
)

func NewCommitMultiStore(db dbm.DB) types.CommitMultiStore {
	return rootmulti.NewStore(db)
}

func NewCommitKVStoreCacheManager() types.MultiStorePersistentCache {
	return cache.NewCommitKVStoreCacheManager(cache.DefaultCommitKVStoreCacheSize)
}
