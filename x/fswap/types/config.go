package types

type Config struct {
	MaxSwaps      int
	UpdateAllowed bool
}

func DefaultConfig() Config {
	return Config{
		MaxSwaps:      1,
		UpdateAllowed: false,
	}
}
