package client

import (
	"github.com/Finschia/finschia-sdk/x/distribution/client/cli"
	govclient "github.com/Finschia/finschia-sdk/x/gov/client"
)

// ProposalHandler is the community spend proposal handler.
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal)
)
