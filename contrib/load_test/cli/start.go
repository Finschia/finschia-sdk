package cli

import (
	"github.com/line/link/contrib/load_test/master"
	"github.com/line/link/contrib/load_test/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start load test",
		RunE:  start,
	}
	return cmd
}

func start(cmd *cobra.Command, args []string) error {
	outputDir := viper.GetString(FlagOutputDir)
	slaves := getSlaves()
	config := types.Config{
		MsgsPerTxLoadTest: viper.GetInt(FlagMsgsPerTxLoadTest),
		TPS:               viper.GetInt(FlagTPS),
		Duration:          viper.GetInt(FlagDuration),
		MaxWorkers:        viper.GetInt(FlagMaxWorkers),
		PacerType:         viper.GetString(FlagPacerType),
		TargetURL:         viper.GetString(FlagTargetURL),
		ChainID:           viper.GetString(FlagChainID),
		CoinName:          viper.GetString(FlagCoinName),
	}

	c := master.NewController(slaves, config)
	if err := c.StartLoadTest(); err != nil {
		return err
	}

	r := master.NewReporter(c.Results)
	if outputDir != "" {
		if err := r.WriteResult(outputDir); err != nil {
			return err
		}
	}
	return nil
}
