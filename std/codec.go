package std

import (
	"github.com/Finschia/finschia-rdk/codec"
	"github.com/Finschia/finschia-rdk/codec/types"
	cryptocodec "github.com/Finschia/finschia-rdk/crypto/codec"
	sdk "github.com/Finschia/finschia-rdk/types"
	txtypes "github.com/Finschia/finschia-rdk/types/tx"
)

// RegisterLegacyAminoCodec registers types with the Amino codec.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	sdk.RegisterLegacyAminoCodec(cdc)
	cryptocodec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)
}

// RegisterInterfaces registers Interfaces from sdk/types, vesting, crypto, tx.
func RegisterInterfaces(interfaceRegistry types.InterfaceRegistry) {
	sdk.RegisterInterfaces(interfaceRegistry)
	txtypes.RegisterInterfaces(interfaceRegistry)
	cryptocodec.RegisterInterfaces(interfaceRegistry)
}
