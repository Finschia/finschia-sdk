package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/iam/exported"
)

type AccountPermissionI interface {
	GetAddress() sdk.AccAddress
	HasPermission(exported.PermissionI) bool
	AddPermission(exported.PermissionI)
	RemovePermission(exported.PermissionI)
	String() string
	GetPermissions() []exported.PermissionI
	InheritAccountPermission(AccountPermissionI) InheritedAccountPermissionI
}

type InheritedAccountPermissionI interface {
	AccountPermissionI
	SetParent(parent AccountPermissionI)
}

type Permissions []exported.PermissionI

func NewPermissions(perms ...exported.PermissionI) Permissions {
	pms := Permissions{}
	for _, perm := range perms {
		pms.AddPermission(perm)
	}
	return pms
}

func (pms *Permissions) GetPermissions() []exported.PermissionI {
	return []exported.PermissionI(*pms)
}

func (pms *Permissions) RemoveElement(idx int) {
	*pms = append((*pms)[:idx], (*pms)[idx+1:]...)
}

func (pms *Permissions) AddPermission(p exported.PermissionI) {
	for _, pin := range *pms {
		if pin.Equal(p.GetResource(), p.GetAction()) {
			return
		}
	}
	*pms = append(*pms, p)
}

func (pms *Permissions) RemovePermission(p exported.PermissionI) {
	for idx, pin := range *pms {
		if pin.Equal(p.GetResource(), p.GetAction()) {
			pms.RemoveElement(idx)
			return
		}
	}
}

func (pms Permissions) HasPermission(p exported.PermissionI) bool {
	for _, pin := range pms {
		if pin.Equal(p.GetResource(), p.GetAction()) {
			return true
		}
	}
	return false
}
func (pms Permissions) String() string {
	return fmt.Sprintf("%#v", pms)
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

func (ap *AccountPermission) GetPermissions() []exported.PermissionI {
	return ap.Permissions.GetPermissions()
}

func (ap *AccountPermission) GetAddress() sdk.AccAddress {
	return ap.Address
}

func (ap *AccountPermission) HasPermission(p exported.PermissionI) bool {
	return ap.Permissions.HasPermission(p)
}
func (ap *AccountPermission) AddPermission(p exported.PermissionI) {
	ap.Permissions.AddPermission(p)
}
func (ap *AccountPermission) RemovePermission(p exported.PermissionI) {
	ap.Permissions.RemovePermission(p)
}

func (ap *AccountPermission) InheritAccountPermission(parent AccountPermissionI) InheritedAccountPermissionI {
	return NewInheritedAccountPermission(ap, parent)
}

type InheritedAccountPermission struct {
	AccountPermission
	ParentAddr sdk.AccAddress
	parent     AccountPermissionI
}

func NewInheritedAccountPermission(self, parent AccountPermissionI) InheritedAccountPermissionI {
	return &InheritedAccountPermission{
		AccountPermission: *self.(*AccountPermission),
		ParentAddr:        parent.GetAddress(),
		parent:            parent,
	}
}

func (iap *InheritedAccountPermission) HasPermission(p exported.PermissionI) bool {
	has := iap.AccountPermission.HasPermission(p)
	if !has {
		has = iap.parent.HasPermission(p)
	}
	return has
}

func (iap *InheritedAccountPermission) SetParent(parent AccountPermissionI) {
	iap.parent = parent
	iap.ParentAddr = parent.GetAddress()
}

func (iap *InheritedAccountPermission) GetPermissions() []exported.PermissionI {
	pms := iap.AccountPermission.Permissions.GetPermissions()
	pms = append(pms, iap.parent.GetPermissions()...)
	return pms
}
