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

	// StoreKey is the store key string for distribution
	StoreKey = ModuleName
)

// - 0x01: params
// - 0x02: next sequence number for bridge send
//
// - 0x10: next proposal ID
// 	 0x11<proposalID (8-byte)>: proposal
//   0x12<proposalID (8-byte)><voterAddrLen (1-byte)><voterAddr>: vote
// - 0x13<addrLen (1-byte)><targetAddr>: role

var (
	KeyParams      = []byte{0x01} // key for fbridge module params
	KeyNextSeqSend = []byte{0x02} // key for the next bridge send sequence

	KeyNextProposalID     = []byte{0x10} // key for the next role proposal ID
	KeyProposalPrefix     = []byte{0x11} // key prefix for the role proposal
	KeyProposalVotePrefix = []byte{0x12} // key prefix for the role proposal vote
	KeyRolePrefix         = []byte{0x13} // key prefix for the role of an address
	KeyRoleMetadata       = []byte{0x14}
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

// SplitProposalKey split the proposal key and returns the proposal id
func SplitProposalKey(key []byte) (proposalID uint64) {
	kv.AssertKeyLength(key[1:], 8)

	return binary.BigEndian.Uint64(key[1:])
}

// VotesKey gets the first part of the votes key based on the proposalID
func VotesKey(proposalID uint64) []byte {
	return append(KeyProposalVotePrefix, GetProposalIDBytes(proposalID)...)
}

// VoterVoteKey key of a specific vote from the store
func VoterVoteKey(proposalID uint64, voterAddr sdk.AccAddress) []byte {
	return append(VotesKey(proposalID), address.MustLengthPrefix(voterAddr.Bytes())...)
}

// RoleKey key of a specific role of the address from the store
func RoleKey(target sdk.AccAddress) []byte {
	return append(KeyRolePrefix, address.MustLengthPrefix(target.Bytes())...)
}
