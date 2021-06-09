package v036_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lfb-sdk/codec"
	"github.com/line/lfb-sdk/crypto/keys/ed25519"
	"github.com/line/lfb-sdk/crypto/keys/secp256k1"
	v034staking "github.com/line/lfb-sdk/x/staking/legacy/v034"
	v036staking "github.com/line/lfb-sdk/x/staking/legacy/v036"
)

func TestMigrate(t *testing.T) {
	aminoCdc := codec.NewLegacyAmino()
	consPubKeyEd := ed25519.GenPrivKeyFromSecret([]byte("val0")).PubKey()
	consPubKeySecp := secp256k1.GenPrivKeyFromSecret([]byte("val1")).PubKey()
	stakingGenState := v034staking.GenesisState{
		Validators: v034staking.Validators{
			v034staking.Validator{
				ConsPubKey: consPubKeyEd,
				Status:     v034staking.Unbonded,
			}, v034staking.Validator{
				ConsPubKey: consPubKeySecp,
				Status:     v034staking.Unbonded,
			},
		},
	}

	migrated := v036staking.Migrate(stakingGenState)

	json, err := aminoCdc.MarshalJSONIndent(migrated, "", "  ")
	require.NoError(t, err)

	expectedJSON := `{
  "params": {
    "unbonding_time": "0",
    "max_validators": 0,
    "max_entries": 0,
    "bond_denom": ""
  },
  "last_total_power": "0",
  "last_validator_powers": null,
  "validators": [
    {
      "operator_address": "",
      "consensus_pubkey": "linkvalconspub1ef0jhvpq9ymett3nlv6fytn7lqxzd3q3ckvd79eqlcf3wkhgamcl4rzghesqd3x3xc",
      "jailed": false,
      "status": 0,
      "tokens": "0",
      "delegator_shares": "0",
      "description": {
        "moniker": "",
        "identity": "",
        "website": "",
        "details": ""
      },
      "unbonding_height": "0",
      "unbonding_time": "0001-01-01T00:00:00Z",
      "commission": {
        "commission_rates": {
          "rate": "0",
          "max_rate": "0",
          "max_change_rate": "0"
        },
        "update_time": "0001-01-01T00:00:00Z"
      },
      "min_self_delegation": "0"
    },
    {
      "operator_address": "",
      "consensus_pubkey": "linkvalconspub1cqmsrdepqwfxk5k5pugwz3quqyzvzupefm3589tw6x9dkzjdkuzn7hgpz33ag9gvw6q",
      "jailed": false,
      "status": 0,
      "tokens": "0",
      "delegator_shares": "0",
      "description": {
        "moniker": "",
        "identity": "",
        "website": "",
        "details": ""
      },
      "unbonding_height": "0",
      "unbonding_time": "0001-01-01T00:00:00Z",
      "commission": {
        "commission_rates": {
          "rate": "0",
          "max_rate": "0",
          "max_change_rate": "0"
        },
        "update_time": "0001-01-01T00:00:00Z"
      },
      "min_self_delegation": "0"
    }
  ],
  "delegations": null,
  "unbonding_delegations": null,
  "redelegations": null,
  "exported": false
}`

	require.Equal(t, expectedJSON, string(json))
}
