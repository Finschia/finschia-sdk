package k8s

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	tmconfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
)

type Node struct {
	Name     string
	Idx      int
	MetaData *BuildMetaData
}

func NewNode(m *BuildMetaData, nidx int) *Node {
	name := fmt.Sprintf("%s%d", nodeDirPrefix, nidx)
	m.NodeNickNames[nidx] = name
	return &Node{name, nidx, m}
}
func (n *Node) process(tmConfig *tmconfig.Config, cosConfig *srvconfig.Config, cdc *codec.Codec) *Node {
	tmConfig.SetRoot(n.linkdDirFullPath())
	tmConfig.Moniker = fmt.Sprintf("%s%s", n.MetaData.InputNodeIp[n.Idx], n.Name)

	n.prepareOutputDir(n.MetaData.ConfHomePath, n.linkdDirFullPath(), n.cliBinDirNameFullPath())

	var err error
	n.MetaData.ValidatorIDs[n.Idx], n.MetaData.PubKeys[n.Idx], err = genutil.InitializeNodeValidatorFiles(tmConfig)
	if err != nil {
		_ = os.RemoveAll(n.MetaData.ConfHomePath)
		panic(err)
	}
	n.MetaData.GenFiles[n.Idx] = tmConfig.GenesisFile()
	addr, secret, err := server.GenerateSaveCoinKey(n.cliBinDirNameFullPath(), n.Name, client.DefaultKeyPass,
		true)
	if err != nil {
		_ = os.RemoveAll(n.MetaData.ConfHomePath)
		panic(err)
	}
	cliPrint, err := json.Marshal(map[string]string{"secret": secret})
	if err != nil {
		panic(err)
	}
	if err := WriteFile(n.cliBinDirNameFullPath(), fmt.Sprintf("%v%s",
		defPrivateKeySeedFileName, defConfigurationFileExt), cliPrint); err != nil {
		panic(err)
	}
	n.MetaData.Accs[n.Idx] = buildGenesisAcc(n.Name, addr)
	if err := writeGenTx(n, addr, cdc); err != nil {
		panic(err)
	}
	srvconfig.WriteConfigFile(filepath.Join(n.linkdDirFullPath(), n.MetaData.ConfDirName)+"/app.toml", cosConfig)
	return n
}

func (n *Node) PubKey() crypto.PubKey {
	return n.MetaData.PubKeys[n.Idx]
}
func (n *Node) InputNodeIp() string {
	return n.MetaData.InputNodeIp[n.Idx]
}

func (n *Node) cliBinDirNameFullPath() string {
	return filepath.Join(n.MetaData.ConfHomePath, n.Name, n.MetaData.CliBinDirName)
}

func (n *Node) linkdDirFullPath() string {
	return filepath.Join(n.MetaData.ConfHomePath, n.Name, n.MetaData.LinkdBinDirName)
}
func (n *Node) gentxsDirFullPath() string {
	return filepath.Join(n.MetaData.ConfHomePath, genTxsDefaultDir)
}

func (n *Node) prepareOutputDir(confHomePath string, nodeDirFullPath string, nodeCLIDirFullPath string) {
	if err := os.MkdirAll(filepath.Join(nodeDirFullPath, defConfDirName), nodeDirPerm); err != nil {
		_ = os.RemoveAll(confHomePath)
		panic(err)
	}
	if err := os.MkdirAll(nodeCLIDirFullPath, nodeDirPerm); err != nil {
		_ = os.RemoveAll(confHomePath)
		panic(err)
	}
}

func (n *Node) writeK8STemplate() *Deployment {
	deploymentTemplate := DeploymentTemplate{Node: n}
	deployment, err := deploymentTemplate.Write()
	if err != nil {
		panic(err)
	}
	return deployment
}

type BuildMetaData struct {
	InputNodeIp, ValidatorIDs, GenFiles, NodeNickNames                 []string
	PubKeys                                                            []crypto.PubKey
	Accs                                                               []genaccounts.GenesisAccount
	NumNodes                                                           int
	ConfHomePath, ChainID, ConfDirName, CliBinDirName, LinkdBinDirName string
	k8STemplateFilePath, filebeatTemplateFilePath, linkDockerImageUrl  string
	NodeP2PPort, NodeRestAPIPort, NodeABCIPort, PrometheusListenPort   int
	TmConfig                                                           *tmconfig.Config
}

func NewBuildMetaData(inputNodes []string, confHomePath, chainID, confDirName, cliBinDirName, linkdBinDirName string,
	nodeP2PPort, nodeRestAPIPort, nodeABCIPort, prometheusListenPort int, tmConfig *tmconfig.Config,
	k8STemplateFilePath string, filebeatTemplateFilePath string, linkDockerImageUrl string) BuildMetaData {
	numNodes := len(inputNodes)
	return BuildMetaData{inputNodes, make([]string, numNodes), make([]string, numNodes),
		make([]string, numNodes), make([]crypto.PubKey, numNodes),
		make([]genaccounts.GenesisAccount, numNodes), numNodes, confHomePath,
		chainID, confDirName, cliBinDirName, linkdBinDirName,
		k8STemplateFilePath, filebeatTemplateFilePath,
		linkDockerImageUrl, nodeP2PPort, nodeRestAPIPort,
		nodeABCIPort, prometheusListenPort, tmConfig}
}
