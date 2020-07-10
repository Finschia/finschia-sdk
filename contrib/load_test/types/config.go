package types

import (
	"log"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	linktypes "github.com/line/link/types"
)

type Config struct {
	MsgsPerTxPrepare  int
	MaxGasPrepare     int
	MsgsPerTxLoadTest int
	MaxGasLoadTest    int
	TPS               int
	Duration          int
	RampUpTime        int
	MaxWorkers        int
	TargetURL         string
	ChainID           string
	CoinName          string
	Testnet           bool
	Mnemonic          string
}

var mutex = &sync.Mutex{}

func SetBech32Prefix(config *sdk.Config, testnet bool) {
	defer func() {
		mutex.Unlock()
		v := recover()
		if v != nil {
			log.Println(v, "The current runtime should be an integration test")
		}
	}()

	mutex.Lock()
	config.SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(testnet), linktypes.Bech32PrefixAccPub(testnet))
	config.SetBech32PrefixForValidator(linktypes.Bech32PrefixValAddr(testnet), linktypes.Bech32PrefixValPub(testnet))
	config.SetBech32PrefixForConsensusNode(linktypes.Bech32PrefixConsAddr(testnet), linktypes.Bech32PrefixConsPub(testnet))
}
