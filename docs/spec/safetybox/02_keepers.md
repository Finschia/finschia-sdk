# Keepers


## Common Types

### SafetyBox

```golang
type SafetyBox struct {
	ID                   string         `json:"id"`
	Owner                sdk.AccAddress `json:"owner"`
	Address              sdk.AccAddress `json:"address"`
	TotalAllocation      sdk.Coins      `json:"total_allocation"`
	CumulativeAllocation sdk.Coins      `json:"cumulative_allocation"`
	TotalIssuance        sdk.Coins      `json:"total_issuance"`
}
```

### Permission

```golang
type Permission struct {
   	Action   string
   	Resource string
}
```

## Keeper

Note: SafetyBoxKeeper is kept internal

```golang
type SafetyBoxKeeper interface {
    NewSafetyBox(sdk.Context, types.MsgSafetyBoxCreate) (types.SafetyBox, sdk.Error)
    GetSafetyBox(sdk.Context, string) (types.SafetyBox, sdk.Error)
    Allocate(sdk.Context, types.MsgSafetyBoxAllocateCoins) sdk.Error
    Recall(sdk.Context, types.MsgSafetyBoxRecallCoins) sdk.Error
    Issue(sdk.Context, types.MsgSafetyBoxIssueCoins) sdk.Error
    Return(sdk.Context, types.MsgSafetyBoxReturnCoins) sdk.Error
    GrantPermission(sdk.Context, string, sdk.AccAddress, sdk.AccAddress, string) sdk.Error
    RevokePermission(sdk.Context, string, sdk.AccAddress, sdk.AccAddress, string) sdk.Error
    GetPermissions(sdk.Context, string, string, sdk.AccAddress) (types.MsgSafetyBoxRoleResponse, sdk.Error)
    IsOwner(sdk.Context, string, sdk.AccAddress) bool
    IsOperator(sdk.Context, string, sdk.AccAddress) bool
    IsAllocator(sdk.Context, string, sdk.AccAddress) bool
    IsIssuer(sdk.Context, string, sdk.AccAddress) bool
    IsReturner(sdk.Context, string, sdk.AccAddress) bool
}
```
