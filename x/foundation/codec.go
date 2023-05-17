package foundation

import (
	"github.com/Finschia/finschia-sdk/codec"
	"github.com/Finschia/finschia-sdk/codec/legacy"
	"github.com/Finschia/finschia-sdk/codec/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/msgservice"
	authzcodec "github.com/Finschia/finschia-sdk/x/authz/codec"
	fdncodec "github.com/Finschia/finschia-sdk/x/foundation/codec"
	govcodec "github.com/Finschia/finschia-sdk/x/gov/codec"
	govtypes "github.com/Finschia/finschia-sdk/x/gov/types"
)

// RegisterLegacyAminoCodec registers concrete types on the LegacyAmino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgFundTreasury{}, "lbm-sdk/MsgFundTreasury")
	legacy.RegisterAminoMsg(cdc, &MsgSubmitProposal{}, "lbm-sdk/MsgSubmitProposal")
	legacy.RegisterAminoMsg(cdc, &MsgVote{}, "lbm-sdk/MsgVote")
	legacy.RegisterAminoMsg(cdc, &MsgExec{}, "lbm-sdk/MsgExec")
	legacy.RegisterAminoMsg(cdc, &MsgLeaveFoundation{}, "lbm-sdk/MsgLeaveFoundation")
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawProposal{}, "lbm-sdk/MsgWithdrawProposal")

	// proposal from foundation operator
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawFromTreasury{}, "lbm-sdk/MsgWithdrawFromTreasury")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateMembers{}, "lbm-sdk/MsgUpdateMembers")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateDecisionPolicy{}, "lbm-sdk/MsgUpdateDecisionPolicy")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateCensorship{}, "lbm-sdk/MsgUpdateCensorship")
	legacy.RegisterAminoMsg(cdc, &MsgGrant{}, "lbm-sdk/MsgGrant")
	legacy.RegisterAminoMsg(cdc, &MsgRevoke{}, "lbm-sdk/MsgRevoke")

	cdc.RegisterInterface((*Authorization)(nil), nil)
	cdc.RegisterInterface((*DecisionPolicy)(nil), nil)
	cdc.RegisterConcrete(&ThresholdDecisionPolicy{}, "lbm-sdk/ThresholdDecisionPolicy", nil)
	cdc.RegisterConcrete(&PercentageDecisionPolicy{}, "lbm-sdk/PercentageDecisionPolicy", nil)
	cdc.RegisterConcrete(&ReceiveFromTreasuryAuthorization{}, "lbm-sdk/ReceiveFromTreasuryAuthorization", nil)

	cdc.RegisterConcrete(&FoundationExecProposal{}, "lbm-sdk/FoundationExecProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
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
		&MsgUpdateCensorship{},
		&MsgGrant{},
		&MsgRevoke{},
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

	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&FoundationExecProposal{},
	)
}

func init() {
	// Register all Amino interfaces and concrete types on the authz and gov Amino codec so that this can later be
	// used to properly serialize MsgGrant, MsgExec and MsgSubmitProposal instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
	RegisterLegacyAminoCodec(govcodec.Amino)
	RegisterLegacyAminoCodec(fdncodec.Amino)
}
