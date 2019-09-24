# genesis_config
If you want to change the genesis state, consider using this script.
First create `genesis.toml` and then run this script to configure `genesis.json`.

### Quick Start
```bash
$ make install
go install -mod=readonly -tags genesis_config genesis_config.go

$ genesis_config /path/to/genesis.toml -o
2019/09/23 14:32:17 staking.params.bond_denom: stake -> cony
2019/09/23 14:32:17 crisis.constant_fee.denom: stake -> cony
2019/09/23 14:32:17 gov.deposit_params.min_deposit[0].denom : stake -> cony
2019/09/23 14:32:17 mint.params.mint_denom: stake -> cony
2019/09/23 14:32:17 overwrite /home1/irteam/.linkd/config/genesis.json with 4 changes
```

### Build & Install
```bash
$ make build
go build -mod=readonly -tags genesis_config -o build/genesis_config genesis_config.go

$ make install
go install -mod=readonly -tags genesis_config genesis_config.go

$ make clean
rm -rf build
```

### Usage
```bash
$ genesis_config --help
Configure genesis.json with genesis.toml

Usage:
  genesis_config [genesis.toml] [flags]

Flags:
  -h, --help          help for genesis_config
      --home string   node's home directory (default "/home/user/.linkd")
  -o, --overwrite     overwrite the genesis.json file
```

### [genesis.toml](../../genesis.toml)
```bash
# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

##### genesis state configuration options #####

# default denomination
#
# affected modules:
# - crisis: coin denom for constant fee.
# - governance: coin denom of minimum deposit for a proposal to enter voting period.
# - mint: coin denom to mint.
# - staking: bondable coin denom.
denom = "cony"
```
