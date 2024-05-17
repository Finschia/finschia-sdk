package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/client/flags"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group fswap queries under a subcommand
	cmd := &cobra.Command{
		Use:                        queryRoute,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQuerySwapped(),
		CmdQueryTotalSwappableAmount(),
		CmdQuerySwaps(),
	)
	return cmd
}

func CmdQuerySwapped() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swapped [from_denom] [to_denom]",
		Short: "shows the current swap status, including both old and new coin amount",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			req := &types.QuerySwappedRequest{
				FromDenom: args[0],
				ToDenom:   args[1],
			}
			res, err := queryClient.Swapped(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryTotalSwappableAmount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-swappable-amount [from_denom] [to_denom]",
		Short: "shows the current total amount of new coin that're swappable",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryTotalSwappableToCoinAmountRequest{
				FromDenom: args[0],
				ToDenom:   args[1],
			}
			res, err := queryClient.TotalSwappableToCoinAmount(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQuerySwaps() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swaps",
		Short: "shows the all the swaps that proposed",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QuerySwapsRequest{
				Pagination: pageReq,
			}
			res, err := queryClient.Swaps(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
