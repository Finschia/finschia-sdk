package cli

import (
	"github.com/spf13/cobra"

	controllercli "github.com/line/lbm-sdk/x/ibc/applications/27-interchain-accounts/controller/client/cli"
	hostcli "github.com/line/lbm-sdk/x/ibc/applications/27-interchain-accounts/host/client/cli"
)

// GetQueryCmd returns the query commands for the interchain-accounts submodule
func GetQueryCmd() *cobra.Command {
	icaQueryCmd := &cobra.Command{
		Use:                        "interchain-accounts",
		Aliases:                    []string{"ica"},
		Short:                      "interchain-accounts subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
	}

	icaQueryCmd.AddCommand(
		controllercli.GetQueryCmd(),
		hostcli.GetQueryCmd(),
	)

	return icaQueryCmd
}
