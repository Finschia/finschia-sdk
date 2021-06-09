package client

import (
	govclient "github.com/line/lfb-sdk/x/gov/client"
	"github.com/line/lfb-sdk/x/params/client/cli"
	"github.com/line/lfb-sdk/x/params/client/rest"
)

// ProposalHandler is the param change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewSubmitParamChangeProposalTxCmd, rest.ProposalRESTHandler)
