package foundation_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

func TestUpdateFoundationParamsProposal(t *testing.T) {
	testCases := map[string]struct {
		params *foundation.Params
		valid  bool
	}{
		"valid proposal": {
			params: foundation.DefaultParams(),
			valid:  true,
		},
		"attempt to enable foundation": {
			params: &foundation.Params{
				Enabled:       true,
				FoundationTax: sdk.ZeroDec(),
			},
		},
		"invalid tax rate": {
			params: &foundation.Params{
				FoundationTax: sdk.NewDec(2),
			},
		},
	}

	for name, tc := range testCases {
		proposal := foundation.NewUpdateFoundationParamsProposal("", "", tc.params)

		require.Empty(t, proposal.GetTitle(), name)
		require.Empty(t, proposal.GetDescription(), name)
		require.Equal(t, foundation.RouterKey, proposal.ProposalRoute(), name)
		require.Equal(t, foundation.ProposalTypeUpdateFoundationParams, proposal.ProposalType(), name)

		err := proposal.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestUpdateValidatorAuthsProposal(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		updates []foundation.ValidatorAuth
		valid   bool
	}{
		"valid proposal": {
			updates: []foundation.ValidatorAuth{{
				OperatorAddress: sdk.ValAddress(addrs[0]).String(),
				CreationAllowed: true,
			}},
			valid: true,
		},
		"empty auths": {
			updates: []foundation.ValidatorAuth{},
		},
		"invalid address": {
			updates: []foundation.ValidatorAuth{{}},
		},
		"duplicate addresses": {
			updates: []foundation.ValidatorAuth{
				{
					OperatorAddress: sdk.ValAddress(addrs[0]).String(),
					CreationAllowed: true,
				},
				{
					OperatorAddress: sdk.ValAddress(addrs[0]).String(),
				},
			},
		},
	}

	for name, tc := range testCases {
		proposal := foundation.NewUpdateValidatorAuthsProposal("", "", tc.updates)

		require.Empty(t, proposal.GetTitle(), name)
		require.Empty(t, proposal.GetDescription(), name)
		require.Equal(t, foundation.RouterKey, proposal.ProposalRoute(), name)
		require.Equal(t, foundation.ProposalTypeUpdateValidatorAuths, proposal.ProposalType(), name)

		err := proposal.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}
