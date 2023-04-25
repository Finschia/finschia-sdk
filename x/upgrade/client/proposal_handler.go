package client

import (
	govclient "github.com/Finschia/finschia-sdk/x/gov/client"
	"github.com/Finschia/finschia-sdk/x/upgrade/client/cli"
)

var (
	ProposalHandler       = govclient.NewProposalHandler(cli.NewCmdSubmitUpgradeProposal)
	CancelProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitCancelUpgradeProposal)
)
