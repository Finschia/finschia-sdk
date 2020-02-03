package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

/**
 * This test is for test coverage
 * errors.go is not executed by the tests of this package (tested by keeper.go)
 * So these are needed to raise test coverage
 */
func TestErrors(t *testing.T) {
	require.Error(t, ErrTokenNotNFT(DefaultCodespace, ""))
	require.Error(t, ErrTokenNotCNFT(DefaultCodespace, ""))
}
