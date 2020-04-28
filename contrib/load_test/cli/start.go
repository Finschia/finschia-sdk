package cli

import (
	"net/http"

	"github.com/line/link/app"
	"github.com/line/link/contrib/load_test/master"
	"github.com/line/link/contrib/load_test/service"
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
		RampUpTime:        viper.GetInt(FlagRampUpTime),
		MaxWorkers:        viper.GetInt(FlagMaxWorkers),
		TargetURL:         viper.GetString(FlagTargetURL),
		ChainID:           viper.GetString(FlagChainID),
		CoinName:          viper.GetString(FlagCoinName),
	}

	c := master.NewController(slaves, config)
	startHeight := getCurrentHeight(config.TargetURL)
	if err := c.StartLoadTest(); err != nil {
		return err
	}
	endHeight := getCurrentHeight(config.TargetURL)

	r := master.NewReporter(c.Results, config, startHeight, endHeight)
	if err := r.Report(outputDir); err != nil {
		return err
	}
	return nil
}

func getCurrentHeight(targetURL string) int64 {
	block, err := service.NewLinkService(&http.Client{}, app.MakeCodec(), targetURL).GetLatestBlock()
	if err != nil {
		panic(err)
	}
	return block.Block.Height
}
