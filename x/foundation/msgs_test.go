package foundation_test

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

func NewMsgFundTreasury(fromAddr sdk.AccAddress) *foundation.MsgFundTreasury {
	return &foundation.MsgFundTreasury{From: fromAddr.String()}
}

func NewMsgWithdrawFromTreasury(fromAddr sdk.AccAddress) *foundation.MsgWithdrawFromTreasury {
	return &foundation.MsgWithdrawFromTreasury{Operator: fromAddr.String()}
}

func NewMsgUpdateMembers(fromAddr sdk.AccAddress) *foundation.MsgUpdateMembers {
	return &foundation.MsgUpdateMembers{Operator: fromAddr.String()}
}

func NewMsgUpdateDecisionPolicy(fromAddr sdk.AccAddress) *foundation.MsgUpdateDecisionPolicy {
	return &foundation.MsgUpdateDecisionPolicy{Operator: fromAddr.String()}
}

func NewMsgSubmitProposal(fromAddr sdk.AccAddress) *foundation.MsgSubmitProposal {
	return &foundation.MsgSubmitProposal{Proposers: []string{fromAddr.String()}}
}

func NewMsgWithdrawProposal(fromAddr sdk.AccAddress) *foundation.MsgWithdrawProposal {
	return &foundation.MsgWithdrawProposal{Address: fromAddr.String()}
}

func NewMsgVote(fromAddr sdk.AccAddress) *foundation.MsgVote {
	return &foundation.MsgVote{Voter: fromAddr.String()}
}

func NewMsgExec(fromAddr sdk.AccAddress) *foundation.MsgExec {
	return &foundation.MsgExec{Signer: fromAddr.String()}
}

func NewMsgLeaveFoundation(fromAddr sdk.AccAddress) *foundation.MsgLeaveFoundation {
	return &foundation.MsgLeaveFoundation{Address: fromAddr.String()}
}

func NewMsgGrant(fromAddr sdk.AccAddress) *foundation.MsgGrant {
	return &foundation.MsgGrant{Operator: fromAddr.String()}
}

func NewMsgRevoke(fromAddr sdk.AccAddress) *foundation.MsgRevoke {
	return &foundation.MsgRevoke{Operator: fromAddr.String()}
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
		msg := foundation.MsgFundTreasury{
			From:   tc.from.String(),
			Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, tc.amount)),
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.from}, msg.GetSigners(), name)
	}
}

func TestMsgWithdrawFromTreasury(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		operator sdk.AccAddress
		to       sdk.AccAddress
		amount   sdk.Int
		valid    bool
	}{
		"valid msg": {
			operator: addrs[0],
			to:       addrs[1],
			amount:   sdk.OneInt(),
			valid:    true,
		},
		"empty operator": {
			to:     addrs[1],
			amount: sdk.OneInt(),
		},
		"empty to": {
			operator: addrs[0],
			amount:   sdk.OneInt(),
		},
		"zero amount": {
			operator: addrs[0],
			to:       addrs[1],
			amount:   sdk.ZeroInt(),
		},
	}

	for name, tc := range testCases {
		msg := foundation.MsgWithdrawFromTreasury{
			Operator: tc.operator.String(),
			To:       tc.to.String(),
			Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, tc.amount)),
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners(), name)
	}
}

func TestMsgUpdateMembers(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		operator sdk.AccAddress
		members  []foundation.Member
		valid    bool
	}{
		"valid msg": {
			operator: addrs[0],
			members: []foundation.Member{{
				Address:       addrs[1].String(),
				Participating: true,
			}},
			valid: true,
		},
		"empty operator": {
			members: []foundation.Member{{
				Address:       addrs[1].String(),
				Participating: true,
			}},
		},
		"empty members": {
			operator: addrs[0],
		},
		"empty member address": {
			operator: addrs[0],
			members: []foundation.Member{{
				Participating: true,
			}},
		},
		"duplicate updates": {
			operator: addrs[0],
			members: []foundation.Member{
				{
					Address:       addrs[1].String(),
					Participating: true,
				},
				{
					Address: addrs[1].String(),
				},
			},
		},
	}

	for name, tc := range testCases {
		msg := foundation.MsgUpdateMembers{
			Operator:      tc.operator.String(),
			MemberUpdates: tc.members,
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners(), name)
	}
}

func TestMsgUpdateDecisionPolicy(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		operator sdk.AccAddress
		policy   foundation.DecisionPolicy
		valid    bool
	}{
		"valid threshold policy": {
			operator: addrs[0],
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: sdk.NewDec(3),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
			valid: true,
		},
		"valid percentage policy": {
			operator: addrs[0],
			policy: &foundation.PercentageDecisionPolicy{
				Percentage: sdk.OneDec(),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
			valid: true,
		},
		"empty operator": {
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: sdk.NewDec(3),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
		},
		"empty policy": {
			operator: addrs[0],
		},
		"zero threshold": {
			operator: addrs[0],
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: sdk.ZeroDec(),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
		},
		"zero voting period": {
			operator: addrs[0],
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: sdk.NewDec(3),
				Windows:   &foundation.DecisionPolicyWindows{},
			},
		},
		"invalid percentage": {
			operator: addrs[0],
			policy: &foundation.PercentageDecisionPolicy{
				Percentage: sdk.NewDec(2),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
		},
	}

	for name, tc := range testCases {
		msg := foundation.MsgUpdateDecisionPolicy{
			Operator: tc.operator.String(),
		}
		if tc.policy != nil {
			err := msg.SetDecisionPolicy(tc.policy)
			require.NoError(t, err, name)
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners(), name)
	}
}

func TestMsgSubmitProposal(t *testing.T) {
	addrs := make([]sdk.AccAddress, 3)
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
			msgs: []sdk.Msg{&foundation.MsgWithdrawFromTreasury{
				Operator: addrs[1].String(),
				To:       addrs[2].String(),
				Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			}},
			valid: true,
		},
		"empty proposers": {
			msgs: []sdk.Msg{&foundation.MsgWithdrawFromTreasury{
				Operator: addrs[1].String(),
				To:       addrs[2].String(),
				Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			}},
		},
		"invalid proposer": {
			proposers: []sdk.AccAddress{nil},
			msgs: []sdk.Msg{&foundation.MsgWithdrawFromTreasury{
				Operator: addrs[1].String(),
				To:       addrs[2].String(),
				Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			}},
		},
		"duplicate proposers": {
			proposers: []sdk.AccAddress{addrs[0], addrs[0]},
			msgs: []sdk.Msg{&foundation.MsgWithdrawFromTreasury{
				Operator: addrs[1].String(),
				To:       addrs[2].String(),
				Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			}},
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
			msgs: []sdk.Msg{&foundation.MsgWithdrawFromTreasury{
				Operator: addrs[1].String(),
				To:       addrs[2].String(),
				Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			}},
			exec: -1,
		},
	}

	for name, tc := range testCases {
		var proposers []string
		for _, proposer := range tc.proposers {
			proposers = append(proposers, proposer.String())
		}

		msg := foundation.MsgSubmitProposal{
			Proposers: proposers,
			Exec:      tc.exec,
		}
		err := msg.SetMsgs(tc.msgs)
		require.NoError(t, err, name)

		err = msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, tc.proposers, msg.GetSigners(), name)
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
		msg := foundation.MsgWithdrawProposal{
			ProposalId: tc.id,
			Address:    tc.address.String(),
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.address}, msg.GetSigners(), name)
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
		msg := foundation.MsgVote{
			ProposalId: tc.id,
			Voter:      tc.voter.String(),
			Option:     tc.option,
			Exec:       tc.exec,
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.voter}, msg.GetSigners(), name)
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
		msg := foundation.MsgExec{
			ProposalId: tc.id,
			Signer:     tc.signer.String(),
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.signer}, msg.GetSigners(), name)
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
		msg := foundation.MsgLeaveFoundation{
			Address: tc.address.String(),
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.address}, msg.GetSigners(), name)
	}
}

func TestMsgGrant(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		operator      sdk.AccAddress
		grantee       sdk.AccAddress
		authorization foundation.Authorization
		valid         bool
	}{
		"valid msg": {
			operator:      addrs[0],
			grantee:       addrs[1],
			authorization: &foundation.ReceiveFromTreasuryAuthorization{},
			valid:         true,
		},
		"empty operator": {
			grantee:       addrs[1],
			authorization: &foundation.ReceiveFromTreasuryAuthorization{},
		},
		"empty grantee": {
			operator:      addrs[0],
			authorization: &foundation.ReceiveFromTreasuryAuthorization{},
		},
		"empty authorization": {
			operator: addrs[0],
			grantee:  addrs[1],
		},
	}

	for name, tc := range testCases {
		msg := foundation.MsgGrant{
			Operator: tc.operator.String(),
			Grantee:  tc.grantee.String(),
		}
		if tc.authorization != nil {
			msg.SetAuthorization(tc.authorization)
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners(), name)
	}
}

func TestMsgRevoke(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		operator   sdk.AccAddress
		grantee    sdk.AccAddress
		msgTypeURL string
		valid      bool
	}{
		"valid msg": {
			operator:   addrs[0],
			grantee:    addrs[1],
			msgTypeURL: foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
			valid:      true,
		},
		"empty operator": {
			grantee:    addrs[1],
			msgTypeURL: foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
		},
		"empty grantee": {
			operator:   addrs[0],
			msgTypeURL: foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
		},
		"empty url": {
			operator: addrs[0],
			grantee:  addrs[1],
		},
	}

	for name, tc := range testCases {
		msg := foundation.MsgRevoke{
			Operator:   tc.operator.String(),
			Grantee:    tc.grantee.String(),
			MsgTypeUrl: tc.msgTypeURL,
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners(), name)
	}
}

func TestMsgFundTreasuryGetSigners(t *testing.T) {
	res := NewMsgFundTreasury(sdk.AccAddress([]byte("input111111111111111"))).GetSigners()
	bytes := sdk.MustAccAddressFromBech32(res[0].String())
	require.Equal(t, "696e707574313131313131313131313131313131", fmt.Sprintf("%v", hex.EncodeToString(bytes)))
}

func TestMsgWithdrawFromTreasuryGetSigners(t *testing.T) {
	res := NewMsgFundTreasury(sdk.AccAddress([]byte("input111111111111111"))).GetSigners()
	bytes := sdk.MustAccAddressFromBech32(res[0].String())
	require.Equal(t, "696e707574313131313131313131313131313131", fmt.Sprintf("%v", hex.EncodeToString(bytes)))
}

func TestMsgUpdateMembersGetSigners(t *testing.T) {
	res := NewMsgUpdateMembers(sdk.AccAddress([]byte("input111111111111111"))).GetSigners()
	bytes := sdk.MustAccAddressFromBech32(res[0].String())
	require.Equal(t, "696e707574313131313131313131313131313131", fmt.Sprintf("%v", hex.EncodeToString(bytes)))
}

func TestMsgUpdateDecisionPolicyGetSigners(t *testing.T) {
	res := NewMsgUpdateDecisionPolicy(sdk.AccAddress([]byte("input111111111111111"))).GetSigners()
	bytes := sdk.MustAccAddressFromBech32(res[0].String())
	require.Equal(t, "696e707574313131313131313131313131313131", fmt.Sprintf("%v", hex.EncodeToString(bytes)))
}

func TestMsgSubmitProposalGetSigners(t *testing.T) {
	res := NewMsgSubmitProposal(sdk.AccAddress([]byte("input111111111111111"))).GetSigners()
	bytes := sdk.MustAccAddressFromBech32(res[0].String())
	require.Equal(t, "696e707574313131313131313131313131313131", fmt.Sprintf("%v", hex.EncodeToString(bytes)))
}

func TestMsgWithdrawProposalGetSigners(t *testing.T) {
	res := NewMsgWithdrawProposal(sdk.AccAddress([]byte("input111111111111111"))).GetSigners()
	bytes := sdk.MustAccAddressFromBech32(res[0].String())
	require.Equal(t, "696e707574313131313131313131313131313131", fmt.Sprintf("%v", hex.EncodeToString(bytes)))
}

func TestMsgVoteGetSigners(t *testing.T) {
	res := NewMsgVote(sdk.AccAddress([]byte("input111111111111111"))).GetSigners()
	bytes := sdk.MustAccAddressFromBech32(res[0].String())
	require.Equal(t, "696e707574313131313131313131313131313131", fmt.Sprintf("%v", hex.EncodeToString(bytes)))
}

func TestMsgExecGetSigners(t *testing.T) {
	res := NewMsgExec(sdk.AccAddress([]byte("input111111111111111"))).GetSigners()
	bytes := sdk.MustAccAddressFromBech32(res[0].String())
	require.Equal(t, "696e707574313131313131313131313131313131", fmt.Sprintf("%v", hex.EncodeToString(bytes)))
}

func TestMsgLeaveFoundationGetSigners(t *testing.T) {
	res := NewMsgLeaveFoundation(sdk.AccAddress([]byte("input111111111111111"))).GetSigners()
	bytes := sdk.MustAccAddressFromBech32(res[0].String())
	require.Equal(t, "696e707574313131313131313131313131313131", fmt.Sprintf("%v", hex.EncodeToString(bytes)))
}

func TestMsgGrantGetSigners(t *testing.T) {
	res := NewMsgGrant(sdk.AccAddress([]byte("input111111111111111"))).GetSigners()
	bytes := sdk.MustAccAddressFromBech32(res[0].String())
	require.Equal(t, "696e707574313131313131313131313131313131", fmt.Sprintf("%v", hex.EncodeToString(bytes)))
}

func TestMsgRevokeGetSigners(t *testing.T) {
	res := NewMsgRevoke(sdk.AccAddress([]byte("input111111111111111"))).GetSigners()
	bytes := sdk.MustAccAddressFromBech32(res[0].String())
	require.Equal(t, "696e707574313131313131313131313131313131", fmt.Sprintf("%v", hex.EncodeToString(bytes)))
}
