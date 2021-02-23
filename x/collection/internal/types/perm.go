package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	MintAction   = "mint"
	BurnAction   = "burn"
	IssueAction  = "issue"
	ModifyAction = "modify"
)

type Permission string

func NewMintPermission() Permission {
	return MintAction
}

func NewBurnPermission() Permission {
	return BurnAction
}

func NewIssuePermission() Permission {
	return IssueAction
}

func NewModifyPermission() Permission {
	return ModifyAction
}

func (p Permission) Equal(p2 Permission) bool {
	return p == p2
}
func (p Permission) String() string {
	return string(p)
}

func (p Permission) Validate() bool {
	if p == MintAction {
		return true
	}
	if p == BurnAction {
		return true
	}
	if p == IssueAction {
		return true
	}
	if p == ModifyAction {
		return true
	}
	return false
}

type Permissions []Permission

func NewPermissions(perms ...Permission) Permissions {
	pms := Permissions{}
	for _, perm := range perms {
		pms.AddPermission(perm)
	}
	return pms
}

func (pms *Permissions) GetPermissions() []Permission {
	return []Permission(*pms)
}

func (pms *Permissions) RemoveElement(idx int) {
	*pms = append((*pms)[:idx], (*pms)[idx+1:]...)
}

func (pms *Permissions) AddPermission(p Permission) {
	for _, pin := range *pms {
		if pin.Equal(p) {
			return
		}
	}
	*pms = append(*pms, p)
}

func (pms *Permissions) RemovePermission(p Permission) {
	for idx, pin := range *pms {
		if pin.Equal(p) {
			pms.RemoveElement(idx)
			return
		}
	}
}

func (pms Permissions) HasPermission(p Permission) bool {
	for _, pin := range pms {
		if pin.Equal(p) {
			return true
		}
	}
	return false
}
func (pms Permissions) String() string {
	return fmt.Sprintf("%#v", pms)
}

type AccountPermissionI interface {
	GetAddress() sdk.AccAddress
	HasPermission(Permission) bool
	AddPermission(Permission)
	RemovePermission(Permission)
	String() string
	GetPermissions() Permissions
}

type AccountPermission struct {
	Address     sdk.AccAddress
	Permissions Permissions
}

func NewAccountPermission(addr sdk.AccAddress) AccountPermissionI {
	return &AccountPermission{
		Address: addr,
	}
}

func (ap *AccountPermission) String() string {
	return fmt.Sprintf("%#v", ap)
}

func (ap *AccountPermission) GetPermissions() Permissions {
	return ap.Permissions.GetPermissions()
}

func (ap *AccountPermission) GetAddress() sdk.AccAddress {
	return ap.Address
}

func (ap *AccountPermission) HasPermission(p Permission) bool {
	return ap.Permissions.HasPermission(p)
}
func (ap *AccountPermission) AddPermission(p Permission) {
	ap.Permissions.AddPermission(p)
}
func (ap *AccountPermission) RemovePermission(p Permission) {
	ap.Permissions.RemovePermission(p)
}
