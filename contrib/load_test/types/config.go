package types

type Config struct {
	MsgsPerTxPrepare  int
	MsgsPerTxLoadTest int
	TPS               int
	Duration          int
	MaxWorkers        int
	PacerType         string
	TargetURL         string
	ChainID           string
	CoinName          string
	Mnemonic          string
}
