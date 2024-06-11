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
		"MsgProvision": {
			&fbridgetypes.MsgProvision{
				From:     addrs[0],
				Seq:      uint64(1),
				Sender:   addrs[0],
				Receiver: addrs[1],
				Amount:   sdk.NewInt(1000000),
			},
			"/lbm.fbridge.v1.MsgProvision",
			"{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgProvision\",\"value\":{\"amount\":\"1000000\",\"from\":\"link1zf469e6y5zvsvkjz8vpr27j6txseyfnsh3ydze\",\"receiver\":\"link1aydj7vdxljxq7cn2dlkrhrcwf5py8dnuqp32qd\",\"sender\":\"link1zf469e6y5zvsvkjz8vpr27j6txseyfnsh3ydze\",\"seq\":\"1\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}",
		},
		"MsgHoldTransfer": {
			&fbridgetypes.MsgHoldTransfer{
				From: addrs[0],
				Seq:  100,
			},
			"/lbm.fbridge.v1.MsgHoldTransfer",
			"{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgHoldTransfer\",\"value\":{\"from\":\"link1zf469e6y5zvsvkjz8vpr27j6txseyfnsh3ydze\",\"seq\":\"100\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}",
		},
		"MsgReleaseTransfer": {
			&fbridgetypes.MsgReleaseTransfer{
				From: addrs[0],
				Seq:  200,
			},
			"/lbm.fbridge.v1.MsgReleaseTransfer",
			"{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgReleaseTransfer\",\"value\":{\"from\":\"link1zf469e6y5zvsvkjz8vpr27j6txseyfnsh3ydze\",\"seq\":\"200\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}",
		},
		"MsgRemoveProvision": {
			&fbridgetypes.MsgRemoveProvision{
				From: addrs[0],
				Seq:  300,
			},
			"/lbm.fbridge.v1.MsgRemoveProvision",
			"{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgRemoveProvision\",\"value\":{\"from\":\"link1zf469e6y5zvsvkjz8vpr27j6txseyfnsh3ydze\",\"seq\":\"300\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}",
		},
		"MsgClaimBatch": {
			&fbridgetypes.MsgClaimBatch{
				From:      addrs[0],
				MaxClaims: 50,
			},
			"/lbm.fbridge.v1.MsgClaimBatch",
			"{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgClaimBatch\",\"value\":{\"from\":\"link1zf469e6y5zvsvkjz8vpr27j6txseyfnsh3ydze\",\"max_claims\":\"50\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}",
		},
		"MsgClaim": {
			&fbridgetypes.MsgClaim{
				From: addrs[0],
				Seq:  400,
			},
			"/lbm.fbridge.v1.MsgClaim",
			"{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgClaim\",\"value\":{\"from\":\"link1zf469e6y5zvsvkjz8vpr27j6txseyfnsh3ydze\",\"seq\":\"400\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}",
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
		t.Run(name, func(t *testing.T) {
			tx.Msgs = []sdk.Msg{tc.msg}
			require.Equal(t, fbridgetypes.RouterKey, tc.msg.Route())
			require.Equal(t, tc.expectedType, tc.msg.Type())
			require.Equal(t, tc.expected, string(legacytx.StdSignBytes("foo", 1, 1, 1, legacytx.StdFee{}, []sdk.Msg{tc.msg}, "memo")))
		})
	}
}
