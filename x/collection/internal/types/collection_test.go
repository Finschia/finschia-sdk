package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetCollection(t *testing.T) {
	collection := NewCollection(defaultSymbol, defaultName, defaultBaseImgURI)

	require.Equal(t, "new_name", collection.SetName("new_name").GetName())
	require.Equal(t, "new_uri", collection.SetBaseImgURI("new_uri").GetBaseImgURI())
}
