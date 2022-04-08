package foundation

import (
	"github.com/line/lbm-sdk/codec/types"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&UpdateFoundationParamsProposal{},
		&UpdateValidatorAuthsProposal{},
	)
}
