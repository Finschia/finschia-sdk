package client

import (
	"github.com/line/lbm-sdk/x/distribution/client/cli"
	"github.com/line/lbm-sdk/x/distribution/client/rest"
	govclient "github.com/line/lbm-sdk/x/gov/client"
)

// param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
