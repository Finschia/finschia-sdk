```
make build
make install
zsh init_node.sh sim 1
```

open the `./.simapp/simapp0/config/genesis.json` change the value of  `voting_params.voting_period` to `10s`
```
simd start --home ~/.simapp/simapp0 
simd tx gov submit-proposal swap-init --title "test" --description "test" --from link146asaycmtydq45kxc8evntqfgepagygelel00h --from-denom "cony" --to-denom "PDT" --swap-rate 123 --amount-limit 1000 --deposit 10000000stake --chain-id=sim --keyring-backend=test --gas-prices 1000stake --gas 10000000 --gas-adjustment 1.5 --home ~/.simapp/simapp0 -b block -y
simd tx gov vote 1 yes  --from link146asaycmtydq45kxc8evntqfgepagygelel00h --chain-id=sim --keyring-backend=test --home ~/.simapp/simapp0 -b block -y
```

```
simd query fswap swapped --chain-id=sim
```