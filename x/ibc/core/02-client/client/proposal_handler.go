package client

import (
	govclient "github.com/line/lbm-sdk/x/gov/client"

	"github.com/line/lbm-sdk/x/ibc/core/02-client/client/cli"
)

var (
	UpdateClientProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitUpdateClientProposal)
	UpgradeProposalHandler      = govclient.NewProposalHandler(cli.NewCmdSubmitUpgradeProposal)
)
