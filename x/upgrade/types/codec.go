package types

import (
	"github.com/Finschia/finschia-sdk/codec"
	"github.com/Finschia/finschia-sdk/codec/types"
	cryptocodec "github.com/Finschia/finschia-sdk/crypto/codec"
	sdk "github.com/Finschia/finschia-sdk/types"
	authzcodec "github.com/Finschia/finschia-sdk/x/authz/codec"
	fdncodec "github.com/Finschia/finschia-sdk/x/foundation/codec"
	govcodec "github.com/Finschia/finschia-sdk/x/gov/codec"
	govtypes "github.com/Finschia/finschia-sdk/x/gov/types"
)

// RegisterLegacyAminoCodec registers concrete types on the LegacyAmino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(Plan{}, "cosmos-sdk/Plan", nil)
	cdc.RegisterConcrete(&SoftwareUpgradeProposal{}, "cosmos-sdk/SoftwareUpgradeProposal", nil)
	cdc.RegisterConcrete(&CancelSoftwareUpgradeProposal{}, "cosmos-sdk/CancelSoftwareUpgradeProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&SoftwareUpgradeProposal{},
		&CancelSoftwareUpgradeProposal{},
	)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)

	// Register all Amino interfaces and concrete types on the authz  and gov Amino codec
	// so that this can later be used to properly serialize MsgGrant and MsgExec
	// instances.
	RegisterLegacyAminoCodec(authzcodec.Amino)
	RegisterLegacyAminoCodec(govcodec.Amino)
	RegisterLegacyAminoCodec(fdncodec.Amino)
}
