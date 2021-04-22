package v039_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/v2/client"
	v038auth "github.com/line/lbm-sdk/v2/x/auth/legacy/v038"
	v039auth "github.com/line/lbm-sdk/v2/x/auth/legacy/v039"
	v039 "github.com/line/lbm-sdk/v2/x/genutil/legacy/v039"
	"github.com/line/lbm-sdk/v2/x/genutil/types"
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
      "type": "lbm-sdk/Account",
      "value": {
        "address": "link1fgc3542nakex4w8rl5n77zelpdguppfc3f6wdr",
        "coins": [
          {
            "denom": "stake",
            "amount": "400000"
          }
        ],
        "public_key": "linkpub1addwnpepqf2nnrpt3jkwsv6wgs2ndc52y3hyfhwrtldlm80g2hr5z4lzwudgxll3drm",
        "account_number": 1,
        "sequence": 1
      }
    },
    {
      "type": "lbm-sdk/ModuleAccount",
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
      "type": "lbm-sdk/ContinuousVestingAccount",
      "value": {
        "address": "link17kxlwdwmjuzh2y9vnfwyhz2fhfkznzuwt87sc6",
        "coins": [
          {
            "denom": "stake",
            "amount": "10000205"
          }
        ],
        "public_key": "linkpub1addwnpepqw6mp5utlyp85sezqgrvclta80j8zn0p56zvxdexgpgpxh9wvyfeyu683cr",
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
      "type": "lbm-sdk/DelayedVestingAccount",
      "value": {
        "address": "link15q7fcsjxq3j635dk02hqalfu3d2plcp2tw9ngt",
        "coins": [
          {
            "denom": "stake",
            "amount": "10000205"
          }
        ],
        "public_key": "linkpub1addwnpepqg7yvmjtly3kn4qwctxghg8zn8mplydedydww904uf33xca8y4zn2temgn4",
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

var expectedGenAuthState = []byte(`{"params":{"max_memo_characters":"10","tx_sig_limit":"10","tx_size_cost_per_byte":"10","sig_verify_cost_ed25519":"10","sig_verify_cost_secp256k1":"10"},"accounts":[{"type":"lbm-sdk/Account","value":{"address":"link1fgc3542nakex4w8rl5n77zelpdguppfc3f6wdr","coins":[{"denom":"stake","amount":"400000"}],"public_key":{"type":"tendermint/PubKeySecp256k1","value":"AlU5jCuMrOgzTkQVNuKKJG5E3cNf2/2d6FXHQVfidxqD"},"account_number":"1","sequence":"1"}},{"type":"lbm-sdk/ModuleAccount","value":{"address":"link1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3q4fdzl","coins":[{"denom":"stake","amount":"400000000"}],"public_key":"","account_number":"2","sequence":"4","name":"bonded_tokens_pool","permissions":["burner","staking"]}},{"type":"lbm-sdk/ContinuousVestingAccount","value":{"address":"link17kxlwdwmjuzh2y9vnfwyhz2fhfkznzuwt87sc6","coins":[{"denom":"stake","amount":"10000205"}],"public_key":{"type":"tendermint/PubKeySecp256k1","value":"A7Ww04v5AnpDIgIGzH19O+RxTeGmhMM3JkBQE1yuYROS"},"account_number":"3","sequence":"5","original_vesting":[{"denom":"stake","amount":"10000205"}],"delegated_free":[],"delegated_vesting":[],"end_time":"1596125048","start_time":"1595952248"}},{"type":"lbm-sdk/DelayedVestingAccount","value":{"address":"link15q7fcsjxq3j635dk02hqalfu3d2plcp2tw9ngt","coins":[{"denom":"stake","amount":"10000205"}],"public_key":{"type":"tendermint/PubKeySecp256k1","value":"AjxGbkv5I2nUDsLMi6DimfYfkblpGucV9eJjE2OnJUU1"},"account_number":"4","sequence":"15","original_vesting":[{"denom":"stake","amount":"10000205"}],"delegated_free":[],"delegated_vesting":[],"end_time":"1596125048"}}]}`)

func TestMigrate(t *testing.T) {
	genesis := types.AppMap{
		v038auth.ModuleName: genAuthState,
	}

	var migrated types.AppMap
	require.NotPanics(t, func() { migrated = v039.Migrate(genesis, client.Context{}) })
	require.Equal(t, string(expectedGenAuthState), string(migrated[v039auth.ModuleName]))
}
