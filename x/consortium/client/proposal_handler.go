package client

import (
	"github.com/line/lbm-sdk/x/consortium/client/cli"
	"github.com/line/lbm-sdk/x/consortium/client/rest"
	govclient "github.com/line/lbm-sdk/x/gov/client"
)

var UpdateConsortiumParamsProposalHandler = govclient.NewProposalHandler(cli.NewProposalCmdUpdateConsortiumParams, rest.DummyRESTHandler)
var UpdateValidatorAuthsProposalHandler = govclient.NewProposalHandler(cli.NewProposalCmdUpdateValidatorAuths, rest.DummyRESTHandler)
