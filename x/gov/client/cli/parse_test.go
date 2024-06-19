package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-sdk/testutil"
)

func TestParseSubmitProposalFlags(t *testing.T) {
	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Test Proposal",
  "description": "My awesome proposal",
  "type": "Text",
  "deposit": "1000test"
}
`)

	badJSON := testutil.WriteToNewTempFile(t, "bad json")
	fs := NewCmdSubmitProposal().Flags()

	// nonexistent json
	_ = fs.Set(FlagProposal, "fileDoesNotExist")
	_, err := parseSubmitProposalFlags(fs)
	require.Error(t, err)

	// invalid json
	_ = fs.Set(FlagProposal, badJSON.Name())
	_, err = parseSubmitProposalFlags(fs)
	require.Error(t, err)

	// ok json
	_ = fs.Set(FlagProposal, okJSON.Name())
	proposal1, err := parseSubmitProposalFlags(fs)
	require.Nil(t, err, "unexpected error")
	require.Equal(t, "Test Proposal", proposal1.Title)
	require.Equal(t, "My awesome proposal", proposal1.Description)
	require.Equal(t, "Text", proposal1.Type)
	require.Equal(t, "1000test", proposal1.Deposit)

	// flags that can't be used with --proposal
	for _, incompatibleFlag := range ProposalFlags {
		_ = fs.Set(incompatibleFlag, "some value")
		_, err := parseSubmitProposalFlags(fs)
		require.Error(t, err)
		_ = fs.Set(incompatibleFlag, "")
	}

	// no --proposal, only flags
	_ = fs.Set(FlagProposal, "")
	_ = fs.Set(FlagTitle, proposal1.Title)
	_ = fs.Set(FlagDescription, proposal1.Description)
	_ = fs.Set(FlagProposalType, proposal1.Type)
	_ = fs.Set(FlagDeposit, proposal1.Deposit)
	proposal2, err := parseSubmitProposalFlags(fs)

	require.Nil(t, err, "unexpected error")
	require.Equal(t, proposal1.Title, proposal2.Title)
	require.Equal(t, proposal1.Description, proposal2.Description)
	require.Equal(t, proposal1.Type, proposal2.Type)
	require.Equal(t, proposal1.Deposit, proposal2.Deposit)

	err = okJSON.Close()
	require.Nil(t, err, "unexpected error")
	err = badJSON.Close()
	require.Nil(t, err, "unexpected error")
}
