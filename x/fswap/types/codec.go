package types

import (
	"github.com/Finschia/finschia-sdk/codec"
	"github.com/Finschia/finschia-sdk/codec/legacy"
	"github.com/Finschia/finschia-sdk/codec/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/msgservice"
	fscodec "github.com/Finschia/finschia-sdk/x/fswap/codec"
	govcodec "github.com/Finschia/finschia-sdk/x/gov/codec"
	govtypes "github.com/Finschia/finschia-sdk/x/gov/types"
)

// RegisterLegacyAminoCodec registers concrete types on the LegacyAmino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgSwap{}, "lbm-sdk/MsgSwap")
	legacy.RegisterAminoMsg(cdc, &MsgSwapAll{}, "lbm-sdk/MsgSwapAll")

	cdc.RegisterConcrete(&MakeSwapProposal{}, "lbm-sdk/MakeSwapProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSwap{},
		&MsgSwapAll{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)

	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&MakeSwapProposal{},
	)
}

func init() {
	RegisterLegacyAminoCodec(govcodec.Amino)
	RegisterLegacyAminoCodec(fscodec.Amino)
}
