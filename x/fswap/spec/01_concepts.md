<!--
order: 1
-->

# Concepts

## Swap


The `x/fswap` module defines a `Swap` type in which a coin is allowed to be swapped into another coin on the chain.
You could find detailed information in the [Protobuf reference](../../../proto/lbm/fswap/v1/fswap.proto#L9-L16) 

```go
type Swap struct {
  FromDenom             string
  ToDenom               string
  AmountCapForToDenom   sdk.Int
  SwapRate              sdk.Dec
}
```

Anyone could use one of the following two transcations to swap `FromDedenom` to `ToDenom`.
1. `simd tx fswap swap [from] [from_coin_amount] [to_denom]`
    - this transcation could swap a specified amount of `from_denom` via [`MsgSwap`](../../../proto/lbm/fswap/v1/tx.proto#L17-L24)
2. `simd tx fswap swap-all [from_address] [from_denom] [to_denom]`
    - this transcation could swap all of `from_denom` under `from_address` via [`MsgSwapAll`](../../../proto/lbm/fswap/v1/tx.proto#L28-L33)

When the swap is triggered, the following event will occur:
1. `from_denom` will be sent from `from_address` to `x/fswap` module
2. `x/fswap` module will burn `from_denom`
3. `x/fswap` module will mint `to_denom` as amount as `from_denom * swapRate`
4. these `to_denom` will sent to `from_address`
5. `EventSwapCoins` will be emitted

## Config

The `x/fswap` module defines a `Config` type for managing the maximum number of Swaps allowed on chain through `MaxSwaps`. Additionally, `UpdateAllowed` specifies whether `Swap` can be modified.

```go
type Config struct {
	MaxSwaps      int
	UpdateAllowed bool
}
```

## MsgSetSwap

Other modules can include `MsgSetSwap` in their proposals to set `Swap`. If the proposal passes, the `Swap` can be used on chain.

`ToDenomMetadata` is [`Metadata`](../../bank/types/bank.pb.go#L325) in `x/bank` module, and it MUST meet these [limitations](../../bank/types/metadata.go#L11). 
In addition, `ToDenomMetadata` should also meet the following two additional constraints by x/swap.
1. `Base` should be consistent with `ToDenom` in `Swap` ([valiation](../types/msgs.go#L121-L123))
2. It cannot override existing denom metadata ([valiation](../keeper/keeper.go#L169))

The following example illustrates the use of `MsgSetSwap` within the `x/foundation` module. `Authority` is a spec in the `x/foundation` module, and you can get more information [here](../../foundation/README.md#L54).

```go
type MsgSetSwap struct {
    Authority           string
    Swap                Swap
    ToDenomMetadata     bank.Metadata
}
```
