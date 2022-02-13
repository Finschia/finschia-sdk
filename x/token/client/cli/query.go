package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/flags"
	"github.com/line/lbm-sdk/version"
	"github.com/line/lbm-sdk/x/token"
)

// Flag names and values
const (
	FlagClassID = "class-id"
)

// NewQueryCmd returns the cli query commands for this module
func NewQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        token.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", token.ModuleName),
		Long:                       "",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		NewQueryCmdBalance(),
		NewQueryCmdSupply(),
		NewQueryCmdToken(),
		NewQueryCmdTokens(),
		NewQueryCmdGrants(),
		NewQueryCmdApprove(),
		NewQueryCmdApproves(),
	)

	return queryCmd
}

func NewQueryCmdBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "balance [class-id] [address]",
		Args:    cobra.ExactArgs(2),
		Short:   "query for token balances by a given address",
		Example: fmt.Sprintf(`$ %s query %s balance <class-id> <address>`, version.AppName, token.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := token.NewQueryClient(clientCtx)
			res, err := queryClient.Balance(cmd.Context(), &token.QueryBalanceRequest{
				ClassId: args[0],
				Address: args[1],
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

func NewQueryCmdSupply() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "supply [class-id] [type]",
		Args:    cobra.ExactArgs(2),
		Short:   "query the supply of tokens of the class",
		Example: fmt.Sprintf(`$ %s query %s supply <class-id> <type>`, version.AppName, token.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := token.NewQueryClient(clientCtx)
			res, err := queryClient.Supply(cmd.Context(), &token.QuerySupplyRequest{
				ClassId: args[0],
				Type:    args[1],
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

func NewQueryCmdToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token [class-id]",
		Args:    cobra.ExactArgs(1),
		Short:   "query token metadata based on its id",
		Example: fmt.Sprintf(`$ %s query %s token <class-id>`, version.AppName, token.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := token.NewQueryClient(clientCtx)
			res, err := queryClient.Token(cmd.Context(), &token.QueryTokenRequest{
				ClassId: args[0],
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

func NewQueryCmdTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tokens",
		Args:    cobra.NoArgs,
		Short:   "query all token metadata",
		Example: fmt.Sprintf(`$ %s query %s tokens`, version.AppName, token.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := token.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			res, err := queryClient.Tokens(cmd.Context(), &token.QueryTokensRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "tokens")
	return cmd
}

func NewQueryCmdGrants() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "grant [class-id] [grantee]",
		Args:    cobra.ExactArgs(2),
		Short:   "query grants on a given grantee",
		Example: fmt.Sprintf(`$ %s query %s grant <class-id> <grantee>`, version.AppName, token.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := token.NewQueryClient(clientCtx)
			res, err := queryClient.Grants(cmd.Context(), &token.QueryGrantsRequest{
				ClassId: args[0],
				Grantee: args[1],
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

func NewQueryCmdApprove() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "approve [class-id] [proxy] [approver]",
		Args:    cobra.ExactArgs(3),
		Short:   "query approve on its proxy and approver",
		Example: fmt.Sprintf(`$ %s query %s approve <class-id> <proxy> <approver>`, version.AppName, token.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := token.NewQueryClient(clientCtx)
			res, err := queryClient.Approve(cmd.Context(), &token.QueryApproveRequest{
				ClassId:  args[0],
				Proxy:    args[1],
				Approver: args[2],
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

func NewQueryCmdApproves() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "approves [class-id] [proxy]",
		Args:    cobra.ExactArgs(2),
		Short:   "query all approves on a given proxy",
		Example: fmt.Sprintf(`$ %s query %s approves <class-id> <proxy>`, version.AppName, token.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := token.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			res, err := queryClient.Approves(cmd.Context(), &token.QueryApprovesRequest{
				ClassId:    args[0],
				Proxy:      args[1],
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "approves")
	return cmd
}
