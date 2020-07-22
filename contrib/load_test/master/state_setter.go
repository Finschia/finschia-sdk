package master

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/app"
	"github.com/line/link/contrib/load_test/scenario"
	"github.com/line/link/contrib/load_test/service"
	"github.com/line/link/contrib/load_test/transaction"
	"github.com/line/link/contrib/load_test/types"
	"github.com/line/link/contrib/load_test/wallet"
)

const (
	ParamsFileName = "test_params.json"
)

type StateSetter struct {
	linkService     *service.LinkService
	txBuilder       transaction.TxBuilderWithoutKeybase
	masterKeyWallet *wallet.KeyWallet
	scenarios       scenario.Scenarios
	config          types.Config
}

func NewStateSetter(masterMnemonic string, config types.Config) (*StateSetter, error) {
	masterHDWallet, err := wallet.NewHDWallet(masterMnemonic)
	if err != nil {
		return nil, types.InvalidMasterMnemonic{Mnemonic: masterMnemonic}
	}
	masterWallet, err := masterHDWallet.GetKeyWallet(0, 0)
	if err != nil {
		return nil, err
	}
	return &StateSetter{
		linkService:     service.NewLinkService(&http.Client{}, app.MakeCodec(), config.TargetURL),
		txBuilder:       transaction.NewTxBuilder(uint64(config.MaxGasPrepare)).WithChainID(config.ChainID),
		masterKeyWallet: masterWallet,
		scenarios:       scenario.NewScenarios(config, nil, nil),
		config:          config,
	}, nil
}

func (ss *StateSetter) PrepareTestState(slaves []types.Slave, outputDir string) error {
	allParams := make(map[string]map[string]string)
	shouldSave := false
	for _, slave := range slaves {
		hdWallet, err := wallet.NewHDWallet(slave.Mnemonic)
		if err != nil {
			return types.InvalidMnemonic{Mnemonic: slave.Mnemonic}
		}
		log.Println("Start to generate state setting msgs")
		msgs, params, err := ss.scenarios[slave.Scenario].GenerateStateSettingMsgs(ss.masterKeyWallet, hdWallet, slave.Params)
		if err != nil {
			return err
		}
		if msgs == nil {
			return nil
		}
		if err := ss.broadcastMsgs(msgs, scenario.GetPrepareBroadcastMode(slave.Scenario)); err != nil {
			return err
		}
		if params != nil {
			allParams[slave.URL] = params
			shouldSave = true
		}
	}
	if shouldSave {
		if err := ss.saveParams(allParams, outputDir); err != nil {
			return err
		}
	}
	return nil
}

func (ss *StateSetter) broadcastMsgs(msgs []sdk.Msg, mode string) error {
	numPrepareTxs := (len(msgs) + ss.config.MsgsPerTxPrepare - 1) / ss.config.MsgsPerTxPrepare // round up
	account, err := ss.linkService.GetAccount(ss.masterKeyWallet.Address().String())
	if err != nil {
		return err
	}

	for i := 0; i < numPrepareTxs; i++ {
		start := ss.config.MsgsPerTxPrepare * i
		end := min(ss.config.MsgsPerTxPrepare*(i+1), len(msgs))
		res, err := ss.linkService.BroadcastMsgs(ss.txBuilder, account, ss.masterKeyWallet, msgs[start:end], mode)
		account.Sequence++
		if err != nil {
			return err
		}
		if res.Code != 0 {
			return types.FailedTxError{Tx: res}
		}
		log.Printf("Send state setting tx successfully (%d/%d) : %s\n", i+1, numPrepareTxs, res.TxHash)
	}
	log.Println("Preparation is completed")
	return nil
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (ss *StateSetter) saveParams(params map[string]map[string]string, outputDir string) error {
	log.Println("Save file of test parameters")
	file, err := json.Marshal(params)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(fmt.Sprintf("%s/%s", outputDir, ParamsFileName), file, 0600); err != nil {
		return err
	}
	return nil
}
