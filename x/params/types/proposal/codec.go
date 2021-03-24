package proposal

import (
	"github.com/line/lbm-sdk/v2/codec"
	"github.com/line/lbm-sdk/v2/codec/types"
	govtypes "github.com/line/lbm-sdk/v2/x/gov/types"
)

// RegisterLegacyAminoCodec registers all necessary param module types with a given LegacyAmino codec.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&ParameterChangeProposal{}, "lbm-sdk/ParameterChangeProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&ParameterChangeProposal{},
	)
}
