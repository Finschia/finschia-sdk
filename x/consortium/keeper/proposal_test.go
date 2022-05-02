package keeper_test

import (
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/consortium"
	"github.com/line/lbm-sdk/x/consortium/keeper"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

func newParams(enabled bool) *consortium.Params {
	return &consortium.Params{Enabled: enabled}
}

func newUpdateConsortiumParamsProposal(params *consortium.Params) govtypes.Content {
	return consortium.NewUpdateConsortiumParamsProposal("Test", "description", params)
}

func newValidatorAuths(addrs []sdk.ValAddress, allow bool) []*consortium.ValidatorAuth {
	auths := []*consortium.ValidatorAuth{}
	for _, addr := range addrs {
		auth := &consortium.ValidatorAuth{
			OperatorAddress: addr.String(),
			CreationAllowed: allow,
		}
		auths = append(auths, auth)
	}

	return auths
}

func newUpdateValidatorAuthsProposal(auths []*consortium.ValidatorAuth) govtypes.Content {
	return consortium.NewUpdateValidatorAuthsProposal("Test", "description", auths)
}

func TestProposalHandler(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	// turn on the module
	k := app.ConsortiumKeeper
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

	// disable consortium
	params_off := newParams(false)
	pp := newUpdateConsortiumParamsProposal(params_off)
	require.NoError(t, pp.ValidateBasic())
	require.NoError(t, handler(ctx, pp))
	require.Equal(t, []*consortium.ValidatorAuth{}, k.GetValidatorAuths(ctx))
	require.Equal(t, params_off, k.GetParams(ctx))

	// attempt to enable consortium, which fails
	pp = newUpdateConsortiumParamsProposal(params_on)
	require.Error(t, pp.ValidateBasic())
}
