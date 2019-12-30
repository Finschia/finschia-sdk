package types

import (
	"fmt"
)

const (
	ActionWhitelistOperators  = "whitelistOperators"
	ActionWhitelistOtherRoles = "whitelistOtherRoles"

	RoleOwner     = "owner"
	RoleOperator  = "operator"
	RoleAllocator = "allocator"
	RoleIssuer    = "issuer"
	RoleReturner  = "returner"

	ActionAllocate = "allocate"
	ActionRecall   = "recall"
	ActionIssue    = "issue"
	ActionReturn   = "return"

	RegisterRole   = "register"
	DeregisterRole = "deregister"
)

type PermissionI interface {
	GetResource() string
	GetAction() string
	Equal(string, string) bool
}

type Permissions []PermissionI

func (pms Permissions) String() string {
	return fmt.Sprintf("%#v", pms)
}

type Permission struct {
	Action   string `json:"action"`
	Resource string `json:"resource"`
}

func (p Permission) Validate() bool {
	return len(p.GetResource()) > 0 && len(p.GetAction()) > 0
}

func (p Permission) GetResource() string {
	return p.Resource
}

func (p Permission) GetAction() string {
	return p.Action
}

func (p Permission) Equal(res, act string) bool {
	return p.GetResource() == res && p.GetAction() == act
}

func NewWhitelistOperatorsPermission(resource string) PermissionI {
	return &Permission{
		Action:   ActionWhitelistOperators,
		Resource: resource,
	}
}

func NewWhitelistOtherRolesPermission(resource string) PermissionI {
	return &Permission{
		Action:   ActionWhitelistOtherRoles,
		Resource: resource,
	}
}

func NewAllocatePermission(resource string) PermissionI {
	return &Permission{
		Action:   RoleAllocator,
		Resource: resource,
	}
}

func NewRecallPermission(resource string) PermissionI {
	return &Permission{
		Action:   RoleOperator,
		Resource: resource,
	}
}

func NewIssuePermission(resource string) PermissionI {
	return &Permission{
		Action:   RoleIssuer,
		Resource: resource,
	}
}

func NewReturnPermission(resource string) PermissionI {
	return &Permission{
		Action:   RoleReturner,
		Resource: resource,
	}
}
