package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-sdk/crypto/keys/secp256k1"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/bank/types"
)

func TestMsgExecutionValidateBasic(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	const proofStr = "{\n \"pi_a\": [\n  \"7575287679446209007446416020137456670042570578978230730578011103770415897062\",\n  \"20469978368515629364541212704109752583692706286549284712208570249653184893207\",\n  \"1\"\n ],\n \"pi_b\": [\n  [\n   \"4001119070037193619600086014535210556571209449080681376392853276923728808564\",\n   \"18475391841797083641468254159150812922259839776046448499150732610021959794558\"\n  ],\n  [\n   \"19781252109528278034156073207688818205850783935629584279449144780221040670063\",\n   \"5873714313814830719712095806732872482213125567325442209795797618441438990229\"\n  ],\n  [\n   \"1\",\n   \"0\"\n  ]\n ],\n \"pi_c\": [\n  \"18920522434978516095250248740518039198650690968720755259416280639852277665022\",\n  \"1945774583580804632084048753815901730674007769630810705050114062476636502591\",\n  \"1\"\n ],\n \"protocol\": \"groth16\",\n \"curve\": \"bn128\"\n}"
	validZKAuthInputs := &ZKAuthInputs{
		ProofPoints:  []byte(proofStr),
		IssBase64:    "aHR0cHM6Ly9hY2NvdW50cy5nb29nbGUuY29t",
		HeaderBase64: "eyJhbGciOiJSUzI1NiIsImtpZCI6IjU1YzE4OGE4MzU0NmZjMTg4ZTUxNTc2YmE3MjgzNmUwNjAwZThiNzMiLCJ0eXAiOiJKV1QifQ",
		AddressSeed:  "15035161560159971633800983619931498696152633426768016966057770643262022096073",
	}

	type TestMsg struct {
		addr string
	}

	testcase := map[string]struct {
		msgs            []sdk.Msg
		zkAuthSignature ZKAuthSignature
		valid           bool
	}{
		"valid msg": {
			[]sdk.Msg{
				&types.MsgSend{
					FromAddress: addrs[0].String(),
					ToAddress:   addrs[1].String(),
					Amount:      sdk.NewCoins(sdk.NewInt64Coin("cony", 1000)),
				},
			},
			ZKAuthSignature{
				ZkAuthInputs:   validZKAuthInputs,
				MaxBlockHeight: 32754,
			},
			true,
		},
		"no msg": {
			nil,
			ZKAuthSignature{
				ZkAuthInputs:   validZKAuthInputs,
				MaxBlockHeight: 32754,
			},
			false,
		},
		"empty signature": {
			[]sdk.Msg{
				&types.MsgSend{
					FromAddress: addrs[0].String(),
					ToAddress:   addrs[1].String(),
					Amount:      sdk.NewCoins(sdk.NewInt64Coin("cony", 1000)),
				},
			},
			ZKAuthSignature{},
			false,
		},
		"max block height is zero": {
			[]sdk.Msg{
				&types.MsgSend{
					FromAddress: addrs[0].String(),
					ToAddress:   addrs[1].String(),
					Amount:      sdk.NewCoins(sdk.NewInt64Coin("cony", 1000)),
				},
			},
			ZKAuthSignature{
				ZkAuthInputs:   validZKAuthInputs,
				MaxBlockHeight: 0,
			},
			false,
		},
	}

	for name, tc := range testcase {
		t.Run(name, func(t *testing.T) {
			msg := MsgExecution{
				ZkAuthSignature: tc.zkAuthSignature,
			}
			err := msg.SetMsgs(tc.msgs)
			require.NoError(t, err)

			err = msg.ValidateBasic()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}

}
