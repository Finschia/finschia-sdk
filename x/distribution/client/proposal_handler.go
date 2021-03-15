package client

import (
	"github.com/line/lbm-sdk/v2/x/distribution/client/cli"
	"github.com/line/lbm-sdk/v2/x/distribution/client/rest"
	govclient "github.com/line/lbm-sdk/v2/x/gov/client"
)

// ProposalHandler is the community spend proposal handler.
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
