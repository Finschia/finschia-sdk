package lbmtypes

import (
	"github.com/line/lbm-sdk/codec"
	"github.com/line/lbm-sdk/codec/legacy"
	"github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/msgservice"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
	wasmTypes "github.com/line/lbm-sdk/x/wasm/types"
)

// RegisterLegacyAminoCodec registers the account types and interface
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) { //nolint:staticcheck
	cdc.RegisterConcrete(&MsgStoreCodeAndInstantiateContract{}, "wasm/StoreCodeAndInstantiateContract", nil)

	cdc.RegisterConcrete(&DeactivateContractProposal{}, "wasm/DeactivateContractProposal", nil)
	cdc.RegisterConcrete(&ActivateContractProposal{}, "wasm/ActivateContractProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	wasmTypes.RegisterInterfaces(registry)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgStoreCodeAndInstantiateContract{},
	)
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&DeactivateContractProposal{},
		&ActivateContractProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func init() {
	RegisterLegacyAminoCodec(legacy.Cdc)
}
