package keeper_test

import (
	"fmt"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

func (s *KeeperTestSuite) TestInitAndExportGenesis() {
	ctx, _ := s.ctx.CacheContext()
	testSwapRate, _ := sdk.NewDecFromStr("1234567890")
	testGenesis := &types.GenesisState{
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
	err := s.keeper.InitGenesis(ctx, testGenesis)
	s.Require().NoError(err)

	exportGenesis := s.keeper.ExportGenesis(ctx)
	fmt.Println(len(exportGenesis.GetSwaps()))
	s.Require().Equal(testGenesis, exportGenesis)
	s.Require().Equal(testGenesis.GetSwaps(), exportGenesis.GetSwaps())
	s.Require().Equal(testGenesis.GetSwapStats(), exportGenesis.GetSwapStats())
	s.Require().Equal(testGenesis.GetSwappeds(), exportGenesis.GetSwappeds())
}
