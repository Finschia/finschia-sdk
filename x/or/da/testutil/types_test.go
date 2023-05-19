package testutil

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccAddress(t *testing.T) {
	addr := AccAddress()
	require.NotPanics(t, func() {
		sdk.MustAccAddressFromBech32(addr)
	})
}
