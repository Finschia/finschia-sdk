# Messages

## Create a Safety Box

Create a safety box with `safety_box_id`. Given `safety_box_owner` will be the owner of the safety box.

### MsgSafetyBoxCreate
```golang
type MsgSafetyBoxCreate struct {
	SafetyBoxId    string         `json:"safety_box_id"`
	SafetyBoxOwner sdk.AccAddress `json:"safety_box_owner"`
}
```

## Transfer Coins

Transfer coins to/from the safety box. Require valid roles.

### MsgSafetyBoxAllocateCoins

```golang
type MsgSafetyBoxAllocateCoins struct {
	SafetyBoxId      string         `json:"safety_box_id"`
	AllocatorAddress sdk.AccAddress `json:"allocator_address"`
	Coins            sdk.Coins      `json:"coins"`
}
```

* `AllocatorAddress` allocates `Coins` to the safety box `SafetyBoxId`

### MsgSafetyBoxRecallCoins

```golang
type MsgSafetyBoxRecallCoins struct {
	SafetyBoxId      string         `json:"safety_box_id"`
	AllocatorAddress sdk.AccAddress `json:"allocator_address"`
	Coins            sdk.Coins      `json:"coins"`
}
```

* `AllocatorAddress` recalls `Coins` from the safety box `SafetyBoxId`

### MsgSafetyBoxIssueCoins

```golang
type MsgSafetyBoxIssueCoins struct {
	SafetyBoxId string         `json:"safety_box_id"`
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	Coins       sdk.Coins      `json:"coins"`
}
```

* `FromAddress` request `Coins` issuance to the safety box `SafetyBoxId`.
* `ToAddress` will receive the `Coins` if succeed. 
* `FromAddress` and `ToAddress` may be the same.

### MsgSafetyBoxReturnCoins 

```golang
type MsgSafetyBoxReturnCoins struct {
	SafetyBoxId     string         `json:"safety_box_id"`
	ReturnerAddress sdk.AccAddress `json:"returner_address"`
	Coins           sdk.Coins      `json:"coins"`
}
```

* `ReturnerAddress` returns `Coins` to the safety box `SafetyBoxId`

## Grant/Revoke Roles

Grant or revoke roles. Role is required to perform actions.

### MsgSafetyBoxRegisterOperator

```golang
type MsgSafetyBoxRegisterOperator struct {
	SafetyBoxId    string         `json:"safety_box_id"`
	SafetyBoxOwner sdk.AccAddress `json:"safety_box_owner"`
	Address        sdk.AccAddress `json:"address"`
}
```

### MsgSafetyBoxDeregisterOperator

```golang
type MsgSafetyBoxDeregisterOperator struct {
	SafetyBoxId    string         `json:"safety_box_id"`
	SafetyBoxOwner sdk.AccAddress `json:"safety_box_owner"`
	Address        sdk.AccAddress `json:"address"`
}
```

* `SafetyBoxOwner` register or deregister `Address` to the safety box `SafetyBoxId` as operator. 

### MsgSafetyBoxRegisterAllocator

```golang
type MsgSafetyBoxRegisterAllocator struct {
	SafetyBoxId string         `json:"safety_box_id"`
	Operator    sdk.AccAddress `json:"operator"`
	Address     sdk.AccAddress `json:"address"`
}
``` 

### MsgSafetyBoxDeregisterAllocator

```golang
type MsgSafetyBoxDeregisterAllocator struct {
	SafetyBoxId string         `json:"safety_box_id"`
	Operator    sdk.AccAddress `json:"operator"`
	Address     sdk.AccAddress `json:"address"`
}
```

* `Operator` register or deregister `Address` to the safety box `SafetyBoxId` as allocator.
* Allocator gets permission to `allocate` and `recall`.

### MsgSafetyBoxRegisterIssuer

```golang
type MsgSafetyBoxRegisterIssuer struct {
	SafetyBoxId string         `json:"safety_box_id"`
	Operator    sdk.AccAddress `json:"operator"`
	Address     sdk.AccAddress `json:"address"`
}
```

### MsgSafetyBoxDeregisterIssuer

```golang
type MsgSafetyBoxDeregisterIssuer struct {
	SafetyBoxId string         `json:"safety_box_id"`
	Operator    sdk.AccAddress `json:"operator"`
	Address     sdk.AccAddress `json:"address"`
}
```

* `Operator` register or deregister `Address` to the safety box `SafetyBoxId` as issuer.

### MsgSafetyBoxRegisterReturner

```golang
type MsgSafetyBoxRegisterReturner struct {
	SafetyBoxId string         `json:"safety_box_id"`
	Operator    sdk.AccAddress `json:"operator"`
	Address     sdk.AccAddress `json:"address"`
}
```

### MsgSafetyBoxDeregisterReturner

```golang
type MsgSafetyBoxDeregisterReturner struct {
	SafetyBoxId string         `json:"safety_box_id"`
	Operator    sdk.AccAddress `json:"operator"`
	Address     sdk.AccAddress `json:"address"`
}
```

* `Operator` register or deregister `Address` to the safety box `SafetyBoxId` as returner.