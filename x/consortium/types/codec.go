package types

import (
	"github.com/line/lbm-sdk/codec"
	"github.com/line/lbm-sdk/codec/types"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

// RegisterLegacyAminoCodec registers concrete types on the LegacyAmino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&DisableConsortiumProposal{}, "lbm-sdk/DisableConsortiumProposal", nil)
	cdc.RegisterConcrete(&EditAllowedValidatorsProposal{}, "lbm-sdk/EditAllowedValidatorsProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&DisableConsortiumProposal{},
		&EditAllowedValidatorsProposal{},
	)
}
