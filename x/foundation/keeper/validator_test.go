package keeper_test

import (
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/simapp"
	"github.com/line/lbm-sdk/x/foundation"
)

func TestGetSetValidatorAuth(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	k := app.FoundationKeeper

	// not added yet
	_, err := k.GetValidatorAuth(ctx, valAddr)
	require.Error(t, err)

	// test adding creation allowed validators
	expected := &foundation.ValidatorAuth{
		OperatorAddress: valAddr.String(),
		CreationAllowed: true,
	}
	require.NoError(t, k.SetValidatorAuth(ctx, expected))
	actual, err := k.GetValidatorAuth(ctx, valAddr)
	require.Equal(t, expected, actual)

	require.Equal(t, []*foundation.ValidatorAuth{expected}, k.GetValidatorAuths(ctx))
}
