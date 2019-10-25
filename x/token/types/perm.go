package types

import (
	"fmt"
)

var _ PermissionI = (*Permission)(nil)

const (
	MintAction = "mint"
	BurnAction = "burn"
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

type ActionPermissionI interface {
	GetAction() string
}

type Permission struct {
	Action   string
	Resource string
}

func (p Permission) Validate() bool {
	if len(p.GetResource()) == 0 || len(p.GetAction()) == 0 {
		return false
	}
	return true
}

func (p Permission) GetResource() string {
	return p.Resource
}

func (p Permission) GetAction() string {
	return p.Action
}

func (p Permission) Equal(res, act string) bool {
	if p.GetResource() == res && p.GetAction() == act {
		return true
	}
	return false
}
func NewMintPermission(resource string) PermissionI {
	return &Permission{
		Action:   MintAction,
		Resource: resource,
	}
}

func NewBurnPermission(resource string) PermissionI {
	return &Permission{
		Action:   BurnAction,
		Resource: resource,
	}
}
