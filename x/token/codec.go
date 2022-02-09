package token

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/codec/types"
	"github.com/line/lbm-sdk/types/msgservice"
)

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgTransfer{},
		&MsgTransferFrom{},
		&MsgApprove{},
		&MsgIssue{},
		&MsgGrant{},
		&MsgRevoke{},
		&MsgMint{},
		&MsgBurn{},
		&MsgBurnFrom{},
		&MsgModify{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
