package mem_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/v2/store/mem"
	"github.com/line/lbm-sdk/v2/store/types"
)

func TestStore(t *testing.T) {
	db := mem.NewStore()
	key, value := []byte("key"), []byte("value")

	require.Equal(t, types.StoreTypeMemory, db.GetStoreType())

	require.Nil(t, db.Get(key, types.GetBytesUnmarshalFunc()))
	db.Set(key, value, types.GetBytesMarshalFunc())
	require.Equal(t, value, db.Get(key, types.GetBytesUnmarshalFunc()))

	newValue := []byte("newValue")
	db.Set(key, newValue, types.GetBytesMarshalFunc())
	require.Equal(t, newValue, db.Get(key, types.GetBytesUnmarshalFunc()))

	db.Delete(key)
	require.Nil(t, db.Get(key, types.GetBytesUnmarshalFunc()))
}

func TestCommit(t *testing.T) {
	db := mem.NewStore()
	key, value := []byte("key"), []byte("value")

	db.Set(key, value, types.GetBytesMarshalFunc())
	id := db.Commit()
	require.True(t, id.IsZero())
	require.True(t, db.LastCommitID().IsZero())
	require.Equal(t, value, db.Get(key, types.GetBytesUnmarshalFunc()))
}
