package client

import (
	govclient "github.com/line/lbm-sdk/x/gov/client"
	"github.com/line/lbm-sdk/x/foundation/client/cli"
)

// ProposalHandler is the param change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewProposalCmdFoundationExec)
