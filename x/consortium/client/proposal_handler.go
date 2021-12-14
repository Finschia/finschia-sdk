package client

import (
	govclient "github.com/line/lbm-sdk/x/gov/client"
	"github.com/line/lbm-sdk/x/consortium/client/cli"
	"github.com/line/lbm-sdk/x/consortium/client/rest"
)

var UpdateConsortiumParamsProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitUpdateConsortiumParamsProposal, rest.ProposalUpdateConsortiumParamsRESTHandler)
var UpdateValidatorAuthsProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitUpdateValidatorAuthsProposal, rest.ProposalUpdateValidatorAuthsRESTHandler)
