package loadgenerator

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/app"
	"github.com/line/link/contrib/load_test/service"
	"github.com/line/link/contrib/load_test/transaction"
	"github.com/line/link/contrib/load_test/types"
	"github.com/line/link/contrib/load_test/wallet"
	"github.com/line/link/x/bank"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"golang.org/x/sync/errgroup"
)

const (
	BroadcastMode = "sync"
	GAS           = 10000000000000000
)

type LoadGenerator struct {
	hdWallet      *wallet.HDWallet
	linkService   *service.LinkService
	targetBuilder *TargetBuilder
	txBuilder     transaction.TxBuilderWithoutKeybase
	targets       []vegeta.Target
	config        types.Config
	numTargets    int
}

func NewLoadGenerator() *LoadGenerator {
	return &LoadGenerator{
		txBuilder: transaction.NewTxBuilder().WithGas(GAS),
	}
}

func (lg *LoadGenerator) ApplyConfig(config types.Config) (err error) {
	lg.config = config

	lg.linkService = service.NewLinkService(&http.Client{}, app.MakeCodec(), config.TargetURL)
	lg.targetBuilder = NewTargetBuilder(app.MakeCodec(), config.TargetURL)
	lg.txBuilder = lg.txBuilder.WithChainID(config.ChainID)
	lg.numTargets = config.TPS * config.Duration
	lg.targets = make([]vegeta.Target, lg.numTargets)
	lg.hdWallet, err = wallet.NewHDWallet(config.Mnemonic)
	return
}

func (lg *LoadGenerator) RunWithGoroutines(generateTargetFunc func(*chan int, int) error) error {
	// Semaphore is used to limit the maximum number of executable goroutines.
	sem := make(chan int, lg.config.MaxWorkers)
	var eg errgroup.Group
	for i := 0; i < lg.numTargets; i++ {
		i := i
		sem <- 1
		eg.Go(func() error {
			return generateTargetFunc(&sem, i)
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (lg *LoadGenerator) Fire(lgURL string) <-chan *vegeta.Result {
	duration := time.Duration(lg.config.Duration) * time.Second

	var pacer vegeta.Pacer
	switch lg.config.PacerType {
	case types.ConstantPacer:
		pacer = vegeta.Rate{Freq: lg.config.TPS, Per: time.Second}
	case types.LinearPacer:
		slope := float64(lg.config.TPS) / float64(lg.config.Duration)
		pacer = vegeta.LinearPacer{
			StartAt: vegeta.Rate{Freq: 1, Per: time.Second},
			Slope:   slope,
		}
	}
	targeter := vegeta.NewStaticTargeter(lg.targets...)
	attacker := vegeta.NewAttacker()

	return attacker.Attack(targeter, pacer, duration, "LINK v2 load test: "+lgURL)
}

func (lg *LoadGenerator) GenerateAccountQueryTarget(sem *chan int, i int) error {
	defer CompleteGoroutine(sem)

	keyWallet, err := lg.hdWallet.GetKeyWallet(0, uint32(i))
	if err != nil {
		return err
	}

	lg.targets[i] = lg.targetBuilder.MakeQueryTarget("/auth/accounts/" + keyWallet.Address().String())
	return nil
}

func (lg *LoadGenerator) GenerateMsgSendTxTarget(sem *chan int, i int) error {
	defer CompleteGoroutine(sem)

	keyWallet, err := lg.hdWallet.GetKeyWallet(0, uint32(i))
	if err != nil {
		return err
	}

	account, err := lg.linkService.GetAccount(keyWallet.Address().String())
	if err != nil {
		return err
	}

	from := account.Address
	to := secp256k1.GenPrivKey().PubKey().Address().Bytes()
	coins := sdk.NewCoins(sdk.NewCoin(lg.config.CoinName, sdk.NewInt(1)))
	msgs := make([]sdk.Msg, lg.config.MsgsPerTxLoadTest)
	for i := 0; i < lg.config.MsgsPerTxLoadTest; i++ {
		msgs[i] = bank.NewMsgSend(from, to, coins)
	}

	stdTx, err := lg.txBuilder.WithAccountNumber(account.AccountNumber).WithSequence(account.Sequence).
		BuildAndSign(keyWallet.PrivateKey(), msgs)
	if err != nil {
		return err
	}

	lg.targets[i], err = lg.targetBuilder.MakeTxTarget(stdTx, BroadcastMode)
	if err != nil {
		return err
	}
	return nil
}

func CompleteGoroutine(sem *chan int) {
	<-*sem
	if err := recover(); err != nil {
		log.Println("Failed to generate target:", err)
		log.Println(string(debug.Stack()))
	}
}
