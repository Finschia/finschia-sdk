package types

import (
	"github.com/line/lfb-sdk/codec"
	"github.com/line/lfb-sdk/codec/types"
	cryptocodec "github.com/line/lfb-sdk/crypto/codec"
	sdk "github.com/line/lfb-sdk/types"
	"github.com/line/lfb-sdk/types/msgservice"
	govtypes "github.com/line/lfb-sdk/x/gov/types"
)

// RegisterLegacyAminoCodec registers the account types and interface
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgStoreCode{}, "wasm/MsgStoreCode", nil)
	cdc.RegisterConcrete(&MsgInstantiateContract{}, "wasm/MsgInstantiateContract", nil)
	cdc.RegisterConcrete(&MsgExecuteContract{}, "wasm/MsgExecuteContract", nil)
	cdc.RegisterConcrete(&MsgMigrateContract{}, "wasm/MsgMigrateContract", nil)
	cdc.RegisterConcrete(&MsgUpdateAdmin{}, "wasm/MsgUpdateAdmin", nil)
	cdc.RegisterConcrete(&MsgClearAdmin{}, "wasm/MsgClearAdmin", nil)
	cdc.RegisterConcrete(&PinCodesProposal{}, "wasm/PinCodesProposal", nil)
	cdc.RegisterConcrete(&UnpinCodesProposal{}, "wasm/UnpinCodesProposal", nil)

	cdc.RegisterConcrete(&StoreCodeProposal{}, "wasm/StoreCodeProposal", nil)
	cdc.RegisterConcrete(&InstantiateContractProposal{}, "wasm/InstantiateContractProposal", nil)
	cdc.RegisterConcrete(&MigrateContractProposal{}, "wasm/MigrateContractProposal", nil)
	cdc.RegisterConcrete(&UpdateAdminProposal{}, "wasm/UpdateAdminProposal", nil)
	cdc.RegisterConcrete(&ClearAdminProposal{}, "wasm/ClearAdminProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgStoreCode{},
		&MsgInstantiateContract{},
		&MsgExecuteContract{},
		&MsgMigrateContract{},
		&MsgUpdateAdmin{},
		&MsgClearAdmin{},
		&MsgIBCCloseChannel{},
		&MsgIBCSend{},
	)
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&StoreCodeProposal{},
		&InstantiateContractProposal{},
		&MigrateContractProposal{},
		&UpdateAdminProposal{},
		&ClearAdminProposal{},
		&PinCodesProposal{},
		&UnpinCodesProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/wasm module codec.

	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
