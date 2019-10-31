package main

import (
	"fmt"
	"github.com/link-chain/link/contrib/provisioning/k8s"
	"github.com/link-chain/link/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/user"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var cfgFile string

func main() {

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(types.Bech32PrefixAccAddr, types.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(types.Bech32PrefixValAddr, types.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(types.Bech32PrefixConsAddr, types.Bech32PrefixConsPub)
	config.Seal()
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "provisioning",
		Short: "It can be used for provisioning LINK BlockChain Nodes to anywhere",
		Long: `It can be used for provisioning LINK BlockChain Nodes to anywhere
	e.g) provisioning  k8s -a build -c link-chain -m 0.000003stake -n 24656 -e 24658 -r 24657 -i 10.231.253.192,10.231.253.193,10.231.253.195,10.231.224.247`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		//	Run: func(cmd *cobra.Command, args []string) { },
	}

	rootCmd.AddCommand(k8s.Init())

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.link-provisioning.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		// Search config in home directory with Name ".link-provisioning" (without extension).
		viper.AddConfigPath(usr.HomeDir)
		viper.SetConfigName(".provisioning")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file: ", viper.ConfigFileUsed())
	}
}
