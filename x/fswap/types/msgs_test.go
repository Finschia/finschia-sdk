package types_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-sdk/crypto/keys/secp256k1"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
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
					Display: "kaia",
					Name:    "Kaia",
					Symbol:  "KAIA",
				},
			},
			"/lbm.fswap.v1.MsgSetSwap",
			"{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSetSwap\",\"value\":{\"authority\":\"link15sdc7wdajsps42fky3j6mnvr4tj9fv6w3hkqkj\",\"swap\":{\"amount_cap_for_to_denom\":\"1000000000000000\",\"from_denom\":\"cony\",\"swap_rate\":\"148.079656000000000000\",\"to_denom\":\"pdt\"},\"to_denom_metadata\":{\"base\":\"kaia\",\"denom_units\":[{\"denom\":\"kaia\"}],\"description\":\"test\",\"display\":\"kaia\",\"name\":\"Kaia\",\"symbol\":\"KAIA\"}}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}",
		},
	}

	for name, tc := range testCase {
		t.Run(name, func(t *testing.T) {
			tx.Msgs = []sdk.Msg{tc.msg}
			require.Equal(t, fswaptypes.RouterKey, tc.msg.Route())
			require.Equal(t, tc.expectedType, tc.msg.Type())
			require.Equal(t, tc.expected, string(legacytx.StdSignBytes("foo", 1, 1, 1, legacytx.StdFee{}, []sdk.Msg{tc.msg}, "memo")))
		})
	}
}

func TestMsgSwapValidate(t *testing.T) {
	pk := secp256k1.GenPrivKey().PubKey()
	address, err := sdk.Bech32ifyAddressBytes("link", pk.Address())
	if err != nil {
		return
	}
	tests := []struct {
		name          string
		msg           *fswaptypes.MsgSwap
		expectedError error
	}{
		{
			name: "valid",
			msg: &fswaptypes.MsgSwap{
				FromAddress:    address,
				FromCoinAmount: sdk.NewCoin("fromDenom", sdk.OneInt()),
				ToDenom:        "kei",
			},
			expectedError: nil,
		},
		{
			name: "invalid: address",
			msg: &fswaptypes.MsgSwap{
				FromAddress:    "invalid-address",
				FromCoinAmount: sdk.NewCoin("fromDenom", sdk.OneInt()),
				ToDenom:        "kei",
			},
			expectedError: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid: FromCoinAmount empty denom",
			msg: &fswaptypes.MsgSwap{
				FromAddress: address,
				FromCoinAmount: sdk.Coin{
					Denom:  "",
					Amount: sdk.OneInt(),
				},
				ToDenom: "kei",
			},
			expectedError: sdkerrors.ErrInvalidCoins,
		},
		{
			name: "invalid: FromCoinAmount zero amount",
			msg: &fswaptypes.MsgSwap{
				FromAddress: address,
				FromCoinAmount: sdk.Coin{
					Denom:  "cony",
					Amount: sdk.ZeroInt(),
				},
				ToDenom: "kei",
			},
			expectedError: sdkerrors.ErrInvalidCoins,
		},
		{
			name: "invalid: ToDenom",
			msg: &fswaptypes.MsgSwap{
				FromAddress:    address,
				FromCoinAmount: sdk.NewCoin("fromDenom", sdk.OneInt()),
				ToDenom:        "",
			},
			expectedError: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			require.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestMsgSwapAllValidate(t *testing.T) {
	pk := secp256k1.GenPrivKey().PubKey()
	address, err := sdk.Bech32ifyAddressBytes("link", pk.Address())
	if err != nil {
		return
	}
	tests := []struct {
		name          string
		msg           *fswaptypes.MsgSwapAll
		expectedError error
	}{
		{
			name: "valid",
			msg: &fswaptypes.MsgSwapAll{
				FromAddress: address,
				FromDenom:   "cony",
				ToDenom:     "kei",
			},
			expectedError: nil,
		},
		{
			name: "invalid: address",
			msg: &fswaptypes.MsgSwapAll{
				FromAddress: "invalid-address",
				FromDenom:   "cony",
				ToDenom:     "kei",
			},
			expectedError: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid: FromDenom",
			msg: &fswaptypes.MsgSwapAll{
				FromAddress: address,
				FromDenom:   "",
				ToDenom:     "kei",
			},
			expectedError: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid: ToDenom",
			msg: &fswaptypes.MsgSwapAll{
				FromAddress: address,
				FromDenom:   "cony",
				ToDenom:     "",
			},
			expectedError: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			require.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestMsgSetSwapValidate(t *testing.T) {
	pk := secp256k1.GenPrivKey().PubKey()
	address, err := sdk.Bech32ifyAddressBytes("link", pk.Address())
	if err != nil {
		return
	}
	tests := []struct {
		name          string
		msg           *fswaptypes.MsgSetSwap
		expectedError error
	}{
		{
			name: "valid",
			msg: &fswaptypes.MsgSetSwap{
				Authority: address,
				Swap: fswaptypes.Swap{
					FromDenom:           "cony",
					ToDenom:             "kei",
					AmountCapForToDenom: sdk.OneInt(),
					SwapRate:            sdk.NewDec(123),
				},
				ToDenomMetadata: banktypes.Metadata{
					Description: "desc",
					DenomUnits: []*banktypes.DenomUnit{
						{
							Denom:    "kei",
							Exponent: 0,
							Aliases:  nil,
						},
					},
					Base:    "kei",
					Display: "kei",
					Name:    "kei",
					Symbol:  "KAIA",
				},
			},
			expectedError: nil,
		},
		{
			name: "invalid: address",
			msg: &fswaptypes.MsgSetSwap{
				Authority: "invalid-address",
				Swap: fswaptypes.Swap{
					FromDenom:           "cony",
					ToDenom:             "kei",
					AmountCapForToDenom: sdk.OneInt(),
					SwapRate:            sdk.NewDec(123),
				},
				ToDenomMetadata: banktypes.Metadata{
					Description: "desc",
					DenomUnits: []*banktypes.DenomUnit{
						{
							Denom:    "kei",
							Exponent: 0,
							Aliases:  nil,
						},
					},
					Base:    "kei",
					Display: "kei",
					Name:    "kei",
					Symbol:  "KAIA",
				},
			},
			expectedError: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid: Swap",
			msg: &fswaptypes.MsgSetSwap{
				Authority: address,
				Swap:      fswaptypes.Swap{},
				ToDenomMetadata: banktypes.Metadata{
					Description: "desc",
					DenomUnits: []*banktypes.DenomUnit{
						{
							Denom:    "kei",
							Exponent: 0,
							Aliases:  nil,
						},
					},
					Base:    "kei",
					Display: "kei",
					Name:    "kei",
					Symbol:  "KAIA",
				},
			},
			expectedError: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid: ToDenomMetadata",
			msg: &fswaptypes.MsgSetSwap{
				Authority: address,
				Swap: fswaptypes.Swap{
					FromDenom:           "cony",
					ToDenom:             "kei",
					AmountCapForToDenom: sdk.OneInt(),
					SwapRate:            sdk.NewDec(123),
				},
				ToDenomMetadata: banktypes.Metadata{
					Description: "",
					DenomUnits: []*banktypes.DenomUnit{
						{
							Denom:    "kei",
							Exponent: 0,
							Aliases:  nil,
						},
					},
				},
			},
			expectedError: errors.New("name field cannot be blank"),
		},

		{
			name: "invalid: mismatched toDenom",
			msg: &fswaptypes.MsgSetSwap{
				Authority: address,
				Swap: fswaptypes.Swap{
					FromDenom:           "cony",
					ToDenom:             "kei",
					AmountCapForToDenom: sdk.OneInt(),
					SwapRate:            sdk.NewDec(123),
				},
				ToDenomMetadata: banktypes.Metadata{
					Description: "desc",
					DenomUnits: []*banktypes.DenomUnit{
						{
							Denom:    "fkei",
							Exponent: 0,
							Aliases:  nil,
						},
					},
					Base:    "fkei",
					Display: "fkei",
					Name:    "fkei",
					Symbol:  "KAIA",
				},
			},
			expectedError: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectedError != nil {
				require.Contains(t, err.Error(), tc.expectedError.Error())
			}
		})
	}
}
