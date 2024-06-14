<!--
order: 2
-->

# State

The `x/fswap` module keeps state of three primary objects, Swap, SwapStats and Swapped.

## Swap

- Swap: `0x01 + (lengthPrefixed+)fromDenom + (lengthPrefixed+)toDenom`


## SwapStats

- SwapStats: `0x02`

## Swapped

- Swapped: `0x03 + (lengthPrefixed+)fromDenom + (lengthPrefixed+)toDenom`

