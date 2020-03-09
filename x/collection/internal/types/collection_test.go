package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetCollection(t *testing.T) {
	collection := NewCollection(defaultContractID, defaultName, defaultMeta, defaultBaseImgURI)
	collection.SetName("new_name")
	collection.SetBaseImgURI("new_uri")
	collection.SetMeta("new_meta")
	require.Equal(t, "new_name", collection.GetName())
	require.Equal(t, "new_uri", collection.GetBaseImgURI())
	require.Equal(t, "new_meta", collection.GetMeta())
}
