package foundation

import (
	"github.com/line/lbm-sdk/codec"
	"github.com/line/lbm-sdk/codec/legacy"
	"github.com/line/lbm-sdk/codec/types"
	cryptocodec "github.com/line/lbm-sdk/crypto/codec"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/msgservice"
	authzcodec "github.com/line/lbm-sdk/x/authz/codec"
	govcodec "github.com/line/lbm-sdk/x/gov/codec"
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
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "lbm-sdk/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawFromTreasury{}, "lbm-sdk/MsgWithdrawFromTreasury")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateMembers{}, "lbm-sdk/MsgUpdateMembers")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateDecisionPolicy{}, "lbm-sdk/MsgUpdateDecisionPolicy")
	legacy.RegisterAminoMsg(cdc, &MsgGrant{}, "lbm-sdk/MsgGrant")
	legacy.RegisterAminoMsg(cdc, &MsgRevoke{}, "lbm-sdk/MsgRevoke")
	legacy.RegisterAminoMsg(cdc, &MsgGovMint{}, "lbm-sdk/MsgGovMint")

	cdc.RegisterInterface((*Authorization)(nil), nil)
	cdc.RegisterInterface((*DecisionPolicy)(nil), nil)
	cdc.RegisterConcrete(&ThresholdDecisionPolicy{}, "lbm-sdk/ThresholdDecisionPolicy", nil)
	cdc.RegisterConcrete(&PercentageDecisionPolicy{}, "lbm-sdk/PercentageDecisionPolicy", nil)
	cdc.RegisterConcrete(&ReceiveFromTreasuryAuthorization{}, "lbm-sdk/ReceiveFromTreasuryAuthorization", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
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

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(Amino)
)

func init() {
	RegisterLegacyAminoCodec(Amino)
	cryptocodec.RegisterCrypto(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)

	// Register all Amino interfaces and concrete types on the authz and gov Amino codec so that this can later be
	// used to properly serialize MsgGrant, MsgExec and MsgSubmitProposal instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
	RegisterLegacyAminoCodec(govcodec.Amino)
}
