package server

// DONTCOVER

import (
	"fmt"

	"github.com/spf13/cobra"
	cfg "github.com/tendermint/tendermint/config"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmos "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	pvm "github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"
	tmversion "github.com/tendermint/tendermint/version"
	yaml "gopkg.in/yaml.v2"

	sdk "github.com/Finschia/finschia-sdk/types"
)

// ShowNodeIDCmd - ported from tendermint, dump node ID to stdout
func ShowNodeIDCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show-node-id",
		Short: "Show this node's ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := GetServerContextFromCmd(cmd)
			cfg := serverCtx.Config

			nodeKey, err := p2p.LoadNodeKey(cfg.NodeKeyFile())
			if err != nil {
				return err
			}
			fmt.Println(nodeKey.ID())
			return nil
		},
	}
}

// ShowValidatorCmd - ported from tendermint, show this node's validator info
func ShowValidatorCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "show-validator",
		Short: "Show this node's tendermint validator info",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			return showValidator(cmd, config)
		},
	}

	return &cmd
}

func showValidator(_ *cobra.Command, config *cfg.Config) error {
	var pv types.PrivValidator
	//nolint:gocritic
	// todo: need to fix it when node.createAndStartPrivValidatorSocketClient change to public function in cometbft
	// This changes is applied in finschia-sdk.https://github.com/Finschia/finschia-sdk/pull/821
	//if config.PrivValidatorListenAddr != "" {
	//	chainID, err := loadChainID(config)
	//	if err != nil {
	//		return err
	//	}
	//	serverCtx := GetServerContextFromCmd(cmd)
	//	log := serverCtx.Logger
	//	node.
	//		pv, err = node.CreateAndStartPrivValidatorSocketClient(config, chainID, log)
	//	if err != nil {
	//		return err
	//	}
	//} else {
	keyFilePath := config.PrivValidatorKeyFile()
	if !tmos.FileExists(keyFilePath) {
		return fmt.Errorf("private validator file %s does not exist", keyFilePath)
	}
	pv = pvm.LoadFilePV(keyFilePath, config.PrivValidatorStateFile())
	//}

	pubKey, err := pv.GetPubKey()
	if err != nil {
		return fmt.Errorf("can't get pubkey: %w", err)
	}

	bz, err := tmjson.Marshal(pubKey)
	if err != nil {
		return fmt.Errorf("failed to marshal private validator pubkey: %w", err)
	}

	fmt.Println(string(bz))
	return nil
}

func loadChainID(config *cfg.Config) (string, error) {
	stateDB, err := node.DefaultDBProvider(&node.DBContext{ID: "state", Config: config})
	if err != nil {
		return "", err
	}
	defer func() {
		_ = stateDB.Close()
	}()
	genesisDocProvider := node.DefaultGenesisDocProviderFunc(config)
	_, genDoc, err := node.LoadStateFromDBOrGenesisDocProvider(stateDB, genesisDocProvider)
	if err != nil {
		return "", err
	}
	return genDoc.ChainID, nil
}

// ShowAddressCmd - show this node's validator address
func ShowAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-address",
		Short: "Shows this node's tendermint validator consensus address",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := GetServerContextFromCmd(cmd)
			cfg := serverCtx.Config

			privValidator := pvm.LoadFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())
			valConsAddr := (sdk.ConsAddress)(privValidator.GetAddress())
			fmt.Println(valConsAddr.String())
			return nil
		},
	}

	return cmd
}

// VersionCmd prints tendermint and ABCI version numbers.
func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print tendermint libraries' version",
		Long: `Print protocols' and libraries' version numbers
against which this app has been compiled.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			bs, err := yaml.Marshal(&struct {
				Tendermint    string
				ABCI          string
				BlockProtocol uint64
				P2PProtocol   uint64
			}{
				Tendermint:    tmversion.TMCoreSemVer,
				ABCI:          tmversion.ABCIVersion,
				BlockProtocol: tmversion.BlockProtocol,
				P2PProtocol:   tmversion.P2PProtocol,
			})
			if err != nil {
				return err
			}

			fmt.Println(string(bs))
			return nil
		},
	}
}
