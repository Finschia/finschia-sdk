package snapshots_test

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-sdk/snapshots"
)

func TestStreamWriter(t *testing.T) {
	ch := make(chan io.ReadCloser, 1)
	writer := snapshots.NewStreamWriter(ch)
	writer.CloseWithError(errors.New("test error"))
	err := writer.Close()
	require.Error(t, err)
}
