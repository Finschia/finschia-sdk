package types

const (
	Bech32MainPrefix    = "link"
	Bech32TestnetPrefix = "tlink"

	// LINK in [SLIP-044](https://github.com/satoshilabs/slips/blob/master/slip-0044.md)
	CoinType           = 438
	FullFundraiserPath = "44'/438'/0'/0/0"
)

var (
	testnetMode = false
)

func SetTestnetMode() {
	testnetMode = true
}

func IsTestnetMode() bool {
	return testnetMode
}

func prefix() string {
	if testnetMode {
		return Bech32TestnetPrefix
	}
	return Bech32MainPrefix
}

func Bech32PrefixAccAddr() string {
	return prefix()
}

func Bech32PrefixAccPub() string {
	return prefix() + PrefixPublic
}

func Bech32PrefixValAddr() string {
	return prefix() + PrefixValidator + PrefixOperator
}

func Bech32PrefixValPub() string {
	return prefix() + PrefixValidator + PrefixOperator + PrefixPublic
}

func Bech32PrefixConsAddr() string {
	return prefix() + PrefixValidator + PrefixConsensus
}

func Bech32PrefixConsPub() string {
	return prefix() + PrefixValidator + PrefixConsensus + PrefixPublic
}
