package types

import (
	"github.com/Finschia/finschia-sdk/codec"
	"github.com/Finschia/finschia-sdk/codec/legacy"
	"github.com/Finschia/finschia-sdk/codec/types"
	cryptocodec "github.com/Finschia/finschia-sdk/crypto/codec"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/msgservice"
	authzcodec "github.com/Finschia/finschia-sdk/x/authz/codec"
	fdncodec "github.com/Finschia/finschia-sdk/x/foundation/codec"
	govcodec "github.com/Finschia/finschia-sdk/x/gov/codec"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "lbm-sdk/fbridge/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgTransfer{}, "lbm-sdk/MsgTransfer")
	legacy.RegisterAminoMsg(cdc, &MsgProvision{}, "lbm-sdk/MsgProvision")
	legacy.RegisterAminoMsg(cdc, &MsgHoldTransfer{}, "lbm-sdk/MsgHoldTransfer")
	legacy.RegisterAminoMsg(cdc, &MsgReleaseTransfer{}, "lbm-sdk/MsgReleaseTransfer")
	legacy.RegisterAminoMsg(cdc, &MsgRemoveProvision{}, "lbm-sdk/MsgRemoveProvision")
	legacy.RegisterAminoMsg(cdc, &MsgClaimBatch{}, "lbm-sdk/MsgClaimBatch")
	legacy.RegisterAminoMsg(cdc, &MsgClaim{}, "lbm-sdk/MsgClaim")
	legacy.RegisterAminoMsg(cdc, &MsgSuggestRole{}, "lbm-sdk/MsgSuggestRole")
	legacy.RegisterAminoMsg(cdc, &MsgAddVoteForRole{}, "lbm-sdk/MsgAddVoteForRole")
	legacy.RegisterAminoMsg(cdc, &MsgSetBridgeStatus{}, "lbm-sdk/MsgSetBridgeStatus")
}

func RegisterInterfaces(registrar types.InterfaceRegistry) {
	registrar.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgUpdateParams{},
		&MsgTransfer{},
		&MsgProvision{},
		&MsgHoldTransfer{},
		&MsgReleaseTransfer{},
		&MsgRemoveProvision{},
		&MsgClaimBatch{},
		&MsgClaim{},
		&MsgSuggestRole{},
		&MsgAddVoteForRole{},
		&MsgSetBridgeStatus{},
	)

	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(Amino)
)

func init() {
	cryptocodec.RegisterCrypto(Amino)
	codec.RegisterEvidences(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)

	// Register all Amino interfaces and concrete types on the authz and gov Amino codec so that this can later be
	// used to properly serialize MsgGrant, MsgExec and MsgSubmitProposal instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
	RegisterLegacyAminoCodec(govcodec.Amino)
	RegisterLegacyAminoCodec(fdncodec.Amino)
}
