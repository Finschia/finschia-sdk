package types

import (
	"fmt"
	"reflect"
	"sync/atomic"
	"unsafe"

	"github.com/line/lbm-sdk/v2/codec"
	"github.com/line/lbm-sdk/v2/store/prefix"
	sdk "github.com/line/lbm-sdk/v2/types"
)

const (
	// StoreKey is the string store key for the param store
	StoreKey = "params"

	// TStoreKey is the string store key for the param transient store
	TStoreKey = "transient_params"
)

// Individual parameter store for each keeper
// Transient store persists for a block, so we use it for
// recording whether the parameter has been changed or not
type Subspace struct {
	cdc         codec.BinaryMarshaler
	legacyAmino *codec.LegacyAmino
	key         sdk.StoreKey // []byte -> []byte, stores parameter
	tkey        sdk.StoreKey // []byte -> bool, stores parameter change
	name        []byte
	table       KeyTable
}

// NewSubspace constructs a store with namestore
func NewSubspace(cdc codec.BinaryMarshaler, legacyAmino *codec.LegacyAmino, key sdk.StoreKey, tkey sdk.StoreKey, name string) *Subspace {
	return &Subspace{
		cdc:         cdc,
		legacyAmino: legacyAmino,
		key:         key,
		tkey:        tkey,
		name:        []byte(name),
		table:       NewKeyTable(),
	}
}

// HasKeyTable returns if the Subspace has a KeyTable registered.
func (s *Subspace) HasKeyTable() bool {
	return len(s.table.m) > 0
}

// WithKeyTable initializes KeyTable and returns modified Subspace
func (s *Subspace) WithKeyTable(table KeyTable) *Subspace {
	if table.m == nil {
		panic("SetKeyTable() called with nil KeyTable")
	}
	if len(s.table.m) != 0 {
		panic("SetKeyTable() called on already initialized Subspace")
	}

	for k, v := range table.m {
		s.table.m[k] = v
	}

	// Allocate additional capacity for Subspace.name
	// So we don't have to allocate extra space each time appending to the key
	name := s.name
	s.name = make([]byte, len(name), len(name)+table.maxKeyLength())
	copy(s.name, name)

	return s
}

// Returns a KVStore identical with ctx.KVStore(s.key).Prefix()
func (s *Subspace) kvStore(ctx sdk.Context) sdk.KVStore {
	// this function can be called concurrently so we should not call append on s.name directly
	name := make([]byte, len(s.name))
	copy(name, s.name)
	return prefix.NewStore(ctx.KVStore(s.key), append(name, '/'))
}

// Returns a transient store for modification
func (s *Subspace) transientStore(ctx sdk.Context) sdk.KVStore {
	// this function can be called concurrently so we should not call append on s.name directly
	name := make([]byte, len(s.name))
	copy(name, s.name)
	return prefix.NewStore(ctx.TransientStore(s.tkey), append(name, '/'))
}

// Validate attempts to validate a parameter value by its key. If the key is not
// registered or if the validation of the value fails, an error is returned.
func (s *Subspace) Validate(ctx sdk.Context, key []byte, value interface{}) error {
	attr, ok := s.table.m[string(key)]
	if !ok {
		return fmt.Errorf("parameter %s not registered", string(key))
	}

	if err := attr.vfn(value); err != nil {
		return fmt.Errorf("invalid parameter value: %s", err)
	}

	return nil
}

// Get queries for a parameter by key from the Subspace's KVStore and sets the
// value to the provided pointer. If the value does not exist, it will panic.
func (s *Subspace) Get(ctx sdk.Context, key []byte, ptr interface{}) {
	s.checkType(key, ptr)

	if s.loadFromCache(key, ptr) {
		// cache hit
		return
	}
	store := s.kvStore(ctx)
	bz := store.Get(key)

	if err := s.legacyAmino.UnmarshalJSON(bz, ptr); err != nil {
		panic(err)
	}
	s.cacheValue(key, ptr)
}

// GetIfExists queries for a parameter by key from the Subspace's KVStore and
// sets the value to the provided pointer. If the value does not exist, it will
// perform a no-op.
func (s *Subspace) GetIfExists(ctx sdk.Context, key []byte, ptr interface{}) {
	if s.loadFromCache(key, ptr) {
		// cache hit
		return
	}
	store := s.kvStore(ctx)
	bz := store.Get(key)
	if bz == nil {
		return
	}

	s.checkType(key, ptr)

	if err := s.legacyAmino.UnmarshalJSON(bz, ptr); err != nil {
		panic(err)
	}
	s.cacheValue(key, ptr)
}

// GetRaw queries for the raw values bytes for a parameter by key.
func (s *Subspace) GetRaw(ctx sdk.Context, key []byte) []byte {
	store := s.kvStore(ctx)
	return store.Get(key)
}

// Has returns if a parameter key exists or not in the Subspace's KVStore.
func (s *Subspace) Has(ctx sdk.Context, key []byte) bool {
	if s.hasCache(key) {
		return true
	}
	store := s.kvStore(ctx)
	return store.Has(key)
}

// Modified returns true if the parameter key is set in the Subspace's transient
// KVStore.
func (s *Subspace) Modified(ctx sdk.Context, key []byte) bool {
	tstore := s.transientStore(ctx)
	return tstore.Has(key)
}

// checkType verifies that the provided key and value are comptable and registered.
func (s *Subspace) checkType(key []byte, value interface{}) {
	attr, ok := s.table.m[string(key)]
	if !ok {
		panic(fmt.Sprintf("parameter %s not registered", string(key)))
	}

	ty := attr.ty
	pty := reflect.TypeOf(value)
	if pty.Kind() == reflect.Ptr {
		pty = pty.Elem()
	}

	if pty != ty {
		panic("type mismatch with registered table")
	}
}

// All the cache-related functions here are thread-safe.
// Currently, since `CheckTx` and `DeliverTx` can run without abci locking,
// these functions must be thread-safe as tx can run concurrently.
// The map data type is not thread-safe by itself, but concurrent access is
// possible with entry fixed. If we access the subspace with an unregistered key,
// it panics, ensuring that the entry of the map is not extended after the server runs.
// Value update and read operations for a single entry of a map can be performed concurrently by
// `atomic.StorePointer` and `atomic.LoadPointer`.
func (s *Subspace) cacheValue(key []byte, value interface{}) {
	attr, ok := s.table.m[string(key)]
	if !ok {
		panic(fmt.Sprintf("parameter %s not registered", string(key)))
	}
	val := reflect.ValueOf(value)
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		val = val.Elem()
	}
	valueToBeCached := reflect.New(val.Type())
	valueToBeCached.Elem().Set(val)
	atomic.StorePointer(&attr.cachedValue, unsafe.Pointer(&valueToBeCached))
}

func (s *Subspace) hasCache(key []byte) bool {
	attr, ok := s.table.m[string(key)]
	if !ok {
		panic(fmt.Sprintf("parameter %s not registered", string(key)))
	}
	cachedValuePtr := (*reflect.Value)(atomic.LoadPointer(&attr.cachedValue))
	return cachedValuePtr != nil
}

func (s *Subspace) loadFromCache(key []byte, value interface{}) bool {
	attr, ok := s.table.m[string(key)]
	if !ok {
		return false
	}
	if reflect.TypeOf(value).Kind() != reflect.Ptr {
		panic("value should be a Pointer")
	}

	cachedValuePtr := (*reflect.Value)(atomic.LoadPointer(&attr.cachedValue))
	if cachedValuePtr == nil {
		return false
	}
	reflect.ValueOf(value).Elem().Set((*cachedValuePtr).Elem())
	return true
}

// Only for test
func (s *Subspace) GetCachedValueForTesting(key []byte, value interface{}) bool {
	return s.loadFromCache(key, value)
}

// Only for test
func (s *Subspace) HasCacheForTesting(key []byte) bool {
	return s.hasCache(key)
}

// Only for test
func (s *Subspace) SetCacheForTesting(key []byte, value interface{}) {
	s.cacheValue(key, value)
}

// Set stores a value for given a parameter key assuming the parameter type has
// been registered. It will panic if the parameter type has not been registered
// or if the value cannot be encoded. A change record is also set in the Subspace's
// transient KVStore to mark the parameter as modified.
func (s *Subspace) Set(ctx sdk.Context, key []byte, value interface{}) {
	s.checkType(key, value)
	store := s.kvStore(ctx)

	bz, err := s.legacyAmino.MarshalJSON(value)
	if err != nil {
		panic(err)
	}

	store.Set(key, bz)

	tstore := s.transientStore(ctx)
	tstore.Set(key, []byte{})

	s.cacheValue(key, value)
}

// Update stores an updated raw value for a given parameter key assuming the
// parameter type has been registered. It will panic if the parameter type has
// not been registered or if the value cannot be encoded. An error is returned
// if the raw value is not compatible with the registered type for the parameter
// key or if the new value is invalid as determined by the registered type's
// validation function.
func (s *Subspace) Update(ctx sdk.Context, key, value []byte) error {
	attr, ok := s.table.m[string(key)]
	if !ok {
		panic(fmt.Sprintf("parameter %s not registered", string(key)))
	}

	ty := attr.ty
	dest := reflect.New(ty).Interface()
	s.GetIfExists(ctx, key, dest)

	if err := s.legacyAmino.UnmarshalJSON(value, dest); err != nil {
		return err
	}

	// destValue contains the dereferenced value of dest so validation function do
	// not have to operate on pointers.
	destValue := reflect.Indirect(reflect.ValueOf(dest)).Interface()
	if err := s.Validate(ctx, key, destValue); err != nil {
		return err
	}

	s.Set(ctx, key, dest)
	return nil
}

// GetParamSet iterates through each ParamSetPair where for each pair, it will
// retrieve the value and set it to the corresponding value pointer provided
// in the ParamSetPair by calling Subspace#Get.
func (s *Subspace) GetParamSet(ctx sdk.Context, ps ParamSet) {
	for _, pair := range ps.ParamSetPairs() {
		s.Get(ctx, pair.Key, pair.Value)
	}
}

// SetParamSet iterates through each ParamSetPair and sets the value with the
// corresponding parameter key in the Subspace's KVStore.
func (s *Subspace) SetParamSet(ctx sdk.Context, ps ParamSet) {
	for _, pair := range ps.ParamSetPairs() {
		// pair.Field is a pointer to the field, so indirecting the ptr.
		// go-amino automatically handles it but just for sure,
		// since SetStruct is meant to be used in InitGenesis
		// so this method will not be called frequently
		v := reflect.Indirect(reflect.ValueOf(pair.Value)).Interface()

		if err := pair.ValidatorFn(v); err != nil {
			panic(fmt.Sprintf("value from ParamSetPair is invalid: %s", err))
		}

		s.Set(ctx, pair.Key, v)
	}
}

// Name returns the name of the Subspace.
func (s *Subspace) Name() string {
	return string(s.name)
}

// Wrapper of Subspace, provides immutable functions only
type ReadOnlySubspace struct {
	s *Subspace
}

// Get delegates a read-only Get call to the Subspace.
func (ros ReadOnlySubspace) Get(ctx sdk.Context, key []byte, ptr interface{}) {
	ros.s.Get(ctx, key, ptr)
}

// GetRaw delegates a read-only GetRaw call to the Subspace.
func (ros ReadOnlySubspace) GetRaw(ctx sdk.Context, key []byte) []byte {
	return ros.s.GetRaw(ctx, key)
}

// Has delegates a read-only Has call to the Subspace.
func (ros ReadOnlySubspace) Has(ctx sdk.Context, key []byte) bool {
	return ros.s.Has(ctx, key)
}

// Modified delegates a read-only Modified call to the Subspace.
func (ros ReadOnlySubspace) Modified(ctx sdk.Context, key []byte) bool {
	return ros.s.Modified(ctx, key)
}

// Name delegates a read-only Name call to the Subspace.
func (ros ReadOnlySubspace) Name() string {
	return ros.s.Name()
}
