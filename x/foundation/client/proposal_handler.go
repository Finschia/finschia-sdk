package client

import (
	"github.com/Finschia/finschia-rdk/x/foundation/client/cli"
	govclient "github.com/Finschia/finschia-rdk/x/gov/client"
)

// ProposalHandler is the param change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewProposalCmdFoundationExec)
