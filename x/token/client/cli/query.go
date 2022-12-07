package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/flags"
	"github.com/line/lbm-sdk/version"
	"github.com/line/lbm-sdk/x/token"
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
		NewQueryCmdMinted(),
		NewQueryCmdBurnt(),
		NewQueryCmdTokenClass(),
		NewQueryCmdTokenClasses(),
		NewQueryCmdGranteeGrants(),
		NewQueryCmdApproved(),
		NewQueryCmdApprovers(),
	)

	return queryCmd
}

func NewQueryCmdBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token-balance [class-id] [address]",
		Args:    cobra.ExactArgs(2),
		Short:   "query for token balances by a given address",
		Example: fmt.Sprintf(`$ %s query %s token-balance <class-id> <address>`, version.AppName, token.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := token.NewQueryClient(clientCtx)
			res, err := queryClient.Balance(cmd.Context(), &token.QueryBalanceRequest{
				ContractId: args[0],
				Address:    args[1],
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
		Use:     "supply [class-id]",
		Args:    cobra.ExactArgs(1),
		Short:   "query the supply of tokens of the class",
		Example: fmt.Sprintf(`$ %s query %s supply <class-id>`, version.AppName, token.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := token.NewQueryClient(clientCtx)
			res, err := queryClient.Supply(cmd.Context(), &token.QuerySupplyRequest{
				ContractId: args[0],
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

func NewQueryCmdMinted() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "minted [class-id]",
		Args:    cobra.ExactArgs(1),
		Short:   "query the minted tokens of the class",
		Example: fmt.Sprintf(`$ %s query %s supply <class-id>`, version.AppName, token.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := token.NewQueryClient(clientCtx)
			res, err := queryClient.Minted(cmd.Context(), &token.QueryMintedRequest{
				ContractId: args[0],
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

func NewQueryCmdBurnt() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "burnt [class-id]",
		Args:    cobra.ExactArgs(1),
		Short:   "query the burnt tokens of the class",
		Example: fmt.Sprintf(`$ %s query %s supply <class-id>`, version.AppName, token.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := token.NewQueryClient(clientCtx)
			res, err := queryClient.Burnt(cmd.Context(), &token.QueryBurntRequest{
				ContractId: args[0],
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

func NewQueryCmdTokenClass() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token [contract-id]",
		Args:    cobra.ExactArgs(1),
		Short:   "query token metadata based on its id",
		Example: fmt.Sprintf(`$ %s query %s token <contract-id>`, version.AppName, token.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := token.NewQueryClient(clientCtx)
			res, err := queryClient.TokenClass(cmd.Context(), &token.QueryTokenClassRequest{
				ContractId: args[0],
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

func NewQueryCmdTokenClasses() *cobra.Command {
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
			res, err := queryClient.TokenClasses(cmd.Context(), &token.QueryTokenClassesRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "classes")
	return cmd
}

func NewQueryCmdGranteeGrants() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "grantee-grants [class-id] [grantee]",
		Args:    cobra.ExactArgs(2),
		Short:   "query grants on a given grantee",
		Example: fmt.Sprintf(`$ %s query %s grantee-grants <class-id> <grantee>`, version.AppName, token.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := token.NewQueryClient(clientCtx)
			res, err := queryClient.GranteeGrants(cmd.Context(), &token.QueryGranteeGrantsRequest{
				ContractId: args[0],
				Grantee:    args[1],
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

func NewQueryCmdApproved() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "approved [class-id] [operator] [holder]",
		Args:    cobra.ExactArgs(3),
		Short:   "query authorization on its operator and the token holder",
		Example: fmt.Sprintf(`$ %s query %s approved <class-id> <operator> <holder>`, version.AppName, token.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := token.NewQueryClient(clientCtx)
			res, err := queryClient.Approved(cmd.Context(), &token.QueryApprovedRequest{
				ContractId: args[0],
				Proxy:      args[1],
				Approver:   args[2],
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

func NewQueryCmdApprovers() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "approvers [class-id] [address]",
		Args:    cobra.ExactArgs(2),
		Short:   "query all authorizations on a given operator",
		Example: fmt.Sprintf(`$ %s query %s approvers <class-id> <address>`, version.AppName, token.ModuleName),
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
			res, err := queryClient.Approvers(cmd.Context(), &token.QueryApproversRequest{
				ContractId: args[0],
				Address:    args[1],
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "authorizations")
	return cmd
}
