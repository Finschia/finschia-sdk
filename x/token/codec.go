package token

import (
	"github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/msgservice"
)

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSend{},
		&MsgRevokeOperator{},
		&MsgIssue{},
		&MsgMint{},
		&MsgBurn{},
		&MsgModify{},
		&MsgTransferFrom{},
		&MsgApprove{},
		&MsgBurnFrom{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
