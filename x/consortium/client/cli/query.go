package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/flags"
	"github.com/line/lbm-sdk/x/consortium/types"
	sdk "github.com/line/lbm-sdk/types"
)

// GetQueryCmd returns the parent command for all x/consortium CLi query commands.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the consortium module",
	}

	cmd.AddCommand(
		GetEnabledCmd(),
		GetAllowedValidatorsCmd(),
		GetAllowedValidatorCmd(),
	)

	return cmd
}

// GetEnabledCmd returns the query consortium status command.
func GetEnabledCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enabled",
		Short: "get consortium status",
		Long:  "Gets the current status of consortium",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			params := types.QueryEnabledRequest{}
			res, err := queryClient.Enabled(context.Background(), &params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetAllowedValidatorsCmd returns allowed validators by consortium
func GetAllowedValidatorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "allowed-validators",
		Short: "allowed validators",
		Long: "Gets allowed validators by consortium",
		Args: cobra.NoArgs,
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

			params := types.QueryAllowedValidatorsRequest{Pagination: pageReq}
			res, err := queryClient.AllowedValidators(context.Background(), &params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "allowed validators")

	return cmd
}

// GetAllowedValidatorCmd returns allowed validators by consortium
func GetAllowedValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "allowed-validator [validator-address]",
		Short: "allowed validator",
		Long: "Gets allowed validator by consortium",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			valAddr := args[0]
			if err = sdk.ValidateValAddress(valAddr); err != nil {
				return err
			}

			params := types.QueryAllowedValidatorRequest{ValidatorAddress: valAddr}
			res, err := queryClient.AllowedValidator(context.Background(), &params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
