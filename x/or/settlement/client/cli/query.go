package cli

import (
	"context"
	"fmt"
	"strconv"

	// "strings"

	"github.com/spf13/cobra"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/client/flags"
	"github.com/Finschia/finschia-sdk/version"

	// "github.com/Finschia/finschia-sdk/client/flags"
	// sdk "github.com/Finschia/finschia-sdk/types"

	"github.com/Finschia/finschia-sdk/x/or/settlement/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group settlement queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(NewQueryCmdChallenge())

	return cmd
}

func NewQueryCmdChallenge() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "challenge",
		Short:   "shows the challenge of the module",
		Args:    cobra.ExactArgs(1),
		Example: fmt.Sprintf(`$ %s query %s challenge <challenge-id>`, version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			ci, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			req := &types.QueryChallengeRequest{
				ChallengeId: ci,
			}
			res, err := queryClient.Challenge(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
