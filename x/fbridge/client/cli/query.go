package cli

import (
	"github.com/spf13/cobra"

	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func NewQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the fbridge module",
	}

	cmd.AddCommand()

	return cmd
}
