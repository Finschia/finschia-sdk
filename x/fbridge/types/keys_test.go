package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func TestFbridgeKeys(t *testing.T) {
	require.Equal(t, []byte{types.KeySeqToBlocknumPrefix[0], 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}, types.SeqToBlocknumKey(1))
	require.Equal(t, []byte{types.KeyProposalPrefix[0], 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}, types.ProposalKey(1))
	require.Equal(t, []byte{types.KeyProposalVotePrefix[0], 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}, types.VotesKey(1))

	vvkey := types.VoterVoteKey(1, []byte("voter"))
	require.Equal(t, []byte{types.KeyProposalVotePrefix[0], 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x5, 0x76, 0x6f, 0x74, 0x65, 0x72}, vvkey)
	pid, voter := types.SplitVoterVoteKey(vvkey)
	require.Equal(t, uint64(1), pid)
	require.Equal(t, sdk.AccAddress("voter"), voter)

	rkey := types.RoleKey([]byte("assignee"))
	require.Equal(t, []byte{types.KeyRolePrefix[0], 0x8, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x65}, rkey)
	role := types.SplitRoleKey(rkey)
	require.Equal(t, sdk.AccAddress("assignee"), role)

	bskey := types.BridgeSwitchKey([]byte("guardian"))
	require.Equal(t, []byte{types.KeyBridgeSwitchPrefix[0], 0x8, 0x67, 0x75, 0x61, 0x72, 0x64, 0x69, 0x61, 0x6e}, bskey)
	guardian := types.SplitBridgeSwitchKey(bskey)
	require.Equal(t, sdk.AccAddress("guardian"), guardian)
}
