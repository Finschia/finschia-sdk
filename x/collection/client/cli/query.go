package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/flags"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/version"
	"github.com/line/lbm-sdk/x/collection"
)

const (
	FlagTokenID = "token-id"
)

// NewQueryCmd returns the cli query commands for this module
func NewQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        collection.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", collection.ModuleName),
		Long:                       "",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		NewQueryCmdBalances(),
		NewQueryCmdFTSupply(),
		NewQueryCmdFTMinted(),
		NewQueryCmdFTBurnt(),
		NewQueryCmdNFTSupply(),
		NewQueryCmdNFTMinted(),
		NewQueryCmdNFTBurnt(),
		NewQueryCmdContract(),
		NewQueryCmdToken(),
		// NewQueryCmdTokensWithTokenType(),
		NewQueryCmdTokens(),
		NewQueryCmdTokenType(),
		NewQueryCmdTokenTypes(),
		NewQueryCmdRoot(),
		NewQueryCmdParent(),
		NewQueryCmdChildren(),
		NewQueryCmdGranteeGrants(),
		NewQueryCmdIsOperatorFor(),
		NewQueryCmdHoldersByOperator(),
	)

	return queryCmd
}

func NewQueryCmdBalances() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "balances [contract-id] [address]",
		Args:    cobra.ExactArgs(2),
		Short:   "query for token balances by a given address",
		Example: fmt.Sprintf(`$ %s query %s balances [contract-id] [address]`, version.AppName, collection.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			contractID := args[0]
			if err := collection.ValidateContractID(contractID); err != nil {
				return err
			}

			address := args[1]
			if _, err := sdk.AccAddressFromBech32(address); err != nil {
				return err
			}

			tokenID, err := cmd.Flags().GetString(FlagTokenID)
			if err != nil {
				return err
			}

			queryClient := collection.NewQueryClient(clientCtx)
			if len(tokenID) == 0 {
				pageReq, err := client.ReadPageRequest(cmd.Flags())
				if err != nil {
					return err
				}

				req := &collection.QueryAllBalancesRequest{
					ContractId: contractID,
					Address:    address,
					Pagination: pageReq,
				}
				res, err := queryClient.AllBalances(cmd.Context(), req)
				if err != nil {
					return err
				}
				return clientCtx.PrintProto(res)
			}

			if err := collection.ValidateTokenID(tokenID); err != nil {
				return err
			}

			req := &collection.QueryBalanceRequest{
				ContractId: contractID,
				Address:    address,
				TokenId:    tokenID,
			}
			res, err := queryClient.Balance(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().String(FlagTokenID, "", "Token ID to query for")
	flags.AddPaginationFlagsToCmd(cmd, "all balances")

	return cmd
}

func NewQueryCmdFTSupply() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ft-supply [contract-id] [token-id]",
		Args:    cobra.ExactArgs(2),
		Short:   "query the supply of tokens",
		Example: fmt.Sprintf(`$ %s ft-query %s supply [contract-id] [token-id]`, version.AppName, collection.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			contractID := args[0]
			if err := collection.ValidateContractID(contractID); err != nil {
				return err
			}

			tokenID := args[1]
			if err := collection.ValidateTokenID(tokenID); err != nil {
				return err
			}

			queryClient := collection.NewQueryClient(clientCtx)
			req := &collection.QueryFTSupplyRequest{
				ContractId: contractID,
				TokenId:    tokenID,
			}
			res, err := queryClient.FTSupply(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryCmdFTMinted() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ft-minted [contract-id] [token-id]",
		Args:    cobra.ExactArgs(2),
		Short:   "query the minted tokens",
		Example: fmt.Sprintf(`$ %s query %s ft-minted [contract-id] [token-id]`, version.AppName, collection.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			contractID := args[0]
			if err := collection.ValidateContractID(contractID); err != nil {
				return err
			}

			tokenID := args[1]
			if err := collection.ValidateTokenID(tokenID); err != nil {
				return err
			}

			queryClient := collection.NewQueryClient(clientCtx)
			req := &collection.QueryFTMintedRequest{
				ContractId: contractID,
				TokenId:    tokenID,
			}
			res, err := queryClient.FTMinted(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryCmdFTBurnt() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ft-burnt [contract-id] [token-id]",
		Args:    cobra.ExactArgs(2),
		Short:   "query the burnt tokens",
		Example: fmt.Sprintf(`$ %s query %s ft-burnt [contract-id] [token-id]`, version.AppName, collection.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			contractID := args[0]
			if err := collection.ValidateContractID(contractID); err != nil {
				return err
			}

			tokenID := args[1]
			if err := collection.ValidateTokenID(tokenID); err != nil {
				return err
			}

			queryClient := collection.NewQueryClient(clientCtx)
			req := &collection.QueryFTBurntRequest{
				ContractId: contractID,
				TokenId:    tokenID,
			}
			res, err := queryClient.FTBurnt(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryCmdNFTSupply() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "nft-supply [contract-id] [token-type]",
		Args:    cobra.ExactArgs(2),
		Short:   "query the supply of tokens",
		Example: fmt.Sprintf(`$ %s query %s nft-supply [contract-id] [token-type]`, version.AppName, collection.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			contractID := args[0]
			if err := collection.ValidateContractID(contractID); err != nil {
				return err
			}

			tokenType := args[1]
			if err := collection.ValidateClassID(tokenType); err != nil {
				return err
			}

			queryClient := collection.NewQueryClient(clientCtx)
			req := &collection.QueryNFTSupplyRequest{
				ContractId: contractID,
				TokenType:  tokenType,
			}
			res, err := queryClient.NFTSupply(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryCmdNFTMinted() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "nft-minted [contract-id] [token-type]",
		Args:    cobra.ExactArgs(2),
		Short:   "query the minted tokens of the class",
		Example: fmt.Sprintf(`$ %s query %s nft-minted [contract-id] [token-type]`, version.AppName, collection.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			contractID := args[0]
			if err := collection.ValidateContractID(contractID); err != nil {
				return err
			}

			tokenType := args[1]
			if err := collection.ValidateClassID(tokenType); err != nil {
				return err
			}

			queryClient := collection.NewQueryClient(clientCtx)
			req := &collection.QueryNFTMintedRequest{
				ContractId: contractID,
				TokenType:  tokenType,
			}
			res, err := queryClient.NFTMinted(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryCmdNFTBurnt() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "nft-burnt [contract-id] [token-type]",
		Args:    cobra.ExactArgs(2),
		Short:   "query the burnt tokens of the class",
		Example: fmt.Sprintf(`$ %s query %s nft-burnt [contract-id] [token-type]`, version.AppName, collection.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			contractID := args[0]
			if err := collection.ValidateContractID(contractID); err != nil {
				return err
			}

			tokenType := args[1]
			if err := collection.ValidateClassID(tokenType); err != nil {
				return err
			}

			queryClient := collection.NewQueryClient(clientCtx)
			req := &collection.QueryNFTBurntRequest{
				ContractId: contractID,
				TokenType:  tokenType,
			}
			res, err := queryClient.NFTBurnt(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryCmdContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "contract [contract-id]",
		Args:    cobra.ExactArgs(1),
		Short:   "query token metadata based on its id",
		Example: fmt.Sprintf(`$ %s query %s contract [contract-id]`, version.AppName, collection.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			contractID := args[0]
			if err := collection.ValidateContractID(contractID); err != nil {
				return err
			}

			queryClient := collection.NewQueryClient(clientCtx)
			req := &collection.QueryContractRequest{
				ContractId: contractID,
			}
			res, err := queryClient.Contract(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryCmdTokenType() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token-type [contract-id] [token-type]",
		Args:    cobra.ExactArgs(2),
		Short:   "query token type",
		Example: fmt.Sprintf(`$ %s query %s token-type [contract-id] [token-type]`, version.AppName, collection.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			contractID := args[0]
			if err := collection.ValidateContractID(contractID); err != nil {
				return err
			}

			classID := args[1]
			if err := collection.ValidateClassID(classID); err != nil {
				return err
			}

			queryClient := collection.NewQueryClient(clientCtx)
			req := &collection.QueryTokenTypeRequest{
				ContractId: contractID,
				TokenType:  classID,
			}
			res, err := queryClient.TokenType(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryCmdTokenTypes() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token-types [contract-id]",
		Args:    cobra.ExactArgs(1),
		Short:   "query all token types metadata",
		Example: fmt.Sprintf(`$ %s query %s token-types [contract-id]`, version.AppName, collection.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			contractID := args[0]
			if err := collection.ValidateContractID(contractID); err != nil {
				return err
			}

			queryClient := collection.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			req := &collection.QueryTokenTypesRequest{
				ContractId: contractID,
				Pagination: pageReq,
			}
			res, err := queryClient.TokenTypes(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "tokentypes")
	return cmd
}

func NewQueryCmdToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token [contract-id] [token-id]",
		Args:    cobra.ExactArgs(2),
		Short:   "query token metadata",
		Example: fmt.Sprintf(`$ %s query %s token [contract-id] [token-id]`, version.AppName, collection.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			contractID := args[0]
			if err := collection.ValidateContractID(contractID); err != nil {
				return err
			}

			tokenID := args[1]
			if err := collection.ValidateTokenID(tokenID); err != nil {
				return err
			}

			queryClient := collection.NewQueryClient(clientCtx)
			req := &collection.QueryTokenRequest{
				ContractId: contractID,
				TokenId:    tokenID,
			}
			res, err := queryClient.Token(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryCmdTokensWithTokenType() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tokens-with-token-type [contract-id] [token-type]",
		Args:    cobra.ExactArgs(2),
		Short:   "query all tokens with token type",
		Example: fmt.Sprintf(`$ %s query %s tokens [contract-id] [token-type]`, version.AppName, collection.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			contractID := args[0]
			if err := collection.ValidateContractID(contractID); err != nil {
				return err
			}

			classID := args[1]
			if err := collection.ValidateClassID(classID); err != nil {
				return err
			}

			queryClient := collection.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			req := &collection.QueryTokensWithTokenTypeRequest{
				ContractId: contractID,
				TokenType:  classID,
				Pagination: pageReq,
			}
			res, err := queryClient.TokensWithTokenType(cmd.Context(), req)
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

func NewQueryCmdTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tokens [contract-id]",
		Args:    cobra.ExactArgs(1),
		Short:   "query all token metadata",
		Example: fmt.Sprintf(`$ %s query %s tokens [contract-id]`, version.AppName, collection.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			contractID := args[0]
			if err := collection.ValidateContractID(contractID); err != nil {
				return err
			}

			queryClient := collection.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			req := &collection.QueryTokensRequest{
				ContractId: contractID,
				Pagination: pageReq,
			}
			res, err := queryClient.Tokens(cmd.Context(), req)
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

func NewQueryCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "root [contract-id] [token-id]",
		Args:    cobra.ExactArgs(2),
		Short:   "query root of an nft",
		Example: fmt.Sprintf(`$ %s query %s root [contract-id] [token-id]`, version.AppName, collection.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			contractID := args[0]
			if err := collection.ValidateContractID(contractID); err != nil {
				return err
			}

			tokenID := args[1]
			if err := collection.ValidateNFTID(tokenID); err != nil {
				return err
			}

			queryClient := collection.NewQueryClient(clientCtx)
			req := &collection.QueryRootRequest{
				ContractId: contractID,
				TokenId:    tokenID,
			}
			res, err := queryClient.Root(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryCmdParent() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "parent [contract-id] [token-id]",
		Args:    cobra.ExactArgs(2),
		Short:   "query parent of an nft",
		Example: fmt.Sprintf(`$ %s query %s parent [contract-id] [token-id]`, version.AppName, collection.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			contractID := args[0]
			if err := collection.ValidateContractID(contractID); err != nil {
				return err
			}

			tokenID := args[1]
			if err := collection.ValidateNFTID(tokenID); err != nil {
				return err
			}

			queryClient := collection.NewQueryClient(clientCtx)
			req := &collection.QueryParentRequest{
				ContractId: contractID,
				TokenId:    tokenID,
			}
			res, err := queryClient.Parent(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryCmdChildren() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "children [contract-id] [token-id]",
		Args:    cobra.ExactArgs(2),
		Short:   "query children of an nft",
		Example: fmt.Sprintf(`$ %s query %s children [contract-id] [token-id]`, version.AppName, collection.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			contractID := args[0]
			if err := collection.ValidateContractID(contractID); err != nil {
				return err
			}

			tokenID := args[1]
			if err := collection.ValidateNFTID(tokenID); err != nil {
				return err
			}

			queryClient := collection.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			req := &collection.QueryChildrenRequest{
				ContractId: contractID,
				TokenId:    tokenID,
				Pagination: pageReq,
			}
			res, err := queryClient.Children(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "children")
	return cmd
}

func NewQueryCmdGranteeGrants() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "grantee-grants [contract-id] [grantee]",
		Args:    cobra.ExactArgs(2),
		Short:   "query grants on a given grantee",
		Example: fmt.Sprintf(`$ %s query %s grantee-grants [contract-id] [grantee]`, version.AppName, collection.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			contractID := args[0]
			if err := collection.ValidateContractID(contractID); err != nil {
				return err
			}

			grantee := args[1]
			if _, err := sdk.AccAddressFromBech32(grantee); err != nil {
				return err
			}

			queryClient := collection.NewQueryClient(clientCtx)
			req := &collection.QueryGranteeGrantsRequest{
				ContractId: contractID,
				Grantee:    grantee,
			}
			res, err := queryClient.GranteeGrants(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryCmdIsOperatorFor() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "approved [contract-id] [operator] [holder]",
		Args:    cobra.ExactArgs(3),
		Short:   "query authorization on its operator and the token holder",
		Example: fmt.Sprintf(`$ %s query %s approved [contract-id] [operator] [holder]`, version.AppName, collection.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			contractID := args[0]
			if err := collection.ValidateContractID(contractID); err != nil {
				return err
			}

			operator := args[1]
			if _, err := sdk.AccAddressFromBech32(operator); err != nil {
				return err
			}

			holder := args[2]
			if _, err := sdk.AccAddressFromBech32(holder); err != nil {
				return err
			}

			queryClient := collection.NewQueryClient(clientCtx)
			req := &collection.QueryIsOperatorForRequest{
				ContractId: contractID,
				Operator:   operator,
				Holder:     holder,
			}
			res, err := queryClient.IsOperatorFor(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryCmdHoldersByOperator() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "approvers [contract-id] [operator]",
		Args:    cobra.ExactArgs(2),
		Short:   "query all authorizations on a given operator",
		Example: fmt.Sprintf(`$ %s query %s approvers [contract-id] [operator]`, version.AppName, collection.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			contractID := args[0]
			if err := collection.ValidateContractID(contractID); err != nil {
				return err
			}

			operator := args[1]
			if _, err := sdk.AccAddressFromBech32(operator); err != nil {
				return err
			}

			queryClient := collection.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			req := &collection.QueryHoldersByOperatorRequest{
				ContractId: contractID,
				Operator:   operator,
				Pagination: pageReq,
			}
			res, err := queryClient.HoldersByOperator(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "approvers")
	return cmd
}
