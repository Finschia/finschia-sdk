package keeper_test

import (
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
	"github.com/line/lbm-sdk/x/foundation/keeper"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

func newParams(enabled bool) *foundation.Params {
	return &foundation.Params{Enabled: enabled}
}

func newUpdateFoundationParamsProposal(params *foundation.Params) govtypes.Content {
	return foundation.NewUpdateFoundationParamsProposal("Test", "description", params)
}

func newValidatorAuths(addrs []sdk.ValAddress, allow bool) []*foundation.ValidatorAuth {
	auths := []*foundation.ValidatorAuth{}
	for _, addr := range addrs {
		auth := &foundation.ValidatorAuth{
			OperatorAddress: addr.String(),
			CreationAllowed: allow,
		}
		auths = append(auths, auth)
	}

	return auths
}

func newUpdateValidatorAuthsProposal(auths []*foundation.ValidatorAuth) govtypes.Content {
	return foundation.NewUpdateValidatorAuthsProposal("Test", "description", auths)
}

func TestProposalHandler(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	// turn on the module
	k := app.FoundationKeeper
	params_on := newParams(true)
	k.SetParams(ctx, params_on)

	handler := keeper.NewProposalHandler(k)

	// test adding creation allowed validators
	adding := newValidatorAuths([]sdk.ValAddress{valAddr}, true)
	ap := newUpdateValidatorAuthsProposal(adding)
	require.NoError(t, ap.ValidateBasic())
	require.NoError(t, handler(ctx, ap))
	require.Equal(t, adding, k.GetValidatorAuths(ctx))

	// test deleting creation allowed validators
	deleting := newValidatorAuths([]sdk.ValAddress{valAddr}, false)
	dp := newUpdateValidatorAuthsProposal(deleting)
	require.NoError(t, dp.ValidateBasic())
	require.NoError(t, handler(ctx, dp))
	require.Equal(t, deleting, k.GetValidatorAuths(ctx))

	// disable foundation
	params_off := newParams(false)
	pp := newUpdateFoundationParamsProposal(params_off)
	require.NoError(t, pp.ValidateBasic())
	require.NoError(t, handler(ctx, pp))
	require.Equal(t, []*foundation.ValidatorAuth{}, k.GetValidatorAuths(ctx))
	require.Equal(t, params_off, k.GetParams(ctx))

	// attempt to enable foundation, which fails
	pp = newUpdateFoundationParamsProposal(params_on)
	require.Error(t, pp.ValidateBasic())
}
