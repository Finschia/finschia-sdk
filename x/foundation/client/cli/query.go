package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/flags"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

// NewQueryCmd returns the parent command for all x/foundation CLi query commands.
func NewQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   foundation.ModuleName,
		Short: "Querying commands for the foundation module",
	}

	cmd.AddCommand(
		NewQueryCmdParams(),
		NewQueryCmdValidatorAuth(),
		NewQueryCmdValidatorAuths(),
		NewQueryCmdTreasury(),
	)

	return cmd
}

// NewQueryCmdParams returns the query foundation parameters command.
func NewQueryCmdParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query foundation params",
		Long:  "Gets the current parameters of foundation",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := foundation.NewQueryClient(clientCtx)

			params := foundation.QueryParamsRequest{}
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

// NewQueryCmdValidatorAuth returns validator authorization by foundation
func NewQueryCmdValidatorAuth() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-auth [validator-address]",
		Short: "Query validator authorization",
		Long:  "Gets validator authorization by foundation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := foundation.NewQueryClient(clientCtx)

			valAddr := args[0]
			if err = sdk.ValidateValAddress(valAddr); err != nil {
				return err
			}

			params := foundation.QueryValidatorAuthRequest{ValidatorAddress: valAddr}
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

// NewQueryCmdValidatorAuths returns validator authorizations by foundation
func NewQueryCmdValidatorAuths() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-auths",
		Short: "Query validator authorizations",
		Long:  "Gets validator authorizations by foundation",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := foundation.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := foundation.QueryValidatorAuthsRequest{Pagination: pageReq}
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

// NewQueryCmdTreasury returns the amount of coins in the foundation treasury
func NewQueryCmdTreasury() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "treasury",
		Short: "Query foundation treasury",
		Long:  "Gets the amount of coins in the foundation treasury",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := foundation.NewQueryClient(clientCtx)

			req := foundation.QueryTreasuryRequest{}
			res, err := queryClient.Treasury(context.Background(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}


// NewQueryCmdFoundationInfo returns the information of the foundation.
func NewQueryCmdFoundationInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "foundation-info",
		Args:  cobra.NoArgs,
		Short: "Query the foundation information",
		Long: `Query the foundation information
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := foundation.NewQueryClient(clientCtx)

			req := foundation.QueryFoundationInfoRequest{}
			res, err := queryClient.FoundationInfo(context.Background(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
