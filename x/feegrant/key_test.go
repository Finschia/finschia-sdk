package feegrant_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/feegrant"
)

func TestMarshalAndUnmarshalFeegrantKey(t *testing.T) {
	grantee, err := sdk.AccAddressFromBech32("link1qk93t4j0yyzgqgt6k5qf8deh8fq6smpnyat72w")
	require.NoError(t, err)
	granter, err := sdk.AccAddressFromBech32("link1p9qh4ldfd6n0qehujsal4k7g0e37kel96dchsc")
	require.NoError(t, err)

	key := feegrant.FeeAllowanceKey(granter, grantee)
	require.Len(t, key, len(grantee.Bytes())+len(granter.Bytes())+3)
	require.Equal(t, feegrant.FeeAllowancePrefixByGrantee(grantee), key[:len(grantee.Bytes())+2])

	g1, g2 := feegrant.ParseAddressesFromFeeAllowanceKey(key)
	require.Equal(t, granter, g1)
	require.Equal(t, grantee, g2)
}
