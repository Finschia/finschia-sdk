package solomachine

import (
	"github.com/spf13/cobra"

	"github.com/line/lfb-sdk/x/ibc/light-clients/06-solomachine/client/cli"
	"github.com/line/lfb-sdk/x/ibc/light-clients/06-solomachine/types"
)

// Name returns the solo machine client name.
func Name() string {
	return types.SubModuleName
}

// GetTxCmd returns the root tx command for the solo machine client.
func GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}
