# Sample tx commands for modules

## Auth

**Create new account**
```
simd keys add user0 --keyring-backend test --home ~/.simapp/simapp0

# check if new account was added successfully
simd keys list --keyring-backend test --home ~/.simapp/simapp0               
```

Let the user0 and validator0 **account address** be each 
* **user0: link1lu5hgjp2gyvgdpf674aklzrpdeuyhjr4fsuqrj**
* **validator0: link146asaycmtydq45kxc8evntqfgepagygelel00h**

If you run multi node, home option's value can be ~/.simapp/simapp1, ~/.simapp/simapp2, ...
You can get same result whatever --home option you use

&nbsp;

## Bank

**Send funds(Bank)**
```
# user0 balances: "0"
simd query bank balances link1lu5hgjp2gyvgdpf674aklzrpdeuyhjr4fsuqrj --home ~/.simapp/simapp0

# validator0 balances: 90000000000stake, 100000000000ukrw
simd query bank balances link146asaycmtydq45kxc8evntqfgepagygelel00h --home ~/.simapp/simapp0

# send 10000stake from validator0 to user0
simd tx bank send link146asaycmtydq45kxc8evntqfgepagygelel00h link1lu5hgjp2gyvgdpf674aklzrpdeuyhjr4fsuqrj 10000000000stake --keyring-backend test --chain-id sim --home ~/.simapp/simapp0

# user0 balances: 10000000000stake
simd query bank balances link1lu5hgjp2gyvgdpf674aklzrpdeuyhjr4fsuqrj --home ~/.simapp/simapp0

# validator0 balances: 80000000000stake, 100000000000ukrw
simd query bank balances link146asaycmtydq45kxc8evntqfgepagygelel00h --home ~/.simapp/simapp0
```

&nbsp;

## Staking

**Staking(deligate)**
```
# Bech32 Val is operator address of validator0
simd debug addr link146asaycmtydq45kxc8evntqfgepagygelel00h --home ~/.simapp/simapp0
```
Let the **validator0 operator address** be **linkvaloper146asaycmtydq45kxc8evntqfgepagygeddajpy**

&nbsp;

```
# deligate 10000000000stake to validator0
simd tx staking delegate linkvaloper146asaycmtydq45kxc8evntqfgepagygeddajpy 10000000000stake 
--from link1lu5hgjp2gyvgdpf674aklzrpdeuyhjr4fsuqrj --keyring-backend test --chain-id sim --home ~/.simapp/simapp0

# check if deligation was successful
simd query staking validators --chain-id sim --home ~/.simapp/simapp0

# undeligate 10000000000stake from validator
simd tx staking unbond linkvaloper146asaycmtydq45kxc8evntqfgepagygeddajpy 10000000000stake --from link1lu5hgjp2gyvgdpf674aklzrpdeuyhjr4fsuqrj --keyring-backend test --chain-id sim --home ~/.simapp/simapp0

# check if undeligation was successful
simd query staking validators --chain-id sim --home ~/.simapp/simapp0
```

&nbsp;

## Gov

**Submit proposal**
```
# genesis config for more efficient testing gov txs
# "max_deposit_period": "120s"
# "voting_period": "120s"

simd tx gov submit-proposal --title="Test Proposal" --description="testing, testing ..." --type="Text" --deposit="10000000stake" --from link146asaycmtydq45kxc8evntqfgepagygelel00h --keyring-backend test --chain-id sim --home ~/.simapp/simapp0 --yes

```
For confirming the proposal
```
simd query gov proposal 1 --chain-id sim --home ~/.simapp/simapp0
```

**Voting a proposal**
```
simd tx gov vote 1 Yes --from link146asaycmtydq45kxc8evntqfgepagygelel00h --keyring-backend test --chain-id sim --home ~/.simapp/simapp0 --yes
simd tx gov vote 1 Yes --from link1twsfmuj28ndph54k4nw8crwu8h9c8mh3rtx705 --keyring-backend test --chain-id sim --home ~/.simapp/simapp0 --yes
simd tx gov vote 1 Yes --from link1008wengr28z5quat2dzrprt9h8euav4herfyum --keyring-backend test --chain-id sim --home ~/.simapp/simapp0 --yes
simd tx gov vote 1 No --from link1h82llw7m5rv05nal6nj92ce7wm6tkq4c4xsk99 --keyring-backend test --chain-id sim --home ~/.simapp/simapp0 --yes
```

Or you can use weighted voting
```
simd tx gov weighted-vote 1 yes=0.5,no=0.3,no_with_veto=0.2 --from link146asaycmtydq45kxc8evntqfgepagygelel00h --keyring-backend test --chain-id sim --home ~/.simapp/simapp0
```

And you can see the voting status
```
simd query gov votes 1 --chain-id sim --home ~/.simapp/simapp0
```

## Upgrade

**Submit update proposal**
```
# genesis config for more efficient testing gov txs
# "max_deposit_period": "120s"
# "voting_period": "120s"
simd tx gov submit-proposal software-upgrade ebony --upgrade-height 200 --upgrade-info "merong" --deposit 100stake --from link146asaycmtydq45kxc8evntqfgepagygelel00h --chain-id sim --home ~/.simapp/simapp0 --keyring-backend test --title "first_time" --description "this is sample upgrade"

# query the proposal
simd query gov proposals

# fulfill the deposit
simd tx gov deposit 1 10000000stake --from link146asaycmtydq45kxc8evntqfgepagygelel00h --keyring-backend test --chain-id sim --home ~/.simapp/simapp0

```

**Vote the proposal**
```
simd tx gov vote 1 yes --from link146asaycmtydq45kxc8evntqfgepagygelel00h --keyring-backend test --chain-id sim --home ~/.simapp/simapp0/
simd tx gov vote 1 yes --from link1twsfmuj28ndph54k4nw8crwu8h9c8mh3rtx705 --keyring-backend test --chain-id sim --home ~/.simapp/simapp0/
simd tx gov vote 1 yes --from link1008wengr28z5quat2dzrprt9h8euav4herfyum --keyring-backend test --chain-id sim --home ~/.simapp/simapp0/
simd tx gov vote 1 yes --from link1h82llw7m5rv05nal6nj92ce7wm6tkq4c4xsk99 --keyring-backend test --chain-id sim --home ~/.simapp/simapp0/
```

**Querying the scheduled plan**
```
# You can query the plan if the proposal would be satisfied the quorum
simd query upgrade plan
```

**Cancel the software-upgrade**
```
# You can cancel the scheduled software-upgrade plan
simd tx gov submit-proposal cancel-software-upgrade --title "first_time" --description "this is sample upgrade" --deposit 100stake --from link146asaycmtydq45kxc8evntqfgepagygelel00h --chain-id sim --home ~/.simapp/simapp0 --keyring-backend test

```

