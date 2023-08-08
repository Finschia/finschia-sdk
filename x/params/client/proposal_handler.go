package client

import (
	govclient "github.com/Finschia/finschia-rdk/x/gov/client"
	"github.com/Finschia/finschia-rdk/x/params/client/cli"
)

// ProposalHandler is the param change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewSubmitParamChangeProposalTxCmd)
