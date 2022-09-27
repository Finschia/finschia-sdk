package foundation_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

func TestUpdateFoundationParamsProposal(t *testing.T) {
	testCases := map[string]struct {
		params foundation.Params
		valid  bool
	}{
		"valid proposal": {
			params: foundation.DefaultParams(),
			valid:  true,
		},
		"invalid tax rate": {
			params: foundation.Params{
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
