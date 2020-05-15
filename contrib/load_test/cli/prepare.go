package cli

import (
	"github.com/line/link/contrib/load_test/master"
	"github.com/line/link/contrib/load_test/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func PrepareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prepare",
		Short: "Prepare test state",
		RunE:  prepare,
	}
	return cmd
}

func prepare(cmd *cobra.Command, args []string) error {
	masterMnemonic := viper.GetString(FlagMasterMnemonic)
	slaves := getSlaves()
	config := types.Config{
		MsgsPerTxPrepare: viper.GetInt(FlagMsgsPerTxPrepare),
		TPS:              viper.GetInt(FlagTPS),
		Duration:         viper.GetInt(FlagDuration),
		TargetURL:        viper.GetString(FlagTargetURL),
		ChainID:          viper.GetString(FlagChainID),
		CoinName:         viper.GetString(FlagCoinName),
	}

	ss, err := master.NewStateSetter(masterMnemonic, config)
	if err != nil {
		return err
	}
	if err := ss.PrepareTestState(slaves); err != nil {
		return err
	}

	return nil
}
