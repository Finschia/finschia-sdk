package types

import (
	"github.com/line/lbm-sdk/codec"
	"github.com/line/lbm-sdk/codec/types"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

// RegisterLegacyAminoCodec registers concrete types on the LegacyAmino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(Plan{}, "lbm-sdk/Plan", nil)
	cdc.RegisterConcrete(&SoftwareUpgradeProposal{}, "lbm-sdk/SoftwareUpgradeProposal", nil)
	cdc.RegisterConcrete(&CancelSoftwareUpgradeProposal{}, "lbm-sdk/CancelSoftwareUpgradeProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&SoftwareUpgradeProposal{},
		&CancelSoftwareUpgradeProposal{},
	)
}
