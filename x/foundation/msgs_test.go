package foundation_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	"github.com/line/lbm-sdk/testutil/testdata"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/auth/legacy/legacytx"
	"github.com/line/lbm-sdk/x/foundation"
)

func TestMsgUpdateParams(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		authority sdk.AccAddress
		params    foundation.Params
		valid     bool
	}{
		"valid msg": {
			authority: addrs[0],
			params: foundation.Params{
				FoundationTax: sdk.ZeroDec(),
			},
			valid: true,
		},
		"invalid authority": {
			params: foundation.Params{
				FoundationTax: sdk.ZeroDec(),
			},
		},
		"invalid params": {
			authority: addrs[0],
			params:    foundation.Params{},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := foundation.MsgUpdateParams{
				Authority: tc.authority.String(),
				Params:    tc.params,
			}

			err := msg.ValidateBasic()
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, []sdk.AccAddress{tc.authority}, msg.GetSigners())
		})
	}
}

func TestMsgFundTreasury(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		from   sdk.AccAddress
		amount sdk.Int
		valid  bool
	}{
		"valid msg": {
			from:   addrs[0],
			amount: sdk.OneInt(),
			valid:  true,
		},
		"empty from": {
			amount: sdk.OneInt(),
		},
		"zero amount": {
			from:   addrs[0],
			amount: sdk.ZeroInt(),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := foundation.MsgFundTreasury{
				From:   tc.from.String(),
				Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, tc.amount)),
			}

			err := msg.ValidateBasic()
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, []sdk.AccAddress{tc.from}, msg.GetSigners())
		})
	}
}

func TestMsgWithdrawFromTreasury(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		authority sdk.AccAddress
		to        sdk.AccAddress
		amount    sdk.Int
		valid     bool
	}{
		"valid msg": {
			authority: addrs[0],
			to:        addrs[1],
			amount:    sdk.OneInt(),
			valid:     true,
		},
		"empty authority": {
			to:     addrs[1],
			amount: sdk.OneInt(),
		},
		"empty to": {
			authority: addrs[0],
			amount:    sdk.OneInt(),
		},
		"zero amount": {
			authority: addrs[0],
			to:        addrs[1],
			amount:    sdk.ZeroInt(),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := foundation.MsgWithdrawFromTreasury{
				Authority: tc.authority.String(),
				To:        tc.to.String(),
				Amount:    sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, tc.amount)),
			}

			err := msg.ValidateBasic()
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, []sdk.AccAddress{tc.authority}, msg.GetSigners())
		})
	}
}

func TestMsgUpdateMembers(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		authority sdk.AccAddress
		members   []foundation.MemberRequest
		valid     bool
	}{
		"valid msg": {
			authority: addrs[0],
			members: []foundation.MemberRequest{{
				Address: addrs[1].String(),
			}},
			valid: true,
		},
		"empty authority": {
			members: []foundation.MemberRequest{{
				Address: addrs[1].String(),
			}},
		},
		"empty requests": {
			authority: addrs[0],
		},
		"invalid requests": {
			authority: addrs[0],
			members:   []foundation.MemberRequest{{}},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := foundation.MsgUpdateMembers{
				Authority:     tc.authority.String(),
				MemberUpdates: tc.members,
			}

			err := msg.ValidateBasic()
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, []sdk.AccAddress{tc.authority}, msg.GetSigners())
		})
	}
}

func TestMsgUpdateDecisionPolicy(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		authority sdk.AccAddress
		policy    foundation.DecisionPolicy
		valid     bool
	}{
		"valid threshold policy": {
			authority: addrs[0],
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: sdk.OneDec(),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
			valid: true,
		},
		"valid percentage policy": {
			authority: addrs[0],
			policy: &foundation.PercentageDecisionPolicy{
				Percentage: sdk.OneDec(),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
			valid: true,
		},
		"empty authority": {
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: sdk.OneDec(),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
		},
		"empty policy": {
			authority: addrs[0],
		},
		"zero threshold": {
			authority: addrs[0],
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: sdk.ZeroDec(),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
		},
		"zero voting period": {
			authority: addrs[0],
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: sdk.OneDec(),
				Windows:   &foundation.DecisionPolicyWindows{},
			},
		},
		"invalid percentage": {
			authority: addrs[0],
			policy: &foundation.PercentageDecisionPolicy{
				Percentage: sdk.NewDec(2),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := foundation.MsgUpdateDecisionPolicy{
				Authority: tc.authority.String(),
			}
			if tc.policy != nil {
				err := msg.SetDecisionPolicy(tc.policy)
				require.NoError(t, err)
			}

			err := msg.ValidateBasic()
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, []sdk.AccAddress{tc.authority}, msg.GetSigners())
		})
	}
}

func TestMsgSubmitProposal(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		proposers []sdk.AccAddress
		msgs      []sdk.Msg
		exec      foundation.Exec
		valid     bool
	}{
		"valid msg": {
			proposers: []sdk.AccAddress{addrs[0]},
			msgs:      []sdk.Msg{testdata.NewTestMsg()},
			valid:     true,
		},
		"empty proposers": {
			msgs: []sdk.Msg{testdata.NewTestMsg()},
		},
		"invalid proposer": {
			proposers: []sdk.AccAddress{nil},
			msgs:      []sdk.Msg{testdata.NewTestMsg()},
		},
		"duplicate proposers": {
			proposers: []sdk.AccAddress{addrs[0], addrs[0]},
			msgs:      []sdk.Msg{testdata.NewTestMsg()},
		},
		"empty msgs": {
			proposers: []sdk.AccAddress{addrs[0]},
		},
		"invalid msg": {
			proposers: []sdk.AccAddress{addrs[0]},
			msgs:      []sdk.Msg{&foundation.MsgWithdrawFromTreasury{}},
		},
		"invalid exec": {
			proposers: []sdk.AccAddress{addrs[0]},
			msgs:      []sdk.Msg{testdata.NewTestMsg()},
			exec:      -1,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var proposers []string
			for _, proposer := range tc.proposers {
				proposers = append(proposers, proposer.String())
			}

			msg := foundation.MsgSubmitProposal{
				Proposers: proposers,
				Exec:      tc.exec,
			}
			err := msg.SetMsgs(tc.msgs)
			require.NoError(t, err)

			err = msg.ValidateBasic()
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, tc.proposers, msg.GetSigners())
		})
	}
}

func TestMsgWithdrawProposal(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		id      uint64
		address sdk.AccAddress
		valid   bool
	}{
		"valid msg": {
			id:      1,
			address: addrs[0],
			valid:   true,
		},
		"empty proposal id": {
			address: addrs[0],
		},
		"empty address": {
			id: 1,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := foundation.MsgWithdrawProposal{
				ProposalId: tc.id,
				Address:    tc.address.String(),
			}

			err := msg.ValidateBasic()
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, []sdk.AccAddress{tc.address}, msg.GetSigners())
		})
	}
}

func TestMsgVote(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		id     uint64
		voter  sdk.AccAddress
		option foundation.VoteOption
		exec   foundation.Exec
		valid  bool
	}{
		"valid msg": {
			id:     1,
			voter:  addrs[0],
			option: foundation.VOTE_OPTION_YES,
			valid:  true,
		},
		"empty proposal id": {
			voter:  addrs[0],
			option: foundation.VOTE_OPTION_YES,
		},
		"empty voter": {
			id:     1,
			option: foundation.VOTE_OPTION_YES,
		},
		"empty option": {
			id:    1,
			voter: addrs[0],
		},
		"invalid option": {
			id:     1,
			voter:  addrs[0],
			option: -1,
		},
		"invalid exec": {
			id:     1,
			voter:  addrs[0],
			option: foundation.VOTE_OPTION_YES,
			exec:   -1,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := foundation.MsgVote{
				ProposalId: tc.id,
				Voter:      tc.voter.String(),
				Option:     tc.option,
				Exec:       tc.exec,
			}

			err := msg.ValidateBasic()
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, []sdk.AccAddress{tc.voter}, msg.GetSigners())
		})
	}
}

func TestMsgExec(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		id     uint64
		signer sdk.AccAddress
		valid  bool
	}{
		"valid msg": {
			id:     1,
			signer: addrs[0],
			valid:  true,
		},
		"empty proposal id": {
			signer: addrs[0],
		},
		"empty signer": {
			id: 1,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := foundation.MsgExec{
				ProposalId: tc.id,
				Signer:     tc.signer.String(),
			}

			err := msg.ValidateBasic()
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, []sdk.AccAddress{tc.signer}, msg.GetSigners())
		})
	}
}

func TestMsgLeaveFoundation(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		address sdk.AccAddress
		valid   bool
	}{
		"valid msg": {
			address: addrs[0],
			valid:   true,
		},
		"empty address": {},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := foundation.MsgLeaveFoundation{
				Address: tc.address.String(),
			}

			err := msg.ValidateBasic()
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, []sdk.AccAddress{tc.address}, msg.GetSigners())
		})
	}
}

func TestMsgGrant(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		authority     sdk.AccAddress
		grantee       sdk.AccAddress
		authorization foundation.Authorization
		valid         bool
	}{
		"valid msg": {
			authority:     addrs[0],
			grantee:       addrs[1],
			authorization: &foundation.ReceiveFromTreasuryAuthorization{},
			valid:         true,
		},
		"empty authority": {
			grantee:       addrs[1],
			authorization: &foundation.ReceiveFromTreasuryAuthorization{},
		},
		"empty grantee": {
			authority:     addrs[0],
			authorization: &foundation.ReceiveFromTreasuryAuthorization{},
		},
		"empty authorization": {
			authority: addrs[0],
			grantee:   addrs[1],
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := foundation.MsgGrant{
				Authority: tc.authority.String(),
				Grantee:   tc.grantee.String(),
			}
			if tc.authorization != nil {
				msg.SetAuthorization(tc.authorization)
			}

			err := msg.ValidateBasic()
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, []sdk.AccAddress{tc.authority}, msg.GetSigners())
		})
	}
}

func TestMsgRevoke(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		authority  sdk.AccAddress
		grantee    sdk.AccAddress
		msgTypeURL string
		valid      bool
	}{
		"valid msg": {
			authority:  addrs[0],
			grantee:    addrs[1],
			msgTypeURL: foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
			valid:      true,
		},
		"empty authority": {
			grantee:    addrs[1],
			msgTypeURL: foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
		},
		"empty grantee": {
			authority:  addrs[0],
			msgTypeURL: foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
		},
		"empty url": {
			authority: addrs[0],
			grantee:   addrs[1],
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := foundation.MsgRevoke{
				Authority:  tc.authority.String(),
				Grantee:    tc.grantee.String(),
				MsgTypeUrl: tc.msgTypeURL,
			}

			err := msg.ValidateBasic()
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, []sdk.AccAddress{tc.authority}, msg.GetSigners())
		})
	}
}

func TestAminoJSON(t *testing.T) {
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
				From:   addrs[0].String(),
				Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgFundTreasury\",\"value\":{\"amount\":[{\"amount\":\"1\",\"denom\":\"stake\"}],\"from\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
		"MsgVote": {
			&foundation.MsgVote{
				ProposalId: 1,
				Voter:      addrs[0].String(),
				Option:     foundation.VOTE_OPTION_YES,
				Metadata:   "I'm YES",
				Exec:       foundation.Exec_EXEC_UNSPECIFIED,
			},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgVote\",\"value\":{\"metadata\":\"I'm YES\",\"option\":1,\"proposal_id\":\"1\",\"voter\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
		"MsgExec": {
			&foundation.MsgExec{
				ProposalId: 1,
				Signer:     addrs[0].String(),
			},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgExec\",\"value\":{\"proposal_id\":\"1\",\"signer\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
		"MsgLeaveFoundation": {
			&foundation.MsgLeaveFoundation{Address: addrs[0].String()},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgLeaveFoundation\",\"value\":{\"address\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
		"MsgWithdrawProposal": {
			&foundation.MsgWithdrawProposal{
				ProposalId: 1,
				Address:    addrs[0].String(),
			},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgWithdrawProposal\",\"value\":{\"address\":\"%s\",\"proposal_id\":\"1\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
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
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	var proposer = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	testCases := map[string]struct {
		msg      sdk.Msg
		expected string
	}{
		"MsgUpdateParams": {
			&foundation.MsgUpdateParams{
				Authority: addrs[0].String(),
				Params:    foundation.Params{FoundationTax: sdk.ZeroDec()},
			},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSubmitProposal\",\"value\":{\"exec\":1,\"messages\":[{\"type\":\"lbm-sdk/MsgUpdateParams\",\"value\":{\"authority\":\"%s\",\"params\":{\"foundation_tax\":\"0.000000000000000000\"}}}],\"metadata\":\"MsgUpdateParams\",\"proposers\":[\"%s\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), proposer.String()),
		},
		"MsgWithdrawFromTreasury": {
			&foundation.MsgWithdrawFromTreasury{
				Authority: addrs[0].String(),
				To:        addrs[1].String(),
				Amount:    sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000000))),
			},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSubmitProposal\",\"value\":{\"exec\":1,\"messages\":[{\"type\":\"lbm-sdk/MsgWithdrawFromTreasury\",\"value\":{\"amount\":[{\"amount\":\"1000000\",\"denom\":\"stake\"}],\"authority\":\"%s\",\"to\":\"%s\"}}],\"metadata\":\"MsgWithdrawFromTreasury\",\"proposers\":[\"%s\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String(), proposer.String()),
		},
		"MsgUpdateMembers": {
			&foundation.MsgUpdateMembers{
				Authority: addrs[0].String(),
				MemberUpdates: []foundation.MemberRequest{{
					Address: addrs[1].String(),
				}},
			},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSubmitProposal\",\"value\":{\"exec\":1,\"messages\":[{\"type\":\"lbm-sdk/MsgUpdateMembers\",\"value\":{\"authority\":\"%s\",\"member_updates\":[{\"address\":\"%s\"}]}}],\"metadata\":\"MsgUpdateMembers\",\"proposers\":[\"%s\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String(), proposer.String()),
		},
		"MsgRevoke": {
			&foundation.MsgRevoke{
				Authority:  addrs[0].String(),
				Grantee:    addrs[1].String(),
				MsgTypeUrl: foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
			},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSubmitProposal\",\"value\":{\"exec\":1,\"messages\":[{\"type\":\"lbm-sdk/MsgRevoke\",\"value\":{\"authority\":\"%s\",\"grantee\":\"%s\",\"msg_type_url\":\"/lbm.foundation.v1.MsgWithdrawFromTreasury\"}}],\"metadata\":\"MsgRevoke\",\"proposers\":[\"%s\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String(), proposer.String()),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			proposalMsg := &foundation.MsgSubmitProposal{
				Proposers: []string{proposer.String()},
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

	testCases := map[string]struct {
		policy   foundation.DecisionPolicy
		expected string
	}{
		"ThresholdDecisionPolicy": {
			&foundation.ThresholdDecisionPolicy{
				Threshold: sdk.OneDec(),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSubmitProposal\",\"value\":{\"exec\":1,\"messages\":[{\"type\":\"lbm-sdk/MsgUpdateDecisionPolicy\",\"value\":{\"authority\":\"%s\",\"decision_policy\":{\"type\":\"lbm-sdk/ThresholdDecisionPolicy\",\"value\":{\"threshold\":\"1.000000000000000000\",\"windows\":{\"min_execution_period\":\"0\",\"voting_period\":\"3600000000000\"}}}}}],\"metadata\":\"ThresholdDecisionPolicy\",\"proposers\":[\"%s\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", authority, proposer),
		},
		"PercentageDecisionPolicy": {
			&foundation.PercentageDecisionPolicy{
				Percentage: sdk.OneDec(),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSubmitProposal\",\"value\":{\"exec\":1,\"messages\":[{\"type\":\"lbm-sdk/MsgUpdateDecisionPolicy\",\"value\":{\"authority\":\"%s\",\"decision_policy\":{\"type\":\"lbm-sdk/PercentageDecisionPolicy\",\"value\":{\"percentage\":\"1.000000000000000000\",\"windows\":{\"min_execution_period\":\"0\",\"voting_period\":\"3600000000000\"}}}}}],\"metadata\":\"PercentageDecisionPolicy\",\"proposers\":[\"%s\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", authority, proposer),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			policyMsg := &foundation.MsgUpdateDecisionPolicy{
				Authority: authority.String(),
			}
			err := policyMsg.SetDecisionPolicy(tc.policy)
			require.NoError(t, err)

			err = policyMsg.ValidateBasic()
			require.NoError(t, err)

			proposalMsg := &foundation.MsgSubmitProposal{
				Proposers: []string{proposer.String()},
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
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSubmitProposal\",\"value\":{\"exec\":1,\"messages\":[{\"type\":\"lbm-sdk/MsgGrant\",\"value\":{\"authority\":\"%s\",\"authorization\":{\"type\":\"lbm-sdk/ReceiveFromTreasuryAuthorization\",\"value\":{}},\"grantee\":\"%s\"}}],\"metadata\":\"ReceiveFromTreasuryAuthorization\",\"proposers\":[\"%s\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", operator.String(), grantee.String(), proposer.String()),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			grantMsg := &foundation.MsgGrant{
				Authority: operator.String(),
				Grantee:   grantee.String(),
			}
			err := grantMsg.SetAuthorization(tc.authorization)
			require.NoError(t, err)

			err = grantMsg.ValidateBasic()
			require.NoError(t, err)

			proposalMsg := &foundation.MsgSubmitProposal{
				Proposers: []string{proposer.String()},
				Metadata:  name,
				Exec:      foundation.Exec_EXEC_TRY,
			}
			err = proposalMsg.SetMsgs([]sdk.Msg{grantMsg})
			require.NoError(t, err)

			require.Equal(t, tc.expected, string(legacytx.StdSignBytes("foo", 1, 1, 1, legacytx.StdFee{}, []sdk.Msg{proposalMsg}, "memo")))
		})
	}
}
