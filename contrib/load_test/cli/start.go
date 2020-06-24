package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
		MaxGasLoadTest:    viper.GetInt(FlagMaxGasLoadTest),
		TPS:               viper.GetInt(FlagTPS),
		Duration:          viper.GetInt(FlagDuration),
		RampUpTime:        viper.GetInt(FlagRampUpTime),
		MaxWorkers:        viper.GetInt(FlagMaxWorkers),
		TargetURL:         viper.GetString(FlagTargetURL),
		ChainID:           viper.GetString(FlagChainID),
		CoinName:          viper.GetString(FlagCoinName),
		Testnet:           viper.GetBool(FlagTestnet),
	}
	types.SetBech32Prefix(sdk.GetConfig(), config.Testnet)
	params, err := readParams(outputDir)
	if err != nil {
		return err
	}

	c := master.NewController(slaves, config, params)
	if err := c.StartLoadTest(); err != nil {
		return err
	}

	if outputDir != "" {
		if err := writeResultData(c, outputDir); err != nil {
			return err
		}
	}
	return nil
}

func readParams(outputDir string) (map[string]map[string]string, error) {
	paramsFilePath := fmt.Sprintf("%s/%s", outputDir, master.ParamsFileName)
	if _, err := os.Stat(paramsFilePath); os.IsNotExist(err) {
		return nil, nil
	}
	log.Println("Load file of test parameters")
	data, err := ioutil.ReadFile(paramsFilePath)
	if err != nil {
		return nil, err
	}
	var params map[string]map[string]string
	if err := json.Unmarshal(data, &params); err != nil {
		return nil, err
	}
	return params, nil
}

func writeResultData(c *master.Controller, outputDir string) error {
	log.Println("Write result data of load test")
	file, err := json.Marshal(c)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(fmt.Sprintf("%s/%s", outputDir, "result_data.json"), file, 0600); err != nil {
		return err
	}
	return nil
}
