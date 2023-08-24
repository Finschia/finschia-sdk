package compat

import (
	osconfig "github.com/Finschia/ostracon/config"
	oscrypto "github.com/Finschia/ostracon/crypto/ed25519"
	ocmerkle "github.com/Finschia/ostracon/crypto/merkle"
	osbytes "github.com/Finschia/ostracon/libs/bytes"
	osflow "github.com/Finschia/ostracon/libs/flowrate"
	oslog "github.com/Finschia/ostracon/libs/log"
	osp2p "github.com/Finschia/ostracon/p2p"
	osp2pconn "github.com/Finschia/ostracon/p2p/conn"
	ocrpctypes "github.com/Finschia/ostracon/rpc/core/types"
	ostypes "github.com/Finschia/ostracon/types"
	tmabcicli "github.com/tendermint/tendermint/abci/client"
	tmabcitypes "github.com/tendermint/tendermint/abci/types"
	tmconfig "github.com/tendermint/tendermint/config"
	tmcrypto "github.com/tendermint/tendermint/crypto/ed25519"
	tmlog "github.com/tendermint/tendermint/libs/log"
	tmp2p "github.com/tendermint/tendermint/p2p"
	tmproxy "github.com/tendermint/tendermint/proxy"
	tmrpctypes "github.com/tendermint/tendermint/rpc/core/types"
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

func OCBlockIDFrom(block tmtypes.BlockID) ostypes.BlockID {
	return ostypes.BlockID{
		Hash: osbytes.HexBytes(block.Hash),
		PartSetHeader: ostypes.PartSetHeader{
			Total: block.PartSetHeader.Total,
			Hash:  osbytes.HexBytes(block.PartSetHeader.Hash),
		},
	}
}

func OCBlockFrom(block *tmtypes.Block) *ostypes.Block {
	var osEvidenceList ostypes.EvidenceList
	for _, v := range block.Evidence.Evidence {
		osEvidenceList = append(osEvidenceList, ostypes.Evidence(v))
	}

	return &ostypes.Block{
		Header: OCHeaderFrom(block.Header),
		Data: ostypes.Data{
			Txs: OCTxsFrom(block.Txs),
		},
		Evidence: ostypes.EvidenceData{
			Evidence: osEvidenceList,
		},
		LastCommit: OCCommitFrom(block.LastCommit),
	}
}

func OCBlockMetasFrom(blockMetas []*tmtypes.BlockMeta) []*ostypes.BlockMeta {
	var osBlockMetas []*ostypes.BlockMeta
	for _, v := range blockMetas {
		osBlockMeta := &ostypes.BlockMeta{
			BlockID:   OCBlockIDFrom(v.BlockID),
			BlockSize: v.BlockSize,
			Header:    OCHeaderFrom(v.Header),
			NumTxs:    v.NumTxs,
		}
		osBlockMetas = append(osBlockMetas, osBlockMeta)
	}

	return osBlockMetas
}

func OCHeaderFrom(header tmtypes.Header) ostypes.Header {
	return ostypes.Header{
		Version:            header.Version,
		ChainID:            header.ChainID,
		Height:             header.Height,
		Time:               header.Time,
		LastBlockID:        OCBlockIDFrom(header.LastBlockID),
		LastCommitHash:     osbytes.HexBytes(header.LastCommitHash),
		DataHash:           osbytes.HexBytes(header.DataHash),
		ValidatorsHash:     osbytes.HexBytes(header.ValidatorsHash),
		NextValidatorsHash: osbytes.HexBytes(header.NextValidatorsHash),
		ConsensusHash:      osbytes.HexBytes(header.ConsensusHash),
		AppHash:            osbytes.HexBytes(header.AppHash),
		LastResultsHash:    osbytes.HexBytes(header.LastResultsHash),
		EvidenceHash:       osbytes.HexBytes(header.EvidenceHash),
		ProposerAddress:    osbytes.HexBytes(header.ProposerAddress),
	}
}

func OCTxsFrom(txs tmtypes.Txs) ostypes.Txs {
	var osTxs ostypes.Txs
	for _, v := range txs {
		osTxs = append(osTxs, ostypes.Tx(v))
	}
	return osTxs
}

func OCCommitFrom(commit *tmtypes.Commit) *ostypes.Commit {
	var osSignatures []ostypes.CommitSig
	for _, v := range commit.Signatures {
		osCommitSig := ostypes.CommitSig{
			BlockIDFlag:      ostypes.BlockIDFlag(v.BlockIDFlag),
			ValidatorAddress: ostypes.Address(v.ValidatorAddress),
			Timestamp:        v.Timestamp,
			Signature:        v.Signature,
		}
		osSignatures = append(osSignatures, osCommitSig)
	}

	return &ostypes.Commit{
		Height:     commit.Height,
		Round:      commit.Round,
		BlockID:    OCBlockIDFrom(commit.BlockID),
		Signatures: osSignatures,
	}
}

func OCValidatorsFrom(validators []*tmtypes.Validator) []*ostypes.Validator {
	var osValidators []*ostypes.Validator
	for _, v := range validators {
		var pubKey oscrypto.PubKey
		copy(pubKey, v.PubKey.Bytes())
		osValidator := &ostypes.Validator{
			Address:          ostypes.Address(v.Address),
			PubKey:           pubKey,
			VotingPower:      v.VotingPower,
			ProposerPriority: v.ProposerPriority,
		}
		osValidators = append(osValidators, osValidator)
	}
	return osValidators
}

func OCNodeInfoFrom(nodeInfo tmp2p.DefaultNodeInfo) osp2p.DefaultNodeInfo {
	return osp2p.DefaultNodeInfo{
		ProtocolVersion: osp2p.ProtocolVersion{
			P2P:   nodeInfo.ProtocolVersion.P2P,
			Block: nodeInfo.ProtocolVersion.Block,
			App:   nodeInfo.ProtocolVersion.App,
		},
		DefaultNodeID: osp2p.ID(nodeInfo.DefaultNodeID),
		ListenAddr:    nodeInfo.ListenAddr,
		Network:       nodeInfo.Network,
		Version:       nodeInfo.Version,
		Channels:      osbytes.HexBytes(nodeInfo.Channels),
		Moniker:       nodeInfo.Moniker,
		Other: osp2p.DefaultNodeInfoOther{
			TxIndex:    nodeInfo.Other.TxIndex,
			RPCAddress: nodeInfo.Other.RPCAddress,
		},
	}
}

func OCTxProofFrom(proof tmtypes.TxProof) ostypes.TxProof {
	return ostypes.TxProof{
		RootHash: osbytes.HexBytes(proof.RootHash),
		Data:     ostypes.Tx(proof.Data),
		Proof: ocmerkle.Proof{
			Total:    proof.Proof.Total,
			Index:    proof.Proof.Index,
			LeafHash: proof.Proof.LeafHash,
			Aunts:    proof.Proof.Aunts,
		},
	}
}

func OCConnectionStatusFrom(connectionStatus tmp2p.ConnectionStatus) osp2p.ConnectionStatus {
	var osChannels []osp2pconn.ChannelStatus
	for _, v := range connectionStatus.Channels {
		osChannel := osp2pconn.ChannelStatus{
			ID:                v.ID,
			SendQueueCapacity: v.SendQueueCapacity,
			SendQueueSize:     v.SendQueueSize,
			Priority:          v.Priority,
			RecentlySent:      v.RecentlySent,
		}
		osChannels = append(osChannels, osChannel)
	}
	return osp2p.ConnectionStatus{
		Duration: connectionStatus.Duration,
		SendMonitor: osflow.Status{
			Start:    connectionStatus.SendMonitor.Start,
			Bytes:    connectionStatus.SendMonitor.Bytes,
			Samples:  connectionStatus.SendMonitor.Samples,
			InstRate: connectionStatus.SendMonitor.InstRate,
			CurRate:  connectionStatus.SendMonitor.CurRate,
			AvgRate:  connectionStatus.SendMonitor.AvgRate,
			PeakRate: connectionStatus.SendMonitor.PeakRate,
			BytesRem: connectionStatus.SendMonitor.BytesRem,
			Duration: connectionStatus.SendMonitor.Duration,
			Idle:     connectionStatus.SendMonitor.Idle,
			TimeRem:  connectionStatus.SendMonitor.TimeRem,
			Progress: osflow.Percent(connectionStatus.SendMonitor.Progress),
			Active:   connectionStatus.SendMonitor.Active,
		},
		RecvMonitor: osflow.Status{
			Start:    connectionStatus.RecvMonitor.Start,
			Bytes:    connectionStatus.RecvMonitor.Bytes,
			Samples:  connectionStatus.RecvMonitor.Samples,
			InstRate: connectionStatus.RecvMonitor.InstRate,
			CurRate:  connectionStatus.RecvMonitor.CurRate,
			AvgRate:  connectionStatus.RecvMonitor.AvgRate,
			PeakRate: connectionStatus.RecvMonitor.PeakRate,
			BytesRem: connectionStatus.RecvMonitor.BytesRem,
			Duration: connectionStatus.RecvMonitor.Duration,
			Idle:     connectionStatus.RecvMonitor.Idle,
			TimeRem:  connectionStatus.RecvMonitor.TimeRem,
			Progress: osflow.Percent(connectionStatus.RecvMonitor.Progress),
			Active:   connectionStatus.RecvMonitor.Active,
		},
		Channels: osChannels,
	}
}

func OCPeersFrom(peers []tmrpctypes.Peer) []ocrpctypes.Peer {
	var osPeers []ocrpctypes.Peer
	for _, v := range peers {
		osPeer := ocrpctypes.Peer{
			NodeInfo:         OCNodeInfoFrom(v.NodeInfo),
			IsOutbound:       v.IsOutbound,
			ConnectionStatus: OCConnectionStatusFrom(v.ConnectionStatus),
			RemoteIP:         v.RemoteIP,
		}
		osPeers = append(osPeers, osPeer)
	}
	return osPeers
}

func NewOCGenesisDoc(d *tmtypes.GenesisDoc) *ostypes.GenesisDoc {
	vals := make([]ostypes.GenesisValidator, len(d.Validators))
	for i := range d.Validators {
		var pubKey oscrypto.PubKey
		copy(pubKey, d.Validators[i].PubKey.Bytes())
		vals[i].Address = d.Validators[i].Address.Bytes()
		vals[i].PubKey = pubKey
		vals[i].Power = d.Validators[i].Power
		vals[i].Name = d.Validators[i].Name
	}

	var doc ostypes.GenesisDoc
	doc.GenesisTime = d.GenesisTime
	doc.ChainID = d.ChainID
	doc.InitialHeight = d.InitialHeight
	doc.ConsensusParams = d.ConsensusParams
	doc.Validators = vals
	doc.AppHash = d.AppHash.Bytes()
	doc.AppState = d.AppState
	return &doc
}
