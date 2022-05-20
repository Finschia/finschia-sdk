package client

import (
	"github.com/line/lbm-sdk/x/foundation/client/cli"
	"github.com/line/lbm-sdk/x/foundation/client/rest"
	govclient "github.com/line/lbm-sdk/x/gov/client"
)

var UpdateFoundationParamsProposalHandler = govclient.NewProposalHandler(cli.NewProposalCmdUpdateFoundationParams, rest.DummyRESTHandler)
var UpdateValidatorAuthsProposalHandler = govclient.NewProposalHandler(cli.NewProposalCmdUpdateValidatorAuths, rest.DummyRESTHandler)
