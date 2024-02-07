package testutil

import (
	"testing"

	sdk "github.com/Finschia/finschia-sdk/types"

	"github.com/stretchr/testify/require"
)

func TestAccAddress(t *testing.T) {
	addr := AccAddress()
	require.NotPanics(t, func() {
		sdk.MustAccAddressFromBech32(addr)
	})
}
