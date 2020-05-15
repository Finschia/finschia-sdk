package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"

	linklcd "github.com/line/link/client/lcd"
	linkrpc "github.com/line/link/client/rpc"
	"github.com/line/link/client/rpc/link/block"
	"github.com/line/link/client/rpc/link/genesis"
	"github.com/line/link/client/rpc/link/mempool"
)

const (
	DefaultGasAdjustment   = flags.DefaultGasAdjustment
	DefaultGasLimit        = flags.DefaultGasLimit
	GasFlagAuto            = flags.GasFlagAuto
	BroadcastBlock         = flags.BroadcastBlock
	BroadcastSync          = flags.BroadcastSync
	BroadcastAsync         = flags.BroadcastAsync
	FlagHome               = flags.FlagHome
	FlagUseLedger          = flags.FlagUseLedger
	FlagChainID            = flags.FlagChainID
	FlagNode               = flags.FlagNode
	FlagHeight             = flags.FlagHeight
	FlagGasAdjustment      = flags.FlagGasAdjustment
	FlagTrustNode          = flags.FlagTrustNode
	FlagFrom               = flags.FlagFrom
	FlagName               = flags.FlagName
	FlagAccountNumber      = flags.FlagAccountNumber
	FlagSequence           = flags.FlagSequence
	FlagMemo               = flags.FlagMemo
	FlagFees               = flags.FlagFees
	FlagGasPrices          = flags.FlagGasPrices
	FlagBroadcastMode      = flags.FlagBroadcastMode
	FlagDryRun             = flags.FlagDryRun
	FlagGenerateOnly       = flags.FlagGenerateOnly
	FlagIndentResponse     = flags.FlagIndentResponse
	FlagListenAddr         = flags.FlagListenAddr
	FlagMaxOpenConnections = flags.FlagMaxOpenConnections
	FlagRPCReadTimeout     = flags.FlagRPCReadTimeout
	FlagRPCWriteTimeout    = flags.FlagRPCWriteTimeout
	FlagOutputDocument     = flags.FlagOutputDocument
	FlagSkipConfirmation   = flags.FlagSkipConfirmation
	DefaultKeyPass         = keys.DefaultKeyPass
	FlagAddress            = keys.FlagAddress
	FlagPublicKey          = keys.FlagPublicKey
	FlagBechPrefix         = keys.FlagBechPrefix
	FlagDevice             = keys.FlagDevice
	OutputFormatText       = keys.OutputFormatText
	OutputFormatJSON       = keys.OutputFormatJSON
	MinPassLength          = input.MinPassLength
)

var (
	NewCLIContextWithFrom              = context.NewCLIContextWithFrom
	NewCLIContextWithInputAndFrom      = context.NewCLIContextWithInputAndFrom
	NewCLIContext                      = context.NewCLIContext
	GetFromFields                      = context.GetFromFields
	ErrInvalidAccount                  = context.ErrInvalidAccount
	ErrVerifyCommit                    = context.ErrVerifyCommit
	GetCommands                        = flags.GetCommands
	PostCommands                       = flags.PostCommands
	RegisterRestServerFlags            = flags.RegisterRestServerFlags
	ParseGas                           = flags.ParseGas
	NewCompletionCmd                   = flags.NewCompletionCmd
	MarshalJSON                        = keys.MarshalJSON
	UnmarshalJSON                      = keys.UnmarshalJSON
	Commands                           = keys.Commands
	NewAddNewKey                       = keys.NewAddNewKey
	NewRecoverKey                      = keys.NewRecoverKey
	NewUpdateKeyReq                    = keys.NewUpdateKeyReq
	NewDeleteKeyReq                    = keys.NewDeleteKeyReq
	NewKeyBaseFromDir                  = keys.NewKeyBaseFromDir
	NewInMemoryKeyBase                 = keys.NewInMemoryKeyBase
	NewRestServer                      = lcd.NewRestServer
	ServeCommand                       = linklcd.ServeCommand
	BlockCommand                       = rpc.BlockCommand
	BlockWithResultCommand             = block.WithTxResultCommand
	QueryGenesisAccountCmd             = genesis.QueryGenesisAccountCmd
	QueryGenesisTxCmd                  = genesis.QueryGenesisTxCmd
	BlockRequestHandlerFn              = rpc.BlockRequestHandlerFn
	LatestBlockRequestHandlerFn        = rpc.LatestBlockRequestHandlerFn
	MempoolCmd                         = mempool.Command
	RegisterRPCRoutes                  = linkrpc.RegisterRPCRoutes
	StatusCommand                      = rpc.StatusCommand
	NodeInfoRequestHandlerFn           = rpc.NodeInfoRequestHandlerFn
	NodeSyncingRequestHandlerFn        = rpc.NodeSyncingRequestHandlerFn
	ValidatorCommand                   = rpc.ValidatorCommand
	GetValidators                      = rpc.GetValidators
	ValidatorSetRequestHandlerFn       = rpc.ValidatorSetRequestHandlerFn
	LatestValidatorSetRequestHandlerFn = rpc.LatestValidatorSetRequestHandlerFn
	GetPassword                        = input.GetPassword
	GetCheckPassword                   = input.GetCheckPassword
	GetConfirmation                    = input.GetConfirmation
	GetString                          = input.GetString
	PrintPrefixed                      = input.PrintPrefixed

	LineBreak  = flags.LineBreak
	GasFlagVar = flags.GasFlagVar

	ConfigCmd   = client.ConfigCmd
	ValidateCmd = client.ValidateCmd
)

type (
	CLIContext             = context.CLIContext
	GasSetting             = flags.GasSetting
	AddNewKey              = keys.AddNewKey
	RecoverKey             = keys.RecoverKey
	UpdateKeyReq           = keys.UpdateKeyReq
	DeleteKeyReq           = keys.DeleteKeyReq
	RestServer             = lcd.RestServer
	ValidatorOutput        = rpc.ValidatorOutput
	ResultValidatorsOutput = rpc.ResultValidatorsOutput
)
