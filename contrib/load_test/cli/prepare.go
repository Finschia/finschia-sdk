package cli

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	outputDir := viper.GetString(FlagOutputDir)
	masterMnemonic := viper.GetString(FlagMasterMnemonic)
	slaves := getSlaves()
	config := types.Config{
		MsgsPerTxPrepare:  viper.GetInt(FlagMsgsPerTxPrepare),
		MaxGasPrepare:     viper.GetInt(FlagMaxGasPrepare),
		MsgsPerTxLoadTest: viper.GetInt(FlagMsgsPerTxLoadTest),
		MaxGasLoadTest:    viper.GetInt(FlagMaxGasLoadTest),
		TPS:               viper.GetInt(FlagTPS),
		Duration:          viper.GetInt(FlagDuration),
		TargetURL:         viper.GetString(FlagPrepareTargetURL),
		ChainID:           viper.GetString(FlagChainID),
		CoinName:          viper.GetString(FlagCoinName),
		Testnet:           viper.GetBool(FlagTestnet),
	}
	types.SetBech32Prefix(sdk.GetConfig(), config.Testnet)

	ss, err := master.NewStateSetter(masterMnemonic, config)
	if err != nil {
		return err
	}
	if err := ss.PrepareTestState(slaves, outputDir); err != nil {
		return err
	}

	return nil
}
