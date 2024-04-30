package client

import (
	"github.com/Finschia/finschia-sdk/x/fswap/client/cli"
	govclient "github.com/Finschia/finschia-sdk/x/gov/client"
)

var ProposalHandler = govclient.NewProposalHandler(cli.NewCmdSwapInitProposal)
