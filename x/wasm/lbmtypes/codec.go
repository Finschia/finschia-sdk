package lbmtypes

import (
	"github.com/line/lbm-sdk/codec"
	"github.com/line/lbm-sdk/codec/types"
	cryptocodec "github.com/line/lbm-sdk/crypto/codec"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/msgservice"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
	wasmTypes "github.com/line/lbm-sdk/x/wasm/types"
)

// RegisterLegacyAminoCodec registers the account types and interface
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) { //nolint:staticcheck
	wasmTypes.RegisterLegacyAminoCodec(cdc)
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
