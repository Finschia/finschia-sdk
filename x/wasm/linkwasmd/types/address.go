package types

const (
	Bech32MainPrefix    = "link"
	Bech32TestnetPrefix = "tlink"

	// LINK in [SLIP-044](https://github.com/satoshilabs/slips/blob/master/slip-0044.md)
	CoinType           = 438
	FullFundraiserPath = "44'/438'/0'/0/0"
)

func Bech32PrefixAcc(testnet bool) string {
	if testnet {
		return Bech32TestnetPrefix
	}
	return Bech32MainPrefix
}

func Bech32PrefixAccPub(testnet bool) string {
	return Bech32PrefixAcc(testnet) + PrefixPublic
}

func Bech32PrefixValAddr(testnet bool) string {
	return Bech32PrefixAcc(testnet) + PrefixValidator + PrefixOperator
}

func Bech32PrefixValPub(testnet bool) string {
	return Bech32PrefixAcc(testnet) + PrefixValidator + PrefixOperator + PrefixPublic
}

func Bech32PrefixConsAddr(testnet bool) string {
	return Bech32PrefixAcc(testnet) + PrefixValidator + PrefixConsensus
}

func Bech32PrefixConsPub(testnet bool) string {
	return Bech32PrefixAcc(testnet) + PrefixValidator + PrefixConsensus + PrefixPublic
}
