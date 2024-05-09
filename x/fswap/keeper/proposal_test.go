package keeper_test

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	bank "github.com/Finschia/finschia-sdk/x/bank/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

func (s *KeeperTestSuite) TestMakeSwapProposal() {
	testCases := map[string]struct {
		swap             types.Swap
		toDenomMeta      bank.Metadata
		existingMetadata bool
		expectedError    error
		expectedEvents   sdk.Events
	}{
		"valid": {
			types.Swap{
				FromDenom:           "fromD",
				ToDenom:             "toD",
				AmountCapForToDenom: sdk.OneInt(),
				SwapRate:            sdk.OneDec(),
			},
			s.toDenomMetadata,
			false,
			nil,
			sdk.Events{
				sdk.Event{
					Type: "lbm.fswap.v1.EventMakeSwap",
					Attributes: []abci.EventAttribute{
						{
							Key:   []byte("swap"),
							Value: []uint8{0x7b, 0x22, 0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x22, 0x3a, 0x22, 0x66, 0x72, 0x6f, 0x6d, 0x44, 0x22, 0x2c, 0x22, 0x74, 0x6f, 0x5f, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x22, 0x3a, 0x22, 0x74, 0x6f, 0x44, 0x22, 0x2c, 0x22, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x63, 0x61, 0x70, 0x5f, 0x66, 0x6f, 0x72, 0x5f, 0x74, 0x6f, 0x5f, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x22, 0x3a, 0x22, 0x31, 0x22, 0x2c, 0x22, 0x73, 0x77, 0x61, 0x70, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x22, 0x3a, 0x22, 0x31, 0x2e, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x22, 0x7d},
							Index: false,
						},
					},
				},
				sdk.Event{
					Type: "lbm.fswap.v1.EventAddDenomMetadata",
					Attributes: []abci.EventAttribute{
						{
							Key:   []byte("metadata"),
							Value: []uint8{0x7b, 0x22, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x3a, 0x22, 0x54, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x20, 0x66, 0x6f, 0x72, 0x20, 0x74, 0x6f, 0x2d, 0x63, 0x6f, 0x69, 0x6e, 0x22, 0x2c, 0x22, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x5f, 0x75, 0x6e, 0x69, 0x74, 0x73, 0x22, 0x3a, 0x5b, 0x7b, 0x22, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x22, 0x3a, 0x22, 0x74, 0x6f, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x22, 0x2c, 0x22, 0x65, 0x78, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x22, 0x3a, 0x30, 0x2c, 0x22, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x65, 0x73, 0x22, 0x3a, 0x5b, 0x5d, 0x7d, 0x5d, 0x2c, 0x22, 0x62, 0x61, 0x73, 0x65, 0x22, 0x3a, 0x22, 0x64, 0x75, 0x6d, 0x6d, 0x79, 0x22, 0x2c, 0x22, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x22, 0x3a, 0x22, 0x64, 0x75, 0x6d, 0x6d, 0x79, 0x63, 0x6f, 0x69, 0x6e, 0x22, 0x2c, 0x22, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x3a, 0x22, 0x44, 0x55, 0x4d, 0x4d, 0x59, 0x22, 0x2c, 0x22, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x22, 0x3a, 0x22, 0x44, 0x55, 0x4d, 0x22, 0x7d},
							Index: false,
						},
					},
				},
			},
		},
		"to-denom metadata change not allowed": {
			types.Swap{
				FromDenom:           "fromD",
				ToDenom:             "toD",
				AmountCapForToDenom: sdk.OneInt(),
				SwapRate:            sdk.OneDec(),
			},
			bank.Metadata{
				Description: s.toDenomMetadata.Description,
				DenomUnits:  s.toDenomMetadata.DenomUnits,
				Base:        "change",
				Display:     s.toDenomMetadata.Display,
				Name:        s.toDenomMetadata.Name,
				Symbol:      s.toDenomMetadata.Symbol,
			},
			true,
			sdkerrors.ErrInvalidRequest,
			sdk.Events{},
		},
	}
	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			err := s.keeper.MakeSwap(ctx, tc.swap, s.toDenomMetadata)
			if tc.existingMetadata {
				err := s.keeper.MakeSwap(ctx, tc.swap, s.toDenomMetadata)
				s.Require().ErrorIs(err, tc.expectedError)
				return
			}

			s.Require().ErrorIs(err, tc.expectedError)
			events := ctx.EventManager().Events()
			s.Require().Equal(tc.expectedEvents, events)
		})
	}
}

func (s *KeeperTestSuite) TestSwapValidateBasic() {
	testCases := map[string]struct {
		swap             types.Swap
		shouldThrowError bool
		expectedError    error
	}{
		"valid": {
			types.Swap{
				FromDenom:           "fromD",
				ToDenom:             "toD",
				AmountCapForToDenom: sdk.OneInt(),
				SwapRate:            sdk.OneDec(),
			},
			false,
			nil,
		},
		"invalid empty from-denom": {
			types.Swap{
				FromDenom:           "",
				ToDenom:             "toD",
				AmountCapForToDenom: sdk.OneInt(),
				SwapRate:            sdk.OneDec(),
			},
			true,
			sdkerrors.ErrInvalidRequest,
		},
		"invalid empty to-denom": {
			types.Swap{
				FromDenom:           "fromD",
				ToDenom:             "",
				AmountCapForToDenom: sdk.OneInt(),
				SwapRate:            sdk.OneDec(),
			},
			true,
			sdkerrors.ErrInvalidRequest,
		},
		"invalid zero amount cap for to-denom": {
			types.Swap{
				FromDenom:           "fromD",
				ToDenom:             "toD",
				AmountCapForToDenom: sdk.ZeroInt(),
				SwapRate:            sdk.OneDec(),
			},
			true,
			sdkerrors.ErrInvalidRequest,
		},
		"invalid zero swap-rate": {
			types.Swap{
				FromDenom:           "fromD",
				ToDenom:             "toD",
				AmountCapForToDenom: sdk.OneInt(),
				SwapRate:            sdk.ZeroDec(),
			},
			true,
			sdkerrors.ErrInvalidRequest,
		},
		"invalid the same from-denom and to-denom": {
			types.Swap{
				FromDenom:           "same",
				ToDenom:             "same",
				AmountCapForToDenom: sdk.OneInt(),
				SwapRate:            sdk.OneDec(),
			},
			true,
			sdkerrors.ErrInvalidRequest,
		},
	}
	for name, tc := range testCases {
		s.Run(name, func() {
			err := tc.swap.ValidateBasic()
			if tc.shouldThrowError {
				s.Require().ErrorIs(err, tc.expectedError)
				return
			}
			s.Require().NoError(err)
		})
	}
}
