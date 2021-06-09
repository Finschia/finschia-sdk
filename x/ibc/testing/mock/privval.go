package mock

import (
	"github.com/line/ostracon/crypto"
	ostproto "github.com/line/ostracon/proto/ostracon/types"
	osttypes "github.com/line/ostracon/types"

	cryptocodec "github.com/line/lfb-sdk/crypto/codec"
	"github.com/line/lfb-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/line/lfb-sdk/crypto/types"
)

var _ osttypes.PrivValidator = PV{}

// MockPV implements PrivValidator without any safety or persistence.
// Only use it for testing.
type PV struct {
	PrivKey cryptotypes.PrivKey
}

func NewPV() PV {
	return PV{ed25519.GenPrivKey()}
}

// GetPubKey implements PrivValidator interface
func (pv PV) GetPubKey() (crypto.PubKey, error) {
	return cryptocodec.ToTmPubKeyInterface(pv.PrivKey.PubKey())
}

// SignVote implements PrivValidator interface
func (pv PV) SignVote(chainID string, vote *ostproto.Vote) error {
	signBytes := osttypes.VoteSignBytes(chainID, vote)
	sig, err := pv.PrivKey.Sign(signBytes)
	if err != nil {
		return err
	}
	vote.Signature = sig
	return nil
}

// SignProposal implements PrivValidator interface
func (pv PV) SignProposal(chainID string, proposal *ostproto.Proposal) error {
	signBytes := osttypes.ProposalSignBytes(chainID, proposal)
	sig, err := pv.PrivKey.Sign(signBytes)
	if err != nil {
		return err
	}
	proposal.Signature = sig
	return nil
}
