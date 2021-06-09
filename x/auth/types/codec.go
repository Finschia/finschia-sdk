package types

import (
	"github.com/line/lfb-sdk/codec"
	"github.com/line/lfb-sdk/codec/types"
	cryptocodec "github.com/line/lfb-sdk/crypto/codec"
	sdk "github.com/line/lfb-sdk/types"
	"github.com/line/lfb-sdk/types/msgservice"
	"github.com/line/lfb-sdk/x/auth/legacy/legacytx"
)

// RegisterLegacyAminoCodec registers the account interfaces and concrete types on the
// provided LegacyAmino codec. These types are used for Amino JSON serialization
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*ModuleAccountI)(nil), nil)
	cdc.RegisterInterface((*GenesisAccount)(nil), nil)
	cdc.RegisterInterface((*AccountI)(nil), nil)
	cdc.RegisterConcrete(&BaseAccount{}, "lfb-sdk/BaseAccount", nil)
	cdc.RegisterConcrete(&ModuleAccount{}, "lfb-sdk/ModuleAccount", nil)
	cdc.RegisterConcrete(&MsgEmpty{}, "lfb-sdk/MsgEmpty", nil)

	legacytx.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces associates protoName with AccountI interface
// and creates a registry of it's concrete implementations
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgEmpty{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)

	registry.RegisterInterface(
		"lfb.auth.v1beta1.AccountI",
		(*AccountI)(nil),
		&BaseAccount{},
		&ModuleAccount{},
	)

	registry.RegisterInterface(
		"lfb.auth.v1beta1.GenesisAccount",
		(*GenesisAccount)(nil),
		&BaseAccount{},
		&ModuleAccount{},
	)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
}
