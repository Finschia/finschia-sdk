package client

import (
	"github.com/line/lbm-sdk/x/distribution/client/cli"
	govclient "github.com/line/lbm-sdk/x/gov/client"
)

// ProposalHandler is the community spend proposal handler.
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal)
)
