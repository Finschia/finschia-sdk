package types

import (
	"bytes"
	"time"

	octypes "github.com/line/ostracon/types"

	sdkerrors "github.com/line/lbm-sdk/types/errors"
	clienttypes "github.com/line/lbm-sdk/x/ibc/core/02-client/types"
	commitmenttypes "github.com/line/lbm-sdk/x/ibc/core/23-commitment/types"
	"github.com/line/lbm-sdk/x/ibc/core/exported"
)

var _ exported.Header = &Header{}

// ConsensusState returns the updated consensus state associated with the header
func (h Header) ConsensusState() *ConsensusState {
	return &ConsensusState{
		Timestamp:          h.GetTime(),
		Root:               commitmenttypes.NewMerkleRoot(h.Header.GetAppHash()),
		NextValidatorsHash: h.Header.NextValidatorsHash,
	}
}

// ClientType defines that the Header is a Ostracon consensus algorithm
func (h Header) ClientType() string {
	return exported.Ostracon
}

// GetHeight returns the current height. It returns 0 if the ostracon
// header is nil.
// NOTE: the header.Header is checked to be non nil in ValidateBasic.
func (h Header) GetHeight() exported.Height {
	revision := clienttypes.ParseChainID(h.Header.ChainID)
	return clienttypes.NewHeight(revision, uint64(h.Header.Height))
}

// GetTime returns the current block timestamp. It returns a zero time if
// the ostracon header is nil.
// NOTE: the header.Header is checked to be non nil in ValidateBasic.
func (h Header) GetTime() time.Time {
	return h.Header.Time
}

// ValidateBasic calls the SignedHeader ValidateBasic function and checks
// that validatorsets are not nil.
// NOTE: TrustedHeight and TrustedValidators may be empty when creating client
// with MsgCreateClient
func (h Header) ValidateBasic() error {
	if h.SignedHeader == nil {
		return sdkerrors.Wrap(clienttypes.ErrInvalidHeader, "ostracon signed header cannot be nil")
	}
	if h.Header == nil {
		return sdkerrors.Wrap(clienttypes.ErrInvalidHeader, "ostracon header cannot be nil")
	}
	ocSignedHeader, err := octypes.SignedHeaderFromProto(h.SignedHeader)
	if err != nil {
		return sdkerrors.Wrap(err, "header is not a ostracon header")
	}
	if err := ocSignedHeader.ValidateBasic(h.Header.GetChainID()); err != nil {
		return sdkerrors.Wrap(err, "header failed basic validation")
	}

	// TrustedHeight is less than Header for updates
	// and less than or equal to Header for misbehaviour
	if h.TrustedHeight.GT(h.GetHeight()) {
		return sdkerrors.Wrapf(ErrInvalidHeaderHeight, "TrustedHeight %d must be less than or equal to header height %d",
			h.TrustedHeight, h.GetHeight())
	}

	if h.ValidatorSet == nil {
		return sdkerrors.Wrap(clienttypes.ErrInvalidHeader, "validator set is nil")
	}
	ocValset, err := octypes.ValidatorSetFromProto(h.ValidatorSet)
	if err != nil {
		return sdkerrors.Wrap(err, "validator set is not ostracon validator set")
	}
	if !bytes.Equal(h.Header.ValidatorsHash, ocValset.Hash()) {
		return sdkerrors.Wrap(clienttypes.ErrInvalidHeader, "validator set does not match hash")
	}
	return nil
}
