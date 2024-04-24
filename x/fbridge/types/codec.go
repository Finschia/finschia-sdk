package types

import (
	"github.com/Finschia/finschia-sdk/codec/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/msgservice"
)

func RegisterInterfaces(registrar types.InterfaceRegistry) {
	registrar.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgTransfer{},
		&MsgProvision{},
		&MsgHoldTransfer{},
		&MsgReleaseTransfer{},
		&MsgRemoveProvision{},
		&MsgClaimBatch{},
		&MsgClaim{},
		&MsgUpdateRole{},
	)

	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
