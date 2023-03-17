package types

import (
	wasmvm "github.com/line/wasmvm"

	storetypes "github.com/line/lbm-sdk/store/types"
)

var _ wasmvm.KVStore = (*WasmStore)(nil)

// WasmStore is a wrapper struct of `KVStore`
// It translates from cosmos KVStore to wasmvm-defined KVStore.
// The spec of interface `Iterator` is a bit different so we cannot use cosmos KVStore directly.
type WasmStore struct {
	storetypes.KVStore
}

// NewWasmStore creates a instance of WasmStore
func NewWasmStore(kvStore storetypes.KVStore) WasmStore {
	return WasmStore{kvStore}
}
