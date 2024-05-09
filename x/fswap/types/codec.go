package types

import (
	"github.com/Finschia/finschia-sdk/codec"
	"github.com/Finschia/finschia-sdk/codec/legacy"
	"github.com/Finschia/finschia-sdk/codec/types"
	cryptocodec "github.com/Finschia/finschia-sdk/crypto/codec"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/msgservice"
	fdncodec "github.com/Finschia/finschia-sdk/x/foundation/codec"
)

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(Amino)
)

func init() {
	cryptocodec.RegisterCrypto(Amino)
	codec.RegisterEvidences(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)

	RegisterLegacyAminoCodec(Amino)
	RegisterLegacyAminoCodec(fdncodec.Amino)
}

// RegisterLegacyAminoCodec registers concrete types on the LegacyAmino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgSwap{}, "lbm-sdk/MsgSwap")
	legacy.RegisterAminoMsg(cdc, &MsgSwapAll{}, "lbm-sdk/MsgSwapAll")
	legacy.RegisterAminoMsg(cdc, &MsgMakeSwapProposal{}, "lbm-sdk/MsgMakeSwapProposal")
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSwap{},
		&MsgSwapAll{},
		&MsgMakeSwapProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
