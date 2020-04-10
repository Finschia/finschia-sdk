package contract

import (
	"github.com/line/link/x/contract/internal/keeper"
	"github.com/line/link/x/contract/internal/types"
)

const (
	ModuleName       = types.ModuleName
	StoreKey         = types.StoreKey
	SampleContractID = "abcde012"
)

type (
	Msg    = types.ContractMsg
	Keeper = keeper.ContractKeeper
)

var (
	ErrInvalidContractID = types.ErrInvalidContractID
	ErrContractNotExist  = types.ErrContractNotExist
	NewContractKeeper    = keeper.NewContractKeeper
)
