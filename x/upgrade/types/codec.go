package types

import (
	"github.com/line/lfb-sdk/codec"
	"github.com/line/lfb-sdk/codec/types"
	govtypes "github.com/line/lfb-sdk/x/gov/types"
)

// RegisterLegacyAminoCodec registers concrete types on the LegacyAmino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(Plan{}, "lfb-sdk/Plan", nil)
	cdc.RegisterConcrete(&SoftwareUpgradeProposal{}, "lfb-sdk/SoftwareUpgradeProposal", nil)
	cdc.RegisterConcrete(&CancelSoftwareUpgradeProposal{}, "lfb-sdk/CancelSoftwareUpgradeProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&SoftwareUpgradeProposal{},
		&CancelSoftwareUpgradeProposal{},
	)
}
