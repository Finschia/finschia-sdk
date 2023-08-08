package client

import (
	"github.com/Finschia/finschia-rdk/x/distribution/client/cli"
	govclient "github.com/Finschia/finschia-rdk/x/gov/client"
)

// ProposalHandler is the community spend proposal handler.
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal)
)
