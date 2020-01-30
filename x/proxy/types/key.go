package types

const (
	ModuleName = "proxy"

	StoreKey  = ModuleName
	RouterKey = ModuleName
)

var (
	ProxyKeyPrefix = []byte{0x00}
)

func ProxyKey(proxyAddress string) []byte {
	return append(ProxyKeyPrefix, []byte(proxyAddress)...)
}
