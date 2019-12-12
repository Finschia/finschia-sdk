package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/link-chain/link/client/lcd"
	"github.com/link-chain/link/client/rpc"
)

const (
	DefaultGasAdjustment   = client.DefaultGasAdjustment
	DefaultGasLimit        = client.DefaultGasLimit
	GasFlagAuto            = client.GasFlagAuto
	BroadcastBlock         = client.BroadcastBlock
	BroadcastSync          = client.BroadcastSync
	BroadcastAsync         = client.BroadcastAsync
	FlagHome               = client.FlagHome
	FlagUseLedger          = client.FlagUseLedger
	FlagChainID            = client.FlagChainID
	FlagNode               = client.FlagNode
	FlagHeight             = client.FlagHeight
	FlagGasAdjustment      = client.FlagGasAdjustment
	FlagTrustNode          = client.FlagTrustNode
	FlagFrom               = client.FlagFrom
	FlagName               = client.FlagName
	FlagAccountNumber      = client.FlagAccountNumber
	FlagSequence           = client.FlagSequence
	FlagMemo               = client.FlagMemo
	FlagFees               = client.FlagFees
	FlagGasPrices          = client.FlagGasPrices
	FlagBroadcastMode      = client.FlagBroadcastMode
	FlagDryRun             = client.FlagDryRun
	FlagGenerateOnly       = client.FlagGenerateOnly
	FlagIndentResponse     = client.FlagIndentResponse
	FlagListenAddr         = client.FlagListenAddr
	FlagMaxOpenConnections = client.FlagMaxOpenConnections
	FlagRPCReadTimeout     = client.FlagRPCReadTimeout
	FlagRPCWriteTimeout    = client.FlagRPCWriteTimeout
	FlagOutputDocument     = client.FlagOutputDocument
	FlagSkipConfirmation   = client.FlagSkipConfirmation
	DefaultKeyPass         = client.DefaultKeyPass
	FlagAddress            = client.FlagAddress
	FlagPublicKey          = client.FlagPublicKey
	FlagBechPrefix         = client.FlagBechPrefix
	FlagDevice             = client.FlagDevice
	OutputFormatText       = client.OutputFormatText
	OutputFormatJSON       = client.OutputFormatJSON
	MinPassLength          = client.MinPassLength
)

var (
	// functions aliases
	NewCLIContextWithFrom              = client.NewCLIContextWithFrom
	NewCLIContext                      = client.NewCLIContext
	GetFromFields                      = client.GetFromFields
	ErrInvalidAccount                  = client.ErrInvalidAccount
	ErrVerifyCommit                    = client.ErrVerifyCommit
	GetCommands                        = client.GetCommands
	PostCommands                       = client.PostCommands
	RegisterRestServerFlags            = client.RegisterRestServerFlags
	ParseGas                           = client.ParseGas
	NewCompletionCmd                   = client.NewCompletionCmd
	MarshalJSON                        = client.MarshalJSON
	UnmarshalJSON                      = client.UnmarshalJSON
	Commands                           = client.Commands
	NewAddNewKey                       = client.NewAddNewKey
	NewRecoverKey                      = client.NewRecoverKey
	NewUpdateKeyReq                    = client.NewUpdateKeyReq
	NewDeleteKeyReq                    = client.NewDeleteKeyReq
	GetKeyInfo                         = client.GetKeyInfo
	GetPassphrase                      = client.GetPassphrase
	ReadPassphraseFromStdin            = client.ReadPassphraseFromStdin
	NewKeyBaseFromHomeFlag             = client.NewKeyBaseFromHomeFlag
	NewKeyBaseFromDir                  = client.NewKeyBaseFromDir
	NewInMemoryKeyBase                 = client.NewInMemoryKeyBase
	NewRestServer                      = client.NewRestServer
	ServeCommand                       = lcd.ServeCommand
	BlockCommand                       = rpc.BlockCommand
	QueryGenesisAccountCmd             = rpc.QueryGenesisAccountCmd
	QueryGenesisTxCmd                  = rpc.QueryGenesisTxCmd
	GetChainHeight                     = rpc.GetChainHeight
	BlockRequestHandlerFn              = rpc.BlockRequestHandlerFn
	LatestBlockRequestHandlerFn        = rpc.LatestBlockRequestHandlerFn
	MempoolCmd                         = rpc.MempoolCmd
	RegisterRPCRoutes                  = rpc.RegisterRPCRoutes
	StatusCommand                      = client.StatusCommand
	NodeInfoRequestHandlerFn           = client.NodeInfoRequestHandlerFn
	NodeSyncingRequestHandlerFn        = client.NodeSyncingRequestHandlerFn
	ValidatorCommand                   = client.ValidatorCommand
	GetValidators                      = client.GetValidators
	ValidatorSetRequestHandlerFn       = client.ValidatorSetRequestHandlerFn
	LatestValidatorSetRequestHandlerFn = client.LatestValidatorSetRequestHandlerFn
	GetPassword                        = client.GetPassword
	GetCheckPassword                   = client.GetCheckPassword
	GetConfirmation                    = client.GetConfirmation
	GetString                          = client.GetString
	PrintPrefixed                      = client.PrintPrefixed

	// variable aliases
	LineBreak  = client.LineBreak
	GasFlagVar = client.GasFlagVar

	ConfigCmd   = client.ConfigCmd
	ValidateCmd = client.ValidateCmd
)

type (
	CLIContext             = client.CLIContext
	GasSetting             = client.GasSetting
	AddNewKey              = client.AddNewKey
	RecoverKey             = client.RecoverKey
	UpdateKeyReq           = client.UpdateKeyReq
	DeleteKeyReq           = client.DeleteKeyReq
	RestServer             = client.RestServer
	ValidatorOutput        = client.ValidatorOutput
	ResultValidatorsOutput = client.ResultValidatorsOutput
)
