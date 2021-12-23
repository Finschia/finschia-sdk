package keeper_test

import (
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/simapp"
	"github.com/line/lbm-sdk/x/consortium/types"
)

func TestGetSetValidatorAuth(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	k := app.ConsortiumKeeper

	// not added yet
	_, err := k.GetValidatorAuth(ctx, valAddr)
	require.Error(t, err)

	// test adding creation allowed validators
	expected := &types.ValidatorAuth{
		OperatorAddress: string(valAddr),
		CreationAllowed: true,
	}
	require.NoError(t, k.SetValidatorAuth(ctx, expected))
	actual, err := k.GetValidatorAuth(ctx, valAddr)
	require.Equal(t, expected, actual)

	require.Equal(t, []*types.ValidatorAuth{expected}, k.GetValidatorAuths(ctx))
}
