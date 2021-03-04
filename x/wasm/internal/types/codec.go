package types

import (
	"github.com/line/lbm-sdk/codec"
)

// RegisterCodec registers the account types and interface
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgStoreCode{}, "wasm/MsgStoreCode", nil)
	cdc.RegisterConcrete(MsgInstantiateContract{}, "wasm/MsgInstantiateContract", nil)
	cdc.RegisterConcrete(MsgExecuteContract{}, "wasm/MsgExecuteContract", nil)
	cdc.RegisterConcrete(MsgMigrateContract{}, "wasm/MsgMigrateContract", nil)
	cdc.RegisterConcrete(MsgUpdateAdmin{}, "wasm/MsgUpdateAdmin", nil)
	cdc.RegisterConcrete(MsgClearAdmin{}, "wasm/MsgClearAdmin", nil)

	cdc.RegisterConcrete(StoreCodeProposal{}, "wasm/StoreCodeProposal", nil)
	cdc.RegisterConcrete(InstantiateContractProposal{}, "wasm/InstantiateContractProposal", nil)
	cdc.RegisterConcrete(MigrateContractProposal{}, "wasm/MigrateContractProposal", nil)
	cdc.RegisterConcrete(UpdateAdminProposal{}, "wasm/UpdateAdminProposal", nil)
	cdc.RegisterConcrete(ClearAdminProposal{}, "wasm/ClearAdminProposal", nil)

	// query responses

	// For the type-tags in case of a slice item or a nested property.
	cdc.RegisterInterface((*CodeInfoResponse)(nil), nil)
	cdc.RegisterInterface((*ContractInfoResponse)(nil), nil)
	cdc.RegisterInterface((*ContractHistoryResponse)(nil), nil)

	cdc.RegisterConcrete(codeInfo{}, "wasm/CodeInfo", nil)
	cdc.RegisterConcrete(contractInfo{}, "wasm/ContractInfo", nil)
	cdc.RegisterConcrete(contractHistory{}, "wasm/ContractHistory", nil)
}

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	ModuleCdc = cdc.Seal()
}
