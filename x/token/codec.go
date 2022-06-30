package token

import (
	"github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/msgservice"
)

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSend{},
		&MsgOperatorSend{},
		&MsgAuthorizeOperator{},
		&MsgRevokeOperator{},
		&MsgIssue{},
		&MsgGrant{},
		&MsgAbandon{},
		&MsgMint{},
		&MsgBurn{},
		&MsgOperatorBurn{},
		&MsgModify{},
		// deprecated messages
		&MsgTransferFrom{},
		&MsgApprove{},
		&MsgBurnFrom{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
