package v043_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/simapp"
	sdk "github.com/Finschia/finschia-sdk/types"
	v043bank "github.com/Finschia/finschia-sdk/x/bank/legacy/v043"
	"github.com/Finschia/finschia-sdk/x/bank/types"
)

func TestMigrateJSON(t *testing.T) {
	encodingConfig := simapp.MakeTestEncodingConfig()
	clientCtx := client.Context{}.
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithCodec(encodingConfig.Marshaler)

	voter, err := sdk.AccAddressFromBech32("link1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3q4fdzl")
	require.NoError(t, err)
	bankGenState := &types.GenesisState{
		Balances: []types.Balance{
			{
				Address: voter.String(),
				Coins: sdk.Coins{
					sdk.NewCoin("foo", sdk.NewInt(10)),
					sdk.NewCoin("bar", sdk.NewInt(20)),
					sdk.NewCoin("foobar", sdk.NewInt(0)),
				},
			},
		},
		Supply: sdk.Coins{
			sdk.NewCoin("foo", sdk.NewInt(10)),
			sdk.NewCoin("bar", sdk.NewInt(20)),
			sdk.NewCoin("foobar", sdk.NewInt(0)),
			sdk.NewCoin("barfoo", sdk.NewInt(0)),
		},
	}

	migrated := v043bank.MigrateJSON(bankGenState)

	bz, err := clientCtx.Codec.MarshalJSON(migrated)
	require.NoError(t, err)

	// Indent the JSON bz correctly.
	var jsonObj map[string]interface{}
	err = json.Unmarshal(bz, &jsonObj)
	require.NoError(t, err)
	indentedBz, err := json.MarshalIndent(jsonObj, "", "\t")
	require.NoError(t, err)

	// Make sure about:
	// - zero coin balances pruned.
	// - zero supply denoms pruned.
	expected := `{
	"balances": [
		{
			"address": "link1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3q4fdzl",
			"coins": [
				{
					"amount": "20",
					"denom": "bar"
				},
				{
					"amount": "10",
					"denom": "foo"
				}
			]
		}
	],
	"denom_metadata": [],
	"params": {
		"default_send_enabled": false,
		"send_enabled": []
	},
	"supply": [
		{
			"amount": "20",
			"denom": "bar"
		},
		{
			"amount": "10",
			"denom": "foo"
		}
	]
}`

	require.Equal(t, expected, string(indentedBz))
}
