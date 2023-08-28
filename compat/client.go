package compat

import (
	"context"

	oscrypto "github.com/Finschia/ostracon/crypto/ed25519"
	ocmerkle "github.com/Finschia/ostracon/crypto/merkle"
	ocbytes "github.com/Finschia/ostracon/libs/bytes"
	oclog "github.com/Finschia/ostracon/libs/log"
	occlient "github.com/Finschia/ostracon/rpc/client"
	ocrpctypes "github.com/Finschia/ostracon/rpc/core/types"
	octypes "github.com/Finschia/ostracon/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"
)

type TMClientWrapper struct {
	client tmclient.Client
}

func NewOCRpcClient(client tmclient.Client) occlient.Client {
	return TMClientWrapper{client: client}
}

// Service

func (t TMClientWrapper) Start() error {
	return t.client.Start()
}
func (t TMClientWrapper) OnStart() error {
	return t.client.OnStart()
}
func (t TMClientWrapper) Stop() error {
	return t.client.Stop()
}
func (t TMClientWrapper) OnStop() {
	t.client.OnStop()
}
func (t TMClientWrapper) Reset() error {
	return t.client.Reset()
}
func (t TMClientWrapper) OnReset() error {
	return t.client.OnReset()
}
func (t TMClientWrapper) IsRunning() bool {
	return t.client.IsRunning()
}
func (t TMClientWrapper) Quit() <-chan struct{} {
	return t.client.Quit()
}
func (t TMClientWrapper) String() string {
	return t.client.String()
}
func (t TMClientWrapper) SetLogger(logger oclog.Logger) {
	t.client.SetLogger(NewTMLogger(logger))
}

// ABCI Client

func (t TMClientWrapper) ABCIInfo(ctx context.Context) (*ocrpctypes.ResultABCIInfo, error) {
	tmResultABCIInfo, err := t.client.ABCIInfo(ctx)
	return &ocrpctypes.ResultABCIInfo{
		Response: tmResultABCIInfo.Response,
	}, err
}

func (t TMClientWrapper) ABCIQuery(ctx context.Context, path string, data ocbytes.HexBytes) (*ocrpctypes.ResultABCIQuery, error) {
	tmResultABCIQuery, err := t.client.ABCIQuery(ctx, path, tmbytes.HexBytes(data))
	return &ocrpctypes.ResultABCIQuery{
		Response: tmResultABCIQuery.Response,
	}, err
}

func (t TMClientWrapper) ABCIQueryWithOptions(ctx context.Context, path string, data ocbytes.HexBytes,
	opts occlient.ABCIQueryOptions) (*ocrpctypes.ResultABCIQuery, error) {
	tmABCIQueryOptions := tmclient.ABCIQueryOptions{
		Height: opts.Height,
		Prove:  opts.Prove,
	}
	tmResultABCIQuery, err := t.client.ABCIQueryWithOptions(ctx, path, tmbytes.HexBytes(data), tmABCIQueryOptions)
	return &ocrpctypes.ResultABCIQuery{
		Response: tmResultABCIQuery.Response,
	}, err
}

func (t TMClientWrapper) BroadcastTxCommit(ctx context.Context, tx octypes.Tx) (*ocrpctypes.ResultBroadcastTxCommit, error) {
	tmResultBroadcastTxCommit, err := t.client.BroadcastTxCommit(ctx, tmtypes.Tx(tx))
	return &ocrpctypes.ResultBroadcastTxCommit{
		CheckTx:   tmResultBroadcastTxCommit.CheckTx,
		DeliverTx: tmResultBroadcastTxCommit.DeliverTx,
		Hash:      ocbytes.HexBytes(tmResultBroadcastTxCommit.Hash),
		Height:    tmResultBroadcastTxCommit.Height,
	}, err
}

func (t TMClientWrapper) BroadcastTxAsync(ctx context.Context, tx octypes.Tx) (*ocrpctypes.ResultBroadcastTx, error) {
	tmResultBroadcastTx, err := t.client.BroadcastTxAsync(ctx, tmtypes.Tx(tx))
	return &ocrpctypes.ResultBroadcastTx{
		Code:      tmResultBroadcastTx.Code,
		Data:      ocbytes.HexBytes(tmResultBroadcastTx.Data),
		Log:       tmResultBroadcastTx.Log,
		Codespace: tmResultBroadcastTx.Codespace,
		// MempoolError does not exist on tm side
		Hash: ocbytes.HexBytes(tmResultBroadcastTx.Hash),
	}, err
}
func (t TMClientWrapper) BroadcastTxSync(ctx context.Context, tx octypes.Tx) (*ocrpctypes.ResultBroadcastTx, error) {
	tmResultBroadcastTx, err := t.client.BroadcastTxSync(ctx, tmtypes.Tx(tx))
	return &ocrpctypes.ResultBroadcastTx{
		Code:      tmResultBroadcastTx.Code,
		Data:      ocbytes.HexBytes(tmResultBroadcastTx.Data),
		Log:       tmResultBroadcastTx.Log,
		Codespace: tmResultBroadcastTx.Codespace,
		// MempoolError does not exist on tm side
		Hash: ocbytes.HexBytes(tmResultBroadcastTx.Hash),
	}, err
}

// SignClient
func (t TMClientWrapper) Block(ctx context.Context, height *int64) (*ocrpctypes.ResultBlock, error) {
	tmResultBlock, err := t.client.Block(ctx, height)
	return &ocrpctypes.ResultBlock{
		BlockID: OCBlockIDFrom(tmResultBlock.BlockID),
		Block:   OCBlockFrom(tmResultBlock.Block),
	}, err
}

func (t TMClientWrapper) BlockByHash(ctx context.Context, hash []byte) (*ocrpctypes.ResultBlock, error) {
	tmResultBlock, err := t.client.BlockByHash(ctx, hash)
	return &ocrpctypes.ResultBlock{
		BlockID: OCBlockIDFrom(tmResultBlock.BlockID),
		Block:   OCBlockFrom(tmResultBlock.Block),
	}, err
}

func (t TMClientWrapper) BlockResults(ctx context.Context, height *int64) (*ocrpctypes.ResultBlockResults, error) {
	tmResultBlockResults, err := t.client.BlockResults(ctx, height)
	return &ocrpctypes.ResultBlockResults{
		Height:                tmResultBlockResults.Height,
		TxsResults:            tmResultBlockResults.TxsResults,
		BeginBlockEvents:      tmResultBlockResults.BeginBlockEvents,
		EndBlockEvents:        tmResultBlockResults.EndBlockEvents,
		ValidatorUpdates:      tmResultBlockResults.ValidatorUpdates,
		ConsensusParamUpdates: tmResultBlockResults.ConsensusParamUpdates,
	}, err
}

func (t TMClientWrapper) Commit(ctx context.Context, height *int64) (*ocrpctypes.ResultCommit, error) {
	tmResultCommit, err := t.client.Commit(ctx, height)
	ocHeader := OCHeaderFrom(*tmResultCommit.Header)
	return &ocrpctypes.ResultCommit{
		SignedHeader: octypes.SignedHeader{
			Header: &ocHeader,
			Commit: OCCommitFrom(tmResultCommit.Commit),
		},
		CanonicalCommit: tmResultCommit.CanonicalCommit,
	}, err
}

func (t TMClientWrapper) Validators(ctx context.Context, height *int64, page, perPage *int) (*ocrpctypes.ResultValidators, error) {
	tmResultValidators, err := t.client.Validators(ctx, height, page, perPage)
	return &ocrpctypes.ResultValidators{
		BlockHeight: tmResultValidators.BlockHeight,
		Validators:  OCValidatorsFrom(tmResultValidators.Validators),
		Count:       tmResultValidators.Count,
		Total:       tmResultValidators.Total,
	}, err
}

func (t TMClientWrapper) Tx(ctx context.Context, hash []byte, prove bool) (*ocrpctypes.ResultTx, error) {
	tmResultTx, err := t.client.Tx(ctx, hash, prove)
	return &ocrpctypes.ResultTx{
		Hash:     ocbytes.HexBytes(tmResultTx.Hash),
		Height:   tmResultTx.Height,
		Index:    tmResultTx.Index,
		TxResult: tmResultTx.TxResult,
		Tx:       octypes.Tx(tmResultTx.Tx),
		Proof: octypes.TxProof{
			RootHash: ocbytes.HexBytes(tmResultTx.Proof.RootHash),
			Data:     octypes.Tx(tmResultTx.Proof.Data),
			Proof: ocmerkle.Proof{
				Total:    tmResultTx.Proof.Proof.Total,
				Index:    tmResultTx.Proof.Proof.Index,
				LeafHash: tmResultTx.Proof.Proof.LeafHash,
				Aunts:    tmResultTx.Proof.Proof.Aunts,
			},
		},
	}, err
}

func (t TMClientWrapper) TxSearch(
	ctx context.Context,
	query string,
	prove bool,
	page, perPage *int,
	orderBy string,
) (*ocrpctypes.ResultTxSearch, error) {
	tmResultTxSearch, err := t.client.TxSearch(ctx, query, prove, page, perPage, orderBy)
	ocResultTxs := make([]*ocrpctypes.ResultTx, len(tmResultTxSearch.Txs))
	for _, v := range tmResultTxSearch.Txs {
		ocResultTx := &ocrpctypes.ResultTx{
			Hash:     ocbytes.HexBytes(v.Hash),
			Height:   v.Height,
			Index:    v.Index,
			TxResult: v.TxResult,
			Tx:       octypes.Tx(v.Tx),
			Proof:    OCTxProofFrom(v.Proof),
		}
		ocResultTxs = append(ocResultTxs, ocResultTx)
	}
	return &ocrpctypes.ResultTxSearch{
		Txs:        ocResultTxs,
		TotalCount: tmResultTxSearch.TotalCount,
	}, err
}

func (t TMClientWrapper) BlockSearch(
	ctx context.Context,
	query string,
	page, perPage *int,
	orderBy string,
) (*ocrpctypes.ResultBlockSearch, error) {
	tmResultBlockSearch, err := t.client.BlockSearch(ctx, query, page, perPage, orderBy)
	ocResultBlocks := make([]*ocrpctypes.ResultBlock, len(tmResultBlockSearch.Blocks))
	for _, v := range tmResultBlockSearch.Blocks {
		ocResultBlock := &ocrpctypes.ResultBlock{
			BlockID: OCBlockIDFrom(v.BlockID),
			Block:   OCBlockFrom(v.Block),
		}
		ocResultBlocks = append(ocResultBlocks, ocResultBlock)
	}
	return &ocrpctypes.ResultBlockSearch{
		Blocks:     ocResultBlocks,
		TotalCount: tmResultBlockSearch.TotalCount,
	}, err
}

// HistoryClient

func (t TMClientWrapper) Genesis(ctx context.Context) (*ocrpctypes.ResultGenesis, error) {
	tmResultGenesis, err := t.client.Genesis(ctx)
	return &ocrpctypes.ResultGenesis{
		Genesis: NewOCGenesisDoc(tmResultGenesis.Genesis),
	}, err
}

func (t TMClientWrapper) GenesisChunked(ctx context.Context, uint uint) (*ocrpctypes.ResultGenesisChunk, error) {
	tmResultGenesisChunk, err := t.client.GenesisChunked(ctx, uint)
	return &ocrpctypes.ResultGenesisChunk{
		ChunkNumber: tmResultGenesisChunk.ChunkNumber,
		TotalChunks: tmResultGenesisChunk.TotalChunks,
		Data:        tmResultGenesisChunk.Data,
	}, err
}

func (t TMClientWrapper) BlockchainInfo(ctx context.Context, minHeight, maxHeight int64) (*ocrpctypes.ResultBlockchainInfo, error) {
	tmResultBlockchainInfo, err := t.client.BlockchainInfo(ctx, minHeight, maxHeight)
	return &ocrpctypes.ResultBlockchainInfo{
		LastHeight: tmResultBlockchainInfo.LastHeight,
		BlockMetas: OCBlockMetasFrom(tmResultBlockchainInfo.BlockMetas),
	}, err
}

// StatusClient

func (t TMClientWrapper) Status(ctx context.Context) (*ocrpctypes.ResultStatus, error) {
	tmResultStatus, err := t.client.Status(ctx)
	var pubKey oscrypto.PubKey
	copy(pubKey, tmResultStatus.ValidatorInfo.PubKey.Bytes())
	return &ocrpctypes.ResultStatus{
		NodeInfo: OCNodeInfoFrom(tmResultStatus.NodeInfo),
		SyncInfo: ocrpctypes.SyncInfo{
			LatestBlockHash:     ocbytes.HexBytes(tmResultStatus.SyncInfo.LatestBlockHash),
			LatestAppHash:       ocbytes.HexBytes(tmResultStatus.SyncInfo.LatestAppHash),
			LatestBlockHeight:   tmResultStatus.SyncInfo.LatestBlockHeight,
			LatestBlockTime:     tmResultStatus.SyncInfo.LatestBlockTime,
			EarliestBlockHash:   ocbytes.HexBytes(tmResultStatus.SyncInfo.EarliestBlockHash),
			EarliestAppHash:     ocbytes.HexBytes(tmResultStatus.SyncInfo.EarliestAppHash),
			EarliestBlockHeight: tmResultStatus.SyncInfo.EarliestBlockHeight,
			EarliestBlockTime:   tmResultStatus.SyncInfo.EarliestBlockTime,
			CatchingUp:          tmResultStatus.SyncInfo.CatchingUp,
		},
		ValidatorInfo: ocrpctypes.ValidatorInfo{
			Address:     ocbytes.HexBytes(tmResultStatus.ValidatorInfo.Address),
			PubKey:      pubKey,
			VotingPower: tmResultStatus.ValidatorInfo.VotingPower,
		},
	}, err
}

// NetworkClient

func (t TMClientWrapper) NetInfo(ctx context.Context) (*ocrpctypes.ResultNetInfo, error) {
	tmResultNetInfo, err := t.client.NetInfo(ctx)
	return &ocrpctypes.ResultNetInfo{
		Listening: tmResultNetInfo.Listening,
		Listeners: tmResultNetInfo.Listeners,
		NPeers:    tmResultNetInfo.NPeers,
		Peers:     OCPeersFrom(tmResultNetInfo.Peers),
	}, err
}

func (t TMClientWrapper) DumpConsensusState(ctx context.Context) (*ocrpctypes.ResultDumpConsensusState, error) {
	tmResultDumpConsensusState, err := t.client.DumpConsensusState(ctx)
	ocPeers := make([]ocrpctypes.PeerStateInfo, len(tmResultDumpConsensusState.Peers))
	for _, v := range tmResultDumpConsensusState.Peers {
		ocPeer := ocrpctypes.PeerStateInfo{
			NodeAddress: v.NodeAddress,
			PeerState:   v.PeerState,
		}
		ocPeers = append(ocPeers, ocPeer)
	}
	return &ocrpctypes.ResultDumpConsensusState{
		RoundState: tmResultDumpConsensusState.RoundState,
		Peers:      ocPeers,
	}, err
}

func (t TMClientWrapper) ConsensusState(ctx context.Context) (*ocrpctypes.ResultConsensusState, error) {
	tmResultConsensusState, err := t.client.ConsensusState(ctx)
	return &ocrpctypes.ResultConsensusState{
		RoundState: tmResultConsensusState.RoundState,
	}, err
}

func (t TMClientWrapper) ConsensusParams(ctx context.Context, height *int64) (*ocrpctypes.ResultConsensusParams, error) {
	tmResultConsensusParams, err := t.client.ConsensusParams(ctx, height)
	return &ocrpctypes.ResultConsensusParams{
		BlockHeight:     tmResultConsensusParams.BlockHeight,
		ConsensusParams: tmResultConsensusParams.ConsensusParams,
	}, err
}

func (t TMClientWrapper) Health(ctx context.Context) (*ocrpctypes.ResultHealth, error) {
	_, err := t.client.Health(ctx)
	return &ocrpctypes.ResultHealth{}, err
}

// EventsClient

func (t TMClientWrapper) Subscribe(ctx context.Context, subscriber, query string, outCapacity ...int) (<-chan ocrpctypes.ResultEvent, error) {
	tmResultUnconfirmedTxs, err := t.client.Subscribe(ctx, subscriber, query, outCapacity...)
	out := make(chan ocrpctypes.ResultEvent)
	for msg := range tmResultUnconfirmedTxs {
		out <- ocrpctypes.ResultEvent{
			Query:  msg.Query,
			Data:   msg.Data,
			Events: msg.Events,
		}
	}
	return out, err
}

func (t TMClientWrapper) Unsubscribe(ctx context.Context, subscriber, query string) error {
	return t.client.Unsubscribe(ctx, subscriber, query)
}

func (t TMClientWrapper) UnsubscribeAll(ctx context.Context, subscriber string) error {
	return t.client.UnsubscribeAll(ctx, subscriber)
}

// MempoolClient

func (t TMClientWrapper) UnconfirmedTxs(ctx context.Context, limit *int) (*ocrpctypes.ResultUnconfirmedTxs, error) {
	tmResultUnconfirmedTxs, err := t.client.UnconfirmedTxs(ctx, limit)
	return &ocrpctypes.ResultUnconfirmedTxs{
		Count:      tmResultUnconfirmedTxs.Count,
		Total:      tmResultUnconfirmedTxs.Total,
		TotalBytes: tmResultUnconfirmedTxs.TotalBytes,
		Txs:        OCTxsFrom(tmResultUnconfirmedTxs.Txs),
	}, err
}

func (t TMClientWrapper) NumUnconfirmedTxs(ctx context.Context) (*ocrpctypes.ResultUnconfirmedTxs, error) {
	tmResultUnconfirmedTxs, err := t.client.NumUnconfirmedTxs(ctx)
	return &ocrpctypes.ResultUnconfirmedTxs{
		Count:      tmResultUnconfirmedTxs.Count,
		Total:      tmResultUnconfirmedTxs.Total,
		TotalBytes: tmResultUnconfirmedTxs.TotalBytes,
		Txs:        OCTxsFrom(tmResultUnconfirmedTxs.Txs),
	}, err
}

func (t TMClientWrapper) CheckTx(ctx context.Context, tx octypes.Tx) (*ocrpctypes.ResultCheckTx, error) {
	tmResultCheckTx, err := t.client.CheckTx(ctx, tmtypes.Tx(tx))
	return &ocrpctypes.ResultCheckTx{
		ResponseCheckTx: tmResultCheckTx.ResponseCheckTx,
	}, err
}

// EvidenceClient

func (t TMClientWrapper) BroadcastEvidence(ctx context.Context, evidence octypes.Evidence) (*ocrpctypes.ResultBroadcastEvidence, error) {
	tmResultBroadcastEvidence, err := t.client.BroadcastEvidence(ctx, evidence)
	return &ocrpctypes.ResultBroadcastEvidence{
		Hash: tmResultBroadcastEvidence.Hash,
	}, err
}
