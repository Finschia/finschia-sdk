package foundation_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

func TestUpdateFoundationParamsProposal(t *testing.T) {
	title := "hello"
	description := "world"

	testCases := map[string]struct {
		params      foundation.Params
		valid       bool
		stringified string
	}{
		"valid proposal": {
			params: foundation.Params{
				FoundationTax: sdk.OneDec(),
				CensoredMsgTypeUrls: []string{
					sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
				},
			},
			valid: true,
			stringified: fmt.Sprintf(`title: %s
description: %s
params:
  foundationtax: "%s"
  censoredmsgtypeurls:
  - %s
`, title, description, sdk.OneDec(), sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil))),
		},
		"invalid tax rate": {
			params: foundation.Params{
				FoundationTax: sdk.NewDec(2),
			},
		},
	}

	for name, tc := range testCases {
		proposal := foundation.NewUpdateFoundationParamsProposal(title, description, tc.params)

		require.Equal(t, title, proposal.GetTitle(), name)
		require.Equal(t, description, proposal.GetDescription(), name)
		require.Equal(t, foundation.RouterKey, proposal.ProposalRoute(), name)
		require.Equal(t, foundation.ProposalTypeUpdateFoundationParams, proposal.ProposalType(), name)

		err := proposal.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			continue
		}
		require.NoError(t, err, name)

		require.Equal(t, tc.stringified, proposal.String(), name)
	}
}
