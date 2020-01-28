package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCollections_String(t *testing.T) {
	collections := Collections{NewCollection("1"), NewCollection("2"), NewCollection("3"), NewCollection("4")}
	expected := "{\n  \"symbol\": \"1\"\n},{\n  \"symbol\": \"2\"\n},{\n  \"symbol\": \"3\"\n},{\n  \"symbol\": \"4\"\n}"
	require.Equal(t, expected, collections.String())
}
