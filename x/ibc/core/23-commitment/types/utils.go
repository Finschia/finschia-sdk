package types

import (
	ics23 "github.com/confio/ics23/go"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	crypto "github.com/line/ostracon/proto/ostracon/crypto"
)

// ConvertProofs converts crypto.ProofOps into MerkleProof
func ConvertProofs(tmProof *crypto.ProofOps) (MerkleProof, error) {
	if tmProof == nil {
		return MerkleProof{}, sdkerrors.Wrapf(ErrInvalidMerkleProof, "ostracon proof is nil")
	}
	// Unmarshal all proof ops to CommitmentProof
	proofs := make([]*ics23.CommitmentProof, len(tmProof.Ops))
	for i, op := range tmProof.Ops {
		var p ics23.CommitmentProof
		err := p.Unmarshal(op.Data)
		if err != nil || p.Proof == nil {
			return MerkleProof{}, sdkerrors.Wrapf(ErrInvalidMerkleProof, "could not unmarshal proof op into CommitmentProof at index %d: %v", i, err)
		}
		proofs[i] = &p
	}
	return MerkleProof{
		Proofs: proofs,
	}, nil
}
