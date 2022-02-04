package cli

import (
	// "fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/line/lbm-sdk/client"
	// "github.com/line/lbm-sdk/client/tx"
	// sdk "github.com/line/lbm-sdk/types"
	// "github.com/line/lbm-sdk/version"
	"github.com/line/lbm-sdk/x/token"
)

// NewTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
	nftTxCmd := &cobra.Command{
		Use:                        token.ModuleName,
		Short:                      "token transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	nftTxCmd.AddCommand(
	)
	panic("Not implemented")
	return nftTxCmd
}
