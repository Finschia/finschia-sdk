package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/master"
	"github.com/line/link/contrib/load_test/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ReportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Get report of load test",
		RunE:  report,
	}
	return cmd
}

func report(cmd *cobra.Command, args []string) error {
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
		Testnet:           viper.GetBool(FlagTestnet),
	}
	thresholds := types.Thresholds{
		Latency: time.Duration(viper.GetInt(FlagLatencyThreshold)) * time.Millisecond,
		TPS:     viper.GetInt(FlagTPSThreshold),
	}
	types.SetBech32Prefix(sdk.GetConfig(), config.Testnet)

	data, err := readResultData(outputDir)
	if err != nil {
		return err
	}

	r := master.NewReporter(data.Results, slaves, config, data.StartHeight, data.EndHeight, thresholds)
	if err := r.Report(outputDir); err != nil {
		return err
	}
	return nil
}

func readResultData(outputDir string) (*master.Controller, error) {
	if outputDir == "" {
		return nil, nil
	}
	log.Println("Read result data of load test")
	filePath := fmt.Sprintf("%s/%s", outputDir, "result_data.json")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var resultData *master.Controller
	if err := json.Unmarshal(data, &resultData); err != nil {
		return nil, err
	}
	return resultData, nil
}
