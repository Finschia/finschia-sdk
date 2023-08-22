package cli

import (
	"context"
	"fmt"

	"github.com/Finschia/finschia-rdk/client"
	"github.com/Finschia/finschia-rdk/client/flags"
	"github.com/Finschia/finschia-rdk/x/or/rollup/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group rollup queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(NewQueryCmdShowRollup())
	cmd.AddCommand(NewQueryCmdShowAllRollup())
	cmd.AddCommand(NewQueryCmdShowSequencer())
	cmd.AddCommand(NewQueryCmdShowSequencersByRollup())

	return cmd
}

func NewQueryCmdShowRollup() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-rollup [rollup-name]",
		Short: "shows a rollup",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argRollupName := args[0]

			params := &types.QueryRollupRequest{
				RollupName: argRollupName,
			}

			res, err := queryClient.Rollup(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func NewQueryCmdShowSequencersByRollup() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-sequencers-by-rollup [rollup-name]",
		Short: "show sequencers by rollup",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argRollupName := args[0]

			params := &types.QuerySequencersByRollupRequest{
				RollupName: argRollupName,
			}

			res, err := queryClient.SequencersByRollup(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func NewQueryCmdShowAllRollup() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "shows a all rollup",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryAllRollupRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.AllRollup(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func NewQueryCmdShowSequencer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-sequencer [sequencer-address]",
		Short: "shows a sequencer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argSequencerAddress := args[0]

			params := &types.QuerySequencerRequest{
				SequencerAddress: argSequencerAddress,
			}

			res, err := queryClient.Sequencer(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
