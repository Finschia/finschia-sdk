package v039_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lfb-sdk/client"
	v038auth "github.com/line/lfb-sdk/x/auth/legacy/v038"
	v039auth "github.com/line/lfb-sdk/x/auth/legacy/v039"
	v039 "github.com/line/lfb-sdk/x/genutil/legacy/v039"
	"github.com/line/lfb-sdk/x/genutil/types"
)

var genAuthState = []byte(`{
  "params": {
    "max_memo_characters": "10",
    "tx_sig_limit": "10",
    "tx_size_cost_per_byte": "10",
    "sig_verify_cost_ed25519": "10",
    "sig_verify_cost_secp256k1": "10"
  },
  "accounts": [
    {
      "type": "lfb-sdk/Account",
      "value": {
        "address": "link1vncp8z0kqt52406m5aq8f5tgw7r62hy9sdc7ts",
        "coins": [
          {
            "denom": "stake",
            "amount": "400000"
          }
        ],
        "public_key": "linkpub1cqmsrdepqwygwv232a90sgk5k5wkdq990sg2r27wn5p7kc2cemm2yq50fvh52j2swpu",
        "account_number": 1,
        "sequence": 1
      }
    },
    {
      "type": "lfb-sdk/ModuleAccount",
      "value": {
        "address": "link1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3q4fdzl",
        "coins": [
          {
            "denom": "stake",
            "amount": "400000000"
          }
        ],
        "public_key": "",
        "account_number": 2,
        "sequence": 4,
        "name": "bonded_tokens_pool",
        "permissions": [
          "burner",
          "staking"
        ]
      }
    },
    {
      "type": "lfb-sdk/ContinuousVestingAccount",
      "value": {
        "address": "link1u3kdlk3ygg54x8xargfaj5qsxr94rzrpslcjla",
        "coins": [
          {
            "denom": "stake",
            "amount": "10000205"
          }
        ],
        "public_key": "linkpub1cqmsrdepqfh73u2uw0qe4kvptrajmvnep7gvhzc3fuh745scvg7zewl4l76zct2qge9",
        "account_number": 3,
        "sequence": 5,
        "original_vesting": [
          {
            "denom": "stake",
            "amount": "10000205"
          }
        ],
        "delegated_free": [],
        "delegated_vesting": [],
        "end_time": 1596125048,
        "start_time": 1595952248
      }
    },
    {
      "type": "lfb-sdk/DelayedVestingAccount",
      "value": {
        "address": "link1yxyannuykr395nhee4rnzkxq9keexd7xgmqrqd",
        "coins": [
          {
            "denom": "stake",
            "amount": "10000205"
          }
        ],
        "public_key": "linkpub1cqmsrdepqt2qvn9hxdpnqpwp33hw4znk5stakfp67auztzrqhs5xxvfah6fagud2wcd",
        "account_number": 4,
        "sequence": 15,
        "original_vesting": [
          {
            "denom": "stake",
            "amount": "10000205"
          }
        ],
        "delegated_free": [],
        "delegated_vesting": [],
        "end_time": 1596125048
      }
    }
  ]
}`)

var expectedGenAuthState = []byte(`{"params":{"max_memo_characters":"10","tx_sig_limit":"10","tx_size_cost_per_byte":"10","sig_verify_cost_ed25519":"10","sig_verify_cost_secp256k1":"10"},"accounts":[{"type":"lfb-sdk/Account","value":{"address":"link1vncp8z0kqt52406m5aq8f5tgw7r62hy9sdc7ts","coins":[{"denom":"stake","amount":"400000"}],"public_key":{"type":"ostracon/PubKeySecp256k1","value":"A4iHMVFXSvgi1LUdZoClfBChq86dA+thWM72ogKPSy9F"},"account_number":"1","sequence":"1"}},{"type":"lfb-sdk/ModuleAccount","value":{"address":"link1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3q4fdzl","coins":[{"denom":"stake","amount":"400000000"}],"public_key":"","account_number":"2","sequence":"4","name":"bonded_tokens_pool","permissions":["burner","staking"]}},{"type":"lfb-sdk/ContinuousVestingAccount","value":{"address":"link1u3kdlk3ygg54x8xargfaj5qsxr94rzrpslcjla","coins":[{"denom":"stake","amount":"10000205"}],"public_key":{"type":"ostracon/PubKeySecp256k1","value":"Am/o8VxzwZrZgVj7LbJ5D5DLixFPL+rSGGI8LLv1/7Qs"},"account_number":"3","sequence":"5","original_vesting":[{"denom":"stake","amount":"10000205"}],"delegated_free":[],"delegated_vesting":[],"end_time":"1596125048","start_time":"1595952248"}},{"type":"lfb-sdk/DelayedVestingAccount","value":{"address":"link1yxyannuykr395nhee4rnzkxq9keexd7xgmqrqd","coins":[{"denom":"stake","amount":"10000205"}],"public_key":{"type":"ostracon/PubKeySecp256k1","value":"AtQGTLczQzAFwYxu6op2pBfbJDr3eCWIYLwoYzE9vpPU"},"account_number":"4","sequence":"15","original_vesting":[{"denom":"stake","amount":"10000205"}],"delegated_free":[],"delegated_vesting":[],"end_time":"1596125048"}}]}`)

func TestMigrate(t *testing.T) {
	genesis := types.AppMap{
		v038auth.ModuleName: genAuthState,
	}

	var migrated types.AppMap
	require.NotPanics(t, func() { migrated = v039.Migrate(genesis, client.Context{}) })
	require.Equal(t, string(expectedGenAuthState), string(migrated[v039auth.ModuleName]))
}
