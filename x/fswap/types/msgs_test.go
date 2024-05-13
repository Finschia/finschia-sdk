package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/auth/legacy/legacytx"
	banktypes "github.com/Finschia/finschia-sdk/x/bank/types"
	fswaptypes "github.com/Finschia/finschia-sdk/x/fswap/types"
)

func TestAminoJSON(t *testing.T) {
	tx := legacytx.StdTx{}

	sender := "link15sdc7wdajsps42fky3j6mnvr4tj9fv6w3hkqkj"

	swapRate, _ := sdk.NewDecFromStr("148.079656")

	testCase := map[string]struct {
		msg          legacytx.LegacyMsg
		expectedType string
		expected     string
	}{
		"MsgSwap": {
			&fswaptypes.MsgSwap{
				FromAddress: sender,
				FromCoinAmount: sdk.Coin{
					Denom:  "cony",
					Amount: sdk.NewInt(100000),
				},
				ToDenom: "pdt",
			},
			"/lbm.fswap.v1.MsgSwap",
			"{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSwap\",\"value\":{\"from_address\":\"link15sdc7wdajsps42fky3j6mnvr4tj9fv6w3hkqkj\",\"from_coin_amount\":{\"amount\":\"100000\",\"denom\":\"cony\"},\"to_denom\":\"pdt\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}",
		},
		"MsgSwapAll": {
			&fswaptypes.MsgSwapAll{
				FromAddress: sender,
				FromDenom:   "cony",
				ToDenom:     "pdt",
			},
			"/lbm.fswap.v1.MsgSwapAll",
			"{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSwapAll\",\"value\":{\"from_address\":\"link15sdc7wdajsps42fky3j6mnvr4tj9fv6w3hkqkj\",\"from_denom\":\"cony\",\"to_denom\":\"pdt\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}",
		},
		"MsgSetSwap": {
			&fswaptypes.MsgSetSwap{
				Authority: sender,
				Swap: fswaptypes.Swap{
					FromDenom:           "cony",
					ToDenom:             "pdt",
					AmountCapForToDenom: sdk.NewInt(1000000000000000),
					SwapRate:            swapRate,
				},
				ToDenomMetadata: banktypes.Metadata{
					Description: "test",
					DenomUnits: []*banktypes.DenomUnit{
						{
							Denom:    "kaia",
							Exponent: 0,
							Aliases:  nil,
						},
					},
					Base:    "kaia",
					Display: "KAIA",
					Name:    "Kaia",
					Symbol:  "KAIA",
				},
			},
			"/lbm.fswap.v1.MsgSetSwap",
			"{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSetSwap\",\"value\":{\"authority\":\"link15sdc7wdajsps42fky3j6mnvr4tj9fv6w3hkqkj\",\"swap\":{\"amount_cap_for_to_denom\":\"1000000000000000\",\"from_denom\":\"cony\",\"swap_rate\":\"148.079656000000000000\",\"to_denom\":\"pdt\"},\"to_denom_metadata\":{\"base\":\"kaia\",\"denom_units\":[{\"denom\":\"kaia\"}],\"description\":\"test\",\"display\":\"KAIA\",\"name\":\"Kaia\",\"symbol\":\"KAIA\"}}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}",
		},
	}

	for name, tc := range testCase {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tx.Msgs = []sdk.Msg{tc.msg}
			require.Equal(t, fswaptypes.RouterKey, tc.msg.Route())
			require.Equal(t, tc.expectedType, tc.msg.Type())
			require.Equal(t, tc.expected, string(legacytx.StdSignBytes("foo", 1, 1, 1, legacytx.StdFee{}, []sdk.Msg{tc.msg}, "memo")))
		})
	}
}
