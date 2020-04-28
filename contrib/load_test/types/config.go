package types

type Config struct {
	MsgsPerTxPrepare  int
	MsgsPerTxLoadTest int
	TPS               int
	Duration          int
	RampUpTime        int
	MaxWorkers        int
	TargetURL         string
	ChainID           string
	CoinName          string
	Mnemonic          string
}
