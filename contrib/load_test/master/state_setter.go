package master

import (
	"log"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/app"
	"github.com/line/link/contrib/load_test/service"
	"github.com/line/link/contrib/load_test/transaction"
	"github.com/line/link/contrib/load_test/types"
	"github.com/line/link/contrib/load_test/wallet"
	"github.com/line/link/x/coin"
)

const (
	BroadcastMode = "block"
	GAS           = 10000000000000000
)

type StateSetter struct {
	masterKeyWallet *wallet.KeyWallet
	txBuilder       transaction.TxBuilderWithoutKeybase
	linkService     *service.LinkService

	config types.Config
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
		masterKeyWallet: masterWallet,
		txBuilder:       transaction.NewTxBuilder().WithGas(GAS).WithChainID(config.ChainID),
		linkService:     service.NewLinkService(&http.Client{}, app.MakeCodec(), config.TargetURL),
		config:          config,
	}, nil
}

func (ss *StateSetter) PrepareTestState(slaves []types.Slave) error {
	for _, slave := range slaves {
		if err := ss.RegisterAccounts(slave.Mnemonic); err != nil {
			return err
		}
	}
	return nil
}

func (ss *StateSetter) RegisterAccounts(mnemonic string) error {
	hdWallet, err := wallet.NewHDWallet(mnemonic)
	if err != nil {
		return types.InvalidMnemonic{Mnemonic: mnemonic}
	}

	numTargets := ss.config.TPS * ss.config.Duration
	numMsgsSent := 0
	for i := 0; numMsgsSent < numTargets; i++ {
		numMsgs := min(ss.config.MsgsPerTxPrepare, numTargets-numMsgsSent)
		if err := ss.broadcastRegistrationTx(hdWallet, numMsgsSent, numMsgs); err != nil {
			return err
		}
		numMsgsSent += numMsgs
	}
	return nil
}

func (ss *StateSetter) broadcastRegistrationTx(hdWallet *wallet.HDWallet, startAccountIndex, numMsgs int) error {
	account, err := ss.linkService.GetAccount(ss.masterKeyWallet.Address().String())
	if err != nil {
		return err
	}

	coins := sdk.NewCoins(sdk.NewCoin(ss.config.CoinName, sdk.NewInt(1)))
	msgs := make([]sdk.Msg, numMsgs)
	for i := 0; i < numMsgs; i++ {
		keyWallet, err := hdWallet.GetKeyWallet(0, uint32(startAccountIndex+i))
		if err != nil {
			return err
		}
		to := keyWallet.Address()
		msgs[i] = coin.NewMsgSend(ss.masterKeyWallet.Address(), to, coins)
	}

	stdTx, err := ss.txBuilder.WithAccountNumber(account.AccountNumber).WithSequence(account.Sequence).
		BuildAndSign(ss.masterKeyWallet.PrivateKey(), msgs)
	if err != nil {
		return err
	}

	res, err := ss.linkService.BroadcastTx(stdTx, BroadcastMode)
	if err != nil {
		return err
	}
	log.Printf("Send tx to register accounts : %s\n", res.TxHash)
	return nil
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
