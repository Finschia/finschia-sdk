package client

import (
	govclient "github.com/line/lbm-sdk/x/gov/client"
	"github.com/line/lbm-sdk/x/consortium/client/cli"
	"github.com/line/lbm-sdk/x/consortium/client/rest"
)

var DisableConsortiumProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitDisableConsortiumProposal, rest.ProposalDisableConsortiumRESTHandler)
var EditAllowedValidatorsProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitEditAllowedValidatorsProposal, rest.ProposalEditAllowedValidatorsRESTHandler)
