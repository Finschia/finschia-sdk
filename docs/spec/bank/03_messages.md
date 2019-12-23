# Messages

## MsgSend

```golang
type MsgSend struct {
	FromAddress sdk.AccAddress `json:"from_address" yaml:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address" yaml:"to_address"`
	Amount      sdk.Coins      `json:"amount" yaml:"amount"`
}
```



## MsgMultiSend
```golang
type MsgMultiSend struct {
	Inputs  []Input  `json:"inputs" yaml:"inputs"`
	Outputs []Output `json:"outputs" yaml:"outputs"`
}
```

```golang
type Input struct {
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Coins   sdk.Coins      `json:"coins" yaml:"coins"`
}
```
```golang
type Output struct {
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Coins   sdk.Coins      `json:"coins" yaml:"coins"`
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

