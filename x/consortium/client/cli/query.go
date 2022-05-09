package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/flags"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/consortium"
)

// NewQueryCmd returns the parent command for all x/consortium CLi query commands.
func NewQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   consortium.ModuleName,
		Short: "Querying commands for the consortium module",
	}

	cmd.AddCommand(
		NewQueryCmdParams(),
		NewQueryCmdValidatorAuth(),
		NewQueryCmdValidatorAuths(),
	)

	return cmd
}

// NewQueryCmdParams returns the query consortium parameters command.
func NewQueryCmdParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query consortium params",
		Long:  "Gets the current parameters of consortium",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := consortium.NewQueryClient(clientCtx)

			params := consortium.QueryParamsRequest{}
			res, err := queryClient.Params(context.Background(), &params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryCmdValidatorAuth returns validator authorization by consortium
func NewQueryCmdValidatorAuth() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-auth [validator-address]",
		Short: "Query validator authorization",
		Long:  "Gets validator authorization by consortium",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := consortium.NewQueryClient(clientCtx)

			valAddr := args[0]
			if err = sdk.ValidateValAddress(valAddr); err != nil {
				return err
			}

			params := consortium.QueryValidatorAuthRequest{ValidatorAddress: valAddr}
			res, err := queryClient.ValidatorAuth(context.Background(), &params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryCmdValidatorAuths returns validator authorizations by consortium
func NewQueryCmdValidatorAuths() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-auths",
		Short: "Query validator authorizations",
		Long:  "Gets validator authorizations by consortium",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := consortium.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := consortium.QueryValidatorAuthsRequest{Pagination: pageReq}
			res, err := queryClient.ValidatorAuths(context.Background(), &params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "validator auths")

	return cmd
}
