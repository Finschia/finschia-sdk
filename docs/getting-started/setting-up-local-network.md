# Getting started
For your easy understanding of LINK Network v2, we'll show how to run a blockchain network on your local machine. The goal of this document is to make a network with a single node and then transfer coins on it.

---

You need these two binaries for running a LINK Network:

- `linkd`: LINK full node daemon
- `linkcli`: A tool to run CLI commands including LCD (Light Client Daemon)

Build the source code. Your machine should have Go 1.13+ installed in advance.

```bash
make install
```

The tasks you're going to do here are:

1. Setting up a local network
2. Using the local network

## Setting up a local network
In this section, you will set up and run a blockchain network in 6 steps.

### Create initial files for your network
First, you need initial files that contains network configurations and startup information. `linkd init` command helps you to do it at one stroke.

Create files for your network with the following command.

```bash
linkd init localnetnode --chain-id linklocal
           ------------            ---------
             moniker                chain ID
```

That marked as moniker is the name of the node while chain ID is an identifier of the network. The "--chain-id" flag is optional.
You can find the `config` directory created under the chain data home.

> **Tip**
>
> The chain data home is `$HOME/.linkd` by default. You can change it with "--home" flag.

Each file in the `config` directory contains information to run a local network. Refer to the following table for details.

| File | Description |
|---|---|
| `app.toml` | App configuration |
| `config.toml` | Node configuration. Displayed on the terminal when running `linkd init` |
| `genesis.json` | Onchain data. Common data that every node should have in the beginning. |
| `node_key.json` | Private key of the node |
| `priv_validator_key.json` | Private key for validator |

### Create accounts
You also need at least one account for the network. Account consists of a private/public key pair for signature and an address for asset storage. Use `linkcli keys` command to create accounts.

Type the following commands to create accounts with name of Brown and Sally.

```bash
linkcli keys add brown
linkcli keys add sally
                 -----
                 account name
```

The keys generated here are stored in the `keyring-*` directory under the chain data home. The name of the directory varies with "--keyring-backend" flag.

> **Tip**
>
> The chain data home is `$HOME/.linkd` by default. You can change it with "--home" flag.

### Register genesis accounts
Every blockchain network requires at least one account for initiallly minted coins. We call it genesis account. `linkd add-genesis-account` command registers an existing account as a genesis account.
The first argument is an address or name of the existing account. The second is the coin you want to mint, it's a comma separated coin string which concatenates amount and denom.

Make Brown's account as a genesis account, and mint 1,000 LINK and 100,000,000 stake.

```bash
linkd add-genesis-account $(linkcli keys show brown -a) 1000link,100000000stake
```

> **Tip** 
>
> - The account must be in your local keybase if you use an account name as the first argument. We strongly advise using `link cli keys show` instead of the name itself in case something wrong.
> - Stake is a token made in [Cosmos SDK](https://cosmos.network) to support the advanced PoS. Accounts need to delegate these tokens to bond a validator. 

Now, you can find a genesis account with Brown's key added in `config/genesis.json`.

### Become a validator
For the network to run, one or more validators should exist. Validator is a node that has a permission to consent. To be a validator, an account should be bonded by delegating coins.

We use `linkd gentx` (alias of `linkd create-validator`) for this end. The command generates a genesis transaction, which will be executed right after the network starts. "--name" flag enables you choose the account who bonds to, and the account will sign the transaction with its private key.

Generate a genesis transaction to become a validator, bonding with Brown's private key.

```bash
linkd gentx --name brown
```

Without any other flags, the default parameters are used for delegation amount and commission rate. Check the detailed options with "--help" flag if you want to adjust it.

### Validate the genesis file
Before running a network, you must check whether the genesis file (`config/genesis.json` under the chain data home) is valid or not. The genesis file, which stores onchain data, should have genesis transactions on it.
Gather all genesis transactions generated above with the following command:

```bash
linkd collect-gentxs
```

This will put all genesis transactions in the genesis file. Now, let's validate it.

```bash
linkd validate-genesis
```

Correct it if there's any error occurred.

### Start the network
Finally, everything is ready for running a blockchain network. Type the following command.

```bash
linkd start
```

You can see runtime messages of the network's start.

## Using the local network 
You would concern that your local network works well. This section helps you verify it by querying and transferring on the network.

### Query account information
Let's try a simple query that gets account information.
Check the accounts of Brown and Sally with the following command.

```bash
linkcli query account $(linkcli keys show brown -a)
linkcli query account $(linkcli keys show sally -a)
```

You can see the detailed information of the account you configured in the previous section.

> **Tip**
>
> `linkcli` decides which network to access referring to `config/config.toml` which is in the client home (`$HOME/.linkcli/`).

### Transfer coins
Now, it's able to transfer coins between the accounts shown above.
Send 100 LINK from Brown to Sally with the following command.

```bash
linkcli tx send $(linkcli keys show brown -a) $(linkcli keys show sally -a) 1link
```

Transferring always generates a transaction. The command will show the transaction information, which contains the following attributes:

| Attribute | Description |
|---|---|
| `chain_id` | Chain identifier for the transaction |
| `account_number` | Account number increased whenever an account created |
| `sequence` | Transaction sequence. Used to avoid duplicate transaction |
| `fee` | Fee of the transaction |
| `msg` | Message contained in the transaction. Transaction can have one or more messages in it. See [Message - LINE Blockchain Docs](https://docs-blockchain.line.biz/api-guide/Callback-Response#message) for more information. |
| `memo` | Additional data |

`linkcli` will ask you whether to confirm and broadcast the transaction. Type "yes", then you get a txhash.

### Check the result of transactions
You can check the result of transactions with the txhash you got in the above. Type the following command.

```bash
linkcli query tx [txhash]
```

Then, you can see the status and information of the transaction. Refer to [Raw transaction - LINE Blockchain Docs](https://docs-blockchain.line.biz/api-guide/Callback-Response#raw-transaction) to learn more about transaction information.

### Check the balance of accounts
The final step of transferring is checking the balance of the recipient. Query Sally's account balance to be sure the coins are transferred well.

```bash
linkcli query account $(linkcli keys show sally -a)
```

You can see 100 LINK hold by Sally.
