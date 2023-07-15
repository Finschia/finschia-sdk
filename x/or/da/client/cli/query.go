package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/client/flags"
	"github.com/Finschia/finschia-sdk/x/or/da/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		CmdQueryCCState(),
		CmdCCRef(),
		CmdCCRefs(),
		CmdQueueTxState(),
		CmdQueueTx(),
		CmdQueueTxs(),
		CmdMappedBatch(),
	)
	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "shows the parameters of the module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryCCState() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cc-state [rollup-name]",
		Short: "shows the state of the specific rollup's canonical chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CCState(cmd.Context(), &types.QueryCCStateRequest{
				RollupName: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdCCRef() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cc-ref [rollup-name] [batch-height]",
		Short: "shows the reference of the specific batch in the specific rollup's canonical chain",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			h, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			res, err := queryClient.CCRef(cmd.Context(), &types.QueryCCRefRequest{
				RollupName:  args[0],
				BatchHeight: h,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdCCRefs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cc-refs [rollup-name]",
		Short: "shows all references of the specific rollup's canonical chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.CCRefs(cmd.Context(), &types.QueryCCRefsRequest{
				RollupName: args[0],
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "canonical chain references")

	return cmd
}

func CmdQueueTxState() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "queue-tx-state [rollup-name]",
		Short: "shows the state of the specific rollup's tx queue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.QueueTxState(cmd.Context(), &types.QueryQueueTxStateRequest{
				RollupName: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueueTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "queue-tx [rollup-name] [tx-index]",
		Short: "shows the specific tx in the specific rollup's tx queue",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			idx, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.QueueTx(cmd.Context(), &types.QueryQueueTxRequest{
				RollupName: args[0],
				QueueIndex: idx,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueueTxs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "queue-txs [rollup-name]",
		Short: "shows all txs in the specific rollup's tx queue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.QueueTxs(cmd.Context(), &types.QueryQueueTxsRequest{
				RollupName: args[0],
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "rollup queue transactions")

	return cmd
}

func CmdMappedBatch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mapped-batch [rollup-name] [rollup-height]",
		Short: "shows the specific batch reference which is mapped to the specific rollup height",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			l2h, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.MappedBatch(cmd.Context(), &types.QueryMappedBatchRequest{
				RollupName: args[0],
				L2Height:   l2h,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
