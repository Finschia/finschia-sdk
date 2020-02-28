package types

import (
	"fmt"
)

var _ PermissionI = (*Permission)(nil)

const (
	MintAction   = "mint"
	BurnAction   = "burn"
	IssueAction  = "issue"
	ModifyAction = "modify"
)

type PermissionI interface {
	GetResource() string
	GetAction() string
	Equal(string, string) bool
	String() string
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

func (p Permission) String() string {
	return fmt.Sprintf("%s-%s", p.GetResource(), p.GetAction())
}

func NewMintPermission(contractID string) Permission {
	return Permission{
		Action:   MintAction,
		Resource: contractID,
	}
}

func NewBurnPermission(contractID string) Permission {
	return Permission{
		Action:   BurnAction,
		Resource: contractID,
	}
}

func NewIssuePermission(contractID string) Permission {
	return Permission{
		Action:   IssueAction,
		Resource: contractID,
	}
}

func NewModifyPermission(contractID string) Permission {
	return Permission{
		Action:   ModifyAction,
		Resource: contractID,
	}
}
