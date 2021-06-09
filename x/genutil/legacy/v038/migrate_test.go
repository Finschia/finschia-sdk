package v038_test

import (
	"testing"

	"github.com/line/lfb-sdk/client"
	v036auth "github.com/line/lfb-sdk/x/auth/legacy/v036"
	v036genaccounts "github.com/line/lfb-sdk/x/genaccounts/legacy/v036"
	v038 "github.com/line/lfb-sdk/x/genutil/legacy/v038"
	"github.com/line/lfb-sdk/x/genutil/types"
	v036staking "github.com/line/lfb-sdk/x/staking/legacy/v036"

	"github.com/stretchr/testify/require"
)

var genAccountsState = []byte(`[
		{
			"account_number": "0",
			"address": "link1vu6qhekl7ve5gk5exzy36tauyrk96mkphrf5u8",
			"coins": [
				{
					"amount": "1000000000",
					"denom": "node0token"
				},
				{
					"amount": "400000198",
					"denom": "stake"
				}
			],
			"delegated_free": [],
			"delegated_vesting": [],
			"end_time": "0",
			"module_name": "",
			"module_permissions": [],
			"original_vesting": [],
			"sequence_number": "1",
			"start_time": "0"
		},
		{
			"account_number": "0",
			"address": "link1tygms3xhhs3yv487phx3dw4a95jn7t7l544u5t",
			"coins": [],
			"delegated_free": [],
			"delegated_vesting": [],
			"end_time": "0",
			"module_name": "not_bonded_tokens_pool",
			"module_permissions": [
				"burner",
				"staking"
			],
			"original_vesting": [],
			"sequence_number": "0",
			"start_time": "0"
		},
		{
			"account_number": "0",
			"address": "link1m3h30wlvsf8llruxtpukdvsy0km2kum8al86ug",
			"coins": [],
			"delegated_free": [],
			"delegated_vesting": [],
			"end_time": "0",
			"module_name": "mint",
			"module_permissions": [
				"minter"
			],
			"original_vesting": [],
			"sequence_number": "0",
			"start_time": "0"
		}
	]`)

var genAuthState = []byte(`{
  "params": {
    "max_memo_characters": "256",
    "sig_verify_cost_ed25519": "590",
    "sig_verify_cost_secp256k1": "1000",
    "tx_sig_limit": "7",
    "tx_size_cost_per_byte": "10"
  }
}`)

var genStakingState = []byte(`{
  "delegations": [
    {
      "delegator_address": "link1vu6qhekl7ve5gk5exzy36tauyrk96mkphrf5u8",
      "shares": "100000000.000000000000000000",
      "validator_address": "linkvaloper19ex3qzs2yy73xjfth9dgxxah6pqjdyjtckzx7y"
    }
  ],
  "exported": true,
  "last_total_power": "400",
  "last_validator_powers": [
    {
      "Address": "linkvaloper1twsfmuj28ndph54k4nw8crwu8h9c8mh33lyrp8",
      "Power": "100"
    }
  ],
  "params": {
    "bond_denom": "stake",
    "max_entries": 7,
    "max_validators": 100,
    "unbonding_time": "259200000000000"
  },
  "redelegations": null,
  "unbonding_delegations": null,
  "validators": [
    {
      "commission": {
        "commission_rates": {
          "max_change_rate": "0.000000000000000000",
          "max_rate": "0.000000000000000000",
          "rate": "0.000000000000000000"
        },
        "update_time": "2019-09-24T23:11:22.9692177Z"
      },
      "consensus_pubkey": "linkvalconspub1cqmsrdepqt2qvn9hxdpnqpwp33hw4znk5stakfp67auztzrqhs5xxvfah6fag55tu22",
      "delegator_shares": "100000000.000000000000000000",
      "description": {
        "details": "",
        "identity": "",
        "moniker": "node0",
        "website": ""
      },
      "jailed": false,
      "min_self_delegation": "1",
      "operator_address": "linkvaloper19ex3qzs2yy73xjfth9dgxxah6pqjdyjtckzx7y",
      "status": 2,
      "tokens": "100000000",
      "unbonding_height": "0",
      "unbonding_time": "1970-01-01T00:00:00Z"
    }
  ]
}`)

func TestMigrate(t *testing.T) {
	genesis := types.AppMap{
		v036auth.ModuleName:        genAuthState,
		v036genaccounts.ModuleName: genAccountsState,
		v036staking.ModuleName:     genStakingState,
	}

	require.NotPanics(t, func() { v038.Migrate(genesis, client.Context{}) })
}
