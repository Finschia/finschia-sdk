package cli

import (
	"github.com/Finschia/finschia-rdk/codec"
	"github.com/Finschia/finschia-rdk/internal/os"
	"github.com/Finschia/finschia-rdk/x/distribution/types"
)

// ParseCommunityPoolSpendProposalWithDeposit reads and parses a CommunityPoolSpendProposalWithDeposit from a file.
func ParseCommunityPoolSpendProposalWithDeposit(cdc codec.JSONCodec, proposalFile string) (types.CommunityPoolSpendProposalWithDeposit, error) {
	proposal := types.CommunityPoolSpendProposalWithDeposit{}

	// 2M size limit is enough for a proposal.
	// Check the proposals:
	// https://hubble.figment.io/cosmos/chains/cosmoshub-4/governance
	contents, err := os.ReadFileWithSizeLimit(proposalFile, 2*1024*1024)
	if err != nil {
		return proposal, err
	}

	if err = cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
