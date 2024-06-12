<!--
order: 1
-->

# Concepts

## Swap


The `x/fswap` module defines a `Swap` type in which a coin is allowed to be swapped into another coin on the chain.

```go
type Swap struct {
  FromDenom             string
  ToDenom               string
  AmountCapForToDenom   sdk.Int
  SwapRate              sdk.Dec
}
```

## Config

The `x/fswap` module defines a `Config` type for managing the maximum number of Swaps allowed on chain  through `MaxSwaps`. Additionally, `UpdateAllowed` specifies whether `Swap` can be modified.

```go
type Config struct {
	MaxSwaps      int
	UpdateAllowed bool
}
```

## Proposal

Typically, a `Swap` is proposed and submitted through foundation via a `MsgSetSwap`.
This proposal prescribes to the standard foundation process. If the proposal passes, the `Swap` can be used on chain.

```go
type MsgSetSwap struct {
	Authority           string
    Swap                Swap
    ToDenomMetadata     bank.Metadata
}
```
