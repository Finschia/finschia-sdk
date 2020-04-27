# Messages

## Create a Safety Box

Create a safety box with `safety_box_id`. Given `safety_box_owner` will be the owner of the safety box.

### MsgSafetyBoxCreate
```golang
type MsgSafetyBoxCreate struct {
	SafetyBoxId    string         `json:"safety_box_id"`
	SafetyBoxOwner sdk.AccAddress `json:"safety_box_owner"`
	ContractID     string         `json:"contract_id"`
}
```

## Transfer Token

Transfer token to/from the safety box. Require valid roles.

### MsgSafetyBoxAllocateToken

```golang
type MsgSafetyBoxAllocateToken struct {
	SafetyBoxId      string         `json:"safety_box_id"`
	AllocatorAddress sdk.AccAddress `json:"allocator_address"`
	ContractID       string         `json:"contract_id"`
	Amount           sdk.Int        `json:"amount"`
}
```

* `AllocatorAddress` allocates `token` to the safety box `SafetyBoxId`

### MsgSafetyBoxRecallToken

```golang
type MsgSafetyBoxRecallToken struct {
	SafetyBoxId      string         `json:"safety_box_id"`
	AllocatorAddress sdk.AccAddress `json:"allocator_address"`
	ContractID       string         `json:"contract_id"`
	Amount           sdk.Int        `json:"amount"`
}
```

* `AllocatorAddress` recalls `token` from the safety box `SafetyBoxId`

### MsgSafetyBoxIssueToken

```golang
type MsgSafetyBoxIssueToken struct {
	SafetyBoxId string         `json:"safety_box_id"`
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	ContractID  string         `json:"contract_id"`
	Amount      sdk.Int        `json:"amount"`
}
```

* `FromAddress` request `token` issuance to the safety box `SafetyBoxId`.
* `ToAddress` will receive the `token` if succeed. 
* `FromAddress` and `ToAddress` may be the same.

### MsgSafetyBoxReturnToken

```golang
type MsgSafetyBoxReturnToken struct {
	SafetyBoxId     string         `json:"safety_box_id"`
	ReturnerAddress sdk.AccAddress `json:"returner_address"`
	ContractID      string         `json:"contract_id"`
	Amount          sdk.Int        `json:"amount"`
}
```

* `ReturnerAddress` returns `token` to the safety box `SafetyBoxId`

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
