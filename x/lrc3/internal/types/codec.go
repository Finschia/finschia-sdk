package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	nft "github.com/link-chain/link/x/nft"
	"github.com/link-chain/link/x/nft/exported"
)

// RegisterCodec concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.NFT)(nil), nil)
	cdc.RegisterConcrete(&nft.BaseNFT{}, "lrc3/BaseNFT", nil)

	cdc.RegisterConcrete(MsgInit{}, "lrc3/MsgInit", nil)
	cdc.RegisterConcrete(MsgMintNFT{}, "lrc3/MsgMintNFT", nil)
	cdc.RegisterConcrete(MsgBurn{}, "lrc3/MsgBurn", nil)
	cdc.RegisterConcrete(MsgTransfer{}, "lrc3/MsgTransfer", nil)
	cdc.RegisterConcrete(MsgEditMetadata{}, "lrc3/MsgEditMetadata", nil)
	cdc.RegisterConcrete(MsgApprove{}, "lrc3/MsgApprove", nil)
	cdc.RegisterConcrete(MsgSetApprovalForAll{}, "lrc3/MsgSetApprovalForAll", nil)
}

// ModuleCdc generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	codec.RegisterCrypto(ModuleCdc)
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
