package client

import (
	govclient "github.com/Finschia/finschia-rdk/x/gov/client"
	"github.com/Finschia/finschia-rdk/x/upgrade/client/cli"
)

var (
	ProposalHandler       = govclient.NewProposalHandler(cli.NewCmdSubmitUpgradeProposal)
	CancelProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitCancelUpgradeProposal)
)
