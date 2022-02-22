package consortium_test

import (
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/ed25519"
	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/consortium"
	"github.com/line/lbm-sdk/x/consortium/types"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

var (
	delPk   = ed25519.GenPrivKey().PubKey()
	delAddr = sdk.BytesToAccAddress(delPk.Address())
	valAddr = delAddr.ToValAddress()
)

func newParams(enabled bool) *types.Params {
	return &types.Params{Enabled: enabled}
}

func newUpdateConsortiumParamsProposal(params *types.Params) govtypes.Content {
	return types.NewUpdateConsortiumParamsProposal("Test", "description", params)
}

func newValidatorAuths(addrs []sdk.ValAddress, allow bool) []*types.ValidatorAuth {
	auths := []*types.ValidatorAuth{}
	for _, addr := range addrs {
		auth := &types.ValidatorAuth{
			OperatorAddress: addr.String(),
			CreationAllowed: allow,
		}
		auths = append(auths, auth)
	}

	return auths
}

func newUpdateValidatorAuthsProposal(auths []*types.ValidatorAuth) govtypes.Content {
	return types.NewUpdateValidatorAuthsProposal("Test", "description", auths)
}

func TestProposalHandler(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	// turn on the module
	keeper := app.ConsortiumKeeper
	params_on := newParams(true)
	keeper.SetParams(ctx, params_on)

	handler := consortium.NewProposalHandler(keeper)

	// test adding creation allowed validators
	adding := newValidatorAuths([]sdk.ValAddress{valAddr}, true)
	ap := newUpdateValidatorAuthsProposal(adding)
	require.NoError(t, ap.ValidateBasic())
	require.NoError(t, handler(ctx, ap))
	require.Equal(t, adding, keeper.GetValidatorAuths(ctx))

	// test deleting creation allowed validators
	deleting := newValidatorAuths([]sdk.ValAddress{valAddr}, false)
	dp := newUpdateValidatorAuthsProposal(deleting)
	require.NoError(t, dp.ValidateBasic())
	require.NoError(t, handler(ctx, dp))
	require.Equal(t, deleting, keeper.GetValidatorAuths(ctx))

	// disable consortium
	params_off := newParams(false)
	pp := newUpdateConsortiumParamsProposal(params_off)
	require.NoError(t, pp.ValidateBasic())
	require.NoError(t, handler(ctx, pp))
	require.Equal(t, []*types.ValidatorAuth{}, keeper.GetValidatorAuths(ctx))
	require.Equal(t, params_off, keeper.GetParams(ctx))

	// attempt to enable consortium, which fails
	pp = newUpdateConsortiumParamsProposal(params_on)
	require.Error(t, pp.ValidateBasic())
}
