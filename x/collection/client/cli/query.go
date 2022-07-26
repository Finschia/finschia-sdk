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
		NewQueryCmdSupply(),
		NewQueryCmdMinted(),
		NewQueryCmdBurnt(),
		NewQueryCmdContract(),
		NewQueryCmdContracts(),
		NewQueryCmdNFT(),
		// NewQueryCmdNFTs(),
		NewQueryCmdOwner(),
		NewQueryCmdRoot(),
		NewQueryCmdParent(),
		NewQueryCmdChildren(),
		NewQueryCmdGrant(),
		NewQueryCmdGranteeGrants(),
		NewQueryCmdAuthorization(),
		NewQueryCmdOperatorAuthorizations(),
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

func NewQueryCmdSupply() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "supply [contract-id] [class-id]",
		Args:    cobra.ExactArgs(2),
		Short:   "query the supply of tokens of the class",
		Example: fmt.Sprintf(`$ %s query %s supply [contract-id] [class-id]`, version.AppName, collection.ModuleName),
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
			req := &collection.QuerySupplyRequest{
				ContractId: contractID,
				ClassId:    classID,
			}
			res, err := queryClient.Supply(cmd.Context(), req)
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
		Use:     "minted [contract-id] [class-id]",
		Args:    cobra.ExactArgs(2),
		Short:   "query the minted tokens of the class",
		Example: fmt.Sprintf(`$ %s query %s minted [contract-id] [class-id]`, version.AppName, collection.ModuleName),
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
			req := &collection.QueryMintedRequest{
				ContractId: contractID,
				ClassId:    classID,
			}
			res, err := queryClient.Minted(cmd.Context(), req)
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
		Use:     "burnt [contract-id] [class-id]",
		Args:    cobra.ExactArgs(2),
		Short:   "query the burnt tokens of the class",
		Example: fmt.Sprintf(`$ %s query %s burnt [contract-id] [class-id]`, version.AppName, collection.ModuleName),
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
			req := &collection.QueryBurntRequest{
				ContractId: contractID,
				ClassId:    classID,
			}
			res, err := queryClient.Burnt(cmd.Context(), req)
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

func NewQueryCmdContracts() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "contracts",
		Args:    cobra.NoArgs,
		Short:   "query all contract metadata",
		Example: fmt.Sprintf(`$ %s query %s contracts`, version.AppName, collection.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := collection.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			req := &collection.QueryContractsRequest{
				Pagination: pageReq,
			}
			res, err := queryClient.Contracts(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "contracts")
	return cmd
}

func NewQueryCmdFTClass() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ft-class [contract-id] [class-id]",
		Args:    cobra.ExactArgs(2),
		Short:   "query ft class metadata based on its id",
		Example: fmt.Sprintf(`$ %s query %s ft-class [contract-id] [class-id]`, version.AppName, collection.ModuleName),
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
			req := &collection.QueryFTClassRequest{
				ContractId: contractID,
				ClassId:    classID,
			}
			res, err := queryClient.FTClass(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryCmdNFTClass() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "nft-class [contract-id] [class-id]",
		Args:    cobra.ExactArgs(2),
		Short:   "query nft class metadata based on its id",
		Example: fmt.Sprintf(`$ %s query %s nft-class [contract-id] [class-id]`, version.AppName, collection.ModuleName),
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
			req := &collection.QueryNFTClassRequest{
				ContractId: contractID,
				ClassId:    classID,
			}
			res, err := queryClient.NFTClass(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryCmdTokenClassTypeName() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token-class-type-name [contract-id] [class-id]",
		Args:    cobra.ExactArgs(2),
		Short:   "query token class type name based on its id",
		Example: fmt.Sprintf(`$ %s query %s token-class-type-name [contract-id] [class-id]`, version.AppName, collection.ModuleName),
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
			req := &collection.QueryTokenClassTypeNameRequest{
				ContractId: contractID,
				ClassId:    classID,
			}
			res, err := queryClient.TokenClassTypeName(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// func NewQueryCmdTokenClasses() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:     "classes [contract-id]",
// 		Args:    cobra.ExactArgs(1),
// 		Short:   "query all token class metadata",
// 		Example: fmt.Sprintf(`$ %s query %s classes [contract-id]`, version.AppName, collection.ModuleName),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			clientCtx, err := client.GetClientQueryContext(cmd)
// 			if err != nil {
// 				return err
// 			}

// 			contractID := args[0]
// 			if err := collection.ValidateContractID(contractID); err != nil {
// 				return err
// 			}

// 			queryClient := collection.NewQueryClient(clientCtx)
// 			pageReq, err := client.ReadPageRequest(cmd.Flags())
// 			if err != nil {
// 				return err
// 			}
// 			req := &collection.QueryTokenClassesRequest{
// 				ContractId: contractID,
// 				Pagination: pageReq,
// 			}
// 			res, err := queryClient.TokenClasses(cmd.Context(), req)
// 			if err != nil {
// 				return err
// 			}
// 			return clientCtx.PrintProto(res)
// 		},
// 	}

// 	flags.AddQueryFlagsToCmd(cmd)
// 	flags.AddPaginationFlagsToCmd(cmd, "classes")
// 	return cmd
// }

func NewQueryCmdNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "nft [contract-id] [token-id]",
		Args:    cobra.ExactArgs(2),
		Short:   "query nft metadata based on its id",
		Example: fmt.Sprintf(`$ %s query %s nft [contract-id] [token-id]`, version.AppName, collection.ModuleName),
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
			req := &collection.QueryNFTRequest{
				ContractId: contractID,
				TokenId:    tokenID,
			}
			res, err := queryClient.NFT(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// func NewQueryCmdNFTs() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:     "nfts [contract-id]",
// 		Args:    cobra.ExactArgs(1),
// 		Short:   "query all nft metadata",
// 		Example: fmt.Sprintf(`$ %s query %s nfts [contract-id]`, version.AppName, collection.ModuleName),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			clientCtx, err := client.GetClientQueryContext(cmd)
// 			if err != nil {
// 				return err
// 			}

// 			contractID := args[0]
// 			if err := collection.ValidateContractID(contractID); err != nil {
// 				return err
// 			}

// 			queryClient := collection.NewQueryClient(clientCtx)
// 			pageReq, err := client.ReadPageRequest(cmd.Flags())
// 			if err != nil {
// 				return err
// 			}
// 			req := &collection.QueryNFTsRequest{
// 				ContractId: contractID,
// 				Pagination: pageReq,
// 			}
// 			res, err := queryClient.NFTs(cmd.Context(), req)
// 			if err != nil {
// 				return err
// 			}
// 			return clientCtx.PrintProto(res)
// 		},
// 	}

// 	flags.AddQueryFlagsToCmd(cmd)
// 	flags.AddPaginationFlagsToCmd(cmd, "nfts")
// 	return cmd
// }

func NewQueryCmdOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "owner [contract-id] [token-id]",
		Args:    cobra.ExactArgs(2),
		Short:   "query owner of an nft",
		Example: fmt.Sprintf(`$ %s query %s owner [contract-id] [token-id]`, version.AppName, collection.ModuleName),
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
			req := &collection.QueryOwnerRequest{
				ContractId: contractID,
				TokenId:    tokenID,
			}
			res, err := queryClient.Owner(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
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

func NewQueryCmdGrant() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "grant [contract-id] [grantee] [permission]",
		Args:    cobra.ExactArgs(3),
		Short:   "query a permission on a given grantee",
		Example: fmt.Sprintf(`$ %s query %s grant [contract-id] [grantee] [permission]`, version.AppName, collection.ModuleName),
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

			permission := collection.Permission(collection.Permission_value[args[2]])

			queryClient := collection.NewQueryClient(clientCtx)
			req := &collection.QueryGrantRequest{
				ContractId: contractID,
				Grantee:    grantee,
				Permission: permission,
			}
			res, err := queryClient.Grant(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
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

func NewQueryCmdAuthorization() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "authorization [contract-id] [operator] [holder]",
		Args:    cobra.ExactArgs(3),
		Short:   "query authorization on its operator and the token holder",
		Example: fmt.Sprintf(`$ %s query %s authorization [contract-id] [operator] [holder]`, version.AppName, collection.ModuleName),
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
			req := &collection.QueryAuthorizationRequest{
				ContractId: contractID,
				Operator:   operator,
				Holder:     holder,
			}
			res, err := queryClient.Authorization(cmd.Context(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryCmdOperatorAuthorizations() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "operator-authorizations [contract-id] [operator]",
		Args:    cobra.ExactArgs(2),
		Short:   "query all authorizations on a given operator",
		Example: fmt.Sprintf(`$ %s query %s operator-authorizations [contract-id] [operator]`, version.AppName, collection.ModuleName),
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
			req := &collection.QueryOperatorAuthorizationsRequest{
				ContractId: contractID,
				Operator:   operator,
				Pagination: pageReq,
			}
			res, err := queryClient.OperatorAuthorizations(cmd.Context(), req)
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
