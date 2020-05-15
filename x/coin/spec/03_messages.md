# Messages

## MsgSend

```golang
type MsgSend struct {
	From    sdk.AccAddress `json:"from"`
	To      sdk.AccAddress `json:"to"`
	Amount  sdk.Coins      `json:"amount"`
}
```

## MsgMultiSend
```golang
type MsgMultiSend struct {
	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`
}
```

```golang
type Input struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}
```
```golang
type Output struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}
```

```
handleMsgSend(msg MsgSend)
  inputSum = 0
  for input in inputs
    inputSum += input.Amount
  outputSum = 0
  for output in outputs
    outputSum += output.Amount
  if inputSum != outputSum:
    fail with "input/output amount mismatch"

  return inputOutputCoins(msg.Inputs, msg.Outputs)
```

# Syntax
| Message/Attributes | Tag | Type |
| ---- | ---- | ---- |
| Message | coin/MsgSend | github.com/line/link/x/coin/internal/types.MsgSend |  
 | Attributes | from | []uint8 |  
 | Attributes | to | []uint8 |  
 | Attributes | amount | []github.com/cosmos/cosmos-sdk/types.Coin |  
| Message | coin/MsgMultiSend | github.com/line/link/x/coin/internal/types.MsgMultiSend |  
 | Attributes | inputs | []github.com/line/link/x/coin/internal/types.Input |  
 | Attributes | outputs | []github.com/line/link/x/coin/internal/types.Output |  