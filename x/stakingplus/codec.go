package stakingplus

import (
	"github.com/line/lbm-sdk/codec/types"
	"github.com/line/lbm-sdk/x/foundation"
)

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*foundation.Authorization)(nil),
		&CreateValidatorAuthorization{},
	)
}
