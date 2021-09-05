package cli

import (
	"github.com/spf13/cobra"

	"github.com/line/lfb-sdk/x/ibc/light-clients/99-ostracon/types"
)

// NewTxCmd returns a root CLI command handler for all x/ibc/light-clients/99-ostracon transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.SubModuleName,
		Short:                      "Tendermint client transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
	}

	txCmd.AddCommand(
		NewCreateClientCmd(),
		NewUpdateClientCmd(),
		NewSubmitMisbehaviourCmd(),
	)

	return txCmd
}
