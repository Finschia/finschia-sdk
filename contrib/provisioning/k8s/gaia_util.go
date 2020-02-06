package k8s

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	types2 "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/staking"
	types3 "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmconfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
)

func InitGenFiles(cdc *codec.Codec, mbm module.BasicManager, chainID string,
	accs []genaccounts.GenesisAccount, genFiles []string, numValidators int) error {

	appGenState := mbm.DefaultGenesis()

	// set the accounts in the genesis state
	appGenState = genaccounts.SetGenesisStateInAppState(cdc, appGenState, accs)

	appGenStateJSON, err := codec.MarshalJSONIndent(cdc, appGenState)
	if err != nil {
		return err
	}

	genDoc := types.GenesisDoc{
		ChainID:    chainID,
		AppState:   appGenStateJSON,
		Validators: nil,
	}

	// generate empty genesis files for each validator and save
	for i := 0; i < numValidators; i++ {
		if err := genDoc.SaveAs(genFiles[i]); err != nil {
			return err
		}
	}
	return nil
}

func CollectGenFiles(
	cdc *codec.Codec, config *tmconfig.Config, chainID string,
	monikers, nodeIDs []string, valPubKeys []crypto.PubKey,
	numValidators int, outputDir, nodeDirPrefix, nodeDaemonHome string,
	genAccIterator genutiltypes.GenesisAccountsIterator) error {

	var appState json.RawMessage
	genTime := tmtime.Now()

	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
		nodeDir := filepath.Join(outputDir, nodeDirName, nodeDaemonHome)
		gentxsDir := filepath.Join(outputDir, "gentxs")
		moniker := monikers[i]
		config.Moniker = nodeDirName

		config.SetRoot(nodeDir)

		nodeID, valPubKey := nodeIDs[i], valPubKeys[i]
		initCfg := genutil.NewInitConfig(chainID, gentxsDir, moniker, nodeID, valPubKey)

		genDoc, err := types.GenesisDocFromFile(config.GenesisFile())
		if err != nil {
			return err
		}

		nodeAppState, err := genutil.GenAppStateFromConfig(cdc, config, initCfg, *genDoc, genAccIterator)
		if err != nil {
			return err
		}

		if appState == nil {
			// set the canonical application state (they should not differ)
			appState = nodeAppState
		}

		genFile := config.GenesisFile()

		// overwrite each validator's genesis file to have a canonical genesis time
		if err := genutil.ExportGenesisFileWithTime(genFile, chainID, nil, appState, genTime); err != nil {
			return err
		}
	}

	return nil
}

func WriteFile(dir string, name string, contents []byte) error {
	writePath := filepath.Join(dir)
	file := filepath.Join(writePath, name)

	err := cmn.EnsureDir(writePath, 0700)
	if err != nil {
		return err
	}

	err = cmn.WriteFile(file, contents, 0600)
	if err != nil {
		return err
	}

	return nil
}

func writeGenTx(n *Node, addr types2.AccAddress, cdc *codec.Codec) error {
	kb, err := keys.NewKeyBaseFromDir(n.cliBinDirNameFullPath())
	if err != nil {
		panic(err)
	}
	tx := auth.NewStdTx([]types2.Msg{types3.MsgCreateValidator{
		Description:       staking.NewDescription(n.Name, "", "", ""),
		Commission:        staking.NewCommissionRates(types2.ZeroDec(), types2.ZeroDec(), types2.ZeroDec()),
		MinSelfDelegation: types2.OneInt(),
		DelegatorAddress:  addr,
		ValidatorAddress:  types2.ValAddress(addr),
		PubKey:            n.PubKey(),
		Value:             types2.NewCoin(types2.DefaultBondDenom, types2.TokensFromConsensusPower(100)),
	}}, auth.StdFee{}, []auth.StdSignature{}, fmt.Sprintf("%s@%s:%d", n.MetaData.ValidatorIDs[n.Idx],
		n.InputNodeIP(), n.MetaData.NodeP2PPort))

	signedTx, err := auth.NewTxBuilderFromCLI().WithChainID(n.MetaData.ChainID).WithMemo(tx.Memo).WithKeybase(kb).
		SignStdTx(n.Name, client.DefaultKeyPass, tx, false)
	if err != nil {
		_ = os.RemoveAll(n.MetaData.ConfHomePath)
		panic(err)
	}
	txBytes, err := cdc.MarshalJSON(signedTx)
	if err != nil {
		_ = os.RemoveAll(n.MetaData.ConfHomePath)
		panic(err)
	}
	if err := WriteFile(n.gentxsDirFullPath(), fmt.Sprintf("%v%s", n.Name,
		defConfigurationFileExt), txBytes); err != nil {
		_ = os.RemoveAll(n.MetaData.ConfHomePath)
		panic(err)
	}
	return nil
}

func buildGenesisAcc(nodeDirName string, addr types2.AccAddress) genaccounts.GenesisAccount {
	accTokens := Tokens{types2.TokensFromConsensusPower(1000),
		types2.TokensFromConsensusPower(500)}
	coins := Coins{types2.NewCoin(fmt.Sprintf("%stoken", nodeDirName), accTokens.holding),
		types2.NewCoin(types2.DefaultBondDenom, accTokens.staking)}
	return genaccounts.GenesisAccount{
		Address: addr,
		Coins: types2.Coins{
			coins.node,
			coins.defBondDenomStake,
		},
	}
}
