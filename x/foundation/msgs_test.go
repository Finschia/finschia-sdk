package foundation_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

// func TestMsgUpdateParams(t *testing.T) {
// 	addressCodec := addresscodec.NewBech32Codec("link")
// 	bytesToString := func(addr sdk.AccAddress) string {
// 		str, err := addressCodec.BytesToString(addr)
// 		require.NoError(t, err)
// 		return str
// 	}

// 	addrs := make([]sdk.AccAddress, 1)
// 	for i := range addrs {
// 		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
// 	}

// 	testCases := map[string]struct {
// 		authority sdk.AccAddress
// 		params    foundation.Params
// 		valid     bool
// 	}{
// 		"handler for MsgUpdateParams removed, ValidateBasic should throw error always": {
// 			authority: addrs[0],
// 			params: foundation.Params{
// 				FoundationTax: math.LegacyZeroDec(),
// 			},
// 			valid: false,
// 		},
// 	}

// 	for name, tc := range testCases {
// 		t.Run(name, func(t *testing.T) {
// 			msg := foundation.MsgUpdateParams{
// 				Authority: bytesToString(tc.authority),
// 				Params:    tc.params,
// 			}

// 			err := msg.ValidateBasic()
// 			require.Error(t, err)
// 			require.ErrorIs(t, err, sdkerrors.ErrUnknownRequest)
// 		})
// 		msg := foundation.MsgUpdateParams{
// 			addrs[0].String(),
// 			foundation.Params{},
// 		}
// 		// Note: Dummy test for coverage of deprecated message
// 		_ = msg.String()
// 		_, _ = msg.Descriptor()
// 		_, _ = msg.Marshal()
// 		msg.ProtoMessage()
// 		msg.Reset()
// 		_ = msg.Size()
// 		_ = msg.XXX_Size()
// 	}
// }

func TestAminoJSON(t *testing.T) {
	legacyAmino := codec.NewLegacyAmino()
	foundation.RegisterLegacyAminoCodec(legacyAmino)
	legacytx.RegressionTestingAminoCodec = legacyAmino
	
	addressCodec := addresscodec.NewBech32Codec("link")
	bytesToString := func(addr sdk.AccAddress) string {
		str, err := addressCodec.BytesToString(addr)
		require.NoError(t, err)
		return str
	}

	addrs := make([]sdk.AccAddress, 3)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		msg      legacytx.LegacyMsg
		expected string
	}{
		"MsgFundTreasury": {
			&foundation.MsgFundTreasury{
				From:   bytesToString(addrs[0]),
				Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.OneInt())),
			},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgFundTreasury\",\"value\":{\"amount\":[{\"amount\":\"1\",\"denom\":\"stake\"}],\"from\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", bytesToString(addrs[0])),
		},
		"MsgVote": {
			&foundation.MsgVote{
				ProposalId: 1,
				Voter:      bytesToString(addrs[0]),
				Option:     foundation.VOTE_OPTION_YES,
				Metadata:   "I'm YES",
				Exec:       foundation.Exec_EXEC_UNSPECIFIED,
			},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgVote\",\"value\":{\"metadata\":\"I'm YES\",\"option\":1,\"proposal_id\":\"1\",\"voter\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", bytesToString(addrs[0])),
		},
		"MsgExec": {
			&foundation.MsgExec{
				ProposalId: 1,
				Signer:     bytesToString(addrs[0]),
			},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgExec\",\"value\":{\"proposal_id\":\"1\",\"signer\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", bytesToString(addrs[0])),
		},
		"MsgLeaveFoundation": {
			&foundation.MsgLeaveFoundation{Address: bytesToString(addrs[0])},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgLeaveFoundation\",\"value\":{\"address\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", bytesToString(addrs[0])),
		},
		"MsgWithdrawProposal": {
			&foundation.MsgWithdrawProposal{
				ProposalId: 1,
				Address:    bytesToString(addrs[0]),
			},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgWithdrawProposal\",\"value\":{\"address\":\"%s\",\"proposal_id\":\"1\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", bytesToString(addrs[0])),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expected, string(legacytx.StdSignBytes("foo", 1, 1, 1, legacytx.StdFee{}, []sdk.Msg{tc.msg}, "memo")))
		})
	}
}

func TestMsgSubmitProposalAminoJSON(t *testing.T) {
	addressCodec := addresscodec.NewBech32Codec("link")
	bytesToString := func(addr sdk.AccAddress) string {
		str, err := addressCodec.BytesToString(addr)
		require.NoError(t, err)
		return str
	}

	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	proposer := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	testCases := map[string]struct {
		msg      sdk.Msg
		expected string
	}{
		"MsgUpdateParams": {
			&foundation.MsgUpdateParams{
				Authority: bytesToString(addrs[0]),
				Params:    foundation.Params{FoundationTax: math.LegacyZeroDec()},
			},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSubmitProposal\",\"value\":{\"exec\":1,\"messages\":[{\"type\":\"lbm-sdk/MsgUpdateParams\",\"value\":{\"authority\":\"%s\",\"params\":{\"foundation_tax\":\"0.000000000000000000\"}}}],\"metadata\":\"MsgUpdateParams\",\"proposers\":[\"%s\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", bytesToString(addrs[0]), bytesToString(proposer)),
		},
		"MsgWithdrawFromTreasury": {
			&foundation.MsgWithdrawFromTreasury{
				Authority: bytesToString(addrs[0]),
				To:        bytesToString(addrs[1]),
				Amount:    sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1000000))),
			},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSubmitProposal\",\"value\":{\"exec\":1,\"messages\":[{\"type\":\"lbm-sdk/MsgWithdrawFromTreasury\",\"value\":{\"amount\":[{\"amount\":\"1000000\",\"denom\":\"stake\"}],\"authority\":\"%s\",\"to\":\"%s\"}}],\"metadata\":\"MsgWithdrawFromTreasury\",\"proposers\":[\"%s\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", bytesToString(addrs[0]), bytesToString(addrs[1]), bytesToString(proposer)),
		},
		"MsgUpdateMembers": {
			&foundation.MsgUpdateMembers{
				Authority: bytesToString(addrs[0]),
				MemberUpdates: []foundation.MemberRequest{{
					Address: bytesToString(addrs[1]),
				}},
			},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSubmitProposal\",\"value\":{\"exec\":1,\"messages\":[{\"type\":\"lbm-sdk/MsgUpdateMembers\",\"value\":{\"authority\":\"%s\",\"member_updates\":[{\"address\":\"%s\"}]}}],\"metadata\":\"MsgUpdateMembers\",\"proposers\":[\"%s\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", bytesToString(addrs[0]), bytesToString(addrs[1]), bytesToString(proposer)),
		},
		"MsgUpdateCensorship": {
			&foundation.MsgUpdateCensorship{
				Authority: bytesToString(addrs[0]),
				Censorship: foundation.Censorship{
					MsgTypeUrl: sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
					Authority:  foundation.CensorshipAuthorityGovernance,
				},
			},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSubmitProposal\",\"value\":{\"exec\":1,\"messages\":[{\"type\":\"lbm-sdk/MsgUpdateCensorship\",\"value\":{\"authority\":\"%s\",\"censorship\":{\"authority\":1,\"msg_type_url\":\"/lbm.foundation.v1.MsgWithdrawFromTreasury\"}}}],\"metadata\":\"MsgUpdateCensorship\",\"proposers\":[\"%s\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", bytesToString(addrs[0]), bytesToString(proposer)),
		},
		"MsgRevoke": {
			&foundation.MsgRevoke{
				Authority:  bytesToString(addrs[0]),
				Grantee:    bytesToString(addrs[1]),
				MsgTypeUrl: foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
			},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSubmitProposal\",\"value\":{\"exec\":1,\"messages\":[{\"type\":\"lbm-sdk/MsgRevoke\",\"value\":{\"authority\":\"%s\",\"grantee\":\"%s\",\"msg_type_url\":\"/lbm.foundation.v1.MsgWithdrawFromTreasury\"}}],\"metadata\":\"MsgRevoke\",\"proposers\":[\"%s\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", bytesToString(addrs[0]), bytesToString(addrs[1]), bytesToString(proposer)),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			proposalMsg := &foundation.MsgSubmitProposal{
				Proposers: []string{bytesToString(proposer)},
				Metadata:  name,
				Exec:      foundation.Exec_EXEC_TRY,
			}
			err := proposalMsg.SetMsgs([]sdk.Msg{tc.msg})
			require.NoError(t, err)
			require.Equal(t, tc.expected, string(legacytx.StdSignBytes("foo", 1, 1, 1, legacytx.StdFee{}, []sdk.Msg{proposalMsg}, "memo")))
		})
	}
}

func TestMsgUpdateDecisionPolicyAminoJson(t *testing.T) {
	var (
		authority = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
		proposer  = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	)

	addressCodec := addresscodec.NewBech32Codec("link")
	bytesToString := func(addr sdk.AccAddress) string {
		str, err := addressCodec.BytesToString(addr)
		require.NoError(t, err)
		return str
	}

	testCases := map[string]struct {
		policy   foundation.DecisionPolicy
		expected string
	}{
		"ThresholdDecisionPolicy": {
			&foundation.ThresholdDecisionPolicy{
				Threshold: math.LegacyOneDec(),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSubmitProposal\",\"value\":{\"exec\":1,\"messages\":[{\"type\":\"lbm-sdk/MsgUpdateDecisionPolicy\",\"value\":{\"authority\":\"%s\",\"decision_policy\":{\"type\":\"lbm-sdk/ThresholdDecisionPolicy\",\"value\":{\"threshold\":\"1.000000000000000000\",\"windows\":{\"min_execution_period\":\"0\",\"voting_period\":\"3600000000000\"}}}}}],\"metadata\":\"ThresholdDecisionPolicy\",\"proposers\":[\"%s\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", bytesToString(authority), bytesToString(proposer)),
		},
		"PercentageDecisionPolicy": {
			&foundation.PercentageDecisionPolicy{
				Percentage: math.LegacyOneDec(),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSubmitProposal\",\"value\":{\"exec\":1,\"messages\":[{\"type\":\"lbm-sdk/MsgUpdateDecisionPolicy\",\"value\":{\"authority\":\"%s\",\"decision_policy\":{\"type\":\"lbm-sdk/PercentageDecisionPolicy\",\"value\":{\"percentage\":\"1.000000000000000000\",\"windows\":{\"min_execution_period\":\"0\",\"voting_period\":\"3600000000000\"}}}}}],\"metadata\":\"PercentageDecisionPolicy\",\"proposers\":[\"%s\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", bytesToString(authority), bytesToString(proposer)),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			policyMsg := &foundation.MsgUpdateDecisionPolicy{
				Authority: bytesToString(authority),
			}
			err := policyMsg.SetDecisionPolicy(tc.policy)
			require.NoError(t, err)

			proposalMsg := &foundation.MsgSubmitProposal{
				Proposers: []string{bytesToString(proposer)},
				Metadata:  name,
				Exec:      foundation.Exec_EXEC_TRY,
			}
			err = proposalMsg.SetMsgs([]sdk.Msg{policyMsg})
			require.NoError(t, err)

			require.Equal(t, tc.expected, string(legacytx.StdSignBytes("foo", 1, 1, 1, legacytx.StdFee{}, []sdk.Msg{proposalMsg}, "memo")))
		})
	}
}

func TestMsgGrantAminoJson(t *testing.T) {
	addressCodec := addresscodec.NewBech32Codec("link")
	bytesToString := func(addr sdk.AccAddress) string {
		str, err := addressCodec.BytesToString(addr)
		require.NoError(t, err)
		return str
	}

	var (
		operator = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
		grantee  = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
		proposer = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	)

	testCases := map[string]struct {
		authorization foundation.Authorization
		expected      string
	}{
		"ReceiveFromTreasuryAuthorization": {
			&foundation.ReceiveFromTreasuryAuthorization{},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSubmitProposal\",\"value\":{\"exec\":1,\"messages\":[{\"type\":\"lbm-sdk/MsgGrant\",\"value\":{\"authority\":\"%s\",\"authorization\":{\"type\":\"lbm-sdk/ReceiveFromTreasuryAuthorization\",\"value\":{}},\"grantee\":\"%s\"}}],\"metadata\":\"ReceiveFromTreasuryAuthorization\",\"proposers\":[\"%s\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", bytesToString(operator), bytesToString(grantee), bytesToString(proposer)),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			grantMsg := &foundation.MsgGrant{
				Authority: bytesToString(operator),
				Grantee:   bytesToString(grantee),
			}
			err := grantMsg.SetAuthorization(tc.authorization)
			require.NoError(t, err)

			proposalMsg := &foundation.MsgSubmitProposal{
				Proposers: []string{bytesToString(proposer)},
				Metadata:  name,
				Exec:      foundation.Exec_EXEC_TRY,
			}
			err = proposalMsg.SetMsgs([]sdk.Msg{grantMsg})
			require.NoError(t, err)

			require.Equal(t, tc.expected, string(legacytx.StdSignBytes("foo", 1, 1, 1, legacytx.StdFee{}, []sdk.Msg{proposalMsg}, "memo")))
		})
	}
}
