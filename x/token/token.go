package token

import (
	"strings"
)

const (
	prefixLegacyPermission = "LEGACY_PERMISSION_"
)

func (x LegacyPermission) String() string {
	lenPrefix := len(prefixLegacyPermission)
	return strings.ToLower(LegacyPermission_name[int32(x)][lenPrefix:])
}

func LegacyPermissionFromString(name string) LegacyPermission {
	legacyPermissionName := prefixLegacyPermission + strings.ToUpper(name)
	return LegacyPermission(LegacyPermission_value[legacyPermissionName])
}
