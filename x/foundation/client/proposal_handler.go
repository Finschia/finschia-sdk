package client

import (
	"github.com/line/lbm-sdk/x/foundation/client/cli"
	govclient "github.com/line/lbm-sdk/x/gov/client"
)

// ProposalHandler is the param change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewProposalCmdFoundationExec)
