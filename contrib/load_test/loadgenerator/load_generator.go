package loadgenerator

import (
	"log"
	"runtime/debug"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/types"
	"github.com/line/link/contrib/load_test/wallet"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"golang.org/x/sync/errgroup"
)

const (
	maxTargetsPerUser = 100
)

type LoadGenerator struct {
	hdWallet          *wallet.HDWallet
	targets           []vegeta.Target
	config            types.Config
	numUsers          int
	numTargetsPerUser int
}

func NewLoadGenerator() *LoadGenerator {
	return &LoadGenerator{}
}

func (lg *LoadGenerator) ApplyConfig(config types.Config) (err error) {
	lg.config = config

	types.SetBech32Prefix(sdk.GetConfig(), config.Testnet)
	lg.numUsers = config.TPS * config.Duration
	lg.targets = make([]vegeta.Target, 0, lg.numUsers*maxTargetsPerUser)
	lg.hdWallet, err = wallet.NewHDWallet(config.Mnemonic)
	return
}

func (lg *LoadGenerator) RunWithGoroutines(generateTargetFunc func(*wallet.KeyWallet, int) (*[]*vegeta.Target,
	int, error)) (err error) {
	log.Println("Start to generate target")
	tmpTargets := make([]*vegeta.Target, lg.numUsers*maxTargetsPerUser)
	// Semaphore is used to limit the maximum number of executable goroutines.
	sem := make(chan int, lg.config.MaxWorkers)
	var eg errgroup.Group
	numTargetsChan := make(chan int, lg.numUsers)
	for i := 0; i < lg.numUsers; i++ {
		i := i
		sem <- 1
		eg.Go(func() error {
			defer CompleteGoroutine(&sem)
			keyWallet, err := lg.hdWallet.GetKeyWallet(0, uint32(i))
			if err != nil {
				return err
			}
			targets, numTargets, err := generateTargetFunc(keyWallet, i)
			if err != nil {
				return err
			}
			if numTargets > maxTargetsPerUser {
				return types.ExceedMaxNumTargetsError{NumTargets: numTargets, MaxTargetsPerUser: maxTargetsPerUser}
			}
			for j := 0; j < numTargets; j++ {
				tmpTargets[numTargets*i+j] = (*targets)[j]
			}
			numTargetsChan <- numTargets
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	close(numTargetsChan)
	lg.targets = removeEmpty(tmpTargets)
	lg.numTargetsPerUser = getNumTargets(numTargetsChan)
	return nil
}

func (lg *LoadGenerator) Fire(lgURL string) <-chan *vegeta.Result {
	log.Println("Start to fire")
	duration := time.Duration(lg.config.Duration) * time.Second
	pacer := RampUpPacer{
		Constant:   vegeta.Rate{Freq: lg.config.TPS * lg.numTargetsPerUser, Per: time.Second},
		RampUpTime: time.Duration(lg.config.RampUpTime) * time.Second,
	}
	targeter := vegeta.NewStaticTargeter(lg.targets...)
	attacker := vegeta.NewAttacker()

	return attacker.Attack(targeter, pacer, duration, "LINK v2 load test: "+lgURL)
}

func CompleteGoroutine(sem *chan int) {
	<-*sem
	if err := recover(); err != nil {
		log.Println("Failed to generate target:", err)
		log.Println(string(debug.Stack()))
	}
}

func removeEmpty(targets []*vegeta.Target) []vegeta.Target {
	var r []vegeta.Target
	for _, target := range targets {
		if target != nil {
			r = append(r, *target)
		}
	}
	return r
}

func getNumTargets(numTargetsChan chan int) int {
	var numTargets int
	for nt := range numTargetsChan {
		switch {
		case numTargets == 0:
			numTargets = nt
		case numTargets != nt:
			panic("The numTargets in the same scenario should be the same.")
		}
	}
	return numTargets
}
