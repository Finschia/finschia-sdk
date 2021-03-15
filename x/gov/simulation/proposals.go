package simulation

import (
	"math/rand"

	simappparams "github.com/line/lbm-sdk/v2/simapp/params"
	sdk "github.com/line/lbm-sdk/v2/types"
	simtypes "github.com/line/lbm-sdk/v2/types/simulation"
	"github.com/line/lbm-sdk/v2/x/gov/types"
	"github.com/line/lbm-sdk/v2/x/simulation"
)

// OpWeightSubmitTextProposal app params key for text proposal
const OpWeightSubmitTextProposal = "op_weight_submit_text_proposal"

// ProposalContents defines the module weighted proposals' contents
func ProposalContents() []simtypes.WeightedProposalContent {
	return []simtypes.WeightedProposalContent{
		simulation.NewWeightedProposalContent(
			OpWeightMsgDeposit,
			simappparams.DefaultWeightTextProposal,
			SimulateTextProposalContent,
		),
	}
}

// SimulateTextProposalContent returns a random text proposal content.
func SimulateTextProposalContent(r *rand.Rand, _ sdk.Context, _ []simtypes.Account) simtypes.Content {
	return types.NewTextProposal(
		simtypes.RandStringOfLength(r, 140),
		simtypes.RandStringOfLength(r, 5000),
	)
}
