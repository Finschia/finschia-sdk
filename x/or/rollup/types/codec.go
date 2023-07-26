package types

import (
	"github.com/Finschia/finschia-sdk/codec"
	cdctypes "github.com/Finschia/finschia-sdk/codec/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateRollup{}, "rollup/CreateRollup", nil)
	cdc.RegisterConcrete(&MsgRegisterSequencer{}, "rollup/RegisterSequencer", nil)
	cdc.RegisterConcrete(&MsgRemoveSequencer{}, "rollup/RemoveSequencer", nil)
	cdc.RegisterConcrete(&MsgDeposit{}, "rollup/Deposit", nil)
	cdc.RegisterConcrete(&MsgWithdraw{}, "rollup/Withdraw", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateRollup{},
		&MsgRegisterSequencer{},
		&MsgRemoveSequencer{},
		&MsgDeposit{},
		&MsgWithdraw{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
