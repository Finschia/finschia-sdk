<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [ibc/applications/transfer/v1/transfer.proto](#ibc/applications/transfer/v1/transfer.proto)
    - [DenomTrace](#ibc.applications.transfer.v1.DenomTrace)
    - [FungibleTokenPacketData](#ibc.applications.transfer.v1.FungibleTokenPacketData)
    - [Params](#ibc.applications.transfer.v1.Params)
  
- [ibc/applications/transfer/v1/genesis.proto](#ibc/applications/transfer/v1/genesis.proto)
    - [GenesisState](#ibc.applications.transfer.v1.GenesisState)
  
- [lbm/base/query/v1/pagination.proto](#lbm/base/query/v1/pagination.proto)
    - [PageRequest](#lbm.base.query.v1.PageRequest)
    - [PageResponse](#lbm.base.query.v1.PageResponse)
  
- [ibc/applications/transfer/v1/query.proto](#ibc/applications/transfer/v1/query.proto)
    - [QueryDenomTraceRequest](#ibc.applications.transfer.v1.QueryDenomTraceRequest)
    - [QueryDenomTraceResponse](#ibc.applications.transfer.v1.QueryDenomTraceResponse)
    - [QueryDenomTracesRequest](#ibc.applications.transfer.v1.QueryDenomTracesRequest)
    - [QueryDenomTracesResponse](#ibc.applications.transfer.v1.QueryDenomTracesResponse)
    - [QueryParamsRequest](#ibc.applications.transfer.v1.QueryParamsRequest)
    - [QueryParamsResponse](#ibc.applications.transfer.v1.QueryParamsResponse)
  
    - [Query](#ibc.applications.transfer.v1.Query)
  
- [lbm/base/v1/coin.proto](#lbm/base/v1/coin.proto)
    - [Coin](#lbm.base.v1.Coin)
    - [DecCoin](#lbm.base.v1.DecCoin)
    - [DecProto](#lbm.base.v1.DecProto)
    - [IntProto](#lbm.base.v1.IntProto)
  
- [ibc/core/client/v1/client.proto](#ibc/core/client/v1/client.proto)
    - [ClientConsensusStates](#ibc.core.client.v1.ClientConsensusStates)
    - [ClientUpdateProposal](#ibc.core.client.v1.ClientUpdateProposal)
    - [ConsensusStateWithHeight](#ibc.core.client.v1.ConsensusStateWithHeight)
    - [Height](#ibc.core.client.v1.Height)
    - [IdentifiedClientState](#ibc.core.client.v1.IdentifiedClientState)
    - [Params](#ibc.core.client.v1.Params)
  
- [ibc/applications/transfer/v1/tx.proto](#ibc/applications/transfer/v1/tx.proto)
    - [MsgTransfer](#ibc.applications.transfer.v1.MsgTransfer)
    - [MsgTransferResponse](#ibc.applications.transfer.v1.MsgTransferResponse)
  
    - [Msg](#ibc.applications.transfer.v1.Msg)
  
- [ibc/core/channel/v1/channel.proto](#ibc/core/channel/v1/channel.proto)
    - [Acknowledgement](#ibc.core.channel.v1.Acknowledgement)
    - [Channel](#ibc.core.channel.v1.Channel)
    - [Counterparty](#ibc.core.channel.v1.Counterparty)
    - [IdentifiedChannel](#ibc.core.channel.v1.IdentifiedChannel)
    - [Packet](#ibc.core.channel.v1.Packet)
    - [PacketState](#ibc.core.channel.v1.PacketState)
  
    - [Order](#ibc.core.channel.v1.Order)
    - [State](#ibc.core.channel.v1.State)
  
- [ibc/core/channel/v1/genesis.proto](#ibc/core/channel/v1/genesis.proto)
    - [GenesisState](#ibc.core.channel.v1.GenesisState)
    - [PacketSequence](#ibc.core.channel.v1.PacketSequence)
  
- [ibc/core/channel/v1/query.proto](#ibc/core/channel/v1/query.proto)
    - [QueryChannelClientStateRequest](#ibc.core.channel.v1.QueryChannelClientStateRequest)
    - [QueryChannelClientStateResponse](#ibc.core.channel.v1.QueryChannelClientStateResponse)
    - [QueryChannelConsensusStateRequest](#ibc.core.channel.v1.QueryChannelConsensusStateRequest)
    - [QueryChannelConsensusStateResponse](#ibc.core.channel.v1.QueryChannelConsensusStateResponse)
    - [QueryChannelRequest](#ibc.core.channel.v1.QueryChannelRequest)
    - [QueryChannelResponse](#ibc.core.channel.v1.QueryChannelResponse)
    - [QueryChannelsRequest](#ibc.core.channel.v1.QueryChannelsRequest)
    - [QueryChannelsResponse](#ibc.core.channel.v1.QueryChannelsResponse)
    - [QueryConnectionChannelsRequest](#ibc.core.channel.v1.QueryConnectionChannelsRequest)
    - [QueryConnectionChannelsResponse](#ibc.core.channel.v1.QueryConnectionChannelsResponse)
    - [QueryNextSequenceReceiveRequest](#ibc.core.channel.v1.QueryNextSequenceReceiveRequest)
    - [QueryNextSequenceReceiveResponse](#ibc.core.channel.v1.QueryNextSequenceReceiveResponse)
    - [QueryPacketAcknowledgementRequest](#ibc.core.channel.v1.QueryPacketAcknowledgementRequest)
    - [QueryPacketAcknowledgementResponse](#ibc.core.channel.v1.QueryPacketAcknowledgementResponse)
    - [QueryPacketAcknowledgementsRequest](#ibc.core.channel.v1.QueryPacketAcknowledgementsRequest)
    - [QueryPacketAcknowledgementsResponse](#ibc.core.channel.v1.QueryPacketAcknowledgementsResponse)
    - [QueryPacketCommitmentRequest](#ibc.core.channel.v1.QueryPacketCommitmentRequest)
    - [QueryPacketCommitmentResponse](#ibc.core.channel.v1.QueryPacketCommitmentResponse)
    - [QueryPacketCommitmentsRequest](#ibc.core.channel.v1.QueryPacketCommitmentsRequest)
    - [QueryPacketCommitmentsResponse](#ibc.core.channel.v1.QueryPacketCommitmentsResponse)
    - [QueryPacketReceiptRequest](#ibc.core.channel.v1.QueryPacketReceiptRequest)
    - [QueryPacketReceiptResponse](#ibc.core.channel.v1.QueryPacketReceiptResponse)
    - [QueryUnreceivedAcksRequest](#ibc.core.channel.v1.QueryUnreceivedAcksRequest)
    - [QueryUnreceivedAcksResponse](#ibc.core.channel.v1.QueryUnreceivedAcksResponse)
    - [QueryUnreceivedPacketsRequest](#ibc.core.channel.v1.QueryUnreceivedPacketsRequest)
    - [QueryUnreceivedPacketsResponse](#ibc.core.channel.v1.QueryUnreceivedPacketsResponse)
  
    - [Query](#ibc.core.channel.v1.Query)
  
- [ibc/core/channel/v1/tx.proto](#ibc/core/channel/v1/tx.proto)
    - [MsgAcknowledgement](#ibc.core.channel.v1.MsgAcknowledgement)
    - [MsgAcknowledgementResponse](#ibc.core.channel.v1.MsgAcknowledgementResponse)
    - [MsgChannelCloseConfirm](#ibc.core.channel.v1.MsgChannelCloseConfirm)
    - [MsgChannelCloseConfirmResponse](#ibc.core.channel.v1.MsgChannelCloseConfirmResponse)
    - [MsgChannelCloseInit](#ibc.core.channel.v1.MsgChannelCloseInit)
    - [MsgChannelCloseInitResponse](#ibc.core.channel.v1.MsgChannelCloseInitResponse)
    - [MsgChannelOpenAck](#ibc.core.channel.v1.MsgChannelOpenAck)
    - [MsgChannelOpenAckResponse](#ibc.core.channel.v1.MsgChannelOpenAckResponse)
    - [MsgChannelOpenConfirm](#ibc.core.channel.v1.MsgChannelOpenConfirm)
    - [MsgChannelOpenConfirmResponse](#ibc.core.channel.v1.MsgChannelOpenConfirmResponse)
    - [MsgChannelOpenInit](#ibc.core.channel.v1.MsgChannelOpenInit)
    - [MsgChannelOpenInitResponse](#ibc.core.channel.v1.MsgChannelOpenInitResponse)
    - [MsgChannelOpenTry](#ibc.core.channel.v1.MsgChannelOpenTry)
    - [MsgChannelOpenTryResponse](#ibc.core.channel.v1.MsgChannelOpenTryResponse)
    - [MsgRecvPacket](#ibc.core.channel.v1.MsgRecvPacket)
    - [MsgRecvPacketResponse](#ibc.core.channel.v1.MsgRecvPacketResponse)
    - [MsgTimeout](#ibc.core.channel.v1.MsgTimeout)
    - [MsgTimeoutOnClose](#ibc.core.channel.v1.MsgTimeoutOnClose)
    - [MsgTimeoutOnCloseResponse](#ibc.core.channel.v1.MsgTimeoutOnCloseResponse)
    - [MsgTimeoutResponse](#ibc.core.channel.v1.MsgTimeoutResponse)
  
    - [Msg](#ibc.core.channel.v1.Msg)
  
- [ibc/core/client/v1/genesis.proto](#ibc/core/client/v1/genesis.proto)
    - [GenesisMetadata](#ibc.core.client.v1.GenesisMetadata)
    - [GenesisState](#ibc.core.client.v1.GenesisState)
    - [IdentifiedGenesisMetadata](#ibc.core.client.v1.IdentifiedGenesisMetadata)
  
- [ibc/core/client/v1/query.proto](#ibc/core/client/v1/query.proto)
    - [QueryClientParamsRequest](#ibc.core.client.v1.QueryClientParamsRequest)
    - [QueryClientParamsResponse](#ibc.core.client.v1.QueryClientParamsResponse)
    - [QueryClientStateRequest](#ibc.core.client.v1.QueryClientStateRequest)
    - [QueryClientStateResponse](#ibc.core.client.v1.QueryClientStateResponse)
    - [QueryClientStatesRequest](#ibc.core.client.v1.QueryClientStatesRequest)
    - [QueryClientStatesResponse](#ibc.core.client.v1.QueryClientStatesResponse)
    - [QueryConsensusStateRequest](#ibc.core.client.v1.QueryConsensusStateRequest)
    - [QueryConsensusStateResponse](#ibc.core.client.v1.QueryConsensusStateResponse)
    - [QueryConsensusStatesRequest](#ibc.core.client.v1.QueryConsensusStatesRequest)
    - [QueryConsensusStatesResponse](#ibc.core.client.v1.QueryConsensusStatesResponse)
  
    - [Query](#ibc.core.client.v1.Query)
  
- [ibc/core/client/v1/tx.proto](#ibc/core/client/v1/tx.proto)
    - [MsgCreateClient](#ibc.core.client.v1.MsgCreateClient)
    - [MsgCreateClientResponse](#ibc.core.client.v1.MsgCreateClientResponse)
    - [MsgSubmitMisbehaviour](#ibc.core.client.v1.MsgSubmitMisbehaviour)
    - [MsgSubmitMisbehaviourResponse](#ibc.core.client.v1.MsgSubmitMisbehaviourResponse)
    - [MsgUpdateClient](#ibc.core.client.v1.MsgUpdateClient)
    - [MsgUpdateClientResponse](#ibc.core.client.v1.MsgUpdateClientResponse)
    - [MsgUpgradeClient](#ibc.core.client.v1.MsgUpgradeClient)
    - [MsgUpgradeClientResponse](#ibc.core.client.v1.MsgUpgradeClientResponse)
  
    - [Msg](#ibc.core.client.v1.Msg)
  
- [ibc/core/commitment/v1/commitment.proto](#ibc/core/commitment/v1/commitment.proto)
    - [MerklePath](#ibc.core.commitment.v1.MerklePath)
    - [MerklePrefix](#ibc.core.commitment.v1.MerklePrefix)
    - [MerkleProof](#ibc.core.commitment.v1.MerkleProof)
    - [MerkleRoot](#ibc.core.commitment.v1.MerkleRoot)
  
- [ibc/core/connection/v1/connection.proto](#ibc/core/connection/v1/connection.proto)
    - [ClientPaths](#ibc.core.connection.v1.ClientPaths)
    - [ConnectionEnd](#ibc.core.connection.v1.ConnectionEnd)
    - [ConnectionPaths](#ibc.core.connection.v1.ConnectionPaths)
    - [Counterparty](#ibc.core.connection.v1.Counterparty)
    - [IdentifiedConnection](#ibc.core.connection.v1.IdentifiedConnection)
    - [Version](#ibc.core.connection.v1.Version)
  
    - [State](#ibc.core.connection.v1.State)
  
- [ibc/core/connection/v1/genesis.proto](#ibc/core/connection/v1/genesis.proto)
    - [GenesisState](#ibc.core.connection.v1.GenesisState)
  
- [ibc/core/connection/v1/query.proto](#ibc/core/connection/v1/query.proto)
    - [QueryClientConnectionsRequest](#ibc.core.connection.v1.QueryClientConnectionsRequest)
    - [QueryClientConnectionsResponse](#ibc.core.connection.v1.QueryClientConnectionsResponse)
    - [QueryConnectionClientStateRequest](#ibc.core.connection.v1.QueryConnectionClientStateRequest)
    - [QueryConnectionClientStateResponse](#ibc.core.connection.v1.QueryConnectionClientStateResponse)
    - [QueryConnectionConsensusStateRequest](#ibc.core.connection.v1.QueryConnectionConsensusStateRequest)
    - [QueryConnectionConsensusStateResponse](#ibc.core.connection.v1.QueryConnectionConsensusStateResponse)
    - [QueryConnectionRequest](#ibc.core.connection.v1.QueryConnectionRequest)
    - [QueryConnectionResponse](#ibc.core.connection.v1.QueryConnectionResponse)
    - [QueryConnectionsRequest](#ibc.core.connection.v1.QueryConnectionsRequest)
    - [QueryConnectionsResponse](#ibc.core.connection.v1.QueryConnectionsResponse)
  
    - [Query](#ibc.core.connection.v1.Query)
  
- [ibc/core/connection/v1/tx.proto](#ibc/core/connection/v1/tx.proto)
    - [MsgConnectionOpenAck](#ibc.core.connection.v1.MsgConnectionOpenAck)
    - [MsgConnectionOpenAckResponse](#ibc.core.connection.v1.MsgConnectionOpenAckResponse)
    - [MsgConnectionOpenConfirm](#ibc.core.connection.v1.MsgConnectionOpenConfirm)
    - [MsgConnectionOpenConfirmResponse](#ibc.core.connection.v1.MsgConnectionOpenConfirmResponse)
    - [MsgConnectionOpenInit](#ibc.core.connection.v1.MsgConnectionOpenInit)
    - [MsgConnectionOpenInitResponse](#ibc.core.connection.v1.MsgConnectionOpenInitResponse)
    - [MsgConnectionOpenTry](#ibc.core.connection.v1.MsgConnectionOpenTry)
    - [MsgConnectionOpenTryResponse](#ibc.core.connection.v1.MsgConnectionOpenTryResponse)
  
    - [Msg](#ibc.core.connection.v1.Msg)
  
- [ibc/core/types/v1/genesis.proto](#ibc/core/types/v1/genesis.proto)
    - [GenesisState](#ibc.core.types.v1.GenesisState)
  
- [ibc/lightclients/localhost/v1/localhost.proto](#ibc/lightclients/localhost/v1/localhost.proto)
    - [ClientState](#ibc.lightclients.localhost.v1.ClientState)
  
- [ibc/lightclients/ostracon/v1/ostracon.proto](#ibc/lightclients/ostracon/v1/ostracon.proto)
    - [ClientState](#ibc.lightclients.ostracon.v1.ClientState)
    - [ConsensusState](#ibc.lightclients.ostracon.v1.ConsensusState)
    - [Fraction](#ibc.lightclients.ostracon.v1.Fraction)
    - [Header](#ibc.lightclients.ostracon.v1.Header)
    - [Misbehaviour](#ibc.lightclients.ostracon.v1.Misbehaviour)
  
- [ibc/lightclients/solomachine/v1/solomachine.proto](#ibc/lightclients/solomachine/v1/solomachine.proto)
    - [ChannelStateData](#ibc.lightclients.solomachine.v1.ChannelStateData)
    - [ClientState](#ibc.lightclients.solomachine.v1.ClientState)
    - [ClientStateData](#ibc.lightclients.solomachine.v1.ClientStateData)
    - [ConnectionStateData](#ibc.lightclients.solomachine.v1.ConnectionStateData)
    - [ConsensusState](#ibc.lightclients.solomachine.v1.ConsensusState)
    - [ConsensusStateData](#ibc.lightclients.solomachine.v1.ConsensusStateData)
    - [Header](#ibc.lightclients.solomachine.v1.Header)
    - [HeaderData](#ibc.lightclients.solomachine.v1.HeaderData)
    - [Misbehaviour](#ibc.lightclients.solomachine.v1.Misbehaviour)
    - [NextSequenceRecvData](#ibc.lightclients.solomachine.v1.NextSequenceRecvData)
    - [PacketAcknowledgementData](#ibc.lightclients.solomachine.v1.PacketAcknowledgementData)
    - [PacketCommitmentData](#ibc.lightclients.solomachine.v1.PacketCommitmentData)
    - [PacketReceiptAbsenceData](#ibc.lightclients.solomachine.v1.PacketReceiptAbsenceData)
    - [SignBytes](#ibc.lightclients.solomachine.v1.SignBytes)
    - [SignatureAndData](#ibc.lightclients.solomachine.v1.SignatureAndData)
    - [TimestampedSignatureData](#ibc.lightclients.solomachine.v1.TimestampedSignatureData)
  
    - [DataType](#ibc.lightclients.solomachine.v1.DataType)
  
- [lbm/crypto/ed25519/keys.proto](#lbm/crypto/ed25519/keys.proto)
    - [PrivKey](#lbm.crypto.ed25519.PrivKey)
    - [PubKey](#lbm.crypto.ed25519.PubKey)
  
- [lbm/crypto/multisig/keys.proto](#lbm/crypto/multisig/keys.proto)
    - [LegacyAminoPubKey](#lbm.crypto.multisig.LegacyAminoPubKey)
  
- [lbm/crypto/secp256k1/keys.proto](#lbm/crypto/secp256k1/keys.proto)
    - [PrivKey](#lbm.crypto.secp256k1.PrivKey)
    - [PubKey](#lbm.crypto.secp256k1.PubKey)
  
- [lbm/auth/v1/auth.proto](#lbm/auth/v1/auth.proto)
    - [BaseAccount](#lbm.auth.v1.BaseAccount)
    - [ModuleAccount](#lbm.auth.v1.ModuleAccount)
    - [Params](#lbm.auth.v1.Params)
  
- [lbm/auth/v1/genesis.proto](#lbm/auth/v1/genesis.proto)
    - [GenesisState](#lbm.auth.v1.GenesisState)
  
- [lbm/auth/v1/query.proto](#lbm/auth/v1/query.proto)
    - [QueryAccountRequest](#lbm.auth.v1.QueryAccountRequest)
    - [QueryAccountResponse](#lbm.auth.v1.QueryAccountResponse)
    - [QueryParamsRequest](#lbm.auth.v1.QueryParamsRequest)
    - [QueryParamsResponse](#lbm.auth.v1.QueryParamsResponse)
  
    - [Query](#lbm.auth.v1.Query)
  
- [lbm/auth/v1/tx.proto](#lbm/auth/v1/tx.proto)
    - [MsgEmpty](#lbm.auth.v1.MsgEmpty)
    - [MsgEmptyResponse](#lbm.auth.v1.MsgEmptyResponse)
  
    - [Msg](#lbm.auth.v1.Msg)
  
- [lbm/bank/v1/bank.proto](#lbm/bank/v1/bank.proto)
    - [DenomUnit](#lbm.bank.v1.DenomUnit)
    - [Input](#lbm.bank.v1.Input)
    - [Metadata](#lbm.bank.v1.Metadata)
    - [Output](#lbm.bank.v1.Output)
    - [Params](#lbm.bank.v1.Params)
    - [SendEnabled](#lbm.bank.v1.SendEnabled)
    - [Supply](#lbm.bank.v1.Supply)
  
- [lbm/bank/v1/genesis.proto](#lbm/bank/v1/genesis.proto)
    - [Balance](#lbm.bank.v1.Balance)
    - [GenesisState](#lbm.bank.v1.GenesisState)
  
- [lbm/bank/v1/query.proto](#lbm/bank/v1/query.proto)
    - [QueryAllBalancesRequest](#lbm.bank.v1.QueryAllBalancesRequest)
    - [QueryAllBalancesResponse](#lbm.bank.v1.QueryAllBalancesResponse)
    - [QueryBalanceRequest](#lbm.bank.v1.QueryBalanceRequest)
    - [QueryBalanceResponse](#lbm.bank.v1.QueryBalanceResponse)
    - [QueryDenomMetadataRequest](#lbm.bank.v1.QueryDenomMetadataRequest)
    - [QueryDenomMetadataResponse](#lbm.bank.v1.QueryDenomMetadataResponse)
    - [QueryDenomsMetadataRequest](#lbm.bank.v1.QueryDenomsMetadataRequest)
    - [QueryDenomsMetadataResponse](#lbm.bank.v1.QueryDenomsMetadataResponse)
    - [QueryParamsRequest](#lbm.bank.v1.QueryParamsRequest)
    - [QueryParamsResponse](#lbm.bank.v1.QueryParamsResponse)
    - [QuerySupplyOfRequest](#lbm.bank.v1.QuerySupplyOfRequest)
    - [QuerySupplyOfResponse](#lbm.bank.v1.QuerySupplyOfResponse)
    - [QueryTotalSupplyRequest](#lbm.bank.v1.QueryTotalSupplyRequest)
    - [QueryTotalSupplyResponse](#lbm.bank.v1.QueryTotalSupplyResponse)
  
    - [Query](#lbm.bank.v1.Query)
  
- [lbm/bank/v1/tx.proto](#lbm/bank/v1/tx.proto)
    - [MsgMultiSend](#lbm.bank.v1.MsgMultiSend)
    - [MsgMultiSendResponse](#lbm.bank.v1.MsgMultiSendResponse)
    - [MsgSend](#lbm.bank.v1.MsgSend)
    - [MsgSendResponse](#lbm.bank.v1.MsgSendResponse)
  
    - [Msg](#lbm.bank.v1.Msg)
  
- [lbm/bankplus/v1/bankplus.proto](#lbm/bankplus/v1/bankplus.proto)
    - [InactiveAddr](#lbm.bankplus.v1.InactiveAddr)
  
- [lbm/base/abci/v1/abci.proto](#lbm/base/abci/v1/abci.proto)
    - [ABCIMessageLog](#lbm.base.abci.v1.ABCIMessageLog)
    - [Attribute](#lbm.base.abci.v1.Attribute)
    - [GasInfo](#lbm.base.abci.v1.GasInfo)
    - [MsgData](#lbm.base.abci.v1.MsgData)
    - [Result](#lbm.base.abci.v1.Result)
    - [SearchTxsResult](#lbm.base.abci.v1.SearchTxsResult)
    - [SimulationResponse](#lbm.base.abci.v1.SimulationResponse)
    - [StringEvent](#lbm.base.abci.v1.StringEvent)
    - [TxMsgData](#lbm.base.abci.v1.TxMsgData)
    - [TxResponse](#lbm.base.abci.v1.TxResponse)
  
- [lbm/base/kv/v1/kv.proto](#lbm/base/kv/v1/kv.proto)
    - [Pair](#lbm.base.kv.v1.Pair)
    - [Pairs](#lbm.base.kv.v1.Pairs)
  
- [lbm/base/ostracon/v1/query.proto](#lbm/base/ostracon/v1/query.proto)
    - [GetBlockByHashRequest](#lbm.base.ostracon.v1.GetBlockByHashRequest)
    - [GetBlockByHashResponse](#lbm.base.ostracon.v1.GetBlockByHashResponse)
    - [GetBlockByHeightRequest](#lbm.base.ostracon.v1.GetBlockByHeightRequest)
    - [GetBlockByHeightResponse](#lbm.base.ostracon.v1.GetBlockByHeightResponse)
    - [GetBlockResultsByHeightRequest](#lbm.base.ostracon.v1.GetBlockResultsByHeightRequest)
    - [GetBlockResultsByHeightResponse](#lbm.base.ostracon.v1.GetBlockResultsByHeightResponse)
    - [GetLatestBlockRequest](#lbm.base.ostracon.v1.GetLatestBlockRequest)
    - [GetLatestBlockResponse](#lbm.base.ostracon.v1.GetLatestBlockResponse)
    - [GetLatestValidatorSetRequest](#lbm.base.ostracon.v1.GetLatestValidatorSetRequest)
    - [GetLatestValidatorSetResponse](#lbm.base.ostracon.v1.GetLatestValidatorSetResponse)
    - [GetNodeInfoRequest](#lbm.base.ostracon.v1.GetNodeInfoRequest)
    - [GetNodeInfoResponse](#lbm.base.ostracon.v1.GetNodeInfoResponse)
    - [GetSyncingRequest](#lbm.base.ostracon.v1.GetSyncingRequest)
    - [GetSyncingResponse](#lbm.base.ostracon.v1.GetSyncingResponse)
    - [GetValidatorSetByHeightRequest](#lbm.base.ostracon.v1.GetValidatorSetByHeightRequest)
    - [GetValidatorSetByHeightResponse](#lbm.base.ostracon.v1.GetValidatorSetByHeightResponse)
    - [Module](#lbm.base.ostracon.v1.Module)
    - [Validator](#lbm.base.ostracon.v1.Validator)
    - [VersionInfo](#lbm.base.ostracon.v1.VersionInfo)
  
    - [Service](#lbm.base.ostracon.v1.Service)
  
- [lbm/base/reflection/v1/reflection.proto](#lbm/base/reflection/v1/reflection.proto)
    - [ListAllInterfacesRequest](#lbm.base.reflection.v1.ListAllInterfacesRequest)
    - [ListAllInterfacesResponse](#lbm.base.reflection.v1.ListAllInterfacesResponse)
    - [ListImplementationsRequest](#lbm.base.reflection.v1.ListImplementationsRequest)
    - [ListImplementationsResponse](#lbm.base.reflection.v1.ListImplementationsResponse)
  
    - [ReflectionService](#lbm.base.reflection.v1.ReflectionService)
  
- [lbm/base/snapshots/v1/snapshot.proto](#lbm/base/snapshots/v1/snapshot.proto)
    - [Metadata](#lbm.base.snapshots.v1.Metadata)
    - [Snapshot](#lbm.base.snapshots.v1.Snapshot)
  
- [lbm/base/store/v1/commit_info.proto](#lbm/base/store/v1/commit_info.proto)
    - [CommitID](#lbm.base.store.v1.CommitID)
    - [CommitInfo](#lbm.base.store.v1.CommitInfo)
    - [StoreInfo](#lbm.base.store.v1.StoreInfo)
  
- [lbm/base/store/v1/snapshot.proto](#lbm/base/store/v1/snapshot.proto)
    - [SnapshotIAVLItem](#lbm.base.store.v1.SnapshotIAVLItem)
    - [SnapshotItem](#lbm.base.store.v1.SnapshotItem)
    - [SnapshotStoreItem](#lbm.base.store.v1.SnapshotStoreItem)
  
- [lbm/capability/v1/capability.proto](#lbm/capability/v1/capability.proto)
    - [Capability](#lbm.capability.v1.Capability)
    - [CapabilityOwners](#lbm.capability.v1.CapabilityOwners)
    - [Owner](#lbm.capability.v1.Owner)
  
- [lbm/capability/v1/genesis.proto](#lbm/capability/v1/genesis.proto)
    - [GenesisOwners](#lbm.capability.v1.GenesisOwners)
    - [GenesisState](#lbm.capability.v1.GenesisState)
  
- [lbm/consortium/v1/consortium.proto](#lbm/consortium/v1/consortium.proto)
    - [Params](#lbm.consortium.v1.Params)
    - [UpdateConsortiumParamsProposal](#lbm.consortium.v1.UpdateConsortiumParamsProposal)
    - [UpdateValidatorAuthsProposal](#lbm.consortium.v1.UpdateValidatorAuthsProposal)
    - [ValidatorAuth](#lbm.consortium.v1.ValidatorAuth)
  
- [lbm/consortium/v1/event.proto](#lbm/consortium/v1/event.proto)
    - [EventUpdateConsortiumParams](#lbm.consortium.v1.EventUpdateConsortiumParams)
    - [EventUpdateValidatorAuths](#lbm.consortium.v1.EventUpdateValidatorAuths)
  
- [lbm/consortium/v1/genesis.proto](#lbm/consortium/v1/genesis.proto)
    - [GenesisState](#lbm.consortium.v1.GenesisState)
  
- [lbm/consortium/v1/query.proto](#lbm/consortium/v1/query.proto)
    - [QueryParamsRequest](#lbm.consortium.v1.QueryParamsRequest)
    - [QueryParamsResponse](#lbm.consortium.v1.QueryParamsResponse)
    - [QueryValidatorAuthRequest](#lbm.consortium.v1.QueryValidatorAuthRequest)
    - [QueryValidatorAuthResponse](#lbm.consortium.v1.QueryValidatorAuthResponse)
    - [QueryValidatorAuthsRequest](#lbm.consortium.v1.QueryValidatorAuthsRequest)
    - [QueryValidatorAuthsResponse](#lbm.consortium.v1.QueryValidatorAuthsResponse)
  
    - [Query](#lbm.consortium.v1.Query)
  
- [lbm/crisis/v1/genesis.proto](#lbm/crisis/v1/genesis.proto)
    - [GenesisState](#lbm.crisis.v1.GenesisState)
  
- [lbm/crisis/v1/tx.proto](#lbm/crisis/v1/tx.proto)
    - [MsgVerifyInvariant](#lbm.crisis.v1.MsgVerifyInvariant)
    - [MsgVerifyInvariantResponse](#lbm.crisis.v1.MsgVerifyInvariantResponse)
  
    - [Msg](#lbm.crisis.v1.Msg)
  
- [lbm/crypto/multisig/v1/multisig.proto](#lbm/crypto/multisig/v1/multisig.proto)
    - [CompactBitArray](#lbm.crypto.multisig.v1.CompactBitArray)
    - [MultiSignature](#lbm.crypto.multisig.v1.MultiSignature)
  
- [lbm/distribution/v1/distribution.proto](#lbm/distribution/v1/distribution.proto)
    - [CommunityPoolSpendProposal](#lbm.distribution.v1.CommunityPoolSpendProposal)
    - [CommunityPoolSpendProposalWithDeposit](#lbm.distribution.v1.CommunityPoolSpendProposalWithDeposit)
    - [DelegationDelegatorReward](#lbm.distribution.v1.DelegationDelegatorReward)
    - [DelegatorStartingInfo](#lbm.distribution.v1.DelegatorStartingInfo)
    - [FeePool](#lbm.distribution.v1.FeePool)
    - [Params](#lbm.distribution.v1.Params)
    - [ValidatorAccumulatedCommission](#lbm.distribution.v1.ValidatorAccumulatedCommission)
    - [ValidatorCurrentRewards](#lbm.distribution.v1.ValidatorCurrentRewards)
    - [ValidatorHistoricalRewards](#lbm.distribution.v1.ValidatorHistoricalRewards)
    - [ValidatorOutstandingRewards](#lbm.distribution.v1.ValidatorOutstandingRewards)
    - [ValidatorSlashEvent](#lbm.distribution.v1.ValidatorSlashEvent)
    - [ValidatorSlashEvents](#lbm.distribution.v1.ValidatorSlashEvents)
  
- [lbm/distribution/v1/genesis.proto](#lbm/distribution/v1/genesis.proto)
    - [DelegatorStartingInfoRecord](#lbm.distribution.v1.DelegatorStartingInfoRecord)
    - [DelegatorWithdrawInfo](#lbm.distribution.v1.DelegatorWithdrawInfo)
    - [GenesisState](#lbm.distribution.v1.GenesisState)
    - [ValidatorAccumulatedCommissionRecord](#lbm.distribution.v1.ValidatorAccumulatedCommissionRecord)
    - [ValidatorCurrentRewardsRecord](#lbm.distribution.v1.ValidatorCurrentRewardsRecord)
    - [ValidatorHistoricalRewardsRecord](#lbm.distribution.v1.ValidatorHistoricalRewardsRecord)
    - [ValidatorOutstandingRewardsRecord](#lbm.distribution.v1.ValidatorOutstandingRewardsRecord)
    - [ValidatorSlashEventRecord](#lbm.distribution.v1.ValidatorSlashEventRecord)
  
- [lbm/distribution/v1/query.proto](#lbm/distribution/v1/query.proto)
    - [QueryCommunityPoolRequest](#lbm.distribution.v1.QueryCommunityPoolRequest)
    - [QueryCommunityPoolResponse](#lbm.distribution.v1.QueryCommunityPoolResponse)
    - [QueryDelegationRewardsRequest](#lbm.distribution.v1.QueryDelegationRewardsRequest)
    - [QueryDelegationRewardsResponse](#lbm.distribution.v1.QueryDelegationRewardsResponse)
    - [QueryDelegationTotalRewardsRequest](#lbm.distribution.v1.QueryDelegationTotalRewardsRequest)
    - [QueryDelegationTotalRewardsResponse](#lbm.distribution.v1.QueryDelegationTotalRewardsResponse)
    - [QueryDelegatorValidatorsRequest](#lbm.distribution.v1.QueryDelegatorValidatorsRequest)
    - [QueryDelegatorValidatorsResponse](#lbm.distribution.v1.QueryDelegatorValidatorsResponse)
    - [QueryDelegatorWithdrawAddressRequest](#lbm.distribution.v1.QueryDelegatorWithdrawAddressRequest)
    - [QueryDelegatorWithdrawAddressResponse](#lbm.distribution.v1.QueryDelegatorWithdrawAddressResponse)
    - [QueryParamsRequest](#lbm.distribution.v1.QueryParamsRequest)
    - [QueryParamsResponse](#lbm.distribution.v1.QueryParamsResponse)
    - [QueryValidatorCommissionRequest](#lbm.distribution.v1.QueryValidatorCommissionRequest)
    - [QueryValidatorCommissionResponse](#lbm.distribution.v1.QueryValidatorCommissionResponse)
    - [QueryValidatorOutstandingRewardsRequest](#lbm.distribution.v1.QueryValidatorOutstandingRewardsRequest)
    - [QueryValidatorOutstandingRewardsResponse](#lbm.distribution.v1.QueryValidatorOutstandingRewardsResponse)
    - [QueryValidatorSlashesRequest](#lbm.distribution.v1.QueryValidatorSlashesRequest)
    - [QueryValidatorSlashesResponse](#lbm.distribution.v1.QueryValidatorSlashesResponse)
  
    - [Query](#lbm.distribution.v1.Query)
  
- [lbm/distribution/v1/tx.proto](#lbm/distribution/v1/tx.proto)
    - [MsgFundCommunityPool](#lbm.distribution.v1.MsgFundCommunityPool)
    - [MsgFundCommunityPoolResponse](#lbm.distribution.v1.MsgFundCommunityPoolResponse)
    - [MsgSetWithdrawAddress](#lbm.distribution.v1.MsgSetWithdrawAddress)
    - [MsgSetWithdrawAddressResponse](#lbm.distribution.v1.MsgSetWithdrawAddressResponse)
    - [MsgWithdrawDelegatorReward](#lbm.distribution.v1.MsgWithdrawDelegatorReward)
    - [MsgWithdrawDelegatorRewardResponse](#lbm.distribution.v1.MsgWithdrawDelegatorRewardResponse)
    - [MsgWithdrawValidatorCommission](#lbm.distribution.v1.MsgWithdrawValidatorCommission)
    - [MsgWithdrawValidatorCommissionResponse](#lbm.distribution.v1.MsgWithdrawValidatorCommissionResponse)
  
    - [Msg](#lbm.distribution.v1.Msg)
  
- [lbm/evidence/v1/evidence.proto](#lbm/evidence/v1/evidence.proto)
    - [Equivocation](#lbm.evidence.v1.Equivocation)
  
- [lbm/evidence/v1/genesis.proto](#lbm/evidence/v1/genesis.proto)
    - [GenesisState](#lbm.evidence.v1.GenesisState)
  
- [lbm/evidence/v1/query.proto](#lbm/evidence/v1/query.proto)
    - [QueryAllEvidenceRequest](#lbm.evidence.v1.QueryAllEvidenceRequest)
    - [QueryAllEvidenceResponse](#lbm.evidence.v1.QueryAllEvidenceResponse)
    - [QueryEvidenceRequest](#lbm.evidence.v1.QueryEvidenceRequest)
    - [QueryEvidenceResponse](#lbm.evidence.v1.QueryEvidenceResponse)
  
    - [Query](#lbm.evidence.v1.Query)
  
- [lbm/evidence/v1/tx.proto](#lbm/evidence/v1/tx.proto)
    - [MsgSubmitEvidence](#lbm.evidence.v1.MsgSubmitEvidence)
    - [MsgSubmitEvidenceResponse](#lbm.evidence.v1.MsgSubmitEvidenceResponse)
  
    - [Msg](#lbm.evidence.v1.Msg)
  
- [lbm/feegrant/v1/feegrant.proto](#lbm/feegrant/v1/feegrant.proto)
    - [AllowedMsgAllowance](#lbm.feegrant.v1.AllowedMsgAllowance)
    - [BasicAllowance](#lbm.feegrant.v1.BasicAllowance)
    - [Grant](#lbm.feegrant.v1.Grant)
    - [PeriodicAllowance](#lbm.feegrant.v1.PeriodicAllowance)
  
- [lbm/feegrant/v1/genesis.proto](#lbm/feegrant/v1/genesis.proto)
    - [GenesisState](#lbm.feegrant.v1.GenesisState)
  
- [lbm/feegrant/v1/query.proto](#lbm/feegrant/v1/query.proto)
    - [QueryAllowanceRequest](#lbm.feegrant.v1.QueryAllowanceRequest)
    - [QueryAllowanceResponse](#lbm.feegrant.v1.QueryAllowanceResponse)
    - [QueryAllowancesRequest](#lbm.feegrant.v1.QueryAllowancesRequest)
    - [QueryAllowancesResponse](#lbm.feegrant.v1.QueryAllowancesResponse)
  
    - [Query](#lbm.feegrant.v1.Query)
  
- [lbm/feegrant/v1/tx.proto](#lbm/feegrant/v1/tx.proto)
    - [MsgGrantAllowance](#lbm.feegrant.v1.MsgGrantAllowance)
    - [MsgGrantAllowanceResponse](#lbm.feegrant.v1.MsgGrantAllowanceResponse)
    - [MsgRevokeAllowance](#lbm.feegrant.v1.MsgRevokeAllowance)
    - [MsgRevokeAllowanceResponse](#lbm.feegrant.v1.MsgRevokeAllowanceResponse)
  
    - [Msg](#lbm.feegrant.v1.Msg)
  
- [lbm/genutil/v1/genesis.proto](#lbm/genutil/v1/genesis.proto)
    - [GenesisState](#lbm.genutil.v1.GenesisState)
  
- [lbm/gov/v1/gov.proto](#lbm/gov/v1/gov.proto)
    - [Deposit](#lbm.gov.v1.Deposit)
    - [DepositParams](#lbm.gov.v1.DepositParams)
    - [Proposal](#lbm.gov.v1.Proposal)
    - [TallyParams](#lbm.gov.v1.TallyParams)
    - [TallyResult](#lbm.gov.v1.TallyResult)
    - [TextProposal](#lbm.gov.v1.TextProposal)
    - [Vote](#lbm.gov.v1.Vote)
    - [VotingParams](#lbm.gov.v1.VotingParams)
    - [WeightedVoteOption](#lbm.gov.v1.WeightedVoteOption)
  
    - [ProposalStatus](#lbm.gov.v1.ProposalStatus)
    - [VoteOption](#lbm.gov.v1.VoteOption)
  
- [lbm/gov/v1/genesis.proto](#lbm/gov/v1/genesis.proto)
    - [GenesisState](#lbm.gov.v1.GenesisState)
  
- [lbm/gov/v1/query.proto](#lbm/gov/v1/query.proto)
    - [QueryDepositRequest](#lbm.gov.v1.QueryDepositRequest)
    - [QueryDepositResponse](#lbm.gov.v1.QueryDepositResponse)
    - [QueryDepositsRequest](#lbm.gov.v1.QueryDepositsRequest)
    - [QueryDepositsResponse](#lbm.gov.v1.QueryDepositsResponse)
    - [QueryParamsRequest](#lbm.gov.v1.QueryParamsRequest)
    - [QueryParamsResponse](#lbm.gov.v1.QueryParamsResponse)
    - [QueryProposalRequest](#lbm.gov.v1.QueryProposalRequest)
    - [QueryProposalResponse](#lbm.gov.v1.QueryProposalResponse)
    - [QueryProposalsRequest](#lbm.gov.v1.QueryProposalsRequest)
    - [QueryProposalsResponse](#lbm.gov.v1.QueryProposalsResponse)
    - [QueryTallyResultRequest](#lbm.gov.v1.QueryTallyResultRequest)
    - [QueryTallyResultResponse](#lbm.gov.v1.QueryTallyResultResponse)
    - [QueryVoteRequest](#lbm.gov.v1.QueryVoteRequest)
    - [QueryVoteResponse](#lbm.gov.v1.QueryVoteResponse)
    - [QueryVotesRequest](#lbm.gov.v1.QueryVotesRequest)
    - [QueryVotesResponse](#lbm.gov.v1.QueryVotesResponse)
  
    - [Query](#lbm.gov.v1.Query)
  
- [lbm/gov/v1/tx.proto](#lbm/gov/v1/tx.proto)
    - [MsgDeposit](#lbm.gov.v1.MsgDeposit)
    - [MsgDepositResponse](#lbm.gov.v1.MsgDepositResponse)
    - [MsgSubmitProposal](#lbm.gov.v1.MsgSubmitProposal)
    - [MsgSubmitProposalResponse](#lbm.gov.v1.MsgSubmitProposalResponse)
    - [MsgVote](#lbm.gov.v1.MsgVote)
    - [MsgVoteResponse](#lbm.gov.v1.MsgVoteResponse)
    - [MsgVoteWeighted](#lbm.gov.v1.MsgVoteWeighted)
    - [MsgVoteWeightedResponse](#lbm.gov.v1.MsgVoteWeightedResponse)
  
    - [Msg](#lbm.gov.v1.Msg)
  
- [lbm/mint/v1/mint.proto](#lbm/mint/v1/mint.proto)
    - [Minter](#lbm.mint.v1.Minter)
    - [Params](#lbm.mint.v1.Params)
  
- [lbm/mint/v1/genesis.proto](#lbm/mint/v1/genesis.proto)
    - [GenesisState](#lbm.mint.v1.GenesisState)
  
- [lbm/mint/v1/query.proto](#lbm/mint/v1/query.proto)
    - [QueryAnnualProvisionsRequest](#lbm.mint.v1.QueryAnnualProvisionsRequest)
    - [QueryAnnualProvisionsResponse](#lbm.mint.v1.QueryAnnualProvisionsResponse)
    - [QueryInflationRequest](#lbm.mint.v1.QueryInflationRequest)
    - [QueryInflationResponse](#lbm.mint.v1.QueryInflationResponse)
    - [QueryParamsRequest](#lbm.mint.v1.QueryParamsRequest)
    - [QueryParamsResponse](#lbm.mint.v1.QueryParamsResponse)
  
    - [Query](#lbm.mint.v1.Query)
  
- [lbm/params/v1/params.proto](#lbm/params/v1/params.proto)
    - [ParamChange](#lbm.params.v1.ParamChange)
    - [ParameterChangeProposal](#lbm.params.v1.ParameterChangeProposal)
  
- [lbm/params/v1/query.proto](#lbm/params/v1/query.proto)
    - [QueryParamsRequest](#lbm.params.v1.QueryParamsRequest)
    - [QueryParamsResponse](#lbm.params.v1.QueryParamsResponse)
  
    - [Query](#lbm.params.v1.Query)
  
- [lbm/slashing/v1/slashing.proto](#lbm/slashing/v1/slashing.proto)
    - [Params](#lbm.slashing.v1.Params)
    - [ValidatorSigningInfo](#lbm.slashing.v1.ValidatorSigningInfo)
  
- [lbm/slashing/v1/genesis.proto](#lbm/slashing/v1/genesis.proto)
    - [GenesisState](#lbm.slashing.v1.GenesisState)
    - [MissedBlock](#lbm.slashing.v1.MissedBlock)
    - [SigningInfo](#lbm.slashing.v1.SigningInfo)
    - [ValidatorMissedBlocks](#lbm.slashing.v1.ValidatorMissedBlocks)
  
- [lbm/slashing/v1/query.proto](#lbm/slashing/v1/query.proto)
    - [QueryParamsRequest](#lbm.slashing.v1.QueryParamsRequest)
    - [QueryParamsResponse](#lbm.slashing.v1.QueryParamsResponse)
    - [QuerySigningInfoRequest](#lbm.slashing.v1.QuerySigningInfoRequest)
    - [QuerySigningInfoResponse](#lbm.slashing.v1.QuerySigningInfoResponse)
    - [QuerySigningInfosRequest](#lbm.slashing.v1.QuerySigningInfosRequest)
    - [QuerySigningInfosResponse](#lbm.slashing.v1.QuerySigningInfosResponse)
  
    - [Query](#lbm.slashing.v1.Query)
  
- [lbm/slashing/v1/tx.proto](#lbm/slashing/v1/tx.proto)
    - [MsgUnjail](#lbm.slashing.v1.MsgUnjail)
    - [MsgUnjailResponse](#lbm.slashing.v1.MsgUnjailResponse)
  
    - [Msg](#lbm.slashing.v1.Msg)
  
- [lbm/staking/v1/staking.proto](#lbm/staking/v1/staking.proto)
    - [Commission](#lbm.staking.v1.Commission)
    - [CommissionRates](#lbm.staking.v1.CommissionRates)
    - [DVPair](#lbm.staking.v1.DVPair)
    - [DVPairs](#lbm.staking.v1.DVPairs)
    - [DVVTriplet](#lbm.staking.v1.DVVTriplet)
    - [DVVTriplets](#lbm.staking.v1.DVVTriplets)
    - [Delegation](#lbm.staking.v1.Delegation)
    - [DelegationResponse](#lbm.staking.v1.DelegationResponse)
    - [Description](#lbm.staking.v1.Description)
    - [HistoricalInfo](#lbm.staking.v1.HistoricalInfo)
    - [Params](#lbm.staking.v1.Params)
    - [Pool](#lbm.staking.v1.Pool)
    - [Redelegation](#lbm.staking.v1.Redelegation)
    - [RedelegationEntry](#lbm.staking.v1.RedelegationEntry)
    - [RedelegationEntryResponse](#lbm.staking.v1.RedelegationEntryResponse)
    - [RedelegationResponse](#lbm.staking.v1.RedelegationResponse)
    - [UnbondingDelegation](#lbm.staking.v1.UnbondingDelegation)
    - [UnbondingDelegationEntry](#lbm.staking.v1.UnbondingDelegationEntry)
    - [ValAddresses](#lbm.staking.v1.ValAddresses)
    - [Validator](#lbm.staking.v1.Validator)
  
    - [BondStatus](#lbm.staking.v1.BondStatus)
  
- [lbm/staking/v1/genesis.proto](#lbm/staking/v1/genesis.proto)
    - [GenesisState](#lbm.staking.v1.GenesisState)
    - [LastValidatorPower](#lbm.staking.v1.LastValidatorPower)
  
- [lbm/staking/v1/query.proto](#lbm/staking/v1/query.proto)
    - [QueryDelegationRequest](#lbm.staking.v1.QueryDelegationRequest)
    - [QueryDelegationResponse](#lbm.staking.v1.QueryDelegationResponse)
    - [QueryDelegatorDelegationsRequest](#lbm.staking.v1.QueryDelegatorDelegationsRequest)
    - [QueryDelegatorDelegationsResponse](#lbm.staking.v1.QueryDelegatorDelegationsResponse)
    - [QueryDelegatorUnbondingDelegationsRequest](#lbm.staking.v1.QueryDelegatorUnbondingDelegationsRequest)
    - [QueryDelegatorUnbondingDelegationsResponse](#lbm.staking.v1.QueryDelegatorUnbondingDelegationsResponse)
    - [QueryDelegatorValidatorRequest](#lbm.staking.v1.QueryDelegatorValidatorRequest)
    - [QueryDelegatorValidatorResponse](#lbm.staking.v1.QueryDelegatorValidatorResponse)
    - [QueryDelegatorValidatorsRequest](#lbm.staking.v1.QueryDelegatorValidatorsRequest)
    - [QueryDelegatorValidatorsResponse](#lbm.staking.v1.QueryDelegatorValidatorsResponse)
    - [QueryHistoricalInfoRequest](#lbm.staking.v1.QueryHistoricalInfoRequest)
    - [QueryHistoricalInfoResponse](#lbm.staking.v1.QueryHistoricalInfoResponse)
    - [QueryParamsRequest](#lbm.staking.v1.QueryParamsRequest)
    - [QueryParamsResponse](#lbm.staking.v1.QueryParamsResponse)
    - [QueryPoolRequest](#lbm.staking.v1.QueryPoolRequest)
    - [QueryPoolResponse](#lbm.staking.v1.QueryPoolResponse)
    - [QueryRedelegationsRequest](#lbm.staking.v1.QueryRedelegationsRequest)
    - [QueryRedelegationsResponse](#lbm.staking.v1.QueryRedelegationsResponse)
    - [QueryUnbondingDelegationRequest](#lbm.staking.v1.QueryUnbondingDelegationRequest)
    - [QueryUnbondingDelegationResponse](#lbm.staking.v1.QueryUnbondingDelegationResponse)
    - [QueryValidatorDelegationsRequest](#lbm.staking.v1.QueryValidatorDelegationsRequest)
    - [QueryValidatorDelegationsResponse](#lbm.staking.v1.QueryValidatorDelegationsResponse)
    - [QueryValidatorRequest](#lbm.staking.v1.QueryValidatorRequest)
    - [QueryValidatorResponse](#lbm.staking.v1.QueryValidatorResponse)
    - [QueryValidatorUnbondingDelegationsRequest](#lbm.staking.v1.QueryValidatorUnbondingDelegationsRequest)
    - [QueryValidatorUnbondingDelegationsResponse](#lbm.staking.v1.QueryValidatorUnbondingDelegationsResponse)
    - [QueryValidatorsRequest](#lbm.staking.v1.QueryValidatorsRequest)
    - [QueryValidatorsResponse](#lbm.staking.v1.QueryValidatorsResponse)
  
    - [Query](#lbm.staking.v1.Query)
  
- [lbm/staking/v1/tx.proto](#lbm/staking/v1/tx.proto)
    - [MsgBeginRedelegate](#lbm.staking.v1.MsgBeginRedelegate)
    - [MsgBeginRedelegateResponse](#lbm.staking.v1.MsgBeginRedelegateResponse)
    - [MsgCreateValidator](#lbm.staking.v1.MsgCreateValidator)
    - [MsgCreateValidatorResponse](#lbm.staking.v1.MsgCreateValidatorResponse)
    - [MsgDelegate](#lbm.staking.v1.MsgDelegate)
    - [MsgDelegateResponse](#lbm.staking.v1.MsgDelegateResponse)
    - [MsgEditValidator](#lbm.staking.v1.MsgEditValidator)
    - [MsgEditValidatorResponse](#lbm.staking.v1.MsgEditValidatorResponse)
    - [MsgUndelegate](#lbm.staking.v1.MsgUndelegate)
    - [MsgUndelegateResponse](#lbm.staking.v1.MsgUndelegateResponse)
  
    - [Msg](#lbm.staking.v1.Msg)
  
- [lbm/token/v1/event.proto](#lbm/token/v1/event.proto)
    - [EventApprove](#lbm.token.v1.EventApprove)
    - [EventBurn](#lbm.token.v1.EventBurn)
    - [EventGrant](#lbm.token.v1.EventGrant)
    - [EventIssue](#lbm.token.v1.EventIssue)
    - [EventMint](#lbm.token.v1.EventMint)
    - [EventModify](#lbm.token.v1.EventModify)
    - [EventRevoke](#lbm.token.v1.EventRevoke)
    - [EventTransfer](#lbm.token.v1.EventTransfer)
  
- [lbm/token/v1/token.proto](#lbm/token/v1/token.proto)
    - [Approve](#lbm.token.v1.Approve)
    - [FT](#lbm.token.v1.FT)
    - [Grant](#lbm.token.v1.Grant)
    - [Pair](#lbm.token.v1.Pair)
    - [Params](#lbm.token.v1.Params)
    - [Token](#lbm.token.v1.Token)
  
- [lbm/token/v1/genesis.proto](#lbm/token/v1/genesis.proto)
    - [Balance](#lbm.token.v1.Balance)
    - [ClassGenesisState](#lbm.token.v1.ClassGenesisState)
    - [GenesisState](#lbm.token.v1.GenesisState)
  
- [lbm/token/v1/query.proto](#lbm/token/v1/query.proto)
    - [QueryApproveRequest](#lbm.token.v1.QueryApproveRequest)
    - [QueryApproveResponse](#lbm.token.v1.QueryApproveResponse)
    - [QueryApprovesRequest](#lbm.token.v1.QueryApprovesRequest)
    - [QueryApprovesResponse](#lbm.token.v1.QueryApprovesResponse)
    - [QueryGrantsRequest](#lbm.token.v1.QueryGrantsRequest)
    - [QueryGrantsResponse](#lbm.token.v1.QueryGrantsResponse)
    - [QuerySupplyRequest](#lbm.token.v1.QuerySupplyRequest)
    - [QuerySupplyResponse](#lbm.token.v1.QuerySupplyResponse)
    - [QueryTokenBalanceRequest](#lbm.token.v1.QueryTokenBalanceRequest)
    - [QueryTokenBalanceResponse](#lbm.token.v1.QueryTokenBalanceResponse)
    - [QueryTokenRequest](#lbm.token.v1.QueryTokenRequest)
    - [QueryTokenResponse](#lbm.token.v1.QueryTokenResponse)
    - [QueryTokensRequest](#lbm.token.v1.QueryTokensRequest)
    - [QueryTokensResponse](#lbm.token.v1.QueryTokensResponse)
  
    - [Query](#lbm.token.v1.Query)
  
- [lbm/token/v1/tx.proto](#lbm/token/v1/tx.proto)
    - [MsgApprove](#lbm.token.v1.MsgApprove)
    - [MsgApproveResponse](#lbm.token.v1.MsgApproveResponse)
    - [MsgBurn](#lbm.token.v1.MsgBurn)
    - [MsgBurnFrom](#lbm.token.v1.MsgBurnFrom)
    - [MsgBurnFromResponse](#lbm.token.v1.MsgBurnFromResponse)
    - [MsgBurnResponse](#lbm.token.v1.MsgBurnResponse)
    - [MsgGrant](#lbm.token.v1.MsgGrant)
    - [MsgGrantResponse](#lbm.token.v1.MsgGrantResponse)
    - [MsgIssue](#lbm.token.v1.MsgIssue)
    - [MsgIssueResponse](#lbm.token.v1.MsgIssueResponse)
    - [MsgMint](#lbm.token.v1.MsgMint)
    - [MsgMintResponse](#lbm.token.v1.MsgMintResponse)
    - [MsgModify](#lbm.token.v1.MsgModify)
    - [MsgModifyResponse](#lbm.token.v1.MsgModifyResponse)
    - [MsgRevoke](#lbm.token.v1.MsgRevoke)
    - [MsgRevokeResponse](#lbm.token.v1.MsgRevokeResponse)
    - [MsgTransfer](#lbm.token.v1.MsgTransfer)
    - [MsgTransferFrom](#lbm.token.v1.MsgTransferFrom)
    - [MsgTransferFromResponse](#lbm.token.v1.MsgTransferFromResponse)
    - [MsgTransferResponse](#lbm.token.v1.MsgTransferResponse)
  
    - [Msg](#lbm.token.v1.Msg)
  
- [lbm/tx/signing/v1/signing.proto](#lbm/tx/signing/v1/signing.proto)
    - [SignatureDescriptor](#lbm.tx.signing.v1.SignatureDescriptor)
    - [SignatureDescriptor.Data](#lbm.tx.signing.v1.SignatureDescriptor.Data)
    - [SignatureDescriptor.Data.Multi](#lbm.tx.signing.v1.SignatureDescriptor.Data.Multi)
    - [SignatureDescriptor.Data.Single](#lbm.tx.signing.v1.SignatureDescriptor.Data.Single)
    - [SignatureDescriptors](#lbm.tx.signing.v1.SignatureDescriptors)
  
    - [SignMode](#lbm.tx.signing.v1.SignMode)
  
- [lbm/tx/v1/tx.proto](#lbm/tx/v1/tx.proto)
    - [AuthInfo](#lbm.tx.v1.AuthInfo)
    - [Fee](#lbm.tx.v1.Fee)
    - [ModeInfo](#lbm.tx.v1.ModeInfo)
    - [ModeInfo.Multi](#lbm.tx.v1.ModeInfo.Multi)
    - [ModeInfo.Single](#lbm.tx.v1.ModeInfo.Single)
    - [SignDoc](#lbm.tx.v1.SignDoc)
    - [SignerInfo](#lbm.tx.v1.SignerInfo)
    - [Tx](#lbm.tx.v1.Tx)
    - [TxBody](#lbm.tx.v1.TxBody)
    - [TxRaw](#lbm.tx.v1.TxRaw)
  
- [lbm/tx/v1/service.proto](#lbm/tx/v1/service.proto)
    - [BroadcastTxRequest](#lbm.tx.v1.BroadcastTxRequest)
    - [BroadcastTxResponse](#lbm.tx.v1.BroadcastTxResponse)
    - [GetTxRequest](#lbm.tx.v1.GetTxRequest)
    - [GetTxResponse](#lbm.tx.v1.GetTxResponse)
    - [GetTxsEventRequest](#lbm.tx.v1.GetTxsEventRequest)
    - [GetTxsEventResponse](#lbm.tx.v1.GetTxsEventResponse)
    - [SimulateRequest](#lbm.tx.v1.SimulateRequest)
    - [SimulateResponse](#lbm.tx.v1.SimulateResponse)
  
    - [BroadcastMode](#lbm.tx.v1.BroadcastMode)
    - [OrderBy](#lbm.tx.v1.OrderBy)
  
    - [Service](#lbm.tx.v1.Service)
  
- [lbm/upgrade/v1/upgrade.proto](#lbm/upgrade/v1/upgrade.proto)
    - [CancelSoftwareUpgradeProposal](#lbm.upgrade.v1.CancelSoftwareUpgradeProposal)
    - [ModuleVersion](#lbm.upgrade.v1.ModuleVersion)
    - [Plan](#lbm.upgrade.v1.Plan)
    - [SoftwareUpgradeProposal](#lbm.upgrade.v1.SoftwareUpgradeProposal)
  
- [lbm/upgrade/v1/query.proto](#lbm/upgrade/v1/query.proto)
    - [QueryAppliedPlanRequest](#lbm.upgrade.v1.QueryAppliedPlanRequest)
    - [QueryAppliedPlanResponse](#lbm.upgrade.v1.QueryAppliedPlanResponse)
    - [QueryCurrentPlanRequest](#lbm.upgrade.v1.QueryCurrentPlanRequest)
    - [QueryCurrentPlanResponse](#lbm.upgrade.v1.QueryCurrentPlanResponse)
    - [QueryModuleVersionsRequest](#lbm.upgrade.v1.QueryModuleVersionsRequest)
    - [QueryModuleVersionsResponse](#lbm.upgrade.v1.QueryModuleVersionsResponse)
    - [QueryUpgradedConsensusStateRequest](#lbm.upgrade.v1.QueryUpgradedConsensusStateRequest)
    - [QueryUpgradedConsensusStateResponse](#lbm.upgrade.v1.QueryUpgradedConsensusStateResponse)
  
    - [Query](#lbm.upgrade.v1.Query)
  
- [lbm/vesting/v1/tx.proto](#lbm/vesting/v1/tx.proto)
    - [MsgCreateVestingAccount](#lbm.vesting.v1.MsgCreateVestingAccount)
    - [MsgCreateVestingAccountResponse](#lbm.vesting.v1.MsgCreateVestingAccountResponse)
  
    - [Msg](#lbm.vesting.v1.Msg)
  
- [lbm/vesting/v1/vesting.proto](#lbm/vesting/v1/vesting.proto)
    - [BaseVestingAccount](#lbm.vesting.v1.BaseVestingAccount)
    - [ContinuousVestingAccount](#lbm.vesting.v1.ContinuousVestingAccount)
    - [DelayedVestingAccount](#lbm.vesting.v1.DelayedVestingAccount)
    - [Period](#lbm.vesting.v1.Period)
    - [PeriodicVestingAccount](#lbm.vesting.v1.PeriodicVestingAccount)
  
- [lbm/wasm/v1/types.proto](#lbm/wasm/v1/types.proto)
    - [AbsoluteTxPosition](#lbm.wasm.v1.AbsoluteTxPosition)
    - [AccessConfig](#lbm.wasm.v1.AccessConfig)
    - [AccessTypeParam](#lbm.wasm.v1.AccessTypeParam)
    - [CodeInfo](#lbm.wasm.v1.CodeInfo)
    - [ContractCodeHistoryEntry](#lbm.wasm.v1.ContractCodeHistoryEntry)
    - [ContractInfo](#lbm.wasm.v1.ContractInfo)
    - [Model](#lbm.wasm.v1.Model)
    - [Params](#lbm.wasm.v1.Params)
  
    - [AccessType](#lbm.wasm.v1.AccessType)
    - [ContractCodeHistoryOperationType](#lbm.wasm.v1.ContractCodeHistoryOperationType)
    - [ContractStatus](#lbm.wasm.v1.ContractStatus)
  
- [lbm/wasm/v1/tx.proto](#lbm/wasm/v1/tx.proto)
    - [MsgClearAdmin](#lbm.wasm.v1.MsgClearAdmin)
    - [MsgClearAdminResponse](#lbm.wasm.v1.MsgClearAdminResponse)
    - [MsgExecuteContract](#lbm.wasm.v1.MsgExecuteContract)
    - [MsgExecuteContractResponse](#lbm.wasm.v1.MsgExecuteContractResponse)
    - [MsgInstantiateContract](#lbm.wasm.v1.MsgInstantiateContract)
    - [MsgInstantiateContractResponse](#lbm.wasm.v1.MsgInstantiateContractResponse)
    - [MsgMigrateContract](#lbm.wasm.v1.MsgMigrateContract)
    - [MsgMigrateContractResponse](#lbm.wasm.v1.MsgMigrateContractResponse)
    - [MsgStoreCode](#lbm.wasm.v1.MsgStoreCode)
    - [MsgStoreCodeAndInstantiateContract](#lbm.wasm.v1.MsgStoreCodeAndInstantiateContract)
    - [MsgStoreCodeAndInstantiateContractResponse](#lbm.wasm.v1.MsgStoreCodeAndInstantiateContractResponse)
    - [MsgStoreCodeResponse](#lbm.wasm.v1.MsgStoreCodeResponse)
    - [MsgUpdateAdmin](#lbm.wasm.v1.MsgUpdateAdmin)
    - [MsgUpdateAdminResponse](#lbm.wasm.v1.MsgUpdateAdminResponse)
    - [MsgUpdateContractStatus](#lbm.wasm.v1.MsgUpdateContractStatus)
    - [MsgUpdateContractStatusResponse](#lbm.wasm.v1.MsgUpdateContractStatusResponse)
  
    - [Msg](#lbm.wasm.v1.Msg)
  
- [lbm/wasm/v1/genesis.proto](#lbm/wasm/v1/genesis.proto)
    - [Code](#lbm.wasm.v1.Code)
    - [Contract](#lbm.wasm.v1.Contract)
    - [GenesisState](#lbm.wasm.v1.GenesisState)
    - [GenesisState.GenMsgs](#lbm.wasm.v1.GenesisState.GenMsgs)
    - [Sequence](#lbm.wasm.v1.Sequence)
  
- [lbm/wasm/v1/ibc.proto](#lbm/wasm/v1/ibc.proto)
    - [MsgIBCCloseChannel](#lbm.wasm.v1.MsgIBCCloseChannel)
    - [MsgIBCSend](#lbm.wasm.v1.MsgIBCSend)
  
- [lbm/wasm/v1/proposal.proto](#lbm/wasm/v1/proposal.proto)
    - [ClearAdminProposal](#lbm.wasm.v1.ClearAdminProposal)
    - [InstantiateContractProposal](#lbm.wasm.v1.InstantiateContractProposal)
    - [MigrateContractProposal](#lbm.wasm.v1.MigrateContractProposal)
    - [PinCodesProposal](#lbm.wasm.v1.PinCodesProposal)
    - [StoreCodeProposal](#lbm.wasm.v1.StoreCodeProposal)
    - [UnpinCodesProposal](#lbm.wasm.v1.UnpinCodesProposal)
    - [UpdateAdminProposal](#lbm.wasm.v1.UpdateAdminProposal)
    - [UpdateContractStatusProposal](#lbm.wasm.v1.UpdateContractStatusProposal)
  
- [lbm/wasm/v1/query.proto](#lbm/wasm/v1/query.proto)
    - [CodeInfoResponse](#lbm.wasm.v1.CodeInfoResponse)
    - [QueryAllContractStateRequest](#lbm.wasm.v1.QueryAllContractStateRequest)
    - [QueryAllContractStateResponse](#lbm.wasm.v1.QueryAllContractStateResponse)
    - [QueryCodeRequest](#lbm.wasm.v1.QueryCodeRequest)
    - [QueryCodeResponse](#lbm.wasm.v1.QueryCodeResponse)
    - [QueryCodesRequest](#lbm.wasm.v1.QueryCodesRequest)
    - [QueryCodesResponse](#lbm.wasm.v1.QueryCodesResponse)
    - [QueryContractHistoryRequest](#lbm.wasm.v1.QueryContractHistoryRequest)
    - [QueryContractHistoryResponse](#lbm.wasm.v1.QueryContractHistoryResponse)
    - [QueryContractInfoRequest](#lbm.wasm.v1.QueryContractInfoRequest)
    - [QueryContractInfoResponse](#lbm.wasm.v1.QueryContractInfoResponse)
    - [QueryContractsByCodeRequest](#lbm.wasm.v1.QueryContractsByCodeRequest)
    - [QueryContractsByCodeResponse](#lbm.wasm.v1.QueryContractsByCodeResponse)
    - [QueryRawContractStateRequest](#lbm.wasm.v1.QueryRawContractStateRequest)
    - [QueryRawContractStateResponse](#lbm.wasm.v1.QueryRawContractStateResponse)
    - [QuerySmartContractStateRequest](#lbm.wasm.v1.QuerySmartContractStateRequest)
    - [QuerySmartContractStateResponse](#lbm.wasm.v1.QuerySmartContractStateResponse)
  
    - [Query](#lbm.wasm.v1.Query)
  
- [Scalar Value Types](#scalar-value-types)



<a name="ibc/applications/transfer/v1/transfer.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/transfer/v1/transfer.proto



<a name="ibc.applications.transfer.v1.DenomTrace"></a>

### DenomTrace
DenomTrace contains the base denomination for ICS20 fungible tokens and the
source tracing information path.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [string](#string) |  | path defines the chain of port/channel identifiers used for tracing the source of the fungible token. |
| `base_denom` | [string](#string) |  | base denomination of the relayed fungible token. |






<a name="ibc.applications.transfer.v1.FungibleTokenPacketData"></a>

### FungibleTokenPacketData
FungibleTokenPacketData defines a struct for the packet payload
See FungibleTokenPacketData spec:
https://github.com/cosmos/ics/tree/master/spec/ics-020-fungible-token-transfer#data-structures


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  | the token denomination to be transferred |
| `amount` | [uint64](#uint64) |  | the token amount to be transferred |
| `sender` | [string](#string) |  | the sender address |
| `receiver` | [string](#string) |  | the recipient address on the destination chain |






<a name="ibc.applications.transfer.v1.Params"></a>

### Params
Params defines the set of IBC transfer parameters.
NOTE: To prevent a single token from being transferred, set the
TransfersEnabled parameter to true and then set the bank module's SendEnabled
parameter for the denomination to false.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `send_enabled` | [bool](#bool) |  | send_enabled enables or disables all cross-chain token transfers from this chain. |
| `receive_enabled` | [bool](#bool) |  | receive_enabled enables or disables all cross-chain token transfers to this chain. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/applications/transfer/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/transfer/v1/genesis.proto



<a name="ibc.applications.transfer.v1.GenesisState"></a>

### GenesisState
GenesisState defines the ibc-transfer genesis state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  |  |
| `denom_traces` | [DenomTrace](#ibc.applications.transfer.v1.DenomTrace) | repeated |  |
| `params` | [Params](#ibc.applications.transfer.v1.Params) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/base/query/v1/pagination.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/base/query/v1/pagination.proto



<a name="lbm.base.query.v1.PageRequest"></a>

### PageRequest
PageRequest is to be embedded in gRPC request messages for efficient
pagination. Ex:

 message SomeRequest {
         Foo some_parameter = 1;
         PageRequest pagination = 2;
 }


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  | key is a value returned in PageResponse.next_key to begin querying the next page most efficiently. Only one of offset or key should be set. |
| `offset` | [uint64](#uint64) |  | offset is a numeric offset that can be used when key is unavailable. It is less efficient than using key. Only one of offset or key should be set. |
| `limit` | [uint64](#uint64) |  | limit is the total number of results to be returned in the result page. If left empty it will default to a value to be set by each app. |
| `count_total` | [bool](#bool) |  | count_total is set to true to indicate that the result set should include a count of the total number of items available for pagination in UIs. count_total is only respected when offset is used. It is ignored when key is set. |






<a name="lbm.base.query.v1.PageResponse"></a>

### PageResponse
PageResponse is to be embedded in gRPC response messages where the
corresponding request message has used PageRequest.

 message SomeResponse {
         repeated Bar results = 1;
         PageResponse page = 2;
 }


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `next_key` | [bytes](#bytes) |  | next_key is the key to be passed to PageRequest.key to query the next page most efficiently |
| `total` | [uint64](#uint64) |  | total is total number of results available if PageRequest.count_total was set, its value is undefined otherwise |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/applications/transfer/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/transfer/v1/query.proto



<a name="ibc.applications.transfer.v1.QueryDenomTraceRequest"></a>

### QueryDenomTraceRequest
QueryDenomTraceRequest is the request type for the Query/DenomTrace RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [string](#string) |  | hash (in hex format) of the denomination trace information. |






<a name="ibc.applications.transfer.v1.QueryDenomTraceResponse"></a>

### QueryDenomTraceResponse
QueryDenomTraceResponse is the response type for the Query/DenomTrace RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom_trace` | [DenomTrace](#ibc.applications.transfer.v1.DenomTrace) |  | denom_trace returns the requested denomination trace information. |






<a name="ibc.applications.transfer.v1.QueryDenomTracesRequest"></a>

### QueryDenomTracesRequest
QueryConnectionsRequest is the request type for the Query/DenomTraces RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="ibc.applications.transfer.v1.QueryDenomTracesResponse"></a>

### QueryDenomTracesResponse
QueryConnectionsResponse is the response type for the Query/DenomTraces RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom_traces` | [DenomTrace](#ibc.applications.transfer.v1.DenomTrace) | repeated | denom_traces returns all denominations trace information. |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="ibc.applications.transfer.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="ibc.applications.transfer.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#ibc.applications.transfer.v1.Params) |  | params defines the parameters of the module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ibc.applications.transfer.v1.Query"></a>

### Query
Query provides defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `DenomTrace` | [QueryDenomTraceRequest](#ibc.applications.transfer.v1.QueryDenomTraceRequest) | [QueryDenomTraceResponse](#ibc.applications.transfer.v1.QueryDenomTraceResponse) | DenomTrace queries a denomination trace information. | GET|/ibc/applications/transfer/v1/denom_traces/{hash}|
| `DenomTraces` | [QueryDenomTracesRequest](#ibc.applications.transfer.v1.QueryDenomTracesRequest) | [QueryDenomTracesResponse](#ibc.applications.transfer.v1.QueryDenomTracesResponse) | DenomTraces queries all denomination traces. | GET|/ibc/applications/transfer/v1/denom_traces|
| `Params` | [QueryParamsRequest](#ibc.applications.transfer.v1.QueryParamsRequest) | [QueryParamsResponse](#ibc.applications.transfer.v1.QueryParamsResponse) | Params queries all parameters of the ibc-transfer module. | GET|/ibc/applications/transfer/v1/params|

 <!-- end services -->



<a name="lbm/base/v1/coin.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/base/v1/coin.proto



<a name="lbm.base.v1.Coin"></a>

### Coin
Coin defines a token with a denomination and an amount.

NOTE: The amount field is an Int which implements the custom method
signatures required by gogoproto.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `amount` | [string](#string) |  |  |






<a name="lbm.base.v1.DecCoin"></a>

### DecCoin
DecCoin defines a token with a denomination and a decimal amount.

NOTE: The amount field is an Dec which implements the custom method
signatures required by gogoproto.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `amount` | [string](#string) |  |  |






<a name="lbm.base.v1.DecProto"></a>

### DecProto
DecProto defines a Protobuf wrapper around a Dec object.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dec` | [string](#string) |  |  |






<a name="lbm.base.v1.IntProto"></a>

### IntProto
IntProto defines a Protobuf wrapper around an Int object.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `int` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/core/client/v1/client.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/client/v1/client.proto



<a name="ibc.core.client.v1.ClientConsensusStates"></a>

### ClientConsensusStates
ClientConsensusStates defines all the stored consensus states for a given
client.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client identifier |
| `consensus_states` | [ConsensusStateWithHeight](#ibc.core.client.v1.ConsensusStateWithHeight) | repeated | consensus states and their heights associated with the client |






<a name="ibc.core.client.v1.ClientUpdateProposal"></a>

### ClientUpdateProposal
ClientUpdateProposal is a governance proposal. If it passes, the client is
updated with the provided header. The update may fail if the header is not
valid given certain conditions specified by the client implementation.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | the title of the update proposal |
| `description` | [string](#string) |  | the description of the proposal |
| `client_id` | [string](#string) |  | the client identifier for the client to be updated if the proposal passes |
| `header` | [google.protobuf.Any](#google.protobuf.Any) |  | the header used to update the client if the proposal passes |






<a name="ibc.core.client.v1.ConsensusStateWithHeight"></a>

### ConsensusStateWithHeight
ConsensusStateWithHeight defines a consensus state with an additional height field.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [Height](#ibc.core.client.v1.Height) |  | consensus state height |
| `consensus_state` | [google.protobuf.Any](#google.protobuf.Any) |  | consensus state |






<a name="ibc.core.client.v1.Height"></a>

### Height
Height is a monotonically increasing data type
that can be compared against another Height for the purposes of updating and
freezing clients

Normally the RevisionHeight is incremented at each height while keeping RevisionNumber
the same. However some consensus algorithms may choose to reset the
height in certain conditions e.g. hard forks, state-machine breaking changes
In these cases, the RevisionNumber is incremented so that height continues to
be monitonically increasing even as the RevisionHeight gets reset


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `revision_number` | [uint64](#uint64) |  | the revision that the client is currently on |
| `revision_height` | [uint64](#uint64) |  | the height within the given revision |






<a name="ibc.core.client.v1.IdentifiedClientState"></a>

### IdentifiedClientState
IdentifiedClientState defines a client state with an additional client
identifier field.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client identifier |
| `client_state` | [google.protobuf.Any](#google.protobuf.Any) |  | client state |






<a name="ibc.core.client.v1.Params"></a>

### Params
Params defines the set of IBC light client parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `allowed_clients` | [string](#string) | repeated | allowed_clients defines the list of allowed client state types. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/applications/transfer/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/transfer/v1/tx.proto



<a name="ibc.applications.transfer.v1.MsgTransfer"></a>

### MsgTransfer
MsgTransfer defines a msg to transfer fungible tokens (i.e Coins) between
ICS20 enabled chains. See ICS Spec here:
https://github.com/cosmos/ics/tree/master/spec/ics-020-fungible-token-transfer#data-structures


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `source_port` | [string](#string) |  | the port on which the packet will be sent |
| `source_channel` | [string](#string) |  | the channel by which the packet will be sent |
| `token` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) |  | the tokens to be transferred |
| `sender` | [string](#string) |  | the sender address |
| `receiver` | [string](#string) |  | the recipient address on the destination chain |
| `timeout_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | Timeout height relative to the current block height. The timeout is disabled when set to 0. |
| `timeout_timestamp` | [uint64](#uint64) |  | Timeout timestamp (in nanoseconds) relative to the current block timestamp. The timeout is disabled when set to 0. |






<a name="ibc.applications.transfer.v1.MsgTransferResponse"></a>

### MsgTransferResponse
MsgTransferResponse defines the Msg/Transfer response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ibc.applications.transfer.v1.Msg"></a>

### Msg
Msg defines the ibc/transfer Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Transfer` | [MsgTransfer](#ibc.applications.transfer.v1.MsgTransfer) | [MsgTransferResponse](#ibc.applications.transfer.v1.MsgTransferResponse) | Transfer defines a rpc handler method for MsgTransfer. | |

 <!-- end services -->



<a name="ibc/core/channel/v1/channel.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/channel/v1/channel.proto



<a name="ibc.core.channel.v1.Acknowledgement"></a>

### Acknowledgement
Acknowledgement is the recommended acknowledgement format to be used by
app-specific protocols.
NOTE: The field numbers 21 and 22 were explicitly chosen to avoid accidental
conflicts with other protobuf message formats used for acknowledgements.
The first byte of any message with this format will be the non-ASCII values
`0xaa` (result) or `0xb2` (error). Implemented as defined by ICS:
https://github.com/cosmos/ics/tree/master/spec/ics-004-channel-and-packet-semantics#acknowledgement-envelope


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `result` | [bytes](#bytes) |  |  |
| `error` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.Channel"></a>

### Channel
Channel defines pipeline for exactly-once packet delivery between specific
modules on separate blockchains, which has at least one end capable of
sending packets and one end capable of receiving packets.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `state` | [State](#ibc.core.channel.v1.State) |  | current state of the channel end |
| `ordering` | [Order](#ibc.core.channel.v1.Order) |  | whether the channel is ordered or unordered |
| `counterparty` | [Counterparty](#ibc.core.channel.v1.Counterparty) |  | counterparty channel end |
| `connection_hops` | [string](#string) | repeated | list of connection identifiers, in order, along which packets sent on this channel will travel |
| `version` | [string](#string) |  | opaque channel version, which is agreed upon during the handshake |






<a name="ibc.core.channel.v1.Counterparty"></a>

### Counterparty
Counterparty defines a channel end counterparty


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port on the counterparty chain which owns the other end of the channel. |
| `channel_id` | [string](#string) |  | channel end on the counterparty chain |






<a name="ibc.core.channel.v1.IdentifiedChannel"></a>

### IdentifiedChannel
IdentifiedChannel defines a channel with additional port and channel
identifier fields.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `state` | [State](#ibc.core.channel.v1.State) |  | current state of the channel end |
| `ordering` | [Order](#ibc.core.channel.v1.Order) |  | whether the channel is ordered or unordered |
| `counterparty` | [Counterparty](#ibc.core.channel.v1.Counterparty) |  | counterparty channel end |
| `connection_hops` | [string](#string) | repeated | list of connection identifiers, in order, along which packets sent on this channel will travel |
| `version` | [string](#string) |  | opaque channel version, which is agreed upon during the handshake |
| `port_id` | [string](#string) |  | port identifier |
| `channel_id` | [string](#string) |  | channel identifier |






<a name="ibc.core.channel.v1.Packet"></a>

### Packet
Packet defines a type that carries data across different chains through IBC


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sequence` | [uint64](#uint64) |  | number corresponds to the order of sends and receives, where a Packet with an earlier sequence number must be sent and received before a Packet with a later sequence number. |
| `source_port` | [string](#string) |  | identifies the port on the sending chain. |
| `source_channel` | [string](#string) |  | identifies the channel end on the sending chain. |
| `destination_port` | [string](#string) |  | identifies the port on the receiving chain. |
| `destination_channel` | [string](#string) |  | identifies the channel end on the receiving chain. |
| `data` | [bytes](#bytes) |  | actual opaque bytes transferred directly to the application module |
| `timeout_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | block height after which the packet times out |
| `timeout_timestamp` | [uint64](#uint64) |  | block timestamp (in nanoseconds) after which the packet times out |






<a name="ibc.core.channel.v1.PacketState"></a>

### PacketState
PacketState defines the generic type necessary to retrieve and store
packet commitments, acknowledgements, and receipts.
Caller is responsible for knowing the context necessary to interpret this
state as a commitment, acknowledgement, or a receipt.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | channel port identifier. |
| `channel_id` | [string](#string) |  | channel unique identifier. |
| `sequence` | [uint64](#uint64) |  | packet sequence. |
| `data` | [bytes](#bytes) |  | embedded data that represents packet state. |





 <!-- end messages -->


<a name="ibc.core.channel.v1.Order"></a>

### Order
Order defines if a channel is ORDERED or UNORDERED

| Name | Number | Description |
| ---- | ------ | ----------- |
| ORDER_NONE_UNSPECIFIED | 0 | zero-value for channel ordering |
| ORDER_UNORDERED | 1 | packets can be delivered in any order, which may differ from the order in which they were sent. |
| ORDER_ORDERED | 2 | packets are delivered exactly in the order which they were sent |



<a name="ibc.core.channel.v1.State"></a>

### State
State defines if a channel is in one of the following states:
CLOSED, INIT, TRYOPEN, OPEN or UNINITIALIZED.

| Name | Number | Description |
| ---- | ------ | ----------- |
| STATE_UNINITIALIZED_UNSPECIFIED | 0 | Default State |
| STATE_INIT | 1 | A channel has just started the opening handshake. |
| STATE_TRYOPEN | 2 | A channel has acknowledged the handshake step on the counterparty chain. |
| STATE_OPEN | 3 | A channel has completed the handshake. Open channels are ready to send and receive packets. |
| STATE_CLOSED | 4 | A channel has been closed and can no longer be used to send or receive packets. |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/core/channel/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/channel/v1/genesis.proto



<a name="ibc.core.channel.v1.GenesisState"></a>

### GenesisState
GenesisState defines the ibc channel submodule's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `channels` | [IdentifiedChannel](#ibc.core.channel.v1.IdentifiedChannel) | repeated |  |
| `acknowledgements` | [PacketState](#ibc.core.channel.v1.PacketState) | repeated |  |
| `commitments` | [PacketState](#ibc.core.channel.v1.PacketState) | repeated |  |
| `receipts` | [PacketState](#ibc.core.channel.v1.PacketState) | repeated |  |
| `send_sequences` | [PacketSequence](#ibc.core.channel.v1.PacketSequence) | repeated |  |
| `recv_sequences` | [PacketSequence](#ibc.core.channel.v1.PacketSequence) | repeated |  |
| `ack_sequences` | [PacketSequence](#ibc.core.channel.v1.PacketSequence) | repeated |  |
| `next_channel_sequence` | [uint64](#uint64) |  | the sequence for the next generated channel identifier |






<a name="ibc.core.channel.v1.PacketSequence"></a>

### PacketSequence
PacketSequence defines the genesis type necessary to retrieve and store
next send and receive sequences.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  |  |
| `channel_id` | [string](#string) |  |  |
| `sequence` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/core/channel/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/channel/v1/query.proto



<a name="ibc.core.channel.v1.QueryChannelClientStateRequest"></a>

### QueryChannelClientStateRequest
QueryChannelClientStateRequest is the request type for the Query/ClientState
RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port unique identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |






<a name="ibc.core.channel.v1.QueryChannelClientStateResponse"></a>

### QueryChannelClientStateResponse
QueryChannelClientStateResponse is the Response type for the
Query/QueryChannelClientState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `identified_client_state` | [ibc.core.client.v1.IdentifiedClientState](#ibc.core.client.v1.IdentifiedClientState) |  | client state associated with the channel |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.channel.v1.QueryChannelConsensusStateRequest"></a>

### QueryChannelConsensusStateRequest
QueryChannelConsensusStateRequest is the request type for the
Query/ConsensusState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port unique identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |
| `revision_number` | [uint64](#uint64) |  | revision number of the consensus state |
| `revision_height` | [uint64](#uint64) |  | revision height of the consensus state |






<a name="ibc.core.channel.v1.QueryChannelConsensusStateResponse"></a>

### QueryChannelConsensusStateResponse
QueryChannelClientStateResponse is the Response type for the
Query/QueryChannelClientState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `consensus_state` | [google.protobuf.Any](#google.protobuf.Any) |  | consensus state associated with the channel |
| `client_id` | [string](#string) |  | client ID associated with the consensus state |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.channel.v1.QueryChannelRequest"></a>

### QueryChannelRequest
QueryChannelRequest is the request type for the Query/Channel RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port unique identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |






<a name="ibc.core.channel.v1.QueryChannelResponse"></a>

### QueryChannelResponse
QueryChannelResponse is the response type for the Query/Channel RPC method.
Besides the Channel end, it includes a proof and the height from which the
proof was retrieved.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `channel` | [Channel](#ibc.core.channel.v1.Channel) |  | channel associated with the request identifiers |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.channel.v1.QueryChannelsRequest"></a>

### QueryChannelsRequest
QueryChannelsRequest is the request type for the Query/Channels RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination request |






<a name="ibc.core.channel.v1.QueryChannelsResponse"></a>

### QueryChannelsResponse
QueryChannelsResponse is the response type for the Query/Channels RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `channels` | [IdentifiedChannel](#ibc.core.channel.v1.IdentifiedChannel) | repeated | list of stored channels of the chain. |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination response |
| `height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | query block height |






<a name="ibc.core.channel.v1.QueryConnectionChannelsRequest"></a>

### QueryConnectionChannelsRequest
QueryConnectionChannelsRequest is the request type for the
Query/QueryConnectionChannels RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connection` | [string](#string) |  | connection unique identifier |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination request |






<a name="ibc.core.channel.v1.QueryConnectionChannelsResponse"></a>

### QueryConnectionChannelsResponse
QueryConnectionChannelsResponse is the Response type for the
Query/QueryConnectionChannels RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `channels` | [IdentifiedChannel](#ibc.core.channel.v1.IdentifiedChannel) | repeated | list of channels associated with a connection. |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination response |
| `height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | query block height |






<a name="ibc.core.channel.v1.QueryNextSequenceReceiveRequest"></a>

### QueryNextSequenceReceiveRequest
QueryNextSequenceReceiveRequest is the request type for the
Query/QueryNextSequenceReceiveRequest RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port unique identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |






<a name="ibc.core.channel.v1.QueryNextSequenceReceiveResponse"></a>

### QueryNextSequenceReceiveResponse
QuerySequenceResponse is the request type for the
Query/QueryNextSequenceReceiveResponse RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `next_sequence_receive` | [uint64](#uint64) |  | next sequence receive number |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.channel.v1.QueryPacketAcknowledgementRequest"></a>

### QueryPacketAcknowledgementRequest
QueryPacketAcknowledgementRequest is the request type for the
Query/PacketAcknowledgement RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port unique identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |
| `sequence` | [uint64](#uint64) |  | packet sequence |






<a name="ibc.core.channel.v1.QueryPacketAcknowledgementResponse"></a>

### QueryPacketAcknowledgementResponse
QueryPacketAcknowledgementResponse defines the client query response for a
packet which also includes a proof and the height from which the
proof was retrieved


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `acknowledgement` | [bytes](#bytes) |  | packet associated with the request fields |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.channel.v1.QueryPacketAcknowledgementsRequest"></a>

### QueryPacketAcknowledgementsRequest
QueryPacketAcknowledgementsRequest is the request type for the
Query/QueryPacketCommitments RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port unique identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination request |






<a name="ibc.core.channel.v1.QueryPacketAcknowledgementsResponse"></a>

### QueryPacketAcknowledgementsResponse
QueryPacketAcknowledgemetsResponse is the request type for the
Query/QueryPacketAcknowledgements RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `acknowledgements` | [PacketState](#ibc.core.channel.v1.PacketState) | repeated |  |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination response |
| `height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | query block height |






<a name="ibc.core.channel.v1.QueryPacketCommitmentRequest"></a>

### QueryPacketCommitmentRequest
QueryPacketCommitmentRequest is the request type for the
Query/PacketCommitment RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port unique identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |
| `sequence` | [uint64](#uint64) |  | packet sequence |






<a name="ibc.core.channel.v1.QueryPacketCommitmentResponse"></a>

### QueryPacketCommitmentResponse
QueryPacketCommitmentResponse defines the client query response for a packet
which also includes a proof and the height from which the proof was
retrieved


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commitment` | [bytes](#bytes) |  | packet associated with the request fields |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.channel.v1.QueryPacketCommitmentsRequest"></a>

### QueryPacketCommitmentsRequest
QueryPacketCommitmentsRequest is the request type for the
Query/QueryPacketCommitments RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port unique identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination request |






<a name="ibc.core.channel.v1.QueryPacketCommitmentsResponse"></a>

### QueryPacketCommitmentsResponse
QueryPacketCommitmentsResponse is the request type for the
Query/QueryPacketCommitments RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commitments` | [PacketState](#ibc.core.channel.v1.PacketState) | repeated |  |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination response |
| `height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | query block height |






<a name="ibc.core.channel.v1.QueryPacketReceiptRequest"></a>

### QueryPacketReceiptRequest
QueryPacketReceiptRequest is the request type for the
Query/PacketReceipt RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port unique identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |
| `sequence` | [uint64](#uint64) |  | packet sequence |






<a name="ibc.core.channel.v1.QueryPacketReceiptResponse"></a>

### QueryPacketReceiptResponse
QueryPacketReceiptResponse defines the client query response for a packet receipt
which also includes a proof, and the height from which the proof was
retrieved


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `received` | [bool](#bool) |  | success flag for if receipt exists |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.channel.v1.QueryUnreceivedAcksRequest"></a>

### QueryUnreceivedAcksRequest
QueryUnreceivedAcks is the request type for the
Query/UnreceivedAcks RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port unique identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |
| `packet_ack_sequences` | [uint64](#uint64) | repeated | list of acknowledgement sequences |






<a name="ibc.core.channel.v1.QueryUnreceivedAcksResponse"></a>

### QueryUnreceivedAcksResponse
QueryUnreceivedAcksResponse is the response type for the
Query/UnreceivedAcks RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sequences` | [uint64](#uint64) | repeated | list of unreceived acknowledgement sequences |
| `height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | query block height |






<a name="ibc.core.channel.v1.QueryUnreceivedPacketsRequest"></a>

### QueryUnreceivedPacketsRequest
QueryUnreceivedPacketsRequest is the request type for the
Query/UnreceivedPackets RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port unique identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |
| `packet_commitment_sequences` | [uint64](#uint64) | repeated | list of packet sequences |






<a name="ibc.core.channel.v1.QueryUnreceivedPacketsResponse"></a>

### QueryUnreceivedPacketsResponse
QueryUnreceivedPacketsResponse is the response type for the
Query/UnreceivedPacketCommitments RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sequences` | [uint64](#uint64) | repeated | list of unreceived packet sequences |
| `height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | query block height |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ibc.core.channel.v1.Query"></a>

### Query
Query provides defines the gRPC querier service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Channel` | [QueryChannelRequest](#ibc.core.channel.v1.QueryChannelRequest) | [QueryChannelResponse](#ibc.core.channel.v1.QueryChannelResponse) | Channel queries an IBC Channel. | GET|/ibc/core/channel/v1/channels/{channel_id}/ports/{port_id}|
| `Channels` | [QueryChannelsRequest](#ibc.core.channel.v1.QueryChannelsRequest) | [QueryChannelsResponse](#ibc.core.channel.v1.QueryChannelsResponse) | Channels queries all the IBC channels of a chain. | GET|/ibc/core/channel/v1/channels|
| `ConnectionChannels` | [QueryConnectionChannelsRequest](#ibc.core.channel.v1.QueryConnectionChannelsRequest) | [QueryConnectionChannelsResponse](#ibc.core.channel.v1.QueryConnectionChannelsResponse) | ConnectionChannels queries all the channels associated with a connection end. | GET|/ibc/core/channel/v1/connections/{connection}/channels|
| `ChannelClientState` | [QueryChannelClientStateRequest](#ibc.core.channel.v1.QueryChannelClientStateRequest) | [QueryChannelClientStateResponse](#ibc.core.channel.v1.QueryChannelClientStateResponse) | ChannelClientState queries for the client state for the channel associated with the provided channel identifiers. | GET|/ibc/core/channel/v1/channels/{channel_id}/ports/{port_id}/client_state|
| `ChannelConsensusState` | [QueryChannelConsensusStateRequest](#ibc.core.channel.v1.QueryChannelConsensusStateRequest) | [QueryChannelConsensusStateResponse](#ibc.core.channel.v1.QueryChannelConsensusStateResponse) | ChannelConsensusState queries for the consensus state for the channel associated with the provided channel identifiers. | GET|/ibc/core/channel/v1/channels/{channel_id}/ports/{port_id}/consensus_state/revision/{revision_number}/height/{revision_height}|
| `PacketCommitment` | [QueryPacketCommitmentRequest](#ibc.core.channel.v1.QueryPacketCommitmentRequest) | [QueryPacketCommitmentResponse](#ibc.core.channel.v1.QueryPacketCommitmentResponse) | PacketCommitment queries a stored packet commitment hash. | GET|/ibc/core/channel/v1/channels/{channel_id}/ports/{port_id}/packet_commitments/{sequence}|
| `PacketCommitments` | [QueryPacketCommitmentsRequest](#ibc.core.channel.v1.QueryPacketCommitmentsRequest) | [QueryPacketCommitmentsResponse](#ibc.core.channel.v1.QueryPacketCommitmentsResponse) | PacketCommitments returns all the packet commitments hashes associated with a channel. | GET|/ibc/core/channel/v1/channels/{channel_id}/ports/{port_id}/packet_commitments|
| `PacketReceipt` | [QueryPacketReceiptRequest](#ibc.core.channel.v1.QueryPacketReceiptRequest) | [QueryPacketReceiptResponse](#ibc.core.channel.v1.QueryPacketReceiptResponse) | PacketReceipt queries if a given packet sequence has been received on the queried chain | GET|/ibc/core/channel/v1/channels/{channel_id}/ports/{port_id}/packet_receipts/{sequence}|
| `PacketAcknowledgement` | [QueryPacketAcknowledgementRequest](#ibc.core.channel.v1.QueryPacketAcknowledgementRequest) | [QueryPacketAcknowledgementResponse](#ibc.core.channel.v1.QueryPacketAcknowledgementResponse) | PacketAcknowledgement queries a stored packet acknowledgement hash. | GET|/ibc/core/channel/v1/channels/{channel_id}/ports/{port_id}/packet_acks/{sequence}|
| `PacketAcknowledgements` | [QueryPacketAcknowledgementsRequest](#ibc.core.channel.v1.QueryPacketAcknowledgementsRequest) | [QueryPacketAcknowledgementsResponse](#ibc.core.channel.v1.QueryPacketAcknowledgementsResponse) | PacketAcknowledgements returns all the packet acknowledgements associated with a channel. | GET|/ibc/core/channel/v1/channels/{channel_id}/ports/{port_id}/packet_acknowledgements|
| `UnreceivedPackets` | [QueryUnreceivedPacketsRequest](#ibc.core.channel.v1.QueryUnreceivedPacketsRequest) | [QueryUnreceivedPacketsResponse](#ibc.core.channel.v1.QueryUnreceivedPacketsResponse) | UnreceivedPackets returns all the unreceived IBC packets associated with a channel and sequences. | GET|/ibc/core/channel/v1/channels/{channel_id}/ports/{port_id}/packet_commitments/{packet_commitment_sequences}/unreceived_packets|
| `UnreceivedAcks` | [QueryUnreceivedAcksRequest](#ibc.core.channel.v1.QueryUnreceivedAcksRequest) | [QueryUnreceivedAcksResponse](#ibc.core.channel.v1.QueryUnreceivedAcksResponse) | UnreceivedAcks returns all the unreceived IBC acknowledgements associated with a channel and sequences. | GET|/ibc/core/channel/v1/channels/{channel_id}/ports/{port_id}/packet_commitments/{packet_ack_sequences}/unreceived_acks|
| `NextSequenceReceive` | [QueryNextSequenceReceiveRequest](#ibc.core.channel.v1.QueryNextSequenceReceiveRequest) | [QueryNextSequenceReceiveResponse](#ibc.core.channel.v1.QueryNextSequenceReceiveResponse) | NextSequenceReceive returns the next receive sequence for a given channel. | GET|/ibc/core/channel/v1/channels/{channel_id}/ports/{port_id}/next_sequence|

 <!-- end services -->



<a name="ibc/core/channel/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/channel/v1/tx.proto



<a name="ibc.core.channel.v1.MsgAcknowledgement"></a>

### MsgAcknowledgement
MsgAcknowledgement receives incoming IBC acknowledgement


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `packet` | [Packet](#ibc.core.channel.v1.Packet) |  |  |
| `acknowledgement` | [bytes](#bytes) |  |  |
| `proof_acked` | [bytes](#bytes) |  |  |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.MsgAcknowledgementResponse"></a>

### MsgAcknowledgementResponse
MsgAcknowledgementResponse defines the Msg/Acknowledgement response type.






<a name="ibc.core.channel.v1.MsgChannelCloseConfirm"></a>

### MsgChannelCloseConfirm
MsgChannelCloseConfirm defines a msg sent by a Relayer to Chain B
to acknowledge the change of channel state to CLOSED on Chain A.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  |  |
| `channel_id` | [string](#string) |  |  |
| `proof_init` | [bytes](#bytes) |  |  |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.MsgChannelCloseConfirmResponse"></a>

### MsgChannelCloseConfirmResponse
MsgChannelCloseConfirmResponse defines the Msg/ChannelCloseConfirm response type.






<a name="ibc.core.channel.v1.MsgChannelCloseInit"></a>

### MsgChannelCloseInit
MsgChannelCloseInit defines a msg sent by a Relayer to Chain A
to close a channel with Chain B.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  |  |
| `channel_id` | [string](#string) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.MsgChannelCloseInitResponse"></a>

### MsgChannelCloseInitResponse
MsgChannelCloseInitResponse defines the Msg/ChannelCloseInit response type.






<a name="ibc.core.channel.v1.MsgChannelOpenAck"></a>

### MsgChannelOpenAck
MsgChannelOpenAck defines a msg sent by a Relayer to Chain A to acknowledge
the change of channel state to TRYOPEN on Chain B.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  |  |
| `channel_id` | [string](#string) |  |  |
| `counterparty_channel_id` | [string](#string) |  |  |
| `counterparty_version` | [string](#string) |  |  |
| `proof_try` | [bytes](#bytes) |  |  |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.MsgChannelOpenAckResponse"></a>

### MsgChannelOpenAckResponse
MsgChannelOpenAckResponse defines the Msg/ChannelOpenAck response type.






<a name="ibc.core.channel.v1.MsgChannelOpenConfirm"></a>

### MsgChannelOpenConfirm
MsgChannelOpenConfirm defines a msg sent by a Relayer to Chain B to
acknowledge the change of channel state to OPEN on Chain A.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  |  |
| `channel_id` | [string](#string) |  |  |
| `proof_ack` | [bytes](#bytes) |  |  |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.MsgChannelOpenConfirmResponse"></a>

### MsgChannelOpenConfirmResponse
MsgChannelOpenConfirmResponse defines the Msg/ChannelOpenConfirm response type.






<a name="ibc.core.channel.v1.MsgChannelOpenInit"></a>

### MsgChannelOpenInit
MsgChannelOpenInit defines an sdk.Msg to initialize a channel handshake. It
is called by a relayer on Chain A.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  |  |
| `channel` | [Channel](#ibc.core.channel.v1.Channel) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.MsgChannelOpenInitResponse"></a>

### MsgChannelOpenInitResponse
MsgChannelOpenInitResponse defines the Msg/ChannelOpenInit response type.






<a name="ibc.core.channel.v1.MsgChannelOpenTry"></a>

### MsgChannelOpenTry
MsgChannelOpenInit defines a msg sent by a Relayer to try to open a channel
on Chain B.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  |  |
| `previous_channel_id` | [string](#string) |  | in the case of crossing hello's, when both chains call OpenInit, we need the channel identifier of the previous channel in state INIT |
| `channel` | [Channel](#ibc.core.channel.v1.Channel) |  |  |
| `counterparty_version` | [string](#string) |  |  |
| `proof_init` | [bytes](#bytes) |  |  |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.MsgChannelOpenTryResponse"></a>

### MsgChannelOpenTryResponse
MsgChannelOpenTryResponse defines the Msg/ChannelOpenTry response type.






<a name="ibc.core.channel.v1.MsgRecvPacket"></a>

### MsgRecvPacket
MsgRecvPacket receives incoming IBC packet


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `packet` | [Packet](#ibc.core.channel.v1.Packet) |  |  |
| `proof_commitment` | [bytes](#bytes) |  |  |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.MsgRecvPacketResponse"></a>

### MsgRecvPacketResponse
MsgRecvPacketResponse defines the Msg/RecvPacket response type.






<a name="ibc.core.channel.v1.MsgTimeout"></a>

### MsgTimeout
MsgTimeout receives timed-out packet


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `packet` | [Packet](#ibc.core.channel.v1.Packet) |  |  |
| `proof_unreceived` | [bytes](#bytes) |  |  |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `next_sequence_recv` | [uint64](#uint64) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.MsgTimeoutOnClose"></a>

### MsgTimeoutOnClose
MsgTimeoutOnClose timed-out packet upon counterparty channel closure.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `packet` | [Packet](#ibc.core.channel.v1.Packet) |  |  |
| `proof_unreceived` | [bytes](#bytes) |  |  |
| `proof_close` | [bytes](#bytes) |  |  |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `next_sequence_recv` | [uint64](#uint64) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.MsgTimeoutOnCloseResponse"></a>

### MsgTimeoutOnCloseResponse
MsgTimeoutOnCloseResponse defines the Msg/TimeoutOnClose response type.






<a name="ibc.core.channel.v1.MsgTimeoutResponse"></a>

### MsgTimeoutResponse
MsgTimeoutResponse defines the Msg/Timeout response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ibc.core.channel.v1.Msg"></a>

### Msg
Msg defines the ibc/channel Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `ChannelOpenInit` | [MsgChannelOpenInit](#ibc.core.channel.v1.MsgChannelOpenInit) | [MsgChannelOpenInitResponse](#ibc.core.channel.v1.MsgChannelOpenInitResponse) | ChannelOpenInit defines a rpc handler method for MsgChannelOpenInit. | |
| `ChannelOpenTry` | [MsgChannelOpenTry](#ibc.core.channel.v1.MsgChannelOpenTry) | [MsgChannelOpenTryResponse](#ibc.core.channel.v1.MsgChannelOpenTryResponse) | ChannelOpenTry defines a rpc handler method for MsgChannelOpenTry. | |
| `ChannelOpenAck` | [MsgChannelOpenAck](#ibc.core.channel.v1.MsgChannelOpenAck) | [MsgChannelOpenAckResponse](#ibc.core.channel.v1.MsgChannelOpenAckResponse) | ChannelOpenAck defines a rpc handler method for MsgChannelOpenAck. | |
| `ChannelOpenConfirm` | [MsgChannelOpenConfirm](#ibc.core.channel.v1.MsgChannelOpenConfirm) | [MsgChannelOpenConfirmResponse](#ibc.core.channel.v1.MsgChannelOpenConfirmResponse) | ChannelOpenConfirm defines a rpc handler method for MsgChannelOpenConfirm. | |
| `ChannelCloseInit` | [MsgChannelCloseInit](#ibc.core.channel.v1.MsgChannelCloseInit) | [MsgChannelCloseInitResponse](#ibc.core.channel.v1.MsgChannelCloseInitResponse) | ChannelCloseInit defines a rpc handler method for MsgChannelCloseInit. | |
| `ChannelCloseConfirm` | [MsgChannelCloseConfirm](#ibc.core.channel.v1.MsgChannelCloseConfirm) | [MsgChannelCloseConfirmResponse](#ibc.core.channel.v1.MsgChannelCloseConfirmResponse) | ChannelCloseConfirm defines a rpc handler method for MsgChannelCloseConfirm. | |
| `RecvPacket` | [MsgRecvPacket](#ibc.core.channel.v1.MsgRecvPacket) | [MsgRecvPacketResponse](#ibc.core.channel.v1.MsgRecvPacketResponse) | RecvPacket defines a rpc handler method for MsgRecvPacket. | |
| `Timeout` | [MsgTimeout](#ibc.core.channel.v1.MsgTimeout) | [MsgTimeoutResponse](#ibc.core.channel.v1.MsgTimeoutResponse) | Timeout defines a rpc handler method for MsgTimeout. | |
| `TimeoutOnClose` | [MsgTimeoutOnClose](#ibc.core.channel.v1.MsgTimeoutOnClose) | [MsgTimeoutOnCloseResponse](#ibc.core.channel.v1.MsgTimeoutOnCloseResponse) | TimeoutOnClose defines a rpc handler method for MsgTimeoutOnClose. | |
| `Acknowledgement` | [MsgAcknowledgement](#ibc.core.channel.v1.MsgAcknowledgement) | [MsgAcknowledgementResponse](#ibc.core.channel.v1.MsgAcknowledgementResponse) | Acknowledgement defines a rpc handler method for MsgAcknowledgement. | |

 <!-- end services -->



<a name="ibc/core/client/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/client/v1/genesis.proto



<a name="ibc.core.client.v1.GenesisMetadata"></a>

### GenesisMetadata
GenesisMetadata defines the genesis type for metadata that clients may return
with ExportMetadata


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  | store key of metadata without clientID-prefix |
| `value` | [bytes](#bytes) |  | metadata value |






<a name="ibc.core.client.v1.GenesisState"></a>

### GenesisState
GenesisState defines the ibc client submodule's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `clients` | [IdentifiedClientState](#ibc.core.client.v1.IdentifiedClientState) | repeated | client states with their corresponding identifiers |
| `clients_consensus` | [ClientConsensusStates](#ibc.core.client.v1.ClientConsensusStates) | repeated | consensus states from each client |
| `clients_metadata` | [IdentifiedGenesisMetadata](#ibc.core.client.v1.IdentifiedGenesisMetadata) | repeated | metadata from each client |
| `params` | [Params](#ibc.core.client.v1.Params) |  |  |
| `create_localhost` | [bool](#bool) |  | create localhost on initialization |
| `next_client_sequence` | [uint64](#uint64) |  | the sequence for the next generated client identifier |






<a name="ibc.core.client.v1.IdentifiedGenesisMetadata"></a>

### IdentifiedGenesisMetadata
IdentifiedGenesisMetadata has the client metadata with the corresponding client id.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  |  |
| `client_metadata` | [GenesisMetadata](#ibc.core.client.v1.GenesisMetadata) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/core/client/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/client/v1/query.proto



<a name="ibc.core.client.v1.QueryClientParamsRequest"></a>

### QueryClientParamsRequest
QueryClientParamsRequest is the request type for the Query/ClientParams RPC method.






<a name="ibc.core.client.v1.QueryClientParamsResponse"></a>

### QueryClientParamsResponse
QueryClientParamsResponse is the response type for the Query/ClientParams RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#ibc.core.client.v1.Params) |  | params defines the parameters of the module. |






<a name="ibc.core.client.v1.QueryClientStateRequest"></a>

### QueryClientStateRequest
QueryClientStateRequest is the request type for the Query/ClientState RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client state unique identifier |






<a name="ibc.core.client.v1.QueryClientStateResponse"></a>

### QueryClientStateResponse
QueryClientStateResponse is the response type for the Query/ClientState RPC
method. Besides the client state, it includes a proof and the height from
which the proof was retrieved.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_state` | [google.protobuf.Any](#google.protobuf.Any) |  | client state associated with the request identifier |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.client.v1.QueryClientStatesRequest"></a>

### QueryClientStatesRequest
QueryClientStatesRequest is the request type for the Query/ClientStates RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination request |






<a name="ibc.core.client.v1.QueryClientStatesResponse"></a>

### QueryClientStatesResponse
QueryClientStatesResponse is the response type for the Query/ClientStates RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_states` | [IdentifiedClientState](#ibc.core.client.v1.IdentifiedClientState) | repeated | list of stored ClientStates of the chain. |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination response |






<a name="ibc.core.client.v1.QueryConsensusStateRequest"></a>

### QueryConsensusStateRequest
QueryConsensusStateRequest is the request type for the Query/ConsensusState
RPC method. Besides the consensus state, it includes a proof and the height
from which the proof was retrieved.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client identifier |
| `revision_number` | [uint64](#uint64) |  | consensus state revision number |
| `revision_height` | [uint64](#uint64) |  | consensus state revision height |
| `latest_height` | [bool](#bool) |  | latest_height overrrides the height field and queries the latest stored ConsensusState |






<a name="ibc.core.client.v1.QueryConsensusStateResponse"></a>

### QueryConsensusStateResponse
QueryConsensusStateResponse is the response type for the Query/ConsensusState
RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `consensus_state` | [google.protobuf.Any](#google.protobuf.Any) |  | consensus state associated with the client identifier at the given height |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.client.v1.QueryConsensusStatesRequest"></a>

### QueryConsensusStatesRequest
QueryConsensusStatesRequest is the request type for the Query/ConsensusStates
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client identifier |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination request |






<a name="ibc.core.client.v1.QueryConsensusStatesResponse"></a>

### QueryConsensusStatesResponse
QueryConsensusStatesResponse is the response type for the
Query/ConsensusStates RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `consensus_states` | [ConsensusStateWithHeight](#ibc.core.client.v1.ConsensusStateWithHeight) | repeated | consensus states associated with the identifier |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination response |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ibc.core.client.v1.Query"></a>

### Query
Query provides defines the gRPC querier service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `ClientState` | [QueryClientStateRequest](#ibc.core.client.v1.QueryClientStateRequest) | [QueryClientStateResponse](#ibc.core.client.v1.QueryClientStateResponse) | ClientState queries an IBC light client. | GET|/ibc/core/client/v1/client_states/{client_id}|
| `ClientStates` | [QueryClientStatesRequest](#ibc.core.client.v1.QueryClientStatesRequest) | [QueryClientStatesResponse](#ibc.core.client.v1.QueryClientStatesResponse) | ClientStates queries all the IBC light clients of a chain. | GET|/ibc/core/client/v1/client_states|
| `ConsensusState` | [QueryConsensusStateRequest](#ibc.core.client.v1.QueryConsensusStateRequest) | [QueryConsensusStateResponse](#ibc.core.client.v1.QueryConsensusStateResponse) | ConsensusState queries a consensus state associated with a client state at a given height. | GET|/ibc/core/client/v1/consensus_states/{client_id}/revision/{revision_number}/height/{revision_height}|
| `ConsensusStates` | [QueryConsensusStatesRequest](#ibc.core.client.v1.QueryConsensusStatesRequest) | [QueryConsensusStatesResponse](#ibc.core.client.v1.QueryConsensusStatesResponse) | ConsensusStates queries all the consensus state associated with a given client. | GET|/ibc/core/client/v1/consensus_states/{client_id}|
| `ClientParams` | [QueryClientParamsRequest](#ibc.core.client.v1.QueryClientParamsRequest) | [QueryClientParamsResponse](#ibc.core.client.v1.QueryClientParamsResponse) | ClientParams queries all parameters of the ibc client. | GET|/ibc/client/v1/params|

 <!-- end services -->



<a name="ibc/core/client/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/client/v1/tx.proto



<a name="ibc.core.client.v1.MsgCreateClient"></a>

### MsgCreateClient
MsgCreateClient defines a message to create an IBC client


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_state` | [google.protobuf.Any](#google.protobuf.Any) |  | light client state |
| `consensus_state` | [google.protobuf.Any](#google.protobuf.Any) |  | consensus state associated with the client that corresponds to a given height. |
| `signer` | [string](#string) |  | signer address |






<a name="ibc.core.client.v1.MsgCreateClientResponse"></a>

### MsgCreateClientResponse
MsgCreateClientResponse defines the Msg/CreateClient response type.






<a name="ibc.core.client.v1.MsgSubmitMisbehaviour"></a>

### MsgSubmitMisbehaviour
MsgSubmitMisbehaviour defines an sdk.Msg type that submits Evidence for
light client misbehaviour.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client unique identifier |
| `misbehaviour` | [google.protobuf.Any](#google.protobuf.Any) |  | misbehaviour used for freezing the light client |
| `signer` | [string](#string) |  | signer address |






<a name="ibc.core.client.v1.MsgSubmitMisbehaviourResponse"></a>

### MsgSubmitMisbehaviourResponse
MsgSubmitMisbehaviourResponse defines the Msg/SubmitMisbehaviour response type.






<a name="ibc.core.client.v1.MsgUpdateClient"></a>

### MsgUpdateClient
MsgUpdateClient defines an sdk.Msg to update a IBC client state using
the given header.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client unique identifier |
| `header` | [google.protobuf.Any](#google.protobuf.Any) |  | header to update the light client |
| `signer` | [string](#string) |  | signer address |






<a name="ibc.core.client.v1.MsgUpdateClientResponse"></a>

### MsgUpdateClientResponse
MsgUpdateClientResponse defines the Msg/UpdateClient response type.






<a name="ibc.core.client.v1.MsgUpgradeClient"></a>

### MsgUpgradeClient
MsgUpgradeClient defines an sdk.Msg to upgrade an IBC client to a new client state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client unique identifier |
| `client_state` | [google.protobuf.Any](#google.protobuf.Any) |  | upgraded client state |
| `consensus_state` | [google.protobuf.Any](#google.protobuf.Any) |  | upgraded consensus state, only contains enough information to serve as a basis of trust in update logic |
| `proof_upgrade_client` | [bytes](#bytes) |  | proof that old chain committed to new client |
| `proof_upgrade_consensus_state` | [bytes](#bytes) |  | proof that old chain committed to new consensus state |
| `signer` | [string](#string) |  | signer address |






<a name="ibc.core.client.v1.MsgUpgradeClientResponse"></a>

### MsgUpgradeClientResponse
MsgUpgradeClientResponse defines the Msg/UpgradeClient response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ibc.core.client.v1.Msg"></a>

### Msg
Msg defines the ibc/client Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `CreateClient` | [MsgCreateClient](#ibc.core.client.v1.MsgCreateClient) | [MsgCreateClientResponse](#ibc.core.client.v1.MsgCreateClientResponse) | CreateClient defines a rpc handler method for MsgCreateClient. | |
| `UpdateClient` | [MsgUpdateClient](#ibc.core.client.v1.MsgUpdateClient) | [MsgUpdateClientResponse](#ibc.core.client.v1.MsgUpdateClientResponse) | UpdateClient defines a rpc handler method for MsgUpdateClient. | |
| `UpgradeClient` | [MsgUpgradeClient](#ibc.core.client.v1.MsgUpgradeClient) | [MsgUpgradeClientResponse](#ibc.core.client.v1.MsgUpgradeClientResponse) | UpgradeClient defines a rpc handler method for MsgUpgradeClient. | |
| `SubmitMisbehaviour` | [MsgSubmitMisbehaviour](#ibc.core.client.v1.MsgSubmitMisbehaviour) | [MsgSubmitMisbehaviourResponse](#ibc.core.client.v1.MsgSubmitMisbehaviourResponse) | SubmitMisbehaviour defines a rpc handler method for MsgSubmitMisbehaviour. | |

 <!-- end services -->



<a name="ibc/core/commitment/v1/commitment.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/commitment/v1/commitment.proto



<a name="ibc.core.commitment.v1.MerklePath"></a>

### MerklePath
MerklePath is the path used to verify commitment proofs, which can be an
arbitrary structured object (defined by a commitment type).
MerklePath is represented from root-to-leaf


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key_path` | [string](#string) | repeated |  |






<a name="ibc.core.commitment.v1.MerklePrefix"></a>

### MerklePrefix
MerklePrefix is merkle path prefixed to the key.
The constructed key from the Path and the key will be append(Path.KeyPath,
append(Path.KeyPrefix, key...))


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key_prefix` | [bytes](#bytes) |  |  |






<a name="ibc.core.commitment.v1.MerkleProof"></a>

### MerkleProof
MerkleProof is a wrapper type over a chain of CommitmentProofs.
It demonstrates membership or non-membership for an element or set of
elements, verifiable in conjunction with a known commitment root. Proofs
should be succinct.
MerkleProofs are ordered from leaf-to-root


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proofs` | [ics23.CommitmentProof](#ics23.CommitmentProof) | repeated |  |






<a name="ibc.core.commitment.v1.MerkleRoot"></a>

### MerkleRoot
MerkleRoot defines a merkle root hash.
In the Cosmos SDK, the AppHash of a block header becomes the root.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [bytes](#bytes) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/core/connection/v1/connection.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/connection/v1/connection.proto



<a name="ibc.core.connection.v1.ClientPaths"></a>

### ClientPaths
ClientPaths define all the connection paths for a client state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `paths` | [string](#string) | repeated | list of connection paths |






<a name="ibc.core.connection.v1.ConnectionEnd"></a>

### ConnectionEnd
ConnectionEnd defines a stateful object on a chain connected to another
separate one.
NOTE: there must only be 2 defined ConnectionEnds to establish
a connection between two chains.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client associated with this connection. |
| `versions` | [Version](#ibc.core.connection.v1.Version) | repeated | IBC version which can be utilised to determine encodings or protocols for channels or packets utilising this connection. |
| `state` | [State](#ibc.core.connection.v1.State) |  | current state of the connection end. |
| `counterparty` | [Counterparty](#ibc.core.connection.v1.Counterparty) |  | counterparty chain associated with this connection. |
| `delay_period` | [uint64](#uint64) |  | delay period that must pass before a consensus state can be used for packet-verification NOTE: delay period logic is only implemented by some clients. |






<a name="ibc.core.connection.v1.ConnectionPaths"></a>

### ConnectionPaths
ConnectionPaths define all the connection paths for a given client state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client state unique identifier |
| `paths` | [string](#string) | repeated | list of connection paths |






<a name="ibc.core.connection.v1.Counterparty"></a>

### Counterparty
Counterparty defines the counterparty chain associated with a connection end.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | identifies the client on the counterparty chain associated with a given connection. |
| `connection_id` | [string](#string) |  | identifies the connection end on the counterparty chain associated with a given connection. |
| `prefix` | [ibc.core.commitment.v1.MerklePrefix](#ibc.core.commitment.v1.MerklePrefix) |  | commitment merkle prefix of the counterparty chain. |






<a name="ibc.core.connection.v1.IdentifiedConnection"></a>

### IdentifiedConnection
IdentifiedConnection defines a connection with additional connection
identifier field.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  | connection identifier. |
| `client_id` | [string](#string) |  | client associated with this connection. |
| `versions` | [Version](#ibc.core.connection.v1.Version) | repeated | IBC version which can be utilised to determine encodings or protocols for channels or packets utilising this connection |
| `state` | [State](#ibc.core.connection.v1.State) |  | current state of the connection end. |
| `counterparty` | [Counterparty](#ibc.core.connection.v1.Counterparty) |  | counterparty chain associated with this connection. |
| `delay_period` | [uint64](#uint64) |  | delay period associated with this connection. |






<a name="ibc.core.connection.v1.Version"></a>

### Version
Version defines the versioning scheme used to negotiate the IBC verison in
the connection handshake.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `identifier` | [string](#string) |  | unique version identifier |
| `features` | [string](#string) | repeated | list of features compatible with the specified identifier |





 <!-- end messages -->


<a name="ibc.core.connection.v1.State"></a>

### State
State defines if a connection is in one of the following states:
INIT, TRYOPEN, OPEN or UNINITIALIZED.

| Name | Number | Description |
| ---- | ------ | ----------- |
| STATE_UNINITIALIZED_UNSPECIFIED | 0 | Default State |
| STATE_INIT | 1 | A connection end has just started the opening handshake. |
| STATE_TRYOPEN | 2 | A connection end has acknowledged the handshake step on the counterparty chain. |
| STATE_OPEN | 3 | A connection end has completed the handshake. |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/core/connection/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/connection/v1/genesis.proto



<a name="ibc.core.connection.v1.GenesisState"></a>

### GenesisState
GenesisState defines the ibc connection submodule's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connections` | [IdentifiedConnection](#ibc.core.connection.v1.IdentifiedConnection) | repeated |  |
| `client_connection_paths` | [ConnectionPaths](#ibc.core.connection.v1.ConnectionPaths) | repeated |  |
| `next_connection_sequence` | [uint64](#uint64) |  | the sequence for the next generated connection identifier |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/core/connection/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/connection/v1/query.proto



<a name="ibc.core.connection.v1.QueryClientConnectionsRequest"></a>

### QueryClientConnectionsRequest
QueryClientConnectionsRequest is the request type for the
Query/ClientConnections RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client identifier associated with a connection |






<a name="ibc.core.connection.v1.QueryClientConnectionsResponse"></a>

### QueryClientConnectionsResponse
QueryClientConnectionsResponse is the response type for the
Query/ClientConnections RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connection_paths` | [string](#string) | repeated | slice of all the connection paths associated with a client. |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | height at which the proof was generated |






<a name="ibc.core.connection.v1.QueryConnectionClientStateRequest"></a>

### QueryConnectionClientStateRequest
QueryConnectionClientStateRequest is the request type for the
Query/ConnectionClientState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connection_id` | [string](#string) |  | connection identifier |






<a name="ibc.core.connection.v1.QueryConnectionClientStateResponse"></a>

### QueryConnectionClientStateResponse
QueryConnectionClientStateResponse is the response type for the
Query/ConnectionClientState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `identified_client_state` | [ibc.core.client.v1.IdentifiedClientState](#ibc.core.client.v1.IdentifiedClientState) |  | client state associated with the channel |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.connection.v1.QueryConnectionConsensusStateRequest"></a>

### QueryConnectionConsensusStateRequest
QueryConnectionConsensusStateRequest is the request type for the
Query/ConnectionConsensusState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connection_id` | [string](#string) |  | connection identifier |
| `revision_number` | [uint64](#uint64) |  |  |
| `revision_height` | [uint64](#uint64) |  |  |






<a name="ibc.core.connection.v1.QueryConnectionConsensusStateResponse"></a>

### QueryConnectionConsensusStateResponse
QueryConnectionConsensusStateResponse is the response type for the
Query/ConnectionConsensusState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `consensus_state` | [google.protobuf.Any](#google.protobuf.Any) |  | consensus state associated with the channel |
| `client_id` | [string](#string) |  | client ID associated with the consensus state |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.connection.v1.QueryConnectionRequest"></a>

### QueryConnectionRequest
QueryConnectionRequest is the request type for the Query/Connection RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connection_id` | [string](#string) |  | connection unique identifier |






<a name="ibc.core.connection.v1.QueryConnectionResponse"></a>

### QueryConnectionResponse
QueryConnectionResponse is the response type for the Query/Connection RPC
method. Besides the connection end, it includes a proof and the height from
which the proof was retrieved.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connection` | [ConnectionEnd](#ibc.core.connection.v1.ConnectionEnd) |  | connection associated with the request identifier |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.connection.v1.QueryConnectionsRequest"></a>

### QueryConnectionsRequest
QueryConnectionsRequest is the request type for the Query/Connections RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  |  |






<a name="ibc.core.connection.v1.QueryConnectionsResponse"></a>

### QueryConnectionsResponse
QueryConnectionsResponse is the response type for the Query/Connections RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connections` | [IdentifiedConnection](#ibc.core.connection.v1.IdentifiedConnection) | repeated | list of stored connections of the chain. |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination response |
| `height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | query block height |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ibc.core.connection.v1.Query"></a>

### Query
Query provides defines the gRPC querier service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Connection` | [QueryConnectionRequest](#ibc.core.connection.v1.QueryConnectionRequest) | [QueryConnectionResponse](#ibc.core.connection.v1.QueryConnectionResponse) | Connection queries an IBC connection end. | GET|/ibc/core/connection/v1/connections/{connection_id}|
| `Connections` | [QueryConnectionsRequest](#ibc.core.connection.v1.QueryConnectionsRequest) | [QueryConnectionsResponse](#ibc.core.connection.v1.QueryConnectionsResponse) | Connections queries all the IBC connections of a chain. | GET|/ibc/core/connection/v1/connections|
| `ClientConnections` | [QueryClientConnectionsRequest](#ibc.core.connection.v1.QueryClientConnectionsRequest) | [QueryClientConnectionsResponse](#ibc.core.connection.v1.QueryClientConnectionsResponse) | ClientConnections queries the connection paths associated with a client state. | GET|/ibc/core/connection/v1/client_connections/{client_id}|
| `ConnectionClientState` | [QueryConnectionClientStateRequest](#ibc.core.connection.v1.QueryConnectionClientStateRequest) | [QueryConnectionClientStateResponse](#ibc.core.connection.v1.QueryConnectionClientStateResponse) | ConnectionClientState queries the client state associated with the connection. | GET|/ibc/core/connection/v1/connections/{connection_id}/client_state|
| `ConnectionConsensusState` | [QueryConnectionConsensusStateRequest](#ibc.core.connection.v1.QueryConnectionConsensusStateRequest) | [QueryConnectionConsensusStateResponse](#ibc.core.connection.v1.QueryConnectionConsensusStateResponse) | ConnectionConsensusState queries the consensus state associated with the connection. | GET|/ibc/core/connection/v1/connections/{connection_id}/consensus_state/revision/{revision_number}/height/{revision_height}|

 <!-- end services -->



<a name="ibc/core/connection/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/connection/v1/tx.proto



<a name="ibc.core.connection.v1.MsgConnectionOpenAck"></a>

### MsgConnectionOpenAck
MsgConnectionOpenAck defines a msg sent by a Relayer to Chain A to
acknowledge the change of connection state to TRYOPEN on Chain B.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connection_id` | [string](#string) |  |  |
| `counterparty_connection_id` | [string](#string) |  |  |
| `version` | [Version](#ibc.core.connection.v1.Version) |  |  |
| `client_state` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `proof_try` | [bytes](#bytes) |  | proof of the initialization the connection on Chain B: `UNITIALIZED -> TRYOPEN` |
| `proof_client` | [bytes](#bytes) |  | proof of client state included in message |
| `proof_consensus` | [bytes](#bytes) |  | proof of client consensus state |
| `consensus_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.connection.v1.MsgConnectionOpenAckResponse"></a>

### MsgConnectionOpenAckResponse
MsgConnectionOpenAckResponse defines the Msg/ConnectionOpenAck response type.






<a name="ibc.core.connection.v1.MsgConnectionOpenConfirm"></a>

### MsgConnectionOpenConfirm
MsgConnectionOpenConfirm defines a msg sent by a Relayer to Chain B to
acknowledge the change of connection state to OPEN on Chain A.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connection_id` | [string](#string) |  |  |
| `proof_ack` | [bytes](#bytes) |  | proof for the change of the connection state on Chain A: `INIT -> OPEN` |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.connection.v1.MsgConnectionOpenConfirmResponse"></a>

### MsgConnectionOpenConfirmResponse
MsgConnectionOpenConfirmResponse defines the Msg/ConnectionOpenConfirm response type.






<a name="ibc.core.connection.v1.MsgConnectionOpenInit"></a>

### MsgConnectionOpenInit
MsgConnectionOpenInit defines the msg sent by an account on Chain A to
initialize a connection with Chain B.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  |  |
| `counterparty` | [Counterparty](#ibc.core.connection.v1.Counterparty) |  |  |
| `version` | [Version](#ibc.core.connection.v1.Version) |  |  |
| `delay_period` | [uint64](#uint64) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.connection.v1.MsgConnectionOpenInitResponse"></a>

### MsgConnectionOpenInitResponse
MsgConnectionOpenInitResponse defines the Msg/ConnectionOpenInit response type.






<a name="ibc.core.connection.v1.MsgConnectionOpenTry"></a>

### MsgConnectionOpenTry
MsgConnectionOpenTry defines a msg sent by a Relayer to try to open a
connection on Chain B.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  |  |
| `previous_connection_id` | [string](#string) |  | in the case of crossing hello's, when both chains call OpenInit, we need the connection identifier of the previous connection in state INIT |
| `client_state` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `counterparty` | [Counterparty](#ibc.core.connection.v1.Counterparty) |  |  |
| `delay_period` | [uint64](#uint64) |  |  |
| `counterparty_versions` | [Version](#ibc.core.connection.v1.Version) | repeated |  |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `proof_init` | [bytes](#bytes) |  | proof of the initialization the connection on Chain A: `UNITIALIZED -> INIT` |
| `proof_client` | [bytes](#bytes) |  | proof of client state included in message |
| `proof_consensus` | [bytes](#bytes) |  | proof of client consensus state |
| `consensus_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.connection.v1.MsgConnectionOpenTryResponse"></a>

### MsgConnectionOpenTryResponse
MsgConnectionOpenTryResponse defines the Msg/ConnectionOpenTry response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ibc.core.connection.v1.Msg"></a>

### Msg
Msg defines the ibc/connection Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `ConnectionOpenInit` | [MsgConnectionOpenInit](#ibc.core.connection.v1.MsgConnectionOpenInit) | [MsgConnectionOpenInitResponse](#ibc.core.connection.v1.MsgConnectionOpenInitResponse) | ConnectionOpenInit defines a rpc handler method for MsgConnectionOpenInit. | |
| `ConnectionOpenTry` | [MsgConnectionOpenTry](#ibc.core.connection.v1.MsgConnectionOpenTry) | [MsgConnectionOpenTryResponse](#ibc.core.connection.v1.MsgConnectionOpenTryResponse) | ConnectionOpenTry defines a rpc handler method for MsgConnectionOpenTry. | |
| `ConnectionOpenAck` | [MsgConnectionOpenAck](#ibc.core.connection.v1.MsgConnectionOpenAck) | [MsgConnectionOpenAckResponse](#ibc.core.connection.v1.MsgConnectionOpenAckResponse) | ConnectionOpenAck defines a rpc handler method for MsgConnectionOpenAck. | |
| `ConnectionOpenConfirm` | [MsgConnectionOpenConfirm](#ibc.core.connection.v1.MsgConnectionOpenConfirm) | [MsgConnectionOpenConfirmResponse](#ibc.core.connection.v1.MsgConnectionOpenConfirmResponse) | ConnectionOpenConfirm defines a rpc handler method for MsgConnectionOpenConfirm. | |

 <!-- end services -->



<a name="ibc/core/types/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/types/v1/genesis.proto



<a name="ibc.core.types.v1.GenesisState"></a>

### GenesisState
GenesisState defines the ibc module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_genesis` | [ibc.core.client.v1.GenesisState](#ibc.core.client.v1.GenesisState) |  | ICS002 - Clients genesis state |
| `connection_genesis` | [ibc.core.connection.v1.GenesisState](#ibc.core.connection.v1.GenesisState) |  | ICS003 - Connections genesis state |
| `channel_genesis` | [ibc.core.channel.v1.GenesisState](#ibc.core.channel.v1.GenesisState) |  | ICS004 - Channel genesis state |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/lightclients/localhost/v1/localhost.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/lightclients/localhost/v1/localhost.proto



<a name="ibc.lightclients.localhost.v1.ClientState"></a>

### ClientState
ClientState defines a loopback (localhost) client. It requires (read-only)
access to keys outside the client prefix.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [string](#string) |  | self chain ID |
| `height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | self latest block height |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/lightclients/ostracon/v1/ostracon.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/lightclients/ostracon/v1/ostracon.proto



<a name="ibc.lightclients.ostracon.v1.ClientState"></a>

### ClientState
ClientState from Ostracon tracks the current validator set, latest height,
and a possible frozen height.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [string](#string) |  |  |
| `trust_level` | [Fraction](#ibc.lightclients.ostracon.v1.Fraction) |  |  |
| `trusting_period` | [google.protobuf.Duration](#google.protobuf.Duration) |  | duration of the period since the LastestTimestamp during which the submitted headers are valid for upgrade |
| `unbonding_period` | [google.protobuf.Duration](#google.protobuf.Duration) |  | duration of the staking unbonding period |
| `max_clock_drift` | [google.protobuf.Duration](#google.protobuf.Duration) |  | defines how much new (untrusted) header's Time can drift into the future. |
| `frozen_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | Block height when the client was frozen due to a misbehaviour |
| `latest_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | Latest height the client was updated to |
| `proof_specs` | [ics23.ProofSpec](#ics23.ProofSpec) | repeated | Proof specifications used in verifying counterparty state |
| `upgrade_path` | [string](#string) | repeated | Path at which next upgraded client will be committed. Each element corresponds to the key for a single CommitmentProof in the chained proof. NOTE: ClientState must stored under `{upgradePath}/{upgradeHeight}/clientState` ConsensusState must be stored under `{upgradepath}/{upgradeHeight}/consensusState` For SDK chains using the default upgrade module, upgrade_path should be []string{"upgrade", "upgradedIBCState"}` |
| `allow_update_after_expiry` | [bool](#bool) |  | This flag, when set to true, will allow governance to recover a client which has expired |
| `allow_update_after_misbehaviour` | [bool](#bool) |  | This flag, when set to true, will allow governance to unfreeze a client whose chain has experienced a misbehaviour event |






<a name="ibc.lightclients.ostracon.v1.ConsensusState"></a>

### ConsensusState
ConsensusState defines the consensus state from Ostracon.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `timestamp` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | timestamp that corresponds to the block height in which the ConsensusState was stored. |
| `root` | [ibc.core.commitment.v1.MerkleRoot](#ibc.core.commitment.v1.MerkleRoot) |  | commitment root (i.e app hash) |
| `next_validators_hash` | [bytes](#bytes) |  |  |






<a name="ibc.lightclients.ostracon.v1.Fraction"></a>

### Fraction
Fraction defines the protobuf message type for tmmath.Fraction that only supports positive values.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `numerator` | [uint64](#uint64) |  |  |
| `denominator` | [uint64](#uint64) |  |  |






<a name="ibc.lightclients.ostracon.v1.Header"></a>

### Header
Header defines the Ostracon client consensus Header.
It encapsulates all the information necessary to update from a trusted
Ostracon ConsensusState. The inclusion of TrustedHeight and
TrustedValidators allows this update to process correctly, so long as the
ConsensusState for the TrustedHeight exists, this removes race conditions
among relayers The SignedHeader and ValidatorSet are the new untrusted update
fields for the client. The TrustedHeight is the height of a stored
ConsensusState on the client that will be used to verify the new untrusted
header. The Trusted ConsensusState must be within the unbonding period of
current time in order to correctly verify, and the TrustedValidators must
hash to TrustedConsensusState.NextValidatorsHash since that is the last
trusted validator set at the TrustedHeight.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signed_header` | [ostracon.types.SignedHeader](#ostracon.types.SignedHeader) |  |  |
| `validator_set` | [ostracon.types.ValidatorSet](#ostracon.types.ValidatorSet) |  |  |
| `voter_set` | [ostracon.types.VoterSet](#ostracon.types.VoterSet) |  |  |
| `trusted_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `trusted_validators` | [ostracon.types.ValidatorSet](#ostracon.types.ValidatorSet) |  |  |
| `trusted_voters` | [ostracon.types.VoterSet](#ostracon.types.VoterSet) |  |  |






<a name="ibc.lightclients.ostracon.v1.Misbehaviour"></a>

### Misbehaviour
Misbehaviour is a wrapper over two conflicting Headers
that implements Misbehaviour interface expected by ICS-02


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  |  |
| `header_1` | [Header](#ibc.lightclients.ostracon.v1.Header) |  |  |
| `header_2` | [Header](#ibc.lightclients.ostracon.v1.Header) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/lightclients/solomachine/v1/solomachine.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/lightclients/solomachine/v1/solomachine.proto



<a name="ibc.lightclients.solomachine.v1.ChannelStateData"></a>

### ChannelStateData
ChannelStateData returns the SignBytes data for channel state
verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |
| `channel` | [ibc.core.channel.v1.Channel](#ibc.core.channel.v1.Channel) |  |  |






<a name="ibc.lightclients.solomachine.v1.ClientState"></a>

### ClientState
ClientState defines a solo machine client that tracks the current consensus
state and if the client is frozen.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sequence` | [uint64](#uint64) |  | latest sequence of the client state |
| `frozen_sequence` | [uint64](#uint64) |  | frozen sequence of the solo machine |
| `consensus_state` | [ConsensusState](#ibc.lightclients.solomachine.v1.ConsensusState) |  |  |
| `allow_update_after_proposal` | [bool](#bool) |  | when set to true, will allow governance to update a solo machine client. The client will be unfrozen if it is frozen. |






<a name="ibc.lightclients.solomachine.v1.ClientStateData"></a>

### ClientStateData
ClientStateData returns the SignBytes data for client state verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |
| `client_state` | [google.protobuf.Any](#google.protobuf.Any) |  |  |






<a name="ibc.lightclients.solomachine.v1.ConnectionStateData"></a>

### ConnectionStateData
ConnectionStateData returns the SignBytes data for connection state
verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |
| `connection` | [ibc.core.connection.v1.ConnectionEnd](#ibc.core.connection.v1.ConnectionEnd) |  |  |






<a name="ibc.lightclients.solomachine.v1.ConsensusState"></a>

### ConsensusState
ConsensusState defines a solo machine consensus state. The sequence of a consensus state
is contained in the "height" key used in storing the consensus state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `public_key` | [google.protobuf.Any](#google.protobuf.Any) |  | public key of the solo machine |
| `diversifier` | [string](#string) |  | diversifier allows the same public key to be re-used across different solo machine clients (potentially on different chains) without being considered misbehaviour. |
| `timestamp` | [uint64](#uint64) |  |  |






<a name="ibc.lightclients.solomachine.v1.ConsensusStateData"></a>

### ConsensusStateData
ConsensusStateData returns the SignBytes data for consensus state
verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |
| `consensus_state` | [google.protobuf.Any](#google.protobuf.Any) |  |  |






<a name="ibc.lightclients.solomachine.v1.Header"></a>

### Header
Header defines a solo machine consensus header


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sequence` | [uint64](#uint64) |  | sequence to update solo machine public key at |
| `timestamp` | [uint64](#uint64) |  |  |
| `signature` | [bytes](#bytes) |  |  |
| `new_public_key` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `new_diversifier` | [string](#string) |  |  |






<a name="ibc.lightclients.solomachine.v1.HeaderData"></a>

### HeaderData
HeaderData returns the SignBytes data for update verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `new_pub_key` | [google.protobuf.Any](#google.protobuf.Any) |  | header public key |
| `new_diversifier` | [string](#string) |  | header diversifier |






<a name="ibc.lightclients.solomachine.v1.Misbehaviour"></a>

### Misbehaviour
Misbehaviour defines misbehaviour for a solo machine which consists
of a sequence and two signatures over different messages at that sequence.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  |  |
| `sequence` | [uint64](#uint64) |  |  |
| `signature_one` | [SignatureAndData](#ibc.lightclients.solomachine.v1.SignatureAndData) |  |  |
| `signature_two` | [SignatureAndData](#ibc.lightclients.solomachine.v1.SignatureAndData) |  |  |






<a name="ibc.lightclients.solomachine.v1.NextSequenceRecvData"></a>

### NextSequenceRecvData
NextSequenceRecvData returns the SignBytes data for verification of the next
sequence to be received.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |
| `next_seq_recv` | [uint64](#uint64) |  |  |






<a name="ibc.lightclients.solomachine.v1.PacketAcknowledgementData"></a>

### PacketAcknowledgementData
PacketAcknowledgementData returns the SignBytes data for acknowledgement
verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |
| `acknowledgement` | [bytes](#bytes) |  |  |






<a name="ibc.lightclients.solomachine.v1.PacketCommitmentData"></a>

### PacketCommitmentData
PacketCommitmentData returns the SignBytes data for packet commitment
verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |
| `commitment` | [bytes](#bytes) |  |  |






<a name="ibc.lightclients.solomachine.v1.PacketReceiptAbsenceData"></a>

### PacketReceiptAbsenceData
PacketReceiptAbsenceData returns the SignBytes data for
packet receipt absence verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |






<a name="ibc.lightclients.solomachine.v1.SignBytes"></a>

### SignBytes
SignBytes defines the signed bytes used for signature verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sequence` | [uint64](#uint64) |  |  |
| `timestamp` | [uint64](#uint64) |  |  |
| `diversifier` | [string](#string) |  |  |
| `data_type` | [DataType](#ibc.lightclients.solomachine.v1.DataType) |  | type of the data used |
| `data` | [bytes](#bytes) |  | marshaled data |






<a name="ibc.lightclients.solomachine.v1.SignatureAndData"></a>

### SignatureAndData
SignatureAndData contains a signature and the data signed over to create that
signature.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signature` | [bytes](#bytes) |  |  |
| `data_type` | [DataType](#ibc.lightclients.solomachine.v1.DataType) |  |  |
| `data` | [bytes](#bytes) |  |  |
| `timestamp` | [uint64](#uint64) |  |  |






<a name="ibc.lightclients.solomachine.v1.TimestampedSignatureData"></a>

### TimestampedSignatureData
TimestampedSignatureData contains the signature data and the timestamp of the
signature.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signature_data` | [bytes](#bytes) |  |  |
| `timestamp` | [uint64](#uint64) |  |  |





 <!-- end messages -->


<a name="ibc.lightclients.solomachine.v1.DataType"></a>

### DataType
DataType defines the type of solo machine proof being created. This is done to preserve uniqueness of different
data sign byte encodings.

| Name | Number | Description |
| ---- | ------ | ----------- |
| DATA_TYPE_UNINITIALIZED_UNSPECIFIED | 0 | Default State |
| DATA_TYPE_CLIENT_STATE | 1 | Data type for client state verification |
| DATA_TYPE_CONSENSUS_STATE | 2 | Data type for consensus state verification |
| DATA_TYPE_CONNECTION_STATE | 3 | Data type for connection state verification |
| DATA_TYPE_CHANNEL_STATE | 4 | Data type for channel state verification |
| DATA_TYPE_PACKET_COMMITMENT | 5 | Data type for packet commitment verification |
| DATA_TYPE_PACKET_ACKNOWLEDGEMENT | 6 | Data type for packet acknowledgement verification |
| DATA_TYPE_PACKET_RECEIPT_ABSENCE | 7 | Data type for packet receipt absence verification |
| DATA_TYPE_NEXT_SEQUENCE_RECV | 8 | Data type for next sequence recv verification |
| DATA_TYPE_HEADER | 9 | Data type for header verification |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/crypto/ed25519/keys.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/crypto/ed25519/keys.proto



<a name="lbm.crypto.ed25519.PrivKey"></a>

### PrivKey
PrivKey defines a ed25519 private key.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  |  |






<a name="lbm.crypto.ed25519.PubKey"></a>

### PubKey
PubKey defines a ed25519 public key
Key is the compressed form of the pubkey. The first byte depends is a 0x02 byte
if the y-coordinate is the lexicographically largest of the two associated with
the x-coordinate. Otherwise the first byte is a 0x03.
This prefix is followed with the x-coordinate.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/crypto/multisig/keys.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/crypto/multisig/keys.proto



<a name="lbm.crypto.multisig.LegacyAminoPubKey"></a>

### LegacyAminoPubKey
LegacyAminoPubKey specifies a public key type
which nests multiple public keys and a threshold,
it uses legacy amino address rules.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `threshold` | [uint32](#uint32) |  |  |
| `public_keys` | [google.protobuf.Any](#google.protobuf.Any) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/crypto/secp256k1/keys.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/crypto/secp256k1/keys.proto



<a name="lbm.crypto.secp256k1.PrivKey"></a>

### PrivKey
PrivKey defines a secp256k1 private key.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  |  |






<a name="lbm.crypto.secp256k1.PubKey"></a>

### PubKey
PubKey defines a secp256k1 public key
Key is the compressed form of the pubkey. The first byte depends is a 0x02 byte
if the y-coordinate is the lexicographically largest of the two associated with
the x-coordinate. Otherwise the first byte is a 0x03.
This prefix is followed with the x-coordinate.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/auth/v1/auth.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/auth/v1/auth.proto



<a name="lbm.auth.v1.BaseAccount"></a>

### BaseAccount
BaseAccount defines a base account type. It contains all the necessary fields
for basic account functionality. Any custom account type should extend this
type for additional functionality (e.g. vesting).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |
| `ed25519_pub_key` | [lbm.crypto.ed25519.PubKey](#lbm.crypto.ed25519.PubKey) |  |  |
| `secp256k1_pub_key` | [lbm.crypto.secp256k1.PubKey](#lbm.crypto.secp256k1.PubKey) |  |  |
| `multisig_pub_key` | [lbm.crypto.multisig.LegacyAminoPubKey](#lbm.crypto.multisig.LegacyAminoPubKey) |  |  |
| `account_number` | [uint64](#uint64) |  |  |
| `sequence` | [uint64](#uint64) |  |  |






<a name="lbm.auth.v1.ModuleAccount"></a>

### ModuleAccount
ModuleAccount defines an account for modules that holds coins on a pool.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_account` | [BaseAccount](#lbm.auth.v1.BaseAccount) |  |  |
| `name` | [string](#string) |  |  |
| `permissions` | [string](#string) | repeated |  |






<a name="lbm.auth.v1.Params"></a>

### Params
Params defines the parameters for the auth module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `max_memo_characters` | [uint64](#uint64) |  |  |
| `tx_sig_limit` | [uint64](#uint64) |  |  |
| `tx_size_cost_per_byte` | [uint64](#uint64) |  |  |
| `sig_verify_cost_ed25519` | [uint64](#uint64) |  |  |
| `sig_verify_cost_secp256k1` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/auth/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/auth/v1/genesis.proto



<a name="lbm.auth.v1.GenesisState"></a>

### GenesisState
GenesisState defines the auth module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#lbm.auth.v1.Params) |  | params defines all the paramaters of the module. |
| `accounts` | [google.protobuf.Any](#google.protobuf.Any) | repeated | accounts are the accounts present at genesis. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/auth/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/auth/v1/query.proto



<a name="lbm.auth.v1.QueryAccountRequest"></a>

### QueryAccountRequest
QueryAccountRequest is the request type for the Query/Account RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address defines the address to query for. |






<a name="lbm.auth.v1.QueryAccountResponse"></a>

### QueryAccountResponse
QueryAccountResponse is the response type for the Query/Account RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account` | [google.protobuf.Any](#google.protobuf.Any) |  | account defines the account of the corresponding address. |






<a name="lbm.auth.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="lbm.auth.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#lbm.auth.v1.Params) |  | params defines the parameters of the module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.auth.v1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Account` | [QueryAccountRequest](#lbm.auth.v1.QueryAccountRequest) | [QueryAccountResponse](#lbm.auth.v1.QueryAccountResponse) | Account returns account details based on address. | GET|/lbm/auth/v1/accounts/{address}|
| `Params` | [QueryParamsRequest](#lbm.auth.v1.QueryParamsRequest) | [QueryParamsResponse](#lbm.auth.v1.QueryParamsResponse) | Params queries all parameters. | GET|/lbm/auth/v1/params|

 <!-- end services -->



<a name="lbm/auth/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/auth/v1/tx.proto



<a name="lbm.auth.v1.MsgEmpty"></a>

### MsgEmpty
MsgEmpty represents a message that doesn't do anything. Used to measure performance.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `from_address` | [string](#string) |  |  |






<a name="lbm.auth.v1.MsgEmptyResponse"></a>

### MsgEmptyResponse
MsgEmptyResponse defines the Msg/Empty response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.auth.v1.Msg"></a>

### Msg
Msg defines the auth Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Empty` | [MsgEmpty](#lbm.auth.v1.MsgEmpty) | [MsgEmptyResponse](#lbm.auth.v1.MsgEmptyResponse) | Empty defines a method that doesn't do anything. Used to measure performance. | |

 <!-- end services -->



<a name="lbm/bank/v1/bank.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/bank/v1/bank.proto



<a name="lbm.bank.v1.DenomUnit"></a>

### DenomUnit
DenomUnit represents a struct that describes a given
denomination unit of the basic token.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  | denom represents the string name of the given denom unit (e.g uatom). |
| `exponent` | [uint32](#uint32) |  | exponent represents power of 10 exponent that one must raise the base_denom to in order to equal the given DenomUnit's denom 1 denom = 1^exponent base_denom (e.g. with a base_denom of uatom, one can create a DenomUnit of 'atom' with exponent = 6, thus: 1 atom = 10^6 uatom). |
| `aliases` | [string](#string) | repeated | aliases is a list of string aliases for the given denom |






<a name="lbm.bank.v1.Input"></a>

### Input
Input models transaction input.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |
| `coins` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated |  |






<a name="lbm.bank.v1.Metadata"></a>

### Metadata
Metadata represents a struct that describes
a basic token.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `description` | [string](#string) |  |  |
| `denom_units` | [DenomUnit](#lbm.bank.v1.DenomUnit) | repeated | denom_units represents the list of DenomUnit's for a given coin |
| `base` | [string](#string) |  | base represents the base denom (should be the DenomUnit with exponent = 0). |
| `display` | [string](#string) |  | display indicates the suggested denom that should be displayed in clients. |






<a name="lbm.bank.v1.Output"></a>

### Output
Output models transaction outputs.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |
| `coins` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated |  |






<a name="lbm.bank.v1.Params"></a>

### Params
Params defines the parameters for the bank module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `send_enabled` | [SendEnabled](#lbm.bank.v1.SendEnabled) | repeated |  |
| `default_send_enabled` | [bool](#bool) |  |  |






<a name="lbm.bank.v1.SendEnabled"></a>

### SendEnabled
SendEnabled maps coin denom to a send_enabled status (whether a denom is
sendable).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `enabled` | [bool](#bool) |  |  |






<a name="lbm.bank.v1.Supply"></a>

### Supply
Supply represents a struct that passively keeps track of the total supply
amounts in the network.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `total` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/bank/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/bank/v1/genesis.proto



<a name="lbm.bank.v1.Balance"></a>

### Balance
Balance defines an account address and balance pair used in the bank module's
genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address of the balance holder. |
| `coins` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated | coins defines the different coins this balance holds. |






<a name="lbm.bank.v1.GenesisState"></a>

### GenesisState
GenesisState defines the bank module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#lbm.bank.v1.Params) |  | params defines all the paramaters of the module. |
| `balances` | [Balance](#lbm.bank.v1.Balance) | repeated | balances is an array containing the balances of all the accounts. |
| `supply` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated | supply represents the total supply. |
| `denom_metadata` | [Metadata](#lbm.bank.v1.Metadata) | repeated | denom_metadata defines the metadata of the differents coins. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/bank/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/bank/v1/query.proto



<a name="lbm.bank.v1.QueryAllBalancesRequest"></a>

### QueryAllBalancesRequest
QueryBalanceRequest is the request type for the Query/AllBalances RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address to query balances for. |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.bank.v1.QueryAllBalancesResponse"></a>

### QueryAllBalancesResponse
QueryAllBalancesResponse is the response type for the Query/AllBalances RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `balances` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated | balances is the balances of all the coins. |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.bank.v1.QueryBalanceRequest"></a>

### QueryBalanceRequest
QueryBalanceRequest is the request type for the Query/Balance RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address to query balances for. |
| `denom` | [string](#string) |  | denom is the coin denom to query balances for. |






<a name="lbm.bank.v1.QueryBalanceResponse"></a>

### QueryBalanceResponse
QueryBalanceResponse is the response type for the Query/Balance RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `balance` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) |  | balance is the balance of the coin. |






<a name="lbm.bank.v1.QueryDenomMetadataRequest"></a>

### QueryDenomMetadataRequest
QueryDenomMetadataRequest is the request type for the Query/DenomMetadata RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  | denom is the coin denom to query the metadata for. |






<a name="lbm.bank.v1.QueryDenomMetadataResponse"></a>

### QueryDenomMetadataResponse
QueryDenomMetadataResponse is the response type for the Query/DenomMetadata RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `metadata` | [Metadata](#lbm.bank.v1.Metadata) |  | metadata describes and provides all the client information for the requested token. |






<a name="lbm.bank.v1.QueryDenomsMetadataRequest"></a>

### QueryDenomsMetadataRequest
QueryDenomsMetadataRequest is the request type for the Query/DenomsMetadata RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.bank.v1.QueryDenomsMetadataResponse"></a>

### QueryDenomsMetadataResponse
QueryDenomsMetadataResponse is the response type for the Query/DenomsMetadata RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `metadatas` | [Metadata](#lbm.bank.v1.Metadata) | repeated | metadata provides the client information for all the registered tokens. |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.bank.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest defines the request type for querying x/bank parameters.






<a name="lbm.bank.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse defines the response type for querying x/bank parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#lbm.bank.v1.Params) |  |  |






<a name="lbm.bank.v1.QuerySupplyOfRequest"></a>

### QuerySupplyOfRequest
QuerySupplyOfRequest is the request type for the Query/SupplyOf RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  | denom is the coin denom to query balances for. |






<a name="lbm.bank.v1.QuerySupplyOfResponse"></a>

### QuerySupplyOfResponse
QuerySupplyOfResponse is the response type for the Query/SupplyOf RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `amount` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) |  | amount is the supply of the coin. |






<a name="lbm.bank.v1.QueryTotalSupplyRequest"></a>

### QueryTotalSupplyRequest
QueryTotalSupplyRequest is the request type for the Query/TotalSupply RPC
method.






<a name="lbm.bank.v1.QueryTotalSupplyResponse"></a>

### QueryTotalSupplyResponse
QueryTotalSupplyResponse is the response type for the Query/TotalSupply RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `supply` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated | supply is the supply of the coins |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.bank.v1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Balance` | [QueryBalanceRequest](#lbm.bank.v1.QueryBalanceRequest) | [QueryBalanceResponse](#lbm.bank.v1.QueryBalanceResponse) | Balance queries the balance of a single coin for a single account. | GET|/lbm/bank/v1/balances/{address}/by_denom|
| `AllBalances` | [QueryAllBalancesRequest](#lbm.bank.v1.QueryAllBalancesRequest) | [QueryAllBalancesResponse](#lbm.bank.v1.QueryAllBalancesResponse) | AllBalances queries the balance of all coins for a single account. | GET|/lbm/bank/v1/balances/{address}|
| `TotalSupply` | [QueryTotalSupplyRequest](#lbm.bank.v1.QueryTotalSupplyRequest) | [QueryTotalSupplyResponse](#lbm.bank.v1.QueryTotalSupplyResponse) | TotalSupply queries the total supply of all coins. | GET|/lbm/bank/v1/supply|
| `SupplyOf` | [QuerySupplyOfRequest](#lbm.bank.v1.QuerySupplyOfRequest) | [QuerySupplyOfResponse](#lbm.bank.v1.QuerySupplyOfResponse) | SupplyOf queries the supply of a single coin. | GET|/lbm/bank/v1/supply/{denom}|
| `Params` | [QueryParamsRequest](#lbm.bank.v1.QueryParamsRequest) | [QueryParamsResponse](#lbm.bank.v1.QueryParamsResponse) | Params queries the parameters of x/bank module. | GET|/lbm/bank/v1/params|
| `DenomMetadata` | [QueryDenomMetadataRequest](#lbm.bank.v1.QueryDenomMetadataRequest) | [QueryDenomMetadataResponse](#lbm.bank.v1.QueryDenomMetadataResponse) | DenomsMetadata queries the client metadata of a given coin denomination. | GET|/lbm/bank/v1/denoms_metadata/{denom}|
| `DenomsMetadata` | [QueryDenomsMetadataRequest](#lbm.bank.v1.QueryDenomsMetadataRequest) | [QueryDenomsMetadataResponse](#lbm.bank.v1.QueryDenomsMetadataResponse) | DenomsMetadata queries the client metadata for all registered coin denominations. | GET|/lbm/bank/v1/denoms_metadata|

 <!-- end services -->



<a name="lbm/bank/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/bank/v1/tx.proto



<a name="lbm.bank.v1.MsgMultiSend"></a>

### MsgMultiSend
MsgMultiSend represents an arbitrary multi-in, multi-out send message.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `inputs` | [Input](#lbm.bank.v1.Input) | repeated |  |
| `outputs` | [Output](#lbm.bank.v1.Output) | repeated |  |






<a name="lbm.bank.v1.MsgMultiSendResponse"></a>

### MsgMultiSendResponse
MsgMultiSendResponse defines the Msg/MultiSend response type.






<a name="lbm.bank.v1.MsgSend"></a>

### MsgSend
MsgSend represents a message to send coins from one account to another.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `from_address` | [string](#string) |  |  |
| `to_address` | [string](#string) |  |  |
| `amount` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated |  |






<a name="lbm.bank.v1.MsgSendResponse"></a>

### MsgSendResponse
MsgSendResponse defines the Msg/Send response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.bank.v1.Msg"></a>

### Msg
Msg defines the bank Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Send` | [MsgSend](#lbm.bank.v1.MsgSend) | [MsgSendResponse](#lbm.bank.v1.MsgSendResponse) | Send defines a method for sending coins from one account to another account. | |
| `MultiSend` | [MsgMultiSend](#lbm.bank.v1.MsgMultiSend) | [MsgMultiSendResponse](#lbm.bank.v1.MsgMultiSendResponse) | MultiSend defines a method for sending coins from some accounts to other accounts. | |

 <!-- end services -->



<a name="lbm/bankplus/v1/bankplus.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/bankplus/v1/bankplus.proto



<a name="lbm.bankplus.v1.InactiveAddr"></a>

### InactiveAddr
InactiveAddr models the blocked address for the bankplus module


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/base/abci/v1/abci.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/base/abci/v1/abci.proto



<a name="lbm.base.abci.v1.ABCIMessageLog"></a>

### ABCIMessageLog
ABCIMessageLog defines a structure containing an indexed tx ABCI message log.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `msg_index` | [uint32](#uint32) |  |  |
| `log` | [string](#string) |  |  |
| `events` | [StringEvent](#lbm.base.abci.v1.StringEvent) | repeated | Events contains a slice of Event objects that were emitted during some execution. |






<a name="lbm.base.abci.v1.Attribute"></a>

### Attribute
Attribute defines an attribute wrapper where the key and value are
strings instead of raw bytes.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |






<a name="lbm.base.abci.v1.GasInfo"></a>

### GasInfo
GasInfo defines tx execution gas context.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gas_wanted` | [uint64](#uint64) |  | GasWanted is the maximum units of work we allow this tx to perform. |
| `gas_used` | [uint64](#uint64) |  | GasUsed is the amount of gas actually consumed. |






<a name="lbm.base.abci.v1.MsgData"></a>

### MsgData
MsgData defines the data returned in a Result object during message
execution.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `msg_type` | [string](#string) |  |  |
| `data` | [bytes](#bytes) |  |  |






<a name="lbm.base.abci.v1.Result"></a>

### Result
Result is the union of ResponseFormat and ResponseCheckTx.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  | Data is any data returned from message or handler execution. It MUST be length prefixed in order to separate data from multiple message executions. |
| `log` | [string](#string) |  | Log contains the log information from message or handler execution. |
| `events` | [ostracon.abci.Event](#ostracon.abci.Event) | repeated | Events contains a slice of Event objects that were emitted during message or handler execution. |






<a name="lbm.base.abci.v1.SearchTxsResult"></a>

### SearchTxsResult
SearchTxsResult defines a structure for querying txs pageable


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `total_count` | [uint64](#uint64) |  | Count of all txs |
| `count` | [uint64](#uint64) |  | Count of txs in current page |
| `page_number` | [uint64](#uint64) |  | Index of current page, start from 1 |
| `page_total` | [uint64](#uint64) |  | Count of total pages |
| `limit` | [uint64](#uint64) |  | Max count txs per page |
| `txs` | [TxResponse](#lbm.base.abci.v1.TxResponse) | repeated | List of txs in current page |






<a name="lbm.base.abci.v1.SimulationResponse"></a>

### SimulationResponse
SimulationResponse defines the response generated when a transaction is
successfully simulated.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gas_info` | [GasInfo](#lbm.base.abci.v1.GasInfo) |  |  |
| `result` | [Result](#lbm.base.abci.v1.Result) |  |  |






<a name="lbm.base.abci.v1.StringEvent"></a>

### StringEvent
StringEvent defines en Event object wrapper where all the attributes
contain key/value pairs that are strings instead of raw bytes.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [string](#string) |  |  |
| `attributes` | [Attribute](#lbm.base.abci.v1.Attribute) | repeated |  |






<a name="lbm.base.abci.v1.TxMsgData"></a>

### TxMsgData
TxMsgData defines a list of MsgData. A transaction will have a MsgData object
for each message.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [MsgData](#lbm.base.abci.v1.MsgData) | repeated |  |






<a name="lbm.base.abci.v1.TxResponse"></a>

### TxResponse
TxResponse defines a structure containing relevant tx data and metadata. The
tags are stringified and the log is JSON decoded.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  | The block height |
| `txhash` | [string](#string) |  | The transaction hash. |
| `codespace` | [string](#string) |  | Namespace for the Code |
| `code` | [uint32](#uint32) |  | Response code. |
| `data` | [string](#string) |  | Result bytes, if any. |
| `raw_log` | [string](#string) |  | The output of the application's logger (raw string). May be non-deterministic. |
| `logs` | [ABCIMessageLog](#lbm.base.abci.v1.ABCIMessageLog) | repeated | The output of the application's logger (typed). May be non-deterministic. |
| `info` | [string](#string) |  | Additional information. May be non-deterministic. |
| `gas_wanted` | [int64](#int64) |  | Amount of gas requested for transaction. |
| `gas_used` | [int64](#int64) |  | Amount of gas consumed by transaction. |
| `tx` | [google.protobuf.Any](#google.protobuf.Any) |  | The request transaction bytes. |
| `timestamp` | [string](#string) |  | Time of the previous block. For heights > 1, it's the weighted median of the timestamps of the valid votes in the block.LastCommit. For height == 1, it's genesis time. |
| `events` | [ostracon.abci.Event](#ostracon.abci.Event) | repeated | Events defines all the events emitted by processing a transaction. Note, these events include those emitted by processing all the messages and those emitted from the ante handler. Whereas Logs contains the events, with additional metadata, emitted only by processing the messages.

Since: cosmos-sdk 0.42.11, 0.44.5, 0.45 |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/base/kv/v1/kv.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/base/kv/v1/kv.proto



<a name="lbm.base.kv.v1.Pair"></a>

### Pair
Pair defines a key/value bytes tuple.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  |  |
| `value` | [bytes](#bytes) |  |  |






<a name="lbm.base.kv.v1.Pairs"></a>

### Pairs
Pairs defines a repeated slice of Pair objects.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pairs` | [Pair](#lbm.base.kv.v1.Pair) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/base/ostracon/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/base/ostracon/v1/query.proto



<a name="lbm.base.ostracon.v1.GetBlockByHashRequest"></a>

### GetBlockByHashRequest
GetBlockByHashRequest is the request type for the Query/GetBlockByHash RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [bytes](#bytes) |  |  |






<a name="lbm.base.ostracon.v1.GetBlockByHashResponse"></a>

### GetBlockByHashResponse
GetBlockByHashResponse is the response type for the Query/GetBlockByHash RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_id` | [ostracon.types.BlockID](#ostracon.types.BlockID) |  |  |
| `block` | [ostracon.types.Block](#ostracon.types.Block) |  |  |






<a name="lbm.base.ostracon.v1.GetBlockByHeightRequest"></a>

### GetBlockByHeightRequest
GetBlockByHeightRequest is the request type for the Query/GetBlockByHeight RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  |  |






<a name="lbm.base.ostracon.v1.GetBlockByHeightResponse"></a>

### GetBlockByHeightResponse
GetBlockByHeightResponse is the response type for the Query/GetBlockByHeight RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_id` | [ostracon.types.BlockID](#ostracon.types.BlockID) |  |  |
| `block` | [ostracon.types.Block](#ostracon.types.Block) |  |  |






<a name="lbm.base.ostracon.v1.GetBlockResultsByHeightRequest"></a>

### GetBlockResultsByHeightRequest
GetBlockResultsByHeightRequest is the request type for the Query/GetBlockResultsByHeight RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  |  |






<a name="lbm.base.ostracon.v1.GetBlockResultsByHeightResponse"></a>

### GetBlockResultsByHeightResponse
GetBlockResultsByHeightResponse is the response type for the Query/GetBlockResultsByHeight RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  |  |
| `txs_results` | [ostracon.abci.ResponseDeliverTx](#ostracon.abci.ResponseDeliverTx) | repeated |  |
| `res_begin_block` | [ostracon.abci.ResponseBeginBlock](#ostracon.abci.ResponseBeginBlock) |  |  |
| `res_end_block` | [ostracon.abci.ResponseEndBlock](#ostracon.abci.ResponseEndBlock) |  |  |






<a name="lbm.base.ostracon.v1.GetLatestBlockRequest"></a>

### GetLatestBlockRequest
GetLatestBlockRequest is the request type for the Query/GetLatestBlock RPC method.






<a name="lbm.base.ostracon.v1.GetLatestBlockResponse"></a>

### GetLatestBlockResponse
GetLatestBlockResponse is the response type for the Query/GetLatestBlock RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_id` | [ostracon.types.BlockID](#ostracon.types.BlockID) |  |  |
| `block` | [ostracon.types.Block](#ostracon.types.Block) |  |  |






<a name="lbm.base.ostracon.v1.GetLatestValidatorSetRequest"></a>

### GetLatestValidatorSetRequest
GetLatestValidatorSetRequest is the request type for the Query/GetValidatorSetByHeight RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an pagination for the request. |






<a name="lbm.base.ostracon.v1.GetLatestValidatorSetResponse"></a>

### GetLatestValidatorSetResponse
GetLatestValidatorSetResponse is the response type for the Query/GetValidatorSetByHeight RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_height` | [int64](#int64) |  |  |
| `validators` | [Validator](#lbm.base.ostracon.v1.Validator) | repeated |  |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines an pagination for the response. |






<a name="lbm.base.ostracon.v1.GetNodeInfoRequest"></a>

### GetNodeInfoRequest
GetNodeInfoRequest is the request type for the Query/GetNodeInfo RPC method.






<a name="lbm.base.ostracon.v1.GetNodeInfoResponse"></a>

### GetNodeInfoResponse
GetNodeInfoResponse is the request type for the Query/GetNodeInfo RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `default_node_info` | [ostracon.p2p.DefaultNodeInfo](#ostracon.p2p.DefaultNodeInfo) |  |  |
| `application_version` | [VersionInfo](#lbm.base.ostracon.v1.VersionInfo) |  |  |






<a name="lbm.base.ostracon.v1.GetSyncingRequest"></a>

### GetSyncingRequest
GetSyncingRequest is the request type for the Query/GetSyncing RPC method.






<a name="lbm.base.ostracon.v1.GetSyncingResponse"></a>

### GetSyncingResponse
GetSyncingResponse is the response type for the Query/GetSyncing RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `syncing` | [bool](#bool) |  |  |






<a name="lbm.base.ostracon.v1.GetValidatorSetByHeightRequest"></a>

### GetValidatorSetByHeightRequest
GetValidatorSetByHeightRequest is the request type for the Query/GetValidatorSetByHeight RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  |  |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an pagination for the request. |






<a name="lbm.base.ostracon.v1.GetValidatorSetByHeightResponse"></a>

### GetValidatorSetByHeightResponse
GetValidatorSetByHeightResponse is the response type for the Query/GetValidatorSetByHeight RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_height` | [int64](#int64) |  |  |
| `validators` | [Validator](#lbm.base.ostracon.v1.Validator) | repeated |  |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines an pagination for the response. |






<a name="lbm.base.ostracon.v1.Module"></a>

### Module
Module is the type for VersionInfo


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [string](#string) |  | module path |
| `version` | [string](#string) |  | module version |
| `sum` | [string](#string) |  | checksum |






<a name="lbm.base.ostracon.v1.Validator"></a>

### Validator
Validator is the type for the validator-set.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |
| `pub_key` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `voting_power` | [int64](#int64) |  |  |
| `proposer_priority` | [int64](#int64) |  |  |






<a name="lbm.base.ostracon.v1.VersionInfo"></a>

### VersionInfo
VersionInfo is the type for the GetNodeInfoResponse message.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `app_name` | [string](#string) |  |  |
| `version` | [string](#string) |  |  |
| `git_commit` | [string](#string) |  |  |
| `build_tags` | [string](#string) |  |  |
| `go_version` | [string](#string) |  |  |
| `build_deps` | [Module](#lbm.base.ostracon.v1.Module) | repeated |  |
| `lbm_sdk_version` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.base.ostracon.v1.Service"></a>

### Service
Service defines the gRPC querier service for ostracon queries.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `GetNodeInfo` | [GetNodeInfoRequest](#lbm.base.ostracon.v1.GetNodeInfoRequest) | [GetNodeInfoResponse](#lbm.base.ostracon.v1.GetNodeInfoResponse) | GetNodeInfo queries the current node info. | GET|/lbm/base/ostracon/v1/node_info|
| `GetSyncing` | [GetSyncingRequest](#lbm.base.ostracon.v1.GetSyncingRequest) | [GetSyncingResponse](#lbm.base.ostracon.v1.GetSyncingResponse) | GetSyncing queries node syncing. | GET|/lbm/base/ostracon/v1/syncing|
| `GetLatestBlock` | [GetLatestBlockRequest](#lbm.base.ostracon.v1.GetLatestBlockRequest) | [GetLatestBlockResponse](#lbm.base.ostracon.v1.GetLatestBlockResponse) | GetLatestBlock returns the latest block. | GET|/lbm/base/ostracon/v1/blocks/latest|
| `GetBlockByHeight` | [GetBlockByHeightRequest](#lbm.base.ostracon.v1.GetBlockByHeightRequest) | [GetBlockByHeightResponse](#lbm.base.ostracon.v1.GetBlockByHeightResponse) | GetBlockByHeight queries block for given height. | GET|/lbm/base/ostracon/v1/blocks/{height}|
| `GetBlockByHash` | [GetBlockByHashRequest](#lbm.base.ostracon.v1.GetBlockByHashRequest) | [GetBlockByHashResponse](#lbm.base.ostracon.v1.GetBlockByHashResponse) | GetBlockByHash queries block for given hash. | GET|/lbm/base/ostracon/v1/blocks/{hash}|
| `GetBlockResultsByHeight` | [GetBlockResultsByHeightRequest](#lbm.base.ostracon.v1.GetBlockResultsByHeightRequest) | [GetBlockResultsByHeightResponse](#lbm.base.ostracon.v1.GetBlockResultsByHeightResponse) | GetBlockResultsByHeight queries block results for given height. | GET|/lbm/base/ostracon/v1/blockresults/{height}|
| `GetLatestValidatorSet` | [GetLatestValidatorSetRequest](#lbm.base.ostracon.v1.GetLatestValidatorSetRequest) | [GetLatestValidatorSetResponse](#lbm.base.ostracon.v1.GetLatestValidatorSetResponse) | GetLatestValidatorSet queries latest validator-set. | GET|/lbm/base/ostracon/v1/validatorsets/latest|
| `GetValidatorSetByHeight` | [GetValidatorSetByHeightRequest](#lbm.base.ostracon.v1.GetValidatorSetByHeightRequest) | [GetValidatorSetByHeightResponse](#lbm.base.ostracon.v1.GetValidatorSetByHeightResponse) | GetValidatorSetByHeight queries validator-set at a given height. | GET|/lbm/base/ostracon/v1/validatorsets/{height}|

 <!-- end services -->



<a name="lbm/base/reflection/v1/reflection.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/base/reflection/v1/reflection.proto



<a name="lbm.base.reflection.v1.ListAllInterfacesRequest"></a>

### ListAllInterfacesRequest
ListAllInterfacesRequest is the request type of the ListAllInterfaces RPC.






<a name="lbm.base.reflection.v1.ListAllInterfacesResponse"></a>

### ListAllInterfacesResponse
ListAllInterfacesResponse is the response type of the ListAllInterfaces RPC.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `interface_names` | [string](#string) | repeated | interface_names is an array of all the registered interfaces. |






<a name="lbm.base.reflection.v1.ListImplementationsRequest"></a>

### ListImplementationsRequest
ListImplementationsRequest is the request type of the ListImplementations
RPC.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `interface_name` | [string](#string) |  | interface_name defines the interface to query the implementations for. |






<a name="lbm.base.reflection.v1.ListImplementationsResponse"></a>

### ListImplementationsResponse
ListImplementationsResponse is the response type of the ListImplementations
RPC.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `implementation_message_names` | [string](#string) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.base.reflection.v1.ReflectionService"></a>

### ReflectionService
ReflectionService defines a service for interface reflection.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `ListAllInterfaces` | [ListAllInterfacesRequest](#lbm.base.reflection.v1.ListAllInterfacesRequest) | [ListAllInterfacesResponse](#lbm.base.reflection.v1.ListAllInterfacesResponse) | ListAllInterfaces lists all the interfaces registered in the interface registry. | GET|/lbm/base/reflection/v1/interfaces|
| `ListImplementations` | [ListImplementationsRequest](#lbm.base.reflection.v1.ListImplementationsRequest) | [ListImplementationsResponse](#lbm.base.reflection.v1.ListImplementationsResponse) | ListImplementations list all the concrete types that implement a given interface. | GET|/lbm/base/reflection/v1/interfaces/{interface_name}/implementations|

 <!-- end services -->



<a name="lbm/base/snapshots/v1/snapshot.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/base/snapshots/v1/snapshot.proto



<a name="lbm.base.snapshots.v1.Metadata"></a>

### Metadata
Metadata contains SDK-specific snapshot metadata.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chunk_hashes` | [bytes](#bytes) | repeated | SHA-256 chunk hashes |






<a name="lbm.base.snapshots.v1.Snapshot"></a>

### Snapshot
Snapshot contains Tendermint state sync snapshot info.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [uint64](#uint64) |  |  |
| `format` | [uint32](#uint32) |  |  |
| `chunks` | [uint32](#uint32) |  |  |
| `hash` | [bytes](#bytes) |  |  |
| `metadata` | [Metadata](#lbm.base.snapshots.v1.Metadata) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/base/store/v1/commit_info.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/base/store/v1/commit_info.proto



<a name="lbm.base.store.v1.CommitID"></a>

### CommitID
CommitID defines the committment information when a specific store is
committed.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `version` | [int64](#int64) |  |  |
| `hash` | [bytes](#bytes) |  |  |






<a name="lbm.base.store.v1.CommitInfo"></a>

### CommitInfo
CommitInfo defines commit information used by the multi-store when committing
a version/height.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `version` | [int64](#int64) |  |  |
| `store_infos` | [StoreInfo](#lbm.base.store.v1.StoreInfo) | repeated |  |






<a name="lbm.base.store.v1.StoreInfo"></a>

### StoreInfo
StoreInfo defines store-specific commit information. It contains a reference
between a store name and the commit ID.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `commit_id` | [CommitID](#lbm.base.store.v1.CommitID) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/base/store/v1/snapshot.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/base/store/v1/snapshot.proto



<a name="lbm.base.store.v1.SnapshotIAVLItem"></a>

### SnapshotIAVLItem
SnapshotIAVLItem is an exported IAVL node.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  |  |
| `value` | [bytes](#bytes) |  |  |
| `version` | [int64](#int64) |  |  |
| `height` | [int32](#int32) |  |  |






<a name="lbm.base.store.v1.SnapshotItem"></a>

### SnapshotItem
SnapshotItem is an item contained in a rootmulti.Store snapshot.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `store` | [SnapshotStoreItem](#lbm.base.store.v1.SnapshotStoreItem) |  |  |
| `iavl` | [SnapshotIAVLItem](#lbm.base.store.v1.SnapshotIAVLItem) |  |  |






<a name="lbm.base.store.v1.SnapshotStoreItem"></a>

### SnapshotStoreItem
SnapshotStoreItem contains metadata about a snapshotted store.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/capability/v1/capability.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/capability/v1/capability.proto



<a name="lbm.capability.v1.Capability"></a>

### Capability
Capability defines an implementation of an object capability. The index
provided to a Capability must be globally unique.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `index` | [uint64](#uint64) |  |  |






<a name="lbm.capability.v1.CapabilityOwners"></a>

### CapabilityOwners
CapabilityOwners defines a set of owners of a single Capability. The set of
owners must be unique.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `owners` | [Owner](#lbm.capability.v1.Owner) | repeated |  |






<a name="lbm.capability.v1.Owner"></a>

### Owner
Owner defines a single capability owner. An owner is defined by the name of
capability and the module name.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `module` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/capability/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/capability/v1/genesis.proto



<a name="lbm.capability.v1.GenesisOwners"></a>

### GenesisOwners
GenesisOwners defines the capability owners with their corresponding index.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `index` | [uint64](#uint64) |  | index is the index of the capability owner. |
| `index_owners` | [CapabilityOwners](#lbm.capability.v1.CapabilityOwners) |  | index_owners are the owners at the given index. |






<a name="lbm.capability.v1.GenesisState"></a>

### GenesisState
GenesisState defines the capability module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `index` | [uint64](#uint64) |  | index is the capability global index. |
| `owners` | [GenesisOwners](#lbm.capability.v1.GenesisOwners) | repeated | owners represents a map from index to owners of the capability index index key is string to allow amino marshalling. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/consortium/v1/consortium.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/consortium/v1/consortium.proto



<a name="lbm.consortium.v1.Params"></a>

### Params
Params defines the parameters for the consortium module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `enabled` | [bool](#bool) |  |  |






<a name="lbm.consortium.v1.UpdateConsortiumParamsProposal"></a>

### UpdateConsortiumParamsProposal
UpdateConsortiumParamsProposal details a proposal to update params of cosortium module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `params` | [Params](#lbm.consortium.v1.Params) |  |  |






<a name="lbm.consortium.v1.UpdateValidatorAuthsProposal"></a>

### UpdateValidatorAuthsProposal
UpdateValidatorAuthsProposal details a proposal to update validator auths on consortium.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `auths` | [ValidatorAuth](#lbm.consortium.v1.ValidatorAuth) | repeated |  |






<a name="lbm.consortium.v1.ValidatorAuth"></a>

### ValidatorAuth
ValidatorAuth defines authorization info of a validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `operator_address` | [string](#string) |  |  |
| `creation_allowed` | [bool](#bool) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/consortium/v1/event.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/consortium/v1/event.proto



<a name="lbm.consortium.v1.EventUpdateConsortiumParams"></a>

### EventUpdateConsortiumParams
EventUpdateConsortiumParams is emitted after updating consortium parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#lbm.consortium.v1.Params) |  |  |






<a name="lbm.consortium.v1.EventUpdateValidatorAuths"></a>

### EventUpdateValidatorAuths
EventUpdateValidatorAuths is emitted after updating validator auth info.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auths` | [ValidatorAuth](#lbm.consortium.v1.ValidatorAuth) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/consortium/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/consortium/v1/genesis.proto



<a name="lbm.consortium.v1.GenesisState"></a>

### GenesisState
GenesisState defines the consortium module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#lbm.consortium.v1.Params) |  | params defines the module parameters at genesis. |
| `validator_auths` | [ValidatorAuth](#lbm.consortium.v1.ValidatorAuth) | repeated | allowed_validators defines the allowed validator addresses at genesis. provided empty, the module gathers information from staking module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/consortium/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/consortium/v1/query.proto



<a name="lbm.consortium.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="lbm.consortium.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#lbm.consortium.v1.Params) |  |  |






<a name="lbm.consortium.v1.QueryValidatorAuthRequest"></a>

### QueryValidatorAuthRequest
QueryValidatorAuthRequest is the request type for the
Query/ValidatorAuth RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address defines the validator address to query for. |






<a name="lbm.consortium.v1.QueryValidatorAuthResponse"></a>

### QueryValidatorAuthResponse
QueryValidatorAuthResponse is the request type for the
Query/ValidatorAuth RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auth` | [ValidatorAuth](#lbm.consortium.v1.ValidatorAuth) |  |  |






<a name="lbm.consortium.v1.QueryValidatorAuthsRequest"></a>

### QueryValidatorAuthsRequest
QueryValidatorAuthsRequest is the request type for the
Query/ValidatorAuths RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.consortium.v1.QueryValidatorAuthsResponse"></a>

### QueryValidatorAuthsResponse
QueryValidatorAuthsResponse is the response type for the
Query/ValidatorAuths RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auths` | [ValidatorAuth](#lbm.consortium.v1.ValidatorAuth) | repeated |  |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines the pagination in the response. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.consortium.v1.Query"></a>

### Query
Query defines the gRPC querier service for consortium module.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#lbm.consortium.v1.QueryParamsRequest) | [QueryParamsResponse](#lbm.consortium.v1.QueryParamsResponse) | Params queries the module params. | GET|/lbm/consortium/v1/params|
| `ValidatorAuth` | [QueryValidatorAuthRequest](#lbm.consortium.v1.QueryValidatorAuthRequest) | [QueryValidatorAuthResponse](#lbm.consortium.v1.QueryValidatorAuthResponse) | ValidatorAuth queries authorization info of a validator. | GET|/lbm/consortium/v1/validators/{validator_address}|
| `ValidatorAuths` | [QueryValidatorAuthsRequest](#lbm.consortium.v1.QueryValidatorAuthsRequest) | [QueryValidatorAuthsResponse](#lbm.consortium.v1.QueryValidatorAuthsResponse) | ValidatorAuths queries authorization infos of validators. | GET|/lbm/consortium/v1/validators|

 <!-- end services -->



<a name="lbm/crisis/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/crisis/v1/genesis.proto



<a name="lbm.crisis.v1.GenesisState"></a>

### GenesisState
GenesisState defines the crisis module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `constant_fee` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) |  | constant_fee is the fee used to verify the invariant in the crisis module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/crisis/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/crisis/v1/tx.proto



<a name="lbm.crisis.v1.MsgVerifyInvariant"></a>

### MsgVerifyInvariant
MsgVerifyInvariant represents a message to verify a particular invariance.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `invariant_module_name` | [string](#string) |  |  |
| `invariant_route` | [string](#string) |  |  |






<a name="lbm.crisis.v1.MsgVerifyInvariantResponse"></a>

### MsgVerifyInvariantResponse
MsgVerifyInvariantResponse defines the Msg/VerifyInvariant response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.crisis.v1.Msg"></a>

### Msg
Msg defines the bank Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `VerifyInvariant` | [MsgVerifyInvariant](#lbm.crisis.v1.MsgVerifyInvariant) | [MsgVerifyInvariantResponse](#lbm.crisis.v1.MsgVerifyInvariantResponse) | VerifyInvariant defines a method to verify a particular invariance. | |

 <!-- end services -->



<a name="lbm/crypto/multisig/v1/multisig.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/crypto/multisig/v1/multisig.proto



<a name="lbm.crypto.multisig.v1.CompactBitArray"></a>

### CompactBitArray
CompactBitArray is an implementation of a space efficient bit array.
This is used to ensure that the encoded data takes up a minimal amount of
space after proto encoding.
This is not thread safe, and is not intended for concurrent usage.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `extra_bits_stored` | [uint32](#uint32) |  |  |
| `elems` | [bytes](#bytes) |  |  |






<a name="lbm.crypto.multisig.v1.MultiSignature"></a>

### MultiSignature
MultiSignature wraps the signatures from a multisig.LegacyAminoPubKey.
See lbm.tx.v1betata1.ModeInfo.Multi for how to specify which signers
signed and with which modes.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signatures` | [bytes](#bytes) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/distribution/v1/distribution.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/distribution/v1/distribution.proto



<a name="lbm.distribution.v1.CommunityPoolSpendProposal"></a>

### CommunityPoolSpendProposal
CommunityPoolSpendProposal details a proposal for use of community funds,
together with how many coins are proposed to be spent, and to which
recipient account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `recipient` | [string](#string) |  |  |
| `amount` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated |  |






<a name="lbm.distribution.v1.CommunityPoolSpendProposalWithDeposit"></a>

### CommunityPoolSpendProposalWithDeposit
CommunityPoolSpendProposalWithDeposit defines a CommunityPoolSpendProposal
with a deposit


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `recipient` | [string](#string) |  |  |
| `amount` | [string](#string) |  |  |
| `deposit` | [string](#string) |  |  |






<a name="lbm.distribution.v1.DelegationDelegatorReward"></a>

### DelegationDelegatorReward
DelegationDelegatorReward represents the properties
of a delegator's delegation reward.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  |  |
| `reward` | [lbm.base.v1.DecCoin](#lbm.base.v1.DecCoin) | repeated |  |






<a name="lbm.distribution.v1.DelegatorStartingInfo"></a>

### DelegatorStartingInfo
DelegatorStartingInfo represents the starting info for a delegator reward
period. It tracks the previous validator period, the delegation's amount of
staking token, and the creation height (to check later on if any slashes have
occurred). NOTE: Even though validators are slashed to whole staking tokens,
the delegators within the validator may be left with less than a full token,
thus sdk.Dec is used.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `previous_period` | [uint64](#uint64) |  |  |
| `stake` | [string](#string) |  |  |
| `height` | [uint64](#uint64) |  |  |






<a name="lbm.distribution.v1.FeePool"></a>

### FeePool
FeePool is the global fee pool for distribution.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `community_pool` | [lbm.base.v1.DecCoin](#lbm.base.v1.DecCoin) | repeated |  |






<a name="lbm.distribution.v1.Params"></a>

### Params
Params defines the set of params for the distribution module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `community_tax` | [string](#string) |  |  |
| `base_proposer_reward` | [string](#string) |  |  |
| `bonus_proposer_reward` | [string](#string) |  |  |
| `withdraw_addr_enabled` | [bool](#bool) |  |  |






<a name="lbm.distribution.v1.ValidatorAccumulatedCommission"></a>

### ValidatorAccumulatedCommission
ValidatorAccumulatedCommission represents accumulated commission
for a validator kept as a running counter, can be withdrawn at any time.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commission` | [lbm.base.v1.DecCoin](#lbm.base.v1.DecCoin) | repeated |  |






<a name="lbm.distribution.v1.ValidatorCurrentRewards"></a>

### ValidatorCurrentRewards
ValidatorCurrentRewards represents current rewards and current
period for a validator kept as a running counter and incremented
each block as long as the validator's tokens remain constant.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rewards` | [lbm.base.v1.DecCoin](#lbm.base.v1.DecCoin) | repeated |  |
| `period` | [uint64](#uint64) |  |  |






<a name="lbm.distribution.v1.ValidatorHistoricalRewards"></a>

### ValidatorHistoricalRewards
ValidatorHistoricalRewards represents historical rewards for a validator.
Height is implicit within the store key.
Cumulative reward ratio is the sum from the zeroeth period
until this period of rewards / tokens, per the spec.
The reference count indicates the number of objects
which might need to reference this historical entry at any point.
ReferenceCount =
   number of outstanding delegations which ended the associated period (and
   might need to read that record)
 + number of slashes which ended the associated period (and might need to
 read that record)
 + one per validator for the zeroeth period, set on initialization


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cumulative_reward_ratio` | [lbm.base.v1.DecCoin](#lbm.base.v1.DecCoin) | repeated |  |
| `reference_count` | [uint32](#uint32) |  |  |






<a name="lbm.distribution.v1.ValidatorOutstandingRewards"></a>

### ValidatorOutstandingRewards
ValidatorOutstandingRewards represents outstanding (un-withdrawn) rewards
for a validator inexpensive to track, allows simple sanity checks.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rewards` | [lbm.base.v1.DecCoin](#lbm.base.v1.DecCoin) | repeated |  |






<a name="lbm.distribution.v1.ValidatorSlashEvent"></a>

### ValidatorSlashEvent
ValidatorSlashEvent represents a validator slash event.
Height is implicit within the store key.
This is needed to calculate appropriate amount of staking tokens
for delegations which are withdrawn after a slash has occurred.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_period` | [uint64](#uint64) |  |  |
| `fraction` | [string](#string) |  |  |






<a name="lbm.distribution.v1.ValidatorSlashEvents"></a>

### ValidatorSlashEvents
ValidatorSlashEvents is a collection of ValidatorSlashEvent messages.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_slash_events` | [ValidatorSlashEvent](#lbm.distribution.v1.ValidatorSlashEvent) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/distribution/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/distribution/v1/genesis.proto



<a name="lbm.distribution.v1.DelegatorStartingInfoRecord"></a>

### DelegatorStartingInfoRecord
DelegatorStartingInfoRecord used for import / export via genesis json.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address is the address of the delegator. |
| `validator_address` | [string](#string) |  | validator_address is the address of the validator. |
| `starting_info` | [DelegatorStartingInfo](#lbm.distribution.v1.DelegatorStartingInfo) |  | starting_info defines the starting info of a delegator. |






<a name="lbm.distribution.v1.DelegatorWithdrawInfo"></a>

### DelegatorWithdrawInfo
DelegatorWithdrawInfo is the address for where distributions rewards are
withdrawn to by default this struct is only used at genesis to feed in
default withdraw addresses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address is the address of the delegator. |
| `withdraw_address` | [string](#string) |  | withdraw_address is the address to withdraw the delegation rewards to. |






<a name="lbm.distribution.v1.GenesisState"></a>

### GenesisState
GenesisState defines the distribution module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#lbm.distribution.v1.Params) |  | params defines all the paramaters of the module. |
| `fee_pool` | [FeePool](#lbm.distribution.v1.FeePool) |  | fee_pool defines the fee pool at genesis. |
| `delegator_withdraw_infos` | [DelegatorWithdrawInfo](#lbm.distribution.v1.DelegatorWithdrawInfo) | repeated | fee_pool defines the delegator withdraw infos at genesis. |
| `previous_proposer` | [string](#string) |  | fee_pool defines the previous proposer at genesis. |
| `outstanding_rewards` | [ValidatorOutstandingRewardsRecord](#lbm.distribution.v1.ValidatorOutstandingRewardsRecord) | repeated | fee_pool defines the outstanding rewards of all validators at genesis. |
| `validator_accumulated_commissions` | [ValidatorAccumulatedCommissionRecord](#lbm.distribution.v1.ValidatorAccumulatedCommissionRecord) | repeated | fee_pool defines the accumulated commisions of all validators at genesis. |
| `validator_historical_rewards` | [ValidatorHistoricalRewardsRecord](#lbm.distribution.v1.ValidatorHistoricalRewardsRecord) | repeated | fee_pool defines the historical rewards of all validators at genesis. |
| `validator_current_rewards` | [ValidatorCurrentRewardsRecord](#lbm.distribution.v1.ValidatorCurrentRewardsRecord) | repeated | fee_pool defines the current rewards of all validators at genesis. |
| `delegator_starting_infos` | [DelegatorStartingInfoRecord](#lbm.distribution.v1.DelegatorStartingInfoRecord) | repeated | fee_pool defines the delegator starting infos at genesis. |
| `validator_slash_events` | [ValidatorSlashEventRecord](#lbm.distribution.v1.ValidatorSlashEventRecord) | repeated | fee_pool defines the validator slash events at genesis. |






<a name="lbm.distribution.v1.ValidatorAccumulatedCommissionRecord"></a>

### ValidatorAccumulatedCommissionRecord
ValidatorAccumulatedCommissionRecord is used for import / export via genesis
json.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address is the address of the validator. |
| `accumulated` | [ValidatorAccumulatedCommission](#lbm.distribution.v1.ValidatorAccumulatedCommission) |  | accumulated is the accumulated commission of a validator. |






<a name="lbm.distribution.v1.ValidatorCurrentRewardsRecord"></a>

### ValidatorCurrentRewardsRecord
ValidatorCurrentRewardsRecord is used for import / export via genesis json.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address is the address of the validator. |
| `rewards` | [ValidatorCurrentRewards](#lbm.distribution.v1.ValidatorCurrentRewards) |  | rewards defines the current rewards of a validator. |






<a name="lbm.distribution.v1.ValidatorHistoricalRewardsRecord"></a>

### ValidatorHistoricalRewardsRecord
ValidatorHistoricalRewardsRecord is used for import / export via genesis
json.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address is the address of the validator. |
| `period` | [uint64](#uint64) |  | period defines the period the historical rewards apply to. |
| `rewards` | [ValidatorHistoricalRewards](#lbm.distribution.v1.ValidatorHistoricalRewards) |  | rewards defines the historical rewards of a validator. |






<a name="lbm.distribution.v1.ValidatorOutstandingRewardsRecord"></a>

### ValidatorOutstandingRewardsRecord
ValidatorOutstandingRewardsRecord is used for import/export via genesis json.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address is the address of the validator. |
| `outstanding_rewards` | [lbm.base.v1.DecCoin](#lbm.base.v1.DecCoin) | repeated | outstanding_rewards represents the oustanding rewards of a validator. |






<a name="lbm.distribution.v1.ValidatorSlashEventRecord"></a>

### ValidatorSlashEventRecord
ValidatorSlashEventRecord is used for import / export via genesis json.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address is the address of the validator. |
| `height` | [uint64](#uint64) |  | height defines the block height at which the slash event occured. |
| `period` | [uint64](#uint64) |  | period is the period of the slash event. |
| `validator_slash_event` | [ValidatorSlashEvent](#lbm.distribution.v1.ValidatorSlashEvent) |  | validator_slash_event describes the slash event. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/distribution/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/distribution/v1/query.proto



<a name="lbm.distribution.v1.QueryCommunityPoolRequest"></a>

### QueryCommunityPoolRequest
QueryCommunityPoolRequest is the request type for the Query/CommunityPool RPC
method.






<a name="lbm.distribution.v1.QueryCommunityPoolResponse"></a>

### QueryCommunityPoolResponse
QueryCommunityPoolResponse is the response type for the Query/CommunityPool
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pool` | [lbm.base.v1.DecCoin](#lbm.base.v1.DecCoin) | repeated | pool defines community pool's coins. |






<a name="lbm.distribution.v1.QueryDelegationRewardsRequest"></a>

### QueryDelegationRewardsRequest
QueryDelegationRewardsRequest is the request type for the
Query/DelegationRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address defines the delegator address to query for. |
| `validator_address` | [string](#string) |  | validator_address defines the validator address to query for. |






<a name="lbm.distribution.v1.QueryDelegationRewardsResponse"></a>

### QueryDelegationRewardsResponse
QueryDelegationRewardsResponse is the response type for the
Query/DelegationRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rewards` | [lbm.base.v1.DecCoin](#lbm.base.v1.DecCoin) | repeated | rewards defines the rewards accrued by a delegation. |






<a name="lbm.distribution.v1.QueryDelegationTotalRewardsRequest"></a>

### QueryDelegationTotalRewardsRequest
QueryDelegationTotalRewardsRequest is the request type for the
Query/DelegationTotalRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address defines the delegator address to query for. |






<a name="lbm.distribution.v1.QueryDelegationTotalRewardsResponse"></a>

### QueryDelegationTotalRewardsResponse
QueryDelegationTotalRewardsResponse is the response type for the
Query/DelegationTotalRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rewards` | [DelegationDelegatorReward](#lbm.distribution.v1.DelegationDelegatorReward) | repeated | rewards defines all the rewards accrued by a delegator. |
| `total` | [lbm.base.v1.DecCoin](#lbm.base.v1.DecCoin) | repeated | total defines the sum of all the rewards. |






<a name="lbm.distribution.v1.QueryDelegatorValidatorsRequest"></a>

### QueryDelegatorValidatorsRequest
QueryDelegatorValidatorsRequest is the request type for the
Query/DelegatorValidators RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address defines the delegator address to query for. |






<a name="lbm.distribution.v1.QueryDelegatorValidatorsResponse"></a>

### QueryDelegatorValidatorsResponse
QueryDelegatorValidatorsResponse is the response type for the
Query/DelegatorValidators RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validators` | [string](#string) | repeated | validators defines the validators a delegator is delegating for. |






<a name="lbm.distribution.v1.QueryDelegatorWithdrawAddressRequest"></a>

### QueryDelegatorWithdrawAddressRequest
QueryDelegatorWithdrawAddressRequest is the request type for the
Query/DelegatorWithdrawAddress RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address defines the delegator address to query for. |






<a name="lbm.distribution.v1.QueryDelegatorWithdrawAddressResponse"></a>

### QueryDelegatorWithdrawAddressResponse
QueryDelegatorWithdrawAddressResponse is the response type for the
Query/DelegatorWithdrawAddress RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `withdraw_address` | [string](#string) |  | withdraw_address defines the delegator address to query for. |






<a name="lbm.distribution.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="lbm.distribution.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#lbm.distribution.v1.Params) |  | params defines the parameters of the module. |






<a name="lbm.distribution.v1.QueryValidatorCommissionRequest"></a>

### QueryValidatorCommissionRequest
QueryValidatorCommissionRequest is the request type for the
Query/ValidatorCommission RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address defines the validator address to query for. |






<a name="lbm.distribution.v1.QueryValidatorCommissionResponse"></a>

### QueryValidatorCommissionResponse
QueryValidatorCommissionResponse is the response type for the
Query/ValidatorCommission RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commission` | [ValidatorAccumulatedCommission](#lbm.distribution.v1.ValidatorAccumulatedCommission) |  | commission defines the commision the validator received. |






<a name="lbm.distribution.v1.QueryValidatorOutstandingRewardsRequest"></a>

### QueryValidatorOutstandingRewardsRequest
QueryValidatorOutstandingRewardsRequest is the request type for the
Query/ValidatorOutstandingRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address defines the validator address to query for. |






<a name="lbm.distribution.v1.QueryValidatorOutstandingRewardsResponse"></a>

### QueryValidatorOutstandingRewardsResponse
QueryValidatorOutstandingRewardsResponse is the response type for the
Query/ValidatorOutstandingRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rewards` | [ValidatorOutstandingRewards](#lbm.distribution.v1.ValidatorOutstandingRewards) |  |  |






<a name="lbm.distribution.v1.QueryValidatorSlashesRequest"></a>

### QueryValidatorSlashesRequest
QueryValidatorSlashesRequest is the request type for the
Query/ValidatorSlashes RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address defines the validator address to query for. |
| `starting_height` | [uint64](#uint64) |  | starting_height defines the optional starting height to query the slashes. |
| `ending_height` | [uint64](#uint64) |  | starting_height defines the optional ending height to query the slashes. |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.distribution.v1.QueryValidatorSlashesResponse"></a>

### QueryValidatorSlashesResponse
QueryValidatorSlashesResponse is the response type for the
Query/ValidatorSlashes RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `slashes` | [ValidatorSlashEvent](#lbm.distribution.v1.ValidatorSlashEvent) | repeated | slashes defines the slashes the validator received. |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines the pagination in the response. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.distribution.v1.Query"></a>

### Query
Query defines the gRPC querier service for distribution module.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#lbm.distribution.v1.QueryParamsRequest) | [QueryParamsResponse](#lbm.distribution.v1.QueryParamsResponse) | Params queries params of the distribution module. | GET|/lbm/distribution/v1/params|
| `ValidatorOutstandingRewards` | [QueryValidatorOutstandingRewardsRequest](#lbm.distribution.v1.QueryValidatorOutstandingRewardsRequest) | [QueryValidatorOutstandingRewardsResponse](#lbm.distribution.v1.QueryValidatorOutstandingRewardsResponse) | ValidatorOutstandingRewards queries rewards of a validator address. | GET|/lbm/distribution/v1/validators/{validator_address}/outstanding_rewards|
| `ValidatorCommission` | [QueryValidatorCommissionRequest](#lbm.distribution.v1.QueryValidatorCommissionRequest) | [QueryValidatorCommissionResponse](#lbm.distribution.v1.QueryValidatorCommissionResponse) | ValidatorCommission queries accumulated commission for a validator. | GET|/lbm/distribution/v1/validators/{validator_address}/commission|
| `ValidatorSlashes` | [QueryValidatorSlashesRequest](#lbm.distribution.v1.QueryValidatorSlashesRequest) | [QueryValidatorSlashesResponse](#lbm.distribution.v1.QueryValidatorSlashesResponse) | ValidatorSlashes queries slash events of a validator. | GET|/lbm/distribution/v1/validators/{validator_address}/slashes|
| `DelegationRewards` | [QueryDelegationRewardsRequest](#lbm.distribution.v1.QueryDelegationRewardsRequest) | [QueryDelegationRewardsResponse](#lbm.distribution.v1.QueryDelegationRewardsResponse) | DelegationRewards queries the total rewards accrued by a delegation. | GET|/lbm/distribution/v1/delegators/{delegator_address}/rewards/{validator_address}|
| `DelegationTotalRewards` | [QueryDelegationTotalRewardsRequest](#lbm.distribution.v1.QueryDelegationTotalRewardsRequest) | [QueryDelegationTotalRewardsResponse](#lbm.distribution.v1.QueryDelegationTotalRewardsResponse) | DelegationTotalRewards queries the total rewards accrued by a each validator. | GET|/lbm/distribution/v1/delegators/{delegator_address}/rewards|
| `DelegatorValidators` | [QueryDelegatorValidatorsRequest](#lbm.distribution.v1.QueryDelegatorValidatorsRequest) | [QueryDelegatorValidatorsResponse](#lbm.distribution.v1.QueryDelegatorValidatorsResponse) | DelegatorValidators queries the validators of a delegator. | GET|/lbm/distribution/v1/delegators/{delegator_address}/validators|
| `DelegatorWithdrawAddress` | [QueryDelegatorWithdrawAddressRequest](#lbm.distribution.v1.QueryDelegatorWithdrawAddressRequest) | [QueryDelegatorWithdrawAddressResponse](#lbm.distribution.v1.QueryDelegatorWithdrawAddressResponse) | DelegatorWithdrawAddress queries withdraw address of a delegator. | GET|/lbm/distribution/v1/delegators/{delegator_address}/withdraw_address|
| `CommunityPool` | [QueryCommunityPoolRequest](#lbm.distribution.v1.QueryCommunityPoolRequest) | [QueryCommunityPoolResponse](#lbm.distribution.v1.QueryCommunityPoolResponse) | CommunityPool queries the community pool coins. | GET|/lbm/distribution/v1/community_pool|

 <!-- end services -->



<a name="lbm/distribution/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/distribution/v1/tx.proto



<a name="lbm.distribution.v1.MsgFundCommunityPool"></a>

### MsgFundCommunityPool
MsgFundCommunityPool allows an account to directly
fund the community pool.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `amount` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated |  |
| `depositor` | [string](#string) |  |  |






<a name="lbm.distribution.v1.MsgFundCommunityPoolResponse"></a>

### MsgFundCommunityPoolResponse
MsgFundCommunityPoolResponse defines the Msg/FundCommunityPool response type.






<a name="lbm.distribution.v1.MsgSetWithdrawAddress"></a>

### MsgSetWithdrawAddress
MsgSetWithdrawAddress sets the withdraw address for
a delegator (or validator self-delegation).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `withdraw_address` | [string](#string) |  |  |






<a name="lbm.distribution.v1.MsgSetWithdrawAddressResponse"></a>

### MsgSetWithdrawAddressResponse
MsgSetWithdrawAddressResponse defines the Msg/SetWithdrawAddress response type.






<a name="lbm.distribution.v1.MsgWithdrawDelegatorReward"></a>

### MsgWithdrawDelegatorReward
MsgWithdrawDelegatorReward represents delegation withdrawal to a delegator
from a single validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |






<a name="lbm.distribution.v1.MsgWithdrawDelegatorRewardResponse"></a>

### MsgWithdrawDelegatorRewardResponse
MsgWithdrawDelegatorRewardResponse defines the Msg/WithdrawDelegatorReward response type.






<a name="lbm.distribution.v1.MsgWithdrawValidatorCommission"></a>

### MsgWithdrawValidatorCommission
MsgWithdrawValidatorCommission withdraws the full commission to the validator
address.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  |  |






<a name="lbm.distribution.v1.MsgWithdrawValidatorCommissionResponse"></a>

### MsgWithdrawValidatorCommissionResponse
MsgWithdrawValidatorCommissionResponse defines the Msg/WithdrawValidatorCommission response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.distribution.v1.Msg"></a>

### Msg
Msg defines the distribution Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `SetWithdrawAddress` | [MsgSetWithdrawAddress](#lbm.distribution.v1.MsgSetWithdrawAddress) | [MsgSetWithdrawAddressResponse](#lbm.distribution.v1.MsgSetWithdrawAddressResponse) | SetWithdrawAddress defines a method to change the withdraw address for a delegator (or validator self-delegation). | |
| `WithdrawDelegatorReward` | [MsgWithdrawDelegatorReward](#lbm.distribution.v1.MsgWithdrawDelegatorReward) | [MsgWithdrawDelegatorRewardResponse](#lbm.distribution.v1.MsgWithdrawDelegatorRewardResponse) | WithdrawDelegatorReward defines a method to withdraw rewards of delegator from a single validator. | |
| `WithdrawValidatorCommission` | [MsgWithdrawValidatorCommission](#lbm.distribution.v1.MsgWithdrawValidatorCommission) | [MsgWithdrawValidatorCommissionResponse](#lbm.distribution.v1.MsgWithdrawValidatorCommissionResponse) | WithdrawValidatorCommission defines a method to withdraw the full commission to the validator address. | |
| `FundCommunityPool` | [MsgFundCommunityPool](#lbm.distribution.v1.MsgFundCommunityPool) | [MsgFundCommunityPoolResponse](#lbm.distribution.v1.MsgFundCommunityPoolResponse) | FundCommunityPool defines a method to allow an account to directly fund the community pool. | |

 <!-- end services -->



<a name="lbm/evidence/v1/evidence.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/evidence/v1/evidence.proto



<a name="lbm.evidence.v1.Equivocation"></a>

### Equivocation
Equivocation implements the Evidence interface and defines evidence of double
signing misbehavior.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  |  |
| `time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| `power` | [int64](#int64) |  |  |
| `consensus_address` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/evidence/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/evidence/v1/genesis.proto



<a name="lbm.evidence.v1.GenesisState"></a>

### GenesisState
GenesisState defines the evidence module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `evidence` | [google.protobuf.Any](#google.protobuf.Any) | repeated | evidence defines all the evidence at genesis. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/evidence/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/evidence/v1/query.proto



<a name="lbm.evidence.v1.QueryAllEvidenceRequest"></a>

### QueryAllEvidenceRequest
QueryEvidenceRequest is the request type for the Query/AllEvidence RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.evidence.v1.QueryAllEvidenceResponse"></a>

### QueryAllEvidenceResponse
QueryAllEvidenceResponse is the response type for the Query/AllEvidence RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `evidence` | [google.protobuf.Any](#google.protobuf.Any) | repeated | evidence returns all evidences. |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.evidence.v1.QueryEvidenceRequest"></a>

### QueryEvidenceRequest
QueryEvidenceRequest is the request type for the Query/Evidence RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `evidence_hash` | [bytes](#bytes) |  | evidence_hash defines the hash of the requested evidence. |






<a name="lbm.evidence.v1.QueryEvidenceResponse"></a>

### QueryEvidenceResponse
QueryEvidenceResponse is the response type for the Query/Evidence RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `evidence` | [google.protobuf.Any](#google.protobuf.Any) |  | evidence returns the requested evidence. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.evidence.v1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Evidence` | [QueryEvidenceRequest](#lbm.evidence.v1.QueryEvidenceRequest) | [QueryEvidenceResponse](#lbm.evidence.v1.QueryEvidenceResponse) | Evidence queries evidence based on evidence hash. | GET|/lbm/evidence/v1/evidence/{evidence_hash}|
| `AllEvidence` | [QueryAllEvidenceRequest](#lbm.evidence.v1.QueryAllEvidenceRequest) | [QueryAllEvidenceResponse](#lbm.evidence.v1.QueryAllEvidenceResponse) | AllEvidence queries all evidence. | GET|/lbm/evidence/v1/evidence|

 <!-- end services -->



<a name="lbm/evidence/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/evidence/v1/tx.proto



<a name="lbm.evidence.v1.MsgSubmitEvidence"></a>

### MsgSubmitEvidence
MsgSubmitEvidence represents a message that supports submitting arbitrary
Evidence of misbehavior such as equivocation or counterfactual signing.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `submitter` | [string](#string) |  |  |
| `evidence` | [google.protobuf.Any](#google.protobuf.Any) |  |  |






<a name="lbm.evidence.v1.MsgSubmitEvidenceResponse"></a>

### MsgSubmitEvidenceResponse
MsgSubmitEvidenceResponse defines the Msg/SubmitEvidence response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [bytes](#bytes) |  | hash defines the hash of the evidence. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.evidence.v1.Msg"></a>

### Msg
Msg defines the evidence Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `SubmitEvidence` | [MsgSubmitEvidence](#lbm.evidence.v1.MsgSubmitEvidence) | [MsgSubmitEvidenceResponse](#lbm.evidence.v1.MsgSubmitEvidenceResponse) | SubmitEvidence submits an arbitrary Evidence of misbehavior such as equivocation or counterfactual signing. | |

 <!-- end services -->



<a name="lbm/feegrant/v1/feegrant.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/feegrant/v1/feegrant.proto



<a name="lbm.feegrant.v1.AllowedMsgAllowance"></a>

### AllowedMsgAllowance
AllowedMsgAllowance creates allowance only for specified message types.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `allowance` | [google.protobuf.Any](#google.protobuf.Any) |  | allowance can be any of basic and filtered fee allowance. |
| `allowed_messages` | [string](#string) | repeated | allowed_messages are the messages for which the grantee has the access. |






<a name="lbm.feegrant.v1.BasicAllowance"></a>

### BasicAllowance
BasicAllowance implements Allowance with a one-time grant of tokens
that optionally expires. The grantee can use up to SpendLimit to cover fees.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `spend_limit` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated | spend_limit specifies the maximum amount of tokens that can be spent by this allowance and will be updated as tokens are spent. If it is empty, there is no spend limit and any amount of coins can be spent. |
| `expiration` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | expiration specifies an optional time when this allowance expires |






<a name="lbm.feegrant.v1.Grant"></a>

### Grant
Grant is stored in the KVStore to record a grant with full context


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  | granter is the address of the user granting an allowance of their funds. |
| `grantee` | [string](#string) |  | grantee is the address of the user being granted an allowance of another user's funds. |
| `allowance` | [google.protobuf.Any](#google.protobuf.Any) |  | allowance can be any of basic and filtered fee allowance. |






<a name="lbm.feegrant.v1.PeriodicAllowance"></a>

### PeriodicAllowance
PeriodicAllowance extends Allowance to allow for both a maximum cap,
as well as a limit per time period.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `basic` | [BasicAllowance](#lbm.feegrant.v1.BasicAllowance) |  | basic specifies a struct of `BasicAllowance` |
| `period` | [google.protobuf.Duration](#google.protobuf.Duration) |  | period specifies the time duration in which period_spend_limit coins can be spent before that allowance is reset |
| `period_spend_limit` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated | period_spend_limit specifies the maximum number of coins that can be spent in the period |
| `period_can_spend` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated | period_can_spend is the number of coins left to be spent before the period_reset time |
| `period_reset` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | period_reset is the time at which this period resets and a new one begins, it is calculated from the start time of the first transaction after the last period ended |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/feegrant/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/feegrant/v1/genesis.proto



<a name="lbm.feegrant.v1.GenesisState"></a>

### GenesisState
GenesisState contains a set of fee allowances, persisted from the store


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `allowances` | [Grant](#lbm.feegrant.v1.Grant) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/feegrant/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/feegrant/v1/query.proto



<a name="lbm.feegrant.v1.QueryAllowanceRequest"></a>

### QueryAllowanceRequest
QueryAllowanceRequest is the request type for the Query/Allowance RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  | granter is the address of the user granting an allowance of their funds. |
| `grantee` | [string](#string) |  | grantee is the address of the user being granted an allowance of another user's funds. |






<a name="lbm.feegrant.v1.QueryAllowanceResponse"></a>

### QueryAllowanceResponse
QueryAllowanceResponse is the response type for the Query/Allowance RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `allowance` | [Grant](#lbm.feegrant.v1.Grant) |  | allowance is a allowance granted for grantee by granter. |






<a name="lbm.feegrant.v1.QueryAllowancesRequest"></a>

### QueryAllowancesRequest
QueryAllowancesRequest is the request type for the Query/Allowances RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grantee` | [string](#string) |  |  |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an pagination for the request. |






<a name="lbm.feegrant.v1.QueryAllowancesResponse"></a>

### QueryAllowancesResponse
QueryAllowancesResponse is the response type for the Query/Allowances RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `allowances` | [Grant](#lbm.feegrant.v1.Grant) | repeated | allowances are allowance's granted for grantee by granter. |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines an pagination for the response. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.feegrant.v1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Allowance` | [QueryAllowanceRequest](#lbm.feegrant.v1.QueryAllowanceRequest) | [QueryAllowanceResponse](#lbm.feegrant.v1.QueryAllowanceResponse) | Allowance returns fee granted to the grantee by the granter. | GET|/lbm/feegrant/v1/allowance/{granter}/{grantee}|
| `Allowances` | [QueryAllowancesRequest](#lbm.feegrant.v1.QueryAllowancesRequest) | [QueryAllowancesResponse](#lbm.feegrant.v1.QueryAllowancesResponse) | Allowances returns all the grants for address. | GET|/lbm/feegrant/v1/allowances/{grantee}|

 <!-- end services -->



<a name="lbm/feegrant/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/feegrant/v1/tx.proto



<a name="lbm.feegrant.v1.MsgGrantAllowance"></a>

### MsgGrantAllowance
MsgGrantAllowance adds permission for Grantee to spend up to Allowance
of fees from the account of Granter.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  | granter is the address of the user granting an allowance of their funds. |
| `grantee` | [string](#string) |  | grantee is the address of the user being granted an allowance of another user's funds. |
| `allowance` | [google.protobuf.Any](#google.protobuf.Any) |  | allowance can be any of basic and filtered fee allowance. |






<a name="lbm.feegrant.v1.MsgGrantAllowanceResponse"></a>

### MsgGrantAllowanceResponse
MsgGrantAllowanceResponse defines the Msg/GrantAllowanceResponse response type.






<a name="lbm.feegrant.v1.MsgRevokeAllowance"></a>

### MsgRevokeAllowance
MsgRevokeAllowance removes any existing Allowance from Granter to Grantee.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  | granter is the address of the user granting an allowance of their funds. |
| `grantee` | [string](#string) |  | grantee is the address of the user being granted an allowance of another user's funds. |






<a name="lbm.feegrant.v1.MsgRevokeAllowanceResponse"></a>

### MsgRevokeAllowanceResponse
MsgRevokeAllowanceResponse defines the Msg/RevokeAllowanceResponse response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.feegrant.v1.Msg"></a>

### Msg
Msg defines the feegrant msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `GrantAllowance` | [MsgGrantAllowance](#lbm.feegrant.v1.MsgGrantAllowance) | [MsgGrantAllowanceResponse](#lbm.feegrant.v1.MsgGrantAllowanceResponse) | GrantAllowance grants fee allowance to the grantee on the granter's account with the provided expiration time. | |
| `RevokeAllowance` | [MsgRevokeAllowance](#lbm.feegrant.v1.MsgRevokeAllowance) | [MsgRevokeAllowanceResponse](#lbm.feegrant.v1.MsgRevokeAllowanceResponse) | RevokeAllowance revokes any fee allowance of granter's account that has been granted to the grantee. | |

 <!-- end services -->



<a name="lbm/genutil/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/genutil/v1/genesis.proto



<a name="lbm.genutil.v1.GenesisState"></a>

### GenesisState
GenesisState defines the raw genesis transaction in JSON.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gen_txs` | [bytes](#bytes) | repeated | gen_txs defines the genesis transactions. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/gov/v1/gov.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/gov/v1/gov.proto



<a name="lbm.gov.v1.Deposit"></a>

### Deposit
Deposit defines an amount deposited by an account address to an active
proposal.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  |  |
| `depositor` | [string](#string) |  |  |
| `amount` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated |  |






<a name="lbm.gov.v1.DepositParams"></a>

### DepositParams
DepositParams defines the params for deposits on governance proposals.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `min_deposit` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated | Minimum deposit for a proposal to enter voting period. |
| `max_deposit_period` | [google.protobuf.Duration](#google.protobuf.Duration) |  | Maximum period for Atom holders to deposit on a proposal. Initial value: 2 months. |






<a name="lbm.gov.v1.Proposal"></a>

### Proposal
Proposal defines the core field members of a governance proposal.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  |  |
| `content` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `status` | [ProposalStatus](#lbm.gov.v1.ProposalStatus) |  |  |
| `final_tally_result` | [TallyResult](#lbm.gov.v1.TallyResult) |  |  |
| `submit_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| `deposit_end_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| `total_deposit` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated |  |
| `voting_start_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| `voting_end_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |






<a name="lbm.gov.v1.TallyParams"></a>

### TallyParams
TallyParams defines the params for tallying votes on governance proposals.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `quorum` | [bytes](#bytes) |  | Minimum percentage of total stake needed to vote for a result to be considered valid. |
| `threshold` | [bytes](#bytes) |  | Minimum proportion of Yes votes for proposal to pass. Default value: 0.5. |
| `veto_threshold` | [bytes](#bytes) |  | Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Default value: 1/3. |






<a name="lbm.gov.v1.TallyResult"></a>

### TallyResult
TallyResult defines a standard tally for a governance proposal.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `yes` | [string](#string) |  |  |
| `abstain` | [string](#string) |  |  |
| `no` | [string](#string) |  |  |
| `no_with_veto` | [string](#string) |  |  |






<a name="lbm.gov.v1.TextProposal"></a>

### TextProposal
TextProposal defines a standard text proposal whose changes need to be
manually updated in case of approval.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |






<a name="lbm.gov.v1.Vote"></a>

### Vote
Vote defines a vote on a governance proposal.
A Vote consists of a proposal ID, the voter, and the vote option.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  |  |
| `voter` | [string](#string) |  |  |
| `option` | [VoteOption](#lbm.gov.v1.VoteOption) |  | **Deprecated.** Deprecated: Prefer to use `options` instead. This field is set in queries if and only if `len(options) == 1` and that option has weight 1. In all other cases, this field will default to VOTE_OPTION_UNSPECIFIED. |
| `options` | [WeightedVoteOption](#lbm.gov.v1.WeightedVoteOption) | repeated |  |






<a name="lbm.gov.v1.VotingParams"></a>

### VotingParams
VotingParams defines the params for voting on governance proposals.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `voting_period` | [google.protobuf.Duration](#google.protobuf.Duration) |  | Length of the voting period. |






<a name="lbm.gov.v1.WeightedVoteOption"></a>

### WeightedVoteOption
WeightedVoteOption defines a unit of vote for vote split.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `option` | [VoteOption](#lbm.gov.v1.VoteOption) |  |  |
| `weight` | [string](#string) |  |  |





 <!-- end messages -->


<a name="lbm.gov.v1.ProposalStatus"></a>

### ProposalStatus
ProposalStatus enumerates the valid statuses of a proposal.

| Name | Number | Description |
| ---- | ------ | ----------- |
| PROPOSAL_STATUS_UNSPECIFIED | 0 | PROPOSAL_STATUS_UNSPECIFIED defines the default propopsal status. |
| PROPOSAL_STATUS_DEPOSIT_PERIOD | 1 | PROPOSAL_STATUS_DEPOSIT_PERIOD defines a proposal status during the deposit period. |
| PROPOSAL_STATUS_VOTING_PERIOD | 2 | PROPOSAL_STATUS_VOTING_PERIOD defines a proposal status during the voting period. |
| PROPOSAL_STATUS_PASSED | 3 | PROPOSAL_STATUS_PASSED defines a proposal status of a proposal that has passed. |
| PROPOSAL_STATUS_REJECTED | 4 | PROPOSAL_STATUS_REJECTED defines a proposal status of a proposal that has been rejected. |
| PROPOSAL_STATUS_FAILED | 5 | PROPOSAL_STATUS_FAILED defines a proposal status of a proposal that has failed. |



<a name="lbm.gov.v1.VoteOption"></a>

### VoteOption
VoteOption enumerates the valid vote options for a given governance proposal.

| Name | Number | Description |
| ---- | ------ | ----------- |
| VOTE_OPTION_UNSPECIFIED | 0 | VOTE_OPTION_UNSPECIFIED defines a no-op vote option. |
| VOTE_OPTION_YES | 1 | VOTE_OPTION_YES defines a yes vote option. |
| VOTE_OPTION_ABSTAIN | 2 | VOTE_OPTION_ABSTAIN defines an abstain vote option. |
| VOTE_OPTION_NO | 3 | VOTE_OPTION_NO defines a no vote option. |
| VOTE_OPTION_NO_WITH_VETO | 4 | VOTE_OPTION_NO_WITH_VETO defines a no with veto vote option. |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/gov/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/gov/v1/genesis.proto



<a name="lbm.gov.v1.GenesisState"></a>

### GenesisState
GenesisState defines the gov module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `starting_proposal_id` | [uint64](#uint64) |  | starting_proposal_id is the ID of the starting proposal. |
| `deposits` | [Deposit](#lbm.gov.v1.Deposit) | repeated | deposits defines all the deposits present at genesis. |
| `votes` | [Vote](#lbm.gov.v1.Vote) | repeated | votes defines all the votes present at genesis. |
| `proposals` | [Proposal](#lbm.gov.v1.Proposal) | repeated | proposals defines all the proposals present at genesis. |
| `deposit_params` | [DepositParams](#lbm.gov.v1.DepositParams) |  | params defines all the paramaters of related to deposit. |
| `voting_params` | [VotingParams](#lbm.gov.v1.VotingParams) |  | params defines all the paramaters of related to voting. |
| `tally_params` | [TallyParams](#lbm.gov.v1.TallyParams) |  | params defines all the paramaters of related to tally. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/gov/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/gov/v1/query.proto



<a name="lbm.gov.v1.QueryDepositRequest"></a>

### QueryDepositRequest
QueryDepositRequest is the request type for the Query/Deposit RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id defines the unique id of the proposal. |
| `depositor` | [string](#string) |  | depositor defines the deposit addresses from the proposals. |






<a name="lbm.gov.v1.QueryDepositResponse"></a>

### QueryDepositResponse
QueryDepositResponse is the response type for the Query/Deposit RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `deposit` | [Deposit](#lbm.gov.v1.Deposit) |  | deposit defines the requested deposit. |






<a name="lbm.gov.v1.QueryDepositsRequest"></a>

### QueryDepositsRequest
QueryDepositsRequest is the request type for the Query/Deposits RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id defines the unique id of the proposal. |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.gov.v1.QueryDepositsResponse"></a>

### QueryDepositsResponse
QueryDepositsResponse is the response type for the Query/Deposits RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `deposits` | [Deposit](#lbm.gov.v1.Deposit) | repeated |  |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.gov.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params_type` | [string](#string) |  | params_type defines which parameters to query for, can be one of "voting", "tallying" or "deposit". |






<a name="lbm.gov.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `voting_params` | [VotingParams](#lbm.gov.v1.VotingParams) |  | voting_params defines the parameters related to voting. |
| `deposit_params` | [DepositParams](#lbm.gov.v1.DepositParams) |  | deposit_params defines the parameters related to deposit. |
| `tally_params` | [TallyParams](#lbm.gov.v1.TallyParams) |  | tally_params defines the parameters related to tally. |






<a name="lbm.gov.v1.QueryProposalRequest"></a>

### QueryProposalRequest
QueryProposalRequest is the request type for the Query/Proposal RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id defines the unique id of the proposal. |






<a name="lbm.gov.v1.QueryProposalResponse"></a>

### QueryProposalResponse
QueryProposalResponse is the response type for the Query/Proposal RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal` | [Proposal](#lbm.gov.v1.Proposal) |  |  |






<a name="lbm.gov.v1.QueryProposalsRequest"></a>

### QueryProposalsRequest
QueryProposalsRequest is the request type for the Query/Proposals RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_status` | [ProposalStatus](#lbm.gov.v1.ProposalStatus) |  | proposal_status defines the status of the proposals. |
| `voter` | [string](#string) |  | voter defines the voter address for the proposals. |
| `depositor` | [string](#string) |  | depositor defines the deposit addresses from the proposals. |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.gov.v1.QueryProposalsResponse"></a>

### QueryProposalsResponse
QueryProposalsResponse is the response type for the Query/Proposals RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposals` | [Proposal](#lbm.gov.v1.Proposal) | repeated |  |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.gov.v1.QueryTallyResultRequest"></a>

### QueryTallyResultRequest
QueryTallyResultRequest is the request type for the Query/Tally RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id defines the unique id of the proposal. |






<a name="lbm.gov.v1.QueryTallyResultResponse"></a>

### QueryTallyResultResponse
QueryTallyResultResponse is the response type for the Query/Tally RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tally` | [TallyResult](#lbm.gov.v1.TallyResult) |  | tally defines the requested tally. |






<a name="lbm.gov.v1.QueryVoteRequest"></a>

### QueryVoteRequest
QueryVoteRequest is the request type for the Query/Vote RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id defines the unique id of the proposal. |
| `voter` | [string](#string) |  | voter defines the oter address for the proposals. |






<a name="lbm.gov.v1.QueryVoteResponse"></a>

### QueryVoteResponse
QueryVoteResponse is the response type for the Query/Vote RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vote` | [Vote](#lbm.gov.v1.Vote) |  | vote defined the queried vote. |






<a name="lbm.gov.v1.QueryVotesRequest"></a>

### QueryVotesRequest
QueryVotesRequest is the request type for the Query/Votes RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id defines the unique id of the proposal. |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.gov.v1.QueryVotesResponse"></a>

### QueryVotesResponse
QueryVotesResponse is the response type for the Query/Votes RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `votes` | [Vote](#lbm.gov.v1.Vote) | repeated | votes defined the queried votes. |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines the pagination in the response. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.gov.v1.Query"></a>

### Query
Query defines the gRPC querier service for gov module

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Proposal` | [QueryProposalRequest](#lbm.gov.v1.QueryProposalRequest) | [QueryProposalResponse](#lbm.gov.v1.QueryProposalResponse) | Proposal queries proposal details based on ProposalID. | GET|/lbm/gov/v1/proposals/{proposal_id}|
| `Proposals` | [QueryProposalsRequest](#lbm.gov.v1.QueryProposalsRequest) | [QueryProposalsResponse](#lbm.gov.v1.QueryProposalsResponse) | Proposals queries all proposals based on given status. | GET|/lbm/gov/v1/proposals|
| `Vote` | [QueryVoteRequest](#lbm.gov.v1.QueryVoteRequest) | [QueryVoteResponse](#lbm.gov.v1.QueryVoteResponse) | Vote queries voted information based on proposalID, voterAddr. | GET|/lbm/gov/v1/proposals/{proposal_id}/votes/{voter}|
| `Votes` | [QueryVotesRequest](#lbm.gov.v1.QueryVotesRequest) | [QueryVotesResponse](#lbm.gov.v1.QueryVotesResponse) | Votes queries votes of a given proposal. | GET|/lbm/gov/v1/proposals/{proposal_id}/votes|
| `Params` | [QueryParamsRequest](#lbm.gov.v1.QueryParamsRequest) | [QueryParamsResponse](#lbm.gov.v1.QueryParamsResponse) | Params queries all parameters of the gov module. | GET|/lbm/gov/v1/params/{params_type}|
| `Deposit` | [QueryDepositRequest](#lbm.gov.v1.QueryDepositRequest) | [QueryDepositResponse](#lbm.gov.v1.QueryDepositResponse) | Deposit queries single deposit information based proposalID, depositAddr. | GET|/lbm/gov/v1/proposals/{proposal_id}/deposits/{depositor}|
| `Deposits` | [QueryDepositsRequest](#lbm.gov.v1.QueryDepositsRequest) | [QueryDepositsResponse](#lbm.gov.v1.QueryDepositsResponse) | Deposits queries all deposits of a single proposal. | GET|/lbm/gov/v1/proposals/{proposal_id}/deposits|
| `TallyResult` | [QueryTallyResultRequest](#lbm.gov.v1.QueryTallyResultRequest) | [QueryTallyResultResponse](#lbm.gov.v1.QueryTallyResultResponse) | TallyResult queries the tally of a proposal vote. | GET|/lbm/gov/v1/proposals/{proposal_id}/tally|

 <!-- end services -->



<a name="lbm/gov/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/gov/v1/tx.proto



<a name="lbm.gov.v1.MsgDeposit"></a>

### MsgDeposit
MsgDeposit defines a message to submit a deposit to an existing proposal.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  |  |
| `depositor` | [string](#string) |  |  |
| `amount` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated |  |






<a name="lbm.gov.v1.MsgDepositResponse"></a>

### MsgDepositResponse
MsgDepositResponse defines the Msg/Deposit response type.






<a name="lbm.gov.v1.MsgSubmitProposal"></a>

### MsgSubmitProposal
MsgSubmitProposal defines an sdk.Msg type that supports submitting arbitrary
proposal Content.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `content` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `initial_deposit` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated |  |
| `proposer` | [string](#string) |  |  |






<a name="lbm.gov.v1.MsgSubmitProposalResponse"></a>

### MsgSubmitProposalResponse
MsgSubmitProposalResponse defines the Msg/SubmitProposal response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  |  |






<a name="lbm.gov.v1.MsgVote"></a>

### MsgVote
MsgVote defines a message to cast a vote.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  |  |
| `voter` | [string](#string) |  |  |
| `option` | [VoteOption](#lbm.gov.v1.VoteOption) |  |  |






<a name="lbm.gov.v1.MsgVoteResponse"></a>

### MsgVoteResponse
MsgVoteResponse defines the Msg/Vote response type.






<a name="lbm.gov.v1.MsgVoteWeighted"></a>

### MsgVoteWeighted
MsgVoteWeighted defines a message to cast a vote.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  |  |
| `voter` | [string](#string) |  |  |
| `options` | [WeightedVoteOption](#lbm.gov.v1.WeightedVoteOption) | repeated |  |






<a name="lbm.gov.v1.MsgVoteWeightedResponse"></a>

### MsgVoteWeightedResponse
MsgVoteWeightedResponse defines the Msg/VoteWeighted response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.gov.v1.Msg"></a>

### Msg
Msg defines the bank Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `SubmitProposal` | [MsgSubmitProposal](#lbm.gov.v1.MsgSubmitProposal) | [MsgSubmitProposalResponse](#lbm.gov.v1.MsgSubmitProposalResponse) | SubmitProposal defines a method to create new proposal given a content. | |
| `Vote` | [MsgVote](#lbm.gov.v1.MsgVote) | [MsgVoteResponse](#lbm.gov.v1.MsgVoteResponse) | Vote defines a method to add a vote on a specific proposal. | |
| `VoteWeighted` | [MsgVoteWeighted](#lbm.gov.v1.MsgVoteWeighted) | [MsgVoteWeightedResponse](#lbm.gov.v1.MsgVoteWeightedResponse) | VoteWeighted defines a method to add a weighted vote on a specific proposal. | |
| `Deposit` | [MsgDeposit](#lbm.gov.v1.MsgDeposit) | [MsgDepositResponse](#lbm.gov.v1.MsgDepositResponse) | Deposit defines a method to add deposit on a specific proposal. | |

 <!-- end services -->



<a name="lbm/mint/v1/mint.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/mint/v1/mint.proto



<a name="lbm.mint.v1.Minter"></a>

### Minter
Minter represents the minting state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `inflation` | [string](#string) |  | current annual inflation rate |
| `annual_provisions` | [string](#string) |  | current annual expected provisions |






<a name="lbm.mint.v1.Params"></a>

### Params
Params holds parameters for the mint module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `mint_denom` | [string](#string) |  | type of coin to mint |
| `inflation_rate_change` | [string](#string) |  | maximum annual change in inflation rate |
| `inflation_max` | [string](#string) |  | maximum inflation rate |
| `inflation_min` | [string](#string) |  | minimum inflation rate |
| `goal_bonded` | [string](#string) |  | goal of percent bonded atoms |
| `blocks_per_year` | [uint64](#uint64) |  | expected blocks per year |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/mint/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/mint/v1/genesis.proto



<a name="lbm.mint.v1.GenesisState"></a>

### GenesisState
GenesisState defines the mint module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `minter` | [Minter](#lbm.mint.v1.Minter) |  | minter is a space for holding current inflation information. |
| `params` | [Params](#lbm.mint.v1.Params) |  | params defines all the paramaters of the module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/mint/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/mint/v1/query.proto



<a name="lbm.mint.v1.QueryAnnualProvisionsRequest"></a>

### QueryAnnualProvisionsRequest
QueryAnnualProvisionsRequest is the request type for the
Query/AnnualProvisions RPC method.






<a name="lbm.mint.v1.QueryAnnualProvisionsResponse"></a>

### QueryAnnualProvisionsResponse
QueryAnnualProvisionsResponse is the response type for the
Query/AnnualProvisions RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `annual_provisions` | [bytes](#bytes) |  | annual_provisions is the current minting annual provisions value. |






<a name="lbm.mint.v1.QueryInflationRequest"></a>

### QueryInflationRequest
QueryInflationRequest is the request type for the Query/Inflation RPC method.






<a name="lbm.mint.v1.QueryInflationResponse"></a>

### QueryInflationResponse
QueryInflationResponse is the response type for the Query/Inflation RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `inflation` | [bytes](#bytes) |  | inflation is the current minting inflation value. |






<a name="lbm.mint.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="lbm.mint.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#lbm.mint.v1.Params) |  | params defines the parameters of the module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.mint.v1.Query"></a>

### Query
Query provides defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#lbm.mint.v1.QueryParamsRequest) | [QueryParamsResponse](#lbm.mint.v1.QueryParamsResponse) | Params returns the total set of minting parameters. | GET|/lbm/mint/v1/params|
| `Inflation` | [QueryInflationRequest](#lbm.mint.v1.QueryInflationRequest) | [QueryInflationResponse](#lbm.mint.v1.QueryInflationResponse) | Inflation returns the current minting inflation value. | GET|/lbm/mint/v1/inflation|
| `AnnualProvisions` | [QueryAnnualProvisionsRequest](#lbm.mint.v1.QueryAnnualProvisionsRequest) | [QueryAnnualProvisionsResponse](#lbm.mint.v1.QueryAnnualProvisionsResponse) | AnnualProvisions current minting annual provisions value. | GET|/lbm/mint/v1/annual_provisions|

 <!-- end services -->



<a name="lbm/params/v1/params.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/params/v1/params.proto



<a name="lbm.params.v1.ParamChange"></a>

### ParamChange
ParamChange defines an individual parameter change, for use in
ParameterChangeProposal.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subspace` | [string](#string) |  |  |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |






<a name="lbm.params.v1.ParameterChangeProposal"></a>

### ParameterChangeProposal
ParameterChangeProposal defines a proposal to change one or more parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `changes` | [ParamChange](#lbm.params.v1.ParamChange) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/params/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/params/v1/query.proto



<a name="lbm.params.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is request type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subspace` | [string](#string) |  | subspace defines the module to query the parameter for. |
| `key` | [string](#string) |  | key defines the key of the parameter in the subspace. |






<a name="lbm.params.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `param` | [ParamChange](#lbm.params.v1.ParamChange) |  | param defines the queried parameter. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.params.v1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#lbm.params.v1.QueryParamsRequest) | [QueryParamsResponse](#lbm.params.v1.QueryParamsResponse) | Params queries a specific parameter of a module, given its subspace and key. | GET|/lbm/params/v1/params|

 <!-- end services -->



<a name="lbm/slashing/v1/slashing.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/slashing/v1/slashing.proto



<a name="lbm.slashing.v1.Params"></a>

### Params
Params represents the parameters used for by the slashing module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signed_blocks_window` | [int64](#int64) |  |  |
| `min_signed_per_window` | [bytes](#bytes) |  |  |
| `downtime_jail_duration` | [google.protobuf.Duration](#google.protobuf.Duration) |  |  |
| `slash_fraction_double_sign` | [bytes](#bytes) |  |  |
| `slash_fraction_downtime` | [bytes](#bytes) |  |  |






<a name="lbm.slashing.v1.ValidatorSigningInfo"></a>

### ValidatorSigningInfo
ValidatorSigningInfo defines a validator's signing info for monitoring their
liveness activity.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |
| `jailed_until` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | timestamp validator cannot be unjailed until |
| `tombstoned` | [bool](#bool) |  | whether or not a validator has been tombstoned (killed out of validator set) |
| `missed_blocks_counter` | [int64](#int64) |  | missed blocks counter (to avoid scanning the array every time) |
| `voter_set_counter` | [int64](#int64) |  | how many times the validator joined to voter set |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/slashing/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/slashing/v1/genesis.proto



<a name="lbm.slashing.v1.GenesisState"></a>

### GenesisState
GenesisState defines the slashing module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#lbm.slashing.v1.Params) |  | params defines all the paramaters of related to deposit. |
| `signing_infos` | [SigningInfo](#lbm.slashing.v1.SigningInfo) | repeated | signing_infos represents a map between validator addresses and their signing infos. |
| `missed_blocks` | [ValidatorMissedBlocks](#lbm.slashing.v1.ValidatorMissedBlocks) | repeated | signing_infos represents a map between validator addresses and their missed blocks. |






<a name="lbm.slashing.v1.MissedBlock"></a>

### MissedBlock
MissedBlock contains height and missed status as boolean.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `index` | [int64](#int64) |  | index is the height at which the block was missed. |
| `missed` | [bool](#bool) |  | missed is the missed status. |






<a name="lbm.slashing.v1.SigningInfo"></a>

### SigningInfo
SigningInfo stores validator signing info of corresponding address.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the validator address. |
| `validator_signing_info` | [ValidatorSigningInfo](#lbm.slashing.v1.ValidatorSigningInfo) |  | validator_signing_info represents the signing info of this validator. |






<a name="lbm.slashing.v1.ValidatorMissedBlocks"></a>

### ValidatorMissedBlocks
ValidatorMissedBlocks contains array of missed blocks of corresponding
address.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the validator address. |
| `missed_blocks` | [MissedBlock](#lbm.slashing.v1.MissedBlock) | repeated | missed_blocks is an array of missed blocks by the validator. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/slashing/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/slashing/v1/query.proto



<a name="lbm.slashing.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method






<a name="lbm.slashing.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#lbm.slashing.v1.Params) |  |  |






<a name="lbm.slashing.v1.QuerySigningInfoRequest"></a>

### QuerySigningInfoRequest
QuerySigningInfoRequest is the request type for the Query/SigningInfo RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cons_address` | [string](#string) |  | cons_address is the address to query signing info of |






<a name="lbm.slashing.v1.QuerySigningInfoResponse"></a>

### QuerySigningInfoResponse
QuerySigningInfoResponse is the response type for the Query/SigningInfo RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `val_signing_info` | [ValidatorSigningInfo](#lbm.slashing.v1.ValidatorSigningInfo) |  | val_signing_info is the signing info of requested val cons address |






<a name="lbm.slashing.v1.QuerySigningInfosRequest"></a>

### QuerySigningInfosRequest
QuerySigningInfosRequest is the request type for the Query/SigningInfos RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  |  |






<a name="lbm.slashing.v1.QuerySigningInfosResponse"></a>

### QuerySigningInfosResponse
QuerySigningInfosResponse is the response type for the Query/SigningInfos RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `info` | [ValidatorSigningInfo](#lbm.slashing.v1.ValidatorSigningInfo) | repeated | info is the signing info of all validators |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.slashing.v1.Query"></a>

### Query
Query provides defines the gRPC querier service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#lbm.slashing.v1.QueryParamsRequest) | [QueryParamsResponse](#lbm.slashing.v1.QueryParamsResponse) | Params queries the parameters of slashing module | GET|/lbm/slashing/v1/params|
| `SigningInfo` | [QuerySigningInfoRequest](#lbm.slashing.v1.QuerySigningInfoRequest) | [QuerySigningInfoResponse](#lbm.slashing.v1.QuerySigningInfoResponse) | SigningInfo queries the signing info of given cons address | GET|/lbm/slashing/v1/signing_infos/{cons_address}|
| `SigningInfos` | [QuerySigningInfosRequest](#lbm.slashing.v1.QuerySigningInfosRequest) | [QuerySigningInfosResponse](#lbm.slashing.v1.QuerySigningInfosResponse) | SigningInfos queries signing info of all validators | GET|/lbm/slashing/v1/signing_infos|

 <!-- end services -->



<a name="lbm/slashing/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/slashing/v1/tx.proto



<a name="lbm.slashing.v1.MsgUnjail"></a>

### MsgUnjail
MsgUnjail defines the Msg/Unjail request type


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_addr` | [string](#string) |  |  |






<a name="lbm.slashing.v1.MsgUnjailResponse"></a>

### MsgUnjailResponse
MsgUnjailResponse defines the Msg/Unjail response type





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.slashing.v1.Msg"></a>

### Msg
Msg defines the slashing Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Unjail` | [MsgUnjail](#lbm.slashing.v1.MsgUnjail) | [MsgUnjailResponse](#lbm.slashing.v1.MsgUnjailResponse) | Unjail defines a method for unjailing a jailed validator, thus returning them into the bonded validator set, so they can begin receiving provisions and rewards again. | |

 <!-- end services -->



<a name="lbm/staking/v1/staking.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/staking/v1/staking.proto



<a name="lbm.staking.v1.Commission"></a>

### Commission
Commission defines commission parameters for a given validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commission_rates` | [CommissionRates](#lbm.staking.v1.CommissionRates) |  | commission_rates defines the initial commission rates to be used for creating a validator. |
| `update_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | update_time is the last time the commission rate was changed. |






<a name="lbm.staking.v1.CommissionRates"></a>

### CommissionRates
CommissionRates defines the initial commission rates to be used for creating
a validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rate` | [string](#string) |  | rate is the commission rate charged to delegators, as a fraction. |
| `max_rate` | [string](#string) |  | max_rate defines the maximum commission rate which validator can ever charge, as a fraction. |
| `max_change_rate` | [string](#string) |  | max_change_rate defines the maximum daily increase of the validator commission, as a fraction. |






<a name="lbm.staking.v1.DVPair"></a>

### DVPair
DVPair is struct that just has a delegator-validator pair with no other data.
It is intended to be used as a marshalable pointer. For example, a DVPair can
be used to construct the key to getting an UnbondingDelegation from state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |






<a name="lbm.staking.v1.DVPairs"></a>

### DVPairs
DVPairs defines an array of DVPair objects.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pairs` | [DVPair](#lbm.staking.v1.DVPair) | repeated |  |






<a name="lbm.staking.v1.DVVTriplet"></a>

### DVVTriplet
DVVTriplet is struct that just has a delegator-validator-validator triplet
with no other data. It is intended to be used as a marshalable pointer. For
example, a DVVTriplet can be used to construct the key to getting a
Redelegation from state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_src_address` | [string](#string) |  |  |
| `validator_dst_address` | [string](#string) |  |  |






<a name="lbm.staking.v1.DVVTriplets"></a>

### DVVTriplets
DVVTriplets defines an array of DVVTriplet objects.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `triplets` | [DVVTriplet](#lbm.staking.v1.DVVTriplet) | repeated |  |






<a name="lbm.staking.v1.Delegation"></a>

### Delegation
Delegation represents the bond with tokens held by an account. It is
owned by one delegator, and is associated with the voting power of one
validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address is the bech32-encoded address of the delegator. |
| `validator_address` | [string](#string) |  | validator_address is the bech32-encoded address of the validator. |
| `shares` | [string](#string) |  | shares define the delegation shares received. |






<a name="lbm.staking.v1.DelegationResponse"></a>

### DelegationResponse
DelegationResponse is equivalent to Delegation except that it contains a
balance in addition to shares which is more suitable for client responses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegation` | [Delegation](#lbm.staking.v1.Delegation) |  |  |
| `balance` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) |  |  |






<a name="lbm.staking.v1.Description"></a>

### Description
Description defines a validator description.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `moniker` | [string](#string) |  | moniker defines a human-readable name for the validator. |
| `identity` | [string](#string) |  | identity defines an optional identity signature (ex. UPort or Keybase). |
| `website` | [string](#string) |  | website defines an optional website link. |
| `security_contact` | [string](#string) |  | security_contact defines an optional email for security contact. |
| `details` | [string](#string) |  | details define other optional details. |






<a name="lbm.staking.v1.HistoricalInfo"></a>

### HistoricalInfo
HistoricalInfo contains header and validator, voter information for a given block.
It is stored as part of staking module's state, which persists the `n` most
recent HistoricalInfo
(`n` is set by the staking module's `historical_entries` parameter).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `header` | [ostracon.types.Header](#ostracon.types.Header) |  |  |
| `valset` | [Validator](#lbm.staking.v1.Validator) | repeated |  |






<a name="lbm.staking.v1.Params"></a>

### Params
Params defines the parameters for the staking module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unbonding_time` | [google.protobuf.Duration](#google.protobuf.Duration) |  | unbonding_time is the time duration of unbonding. |
| `max_validators` | [uint32](#uint32) |  | max_validators is the maximum number of validators. |
| `max_entries` | [uint32](#uint32) |  | max_entries is the max entries for either unbonding delegation or redelegation (per pair/trio). |
| `historical_entries` | [uint32](#uint32) |  | historical_entries is the number of historical entries to persist. |
| `bond_denom` | [string](#string) |  | bond_denom defines the bondable coin denomination. |






<a name="lbm.staking.v1.Pool"></a>

### Pool
Pool is used for tracking bonded and not-bonded token supply of the bond
denomination.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `not_bonded_tokens` | [string](#string) |  |  |
| `bonded_tokens` | [string](#string) |  |  |






<a name="lbm.staking.v1.Redelegation"></a>

### Redelegation
Redelegation contains the list of a particular delegator's redelegating bonds
from a particular source validator to a particular destination validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address is the bech32-encoded address of the delegator. |
| `validator_src_address` | [string](#string) |  | validator_src_address is the validator redelegation source operator address. |
| `validator_dst_address` | [string](#string) |  | validator_dst_address is the validator redelegation destination operator address. |
| `entries` | [RedelegationEntry](#lbm.staking.v1.RedelegationEntry) | repeated | entries are the redelegation entries.

redelegation entries |






<a name="lbm.staking.v1.RedelegationEntry"></a>

### RedelegationEntry
RedelegationEntry defines a redelegation object with relevant metadata.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `creation_height` | [int64](#int64) |  | creation_height defines the height which the redelegation took place. |
| `completion_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | completion_time defines the unix time for redelegation completion. |
| `initial_balance` | [string](#string) |  | initial_balance defines the initial balance when redelegation started. |
| `shares_dst` | [string](#string) |  | shares_dst is the amount of destination-validator shares created by redelegation. |






<a name="lbm.staking.v1.RedelegationEntryResponse"></a>

### RedelegationEntryResponse
RedelegationEntryResponse is equivalent to a RedelegationEntry except that it
contains a balance in addition to shares which is more suitable for client
responses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `redelegation_entry` | [RedelegationEntry](#lbm.staking.v1.RedelegationEntry) |  |  |
| `balance` | [string](#string) |  |  |






<a name="lbm.staking.v1.RedelegationResponse"></a>

### RedelegationResponse
RedelegationResponse is equivalent to a Redelegation except that its entries
contain a balance in addition to shares which is more suitable for client
responses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `redelegation` | [Redelegation](#lbm.staking.v1.Redelegation) |  |  |
| `entries` | [RedelegationEntryResponse](#lbm.staking.v1.RedelegationEntryResponse) | repeated |  |






<a name="lbm.staking.v1.UnbondingDelegation"></a>

### UnbondingDelegation
UnbondingDelegation stores all of a single delegator's unbonding bonds
for a single validator in an time-ordered list.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address is the bech32-encoded address of the delegator. |
| `validator_address` | [string](#string) |  | validator_address is the bech32-encoded address of the validator. |
| `entries` | [UnbondingDelegationEntry](#lbm.staking.v1.UnbondingDelegationEntry) | repeated | entries are the unbonding delegation entries.

unbonding delegation entries |






<a name="lbm.staking.v1.UnbondingDelegationEntry"></a>

### UnbondingDelegationEntry
UnbondingDelegationEntry defines an unbonding object with relevant metadata.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `creation_height` | [int64](#int64) |  | creation_height is the height which the unbonding took place. |
| `completion_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | completion_time is the unix time for unbonding completion. |
| `initial_balance` | [string](#string) |  | initial_balance defines the tokens initially scheduled to receive at completion. |
| `balance` | [string](#string) |  | balance defines the tokens to receive at completion. |






<a name="lbm.staking.v1.ValAddresses"></a>

### ValAddresses
ValAddresses defines a repeated set of validator addresses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `addresses` | [string](#string) | repeated |  |






<a name="lbm.staking.v1.Validator"></a>

### Validator
Validator defines a validator, together with the total amount of the
Validator's bond shares and their exchange rate to coins. Slashing results in
a decrease in the exchange rate, allowing correct calculation of future
undelegations without iterating over delegators. When coins are delegated to
this validator, the validator is credited with a delegation whose number of
bond shares is based on the amount of coins delegated divided by the current
exchange rate. Voting power can be calculated as total bonded shares
multiplied by exchange rate.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `operator_address` | [string](#string) |  | operator_address defines the address of the validator's operator; bech encoded in JSON. |
| `consensus_pubkey` | [google.protobuf.Any](#google.protobuf.Any) |  | consensus_pubkey is the consensus public key of the validator, as a Protobuf Any. |
| `jailed` | [bool](#bool) |  | jailed defined whether the validator has been jailed from bonded status or not. |
| `status` | [BondStatus](#lbm.staking.v1.BondStatus) |  | status is the validator status (bonded/unbonding/unbonded). |
| `tokens` | [string](#string) |  | tokens define the delegated tokens (incl. self-delegation). |
| `delegator_shares` | [string](#string) |  | delegator_shares defines total shares issued to a validator's delegators. |
| `description` | [Description](#lbm.staking.v1.Description) |  | description defines the description terms for the validator. |
| `unbonding_height` | [int64](#int64) |  | unbonding_height defines, if unbonding, the height at which this validator has begun unbonding. |
| `unbonding_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | unbonding_time defines, if unbonding, the min time for the validator to complete unbonding. |
| `commission` | [Commission](#lbm.staking.v1.Commission) |  | commission defines the commission parameters. |
| `min_self_delegation` | [string](#string) |  | min_self_delegation is the validator's self declared minimum self delegation. |





 <!-- end messages -->


<a name="lbm.staking.v1.BondStatus"></a>

### BondStatus
BondStatus is the status of a validator.

| Name | Number | Description |
| ---- | ------ | ----------- |
| BOND_STATUS_UNSPECIFIED | 0 | UNSPECIFIED defines an invalid validator status. |
| BOND_STATUS_UNBONDED | 1 | UNBONDED defines a validator that is not bonded. |
| BOND_STATUS_UNBONDING | 2 | UNBONDING defines a validator that is unbonding. |
| BOND_STATUS_BONDED | 3 | BONDED defines a validator that is bonded. |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/staking/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/staking/v1/genesis.proto



<a name="lbm.staking.v1.GenesisState"></a>

### GenesisState
GenesisState defines the staking module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#lbm.staking.v1.Params) |  | params defines all the paramaters of related to deposit. |
| `last_total_power` | [bytes](#bytes) |  | last_total_power tracks the total amounts of bonded tokens recorded during the previous end block. |
| `last_validator_powers` | [LastValidatorPower](#lbm.staking.v1.LastValidatorPower) | repeated | last_validator_powers is a special index that provides a historical list of the last-block's bonded validators. |
| `validators` | [Validator](#lbm.staking.v1.Validator) | repeated | delegations defines the validator set at genesis. |
| `delegations` | [Delegation](#lbm.staking.v1.Delegation) | repeated | delegations defines the delegations active at genesis. |
| `unbonding_delegations` | [UnbondingDelegation](#lbm.staking.v1.UnbondingDelegation) | repeated | unbonding_delegations defines the unbonding delegations active at genesis. |
| `redelegations` | [Redelegation](#lbm.staking.v1.Redelegation) | repeated | redelegations defines the redelegations active at genesis. |
| `exported` | [bool](#bool) |  |  |






<a name="lbm.staking.v1.LastValidatorPower"></a>

### LastValidatorPower
LastValidatorPower required for validator set update logic.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address of the validator. |
| `power` | [int64](#int64) |  | power defines the power of the validator. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/staking/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/staking/v1/query.proto



<a name="lbm.staking.v1.QueryDelegationRequest"></a>

### QueryDelegationRequest
QueryDelegationRequest is request type for the Query/Delegation RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  | delegator_addr defines the delegator address to query for. |
| `validator_addr` | [string](#string) |  | validator_addr defines the validator address to query for. |






<a name="lbm.staking.v1.QueryDelegationResponse"></a>

### QueryDelegationResponse
QueryDelegationResponse is response type for the Query/Delegation RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegation_response` | [DelegationResponse](#lbm.staking.v1.DelegationResponse) |  | delegation_responses defines the delegation info of a delegation. |






<a name="lbm.staking.v1.QueryDelegatorDelegationsRequest"></a>

### QueryDelegatorDelegationsRequest
QueryDelegatorDelegationsRequest is request type for the
Query/DelegatorDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  | delegator_addr defines the delegator address to query for. |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.staking.v1.QueryDelegatorDelegationsResponse"></a>

### QueryDelegatorDelegationsResponse
QueryDelegatorDelegationsResponse is response type for the
Query/DelegatorDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegation_responses` | [DelegationResponse](#lbm.staking.v1.DelegationResponse) | repeated | delegation_responses defines all the delegations' info of a delegator. |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.staking.v1.QueryDelegatorUnbondingDelegationsRequest"></a>

### QueryDelegatorUnbondingDelegationsRequest
QueryDelegatorUnbondingDelegationsRequest is request type for the
Query/DelegatorUnbondingDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  | delegator_addr defines the delegator address to query for. |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.staking.v1.QueryDelegatorUnbondingDelegationsResponse"></a>

### QueryDelegatorUnbondingDelegationsResponse
QueryUnbondingDelegatorDelegationsResponse is response type for the
Query/UnbondingDelegatorDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unbonding_responses` | [UnbondingDelegation](#lbm.staking.v1.UnbondingDelegation) | repeated |  |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.staking.v1.QueryDelegatorValidatorRequest"></a>

### QueryDelegatorValidatorRequest
QueryDelegatorValidatorRequest is request type for the
Query/DelegatorValidator RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  | delegator_addr defines the delegator address to query for. |
| `validator_addr` | [string](#string) |  | validator_addr defines the validator address to query for. |






<a name="lbm.staking.v1.QueryDelegatorValidatorResponse"></a>

### QueryDelegatorValidatorResponse
QueryDelegatorValidatorResponse response type for the
Query/DelegatorValidator RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [Validator](#lbm.staking.v1.Validator) |  | validator defines the the validator info. |






<a name="lbm.staking.v1.QueryDelegatorValidatorsRequest"></a>

### QueryDelegatorValidatorsRequest
QueryDelegatorValidatorsRequest is request type for the
Query/DelegatorValidators RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  | delegator_addr defines the delegator address to query for. |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.staking.v1.QueryDelegatorValidatorsResponse"></a>

### QueryDelegatorValidatorsResponse
QueryDelegatorValidatorsResponse is response type for the
Query/DelegatorValidators RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validators` | [Validator](#lbm.staking.v1.Validator) | repeated | validators defines the the validators' info of a delegator. |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.staking.v1.QueryHistoricalInfoRequest"></a>

### QueryHistoricalInfoRequest
QueryHistoricalInfoRequest is request type for the Query/HistoricalInfo RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  | height defines at which height to query the historical info. |






<a name="lbm.staking.v1.QueryHistoricalInfoResponse"></a>

### QueryHistoricalInfoResponse
QueryHistoricalInfoResponse is response type for the Query/HistoricalInfo RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hist` | [HistoricalInfo](#lbm.staking.v1.HistoricalInfo) |  | hist defines the historical info at the given height. |






<a name="lbm.staking.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is request type for the Query/Params RPC method.






<a name="lbm.staking.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#lbm.staking.v1.Params) |  | params holds all the parameters of this module. |






<a name="lbm.staking.v1.QueryPoolRequest"></a>

### QueryPoolRequest
QueryPoolRequest is request type for the Query/Pool RPC method.






<a name="lbm.staking.v1.QueryPoolResponse"></a>

### QueryPoolResponse
QueryPoolResponse is response type for the Query/Pool RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pool` | [Pool](#lbm.staking.v1.Pool) |  | pool defines the pool info. |






<a name="lbm.staking.v1.QueryRedelegationsRequest"></a>

### QueryRedelegationsRequest
QueryRedelegationsRequest is request type for the Query/Redelegations RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  | delegator_addr defines the delegator address to query for. |
| `src_validator_addr` | [string](#string) |  | src_validator_addr defines the validator address to redelegate from. |
| `dst_validator_addr` | [string](#string) |  | dst_validator_addr defines the validator address to redelegate to. |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.staking.v1.QueryRedelegationsResponse"></a>

### QueryRedelegationsResponse
QueryRedelegationsResponse is response type for the Query/Redelegations RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `redelegation_responses` | [RedelegationResponse](#lbm.staking.v1.RedelegationResponse) | repeated |  |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.staking.v1.QueryUnbondingDelegationRequest"></a>

### QueryUnbondingDelegationRequest
QueryUnbondingDelegationRequest is request type for the
Query/UnbondingDelegation RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  | delegator_addr defines the delegator address to query for. |
| `validator_addr` | [string](#string) |  | validator_addr defines the validator address to query for. |






<a name="lbm.staking.v1.QueryUnbondingDelegationResponse"></a>

### QueryUnbondingDelegationResponse
QueryDelegationResponse is response type for the Query/UnbondingDelegation
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unbond` | [UnbondingDelegation](#lbm.staking.v1.UnbondingDelegation) |  | unbond defines the unbonding information of a delegation. |






<a name="lbm.staking.v1.QueryValidatorDelegationsRequest"></a>

### QueryValidatorDelegationsRequest
QueryValidatorDelegationsRequest is request type for the
Query/ValidatorDelegations RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_addr` | [string](#string) |  | validator_addr defines the validator address to query for. |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.staking.v1.QueryValidatorDelegationsResponse"></a>

### QueryValidatorDelegationsResponse
QueryValidatorDelegationsResponse is response type for the
Query/ValidatorDelegations RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegation_responses` | [DelegationResponse](#lbm.staking.v1.DelegationResponse) | repeated |  |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.staking.v1.QueryValidatorRequest"></a>

### QueryValidatorRequest
QueryValidatorRequest is response type for the Query/Validator RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_addr` | [string](#string) |  | validator_addr defines the validator address to query for. |






<a name="lbm.staking.v1.QueryValidatorResponse"></a>

### QueryValidatorResponse
QueryValidatorResponse is response type for the Query/Validator RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [Validator](#lbm.staking.v1.Validator) |  | validator defines the the validator info. |






<a name="lbm.staking.v1.QueryValidatorUnbondingDelegationsRequest"></a>

### QueryValidatorUnbondingDelegationsRequest
QueryValidatorUnbondingDelegationsRequest is required type for the
Query/ValidatorUnbondingDelegations RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_addr` | [string](#string) |  | validator_addr defines the validator address to query for. |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.staking.v1.QueryValidatorUnbondingDelegationsResponse"></a>

### QueryValidatorUnbondingDelegationsResponse
QueryValidatorUnbondingDelegationsResponse is response type for the
Query/ValidatorUnbondingDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unbonding_responses` | [UnbondingDelegation](#lbm.staking.v1.UnbondingDelegation) | repeated |  |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.staking.v1.QueryValidatorsRequest"></a>

### QueryValidatorsRequest
QueryValidatorsRequest is request type for Query/Validators RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `status` | [string](#string) |  | status enables to query for validators matching a given status. |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.staking.v1.QueryValidatorsResponse"></a>

### QueryValidatorsResponse
QueryValidatorsResponse is response type for the Query/Validators RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validators` | [Validator](#lbm.staking.v1.Validator) | repeated | validators contains all the queried validators. |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines the pagination in the response. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.staking.v1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Validators` | [QueryValidatorsRequest](#lbm.staking.v1.QueryValidatorsRequest) | [QueryValidatorsResponse](#lbm.staking.v1.QueryValidatorsResponse) | Validators queries all validators that match the given status. | GET|/lbm/staking/v1/validators|
| `Validator` | [QueryValidatorRequest](#lbm.staking.v1.QueryValidatorRequest) | [QueryValidatorResponse](#lbm.staking.v1.QueryValidatorResponse) | Validator queries validator info for given validator address. | GET|/lbm/staking/v1/validators/{validator_addr}|
| `ValidatorDelegations` | [QueryValidatorDelegationsRequest](#lbm.staking.v1.QueryValidatorDelegationsRequest) | [QueryValidatorDelegationsResponse](#lbm.staking.v1.QueryValidatorDelegationsResponse) | ValidatorDelegations queries delegate info for given validator. | GET|/lbm/staking/v1/validators/{validator_addr}/delegations|
| `ValidatorUnbondingDelegations` | [QueryValidatorUnbondingDelegationsRequest](#lbm.staking.v1.QueryValidatorUnbondingDelegationsRequest) | [QueryValidatorUnbondingDelegationsResponse](#lbm.staking.v1.QueryValidatorUnbondingDelegationsResponse) | ValidatorUnbondingDelegations queries unbonding delegations of a validator. | GET|/lbm/staking/v1/validators/{validator_addr}/unbonding_delegations|
| `Delegation` | [QueryDelegationRequest](#lbm.staking.v1.QueryDelegationRequest) | [QueryDelegationResponse](#lbm.staking.v1.QueryDelegationResponse) | Delegation queries delegate info for given validator delegator pair. | GET|/lbm/staking/v1/validators/{validator_addr}/delegations/{delegator_addr}|
| `UnbondingDelegation` | [QueryUnbondingDelegationRequest](#lbm.staking.v1.QueryUnbondingDelegationRequest) | [QueryUnbondingDelegationResponse](#lbm.staking.v1.QueryUnbondingDelegationResponse) | UnbondingDelegation queries unbonding info for given validator delegator pair. | GET|/lbm/staking/v1/validators/{validator_addr}/delegations/{delegator_addr}/unbonding_delegation|
| `DelegatorDelegations` | [QueryDelegatorDelegationsRequest](#lbm.staking.v1.QueryDelegatorDelegationsRequest) | [QueryDelegatorDelegationsResponse](#lbm.staking.v1.QueryDelegatorDelegationsResponse) | DelegatorDelegations queries all delegations of a given delegator address. | GET|/lbm/staking/v1/delegations/{delegator_addr}|
| `DelegatorUnbondingDelegations` | [QueryDelegatorUnbondingDelegationsRequest](#lbm.staking.v1.QueryDelegatorUnbondingDelegationsRequest) | [QueryDelegatorUnbondingDelegationsResponse](#lbm.staking.v1.QueryDelegatorUnbondingDelegationsResponse) | DelegatorUnbondingDelegations queries all unbonding delegations of a given delegator address. | GET|/lbm/staking/v1/delegators/{delegator_addr}/unbonding_delegations|
| `Redelegations` | [QueryRedelegationsRequest](#lbm.staking.v1.QueryRedelegationsRequest) | [QueryRedelegationsResponse](#lbm.staking.v1.QueryRedelegationsResponse) | Redelegations queries redelegations of given address. | GET|/lbm/staking/v1/delegators/{delegator_addr}/redelegations|
| `DelegatorValidators` | [QueryDelegatorValidatorsRequest](#lbm.staking.v1.QueryDelegatorValidatorsRequest) | [QueryDelegatorValidatorsResponse](#lbm.staking.v1.QueryDelegatorValidatorsResponse) | DelegatorValidators queries all validators info for given delegator address. | GET|/lbm/staking/v1/delegators/{delegator_addr}/validators|
| `DelegatorValidator` | [QueryDelegatorValidatorRequest](#lbm.staking.v1.QueryDelegatorValidatorRequest) | [QueryDelegatorValidatorResponse](#lbm.staking.v1.QueryDelegatorValidatorResponse) | DelegatorValidator queries validator info for given delegator validator pair. | GET|/lbm/staking/v1/delegators/{delegator_addr}/validators/{validator_addr}|
| `HistoricalInfo` | [QueryHistoricalInfoRequest](#lbm.staking.v1.QueryHistoricalInfoRequest) | [QueryHistoricalInfoResponse](#lbm.staking.v1.QueryHistoricalInfoResponse) | HistoricalInfo queries the historical info for given height. | GET|/lbm/staking/v1/historical_info/{height}|
| `Pool` | [QueryPoolRequest](#lbm.staking.v1.QueryPoolRequest) | [QueryPoolResponse](#lbm.staking.v1.QueryPoolResponse) | Pool queries the pool info. | GET|/lbm/staking/v1/pool|
| `Params` | [QueryParamsRequest](#lbm.staking.v1.QueryParamsRequest) | [QueryParamsResponse](#lbm.staking.v1.QueryParamsResponse) | Parameters queries the staking parameters. | GET|/lbm/staking/v1/params|

 <!-- end services -->



<a name="lbm/staking/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/staking/v1/tx.proto



<a name="lbm.staking.v1.MsgBeginRedelegate"></a>

### MsgBeginRedelegate
MsgBeginRedelegate defines a SDK message for performing a redelegation
of coins from a delegator and source validator to a destination validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_src_address` | [string](#string) |  |  |
| `validator_dst_address` | [string](#string) |  |  |
| `amount` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) |  |  |






<a name="lbm.staking.v1.MsgBeginRedelegateResponse"></a>

### MsgBeginRedelegateResponse
MsgBeginRedelegateResponse defines the Msg/BeginRedelegate response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `completion_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |






<a name="lbm.staking.v1.MsgCreateValidator"></a>

### MsgCreateValidator
MsgCreateValidator defines a SDK message for creating a new validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `description` | [Description](#lbm.staking.v1.Description) |  |  |
| `commission` | [CommissionRates](#lbm.staking.v1.CommissionRates) |  |  |
| `min_self_delegation` | [string](#string) |  |  |
| `delegator_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |
| `pubkey` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `value` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) |  |  |






<a name="lbm.staking.v1.MsgCreateValidatorResponse"></a>

### MsgCreateValidatorResponse
MsgCreateValidatorResponse defines the Msg/CreateValidator response type.






<a name="lbm.staking.v1.MsgDelegate"></a>

### MsgDelegate
MsgDelegate defines a SDK message for performing a delegation of coins
from a delegator to a validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |
| `amount` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) |  |  |






<a name="lbm.staking.v1.MsgDelegateResponse"></a>

### MsgDelegateResponse
MsgDelegateResponse defines the Msg/Delegate response type.






<a name="lbm.staking.v1.MsgEditValidator"></a>

### MsgEditValidator
MsgEditValidator defines a SDK message for editing an existing validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `description` | [Description](#lbm.staking.v1.Description) |  |  |
| `validator_address` | [string](#string) |  |  |
| `commission_rate` | [string](#string) |  | We pass a reference to the new commission rate and min self delegation as it's not mandatory to update. If not updated, the deserialized rate will be zero with no way to distinguish if an update was intended. REF: #2373 |
| `min_self_delegation` | [string](#string) |  |  |






<a name="lbm.staking.v1.MsgEditValidatorResponse"></a>

### MsgEditValidatorResponse
MsgEditValidatorResponse defines the Msg/EditValidator response type.






<a name="lbm.staking.v1.MsgUndelegate"></a>

### MsgUndelegate
MsgUndelegate defines a SDK message for performing an undelegation from a
delegate and a validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |
| `amount` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) |  |  |






<a name="lbm.staking.v1.MsgUndelegateResponse"></a>

### MsgUndelegateResponse
MsgUndelegateResponse defines the Msg/Undelegate response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `completion_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.staking.v1.Msg"></a>

### Msg
Msg defines the staking Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `CreateValidator` | [MsgCreateValidator](#lbm.staking.v1.MsgCreateValidator) | [MsgCreateValidatorResponse](#lbm.staking.v1.MsgCreateValidatorResponse) | CreateValidator defines a method for creating a new validator. | |
| `EditValidator` | [MsgEditValidator](#lbm.staking.v1.MsgEditValidator) | [MsgEditValidatorResponse](#lbm.staking.v1.MsgEditValidatorResponse) | EditValidator defines a method for editing an existing validator. | |
| `Delegate` | [MsgDelegate](#lbm.staking.v1.MsgDelegate) | [MsgDelegateResponse](#lbm.staking.v1.MsgDelegateResponse) | Delegate defines a method for performing a delegation of coins from a delegator to a validator. | |
| `BeginRedelegate` | [MsgBeginRedelegate](#lbm.staking.v1.MsgBeginRedelegate) | [MsgBeginRedelegateResponse](#lbm.staking.v1.MsgBeginRedelegateResponse) | BeginRedelegate defines a method for performing a redelegation of coins from a delegator and source validator to a destination validator. | |
| `Undelegate` | [MsgUndelegate](#lbm.staking.v1.MsgUndelegate) | [MsgUndelegateResponse](#lbm.staking.v1.MsgUndelegateResponse) | Undelegate defines a method for performing an undelegation from a delegate and a validator. | |

 <!-- end services -->



<a name="lbm/token/v1/event.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/token/v1/event.proto



<a name="lbm.token.v1.EventApprove"></a>

### EventApprove
EventApprove is emitted on Msg/Approve


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  | class id associated with the token. |
| `approver` | [string](#string) |  |  |
| `proxy` | [string](#string) |  |  |






<a name="lbm.token.v1.EventBurn"></a>

### EventBurn
EventBurn is emitted on Msg/Burn


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  | class id associated with the token. |
| `from` | [string](#string) |  |  |
| `amount` | [string](#string) |  |  |






<a name="lbm.token.v1.EventGrant"></a>

### EventGrant
EventGrant is emitted on Msg/Grant


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  | class id associated with the token. |
| `grantee` | [string](#string) |  | address of the granted account. |
| `action` | [string](#string) |  | action on the token class. Must be one of "mint", "burn" and "modify". |






<a name="lbm.token.v1.EventIssue"></a>

### EventIssue
EventIssue is emitted on Msg/Issue


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  | class id associated with the token. |






<a name="lbm.token.v1.EventMint"></a>

### EventMint
EventMint is emitted on Msg/Mint


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  | class id associated with the token. |
| `to` | [string](#string) |  |  |
| `amount` | [string](#string) |  |  |






<a name="lbm.token.v1.EventModify"></a>

### EventModify
EventModify is emitted on Msg/Modify


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  | class id associated with the token. |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |






<a name="lbm.token.v1.EventRevoke"></a>

### EventRevoke
EventRevoke is emitted on Msg/Revoke


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  | class id associated with the token. |
| `grantee` | [string](#string) |  | address of the revoked account. |
| `action` | [string](#string) |  | action on the token class. Must be one of "mint", "burn" and "modify". |






<a name="lbm.token.v1.EventTransfer"></a>

### EventTransfer
EventTransfer is emitted on Msg/Transfer and Msg/TransferFrom


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  | class id associated with the token. |
| `from` | [string](#string) |  |  |
| `to` | [string](#string) |  |  |
| `amount` | [string](#string) |  |  |
| `dummy` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/token/v1/token.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/token/v1/token.proto



<a name="lbm.token.v1.Approve"></a>

### Approve
Approve defines approve information.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `approver` | [string](#string) |  |  |
| `proxy` | [string](#string) |  |  |
| `class_id` | [string](#string) |  | class id associated with the token. |






<a name="lbm.token.v1.FT"></a>

### FT
FT defines a fungible token with a class id and an amount.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  | class id associated with the token. |
| `amount` | [string](#string) |  | amount of the token |






<a name="lbm.token.v1.Grant"></a>

### Grant
Grant defines grant information.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grantee` | [string](#string) |  | address of the granted account. |
| `class_id` | [string](#string) |  | class id associated with the token. |
| `action` | [string](#string) |  | action on the token class. Must be one of "mint", "burn" and "modify". |






<a name="lbm.token.v1.Pair"></a>

### Pair
Pair defines a key-value pair.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |






<a name="lbm.token.v1.Params"></a>

### Params
Params defines the parameters for the token module.






<a name="lbm.token.v1.Token"></a>

### Token
Token defines token information.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  | id defines the unique identifier of the token. |
| `name` | [string](#string) |  | name defines the human-readable name of the token. |
| `symbol` | [string](#string) |  | symbol is an abbreviated name for token. |
| `meta` | [string](#string) |  | meta is a brief description of token. |
| `image_uri` | [string](#string) |  | image_uri is an uri for the token image stored off chain. |
| `decimals` | [int32](#int32) |  | decimals is the number of decimals which one must divide the amount by to get its user representation. |
| `mintable` | [bool](#bool) |  | mintable represents whether the token is allowed to mint. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/token/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/token/v1/genesis.proto



<a name="lbm.token.v1.Balance"></a>

### Balance
Balance defines an account address and balance pair used in the token module's
genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |
| `tokens` | [FT](#lbm.token.v1.FT) | repeated |  |






<a name="lbm.token.v1.ClassGenesisState"></a>

### ClassGenesisState
ClassGenesisState defines the classs keeper's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [string](#string) |  | nonce is the next class nonce to issue. |
| `ids` | [string](#string) | repeated | ids represents the issued ids. |






<a name="lbm.token.v1.GenesisState"></a>

### GenesisState
GenesisState defines the token module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#lbm.token.v1.Params) |  | params defines all the paramaters of the module. |
| `class_state` | [ClassGenesisState](#lbm.token.v1.ClassGenesisState) |  | class_state is the class keeper's genesis state. |
| `balances` | [Balance](#lbm.token.v1.Balance) | repeated | balances is an array containing the balances of all the accounts. |
| `classes` | [Token](#lbm.token.v1.Token) | repeated | classes defines the metadata of the differents tokens. |
| `grants` | [Grant](#lbm.token.v1.Grant) | repeated | grants defines the grant information. |
| `approves` | [Approve](#lbm.token.v1.Approve) | repeated | approves defines the approve information. |
| `supplies` | [FT](#lbm.token.v1.FT) | repeated | supplies represents the total supplies of tokens. |
| `mints` | [FT](#lbm.token.v1.FT) | repeated | mints represents the total mints of tokens. |
| `burns` | [FT](#lbm.token.v1.FT) | repeated | burns represents the total burns of tokens. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/token/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/token/v1/query.proto



<a name="lbm.token.v1.QueryApproveRequest"></a>

### QueryApproveRequest
QueryApproveRequest is the request type for the Query/Approve RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  | class id associated with the token. |
| `proxy` | [string](#string) |  |  |
| `approver` | [string](#string) |  |  |






<a name="lbm.token.v1.QueryApproveResponse"></a>

### QueryApproveResponse
QueryApproveResponse is the response type for the Query/Approve RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `approve` | [Approve](#lbm.token.v1.Approve) |  |  |






<a name="lbm.token.v1.QueryApprovesRequest"></a>

### QueryApprovesRequest
QueryApprovesRequest is the request type for the Query/Approves RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  | class id associated with the token. |
| `proxy` | [string](#string) |  |  |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.token.v1.QueryApprovesResponse"></a>

### QueryApprovesResponse
QueryApprovesResponse is the response type for the Query/Approves RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `approves` | [Approve](#lbm.token.v1.Approve) | repeated |  |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  |  |






<a name="lbm.token.v1.QueryGrantsRequest"></a>

### QueryGrantsRequest
QueryGrantsRequest is the request type for the Query/Grants RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  | class id associated with the token. |
| `grantee` | [string](#string) |  |  |






<a name="lbm.token.v1.QueryGrantsResponse"></a>

### QueryGrantsResponse
QueryGrantsResponse is the response type for the Query/Grants RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grants` | [Grant](#lbm.token.v1.Grant) | repeated |  |






<a name="lbm.token.v1.QuerySupplyRequest"></a>

### QuerySupplyRequest
QuerySupplyRequest is the request type for the Query/Supply RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  | class id associated with the token. |
| `type` | [string](#string) |  |  |






<a name="lbm.token.v1.QuerySupplyResponse"></a>

### QuerySupplyResponse
QuerySupplyResponse is the response type for the Query/Supply RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `amount` | [string](#string) |  |  |






<a name="lbm.token.v1.QueryTokenBalanceRequest"></a>

### QueryTokenBalanceRequest
QueryTokenBalanceRequest is the request type for the Query/TokenBalance RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  | class id associated with the token. |
| `address` | [string](#string) |  |  |






<a name="lbm.token.v1.QueryTokenBalanceResponse"></a>

### QueryTokenBalanceResponse
QueryTokenBalanceResponse is the response type for the Query/TokenBalance RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `amount` | [string](#string) |  |  |






<a name="lbm.token.v1.QueryTokenRequest"></a>

### QueryTokenRequest
QueryTokenRequest is the request type for the Query/Token RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  |  |






<a name="lbm.token.v1.QueryTokenResponse"></a>

### QueryTokenResponse
QueryTokenResponse is the response type for the Query/Token RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token` | [Token](#lbm.token.v1.Token) |  |  |






<a name="lbm.token.v1.QueryTokensRequest"></a>

### QueryTokensRequest
QueryTokensRequest is the request type for the Query/Tokens RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.token.v1.QueryTokensResponse"></a>

### QueryTokensResponse
QueryTokensResponse is the response type for the Query/Tokens RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tokens` | [Token](#lbm.token.v1.Token) | repeated |  |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.token.v1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `TokenBalance` | [QueryTokenBalanceRequest](#lbm.token.v1.QueryTokenBalanceRequest) | [QueryTokenBalanceResponse](#lbm.token.v1.QueryTokenBalanceResponse) | TokenBalance queries the number of tokens of a given class owned by the address. | GET|/lbm/token/v1/balance/{address}/{class_id}|
| `Supply` | [QuerySupplyRequest](#lbm.token.v1.QuerySupplyRequest) | [QuerySupplyResponse](#lbm.token.v1.QuerySupplyResponse) | Supply queries the number of tokens from the given class id. | GET|/lbm/token/v1/supply/{class_id}|
| `Token` | [QueryTokenRequest](#lbm.token.v1.QueryTokenRequest) | [QueryTokenResponse](#lbm.token.v1.QueryTokenResponse) | Token queries an token metadata based on its class id. | GET|/lbm/token/v1/tokens/{class_id}|
| `Tokens` | [QueryTokensRequest](#lbm.token.v1.QueryTokensRequest) | [QueryTokensResponse](#lbm.token.v1.QueryTokensResponse) | Tokens queries all token metadata. | GET|/lbm/token/v1/tokens|
| `Grants` | [QueryGrantsRequest](#lbm.token.v1.QueryGrantsRequest) | [QueryGrantsResponse](#lbm.token.v1.QueryGrantsResponse) | Grants queries grants on a given grantee. | GET|/lbm/token/v1/grants/{grantee}/{class_id}|
| `Approve` | [QueryApproveRequest](#lbm.token.v1.QueryApproveRequest) | [QueryApproveResponse](#lbm.token.v1.QueryApproveResponse) | Approve queries approve on a given proxy approver pair. | GET|/lbm/token/v1/approve/{class_id}/{proxy}/{approver}|
| `Approves` | [QueryApprovesRequest](#lbm.token.v1.QueryApprovesRequest) | [QueryApprovesResponse](#lbm.token.v1.QueryApprovesResponse) | Approves queries all approves on a given proxy. | GET|/lbm/token/v1/approves/{class_id}/{proxy}|

 <!-- end services -->



<a name="lbm/token/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/token/v1/tx.proto



<a name="lbm.token.v1.MsgApprove"></a>

### MsgApprove
MsgApprove represents a message to transfer tokens on behalf of the approver


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  |  |
| `approver` | [string](#string) |  |  |
| `proxy` | [string](#string) |  |  |






<a name="lbm.token.v1.MsgApproveResponse"></a>

### MsgApproveResponse
MsgApproveResponse defines the Msg/Approve response type.






<a name="lbm.token.v1.MsgBurn"></a>

### MsgBurn
MsgBurn represents a message to burn tokens.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  |  |
| `from` | [string](#string) |  |  |
| `amount` | [string](#string) |  |  |






<a name="lbm.token.v1.MsgBurnFrom"></a>

### MsgBurnFrom
MsgBurnFrom represents a message to burn tokens.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  |  |
| `grantee` | [string](#string) |  |  |
| `from` | [string](#string) |  |  |
| `amount` | [string](#string) |  |  |






<a name="lbm.token.v1.MsgBurnFromResponse"></a>

### MsgBurnFromResponse
MsgBurnFromResponse defines the Msg/BurnFrom response type.






<a name="lbm.token.v1.MsgBurnResponse"></a>

### MsgBurnResponse
MsgBurnResponse defines the Msg/Burn response type.






<a name="lbm.token.v1.MsgGrant"></a>

### MsgGrant
MsgGrant represents a message to allow one to mint or burn tokens or modify a token metadata


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  |  |
| `granter` | [string](#string) |  |  |
| `grantee` | [string](#string) |  |  |
| `action` | [string](#string) |  |  |






<a name="lbm.token.v1.MsgGrantResponse"></a>

### MsgGrantResponse
MsgGrantResponse defines the Msg/Grant response type.






<a name="lbm.token.v1.MsgIssue"></a>

### MsgIssue
MsgIssue represents a message to issue a token.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `owner` | [string](#string) |  |  |
| `to` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `symbol` | [string](#string) |  |  |
| `image_uri` | [string](#string) |  |  |
| `meta` | [string](#string) |  |  |
| `amount` | [string](#string) |  |  |
| `mintable` | [bool](#bool) |  |  |
| `decimals` | [int32](#int32) |  |  |






<a name="lbm.token.v1.MsgIssueResponse"></a>

### MsgIssueResponse
MsgIssueResponse defines the Msg/Issue response type.






<a name="lbm.token.v1.MsgMint"></a>

### MsgMint
MsgMint represents a message to mint tokens.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  |  |
| `grantee` | [string](#string) |  |  |
| `to` | [string](#string) |  |  |
| `amount` | [string](#string) |  |  |






<a name="lbm.token.v1.MsgMintResponse"></a>

### MsgMintResponse
MsgMintResponse defines the Msg/Mint response type.






<a name="lbm.token.v1.MsgModify"></a>

### MsgModify
MsgModify represents a message to modify a token metadata.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  |  |
| `grantee` | [string](#string) |  |  |
| `changes` | [Pair](#lbm.token.v1.Pair) | repeated |  |






<a name="lbm.token.v1.MsgModifyResponse"></a>

### MsgModifyResponse
MsgModifyResponse defines the Msg/Modify response type.






<a name="lbm.token.v1.MsgRevoke"></a>

### MsgRevoke
MsgRevoke represents a message to revoke a grant.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  |  |
| `grantee` | [string](#string) |  |  |
| `action` | [string](#string) |  |  |






<a name="lbm.token.v1.MsgRevokeResponse"></a>

### MsgRevokeResponse
MsgRevokeResponse defines the Msg/Revoke response type.






<a name="lbm.token.v1.MsgTransfer"></a>

### MsgTransfer
MsgTransfer represents a message to transfer tokens from one account to another.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  |  |
| `from` | [string](#string) |  |  |
| `to` | [string](#string) |  |  |
| `amount` | [string](#string) |  |  |






<a name="lbm.token.v1.MsgTransferFrom"></a>

### MsgTransferFrom
MsgTransferFrom represents a message to transfer tokens from one account to another by the proxy.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  |  |
| `proxy` | [string](#string) |  |  |
| `from` | [string](#string) |  |  |
| `to` | [string](#string) |  |  |
| `amount` | [string](#string) |  |  |






<a name="lbm.token.v1.MsgTransferFromResponse"></a>

### MsgTransferFromResponse
MsgTransferFromResponse defines the Msg/TransferFrom response type.






<a name="lbm.token.v1.MsgTransferResponse"></a>

### MsgTransferResponse
MsgTransferResponse defines the Msg/Transfer response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.token.v1.Msg"></a>

### Msg
Msg defines the token Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Transfer` | [MsgTransfer](#lbm.token.v1.MsgTransfer) | [MsgTransferResponse](#lbm.token.v1.MsgTransferResponse) | Transfer defines a method to transfer tokens from one account to another account | |
| `TransferFrom` | [MsgTransferFrom](#lbm.token.v1.MsgTransferFrom) | [MsgTransferFromResponse](#lbm.token.v1.MsgTransferFromResponse) | TransferFrom defines a method to transfer tokens from one account to another account by the proxy | |
| `Approve` | [MsgApprove](#lbm.token.v1.MsgApprove) | [MsgApproveResponse](#lbm.token.v1.MsgApproveResponse) | Approve allows one to transfer tokens on behalf of the approver | |
| `Issue` | [MsgIssue](#lbm.token.v1.MsgIssue) | [MsgIssueResponse](#lbm.token.v1.MsgIssueResponse) | Issue defines a method to issue a token | |
| `Grant` | [MsgGrant](#lbm.token.v1.MsgGrant) | [MsgGrantResponse](#lbm.token.v1.MsgGrantResponse) | Grant allows one to mint or burn tokens or modify a token metadata | |
| `Revoke` | [MsgRevoke](#lbm.token.v1.MsgRevoke) | [MsgRevokeResponse](#lbm.token.v1.MsgRevokeResponse) | Revoke revokes the grant | |
| `Mint` | [MsgMint](#lbm.token.v1.MsgMint) | [MsgMintResponse](#lbm.token.v1.MsgMintResponse) | Mint defines a method to mint tokens | |
| `Burn` | [MsgBurn](#lbm.token.v1.MsgBurn) | [MsgBurnResponse](#lbm.token.v1.MsgBurnResponse) | Burn defines a method to burn tokens | |
| `BurnFrom` | [MsgBurnFrom](#lbm.token.v1.MsgBurnFrom) | [MsgBurnFromResponse](#lbm.token.v1.MsgBurnFromResponse) | BurnFrom defines a method to burn tokens | |
| `Modify` | [MsgModify](#lbm.token.v1.MsgModify) | [MsgModifyResponse](#lbm.token.v1.MsgModifyResponse) | Modify defines a method to modify a token metadata | |

 <!-- end services -->



<a name="lbm/tx/signing/v1/signing.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/tx/signing/v1/signing.proto



<a name="lbm.tx.signing.v1.SignatureDescriptor"></a>

### SignatureDescriptor
SignatureDescriptor is a convenience type which represents the full data for
a signature including the public key of the signer, signing modes and the
signature itself. It is primarily used for coordinating signatures between
clients.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `public_key` | [google.protobuf.Any](#google.protobuf.Any) |  | public_key is the public key of the signer |
| `data` | [SignatureDescriptor.Data](#lbm.tx.signing.v1.SignatureDescriptor.Data) |  |  |
| `sequence` | [uint64](#uint64) |  | sequence is the sequence of the account, which describes the number of committed transactions signed by a given address. It is used to prevent replay attacks. |






<a name="lbm.tx.signing.v1.SignatureDescriptor.Data"></a>

### SignatureDescriptor.Data
Data represents signature data


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `single` | [SignatureDescriptor.Data.Single](#lbm.tx.signing.v1.SignatureDescriptor.Data.Single) |  | single represents a single signer |
| `multi` | [SignatureDescriptor.Data.Multi](#lbm.tx.signing.v1.SignatureDescriptor.Data.Multi) |  | multi represents a multisig signer |






<a name="lbm.tx.signing.v1.SignatureDescriptor.Data.Multi"></a>

### SignatureDescriptor.Data.Multi
Multi is the signature data for a multisig public key


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bitarray` | [lbm.crypto.multisig.v1.CompactBitArray](#lbm.crypto.multisig.v1.CompactBitArray) |  | bitarray specifies which keys within the multisig are signing |
| `signatures` | [SignatureDescriptor.Data](#lbm.tx.signing.v1.SignatureDescriptor.Data) | repeated | signatures is the signatures of the multi-signature |






<a name="lbm.tx.signing.v1.SignatureDescriptor.Data.Single"></a>

### SignatureDescriptor.Data.Single
Single is the signature data for a single signer


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `mode` | [SignMode](#lbm.tx.signing.v1.SignMode) |  | mode is the signing mode of the single signer |
| `signature` | [bytes](#bytes) |  | signature is the raw signature bytes |






<a name="lbm.tx.signing.v1.SignatureDescriptors"></a>

### SignatureDescriptors
SignatureDescriptors wraps multiple SignatureDescriptor's.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signatures` | [SignatureDescriptor](#lbm.tx.signing.v1.SignatureDescriptor) | repeated | signatures are the signature descriptors |





 <!-- end messages -->


<a name="lbm.tx.signing.v1.SignMode"></a>

### SignMode
SignMode represents a signing mode with its own security guarantees.

| Name | Number | Description |
| ---- | ------ | ----------- |
| SIGN_MODE_UNSPECIFIED | 0 | SIGN_MODE_UNSPECIFIED specifies an unknown signing mode and will be rejected |
| SIGN_MODE_DIRECT | 1 | SIGN_MODE_DIRECT specifies a signing mode which uses SignDoc and is verified with raw bytes from Tx |
| SIGN_MODE_TEXTUAL | 2 | SIGN_MODE_TEXTUAL is a future signing mode that will verify some human-readable textual representation on top of the binary representation from SIGN_MODE_DIRECT |
| SIGN_MODE_LEGACY_AMINO_JSON | 127 | SIGN_MODE_LEGACY_AMINO_JSON is a backwards compatibility mode which uses Amino JSON and will be removed in the future |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/tx/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/tx/v1/tx.proto



<a name="lbm.tx.v1.AuthInfo"></a>

### AuthInfo
AuthInfo describes the fee and signer modes that are used to sign a
transaction.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signer_infos` | [SignerInfo](#lbm.tx.v1.SignerInfo) | repeated | signer_infos defines the signing modes for the required signers. The number and order of elements must match the required signers from TxBody's messages. The first element is the primary signer and the one which pays the fee. |
| `fee` | [Fee](#lbm.tx.v1.Fee) |  | Fee is the fee and gas limit for the transaction. The first signer is the primary signer and the one which pays the fee. The fee can be calculated based on the cost of evaluating the body and doing signature verification of the signers. This can be estimated via simulation. |






<a name="lbm.tx.v1.Fee"></a>

### Fee
Fee includes the amount of coins paid in fees and the maximum
gas to be used by the transaction. The ratio yields an effective "gasprice",
which must be above some miminum to be accepted into the mempool.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `amount` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated | amount is the amount of coins to be paid as a fee |
| `gas_limit` | [uint64](#uint64) |  | gas_limit is the maximum gas that can be used in transaction processing before an out of gas error occurs |
| `payer` | [string](#string) |  | if unset, the first signer is responsible for paying the fees. If set, the specified account must pay the fees. the payer must be a tx signer (and thus have signed this field in AuthInfo). setting this field does *not* change the ordering of required signers for the transaction. |
| `granter` | [string](#string) |  | if set, the fee payer (either the first signer or the value of the payer field) requests that a fee grant be used to pay fees instead of the fee payer's own balance. If an appropriate fee grant does not exist or the chain does not support fee grants, this will fail |






<a name="lbm.tx.v1.ModeInfo"></a>

### ModeInfo
ModeInfo describes the signing mode of a single or nested multisig signer.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `single` | [ModeInfo.Single](#lbm.tx.v1.ModeInfo.Single) |  | single represents a single signer |
| `multi` | [ModeInfo.Multi](#lbm.tx.v1.ModeInfo.Multi) |  | multi represents a nested multisig signer |






<a name="lbm.tx.v1.ModeInfo.Multi"></a>

### ModeInfo.Multi
Multi is the mode info for a multisig public key


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bitarray` | [lbm.crypto.multisig.v1.CompactBitArray](#lbm.crypto.multisig.v1.CompactBitArray) |  | bitarray specifies which keys within the multisig are signing |
| `mode_infos` | [ModeInfo](#lbm.tx.v1.ModeInfo) | repeated | mode_infos is the corresponding modes of the signers of the multisig which could include nested multisig public keys |






<a name="lbm.tx.v1.ModeInfo.Single"></a>

### ModeInfo.Single
Single is the mode info for a single signer. It is structured as a message
to allow for additional fields such as locale for SIGN_MODE_TEXTUAL in the
future


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `mode` | [lbm.tx.signing.v1.SignMode](#lbm.tx.signing.v1.SignMode) |  | mode is the signing mode of the single signer |






<a name="lbm.tx.v1.SignDoc"></a>

### SignDoc
SignDoc is the type used for generating sign bytes for SIGN_MODE_DIRECT.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `body_bytes` | [bytes](#bytes) |  | body_bytes is protobuf serialization of a TxBody that matches the representation in TxRaw. |
| `auth_info_bytes` | [bytes](#bytes) |  | auth_info_bytes is a protobuf serialization of an AuthInfo that matches the representation in TxRaw. |
| `chain_id` | [string](#string) |  | chain_id is the unique identifier of the chain this transaction targets. It prevents signed transactions from being used on another chain by an attacker |
| `account_number` | [uint64](#uint64) |  | account_number is the account number of the account in state |






<a name="lbm.tx.v1.SignerInfo"></a>

### SignerInfo
SignerInfo describes the public key and signing mode of a single top-level
signer.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `public_key` | [google.protobuf.Any](#google.protobuf.Any) |  | public_key is the public key of the signer. It is optional for accounts that already exist in state. If unset, the verifier can use the required \ signer address for this position and lookup the public key. |
| `mode_info` | [ModeInfo](#lbm.tx.v1.ModeInfo) |  | mode_info describes the signing mode of the signer and is a nested structure to support nested multisig pubkey's |
| `sequence` | [uint64](#uint64) |  | sequence is the sequence of the account, which describes the number of committed transactions signed by a given address. It is used to prevent replay attacks. |






<a name="lbm.tx.v1.Tx"></a>

### Tx
Tx is the standard type used for broadcasting transactions.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `body` | [TxBody](#lbm.tx.v1.TxBody) |  | body is the processable content of the transaction |
| `auth_info` | [AuthInfo](#lbm.tx.v1.AuthInfo) |  | auth_info is the authorization related content of the transaction, specifically signers, signer modes and fee |
| `signatures` | [bytes](#bytes) | repeated | signatures is a list of signatures that matches the length and order of AuthInfo's signer_infos to allow connecting signature meta information like public key and signing mode by position. |






<a name="lbm.tx.v1.TxBody"></a>

### TxBody
TxBody is the body of a transaction that all signers sign over.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `messages` | [google.protobuf.Any](#google.protobuf.Any) | repeated | messages is a list of messages to be executed. The required signers of those messages define the number and order of elements in AuthInfo's signer_infos and Tx's signatures. Each required signer address is added to the list only the first time it occurs. By convention, the first required signer (usually from the first message) is referred to as the primary signer and pays the fee for the whole transaction. |
| `memo` | [string](#string) |  | memo is any arbitrary memo to be added to the transaction |
| `timeout_height` | [uint64](#uint64) |  | timeout is the block height after which this transaction will not be processed by the chain |
| `extension_options` | [google.protobuf.Any](#google.protobuf.Any) | repeated | extension_options are arbitrary options that can be added by chains when the default options are not sufficient. If any of these are present and can't be handled, the transaction will be rejected |
| `non_critical_extension_options` | [google.protobuf.Any](#google.protobuf.Any) | repeated | extension_options are arbitrary options that can be added by chains when the default options are not sufficient. If any of these are present and can't be handled, they will be ignored |






<a name="lbm.tx.v1.TxRaw"></a>

### TxRaw
TxRaw is a variant of Tx that pins the signer's exact binary representation
of body and auth_info. This is used for signing, broadcasting and
verification. The binary `serialize(tx: TxRaw)` is stored in Tendermint and
the hash `sha256(serialize(tx: TxRaw))` becomes the "txhash", commonly used
as the transaction ID.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `body_bytes` | [bytes](#bytes) |  | body_bytes is a protobuf serialization of a TxBody that matches the representation in SignDoc. |
| `auth_info_bytes` | [bytes](#bytes) |  | auth_info_bytes is a protobuf serialization of an AuthInfo that matches the representation in SignDoc. |
| `signatures` | [bytes](#bytes) | repeated | signatures is a list of signatures that matches the length and order of AuthInfo's signer_infos to allow connecting signature meta information like public key and signing mode by position. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/tx/v1/service.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/tx/v1/service.proto



<a name="lbm.tx.v1.BroadcastTxRequest"></a>

### BroadcastTxRequest
BroadcastTxRequest is the request type for the Service.BroadcastTxRequest
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tx_bytes` | [bytes](#bytes) |  | tx_bytes is the raw transaction. |
| `mode` | [BroadcastMode](#lbm.tx.v1.BroadcastMode) |  |  |






<a name="lbm.tx.v1.BroadcastTxResponse"></a>

### BroadcastTxResponse
BroadcastTxResponse is the response type for the
Service.BroadcastTx method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tx_response` | [lbm.base.abci.v1.TxResponse](#lbm.base.abci.v1.TxResponse) |  | tx_response is the queried TxResponses. |






<a name="lbm.tx.v1.GetTxRequest"></a>

### GetTxRequest
GetTxRequest is the request type for the Service.GetTx
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [string](#string) |  | hash is the tx hash to query, encoded as a hex string. |






<a name="lbm.tx.v1.GetTxResponse"></a>

### GetTxResponse
GetTxResponse is the response type for the Service.GetTx method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tx` | [Tx](#lbm.tx.v1.Tx) |  | tx is the queried transaction. |
| `tx_response` | [lbm.base.abci.v1.TxResponse](#lbm.base.abci.v1.TxResponse) |  | tx_response is the queried TxResponses. |






<a name="lbm.tx.v1.GetTxsEventRequest"></a>

### GetTxsEventRequest
GetTxsEventRequest is the request type for the Service.TxsByEvents
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `events` | [string](#string) | repeated | events is the list of transaction event type. |
| `prove` | [bool](#bool) |  | prove is Include proofs of the transactions inclusion in the block |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an pagination for the request. |
| `order_by` | [OrderBy](#lbm.tx.v1.OrderBy) |  |  |






<a name="lbm.tx.v1.GetTxsEventResponse"></a>

### GetTxsEventResponse
GetTxsEventResponse is the response type for the Service.TxsByEvents
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `txs` | [Tx](#lbm.tx.v1.Tx) | repeated | txs is the list of queried transactions. |
| `tx_responses` | [lbm.base.abci.v1.TxResponse](#lbm.base.abci.v1.TxResponse) | repeated | tx_responses is the list of queried TxResponses. |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines an pagination for the response. |






<a name="lbm.tx.v1.SimulateRequest"></a>

### SimulateRequest
SimulateRequest is the request type for the Service.Simulate
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tx` | [Tx](#lbm.tx.v1.Tx) |  | tx is the transaction to simulate. |






<a name="lbm.tx.v1.SimulateResponse"></a>

### SimulateResponse
SimulateResponse is the response type for the
Service.SimulateRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gas_info` | [lbm.base.abci.v1.GasInfo](#lbm.base.abci.v1.GasInfo) |  | gas_info is the information about gas used in the simulation. |
| `result` | [lbm.base.abci.v1.Result](#lbm.base.abci.v1.Result) |  | result is the result of the simulation. |





 <!-- end messages -->


<a name="lbm.tx.v1.BroadcastMode"></a>

### BroadcastMode
BroadcastMode specifies the broadcast mode for the TxService.Broadcast RPC method.

| Name | Number | Description |
| ---- | ------ | ----------- |
| BROADCAST_MODE_UNSPECIFIED | 0 | zero-value for mode ordering |
| BROADCAST_MODE_BLOCK | 1 | BROADCAST_MODE_BLOCK defines a tx broadcasting mode where the client waits for the tx to be committed in a block. |
| BROADCAST_MODE_SYNC | 2 | BROADCAST_MODE_SYNC defines a tx broadcasting mode where the client waits for a CheckTx execution response only. |
| BROADCAST_MODE_ASYNC | 3 | BROADCAST_MODE_ASYNC defines a tx broadcasting mode where the client returns immediately. |



<a name="lbm.tx.v1.OrderBy"></a>

### OrderBy
OrderBy defines the sorting order

| Name | Number | Description |
| ---- | ------ | ----------- |
| ORDER_BY_UNSPECIFIED | 0 | ORDER_BY_UNSPECIFIED specifies an unknown sorting order. OrderBy defaults to ASC in this case. |
| ORDER_BY_ASC | 1 | ORDER_BY_ASC defines ascending order |
| ORDER_BY_DESC | 2 | ORDER_BY_DESC defines descending order |


 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.tx.v1.Service"></a>

### Service
Service defines a gRPC service for interacting with transactions.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Simulate` | [SimulateRequest](#lbm.tx.v1.SimulateRequest) | [SimulateResponse](#lbm.tx.v1.SimulateResponse) | Simulate simulates executing a transaction for estimating gas usage. | POST|/lbm/tx/v1/simulate|
| `GetTx` | [GetTxRequest](#lbm.tx.v1.GetTxRequest) | [GetTxResponse](#lbm.tx.v1.GetTxResponse) | GetTx fetches a tx by hash. | GET|/lbm/tx/v1/txs/{hash}|
| `BroadcastTx` | [BroadcastTxRequest](#lbm.tx.v1.BroadcastTxRequest) | [BroadcastTxResponse](#lbm.tx.v1.BroadcastTxResponse) | BroadcastTx broadcast transaction. | POST|/lbm/tx/v1/txs|
| `GetTxsEvent` | [GetTxsEventRequest](#lbm.tx.v1.GetTxsEventRequest) | [GetTxsEventResponse](#lbm.tx.v1.GetTxsEventResponse) | GetTxsEvent fetches txs by event. | GET|/lbm/tx/v1/txs|

 <!-- end services -->



<a name="lbm/upgrade/v1/upgrade.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/upgrade/v1/upgrade.proto



<a name="lbm.upgrade.v1.CancelSoftwareUpgradeProposal"></a>

### CancelSoftwareUpgradeProposal
CancelSoftwareUpgradeProposal is a gov Content type for cancelling a software
upgrade.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |






<a name="lbm.upgrade.v1.ModuleVersion"></a>

### ModuleVersion
ModuleVersion specifies a module and its consensus version.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | name of the app module |
| `version` | [uint64](#uint64) |  | consensus version of the app module |






<a name="lbm.upgrade.v1.Plan"></a>

### Plan
Plan specifies information about a planned upgrade and when it should occur.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | Sets the name for the upgrade. This name will be used by the upgraded version of the software to apply any special "on-upgrade" commands during the first BeginBlock method after the upgrade is applied. It is also used to detect whether a software version can handle a given upgrade. If no upgrade handler with this name has been set in the software, it will be assumed that the software is out-of-date when the upgrade Time or Height is reached and the software will exit. |
| `time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | The time after which the upgrade must be performed. Leave set to its zero value to use a pre-defined Height instead. |
| `height` | [int64](#int64) |  | The height at which the upgrade must be performed. Only used if Time is not set. |
| `info` | [string](#string) |  | Any application specific upgrade info to be included on-chain such as a git commit that validators could automatically upgrade to |
| `upgraded_client_state` | [google.protobuf.Any](#google.protobuf.Any) |  | IBC-enabled chains can opt-in to including the upgraded client state in its upgrade plan This will make the chain commit to the correct upgraded (self) client state before the upgrade occurs, so that connecting chains can verify that the new upgraded client is valid by verifying a proof on the previous version of the chain. This will allow IBC connections to persist smoothly across planned chain upgrades |






<a name="lbm.upgrade.v1.SoftwareUpgradeProposal"></a>

### SoftwareUpgradeProposal
SoftwareUpgradeProposal is a gov Content type for initiating a software
upgrade.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `plan` | [Plan](#lbm.upgrade.v1.Plan) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/upgrade/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/upgrade/v1/query.proto



<a name="lbm.upgrade.v1.QueryAppliedPlanRequest"></a>

### QueryAppliedPlanRequest
QueryCurrentPlanRequest is the request type for the Query/AppliedPlan RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | name is the name of the applied plan to query for. |






<a name="lbm.upgrade.v1.QueryAppliedPlanResponse"></a>

### QueryAppliedPlanResponse
QueryAppliedPlanResponse is the response type for the Query/AppliedPlan RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  | height is the block height at which the plan was applied. |






<a name="lbm.upgrade.v1.QueryCurrentPlanRequest"></a>

### QueryCurrentPlanRequest
QueryCurrentPlanRequest is the request type for the Query/CurrentPlan RPC
method.






<a name="lbm.upgrade.v1.QueryCurrentPlanResponse"></a>

### QueryCurrentPlanResponse
QueryCurrentPlanResponse is the response type for the Query/CurrentPlan RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `plan` | [Plan](#lbm.upgrade.v1.Plan) |  | plan is the current upgrade plan. |






<a name="lbm.upgrade.v1.QueryModuleVersionsRequest"></a>

### QueryModuleVersionsRequest
QueryModuleVersionsRequest is the request type for the Query/ModuleVersions
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `module_name` | [string](#string) |  | module_name is a field to query a specific module consensus version from state. Leaving this empty will fetch the full list of module versions from state |






<a name="lbm.upgrade.v1.QueryModuleVersionsResponse"></a>

### QueryModuleVersionsResponse
QueryModuleVersionsResponse is the response type for the Query/ModuleVersions
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `module_versions` | [ModuleVersion](#lbm.upgrade.v1.ModuleVersion) | repeated | module_versions is a list of module names with their consensus versions. |






<a name="lbm.upgrade.v1.QueryUpgradedConsensusStateRequest"></a>

### QueryUpgradedConsensusStateRequest
QueryUpgradedConsensusStateRequest is the request type for the Query/UpgradedConsensusState
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `last_height` | [int64](#int64) |  | last height of the current chain must be sent in request as this is the height under which next consensus state is stored |






<a name="lbm.upgrade.v1.QueryUpgradedConsensusStateResponse"></a>

### QueryUpgradedConsensusStateResponse
QueryUpgradedConsensusStateResponse is the response type for the Query/UpgradedConsensusState
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `upgraded_consensus_state` | [google.protobuf.Any](#google.protobuf.Any) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.upgrade.v1.Query"></a>

### Query
Query defines the gRPC upgrade querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `CurrentPlan` | [QueryCurrentPlanRequest](#lbm.upgrade.v1.QueryCurrentPlanRequest) | [QueryCurrentPlanResponse](#lbm.upgrade.v1.QueryCurrentPlanResponse) | CurrentPlan queries the current upgrade plan. | GET|/lbm/upgrade/v1/current_plan|
| `AppliedPlan` | [QueryAppliedPlanRequest](#lbm.upgrade.v1.QueryAppliedPlanRequest) | [QueryAppliedPlanResponse](#lbm.upgrade.v1.QueryAppliedPlanResponse) | AppliedPlan queries a previously applied upgrade plan by its name. | GET|/lbm/upgrade/v1/applied_plan/{name}|
| `UpgradedConsensusState` | [QueryUpgradedConsensusStateRequest](#lbm.upgrade.v1.QueryUpgradedConsensusStateRequest) | [QueryUpgradedConsensusStateResponse](#lbm.upgrade.v1.QueryUpgradedConsensusStateResponse) | UpgradedConsensusState queries the consensus state that will serve as a trusted kernel for the next version of this chain. It will only be stored at the last height of this chain. UpgradedConsensusState RPC not supported with legacy querier | GET|/lbm/upgrade/v1/upgraded_consensus_state/{last_height}|
| `ModuleVersions` | [QueryModuleVersionsRequest](#lbm.upgrade.v1.QueryModuleVersionsRequest) | [QueryModuleVersionsResponse](#lbm.upgrade.v1.QueryModuleVersionsResponse) | ModuleVersions queries the list of module versions from state. | GET|/lbm/upgrade/v1/module_versions|

 <!-- end services -->



<a name="lbm/vesting/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/vesting/v1/tx.proto



<a name="lbm.vesting.v1.MsgCreateVestingAccount"></a>

### MsgCreateVestingAccount
MsgCreateVestingAccount defines a message that enables creating a vesting
account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `from_address` | [string](#string) |  |  |
| `to_address` | [string](#string) |  |  |
| `amount` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated |  |
| `end_time` | [int64](#int64) |  |  |
| `delayed` | [bool](#bool) |  |  |






<a name="lbm.vesting.v1.MsgCreateVestingAccountResponse"></a>

### MsgCreateVestingAccountResponse
MsgCreateVestingAccountResponse defines the Msg/CreateVestingAccount response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.vesting.v1.Msg"></a>

### Msg
Msg defines the bank Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `CreateVestingAccount` | [MsgCreateVestingAccount](#lbm.vesting.v1.MsgCreateVestingAccount) | [MsgCreateVestingAccountResponse](#lbm.vesting.v1.MsgCreateVestingAccountResponse) | CreateVestingAccount defines a method that enables creating a vesting account. | |

 <!-- end services -->



<a name="lbm/vesting/v1/vesting.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/vesting/v1/vesting.proto



<a name="lbm.vesting.v1.BaseVestingAccount"></a>

### BaseVestingAccount
BaseVestingAccount implements the VestingAccount interface. It contains all
the necessary fields needed for any vesting account implementation.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_account` | [lbm.auth.v1.BaseAccount](#lbm.auth.v1.BaseAccount) |  |  |
| `original_vesting` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated |  |
| `delegated_free` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated |  |
| `delegated_vesting` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated |  |
| `end_time` | [int64](#int64) |  |  |






<a name="lbm.vesting.v1.ContinuousVestingAccount"></a>

### ContinuousVestingAccount
ContinuousVestingAccount implements the VestingAccount interface. It
continuously vests by unlocking coins linearly with respect to time.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_vesting_account` | [BaseVestingAccount](#lbm.vesting.v1.BaseVestingAccount) |  |  |
| `start_time` | [int64](#int64) |  |  |






<a name="lbm.vesting.v1.DelayedVestingAccount"></a>

### DelayedVestingAccount
DelayedVestingAccount implements the VestingAccount interface. It vests all
coins after a specific time, but non prior. In other words, it keeps them
locked until a specified time.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_vesting_account` | [BaseVestingAccount](#lbm.vesting.v1.BaseVestingAccount) |  |  |






<a name="lbm.vesting.v1.Period"></a>

### Period
Period defines a length of time and amount of coins that will vest.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `length` | [int64](#int64) |  |  |
| `amount` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated |  |






<a name="lbm.vesting.v1.PeriodicVestingAccount"></a>

### PeriodicVestingAccount
PeriodicVestingAccount implements the VestingAccount interface. It
periodically vests by unlocking coins during each specified period.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_vesting_account` | [BaseVestingAccount](#lbm.vesting.v1.BaseVestingAccount) |  |  |
| `start_time` | [int64](#int64) |  |  |
| `vesting_periods` | [Period](#lbm.vesting.v1.Period) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/wasm/v1/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/wasm/v1/types.proto



<a name="lbm.wasm.v1.AbsoluteTxPosition"></a>

### AbsoluteTxPosition
AbsoluteTxPosition is a unique transaction position that allows for global
ordering of transactions.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_height` | [uint64](#uint64) |  | BlockHeight is the block the contract was created at |
| `tx_index` | [uint64](#uint64) |  | TxIndex is a monotonic counter within the block (actual transaction index, or gas consumed) |






<a name="lbm.wasm.v1.AccessConfig"></a>

### AccessConfig
AccessConfig access control type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `permission` | [AccessType](#lbm.wasm.v1.AccessType) |  |  |
| `address` | [string](#string) |  |  |






<a name="lbm.wasm.v1.AccessTypeParam"></a>

### AccessTypeParam
AccessTypeParam


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `value` | [AccessType](#lbm.wasm.v1.AccessType) |  |  |






<a name="lbm.wasm.v1.CodeInfo"></a>

### CodeInfo
CodeInfo is data for the uploaded contract WASM code


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_hash` | [bytes](#bytes) |  | CodeHash is the unique identifier created by wasmvm |
| `creator` | [string](#string) |  | Creator address who initially stored the code |
| `instantiate_config` | [AccessConfig](#lbm.wasm.v1.AccessConfig) |  | InstantiateConfig access control to apply on contract creation, optional |






<a name="lbm.wasm.v1.ContractCodeHistoryEntry"></a>

### ContractCodeHistoryEntry
ContractCodeHistoryEntry metadata to a contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `operation` | [ContractCodeHistoryOperationType](#lbm.wasm.v1.ContractCodeHistoryOperationType) |  |  |
| `code_id` | [uint64](#uint64) |  | CodeID is the reference to the stored WASM code |
| `updated` | [AbsoluteTxPosition](#lbm.wasm.v1.AbsoluteTxPosition) |  | Updated Tx position when the operation was executed. |
| `msg` | [bytes](#bytes) |  |  |






<a name="lbm.wasm.v1.ContractInfo"></a>

### ContractInfo
ContractInfo stores a WASM contract instance


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_id` | [uint64](#uint64) |  | CodeID is the reference to the stored Wasm code |
| `creator` | [string](#string) |  | Creator address who initially instantiated the contract |
| `admin` | [string](#string) |  | Admin is an optional address that can execute migrations |
| `label` | [string](#string) |  | Label is optional metadata to be stored with a contract instance. |
| `created` | [AbsoluteTxPosition](#lbm.wasm.v1.AbsoluteTxPosition) |  | Created Tx position when the contract was instantiated. This data should kept internal and not be exposed via query results. Just use for sorting |
| `ibc_port_id` | [string](#string) |  |  |
| `status` | [ContractStatus](#lbm.wasm.v1.ContractStatus) |  | Status is a status of a contract |
| `extension` | [google.protobuf.Any](#google.protobuf.Any) |  | Extension is an extension point to store custom metadata within the persistence model. |






<a name="lbm.wasm.v1.Model"></a>

### Model
Model is a struct that holds a KV pair


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  | hex-encode key to read it better (this is often ascii) |
| `value` | [bytes](#bytes) |  | base64-encode raw value |






<a name="lbm.wasm.v1.Params"></a>

### Params
Params defines the set of wasm parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_upload_access` | [AccessConfig](#lbm.wasm.v1.AccessConfig) |  |  |
| `instantiate_default_permission` | [AccessType](#lbm.wasm.v1.AccessType) |  |  |
| `contract_status_access` | [AccessConfig](#lbm.wasm.v1.AccessConfig) |  |  |
| `max_wasm_code_size` | [uint64](#uint64) |  |  |
| `gas_multiplier` | [uint64](#uint64) |  |  |
| `instance_cost` | [uint64](#uint64) |  |  |
| `compile_cost` | [uint64](#uint64) |  |  |





 <!-- end messages -->


<a name="lbm.wasm.v1.AccessType"></a>

### AccessType
AccessType permission types

| Name | Number | Description |
| ---- | ------ | ----------- |
| ACCESS_TYPE_UNSPECIFIED | 0 | AccessTypeUnspecified placeholder for empty value |
| ACCESS_TYPE_NOBODY | 1 | AccessTypeNobody forbidden |
| ACCESS_TYPE_ONLY_ADDRESS | 2 | AccessTypeOnlyAddress restricted to an address |
| ACCESS_TYPE_EVERYBODY | 3 | AccessTypeEverybody unrestricted |



<a name="lbm.wasm.v1.ContractCodeHistoryOperationType"></a>

### ContractCodeHistoryOperationType
ContractCodeHistoryOperationType actions that caused a code change

| Name | Number | Description |
| ---- | ------ | ----------- |
| CONTRACT_CODE_HISTORY_OPERATION_TYPE_UNSPECIFIED | 0 | ContractCodeHistoryOperationTypeUnspecified placeholder for empty value |
| CONTRACT_CODE_HISTORY_OPERATION_TYPE_INIT | 1 | ContractCodeHistoryOperationTypeInit on chain contract instantiation |
| CONTRACT_CODE_HISTORY_OPERATION_TYPE_MIGRATE | 2 | ContractCodeHistoryOperationTypeMigrate code migration |
| CONTRACT_CODE_HISTORY_OPERATION_TYPE_GENESIS | 3 | ContractCodeHistoryOperationTypeGenesis based on genesis data |



<a name="lbm.wasm.v1.ContractStatus"></a>

### ContractStatus
ContractStatus types

| Name | Number | Description |
| ---- | ------ | ----------- |
| CONTRACT_STATUS_UNSPECIFIED | 0 | ContractStatus unspecified |
| CONTRACT_STATUS_ACTIVE | 1 | ContractStatus active |
| CONTRACT_STATUS_INACTIVE | 2 | ContractStatus inactive |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/wasm/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/wasm/v1/tx.proto



<a name="lbm.wasm.v1.MsgClearAdmin"></a>

### MsgClearAdmin
MsgClearAdmin removes any admin stored for a smart contract


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  | Sender is the that actor that signed the messages |
| `contract` | [string](#string) |  | Contract is the address of the smart contract |






<a name="lbm.wasm.v1.MsgClearAdminResponse"></a>

### MsgClearAdminResponse
MsgClearAdminResponse returns empty data






<a name="lbm.wasm.v1.MsgExecuteContract"></a>

### MsgExecuteContract
MsgExecuteContract submits the given message data to a smart contract


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  | Sender is the that actor that signed the messages |
| `contract` | [string](#string) |  | Contract is the address of the smart contract |
| `msg` | [bytes](#bytes) |  | Msg json encoded message to be passed to the contract |
| `funds` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated | Funds coins that are transferred to the contract on execution |






<a name="lbm.wasm.v1.MsgExecuteContractResponse"></a>

### MsgExecuteContractResponse
MsgExecuteContractResponse returns execution result data.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  | Data contains base64-encoded bytes to returned from the contract |






<a name="lbm.wasm.v1.MsgInstantiateContract"></a>

### MsgInstantiateContract
MsgInstantiateContract create a new smart contract instance for the given
code id.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  | Sender is the that actor that signed the messages |
| `admin` | [string](#string) |  | Admin is an optional address that can execute migrations |
| `code_id` | [uint64](#uint64) |  | CodeID is the reference to the stored WASM code |
| `label` | [string](#string) |  | Label is optional metadata to be stored with a contract instance. |
| `msg` | [bytes](#bytes) |  | Msg json encoded message to be passed to the contract on instantiation |
| `funds` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated | Funds coins that are transferred to the contract on instantiation |






<a name="lbm.wasm.v1.MsgInstantiateContractResponse"></a>

### MsgInstantiateContractResponse
MsgInstantiateContractResponse return instantiation result data


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | Address is the bech32 address of the new contract instance. |
| `data` | [bytes](#bytes) |  | Data contains base64-encoded bytes to returned from the contract |






<a name="lbm.wasm.v1.MsgMigrateContract"></a>

### MsgMigrateContract
MsgMigrateContract runs a code upgrade/ downgrade for a smart contract


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  | Sender is the that actor that signed the messages |
| `contract` | [string](#string) |  | Contract is the address of the smart contract |
| `code_id` | [uint64](#uint64) |  | CodeID references the new WASM code |
| `msg` | [bytes](#bytes) |  | Msg json encoded message to be passed to the contract on migration |






<a name="lbm.wasm.v1.MsgMigrateContractResponse"></a>

### MsgMigrateContractResponse
MsgMigrateContractResponse returns contract migration result data.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  | Data contains same raw bytes returned as data from the wasm contract. (May be empty) |






<a name="lbm.wasm.v1.MsgStoreCode"></a>

### MsgStoreCode
MsgStoreCode submit Wasm code to the system


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  | Sender is the that actor that signed the messages |
| `wasm_byte_code` | [bytes](#bytes) |  | WASMByteCode can be raw or gzip compressed |
| `instantiate_permission` | [AccessConfig](#lbm.wasm.v1.AccessConfig) |  | InstantiatePermission access control to apply on contract creation, optional |






<a name="lbm.wasm.v1.MsgStoreCodeAndInstantiateContract"></a>

### MsgStoreCodeAndInstantiateContract
MsgStoreCodeAndInstantiateContract submit Wasm code to the system and instantiate a contract using it.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  | Sender is the that actor that signed the messages |
| `wasm_byte_code` | [bytes](#bytes) |  | WASMByteCode can be raw or gzip compressed |
| `instantiate_permission` | [AccessConfig](#lbm.wasm.v1.AccessConfig) |  |  |
| `admin` | [string](#string) |  | Admin is an optional address that can execute migrations |
| `label` | [string](#string) |  | Label is optional metadata to be stored with a contract instance. |
| `init_msg` | [bytes](#bytes) |  | InitMsg json encoded message to be passed to the contract on instantiation |
| `funds` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated | Funds coins that are transferred to the contract on instantiation |






<a name="lbm.wasm.v1.MsgStoreCodeAndInstantiateContractResponse"></a>

### MsgStoreCodeAndInstantiateContractResponse
MsgStoreCodeAndInstantiateContractResponse returns store and instantiate result data.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_id` | [uint64](#uint64) |  | CodeID is the reference to the stored WASM code |
| `address` | [string](#string) |  | Address is the bech32 address of the new contract instance. |
| `data` | [bytes](#bytes) |  | Data contains base64-encoded bytes to returned from the contract |






<a name="lbm.wasm.v1.MsgStoreCodeResponse"></a>

### MsgStoreCodeResponse
MsgStoreCodeResponse returns store result data.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_id` | [uint64](#uint64) |  | CodeID is the reference to the stored WASM code |






<a name="lbm.wasm.v1.MsgUpdateAdmin"></a>

### MsgUpdateAdmin
MsgUpdateAdmin sets a new admin for a smart contract


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  | Sender is the that actor that signed the messages |
| `new_admin` | [string](#string) |  | NewAdmin address to be set |
| `contract` | [string](#string) |  | Contract is the address of the smart contract |






<a name="lbm.wasm.v1.MsgUpdateAdminResponse"></a>

### MsgUpdateAdminResponse
MsgUpdateAdminResponse returns empty data






<a name="lbm.wasm.v1.MsgUpdateContractStatus"></a>

### MsgUpdateContractStatus
MsgUpdateContractStatus sets a new status for a smart contract


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  | Sender is the that actor that signed the messages |
| `contract` | [string](#string) |  | Contract is the address of the smart contract |
| `status` | [ContractStatus](#lbm.wasm.v1.ContractStatus) |  | Status to be set |






<a name="lbm.wasm.v1.MsgUpdateContractStatusResponse"></a>

### MsgUpdateContractStatusResponse
MsgUpdateContractStatusResponse returns empty data





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.wasm.v1.Msg"></a>

### Msg
Msg defines the wasm Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `StoreCode` | [MsgStoreCode](#lbm.wasm.v1.MsgStoreCode) | [MsgStoreCodeResponse](#lbm.wasm.v1.MsgStoreCodeResponse) | StoreCode to submit Wasm code to the system | |
| `InstantiateContract` | [MsgInstantiateContract](#lbm.wasm.v1.MsgInstantiateContract) | [MsgInstantiateContractResponse](#lbm.wasm.v1.MsgInstantiateContractResponse) | Instantiate creates a new smart contract instance for the given code id. | |
| `StoreCodeAndInstantiateContract` | [MsgStoreCodeAndInstantiateContract](#lbm.wasm.v1.MsgStoreCodeAndInstantiateContract) | [MsgStoreCodeAndInstantiateContractResponse](#lbm.wasm.v1.MsgStoreCodeAndInstantiateContractResponse) | StoreCodeAndInstantiatecontract upload code and instantiate a contract using it. | |
| `ExecuteContract` | [MsgExecuteContract](#lbm.wasm.v1.MsgExecuteContract) | [MsgExecuteContractResponse](#lbm.wasm.v1.MsgExecuteContractResponse) | Execute submits the given message data to a smart contract | |
| `MigrateContract` | [MsgMigrateContract](#lbm.wasm.v1.MsgMigrateContract) | [MsgMigrateContractResponse](#lbm.wasm.v1.MsgMigrateContractResponse) | Migrate runs a code upgrade/ downgrade for a smart contract | |
| `UpdateAdmin` | [MsgUpdateAdmin](#lbm.wasm.v1.MsgUpdateAdmin) | [MsgUpdateAdminResponse](#lbm.wasm.v1.MsgUpdateAdminResponse) | UpdateAdmin sets a new admin for a smart contract | |
| `ClearAdmin` | [MsgClearAdmin](#lbm.wasm.v1.MsgClearAdmin) | [MsgClearAdminResponse](#lbm.wasm.v1.MsgClearAdminResponse) | ClearAdmin removes any admin stored for a smart contract | |
| `UpdateContractStatus` | [MsgUpdateContractStatus](#lbm.wasm.v1.MsgUpdateContractStatus) | [MsgUpdateContractStatusResponse](#lbm.wasm.v1.MsgUpdateContractStatusResponse) | UpdateContractStatus sets a new status for a smart contract | |

 <!-- end services -->



<a name="lbm/wasm/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/wasm/v1/genesis.proto



<a name="lbm.wasm.v1.Code"></a>

### Code
Code struct encompasses CodeInfo and CodeBytes


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_id` | [uint64](#uint64) |  |  |
| `code_info` | [CodeInfo](#lbm.wasm.v1.CodeInfo) |  |  |
| `code_bytes` | [bytes](#bytes) |  |  |
| `pinned` | [bool](#bool) |  | Pinned to wasmvm cache |






<a name="lbm.wasm.v1.Contract"></a>

### Contract
Contract struct encompasses ContractAddress, ContractInfo, and ContractState


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_address` | [string](#string) |  |  |
| `contract_info` | [ContractInfo](#lbm.wasm.v1.ContractInfo) |  |  |
| `contract_state` | [Model](#lbm.wasm.v1.Model) | repeated |  |






<a name="lbm.wasm.v1.GenesisState"></a>

### GenesisState
GenesisState - genesis state of x/wasm


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#lbm.wasm.v1.Params) |  |  |
| `codes` | [Code](#lbm.wasm.v1.Code) | repeated |  |
| `contracts` | [Contract](#lbm.wasm.v1.Contract) | repeated |  |
| `sequences` | [Sequence](#lbm.wasm.v1.Sequence) | repeated |  |
| `gen_msgs` | [GenesisState.GenMsgs](#lbm.wasm.v1.GenesisState.GenMsgs) | repeated |  |






<a name="lbm.wasm.v1.GenesisState.GenMsgs"></a>

### GenesisState.GenMsgs
GenMsgs define the messages that can be executed during genesis phase in order.
The intention is to have more human readable data that is auditable.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `store_code` | [MsgStoreCode](#lbm.wasm.v1.MsgStoreCode) |  |  |
| `instantiate_contract` | [MsgInstantiateContract](#lbm.wasm.v1.MsgInstantiateContract) |  |  |
| `execute_contract` | [MsgExecuteContract](#lbm.wasm.v1.MsgExecuteContract) |  |  |






<a name="lbm.wasm.v1.Sequence"></a>

### Sequence
Sequence key and value of an id generation counter


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id_key` | [bytes](#bytes) |  |  |
| `value` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/wasm/v1/ibc.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/wasm/v1/ibc.proto



<a name="lbm.wasm.v1.MsgIBCCloseChannel"></a>

### MsgIBCCloseChannel
MsgIBCCloseChannel port and channel need to be owned by the contract


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `channel` | [string](#string) |  |  |






<a name="lbm.wasm.v1.MsgIBCSend"></a>

### MsgIBCSend
MsgIBCSend


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `channel` | [string](#string) |  | the channel by which the packet will be sent |
| `timeout_height` | [uint64](#uint64) |  | Timeout height relative to the current block height. The timeout is disabled when set to 0. |
| `timeout_timestamp` | [uint64](#uint64) |  | Timeout timestamp (in nanoseconds) relative to the current block timestamp. The timeout is disabled when set to 0. |
| `data` | [bytes](#bytes) |  | data is the payload to transfer |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/wasm/v1/proposal.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/wasm/v1/proposal.proto



<a name="lbm.wasm.v1.ClearAdminProposal"></a>

### ClearAdminProposal
ClearAdminProposal gov proposal content type to clear the admin of a contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | Title is a short summary |
| `description` | [string](#string) |  | Description is a human readable text |
| `contract` | [string](#string) |  | Contract is the address of the smart contract |






<a name="lbm.wasm.v1.InstantiateContractProposal"></a>

### InstantiateContractProposal
InstantiateContractProposal gov proposal content type to instantiate a contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | Title is a short summary |
| `description` | [string](#string) |  | Description is a human readable text |
| `run_as` | [string](#string) |  | RunAs is the address that is passed to the contract's environment as sender |
| `admin` | [string](#string) |  | Admin is an optional address that can execute migrations |
| `code_id` | [uint64](#uint64) |  | CodeID is the reference to the stored WASM code |
| `label` | [string](#string) |  | Label is optional metadata to be stored with a constract instance. |
| `msg` | [bytes](#bytes) |  | Msg json encoded message to be passed to the contract on instantiation |
| `funds` | [lbm.base.v1.Coin](#lbm.base.v1.Coin) | repeated | Funds coins that are transferred to the contract on instantiation |






<a name="lbm.wasm.v1.MigrateContractProposal"></a>

### MigrateContractProposal
MigrateContractProposal gov proposal content type to migrate a contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | Title is a short summary |
| `description` | [string](#string) |  | Description is a human readable text |
| `run_as` | [string](#string) |  | RunAs is the address that is passed to the contract's environment as sender |
| `contract` | [string](#string) |  | Contract is the address of the smart contract |
| `code_id` | [uint64](#uint64) |  | CodeID references the new WASM code |
| `msg` | [bytes](#bytes) |  | Msg json encoded message to be passed to the contract on migration |






<a name="lbm.wasm.v1.PinCodesProposal"></a>

### PinCodesProposal
PinCodesProposal gov proposal content type to pin a set of code ids in the wasmvm cache.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | Title is a short summary |
| `description` | [string](#string) |  | Description is a human readable text |
| `code_ids` | [uint64](#uint64) | repeated | CodeIDs references the new WASM codes |






<a name="lbm.wasm.v1.StoreCodeProposal"></a>

### StoreCodeProposal
StoreCodeProposal gov proposal content type to submit WASM code to the system


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | Title is a short summary |
| `description` | [string](#string) |  | Description is a human readable text |
| `run_as` | [string](#string) |  | RunAs is the address that is passed to the contract's environment as sender |
| `wasm_byte_code` | [bytes](#bytes) |  | WASMByteCode can be raw or gzip compressed |
| `instantiate_permission` | [AccessConfig](#lbm.wasm.v1.AccessConfig) |  | InstantiatePermission to apply on contract creation, optional |






<a name="lbm.wasm.v1.UnpinCodesProposal"></a>

### UnpinCodesProposal
UnpinCodesProposal gov proposal content type to unpin a set of code ids in the wasmvm cache.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | Title is a short summary |
| `description` | [string](#string) |  | Description is a human readable text |
| `code_ids` | [uint64](#uint64) | repeated | CodeIDs references the WASM codes |






<a name="lbm.wasm.v1.UpdateAdminProposal"></a>

### UpdateAdminProposal
UpdateAdminProposal gov proposal content type to set an admin for a contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | Title is a short summary |
| `description` | [string](#string) |  | Description is a human readable text |
| `new_admin` | [string](#string) |  | NewAdmin address to be set |
| `contract` | [string](#string) |  | Contract is the address of the smart contract |






<a name="lbm.wasm.v1.UpdateContractStatusProposal"></a>

### UpdateContractStatusProposal
UpdateStatusProposal gov proposal content type to update the contract status.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | Title is a short summary |
| `description` | [string](#string) |  | Description is a human readable text |
| `contract` | [string](#string) |  | Contract is the address of the smart contract |
| `status` | [ContractStatus](#lbm.wasm.v1.ContractStatus) |  | Status to be set |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/wasm/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/wasm/v1/query.proto



<a name="lbm.wasm.v1.CodeInfoResponse"></a>

### CodeInfoResponse
CodeInfoResponse contains code meta data from CodeInfo


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_id` | [uint64](#uint64) |  | id for legacy support |
| `creator` | [string](#string) |  |  |
| `data_hash` | [bytes](#bytes) |  |  |
| `instantiate_permission` | [AccessConfig](#lbm.wasm.v1.AccessConfig) |  |  |






<a name="lbm.wasm.v1.QueryAllContractStateRequest"></a>

### QueryAllContractStateRequest
QueryAllContractStateRequest is the request type for the Query/AllContractState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address of the contract |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.wasm.v1.QueryAllContractStateResponse"></a>

### QueryAllContractStateResponse
QueryAllContractStateResponse is the response type for the
Query/AllContractState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `models` | [Model](#lbm.wasm.v1.Model) | repeated |  |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.wasm.v1.QueryCodeRequest"></a>

### QueryCodeRequest
QueryCodeRequest is the request type for the Query/Code RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_id` | [uint64](#uint64) |  | grpc-gateway_out does not support Go style CodID |






<a name="lbm.wasm.v1.QueryCodeResponse"></a>

### QueryCodeResponse
QueryCodeResponse is the response type for the Query/Code RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_info` | [CodeInfoResponse](#lbm.wasm.v1.CodeInfoResponse) |  |  |
| `data` | [bytes](#bytes) |  |  |






<a name="lbm.wasm.v1.QueryCodesRequest"></a>

### QueryCodesRequest
QueryCodesRequest is the request type for the Query/Codes RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.wasm.v1.QueryCodesResponse"></a>

### QueryCodesResponse
QueryCodesResponse is the response type for the Query/Codes RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_infos` | [CodeInfoResponse](#lbm.wasm.v1.CodeInfoResponse) | repeated |  |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.wasm.v1.QueryContractHistoryRequest"></a>

### QueryContractHistoryRequest
QueryContractHistoryRequest is the request type for the Query/ContractHistory RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address of the contract to query |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.wasm.v1.QueryContractHistoryResponse"></a>

### QueryContractHistoryResponse
QueryContractHistoryResponse is the response type for the Query/ContractHistory RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entries` | [ContractCodeHistoryEntry](#lbm.wasm.v1.ContractCodeHistoryEntry) | repeated |  |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.wasm.v1.QueryContractInfoRequest"></a>

### QueryContractInfoRequest
QueryContractInfoRequest is the request type for the Query/ContractInfo RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address of the contract to query |






<a name="lbm.wasm.v1.QueryContractInfoResponse"></a>

### QueryContractInfoResponse
QueryContractInfoResponse is the response type for the Query/ContractInfo RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address of the contract |
| `contract_info` | [ContractInfo](#lbm.wasm.v1.ContractInfo) |  |  |






<a name="lbm.wasm.v1.QueryContractsByCodeRequest"></a>

### QueryContractsByCodeRequest
QueryContractsByCodeRequest is the request type for the Query/ContractsByCode RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_id` | [uint64](#uint64) |  | grpc-gateway_out does not support Go style CodID |
| `pagination` | [lbm.base.query.v1.PageRequest](#lbm.base.query.v1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.wasm.v1.QueryContractsByCodeResponse"></a>

### QueryContractsByCodeResponse
QueryContractsByCodeResponse is the response type for the
Query/ContractsByCode RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contracts` | [string](#string) | repeated | contracts are a set of contract addresses |
| `pagination` | [lbm.base.query.v1.PageResponse](#lbm.base.query.v1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.wasm.v1.QueryRawContractStateRequest"></a>

### QueryRawContractStateRequest
QueryRawContractStateRequest is the request type for the
Query/RawContractState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address of the contract |
| `query_data` | [bytes](#bytes) |  |  |






<a name="lbm.wasm.v1.QueryRawContractStateResponse"></a>

### QueryRawContractStateResponse
QueryRawContractStateResponse is the response type for the
Query/RawContractState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  | Data contains the raw store data |






<a name="lbm.wasm.v1.QuerySmartContractStateRequest"></a>

### QuerySmartContractStateRequest
QuerySmartContractStateRequest is the request type for the
Query/SmartContractState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address of the contract |
| `query_data` | [bytes](#bytes) |  | QueryData contains the query data passed to the contract |






<a name="lbm.wasm.v1.QuerySmartContractStateResponse"></a>

### QuerySmartContractStateResponse
QuerySmartContractStateResponse is the response type for the
Query/SmartContractState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  | Data contains the json data returned from the smart contract |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.wasm.v1.Query"></a>

### Query
Query provides defines the gRPC querier service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `ContractInfo` | [QueryContractInfoRequest](#lbm.wasm.v1.QueryContractInfoRequest) | [QueryContractInfoResponse](#lbm.wasm.v1.QueryContractInfoResponse) | ContractInfo gets the contract meta data | GET|/lbm/wasm/v1/contract/{address}|
| `ContractHistory` | [QueryContractHistoryRequest](#lbm.wasm.v1.QueryContractHistoryRequest) | [QueryContractHistoryResponse](#lbm.wasm.v1.QueryContractHistoryResponse) | ContractHistory gets the contract code history | GET|/lbm/wasm/v1/contract/{address}/history|
| `ContractsByCode` | [QueryContractsByCodeRequest](#lbm.wasm.v1.QueryContractsByCodeRequest) | [QueryContractsByCodeResponse](#lbm.wasm.v1.QueryContractsByCodeResponse) | ContractsByCode lists all smart contracts for a code id | GET|/lbm/wasm/v1/code/{code_id}/contracts|
| `AllContractState` | [QueryAllContractStateRequest](#lbm.wasm.v1.QueryAllContractStateRequest) | [QueryAllContractStateResponse](#lbm.wasm.v1.QueryAllContractStateResponse) | AllContractState gets all raw store data for a single contract | GET|/lbm/wasm/v1/contract/{address}/state|
| `RawContractState` | [QueryRawContractStateRequest](#lbm.wasm.v1.QueryRawContractStateRequest) | [QueryRawContractStateResponse](#lbm.wasm.v1.QueryRawContractStateResponse) | RawContractState gets single key from the raw store data of a contract | GET|/lbm/wasm/v1/contract/{address}/raw/{query_data}|
| `SmartContractState` | [QuerySmartContractStateRequest](#lbm.wasm.v1.QuerySmartContractStateRequest) | [QuerySmartContractStateResponse](#lbm.wasm.v1.QuerySmartContractStateResponse) | SmartContractState get smart query result from the contract | GET|/lbm/wasm/v1/contract/{address}/smart/{query_data}|
| `Code` | [QueryCodeRequest](#lbm.wasm.v1.QueryCodeRequest) | [QueryCodeResponse](#lbm.wasm.v1.QueryCodeResponse) | Code gets the binary code and metadata for a singe wasm code | GET|/lbm/wasm/v1/code/{code_id}|
| `Codes` | [QueryCodesRequest](#lbm.wasm.v1.QueryCodesRequest) | [QueryCodesResponse](#lbm.wasm.v1.QueryCodesResponse) | Codes gets the metadata for all stored wasm codes | GET|/lbm/wasm/v1/code|

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

