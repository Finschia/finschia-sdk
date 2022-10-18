package foundation

import (
	"github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/msgservice"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*govtypes.Content)(nil),
		&UpdateFoundationParamsProposal{},
	)

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgFundTreasury{},
		&MsgWithdrawFromTreasury{},
		&MsgUpdateMembers{},
		&MsgUpdateDecisionPolicy{},
		&MsgSubmitProposal{},
		&MsgWithdrawProposal{},
		&MsgVote{},
		&MsgExec{},
		&MsgLeaveFoundation{},
		&MsgGrant{},
		&MsgRevoke{},
		&MsgGovMint{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)

	registry.RegisterInterface(
		"lbm.foundation.v1.DecisionPolicy",
		(*DecisionPolicy)(nil),
		&ThresholdDecisionPolicy{},
		&PercentageDecisionPolicy{},
		&OutsourcingDecisionPolicy{},
	)

	registry.RegisterImplementations(
		(*Authorization)(nil),
		&ReceiveFromTreasuryAuthorization{},
	)
}
