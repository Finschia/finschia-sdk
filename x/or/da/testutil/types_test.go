package testutil

import (
	"testing"

	sdk "github.com/Finschia/finschia-rdk/types"

	"github.com/stretchr/testify/require"
)

func TestAccAddress(t *testing.T) {
	_ = AccAddress()
	addr := AccAddressString()
	require.NotPanics(t, func() {
		sdk.MustAccAddressFromBech32(addr)
	})
}
