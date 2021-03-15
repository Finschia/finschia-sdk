package client

import (
	govclient "github.com/line/lbm-sdk/v2/x/gov/client"
	"github.com/line/lbm-sdk/v2/x/params/client/cli"
	"github.com/line/lbm-sdk/v2/x/params/client/rest"
)

// ProposalHandler is the param change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewSubmitParamChangeProposalTxCmd, rest.ProposalRESTHandler)
