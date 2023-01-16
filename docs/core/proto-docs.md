<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [cosmos/auth/v1beta1/auth.proto](#cosmos/auth/v1beta1/auth.proto)
    - [BaseAccount](#cosmos.auth.v1beta1.BaseAccount)
    - [ModuleAccount](#cosmos.auth.v1beta1.ModuleAccount)
    - [Params](#cosmos.auth.v1beta1.Params)
  
- [cosmos/auth/v1beta1/genesis.proto](#cosmos/auth/v1beta1/genesis.proto)
    - [GenesisState](#cosmos.auth.v1beta1.GenesisState)
  
- [cosmos/base/query/v1beta1/pagination.proto](#cosmos/base/query/v1beta1/pagination.proto)
    - [PageRequest](#cosmos.base.query.v1beta1.PageRequest)
    - [PageResponse](#cosmos.base.query.v1beta1.PageResponse)
  
- [cosmos/auth/v1beta1/query.proto](#cosmos/auth/v1beta1/query.proto)
    - [QueryAccountRequest](#cosmos.auth.v1beta1.QueryAccountRequest)
    - [QueryAccountResponse](#cosmos.auth.v1beta1.QueryAccountResponse)
    - [QueryAccountsRequest](#cosmos.auth.v1beta1.QueryAccountsRequest)
    - [QueryAccountsResponse](#cosmos.auth.v1beta1.QueryAccountsResponse)
    - [QueryModuleAccountByNameRequest](#cosmos.auth.v1beta1.QueryModuleAccountByNameRequest)
    - [QueryModuleAccountByNameResponse](#cosmos.auth.v1beta1.QueryModuleAccountByNameResponse)
    - [QueryParamsRequest](#cosmos.auth.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#cosmos.auth.v1beta1.QueryParamsResponse)
  
    - [Query](#cosmos.auth.v1beta1.Query)
  
- [cosmos/authz/v1beta1/authz.proto](#cosmos/authz/v1beta1/authz.proto)
    - [GenericAuthorization](#cosmos.authz.v1beta1.GenericAuthorization)
    - [Grant](#cosmos.authz.v1beta1.Grant)
    - [GrantAuthorization](#cosmos.authz.v1beta1.GrantAuthorization)
  
- [cosmos/authz/v1beta1/event.proto](#cosmos/authz/v1beta1/event.proto)
    - [EventGrant](#cosmos.authz.v1beta1.EventGrant)
    - [EventRevoke](#cosmos.authz.v1beta1.EventRevoke)
  
- [cosmos/authz/v1beta1/genesis.proto](#cosmos/authz/v1beta1/genesis.proto)
    - [GenesisState](#cosmos.authz.v1beta1.GenesisState)
  
- [cosmos/authz/v1beta1/query.proto](#cosmos/authz/v1beta1/query.proto)
    - [QueryGranteeGrantsRequest](#cosmos.authz.v1beta1.QueryGranteeGrantsRequest)
    - [QueryGranteeGrantsResponse](#cosmos.authz.v1beta1.QueryGranteeGrantsResponse)
    - [QueryGranterGrantsRequest](#cosmos.authz.v1beta1.QueryGranterGrantsRequest)
    - [QueryGranterGrantsResponse](#cosmos.authz.v1beta1.QueryGranterGrantsResponse)
    - [QueryGrantsRequest](#cosmos.authz.v1beta1.QueryGrantsRequest)
    - [QueryGrantsResponse](#cosmos.authz.v1beta1.QueryGrantsResponse)
  
    - [Query](#cosmos.authz.v1beta1.Query)
  
- [cosmos/base/abci/v1beta1/abci.proto](#cosmos/base/abci/v1beta1/abci.proto)
    - [ABCIMessageLog](#cosmos.base.abci.v1beta1.ABCIMessageLog)
    - [Attribute](#cosmos.base.abci.v1beta1.Attribute)
    - [GasInfo](#cosmos.base.abci.v1beta1.GasInfo)
    - [MsgData](#cosmos.base.abci.v1beta1.MsgData)
    - [Result](#cosmos.base.abci.v1beta1.Result)
    - [SearchTxsResult](#cosmos.base.abci.v1beta1.SearchTxsResult)
    - [SimulationResponse](#cosmos.base.abci.v1beta1.SimulationResponse)
    - [StringEvent](#cosmos.base.abci.v1beta1.StringEvent)
    - [TxMsgData](#cosmos.base.abci.v1beta1.TxMsgData)
    - [TxResponse](#cosmos.base.abci.v1beta1.TxResponse)
  
- [cosmos/authz/v1beta1/tx.proto](#cosmos/authz/v1beta1/tx.proto)
    - [MsgExec](#cosmos.authz.v1beta1.MsgExec)
    - [MsgExecResponse](#cosmos.authz.v1beta1.MsgExecResponse)
    - [MsgGrant](#cosmos.authz.v1beta1.MsgGrant)
    - [MsgGrantResponse](#cosmos.authz.v1beta1.MsgGrantResponse)
    - [MsgRevoke](#cosmos.authz.v1beta1.MsgRevoke)
    - [MsgRevokeResponse](#cosmos.authz.v1beta1.MsgRevokeResponse)
  
    - [Msg](#cosmos.authz.v1beta1.Msg)
  
- [cosmos/base/v1beta1/coin.proto](#cosmos/base/v1beta1/coin.proto)
    - [Coin](#cosmos.base.v1beta1.Coin)
    - [DecCoin](#cosmos.base.v1beta1.DecCoin)
    - [DecProto](#cosmos.base.v1beta1.DecProto)
    - [IntProto](#cosmos.base.v1beta1.IntProto)
  
- [cosmos/bank/v1beta1/authz.proto](#cosmos/bank/v1beta1/authz.proto)
    - [SendAuthorization](#cosmos.bank.v1beta1.SendAuthorization)
  
- [cosmos/bank/v1beta1/bank.proto](#cosmos/bank/v1beta1/bank.proto)
    - [DenomUnit](#cosmos.bank.v1beta1.DenomUnit)
    - [Input](#cosmos.bank.v1beta1.Input)
    - [Metadata](#cosmos.bank.v1beta1.Metadata)
    - [Output](#cosmos.bank.v1beta1.Output)
    - [Params](#cosmos.bank.v1beta1.Params)
    - [SendEnabled](#cosmos.bank.v1beta1.SendEnabled)
    - [Supply](#cosmos.bank.v1beta1.Supply)
  
- [cosmos/bank/v1beta1/genesis.proto](#cosmos/bank/v1beta1/genesis.proto)
    - [Balance](#cosmos.bank.v1beta1.Balance)
    - [GenesisState](#cosmos.bank.v1beta1.GenesisState)
  
- [cosmos/bank/v1beta1/query.proto](#cosmos/bank/v1beta1/query.proto)
    - [QueryAllBalancesRequest](#cosmos.bank.v1beta1.QueryAllBalancesRequest)
    - [QueryAllBalancesResponse](#cosmos.bank.v1beta1.QueryAllBalancesResponse)
    - [QueryBalanceRequest](#cosmos.bank.v1beta1.QueryBalanceRequest)
    - [QueryBalanceResponse](#cosmos.bank.v1beta1.QueryBalanceResponse)
    - [QueryDenomMetadataRequest](#cosmos.bank.v1beta1.QueryDenomMetadataRequest)
    - [QueryDenomMetadataResponse](#cosmos.bank.v1beta1.QueryDenomMetadataResponse)
    - [QueryDenomsMetadataRequest](#cosmos.bank.v1beta1.QueryDenomsMetadataRequest)
    - [QueryDenomsMetadataResponse](#cosmos.bank.v1beta1.QueryDenomsMetadataResponse)
    - [QueryParamsRequest](#cosmos.bank.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#cosmos.bank.v1beta1.QueryParamsResponse)
    - [QuerySpendableBalancesRequest](#cosmos.bank.v1beta1.QuerySpendableBalancesRequest)
    - [QuerySpendableBalancesResponse](#cosmos.bank.v1beta1.QuerySpendableBalancesResponse)
    - [QuerySupplyOfRequest](#cosmos.bank.v1beta1.QuerySupplyOfRequest)
    - [QuerySupplyOfResponse](#cosmos.bank.v1beta1.QuerySupplyOfResponse)
    - [QueryTotalSupplyRequest](#cosmos.bank.v1beta1.QueryTotalSupplyRequest)
    - [QueryTotalSupplyResponse](#cosmos.bank.v1beta1.QueryTotalSupplyResponse)
  
    - [Query](#cosmos.bank.v1beta1.Query)
  
- [cosmos/bank/v1beta1/tx.proto](#cosmos/bank/v1beta1/tx.proto)
    - [MsgMultiSend](#cosmos.bank.v1beta1.MsgMultiSend)
    - [MsgMultiSendResponse](#cosmos.bank.v1beta1.MsgMultiSendResponse)
    - [MsgSend](#cosmos.bank.v1beta1.MsgSend)
    - [MsgSendResponse](#cosmos.bank.v1beta1.MsgSendResponse)
  
    - [Msg](#cosmos.bank.v1beta1.Msg)
  
- [cosmos/base/kv/v1beta1/kv.proto](#cosmos/base/kv/v1beta1/kv.proto)
    - [Pair](#cosmos.base.kv.v1beta1.Pair)
    - [Pairs](#cosmos.base.kv.v1beta1.Pairs)
  
- [cosmos/base/node/v1beta1/query.proto](#cosmos/base/node/v1beta1/query.proto)
    - [ConfigRequest](#cosmos.base.node.v1beta1.ConfigRequest)
    - [ConfigResponse](#cosmos.base.node.v1beta1.ConfigResponse)
  
    - [Service](#cosmos.base.node.v1beta1.Service)
  
- [cosmos/base/reflection/v1beta1/reflection.proto](#cosmos/base/reflection/v1beta1/reflection.proto)
    - [ListAllInterfacesRequest](#cosmos.base.reflection.v1beta1.ListAllInterfacesRequest)
    - [ListAllInterfacesResponse](#cosmos.base.reflection.v1beta1.ListAllInterfacesResponse)
    - [ListImplementationsRequest](#cosmos.base.reflection.v1beta1.ListImplementationsRequest)
    - [ListImplementationsResponse](#cosmos.base.reflection.v1beta1.ListImplementationsResponse)
  
    - [ReflectionService](#cosmos.base.reflection.v1beta1.ReflectionService)
  
- [cosmos/base/reflection/v2alpha1/reflection.proto](#cosmos/base/reflection/v2alpha1/reflection.proto)
    - [AppDescriptor](#cosmos.base.reflection.v2alpha1.AppDescriptor)
    - [AuthnDescriptor](#cosmos.base.reflection.v2alpha1.AuthnDescriptor)
    - [ChainDescriptor](#cosmos.base.reflection.v2alpha1.ChainDescriptor)
    - [CodecDescriptor](#cosmos.base.reflection.v2alpha1.CodecDescriptor)
    - [ConfigurationDescriptor](#cosmos.base.reflection.v2alpha1.ConfigurationDescriptor)
    - [GetAuthnDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetAuthnDescriptorRequest)
    - [GetAuthnDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetAuthnDescriptorResponse)
    - [GetChainDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetChainDescriptorRequest)
    - [GetChainDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetChainDescriptorResponse)
    - [GetCodecDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetCodecDescriptorRequest)
    - [GetCodecDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetCodecDescriptorResponse)
    - [GetConfigurationDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetConfigurationDescriptorRequest)
    - [GetConfigurationDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetConfigurationDescriptorResponse)
    - [GetQueryServicesDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetQueryServicesDescriptorRequest)
    - [GetQueryServicesDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetQueryServicesDescriptorResponse)
    - [GetTxDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetTxDescriptorRequest)
    - [GetTxDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetTxDescriptorResponse)
    - [InterfaceAcceptingMessageDescriptor](#cosmos.base.reflection.v2alpha1.InterfaceAcceptingMessageDescriptor)
    - [InterfaceDescriptor](#cosmos.base.reflection.v2alpha1.InterfaceDescriptor)
    - [InterfaceImplementerDescriptor](#cosmos.base.reflection.v2alpha1.InterfaceImplementerDescriptor)
    - [MsgDescriptor](#cosmos.base.reflection.v2alpha1.MsgDescriptor)
    - [QueryMethodDescriptor](#cosmos.base.reflection.v2alpha1.QueryMethodDescriptor)
    - [QueryServiceDescriptor](#cosmos.base.reflection.v2alpha1.QueryServiceDescriptor)
    - [QueryServicesDescriptor](#cosmos.base.reflection.v2alpha1.QueryServicesDescriptor)
    - [SigningModeDescriptor](#cosmos.base.reflection.v2alpha1.SigningModeDescriptor)
    - [TxDescriptor](#cosmos.base.reflection.v2alpha1.TxDescriptor)
  
    - [ReflectionService](#cosmos.base.reflection.v2alpha1.ReflectionService)
  
- [cosmos/base/snapshots/v1beta1/snapshot.proto](#cosmos/base/snapshots/v1beta1/snapshot.proto)
    - [Metadata](#cosmos.base.snapshots.v1beta1.Metadata)
    - [Snapshot](#cosmos.base.snapshots.v1beta1.Snapshot)
    - [SnapshotExtensionMeta](#cosmos.base.snapshots.v1beta1.SnapshotExtensionMeta)
    - [SnapshotExtensionPayload](#cosmos.base.snapshots.v1beta1.SnapshotExtensionPayload)
    - [SnapshotIAVLItem](#cosmos.base.snapshots.v1beta1.SnapshotIAVLItem)
    - [SnapshotItem](#cosmos.base.snapshots.v1beta1.SnapshotItem)
    - [SnapshotStoreItem](#cosmos.base.snapshots.v1beta1.SnapshotStoreItem)
  
- [cosmos/base/store/v1beta1/commit_info.proto](#cosmos/base/store/v1beta1/commit_info.proto)
    - [CommitID](#cosmos.base.store.v1beta1.CommitID)
    - [CommitInfo](#cosmos.base.store.v1beta1.CommitInfo)
    - [StoreInfo](#cosmos.base.store.v1beta1.StoreInfo)
  
- [cosmos/base/store/v1beta1/listening.proto](#cosmos/base/store/v1beta1/listening.proto)
    - [StoreKVPair](#cosmos.base.store.v1beta1.StoreKVPair)
  
- [cosmos/capability/v1beta1/capability.proto](#cosmos/capability/v1beta1/capability.proto)
    - [Capability](#cosmos.capability.v1beta1.Capability)
    - [CapabilityOwners](#cosmos.capability.v1beta1.CapabilityOwners)
    - [Owner](#cosmos.capability.v1beta1.Owner)
  
- [cosmos/capability/v1beta1/genesis.proto](#cosmos/capability/v1beta1/genesis.proto)
    - [GenesisOwners](#cosmos.capability.v1beta1.GenesisOwners)
    - [GenesisState](#cosmos.capability.v1beta1.GenesisState)
  
- [cosmos/crisis/v1beta1/genesis.proto](#cosmos/crisis/v1beta1/genesis.proto)
    - [GenesisState](#cosmos.crisis.v1beta1.GenesisState)
  
- [cosmos/crisis/v1beta1/tx.proto](#cosmos/crisis/v1beta1/tx.proto)
    - [MsgVerifyInvariant](#cosmos.crisis.v1beta1.MsgVerifyInvariant)
    - [MsgVerifyInvariantResponse](#cosmos.crisis.v1beta1.MsgVerifyInvariantResponse)
  
    - [Msg](#cosmos.crisis.v1beta1.Msg)
  
- [cosmos/crypto/ed25519/keys.proto](#cosmos/crypto/ed25519/keys.proto)
    - [PrivKey](#cosmos.crypto.ed25519.PrivKey)
    - [PubKey](#cosmos.crypto.ed25519.PubKey)
  
- [cosmos/crypto/multisig/keys.proto](#cosmos/crypto/multisig/keys.proto)
    - [LegacyAminoPubKey](#cosmos.crypto.multisig.LegacyAminoPubKey)
  
- [cosmos/crypto/multisig/v1beta1/multisig.proto](#cosmos/crypto/multisig/v1beta1/multisig.proto)
    - [CompactBitArray](#cosmos.crypto.multisig.v1beta1.CompactBitArray)
    - [MultiSignature](#cosmos.crypto.multisig.v1beta1.MultiSignature)
  
- [cosmos/crypto/secp256k1/keys.proto](#cosmos/crypto/secp256k1/keys.proto)
    - [PrivKey](#cosmos.crypto.secp256k1.PrivKey)
    - [PubKey](#cosmos.crypto.secp256k1.PubKey)
  
- [cosmos/crypto/secp256r1/keys.proto](#cosmos/crypto/secp256r1/keys.proto)
    - [PrivKey](#cosmos.crypto.secp256r1.PrivKey)
    - [PubKey](#cosmos.crypto.secp256r1.PubKey)
  
- [cosmos/distribution/v1beta1/distribution.proto](#cosmos/distribution/v1beta1/distribution.proto)
    - [CommunityPoolSpendProposal](#cosmos.distribution.v1beta1.CommunityPoolSpendProposal)
    - [CommunityPoolSpendProposalWithDeposit](#cosmos.distribution.v1beta1.CommunityPoolSpendProposalWithDeposit)
    - [DelegationDelegatorReward](#cosmos.distribution.v1beta1.DelegationDelegatorReward)
    - [DelegatorStartingInfo](#cosmos.distribution.v1beta1.DelegatorStartingInfo)
    - [FeePool](#cosmos.distribution.v1beta1.FeePool)
    - [Params](#cosmos.distribution.v1beta1.Params)
    - [ValidatorAccumulatedCommission](#cosmos.distribution.v1beta1.ValidatorAccumulatedCommission)
    - [ValidatorCurrentRewards](#cosmos.distribution.v1beta1.ValidatorCurrentRewards)
    - [ValidatorHistoricalRewards](#cosmos.distribution.v1beta1.ValidatorHistoricalRewards)
    - [ValidatorOutstandingRewards](#cosmos.distribution.v1beta1.ValidatorOutstandingRewards)
    - [ValidatorSlashEvent](#cosmos.distribution.v1beta1.ValidatorSlashEvent)
    - [ValidatorSlashEvents](#cosmos.distribution.v1beta1.ValidatorSlashEvents)
  
- [cosmos/distribution/v1beta1/genesis.proto](#cosmos/distribution/v1beta1/genesis.proto)
    - [DelegatorStartingInfoRecord](#cosmos.distribution.v1beta1.DelegatorStartingInfoRecord)
    - [DelegatorWithdrawInfo](#cosmos.distribution.v1beta1.DelegatorWithdrawInfo)
    - [GenesisState](#cosmos.distribution.v1beta1.GenesisState)
    - [ValidatorAccumulatedCommissionRecord](#cosmos.distribution.v1beta1.ValidatorAccumulatedCommissionRecord)
    - [ValidatorCurrentRewardsRecord](#cosmos.distribution.v1beta1.ValidatorCurrentRewardsRecord)
    - [ValidatorHistoricalRewardsRecord](#cosmos.distribution.v1beta1.ValidatorHistoricalRewardsRecord)
    - [ValidatorOutstandingRewardsRecord](#cosmos.distribution.v1beta1.ValidatorOutstandingRewardsRecord)
    - [ValidatorSlashEventRecord](#cosmos.distribution.v1beta1.ValidatorSlashEventRecord)
  
- [cosmos/distribution/v1beta1/query.proto](#cosmos/distribution/v1beta1/query.proto)
    - [QueryCommunityPoolRequest](#cosmos.distribution.v1beta1.QueryCommunityPoolRequest)
    - [QueryCommunityPoolResponse](#cosmos.distribution.v1beta1.QueryCommunityPoolResponse)
    - [QueryDelegationRewardsRequest](#cosmos.distribution.v1beta1.QueryDelegationRewardsRequest)
    - [QueryDelegationRewardsResponse](#cosmos.distribution.v1beta1.QueryDelegationRewardsResponse)
    - [QueryDelegationTotalRewardsRequest](#cosmos.distribution.v1beta1.QueryDelegationTotalRewardsRequest)
    - [QueryDelegationTotalRewardsResponse](#cosmos.distribution.v1beta1.QueryDelegationTotalRewardsResponse)
    - [QueryDelegatorValidatorsRequest](#cosmos.distribution.v1beta1.QueryDelegatorValidatorsRequest)
    - [QueryDelegatorValidatorsResponse](#cosmos.distribution.v1beta1.QueryDelegatorValidatorsResponse)
    - [QueryDelegatorWithdrawAddressRequest](#cosmos.distribution.v1beta1.QueryDelegatorWithdrawAddressRequest)
    - [QueryDelegatorWithdrawAddressResponse](#cosmos.distribution.v1beta1.QueryDelegatorWithdrawAddressResponse)
    - [QueryParamsRequest](#cosmos.distribution.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#cosmos.distribution.v1beta1.QueryParamsResponse)
    - [QueryValidatorCommissionRequest](#cosmos.distribution.v1beta1.QueryValidatorCommissionRequest)
    - [QueryValidatorCommissionResponse](#cosmos.distribution.v1beta1.QueryValidatorCommissionResponse)
    - [QueryValidatorOutstandingRewardsRequest](#cosmos.distribution.v1beta1.QueryValidatorOutstandingRewardsRequest)
    - [QueryValidatorOutstandingRewardsResponse](#cosmos.distribution.v1beta1.QueryValidatorOutstandingRewardsResponse)
    - [QueryValidatorSlashesRequest](#cosmos.distribution.v1beta1.QueryValidatorSlashesRequest)
    - [QueryValidatorSlashesResponse](#cosmos.distribution.v1beta1.QueryValidatorSlashesResponse)
  
    - [Query](#cosmos.distribution.v1beta1.Query)
  
- [cosmos/distribution/v1beta1/tx.proto](#cosmos/distribution/v1beta1/tx.proto)
    - [MsgFundCommunityPool](#cosmos.distribution.v1beta1.MsgFundCommunityPool)
    - [MsgFundCommunityPoolResponse](#cosmos.distribution.v1beta1.MsgFundCommunityPoolResponse)
    - [MsgSetWithdrawAddress](#cosmos.distribution.v1beta1.MsgSetWithdrawAddress)
    - [MsgSetWithdrawAddressResponse](#cosmos.distribution.v1beta1.MsgSetWithdrawAddressResponse)
    - [MsgWithdrawDelegatorReward](#cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward)
    - [MsgWithdrawDelegatorRewardResponse](#cosmos.distribution.v1beta1.MsgWithdrawDelegatorRewardResponse)
    - [MsgWithdrawValidatorCommission](#cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission)
    - [MsgWithdrawValidatorCommissionResponse](#cosmos.distribution.v1beta1.MsgWithdrawValidatorCommissionResponse)
  
    - [Msg](#cosmos.distribution.v1beta1.Msg)
  
- [cosmos/evidence/v1beta1/evidence.proto](#cosmos/evidence/v1beta1/evidence.proto)
    - [Equivocation](#cosmos.evidence.v1beta1.Equivocation)
  
- [cosmos/evidence/v1beta1/genesis.proto](#cosmos/evidence/v1beta1/genesis.proto)
    - [GenesisState](#cosmos.evidence.v1beta1.GenesisState)
  
- [cosmos/evidence/v1beta1/query.proto](#cosmos/evidence/v1beta1/query.proto)
    - [QueryAllEvidenceRequest](#cosmos.evidence.v1beta1.QueryAllEvidenceRequest)
    - [QueryAllEvidenceResponse](#cosmos.evidence.v1beta1.QueryAllEvidenceResponse)
    - [QueryEvidenceRequest](#cosmos.evidence.v1beta1.QueryEvidenceRequest)
    - [QueryEvidenceResponse](#cosmos.evidence.v1beta1.QueryEvidenceResponse)
  
    - [Query](#cosmos.evidence.v1beta1.Query)
  
- [cosmos/evidence/v1beta1/tx.proto](#cosmos/evidence/v1beta1/tx.proto)
    - [MsgSubmitEvidence](#cosmos.evidence.v1beta1.MsgSubmitEvidence)
    - [MsgSubmitEvidenceResponse](#cosmos.evidence.v1beta1.MsgSubmitEvidenceResponse)
  
    - [Msg](#cosmos.evidence.v1beta1.Msg)
  
- [cosmos/feegrant/v1beta1/feegrant.proto](#cosmos/feegrant/v1beta1/feegrant.proto)
    - [AllowedMsgAllowance](#cosmos.feegrant.v1beta1.AllowedMsgAllowance)
    - [BasicAllowance](#cosmos.feegrant.v1beta1.BasicAllowance)
    - [Grant](#cosmos.feegrant.v1beta1.Grant)
    - [PeriodicAllowance](#cosmos.feegrant.v1beta1.PeriodicAllowance)
  
- [cosmos/feegrant/v1beta1/genesis.proto](#cosmos/feegrant/v1beta1/genesis.proto)
    - [GenesisState](#cosmos.feegrant.v1beta1.GenesisState)
  
- [cosmos/feegrant/v1beta1/query.proto](#cosmos/feegrant/v1beta1/query.proto)
    - [QueryAllowanceRequest](#cosmos.feegrant.v1beta1.QueryAllowanceRequest)
    - [QueryAllowanceResponse](#cosmos.feegrant.v1beta1.QueryAllowanceResponse)
    - [QueryAllowancesByGranterRequest](#cosmos.feegrant.v1beta1.QueryAllowancesByGranterRequest)
    - [QueryAllowancesByGranterResponse](#cosmos.feegrant.v1beta1.QueryAllowancesByGranterResponse)
    - [QueryAllowancesRequest](#cosmos.feegrant.v1beta1.QueryAllowancesRequest)
    - [QueryAllowancesResponse](#cosmos.feegrant.v1beta1.QueryAllowancesResponse)
  
    - [Query](#cosmos.feegrant.v1beta1.Query)
  
- [cosmos/feegrant/v1beta1/tx.proto](#cosmos/feegrant/v1beta1/tx.proto)
    - [MsgGrantAllowance](#cosmos.feegrant.v1beta1.MsgGrantAllowance)
    - [MsgGrantAllowanceResponse](#cosmos.feegrant.v1beta1.MsgGrantAllowanceResponse)
    - [MsgRevokeAllowance](#cosmos.feegrant.v1beta1.MsgRevokeAllowance)
    - [MsgRevokeAllowanceResponse](#cosmos.feegrant.v1beta1.MsgRevokeAllowanceResponse)
  
    - [Msg](#cosmos.feegrant.v1beta1.Msg)
  
- [cosmos/genutil/v1beta1/genesis.proto](#cosmos/genutil/v1beta1/genesis.proto)
    - [GenesisState](#cosmos.genutil.v1beta1.GenesisState)
  
- [cosmos/gov/v1beta1/gov.proto](#cosmos/gov/v1beta1/gov.proto)
    - [Deposit](#cosmos.gov.v1beta1.Deposit)
    - [DepositParams](#cosmos.gov.v1beta1.DepositParams)
    - [Proposal](#cosmos.gov.v1beta1.Proposal)
    - [TallyParams](#cosmos.gov.v1beta1.TallyParams)
    - [TallyResult](#cosmos.gov.v1beta1.TallyResult)
    - [TextProposal](#cosmos.gov.v1beta1.TextProposal)
    - [Vote](#cosmos.gov.v1beta1.Vote)
    - [VotingParams](#cosmos.gov.v1beta1.VotingParams)
    - [WeightedVoteOption](#cosmos.gov.v1beta1.WeightedVoteOption)
  
    - [ProposalStatus](#cosmos.gov.v1beta1.ProposalStatus)
    - [VoteOption](#cosmos.gov.v1beta1.VoteOption)
  
- [cosmos/gov/v1beta1/genesis.proto](#cosmos/gov/v1beta1/genesis.proto)
    - [GenesisState](#cosmos.gov.v1beta1.GenesisState)
  
- [cosmos/gov/v1beta1/query.proto](#cosmos/gov/v1beta1/query.proto)
    - [QueryDepositRequest](#cosmos.gov.v1beta1.QueryDepositRequest)
    - [QueryDepositResponse](#cosmos.gov.v1beta1.QueryDepositResponse)
    - [QueryDepositsRequest](#cosmos.gov.v1beta1.QueryDepositsRequest)
    - [QueryDepositsResponse](#cosmos.gov.v1beta1.QueryDepositsResponse)
    - [QueryParamsRequest](#cosmos.gov.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#cosmos.gov.v1beta1.QueryParamsResponse)
    - [QueryProposalRequest](#cosmos.gov.v1beta1.QueryProposalRequest)
    - [QueryProposalResponse](#cosmos.gov.v1beta1.QueryProposalResponse)
    - [QueryProposalsRequest](#cosmos.gov.v1beta1.QueryProposalsRequest)
    - [QueryProposalsResponse](#cosmos.gov.v1beta1.QueryProposalsResponse)
    - [QueryTallyResultRequest](#cosmos.gov.v1beta1.QueryTallyResultRequest)
    - [QueryTallyResultResponse](#cosmos.gov.v1beta1.QueryTallyResultResponse)
    - [QueryVoteRequest](#cosmos.gov.v1beta1.QueryVoteRequest)
    - [QueryVoteResponse](#cosmos.gov.v1beta1.QueryVoteResponse)
    - [QueryVotesRequest](#cosmos.gov.v1beta1.QueryVotesRequest)
    - [QueryVotesResponse](#cosmos.gov.v1beta1.QueryVotesResponse)
  
    - [Query](#cosmos.gov.v1beta1.Query)
  
- [cosmos/gov/v1beta1/tx.proto](#cosmos/gov/v1beta1/tx.proto)
    - [MsgDeposit](#cosmos.gov.v1beta1.MsgDeposit)
    - [MsgDepositResponse](#cosmos.gov.v1beta1.MsgDepositResponse)
    - [MsgSubmitProposal](#cosmos.gov.v1beta1.MsgSubmitProposal)
    - [MsgSubmitProposalResponse](#cosmos.gov.v1beta1.MsgSubmitProposalResponse)
    - [MsgVote](#cosmos.gov.v1beta1.MsgVote)
    - [MsgVoteResponse](#cosmos.gov.v1beta1.MsgVoteResponse)
    - [MsgVoteWeighted](#cosmos.gov.v1beta1.MsgVoteWeighted)
    - [MsgVoteWeightedResponse](#cosmos.gov.v1beta1.MsgVoteWeightedResponse)
  
    - [Msg](#cosmos.gov.v1beta1.Msg)
  
- [cosmos/mint/v1beta1/mint.proto](#cosmos/mint/v1beta1/mint.proto)
    - [Minter](#cosmos.mint.v1beta1.Minter)
    - [Params](#cosmos.mint.v1beta1.Params)
  
- [cosmos/mint/v1beta1/genesis.proto](#cosmos/mint/v1beta1/genesis.proto)
    - [GenesisState](#cosmos.mint.v1beta1.GenesisState)
  
- [cosmos/mint/v1beta1/query.proto](#cosmos/mint/v1beta1/query.proto)
    - [QueryAnnualProvisionsRequest](#cosmos.mint.v1beta1.QueryAnnualProvisionsRequest)
    - [QueryAnnualProvisionsResponse](#cosmos.mint.v1beta1.QueryAnnualProvisionsResponse)
    - [QueryInflationRequest](#cosmos.mint.v1beta1.QueryInflationRequest)
    - [QueryInflationResponse](#cosmos.mint.v1beta1.QueryInflationResponse)
    - [QueryParamsRequest](#cosmos.mint.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#cosmos.mint.v1beta1.QueryParamsResponse)
  
    - [Query](#cosmos.mint.v1beta1.Query)
  
- [cosmos/params/v1beta1/params.proto](#cosmos/params/v1beta1/params.proto)
    - [ParamChange](#cosmos.params.v1beta1.ParamChange)
    - [ParameterChangeProposal](#cosmos.params.v1beta1.ParameterChangeProposal)
  
- [cosmos/params/v1beta1/query.proto](#cosmos/params/v1beta1/query.proto)
    - [QueryParamsRequest](#cosmos.params.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#cosmos.params.v1beta1.QueryParamsResponse)
  
    - [Query](#cosmos.params.v1beta1.Query)
  
- [cosmos/slashing/v1beta1/slashing.proto](#cosmos/slashing/v1beta1/slashing.proto)
    - [Params](#cosmos.slashing.v1beta1.Params)
    - [ValidatorSigningInfo](#cosmos.slashing.v1beta1.ValidatorSigningInfo)
  
- [cosmos/slashing/v1beta1/genesis.proto](#cosmos/slashing/v1beta1/genesis.proto)
    - [GenesisState](#cosmos.slashing.v1beta1.GenesisState)
    - [MissedBlock](#cosmos.slashing.v1beta1.MissedBlock)
    - [SigningInfo](#cosmos.slashing.v1beta1.SigningInfo)
    - [ValidatorMissedBlocks](#cosmos.slashing.v1beta1.ValidatorMissedBlocks)
  
- [cosmos/slashing/v1beta1/query.proto](#cosmos/slashing/v1beta1/query.proto)
    - [QueryParamsRequest](#cosmos.slashing.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#cosmos.slashing.v1beta1.QueryParamsResponse)
    - [QuerySigningInfoRequest](#cosmos.slashing.v1beta1.QuerySigningInfoRequest)
    - [QuerySigningInfoResponse](#cosmos.slashing.v1beta1.QuerySigningInfoResponse)
    - [QuerySigningInfosRequest](#cosmos.slashing.v1beta1.QuerySigningInfosRequest)
    - [QuerySigningInfosResponse](#cosmos.slashing.v1beta1.QuerySigningInfosResponse)
  
    - [Query](#cosmos.slashing.v1beta1.Query)
  
- [cosmos/slashing/v1beta1/tx.proto](#cosmos/slashing/v1beta1/tx.proto)
    - [MsgUnjail](#cosmos.slashing.v1beta1.MsgUnjail)
    - [MsgUnjailResponse](#cosmos.slashing.v1beta1.MsgUnjailResponse)
  
    - [Msg](#cosmos.slashing.v1beta1.Msg)
  
- [cosmos/staking/v1beta1/authz.proto](#cosmos/staking/v1beta1/authz.proto)
    - [StakeAuthorization](#cosmos.staking.v1beta1.StakeAuthorization)
    - [StakeAuthorization.Validators](#cosmos.staking.v1beta1.StakeAuthorization.Validators)
  
    - [AuthorizationType](#cosmos.staking.v1beta1.AuthorizationType)
  
- [cosmos/staking/v1beta1/staking.proto](#cosmos/staking/v1beta1/staking.proto)
    - [Commission](#cosmos.staking.v1beta1.Commission)
    - [CommissionRates](#cosmos.staking.v1beta1.CommissionRates)
    - [DVPair](#cosmos.staking.v1beta1.DVPair)
    - [DVPairs](#cosmos.staking.v1beta1.DVPairs)
    - [DVVTriplet](#cosmos.staking.v1beta1.DVVTriplet)
    - [DVVTriplets](#cosmos.staking.v1beta1.DVVTriplets)
    - [Delegation](#cosmos.staking.v1beta1.Delegation)
    - [DelegationResponse](#cosmos.staking.v1beta1.DelegationResponse)
    - [Description](#cosmos.staking.v1beta1.Description)
    - [HistoricalInfo](#cosmos.staking.v1beta1.HistoricalInfo)
    - [Params](#cosmos.staking.v1beta1.Params)
    - [Pool](#cosmos.staking.v1beta1.Pool)
    - [Redelegation](#cosmos.staking.v1beta1.Redelegation)
    - [RedelegationEntry](#cosmos.staking.v1beta1.RedelegationEntry)
    - [RedelegationEntryResponse](#cosmos.staking.v1beta1.RedelegationEntryResponse)
    - [RedelegationResponse](#cosmos.staking.v1beta1.RedelegationResponse)
    - [UnbondingDelegation](#cosmos.staking.v1beta1.UnbondingDelegation)
    - [UnbondingDelegationEntry](#cosmos.staking.v1beta1.UnbondingDelegationEntry)
    - [ValAddresses](#cosmos.staking.v1beta1.ValAddresses)
    - [Validator](#cosmos.staking.v1beta1.Validator)
  
    - [BondStatus](#cosmos.staking.v1beta1.BondStatus)
  
- [cosmos/staking/v1beta1/genesis.proto](#cosmos/staking/v1beta1/genesis.proto)
    - [GenesisState](#cosmos.staking.v1beta1.GenesisState)
    - [LastValidatorPower](#cosmos.staking.v1beta1.LastValidatorPower)
  
- [cosmos/staking/v1beta1/query.proto](#cosmos/staking/v1beta1/query.proto)
    - [QueryDelegationRequest](#cosmos.staking.v1beta1.QueryDelegationRequest)
    - [QueryDelegationResponse](#cosmos.staking.v1beta1.QueryDelegationResponse)
    - [QueryDelegatorDelegationsRequest](#cosmos.staking.v1beta1.QueryDelegatorDelegationsRequest)
    - [QueryDelegatorDelegationsResponse](#cosmos.staking.v1beta1.QueryDelegatorDelegationsResponse)
    - [QueryDelegatorUnbondingDelegationsRequest](#cosmos.staking.v1beta1.QueryDelegatorUnbondingDelegationsRequest)
    - [QueryDelegatorUnbondingDelegationsResponse](#cosmos.staking.v1beta1.QueryDelegatorUnbondingDelegationsResponse)
    - [QueryDelegatorValidatorRequest](#cosmos.staking.v1beta1.QueryDelegatorValidatorRequest)
    - [QueryDelegatorValidatorResponse](#cosmos.staking.v1beta1.QueryDelegatorValidatorResponse)
    - [QueryDelegatorValidatorsRequest](#cosmos.staking.v1beta1.QueryDelegatorValidatorsRequest)
    - [QueryDelegatorValidatorsResponse](#cosmos.staking.v1beta1.QueryDelegatorValidatorsResponse)
    - [QueryHistoricalInfoRequest](#cosmos.staking.v1beta1.QueryHistoricalInfoRequest)
    - [QueryHistoricalInfoResponse](#cosmos.staking.v1beta1.QueryHistoricalInfoResponse)
    - [QueryParamsRequest](#cosmos.staking.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#cosmos.staking.v1beta1.QueryParamsResponse)
    - [QueryPoolRequest](#cosmos.staking.v1beta1.QueryPoolRequest)
    - [QueryPoolResponse](#cosmos.staking.v1beta1.QueryPoolResponse)
    - [QueryRedelegationsRequest](#cosmos.staking.v1beta1.QueryRedelegationsRequest)
    - [QueryRedelegationsResponse](#cosmos.staking.v1beta1.QueryRedelegationsResponse)
    - [QueryUnbondingDelegationRequest](#cosmos.staking.v1beta1.QueryUnbondingDelegationRequest)
    - [QueryUnbondingDelegationResponse](#cosmos.staking.v1beta1.QueryUnbondingDelegationResponse)
    - [QueryValidatorDelegationsRequest](#cosmos.staking.v1beta1.QueryValidatorDelegationsRequest)
    - [QueryValidatorDelegationsResponse](#cosmos.staking.v1beta1.QueryValidatorDelegationsResponse)
    - [QueryValidatorRequest](#cosmos.staking.v1beta1.QueryValidatorRequest)
    - [QueryValidatorResponse](#cosmos.staking.v1beta1.QueryValidatorResponse)
    - [QueryValidatorUnbondingDelegationsRequest](#cosmos.staking.v1beta1.QueryValidatorUnbondingDelegationsRequest)
    - [QueryValidatorUnbondingDelegationsResponse](#cosmos.staking.v1beta1.QueryValidatorUnbondingDelegationsResponse)
    - [QueryValidatorsRequest](#cosmos.staking.v1beta1.QueryValidatorsRequest)
    - [QueryValidatorsResponse](#cosmos.staking.v1beta1.QueryValidatorsResponse)
  
    - [Query](#cosmos.staking.v1beta1.Query)
  
- [cosmos/staking/v1beta1/tx.proto](#cosmos/staking/v1beta1/tx.proto)
    - [MsgBeginRedelegate](#cosmos.staking.v1beta1.MsgBeginRedelegate)
    - [MsgBeginRedelegateResponse](#cosmos.staking.v1beta1.MsgBeginRedelegateResponse)
    - [MsgCreateValidator](#cosmos.staking.v1beta1.MsgCreateValidator)
    - [MsgCreateValidatorResponse](#cosmos.staking.v1beta1.MsgCreateValidatorResponse)
    - [MsgDelegate](#cosmos.staking.v1beta1.MsgDelegate)
    - [MsgDelegateResponse](#cosmos.staking.v1beta1.MsgDelegateResponse)
    - [MsgEditValidator](#cosmos.staking.v1beta1.MsgEditValidator)
    - [MsgEditValidatorResponse](#cosmos.staking.v1beta1.MsgEditValidatorResponse)
    - [MsgUndelegate](#cosmos.staking.v1beta1.MsgUndelegate)
    - [MsgUndelegateResponse](#cosmos.staking.v1beta1.MsgUndelegateResponse)
  
    - [Msg](#cosmos.staking.v1beta1.Msg)
  
- [cosmos/tx/signing/v1beta1/signing.proto](#cosmos/tx/signing/v1beta1/signing.proto)
    - [SignatureDescriptor](#cosmos.tx.signing.v1beta1.SignatureDescriptor)
    - [SignatureDescriptor.Data](#cosmos.tx.signing.v1beta1.SignatureDescriptor.Data)
    - [SignatureDescriptor.Data.Multi](#cosmos.tx.signing.v1beta1.SignatureDescriptor.Data.Multi)
    - [SignatureDescriptor.Data.Single](#cosmos.tx.signing.v1beta1.SignatureDescriptor.Data.Single)
    - [SignatureDescriptors](#cosmos.tx.signing.v1beta1.SignatureDescriptors)
  
    - [SignMode](#cosmos.tx.signing.v1beta1.SignMode)
  
- [cosmos/tx/v1beta1/tx.proto](#cosmos/tx/v1beta1/tx.proto)
    - [AuthInfo](#cosmos.tx.v1beta1.AuthInfo)
    - [Fee](#cosmos.tx.v1beta1.Fee)
    - [ModeInfo](#cosmos.tx.v1beta1.ModeInfo)
    - [ModeInfo.Multi](#cosmos.tx.v1beta1.ModeInfo.Multi)
    - [ModeInfo.Single](#cosmos.tx.v1beta1.ModeInfo.Single)
    - [SignDoc](#cosmos.tx.v1beta1.SignDoc)
    - [SignerInfo](#cosmos.tx.v1beta1.SignerInfo)
    - [Tx](#cosmos.tx.v1beta1.Tx)
    - [TxBody](#cosmos.tx.v1beta1.TxBody)
    - [TxRaw](#cosmos.tx.v1beta1.TxRaw)
  
- [cosmos/tx/v1beta1/service.proto](#cosmos/tx/v1beta1/service.proto)
    - [BroadcastTxRequest](#cosmos.tx.v1beta1.BroadcastTxRequest)
    - [BroadcastTxResponse](#cosmos.tx.v1beta1.BroadcastTxResponse)
    - [GetBlockWithTxsRequest](#cosmos.tx.v1beta1.GetBlockWithTxsRequest)
    - [GetBlockWithTxsResponse](#cosmos.tx.v1beta1.GetBlockWithTxsResponse)
    - [GetTxRequest](#cosmos.tx.v1beta1.GetTxRequest)
    - [GetTxResponse](#cosmos.tx.v1beta1.GetTxResponse)
    - [GetTxsEventRequest](#cosmos.tx.v1beta1.GetTxsEventRequest)
    - [GetTxsEventResponse](#cosmos.tx.v1beta1.GetTxsEventResponse)
    - [SimulateRequest](#cosmos.tx.v1beta1.SimulateRequest)
    - [SimulateResponse](#cosmos.tx.v1beta1.SimulateResponse)
  
    - [BroadcastMode](#cosmos.tx.v1beta1.BroadcastMode)
    - [OrderBy](#cosmos.tx.v1beta1.OrderBy)
  
    - [Service](#cosmos.tx.v1beta1.Service)
  
- [cosmos/upgrade/v1beta1/upgrade.proto](#cosmos/upgrade/v1beta1/upgrade.proto)
    - [CancelSoftwareUpgradeProposal](#cosmos.upgrade.v1beta1.CancelSoftwareUpgradeProposal)
    - [ModuleVersion](#cosmos.upgrade.v1beta1.ModuleVersion)
    - [Plan](#cosmos.upgrade.v1beta1.Plan)
    - [SoftwareUpgradeProposal](#cosmos.upgrade.v1beta1.SoftwareUpgradeProposal)
  
- [cosmos/upgrade/v1beta1/query.proto](#cosmos/upgrade/v1beta1/query.proto)
    - [QueryAppliedPlanRequest](#cosmos.upgrade.v1beta1.QueryAppliedPlanRequest)
    - [QueryAppliedPlanResponse](#cosmos.upgrade.v1beta1.QueryAppliedPlanResponse)
    - [QueryCurrentPlanRequest](#cosmos.upgrade.v1beta1.QueryCurrentPlanRequest)
    - [QueryCurrentPlanResponse](#cosmos.upgrade.v1beta1.QueryCurrentPlanResponse)
    - [QueryModuleVersionsRequest](#cosmos.upgrade.v1beta1.QueryModuleVersionsRequest)
    - [QueryModuleVersionsResponse](#cosmos.upgrade.v1beta1.QueryModuleVersionsResponse)
    - [QueryUpgradedConsensusStateRequest](#cosmos.upgrade.v1beta1.QueryUpgradedConsensusStateRequest)
    - [QueryUpgradedConsensusStateResponse](#cosmos.upgrade.v1beta1.QueryUpgradedConsensusStateResponse)
  
    - [Query](#cosmos.upgrade.v1beta1.Query)
  
- [cosmos/vesting/v1beta1/tx.proto](#cosmos/vesting/v1beta1/tx.proto)
    - [MsgCreateVestingAccount](#cosmos.vesting.v1beta1.MsgCreateVestingAccount)
    - [MsgCreateVestingAccountResponse](#cosmos.vesting.v1beta1.MsgCreateVestingAccountResponse)
  
    - [Msg](#cosmos.vesting.v1beta1.Msg)
  
- [cosmos/vesting/v1beta1/vesting.proto](#cosmos/vesting/v1beta1/vesting.proto)
    - [BaseVestingAccount](#cosmos.vesting.v1beta1.BaseVestingAccount)
    - [ContinuousVestingAccount](#cosmos.vesting.v1beta1.ContinuousVestingAccount)
    - [DelayedVestingAccount](#cosmos.vesting.v1beta1.DelayedVestingAccount)
    - [Period](#cosmos.vesting.v1beta1.Period)
    - [PeriodicVestingAccount](#cosmos.vesting.v1beta1.PeriodicVestingAccount)
    - [PermanentLockedAccount](#cosmos.vesting.v1beta1.PermanentLockedAccount)
  
- [lbm/bankplus/v1/bankplus.proto](#lbm/bankplus/v1/bankplus.proto)
    - [InactiveAddr](#lbm.bankplus.v1.InactiveAddr)
  
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
  
- [lbm/collection/v1/collection.proto](#lbm/collection/v1/collection.proto)
    - [Attribute](#lbm.collection.v1.Attribute)
    - [Authorization](#lbm.collection.v1.Authorization)
    - [Change](#lbm.collection.v1.Change)
    - [Coin](#lbm.collection.v1.Coin)
    - [Contract](#lbm.collection.v1.Contract)
    - [FT](#lbm.collection.v1.FT)
    - [FTClass](#lbm.collection.v1.FTClass)
    - [Grant](#lbm.collection.v1.Grant)
    - [NFT](#lbm.collection.v1.NFT)
    - [NFTClass](#lbm.collection.v1.NFTClass)
    - [OwnerNFT](#lbm.collection.v1.OwnerNFT)
    - [Params](#lbm.collection.v1.Params)
    - [TokenType](#lbm.collection.v1.TokenType)
  
    - [LegacyPermission](#lbm.collection.v1.LegacyPermission)
    - [Permission](#lbm.collection.v1.Permission)
  
- [lbm/collection/v1/event.proto](#lbm/collection/v1/event.proto)
    - [EventAttached](#lbm.collection.v1.EventAttached)
    - [EventAuthorizedOperator](#lbm.collection.v1.EventAuthorizedOperator)
    - [EventBurned](#lbm.collection.v1.EventBurned)
    - [EventCreatedContract](#lbm.collection.v1.EventCreatedContract)
    - [EventCreatedFTClass](#lbm.collection.v1.EventCreatedFTClass)
    - [EventCreatedNFTClass](#lbm.collection.v1.EventCreatedNFTClass)
    - [EventDetached](#lbm.collection.v1.EventDetached)
    - [EventGranted](#lbm.collection.v1.EventGranted)
    - [EventMintedFT](#lbm.collection.v1.EventMintedFT)
    - [EventMintedNFT](#lbm.collection.v1.EventMintedNFT)
    - [EventModifiedContract](#lbm.collection.v1.EventModifiedContract)
    - [EventModifiedNFT](#lbm.collection.v1.EventModifiedNFT)
    - [EventModifiedTokenClass](#lbm.collection.v1.EventModifiedTokenClass)
    - [EventOwnerChanged](#lbm.collection.v1.EventOwnerChanged)
    - [EventRenounced](#lbm.collection.v1.EventRenounced)
    - [EventRevokedOperator](#lbm.collection.v1.EventRevokedOperator)
    - [EventRootChanged](#lbm.collection.v1.EventRootChanged)
    - [EventSent](#lbm.collection.v1.EventSent)
  
    - [AttributeKey](#lbm.collection.v1.AttributeKey)
    - [EventType](#lbm.collection.v1.EventType)
  
- [lbm/collection/v1/genesis.proto](#lbm/collection/v1/genesis.proto)
    - [Balance](#lbm.collection.v1.Balance)
    - [ClassStatistics](#lbm.collection.v1.ClassStatistics)
    - [ContractAuthorizations](#lbm.collection.v1.ContractAuthorizations)
    - [ContractBalances](#lbm.collection.v1.ContractBalances)
    - [ContractClasses](#lbm.collection.v1.ContractClasses)
    - [ContractGrants](#lbm.collection.v1.ContractGrants)
    - [ContractNFTs](#lbm.collection.v1.ContractNFTs)
    - [ContractNextTokenIDs](#lbm.collection.v1.ContractNextTokenIDs)
    - [ContractStatistics](#lbm.collection.v1.ContractStatistics)
    - [ContractTokenRelations](#lbm.collection.v1.ContractTokenRelations)
    - [GenesisState](#lbm.collection.v1.GenesisState)
    - [NextClassIDs](#lbm.collection.v1.NextClassIDs)
    - [NextTokenID](#lbm.collection.v1.NextTokenID)
    - [TokenRelation](#lbm.collection.v1.TokenRelation)
  
- [lbm/collection/v1/query.proto](#lbm/collection/v1/query.proto)
    - [QueryAllBalancesRequest](#lbm.collection.v1.QueryAllBalancesRequest)
    - [QueryAllBalancesResponse](#lbm.collection.v1.QueryAllBalancesResponse)
    - [QueryApprovedRequest](#lbm.collection.v1.QueryApprovedRequest)
    - [QueryApprovedResponse](#lbm.collection.v1.QueryApprovedResponse)
    - [QueryApproversRequest](#lbm.collection.v1.QueryApproversRequest)
    - [QueryApproversResponse](#lbm.collection.v1.QueryApproversResponse)
    - [QueryBalanceRequest](#lbm.collection.v1.QueryBalanceRequest)
    - [QueryBalanceResponse](#lbm.collection.v1.QueryBalanceResponse)
    - [QueryChildrenRequest](#lbm.collection.v1.QueryChildrenRequest)
    - [QueryChildrenResponse](#lbm.collection.v1.QueryChildrenResponse)
    - [QueryContractRequest](#lbm.collection.v1.QueryContractRequest)
    - [QueryContractResponse](#lbm.collection.v1.QueryContractResponse)
    - [QueryFTBurntRequest](#lbm.collection.v1.QueryFTBurntRequest)
    - [QueryFTBurntResponse](#lbm.collection.v1.QueryFTBurntResponse)
    - [QueryFTMintedRequest](#lbm.collection.v1.QueryFTMintedRequest)
    - [QueryFTMintedResponse](#lbm.collection.v1.QueryFTMintedResponse)
    - [QueryFTSupplyRequest](#lbm.collection.v1.QueryFTSupplyRequest)
    - [QueryFTSupplyResponse](#lbm.collection.v1.QueryFTSupplyResponse)
    - [QueryGranteeGrantsRequest](#lbm.collection.v1.QueryGranteeGrantsRequest)
    - [QueryGranteeGrantsResponse](#lbm.collection.v1.QueryGranteeGrantsResponse)
    - [QueryNFTBurntRequest](#lbm.collection.v1.QueryNFTBurntRequest)
    - [QueryNFTBurntResponse](#lbm.collection.v1.QueryNFTBurntResponse)
    - [QueryNFTMintedRequest](#lbm.collection.v1.QueryNFTMintedRequest)
    - [QueryNFTMintedResponse](#lbm.collection.v1.QueryNFTMintedResponse)
    - [QueryNFTSupplyRequest](#lbm.collection.v1.QueryNFTSupplyRequest)
    - [QueryNFTSupplyResponse](#lbm.collection.v1.QueryNFTSupplyResponse)
    - [QueryParentRequest](#lbm.collection.v1.QueryParentRequest)
    - [QueryParentResponse](#lbm.collection.v1.QueryParentResponse)
    - [QueryRootRequest](#lbm.collection.v1.QueryRootRequest)
    - [QueryRootResponse](#lbm.collection.v1.QueryRootResponse)
    - [QueryTokenClassTypeNameRequest](#lbm.collection.v1.QueryTokenClassTypeNameRequest)
    - [QueryTokenClassTypeNameResponse](#lbm.collection.v1.QueryTokenClassTypeNameResponse)
    - [QueryTokenRequest](#lbm.collection.v1.QueryTokenRequest)
    - [QueryTokenResponse](#lbm.collection.v1.QueryTokenResponse)
    - [QueryTokenTypeRequest](#lbm.collection.v1.QueryTokenTypeRequest)
    - [QueryTokenTypeResponse](#lbm.collection.v1.QueryTokenTypeResponse)
    - [QueryTokenTypesRequest](#lbm.collection.v1.QueryTokenTypesRequest)
    - [QueryTokenTypesResponse](#lbm.collection.v1.QueryTokenTypesResponse)
    - [QueryTokensRequest](#lbm.collection.v1.QueryTokensRequest)
    - [QueryTokensResponse](#lbm.collection.v1.QueryTokensResponse)
    - [QueryTokensWithTokenTypeRequest](#lbm.collection.v1.QueryTokensWithTokenTypeRequest)
    - [QueryTokensWithTokenTypeResponse](#lbm.collection.v1.QueryTokensWithTokenTypeResponse)
  
    - [Query](#lbm.collection.v1.Query)
  
- [lbm/collection/v1/tx.proto](#lbm/collection/v1/tx.proto)
    - [MintNFTParam](#lbm.collection.v1.MintNFTParam)
    - [MsgApprove](#lbm.collection.v1.MsgApprove)
    - [MsgApproveResponse](#lbm.collection.v1.MsgApproveResponse)
    - [MsgAttach](#lbm.collection.v1.MsgAttach)
    - [MsgAttachFrom](#lbm.collection.v1.MsgAttachFrom)
    - [MsgAttachFromResponse](#lbm.collection.v1.MsgAttachFromResponse)
    - [MsgAttachResponse](#lbm.collection.v1.MsgAttachResponse)
    - [MsgBurnFT](#lbm.collection.v1.MsgBurnFT)
    - [MsgBurnFTFrom](#lbm.collection.v1.MsgBurnFTFrom)
    - [MsgBurnFTFromResponse](#lbm.collection.v1.MsgBurnFTFromResponse)
    - [MsgBurnFTResponse](#lbm.collection.v1.MsgBurnFTResponse)
    - [MsgBurnNFT](#lbm.collection.v1.MsgBurnNFT)
    - [MsgBurnNFTFrom](#lbm.collection.v1.MsgBurnNFTFrom)
    - [MsgBurnNFTFromResponse](#lbm.collection.v1.MsgBurnNFTFromResponse)
    - [MsgBurnNFTResponse](#lbm.collection.v1.MsgBurnNFTResponse)
    - [MsgCreateContract](#lbm.collection.v1.MsgCreateContract)
    - [MsgCreateContractResponse](#lbm.collection.v1.MsgCreateContractResponse)
    - [MsgDetach](#lbm.collection.v1.MsgDetach)
    - [MsgDetachFrom](#lbm.collection.v1.MsgDetachFrom)
    - [MsgDetachFromResponse](#lbm.collection.v1.MsgDetachFromResponse)
    - [MsgDetachResponse](#lbm.collection.v1.MsgDetachResponse)
    - [MsgDisapprove](#lbm.collection.v1.MsgDisapprove)
    - [MsgDisapproveResponse](#lbm.collection.v1.MsgDisapproveResponse)
    - [MsgGrantPermission](#lbm.collection.v1.MsgGrantPermission)
    - [MsgGrantPermissionResponse](#lbm.collection.v1.MsgGrantPermissionResponse)
    - [MsgIssueFT](#lbm.collection.v1.MsgIssueFT)
    - [MsgIssueFTResponse](#lbm.collection.v1.MsgIssueFTResponse)
    - [MsgIssueNFT](#lbm.collection.v1.MsgIssueNFT)
    - [MsgIssueNFTResponse](#lbm.collection.v1.MsgIssueNFTResponse)
    - [MsgMintFT](#lbm.collection.v1.MsgMintFT)
    - [MsgMintFTResponse](#lbm.collection.v1.MsgMintFTResponse)
    - [MsgMintNFT](#lbm.collection.v1.MsgMintNFT)
    - [MsgMintNFTResponse](#lbm.collection.v1.MsgMintNFTResponse)
    - [MsgModify](#lbm.collection.v1.MsgModify)
    - [MsgModifyResponse](#lbm.collection.v1.MsgModifyResponse)
    - [MsgRevokePermission](#lbm.collection.v1.MsgRevokePermission)
    - [MsgRevokePermissionResponse](#lbm.collection.v1.MsgRevokePermissionResponse)
    - [MsgTransferFT](#lbm.collection.v1.MsgTransferFT)
    - [MsgTransferFTFrom](#lbm.collection.v1.MsgTransferFTFrom)
    - [MsgTransferFTFromResponse](#lbm.collection.v1.MsgTransferFTFromResponse)
    - [MsgTransferFTResponse](#lbm.collection.v1.MsgTransferFTResponse)
    - [MsgTransferNFT](#lbm.collection.v1.MsgTransferNFT)
    - [MsgTransferNFTFrom](#lbm.collection.v1.MsgTransferNFTFrom)
    - [MsgTransferNFTFromResponse](#lbm.collection.v1.MsgTransferNFTFromResponse)
    - [MsgTransferNFTResponse](#lbm.collection.v1.MsgTransferNFTResponse)
  
    - [Msg](#lbm.collection.v1.Msg)
  
- [lbm/foundation/v1/authz.proto](#lbm/foundation/v1/authz.proto)
    - [ReceiveFromTreasuryAuthorization](#lbm.foundation.v1.ReceiveFromTreasuryAuthorization)
  
- [lbm/foundation/v1/foundation.proto](#lbm/foundation/v1/foundation.proto)
    - [DecisionPolicyWindows](#lbm.foundation.v1.DecisionPolicyWindows)
    - [FoundationInfo](#lbm.foundation.v1.FoundationInfo)
    - [Member](#lbm.foundation.v1.Member)
    - [MemberRequest](#lbm.foundation.v1.MemberRequest)
    - [OutsourcingDecisionPolicy](#lbm.foundation.v1.OutsourcingDecisionPolicy)
    - [Params](#lbm.foundation.v1.Params)
    - [PercentageDecisionPolicy](#lbm.foundation.v1.PercentageDecisionPolicy)
    - [Pool](#lbm.foundation.v1.Pool)
    - [Proposal](#lbm.foundation.v1.Proposal)
    - [TallyResult](#lbm.foundation.v1.TallyResult)
    - [ThresholdDecisionPolicy](#lbm.foundation.v1.ThresholdDecisionPolicy)
    - [Vote](#lbm.foundation.v1.Vote)
  
    - [ProposalExecutorResult](#lbm.foundation.v1.ProposalExecutorResult)
    - [ProposalStatus](#lbm.foundation.v1.ProposalStatus)
    - [VoteOption](#lbm.foundation.v1.VoteOption)
  
- [lbm/foundation/v1/event.proto](#lbm/foundation/v1/event.proto)
    - [EventExec](#lbm.foundation.v1.EventExec)
    - [EventFundTreasury](#lbm.foundation.v1.EventFundTreasury)
    - [EventGrant](#lbm.foundation.v1.EventGrant)
    - [EventLeaveFoundation](#lbm.foundation.v1.EventLeaveFoundation)
    - [EventRevoke](#lbm.foundation.v1.EventRevoke)
    - [EventSubmitProposal](#lbm.foundation.v1.EventSubmitProposal)
    - [EventUpdateDecisionPolicy](#lbm.foundation.v1.EventUpdateDecisionPolicy)
    - [EventUpdateMembers](#lbm.foundation.v1.EventUpdateMembers)
    - [EventUpdateParams](#lbm.foundation.v1.EventUpdateParams)
    - [EventVote](#lbm.foundation.v1.EventVote)
    - [EventWithdrawFromTreasury](#lbm.foundation.v1.EventWithdrawFromTreasury)
    - [EventWithdrawProposal](#lbm.foundation.v1.EventWithdrawProposal)
  
- [lbm/foundation/v1/genesis.proto](#lbm/foundation/v1/genesis.proto)
    - [GenesisState](#lbm.foundation.v1.GenesisState)
    - [GrantAuthorization](#lbm.foundation.v1.GrantAuthorization)
  
- [lbm/foundation/v1/query.proto](#lbm/foundation/v1/query.proto)
    - [QueryFoundationInfoRequest](#lbm.foundation.v1.QueryFoundationInfoRequest)
    - [QueryFoundationInfoResponse](#lbm.foundation.v1.QueryFoundationInfoResponse)
    - [QueryGrantsRequest](#lbm.foundation.v1.QueryGrantsRequest)
    - [QueryGrantsResponse](#lbm.foundation.v1.QueryGrantsResponse)
    - [QueryMemberRequest](#lbm.foundation.v1.QueryMemberRequest)
    - [QueryMemberResponse](#lbm.foundation.v1.QueryMemberResponse)
    - [QueryMembersRequest](#lbm.foundation.v1.QueryMembersRequest)
    - [QueryMembersResponse](#lbm.foundation.v1.QueryMembersResponse)
    - [QueryParamsRequest](#lbm.foundation.v1.QueryParamsRequest)
    - [QueryParamsResponse](#lbm.foundation.v1.QueryParamsResponse)
    - [QueryProposalRequest](#lbm.foundation.v1.QueryProposalRequest)
    - [QueryProposalResponse](#lbm.foundation.v1.QueryProposalResponse)
    - [QueryProposalsRequest](#lbm.foundation.v1.QueryProposalsRequest)
    - [QueryProposalsResponse](#lbm.foundation.v1.QueryProposalsResponse)
    - [QueryTallyResultRequest](#lbm.foundation.v1.QueryTallyResultRequest)
    - [QueryTallyResultResponse](#lbm.foundation.v1.QueryTallyResultResponse)
    - [QueryTreasuryRequest](#lbm.foundation.v1.QueryTreasuryRequest)
    - [QueryTreasuryResponse](#lbm.foundation.v1.QueryTreasuryResponse)
    - [QueryVoteRequest](#lbm.foundation.v1.QueryVoteRequest)
    - [QueryVoteResponse](#lbm.foundation.v1.QueryVoteResponse)
    - [QueryVotesRequest](#lbm.foundation.v1.QueryVotesRequest)
    - [QueryVotesResponse](#lbm.foundation.v1.QueryVotesResponse)
  
    - [Query](#lbm.foundation.v1.Query)
  
- [lbm/foundation/v1/tx.proto](#lbm/foundation/v1/tx.proto)
    - [MsgExec](#lbm.foundation.v1.MsgExec)
    - [MsgExecResponse](#lbm.foundation.v1.MsgExecResponse)
    - [MsgFundTreasury](#lbm.foundation.v1.MsgFundTreasury)
    - [MsgFundTreasuryResponse](#lbm.foundation.v1.MsgFundTreasuryResponse)
    - [MsgGrant](#lbm.foundation.v1.MsgGrant)
    - [MsgGrantResponse](#lbm.foundation.v1.MsgGrantResponse)
    - [MsgLeaveFoundation](#lbm.foundation.v1.MsgLeaveFoundation)
    - [MsgLeaveFoundationResponse](#lbm.foundation.v1.MsgLeaveFoundationResponse)
    - [MsgRevoke](#lbm.foundation.v1.MsgRevoke)
    - [MsgRevokeResponse](#lbm.foundation.v1.MsgRevokeResponse)
    - [MsgSubmitProposal](#lbm.foundation.v1.MsgSubmitProposal)
    - [MsgSubmitProposalResponse](#lbm.foundation.v1.MsgSubmitProposalResponse)
    - [MsgUpdateDecisionPolicy](#lbm.foundation.v1.MsgUpdateDecisionPolicy)
    - [MsgUpdateDecisionPolicyResponse](#lbm.foundation.v1.MsgUpdateDecisionPolicyResponse)
    - [MsgUpdateMembers](#lbm.foundation.v1.MsgUpdateMembers)
    - [MsgUpdateMembersResponse](#lbm.foundation.v1.MsgUpdateMembersResponse)
    - [MsgUpdateParams](#lbm.foundation.v1.MsgUpdateParams)
    - [MsgUpdateParamsResponse](#lbm.foundation.v1.MsgUpdateParamsResponse)
    - [MsgVote](#lbm.foundation.v1.MsgVote)
    - [MsgVoteResponse](#lbm.foundation.v1.MsgVoteResponse)
    - [MsgWithdrawFromTreasury](#lbm.foundation.v1.MsgWithdrawFromTreasury)
    - [MsgWithdrawFromTreasuryResponse](#lbm.foundation.v1.MsgWithdrawFromTreasuryResponse)
    - [MsgWithdrawProposal](#lbm.foundation.v1.MsgWithdrawProposal)
    - [MsgWithdrawProposalResponse](#lbm.foundation.v1.MsgWithdrawProposalResponse)
  
    - [Exec](#lbm.foundation.v1.Exec)
  
    - [Msg](#lbm.foundation.v1.Msg)
  
- [lbm/stakingplus/v1/authz.proto](#lbm/stakingplus/v1/authz.proto)
    - [CreateValidatorAuthorization](#lbm.stakingplus.v1.CreateValidatorAuthorization)
  
- [lbm/token/v1/token.proto](#lbm/token/v1/token.proto)
    - [Authorization](#lbm.token.v1.Authorization)
    - [Contract](#lbm.token.v1.Contract)
    - [Grant](#lbm.token.v1.Grant)
    - [Pair](#lbm.token.v1.Pair)
    - [Params](#lbm.token.v1.Params)
  
    - [LegacyPermission](#lbm.token.v1.LegacyPermission)
    - [Permission](#lbm.token.v1.Permission)
  
- [lbm/token/v1/event.proto](#lbm/token/v1/event.proto)
    - [EventAuthorizedOperator](#lbm.token.v1.EventAuthorizedOperator)
    - [EventBurned](#lbm.token.v1.EventBurned)
    - [EventGranted](#lbm.token.v1.EventGranted)
    - [EventIssued](#lbm.token.v1.EventIssued)
    - [EventMinted](#lbm.token.v1.EventMinted)
    - [EventModified](#lbm.token.v1.EventModified)
    - [EventRenounced](#lbm.token.v1.EventRenounced)
    - [EventRevokedOperator](#lbm.token.v1.EventRevokedOperator)
    - [EventSent](#lbm.token.v1.EventSent)
  
    - [AttributeKey](#lbm.token.v1.AttributeKey)
    - [EventType](#lbm.token.v1.EventType)
  
- [lbm/token/v1/genesis.proto](#lbm/token/v1/genesis.proto)
    - [Balance](#lbm.token.v1.Balance)
    - [ClassGenesisState](#lbm.token.v1.ClassGenesisState)
    - [ContractAuthorizations](#lbm.token.v1.ContractAuthorizations)
    - [ContractBalances](#lbm.token.v1.ContractBalances)
    - [ContractCoin](#lbm.token.v1.ContractCoin)
    - [ContractGrants](#lbm.token.v1.ContractGrants)
    - [GenesisState](#lbm.token.v1.GenesisState)
  
- [lbm/token/v1/query.proto](#lbm/token/v1/query.proto)
    - [QueryApprovedRequest](#lbm.token.v1.QueryApprovedRequest)
    - [QueryApprovedResponse](#lbm.token.v1.QueryApprovedResponse)
    - [QueryApproversRequest](#lbm.token.v1.QueryApproversRequest)
    - [QueryApproversResponse](#lbm.token.v1.QueryApproversResponse)
    - [QueryBalanceRequest](#lbm.token.v1.QueryBalanceRequest)
    - [QueryBalanceResponse](#lbm.token.v1.QueryBalanceResponse)
    - [QueryBurntRequest](#lbm.token.v1.QueryBurntRequest)
    - [QueryBurntResponse](#lbm.token.v1.QueryBurntResponse)
    - [QueryContractRequest](#lbm.token.v1.QueryContractRequest)
    - [QueryContractResponse](#lbm.token.v1.QueryContractResponse)
    - [QueryContractsRequest](#lbm.token.v1.QueryContractsRequest)
    - [QueryContractsResponse](#lbm.token.v1.QueryContractsResponse)
    - [QueryGranteeGrantsRequest](#lbm.token.v1.QueryGranteeGrantsRequest)
    - [QueryGranteeGrantsResponse](#lbm.token.v1.QueryGranteeGrantsResponse)
    - [QueryMintedRequest](#lbm.token.v1.QueryMintedRequest)
    - [QueryMintedResponse](#lbm.token.v1.QueryMintedResponse)
    - [QuerySupplyRequest](#lbm.token.v1.QuerySupplyRequest)
    - [QuerySupplyResponse](#lbm.token.v1.QuerySupplyResponse)
  
    - [Query](#lbm.token.v1.Query)
  
- [lbm/token/v1/tx.proto](#lbm/token/v1/tx.proto)
    - [MsgApprove](#lbm.token.v1.MsgApprove)
    - [MsgApproveResponse](#lbm.token.v1.MsgApproveResponse)
    - [MsgBurn](#lbm.token.v1.MsgBurn)
    - [MsgBurnFrom](#lbm.token.v1.MsgBurnFrom)
    - [MsgBurnFromResponse](#lbm.token.v1.MsgBurnFromResponse)
    - [MsgBurnResponse](#lbm.token.v1.MsgBurnResponse)
    - [MsgGrantPermission](#lbm.token.v1.MsgGrantPermission)
    - [MsgGrantPermissionResponse](#lbm.token.v1.MsgGrantPermissionResponse)
    - [MsgIssue](#lbm.token.v1.MsgIssue)
    - [MsgIssueResponse](#lbm.token.v1.MsgIssueResponse)
    - [MsgMint](#lbm.token.v1.MsgMint)
    - [MsgMintResponse](#lbm.token.v1.MsgMintResponse)
    - [MsgModify](#lbm.token.v1.MsgModify)
    - [MsgModifyResponse](#lbm.token.v1.MsgModifyResponse)
    - [MsgRevokeOperator](#lbm.token.v1.MsgRevokeOperator)
    - [MsgRevokeOperatorResponse](#lbm.token.v1.MsgRevokeOperatorResponse)
    - [MsgRevokePermission](#lbm.token.v1.MsgRevokePermission)
    - [MsgRevokePermissionResponse](#lbm.token.v1.MsgRevokePermissionResponse)
    - [MsgSend](#lbm.token.v1.MsgSend)
    - [MsgSendResponse](#lbm.token.v1.MsgSendResponse)
    - [MsgTransferFrom](#lbm.token.v1.MsgTransferFrom)
    - [MsgTransferFromResponse](#lbm.token.v1.MsgTransferFromResponse)
  
    - [Msg](#lbm.token.v1.Msg)
  
- [Scalar Value Types](#scalar-value-types)



<a name="cosmos/auth/v1beta1/auth.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/auth/v1beta1/auth.proto



<a name="cosmos.auth.v1beta1.BaseAccount"></a>

### BaseAccount
BaseAccount defines a base account type. It contains all the necessary fields
for basic account functionality. Any custom account type should extend this
type for additional functionality (e.g. vesting).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |
| `pub_key` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `account_number` | [uint64](#uint64) |  |  |
| `sequence` | [uint64](#uint64) |  |  |






<a name="cosmos.auth.v1beta1.ModuleAccount"></a>

### ModuleAccount
ModuleAccount defines an account for modules that holds coins on a pool.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_account` | [BaseAccount](#cosmos.auth.v1beta1.BaseAccount) |  |  |
| `name` | [string](#string) |  |  |
| `permissions` | [string](#string) | repeated |  |






<a name="cosmos.auth.v1beta1.Params"></a>

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



<a name="cosmos/auth/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/auth/v1beta1/genesis.proto



<a name="cosmos.auth.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the auth module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmos.auth.v1beta1.Params) |  | params defines all the paramaters of the module. |
| `accounts` | [google.protobuf.Any](#google.protobuf.Any) | repeated | accounts are the accounts present at genesis. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/base/query/v1beta1/pagination.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/base/query/v1beta1/pagination.proto



<a name="cosmos.base.query.v1beta1.PageRequest"></a>

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
| `reverse` | [bool](#bool) |  | reverse is set to true if results are to be returned in the descending order.

Since: cosmos-sdk 0.43 |






<a name="cosmos.base.query.v1beta1.PageResponse"></a>

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



<a name="cosmos/auth/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/auth/v1beta1/query.proto



<a name="cosmos.auth.v1beta1.QueryAccountRequest"></a>

### QueryAccountRequest
QueryAccountRequest is the request type for the Query/Account RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address defines the address to query for. |






<a name="cosmos.auth.v1beta1.QueryAccountResponse"></a>

### QueryAccountResponse
QueryAccountResponse is the response type for the Query/Account RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account` | [google.protobuf.Any](#google.protobuf.Any) |  | account defines the account of the corresponding address. |






<a name="cosmos.auth.v1beta1.QueryAccountsRequest"></a>

### QueryAccountsRequest
QueryAccountsRequest is the request type for the Query/Accounts RPC method.

Since: cosmos-sdk 0.43


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.auth.v1beta1.QueryAccountsResponse"></a>

### QueryAccountsResponse
QueryAccountsResponse is the response type for the Query/Accounts RPC method.

Since: cosmos-sdk 0.43


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `accounts` | [google.protobuf.Any](#google.protobuf.Any) | repeated | accounts are the existing accounts |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.auth.v1beta1.QueryModuleAccountByNameRequest"></a>

### QueryModuleAccountByNameRequest
QueryModuleAccountByNameRequest is the request type for the Query/ModuleAccountByName RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |






<a name="cosmos.auth.v1beta1.QueryModuleAccountByNameResponse"></a>

### QueryModuleAccountByNameResponse
QueryModuleAccountByNameResponse is the response type for the Query/ModuleAccountByName RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account` | [google.protobuf.Any](#google.protobuf.Any) |  |  |






<a name="cosmos.auth.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="cosmos.auth.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmos.auth.v1beta1.Params) |  | params defines the parameters of the module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.auth.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Accounts` | [QueryAccountsRequest](#cosmos.auth.v1beta1.QueryAccountsRequest) | [QueryAccountsResponse](#cosmos.auth.v1beta1.QueryAccountsResponse) | Accounts returns all the existing accounts

Since: cosmos-sdk 0.43 | GET|/cosmos/auth/v1beta1/accounts|
| `Account` | [QueryAccountRequest](#cosmos.auth.v1beta1.QueryAccountRequest) | [QueryAccountResponse](#cosmos.auth.v1beta1.QueryAccountResponse) | Account returns account details based on address. | GET|/cosmos/auth/v1beta1/accounts/{address}|
| `Params` | [QueryParamsRequest](#cosmos.auth.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#cosmos.auth.v1beta1.QueryParamsResponse) | Params queries all parameters. | GET|/cosmos/auth/v1beta1/params|
| `ModuleAccountByName` | [QueryModuleAccountByNameRequest](#cosmos.auth.v1beta1.QueryModuleAccountByNameRequest) | [QueryModuleAccountByNameResponse](#cosmos.auth.v1beta1.QueryModuleAccountByNameResponse) | ModuleAccountByName returns the module account info by module name | GET|/cosmos/auth/v1beta1/module_accounts/{name}|

 <!-- end services -->



<a name="cosmos/authz/v1beta1/authz.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/authz/v1beta1/authz.proto
Since: cosmos-sdk 0.43


<a name="cosmos.authz.v1beta1.GenericAuthorization"></a>

### GenericAuthorization
GenericAuthorization gives the grantee unrestricted permissions to execute
the provided method on behalf of the granter's account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `msg` | [string](#string) |  | Msg, identified by it's type URL, to grant unrestricted permissions to execute |






<a name="cosmos.authz.v1beta1.Grant"></a>

### Grant
Grant gives permissions to execute
the provide method with expiration time.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authorization` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `expiration` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |






<a name="cosmos.authz.v1beta1.GrantAuthorization"></a>

### GrantAuthorization
GrantAuthorization extends a grant with both the addresses of the grantee and granter.
It is used in genesis.proto and query.proto

Since: cosmos-sdk 0.45.2


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  |  |
| `grantee` | [string](#string) |  |  |
| `authorization` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `expiration` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/authz/v1beta1/event.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/authz/v1beta1/event.proto
Since: cosmos-sdk 0.43


<a name="cosmos.authz.v1beta1.EventGrant"></a>

### EventGrant
EventGrant is emitted on Msg/Grant


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `msg_type_url` | [string](#string) |  | Msg type URL for which an autorization is granted |
| `granter` | [string](#string) |  | Granter account address |
| `grantee` | [string](#string) |  | Grantee account address |






<a name="cosmos.authz.v1beta1.EventRevoke"></a>

### EventRevoke
EventRevoke is emitted on Msg/Revoke


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `msg_type_url` | [string](#string) |  | Msg type URL for which an autorization is revoked |
| `granter` | [string](#string) |  | Granter account address |
| `grantee` | [string](#string) |  | Grantee account address |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/authz/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/authz/v1beta1/genesis.proto
Since: cosmos-sdk 0.43


<a name="cosmos.authz.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the authz module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authorization` | [GrantAuthorization](#cosmos.authz.v1beta1.GrantAuthorization) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/authz/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/authz/v1beta1/query.proto
Since: cosmos-sdk 0.43


<a name="cosmos.authz.v1beta1.QueryGranteeGrantsRequest"></a>

### QueryGranteeGrantsRequest
QueryGranteeGrantsRequest is the request type for the Query/IssuedGrants RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grantee` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an pagination for the request. |






<a name="cosmos.authz.v1beta1.QueryGranteeGrantsResponse"></a>

### QueryGranteeGrantsResponse
QueryGranteeGrantsResponse is the response type for the Query/GranteeGrants RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grants` | [GrantAuthorization](#cosmos.authz.v1beta1.GrantAuthorization) | repeated | grants is a list of grants granted to the grantee. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines an pagination for the response. |






<a name="cosmos.authz.v1beta1.QueryGranterGrantsRequest"></a>

### QueryGranterGrantsRequest
QueryGranterGrantsRequest is the request type for the Query/GranterGrants RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an pagination for the request. |






<a name="cosmos.authz.v1beta1.QueryGranterGrantsResponse"></a>

### QueryGranterGrantsResponse
QueryGranterGrantsResponse is the response type for the Query/GranterGrants RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grants` | [GrantAuthorization](#cosmos.authz.v1beta1.GrantAuthorization) | repeated | grants is a list of grants granted by the granter. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines an pagination for the response. |






<a name="cosmos.authz.v1beta1.QueryGrantsRequest"></a>

### QueryGrantsRequest
QueryGrantsRequest is the request type for the Query/Grants RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  |  |
| `grantee` | [string](#string) |  |  |
| `msg_type_url` | [string](#string) |  | Optional, msg_type_url, when set, will query only grants matching given msg type. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an pagination for the request. |






<a name="cosmos.authz.v1beta1.QueryGrantsResponse"></a>

### QueryGrantsResponse
QueryGrantsResponse is the response type for the Query/Authorizations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grants` | [Grant](#cosmos.authz.v1beta1.Grant) | repeated | authorizations is a list of grants granted for grantee by granter. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines an pagination for the response. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.authz.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Grants` | [QueryGrantsRequest](#cosmos.authz.v1beta1.QueryGrantsRequest) | [QueryGrantsResponse](#cosmos.authz.v1beta1.QueryGrantsResponse) | Returns list of `Authorization`, granted to the grantee by the granter. | GET|/cosmos/authz/v1beta1/grants|
| `GranterGrants` | [QueryGranterGrantsRequest](#cosmos.authz.v1beta1.QueryGranterGrantsRequest) | [QueryGranterGrantsResponse](#cosmos.authz.v1beta1.QueryGranterGrantsResponse) | GranterGrants returns list of `GrantAuthorization`, granted by granter.

Since: cosmos-sdk 0.45.2 | GET|/cosmos/authz/v1beta1/grants/granter/{granter}|
| `GranteeGrants` | [QueryGranteeGrantsRequest](#cosmos.authz.v1beta1.QueryGranteeGrantsRequest) | [QueryGranteeGrantsResponse](#cosmos.authz.v1beta1.QueryGranteeGrantsResponse) | GranteeGrants returns a list of `GrantAuthorization` by grantee.

Since: cosmos-sdk 0.45.2 | GET|/cosmos/authz/v1beta1/grants/grantee/{grantee}|

 <!-- end services -->



<a name="cosmos/base/abci/v1beta1/abci.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/base/abci/v1beta1/abci.proto



<a name="cosmos.base.abci.v1beta1.ABCIMessageLog"></a>

### ABCIMessageLog
ABCIMessageLog defines a structure containing an indexed tx ABCI message log.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `msg_index` | [uint32](#uint32) |  |  |
| `log` | [string](#string) |  |  |
| `events` | [StringEvent](#cosmos.base.abci.v1beta1.StringEvent) | repeated | Events contains a slice of Event objects that were emitted during some execution. |






<a name="cosmos.base.abci.v1beta1.Attribute"></a>

### Attribute
Attribute defines an attribute wrapper where the key and value are
strings instead of raw bytes.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |






<a name="cosmos.base.abci.v1beta1.GasInfo"></a>

### GasInfo
GasInfo defines tx execution gas context.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gas_wanted` | [uint64](#uint64) |  | GasWanted is the maximum units of work we allow this tx to perform. |
| `gas_used` | [uint64](#uint64) |  | GasUsed is the amount of gas actually consumed. |






<a name="cosmos.base.abci.v1beta1.MsgData"></a>

### MsgData
MsgData defines the data returned in a Result object during message
execution.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `msg_type` | [string](#string) |  |  |
| `data` | [bytes](#bytes) |  |  |






<a name="cosmos.base.abci.v1beta1.Result"></a>

### Result
Result is the union of ResponseFormat and ResponseCheckTx.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  | Data is any data returned from message or handler execution. It MUST be length prefixed in order to separate data from multiple message executions. |
| `log` | [string](#string) |  | Log contains the log information from message or handler execution. |
| `events` | [ostracon.abci.Event](#ostracon.abci.Event) | repeated | Events contains a slice of Event objects that were emitted during message or handler execution. |






<a name="cosmos.base.abci.v1beta1.SearchTxsResult"></a>

### SearchTxsResult
SearchTxsResult defines a structure for querying txs pageable


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `total_count` | [uint64](#uint64) |  | Count of all txs |
| `count` | [uint64](#uint64) |  | Count of txs in current page |
| `page_number` | [uint64](#uint64) |  | Index of current page, start from 1 |
| `page_total` | [uint64](#uint64) |  | Count of total pages |
| `limit` | [uint64](#uint64) |  | Max count txs per page |
| `txs` | [TxResponse](#cosmos.base.abci.v1beta1.TxResponse) | repeated | List of txs in current page |






<a name="cosmos.base.abci.v1beta1.SimulationResponse"></a>

### SimulationResponse
SimulationResponse defines the response generated when a transaction is
successfully simulated.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gas_info` | [GasInfo](#cosmos.base.abci.v1beta1.GasInfo) |  |  |
| `result` | [Result](#cosmos.base.abci.v1beta1.Result) |  |  |






<a name="cosmos.base.abci.v1beta1.StringEvent"></a>

### StringEvent
StringEvent defines en Event object wrapper where all the attributes
contain key/value pairs that are strings instead of raw bytes.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [string](#string) |  |  |
| `attributes` | [Attribute](#cosmos.base.abci.v1beta1.Attribute) | repeated |  |






<a name="cosmos.base.abci.v1beta1.TxMsgData"></a>

### TxMsgData
TxMsgData defines a list of MsgData. A transaction will have a MsgData object
for each message.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [MsgData](#cosmos.base.abci.v1beta1.MsgData) | repeated |  |






<a name="cosmos.base.abci.v1beta1.TxResponse"></a>

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
| `logs` | [ABCIMessageLog](#cosmos.base.abci.v1beta1.ABCIMessageLog) | repeated | The output of the application's logger (typed). May be non-deterministic. |
| `info` | [string](#string) |  | Additional information. May be non-deterministic. |
| `gas_wanted` | [int64](#int64) |  | Amount of gas requested for transaction. |
| `gas_used` | [int64](#int64) |  | Amount of gas consumed by transaction. |
| `tx` | [google.protobuf.Any](#google.protobuf.Any) |  | The request transaction bytes. |
| `timestamp` | [string](#string) |  | Time of the previous block. For heights > 1, it's the weighted median of the timestamps of the valid votes in the block.LastCommit. For height == 1, it's genesis time. |
| `events` | [ostracon.abci.Event](#ostracon.abci.Event) | repeated | Events defines all the events emitted by processing a transaction. Note, these events include those emitted by processing all the messages and those emitted from the ante handler. Whereas Logs contains the events, with additional metadata, emitted only by processing the messages.

Since: cosmos-sdk 0.42.11, 0.44.5, 0.45 |
| `index` | [uint32](#uint32) |  | The transaction index within block |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/authz/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/authz/v1beta1/tx.proto
Since: cosmos-sdk 0.43


<a name="cosmos.authz.v1beta1.MsgExec"></a>

### MsgExec
MsgExec attempts to execute the provided messages using
authorizations granted to the grantee. Each message should have only
one signer corresponding to the granter of the authorization.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grantee` | [string](#string) |  |  |
| `msgs` | [google.protobuf.Any](#google.protobuf.Any) | repeated | Authorization Msg requests to execute. Each msg must implement Authorization interface The x/authz will try to find a grant matching (msg.signers[0], grantee, MsgTypeURL(msg)) triple and validate it. |






<a name="cosmos.authz.v1beta1.MsgExecResponse"></a>

### MsgExecResponse
MsgExecResponse defines the Msg/MsgExecResponse response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `results` | [bytes](#bytes) | repeated |  |






<a name="cosmos.authz.v1beta1.MsgGrant"></a>

### MsgGrant
MsgGrant is a request type for Grant method. It declares authorization to the grantee
on behalf of the granter with the provided expiration time.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  |  |
| `grantee` | [string](#string) |  |  |
| `grant` | [Grant](#cosmos.authz.v1beta1.Grant) |  |  |






<a name="cosmos.authz.v1beta1.MsgGrantResponse"></a>

### MsgGrantResponse
MsgGrantResponse defines the Msg/MsgGrant response type.






<a name="cosmos.authz.v1beta1.MsgRevoke"></a>

### MsgRevoke
MsgRevoke revokes any authorization with the provided sdk.Msg type on the
granter's account with that has been granted to the grantee.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  |  |
| `grantee` | [string](#string) |  |  |
| `msg_type_url` | [string](#string) |  |  |






<a name="cosmos.authz.v1beta1.MsgRevokeResponse"></a>

### MsgRevokeResponse
MsgRevokeResponse defines the Msg/MsgRevokeResponse response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.authz.v1beta1.Msg"></a>

### Msg
Msg defines the authz Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Grant` | [MsgGrant](#cosmos.authz.v1beta1.MsgGrant) | [MsgGrantResponse](#cosmos.authz.v1beta1.MsgGrantResponse) | Grant grants the provided authorization to the grantee on the granter's account with the provided expiration time. If there is already a grant for the given (granter, grantee, Authorization) triple, then the grant will be overwritten. | |
| `Exec` | [MsgExec](#cosmos.authz.v1beta1.MsgExec) | [MsgExecResponse](#cosmos.authz.v1beta1.MsgExecResponse) | Exec attempts to execute the provided messages using authorizations granted to the grantee. Each message should have only one signer corresponding to the granter of the authorization. | |
| `Revoke` | [MsgRevoke](#cosmos.authz.v1beta1.MsgRevoke) | [MsgRevokeResponse](#cosmos.authz.v1beta1.MsgRevokeResponse) | Revoke revokes any authorization corresponding to the provided method name on the granter's account that has been granted to the grantee. | |

 <!-- end services -->



<a name="cosmos/base/v1beta1/coin.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/base/v1beta1/coin.proto



<a name="cosmos.base.v1beta1.Coin"></a>

### Coin
Coin defines a token with a denomination and an amount.

NOTE: The amount field is an Int which implements the custom method
signatures required by gogoproto.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `amount` | [string](#string) |  |  |






<a name="cosmos.base.v1beta1.DecCoin"></a>

### DecCoin
DecCoin defines a token with a denomination and a decimal amount.

NOTE: The amount field is an Dec which implements the custom method
signatures required by gogoproto.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `amount` | [string](#string) |  |  |






<a name="cosmos.base.v1beta1.DecProto"></a>

### DecProto
DecProto defines a Protobuf wrapper around a Dec object.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dec` | [string](#string) |  |  |






<a name="cosmos.base.v1beta1.IntProto"></a>

### IntProto
IntProto defines a Protobuf wrapper around an Int object.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `int` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/bank/v1beta1/authz.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/bank/v1beta1/authz.proto



<a name="cosmos.bank.v1beta1.SendAuthorization"></a>

### SendAuthorization
SendAuthorization allows the grantee to spend up to spend_limit coins from
the granter's account.

Since: cosmos-sdk 0.43


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `spend_limit` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/bank/v1beta1/bank.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/bank/v1beta1/bank.proto



<a name="cosmos.bank.v1beta1.DenomUnit"></a>

### DenomUnit
DenomUnit represents a struct that describes a given
denomination unit of the basic token.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  | denom represents the string name of the given denom unit (e.g uatom). |
| `exponent` | [uint32](#uint32) |  | exponent represents power of 10 exponent that one must raise the base_denom to in order to equal the given DenomUnit's denom 1 denom = 1^exponent base_denom (e.g. with a base_denom of uatom, one can create a DenomUnit of 'atom' with exponent = 6, thus: 1 atom = 10^6 uatom). |
| `aliases` | [string](#string) | repeated | aliases is a list of string aliases for the given denom |






<a name="cosmos.bank.v1beta1.Input"></a>

### Input
Input models transaction input.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |
| `coins` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="cosmos.bank.v1beta1.Metadata"></a>

### Metadata
Metadata represents a struct that describes
a basic token.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `description` | [string](#string) |  |  |
| `denom_units` | [DenomUnit](#cosmos.bank.v1beta1.DenomUnit) | repeated | denom_units represents the list of DenomUnit's for a given coin |
| `base` | [string](#string) |  | base represents the base denom (should be the DenomUnit with exponent = 0). |
| `display` | [string](#string) |  | display indicates the suggested denom that should be displayed in clients. |
| `name` | [string](#string) |  | name defines the name of the token (eg: Cosmos Atom)

Since: cosmos-sdk 0.43 |
| `symbol` | [string](#string) |  | symbol is the token symbol usually shown on exchanges (eg: ATOM). This can be the same as the display.

Since: cosmos-sdk 0.43 |






<a name="cosmos.bank.v1beta1.Output"></a>

### Output
Output models transaction outputs.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |
| `coins` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="cosmos.bank.v1beta1.Params"></a>

### Params
Params defines the parameters for the bank module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `send_enabled` | [SendEnabled](#cosmos.bank.v1beta1.SendEnabled) | repeated |  |
| `default_send_enabled` | [bool](#bool) |  |  |






<a name="cosmos.bank.v1beta1.SendEnabled"></a>

### SendEnabled
SendEnabled maps coin denom to a send_enabled status (whether a denom is
sendable).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `enabled` | [bool](#bool) |  |  |






<a name="cosmos.bank.v1beta1.Supply"></a>

### Supply
Supply represents a struct that passively keeps track of the total supply
amounts in the network.
This message is deprecated now that supply is indexed by denom.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `total` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/bank/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/bank/v1beta1/genesis.proto



<a name="cosmos.bank.v1beta1.Balance"></a>

### Balance
Balance defines an account address and balance pair used in the bank module's
genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address of the balance holder. |
| `coins` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | coins defines the different coins this balance holds. |






<a name="cosmos.bank.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the bank module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmos.bank.v1beta1.Params) |  | params defines all the paramaters of the module. |
| `balances` | [Balance](#cosmos.bank.v1beta1.Balance) | repeated | balances is an array containing the balances of all the accounts. |
| `supply` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | supply represents the total supply. If it is left empty, then supply will be calculated based on the provided balances. Otherwise, it will be used to validate that the sum of the balances equals this amount. |
| `denom_metadata` | [Metadata](#cosmos.bank.v1beta1.Metadata) | repeated | denom_metadata defines the metadata of the differents coins. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/bank/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/bank/v1beta1/query.proto



<a name="cosmos.bank.v1beta1.QueryAllBalancesRequest"></a>

### QueryAllBalancesRequest
QueryBalanceRequest is the request type for the Query/AllBalances RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address to query balances for. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.bank.v1beta1.QueryAllBalancesResponse"></a>

### QueryAllBalancesResponse
QueryAllBalancesResponse is the response type for the Query/AllBalances RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `balances` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | balances is the balances of all the coins. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.bank.v1beta1.QueryBalanceRequest"></a>

### QueryBalanceRequest
QueryBalanceRequest is the request type for the Query/Balance RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address to query balances for. |
| `denom` | [string](#string) |  | denom is the coin denom to query balances for. |






<a name="cosmos.bank.v1beta1.QueryBalanceResponse"></a>

### QueryBalanceResponse
QueryBalanceResponse is the response type for the Query/Balance RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `balance` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | balance is the balance of the coin. |






<a name="cosmos.bank.v1beta1.QueryDenomMetadataRequest"></a>

### QueryDenomMetadataRequest
QueryDenomMetadataRequest is the request type for the Query/DenomMetadata RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  | denom is the coin denom to query the metadata for. |






<a name="cosmos.bank.v1beta1.QueryDenomMetadataResponse"></a>

### QueryDenomMetadataResponse
QueryDenomMetadataResponse is the response type for the Query/DenomMetadata RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `metadata` | [Metadata](#cosmos.bank.v1beta1.Metadata) |  | metadata describes and provides all the client information for the requested token. |






<a name="cosmos.bank.v1beta1.QueryDenomsMetadataRequest"></a>

### QueryDenomsMetadataRequest
QueryDenomsMetadataRequest is the request type for the Query/DenomsMetadata RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.bank.v1beta1.QueryDenomsMetadataResponse"></a>

### QueryDenomsMetadataResponse
QueryDenomsMetadataResponse is the response type for the Query/DenomsMetadata RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `metadatas` | [Metadata](#cosmos.bank.v1beta1.Metadata) | repeated | metadata provides the client information for all the registered tokens. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.bank.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest defines the request type for querying x/bank parameters.






<a name="cosmos.bank.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse defines the response type for querying x/bank parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmos.bank.v1beta1.Params) |  |  |






<a name="cosmos.bank.v1beta1.QuerySpendableBalancesRequest"></a>

### QuerySpendableBalancesRequest
QuerySpendableBalancesRequest defines the gRPC request structure for querying
an account's spendable balances.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address to query spendable balances for. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.bank.v1beta1.QuerySpendableBalancesResponse"></a>

### QuerySpendableBalancesResponse
QuerySpendableBalancesResponse defines the gRPC response structure for querying
an account's spendable balances.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `balances` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | balances is the spendable balances of all the coins. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.bank.v1beta1.QuerySupplyOfRequest"></a>

### QuerySupplyOfRequest
QuerySupplyOfRequest is the request type for the Query/SupplyOf RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  | denom is the coin denom to query balances for. |






<a name="cosmos.bank.v1beta1.QuerySupplyOfResponse"></a>

### QuerySupplyOfResponse
QuerySupplyOfResponse is the response type for the Query/SupplyOf RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | amount is the supply of the coin. |






<a name="cosmos.bank.v1beta1.QueryTotalSupplyRequest"></a>

### QueryTotalSupplyRequest
QueryTotalSupplyRequest is the request type for the Query/TotalSupply RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request.

Since: cosmos-sdk 0.43 |






<a name="cosmos.bank.v1beta1.QueryTotalSupplyResponse"></a>

### QueryTotalSupplyResponse
QueryTotalSupplyResponse is the response type for the Query/TotalSupply RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `supply` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | supply is the supply of the coins |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response.

Since: cosmos-sdk 0.43 |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.bank.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Balance` | [QueryBalanceRequest](#cosmos.bank.v1beta1.QueryBalanceRequest) | [QueryBalanceResponse](#cosmos.bank.v1beta1.QueryBalanceResponse) | Balance queries the balance of a single coin for a single account. | GET|/cosmos/bank/v1beta1/balances/{address}/by_denom|
| `AllBalances` | [QueryAllBalancesRequest](#cosmos.bank.v1beta1.QueryAllBalancesRequest) | [QueryAllBalancesResponse](#cosmos.bank.v1beta1.QueryAllBalancesResponse) | AllBalances queries the balance of all coins for a single account. | GET|/cosmos/bank/v1beta1/balances/{address}|
| `SpendableBalances` | [QuerySpendableBalancesRequest](#cosmos.bank.v1beta1.QuerySpendableBalancesRequest) | [QuerySpendableBalancesResponse](#cosmos.bank.v1beta1.QuerySpendableBalancesResponse) | SpendableBalances queries the spenable balance of all coins for a single account. | GET|/cosmos/bank/v1beta1/spendable_balances/{address}|
| `TotalSupply` | [QueryTotalSupplyRequest](#cosmos.bank.v1beta1.QueryTotalSupplyRequest) | [QueryTotalSupplyResponse](#cosmos.bank.v1beta1.QueryTotalSupplyResponse) | TotalSupply queries the total supply of all coins. | GET|/cosmos/bank/v1beta1/supply|
| `SupplyOf` | [QuerySupplyOfRequest](#cosmos.bank.v1beta1.QuerySupplyOfRequest) | [QuerySupplyOfResponse](#cosmos.bank.v1beta1.QuerySupplyOfResponse) | SupplyOf queries the supply of a single coin. | GET|/cosmos/bank/v1beta1/supply/{denom}|
| `Params` | [QueryParamsRequest](#cosmos.bank.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#cosmos.bank.v1beta1.QueryParamsResponse) | Params queries the parameters of x/bank module. | GET|/cosmos/bank/v1beta1/params|
| `DenomMetadata` | [QueryDenomMetadataRequest](#cosmos.bank.v1beta1.QueryDenomMetadataRequest) | [QueryDenomMetadataResponse](#cosmos.bank.v1beta1.QueryDenomMetadataResponse) | DenomsMetadata queries the client metadata of a given coin denomination. | GET|/cosmos/bank/v1beta1/denoms_metadata/{denom}|
| `DenomsMetadata` | [QueryDenomsMetadataRequest](#cosmos.bank.v1beta1.QueryDenomsMetadataRequest) | [QueryDenomsMetadataResponse](#cosmos.bank.v1beta1.QueryDenomsMetadataResponse) | DenomsMetadata queries the client metadata for all registered coin denominations. | GET|/cosmos/bank/v1beta1/denoms_metadata|

 <!-- end services -->



<a name="cosmos/bank/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/bank/v1beta1/tx.proto



<a name="cosmos.bank.v1beta1.MsgMultiSend"></a>

### MsgMultiSend
MsgMultiSend represents an arbitrary multi-in, multi-out send message.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `inputs` | [Input](#cosmos.bank.v1beta1.Input) | repeated |  |
| `outputs` | [Output](#cosmos.bank.v1beta1.Output) | repeated |  |






<a name="cosmos.bank.v1beta1.MsgMultiSendResponse"></a>

### MsgMultiSendResponse
MsgMultiSendResponse defines the Msg/MultiSend response type.






<a name="cosmos.bank.v1beta1.MsgSend"></a>

### MsgSend
MsgSend represents a message to send coins from one account to another.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `from_address` | [string](#string) |  |  |
| `to_address` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="cosmos.bank.v1beta1.MsgSendResponse"></a>

### MsgSendResponse
MsgSendResponse defines the Msg/Send response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.bank.v1beta1.Msg"></a>

### Msg
Msg defines the bank Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Send` | [MsgSend](#cosmos.bank.v1beta1.MsgSend) | [MsgSendResponse](#cosmos.bank.v1beta1.MsgSendResponse) | Send defines a method for sending coins from one account to another account. | |
| `MultiSend` | [MsgMultiSend](#cosmos.bank.v1beta1.MsgMultiSend) | [MsgMultiSendResponse](#cosmos.bank.v1beta1.MsgMultiSendResponse) | MultiSend defines a method for sending coins from some accounts to other accounts. | |

 <!-- end services -->



<a name="cosmos/base/kv/v1beta1/kv.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/base/kv/v1beta1/kv.proto



<a name="cosmos.base.kv.v1beta1.Pair"></a>

### Pair
Pair defines a key/value bytes tuple.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  |  |
| `value` | [bytes](#bytes) |  |  |






<a name="cosmos.base.kv.v1beta1.Pairs"></a>

### Pairs
Pairs defines a repeated slice of Pair objects.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pairs` | [Pair](#cosmos.base.kv.v1beta1.Pair) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/base/node/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/base/node/v1beta1/query.proto



<a name="cosmos.base.node.v1beta1.ConfigRequest"></a>

### ConfigRequest
ConfigRequest defines the request structure for the Config gRPC query.






<a name="cosmos.base.node.v1beta1.ConfigResponse"></a>

### ConfigResponse
ConfigResponse defines the response structure for the Config gRPC query.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `minimum_gas_price` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.base.node.v1beta1.Service"></a>

### Service
Service defines the gRPC querier service for node related queries.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Config` | [ConfigRequest](#cosmos.base.node.v1beta1.ConfigRequest) | [ConfigResponse](#cosmos.base.node.v1beta1.ConfigResponse) | Config queries for the operator configuration. | GET|/cosmos/base/node/v1beta1/config|

 <!-- end services -->



<a name="cosmos/base/reflection/v1beta1/reflection.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/base/reflection/v1beta1/reflection.proto



<a name="cosmos.base.reflection.v1beta1.ListAllInterfacesRequest"></a>

### ListAllInterfacesRequest
ListAllInterfacesRequest is the request type of the ListAllInterfaces RPC.






<a name="cosmos.base.reflection.v1beta1.ListAllInterfacesResponse"></a>

### ListAllInterfacesResponse
ListAllInterfacesResponse is the response type of the ListAllInterfaces RPC.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `interface_names` | [string](#string) | repeated | interface_names is an array of all the registered interfaces. |






<a name="cosmos.base.reflection.v1beta1.ListImplementationsRequest"></a>

### ListImplementationsRequest
ListImplementationsRequest is the request type of the ListImplementations
RPC.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `interface_name` | [string](#string) |  | interface_name defines the interface to query the implementations for. |






<a name="cosmos.base.reflection.v1beta1.ListImplementationsResponse"></a>

### ListImplementationsResponse
ListImplementationsResponse is the response type of the ListImplementations
RPC.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `implementation_message_names` | [string](#string) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.base.reflection.v1beta1.ReflectionService"></a>

### ReflectionService
ReflectionService defines a service for interface reflection.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `ListAllInterfaces` | [ListAllInterfacesRequest](#cosmos.base.reflection.v1beta1.ListAllInterfacesRequest) | [ListAllInterfacesResponse](#cosmos.base.reflection.v1beta1.ListAllInterfacesResponse) | ListAllInterfaces lists all the interfaces registered in the interface registry. | GET|/cosmos/base/reflection/v1beta1/interfaces|
| `ListImplementations` | [ListImplementationsRequest](#cosmos.base.reflection.v1beta1.ListImplementationsRequest) | [ListImplementationsResponse](#cosmos.base.reflection.v1beta1.ListImplementationsResponse) | ListImplementations list all the concrete types that implement a given interface. | GET|/cosmos/base/reflection/v1beta1/interfaces/{interface_name}/implementations|

 <!-- end services -->



<a name="cosmos/base/reflection/v2alpha1/reflection.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/base/reflection/v2alpha1/reflection.proto
Since: cosmos-sdk 0.43


<a name="cosmos.base.reflection.v2alpha1.AppDescriptor"></a>

### AppDescriptor
AppDescriptor describes a cosmos-sdk based application


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authn` | [AuthnDescriptor](#cosmos.base.reflection.v2alpha1.AuthnDescriptor) |  | AuthnDescriptor provides information on how to authenticate transactions on the application NOTE: experimental and subject to change in future releases. |
| `chain` | [ChainDescriptor](#cosmos.base.reflection.v2alpha1.ChainDescriptor) |  | chain provides the chain descriptor |
| `codec` | [CodecDescriptor](#cosmos.base.reflection.v2alpha1.CodecDescriptor) |  | codec provides metadata information regarding codec related types |
| `configuration` | [ConfigurationDescriptor](#cosmos.base.reflection.v2alpha1.ConfigurationDescriptor) |  | configuration provides metadata information regarding the sdk.Config type |
| `query_services` | [QueryServicesDescriptor](#cosmos.base.reflection.v2alpha1.QueryServicesDescriptor) |  | query_services provides metadata information regarding the available queriable endpoints |
| `tx` | [TxDescriptor](#cosmos.base.reflection.v2alpha1.TxDescriptor) |  | tx provides metadata information regarding how to send transactions to the given application |






<a name="cosmos.base.reflection.v2alpha1.AuthnDescriptor"></a>

### AuthnDescriptor
AuthnDescriptor provides information on how to sign transactions without relying
on the online RPCs GetTxMetadata and CombineUnsignedTxAndSignatures


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sign_modes` | [SigningModeDescriptor](#cosmos.base.reflection.v2alpha1.SigningModeDescriptor) | repeated | sign_modes defines the supported signature algorithm |






<a name="cosmos.base.reflection.v2alpha1.ChainDescriptor"></a>

### ChainDescriptor
ChainDescriptor describes chain information of the application


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  | id is the chain id |






<a name="cosmos.base.reflection.v2alpha1.CodecDescriptor"></a>

### CodecDescriptor
CodecDescriptor describes the registered interfaces and provides metadata information on the types


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `interfaces` | [InterfaceDescriptor](#cosmos.base.reflection.v2alpha1.InterfaceDescriptor) | repeated | interfaces is a list of the registerted interfaces descriptors |






<a name="cosmos.base.reflection.v2alpha1.ConfigurationDescriptor"></a>

### ConfigurationDescriptor
ConfigurationDescriptor contains metadata information on the sdk.Config


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bech32_account_address_prefix` | [string](#string) |  | bech32_account_address_prefix is the account address prefix |






<a name="cosmos.base.reflection.v2alpha1.GetAuthnDescriptorRequest"></a>

### GetAuthnDescriptorRequest
GetAuthnDescriptorRequest is the request used for the GetAuthnDescriptor RPC






<a name="cosmos.base.reflection.v2alpha1.GetAuthnDescriptorResponse"></a>

### GetAuthnDescriptorResponse
GetAuthnDescriptorResponse is the response returned by the GetAuthnDescriptor RPC


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authn` | [AuthnDescriptor](#cosmos.base.reflection.v2alpha1.AuthnDescriptor) |  | authn describes how to authenticate to the application when sending transactions |






<a name="cosmos.base.reflection.v2alpha1.GetChainDescriptorRequest"></a>

### GetChainDescriptorRequest
GetChainDescriptorRequest is the request used for the GetChainDescriptor RPC






<a name="cosmos.base.reflection.v2alpha1.GetChainDescriptorResponse"></a>

### GetChainDescriptorResponse
GetChainDescriptorResponse is the response returned by the GetChainDescriptor RPC


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain` | [ChainDescriptor](#cosmos.base.reflection.v2alpha1.ChainDescriptor) |  | chain describes application chain information |






<a name="cosmos.base.reflection.v2alpha1.GetCodecDescriptorRequest"></a>

### GetCodecDescriptorRequest
GetCodecDescriptorRequest is the request used for the GetCodecDescriptor RPC






<a name="cosmos.base.reflection.v2alpha1.GetCodecDescriptorResponse"></a>

### GetCodecDescriptorResponse
GetCodecDescriptorResponse is the response returned by the GetCodecDescriptor RPC


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `codec` | [CodecDescriptor](#cosmos.base.reflection.v2alpha1.CodecDescriptor) |  | codec describes the application codec such as registered interfaces and implementations |






<a name="cosmos.base.reflection.v2alpha1.GetConfigurationDescriptorRequest"></a>

### GetConfigurationDescriptorRequest
GetConfigurationDescriptorRequest is the request used for the GetConfigurationDescriptor RPC






<a name="cosmos.base.reflection.v2alpha1.GetConfigurationDescriptorResponse"></a>

### GetConfigurationDescriptorResponse
GetConfigurationDescriptorResponse is the response returned by the GetConfigurationDescriptor RPC


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `config` | [ConfigurationDescriptor](#cosmos.base.reflection.v2alpha1.ConfigurationDescriptor) |  | config describes the application's sdk.Config |






<a name="cosmos.base.reflection.v2alpha1.GetQueryServicesDescriptorRequest"></a>

### GetQueryServicesDescriptorRequest
GetQueryServicesDescriptorRequest is the request used for the GetQueryServicesDescriptor RPC






<a name="cosmos.base.reflection.v2alpha1.GetQueryServicesDescriptorResponse"></a>

### GetQueryServicesDescriptorResponse
GetQueryServicesDescriptorResponse is the response returned by the GetQueryServicesDescriptor RPC


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `queries` | [QueryServicesDescriptor](#cosmos.base.reflection.v2alpha1.QueryServicesDescriptor) |  | queries provides information on the available queryable services |






<a name="cosmos.base.reflection.v2alpha1.GetTxDescriptorRequest"></a>

### GetTxDescriptorRequest
GetTxDescriptorRequest is the request used for the GetTxDescriptor RPC






<a name="cosmos.base.reflection.v2alpha1.GetTxDescriptorResponse"></a>

### GetTxDescriptorResponse
GetTxDescriptorResponse is the response returned by the GetTxDescriptor RPC


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tx` | [TxDescriptor](#cosmos.base.reflection.v2alpha1.TxDescriptor) |  | tx provides information on msgs that can be forwarded to the application alongside the accepted transaction protobuf type |






<a name="cosmos.base.reflection.v2alpha1.InterfaceAcceptingMessageDescriptor"></a>

### InterfaceAcceptingMessageDescriptor
InterfaceAcceptingMessageDescriptor describes a protobuf message which contains
an interface represented as a google.protobuf.Any


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fullname` | [string](#string) |  | fullname is the protobuf fullname of the type containing the interface |
| `field_descriptor_names` | [string](#string) | repeated | field_descriptor_names is a list of the protobuf name (not fullname) of the field which contains the interface as google.protobuf.Any (the interface is the same, but it can be in multiple fields of the same proto message) |






<a name="cosmos.base.reflection.v2alpha1.InterfaceDescriptor"></a>

### InterfaceDescriptor
InterfaceDescriptor describes the implementation of an interface


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fullname` | [string](#string) |  | fullname is the name of the interface |
| `interface_accepting_messages` | [InterfaceAcceptingMessageDescriptor](#cosmos.base.reflection.v2alpha1.InterfaceAcceptingMessageDescriptor) | repeated | interface_accepting_messages contains information regarding the proto messages which contain the interface as google.protobuf.Any field |
| `interface_implementers` | [InterfaceImplementerDescriptor](#cosmos.base.reflection.v2alpha1.InterfaceImplementerDescriptor) | repeated | interface_implementers is a list of the descriptors of the interface implementers |






<a name="cosmos.base.reflection.v2alpha1.InterfaceImplementerDescriptor"></a>

### InterfaceImplementerDescriptor
InterfaceImplementerDescriptor describes an interface implementer


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fullname` | [string](#string) |  | fullname is the protobuf queryable name of the interface implementer |
| `type_url` | [string](#string) |  | type_url defines the type URL used when marshalling the type as any this is required so we can provide type safe google.protobuf.Any marshalling and unmarshalling, making sure that we don't accept just 'any' type in our interface fields |






<a name="cosmos.base.reflection.v2alpha1.MsgDescriptor"></a>

### MsgDescriptor
MsgDescriptor describes a cosmos-sdk message that can be delivered with a transaction


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `msg_type_url` | [string](#string) |  | msg_type_url contains the TypeURL of a sdk.Msg. |






<a name="cosmos.base.reflection.v2alpha1.QueryMethodDescriptor"></a>

### QueryMethodDescriptor
QueryMethodDescriptor describes a queryable method of a query service
no other info is provided beside method name and tendermint queryable path
because it would be redundant with the grpc reflection service


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | name is the protobuf name (not fullname) of the method |
| `full_query_path` | [string](#string) |  | full_query_path is the path that can be used to query this method via tendermint abci.Query |






<a name="cosmos.base.reflection.v2alpha1.QueryServiceDescriptor"></a>

### QueryServiceDescriptor
QueryServiceDescriptor describes a cosmos-sdk queryable service


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fullname` | [string](#string) |  | fullname is the protobuf fullname of the service descriptor |
| `is_module` | [bool](#bool) |  | is_module describes if this service is actually exposed by an application's module |
| `methods` | [QueryMethodDescriptor](#cosmos.base.reflection.v2alpha1.QueryMethodDescriptor) | repeated | methods provides a list of query service methods |






<a name="cosmos.base.reflection.v2alpha1.QueryServicesDescriptor"></a>

### QueryServicesDescriptor
QueryServicesDescriptor contains the list of cosmos-sdk queriable services


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `query_services` | [QueryServiceDescriptor](#cosmos.base.reflection.v2alpha1.QueryServiceDescriptor) | repeated | query_services is a list of cosmos-sdk QueryServiceDescriptor |






<a name="cosmos.base.reflection.v2alpha1.SigningModeDescriptor"></a>

### SigningModeDescriptor
SigningModeDescriptor provides information on a signing flow of the application
NOTE(fdymylja): here we could go as far as providing an entire flow on how
to sign a message given a SigningModeDescriptor, but it's better to think about
this another time


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | name defines the unique name of the signing mode |
| `number` | [int32](#int32) |  | number is the unique int32 identifier for the sign_mode enum |
| `authn_info_provider_method_fullname` | [string](#string) |  | authn_info_provider_method_fullname defines the fullname of the method to call to get the metadata required to authenticate using the provided sign_modes |






<a name="cosmos.base.reflection.v2alpha1.TxDescriptor"></a>

### TxDescriptor
TxDescriptor describes the accepted transaction type


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fullname` | [string](#string) |  | fullname is the protobuf fullname of the raw transaction type (for instance the tx.Tx type) it is not meant to support polymorphism of transaction types, it is supposed to be used by reflection clients to understand if they can handle a specific transaction type in an application. |
| `msgs` | [MsgDescriptor](#cosmos.base.reflection.v2alpha1.MsgDescriptor) | repeated | msgs lists the accepted application messages (sdk.Msg) |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.base.reflection.v2alpha1.ReflectionService"></a>

### ReflectionService
ReflectionService defines a service for application reflection.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `GetAuthnDescriptor` | [GetAuthnDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetAuthnDescriptorRequest) | [GetAuthnDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetAuthnDescriptorResponse) | GetAuthnDescriptor returns information on how to authenticate transactions in the application NOTE: this RPC is still experimental and might be subject to breaking changes or removal in future releases of the cosmos-sdk. | GET|/cosmos/base/reflection/v1beta1/app_descriptor/authn|
| `GetChainDescriptor` | [GetChainDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetChainDescriptorRequest) | [GetChainDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetChainDescriptorResponse) | GetChainDescriptor returns the description of the chain | GET|/cosmos/base/reflection/v1beta1/app_descriptor/chain|
| `GetCodecDescriptor` | [GetCodecDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetCodecDescriptorRequest) | [GetCodecDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetCodecDescriptorResponse) | GetCodecDescriptor returns the descriptor of the codec of the application | GET|/cosmos/base/reflection/v1beta1/app_descriptor/codec|
| `GetConfigurationDescriptor` | [GetConfigurationDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetConfigurationDescriptorRequest) | [GetConfigurationDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetConfigurationDescriptorResponse) | GetConfigurationDescriptor returns the descriptor for the sdk.Config of the application | GET|/cosmos/base/reflection/v1beta1/app_descriptor/configuration|
| `GetQueryServicesDescriptor` | [GetQueryServicesDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetQueryServicesDescriptorRequest) | [GetQueryServicesDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetQueryServicesDescriptorResponse) | GetQueryServicesDescriptor returns the available gRPC queryable services of the application | GET|/cosmos/base/reflection/v1beta1/app_descriptor/query_services|
| `GetTxDescriptor` | [GetTxDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetTxDescriptorRequest) | [GetTxDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetTxDescriptorResponse) | GetTxDescriptor returns information on the used transaction object and available msgs that can be used | GET|/cosmos/base/reflection/v1beta1/app_descriptor/tx_descriptor|

 <!-- end services -->



<a name="cosmos/base/snapshots/v1beta1/snapshot.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/base/snapshots/v1beta1/snapshot.proto



<a name="cosmos.base.snapshots.v1beta1.Metadata"></a>

### Metadata
Metadata contains SDK-specific snapshot metadata.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chunk_hashes` | [bytes](#bytes) | repeated | SHA-256 chunk hashes |






<a name="cosmos.base.snapshots.v1beta1.Snapshot"></a>

### Snapshot
Snapshot contains Tendermint state sync snapshot info.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [uint64](#uint64) |  |  |
| `format` | [uint32](#uint32) |  |  |
| `chunks` | [uint32](#uint32) |  |  |
| `hash` | [bytes](#bytes) |  |  |
| `metadata` | [Metadata](#cosmos.base.snapshots.v1beta1.Metadata) |  |  |






<a name="cosmos.base.snapshots.v1beta1.SnapshotExtensionMeta"></a>

### SnapshotExtensionMeta
SnapshotExtensionMeta contains metadata about an external snapshotter.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `format` | [uint32](#uint32) |  |  |






<a name="cosmos.base.snapshots.v1beta1.SnapshotExtensionPayload"></a>

### SnapshotExtensionPayload
SnapshotExtensionPayload contains payloads of an external snapshotter.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `payload` | [bytes](#bytes) |  |  |






<a name="cosmos.base.snapshots.v1beta1.SnapshotIAVLItem"></a>

### SnapshotIAVLItem
SnapshotIAVLItem is an exported IAVL node.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  |  |
| `value` | [bytes](#bytes) |  |  |
| `version` | [int64](#int64) |  | version is block height |
| `height` | [int32](#int32) |  | height is depth of the tree. |






<a name="cosmos.base.snapshots.v1beta1.SnapshotItem"></a>

### SnapshotItem
SnapshotItem is an item contained in a rootmulti.Store snapshot.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `store` | [SnapshotStoreItem](#cosmos.base.snapshots.v1beta1.SnapshotStoreItem) |  |  |
| `iavl` | [SnapshotIAVLItem](#cosmos.base.snapshots.v1beta1.SnapshotIAVLItem) |  |  |
| `extension` | [SnapshotExtensionMeta](#cosmos.base.snapshots.v1beta1.SnapshotExtensionMeta) |  |  |
| `extension_payload` | [SnapshotExtensionPayload](#cosmos.base.snapshots.v1beta1.SnapshotExtensionPayload) |  |  |






<a name="cosmos.base.snapshots.v1beta1.SnapshotStoreItem"></a>

### SnapshotStoreItem
SnapshotStoreItem contains metadata about a snapshotted store.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/base/store/v1beta1/commit_info.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/base/store/v1beta1/commit_info.proto



<a name="cosmos.base.store.v1beta1.CommitID"></a>

### CommitID
CommitID defines the committment information when a specific store is
committed.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `version` | [int64](#int64) |  |  |
| `hash` | [bytes](#bytes) |  |  |






<a name="cosmos.base.store.v1beta1.CommitInfo"></a>

### CommitInfo
CommitInfo defines commit information used by the multi-store when committing
a version/height.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `version` | [int64](#int64) |  |  |
| `store_infos` | [StoreInfo](#cosmos.base.store.v1beta1.StoreInfo) | repeated |  |






<a name="cosmos.base.store.v1beta1.StoreInfo"></a>

### StoreInfo
StoreInfo defines store-specific commit information. It contains a reference
between a store name and the commit ID.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `commit_id` | [CommitID](#cosmos.base.store.v1beta1.CommitID) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/base/store/v1beta1/listening.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/base/store/v1beta1/listening.proto



<a name="cosmos.base.store.v1beta1.StoreKVPair"></a>

### StoreKVPair
StoreKVPair is a KVStore KVPair used for listening to state changes (Sets and Deletes)
It optionally includes the StoreKey for the originating KVStore and a Boolean flag to distinguish between Sets and
Deletes

Since: cosmos-sdk 0.43


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `store_key` | [string](#string) |  | the store key for the KVStore this pair originates from |
| `delete` | [bool](#bool) |  | true indicates a delete operation, false indicates a set operation |
| `key` | [bytes](#bytes) |  |  |
| `value` | [bytes](#bytes) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/capability/v1beta1/capability.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/capability/v1beta1/capability.proto



<a name="cosmos.capability.v1beta1.Capability"></a>

### Capability
Capability defines an implementation of an object capability. The index
provided to a Capability must be globally unique.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `index` | [uint64](#uint64) |  |  |






<a name="cosmos.capability.v1beta1.CapabilityOwners"></a>

### CapabilityOwners
CapabilityOwners defines a set of owners of a single Capability. The set of
owners must be unique.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `owners` | [Owner](#cosmos.capability.v1beta1.Owner) | repeated |  |






<a name="cosmos.capability.v1beta1.Owner"></a>

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



<a name="cosmos/capability/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/capability/v1beta1/genesis.proto



<a name="cosmos.capability.v1beta1.GenesisOwners"></a>

### GenesisOwners
GenesisOwners defines the capability owners with their corresponding index.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `index` | [uint64](#uint64) |  | index is the index of the capability owner. |
| `index_owners` | [CapabilityOwners](#cosmos.capability.v1beta1.CapabilityOwners) |  | index_owners are the owners at the given index. |






<a name="cosmos.capability.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the capability module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `index` | [uint64](#uint64) |  | index is the capability global index. |
| `owners` | [GenesisOwners](#cosmos.capability.v1beta1.GenesisOwners) | repeated | owners represents a map from index to owners of the capability index index key is string to allow amino marshalling. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/crisis/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/crisis/v1beta1/genesis.proto



<a name="cosmos.crisis.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the crisis module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `constant_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | constant_fee is the fee used to verify the invariant in the crisis module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/crisis/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/crisis/v1beta1/tx.proto



<a name="cosmos.crisis.v1beta1.MsgVerifyInvariant"></a>

### MsgVerifyInvariant
MsgVerifyInvariant represents a message to verify a particular invariance.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `invariant_module_name` | [string](#string) |  |  |
| `invariant_route` | [string](#string) |  |  |






<a name="cosmos.crisis.v1beta1.MsgVerifyInvariantResponse"></a>

### MsgVerifyInvariantResponse
MsgVerifyInvariantResponse defines the Msg/VerifyInvariant response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.crisis.v1beta1.Msg"></a>

### Msg
Msg defines the bank Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `VerifyInvariant` | [MsgVerifyInvariant](#cosmos.crisis.v1beta1.MsgVerifyInvariant) | [MsgVerifyInvariantResponse](#cosmos.crisis.v1beta1.MsgVerifyInvariantResponse) | VerifyInvariant defines a method to verify a particular invariance. | |

 <!-- end services -->



<a name="cosmos/crypto/ed25519/keys.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/crypto/ed25519/keys.proto



<a name="cosmos.crypto.ed25519.PrivKey"></a>

### PrivKey
Deprecated: PrivKey defines a ed25519 private key.
NOTE: ed25519 keys must not be used in SDK apps except in a tendermint validator context.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  |  |






<a name="cosmos.crypto.ed25519.PubKey"></a>

### PubKey
PubKey is an ed25519 public key for handling Tendermint keys in SDK.
It's needed for Any serialization and SDK compatibility.
It must not be used in a non Tendermint key context because it doesn't implement
ADR-28. Nevertheless, you will like to use ed25519 in app user level
then you must create a new proto message and follow ADR-28 for Address construction.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/crypto/multisig/keys.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/crypto/multisig/keys.proto



<a name="cosmos.crypto.multisig.LegacyAminoPubKey"></a>

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



<a name="cosmos/crypto/multisig/v1beta1/multisig.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/crypto/multisig/v1beta1/multisig.proto



<a name="cosmos.crypto.multisig.v1beta1.CompactBitArray"></a>

### CompactBitArray
CompactBitArray is an implementation of a space efficient bit array.
This is used to ensure that the encoded data takes up a minimal amount of
space after proto encoding.
This is not thread safe, and is not intended for concurrent usage.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `extra_bits_stored` | [uint32](#uint32) |  |  |
| `elems` | [bytes](#bytes) |  |  |






<a name="cosmos.crypto.multisig.v1beta1.MultiSignature"></a>

### MultiSignature
MultiSignature wraps the signatures from a multisig.LegacyAminoPubKey.
See cosmos.tx.v1beta1.ModeInfo.Multi for how to specify which signers
signed and with which modes.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signatures` | [bytes](#bytes) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/crypto/secp256k1/keys.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/crypto/secp256k1/keys.proto



<a name="cosmos.crypto.secp256k1.PrivKey"></a>

### PrivKey
PrivKey defines a secp256k1 private key.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  |  |






<a name="cosmos.crypto.secp256k1.PubKey"></a>

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



<a name="cosmos/crypto/secp256r1/keys.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/crypto/secp256r1/keys.proto
Since: cosmos-sdk 0.43


<a name="cosmos.crypto.secp256r1.PrivKey"></a>

### PrivKey
PrivKey defines a secp256r1 ECDSA private key.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `secret` | [bytes](#bytes) |  | secret number serialized using big-endian encoding |






<a name="cosmos.crypto.secp256r1.PubKey"></a>

### PubKey
PubKey defines a secp256r1 ECDSA public key.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  | Point on secp256r1 curve in a compressed representation as specified in section 4.3.6 of ANSI X9.62: https://webstore.ansi.org/standards/ascx9/ansix9621998 |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/distribution/v1beta1/distribution.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/distribution/v1beta1/distribution.proto



<a name="cosmos.distribution.v1beta1.CommunityPoolSpendProposal"></a>

### CommunityPoolSpendProposal
CommunityPoolSpendProposal details a proposal for use of community funds,
together with how many coins are proposed to be spent, and to which
recipient account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `recipient` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="cosmos.distribution.v1beta1.CommunityPoolSpendProposalWithDeposit"></a>

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






<a name="cosmos.distribution.v1beta1.DelegationDelegatorReward"></a>

### DelegationDelegatorReward
DelegationDelegatorReward represents the properties
of a delegator's delegation reward.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  |  |
| `reward` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |






<a name="cosmos.distribution.v1beta1.DelegatorStartingInfo"></a>

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






<a name="cosmos.distribution.v1beta1.FeePool"></a>

### FeePool
FeePool is the global fee pool for distribution.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `community_pool` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |






<a name="cosmos.distribution.v1beta1.Params"></a>

### Params
Params defines the set of params for the distribution module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `community_tax` | [string](#string) |  |  |
| `base_proposer_reward` | [string](#string) |  |  |
| `bonus_proposer_reward` | [string](#string) |  |  |
| `withdraw_addr_enabled` | [bool](#bool) |  |  |






<a name="cosmos.distribution.v1beta1.ValidatorAccumulatedCommission"></a>

### ValidatorAccumulatedCommission
ValidatorAccumulatedCommission represents accumulated commission
for a validator kept as a running counter, can be withdrawn at any time.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commission` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |






<a name="cosmos.distribution.v1beta1.ValidatorCurrentRewards"></a>

### ValidatorCurrentRewards
ValidatorCurrentRewards represents current rewards and current
period for a validator kept as a running counter and incremented
each block as long as the validator's tokens remain constant.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rewards` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |
| `period` | [uint64](#uint64) |  |  |






<a name="cosmos.distribution.v1beta1.ValidatorHistoricalRewards"></a>

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
| `cumulative_reward_ratio` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |
| `reference_count` | [uint32](#uint32) |  |  |






<a name="cosmos.distribution.v1beta1.ValidatorOutstandingRewards"></a>

### ValidatorOutstandingRewards
ValidatorOutstandingRewards represents outstanding (un-withdrawn) rewards
for a validator inexpensive to track, allows simple sanity checks.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rewards` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |






<a name="cosmos.distribution.v1beta1.ValidatorSlashEvent"></a>

### ValidatorSlashEvent
ValidatorSlashEvent represents a validator slash event.
Height is implicit within the store key.
This is needed to calculate appropriate amount of staking tokens
for delegations which are withdrawn after a slash has occurred.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_period` | [uint64](#uint64) |  |  |
| `fraction` | [string](#string) |  |  |






<a name="cosmos.distribution.v1beta1.ValidatorSlashEvents"></a>

### ValidatorSlashEvents
ValidatorSlashEvents is a collection of ValidatorSlashEvent messages.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_slash_events` | [ValidatorSlashEvent](#cosmos.distribution.v1beta1.ValidatorSlashEvent) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/distribution/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/distribution/v1beta1/genesis.proto



<a name="cosmos.distribution.v1beta1.DelegatorStartingInfoRecord"></a>

### DelegatorStartingInfoRecord
DelegatorStartingInfoRecord used for import / export via genesis json.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address is the address of the delegator. |
| `validator_address` | [string](#string) |  | validator_address is the address of the validator. |
| `starting_info` | [DelegatorStartingInfo](#cosmos.distribution.v1beta1.DelegatorStartingInfo) |  | starting_info defines the starting info of a delegator. |






<a name="cosmos.distribution.v1beta1.DelegatorWithdrawInfo"></a>

### DelegatorWithdrawInfo
DelegatorWithdrawInfo is the address for where distributions rewards are
withdrawn to by default this struct is only used at genesis to feed in
default withdraw addresses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address is the address of the delegator. |
| `withdraw_address` | [string](#string) |  | withdraw_address is the address to withdraw the delegation rewards to. |






<a name="cosmos.distribution.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the distribution module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmos.distribution.v1beta1.Params) |  | params defines all the paramaters of the module. |
| `fee_pool` | [FeePool](#cosmos.distribution.v1beta1.FeePool) |  | fee_pool defines the fee pool at genesis. |
| `delegator_withdraw_infos` | [DelegatorWithdrawInfo](#cosmos.distribution.v1beta1.DelegatorWithdrawInfo) | repeated | fee_pool defines the delegator withdraw infos at genesis. |
| `previous_proposer` | [string](#string) |  | fee_pool defines the previous proposer at genesis. |
| `outstanding_rewards` | [ValidatorOutstandingRewardsRecord](#cosmos.distribution.v1beta1.ValidatorOutstandingRewardsRecord) | repeated | fee_pool defines the outstanding rewards of all validators at genesis. |
| `validator_accumulated_commissions` | [ValidatorAccumulatedCommissionRecord](#cosmos.distribution.v1beta1.ValidatorAccumulatedCommissionRecord) | repeated | fee_pool defines the accumulated commisions of all validators at genesis. |
| `validator_historical_rewards` | [ValidatorHistoricalRewardsRecord](#cosmos.distribution.v1beta1.ValidatorHistoricalRewardsRecord) | repeated | fee_pool defines the historical rewards of all validators at genesis. |
| `validator_current_rewards` | [ValidatorCurrentRewardsRecord](#cosmos.distribution.v1beta1.ValidatorCurrentRewardsRecord) | repeated | fee_pool defines the current rewards of all validators at genesis. |
| `delegator_starting_infos` | [DelegatorStartingInfoRecord](#cosmos.distribution.v1beta1.DelegatorStartingInfoRecord) | repeated | fee_pool defines the delegator starting infos at genesis. |
| `validator_slash_events` | [ValidatorSlashEventRecord](#cosmos.distribution.v1beta1.ValidatorSlashEventRecord) | repeated | fee_pool defines the validator slash events at genesis. |






<a name="cosmos.distribution.v1beta1.ValidatorAccumulatedCommissionRecord"></a>

### ValidatorAccumulatedCommissionRecord
ValidatorAccumulatedCommissionRecord is used for import / export via genesis
json.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address is the address of the validator. |
| `accumulated` | [ValidatorAccumulatedCommission](#cosmos.distribution.v1beta1.ValidatorAccumulatedCommission) |  | accumulated is the accumulated commission of a validator. |






<a name="cosmos.distribution.v1beta1.ValidatorCurrentRewardsRecord"></a>

### ValidatorCurrentRewardsRecord
ValidatorCurrentRewardsRecord is used for import / export via genesis json.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address is the address of the validator. |
| `rewards` | [ValidatorCurrentRewards](#cosmos.distribution.v1beta1.ValidatorCurrentRewards) |  | rewards defines the current rewards of a validator. |






<a name="cosmos.distribution.v1beta1.ValidatorHistoricalRewardsRecord"></a>

### ValidatorHistoricalRewardsRecord
ValidatorHistoricalRewardsRecord is used for import / export via genesis
json.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address is the address of the validator. |
| `period` | [uint64](#uint64) |  | period defines the period the historical rewards apply to. |
| `rewards` | [ValidatorHistoricalRewards](#cosmos.distribution.v1beta1.ValidatorHistoricalRewards) |  | rewards defines the historical rewards of a validator. |






<a name="cosmos.distribution.v1beta1.ValidatorOutstandingRewardsRecord"></a>

### ValidatorOutstandingRewardsRecord
ValidatorOutstandingRewardsRecord is used for import/export via genesis json.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address is the address of the validator. |
| `outstanding_rewards` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated | outstanding_rewards represents the oustanding rewards of a validator. |






<a name="cosmos.distribution.v1beta1.ValidatorSlashEventRecord"></a>

### ValidatorSlashEventRecord
ValidatorSlashEventRecord is used for import / export via genesis json.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address is the address of the validator. |
| `height` | [uint64](#uint64) |  | height defines the block height at which the slash event occured. |
| `period` | [uint64](#uint64) |  | period is the period of the slash event. |
| `validator_slash_event` | [ValidatorSlashEvent](#cosmos.distribution.v1beta1.ValidatorSlashEvent) |  | validator_slash_event describes the slash event. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/distribution/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/distribution/v1beta1/query.proto



<a name="cosmos.distribution.v1beta1.QueryCommunityPoolRequest"></a>

### QueryCommunityPoolRequest
QueryCommunityPoolRequest is the request type for the Query/CommunityPool RPC
method.






<a name="cosmos.distribution.v1beta1.QueryCommunityPoolResponse"></a>

### QueryCommunityPoolResponse
QueryCommunityPoolResponse is the response type for the Query/CommunityPool
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pool` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated | pool defines community pool's coins. |






<a name="cosmos.distribution.v1beta1.QueryDelegationRewardsRequest"></a>

### QueryDelegationRewardsRequest
QueryDelegationRewardsRequest is the request type for the
Query/DelegationRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address defines the delegator address to query for. |
| `validator_address` | [string](#string) |  | validator_address defines the validator address to query for. |






<a name="cosmos.distribution.v1beta1.QueryDelegationRewardsResponse"></a>

### QueryDelegationRewardsResponse
QueryDelegationRewardsResponse is the response type for the
Query/DelegationRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rewards` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated | rewards defines the rewards accrued by a delegation. |






<a name="cosmos.distribution.v1beta1.QueryDelegationTotalRewardsRequest"></a>

### QueryDelegationTotalRewardsRequest
QueryDelegationTotalRewardsRequest is the request type for the
Query/DelegationTotalRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address defines the delegator address to query for. |






<a name="cosmos.distribution.v1beta1.QueryDelegationTotalRewardsResponse"></a>

### QueryDelegationTotalRewardsResponse
QueryDelegationTotalRewardsResponse is the response type for the
Query/DelegationTotalRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rewards` | [DelegationDelegatorReward](#cosmos.distribution.v1beta1.DelegationDelegatorReward) | repeated | rewards defines all the rewards accrued by a delegator. |
| `total` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated | total defines the sum of all the rewards. |






<a name="cosmos.distribution.v1beta1.QueryDelegatorValidatorsRequest"></a>

### QueryDelegatorValidatorsRequest
QueryDelegatorValidatorsRequest is the request type for the
Query/DelegatorValidators RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address defines the delegator address to query for. |






<a name="cosmos.distribution.v1beta1.QueryDelegatorValidatorsResponse"></a>

### QueryDelegatorValidatorsResponse
QueryDelegatorValidatorsResponse is the response type for the
Query/DelegatorValidators RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validators` | [string](#string) | repeated | validators defines the validators a delegator is delegating for. |






<a name="cosmos.distribution.v1beta1.QueryDelegatorWithdrawAddressRequest"></a>

### QueryDelegatorWithdrawAddressRequest
QueryDelegatorWithdrawAddressRequest is the request type for the
Query/DelegatorWithdrawAddress RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address defines the delegator address to query for. |






<a name="cosmos.distribution.v1beta1.QueryDelegatorWithdrawAddressResponse"></a>

### QueryDelegatorWithdrawAddressResponse
QueryDelegatorWithdrawAddressResponse is the response type for the
Query/DelegatorWithdrawAddress RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `withdraw_address` | [string](#string) |  | withdraw_address defines the delegator address to query for. |






<a name="cosmos.distribution.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="cosmos.distribution.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmos.distribution.v1beta1.Params) |  | params defines the parameters of the module. |






<a name="cosmos.distribution.v1beta1.QueryValidatorCommissionRequest"></a>

### QueryValidatorCommissionRequest
QueryValidatorCommissionRequest is the request type for the
Query/ValidatorCommission RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address defines the validator address to query for. |






<a name="cosmos.distribution.v1beta1.QueryValidatorCommissionResponse"></a>

### QueryValidatorCommissionResponse
QueryValidatorCommissionResponse is the response type for the
Query/ValidatorCommission RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commission` | [ValidatorAccumulatedCommission](#cosmos.distribution.v1beta1.ValidatorAccumulatedCommission) |  | commission defines the commision the validator received. |






<a name="cosmos.distribution.v1beta1.QueryValidatorOutstandingRewardsRequest"></a>

### QueryValidatorOutstandingRewardsRequest
QueryValidatorOutstandingRewardsRequest is the request type for the
Query/ValidatorOutstandingRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address defines the validator address to query for. |






<a name="cosmos.distribution.v1beta1.QueryValidatorOutstandingRewardsResponse"></a>

### QueryValidatorOutstandingRewardsResponse
QueryValidatorOutstandingRewardsResponse is the response type for the
Query/ValidatorOutstandingRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rewards` | [ValidatorOutstandingRewards](#cosmos.distribution.v1beta1.ValidatorOutstandingRewards) |  |  |






<a name="cosmos.distribution.v1beta1.QueryValidatorSlashesRequest"></a>

### QueryValidatorSlashesRequest
QueryValidatorSlashesRequest is the request type for the
Query/ValidatorSlashes RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address defines the validator address to query for. |
| `starting_height` | [uint64](#uint64) |  | starting_height defines the optional starting height to query the slashes. |
| `ending_height` | [uint64](#uint64) |  | starting_height defines the optional ending height to query the slashes. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.distribution.v1beta1.QueryValidatorSlashesResponse"></a>

### QueryValidatorSlashesResponse
QueryValidatorSlashesResponse is the response type for the
Query/ValidatorSlashes RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `slashes` | [ValidatorSlashEvent](#cosmos.distribution.v1beta1.ValidatorSlashEvent) | repeated | slashes defines the slashes the validator received. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.distribution.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service for distribution module.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#cosmos.distribution.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#cosmos.distribution.v1beta1.QueryParamsResponse) | Params queries params of the distribution module. | GET|/cosmos/distribution/v1beta1/params|
| `ValidatorOutstandingRewards` | [QueryValidatorOutstandingRewardsRequest](#cosmos.distribution.v1beta1.QueryValidatorOutstandingRewardsRequest) | [QueryValidatorOutstandingRewardsResponse](#cosmos.distribution.v1beta1.QueryValidatorOutstandingRewardsResponse) | ValidatorOutstandingRewards queries rewards of a validator address. | GET|/cosmos/distribution/v1beta1/validators/{validator_address}/outstanding_rewards|
| `ValidatorCommission` | [QueryValidatorCommissionRequest](#cosmos.distribution.v1beta1.QueryValidatorCommissionRequest) | [QueryValidatorCommissionResponse](#cosmos.distribution.v1beta1.QueryValidatorCommissionResponse) | ValidatorCommission queries accumulated commission for a validator. | GET|/cosmos/distribution/v1beta1/validators/{validator_address}/commission|
| `ValidatorSlashes` | [QueryValidatorSlashesRequest](#cosmos.distribution.v1beta1.QueryValidatorSlashesRequest) | [QueryValidatorSlashesResponse](#cosmos.distribution.v1beta1.QueryValidatorSlashesResponse) | ValidatorSlashes queries slash events of a validator. | GET|/cosmos/distribution/v1beta1/validators/{validator_address}/slashes|
| `DelegationRewards` | [QueryDelegationRewardsRequest](#cosmos.distribution.v1beta1.QueryDelegationRewardsRequest) | [QueryDelegationRewardsResponse](#cosmos.distribution.v1beta1.QueryDelegationRewardsResponse) | DelegationRewards queries the total rewards accrued by a delegation. | GET|/cosmos/distribution/v1beta1/delegators/{delegator_address}/rewards/{validator_address}|
| `DelegationTotalRewards` | [QueryDelegationTotalRewardsRequest](#cosmos.distribution.v1beta1.QueryDelegationTotalRewardsRequest) | [QueryDelegationTotalRewardsResponse](#cosmos.distribution.v1beta1.QueryDelegationTotalRewardsResponse) | DelegationTotalRewards queries the total rewards accrued by a each validator. | GET|/cosmos/distribution/v1beta1/delegators/{delegator_address}/rewards|
| `DelegatorValidators` | [QueryDelegatorValidatorsRequest](#cosmos.distribution.v1beta1.QueryDelegatorValidatorsRequest) | [QueryDelegatorValidatorsResponse](#cosmos.distribution.v1beta1.QueryDelegatorValidatorsResponse) | DelegatorValidators queries the validators of a delegator. | GET|/cosmos/distribution/v1beta1/delegators/{delegator_address}/validators|
| `DelegatorWithdrawAddress` | [QueryDelegatorWithdrawAddressRequest](#cosmos.distribution.v1beta1.QueryDelegatorWithdrawAddressRequest) | [QueryDelegatorWithdrawAddressResponse](#cosmos.distribution.v1beta1.QueryDelegatorWithdrawAddressResponse) | DelegatorWithdrawAddress queries withdraw address of a delegator. | GET|/cosmos/distribution/v1beta1/delegators/{delegator_address}/withdraw_address|
| `CommunityPool` | [QueryCommunityPoolRequest](#cosmos.distribution.v1beta1.QueryCommunityPoolRequest) | [QueryCommunityPoolResponse](#cosmos.distribution.v1beta1.QueryCommunityPoolResponse) | CommunityPool queries the community pool coins. | GET|/cosmos/distribution/v1beta1/community_pool|

 <!-- end services -->



<a name="cosmos/distribution/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/distribution/v1beta1/tx.proto



<a name="cosmos.distribution.v1beta1.MsgFundCommunityPool"></a>

### MsgFundCommunityPool
MsgFundCommunityPool allows an account to directly
fund the community pool.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `depositor` | [string](#string) |  |  |






<a name="cosmos.distribution.v1beta1.MsgFundCommunityPoolResponse"></a>

### MsgFundCommunityPoolResponse
MsgFundCommunityPoolResponse defines the Msg/FundCommunityPool response type.






<a name="cosmos.distribution.v1beta1.MsgSetWithdrawAddress"></a>

### MsgSetWithdrawAddress
MsgSetWithdrawAddress sets the withdraw address for
a delegator (or validator self-delegation).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `withdraw_address` | [string](#string) |  |  |






<a name="cosmos.distribution.v1beta1.MsgSetWithdrawAddressResponse"></a>

### MsgSetWithdrawAddressResponse
MsgSetWithdrawAddressResponse defines the Msg/SetWithdrawAddress response type.






<a name="cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward"></a>

### MsgWithdrawDelegatorReward
MsgWithdrawDelegatorReward represents delegation withdrawal to a delegator
from a single validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |






<a name="cosmos.distribution.v1beta1.MsgWithdrawDelegatorRewardResponse"></a>

### MsgWithdrawDelegatorRewardResponse
MsgWithdrawDelegatorRewardResponse defines the Msg/WithdrawDelegatorReward response type.






<a name="cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission"></a>

### MsgWithdrawValidatorCommission
MsgWithdrawValidatorCommission withdraws the full commission to the validator
address.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  |  |






<a name="cosmos.distribution.v1beta1.MsgWithdrawValidatorCommissionResponse"></a>

### MsgWithdrawValidatorCommissionResponse
MsgWithdrawValidatorCommissionResponse defines the Msg/WithdrawValidatorCommission response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.distribution.v1beta1.Msg"></a>

### Msg
Msg defines the distribution Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `SetWithdrawAddress` | [MsgSetWithdrawAddress](#cosmos.distribution.v1beta1.MsgSetWithdrawAddress) | [MsgSetWithdrawAddressResponse](#cosmos.distribution.v1beta1.MsgSetWithdrawAddressResponse) | SetWithdrawAddress defines a method to change the withdraw address for a delegator (or validator self-delegation). | |
| `WithdrawDelegatorReward` | [MsgWithdrawDelegatorReward](#cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward) | [MsgWithdrawDelegatorRewardResponse](#cosmos.distribution.v1beta1.MsgWithdrawDelegatorRewardResponse) | WithdrawDelegatorReward defines a method to withdraw rewards of delegator from a single validator. | |
| `WithdrawValidatorCommission` | [MsgWithdrawValidatorCommission](#cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission) | [MsgWithdrawValidatorCommissionResponse](#cosmos.distribution.v1beta1.MsgWithdrawValidatorCommissionResponse) | WithdrawValidatorCommission defines a method to withdraw the full commission to the validator address. | |
| `FundCommunityPool` | [MsgFundCommunityPool](#cosmos.distribution.v1beta1.MsgFundCommunityPool) | [MsgFundCommunityPoolResponse](#cosmos.distribution.v1beta1.MsgFundCommunityPoolResponse) | FundCommunityPool defines a method to allow an account to directly fund the community pool. | |

 <!-- end services -->



<a name="cosmos/evidence/v1beta1/evidence.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/evidence/v1beta1/evidence.proto



<a name="cosmos.evidence.v1beta1.Equivocation"></a>

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



<a name="cosmos/evidence/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/evidence/v1beta1/genesis.proto



<a name="cosmos.evidence.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the evidence module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `evidence` | [google.protobuf.Any](#google.protobuf.Any) | repeated | evidence defines all the evidence at genesis. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/evidence/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/evidence/v1beta1/query.proto



<a name="cosmos.evidence.v1beta1.QueryAllEvidenceRequest"></a>

### QueryAllEvidenceRequest
QueryEvidenceRequest is the request type for the Query/AllEvidence RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.evidence.v1beta1.QueryAllEvidenceResponse"></a>

### QueryAllEvidenceResponse
QueryAllEvidenceResponse is the response type for the Query/AllEvidence RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `evidence` | [google.protobuf.Any](#google.protobuf.Any) | repeated | evidence returns all evidences. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.evidence.v1beta1.QueryEvidenceRequest"></a>

### QueryEvidenceRequest
QueryEvidenceRequest is the request type for the Query/Evidence RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `evidence_hash` | [bytes](#bytes) |  | evidence_hash defines the hash of the requested evidence. |






<a name="cosmos.evidence.v1beta1.QueryEvidenceResponse"></a>

### QueryEvidenceResponse
QueryEvidenceResponse is the response type for the Query/Evidence RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `evidence` | [google.protobuf.Any](#google.protobuf.Any) |  | evidence returns the requested evidence. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.evidence.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Evidence` | [QueryEvidenceRequest](#cosmos.evidence.v1beta1.QueryEvidenceRequest) | [QueryEvidenceResponse](#cosmos.evidence.v1beta1.QueryEvidenceResponse) | Evidence queries evidence based on evidence hash. | GET|/cosmos/evidence/v1beta1/evidence/{evidence_hash}|
| `AllEvidence` | [QueryAllEvidenceRequest](#cosmos.evidence.v1beta1.QueryAllEvidenceRequest) | [QueryAllEvidenceResponse](#cosmos.evidence.v1beta1.QueryAllEvidenceResponse) | AllEvidence queries all evidence. | GET|/cosmos/evidence/v1beta1/evidence|

 <!-- end services -->



<a name="cosmos/evidence/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/evidence/v1beta1/tx.proto



<a name="cosmos.evidence.v1beta1.MsgSubmitEvidence"></a>

### MsgSubmitEvidence
MsgSubmitEvidence represents a message that supports submitting arbitrary
Evidence of misbehavior such as equivocation or counterfactual signing.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `submitter` | [string](#string) |  |  |
| `evidence` | [google.protobuf.Any](#google.protobuf.Any) |  |  |






<a name="cosmos.evidence.v1beta1.MsgSubmitEvidenceResponse"></a>

### MsgSubmitEvidenceResponse
MsgSubmitEvidenceResponse defines the Msg/SubmitEvidence response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [bytes](#bytes) |  | hash defines the hash of the evidence. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.evidence.v1beta1.Msg"></a>

### Msg
Msg defines the evidence Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `SubmitEvidence` | [MsgSubmitEvidence](#cosmos.evidence.v1beta1.MsgSubmitEvidence) | [MsgSubmitEvidenceResponse](#cosmos.evidence.v1beta1.MsgSubmitEvidenceResponse) | SubmitEvidence submits an arbitrary Evidence of misbehavior such as equivocation or counterfactual signing. | |

 <!-- end services -->



<a name="cosmos/feegrant/v1beta1/feegrant.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/feegrant/v1beta1/feegrant.proto
Since: cosmos-sdk 0.43


<a name="cosmos.feegrant.v1beta1.AllowedMsgAllowance"></a>

### AllowedMsgAllowance
AllowedMsgAllowance creates allowance only for specified message types.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `allowance` | [google.protobuf.Any](#google.protobuf.Any) |  | allowance can be any of basic and filtered fee allowance. |
| `allowed_messages` | [string](#string) | repeated | allowed_messages are the messages for which the grantee has the access. |






<a name="cosmos.feegrant.v1beta1.BasicAllowance"></a>

### BasicAllowance
BasicAllowance implements Allowance with a one-time grant of tokens
that optionally expires. The grantee can use up to SpendLimit to cover fees.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `spend_limit` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | spend_limit specifies the maximum amount of tokens that can be spent by this allowance and will be updated as tokens are spent. If it is empty, there is no spend limit and any amount of coins can be spent. |
| `expiration` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | expiration specifies an optional time when this allowance expires |






<a name="cosmos.feegrant.v1beta1.Grant"></a>

### Grant
Grant is stored in the KVStore to record a grant with full context


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  | granter is the address of the user granting an allowance of their funds. |
| `grantee` | [string](#string) |  | grantee is the address of the user being granted an allowance of another user's funds. |
| `allowance` | [google.protobuf.Any](#google.protobuf.Any) |  | allowance can be any of basic and filtered fee allowance. |






<a name="cosmos.feegrant.v1beta1.PeriodicAllowance"></a>

### PeriodicAllowance
PeriodicAllowance extends Allowance to allow for both a maximum cap,
as well as a limit per time period.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `basic` | [BasicAllowance](#cosmos.feegrant.v1beta1.BasicAllowance) |  | basic specifies a struct of `BasicAllowance` |
| `period` | [google.protobuf.Duration](#google.protobuf.Duration) |  | period specifies the time duration in which period_spend_limit coins can be spent before that allowance is reset |
| `period_spend_limit` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | period_spend_limit specifies the maximum number of coins that can be spent in the period |
| `period_can_spend` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | period_can_spend is the number of coins left to be spent before the period_reset time |
| `period_reset` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | period_reset is the time at which this period resets and a new one begins, it is calculated from the start time of the first transaction after the last period ended |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/feegrant/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/feegrant/v1beta1/genesis.proto
Since: cosmos-sdk 0.43


<a name="cosmos.feegrant.v1beta1.GenesisState"></a>

### GenesisState
GenesisState contains a set of fee allowances, persisted from the store


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `allowances` | [Grant](#cosmos.feegrant.v1beta1.Grant) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/feegrant/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/feegrant/v1beta1/query.proto
Since: cosmos-sdk 0.43


<a name="cosmos.feegrant.v1beta1.QueryAllowanceRequest"></a>

### QueryAllowanceRequest
QueryAllowanceRequest is the request type for the Query/Allowance RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  | granter is the address of the user granting an allowance of their funds. |
| `grantee` | [string](#string) |  | grantee is the address of the user being granted an allowance of another user's funds. |






<a name="cosmos.feegrant.v1beta1.QueryAllowanceResponse"></a>

### QueryAllowanceResponse
QueryAllowanceResponse is the response type for the Query/Allowance RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `allowance` | [Grant](#cosmos.feegrant.v1beta1.Grant) |  | allowance is a allowance granted for grantee by granter. |






<a name="cosmos.feegrant.v1beta1.QueryAllowancesByGranterRequest"></a>

### QueryAllowancesByGranterRequest
QueryAllowancesByGranterRequest is the request type for the Query/AllowancesByGranter RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an pagination for the request. |






<a name="cosmos.feegrant.v1beta1.QueryAllowancesByGranterResponse"></a>

### QueryAllowancesByGranterResponse
QueryAllowancesByGranterResponse is the response type for the Query/AllowancesByGranter RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `allowances` | [Grant](#cosmos.feegrant.v1beta1.Grant) | repeated | allowances that have been issued by the granter. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines an pagination for the response. |






<a name="cosmos.feegrant.v1beta1.QueryAllowancesRequest"></a>

### QueryAllowancesRequest
QueryAllowancesRequest is the request type for the Query/Allowances RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grantee` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an pagination for the request. |






<a name="cosmos.feegrant.v1beta1.QueryAllowancesResponse"></a>

### QueryAllowancesResponse
QueryAllowancesResponse is the response type for the Query/Allowances RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `allowances` | [Grant](#cosmos.feegrant.v1beta1.Grant) | repeated | allowances are allowance's granted for grantee by granter. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines an pagination for the response. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.feegrant.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Allowance` | [QueryAllowanceRequest](#cosmos.feegrant.v1beta1.QueryAllowanceRequest) | [QueryAllowanceResponse](#cosmos.feegrant.v1beta1.QueryAllowanceResponse) | Allowance returns fee granted to the grantee by the granter. | GET|/cosmos/feegrant/v1beta1/allowance/{granter}/{grantee}|
| `Allowances` | [QueryAllowancesRequest](#cosmos.feegrant.v1beta1.QueryAllowancesRequest) | [QueryAllowancesResponse](#cosmos.feegrant.v1beta1.QueryAllowancesResponse) | Allowances returns all the grants for address. | GET|/cosmos/feegrant/v1beta1/allowances/{grantee}|
| `AllowancesByGranter` | [QueryAllowancesByGranterRequest](#cosmos.feegrant.v1beta1.QueryAllowancesByGranterRequest) | [QueryAllowancesByGranterResponse](#cosmos.feegrant.v1beta1.QueryAllowancesByGranterResponse) | AllowancesByGranter returns all the grants given by an address Since v0.46 | GET|/cosmos/feegrant/v1beta1/issued/{granter}|

 <!-- end services -->



<a name="cosmos/feegrant/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/feegrant/v1beta1/tx.proto
Since: cosmos-sdk 0.43


<a name="cosmos.feegrant.v1beta1.MsgGrantAllowance"></a>

### MsgGrantAllowance
MsgGrantAllowance adds permission for Grantee to spend up to Allowance
of fees from the account of Granter.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  | granter is the address of the user granting an allowance of their funds. |
| `grantee` | [string](#string) |  | grantee is the address of the user being granted an allowance of another user's funds. |
| `allowance` | [google.protobuf.Any](#google.protobuf.Any) |  | allowance can be any of basic and filtered fee allowance. |






<a name="cosmos.feegrant.v1beta1.MsgGrantAllowanceResponse"></a>

### MsgGrantAllowanceResponse
MsgGrantAllowanceResponse defines the Msg/GrantAllowanceResponse response type.






<a name="cosmos.feegrant.v1beta1.MsgRevokeAllowance"></a>

### MsgRevokeAllowance
MsgRevokeAllowance removes any existing Allowance from Granter to Grantee.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  | granter is the address of the user granting an allowance of their funds. |
| `grantee` | [string](#string) |  | grantee is the address of the user being granted an allowance of another user's funds. |






<a name="cosmos.feegrant.v1beta1.MsgRevokeAllowanceResponse"></a>

### MsgRevokeAllowanceResponse
MsgRevokeAllowanceResponse defines the Msg/RevokeAllowanceResponse response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.feegrant.v1beta1.Msg"></a>

### Msg
Msg defines the feegrant msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `GrantAllowance` | [MsgGrantAllowance](#cosmos.feegrant.v1beta1.MsgGrantAllowance) | [MsgGrantAllowanceResponse](#cosmos.feegrant.v1beta1.MsgGrantAllowanceResponse) | GrantAllowance grants fee allowance to the grantee on the granter's account with the provided expiration time. | |
| `RevokeAllowance` | [MsgRevokeAllowance](#cosmos.feegrant.v1beta1.MsgRevokeAllowance) | [MsgRevokeAllowanceResponse](#cosmos.feegrant.v1beta1.MsgRevokeAllowanceResponse) | RevokeAllowance revokes any fee allowance of granter's account that has been granted to the grantee. | |

 <!-- end services -->



<a name="cosmos/genutil/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/genutil/v1beta1/genesis.proto



<a name="cosmos.genutil.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the raw genesis transaction in JSON.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gen_txs` | [bytes](#bytes) | repeated | gen_txs defines the genesis transactions. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/gov/v1beta1/gov.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/gov/v1beta1/gov.proto



<a name="cosmos.gov.v1beta1.Deposit"></a>

### Deposit
Deposit defines an amount deposited by an account address to an active
proposal.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  |  |
| `depositor` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="cosmos.gov.v1beta1.DepositParams"></a>

### DepositParams
DepositParams defines the params for deposits on governance proposals.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `min_deposit` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | Minimum deposit for a proposal to enter voting period. |
| `max_deposit_period` | [google.protobuf.Duration](#google.protobuf.Duration) |  | Maximum period for Atom holders to deposit on a proposal. Initial value: 2 months. |






<a name="cosmos.gov.v1beta1.Proposal"></a>

### Proposal
Proposal defines the core field members of a governance proposal.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  |  |
| `content` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `status` | [ProposalStatus](#cosmos.gov.v1beta1.ProposalStatus) |  |  |
| `final_tally_result` | [TallyResult](#cosmos.gov.v1beta1.TallyResult) |  |  |
| `submit_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| `deposit_end_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| `total_deposit` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `voting_start_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| `voting_end_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |






<a name="cosmos.gov.v1beta1.TallyParams"></a>

### TallyParams
TallyParams defines the params for tallying votes on governance proposals.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `quorum` | [bytes](#bytes) |  | Minimum percentage of total stake needed to vote for a result to be considered valid. |
| `threshold` | [bytes](#bytes) |  | Minimum proportion of Yes votes for proposal to pass. Default value: 0.5. |
| `veto_threshold` | [bytes](#bytes) |  | Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Default value: 1/3. |






<a name="cosmos.gov.v1beta1.TallyResult"></a>

### TallyResult
TallyResult defines a standard tally for a governance proposal.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `yes` | [string](#string) |  |  |
| `abstain` | [string](#string) |  |  |
| `no` | [string](#string) |  |  |
| `no_with_veto` | [string](#string) |  |  |






<a name="cosmos.gov.v1beta1.TextProposal"></a>

### TextProposal
TextProposal defines a standard text proposal whose changes need to be
manually updated in case of approval.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |






<a name="cosmos.gov.v1beta1.Vote"></a>

### Vote
Vote defines a vote on a governance proposal.
A Vote consists of a proposal ID, the voter, and the vote option.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  |  |
| `voter` | [string](#string) |  |  |
| `option` | [VoteOption](#cosmos.gov.v1beta1.VoteOption) |  | **Deprecated.** Deprecated: Prefer to use `options` instead. This field is set in queries if and only if `len(options) == 1` and that option has weight 1. In all other cases, this field will default to VOTE_OPTION_UNSPECIFIED. |
| `options` | [WeightedVoteOption](#cosmos.gov.v1beta1.WeightedVoteOption) | repeated | Since: cosmos-sdk 0.43 |






<a name="cosmos.gov.v1beta1.VotingParams"></a>

### VotingParams
VotingParams defines the params for voting on governance proposals.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `voting_period` | [google.protobuf.Duration](#google.protobuf.Duration) |  | Length of the voting period. |






<a name="cosmos.gov.v1beta1.WeightedVoteOption"></a>

### WeightedVoteOption
WeightedVoteOption defines a unit of vote for vote split.

Since: cosmos-sdk 0.43


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `option` | [VoteOption](#cosmos.gov.v1beta1.VoteOption) |  |  |
| `weight` | [string](#string) |  |  |





 <!-- end messages -->


<a name="cosmos.gov.v1beta1.ProposalStatus"></a>

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



<a name="cosmos.gov.v1beta1.VoteOption"></a>

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



<a name="cosmos/gov/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/gov/v1beta1/genesis.proto



<a name="cosmos.gov.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the gov module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `starting_proposal_id` | [uint64](#uint64) |  | starting_proposal_id is the ID of the starting proposal. |
| `deposits` | [Deposit](#cosmos.gov.v1beta1.Deposit) | repeated | deposits defines all the deposits present at genesis. |
| `votes` | [Vote](#cosmos.gov.v1beta1.Vote) | repeated | votes defines all the votes present at genesis. |
| `proposals` | [Proposal](#cosmos.gov.v1beta1.Proposal) | repeated | proposals defines all the proposals present at genesis. |
| `deposit_params` | [DepositParams](#cosmos.gov.v1beta1.DepositParams) |  | params defines all the paramaters of related to deposit. |
| `voting_params` | [VotingParams](#cosmos.gov.v1beta1.VotingParams) |  | params defines all the paramaters of related to voting. |
| `tally_params` | [TallyParams](#cosmos.gov.v1beta1.TallyParams) |  | params defines all the paramaters of related to tally. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/gov/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/gov/v1beta1/query.proto



<a name="cosmos.gov.v1beta1.QueryDepositRequest"></a>

### QueryDepositRequest
QueryDepositRequest is the request type for the Query/Deposit RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id defines the unique id of the proposal. |
| `depositor` | [string](#string) |  | depositor defines the deposit addresses from the proposals. |






<a name="cosmos.gov.v1beta1.QueryDepositResponse"></a>

### QueryDepositResponse
QueryDepositResponse is the response type for the Query/Deposit RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `deposit` | [Deposit](#cosmos.gov.v1beta1.Deposit) |  | deposit defines the requested deposit. |






<a name="cosmos.gov.v1beta1.QueryDepositsRequest"></a>

### QueryDepositsRequest
QueryDepositsRequest is the request type for the Query/Deposits RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id defines the unique id of the proposal. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.gov.v1beta1.QueryDepositsResponse"></a>

### QueryDepositsResponse
QueryDepositsResponse is the response type for the Query/Deposits RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `deposits` | [Deposit](#cosmos.gov.v1beta1.Deposit) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.gov.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params_type` | [string](#string) |  | params_type defines which parameters to query for, can be one of "voting", "tallying" or "deposit". |






<a name="cosmos.gov.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `voting_params` | [VotingParams](#cosmos.gov.v1beta1.VotingParams) |  | voting_params defines the parameters related to voting. |
| `deposit_params` | [DepositParams](#cosmos.gov.v1beta1.DepositParams) |  | deposit_params defines the parameters related to deposit. |
| `tally_params` | [TallyParams](#cosmos.gov.v1beta1.TallyParams) |  | tally_params defines the parameters related to tally. |






<a name="cosmos.gov.v1beta1.QueryProposalRequest"></a>

### QueryProposalRequest
QueryProposalRequest is the request type for the Query/Proposal RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id defines the unique id of the proposal. |






<a name="cosmos.gov.v1beta1.QueryProposalResponse"></a>

### QueryProposalResponse
QueryProposalResponse is the response type for the Query/Proposal RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal` | [Proposal](#cosmos.gov.v1beta1.Proposal) |  |  |






<a name="cosmos.gov.v1beta1.QueryProposalsRequest"></a>

### QueryProposalsRequest
QueryProposalsRequest is the request type for the Query/Proposals RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_status` | [ProposalStatus](#cosmos.gov.v1beta1.ProposalStatus) |  | proposal_status defines the status of the proposals. |
| `voter` | [string](#string) |  | voter defines the voter address for the proposals. |
| `depositor` | [string](#string) |  | depositor defines the deposit addresses from the proposals. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.gov.v1beta1.QueryProposalsResponse"></a>

### QueryProposalsResponse
QueryProposalsResponse is the response type for the Query/Proposals RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposals` | [Proposal](#cosmos.gov.v1beta1.Proposal) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.gov.v1beta1.QueryTallyResultRequest"></a>

### QueryTallyResultRequest
QueryTallyResultRequest is the request type for the Query/Tally RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id defines the unique id of the proposal. |






<a name="cosmos.gov.v1beta1.QueryTallyResultResponse"></a>

### QueryTallyResultResponse
QueryTallyResultResponse is the response type for the Query/Tally RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tally` | [TallyResult](#cosmos.gov.v1beta1.TallyResult) |  | tally defines the requested tally. |






<a name="cosmos.gov.v1beta1.QueryVoteRequest"></a>

### QueryVoteRequest
QueryVoteRequest is the request type for the Query/Vote RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id defines the unique id of the proposal. |
| `voter` | [string](#string) |  | voter defines the oter address for the proposals. |






<a name="cosmos.gov.v1beta1.QueryVoteResponse"></a>

### QueryVoteResponse
QueryVoteResponse is the response type for the Query/Vote RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vote` | [Vote](#cosmos.gov.v1beta1.Vote) |  | vote defined the queried vote. |






<a name="cosmos.gov.v1beta1.QueryVotesRequest"></a>

### QueryVotesRequest
QueryVotesRequest is the request type for the Query/Votes RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id defines the unique id of the proposal. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.gov.v1beta1.QueryVotesResponse"></a>

### QueryVotesResponse
QueryVotesResponse is the response type for the Query/Votes RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `votes` | [Vote](#cosmos.gov.v1beta1.Vote) | repeated | votes defined the queried votes. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.gov.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service for gov module

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Proposal` | [QueryProposalRequest](#cosmos.gov.v1beta1.QueryProposalRequest) | [QueryProposalResponse](#cosmos.gov.v1beta1.QueryProposalResponse) | Proposal queries proposal details based on ProposalID. | GET|/cosmos/gov/v1beta1/proposals/{proposal_id}|
| `Proposals` | [QueryProposalsRequest](#cosmos.gov.v1beta1.QueryProposalsRequest) | [QueryProposalsResponse](#cosmos.gov.v1beta1.QueryProposalsResponse) | Proposals queries all proposals based on given status. | GET|/cosmos/gov/v1beta1/proposals|
| `Vote` | [QueryVoteRequest](#cosmos.gov.v1beta1.QueryVoteRequest) | [QueryVoteResponse](#cosmos.gov.v1beta1.QueryVoteResponse) | Vote queries voted information based on proposalID, voterAddr. | GET|/cosmos/gov/v1beta1/proposals/{proposal_id}/votes/{voter}|
| `Votes` | [QueryVotesRequest](#cosmos.gov.v1beta1.QueryVotesRequest) | [QueryVotesResponse](#cosmos.gov.v1beta1.QueryVotesResponse) | Votes queries votes of a given proposal. | GET|/cosmos/gov/v1beta1/proposals/{proposal_id}/votes|
| `Params` | [QueryParamsRequest](#cosmos.gov.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#cosmos.gov.v1beta1.QueryParamsResponse) | Params queries all parameters of the gov module. | GET|/cosmos/gov/v1beta1/params/{params_type}|
| `Deposit` | [QueryDepositRequest](#cosmos.gov.v1beta1.QueryDepositRequest) | [QueryDepositResponse](#cosmos.gov.v1beta1.QueryDepositResponse) | Deposit queries single deposit information based proposalID, depositAddr. | GET|/cosmos/gov/v1beta1/proposals/{proposal_id}/deposits/{depositor}|
| `Deposits` | [QueryDepositsRequest](#cosmos.gov.v1beta1.QueryDepositsRequest) | [QueryDepositsResponse](#cosmos.gov.v1beta1.QueryDepositsResponse) | Deposits queries all deposits of a single proposal. | GET|/cosmos/gov/v1beta1/proposals/{proposal_id}/deposits|
| `TallyResult` | [QueryTallyResultRequest](#cosmos.gov.v1beta1.QueryTallyResultRequest) | [QueryTallyResultResponse](#cosmos.gov.v1beta1.QueryTallyResultResponse) | TallyResult queries the tally of a proposal vote. | GET|/cosmos/gov/v1beta1/proposals/{proposal_id}/tally|

 <!-- end services -->



<a name="cosmos/gov/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/gov/v1beta1/tx.proto



<a name="cosmos.gov.v1beta1.MsgDeposit"></a>

### MsgDeposit
MsgDeposit defines a message to submit a deposit to an existing proposal.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  |  |
| `depositor` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="cosmos.gov.v1beta1.MsgDepositResponse"></a>

### MsgDepositResponse
MsgDepositResponse defines the Msg/Deposit response type.






<a name="cosmos.gov.v1beta1.MsgSubmitProposal"></a>

### MsgSubmitProposal
MsgSubmitProposal defines an sdk.Msg type that supports submitting arbitrary
proposal Content.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `content` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `initial_deposit` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `proposer` | [string](#string) |  |  |






<a name="cosmos.gov.v1beta1.MsgSubmitProposalResponse"></a>

### MsgSubmitProposalResponse
MsgSubmitProposalResponse defines the Msg/SubmitProposal response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  |  |






<a name="cosmos.gov.v1beta1.MsgVote"></a>

### MsgVote
MsgVote defines a message to cast a vote.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  |  |
| `voter` | [string](#string) |  |  |
| `option` | [VoteOption](#cosmos.gov.v1beta1.VoteOption) |  |  |






<a name="cosmos.gov.v1beta1.MsgVoteResponse"></a>

### MsgVoteResponse
MsgVoteResponse defines the Msg/Vote response type.






<a name="cosmos.gov.v1beta1.MsgVoteWeighted"></a>

### MsgVoteWeighted
MsgVoteWeighted defines a message to cast a vote.

Since: cosmos-sdk 0.43


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  |  |
| `voter` | [string](#string) |  |  |
| `options` | [WeightedVoteOption](#cosmos.gov.v1beta1.WeightedVoteOption) | repeated |  |






<a name="cosmos.gov.v1beta1.MsgVoteWeightedResponse"></a>

### MsgVoteWeightedResponse
MsgVoteWeightedResponse defines the Msg/VoteWeighted response type.

Since: cosmos-sdk 0.43





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.gov.v1beta1.Msg"></a>

### Msg
Msg defines the bank Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `SubmitProposal` | [MsgSubmitProposal](#cosmos.gov.v1beta1.MsgSubmitProposal) | [MsgSubmitProposalResponse](#cosmos.gov.v1beta1.MsgSubmitProposalResponse) | SubmitProposal defines a method to create new proposal given a content. | |
| `Vote` | [MsgVote](#cosmos.gov.v1beta1.MsgVote) | [MsgVoteResponse](#cosmos.gov.v1beta1.MsgVoteResponse) | Vote defines a method to add a vote on a specific proposal. | |
| `VoteWeighted` | [MsgVoteWeighted](#cosmos.gov.v1beta1.MsgVoteWeighted) | [MsgVoteWeightedResponse](#cosmos.gov.v1beta1.MsgVoteWeightedResponse) | VoteWeighted defines a method to add a weighted vote on a specific proposal.

Since: cosmos-sdk 0.43 | |
| `Deposit` | [MsgDeposit](#cosmos.gov.v1beta1.MsgDeposit) | [MsgDepositResponse](#cosmos.gov.v1beta1.MsgDepositResponse) | Deposit defines a method to add deposit on a specific proposal. | |

 <!-- end services -->



<a name="cosmos/mint/v1beta1/mint.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/mint/v1beta1/mint.proto



<a name="cosmos.mint.v1beta1.Minter"></a>

### Minter
Minter represents the minting state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `inflation` | [string](#string) |  | current annual inflation rate |
| `annual_provisions` | [string](#string) |  | current annual expected provisions |






<a name="cosmos.mint.v1beta1.Params"></a>

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



<a name="cosmos/mint/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/mint/v1beta1/genesis.proto



<a name="cosmos.mint.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the mint module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `minter` | [Minter](#cosmos.mint.v1beta1.Minter) |  | minter is a space for holding current inflation information. |
| `params` | [Params](#cosmos.mint.v1beta1.Params) |  | params defines all the paramaters of the module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/mint/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/mint/v1beta1/query.proto



<a name="cosmos.mint.v1beta1.QueryAnnualProvisionsRequest"></a>

### QueryAnnualProvisionsRequest
QueryAnnualProvisionsRequest is the request type for the
Query/AnnualProvisions RPC method.






<a name="cosmos.mint.v1beta1.QueryAnnualProvisionsResponse"></a>

### QueryAnnualProvisionsResponse
QueryAnnualProvisionsResponse is the response type for the
Query/AnnualProvisions RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `annual_provisions` | [bytes](#bytes) |  | annual_provisions is the current minting annual provisions value. |






<a name="cosmos.mint.v1beta1.QueryInflationRequest"></a>

### QueryInflationRequest
QueryInflationRequest is the request type for the Query/Inflation RPC method.






<a name="cosmos.mint.v1beta1.QueryInflationResponse"></a>

### QueryInflationResponse
QueryInflationResponse is the response type for the Query/Inflation RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `inflation` | [bytes](#bytes) |  | inflation is the current minting inflation value. |






<a name="cosmos.mint.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="cosmos.mint.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmos.mint.v1beta1.Params) |  | params defines the parameters of the module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.mint.v1beta1.Query"></a>

### Query
Query provides defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#cosmos.mint.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#cosmos.mint.v1beta1.QueryParamsResponse) | Params returns the total set of minting parameters. | GET|/cosmos/mint/v1beta1/params|
| `Inflation` | [QueryInflationRequest](#cosmos.mint.v1beta1.QueryInflationRequest) | [QueryInflationResponse](#cosmos.mint.v1beta1.QueryInflationResponse) | Inflation returns the current minting inflation value. | GET|/cosmos/mint/v1beta1/inflation|
| `AnnualProvisions` | [QueryAnnualProvisionsRequest](#cosmos.mint.v1beta1.QueryAnnualProvisionsRequest) | [QueryAnnualProvisionsResponse](#cosmos.mint.v1beta1.QueryAnnualProvisionsResponse) | AnnualProvisions current minting annual provisions value. | GET|/cosmos/mint/v1beta1/annual_provisions|

 <!-- end services -->



<a name="cosmos/params/v1beta1/params.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/params/v1beta1/params.proto



<a name="cosmos.params.v1beta1.ParamChange"></a>

### ParamChange
ParamChange defines an individual parameter change, for use in
ParameterChangeProposal.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subspace` | [string](#string) |  |  |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |






<a name="cosmos.params.v1beta1.ParameterChangeProposal"></a>

### ParameterChangeProposal
ParameterChangeProposal defines a proposal to change one or more parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `changes` | [ParamChange](#cosmos.params.v1beta1.ParamChange) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/params/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/params/v1beta1/query.proto



<a name="cosmos.params.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is request type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subspace` | [string](#string) |  | subspace defines the module to query the parameter for. |
| `key` | [string](#string) |  | key defines the key of the parameter in the subspace. |






<a name="cosmos.params.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `param` | [ParamChange](#cosmos.params.v1beta1.ParamChange) |  | param defines the queried parameter. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.params.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#cosmos.params.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#cosmos.params.v1beta1.QueryParamsResponse) | Params queries a specific parameter of a module, given its subspace and key. | GET|/cosmos/params/v1beta1/params|

 <!-- end services -->



<a name="cosmos/slashing/v1beta1/slashing.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/slashing/v1beta1/slashing.proto



<a name="cosmos.slashing.v1beta1.Params"></a>

### Params
Params represents the parameters used for by the slashing module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signed_blocks_window` | [int64](#int64) |  |  |
| `min_signed_per_window` | [bytes](#bytes) |  |  |
| `downtime_jail_duration` | [google.protobuf.Duration](#google.protobuf.Duration) |  |  |
| `slash_fraction_double_sign` | [bytes](#bytes) |  |  |
| `slash_fraction_downtime` | [bytes](#bytes) |  |  |






<a name="cosmos.slashing.v1beta1.ValidatorSigningInfo"></a>

### ValidatorSigningInfo
ValidatorSigningInfo defines a validator's signing info for monitoring their
liveness activity.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |
| `start_height` | [int64](#int64) |  | Height at which validator was first a candidate OR was unjailed |
| `index_offset` | [int64](#int64) |  | Index which is incremented each time the validator was a bonded in a block and may have signed a precommit or not. This in conjunction with the `SignedBlocksWindow` param determines the index in the `MissedBlocksBitArray`. |
| `jailed_until` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | Timestamp until which the validator is jailed due to liveness downtime. |
| `tombstoned` | [bool](#bool) |  | Whether or not a validator has been tombstoned (killed out of validator set). It is set once the validator commits an equivocation or for any other configured misbehiavor. |
| `missed_blocks_counter` | [int64](#int64) |  | A counter kept to avoid unnecessary array reads. Note that `Sum(MissedBlocksBitArray)` always equals `MissedBlocksCounter`. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/slashing/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/slashing/v1beta1/genesis.proto



<a name="cosmos.slashing.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the slashing module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmos.slashing.v1beta1.Params) |  | params defines all the paramaters of related to deposit. |
| `signing_infos` | [SigningInfo](#cosmos.slashing.v1beta1.SigningInfo) | repeated | signing_infos represents a map between validator addresses and their signing infos. |
| `missed_blocks` | [ValidatorMissedBlocks](#cosmos.slashing.v1beta1.ValidatorMissedBlocks) | repeated | missed_blocks represents a map between validator addresses and their missed blocks. |






<a name="cosmos.slashing.v1beta1.MissedBlock"></a>

### MissedBlock
MissedBlock contains height and missed status as boolean.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `index` | [int64](#int64) |  | index is the height at which the block was missed. |
| `missed` | [bool](#bool) |  | missed is the missed status. |






<a name="cosmos.slashing.v1beta1.SigningInfo"></a>

### SigningInfo
SigningInfo stores validator signing info of corresponding address.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the validator address. |
| `validator_signing_info` | [ValidatorSigningInfo](#cosmos.slashing.v1beta1.ValidatorSigningInfo) |  | validator_signing_info represents the signing info of this validator. |






<a name="cosmos.slashing.v1beta1.ValidatorMissedBlocks"></a>

### ValidatorMissedBlocks
ValidatorMissedBlocks contains array of missed blocks of corresponding
address.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the validator address. |
| `missed_blocks` | [MissedBlock](#cosmos.slashing.v1beta1.MissedBlock) | repeated | missed_blocks is an array of missed blocks by the validator. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/slashing/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/slashing/v1beta1/query.proto



<a name="cosmos.slashing.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method






<a name="cosmos.slashing.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmos.slashing.v1beta1.Params) |  |  |






<a name="cosmos.slashing.v1beta1.QuerySigningInfoRequest"></a>

### QuerySigningInfoRequest
QuerySigningInfoRequest is the request type for the Query/SigningInfo RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cons_address` | [string](#string) |  | cons_address is the address to query signing info of |






<a name="cosmos.slashing.v1beta1.QuerySigningInfoResponse"></a>

### QuerySigningInfoResponse
QuerySigningInfoResponse is the response type for the Query/SigningInfo RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `val_signing_info` | [ValidatorSigningInfo](#cosmos.slashing.v1beta1.ValidatorSigningInfo) |  | val_signing_info is the signing info of requested val cons address |






<a name="cosmos.slashing.v1beta1.QuerySigningInfosRequest"></a>

### QuerySigningInfosRequest
QuerySigningInfosRequest is the request type for the Query/SigningInfos RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="cosmos.slashing.v1beta1.QuerySigningInfosResponse"></a>

### QuerySigningInfosResponse
QuerySigningInfosResponse is the response type for the Query/SigningInfos RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `info` | [ValidatorSigningInfo](#cosmos.slashing.v1beta1.ValidatorSigningInfo) | repeated | info is the signing info of all validators |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.slashing.v1beta1.Query"></a>

### Query
Query provides defines the gRPC querier service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#cosmos.slashing.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#cosmos.slashing.v1beta1.QueryParamsResponse) | Params queries the parameters of slashing module | GET|/cosmos/slashing/v1beta1/params|
| `SigningInfo` | [QuerySigningInfoRequest](#cosmos.slashing.v1beta1.QuerySigningInfoRequest) | [QuerySigningInfoResponse](#cosmos.slashing.v1beta1.QuerySigningInfoResponse) | SigningInfo queries the signing info of given cons address | GET|/cosmos/slashing/v1beta1/signing_infos/{cons_address}|
| `SigningInfos` | [QuerySigningInfosRequest](#cosmos.slashing.v1beta1.QuerySigningInfosRequest) | [QuerySigningInfosResponse](#cosmos.slashing.v1beta1.QuerySigningInfosResponse) | SigningInfos queries signing info of all validators | GET|/cosmos/slashing/v1beta1/signing_infos|

 <!-- end services -->



<a name="cosmos/slashing/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/slashing/v1beta1/tx.proto



<a name="cosmos.slashing.v1beta1.MsgUnjail"></a>

### MsgUnjail
MsgUnjail defines the Msg/Unjail request type


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_addr` | [string](#string) |  |  |






<a name="cosmos.slashing.v1beta1.MsgUnjailResponse"></a>

### MsgUnjailResponse
MsgUnjailResponse defines the Msg/Unjail response type





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.slashing.v1beta1.Msg"></a>

### Msg
Msg defines the slashing Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Unjail` | [MsgUnjail](#cosmos.slashing.v1beta1.MsgUnjail) | [MsgUnjailResponse](#cosmos.slashing.v1beta1.MsgUnjailResponse) | Unjail defines a method for unjailing a jailed validator, thus returning them into the bonded validator set, so they can begin receiving provisions and rewards again. | |

 <!-- end services -->



<a name="cosmos/staking/v1beta1/authz.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/staking/v1beta1/authz.proto



<a name="cosmos.staking.v1beta1.StakeAuthorization"></a>

### StakeAuthorization
StakeAuthorization defines authorization for delegate/undelegate/redelegate.

Since: cosmos-sdk 0.43


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `max_tokens` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | max_tokens specifies the maximum amount of tokens can be delegate to a validator. If it is empty, there is no spend limit and any amount of coins can be delegated. |
| `allow_list` | [StakeAuthorization.Validators](#cosmos.staking.v1beta1.StakeAuthorization.Validators) |  | allow_list specifies list of validator addresses to whom grantee can delegate tokens on behalf of granter's account. |
| `deny_list` | [StakeAuthorization.Validators](#cosmos.staking.v1beta1.StakeAuthorization.Validators) |  | deny_list specifies list of validator addresses to whom grantee can not delegate tokens. |
| `authorization_type` | [AuthorizationType](#cosmos.staking.v1beta1.AuthorizationType) |  | authorization_type defines one of AuthorizationType. |






<a name="cosmos.staking.v1beta1.StakeAuthorization.Validators"></a>

### StakeAuthorization.Validators
Validators defines list of validator addresses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) | repeated |  |





 <!-- end messages -->


<a name="cosmos.staking.v1beta1.AuthorizationType"></a>

### AuthorizationType
AuthorizationType defines the type of staking module authorization type

Since: cosmos-sdk 0.43

| Name | Number | Description |
| ---- | ------ | ----------- |
| AUTHORIZATION_TYPE_UNSPECIFIED | 0 | AUTHORIZATION_TYPE_UNSPECIFIED specifies an unknown authorization type |
| AUTHORIZATION_TYPE_DELEGATE | 1 | AUTHORIZATION_TYPE_DELEGATE defines an authorization type for Msg/Delegate |
| AUTHORIZATION_TYPE_UNDELEGATE | 2 | AUTHORIZATION_TYPE_UNDELEGATE defines an authorization type for Msg/Undelegate |
| AUTHORIZATION_TYPE_REDELEGATE | 3 | AUTHORIZATION_TYPE_REDELEGATE defines an authorization type for Msg/BeginRedelegate |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/staking/v1beta1/staking.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/staking/v1beta1/staking.proto



<a name="cosmos.staking.v1beta1.Commission"></a>

### Commission
Commission defines commission parameters for a given validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commission_rates` | [CommissionRates](#cosmos.staking.v1beta1.CommissionRates) |  | commission_rates defines the initial commission rates to be used for creating a validator. |
| `update_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | update_time is the last time the commission rate was changed. |






<a name="cosmos.staking.v1beta1.CommissionRates"></a>

### CommissionRates
CommissionRates defines the initial commission rates to be used for creating
a validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rate` | [string](#string) |  | rate is the commission rate charged to delegators, as a fraction. |
| `max_rate` | [string](#string) |  | max_rate defines the maximum commission rate which validator can ever charge, as a fraction. |
| `max_change_rate` | [string](#string) |  | max_change_rate defines the maximum daily increase of the validator commission, as a fraction. |






<a name="cosmos.staking.v1beta1.DVPair"></a>

### DVPair
DVPair is struct that just has a delegator-validator pair with no other data.
It is intended to be used as a marshalable pointer. For example, a DVPair can
be used to construct the key to getting an UnbondingDelegation from state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |






<a name="cosmos.staking.v1beta1.DVPairs"></a>

### DVPairs
DVPairs defines an array of DVPair objects.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pairs` | [DVPair](#cosmos.staking.v1beta1.DVPair) | repeated |  |






<a name="cosmos.staking.v1beta1.DVVTriplet"></a>

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






<a name="cosmos.staking.v1beta1.DVVTriplets"></a>

### DVVTriplets
DVVTriplets defines an array of DVVTriplet objects.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `triplets` | [DVVTriplet](#cosmos.staking.v1beta1.DVVTriplet) | repeated |  |






<a name="cosmos.staking.v1beta1.Delegation"></a>

### Delegation
Delegation represents the bond with tokens held by an account. It is
owned by one delegator, and is associated with the voting power of one
validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address is the bech32-encoded address of the delegator. |
| `validator_address` | [string](#string) |  | validator_address is the bech32-encoded address of the validator. |
| `shares` | [string](#string) |  | shares define the delegation shares received. |






<a name="cosmos.staking.v1beta1.DelegationResponse"></a>

### DelegationResponse
DelegationResponse is equivalent to Delegation except that it contains a
balance in addition to shares which is more suitable for client responses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegation` | [Delegation](#cosmos.staking.v1beta1.Delegation) |  |  |
| `balance` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="cosmos.staking.v1beta1.Description"></a>

### Description
Description defines a validator description.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `moniker` | [string](#string) |  | moniker defines a human-readable name for the validator. |
| `identity` | [string](#string) |  | identity defines an optional identity signature (ex. UPort or Keybase). |
| `website` | [string](#string) |  | website defines an optional website link. |
| `security_contact` | [string](#string) |  | security_contact defines an optional email for security contact. |
| `details` | [string](#string) |  | details define other optional details. |






<a name="cosmos.staking.v1beta1.HistoricalInfo"></a>

### HistoricalInfo
HistoricalInfo contains header and validator, voter information for a given block.
It is stored as part of staking module's state, which persists the `n` most
recent HistoricalInfo
(`n` is set by the staking module's `historical_entries` parameter).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `header` | [ostracon.types.Header](#ostracon.types.Header) |  |  |
| `valset` | [Validator](#cosmos.staking.v1beta1.Validator) | repeated |  |






<a name="cosmos.staking.v1beta1.Params"></a>

### Params
Params defines the parameters for the staking module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unbonding_time` | [google.protobuf.Duration](#google.protobuf.Duration) |  | unbonding_time is the time duration of unbonding. |
| `max_validators` | [uint32](#uint32) |  | max_validators is the maximum number of validators. |
| `max_entries` | [uint32](#uint32) |  | max_entries is the max entries for either unbonding delegation or redelegation (per pair/trio). |
| `historical_entries` | [uint32](#uint32) |  | historical_entries is the number of historical entries to persist. |
| `bond_denom` | [string](#string) |  | bond_denom defines the bondable coin denomination. |






<a name="cosmos.staking.v1beta1.Pool"></a>

### Pool
Pool is used for tracking bonded and not-bonded token supply of the bond
denomination.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `not_bonded_tokens` | [string](#string) |  |  |
| `bonded_tokens` | [string](#string) |  |  |






<a name="cosmos.staking.v1beta1.Redelegation"></a>

### Redelegation
Redelegation contains the list of a particular delegator's redelegating bonds
from a particular source validator to a particular destination validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address is the bech32-encoded address of the delegator. |
| `validator_src_address` | [string](#string) |  | validator_src_address is the validator redelegation source operator address. |
| `validator_dst_address` | [string](#string) |  | validator_dst_address is the validator redelegation destination operator address. |
| `entries` | [RedelegationEntry](#cosmos.staking.v1beta1.RedelegationEntry) | repeated | entries are the redelegation entries.

redelegation entries |






<a name="cosmos.staking.v1beta1.RedelegationEntry"></a>

### RedelegationEntry
RedelegationEntry defines a redelegation object with relevant metadata.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `creation_height` | [int64](#int64) |  | creation_height defines the height which the redelegation took place. |
| `completion_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | completion_time defines the unix time for redelegation completion. |
| `initial_balance` | [string](#string) |  | initial_balance defines the initial balance when redelegation started. |
| `shares_dst` | [string](#string) |  | shares_dst is the amount of destination-validator shares created by redelegation. |






<a name="cosmos.staking.v1beta1.RedelegationEntryResponse"></a>

### RedelegationEntryResponse
RedelegationEntryResponse is equivalent to a RedelegationEntry except that it
contains a balance in addition to shares which is more suitable for client
responses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `redelegation_entry` | [RedelegationEntry](#cosmos.staking.v1beta1.RedelegationEntry) |  |  |
| `balance` | [string](#string) |  |  |






<a name="cosmos.staking.v1beta1.RedelegationResponse"></a>

### RedelegationResponse
RedelegationResponse is equivalent to a Redelegation except that its entries
contain a balance in addition to shares which is more suitable for client
responses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `redelegation` | [Redelegation](#cosmos.staking.v1beta1.Redelegation) |  |  |
| `entries` | [RedelegationEntryResponse](#cosmos.staking.v1beta1.RedelegationEntryResponse) | repeated |  |






<a name="cosmos.staking.v1beta1.UnbondingDelegation"></a>

### UnbondingDelegation
UnbondingDelegation stores all of a single delegator's unbonding bonds
for a single validator in an time-ordered list.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address is the bech32-encoded address of the delegator. |
| `validator_address` | [string](#string) |  | validator_address is the bech32-encoded address of the validator. |
| `entries` | [UnbondingDelegationEntry](#cosmos.staking.v1beta1.UnbondingDelegationEntry) | repeated | entries are the unbonding delegation entries.

unbonding delegation entries |






<a name="cosmos.staking.v1beta1.UnbondingDelegationEntry"></a>

### UnbondingDelegationEntry
UnbondingDelegationEntry defines an unbonding object with relevant metadata.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `creation_height` | [int64](#int64) |  | creation_height is the height which the unbonding took place. |
| `completion_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | completion_time is the unix time for unbonding completion. |
| `initial_balance` | [string](#string) |  | initial_balance defines the tokens initially scheduled to receive at completion. |
| `balance` | [string](#string) |  | balance defines the tokens to receive at completion. |






<a name="cosmos.staking.v1beta1.ValAddresses"></a>

### ValAddresses
ValAddresses defines a repeated set of validator addresses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `addresses` | [string](#string) | repeated |  |






<a name="cosmos.staking.v1beta1.Validator"></a>

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
| `status` | [BondStatus](#cosmos.staking.v1beta1.BondStatus) |  | status is the validator status (bonded/unbonding/unbonded). |
| `tokens` | [string](#string) |  | tokens define the delegated tokens (incl. self-delegation). |
| `delegator_shares` | [string](#string) |  | delegator_shares defines total shares issued to a validator's delegators. |
| `description` | [Description](#cosmos.staking.v1beta1.Description) |  | description defines the description terms for the validator. |
| `unbonding_height` | [int64](#int64) |  | unbonding_height defines, if unbonding, the height at which this validator has begun unbonding. |
| `unbonding_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | unbonding_time defines, if unbonding, the min time for the validator to complete unbonding. |
| `commission` | [Commission](#cosmos.staking.v1beta1.Commission) |  | commission defines the commission parameters. |
| `min_self_delegation` | [string](#string) |  | min_self_delegation is the validator's self declared minimum self delegation. |





 <!-- end messages -->


<a name="cosmos.staking.v1beta1.BondStatus"></a>

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



<a name="cosmos/staking/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/staking/v1beta1/genesis.proto



<a name="cosmos.staking.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the staking module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmos.staking.v1beta1.Params) |  | params defines all the paramaters of related to deposit. |
| `last_total_power` | [bytes](#bytes) |  | last_total_power tracks the total amounts of bonded tokens recorded during the previous end block. |
| `last_validator_powers` | [LastValidatorPower](#cosmos.staking.v1beta1.LastValidatorPower) | repeated | last_validator_powers is a special index that provides a historical list of the last-block's bonded validators. |
| `validators` | [Validator](#cosmos.staking.v1beta1.Validator) | repeated | delegations defines the validator set at genesis. |
| `delegations` | [Delegation](#cosmos.staking.v1beta1.Delegation) | repeated | delegations defines the delegations active at genesis. |
| `unbonding_delegations` | [UnbondingDelegation](#cosmos.staking.v1beta1.UnbondingDelegation) | repeated | unbonding_delegations defines the unbonding delegations active at genesis. |
| `redelegations` | [Redelegation](#cosmos.staking.v1beta1.Redelegation) | repeated | redelegations defines the redelegations active at genesis. |
| `exported` | [bool](#bool) |  |  |






<a name="cosmos.staking.v1beta1.LastValidatorPower"></a>

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



<a name="cosmos/staking/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/staking/v1beta1/query.proto



<a name="cosmos.staking.v1beta1.QueryDelegationRequest"></a>

### QueryDelegationRequest
QueryDelegationRequest is request type for the Query/Delegation RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  | delegator_addr defines the delegator address to query for. |
| `validator_addr` | [string](#string) |  | validator_addr defines the validator address to query for. |






<a name="cosmos.staking.v1beta1.QueryDelegationResponse"></a>

### QueryDelegationResponse
QueryDelegationResponse is response type for the Query/Delegation RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegation_response` | [DelegationResponse](#cosmos.staking.v1beta1.DelegationResponse) |  | delegation_responses defines the delegation info of a delegation. |






<a name="cosmos.staking.v1beta1.QueryDelegatorDelegationsRequest"></a>

### QueryDelegatorDelegationsRequest
QueryDelegatorDelegationsRequest is request type for the
Query/DelegatorDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  | delegator_addr defines the delegator address to query for. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.staking.v1beta1.QueryDelegatorDelegationsResponse"></a>

### QueryDelegatorDelegationsResponse
QueryDelegatorDelegationsResponse is response type for the
Query/DelegatorDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegation_responses` | [DelegationResponse](#cosmos.staking.v1beta1.DelegationResponse) | repeated | delegation_responses defines all the delegations' info of a delegator. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.staking.v1beta1.QueryDelegatorUnbondingDelegationsRequest"></a>

### QueryDelegatorUnbondingDelegationsRequest
QueryDelegatorUnbondingDelegationsRequest is request type for the
Query/DelegatorUnbondingDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  | delegator_addr defines the delegator address to query for. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.staking.v1beta1.QueryDelegatorUnbondingDelegationsResponse"></a>

### QueryDelegatorUnbondingDelegationsResponse
QueryUnbondingDelegatorDelegationsResponse is response type for the
Query/UnbondingDelegatorDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unbonding_responses` | [UnbondingDelegation](#cosmos.staking.v1beta1.UnbondingDelegation) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.staking.v1beta1.QueryDelegatorValidatorRequest"></a>

### QueryDelegatorValidatorRequest
QueryDelegatorValidatorRequest is request type for the
Query/DelegatorValidator RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  | delegator_addr defines the delegator address to query for. |
| `validator_addr` | [string](#string) |  | validator_addr defines the validator address to query for. |






<a name="cosmos.staking.v1beta1.QueryDelegatorValidatorResponse"></a>

### QueryDelegatorValidatorResponse
QueryDelegatorValidatorResponse response type for the
Query/DelegatorValidator RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [Validator](#cosmos.staking.v1beta1.Validator) |  | validator defines the the validator info. |






<a name="cosmos.staking.v1beta1.QueryDelegatorValidatorsRequest"></a>

### QueryDelegatorValidatorsRequest
QueryDelegatorValidatorsRequest is request type for the
Query/DelegatorValidators RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  | delegator_addr defines the delegator address to query for. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.staking.v1beta1.QueryDelegatorValidatorsResponse"></a>

### QueryDelegatorValidatorsResponse
QueryDelegatorValidatorsResponse is response type for the
Query/DelegatorValidators RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validators` | [Validator](#cosmos.staking.v1beta1.Validator) | repeated | validators defines the the validators' info of a delegator. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.staking.v1beta1.QueryHistoricalInfoRequest"></a>

### QueryHistoricalInfoRequest
QueryHistoricalInfoRequest is request type for the Query/HistoricalInfo RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  | height defines at which height to query the historical info. |






<a name="cosmos.staking.v1beta1.QueryHistoricalInfoResponse"></a>

### QueryHistoricalInfoResponse
QueryHistoricalInfoResponse is response type for the Query/HistoricalInfo RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hist` | [HistoricalInfo](#cosmos.staking.v1beta1.HistoricalInfo) |  | hist defines the historical info at the given height. |






<a name="cosmos.staking.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is request type for the Query/Params RPC method.






<a name="cosmos.staking.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmos.staking.v1beta1.Params) |  | params holds all the parameters of this module. |






<a name="cosmos.staking.v1beta1.QueryPoolRequest"></a>

### QueryPoolRequest
QueryPoolRequest is request type for the Query/Pool RPC method.






<a name="cosmos.staking.v1beta1.QueryPoolResponse"></a>

### QueryPoolResponse
QueryPoolResponse is response type for the Query/Pool RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pool` | [Pool](#cosmos.staking.v1beta1.Pool) |  | pool defines the pool info. |






<a name="cosmos.staking.v1beta1.QueryRedelegationsRequest"></a>

### QueryRedelegationsRequest
QueryRedelegationsRequest is request type for the Query/Redelegations RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  | delegator_addr defines the delegator address to query for. |
| `src_validator_addr` | [string](#string) |  | src_validator_addr defines the validator address to redelegate from. |
| `dst_validator_addr` | [string](#string) |  | dst_validator_addr defines the validator address to redelegate to. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.staking.v1beta1.QueryRedelegationsResponse"></a>

### QueryRedelegationsResponse
QueryRedelegationsResponse is response type for the Query/Redelegations RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `redelegation_responses` | [RedelegationResponse](#cosmos.staking.v1beta1.RedelegationResponse) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.staking.v1beta1.QueryUnbondingDelegationRequest"></a>

### QueryUnbondingDelegationRequest
QueryUnbondingDelegationRequest is request type for the
Query/UnbondingDelegation RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  | delegator_addr defines the delegator address to query for. |
| `validator_addr` | [string](#string) |  | validator_addr defines the validator address to query for. |






<a name="cosmos.staking.v1beta1.QueryUnbondingDelegationResponse"></a>

### QueryUnbondingDelegationResponse
QueryDelegationResponse is response type for the Query/UnbondingDelegation
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unbond` | [UnbondingDelegation](#cosmos.staking.v1beta1.UnbondingDelegation) |  | unbond defines the unbonding information of a delegation. |






<a name="cosmos.staking.v1beta1.QueryValidatorDelegationsRequest"></a>

### QueryValidatorDelegationsRequest
QueryValidatorDelegationsRequest is request type for the
Query/ValidatorDelegations RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_addr` | [string](#string) |  | validator_addr defines the validator address to query for. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.staking.v1beta1.QueryValidatorDelegationsResponse"></a>

### QueryValidatorDelegationsResponse
QueryValidatorDelegationsResponse is response type for the
Query/ValidatorDelegations RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegation_responses` | [DelegationResponse](#cosmos.staking.v1beta1.DelegationResponse) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.staking.v1beta1.QueryValidatorRequest"></a>

### QueryValidatorRequest
QueryValidatorRequest is response type for the Query/Validator RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_addr` | [string](#string) |  | validator_addr defines the validator address to query for. |






<a name="cosmos.staking.v1beta1.QueryValidatorResponse"></a>

### QueryValidatorResponse
QueryValidatorResponse is response type for the Query/Validator RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [Validator](#cosmos.staking.v1beta1.Validator) |  | validator defines the the validator info. |






<a name="cosmos.staking.v1beta1.QueryValidatorUnbondingDelegationsRequest"></a>

### QueryValidatorUnbondingDelegationsRequest
QueryValidatorUnbondingDelegationsRequest is required type for the
Query/ValidatorUnbondingDelegations RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_addr` | [string](#string) |  | validator_addr defines the validator address to query for. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.staking.v1beta1.QueryValidatorUnbondingDelegationsResponse"></a>

### QueryValidatorUnbondingDelegationsResponse
QueryValidatorUnbondingDelegationsResponse is response type for the
Query/ValidatorUnbondingDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unbonding_responses` | [UnbondingDelegation](#cosmos.staking.v1beta1.UnbondingDelegation) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.staking.v1beta1.QueryValidatorsRequest"></a>

### QueryValidatorsRequest
QueryValidatorsRequest is request type for Query/Validators RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `status` | [string](#string) |  | status enables to query for validators matching a given status. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.staking.v1beta1.QueryValidatorsResponse"></a>

### QueryValidatorsResponse
QueryValidatorsResponse is response type for the Query/Validators RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validators` | [Validator](#cosmos.staking.v1beta1.Validator) | repeated | validators contains all the queried validators. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.staking.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Validators` | [QueryValidatorsRequest](#cosmos.staking.v1beta1.QueryValidatorsRequest) | [QueryValidatorsResponse](#cosmos.staking.v1beta1.QueryValidatorsResponse) | Validators queries all validators that match the given status. | GET|/cosmos/staking/v1beta1/validators|
| `Validator` | [QueryValidatorRequest](#cosmos.staking.v1beta1.QueryValidatorRequest) | [QueryValidatorResponse](#cosmos.staking.v1beta1.QueryValidatorResponse) | Validator queries validator info for given validator address. | GET|/cosmos/staking/v1beta1/validators/{validator_addr}|
| `ValidatorDelegations` | [QueryValidatorDelegationsRequest](#cosmos.staking.v1beta1.QueryValidatorDelegationsRequest) | [QueryValidatorDelegationsResponse](#cosmos.staking.v1beta1.QueryValidatorDelegationsResponse) | ValidatorDelegations queries delegate info for given validator. | GET|/cosmos/staking/v1beta1/validators/{validator_addr}/delegations|
| `ValidatorUnbondingDelegations` | [QueryValidatorUnbondingDelegationsRequest](#cosmos.staking.v1beta1.QueryValidatorUnbondingDelegationsRequest) | [QueryValidatorUnbondingDelegationsResponse](#cosmos.staking.v1beta1.QueryValidatorUnbondingDelegationsResponse) | ValidatorUnbondingDelegations queries unbonding delegations of a validator. | GET|/cosmos/staking/v1beta1/validators/{validator_addr}/unbonding_delegations|
| `Delegation` | [QueryDelegationRequest](#cosmos.staking.v1beta1.QueryDelegationRequest) | [QueryDelegationResponse](#cosmos.staking.v1beta1.QueryDelegationResponse) | Delegation queries delegate info for given validator delegator pair. | GET|/cosmos/staking/v1beta1/validators/{validator_addr}/delegations/{delegator_addr}|
| `UnbondingDelegation` | [QueryUnbondingDelegationRequest](#cosmos.staking.v1beta1.QueryUnbondingDelegationRequest) | [QueryUnbondingDelegationResponse](#cosmos.staking.v1beta1.QueryUnbondingDelegationResponse) | UnbondingDelegation queries unbonding info for given validator delegator pair. | GET|/cosmos/staking/v1beta1/validators/{validator_addr}/delegations/{delegator_addr}/unbonding_delegation|
| `DelegatorDelegations` | [QueryDelegatorDelegationsRequest](#cosmos.staking.v1beta1.QueryDelegatorDelegationsRequest) | [QueryDelegatorDelegationsResponse](#cosmos.staking.v1beta1.QueryDelegatorDelegationsResponse) | DelegatorDelegations queries all delegations of a given delegator address. | GET|/cosmos/staking/v1beta1/delegations/{delegator_addr}|
| `DelegatorUnbondingDelegations` | [QueryDelegatorUnbondingDelegationsRequest](#cosmos.staking.v1beta1.QueryDelegatorUnbondingDelegationsRequest) | [QueryDelegatorUnbondingDelegationsResponse](#cosmos.staking.v1beta1.QueryDelegatorUnbondingDelegationsResponse) | DelegatorUnbondingDelegations queries all unbonding delegations of a given delegator address. | GET|/cosmos/staking/v1beta1/delegators/{delegator_addr}/unbonding_delegations|
| `Redelegations` | [QueryRedelegationsRequest](#cosmos.staking.v1beta1.QueryRedelegationsRequest) | [QueryRedelegationsResponse](#cosmos.staking.v1beta1.QueryRedelegationsResponse) | Redelegations queries redelegations of given address. | GET|/cosmos/staking/v1beta1/delegators/{delegator_addr}/redelegations|
| `DelegatorValidators` | [QueryDelegatorValidatorsRequest](#cosmos.staking.v1beta1.QueryDelegatorValidatorsRequest) | [QueryDelegatorValidatorsResponse](#cosmos.staking.v1beta1.QueryDelegatorValidatorsResponse) | DelegatorValidators queries all validators info for given delegator address. | GET|/cosmos/staking/v1beta1/delegators/{delegator_addr}/validators|
| `DelegatorValidator` | [QueryDelegatorValidatorRequest](#cosmos.staking.v1beta1.QueryDelegatorValidatorRequest) | [QueryDelegatorValidatorResponse](#cosmos.staking.v1beta1.QueryDelegatorValidatorResponse) | DelegatorValidator queries validator info for given delegator validator pair. | GET|/cosmos/staking/v1beta1/delegators/{delegator_addr}/validators/{validator_addr}|
| `HistoricalInfo` | [QueryHistoricalInfoRequest](#cosmos.staking.v1beta1.QueryHistoricalInfoRequest) | [QueryHistoricalInfoResponse](#cosmos.staking.v1beta1.QueryHistoricalInfoResponse) | HistoricalInfo queries the historical info for given height. | GET|/cosmos/staking/v1beta1/historical_info/{height}|
| `Pool` | [QueryPoolRequest](#cosmos.staking.v1beta1.QueryPoolRequest) | [QueryPoolResponse](#cosmos.staking.v1beta1.QueryPoolResponse) | Pool queries the pool info. | GET|/cosmos/staking/v1beta1/pool|
| `Params` | [QueryParamsRequest](#cosmos.staking.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#cosmos.staking.v1beta1.QueryParamsResponse) | Parameters queries the staking parameters. | GET|/cosmos/staking/v1beta1/params|

 <!-- end services -->



<a name="cosmos/staking/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/staking/v1beta1/tx.proto



<a name="cosmos.staking.v1beta1.MsgBeginRedelegate"></a>

### MsgBeginRedelegate
MsgBeginRedelegate defines a SDK message for performing a redelegation
of coins from a delegator and source validator to a destination validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_src_address` | [string](#string) |  |  |
| `validator_dst_address` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="cosmos.staking.v1beta1.MsgBeginRedelegateResponse"></a>

### MsgBeginRedelegateResponse
MsgBeginRedelegateResponse defines the Msg/BeginRedelegate response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `completion_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |






<a name="cosmos.staking.v1beta1.MsgCreateValidator"></a>

### MsgCreateValidator
MsgCreateValidator defines a SDK message for creating a new validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `description` | [Description](#cosmos.staking.v1beta1.Description) |  |  |
| `commission` | [CommissionRates](#cosmos.staking.v1beta1.CommissionRates) |  |  |
| `min_self_delegation` | [string](#string) |  |  |
| `delegator_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |
| `pubkey` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `value` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="cosmos.staking.v1beta1.MsgCreateValidatorResponse"></a>

### MsgCreateValidatorResponse
MsgCreateValidatorResponse defines the Msg/CreateValidator response type.






<a name="cosmos.staking.v1beta1.MsgDelegate"></a>

### MsgDelegate
MsgDelegate defines a SDK message for performing a delegation of coins
from a delegator to a validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="cosmos.staking.v1beta1.MsgDelegateResponse"></a>

### MsgDelegateResponse
MsgDelegateResponse defines the Msg/Delegate response type.






<a name="cosmos.staking.v1beta1.MsgEditValidator"></a>

### MsgEditValidator
MsgEditValidator defines a SDK message for editing an existing validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `description` | [Description](#cosmos.staking.v1beta1.Description) |  |  |
| `validator_address` | [string](#string) |  |  |
| `commission_rate` | [string](#string) |  | We pass a reference to the new commission rate and min self delegation as it's not mandatory to update. If not updated, the deserialized rate will be zero with no way to distinguish if an update was intended. REF: #2373 |
| `min_self_delegation` | [string](#string) |  |  |






<a name="cosmos.staking.v1beta1.MsgEditValidatorResponse"></a>

### MsgEditValidatorResponse
MsgEditValidatorResponse defines the Msg/EditValidator response type.






<a name="cosmos.staking.v1beta1.MsgUndelegate"></a>

### MsgUndelegate
MsgUndelegate defines a SDK message for performing an undelegation from a
delegate and a validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="cosmos.staking.v1beta1.MsgUndelegateResponse"></a>

### MsgUndelegateResponse
MsgUndelegateResponse defines the Msg/Undelegate response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `completion_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.staking.v1beta1.Msg"></a>

### Msg
Msg defines the staking Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `CreateValidator` | [MsgCreateValidator](#cosmos.staking.v1beta1.MsgCreateValidator) | [MsgCreateValidatorResponse](#cosmos.staking.v1beta1.MsgCreateValidatorResponse) | CreateValidator defines a method for creating a new validator. | |
| `EditValidator` | [MsgEditValidator](#cosmos.staking.v1beta1.MsgEditValidator) | [MsgEditValidatorResponse](#cosmos.staking.v1beta1.MsgEditValidatorResponse) | EditValidator defines a method for editing an existing validator. | |
| `Delegate` | [MsgDelegate](#cosmos.staking.v1beta1.MsgDelegate) | [MsgDelegateResponse](#cosmos.staking.v1beta1.MsgDelegateResponse) | Delegate defines a method for performing a delegation of coins from a delegator to a validator. | |
| `BeginRedelegate` | [MsgBeginRedelegate](#cosmos.staking.v1beta1.MsgBeginRedelegate) | [MsgBeginRedelegateResponse](#cosmos.staking.v1beta1.MsgBeginRedelegateResponse) | BeginRedelegate defines a method for performing a redelegation of coins from a delegator and source validator to a destination validator. | |
| `Undelegate` | [MsgUndelegate](#cosmos.staking.v1beta1.MsgUndelegate) | [MsgUndelegateResponse](#cosmos.staking.v1beta1.MsgUndelegateResponse) | Undelegate defines a method for performing an undelegation from a delegate and a validator. | |

 <!-- end services -->



<a name="cosmos/tx/signing/v1beta1/signing.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/tx/signing/v1beta1/signing.proto



<a name="cosmos.tx.signing.v1beta1.SignatureDescriptor"></a>

### SignatureDescriptor
SignatureDescriptor is a convenience type which represents the full data for
a signature including the public key of the signer, signing modes and the
signature itself. It is primarily used for coordinating signatures between
clients.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `public_key` | [google.protobuf.Any](#google.protobuf.Any) |  | public_key is the public key of the signer |
| `data` | [SignatureDescriptor.Data](#cosmos.tx.signing.v1beta1.SignatureDescriptor.Data) |  |  |
| `sequence` | [uint64](#uint64) |  | sequence is the sequence of the account, which describes the number of committed transactions signed by a given address. It is used to prevent replay attacks. |






<a name="cosmos.tx.signing.v1beta1.SignatureDescriptor.Data"></a>

### SignatureDescriptor.Data
Data represents signature data


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `single` | [SignatureDescriptor.Data.Single](#cosmos.tx.signing.v1beta1.SignatureDescriptor.Data.Single) |  | single represents a single signer |
| `multi` | [SignatureDescriptor.Data.Multi](#cosmos.tx.signing.v1beta1.SignatureDescriptor.Data.Multi) |  | multi represents a multisig signer |






<a name="cosmos.tx.signing.v1beta1.SignatureDescriptor.Data.Multi"></a>

### SignatureDescriptor.Data.Multi
Multi is the signature data for a multisig public key


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bitarray` | [cosmos.crypto.multisig.v1beta1.CompactBitArray](#cosmos.crypto.multisig.v1beta1.CompactBitArray) |  | bitarray specifies which keys within the multisig are signing |
| `signatures` | [SignatureDescriptor.Data](#cosmos.tx.signing.v1beta1.SignatureDescriptor.Data) | repeated | signatures is the signatures of the multi-signature |






<a name="cosmos.tx.signing.v1beta1.SignatureDescriptor.Data.Single"></a>

### SignatureDescriptor.Data.Single
Single is the signature data for a single signer


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `mode` | [SignMode](#cosmos.tx.signing.v1beta1.SignMode) |  | mode is the signing mode of the single signer |
| `signature` | [bytes](#bytes) |  | signature is the raw signature bytes |






<a name="cosmos.tx.signing.v1beta1.SignatureDescriptors"></a>

### SignatureDescriptors
SignatureDescriptors wraps multiple SignatureDescriptor's.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signatures` | [SignatureDescriptor](#cosmos.tx.signing.v1beta1.SignatureDescriptor) | repeated | signatures are the signature descriptors |





 <!-- end messages -->


<a name="cosmos.tx.signing.v1beta1.SignMode"></a>

### SignMode
SignMode represents a signing mode with its own security guarantees.

| Name | Number | Description |
| ---- | ------ | ----------- |
| SIGN_MODE_UNSPECIFIED | 0 | SIGN_MODE_UNSPECIFIED specifies an unknown signing mode and will be rejected |
| SIGN_MODE_DIRECT | 1 | SIGN_MODE_DIRECT specifies a signing mode which uses SignDoc and is verified with raw bytes from Tx |
| SIGN_MODE_TEXTUAL | 2 | SIGN_MODE_TEXTUAL is a future signing mode that will verify some human-readable textual representation on top of the binary representation from SIGN_MODE_DIRECT |
| SIGN_MODE_LEGACY_AMINO_JSON | 127 | SIGN_MODE_LEGACY_AMINO_JSON is a backwards compatibility mode which uses Amino JSON and will be removed in the future |
| SIGN_MODE_EIP_191 | 191 | SIGN_MODE_EIP_191 specifies the sign mode for EIP 191 signing on the Cosmos SDK. Ref: https://eips.ethereum.org/EIPS/eip-191

Currently, SIGN_MODE_EIP_191 is registered as a SignMode enum variant, but is not implemented on the SDK by default. To enable EIP-191, you need to pass a custom `TxConfig` that has an implementation of `SignModeHandler` for EIP-191. The SDK may decide to fully support EIP-191 in the future.

Since: cosmos-sdk 0.45.2 |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/tx/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/tx/v1beta1/tx.proto



<a name="cosmos.tx.v1beta1.AuthInfo"></a>

### AuthInfo
AuthInfo describes the fee and signer modes that are used to sign a
transaction.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signer_infos` | [SignerInfo](#cosmos.tx.v1beta1.SignerInfo) | repeated | signer_infos defines the signing modes for the required signers. The number and order of elements must match the required signers from TxBody's messages. The first element is the primary signer and the one which pays the fee. |
| `fee` | [Fee](#cosmos.tx.v1beta1.Fee) |  | Fee is the fee and gas limit for the transaction. The first signer is the primary signer and the one which pays the fee. The fee can be calculated based on the cost of evaluating the body and doing signature verification of the signers. This can be estimated via simulation. |






<a name="cosmos.tx.v1beta1.Fee"></a>

### Fee
Fee includes the amount of coins paid in fees and the maximum
gas to be used by the transaction. The ratio yields an effective "gasprice",
which must be above some miminum to be accepted into the mempool.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | amount is the amount of coins to be paid as a fee |
| `gas_limit` | [uint64](#uint64) |  | gas_limit is the maximum gas that can be used in transaction processing before an out of gas error occurs |
| `payer` | [string](#string) |  | if unset, the first signer is responsible for paying the fees. If set, the specified account must pay the fees. the payer must be a tx signer (and thus have signed this field in AuthInfo). setting this field does *not* change the ordering of required signers for the transaction. |
| `granter` | [string](#string) |  | if set, the fee payer (either the first signer or the value of the payer field) requests that a fee grant be used to pay fees instead of the fee payer's own balance. If an appropriate fee grant does not exist or the chain does not support fee grants, this will fail |






<a name="cosmos.tx.v1beta1.ModeInfo"></a>

### ModeInfo
ModeInfo describes the signing mode of a single or nested multisig signer.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `single` | [ModeInfo.Single](#cosmos.tx.v1beta1.ModeInfo.Single) |  | single represents a single signer |
| `multi` | [ModeInfo.Multi](#cosmos.tx.v1beta1.ModeInfo.Multi) |  | multi represents a nested multisig signer |






<a name="cosmos.tx.v1beta1.ModeInfo.Multi"></a>

### ModeInfo.Multi
Multi is the mode info for a multisig public key


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bitarray` | [cosmos.crypto.multisig.v1beta1.CompactBitArray](#cosmos.crypto.multisig.v1beta1.CompactBitArray) |  | bitarray specifies which keys within the multisig are signing |
| `mode_infos` | [ModeInfo](#cosmos.tx.v1beta1.ModeInfo) | repeated | mode_infos is the corresponding modes of the signers of the multisig which could include nested multisig public keys |






<a name="cosmos.tx.v1beta1.ModeInfo.Single"></a>

### ModeInfo.Single
Single is the mode info for a single signer. It is structured as a message
to allow for additional fields such as locale for SIGN_MODE_TEXTUAL in the
future


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `mode` | [cosmos.tx.signing.v1beta1.SignMode](#cosmos.tx.signing.v1beta1.SignMode) |  | mode is the signing mode of the single signer |






<a name="cosmos.tx.v1beta1.SignDoc"></a>

### SignDoc
SignDoc is the type used for generating sign bytes for SIGN_MODE_DIRECT.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `body_bytes` | [bytes](#bytes) |  | body_bytes is protobuf serialization of a TxBody that matches the representation in TxRaw. |
| `auth_info_bytes` | [bytes](#bytes) |  | auth_info_bytes is a protobuf serialization of an AuthInfo that matches the representation in TxRaw. |
| `chain_id` | [string](#string) |  | chain_id is the unique identifier of the chain this transaction targets. It prevents signed transactions from being used on another chain by an attacker |
| `account_number` | [uint64](#uint64) |  | account_number is the account number of the account in state |






<a name="cosmos.tx.v1beta1.SignerInfo"></a>

### SignerInfo
SignerInfo describes the public key and signing mode of a single top-level
signer.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `public_key` | [google.protobuf.Any](#google.protobuf.Any) |  | public_key is the public key of the signer. It is optional for accounts that already exist in state. If unset, the verifier can use the required \ signer address for this position and lookup the public key. |
| `mode_info` | [ModeInfo](#cosmos.tx.v1beta1.ModeInfo) |  | mode_info describes the signing mode of the signer and is a nested structure to support nested multisig pubkey's |
| `sequence` | [uint64](#uint64) |  | sequence is the sequence of the account, which describes the number of committed transactions signed by a given address. It is used to prevent replay attacks. |






<a name="cosmos.tx.v1beta1.Tx"></a>

### Tx
Tx is the standard type used for broadcasting transactions.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `body` | [TxBody](#cosmos.tx.v1beta1.TxBody) |  | body is the processable content of the transaction |
| `auth_info` | [AuthInfo](#cosmos.tx.v1beta1.AuthInfo) |  | auth_info is the authorization related content of the transaction, specifically signers, signer modes and fee |
| `signatures` | [bytes](#bytes) | repeated | signatures is a list of signatures that matches the length and order of AuthInfo's signer_infos to allow connecting signature meta information like public key and signing mode by position. |






<a name="cosmos.tx.v1beta1.TxBody"></a>

### TxBody
TxBody is the body of a transaction that all signers sign over.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `messages` | [google.protobuf.Any](#google.protobuf.Any) | repeated | messages is a list of messages to be executed. The required signers of those messages define the number and order of elements in AuthInfo's signer_infos and Tx's signatures. Each required signer address is added to the list only the first time it occurs. By convention, the first required signer (usually from the first message) is referred to as the primary signer and pays the fee for the whole transaction. |
| `memo` | [string](#string) |  | memo is any arbitrary note/comment to be added to the transaction. WARNING: in clients, any publicly exposed text should not be called memo, but should be called `note` instead (see https://github.com/cosmos/cosmos-sdk/issues/9122). |
| `timeout_height` | [uint64](#uint64) |  | timeout is the block height after which this transaction will not be processed by the chain |
| `extension_options` | [google.protobuf.Any](#google.protobuf.Any) | repeated | extension_options are arbitrary options that can be added by chains when the default options are not sufficient. If any of these are present and can't be handled, the transaction will be rejected |
| `non_critical_extension_options` | [google.protobuf.Any](#google.protobuf.Any) | repeated | extension_options are arbitrary options that can be added by chains when the default options are not sufficient. If any of these are present and can't be handled, they will be ignored |






<a name="cosmos.tx.v1beta1.TxRaw"></a>

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



<a name="cosmos/tx/v1beta1/service.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/tx/v1beta1/service.proto



<a name="cosmos.tx.v1beta1.BroadcastTxRequest"></a>

### BroadcastTxRequest
BroadcastTxRequest is the request type for the Service.BroadcastTxRequest
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tx_bytes` | [bytes](#bytes) |  | tx_bytes is the raw transaction. |
| `mode` | [BroadcastMode](#cosmos.tx.v1beta1.BroadcastMode) |  |  |






<a name="cosmos.tx.v1beta1.BroadcastTxResponse"></a>

### BroadcastTxResponse
BroadcastTxResponse is the response type for the
Service.BroadcastTx method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tx_response` | [cosmos.base.abci.v1beta1.TxResponse](#cosmos.base.abci.v1beta1.TxResponse) |  | tx_response is the queried TxResponses. |






<a name="cosmos.tx.v1beta1.GetBlockWithTxsRequest"></a>

### GetBlockWithTxsRequest
GetBlockWithTxsRequest is the request type for the Service.GetBlockWithTxs
RPC method.

Since: cosmos-sdk 0.45.2


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  | height is the height of the block to query. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines a pagination for the request. |






<a name="cosmos.tx.v1beta1.GetBlockWithTxsResponse"></a>

### GetBlockWithTxsResponse
GetBlockWithTxsResponse is the response type for the Service.GetBlockWithTxs method.

Since: cosmos-sdk 0.45.2


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `txs` | [Tx](#cosmos.tx.v1beta1.Tx) | repeated | txs are the transactions in the block. |
| `block_id` | [ostracon.types.BlockID](#ostracon.types.BlockID) |  |  |
| `block` | [ostracon.types.Block](#ostracon.types.Block) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines a pagination for the response. |






<a name="cosmos.tx.v1beta1.GetTxRequest"></a>

### GetTxRequest
GetTxRequest is the request type for the Service.GetTx
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [string](#string) |  | hash is the tx hash to query, encoded as a hex string. |






<a name="cosmos.tx.v1beta1.GetTxResponse"></a>

### GetTxResponse
GetTxResponse is the response type for the Service.GetTx method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tx` | [Tx](#cosmos.tx.v1beta1.Tx) |  | tx is the queried transaction. |
| `tx_response` | [cosmos.base.abci.v1beta1.TxResponse](#cosmos.base.abci.v1beta1.TxResponse) |  | tx_response is the queried TxResponses. |






<a name="cosmos.tx.v1beta1.GetTxsEventRequest"></a>

### GetTxsEventRequest
GetTxsEventRequest is the request type for the Service.TxsByEvents
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `events` | [string](#string) | repeated | events is the list of transaction event type. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines a pagination for the request. |
| `order_by` | [OrderBy](#cosmos.tx.v1beta1.OrderBy) |  |  |






<a name="cosmos.tx.v1beta1.GetTxsEventResponse"></a>

### GetTxsEventResponse
GetTxsEventResponse is the response type for the Service.TxsByEvents
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `txs` | [Tx](#cosmos.tx.v1beta1.Tx) | repeated | txs is the list of queried transactions. |
| `tx_responses` | [cosmos.base.abci.v1beta1.TxResponse](#cosmos.base.abci.v1beta1.TxResponse) | repeated | tx_responses is the list of queried TxResponses. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines a pagination for the response. |






<a name="cosmos.tx.v1beta1.SimulateRequest"></a>

### SimulateRequest
SimulateRequest is the request type for the Service.Simulate
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tx` | [Tx](#cosmos.tx.v1beta1.Tx) |  | **Deprecated.** tx is the transaction to simulate. Deprecated. Send raw tx bytes instead. |
| `tx_bytes` | [bytes](#bytes) |  | tx_bytes is the raw transaction.

Since: cosmos-sdk 0.43 |






<a name="cosmos.tx.v1beta1.SimulateResponse"></a>

### SimulateResponse
SimulateResponse is the response type for the
Service.SimulateRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gas_info` | [cosmos.base.abci.v1beta1.GasInfo](#cosmos.base.abci.v1beta1.GasInfo) |  | gas_info is the information about gas used in the simulation. |
| `result` | [cosmos.base.abci.v1beta1.Result](#cosmos.base.abci.v1beta1.Result) |  | result is the result of the simulation. |





 <!-- end messages -->


<a name="cosmos.tx.v1beta1.BroadcastMode"></a>

### BroadcastMode
BroadcastMode specifies the broadcast mode for the TxService.Broadcast RPC method.

| Name | Number | Description |
| ---- | ------ | ----------- |
| BROADCAST_MODE_UNSPECIFIED | 0 | zero-value for mode ordering |
| BROADCAST_MODE_BLOCK | 1 | BROADCAST_MODE_BLOCK defines a tx broadcasting mode where the client waits for the tx to be committed in a block. |
| BROADCAST_MODE_SYNC | 2 | BROADCAST_MODE_SYNC defines a tx broadcasting mode where the client waits for a CheckTx execution response only. |
| BROADCAST_MODE_ASYNC | 3 | BROADCAST_MODE_ASYNC defines a tx broadcasting mode where the client returns immediately. |



<a name="cosmos.tx.v1beta1.OrderBy"></a>

### OrderBy
OrderBy defines the sorting order

| Name | Number | Description |
| ---- | ------ | ----------- |
| ORDER_BY_UNSPECIFIED | 0 | ORDER_BY_UNSPECIFIED specifies an unknown sorting order. OrderBy defaults to ASC in this case. |
| ORDER_BY_ASC | 1 | ORDER_BY_ASC defines ascending order |
| ORDER_BY_DESC | 2 | ORDER_BY_DESC defines descending order |


 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.tx.v1beta1.Service"></a>

### Service
Service defines a gRPC service for interacting with transactions.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Simulate` | [SimulateRequest](#cosmos.tx.v1beta1.SimulateRequest) | [SimulateResponse](#cosmos.tx.v1beta1.SimulateResponse) | Simulate simulates executing a transaction for estimating gas usage. | POST|/cosmos/tx/v1beta1/simulate|
| `GetTx` | [GetTxRequest](#cosmos.tx.v1beta1.GetTxRequest) | [GetTxResponse](#cosmos.tx.v1beta1.GetTxResponse) | GetTx fetches a tx by hash. | GET|/cosmos/tx/v1beta1/txs/{hash}|
| `BroadcastTx` | [BroadcastTxRequest](#cosmos.tx.v1beta1.BroadcastTxRequest) | [BroadcastTxResponse](#cosmos.tx.v1beta1.BroadcastTxResponse) | BroadcastTx broadcast transaction. | POST|/cosmos/tx/v1beta1/txs|
| `GetTxsEvent` | [GetTxsEventRequest](#cosmos.tx.v1beta1.GetTxsEventRequest) | [GetTxsEventResponse](#cosmos.tx.v1beta1.GetTxsEventResponse) | GetTxsEvent fetches txs by event. | GET|/cosmos/tx/v1beta1/txs|
| `GetBlockWithTxs` | [GetBlockWithTxsRequest](#cosmos.tx.v1beta1.GetBlockWithTxsRequest) | [GetBlockWithTxsResponse](#cosmos.tx.v1beta1.GetBlockWithTxsResponse) | GetBlockWithTxs fetches a block with decoded txs.

Since: cosmos-sdk 0.45.2 | GET|/cosmos/tx/v1beta1/txs/block/{height}|

 <!-- end services -->



<a name="cosmos/upgrade/v1beta1/upgrade.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/upgrade/v1beta1/upgrade.proto



<a name="cosmos.upgrade.v1beta1.CancelSoftwareUpgradeProposal"></a>

### CancelSoftwareUpgradeProposal
CancelSoftwareUpgradeProposal is a gov Content type for cancelling a software
upgrade.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |






<a name="cosmos.upgrade.v1beta1.ModuleVersion"></a>

### ModuleVersion
ModuleVersion specifies a module and its consensus version.

Since: cosmos-sdk 0.43


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | name of the app module |
| `version` | [uint64](#uint64) |  | consensus version of the app module |






<a name="cosmos.upgrade.v1beta1.Plan"></a>

### Plan
Plan specifies information about a planned upgrade and when it should occur.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | Sets the name for the upgrade. This name will be used by the upgraded version of the software to apply any special "on-upgrade" commands during the first BeginBlock method after the upgrade is applied. It is also used to detect whether a software version can handle a given upgrade. If no upgrade handler with this name has been set in the software, it will be assumed that the software is out-of-date when the upgrade Time or Height is reached and the software will exit. |
| `time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | **Deprecated.** Deprecated: Time based upgrades have been deprecated. Time based upgrade logic has been removed from the SDK. If this field is not empty, an error will be thrown. |
| `height` | [int64](#int64) |  | The height at which the upgrade must be performed. Only used if Time is not set. |
| `info` | [string](#string) |  | Any application specific upgrade info to be included on-chain such as a git commit that validators could automatically upgrade to |
| `upgraded_client_state` | [google.protobuf.Any](#google.protobuf.Any) |  | **Deprecated.** Deprecated: UpgradedClientState field has been deprecated. IBC upgrade logic has been moved to the IBC module in the sub module 02-client. If this field is not empty, an error will be thrown. |






<a name="cosmos.upgrade.v1beta1.SoftwareUpgradeProposal"></a>

### SoftwareUpgradeProposal
SoftwareUpgradeProposal is a gov Content type for initiating a software
upgrade.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `plan` | [Plan](#cosmos.upgrade.v1beta1.Plan) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/upgrade/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/upgrade/v1beta1/query.proto



<a name="cosmos.upgrade.v1beta1.QueryAppliedPlanRequest"></a>

### QueryAppliedPlanRequest
QueryCurrentPlanRequest is the request type for the Query/AppliedPlan RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | name is the name of the applied plan to query for. |






<a name="cosmos.upgrade.v1beta1.QueryAppliedPlanResponse"></a>

### QueryAppliedPlanResponse
QueryAppliedPlanResponse is the response type for the Query/AppliedPlan RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  | height is the block height at which the plan was applied. |






<a name="cosmos.upgrade.v1beta1.QueryCurrentPlanRequest"></a>

### QueryCurrentPlanRequest
QueryCurrentPlanRequest is the request type for the Query/CurrentPlan RPC
method.






<a name="cosmos.upgrade.v1beta1.QueryCurrentPlanResponse"></a>

### QueryCurrentPlanResponse
QueryCurrentPlanResponse is the response type for the Query/CurrentPlan RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `plan` | [Plan](#cosmos.upgrade.v1beta1.Plan) |  | plan is the current upgrade plan. |






<a name="cosmos.upgrade.v1beta1.QueryModuleVersionsRequest"></a>

### QueryModuleVersionsRequest
QueryModuleVersionsRequest is the request type for the Query/ModuleVersions
RPC method.

Since: cosmos-sdk 0.43


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `module_name` | [string](#string) |  | module_name is a field to query a specific module consensus version from state. Leaving this empty will fetch the full list of module versions from state |






<a name="cosmos.upgrade.v1beta1.QueryModuleVersionsResponse"></a>

### QueryModuleVersionsResponse
QueryModuleVersionsResponse is the response type for the Query/ModuleVersions
RPC method.

Since: cosmos-sdk 0.43


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `module_versions` | [ModuleVersion](#cosmos.upgrade.v1beta1.ModuleVersion) | repeated | module_versions is a list of module names with their consensus versions. |






<a name="cosmos.upgrade.v1beta1.QueryUpgradedConsensusStateRequest"></a>

### QueryUpgradedConsensusStateRequest
QueryUpgradedConsensusStateRequest is the request type for the Query/UpgradedConsensusState
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `last_height` | [int64](#int64) |  | last height of the current chain must be sent in request as this is the height under which next consensus state is stored |






<a name="cosmos.upgrade.v1beta1.QueryUpgradedConsensusStateResponse"></a>

### QueryUpgradedConsensusStateResponse
QueryUpgradedConsensusStateResponse is the response type for the Query/UpgradedConsensusState
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `upgraded_consensus_state` | [bytes](#bytes) |  | Since: cosmos-sdk 0.43 |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.upgrade.v1beta1.Query"></a>

### Query
Query defines the gRPC upgrade querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `CurrentPlan` | [QueryCurrentPlanRequest](#cosmos.upgrade.v1beta1.QueryCurrentPlanRequest) | [QueryCurrentPlanResponse](#cosmos.upgrade.v1beta1.QueryCurrentPlanResponse) | CurrentPlan queries the current upgrade plan. | GET|/cosmos/upgrade/v1beta1/current_plan|
| `AppliedPlan` | [QueryAppliedPlanRequest](#cosmos.upgrade.v1beta1.QueryAppliedPlanRequest) | [QueryAppliedPlanResponse](#cosmos.upgrade.v1beta1.QueryAppliedPlanResponse) | AppliedPlan queries a previously applied upgrade plan by its name. | GET|/cosmos/upgrade/v1beta1/applied_plan/{name}|
| `UpgradedConsensusState` | [QueryUpgradedConsensusStateRequest](#cosmos.upgrade.v1beta1.QueryUpgradedConsensusStateRequest) | [QueryUpgradedConsensusStateResponse](#cosmos.upgrade.v1beta1.QueryUpgradedConsensusStateResponse) | UpgradedConsensusState queries the consensus state that will serve as a trusted kernel for the next version of this chain. It will only be stored at the last height of this chain. UpgradedConsensusState RPC not supported with legacy querier This rpc is deprecated now that IBC has its own replacement (https://github.com/cosmos/ibc-go/blob/2c880a22e9f9cc75f62b527ca94aa75ce1106001/proto/ibc/core/client/v1/query.proto#L54) | GET|/cosmos/upgrade/v1beta1/upgraded_consensus_state/{last_height}|
| `ModuleVersions` | [QueryModuleVersionsRequest](#cosmos.upgrade.v1beta1.QueryModuleVersionsRequest) | [QueryModuleVersionsResponse](#cosmos.upgrade.v1beta1.QueryModuleVersionsResponse) | ModuleVersions queries the list of module versions from state.

Since: cosmos-sdk 0.43 | GET|/cosmos/upgrade/v1beta1/module_versions|

 <!-- end services -->



<a name="cosmos/vesting/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/vesting/v1beta1/tx.proto



<a name="cosmos.vesting.v1beta1.MsgCreateVestingAccount"></a>

### MsgCreateVestingAccount
MsgCreateVestingAccount defines a message that enables creating a vesting
account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `from_address` | [string](#string) |  |  |
| `to_address` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `end_time` | [int64](#int64) |  |  |
| `delayed` | [bool](#bool) |  |  |






<a name="cosmos.vesting.v1beta1.MsgCreateVestingAccountResponse"></a>

### MsgCreateVestingAccountResponse
MsgCreateVestingAccountResponse defines the Msg/CreateVestingAccount response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.vesting.v1beta1.Msg"></a>

### Msg
Msg defines the bank Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `CreateVestingAccount` | [MsgCreateVestingAccount](#cosmos.vesting.v1beta1.MsgCreateVestingAccount) | [MsgCreateVestingAccountResponse](#cosmos.vesting.v1beta1.MsgCreateVestingAccountResponse) | CreateVestingAccount defines a method that enables creating a vesting account. | |

 <!-- end services -->



<a name="cosmos/vesting/v1beta1/vesting.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/vesting/v1beta1/vesting.proto



<a name="cosmos.vesting.v1beta1.BaseVestingAccount"></a>

### BaseVestingAccount
BaseVestingAccount implements the VestingAccount interface. It contains all
the necessary fields needed for any vesting account implementation.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_account` | [cosmos.auth.v1beta1.BaseAccount](#cosmos.auth.v1beta1.BaseAccount) |  |  |
| `original_vesting` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `delegated_free` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `delegated_vesting` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `end_time` | [int64](#int64) |  |  |






<a name="cosmos.vesting.v1beta1.ContinuousVestingAccount"></a>

### ContinuousVestingAccount
ContinuousVestingAccount implements the VestingAccount interface. It
continuously vests by unlocking coins linearly with respect to time.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_vesting_account` | [BaseVestingAccount](#cosmos.vesting.v1beta1.BaseVestingAccount) |  |  |
| `start_time` | [int64](#int64) |  |  |






<a name="cosmos.vesting.v1beta1.DelayedVestingAccount"></a>

### DelayedVestingAccount
DelayedVestingAccount implements the VestingAccount interface. It vests all
coins after a specific time, but non prior. In other words, it keeps them
locked until a specified time.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_vesting_account` | [BaseVestingAccount](#cosmos.vesting.v1beta1.BaseVestingAccount) |  |  |






<a name="cosmos.vesting.v1beta1.Period"></a>

### Period
Period defines a length of time and amount of coins that will vest.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `length` | [int64](#int64) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="cosmos.vesting.v1beta1.PeriodicVestingAccount"></a>

### PeriodicVestingAccount
PeriodicVestingAccount implements the VestingAccount interface. It
periodically vests by unlocking coins during each specified period.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_vesting_account` | [BaseVestingAccount](#cosmos.vesting.v1beta1.BaseVestingAccount) |  |  |
| `start_time` | [int64](#int64) |  |  |
| `vesting_periods` | [Period](#cosmos.vesting.v1beta1.Period) | repeated |  |






<a name="cosmos.vesting.v1beta1.PermanentLockedAccount"></a>

### PermanentLockedAccount
PermanentLockedAccount implements the VestingAccount interface. It does
not ever release coins, locking them indefinitely. Coins in this account can
still be used for delegating and for governance votes even while locked.

Since: cosmos-sdk 0.43


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_vesting_account` | [BaseVestingAccount](#cosmos.vesting.v1beta1.BaseVestingAccount) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

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
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an pagination for the request. |






<a name="lbm.base.ostracon.v1.GetLatestValidatorSetResponse"></a>

### GetLatestValidatorSetResponse
GetLatestValidatorSetResponse is the response type for the Query/GetValidatorSetByHeight RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_height` | [int64](#int64) |  |  |
| `validators` | [Validator](#lbm.base.ostracon.v1.Validator) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines an pagination for the response. |






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
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an pagination for the request. |






<a name="lbm.base.ostracon.v1.GetValidatorSetByHeightResponse"></a>

### GetValidatorSetByHeightResponse
GetValidatorSetByHeightResponse is the response type for the Query/GetValidatorSetByHeight RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_height` | [int64](#int64) |  |  |
| `validators` | [Validator](#lbm.base.ostracon.v1.Validator) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines an pagination for the response. |






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
| `lbm_sdk_version` | [string](#string) |  | Since: cosmos-sdk 0.43 |





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
| `GetBlockByHash` | [GetBlockByHashRequest](#lbm.base.ostracon.v1.GetBlockByHashRequest) | [GetBlockByHashResponse](#lbm.base.ostracon.v1.GetBlockByHashResponse) | GetBlockByHash queries block for given hash. | GET|/lbm/base/ostracon/v1/block/{hash}|
| `GetBlockResultsByHeight` | [GetBlockResultsByHeightRequest](#lbm.base.ostracon.v1.GetBlockResultsByHeightRequest) | [GetBlockResultsByHeightResponse](#lbm.base.ostracon.v1.GetBlockResultsByHeightResponse) | GetBlockResultsByHeight queries block results for given height. | GET|/lbm/base/ostracon/v1/blockresults/{height}|
| `GetLatestValidatorSet` | [GetLatestValidatorSetRequest](#lbm.base.ostracon.v1.GetLatestValidatorSetRequest) | [GetLatestValidatorSetResponse](#lbm.base.ostracon.v1.GetLatestValidatorSetResponse) | GetLatestValidatorSet queries latest validator-set. | GET|/lbm/base/ostracon/v1/validatorsets/latest|
| `GetValidatorSetByHeight` | [GetValidatorSetByHeightRequest](#lbm.base.ostracon.v1.GetValidatorSetByHeightRequest) | [GetValidatorSetByHeightResponse](#lbm.base.ostracon.v1.GetValidatorSetByHeightResponse) | GetValidatorSetByHeight queries validator-set at a given height. | GET|/lbm/base/ostracon/v1/validatorsets/{height}|

 <!-- end services -->



<a name="lbm/collection/v1/collection.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/collection/v1/collection.proto



<a name="lbm.collection.v1.Attribute"></a>

### Attribute
Attribute defines a key and value of the attribute.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |






<a name="lbm.collection.v1.Authorization"></a>

### Authorization
Authorization defines an authorization given to the operator on tokens of the holder.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `holder` | [string](#string) |  | address of the holder which authorizes the manipulation of its tokens. |
| `operator` | [string](#string) |  | address of the operator which the authorization is granted to. |






<a name="lbm.collection.v1.Change"></a>

### Change
Deprecated: use Attribute

Change defines a field-value pair.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `field` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |






<a name="lbm.collection.v1.Coin"></a>

### Coin
Coin defines a token with a token id and an amount.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_id` | [string](#string) |  | token id associated with the token. |
| `amount` | [string](#string) |  | amount of the token. |






<a name="lbm.collection.v1.Contract"></a>

### Contract
Contract defines the information of the contract for the collection.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  | contract_id defines the unique identifier of the contract. |
| `name` | [string](#string) |  | name defines the human-readable name of the contract. |
| `meta` | [string](#string) |  | meta is a brief description of the contract. |
| `base_img_uri` | [string](#string) |  | base img uri is an uri for the contract image stored off chain. |






<a name="lbm.collection.v1.FT"></a>

### FT
Deprecated: use FTClass

FT defines the information of fungible token.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `token_id` | [string](#string) |  | token id defines the unique identifier of the fungible token. |
| `name` | [string](#string) |  | name defines the human-readable name of the fungible token. |
| `meta` | [string](#string) |  | meta is a brief description of the fungible token. |
| `decimals` | [int32](#int32) |  | decimals is the number of decimals which one must divide the amount by to get its user representation. |
| `mintable` | [bool](#bool) |  | mintable represents whether the fungible token is allowed to be minted or burnt. |






<a name="lbm.collection.v1.FTClass"></a>

### FTClass
FTClass defines the class of fungible token.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  | id defines the unique identifier of the token class. Note: size of the class id is 8 in length. Note: token id of the fungible token would be `id` + `00000000`. |
| `name` | [string](#string) |  | name defines the human-readable name of the token class. |
| `meta` | [string](#string) |  | meta is a brief description of the token class. |
| `decimals` | [int32](#int32) |  | decimals is the number of decimals which one must divide the amount by to get its user representation. |
| `mintable` | [bool](#bool) |  | mintable represents whether the token class is allowed to mint or burn its tokens. |






<a name="lbm.collection.v1.Grant"></a>

### Grant
Grant defines permission given to a grantee.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grantee` | [string](#string) |  | address of the grantee. |
| `permission` | [Permission](#lbm.collection.v1.Permission) |  | permission on the contract. |






<a name="lbm.collection.v1.NFT"></a>

### NFT
NFT defines the information of non-fungible token.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_id` | [string](#string) |  | token id defines the unique identifier of the token. |
| `name` | [string](#string) |  | name defines the human-readable name of the token. |
| `meta` | [string](#string) |  | meta is a brief description of the token. |






<a name="lbm.collection.v1.NFTClass"></a>

### NFTClass
NFTClass defines the class of non-fungible token.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  | id defines the unique identifier of the token class. Note: size of the class id is 8 in length. |
| `name` | [string](#string) |  | name defines the human-readable name of the token class. |
| `meta` | [string](#string) |  | meta is a brief description of the token class. |






<a name="lbm.collection.v1.OwnerNFT"></a>

### OwnerNFT
Deprecated: use NFT

OwnerNFT defines the information of non-fungible token.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `token_id` | [string](#string) |  | id defines the unique identifier of the token. |
| `name` | [string](#string) |  | name defines the human-readable name of the token. |
| `meta` | [string](#string) |  | meta is a brief description of the token. |
| `owner` | [string](#string) |  | owner of the token. |






<a name="lbm.collection.v1.Params"></a>

### Params
Params defines the parameters for the collection module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `depth_limit` | [uint32](#uint32) |  |  |
| `width_limit` | [uint32](#uint32) |  |  |






<a name="lbm.collection.v1.TokenType"></a>

### TokenType
Deprecated: use TokenClass

TokenType defines the information of token type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `token_type` | [string](#string) |  | token type defines the unique identifier of the token type. |
| `name` | [string](#string) |  | name defines the human-readable name of the token type. |
| `meta` | [string](#string) |  | meta is a brief description of the token type. |





 <!-- end messages -->


<a name="lbm.collection.v1.LegacyPermission"></a>

### LegacyPermission
Deprecated: use Permission

LegacyPermission enumerates the valid permissions on a contract.

| Name | Number | Description |
| ---- | ------ | ----------- |
| LEGACY_PERMISSION_UNSPECIFIED | 0 | unspecified defines the default permission which is invalid. |
| LEGACY_PERMISSION_ISSUE | 1 | issue defines a permission to create a token class. |
| LEGACY_PERMISSION_MODIFY | 2 | modify defines a permission to modify a contract. |
| LEGACY_PERMISSION_MINT | 3 | mint defines a permission to mint tokens of a contract. |
| LEGACY_PERMISSION_BURN | 4 | burn defines a permission to burn tokens of a contract. |



<a name="lbm.collection.v1.Permission"></a>

### Permission
Permission enumerates the valid permissions on a contract.

| Name | Number | Description |
| ---- | ------ | ----------- |
| PERMISSION_UNSPECIFIED | 0 | unspecified defines the default permission which is invalid. |
| PERMISSION_ISSUE | 1 | PERMISSION_ISSUE defines a permission to create a token class. |
| PERMISSION_MODIFY | 2 | PERMISSION_MODIFY defines a permission to modify a contract. |
| PERMISSION_MINT | 3 | PERMISSION_MINT defines a permission to mint tokens of a contract. |
| PERMISSION_BURN | 4 | PERMISSION_BURN defines a permission to burn tokens of a contract. |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/collection/v1/event.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/collection/v1/event.proto



<a name="lbm.collection.v1.EventAttached"></a>

### EventAttached
EventAttached is emitted when a token is attached to another.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `operator` | [string](#string) |  | address which triggered the attach. |
| `holder` | [string](#string) |  | address which holds the tokens. |
| `subject` | [string](#string) |  | subject of the attach. |
| `target` | [string](#string) |  | target of the attach. |






<a name="lbm.collection.v1.EventAuthorizedOperator"></a>

### EventAuthorizedOperator
EventAuthorizedOperator is emitted when a holder authorizes an operator to manipulate its tokens.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `holder` | [string](#string) |  | address of a holder which authorized the `operator` address as an operator. |
| `operator` | [string](#string) |  | address which became an operator of `holder`. |






<a name="lbm.collection.v1.EventBurned"></a>

### EventBurned
EventBurned is emitted when tokens are burnt.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `operator` | [string](#string) |  | address which triggered the burn. |
| `from` | [string](#string) |  | holder whose tokens were burned. |
| `amount` | [Coin](#lbm.collection.v1.Coin) | repeated | amount of tokens burned. |






<a name="lbm.collection.v1.EventCreatedContract"></a>

### EventCreatedContract
EventCreatedContract is emitted when a new contract is created.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `creator` | [string](#string) |  | address which created the contract. |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `name` | [string](#string) |  | name of the contract. |
| `meta` | [string](#string) |  | metadata of the contract. |
| `base_img_uri` | [string](#string) |  | uri for the contract image stored off chain. |






<a name="lbm.collection.v1.EventCreatedFTClass"></a>

### EventCreatedFTClass
EventCreatedFTClass is emitted when a new fungible token class is created.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `operator` | [string](#string) |  | address which triggered the create. |
| `token_id` | [string](#string) |  | token id associated with the token class. |
| `name` | [string](#string) |  | name of the token class. |
| `meta` | [string](#string) |  | metadata of the token class. |
| `decimals` | [int32](#int32) |  | decimals of the token class. |
| `mintable` | [bool](#bool) |  | mintable represents whether the token class is allowed to mint or burn its tokens. |






<a name="lbm.collection.v1.EventCreatedNFTClass"></a>

### EventCreatedNFTClass
EventCreatedNFTClass is emitted when a new non-fungible token class is created.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `operator` | [string](#string) |  | address which triggered the create. |
| `token_type` | [string](#string) |  | token type associated with the token class. |
| `name` | [string](#string) |  | name of the token class. |
| `meta` | [string](#string) |  | metadata of the token class. |






<a name="lbm.collection.v1.EventDetached"></a>

### EventDetached
EventDetached is emitted when a token is detached from its parent.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `operator` | [string](#string) |  | address which triggered the detach. |
| `holder` | [string](#string) |  | address which holds the token. |
| `subject` | [string](#string) |  | token being detached. |
| `previous_parent` | [string](#string) |  | parent token before the detach. |






<a name="lbm.collection.v1.EventGranted"></a>

### EventGranted
EventGranted is emitted when a granter grants its permission to a grantee.

Info: `granter` would be empty if the permission is granted by an issuance.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `granter` | [string](#string) |  | address of the granter which grants the permission. |
| `grantee` | [string](#string) |  | address of the grantee. |
| `permission` | [Permission](#lbm.collection.v1.Permission) |  | permission on the contract. |






<a name="lbm.collection.v1.EventMintedFT"></a>

### EventMintedFT
EventMintedFT is emitted when fungible tokens are minted.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `operator` | [string](#string) |  | address which triggered the mint. |
| `to` | [string](#string) |  | recipient of the tokens. |
| `amount` | [Coin](#lbm.collection.v1.Coin) | repeated | amount of tokens minted. |






<a name="lbm.collection.v1.EventMintedNFT"></a>

### EventMintedNFT
EventMintedNFT is emitted when non-fungible tokens are minted.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `operator` | [string](#string) |  | address which triggered the mint. |
| `to` | [string](#string) |  | recipient of the tokens. |
| `tokens` | [NFT](#lbm.collection.v1.NFT) | repeated | tokens minted. |






<a name="lbm.collection.v1.EventModifiedContract"></a>

### EventModifiedContract
EventModifiedContract is emitted when the information of a contract is modified.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `operator` | [string](#string) |  | address which triggered the modify. |
| `changes` | [Attribute](#lbm.collection.v1.Attribute) | repeated | changes of the attributes applied. |






<a name="lbm.collection.v1.EventModifiedNFT"></a>

### EventModifiedNFT
EventModifiedNFT is emitted when the information of a non-fungible token is modified.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `operator` | [string](#string) |  | address which triggered the modify. |
| `token_id` | [string](#string) |  | token id associated with the non-fungible token. |
| `changes` | [Attribute](#lbm.collection.v1.Attribute) | repeated | changes of the attributes applied. |






<a name="lbm.collection.v1.EventModifiedTokenClass"></a>

### EventModifiedTokenClass
EventModifiedTokenClass is emitted when the information of a token class is modified.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `operator` | [string](#string) |  | address which triggered the modify. |
| `token_type` | [string](#string) |  | token type associated with the token class. |
| `changes` | [Attribute](#lbm.collection.v1.Attribute) | repeated | changes of the attributes applied. |
| `type_name` | [string](#string) |  | type name of the token class. |






<a name="lbm.collection.v1.EventOwnerChanged"></a>

### EventOwnerChanged
EventOwnerChanged is emitted when the owner of token is changed by operation applied to its ancestor.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `token_id` | [string](#string) |  | token id associated with the token. |
| `from` | [string](#string) |  | address of the previous owner before the change. |
| `to` | [string](#string) |  | address of the new owner. |






<a name="lbm.collection.v1.EventRenounced"></a>

### EventRenounced
EventRenounced is emitted when a grantee renounced its permission.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `grantee` | [string](#string) |  | address of the grantee which abandons its grant. |
| `permission` | [Permission](#lbm.collection.v1.Permission) |  | permission on the contract. |






<a name="lbm.collection.v1.EventRevokedOperator"></a>

### EventRevokedOperator
EventRevokedOperator is emitted when an authorization is revoked.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `holder` | [string](#string) |  | address of a holder which revoked the `operator` address as an operator. |
| `operator` | [string](#string) |  | address which was revoked as an operator of `holder`. |






<a name="lbm.collection.v1.EventRootChanged"></a>

### EventRootChanged
EventRootChanged is emitted when the root of token is changed by operation applied to its ancestor.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `token_id` | [string](#string) |  | token id associated with the token. |
| `from` | [string](#string) |  | token id of the previous root before the change. |
| `to` | [string](#string) |  | token id of the new root. |






<a name="lbm.collection.v1.EventSent"></a>

### EventSent
EventSent is emitted when tokens are transferred.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `operator` | [string](#string) |  | address which triggered the send. |
| `from` | [string](#string) |  | holder whose tokens were sent. |
| `to` | [string](#string) |  | recipient of the tokens. |
| `amount` | [Coin](#lbm.collection.v1.Coin) | repeated | amount of tokens sent. |





 <!-- end messages -->


<a name="lbm.collection.v1.AttributeKey"></a>

### AttributeKey
Deprecated: use typed events.

AttributeKey enumerates the valid attribute keys on x/collection.

| Name | Number | Description |
| ---- | ------ | ----------- |
| ATTRIBUTE_KEY_UNSPECIFIED | 0 |  |
| ATTRIBUTE_KEY_NAME | 1 |  |
| ATTRIBUTE_KEY_META | 2 |  |
| ATTRIBUTE_KEY_CONTRACT_ID | 3 |  |
| ATTRIBUTE_KEY_TOKEN_ID | 4 |  |
| ATTRIBUTE_KEY_OWNER | 5 |  |
| ATTRIBUTE_KEY_AMOUNT | 6 |  |
| ATTRIBUTE_KEY_DECIMALS | 7 |  |
| ATTRIBUTE_KEY_BASE_IMG_URI | 8 |  |
| ATTRIBUTE_KEY_MINTABLE | 9 |  |
| ATTRIBUTE_KEY_TOKEN_TYPE | 10 |  |
| ATTRIBUTE_KEY_FROM | 11 |  |
| ATTRIBUTE_KEY_TO | 12 |  |
| ATTRIBUTE_KEY_PERM | 13 |  |
| ATTRIBUTE_KEY_TO_TOKEN_ID | 14 |  |
| ATTRIBUTE_KEY_FROM_TOKEN_ID | 15 |  |
| ATTRIBUTE_KEY_APPROVER | 16 |  |
| ATTRIBUTE_KEY_PROXY | 17 |  |
| ATTRIBUTE_KEY_OLD_ROOT_TOKEN_ID | 18 |  |
| ATTRIBUTE_KEY_NEW_ROOT_TOKEN_ID | 19 |  |



<a name="lbm.collection.v1.EventType"></a>

### EventType
Deprecated: use typed events.

EventType enumerates the valid event types on x/collection.

| Name | Number | Description |
| ---- | ------ | ----------- |
| EVENT_TYPE_UNSPECIFIED | 0 |  |
| EVENT_TYPE_CREATE_COLLECTION | 1 |  |
| EVENT_TYPE_ISSUE_FT | 2 |  |
| EVENT_TYPE_ISSUE_NFT | 3 |  |
| EVENT_TYPE_MINT_FT | 4 |  |
| EVENT_TYPE_BURN_FT | 5 |  |
| EVENT_TYPE_MINT_NFT | 6 |  |
| EVENT_TYPE_BURN_NFT | 7 |  |
| EVENT_TYPE_BURN_FT_FROM | 8 |  |
| EVENT_TYPE_BURN_NFT_FROM | 9 |  |
| EVENT_TYPE_MODIFY_COLLECTION | 10 |  |
| EVENT_TYPE_MODIFY_TOKEN_TYPE | 11 |  |
| EVENT_TYPE_MODIFY_TOKEN | 12 |  |
| EVENT_TYPE_TRANSFER | 13 |  |
| EVENT_TYPE_TRANSFER_FT | 14 |  |
| EVENT_TYPE_TRANSFER_NFT | 15 |  |
| EVENT_TYPE_TRANSFER_FT_FROM | 16 |  |
| EVENT_TYPE_TRANSFER_NFT_FROM | 17 |  |
| EVENT_TYPE_GRANT_PERM | 18 |  |
| EVENT_TYPE_REVOKE_PERM | 19 |  |
| EVENT_TYPE_ATTACH | 20 |  |
| EVENT_TYPE_DETACH | 21 |  |
| EVENT_TYPE_ATTACH_FROM | 22 |  |
| EVENT_TYPE_DETACH_FROM | 23 |  |
| EVENT_TYPE_APPROVE_COLLECTION | 24 |  |
| EVENT_TYPE_DISAPPROVE_COLLECTION | 25 |  |
| EVENT_TYPE_OPERATION_TRANSFER_NFT | 26 |  |
| EVENT_TYPE_OPERATION_BURN_NFT | 27 |  |
| EVENT_TYPE_OPERATION_ROOT_CHANGED | 28 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/collection/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/collection/v1/genesis.proto



<a name="lbm.collection.v1.Balance"></a>

### Balance
Balance defines a balance of an address.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |
| `amount` | [Coin](#lbm.collection.v1.Coin) | repeated |  |






<a name="lbm.collection.v1.ClassStatistics"></a>

### ClassStatistics
ClassStatistics defines statistics belong to a token class.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  | class id associated with the token class. |
| `amount` | [string](#string) |  | statistics |






<a name="lbm.collection.v1.ContractAuthorizations"></a>

### ContractAuthorizations
ContractAuthorizations defines authorizations belong to a contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `authorizations` | [Authorization](#lbm.collection.v1.Authorization) | repeated | authorizations |






<a name="lbm.collection.v1.ContractBalances"></a>

### ContractBalances
ContractBalances defines balances belong to a contract.
genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `balances` | [Balance](#lbm.collection.v1.Balance) | repeated | balances |






<a name="lbm.collection.v1.ContractClasses"></a>

### ContractClasses
ContractClasses defines token classes belong to a contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `classes` | [google.protobuf.Any](#google.protobuf.Any) | repeated | classes |






<a name="lbm.collection.v1.ContractGrants"></a>

### ContractGrants
ContractGrant defines grants belong to a contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `grants` | [Grant](#lbm.collection.v1.Grant) | repeated | grants |






<a name="lbm.collection.v1.ContractNFTs"></a>

### ContractNFTs
ContractNFTs defines token classes belong to a contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `nfts` | [NFT](#lbm.collection.v1.NFT) | repeated | nfts |






<a name="lbm.collection.v1.ContractNextTokenIDs"></a>

### ContractNextTokenIDs
ContractNextTokenIDs defines the next token ids belong to a contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  |  |
| `token_ids` | [NextTokenID](#lbm.collection.v1.NextTokenID) | repeated |  |






<a name="lbm.collection.v1.ContractStatistics"></a>

### ContractStatistics
ContractStatistics defines statistics belong to a contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `statistics` | [ClassStatistics](#lbm.collection.v1.ClassStatistics) | repeated | statistics |






<a name="lbm.collection.v1.ContractTokenRelations"></a>

### ContractTokenRelations
ContractTokenRelations defines token relations belong to a contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `relations` | [TokenRelation](#lbm.collection.v1.TokenRelation) | repeated | relations |






<a name="lbm.collection.v1.GenesisState"></a>

### GenesisState
GenesisState defines the collection module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#lbm.collection.v1.Params) |  | params defines all the paramaters of the module. |
| `contracts` | [Contract](#lbm.collection.v1.Contract) | repeated | contracts defines the metadata of the contracts. |
| `next_class_ids` | [NextClassIDs](#lbm.collection.v1.NextClassIDs) | repeated | next ids for token classes. |
| `classes` | [ContractClasses](#lbm.collection.v1.ContractClasses) | repeated | classes defines the metadata of the tokens. |
| `next_token_ids` | [ContractNextTokenIDs](#lbm.collection.v1.ContractNextTokenIDs) | repeated | next ids for (non-fungible) tokens. |
| `balances` | [ContractBalances](#lbm.collection.v1.ContractBalances) | repeated | balances is an array containing the balances of all the accounts. |
| `nfts` | [ContractNFTs](#lbm.collection.v1.ContractNFTs) | repeated | nfts is an array containing the nfts. |
| `parents` | [ContractTokenRelations](#lbm.collection.v1.ContractTokenRelations) | repeated | parents represents the parents of (non-fungible) tokens. |
| `grants` | [ContractGrants](#lbm.collection.v1.ContractGrants) | repeated | grants defines the grant information. |
| `authorizations` | [ContractAuthorizations](#lbm.collection.v1.ContractAuthorizations) | repeated | authorizations defines the approve information. |
| `supplies` | [ContractStatistics](#lbm.collection.v1.ContractStatistics) | repeated | supplies represents the total supplies of tokens. |
| `burnts` | [ContractStatistics](#lbm.collection.v1.ContractStatistics) | repeated | burnts represents the total amount of burnt tokens. |






<a name="lbm.collection.v1.NextClassIDs"></a>

### NextClassIDs
NextClassIDs defines the next class ids of the contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `fungible` | [string](#string) |  | id for the fungible tokens. |
| `non_fungible` | [string](#string) |  | id for the non-fungible tokens. |






<a name="lbm.collection.v1.NextTokenID"></a>

### NextTokenID
NextTokenID defines the next (non-fungible) token id of the token class.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `class_id` | [string](#string) |  | class id associated with the token class. |
| `id` | [string](#string) |  | id for the token. |






<a name="lbm.collection.v1.TokenRelation"></a>

### TokenRelation
TokenRelation defines relations between two tokens.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `self` | [string](#string) |  | self |
| `other` | [string](#string) |  | other |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/collection/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/collection/v1/query.proto



<a name="lbm.collection.v1.QueryAllBalancesRequest"></a>

### QueryAllBalancesRequest
QueryAllBalancesRequest is the request type for the Query/AllBalances RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `address` | [string](#string) |  | address is the address to query the balances for. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.collection.v1.QueryAllBalancesResponse"></a>

### QueryAllBalancesResponse
QueryAllBalancesResponse is the response type for the Query/AllBalances RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `balances` | [Coin](#lbm.collection.v1.Coin) | repeated | balances is the balalces of all the tokens. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.collection.v1.QueryApprovedRequest"></a>

### QueryApprovedRequest
QueryApprovedRequest is the request type for the Query/Approved RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `address` | [string](#string) |  | the address of the proxy. |
| `approver` | [string](#string) |  | the address of the token approver. |






<a name="lbm.collection.v1.QueryApprovedResponse"></a>

### QueryApprovedResponse
QueryApprovedResponse is the response type for the Query/Approved RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `approved` | [bool](#bool) |  |  |






<a name="lbm.collection.v1.QueryApproversRequest"></a>

### QueryApproversRequest
QueryApproversRequest is the request type for the Query/Approvers RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `address` | [string](#string) |  | address of the proxy. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.collection.v1.QueryApproversResponse"></a>

### QueryApproversResponse
QueryApproversResponse is the response type for the Query/Approvers RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `approvers` | [string](#string) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.collection.v1.QueryBalanceRequest"></a>

### QueryBalanceRequest
QueryBalanceRequest is the request type for the Query/Balance RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `address` | [string](#string) |  | address is the address to query the balance for. |
| `token_id` | [string](#string) |  | token id associated with the token. |






<a name="lbm.collection.v1.QueryBalanceResponse"></a>

### QueryBalanceResponse
QueryBalanceResponse is the response type for the Query/Balance RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `balance` | [Coin](#lbm.collection.v1.Coin) |  | balance is the balance of the token. |






<a name="lbm.collection.v1.QueryChildrenRequest"></a>

### QueryChildrenRequest
QueryChildrenRequest is the request type for the Query/Children RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `token_id` | [string](#string) |  | token id associated with the non-fungible token. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.collection.v1.QueryChildrenResponse"></a>

### QueryChildrenResponse
QueryChildrenResponse is the response type for the Query/Children RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `children` | [NFT](#lbm.collection.v1.NFT) | repeated | children is the information of the child tokens. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.collection.v1.QueryContractRequest"></a>

### QueryContractRequest
QueryContractRequest is the request type for the Query/Contract RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |






<a name="lbm.collection.v1.QueryContractResponse"></a>

### QueryContractResponse
QueryContractResponse is the response type for the Query/Contract RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract` | [Contract](#lbm.collection.v1.Contract) |  | contract is the information of the contract. |






<a name="lbm.collection.v1.QueryFTBurntRequest"></a>

### QueryFTBurntRequest
QueryFTBurntRequest is the request type for the Query/FTBurnt RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `token_id` | [string](#string) |  | token id associated with the fungible token. |






<a name="lbm.collection.v1.QueryFTBurntResponse"></a>

### QueryFTBurntResponse
QueryFTBurntResponse is the response type for the Query/FTBurnt RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `burnt` | [string](#string) |  | burnt is the amount of the burnt tokens. |






<a name="lbm.collection.v1.QueryFTMintedRequest"></a>

### QueryFTMintedRequest
QueryFTMintedRequest is the request type for the Query/FTMinted RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `token_id` | [string](#string) |  | token id associated with the fungible token. |






<a name="lbm.collection.v1.QueryFTMintedResponse"></a>

### QueryFTMintedResponse
QueryFTMintedResponse is the response type for the Query/FTMinted RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `minted` | [string](#string) |  | minted is the amount of the minted tokens. |






<a name="lbm.collection.v1.QueryFTSupplyRequest"></a>

### QueryFTSupplyRequest
QueryFTSupplyRequest is the request type for the Query/FTSupply RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `token_id` | [string](#string) |  | token id associated with the fungible token. |






<a name="lbm.collection.v1.QueryFTSupplyResponse"></a>

### QueryFTSupplyResponse
QueryFTSupplyResponse is the response type for the Query/FTSupply RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `supply` | [string](#string) |  | supply is the supply of the tokens. |






<a name="lbm.collection.v1.QueryGranteeGrantsRequest"></a>

### QueryGranteeGrantsRequest
QueryGranteeGrantsRequest is the request type for the Query/GranteeGrants RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `grantee` | [string](#string) |  | the address of the grantee. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.collection.v1.QueryGranteeGrantsResponse"></a>

### QueryGranteeGrantsResponse
QueryGranteeGrantsResponse is the response type for the Query/GranteeGrants RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grants` | [Grant](#lbm.collection.v1.Grant) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.collection.v1.QueryNFTBurntRequest"></a>

### QueryNFTBurntRequest
QueryNFTBurntRequest is the request type for the Query/NFTBurnt RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `token_type` | [string](#string) |  | token type associated with the token type. |






<a name="lbm.collection.v1.QueryNFTBurntResponse"></a>

### QueryNFTBurntResponse
QueryNFTBurntResponse is the response type for the Query/NFTBurnt RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `burnt` | [string](#string) |  | burnt is the amount of the burnt tokens. |






<a name="lbm.collection.v1.QueryNFTMintedRequest"></a>

### QueryNFTMintedRequest
QueryNFTMintedRequest is the request type for the Query/NFTMinted RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `token_type` | [string](#string) |  | token type associated with the token type. |






<a name="lbm.collection.v1.QueryNFTMintedResponse"></a>

### QueryNFTMintedResponse
QueryNFTMintedResponse is the response type for the Query/NFTMinted RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `minted` | [string](#string) |  | minted is the amount of minted tokens. |






<a name="lbm.collection.v1.QueryNFTSupplyRequest"></a>

### QueryNFTSupplyRequest
QueryNFTSupplyRequest is the request type for the Query/NFTSupply RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `token_type` | [string](#string) |  | token type associated with the token type. |






<a name="lbm.collection.v1.QueryNFTSupplyResponse"></a>

### QueryNFTSupplyResponse
QueryNFTSupplyResponse is the response type for the Query/NFTSupply RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `supply` | [string](#string) |  | supply is the supply of the non-fungible token. |






<a name="lbm.collection.v1.QueryParentRequest"></a>

### QueryParentRequest
QueryParentRequest is the request type for the Query/Parent RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `token_id` | [string](#string) |  | token id associated wit the non-fungible token. |






<a name="lbm.collection.v1.QueryParentResponse"></a>

### QueryParentResponse
QueryParentResponse is the response type for the Query/Parent RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `parent` | [NFT](#lbm.collection.v1.NFT) |  | parent is the information of the parent token. if there is no parent for the token, it would return nil. |






<a name="lbm.collection.v1.QueryRootRequest"></a>

### QueryRootRequest
QueryRootRequest is the request type for the Query/Root RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `token_id` | [string](#string) |  | token id associated with the non-fungible token. |






<a name="lbm.collection.v1.QueryRootResponse"></a>

### QueryRootResponse
QueryRootResponse is the response type for the Query/Root RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `root` | [NFT](#lbm.collection.v1.NFT) |  | root is the information of the root token. it would return itself if it's the root token. |






<a name="lbm.collection.v1.QueryTokenClassTypeNameRequest"></a>

### QueryTokenClassTypeNameRequest
QueryTokenClassTypeNameRequest is the request type for the Query/TokenClassTypeName RPC method.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `class_id` | [string](#string) |  | class id associated with the token class. |






<a name="lbm.collection.v1.QueryTokenClassTypeNameResponse"></a>

### QueryTokenClassTypeNameResponse
QueryTokenClassTypeNameResponse is the response type for the Query/TokenClassTypeName RPC method.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | type name of the token class. |






<a name="lbm.collection.v1.QueryTokenRequest"></a>

### QueryTokenRequest
QueryTokenRequest is the request type for the Query/Token RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `token_id` | [string](#string) |  | token id associated with the fungible token. |






<a name="lbm.collection.v1.QueryTokenResponse"></a>

### QueryTokenResponse
QueryTokenResponse is the response type for the Query/Token RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token` | [google.protobuf.Any](#google.protobuf.Any) |  | information of the token. |






<a name="lbm.collection.v1.QueryTokenTypeRequest"></a>

### QueryTokenTypeRequest
QueryTokenTypeRequest is the request type for the Query/TokenType RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `token_type` | [string](#string) |  | token type associated with the token type. |






<a name="lbm.collection.v1.QueryTokenTypeResponse"></a>

### QueryTokenTypeResponse
QueryTokenTypeResponse is the response type for the Query/TokenType RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_type` | [TokenType](#lbm.collection.v1.TokenType) |  | token type is the information of the token type. |






<a name="lbm.collection.v1.QueryTokenTypesRequest"></a>

### QueryTokenTypesRequest
QueryTokenTypesRequest is the request type for the Query/TokenTypes RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.collection.v1.QueryTokenTypesResponse"></a>

### QueryTokenTypesResponse
QueryTokenTypesResponse is the response type for the Query/TokenTypes RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_types` | [TokenType](#lbm.collection.v1.TokenType) | repeated | token types is the informations of all the token types. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.collection.v1.QueryTokensRequest"></a>

### QueryTokensRequest
QueryTokensRequest is the request type for the Query/Tokens RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.collection.v1.QueryTokensResponse"></a>

### QueryTokensResponse
QueryTokensResponse is the response type for the Query/Tokens RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tokens` | [google.protobuf.Any](#google.protobuf.Any) | repeated | informations of all the tokens. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.collection.v1.QueryTokensWithTokenTypeRequest"></a>

### QueryTokensWithTokenTypeRequest
QueryTokensWithTokenTypeRequest is the request type for the Query/TokensWithTokenType RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `token_type` | [string](#string) |  | token type associated with the token type. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.collection.v1.QueryTokensWithTokenTypeResponse"></a>

### QueryTokensWithTokenTypeResponse
QueryTokensWithTokenTypeResponse is the response type for the Query/TokensWithTokenType RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tokens` | [google.protobuf.Any](#google.protobuf.Any) | repeated | informations of all the tokens. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.collection.v1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Balance` | [QueryBalanceRequest](#lbm.collection.v1.QueryBalanceRequest) | [QueryBalanceResponse](#lbm.collection.v1.QueryBalanceResponse) | Balance queries the balance of a single token class for a single account. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - ErrInvalidAddress - `address` is of invalid format. | GET|/lbm/collection/v1/contracts/{contract_id}/balances/{address}/{token_id}|
| `AllBalances` | [QueryAllBalancesRequest](#lbm.collection.v1.QueryAllBalancesRequest) | [QueryAllBalancesResponse](#lbm.collection.v1.QueryAllBalancesResponse) | AllBalances queries the balance of all token classes for a single account. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - ErrInvalidAddress - `address` is of invalid format. | GET|/lbm/collection/v1/contracts/{contract_id}/balances/{address}|
| `FTSupply` | [QueryFTSupplyRequest](#lbm.collection.v1.QueryFTSupplyRequest) | [QueryFTSupplyResponse](#lbm.collection.v1.QueryFTSupplyResponse) | FTSupply queries the number of tokens from a given contract id and token id. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - `class_id` is of invalid format. - ErrNotFound - there is no token class of `class_id`. | GET|/lbm/collection/v1/contracts/{contract_id}/fts/{token_id}/supply|
| `FTMinted` | [QueryFTMintedRequest](#lbm.collection.v1.QueryFTMintedRequest) | [QueryFTMintedResponse](#lbm.collection.v1.QueryFTMintedResponse) | FTMinted queries the number of minted tokens from a given contract id and token id. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - `class_id` is of invalid format. - ErrNotFound - there is no token class of `class_id`. | GET|/lbm/collection/v1/contracts/{contract_id}/fts/{token_id}/minted|
| `FTBurnt` | [QueryFTBurntRequest](#lbm.collection.v1.QueryFTBurntRequest) | [QueryFTBurntResponse](#lbm.collection.v1.QueryFTBurntResponse) | FTBurnt queries the number of burnt tokens from a given contract id and token id. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - `class_id` is of invalid format. - ErrNotFound - there is no token class of `class_id`. | GET|/lbm/collection/v1/contracts/{contract_id}/fts/{token_id}/burnt|
| `NFTSupply` | [QueryNFTSupplyRequest](#lbm.collection.v1.QueryNFTSupplyRequest) | [QueryNFTSupplyResponse](#lbm.collection.v1.QueryNFTSupplyResponse) | NFTSupply queries the number of tokens from a given contract id and token type. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - `class_id` is of invalid format. - ErrNotFound - there is no token class of `class_id`. | GET|/lbm/collection/v1/contracts/{contract_id}/token_types/{token_type}/supply|
| `NFTMinted` | [QueryNFTMintedRequest](#lbm.collection.v1.QueryNFTMintedRequest) | [QueryNFTMintedResponse](#lbm.collection.v1.QueryNFTMintedResponse) | NFTMinted queries the number of minted tokens from a given contract id and token type. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - `class_id` is of invalid format. - ErrNotFound - there is no token class of `class_id`. | GET|/lbm/collection/v1/contracts/{contract_id}/token_types/{token_type}/minted|
| `NFTBurnt` | [QueryNFTBurntRequest](#lbm.collection.v1.QueryNFTBurntRequest) | [QueryNFTBurntResponse](#lbm.collection.v1.QueryNFTBurntResponse) | NFTBurnt queries the number of burnt tokens from a given contract id and token type. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - `class_id` is of invalid format. - ErrNotFound - there is no token class of `class_id`. | GET|/lbm/collection/v1/contracts/{contract_id}/token_types/{token_type}/burnt|
| `Contract` | [QueryContractRequest](#lbm.collection.v1.QueryContractRequest) | [QueryContractResponse](#lbm.collection.v1.QueryContractResponse) | Contract queries a contract metadata based on its contract id. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - ErrNotFound - there is no contract of `contract_id`. | GET|/lbm/collection/v1/contracts/{contract_id}|
| `TokenClassTypeName` | [QueryTokenClassTypeNameRequest](#lbm.collection.v1.QueryTokenClassTypeNameRequest) | [QueryTokenClassTypeNameResponse](#lbm.collection.v1.QueryTokenClassTypeNameResponse) | TokenClassTypeName queries the fully qualified message type name of a token class from its class id. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - `class_id` is of invalid format. - ErrNotFound - there is no token class of `class_id`. Since: 0.46.0 (finschia) | GET|/lbm/collection/v1/contracts/{contract_id}/token_classes/{class_id}/type_name|
| `TokenType` | [QueryTokenTypeRequest](#lbm.collection.v1.QueryTokenTypeRequest) | [QueryTokenTypeResponse](#lbm.collection.v1.QueryTokenTypeResponse) | TokenType queries metadata of a token type. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - `class_id` is of invalid format. - ErrNotFound - there is no token class of `class_id`. | GET|/lbm/collection/v1/contracts/{contract_id}/token_types/{token_type}|
| `TokenTypes` | [QueryTokenTypesRequest](#lbm.collection.v1.QueryTokenTypesRequest) | [QueryTokenTypesResponse](#lbm.collection.v1.QueryTokenTypesResponse) | TokenTypes queries metadata of all the token types. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - ErrNotFound - there is no token contract of `contract_id`. | GET|/lbm/collection/v1/contracts/{contract_id}/token_types|
| `Token` | [QueryTokenRequest](#lbm.collection.v1.QueryTokenRequest) | [QueryTokenResponse](#lbm.collection.v1.QueryTokenResponse) | Token queries a metadata of a token from its token id. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - `token_id` is of invalid format. - ErrNotFound - there is no token of `token_id`. | GET|/lbm/collection/v1/contracts/{contract_id}/tokens/{token_id}|
| `TokensWithTokenType` | [QueryTokensWithTokenTypeRequest](#lbm.collection.v1.QueryTokensWithTokenTypeRequest) | [QueryTokensWithTokenTypeResponse](#lbm.collection.v1.QueryTokensWithTokenTypeResponse) | TokensWithTokenType queries all token metadata with token type. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - `token_type` is of invalid format. - ErrNotFound - there is no contract of `contract_id`. | GET|/lbm/collection/v1/contracts/{contract_id}/token_types/{token_type}/tokens|
| `Tokens` | [QueryTokensRequest](#lbm.collection.v1.QueryTokensRequest) | [QueryTokensResponse](#lbm.collection.v1.QueryTokensResponse) | Tokens queries all token metadata. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - ErrNotFound - there is no contract of `contract_id`. | GET|/lbm/collection/v1/contracts/{contract_id}/tokens|
| `Root` | [QueryRootRequest](#lbm.collection.v1.QueryRootRequest) | [QueryRootResponse](#lbm.collection.v1.QueryRootResponse) | Root queries the root of a given nft. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - `token_id` is of invalid format. - ErrNotFound - there is no token of `token_id`. | GET|/lbm/collection/v1/contracts/{contract_id}/nfts/{token_id}/root|
| `Parent` | [QueryParentRequest](#lbm.collection.v1.QueryParentRequest) | [QueryParentResponse](#lbm.collection.v1.QueryParentResponse) | Parent queries the parent of a given nft. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - `token_id` is of invalid format. - ErrNotFound - there is no token of `token_id`. - token is the root. | GET|/lbm/collection/v1/contracts/{contract_id}/nfts/{token_id}/parent|
| `Children` | [QueryChildrenRequest](#lbm.collection.v1.QueryChildrenRequest) | [QueryChildrenResponse](#lbm.collection.v1.QueryChildrenResponse) | Children queries the children of a given nft. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - `token_id` is of invalid format. | GET|/lbm/collection/v1/contracts/{contract_id}/nfts/{token_id}/children|
| `GranteeGrants` | [QueryGranteeGrantsRequest](#lbm.collection.v1.QueryGranteeGrantsRequest) | [QueryGranteeGrantsResponse](#lbm.collection.v1.QueryGranteeGrantsResponse) | GranteeGrants queries all permissions on a given grantee. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - ErrInvalidAddress - `grantee` is of invalid format. | GET|/lbm/collection/v1/contracts/{contract_id}/grants/{grantee}|
| `Approved` | [QueryApprovedRequest](#lbm.collection.v1.QueryApprovedRequest) | [QueryApprovedResponse](#lbm.collection.v1.QueryApprovedResponse) | Approved queries whether the proxy is approved by the approver. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - ErrInvalidAddress - `proxy` is of invalid format. - `approver` is of invalid format. - ErrNotFound - there is no authorization given by `approver` to `proxy`. | GET|/lbm/collection/v1/contracts/{contract_id}/accounts/{address}/proxies/{approver}|
| `Approvers` | [QueryApproversRequest](#lbm.collection.v1.QueryApproversRequest) | [QueryApproversResponse](#lbm.collection.v1.QueryApproversResponse) | Approvers queries approvers of a given proxy. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - ErrInvalidAddress - `proxy` is of invalid format. | GET|/lbm/collection/v1/contracts/{contract_id}/accounts/{address}/approvers|

 <!-- end services -->



<a name="lbm/collection/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/collection/v1/tx.proto



<a name="lbm.collection.v1.MintNFTParam"></a>

### MintNFTParam
MintNFTParam defines a parameter for minting nft.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_type` | [string](#string) |  | token type or class id of the nft. Note: it cannot start with zero. |
| `name` | [string](#string) |  | name defines the human-readable name of the nft (mandatory). Note: it has an app-specific limit in length. |
| `meta` | [string](#string) |  | meta is a brief description of the nft. Note: it has an app-specific limit in length. |






<a name="lbm.collection.v1.MsgApprove"></a>

### MsgApprove
MsgApprove is the Msg/Approve request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `approver` | [string](#string) |  | address of the approver who allows the manipulation of its token. |
| `proxy` | [string](#string) |  | address which the manipulation is allowed to. |






<a name="lbm.collection.v1.MsgApproveResponse"></a>

### MsgApproveResponse
MsgApproveResponse is the Msg/Approve response type.






<a name="lbm.collection.v1.MsgAttach"></a>

### MsgAttach
MsgAttach is the Msg/Attach request type.

Throws:
- ErrInvalidAddress
  - `from` is of invalid format.
- ErrInvalidRequest
  - `contract_id` is of invalid format.
  - `token_id` is of invalid format.
  - `to_token_id` is of invalid format.
  - `token_id` is equal to `to_token_id`.

Signer: `from`


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `from` | [string](#string) |  | address of the owner of the token. |
| `token_id` | [string](#string) |  | token id of the token to attach. |
| `to_token_id` | [string](#string) |  | to token id which one attachs the token to. |






<a name="lbm.collection.v1.MsgAttachFrom"></a>

### MsgAttachFrom
MsgAttachFrom is the Msg/AttachFrom request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `proxy` | [string](#string) |  | address of the proxy. |
| `from` | [string](#string) |  | address of the owner of the token. |
| `token_id` | [string](#string) |  | token id of the token to attach. |
| `to_token_id` | [string](#string) |  | to token id which one attachs the token to. |






<a name="lbm.collection.v1.MsgAttachFromResponse"></a>

### MsgAttachFromResponse
MsgAttachFromResponse is the Msg/AttachFrom response type.






<a name="lbm.collection.v1.MsgAttachResponse"></a>

### MsgAttachResponse
MsgAttachResponse is the Msg/Attach response type.






<a name="lbm.collection.v1.MsgBurnFT"></a>

### MsgBurnFT
MsgBurnFT is the Msg/BurnFT request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `from` | [string](#string) |  | address which the tokens will be burnt from. Note: it must have the permission for the burn. |
| `amount` | [Coin](#lbm.collection.v1.Coin) | repeated | the amount of the burn. Note: amount may be empty. |






<a name="lbm.collection.v1.MsgBurnFTFrom"></a>

### MsgBurnFTFrom
MsgBurnFTFrom is the Msg/BurnFTFrom request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `proxy` | [string](#string) |  | address which triggers the burn. Note: it must have the permission for the burn. Note: it must have been authorized by from. |
| `from` | [string](#string) |  | address which the tokens will be burnt from. |
| `amount` | [Coin](#lbm.collection.v1.Coin) | repeated | the amount of the burn. Note: amount may be empty. |






<a name="lbm.collection.v1.MsgBurnFTFromResponse"></a>

### MsgBurnFTFromResponse
MsgBurnFTFromResponse is the Msg/BurnFTFrom response type.






<a name="lbm.collection.v1.MsgBurnFTResponse"></a>

### MsgBurnFTResponse
MsgBurnFTResponse is the Msg/BurnFT response type.






<a name="lbm.collection.v1.MsgBurnNFT"></a>

### MsgBurnNFT
MsgBurnNFT is the Msg/BurnNFT request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `from` | [string](#string) |  | address which the tokens will be burnt from. Note: it must have the permission for the burn. |
| `token_ids` | [string](#string) | repeated | the token ids to burn. Note: id cannot start with zero. |






<a name="lbm.collection.v1.MsgBurnNFTFrom"></a>

### MsgBurnNFTFrom
MsgBurnNFTFrom is the Msg/BurnNFTFrom request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `proxy` | [string](#string) |  | address which triggers the burn. Note: it must have the permission for the burn. Note: it must have been authorized by from. |
| `from` | [string](#string) |  | address which the tokens will be burnt from. |
| `token_ids` | [string](#string) | repeated | the token ids to burn. Note: id cannot start with zero. |






<a name="lbm.collection.v1.MsgBurnNFTFromResponse"></a>

### MsgBurnNFTFromResponse
MsgBurnNFTFromResponse is the Msg/BurnNFTFrom response type.






<a name="lbm.collection.v1.MsgBurnNFTResponse"></a>

### MsgBurnNFTResponse
MsgBurnNFTResponse is the Msg/BurnNFT response type.






<a name="lbm.collection.v1.MsgCreateContract"></a>

### MsgCreateContract
MsgCreateContract is the Msg/CreateContract request type.

Throws:
- ErrInvalidAddress
  - `owner` is of invalid format.
- ErrInvalidRequest
  - `name` exceeds the app-specific limit in length.
  - `base_img_uri` exceeds the app-specific limit in length.
  - `meta` exceeds the app-specific limit in length.

Signer: `owner`


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `owner` | [string](#string) |  | address which all the permissions on the contract will be granted to (not a permanent property). |
| `name` | [string](#string) |  | name defines the human-readable name of the contract. |
| `base_img_uri` | [string](#string) |  | base img uri is an uri for the contract image stored off chain. |
| `meta` | [string](#string) |  | meta is a brief description of the contract. |






<a name="lbm.collection.v1.MsgCreateContractResponse"></a>

### MsgCreateContractResponse
MsgCreateContractResponse is the Msg/CreateContract response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | id of the new contract. |






<a name="lbm.collection.v1.MsgDetach"></a>

### MsgDetach
MsgDetach is the Msg/Detach request type.

Throws:
- ErrInvalidAddress
  - `from` is of invalid format.
- ErrInvalidRequest
  - `contract_id` is of invalid format.
  - `token_id` is of invalid format.

Signer: `from`


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `from` | [string](#string) |  | address of the owner of the token. |
| `token_id` | [string](#string) |  | token id of the token to detach. |






<a name="lbm.collection.v1.MsgDetachFrom"></a>

### MsgDetachFrom
MsgDetachFrom is the Msg/DetachFrom request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `proxy` | [string](#string) |  | address of the proxy. |
| `from` | [string](#string) |  | address of the owner of the token. |
| `token_id` | [string](#string) |  | token id of the token to detach. |






<a name="lbm.collection.v1.MsgDetachFromResponse"></a>

### MsgDetachFromResponse
MsgDetachFromResponse is the Msg/DetachFrom response type.






<a name="lbm.collection.v1.MsgDetachResponse"></a>

### MsgDetachResponse
MsgDetachResponse is the Msg/Detach response type.






<a name="lbm.collection.v1.MsgDisapprove"></a>

### MsgDisapprove
MsgDisapprove is the Msg/Disapprove request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `approver` | [string](#string) |  | address of the approver who allows the manipulation of its token. |
| `proxy` | [string](#string) |  | address which the manipulation is allowed to. |






<a name="lbm.collection.v1.MsgDisapproveResponse"></a>

### MsgDisapproveResponse
MsgDisapproveResponse is the Msg/Disapprove response type.






<a name="lbm.collection.v1.MsgGrantPermission"></a>

### MsgGrantPermission
MsgGrantPermission is the Msg/GrantPermission request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `from` | [string](#string) |  | address of the granter which must have the permission to give. |
| `to` | [string](#string) |  | address of the grantee. |
| `permission` | [string](#string) |  | permission on the contract. |






<a name="lbm.collection.v1.MsgGrantPermissionResponse"></a>

### MsgGrantPermissionResponse
MsgGrantPermissionResponse is the Msg/GrantPermission response type.






<a name="lbm.collection.v1.MsgIssueFT"></a>

### MsgIssueFT
MsgIssueFT is the Msg/IssueFT request type.

Throws:
- ErrInvalidAddress
  - `owner` is of invalid format.
  - `to` is of invalid format.
- ErrInvalidRequest
  - `contract_id` is of invalid format.
  - `name` is empty.
  - `name` exceeds the app-specific limit in length.
  - `meta` exceeds the app-specific limit in length.
  - `decimals` is lesser than 0 or greater than 18.
  - `amount` is not positive.
  - `mintable` == false, amount == 1 and decimals == 0 (weird, but for the backward compatibility).

Signer: `owner`


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `name` | [string](#string) |  | name defines the human-readable name of the token type. |
| `meta` | [string](#string) |  | meta is a brief description of the token type. |
| `decimals` | [int32](#int32) |  | decimals is the number of decimals which one must divide the amount by to get its user representation. |
| `mintable` | [bool](#bool) |  | mintable represents whether the token is allowed to be minted or burnt. |
| `owner` | [string](#string) |  | the address of the grantee which must have the permission to issue a token. |
| `to` | [string](#string) |  | the address to send the minted tokens to. mandatory. |
| `amount` | [string](#string) |  | the amount of tokens to mint on the issuance. Note: if you provide negative amount, a panic may result. Note: amount may be zero. |






<a name="lbm.collection.v1.MsgIssueFTResponse"></a>

### MsgIssueFTResponse
MsgIssueFTResponse is the Msg/IssueFT response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_id` | [string](#string) |  | id of the token. |






<a name="lbm.collection.v1.MsgIssueNFT"></a>

### MsgIssueNFT
MsgIssueNFT is the Msg/IssueNFT request type.

Throws:
- ErrInvalidAddress
  - `owner` is of invalid format.
- ErrInvalidRequest
  - `contract_id` is of invalid format.
  - `name` exceeds the app-specific limit in length.
  - `meta` exceeds the app-specific limit in length.

Signer: `owner`


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `name` | [string](#string) |  | name defines the human-readable name of the token type. |
| `meta` | [string](#string) |  | meta is a brief description of the token type. |
| `owner` | [string](#string) |  | the address of the grantee which must have the permission to issue a token. |






<a name="lbm.collection.v1.MsgIssueNFTResponse"></a>

### MsgIssueNFTResponse
MsgIssueNFTResponse is the Msg/IssueNFT response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_type` | [string](#string) |  | id of the new token type. |






<a name="lbm.collection.v1.MsgMintFT"></a>

### MsgMintFT
MsgMintFT is the Msg/MintFT request type.

Throws:
- ErrInvalidAddress
  - `from` is of invalid format.
  - `to` is of invalid format.
- ErrInvalidRequest
  - `contract_id` is of invalid format.
  - `amount` is not positive.

Signer: `from`


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `from` | [string](#string) |  | address of the grantee which has the permission for the mint. |
| `to` | [string](#string) |  | address which the minted tokens will be sent to. |
| `amount` | [Coin](#lbm.collection.v1.Coin) | repeated | the amount of the mint. Note: amount may be empty. |






<a name="lbm.collection.v1.MsgMintFTResponse"></a>

### MsgMintFTResponse
MsgMintFTResponse is the Msg/MintFT response type.






<a name="lbm.collection.v1.MsgMintNFT"></a>

### MsgMintNFT
MsgMintNFT is the Msg/MintNFT request type.

Throws:
- ErrInvalidAddress
  - `from` is of invalid format.
  - `to` is of invalid format.
- ErrInvalidRequest
  - `contract_id` is of invalid format.
  - `params` is empty.
  - `params` has an invalid element.

Signer: `from`


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `from` | [string](#string) |  | address of the grantee which has the permission for the mint. |
| `to` | [string](#string) |  | address which the minted token will be sent to. |
| `params` | [MintNFTParam](#lbm.collection.v1.MintNFTParam) | repeated | parameters for the minted tokens. |






<a name="lbm.collection.v1.MsgMintNFTResponse"></a>

### MsgMintNFTResponse
MsgMintNFTResponse is the Msg/MintNFT response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_ids` | [string](#string) | repeated | ids of the new non-fungible tokens. |






<a name="lbm.collection.v1.MsgModify"></a>

### MsgModify
MsgModify is the Msg/Modify request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `owner` | [string](#string) |  | the address of the grantee which must have modify permission. |
| `token_type` | [string](#string) |  | token type of the token. |
| `token_index` | [string](#string) |  | token index of the token. if index is empty, it would modify the corresponding token type. if index is not empty, it would modify the corresponding nft. Note: if token type is of FTs, the index cannot be empty. |
| `changes` | [Change](#lbm.collection.v1.Change) | repeated | changes to apply. on modifying collection: name, base_img_uri, meta. on modifying token type and token: name, meta. |






<a name="lbm.collection.v1.MsgModifyResponse"></a>

### MsgModifyResponse
MsgModifyResponse is the Msg/Modify response type.






<a name="lbm.collection.v1.MsgRevokePermission"></a>

### MsgRevokePermission
MsgRevokePermission is the Msg/RevokePermission request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `from` | [string](#string) |  | address of the grantee which abandons the permission. |
| `permission` | [string](#string) |  | permission on the contract. |






<a name="lbm.collection.v1.MsgRevokePermissionResponse"></a>

### MsgRevokePermissionResponse
MsgRevokePermissionResponse is the Msg/RevokePermission response type.






<a name="lbm.collection.v1.MsgTransferFT"></a>

### MsgTransferFT
MsgTransferFT is the Msg/TransferFT request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `from` | [string](#string) |  | the address which the transfer is from. |
| `to` | [string](#string) |  | the address which the transfer is to. |
| `amount` | [Coin](#lbm.collection.v1.Coin) | repeated | the amount of the transfer. Note: amount may be empty. |






<a name="lbm.collection.v1.MsgTransferFTFrom"></a>

### MsgTransferFTFrom
MsgTransferFTFrom is the Msg/TransferFTFrom request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `proxy` | [string](#string) |  | the address of the proxy. |
| `from` | [string](#string) |  | the address which the transfer is from. |
| `to` | [string](#string) |  | the address which the transfer is to. |
| `amount` | [Coin](#lbm.collection.v1.Coin) | repeated | the amount of the transfer. Note: amount may be empty. |






<a name="lbm.collection.v1.MsgTransferFTFromResponse"></a>

### MsgTransferFTFromResponse
MsgTransferFTFromResponse is the Msg/TransferFTFrom response type.






<a name="lbm.collection.v1.MsgTransferFTResponse"></a>

### MsgTransferFTResponse
MsgTransferFTResponse is the Msg/TransferFT response type.






<a name="lbm.collection.v1.MsgTransferNFT"></a>

### MsgTransferNFT
MsgTransferNFT is the Msg/TransferNFT request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `from` | [string](#string) |  | the address which the transfer is from. |
| `to` | [string](#string) |  | the address which the transfer is to. |
| `token_ids` | [string](#string) | repeated | the token ids to transfer. |






<a name="lbm.collection.v1.MsgTransferNFTFrom"></a>

### MsgTransferNFTFrom
MsgTransferNFTFrom is the Msg/TransferNFTFrom request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `proxy` | [string](#string) |  | the address of the proxy. |
| `from` | [string](#string) |  | the address which the transfer is from. |
| `to` | [string](#string) |  | the address which the transfer is to. |
| `token_ids` | [string](#string) | repeated | the token ids to transfer. |






<a name="lbm.collection.v1.MsgTransferNFTFromResponse"></a>

### MsgTransferNFTFromResponse
MsgTransferNFTFromResponse is the Msg/TransferNFTFrom response type.






<a name="lbm.collection.v1.MsgTransferNFTResponse"></a>

### MsgTransferNFTResponse
MsgTransferNFTResponse is the Msg/TransferNFT response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.collection.v1.Msg"></a>

### Msg
Msg defines the collection Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `TransferFT` | [MsgTransferFT](#lbm.collection.v1.MsgTransferFT) | [MsgTransferFTResponse](#lbm.collection.v1.MsgTransferFTResponse) | TransferFT defines a method to send fungible tokens from one account to another account. Fires: - EventSent - transfer_ft (deprecated, not typed) Throws: - ErrInvalidRequest: - the balance of `from` does not have enough tokens to spend. | |
| `TransferFTFrom` | [MsgTransferFTFrom](#lbm.collection.v1.MsgTransferFTFrom) | [MsgTransferFTFromResponse](#lbm.collection.v1.MsgTransferFTFromResponse) | TransferFTFrom defines a method to send fungible tokens from one account to another account by the proxy. Fires: - EventSent - transfer_ft_from (deprecated, not typed) Throws: - ErrUnauthorized: - the approver has not authorized the proxy. - ErrInvalidRequest: - the balance of `from` does not have enough tokens to spend. | |
| `TransferNFT` | [MsgTransferNFT](#lbm.collection.v1.MsgTransferNFT) | [MsgTransferNFTResponse](#lbm.collection.v1.MsgTransferNFTResponse) | TransferNFT defines a method to send non-fungible tokens from one account to another account. Fires: - EventSent - transfer_nft (deprecated, not typed) - operation_transfer_nft (deprecated, not typed) Throws: - ErrInvalidRequest: - the balance of `from` does not have enough tokens to spend. | |
| `TransferNFTFrom` | [MsgTransferNFTFrom](#lbm.collection.v1.MsgTransferNFTFrom) | [MsgTransferNFTFromResponse](#lbm.collection.v1.MsgTransferNFTFromResponse) | TransferNFTFrom defines a method to send non-fungible tokens from one account to another account by the proxy. Fires: - EventSent - transfer_nft_from (deprecated, not typed) - operation_transfer_nft (deprecated, not typed) Throws: - ErrUnauthorized: - the approver has not authorized the proxy. - ErrInvalidRequest: - the balance of `from` does not have enough tokens to spend. | |
| `Approve` | [MsgApprove](#lbm.collection.v1.MsgApprove) | [MsgApproveResponse](#lbm.collection.v1.MsgApproveResponse) | Approve allows one to send tokens on behalf of the approver. Fires: - EventAuthorizedOperator - approve_collection (deprecated, not typed) Throws: - ErrNotFound: - there is no contract of `contract_id`. - ErrInvalidRequest: - `approver` has already authorized `proxy`. | |
| `Disapprove` | [MsgDisapprove](#lbm.collection.v1.MsgDisapprove) | [MsgDisapproveResponse](#lbm.collection.v1.MsgDisapproveResponse) | Disapprove revokes the authorization of the proxy to send the approver's token. Fires: - EventRevokedOperator - disapprove_collection (deprecated, not typed) Throws: - ErrNotFound: - there is no contract of `contract_id`. - there is no authorization by `approver` to `proxy`. | |
| `CreateContract` | [MsgCreateContract](#lbm.collection.v1.MsgCreateContract) | [MsgCreateContractResponse](#lbm.collection.v1.MsgCreateContractResponse) | CreateContract defines a method to create a contract for collection. it grants `mint`, `burn`, `modify` and `issue` permissions on the contract to its creator. Fires: - EventCreatedContract - create_collection (deprecated, not typed) | |
| `IssueFT` | [MsgIssueFT](#lbm.collection.v1.MsgIssueFT) | [MsgIssueFTResponse](#lbm.collection.v1.MsgIssueFTResponse) | IssueFT defines a method to create a class of fungible token. Fires: - EventCreatedFTClass - EventMintedFT - issue_ft (deprecated, not typed) Note: it does not grant any permissions to its issuer. | |
| `IssueNFT` | [MsgIssueNFT](#lbm.collection.v1.MsgIssueNFT) | [MsgIssueNFTResponse](#lbm.collection.v1.MsgIssueNFTResponse) | IssueNFT defines a method to create a class of non-fungible token. Fires: - EventCreatedNFTClass - issue_nft (deprecated, not typed) Note: it DOES grant `mint` and `burn` permissions to its issuer. | |
| `MintFT` | [MsgMintFT](#lbm.collection.v1.MsgMintFT) | [MsgMintFTResponse](#lbm.collection.v1.MsgMintFTResponse) | MintFT defines a method to mint fungible tokens. Fires: - EventMintedFT - mint_ft (deprecated, not typed) Throws: - ErrUnauthorized - `from` does not have `mint` permission. | |
| `MintNFT` | [MsgMintNFT](#lbm.collection.v1.MsgMintNFT) | [MsgMintNFTResponse](#lbm.collection.v1.MsgMintNFTResponse) | MintNFT defines a method to mint non-fungible tokens. Fires: - EventMintedNFT - mint_nft (deprecated, not typed) Throws: - ErrUnauthorized - `from` does not have `mint` permission. | |
| `BurnFT` | [MsgBurnFT](#lbm.collection.v1.MsgBurnFT) | [MsgBurnFTResponse](#lbm.collection.v1.MsgBurnFTResponse) | BurnFT defines a method to burn fungible tokens. Fires: - EventBurned - burn_ft (deprecated, not typed) - burn_nft (deprecated, not typed) - operation_burn_nft (deprecated, not typed) Throws: - ErrUnauthorized - `from` does not have `burn` permission. - ErrInvalidRequest: - the balance of `from` does not have enough tokens to burn. | |
| `BurnFTFrom` | [MsgBurnFTFrom](#lbm.collection.v1.MsgBurnFTFrom) | [MsgBurnFTFromResponse](#lbm.collection.v1.MsgBurnFTFromResponse) | BurnFTFrom defines a method to burn fungible tokens of the approver by the proxy. Fires: - EventBurned - burn_ft_from (deprecated, not typed) - burn_nft_from (deprecated, not typed) - operation_burn_nft (deprecated, not typed) Throws: - ErrUnauthorized - `proxy` does not have `burn` permission. - the approver has not authorized `proxy`. - ErrInvalidRequest: - the balance of `from` does not have enough tokens to burn. | |
| `BurnNFT` | [MsgBurnNFT](#lbm.collection.v1.MsgBurnNFT) | [MsgBurnNFTResponse](#lbm.collection.v1.MsgBurnNFTResponse) | BurnNFT defines a method to burn non-fungible tokens. Fires: - EventBurned - burn_ft (deprecated, not typed) - burn_nft (deprecated, not typed) - operation_burn_nft (deprecated, not typed) Throws: - ErrUnauthorized - `from` does not have `burn` permission. - ErrInvalidRequest: - the balance of `from` does not have enough tokens to burn. | |
| `BurnNFTFrom` | [MsgBurnNFTFrom](#lbm.collection.v1.MsgBurnNFTFrom) | [MsgBurnNFTFromResponse](#lbm.collection.v1.MsgBurnNFTFromResponse) | BurnNFTFrom defines a method to burn non-fungible tokens of the approver by the proxy. Fires: - EventBurned - burn_ft_from (deprecated, not typed) - burn_nft_from (deprecated, not typed) - operation_burn_nft (deprecated, not typed) Throws: - ErrUnauthorized - `proxy` does not have `burn` permission. - the approver has not authorized `proxy`. - ErrInvalidRequest: - the balance of `from` does not have enough tokens to burn. | |
| `Modify` | [MsgModify](#lbm.collection.v1.MsgModify) | [MsgModifyResponse](#lbm.collection.v1.MsgModifyResponse) | Modify defines a method to modify metadata. Fires: - EventModifiedContract - modify_collection (deprecated, not typed) - EventModifiedTokenClass - modify_token_type (deprecated, not typed) - modify_token (deprecated, not typed) - EventModifiedNFT Throws: - ErrUnauthorized - the proxy does not have `modify` permission. - ErrNotFound - there is no contract of `contract_id`. - there is no token type of `token_type`. - there is no token of `token_id`. | |
| `GrantPermission` | [MsgGrantPermission](#lbm.collection.v1.MsgGrantPermission) | [MsgGrantPermissionResponse](#lbm.collection.v1.MsgGrantPermissionResponse) | GrantPermission allows one to mint or burn tokens or modify metadata. Fires: - EventGranted - grant_perm (deprecated, not typed) Throws: - ErrUnauthorized - `granter` does not have `permission`. - ErrInvalidRequest - `grantee` already has `permission`. | |
| `RevokePermission` | [MsgRevokePermission](#lbm.collection.v1.MsgRevokePermission) | [MsgRevokePermissionResponse](#lbm.collection.v1.MsgRevokePermissionResponse) | RevokePermission abandons a permission. Fires: - EventRenounced - revoke_perm (deprecated, not typed) Throws: - ErrUnauthorized - `grantee` does not have `permission`. | |
| `Attach` | [MsgAttach](#lbm.collection.v1.MsgAttach) | [MsgAttachResponse](#lbm.collection.v1.MsgAttachResponse) | Attach defines a method to attach a token to another token. Fires: - EventAttach - attach (deprecated, not typed) - operation_root_changed (deprecated, not typed) Throws: - ErrInvalidRequest - `owner` does not owns `id`. - `owner` does not owns `to`. - `token_id` is not root. - `token_id` is an ancestor of `to_token_id`, which creates a cycle as a result. - depth of `to_token_id` exceeds an app-specific limit. | |
| `Detach` | [MsgDetach](#lbm.collection.v1.MsgDetach) | [MsgDetachResponse](#lbm.collection.v1.MsgDetachResponse) | Detach defines a method to detach a token from another token. Fires: - EventDetach - detach (deprecated, not typed) - operation_root_changed (deprecated, not typed) Throws: - ErrInvalidRequest - `owner` does not owns `token_id`. | |
| `AttachFrom` | [MsgAttachFrom](#lbm.collection.v1.MsgAttachFrom) | [MsgAttachFromResponse](#lbm.collection.v1.MsgAttachFromResponse) | AttachFrom defines a method to attach a token to another token by proxy. Fires: - EventAttach - attach_from (deprecated, not typed) - operation_root_changed (deprecated, not typed) Throws: - ErrUnauthorized - the approver has not authorized `proxy`. - ErrInvalidRequest - `owner` does not owns `subject`. - `owner` does not owns `target`. - `subject` is not root. - `subject` is an ancestor of `target`, which creates a cycle as a result. - depth of `to` exceeds an app-specific limit. | |
| `DetachFrom` | [MsgDetachFrom](#lbm.collection.v1.MsgDetachFrom) | [MsgDetachFromResponse](#lbm.collection.v1.MsgDetachFromResponse) | DetachFrom defines a method to detach a token from another token by proxy. Fires: - EventDetach - detach_from (deprecated, not typed) - operation_root_changed (deprecated, not typed) Throws: - ErrUnauthorized - the approver has not authorized `proxy`. - ErrInvalidRequest - `owner` does not owns `subject`. | |

 <!-- end services -->



<a name="lbm/foundation/v1/authz.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/foundation/v1/authz.proto



<a name="lbm.foundation.v1.ReceiveFromTreasuryAuthorization"></a>

### ReceiveFromTreasuryAuthorization
ReceiveFromTreasuryAuthorization allows the grantee to receive coins
up to receive_limit from the treasury.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/foundation/v1/foundation.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/foundation/v1/foundation.proto



<a name="lbm.foundation.v1.DecisionPolicyWindows"></a>

### DecisionPolicyWindows
DecisionPolicyWindows defines the different windows for voting and execution.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `voting_period` | [google.protobuf.Duration](#google.protobuf.Duration) |  | voting_period is the duration from submission of a proposal to the end of voting period Within this times votes can be submitted with MsgVote. |
| `min_execution_period` | [google.protobuf.Duration](#google.protobuf.Duration) |  | min_execution_period is the minimum duration after the proposal submission where members can start sending MsgExec. This means that the window for sending a MsgExec transaction is: `[ submission + min_execution_period ; submission + voting_period + max_execution_period]` where max_execution_period is a app-specific config, defined in the keeper. If not set, min_execution_period will default to 0.

Please make sure to set a `min_execution_period` that is smaller than `voting_period + max_execution_period`, or else the above execution window is empty, meaning that all proposals created with this decision policy won't be able to be executed. |






<a name="lbm.foundation.v1.FoundationInfo"></a>

### FoundationInfo
FoundationInfo represents the high-level on-chain information for the foundation.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `version` | [uint64](#uint64) |  | version is used to track changes to the foundation's membership structure that would break existing proposals. Whenever any member is added or removed, this version is incremented and will cause proposals based on older versions of the foundation to fail |
| `total_weight` | [string](#string) |  | total_weight is the number of the foundation members. |
| `decision_policy` | [google.protobuf.Any](#google.protobuf.Any) |  | decision_policy specifies the foundation's decision policy. |






<a name="lbm.foundation.v1.Member"></a>

### Member
Member represents a foundation member with an account address and metadata.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the member's account address. |
| `metadata` | [string](#string) |  | metadata is any arbitrary metadata to attached to the member. |
| `added_at` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | added_at is a timestamp specifying when a member was added. |






<a name="lbm.foundation.v1.MemberRequest"></a>

### MemberRequest
MemberRequest represents a foundation member to be used in Msg server requests.
Contrary to `Member`, it doesn't have any `added_at` field
since this field cannot be set as part of requests.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the member's account address. |
| `remove` | [bool](#bool) |  | remove is the flag which allows one to remove the member by setting the flag to true. |
| `metadata` | [string](#string) |  | metadata is any arbitrary metadata attached to the member. |






<a name="lbm.foundation.v1.OutsourcingDecisionPolicy"></a>

### OutsourcingDecisionPolicy
OutsourcingDecisionPolicy is a dummy decision policy which is set after
the proposal feature has been outsourced to x/group.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `description` | [string](#string) |  |  |






<a name="lbm.foundation.v1.Params"></a>

### Params
Params defines the parameters for the foundation module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `foundation_tax` | [string](#string) |  |  |
| `censored_msg_type_urls` | [string](#string) | repeated |  |






<a name="lbm.foundation.v1.PercentageDecisionPolicy"></a>

### PercentageDecisionPolicy
PercentageDecisionPolicy is a decision policy where a proposal passes when
it satisfies the two following conditions:
1. The percentage of all `YES` voters' weights out of the total group weight
   is greater or equal than the given `percentage`.
2. The voting and execution periods of the proposal respect the parameters
   given by `windows`.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `percentage` | [string](#string) |  | percentage is the minimum percentage the sum of yes votes must meet for a proposal to succeed. |
| `windows` | [DecisionPolicyWindows](#lbm.foundation.v1.DecisionPolicyWindows) |  | windows defines the different windows for voting and execution. |






<a name="lbm.foundation.v1.Pool"></a>

### Pool
Pool is used for tracking treasury.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `treasury` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |






<a name="lbm.foundation.v1.Proposal"></a>

### Proposal
Proposal defines a foundation proposal. Any member of the foundation can submit a proposal
for a group policy to decide upon.
A proposal consists of a set of `sdk.Msg`s that will be executed if the proposal
passes as well as some optional metadata associated with the proposal.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | id is the unique id of the proposal. |
| `metadata` | [string](#string) |  | metadata is any arbitrary metadata to attached to the proposal. |
| `proposers` | [string](#string) | repeated | proposers are the account addresses of the proposers. |
| `submit_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | submit_time is a timestamp specifying when a proposal was submitted. |
| `foundation_version` | [uint64](#uint64) |  | foundation_version tracks the version of the foundation that this proposal corresponds to. When foundation info is changed, existing proposals from previous foundation versions will become invalid. |
| `status` | [ProposalStatus](#lbm.foundation.v1.ProposalStatus) |  | status represents the high level position in the life cycle of the proposal. Initial value is Submitted. |
| `final_tally_result` | [TallyResult](#lbm.foundation.v1.TallyResult) |  | final_tally_result contains the sums of all votes for this proposal for each vote option, after tallying. When querying a proposal via gRPC, this field is not populated until the proposal's voting period has ended. |
| `voting_period_end` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | voting_period_end is the timestamp before which voting must be done. Unless a successfull MsgExec is called before (to execute a proposal whose tally is successful before the voting period ends), tallying will be done at this point, and the `final_tally_result`, as well as `status` and `result` fields will be accordingly updated. |
| `executor_result` | [ProposalExecutorResult](#lbm.foundation.v1.ProposalExecutorResult) |  | executor_result is the final result based on the votes and election rule. Initial value is NotRun. |
| `messages` | [google.protobuf.Any](#google.protobuf.Any) | repeated | messages is a list of Msgs that will be executed if the proposal passes. |






<a name="lbm.foundation.v1.TallyResult"></a>

### TallyResult
TallyResult represents the sum of votes for each vote option.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `yes_count` | [string](#string) |  | yes_count is the sum of yes votes. |
| `abstain_count` | [string](#string) |  | abstain_count is the sum of abstainers. |
| `no_count` | [string](#string) |  | no is the sum of no votes. |
| `no_with_veto_count` | [string](#string) |  | no_with_veto_count is the sum of veto. |






<a name="lbm.foundation.v1.ThresholdDecisionPolicy"></a>

### ThresholdDecisionPolicy
ThresholdDecisionPolicy is a decision policy where a proposal passes when it
satisfies the two following conditions:
1. The sum of all `YES` voters' weights is greater or equal than the defined
   `threshold`.
2. The voting and execution periods of the proposal respect the parameters
   given by `windows`.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `threshold` | [string](#string) |  | threshold is the minimum sum of yes votes that must be met or exceeded for a proposal to succeed. |
| `windows` | [DecisionPolicyWindows](#lbm.foundation.v1.DecisionPolicyWindows) |  | windows defines the different windows for voting and execution. |






<a name="lbm.foundation.v1.Vote"></a>

### Vote
Vote represents a vote for a proposal.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal is the unique ID of the proposal. |
| `voter` | [string](#string) |  | voter is the account address of the voter. |
| `option` | [VoteOption](#lbm.foundation.v1.VoteOption) |  | option is the voter's choice on the proposal. |
| `metadata` | [string](#string) |  | metadata is any arbitrary metadata to attached to the vote. |
| `submit_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | submit_time is the timestamp when the vote was submitted. |





 <!-- end messages -->


<a name="lbm.foundation.v1.ProposalExecutorResult"></a>

### ProposalExecutorResult
ProposalExecutorResult defines types of proposal executor results.

| Name | Number | Description |
| ---- | ------ | ----------- |
| PROPOSAL_EXECUTOR_RESULT_UNSPECIFIED | 0 | An empty value is not allowed. |
| PROPOSAL_EXECUTOR_RESULT_NOT_RUN | 1 | We have not yet run the executor. |
| PROPOSAL_EXECUTOR_RESULT_SUCCESS | 2 | The executor was successful and proposed action updated state. |
| PROPOSAL_EXECUTOR_RESULT_FAILURE | 3 | The executor returned an error and proposed action didn't update state. |



<a name="lbm.foundation.v1.ProposalStatus"></a>

### ProposalStatus
ProposalStatus defines proposal statuses.

| Name | Number | Description |
| ---- | ------ | ----------- |
| PROPOSAL_STATUS_UNSPECIFIED | 0 | An empty value is invalid and not allowed. |
| PROPOSAL_STATUS_SUBMITTED | 1 | Initial status of a proposal when submitted. |
| PROPOSAL_STATUS_ACCEPTED | 2 | Final status of a proposal when the final tally is done and the outcome passes the foundation's decision policy. |
| PROPOSAL_STATUS_REJECTED | 3 | Final status of a proposal when the final tally is done and the outcome is rejected by the foundation's decision policy. |
| PROPOSAL_STATUS_ABORTED | 4 | Final status of a proposal when the decision policy is modified before the final tally. |
| PROPOSAL_STATUS_WITHDRAWN | 5 | A proposal can be withdrawn before the voting start time by the owner. When this happens the final status is Withdrawn. |



<a name="lbm.foundation.v1.VoteOption"></a>

### VoteOption
VoteOption enumerates the valid vote options for a given proposal.

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



<a name="lbm/foundation/v1/event.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/foundation/v1/event.proto



<a name="lbm.foundation.v1.EventExec"></a>

### EventExec
EventExec is an event emitted when a proposal is executed.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id is the unique ID of the proposal. |
| `result` | [ProposalExecutorResult](#lbm.foundation.v1.ProposalExecutorResult) |  | result is the proposal execution result. |
| `logs` | [string](#string) |  | logs contains error logs in case the execution result is FAILURE. |






<a name="lbm.foundation.v1.EventFundTreasury"></a>

### EventFundTreasury
EventFundTreasury is an event emitted when one funds the treasury.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `from` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="lbm.foundation.v1.EventGrant"></a>

### EventGrant
EventGrant is emitted on Msg/Grant


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grantee` | [string](#string) |  | the address of the grantee. |
| `authorization` | [google.protobuf.Any](#google.protobuf.Any) |  | authorization granted. |






<a name="lbm.foundation.v1.EventLeaveFoundation"></a>

### EventLeaveFoundation
EventLeaveFoundation is an event emitted when a foundation member leaves the foundation.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the account address of the foundation member. |






<a name="lbm.foundation.v1.EventRevoke"></a>

### EventRevoke
EventRevoke is emitted on Msg/Revoke


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grantee` | [string](#string) |  | address of the grantee. |
| `msg_type_url` | [string](#string) |  | message type url for which an autorization is revoked. |






<a name="lbm.foundation.v1.EventSubmitProposal"></a>

### EventSubmitProposal
EventSubmitProposal is an event emitted when a proposal is created.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal` | [Proposal](#lbm.foundation.v1.Proposal) |  | proposal is the unique ID of the proposal. |






<a name="lbm.foundation.v1.EventUpdateDecisionPolicy"></a>

### EventUpdateDecisionPolicy
EventUpdateDecisionPolicy is an event emitted when the decision policy have been updated.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `decision_policy` | [google.protobuf.Any](#google.protobuf.Any) |  |  |






<a name="lbm.foundation.v1.EventUpdateMembers"></a>

### EventUpdateMembers
EventUpdateMembers is an event emitted when the members have been updated.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `member_updates` | [MemberRequest](#lbm.foundation.v1.MemberRequest) | repeated |  |






<a name="lbm.foundation.v1.EventUpdateParams"></a>

### EventUpdateParams
EventUpdateParams is emitted after updating foundation parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#lbm.foundation.v1.Params) |  |  |






<a name="lbm.foundation.v1.EventVote"></a>

### EventVote
EventVote is an event emitted when a voter votes on a proposal.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vote` | [Vote](#lbm.foundation.v1.Vote) |  |  |






<a name="lbm.foundation.v1.EventWithdrawFromTreasury"></a>

### EventWithdrawFromTreasury
EventWithdrawFromTreasury is an event emitted when coins are withdrawn from the treasury.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="lbm.foundation.v1.EventWithdrawProposal"></a>

### EventWithdrawProposal
EventWithdrawProposal is an event emitted when a proposal is withdrawn.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id is the unique ID of the proposal. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/foundation/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/foundation/v1/genesis.proto



<a name="lbm.foundation.v1.GenesisState"></a>

### GenesisState
GenesisState defines the foundation module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#lbm.foundation.v1.Params) |  | params defines the module parameters at genesis. |
| `foundation` | [FoundationInfo](#lbm.foundation.v1.FoundationInfo) |  | foundation is the foundation info. |
| `members` | [Member](#lbm.foundation.v1.Member) | repeated | members is the list of the foundation members. |
| `previous_proposal_id` | [uint64](#uint64) |  | it is used to get the next proposal ID. |
| `proposals` | [Proposal](#lbm.foundation.v1.Proposal) | repeated | proposals is the list of proposals. |
| `votes` | [Vote](#lbm.foundation.v1.Vote) | repeated | votes is the list of votes. |
| `authorizations` | [GrantAuthorization](#lbm.foundation.v1.GrantAuthorization) | repeated | grants |
| `pool` | [Pool](#lbm.foundation.v1.Pool) |  | pool |






<a name="lbm.foundation.v1.GrantAuthorization"></a>

### GrantAuthorization
GrantAuthorization defines authorization grant to grantee via route.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grantee` | [string](#string) |  |  |
| `authorization` | [google.protobuf.Any](#google.protobuf.Any) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/foundation/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/foundation/v1/query.proto



<a name="lbm.foundation.v1.QueryFoundationInfoRequest"></a>

### QueryFoundationInfoRequest
QueryFoundationInfoRequest is the Query/FoundationInfo request type.






<a name="lbm.foundation.v1.QueryFoundationInfoResponse"></a>

### QueryFoundationInfoResponse
QueryFoundationInfoResponse is the Query/FoundationInfo response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `info` | [FoundationInfo](#lbm.foundation.v1.FoundationInfo) |  | info is the FoundationInfo for the foundation. |






<a name="lbm.foundation.v1.QueryGrantsRequest"></a>

### QueryGrantsRequest
QueryGrantsRequest is the request type for the Query/Grants RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grantee` | [string](#string) |  |  |
| `msg_type_url` | [string](#string) |  | Optional, msg_type_url, when set, will query only grants matching given msg type. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.foundation.v1.QueryGrantsResponse"></a>

### QueryGrantsResponse
QueryGrantsResponse is the response type for the Query/Grants RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authorizations` | [google.protobuf.Any](#google.protobuf.Any) | repeated | authorizations is a list of grants granted for grantee. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.foundation.v1.QueryMemberRequest"></a>

### QueryMemberRequest
QueryMemberRequest is the Query/Member request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |






<a name="lbm.foundation.v1.QueryMemberResponse"></a>

### QueryMemberResponse
QueryMemberResponse is the Query/MemberResponse response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `member` | [Member](#lbm.foundation.v1.Member) |  | member is the members of the foundation. |






<a name="lbm.foundation.v1.QueryMembersRequest"></a>

### QueryMembersRequest
QueryMembersRequest is the Query/Members request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.foundation.v1.QueryMembersResponse"></a>

### QueryMembersResponse
QueryMembersResponse is the Query/MembersResponse response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `members` | [Member](#lbm.foundation.v1.Member) | repeated | members are the members of the foundation. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.foundation.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="lbm.foundation.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#lbm.foundation.v1.Params) |  |  |






<a name="lbm.foundation.v1.QueryProposalRequest"></a>

### QueryProposalRequest
QueryProposalRequest is the Query/Proposal request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id is the unique ID of a proposal. |






<a name="lbm.foundation.v1.QueryProposalResponse"></a>

### QueryProposalResponse
QueryProposalResponse is the Query/Proposal response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal` | [Proposal](#lbm.foundation.v1.Proposal) |  | proposal is the proposal info. |






<a name="lbm.foundation.v1.QueryProposalsRequest"></a>

### QueryProposalsRequest
QueryProposals is the Query/Proposals request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.foundation.v1.QueryProposalsResponse"></a>

### QueryProposalsResponse
QueryProposalsResponse is the Query/Proposals response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposals` | [Proposal](#lbm.foundation.v1.Proposal) | repeated | proposals are the proposals of the foundation. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.foundation.v1.QueryTallyResultRequest"></a>

### QueryTallyResultRequest
QueryTallyResultRequest is the Query/TallyResult request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id is the unique id of a proposal. |






<a name="lbm.foundation.v1.QueryTallyResultResponse"></a>

### QueryTallyResultResponse
QueryTallyResultResponse is the Query/TallyResult response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tally` | [TallyResult](#lbm.foundation.v1.TallyResult) |  | tally defines the requested tally. |






<a name="lbm.foundation.v1.QueryTreasuryRequest"></a>

### QueryTreasuryRequest
QueryTreasuryRequest is the request type for the
Query/Treasury RPC method.






<a name="lbm.foundation.v1.QueryTreasuryResponse"></a>

### QueryTreasuryResponse
QueryTreasuryResponse is the response type for the
Query/Treasury RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `amount` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |






<a name="lbm.foundation.v1.QueryVoteRequest"></a>

### QueryVoteRequest
QueryVote is the Query/Vote request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id is the unique ID of a proposal. |
| `voter` | [string](#string) |  | voter is a proposal voter account address. |






<a name="lbm.foundation.v1.QueryVoteResponse"></a>

### QueryVoteResponse
QueryVoteResponse is the Query/Vote response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vote` | [Vote](#lbm.foundation.v1.Vote) |  | vote is the vote with given proposal_id and voter. |






<a name="lbm.foundation.v1.QueryVotesRequest"></a>

### QueryVotesRequest
QueryVotes is the Query/Votes request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id is the unique ID of a proposal. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.foundation.v1.QueryVotesResponse"></a>

### QueryVotesResponse
QueryVotesResponse is the Query/Votes response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `votes` | [Vote](#lbm.foundation.v1.Vote) | repeated | votes are the list of votes for given proposal_id. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.foundation.v1.Query"></a>

### Query
Query defines the gRPC querier service for foundation module.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#lbm.foundation.v1.QueryParamsRequest) | [QueryParamsResponse](#lbm.foundation.v1.QueryParamsResponse) | Params queries the module params. | GET|/lbm/foundation/v1/params|
| `Treasury` | [QueryTreasuryRequest](#lbm.foundation.v1.QueryTreasuryRequest) | [QueryTreasuryResponse](#lbm.foundation.v1.QueryTreasuryResponse) | Treasury queries the foundation treasury. | GET|/lbm/foundation/v1/treasury|
| `FoundationInfo` | [QueryFoundationInfoRequest](#lbm.foundation.v1.QueryFoundationInfoRequest) | [QueryFoundationInfoResponse](#lbm.foundation.v1.QueryFoundationInfoResponse) | FoundationInfo queries foundation info. | GET|/lbm/foundation/v1/foundation_info|
| `Member` | [QueryMemberRequest](#lbm.foundation.v1.QueryMemberRequest) | [QueryMemberResponse](#lbm.foundation.v1.QueryMemberResponse) | Member queries a member of the foundation | GET|/lbm/foundation/v1/foundation_members/{address}|
| `Members` | [QueryMembersRequest](#lbm.foundation.v1.QueryMembersRequest) | [QueryMembersResponse](#lbm.foundation.v1.QueryMembersResponse) | Members queries members of the foundation | GET|/lbm/foundation/v1/foundation_members|
| `Proposal` | [QueryProposalRequest](#lbm.foundation.v1.QueryProposalRequest) | [QueryProposalResponse](#lbm.foundation.v1.QueryProposalResponse) | Proposal queries a proposal based on proposal id. | GET|/lbm/foundation/v1/proposals/{proposal_id}|
| `Proposals` | [QueryProposalsRequest](#lbm.foundation.v1.QueryProposalsRequest) | [QueryProposalsResponse](#lbm.foundation.v1.QueryProposalsResponse) | Proposals queries all proposals. | GET|/lbm/foundation/v1/proposals|
| `Vote` | [QueryVoteRequest](#lbm.foundation.v1.QueryVoteRequest) | [QueryVoteResponse](#lbm.foundation.v1.QueryVoteResponse) | Vote queries a vote by proposal id and voter. | GET|/lbm/foundation/v1/proposals/{proposal_id}/votes/{voter}|
| `Votes` | [QueryVotesRequest](#lbm.foundation.v1.QueryVotesRequest) | [QueryVotesResponse](#lbm.foundation.v1.QueryVotesResponse) | Votes queries a vote by proposal. | GET|/lbm/foundation/v1/proposals/{proposal_id}/votes|
| `TallyResult` | [QueryTallyResultRequest](#lbm.foundation.v1.QueryTallyResultRequest) | [QueryTallyResultResponse](#lbm.foundation.v1.QueryTallyResultResponse) | TallyResult queries the tally of a proposal votes. | GET|/lbm/foundation/v1/proposals/{proposal_id}/tally|
| `Grants` | [QueryGrantsRequest](#lbm.foundation.v1.QueryGrantsRequest) | [QueryGrantsResponse](#lbm.foundation.v1.QueryGrantsResponse) | Returns list of authorizations, granted to the grantee. | GET|/lbm/foundation/v1/grants/{grantee}/{msg_type_url}|

 <!-- end services -->



<a name="lbm/foundation/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/foundation/v1/tx.proto



<a name="lbm.foundation.v1.MsgExec"></a>

### MsgExec
MsgExec is the Msg/Exec request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal is the unique ID of the proposal. |
| `signer` | [string](#string) |  | signer is the account address used to execute the proposal. |






<a name="lbm.foundation.v1.MsgExecResponse"></a>

### MsgExecResponse
MsgExecResponse is the Msg/Exec request type.






<a name="lbm.foundation.v1.MsgFundTreasury"></a>

### MsgFundTreasury
MsgFundTreasury is the Msg/FundTreasury request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `from` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="lbm.foundation.v1.MsgFundTreasuryResponse"></a>

### MsgFundTreasuryResponse
MsgFundTreasuryResponse is the Msg/FundTreasury response type.






<a name="lbm.foundation.v1.MsgGrant"></a>

### MsgGrant
MsgGrant is the Msg/Grant request type.
on behalf of the foundation.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  | authority is the address of the privileged account. |
| `grantee` | [string](#string) |  |  |
| `authorization` | [google.protobuf.Any](#google.protobuf.Any) |  |  |






<a name="lbm.foundation.v1.MsgGrantResponse"></a>

### MsgGrantResponse
MsgGrantResponse is the Msg/MsgGrant response type.






<a name="lbm.foundation.v1.MsgLeaveFoundation"></a>

### MsgLeaveFoundation
MsgLeaveFoundation is the Msg/LeaveFoundation request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the account address of the foundation member. |






<a name="lbm.foundation.v1.MsgLeaveFoundationResponse"></a>

### MsgLeaveFoundationResponse
MsgLeaveFoundationResponse is the Msg/LeaveFoundation response type.






<a name="lbm.foundation.v1.MsgRevoke"></a>

### MsgRevoke
MsgRevoke is the Msg/Revoke request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  | authority is the address of the privileged account. |
| `grantee` | [string](#string) |  |  |
| `msg_type_url` | [string](#string) |  |  |






<a name="lbm.foundation.v1.MsgRevokeResponse"></a>

### MsgRevokeResponse
MsgRevokeResponse is the Msg/MsgRevokeResponse response type.






<a name="lbm.foundation.v1.MsgSubmitProposal"></a>

### MsgSubmitProposal
MsgSubmitProposal is the Msg/SubmitProposal request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposers` | [string](#string) | repeated | proposers are the account addresses of the proposers. Proposers signatures will be counted as yes votes. |
| `metadata` | [string](#string) |  | metadata is any arbitrary metadata to attached to the proposal. |
| `messages` | [google.protobuf.Any](#google.protobuf.Any) | repeated | messages is a list of `sdk.Msg`s that will be executed if the proposal passes. |
| `exec` | [Exec](#lbm.foundation.v1.Exec) |  | exec defines the mode of execution of the proposal, whether it should be executed immediately on creation or not. If so, proposers signatures are considered as Yes votes. |






<a name="lbm.foundation.v1.MsgSubmitProposalResponse"></a>

### MsgSubmitProposalResponse
MsgSubmitProposalResponse is the Msg/SubmitProposal response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal is the unique ID of the proposal. |






<a name="lbm.foundation.v1.MsgUpdateDecisionPolicy"></a>

### MsgUpdateDecisionPolicy
MsgUpdateDecisionPolicy is the Msg/UpdateDecisionPolicy request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  | authority is the address of the privileged account. |
| `decision_policy` | [google.protobuf.Any](#google.protobuf.Any) |  | decision_policy is the updated decision policy. |






<a name="lbm.foundation.v1.MsgUpdateDecisionPolicyResponse"></a>

### MsgUpdateDecisionPolicyResponse
MsgUpdateDecisionPolicyResponse is the Msg/UpdateDecisionPolicy response type.






<a name="lbm.foundation.v1.MsgUpdateMembers"></a>

### MsgUpdateMembers
MsgUpdateMembers is the Msg/UpdateMembers request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  | authority is the address of the privileged account. |
| `member_updates` | [MemberRequest](#lbm.foundation.v1.MemberRequest) | repeated | member_updates is the list of members to update, set remove to true to remove a member. |






<a name="lbm.foundation.v1.MsgUpdateMembersResponse"></a>

### MsgUpdateMembersResponse
MsgUpdateMembersResponse is the Msg/UpdateMembers response type.






<a name="lbm.foundation.v1.MsgUpdateParams"></a>

### MsgUpdateParams
MsgUpdateParams is the Msg/UpdateParams request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  | authority is the address of the privileged account. |
| `params` | [Params](#lbm.foundation.v1.Params) |  | params defines the x/foundation parameters to update.

NOTE: All parameters must be supplied. |






<a name="lbm.foundation.v1.MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse
MsgUpdateParamsResponse is the Msg/UpdateParams response type.






<a name="lbm.foundation.v1.MsgVote"></a>

### MsgVote
MsgVote is the Msg/Vote request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal is the unique ID of the proposal. |
| `voter` | [string](#string) |  | voter is the voter account address. |
| `option` | [VoteOption](#lbm.foundation.v1.VoteOption) |  | option is the voter's choice on the proposal. |
| `metadata` | [string](#string) |  | metadata is any arbitrary metadata to attached to the vote. |
| `exec` | [Exec](#lbm.foundation.v1.Exec) |  | exec defines whether the proposal should be executed immediately after voting or not. |






<a name="lbm.foundation.v1.MsgVoteResponse"></a>

### MsgVoteResponse
MsgVoteResponse is the Msg/Vote response type.






<a name="lbm.foundation.v1.MsgWithdrawFromTreasury"></a>

### MsgWithdrawFromTreasury
MsgWithdrawFromTreasury is the Msg/WithdrawFromTreasury request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  | authority is the address of the privileged account. |
| `to` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="lbm.foundation.v1.MsgWithdrawFromTreasuryResponse"></a>

### MsgWithdrawFromTreasuryResponse
MsgWithdrawFromTreasuryResponse is the Msg/WithdrawFromTreasury response type.






<a name="lbm.foundation.v1.MsgWithdrawProposal"></a>

### MsgWithdrawProposal
MsgWithdrawProposal is the Msg/WithdrawProposal request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal is the unique ID of the proposal. |
| `address` | [string](#string) |  | address of one of the proposer of the proposal. |






<a name="lbm.foundation.v1.MsgWithdrawProposalResponse"></a>

### MsgWithdrawProposalResponse
MsgWithdrawProposalResponse is the Msg/WithdrawProposal response type.





 <!-- end messages -->


<a name="lbm.foundation.v1.Exec"></a>

### Exec
Exec defines modes of execution of a proposal on creation or on new vote.

| Name | Number | Description |
| ---- | ------ | ----------- |
| EXEC_UNSPECIFIED | 0 | An empty value means that there should be a separate MsgExec request for the proposal to execute. |
| EXEC_TRY | 1 | Try to execute the proposal immediately. If the proposal is not allowed per the DecisionPolicy, the proposal will still be open and could be executed at a later point. |


 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.foundation.v1.Msg"></a>

### Msg
Msg defines the foundation Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `UpdateParams` | [MsgUpdateParams](#lbm.foundation.v1.MsgUpdateParams) | [MsgUpdateParamsResponse](#lbm.foundation.v1.MsgUpdateParamsResponse) | UpdateParams defines an operation for updating the x/foundation module parameters. | |
| `FundTreasury` | [MsgFundTreasury](#lbm.foundation.v1.MsgFundTreasury) | [MsgFundTreasuryResponse](#lbm.foundation.v1.MsgFundTreasuryResponse) | FundTreasury defines a method to fund the treasury. | |
| `WithdrawFromTreasury` | [MsgWithdrawFromTreasury](#lbm.foundation.v1.MsgWithdrawFromTreasury) | [MsgWithdrawFromTreasuryResponse](#lbm.foundation.v1.MsgWithdrawFromTreasuryResponse) | WithdrawFromTreasury defines a method to withdraw coins from the treasury. | |
| `UpdateMembers` | [MsgUpdateMembers](#lbm.foundation.v1.MsgUpdateMembers) | [MsgUpdateMembersResponse](#lbm.foundation.v1.MsgUpdateMembersResponse) | UpdateMembers updates the foundation members. | |
| `UpdateDecisionPolicy` | [MsgUpdateDecisionPolicy](#lbm.foundation.v1.MsgUpdateDecisionPolicy) | [MsgUpdateDecisionPolicyResponse](#lbm.foundation.v1.MsgUpdateDecisionPolicyResponse) | UpdateDecisionPolicy allows a group policy's decision policy to be updated. | |
| `SubmitProposal` | [MsgSubmitProposal](#lbm.foundation.v1.MsgSubmitProposal) | [MsgSubmitProposalResponse](#lbm.foundation.v1.MsgSubmitProposalResponse) | SubmitProposal submits a new proposal. | |
| `WithdrawProposal` | [MsgWithdrawProposal](#lbm.foundation.v1.MsgWithdrawProposal) | [MsgWithdrawProposalResponse](#lbm.foundation.v1.MsgWithdrawProposalResponse) | WithdrawProposal aborts a proposal. | |
| `Vote` | [MsgVote](#lbm.foundation.v1.MsgVote) | [MsgVoteResponse](#lbm.foundation.v1.MsgVoteResponse) | Vote allows a voter to vote on a proposal. | |
| `Exec` | [MsgExec](#lbm.foundation.v1.MsgExec) | [MsgExecResponse](#lbm.foundation.v1.MsgExecResponse) | Exec executes a proposal. | |
| `LeaveFoundation` | [MsgLeaveFoundation](#lbm.foundation.v1.MsgLeaveFoundation) | [MsgLeaveFoundationResponse](#lbm.foundation.v1.MsgLeaveFoundationResponse) | LeaveFoundation allows a member to leave the foundation. | |
| `Grant` | [MsgGrant](#lbm.foundation.v1.MsgGrant) | [MsgGrantResponse](#lbm.foundation.v1.MsgGrantResponse) | Grant grants the provided authorization to the grantee with authority of the foundation. If there is already a grant for the given (grantee, Authorization) tuple, then the grant will be overwritten. | |
| `Revoke` | [MsgRevoke](#lbm.foundation.v1.MsgRevoke) | [MsgRevokeResponse](#lbm.foundation.v1.MsgRevokeResponse) | Revoke revokes any authorization corresponding to the provided method name that has been granted to the grantee. | |

 <!-- end services -->



<a name="lbm/stakingplus/v1/authz.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/stakingplus/v1/authz.proto



<a name="lbm.stakingplus.v1.CreateValidatorAuthorization"></a>

### CreateValidatorAuthorization
CreateValidatorAuthorization allows the grantee to create a new validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | redundant, but good for the query. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/token/v1/token.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/token/v1/token.proto



<a name="lbm.token.v1.Authorization"></a>

### Authorization
Authorization defines an authorization given to the operator on tokens of the holder.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `holder` | [string](#string) |  | address of the token holder which approves the authorization. |
| `operator` | [string](#string) |  | address of the operator which the authorization is granted to. |






<a name="lbm.token.v1.Contract"></a>

### Contract
Contract defines token information.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  | id defines the unique identifier of the contract. |
| `name` | [string](#string) |  | name defines the human-readable name of the contract. mandatory (not ERC20 compliant). |
| `symbol` | [string](#string) |  | symbol is an abbreviated name for contract. mandatory (not ERC20 compliant). |
| `image_uri` | [string](#string) |  | image_uri is an uri for the image of the contract stored off chain. |
| `meta` | [string](#string) |  | meta is a brief description of contract. |
| `decimals` | [int32](#int32) |  | decimals is the number of decimals which one must divide the amount by to get its user representation. |
| `mintable` | [bool](#bool) |  | mintable represents whether the token is allowed to mint or burn. |






<a name="lbm.token.v1.Grant"></a>

### Grant
Grant defines permission given to a grantee.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grantee` | [string](#string) |  | address of the grantee. |
| `permission` | [Permission](#lbm.token.v1.Permission) |  | permission on the contract. |






<a name="lbm.token.v1.Pair"></a>

### Pair
Pair defines a key-value pair.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `field` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |






<a name="lbm.token.v1.Params"></a>

### Params
Params defines the parameters for the token module.





 <!-- end messages -->


<a name="lbm.token.v1.LegacyPermission"></a>

### LegacyPermission
Deprecated: use Permission

LegacyPermission enumerates the valid permissions on a contract.

| Name | Number | Description |
| ---- | ------ | ----------- |
| LEGACY_PERMISSION_UNSPECIFIED | 0 | unspecified defines the default permission which is invalid. |
| LEGACY_PERMISSION_MODIFY | 1 | modify defines a permission to modify a contract. |
| LEGACY_PERMISSION_MINT | 2 | mint defines a permission to mint tokens of a contract. |
| LEGACY_PERMISSION_BURN | 3 | burn defines a permission to burn tokens of a contract. |



<a name="lbm.token.v1.Permission"></a>

### Permission
Permission enumerates the valid permissions on a contract.

| Name | Number | Description |
| ---- | ------ | ----------- |
| PERMISSION_UNSPECIFIED | 0 | unspecified defines the default permission which is invalid. |
| PERMISSION_MODIFY | 1 | PERMISSION_MODIFY defines a permission to modify a contract. |
| PERMISSION_MINT | 2 | PERMISSION_MINT defines a permission to mint tokens of a contract. |
| PERMISSION_BURN | 3 | PERMISSION_BURN defines a permission to burn tokens of a contract. |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/token/v1/event.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/token/v1/event.proto



<a name="lbm.token.v1.EventAuthorizedOperator"></a>

### EventAuthorizedOperator
EventAuthorizedOperator is emitted when a holder authorizes an operator to manipulate its tokens.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `holder` | [string](#string) |  | address of a holder which authorized the `operator` address as an operator. |
| `operator` | [string](#string) |  | address which became an operator of `holder`. |






<a name="lbm.token.v1.EventBurned"></a>

### EventBurned
EventBurned is emitted when tokens are burnt.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `operator` | [string](#string) |  | address which triggered the burn. |
| `from` | [string](#string) |  | holder whose tokens were burned. |
| `amount` | [string](#string) |  | number of tokens burned. |






<a name="lbm.token.v1.EventGranted"></a>

### EventGranted
EventGranted is emitted when a granter grants its permission to a grantee.

Info: `granter` would be empty if the permission is granted by an issuance.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `granter` | [string](#string) |  | address which granted the permission to `grantee`. it would be empty where the event is triggered by the issuance. |
| `grantee` | [string](#string) |  | address of the grantee. |
| `permission` | [Permission](#lbm.token.v1.Permission) |  | permission on the contract. |






<a name="lbm.token.v1.EventIssued"></a>

### EventIssued
EventIssued is emitted when a new contract is created.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `creator` | [string](#string) |  | address which created the contract. |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `name` | [string](#string) |  | name defines the human-readable name of the contract. |
| `symbol` | [string](#string) |  | symbol is an abbreviated name for contract. |
| `uri` | [string](#string) |  | uri is an uri for the resource of the contract stored off chain. |
| `meta` | [string](#string) |  | meta is a brief description of contract. |
| `decimals` | [int32](#int32) |  | decimals is the number of decimals which one must divide the amount by to get its user representation. |
| `mintable` | [bool](#bool) |  | mintable represents whether the token is allowed to mint. |






<a name="lbm.token.v1.EventMinted"></a>

### EventMinted
EventMinted is emitted when tokens are minted.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `operator` | [string](#string) |  | address which triggered the mint. |
| `to` | [string](#string) |  | recipient of the tokens. |
| `amount` | [string](#string) |  | number of tokens minted. |






<a name="lbm.token.v1.EventModified"></a>

### EventModified
EventModified is emitted when the information of a contract is modified.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `operator` | [string](#string) |  | address which triggered the modify. |
| `changes` | [Pair](#lbm.token.v1.Pair) | repeated | changes on the metadata of the class. |






<a name="lbm.token.v1.EventRenounced"></a>

### EventRenounced
EventRenounced is emitted when a grantee renounces its permission.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `grantee` | [string](#string) |  | address of the grantee which abandons its grant. |
| `permission` | [Permission](#lbm.token.v1.Permission) |  | permission on the contract. |






<a name="lbm.token.v1.EventRevokedOperator"></a>

### EventRevokedOperator
EventRevokedOperator is emitted when an authorization is revoked.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `holder` | [string](#string) |  | address of a holder which revoked the `operator` address as an operator. |
| `operator` | [string](#string) |  | address which was revoked as an operator of `holder`. |






<a name="lbm.token.v1.EventSent"></a>

### EventSent
EventSent is emitted when tokens are transferred.

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `operator` | [string](#string) |  | address which triggered the send. |
| `from` | [string](#string) |  | holder whose tokens were sent. |
| `to` | [string](#string) |  | recipient of the tokens |
| `amount` | [string](#string) |  | number of tokens sent. |





 <!-- end messages -->


<a name="lbm.token.v1.AttributeKey"></a>

### AttributeKey
AttributeKey enumerates the valid attribute keys on x/token.

| Name | Number | Description |
| ---- | ------ | ----------- |
| ATTRIBUTE_KEY_UNSPECIFIED | 0 |  |
| ATTRIBUTE_KEY_NAME | 1 |  |
| ATTRIBUTE_KEY_SYMBOL | 2 |  |
| ATTRIBUTE_KEY_META | 3 |  |
| ATTRIBUTE_KEY_CONTRACT_ID | 4 |  |
| ATTRIBUTE_KEY_OWNER | 5 |  |
| ATTRIBUTE_KEY_AMOUNT | 6 |  |
| ATTRIBUTE_KEY_DECIMALS | 7 |  |
| ATTRIBUTE_KEY_IMG_URI | 8 |  |
| ATTRIBUTE_KEY_MINTABLE | 9 |  |
| ATTRIBUTE_KEY_FROM | 10 |  |
| ATTRIBUTE_KEY_TO | 11 |  |
| ATTRIBUTE_KEY_PERM | 12 |  |
| ATTRIBUTE_KEY_APPROVER | 13 |  |
| ATTRIBUTE_KEY_PROXY | 14 |  |



<a name="lbm.token.v1.EventType"></a>

### EventType
Deprecated: use typed events.

EventType enumerates the valid event types on x/token.

| Name | Number | Description |
| ---- | ------ | ----------- |
| EVENT_TYPE_UNSPECIFIED | 0 |  |
| EVENT_TYPE_ISSUE | 1 |  |
| EVENT_TYPE_MINT | 2 |  |
| EVENT_TYPE_BURN | 3 |  |
| EVENT_TYPE_BURN_FROM | 4 |  |
| EVENT_TYPE_MODIFY_TOKEN | 5 |  |
| EVENT_TYPE_TRANSFER | 6 |  |
| EVENT_TYPE_TRANSFER_FROM | 7 |  |
| EVENT_TYPE_GRANT_PERM | 8 |  |
| EVENT_TYPE_REVOKE_PERM | 9 |  |
| EVENT_TYPE_APPROVE_TOKEN | 10 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/token/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/token/v1/genesis.proto



<a name="lbm.token.v1.Balance"></a>

### Balance
Balance defines a balance of an address.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address of the holder. |
| `amount` | [string](#string) |  | amount of the balance. |






<a name="lbm.token.v1.ClassGenesisState"></a>

### ClassGenesisState
ClassGenesisState defines the classs keeper's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [string](#string) |  | nonce is the next class nonce to issue. |
| `ids` | [string](#string) | repeated | ids represents the issued ids. |






<a name="lbm.token.v1.ContractAuthorizations"></a>

### ContractAuthorizations
ContractAuthorizations defines authorizations belong to a contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the token class. |
| `authorizations` | [Authorization](#lbm.token.v1.Authorization) | repeated | authorizations of the contract. |






<a name="lbm.token.v1.ContractBalances"></a>

### ContractBalances
ContractBalances defines balances belong to a contract.
genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the token class. |
| `balances` | [Balance](#lbm.token.v1.Balance) | repeated | balances of the contract. |






<a name="lbm.token.v1.ContractCoin"></a>

### ContractCoin



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the token class. |
| `amount` | [string](#string) |  | amount of the token. |






<a name="lbm.token.v1.ContractGrants"></a>

### ContractGrants
ContractGrant defines grants belong to a contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the token class. |
| `grants` | [Grant](#lbm.token.v1.Grant) | repeated | grants of the contract. |






<a name="lbm.token.v1.GenesisState"></a>

### GenesisState
GenesisState defines the token module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#lbm.token.v1.Params) |  | params defines all the paramaters of the module. |
| `class_state` | [ClassGenesisState](#lbm.token.v1.ClassGenesisState) |  | class_state is the class keeper's genesis state. |
| `balances` | [ContractBalances](#lbm.token.v1.ContractBalances) | repeated | balances is an array containing the balances of all the accounts. |
| `classes` | [Contract](#lbm.token.v1.Contract) | repeated | classes defines the metadata of the differents tokens. |
| `grants` | [ContractGrants](#lbm.token.v1.ContractGrants) | repeated | grants defines the grant information. |
| `authorizations` | [ContractAuthorizations](#lbm.token.v1.ContractAuthorizations) | repeated | authorizations defines the approve information. |
| `supplies` | [ContractCoin](#lbm.token.v1.ContractCoin) | repeated | supplies represents the total supplies of tokens. |
| `mints` | [ContractCoin](#lbm.token.v1.ContractCoin) | repeated | mints represents the total mints of tokens. |
| `burns` | [ContractCoin](#lbm.token.v1.ContractCoin) | repeated | burns represents the total burns of tokens. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="lbm/token/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/token/v1/query.proto



<a name="lbm.token.v1.QueryApprovedRequest"></a>

### QueryApprovedRequest
QueryApprovedRequest is the request type for the Query/Approved RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `proxy` | [string](#string) |  | address of the proxy which the authorization is granted to. |
| `approver` | [string](#string) |  | approver is the address of the approver of the authorization. |






<a name="lbm.token.v1.QueryApprovedResponse"></a>

### QueryApprovedResponse
QueryApprovedResponse is the response type for the Query/Approved RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `approved` | [bool](#bool) |  |  |






<a name="lbm.token.v1.QueryApproversRequest"></a>

### QueryApproversRequest
QueryApproversRequest is the request type for the Query/Approvers RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `address` | [string](#string) |  | address of the proxy which the authorization is granted to. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.token.v1.QueryApproversResponse"></a>

### QueryApproversResponse
QueryApproversResponse is the response type for the Query/Approvers RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `approvers` | [string](#string) | repeated | all the authorizations on the proxy. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.token.v1.QueryBalanceRequest"></a>

### QueryBalanceRequest
QueryBalanceRequest is the request type for the Query/Balance RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `address` | [string](#string) |  | address is the address to query balance for. |






<a name="lbm.token.v1.QueryBalanceResponse"></a>

### QueryBalanceResponse
QueryBalanceResponse is the response type for the Query/Balance RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `amount` | [string](#string) |  | the balance of the tokens. |






<a name="lbm.token.v1.QueryBurntRequest"></a>

### QueryBurntRequest
QueryBurntRequest is the request type for the Query/Burnt RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |






<a name="lbm.token.v1.QueryBurntResponse"></a>

### QueryBurntResponse
QueryBurntResponse is the response type for the Query/Burnt RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `amount` | [string](#string) |  | the amount of the burnt tokens. |






<a name="lbm.token.v1.QueryContractRequest"></a>

### QueryContractRequest
QueryContractRequest is the request type for the Query/Contract RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |






<a name="lbm.token.v1.QueryContractResponse"></a>

### QueryContractResponse
QueryContractResponse is the response type for the Query/Contract RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract` | [Contract](#lbm.token.v1.Contract) |  |  |






<a name="lbm.token.v1.QueryContractsRequest"></a>

### QueryContractsRequest
QueryContractsRequest is the request type for the Query/Contracts RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.token.v1.QueryContractsResponse"></a>

### QueryContractsResponse
QueryContractsResponse is the response type for the Query/Contracts RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contracts` | [Contract](#lbm.token.v1.Contract) | repeated | information of the contracts. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.token.v1.QueryGranteeGrantsRequest"></a>

### QueryGranteeGrantsRequest
QueryGranteeGrantsRequest is the request type for the Query/GranteeGrants RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `grantee` | [string](#string) |  | grantee which has permissions on the contract. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="lbm.token.v1.QueryGranteeGrantsResponse"></a>

### QueryGranteeGrantsResponse
QueryGranteeGrantsResponse is the response type for the Query/GranteeGrants RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grants` | [Grant](#lbm.token.v1.Grant) | repeated | all the grants on the grantee. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="lbm.token.v1.QueryMintedRequest"></a>

### QueryMintedRequest
QueryMintedRequest is the request type for the Query/Minted RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |






<a name="lbm.token.v1.QueryMintedResponse"></a>

### QueryMintedResponse
QueryMintedResponse is the response type for the Query/Minted RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `amount` | [string](#string) |  | the amount of the minted tokens. |






<a name="lbm.token.v1.QuerySupplyRequest"></a>

### QuerySupplyRequest
QuerySupplyRequest is the request type for the Query/Supply RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |






<a name="lbm.token.v1.QuerySupplyResponse"></a>

### QuerySupplyResponse
QuerySupplyResponse is the response type for the Query/Supply RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `amount` | [string](#string) |  | the supply of the tokens. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.token.v1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Balance` | [QueryBalanceRequest](#lbm.token.v1.QueryBalanceRequest) | [QueryBalanceResponse](#lbm.token.v1.QueryBalanceResponse) | Balance queries the number of tokens of a given contract owned by the address. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - ErrInvalidAddress - `address` is of invalid format. | GET|/lbm/token/v1/token_classes/{contract_id}/balances/{address}|
| `Supply` | [QuerySupplyRequest](#lbm.token.v1.QuerySupplyRequest) | [QuerySupplyResponse](#lbm.token.v1.QuerySupplyResponse) | Supply queries the number of tokens from the given contract id. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - ErrNotFound - there is no token class of `contract_id`. | GET|/lbm/token/v1/token_classes/{contract_id}/supply|
| `Minted` | [QueryMintedRequest](#lbm.token.v1.QueryMintedRequest) | [QueryMintedResponse](#lbm.token.v1.QueryMintedResponse) | Minted queries the number of minted tokens from the given contract id. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - ErrNotFound - there is no token class of `contract_id`. | GET|/lbm/token/v1/token_classes/{contract_id}/minted|
| `Burnt` | [QueryBurntRequest](#lbm.token.v1.QueryBurntRequest) | [QueryBurntResponse](#lbm.token.v1.QueryBurntResponse) | Burnt queries the number of burnt tokens from the given contract id. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - ErrNotFound - there is no token class of `contract_id`. | GET|/lbm/token/v1/token_classes/{contract_id}/burnt|
| `Contract` | [QueryContractRequest](#lbm.token.v1.QueryContractRequest) | [QueryContractResponse](#lbm.token.v1.QueryContractResponse) | Contract queries an token metadata based on its contract id. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - ErrNotFound - there is no token class of `contract_id`. | GET|/lbm/token/v1/token_classes/{contract_id}|
| `Contracts` | [QueryContractsRequest](#lbm.token.v1.QueryContractsRequest) | [QueryContractsResponse](#lbm.token.v1.QueryContractsResponse) | Contracts queries all token metadata. | GET|/lbm/token/v1/token_classes|
| `GranteeGrants` | [QueryGranteeGrantsRequest](#lbm.token.v1.QueryGranteeGrantsRequest) | [QueryGranteeGrantsResponse](#lbm.token.v1.QueryGranteeGrantsResponse) | GranteeGrants queries permissions on a given grantee. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - ErrInvalidAddress - `grantee` is of invalid format. | GET|/lbm/token/v1/token_classes/{contract_id}/grants/{grantee}|
| `Approved` | [QueryApprovedRequest](#lbm.token.v1.QueryApprovedRequest) | [QueryApprovedResponse](#lbm.token.v1.QueryApprovedResponse) | Approved queries authorization on a given proxy approver pair. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - ErrInvalidAddress - `proxy` is of invalid format. - `approver` is of invalid format. | GET|/lbm/token/v1/token_classes/{contract_id}/accounts/{approver}/proxies/{proxy}|
| `Approvers` | [QueryApproversRequest](#lbm.token.v1.QueryApproversRequest) | [QueryApproversResponse](#lbm.token.v1.QueryApproversResponse) | Approvers queries approvers on a given proxy. Throws: - ErrInvalidRequest - `contract_id` is of invalid format. - ErrInvalidAddress - `proxy` is of invalid format. | GET|/lbm/token/v1/token_classes/{contract_id}/accounts/{address}/approvers|

 <!-- end services -->



<a name="lbm/token/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lbm/token/v1/tx.proto



<a name="lbm.token.v1.MsgApprove"></a>

### MsgApprove
MsgApprove defines the Msg/Approve request type.

Throws:
- ErrInvalidAddress
  - `approver` is of invalid format.
  - `proxy` is of invalid format.
- ErrInvalidRequest
  - `contract_id` is of invalid format.

Signer: `approver`


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the token class. |
| `approver` | [string](#string) |  | address of the token approver which approves the authorization. |
| `proxy` | [string](#string) |  | address of the proxy which the authorization is granted to. |






<a name="lbm.token.v1.MsgApproveResponse"></a>

### MsgApproveResponse
MsgApproveResponse defines the Msg/Approve response type.






<a name="lbm.token.v1.MsgBurn"></a>

### MsgBurn
MsgBurn defines the Msg/Burn request type.

Throws:
- ErrInvalidAddress
  - `from` is of invalid format.
- ErrInvalidRequest
  - `contract_id` is of invalid format.
  - `amount` is not positive.

Signer: `from`


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the token class. |
| `from` | [string](#string) |  | address whose tokens are being burned. |
| `amount` | [string](#string) |  | number of tokens to burn. |






<a name="lbm.token.v1.MsgBurnFrom"></a>

### MsgBurnFrom
MsgBurnFrom defines the Msg/BurnFrom request type.

Throws:
- ErrInvalidAddress
  - `proxy` is of invalid format.
  - `from` is of invalid format.
- ErrInvalidRequest
  - `contract_id` is of invalid format.
  - `amount` is not positive.

Signer: `proxy`


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the token class. |
| `proxy` | [string](#string) |  | address which triggers the burn. |
| `from` | [string](#string) |  | address which the tokens will be burnt from. |
| `amount` | [string](#string) |  | the amount of the burn. |






<a name="lbm.token.v1.MsgBurnFromResponse"></a>

### MsgBurnFromResponse
MsgBurnFromResponse defines the Msg/BurnFrom response type.






<a name="lbm.token.v1.MsgBurnResponse"></a>

### MsgBurnResponse
MsgBurnResponse defines the Msg/Burn response type.






<a name="lbm.token.v1.MsgGrantPermission"></a>

### MsgGrantPermission
MsgGrantPermission defines the Msg/GrantPermission request type.

Throws:
- ErrInvalidAddress
  - `granter` is of invalid format.
  - `grantee` is of invalid format.
- ErrInvalidRequest
  - `contract_id` is of invalid format.
  - `permission` is not a valid permission.

Signer: `granter`


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the token class. |
| `from` | [string](#string) |  | address of the granter which must have the permission to give. |
| `to` | [string](#string) |  | address of the grantee. |
| `permission` | [string](#string) |  | permission on the token class. |






<a name="lbm.token.v1.MsgGrantPermissionResponse"></a>

### MsgGrantPermissionResponse
MsgGrantPermissionResponse defines the Msg/GrantPermission response type.






<a name="lbm.token.v1.MsgIssue"></a>

### MsgIssue
MsgIssue defines the Msg/Issue request type.

Throws:
- ErrInvalidAddress
  - `owner` is of invalid format.
  - `to` is of invalid format.
- ErrInvalidRequest
  - `name` is empty.
  - `name` exceeds the app-specific limit in length.
  - `symbol` is of invalid format.
  - `image_uri` exceeds the app-specific limit in length.
  - `meta` exceeds the app-specific limit in length.
  - `decimals` is lesser than 0 or greater than 18.
  - `amount` is not positive.

Signer: `owner`


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | name defines the human-readable name of the token class. mandatory (not ERC20 compliant). |
| `symbol` | [string](#string) |  | symbol is an abbreviated name for token class. mandatory (not ERC20 compliant). |
| `image_uri` | [string](#string) |  | image_uri is an uri for the image of the token class stored off chain. |
| `meta` | [string](#string) |  | meta is a brief description of token class. |
| `decimals` | [int32](#int32) |  | decimals is the number of decimals which one must divide the amount by to get its user representation. |
| `mintable` | [bool](#bool) |  | mintable represents whether the token is allowed to mint. |
| `owner` | [string](#string) |  | the address which all permissions on the token class will be granted to (not a permanent property). |
| `to` | [string](#string) |  | the address to send the minted token to. mandatory. |
| `amount` | [string](#string) |  | amount of tokens to mint on issuance. mandatory. |






<a name="lbm.token.v1.MsgIssueResponse"></a>

### MsgIssueResponse
MsgIssueResponse defines the Msg/Issue response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | id of the new contract. |






<a name="lbm.token.v1.MsgMint"></a>

### MsgMint
MsgMint defines the Msg/Mint request type.

Throws:
- ErrInvalidAddress
  - `from` is of invalid format.
  - `to` is of invalid format.
- ErrInvalidRequest
  - `contract_id` is of invalid format.
  - `amount` is not positive.

Signer: `from`


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the token class. |
| `from` | [string](#string) |  | address which triggers the mint. |
| `to` | [string](#string) |  | recipient of the tokens. |
| `amount` | [string](#string) |  | number of tokens to mint. |






<a name="lbm.token.v1.MsgMintResponse"></a>

### MsgMintResponse
MsgMintResponse defines the Msg/Mint response type.






<a name="lbm.token.v1.MsgModify"></a>

### MsgModify
MsgModify defines the Msg/Modify request type.

Throws:
- ErrInvalidAddress
  - `owner` is of invalid format.
- ErrInvalidRequest
  - `contract_id` is of invalid format.
  - `changes` has duplicate keys.
  - `changes` has a key which is not allowed to modify.
  - `changes` is empty.

Signer: `owner`


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the contract. |
| `owner` | [string](#string) |  | the address of the grantee which must have modify permission. |
| `changes` | [Pair](#lbm.token.v1.Pair) | repeated | changes to apply. |






<a name="lbm.token.v1.MsgModifyResponse"></a>

### MsgModifyResponse
MsgModifyResponse defines the Msg/Modify response type.






<a name="lbm.token.v1.MsgRevokeOperator"></a>

### MsgRevokeOperator
MsgRevokeOperator defines the Msg/RevokeOperator request type.

Throws:
- ErrInvalidAddress
  - `holder` is of invalid format.
  - `operator` is of invalid format.
- ErrInvalidRequest
  - `contract_id` is of invalid format.

Signer: `holder`

Since: 0.46.0 (finschia)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the token class. |
| `holder` | [string](#string) |  | address of a holder which revokes the `operator` address as an operator. |
| `operator` | [string](#string) |  | address to rescind as an operator for `holder`. |






<a name="lbm.token.v1.MsgRevokeOperatorResponse"></a>

### MsgRevokeOperatorResponse
MsgRevokeOperatorResponse defines the Msg/RevokeOperator response type.

Since: 0.46.0 (finschia)






<a name="lbm.token.v1.MsgRevokePermission"></a>

### MsgRevokePermission
MsgRevokePermission defines the Msg/RevokePermission request type.

Throws:
- ErrInvalidAddress
  - `grantee` is of invalid format.
- ErrInvalidRequest
  - `contract_id` is of invalid format.
  - `permission` is not a valid permission.

Signer: `grantee`


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the token class. |
| `from` | [string](#string) |  | address of the grantee which abandons the permission. |
| `permission` | [string](#string) |  | permission on the token class. |






<a name="lbm.token.v1.MsgRevokePermissionResponse"></a>

### MsgRevokePermissionResponse
MsgRevokePermissionResponse defines the Msg/RevokePermission response type.






<a name="lbm.token.v1.MsgSend"></a>

### MsgSend
MsgSend defines the Msg/Send request type.

Throws:
- ErrInvalidAddress
  - `from` is of invalid format.
  - `to` is of invalid format.
- ErrInvalidRequest
  - `contract_id` is of invalid format.
  - `amount` is not positive.

Signer: `from`


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the token class. |
| `from` | [string](#string) |  | approver whose tokens are being sent. |
| `to` | [string](#string) |  | recipient of the tokens. |
| `amount` | [string](#string) |  | number of tokens to send. |






<a name="lbm.token.v1.MsgSendResponse"></a>

### MsgSendResponse
MsgSendResponse defines the Msg/Send response type.






<a name="lbm.token.v1.MsgTransferFrom"></a>

### MsgTransferFrom
MsgTransferFrom defines the Msg/TransferFrom request type.
Throws:
- ErrInvalidAddress
  - `proxy` is of invalid format.
  - `from` is of invalid format.
  - `to` is of invalid format.
- ErrInvalidRequest
  - `contract_id` is of invalid format.
  - `amount` is not positive.

Signer: `proxy`


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_id` | [string](#string) |  | contract id associated with the token class. |
| `proxy` | [string](#string) |  | the address of the proxy. |
| `from` | [string](#string) |  | the address which the transfer is from. |
| `to` | [string](#string) |  | the address which the transfer is to. |
| `amount` | [string](#string) |  | the amount of the transfer. |






<a name="lbm.token.v1.MsgTransferFromResponse"></a>

### MsgTransferFromResponse
MsgTransferFromResponse defines the Msg/TransferFrom response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="lbm.token.v1.Msg"></a>

### Msg
Msg defines the token Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Send` | [MsgSend](#lbm.token.v1.MsgSend) | [MsgSendResponse](#lbm.token.v1.MsgSendResponse) | Send defines a method to send tokens from one account to another account. Fires: - EventSent - transfer (deprecated, not typed) Throws: - ErrInvalidRequest: - the balance of `from` does not have enough tokens to spend. | |
| `TransferFrom` | [MsgTransferFrom](#lbm.token.v1.MsgTransferFrom) | [MsgTransferFromResponse](#lbm.token.v1.MsgTransferFromResponse) | TransferFrom defines a method to send tokens from one account to another account by the proxy. Fires: - EventSent - transfer_from (deprecated, not typed) Throws: - ErrUnauthorized: - the approver has not authorized the proxy. - ErrInvalidRequest: - the balance of `from` does not have enough tokens to spend. Note: the approval has no value of limit (not ERC20 compliant). | |
| `RevokeOperator` | [MsgRevokeOperator](#lbm.token.v1.MsgRevokeOperator) | [MsgRevokeOperatorResponse](#lbm.token.v1.MsgRevokeOperatorResponse) | RevokeOperator revoke the authorization of the operator to send the holder's tokens. Fires: - EventRevokedOperator Throws: - ErrNotFound: - there is no token class of `contract_id`. - there is no authorization by `holder` to `operator`. Note: it introduces breaking change, because the legacy clients cannot track this revocation. Since: 0.46.0 (finschia) | |
| `Approve` | [MsgApprove](#lbm.token.v1.MsgApprove) | [MsgApproveResponse](#lbm.token.v1.MsgApproveResponse) | Approve allows one to send tokens on behalf of the approver. Fires: - EventAuthorizedOperator - approve_token (deprecated, not typed) Throws: - ErrNotFound: - there is no token class of `contract_id`. - ErrInvalidRequest: - `approver` has already authorized `proxy`. | |
| `Issue` | [MsgIssue](#lbm.token.v1.MsgIssue) | [MsgIssueResponse](#lbm.token.v1.MsgIssueResponse) | Issue defines a method to create a class of token. it grants `mint`, `burn` and `modify` permissions on the token class to its creator (see also `mintable`). Fires: - EventIssue - EventMinted - issue (deprecated, not typed) | |
| `GrantPermission` | [MsgGrantPermission](#lbm.token.v1.MsgGrantPermission) | [MsgGrantPermissionResponse](#lbm.token.v1.MsgGrantPermissionResponse) | GrantPermission allows one to mint or burn tokens or modify a token metadata. Fires: - EventGrant - grant_perm (deprecated, not typed) Throws: - ErrUnauthorized - `granter` does not have `permission`. - ErrInvalidRequest - `grantee` already has `permission`. | |
| `RevokePermission` | [MsgRevokePermission](#lbm.token.v1.MsgRevokePermission) | [MsgRevokePermissionResponse](#lbm.token.v1.MsgRevokePermissionResponse) | RevokePermission abandons a permission. Fires: - EventAbandon - revoke_perm (deprecated, not typed) Throws: - ErrUnauthorized - `grantee` does not have `permission`. | |
| `Mint` | [MsgMint](#lbm.token.v1.MsgMint) | [MsgMintResponse](#lbm.token.v1.MsgMintResponse) | Mint defines a method to mint tokens. Fires: - EventMinted - mint (deprecated, not typed) Throws: - ErrUnauthorized - `from` does not have `mint` permission. | |
| `Burn` | [MsgBurn](#lbm.token.v1.MsgBurn) | [MsgBurnResponse](#lbm.token.v1.MsgBurnResponse) | Burn defines a method to burn tokens. Fires: - EventBurned - burn (deprecated, not typed) Throws: - ErrUnauthorized - `from` does not have `burn` permission. - ErrInvalidRequest: - the balance of `from` does not have enough tokens to burn. | |
| `BurnFrom` | [MsgBurnFrom](#lbm.token.v1.MsgBurnFrom) | [MsgBurnFromResponse](#lbm.token.v1.MsgBurnFromResponse) | BurnFrom defines a method to burn tokens by the proxy. Fires: - EventBurned - burn_from (deprecated, not typed) Throws: - ErrUnauthorized - `proxy` does not have `burn` permission. - the approver has not authorized `proxy`. - ErrInvalidRequest: - the balance of `from` does not have enough tokens to burn. | |
| `Modify` | [MsgModify](#lbm.token.v1.MsgModify) | [MsgModifyResponse](#lbm.token.v1.MsgModifyResponse) | Modify defines a method to modify a token class. Fires: - EventModified - modify_token (deprecated, not typed) Throws: - ErrUnauthorized - the proxy does not have `modify` permission. - ErrNotFound - there is no token class of `contract_id`. | |

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

