package v040_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lfb-sdk/client"
	"github.com/line/lfb-sdk/simapp"
	sdk "github.com/line/lfb-sdk/types"
	v039slashing "github.com/line/lfb-sdk/x/slashing/legacy/v039"
	v040slashing "github.com/line/lfb-sdk/x/slashing/legacy/v040"
)

func TestMigrate(t *testing.T) {
	encodingConfig := simapp.MakeTestEncodingConfig()
	clientCtx := client.Context{}.
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithJSONMarshaler(encodingConfig.Marshaler)

	addr1, err := sdk.ConsAddressFromBech32("linkvalcons1hra8nx79ldurl80ddsrlz04q4y8tuenmhw2dua")
	require.NoError(t, err)
	addr2, err := sdk.ConsAddressFromBech32("linkvalcons1twsfmuj28ndph54k4nw8crwu8h9c8mh39vhldx")
	require.NoError(t, err)

	gs := v039slashing.GenesisState{
		Params: v039slashing.DefaultParams(),
		SigningInfos: map[string]v039slashing.ValidatorSigningInfo{
			"linkvalcons1twsfmuj28ndph54k4nw8crwu8h9c8mh39vhldx": {
				Address:             addr2,
				IndexOffset:         615501,
				MissedBlocksCounter: 1,
				Tombstoned:          false,
			},
			"linkvalcons1hra8nx79ldurl80ddsrlz04q4y8tuenmhw2dua": {
				Address:             addr1,
				IndexOffset:         2,
				MissedBlocksCounter: 2,
				Tombstoned:          false,
			},
		},
		MissedBlocks: map[string][]v039slashing.MissedBlock{
			"linkvalcons1twsfmuj28ndph54k4nw8crwu8h9c8mh39vhldx": {
				{
					Index:  2,
					Missed: true,
				},
			},
			"linkvalcons1hra8nx79ldurl80ddsrlz04q4y8tuenmhw2dua": {
				{
					Index:  3,
					Missed: true,
				},
				{
					Index:  4,
					Missed: true,
				},
			},
		},
	}

	migrated := v040slashing.Migrate(gs)
	// Check that in `signing_infos` and `missed_blocks`, the address
	// linkvalcons1hra8nx79ldurl80ddsrlz04q4y8tuenmhw2dua
	// should always come before the address
	// linkvalcons1twsfmuj28ndph54k4nw8crwu8h9c8mh39vhldx
	// (in alphabetic order, basically).
	expected := `{
  "missed_blocks": [
    {
      "address": "linkvalcons1hra8nx79ldurl80ddsrlz04q4y8tuenmhw2dua",
      "missed_blocks": [
        {
          "index": "3",
          "missed": true
        },
        {
          "index": "4",
          "missed": true
        }
      ]
    },
    {
      "address": "linkvalcons1twsfmuj28ndph54k4nw8crwu8h9c8mh39vhldx",
      "missed_blocks": [
        {
          "index": "2",
          "missed": true
        }
      ]
    }
  ],
  "params": {
    "downtime_jail_duration": "600s",
    "min_signed_per_window": "0.500000000000000000",
    "signed_blocks_window": "100",
    "slash_fraction_double_sign": "0.050000000000000000",
    "slash_fraction_downtime": "0.010000000000000000"
  },
  "signing_infos": [
    {
      "address": "linkvalcons1hra8nx79ldurl80ddsrlz04q4y8tuenmhw2dua",
      "validator_signing_info": {
        "address": "linkvalcons1hra8nx79ldurl80ddsrlz04q4y8tuenmhw2dua",
        "index_offset": "2",
        "jailed_until": "0001-01-01T00:00:00Z",
        "missed_blocks_counter": "2",
        "start_height": "0",
        "tombstoned": false
      }
    },
    {
      "address": "linkvalcons1twsfmuj28ndph54k4nw8crwu8h9c8mh39vhldx",
      "validator_signing_info": {
        "address": "linkvalcons1twsfmuj28ndph54k4nw8crwu8h9c8mh39vhldx",
        "index_offset": "615501",
        "jailed_until": "0001-01-01T00:00:00Z",
        "missed_blocks_counter": "1",
        "start_height": "0",
        "tombstoned": false
      }
    }
  ]
}`

	bz, err := clientCtx.JSONMarshaler.MarshalJSON(migrated)
	require.NoError(t, err)

	// Indent the JSON bz correctly.
	var jsonObj map[string]interface{}
	err = json.Unmarshal(bz, &jsonObj)
	require.NoError(t, err)
	indentedBz, err := json.MarshalIndent(jsonObj, "", "  ")
	require.NoError(t, err)

	require.Equal(t, expected, string(indentedBz))
}
