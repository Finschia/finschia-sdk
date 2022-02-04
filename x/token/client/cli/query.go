package cli

import (
	// "fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/line/lbm-sdk/client"
	// "github.com/line/lbm-sdk/client/flags"
	// sdk "github.com/line/lbm-sdk/types"
	// "github.com/line/lbm-sdk/types/errors"
	// "github.com/line/lbm-sdk/version"
	"github.com/line/lbm-sdk/x/token"
)

// Flag names and values
const (
	FlagClassID = "class-id"
)

// NewQueryCmd returns the cli query commands for this module
func NewQueryCmd() *cobra.Command {
	nftQueryCmd := &cobra.Command{
		Use:                        token.ModuleName,
		Short:                      "Querying commands for the nft module",
		Long:                       "",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	nftQueryCmd.AddCommand(
	)
	panic("Not implemented")
	return nftQueryCmd
}
