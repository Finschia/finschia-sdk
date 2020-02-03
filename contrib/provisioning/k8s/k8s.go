package k8s

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/line/link/app"
	"github.com/spf13/cobra"
	tmconfig "github.com/tendermint/tendermint/config"
)

var (
	flagAction                   = "Action"
	flagChainID                  = "ChainID"
	flagCliBinDirName            = "CliBinDirName"
	flagConfDirName              = "ConfDirName"
	flagConfHomePath             = "ConfHomePath"
	flagDBDir                    = "DBDir"
	flagK8STemplateFilePath      = "K8STemplateFilePath"
	flagFilebeatTemplateFilePath = "filebeatTemplateFilePath"
	flagLinkdBinDirName          = "LinkdBinDirName"
	flagLinkDockerImageUrl       = "LinkDockerImageUrl"
	flagNodeABCIPort             = "NodeABCIPort"
	flagNodeIPs                  = "NodeIPs"
	flagNodeP2PPort              = "NodeP2PPort"
	flagNodeRestAPIPort          = "NodeRestAPIPort"
	flagPrometheusListenPort     = "PrometheusListenPort"
	flagPrometheusTurnOn         = "PrometheusTurnOn"
)

const defConfDirName = "/config"
const defConfigurationFileExt = ".json"
const defConsensusTimeoutCommit = 5
const defDBDir = "data"
const defK8STemplateFilePath = "./contrib/provisioning/k8s/deploy-validator-template.yaml"
const defFilebeatTemplateFilePath = "./contrib/provisioning/k8s/filebeat-validator-template.yaml"
const defLinkDockerImageUrl = "docker-registry.linecorp.com/link-network/v2/linkdnode:latest"
const defMinGasPrices = 0.000006
const defNodeABCIPort = 25658
const defNodeP2PPort = 25656
const defNodeRestAPIPort = 25657
const defOutputDir = "./build"
const defPrivateKeySeedFileName = "key_seed"
const defProfilingPort = 25660
const defPrometheusListenPort = 25661
const defPrometheusTurnOn = true
const defTxIndexIndexAllTags = false
const genTxsDefaultDir = "gentxs"
const listenAllIngressPortTemplate = "tcp://0.0.0.0:%d"
const listenLoopbackIngressPortTemplate = "tcp://127.0.0.1:%d"
const nodeDirPerm = 0755
const nodeDirPrefix = "node"
const prefixABCIPort = "abci-"
const prefixForChainId = "k8s-chain"
const prefixForP2PPort = "p2p-"
const prefixPortRestAPIPort = "rpc-"

var availableActions = map[string]bool{
	"build": true,
}

func Init() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "k8s",
		Short: "generate configurations of link for provisioning on K8S",
		RunE: func(cmd *cobra.Command, args []string) error {
			var action = ErrCheckedStrParam(cmd.Flags().GetString(flagAction))
			if availableActions[action] != true {
				return fmt.Errorf("requires a valid action Name %s", action)
			}
			chainID := ErrCheckedStrParam(cmd.Flags().GetString(flagChainID))
			linkdDir := ErrCheckedStrParam(cmd.Flags().GetString(flagLinkdBinDirName))
			linkCliDir := ErrCheckedStrParam(cmd.Flags().GetString(flagCliBinDirName))
			prometheusTurnOn := ErrCheckedBoolParam(cmd.Flags().GetBool(flagPrometheusTurnOn))
			nodeP2PPort := ErrCheckedIntParam(cmd.Flags().GetInt(flagNodeP2PPort))
			nodeABCIPort := ErrCheckedIntParam(cmd.Flags().GetInt(flagNodeABCIPort))
			nodeRestAPIPort := ErrCheckedIntParam(cmd.Flags().GetInt(flagNodeRestAPIPort))
			prometheusListenPort := ErrCheckedIntParam(cmd.Flags().GetInt(flagPrometheusListenPort))
			minGasPrices := ErrCheckedStrParam(cmd.Flags().GetString(server.FlagMinGasPrices))
			nodes := ErrCheckedStrArrParam(cmd.Flags().GetStringSlice(flagNodeIPs))
			confHomePath := ErrCheckedStrParam(cmd.Flags().GetString(flagConfHomePath))
			confDirName := ErrCheckedStrParam(cmd.Flags().GetString(flagConfDirName))
			k8STemplateFilePath := ErrCheckedStrParam(cmd.Flags().GetString(flagK8STemplateFilePath))
			filebeatTemplateFilePath := ErrCheckedStrParam(cmd.Flags().GetString(flagFilebeatTemplateFilePath))
			linkDockerImageUrl := ErrCheckedStrParam(cmd.Flags().GetString(flagLinkDockerImageUrl))
			dbDir := ErrCheckedStrParam(cmd.Flags().GetString(flagDBDir))

			tmConfig := server.NewDefaultContext().Config
			tmConfig.Instrumentation.Prometheus = prometheusTurnOn
			tmConfig.BaseConfig.ProfListenAddress = fmt.Sprintf("localhost:%d", defProfilingPort)
			tmConfig.Consensus.TimeoutCommit = defConsensusTimeoutCommit * time.Second
			tmConfig.TxIndex.IndexAllTags = defTxIndexIndexAllTags
			DefIfEmpty(&tmConfig.DBPath, defDBDir, dbDir)
			prometheusListenPort = DefFormatSetIfLTEZero(&tmConfig.Instrumentation.PrometheusListenAddr, listenLoopbackIngressPortTemplate, defPrometheusListenPort, prometheusListenPort)
			nodeP2PPort = DefFormatSetIfLTEZero(&tmConfig.P2P.ListenAddress, listenAllIngressPortTemplate, defNodeP2PPort, nodeP2PPort)
			nodeRestAPIPort = DefFormatSetIfLTEZero(&tmConfig.RPC.ListenAddress, listenAllIngressPortTemplate, defNodeRestAPIPort, nodeRestAPIPort)
			nodeABCIPort = DefFormatSetIfLTEZero(&tmConfig.BaseConfig.ProxyApp, listenLoopbackIngressPortTemplate, defNodeABCIPort, nodeABCIPort)

			hash, err := RandomHash()
			if err != nil {
				panic(err)
			}
			DefIfEmpty(&chainID, fmt.Sprintf("%s-%s-%s-%s-%s", prefixForChainId,
				prefixForP2PPort+strconv.Itoa(nodeP2PPort), prefixPortRestAPIPort+strconv.Itoa(nodeRestAPIPort),
				prefixABCIPort+strconv.Itoa(nodeABCIPort), hex.EncodeToString(hash.Sum(nil)))[:50], chainID)
			DefIfEmpty(&confHomePath, defOutputDir+"/"+chainID, confHomePath)

			m := NewBuildMetaData(nodes, confHomePath, chainID, confDirName, linkCliDir, linkdDir, nodeP2PPort,
				nodeRestAPIPort, nodeABCIPort, prometheusListenPort, tmConfig, k8STemplateFilePath, filebeatTemplateFilePath, linkDockerImageUrl)

			return buildConfForK8s(cmd, app.MakeCodec(), tmConfig, &m, minGasPrices)
		},
	}

	cmd.Flags().StringP(flagAction, "a", "build", "the action Name what you want to do")
	cmd.Flags().StringP(flagLinkDockerImageUrl, "b", defLinkDockerImageUrl, "input linkd docker image url")
	cmd.Flags().StringP(flagChainID, "c", "", "input ChainID")
	cmd.Flags().StringP(flagCliBinDirName, "d", "linkcli", "input linkcli binary home dir Name in confHomeDir")
	cmd.Flags().IntP(flagNodeABCIPort, "e", defNodeABCIPort, "input ABCI interface communication port")
	cmd.Flags().StringP(flagConfDirName, "f", defConfDirName, "input defConfDirName")
	cmd.Flags().StringSliceP(flagNodeIPs, "i", []string{"192.168.253.192", "192.168.253.193", "192.168.253.195", "192.168.224.247"}, "input node's ip list")
	cmd.Flags().StringP(flagLinkdBinDirName, "k", "linkd", "input linkd binary home dir Name in confHomeDir")
	cmd.Flags().IntP(flagPrometheusListenPort, "l", defPrometheusListenPort, "input Prometheus Listen port")
	cmd.Flags().StringP(server.FlagMinGasPrices, "m", fmt.Sprintf("%f%s", defMinGasPrices, sdk.DefaultBondDenom),
		"Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 0.01photino,0.001stake)")
	cmd.Flags().IntP(flagNodeP2PPort, "n", defNodeP2PPort, "input P2P port for Linkd")
	cmd.Flags().StringP(flagConfHomePath, "o", "", "the configuration home directory path to store configurations")
	cmd.Flags().BoolP(flagPrometheusTurnOn, "p", defPrometheusTurnOn, "turn on prometheus feature")
	cmd.Flags().IntP(flagNodeRestAPIPort, "r", defNodeRestAPIPort, "input RPC port")
	cmd.Flags().StringP(flagK8STemplateFilePath, "t", defK8STemplateFilePath, "input k8s template file path")
	cmd.Flags().StringP(flagFilebeatTemplateFilePath, "z", defFilebeatTemplateFilePath, "input filebeat template file path")
	cmd.Flags().StringP(flagDBDir, "u", defDBDir, "input Database directory path")

	return cmd
}

func buildConfForK8s(cmd *cobra.Command, cdc *codec.Codec, tmConfig *tmconfig.Config, m *BuildMetaData, minGasPrices string) error {

	serverConfig := srvconfig.DefaultConfig()
	DefIfEmpty(&serverConfig.MinGasPrices, minGasPrices, fmt.Sprintf("%f%s", defMinGasPrices, sdk.DefaultBondDenom))

	for nidx := 0; nidx < m.NumNodes; nidx++ {
		NewNode(m, nidx).process(tmConfig, serverConfig, cdc).writeK8STemplate()
		NewTemplateObject("./build/%s/deployments/filebeat-%d.yaml").MakeTemplate(m, nidx)
	}

	if err := InitGenFiles(cdc, app.ModuleBasics, m.ChainID, m.Accs, m.GenFiles, m.NumNodes); err != nil {
		return err
	}
	if err := CollectGenFiles(cdc, tmConfig, m.ChainID, m.NodeNickNames, m.ValidatorIDs, m.PubKeys,
		m.NumNodes, m.ConfHomePath, nodeDirPrefix, m.LinkdBinDirName, genaccounts.AppModuleBasic{},
	); err != nil {
		return err
	}
	cmd.Printf("Successfully initialized for [%d]nodes configuration files at %s\n", m.NumNodes, m.ConfHomePath)
	return nil
}

type Tokens struct {
	holding, staking sdk.Int
}
type Coins struct {
	node              sdk.Coin
	defBondDenomStake sdk.Coin
}
