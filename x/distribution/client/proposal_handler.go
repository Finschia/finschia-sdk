package client

import (
	"github.com/line/lfb-sdk/x/distribution/client/cli"
	"github.com/line/lfb-sdk/x/distribution/client/rest"
	govclient "github.com/line/lfb-sdk/x/gov/client"
)

// ProposalHandler is the community spend proposal handler.
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
