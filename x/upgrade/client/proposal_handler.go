package client

import (
	govclient "github.com/line/lbm-sdk/x/gov/client"
	"github.com/line/lbm-sdk/x/upgrade/client/cli"
)

var (
	ProposalHandler       = govclient.NewProposalHandler(cli.NewCmdSubmitUpgradeProposal)
	CancelProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitCancelUpgradeProposal)
)
