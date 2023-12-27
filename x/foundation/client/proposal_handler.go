package client

import (
	"github.com/Finschia/finschia-sdk/x/foundation/client/cli"
	govclient "github.com/Finschia/finschia-sdk/x/gov/client"
)

// ProposalHandler is the param change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewProposalCmdFoundationExec)
