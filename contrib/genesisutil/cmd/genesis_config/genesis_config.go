package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/staking"
	tmconfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"

	"github.com/link-chain/link/app"
)

const (
	flagHome      = "home"
	flagOverwrite = "overwrite"
)

type GenesisConfig struct {
	Denom string `toml:"denom"`
}

func (gc GenesisConfig) ChangeDenomOf(v interface{}, path string) (cnt int) {
	denom := gc.Denom

	elem := reflect.ValueOf(v).Elem()
	intf := elem.Interface()

	switch elem.Kind() {
	case reflect.String:
		if intf != denom {
			log.Printf("%vdenom: %v -> %v\n", path, intf, denom)
			elem.SetString(denom)
			cnt += 1
		}
	case reflect.Slice:
		for i := 0; i < elem.Len(); i++ {
			item := elem.Index(i)
			if item.Kind() == reflect.Struct {
				field := item.FieldByName("Denom")
				if field.Interface() != denom {
					log.Printf("%v[%v].denom : %v -> %v\n", path, i, field, denom)
					field.SetString(denom)
					cnt += 1
				}
			}
		}
	}

	return cnt
}

func LoadGenesisConfig(configFile string) (config GenesisConfig, err error) {
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("Fatal error config file: %s \n", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("Fatal error config file: %s \n", err)
	}

	return config, nil
}

func main() {
	cobra.EnableCommandSorting = false

	var cmd = &cobra.Command{
		Use:   "genesis_config [genesis.toml]",
		Short: "Configure genesis.json with genesis.toml",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			tmConfig := tmconfig.DefaultBaseConfig()
			tmConfig.RootDir = viper.GetString(flagHome)

			genFile := tmConfig.GenesisFile()
			if !common.FileExists(genFile) {
				return fmt.Errorf("genesis.json file doesn't exist: %v", genFile)
			}

			genDoc := &types.GenesisDoc{}
			if _, err := os.Stat(genFile); err != nil {
				if !os.IsNotExist(err) {
					return err
				}
			} else {
				genDoc, err = types.GenesisDocFromFile(genFile)
				if err != nil {
					return err
				}
			}

			configPath := args[0]

			config, err := LoadGenesisConfig(configPath)
			if err != nil {
				return err
			}

			cdc := app.MakeCodec()

			var appState map[string]json.RawMessage
			cdc.MustUnmarshalJSON(genDoc.AppState, &appState)

			changed := 0

			for name, module := range appState {
				switch name {
				case crisis.ModuleName:
					var state crisis.GenesisState
					crisis.ModuleCdc.MustUnmarshalJSON(module, &state)

					denomPath := "crisis.constant_fee."
					if res := config.ChangeDenomOf(
						&state.ConstantFee.Denom, denomPath); res > 0 {
						changed += res
					}

					appState[name] = crisis.ModuleCdc.MustMarshalJSON(&state)
				case gov.ModuleName:
					var state gov.GenesisState
					gov.ModuleCdc.MustUnmarshalJSON(module, &state)

					denomPath := "gov.deposit_params.min_deposit"
					if res := config.ChangeDenomOf(
						&state.DepositParams.MinDeposit, denomPath); res > 0 {
						changed += res
					}

					appState[name] = gov.ModuleCdc.MustMarshalJSON(&state)
				case mint.ModuleName:
					var state mint.GenesisState
					mint.ModuleCdc.MustUnmarshalJSON(module, &state)

					denomPath := "mint.params.mint_"
					if res := config.ChangeDenomOf(
						&state.Params.MintDenom, denomPath); res > 0 {
						changed += res
					}

					appState[name] = mint.ModuleCdc.MustMarshalJSON(&state)
				case staking.ModuleName:
					var state staking.GenesisState
					staking.ModuleCdc.MustUnmarshalJSON(module, &state)

					denomPath := "staking.params.bond_"
					if res := config.ChangeDenomOf(
						&state.Params.BondDenom, denomPath); res > 0 {
						changed += res
					}

					appState[name] = staking.ModuleCdc.MustMarshalJSON(&state)
				}
			}

			genDoc.AppState = codec.MustMarshalJSONIndent(cdc, &appState)

			if viper.GetBool(flagOverwrite) {
				log.Printf("overwrite %v with %v changes\n", genFile, changed)

				if err = genutil.ExportGenesisFile(genDoc, genFile); err != nil {
					return err
				}
			} else {
				fmt.Println(string(genDoc.AppState))
			}

			return nil
		},
	}

	cmd.Flags().String(flagHome, app.DefaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		panic(err)
	}

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
