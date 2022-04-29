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
			valid:   false,
		},
		"zero amount": {
			from:    addrs[0],
			amount:  sdk.ZeroInt(),
			valid:   false,
		},
	}

	for name, tc := range testCases {
		msg := foundation.MsgFundTreasury{
			From:    tc.from.String(),
			Amount:  sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, tc.amount)),
		}

		require.Equal(t, []sdk.AccAddress{tc.from}, msg.GetSigners())

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgWithdrawFromTreasury(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
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
			to:    addrs[0],
			amount:  sdk.OneInt(),
			valid:   true,
		},
		"empty operator": {
			to:    addrs[0],
			amount:  sdk.OneInt(),
			valid:   false,
		},
		"empty to": {
			operator: addrs[0],
			amount:  sdk.OneInt(),
			valid:   false,
		},
		"zero amount": {
			operator: addrs[0],
			to:    addrs[0],
			amount:  sdk.ZeroInt(),
			valid:   false,
		},
	}

	for name, tc := range testCases {
		msg := foundation.MsgWithdrawFromTreasury{
			Operator: tc.operator.String(),
			To:    tc.to.String(),
			Amount:  sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, tc.amount)),
		}

		require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners())

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgUpdateMembers(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
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
				Address: addrs[0].String(),
				Weight: sdk.OneDec(),
			}},
			valid:   true,
		},
		"empty operator": {
			members: []foundation.Member{{
				Address: addrs[0].String(),
				Weight: sdk.OneDec(),
			}},
			valid:   false,
		},
		"empty members": {
			operator: addrs[0],
			members: []foundation.Member{},
			valid:   false,
		},
		"empty member address": {
			operator: addrs[0],
			members: []foundation.Member{{
				Weight: sdk.OneDec(),
			}},
			valid:   false,
		},
		"invalid member weight": {
			operator: addrs[0],
			members: []foundation.Member{{
				Address: addrs[0].String(),
				Weight: sdk.MustNewDecFromStr("0.5"),
			}},
			valid:   false,
		},
	}

	for name, tc := range testCases {
		msg := foundation.MsgUpdateMembers{
			Operator: tc.operator.String(),
			MemberUpdates: tc.members,
		}

		require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners())

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
		"valid msg": {
			operator: addrs[0],
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: sdk.NewDec(3),
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
			valid:   false,
		},
		"empty policy": {
			operator: addrs[0],
			valid:   false,
		},
		"zero threshold": {
			operator: addrs[0],
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: sdk.ZeroDec(),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
			valid:   false,
		},
		"zero voting period": {
			operator: addrs[0],
			policy: &foundation.ThresholdDecisionPolicy{
				Threshold: sdk.NewDec(3),
				Windows: &foundation.DecisionPolicyWindows{
				},
			},
			valid:   false,
		},
		"zero percentage": {
			operator: addrs[0],
			policy: &foundation.PercentageDecisionPolicy{
				Percentage: sdk.ZeroDec(),
				Windows: &foundation.DecisionPolicyWindows{
					VotingPeriod: time.Hour,
				},
			},
			valid:   false,
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

		require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners())

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgSubmitProposal(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
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
				Operator: addrs[0].String(),
				To: addrs[0].String(),
				Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			}},
			valid:   true,
		},
		"empty proposers": {
			proposers: []sdk.AccAddress{},
			msgs: []sdk.Msg{&foundation.MsgWithdrawFromTreasury{
				Operator: addrs[0].String(),
				To: addrs[0].String(),
				Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			}},
			valid:   false,
		},
		"invalid proposer": {
			proposers: []sdk.AccAddress{""},
			msgs: []sdk.Msg{&foundation.MsgWithdrawFromTreasury{
				Operator: addrs[0].String(),
				To: addrs[0].String(),
				Amount: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			}},
			valid:   false,
		},
		"empty msgs": {
			proposers: []sdk.AccAddress{addrs[0]},
			valid:   false,
		},
		"invalid msg": {
			proposers: []sdk.AccAddress{addrs[0]},
			msgs: []sdk.Msg{&foundation.MsgWithdrawFromTreasury{}},
			valid:   false,
		},
		"invalid exec": {
			proposers: []sdk.AccAddress{addrs[0]},
			msgs: []sdk.Msg{&foundation.MsgWithdrawFromTreasury{}},
			exec: -1,
			valid:   false,
		},
	}

	for name, tc := range testCases {
		var proposers []string
		for _, proposer := range tc.proposers {
			proposers = append(proposers, proposer.String())
		}

		msg := foundation.MsgSubmitProposal{
			Proposers: proposers,
		}
		err := msg.SetMsgs(tc.msgs)
		require.NoError(t, err, name)

		require.Equal(t, tc.proposers, msg.GetSigners())

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
			valid:   false,
		},
		"empty address": {
			id: 1,
			valid:   false,
		},
	}

	for name, tc := range testCases {
		msg := foundation.MsgWithdrawProposal{
			ProposalId: tc.id,
			Address: tc.address.String(),
		}

		require.Equal(t, []sdk.AccAddress{tc.address}, msg.GetSigners())

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
			valid:   false,
		},
		"empty voter": {
			id: 1,
			option: foundation.VOTE_OPTION_YES,
			valid:   false,
		},
		"empty option": {
			id: 1,
			voter: addrs[0],
			valid:   false,
		},
		"invalid option": {
			id: 1,
			voter: addrs[0],
			option: -1,
			valid:   false,
		},
		"invalid exec": {
			id: 1,
			voter: addrs[0],
			exec: -1,
			valid:   false,
		},
	}

	for name, tc := range testCases {
		msg := foundation.MsgVote{
			ProposalId: tc.id,
			Voter: tc.voter.String(),
			Option: tc.option,
			Exec: tc.exec,
		}

		require.Equal(t, []sdk.AccAddress{tc.voter}, msg.GetSigners())

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
			valid:   false,
		},
		"empty signer": {
			id: 1,
			valid:   false,
		},
	}

	for name, tc := range testCases {
		msg := foundation.MsgExec{
			ProposalId: tc.id,
			Signer: tc.signer.String(),
		}

		require.Equal(t, []sdk.AccAddress{tc.signer}, msg.GetSigners())

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
			valid:   false,
		},
	}

	for name, tc := range testCases {
		msg := foundation.MsgLeaveFoundation{
			Address: tc.address.String(),
		}

		require.Equal(t, []sdk.AccAddress{tc.address}, msg.GetSigners())

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}
