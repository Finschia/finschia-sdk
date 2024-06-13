package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Finschia/finschia-sdk/simapp"
	"github.com/Finschia/finschia-sdk/testutil/testdata"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

func (s *KeeperTestSuite) TestInitAndExportGenesis() {
	ctx, _ := s.ctx.CacheContext()
	testGenesis := stubGenesisState()
	err := s.keeper.InitGenesis(ctx, testGenesis)
	s.Require().NoError(err)

	exportGenesis := s.keeper.ExportGenesis(ctx)
	fmt.Println(len(exportGenesis.GetSwaps()))
	s.Require().Equal(testGenesis, exportGenesis)
	s.Require().Equal(testGenesis.GetSwaps(), exportGenesis.GetSwaps())
	s.Require().Equal(testGenesis.GetSwapStats(), exportGenesis.GetSwapStats())
	s.Require().Equal(testGenesis.GetSwappeds(), exportGenesis.GetSwappeds())
}

func TestInitGenesis(t *testing.T) {
	checkTx := false
	app := simapp.Setup(checkTx)
	testdata.RegisterInterfaces(app.InterfaceRegistry())
	testdata.RegisterMsgServer(app.MsgServiceRouter(), testdata.MsgServerImpl{})
	ctx := app.BaseApp.NewContext(checkTx, tmproto.Header{})
	keeper := app.FswapKeeper

	tests := []struct {
		name          string
		genState      *types.GenesisState
		expectedError error
	}{
		{
			name:          "valid",
			genState:      stubGenesisState(),
			expectedError: nil,
		},
		{
			name: "invalid: swapCount",
			genState: func() *types.GenesisState {
				state := stubGenesisState()
				state.SwapStats.SwapCount = -1
				return state
			}(),
			expectedError: types.ErrInvalidState,
		},
		{
			name: "invalid: swaps count exceeds limit",
			genState: func() *types.GenesisState {
				state := stubGenesisState()
				state.Swaps = append(state.Swaps, state.Swaps[0])
				state.Swappeds = append(state.Swappeds, state.Swappeds[0])
				state.SwapStats.SwapCount = 2
				return state
			}(),
			expectedError: types.ErrCanNotHaveMoreSwap,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := keeper.InitGenesis(ctx, tc.genState)
			require.ErrorIs(t, tc.expectedError, err)
		})
	}
}

func stubGenesisState() *types.GenesisState {
	testSwapRate, _ := sdk.NewDecFromStr("1234567890")
	return &types.GenesisState{
		Swaps: []types.Swap{
			{
				FromDenom:           "aaa",
				ToDenom:             "bbb",
				AmountCapForToDenom: sdk.NewInt(1234567890000),
				SwapRate:            testSwapRate,
			},
		},
		SwapStats: types.SwapStats{
			SwapCount: 1,
		},
		Swappeds: []types.Swapped{
			{
				FromCoinAmount: sdk.Coin{
					Denom:  "aaa",
					Amount: sdk.ZeroInt(),
				},
				ToCoinAmount: sdk.Coin{
					Denom:  "bbb",
					Amount: sdk.ZeroInt(),
				},
			},
		},
	}
}
