package types

import (
	"encoding/binary"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/address"
	"github.com/Finschia/finschia-sdk/types/kv"
)

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "fbridge"

	// StoreKey is the store key string for fbridge
	StoreKey = ModuleName

	// MemStoreKey is the in-memory store key string for fbridge
	MemStoreKey = "mem_" + StoreKey
)

// - 0x01: params
// - 0x02: next sequence number for bridge send
//
// - 0x10: next proposal ID
// 	 0x11<proposalID (8-byte)>: proposal
//   0x12<proposalID (8-byte)><voterAddrLen (1-byte)><voterAddr>: vote
// - 0x13<addrLen (1-byte)><targetAddr>: role
// - 0x14<addrLen (1-byte)><guardianAddr>: bridge switch
//
// - 0xF0: memstore initialized
// - 0xF1: role metadata

var (
	KeyParams      = []byte{0x01} // key for fbridge module params
	KeyNextSeqSend = []byte{0x02} // key for the next bridge send sequence

	KeyNextProposalID     = []byte{0x10} // key for the next role proposal ID
	KeyProposalPrefix     = []byte{0x11} // key prefix for the role proposal
	KeyProposalVotePrefix = []byte{0x12} // key prefix for the role proposal vote
	KeyRolePrefix         = []byte{0x13} // key prefix for the role of an address
	KeyBridgeSwitch       = []byte{0x14} // key for the switch to halt

	KeyMemInitialized  = []byte{0xF0}
	KeyMemRoleMetadata = []byte{0xF1} // key for the role metadata
	KeyMemBridgeStatus = []byte{0xF2} // key for the bridge status
)

// GetProposalIDBytes returns the byte representation of the proposalID
func GetProposalIDBytes(proposalID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, proposalID)
	return bz
}

// ProposalKey key of a specific role proposal
func ProposalKey(proposalID uint64) []byte {
	return append(KeyProposalPrefix, GetProposalIDBytes(proposalID)...)
}

// VotesKey gets the first part of the votes key based on the proposalID
func VotesKey(proposalID uint64) []byte {
	return append(KeyProposalVotePrefix, GetProposalIDBytes(proposalID)...)
}

// VoterVoteKey key of a specific vote from the store
func VoterVoteKey(proposalID uint64, voterAddr sdk.AccAddress) []byte {
	return append(VotesKey(proposalID), address.MustLengthPrefix(voterAddr.Bytes())...)
}

// SplitVoterVoteKey split the voter key and returns the proposal id and voter address
func SplitVoterVoteKey(key []byte) (uint64, sdk.AccAddress) {
	kv.AssertKeyAtLeastLength(key, 11)
	proposalID := binary.BigEndian.Uint64(key[1:9])
	voter := sdk.AccAddress(key[10:])
	return proposalID, voter
}

// RoleKey key of a specific role of the address from the store
func RoleKey(target sdk.AccAddress) []byte {
	return append(KeyRolePrefix, address.MustLengthPrefix(target.Bytes())...)
}

// SplitRoleKey split the role key and returns the address
func SplitRoleKey(key []byte) sdk.AccAddress {
	kv.AssertKeyAtLeastLength(key, 3)
	return key[2:]
}

func BridgeSwitchKey(guardian sdk.AccAddress) []byte {
	return append(KeyBridgeSwitch, address.MustLengthPrefix(guardian.Bytes())...)
}

// SplitBridgeSwitchKey split the bridge switch key and returns the guardian address
func SplitBridgeSwitchKey(key []byte) sdk.AccAddress {
	kv.AssertKeyAtLeastLength(key, 3)
	return key[2:]
}
