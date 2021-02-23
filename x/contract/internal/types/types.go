package types

import (
	sdk "github.com/line/lbm-sdk/types"
)

const (
	ModuleName = "contract"

	StoreKey = ModuleName
)

var (
	LastContractCountStoreKeyPrefix = []byte{0x01}
	ContractIDStoreKeyPrefix        = []byte{0x02}
)

func LastContractCountStoreKey() []byte {
	return LastContractCountStoreKeyPrefix
}

func ContractIDStoreKey(contractID string) []byte {
	return append(ContractIDStoreKeyPrefix, []byte(contractID)...)
}

type ContractMsg interface {
	sdk.Msg
	GetContractID() string
}

type CtxKey struct{}
