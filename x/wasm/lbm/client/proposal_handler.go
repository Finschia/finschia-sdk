package client

import (
	govclient "github.com/line/lbm-sdk/x/gov/client"
	wasmclient "github.com/line/lbm-sdk/x/wasm/client"
	"github.com/line/lbm-sdk/x/wasm/lbm/client/cli"
	"github.com/line/lbm-sdk/x/wasm/lbm/client/rest"
)

var ProposalHandlers = append(
	wasmclient.ProposalHandlers, // wasmd's proposals
	[]govclient.ProposalHandler{
		govclient.NewProposalHandler(cli.ProposalDeactivateContractCmd, rest.DeactivateContractProposalHandler),
		govclient.NewProposalHandler(cli.ProposalActivateContractCmd, rest.ActivateContractProposalHandler),
	}...,
)
