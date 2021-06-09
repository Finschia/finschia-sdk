package mock_test

import (
	"testing"

	ostproto "github.com/line/ostracon/proto/ostracon/types"
	osttypes "github.com/line/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lfb-sdk/x/ibc/testing/mock"
)

const chainID = "testChain"

func TestGetPubKey(t *testing.T) {
	pv := mock.NewPV()
	pk, err := pv.GetPubKey()
	require.NoError(t, err)
	require.Equal(t, "ed25519", pk.Type())
}

func TestSignVote(t *testing.T) {
	pv := mock.NewPV()
	pk, _ := pv.GetPubKey()

	vote := &ostproto.Vote{Height: 2}
	pv.SignVote(chainID, vote)

	msg := osttypes.VoteSignBytes(chainID, vote)
	ok := pk.VerifySignature(msg, vote.Signature)
	require.True(t, ok)
}

func TestSignProposal(t *testing.T) {
	pv := mock.NewPV()
	pk, _ := pv.GetPubKey()

	proposal := &ostproto.Proposal{Round: 2}
	pv.SignProposal(chainID, proposal)

	msg := osttypes.ProposalSignBytes(chainID, proposal)
	ok := pk.VerifySignature(msg, proposal.Signature)
	require.True(t, ok)
}
