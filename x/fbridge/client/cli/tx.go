package cli

import (
	"github.com/spf13/cobra"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

// GetTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
	TxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "fbridge transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	TxCmd.AddCommand()

	return TxCmd
}
