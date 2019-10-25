package bank

import (
	cbank "github.com/cosmos/cosmos-sdk/x/bank"
	cbankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	"github.com/link-chain/link/x/bank/internal/types"
)

const (
	DefaultCodespace  = cbank.DefaultCodespace
	ModuleName        = cbank.ModuleName
	DefaultParamspace = cbank.DefaultParamspace
)

var (
	// functions aliases
	RegisterCodec                   = types.RegisterCodec
	NewBaseKeeper                   = cbank.NewBaseKeeper
	NewCosmosAppModule              = cbank.NewAppModule
	SimulateMsgSend                 = cbank.SimulateMsgSend
	SimulateSingleInputMsgMultiSend = cbank.SimulateSingleInputMsgMultiSend

	// functions of client
	SendTxCmd = cbankcli.SendTxCmd
	// variable aliases
	ModuleCdc = types.ModuleCdc
	NewInput  = cbank.NewInput
	NewOutput = cbank.NewOutput
)

type (
	Keeper               = cbank.Keeper
	AccountKeeper        = types.AccountKeeper
	BaseKeeper           = cbank.BaseKeeper
	CosmosAppModuleBasic = cbank.AppModuleBasic
	CosmosAppModule      = cbank.AppModule
	MsgMultiSend         = cbank.MsgMultiSend
	Input                = cbank.Input
	Output               = cbank.Output
)
