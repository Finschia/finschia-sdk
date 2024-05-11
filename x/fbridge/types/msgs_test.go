package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/auth/legacy/legacytx"
	fbridgetypes "github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func TestAminoJSON(t *testing.T) {
	tx := legacytx.StdTx{}

	addrs := []string{
		"link1zf469e6y5zvsvkjz8vpr27j6txseyfnsh3ydze",
		"link1aydj7vdxljxq7cn2dlkrhrcwf5py8dnuqp32qd",
	}
	toAddr := "0xf7bAc63fc7CEaCf0589F25454Ecf5C2ce904997c"

	testCase := map[string]struct {
		msg          legacytx.LegacyMsg
		expectedType string
		expected     string
	}{
		"MsgUpdateParam": {
			&fbridgetypes.MsgUpdateParams{
				Authority: addrs[0],
				Params: fbridgetypes.Params{
					OperatorTrustLevel: fbridgetypes.Fraction{
						Numerator:   uint64(2),
						Denominator: uint64(3),
					},
					GuardianTrustLevel: fbridgetypes.Fraction{
						Numerator:   uint64(2),
						Denominator: uint64(3),
					},
					JudgeTrustLevel: fbridgetypes.Fraction{
						Numerator:   uint64(2),
						Denominator: uint64(3),
					},
					TimelockPeriod: uint64(86400000000000),
					ProposalPeriod: uint64(3600000000000),
					TargetDenom:    "kaia",
				},
			},
			"/lbm.fbridge.v1.MsgUpdateParams",
			"{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/fbridge/MsgUpdateParams\",\"value\":{\"authority\":\"link1zf469e6y5zvsvkjz8vpr27j6txseyfnsh3ydze\",\"params\":{\"guardian_trust_level\":{\"denominator\":\"3\",\"numerator\":\"2\"},\"judge_trust_level\":{\"denominator\":\"3\",\"numerator\":\"2\"},\"operator_trust_level\":{\"denominator\":\"3\",\"numerator\":\"2\"},\"proposal_period\":\"3600000000000\",\"target_denom\":\"kaia\",\"timelock_period\":\"86400000000000\"}}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}",
		},
		"MsgTransfer": {
			&fbridgetypes.MsgTransfer{
				Sender:   addrs[0],
				Receiver: toAddr,
				Amount:   sdk.NewInt(1000000),
			},
			"/lbm.fbridge.v1.MsgTransfer",
			"{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgTransfer\",\"value\":{\"amount\":\"1000000\",\"receiver\":\"0xf7bAc63fc7CEaCf0589F25454Ecf5C2ce904997c\",\"sender\":\"link1zf469e6y5zvsvkjz8vpr27j6txseyfnsh3ydze\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}",
		},
		"MsgSuggestRole": {
			&fbridgetypes.MsgSuggestRole{
				From:   addrs[0],
				Target: addrs[1],
				Role:   fbridgetypes.RoleGuardian,
			},
			"/lbm.fbridge.v1.MsgSuggestRole",
			"{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSuggestRole\",\"value\":{\"from\":\"link1zf469e6y5zvsvkjz8vpr27j6txseyfnsh3ydze\",\"role\":1,\"target\":\"link1aydj7vdxljxq7cn2dlkrhrcwf5py8dnuqp32qd\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}",
		},
		"MsgAddVoteForRole": {
			&fbridgetypes.MsgAddVoteForRole{
				From:       addrs[0],
				ProposalId: 0,
				Option:     fbridgetypes.OptionYes,
			},
			"/lbm.fbridge.v1.MsgAddVoteForRole",
			"{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgAddVoteForRole\",\"value\":{\"from\":\"link1zf469e6y5zvsvkjz8vpr27j6txseyfnsh3ydze\",\"option\":1}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}",
		},
		"MsgSetBridgeStatus": {
			&fbridgetypes.MsgSetBridgeStatus{
				Guardian: addrs[0],
				Status:   fbridgetypes.StatusActive,
			},
			"/lbm.fbridge.v1.MsgSetBridgeStatus",
			"{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSetBridgeStatus\",\"value\":{\"guardian\":\"link1zf469e6y5zvsvkjz8vpr27j6txseyfnsh3ydze\",\"status\":1}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}",
		},
	}

	for name, tc := range testCase {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tx.Msgs = []sdk.Msg{tc.msg}
			require.Equal(t, fbridgetypes.RouterKey, tc.msg.Route())
			require.Equal(t, tc.expectedType, tc.msg.Type())
			require.Equal(t, tc.expected, string(legacytx.StdSignBytes("foo", 1, 1, 1, legacytx.StdFee{}, []sdk.Msg{tc.msg}, "memo")))
		})
	}
}
