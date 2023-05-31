package compat

import (
	osconfig "github.com/Finschia/ostracon/config"
	oslog "github.com/Finschia/ostracon/libs/log"
	ostypes "github.com/Finschia/ostracon/types"
	tmabcicli "github.com/tendermint/tendermint/abci/client"
	tmabcitypes "github.com/tendermint/tendermint/abci/types"
	tmconfig "github.com/tendermint/tendermint/config"
	tmcrypto "github.com/tendermint/tendermint/crypto/ed25519"
	tmlog "github.com/tendermint/tendermint/libs/log"
	tmproxy "github.com/tendermint/tendermint/proxy"
	tmtypes "github.com/tendermint/tendermint/types"
)

var _ tmproxy.ClientCreator = (*clientCreator)(nil)

type clientCreator struct {
	app tmabcitypes.Application
}

func (c clientCreator) NewABCIClient() (tmabcicli.Client, error) {
	return tmabcicli.NewLocalClient(nil, c.app), nil
}

func NewTMClientCreator(app tmabcitypes.Application) tmproxy.ClientCreator {
	return clientCreator{app: app}
}

func NewTMGenesisDoc(d *ostypes.GenesisDoc) *tmtypes.GenesisDoc {
	vals := make([]tmtypes.GenesisValidator, len(d.Validators))
	for i := range d.Validators {
		var pubKey tmcrypto.PubKey
		copy(pubKey, d.Validators[i].PubKey.Bytes())
		vals[i].Address = d.Validators[i].Address.Bytes()
		vals[i].PubKey = pubKey
		vals[i].Power = d.Validators[i].Power
		vals[i].Name = d.Validators[i].Name
	}

	var doc tmtypes.GenesisDoc
	doc.GenesisTime = d.GenesisTime
	doc.ChainID = d.ChainID
	doc.InitialHeight = d.InitialHeight
	doc.ConsensusParams = d.ConsensusParams
	doc.Validators = vals
	doc.AppHash = d.AppHash.Bytes()
	doc.AppState = d.AppState
	return &doc
}

func NewTMRPCConfig(c *osconfig.RPCConfig) *tmconfig.RPCConfig {
	var tmcfg tmconfig.RPCConfig
	tmcfg.RootDir = c.RootDir
	tmcfg.ListenAddress = c.ListenAddress
	tmcfg.CORSAllowedOrigins = c.CORSAllowedOrigins
	tmcfg.CORSAllowedMethods = c.CORSAllowedMethods
	tmcfg.CORSAllowedHeaders = c.CORSAllowedHeaders
	tmcfg.GRPCListenAddress = c.GRPCListenAddress
	tmcfg.GRPCMaxOpenConnections = c.GRPCMaxOpenConnections
	//tmcfg.ReadTimeout = c.ReadTimeout
	//tmcfg.WriteTimeout = c.WriteTimeout
	//tmcfg.IdleTimeout = c.IdleTimeout
	tmcfg.Unsafe = c.Unsafe
	tmcfg.MaxOpenConnections = c.MaxOpenConnections
	tmcfg.MaxSubscriptionClients = c.MaxSubscriptionClients
	tmcfg.MaxSubscriptionsPerClient = c.MaxSubscriptionsPerClient
	tmcfg.SubscriptionBufferSize = c.SubscriptionBufferSize
	tmcfg.WebSocketWriteBufferSize = c.WebSocketWriteBufferSize
	tmcfg.CloseOnSlowClient = c.CloseOnSlowClient
	tmcfg.TimeoutBroadcastTxCommit = c.TimeoutBroadcastTxCommit
	tmcfg.MaxBodyBytes = c.MaxBodyBytes
	tmcfg.MaxHeaderBytes = c.MaxHeaderBytes
	tmcfg.TLSCertFile = c.TLSCertFile
	tmcfg.TLSKeyFile = c.TLSKeyFile
	tmcfg.PprofListenAddress = c.PprofListenAddress
	return &tmcfg
}

var _ tmlog.Logger = (*Logger)(nil)

type Logger struct {
	logger oslog.Logger
}

func (l Logger) Debug(msg string, keyvals ...interface{}) {
	l.logger.Debug(msg, keyvals...)
}

func (l Logger) Info(msg string, keyvals ...interface{}) {
	l.logger.Info(msg, keyvals...)
}

func (l Logger) Error(msg string, keyvals ...interface{}) {
	l.logger.Error(msg, keyvals...)
}

func (l Logger) With(keyvals ...interface{}) tmlog.Logger {
	return NewTMLogger(l.logger.With(keyvals...))
}

func NewTMLogger(l oslog.Logger) tmlog.Logger {
	return Logger{logger: l}
}

func (l *Logger) GetOSLogger() *oslog.Logger {
	return &l.logger
}
