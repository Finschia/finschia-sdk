package client

import (
	govclient "github.com/Finschia/finschia-sdk/x/gov/client"
	"github.com/Finschia/finschia-sdk/x/params/client/cli"
)

// ProposalHandler is the param change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewSubmitParamChangeProposalTxCmd)
