package foundation_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

func TestMsgFundTreasury(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		from    sdk.AccAddress
		amount  sdk.Int
		valid   bool
	}{
		"valid msg": {
			from:    addrs[0],
			amount:  sdk.OneInt(),
			valid:   true,
		},
		"empty from": {
			amount:  sdk.OneInt(),
		},
		"zero amount": {
			from:    addrs[0],
			amount:  sdk.ZeroInt(),
 		},
	}

	for name, tc := range testCases {
		msg := foundation.MsgFundTreasury{
			From:    tc.from.String(),
			Amount:  sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, tc.amount)),
		}

		require.Equal(t, []sdk.AccAddress{tc.from}, msg.GetSigners(), name)

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgWithdrawFromTreasury(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		operator sdk.AccAddress
		to    sdk.AccAddress
		amount  sdk.Int
		valid   bool
	}{
		"valid msg": {
			operator: addrs[0],
			to:    addrs[1],
			amount:  sdk.OneInt(),
			valid:   true,
		},
		"empty operator": {
			to:    addrs[1],
			amount:  sdk.OneInt(),
		},
		"empty to": {
			operator: addrs[0],
			amount:  sdk.OneInt(),
		},
		"zero amount": {
			operator: addrs[0],
			to:    addrs[1],
			amount:  sdk.ZeroInt(),
		},
	}

	for name, tc := range testCases {
		msg := foundation.MsgWithdrawFromTreasury{
			Operator: tc.operator.String(),
			To:    tc.to.String(),
			Amount:  sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, tc.amount)),
		}

		require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners(), name)

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgUpdateMembers(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		operator sdk.AccAddress
		members    []foundation.Member
		valid   bool
	}{
		"valid msg": {
			operator: addrs[0],
			members: []foundation.Member{{
				Address: addrs[1].String(),
				Participating: true,
			}},
			valid:   true,
		},
		"empty operator": {
			members: []foundation.Member{{
				Address: addrs[1].String(),
				Participating: true,
			}},
		},
		"empty members": {
			operator: addrs[0],
			members: []foundation.Member{},
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
					Address: addrs[1].String(),
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
			Operator: tc.operator.String(),
			MemberUpdates: tc.members,
		}

		require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners(), name)

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgUpdateDecisionPolicy(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		operator sdk.AccAddress
		policy    foundation.DecisionPolicy
		valid   bool
	}{
		"valid threshold policy": {
			operator: addrs[0],
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: sdk.NewDec(3),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
			valid:   true,
		},
		"valid percentage policy": {
			operator: addrs[0],
			policy: &foundation.PercentageDecisionPolicy{
				Percentage: sdk.OneDec(),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
			valid:   true,
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
				Windows: &foundation.DecisionPolicyWindows{
				},
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

		require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners(), name)

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgSubmitProposal(t *testing.T) {
	addrs := make([]sdk.AccAddress, 3)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		proposers []sdk.AccAddress
		msgs []sdk.Msg
		exec foundation.Exec
		valid   bool
	}{
		"valid msg": {
			proposers: []sdk.AccAddress{addrs[0]},
			msgs: []sdk.Msg{&foundation.MsgWithdrawFromTreasury{
				Operator: addrs[1].String(),
				To: addrs[2].String(),
				Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			}},
			valid:   true,
		},
		"empty proposers": {
			proposers: []sdk.AccAddress{},
			msgs: []sdk.Msg{&foundation.MsgWithdrawFromTreasury{
				Operator: addrs[1].String(),
				To: addrs[2].String(),
				Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			}},
		},
		"invalid proposer": {
			proposers: []sdk.AccAddress{""},
			msgs: []sdk.Msg{&foundation.MsgWithdrawFromTreasury{
				Operator: addrs[1].String(),
				To: addrs[2].String(),
				Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			}},
		},
		"duplicate proposers": {
			proposers: []sdk.AccAddress{addrs[0], addrs[0]},
			msgs: []sdk.Msg{&foundation.MsgWithdrawFromTreasury{
				Operator: addrs[1].String(),
				To: addrs[2].String(),
				Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			}},
		},
		"empty msgs": {
			proposers: []sdk.AccAddress{addrs[0]},
		},
		"invalid msg": {
			proposers: []sdk.AccAddress{addrs[0]},
			msgs: []sdk.Msg{&foundation.MsgWithdrawFromTreasury{}},
		},
		"invalid exec": {
			proposers: []sdk.AccAddress{addrs[0]},
			msgs: []sdk.Msg{&foundation.MsgWithdrawFromTreasury{
				Operator: addrs[1].String(),
				To: addrs[2].String(),
				Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
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
			Exec: tc.exec,
		}
		err := msg.SetMsgs(tc.msgs)
		require.NoError(t, err, name)

		require.Equal(t, tc.proposers, msg.GetSigners(), name)

		err = msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgWithdrawProposal(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		id uint64
		address sdk.AccAddress
		valid   bool
	}{
		"valid msg": {
			id: 1,
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
			Address: tc.address.String(),
		}

		require.Equal(t, []sdk.AccAddress{tc.address}, msg.GetSigners(), name)

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgVote(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		id uint64
		voter sdk.AccAddress
		option foundation.VoteOption
		exec foundation.Exec
		valid   bool
	}{
		"valid msg": {
			id: 1,
			voter: addrs[0],
			option: foundation.VOTE_OPTION_YES,
			valid:   true,
		},
		"empty proposal id": {
			voter: addrs[0],
			option: foundation.VOTE_OPTION_YES,
		},
		"empty voter": {
			id: 1,
			option: foundation.VOTE_OPTION_YES,
		},
		"empty option": {
			id: 1,
			voter: addrs[0],
		},
		"invalid option": {
			id: 1,
			voter: addrs[0],
			option: -1,
		},
		"invalid exec": {
			id: 1,
			voter: addrs[0],
			option: foundation.VOTE_OPTION_YES,
			exec: -1,
		},
	}

	for name, tc := range testCases {
		msg := foundation.MsgVote{
			ProposalId: tc.id,
			Voter: tc.voter.String(),
			Option: tc.option,
			Exec: tc.exec,
		}

		require.Equal(t, []sdk.AccAddress{tc.voter}, msg.GetSigners(), name)

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgExec(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		id uint64
		signer sdk.AccAddress
		valid   bool
	}{
		"valid msg": {
			id: 1,
			signer: addrs[0],
			valid:   true,
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
			Signer: tc.signer.String(),
		}

		require.Equal(t, []sdk.AccAddress{tc.signer}, msg.GetSigners(), name)

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgLeaveFoundation(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		address sdk.AccAddress
		valid   bool
	}{
		"valid msg": {
			address: addrs[0],
			valid:   true,
		},
		"empty address": {
		},
	}

	for name, tc := range testCases {
		msg := foundation.MsgLeaveFoundation{
			Address: tc.address.String(),
		}

		require.Equal(t, []sdk.AccAddress{tc.address}, msg.GetSigners(), name)

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgGrant(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		operator sdk.AccAddress
		grantee sdk.AccAddress
		authorization foundation.Authorization
		valid   bool
	}{
		"valid msg": {
			operator: addrs[0],
			grantee: addrs[1],
			authorization: &foundation.ReceiveFromTreasuryAuthorization{},
			valid:   true,
		},
		"empty operator": {
			grantee: addrs[1],
			authorization: &foundation.ReceiveFromTreasuryAuthorization{},
		},
		"empty grantee": {
			operator: addrs[0],
			authorization: &foundation.ReceiveFromTreasuryAuthorization{},
		},
		"empty authorization": {
			operator: addrs[0],
			grantee: addrs[1],
		},
	}

	for name, tc := range testCases {
		msg := foundation.MsgGrant{
			Operator: tc.operator.String(),
			Grantee: tc.grantee.String(),
		}
		if tc.authorization != nil {
			msg.SetAuthorization(tc.authorization)
		}

		require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners(), name)

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgRevoke(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		operator sdk.AccAddress
		grantee sdk.AccAddress
		msgTypeURL string
		valid   bool
	}{
		"valid msg": {
			operator: addrs[0],
			grantee: addrs[1],
			msgTypeURL: foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
			valid:   true,
		},
		"empty operator": {
			grantee: addrs[1],
			msgTypeURL: foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
		},
		"empty grantee": {
			operator: addrs[0],
			msgTypeURL: foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
		},
		"empty url": {
			operator: addrs[0],
			grantee: addrs[1],
		},
	}

	for name, tc := range testCases {
		msg := foundation.MsgRevoke{
			Operator: tc.operator.String(),
			Grantee: tc.grantee.String(),
			MsgTypeUrl: tc.msgTypeURL,
		}

		require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners(), name)

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}
