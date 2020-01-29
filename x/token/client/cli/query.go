package cli

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/client"
	clienttypes "github.com/line/link/x/token/client/internal/types"
	"github.com/line/link/x/token/internal/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the token module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		GetTokenCmd(cdc),
		GetTokensCmd(cdc),
		GetCollectionCmd(cdc),
		GetCollectionsCmd(cdc),
		GetSupplyCmd(cdc),
		GetPermsCmd(cdc),
		GetParentCmd(cdc),
		GetRootCmd(cdc),
		GetChildrenCmd(cdc),
	)

	return cmd
}

func GetTokenCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token [symbol] [token-id]",
		Short: "Query token with its symbol and token-id. token-id is optional",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			tokenGetter := clienttypes.NewTokenRetriever(cliCtx)

			symbol := args[0]

			tokenID := ""
			if len(args) == 2 {
				tokenID = args[1]
			}
			if err := tokenGetter.EnsureExists(cliCtx, symbol, tokenID); err != nil {
				return err
			}

			token, err := tokenGetter.GetToken(cliCtx, symbol, tokenID)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(token)
		},
	}

	return client.GetCommands(cmd)[0]
}

func GetTokensCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens",
		Short: "Query all tokens",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			tokenGetter := clienttypes.NewTokenRetriever(cliCtx)

			tokens, err := tokenGetter.GetAllTokens(cliCtx)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(tokens)
		},
	}

	return client.GetCommands(cmd)[0]
}

func GetCollectionCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection [symbol]",
		Short: "Query collection",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			collectionGetter := clienttypes.NewCollectionRetriever(cliCtx)

			symbol := args[0]
			if err := collectionGetter.EnsureExists(cliCtx, symbol); err != nil {
				return err
			}

			collection, err := collectionGetter.GetCollection(cliCtx, symbol)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(collection)
		},
	}

	return client.GetCommands(cmd)[0]
}

func GetCollectionsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collections",
		Short: "Query all collections",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			collectionGetter := clienttypes.NewCollectionRetriever(cliCtx)

			collections, err := collectionGetter.GetAllCollections(cliCtx)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(collections)
		},
	}

	return client.GetCommands(cmd)[0]
}

func GetSupplyCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "supply [symbol]",
		Short: "Query supply",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			supplyGetter := clienttypes.NewSupplyRetriever(cliCtx)

			symbol := args[0]
			if err := supplyGetter.EnsureExists(cliCtx, symbol); err != nil {
				return err
			}

			supply, err := supplyGetter.GetSupply(cliCtx, symbol)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(supply)
		},
	}

	return client.GetCommands(cmd)[0]
}

func GetPermsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "perm [addr]",
		Short: "Get Permission of the Account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			permGetter := clienttypes.NewAccountPermissionRetriever(cliCtx)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			pms, err := permGetter.GetAccountPermission(cliCtx, addr)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(pms)
		},
	}

	return client.GetCommands(cmd)[0]
}

func GetParentCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parent [symbol] [token-id]",
		Short: "Query parent token with symbol and token-id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			tokenGetter := clienttypes.NewTokenRetriever(cliCtx)

			symbol := args[0]
			tokenID := args[1]

			if err := tokenGetter.EnsureExists(cliCtx, symbol, tokenID); err != nil {
				return err
			}

			token, _, err := tokenGetter.GetParent(cliCtx, symbol, tokenID)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(token)
		},
	}

	return client.GetCommands(cmd)[0]
}

func GetRootCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "root [symbol] [token-id]",
		Short: "Query root token with symbol and token-id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			tokenGetter := clienttypes.NewTokenRetriever(cliCtx)

			symbol := args[0]
			tokenID := args[1]

			if err := tokenGetter.EnsureExists(cliCtx, symbol, tokenID); err != nil {
				return err
			}

			token, _, err := tokenGetter.GetRoot(cliCtx, symbol, tokenID)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(token)
		},
	}

	return client.GetCommands(cmd)[0]
}

func GetChildrenCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "children [symbol] [token-id]",
		Short: "Query children tokens with symbol and token-id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			tokenGetter := clienttypes.NewTokenRetriever(cliCtx)

			symbol := args[0]
			tokenID := args[1]

			if err := tokenGetter.EnsureExists(cliCtx, symbol, tokenID); err != nil {
				return err
			}

			tokens, _, err := tokenGetter.GetChildren(cliCtx, symbol, tokenID)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(tokens)
		},
	}

	return client.GetCommands(cmd)[0]
}
