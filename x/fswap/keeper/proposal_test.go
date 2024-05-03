package keeper_test

import (
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
			types.ErrCanNotHaveMoreSwap, // TODO(bjs) check with maxSwap(2)
		},
	}
	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			err := s.keeper.MakeSwap(ctx, tc.swap, s.toDenomMetadata)
			if tc.existingMetadata {
				err := s.keeper.MakeSwap(ctx, tc.swap, s.toDenomMetadata)
				s.Require().ErrorIs(err, tc.expectedError)
			} else {
				s.Require().ErrorIs(err, tc.expectedError)
			}
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
