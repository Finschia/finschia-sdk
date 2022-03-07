package v043_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	v040gov "github.com/line/lbm-sdk/x/gov/legacy/v040"
	v043gov "github.com/line/lbm-sdk/x/gov/legacy/v043"
	"github.com/line/lbm-sdk/x/gov/types"
)

func TestMigrateJSON(t *testing.T) {
	encodingConfig := simapp.MakeTestEncodingConfig()
	clientCtx := client.Context{}.
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithJSONCodec(encodingConfig.Codec)

	voter := sdk.AccAddress("link1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3q4fdzl")
	err := sdk.ValidateAccAddress(voter.String())
	require.NoError(t, err)
	govGenState := &v040gov.GenesisState{
		Votes: v040gov.Votes{
			v040gov.NewVote(1, voter, types.OptionAbstain),
			v040gov.NewVote(2, voter, types.OptionEmpty),
			v040gov.NewVote(3, voter, types.OptionNo),
			v040gov.NewVote(4, voter, types.OptionNoWithVeto),
			v040gov.NewVote(5, voter, types.OptionYes),
		},
	}

	migrated := v043gov.MigrateJSON(govGenState)

	bz, err := clientCtx.JSONCodec.MarshalJSON(migrated)
	require.NoError(t, err)

	// Indent the JSON bz correctly.
	var jsonObj map[string]interface{}
	err = json.Unmarshal(bz, &jsonObj)
	require.NoError(t, err)
	indentedBz, err := json.MarshalIndent(jsonObj, "", "\t")
	require.NoError(t, err)

	// Make sure about:
	// - Votes are all ADR-037 weighted votes with weight 1.
	expected := `{
	"deposit_params": {
		"max_deposit_period": "0s",
		"min_deposit": []
	},
	"deposits": [],
	"proposals": [],
	"starting_proposal_id": "0",
	"tally_params": {
		"quorum": "0",
		"threshold": "0",
		"veto_threshold": "0"
	},
	"votes": [
		{
			"option": "VOTE_OPTION_UNSPECIFIED",
			"options": [
				{
					"option": "VOTE_OPTION_ABSTAIN",
					"weight": "1.000000000000000000"
				}
			],
			"proposal_id": "1",
			"voter": "link1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3q4fdzl"
		},
		{
			"option": "VOTE_OPTION_UNSPECIFIED",
			"options": [
				{
					"option": "VOTE_OPTION_UNSPECIFIED",
					"weight": "1.000000000000000000"
				}
			],
			"proposal_id": "2",
			"voter": "link1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3q4fdzl"
		},
		{
			"option": "VOTE_OPTION_UNSPECIFIED",
			"options": [
				{
					"option": "VOTE_OPTION_NO",
					"weight": "1.000000000000000000"
				}
			],
			"proposal_id": "3",
			"voter": "link1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3q4fdzl"
		},
		{
			"option": "VOTE_OPTION_UNSPECIFIED",
			"options": [
				{
					"option": "VOTE_OPTION_NO_WITH_VETO",
					"weight": "1.000000000000000000"
				}
			],
			"proposal_id": "4",
			"voter": "link1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3q4fdzl"
		},
		{
			"option": "VOTE_OPTION_UNSPECIFIED",
			"options": [
				{
					"option": "VOTE_OPTION_YES",
					"weight": "1.000000000000000000"
				}
			],
			"proposal_id": "5",
			"voter": "link1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3q4fdzl"
		}
	],
	"voting_params": {
		"voting_period": "0s"
	}
}`

	fmt.Println(string(indentedBz))

	require.Equal(t, expected, string(indentedBz))
}
