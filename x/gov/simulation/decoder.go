package simulation

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/line/lbm-sdk/v2/codec"
	types3 "github.com/line/lbm-sdk/v2/store/types"
	types2 "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/x/gov/keeper"
	"github.com/line/lbm-sdk/v2/x/gov/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding gov type.
func NewDecodeStore(cdc codec.Marshaler) types2.StoreDecoder {
	return types2.StoreDecoder{
		Marshal: func(key []byte) func(obj interface{}) []byte {
			switch {
			case bytes.Equal(key[:1], types.ProposalsKeyPrefix):
				return keeper.GetProposerMarshalFunc(cdc)
			case bytes.Equal(key[:1], types.ActiveProposalQueuePrefix),
				 bytes.Equal(key[:1], types.InactiveProposalQueuePrefix),
				 bytes.Equal(key[:1], types.ProposalIDKey):
			 	return types3.GetBytesMarshalFunc()
			case bytes.Equal(key[:1], types.DepositsKeyPrefix):
				return keeper.GetDepositMarshalFunc(cdc)
			case bytes.Equal(key[:1], types.VotesKeyPrefix):
				return keeper.GetVoteMarshalFunc(cdc)
			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, key, key))
			}
		},
		Unmarshal: func(key []byte) func(value []byte) interface{} {
			switch {
			case bytes.Equal(key[:1], types.ProposalsKeyPrefix):
				return keeper.GetProposerUnmarshalFunc(cdc)
			case bytes.Equal(key[:1], types.ActiveProposalQueuePrefix),
				bytes.Equal(key[:1], types.InactiveProposalQueuePrefix),
				bytes.Equal(key[:1], types.ProposalIDKey):
				return types3.GetBytesUnmarshalFunc()
			case bytes.Equal(key[:1], types.DepositsKeyPrefix):
				return keeper.GetDepositUnmarshalFunc(cdc)
			case bytes.Equal(key[:1], types.VotesKeyPrefix):
				return keeper.GetVoteUnmarshalFunc(cdc)
			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, key, key))
			}
		},
		LogPair: func(kvA, kvB types3.KOPair) string {
			switch {
			case bytes.Equal(kvA.Key[:1], types.ProposalsKeyPrefix):
				proposalA := *kvA.Value.(*types.Proposal)
				proposalB := *kvB.Value.(*types.Proposal)
				return fmt.Sprintf("%v\n%v", proposalA, proposalB)

			case bytes.Equal(kvA.Key[:1], types.ActiveProposalQueuePrefix),
				bytes.Equal(kvA.Key[:1], types.InactiveProposalQueuePrefix),
				bytes.Equal(kvA.Key[:1], types.ProposalIDKey):
				proposalIDA := binary.LittleEndian.Uint64(kvA.Value.([]byte))
				proposalIDB := binary.LittleEndian.Uint64(kvB.Value.([]byte))
				return fmt.Sprintf("proposalIDA: %d\nProposalIDB: %d", proposalIDA, proposalIDB)

			case bytes.Equal(kvA.Key[:1], types.DepositsKeyPrefix):
				depositA := *kvA.Value.(*types.Deposit)
				depositB := *kvB.Value.(*types.Deposit)
				return fmt.Sprintf("%v\n%v", depositA, depositB)

			case bytes.Equal(kvA.Key[:1], types.VotesKeyPrefix):
				voteA := *kvA.Value.(*types.Vote)
				voteB := *kvB.Value.(*types.Vote)
				return fmt.Sprintf("%v\n%v", voteA, voteB)

			default:
				panic(fmt.Sprintf("invalid governance key prefix %X", kvA.Key[:1]))
			}
		},
	}
}
