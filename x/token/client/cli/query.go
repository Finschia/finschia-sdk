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
		GetCollectionTokenCmd(cdc),
		GetCollectionTokensCmd(cdc),
		GetCollectionCmd(cdc),
		GetCollectionsCmd(cdc),
		GetSupplyCmd(cdc),
		GetCollectionTokenSupplyCmd(cdc),
		GetCollectionTokenCountCmd(cdc),
		GetPermsCmd(cdc),
		GetParentCmd(cdc),
		GetRootCmd(cdc),
		GetChildrenCmd(cdc),
	)

	return cmd
}

func GetTokenCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token [symbol]",
		Short: "Query token with its symbol",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			symbol := args[0]
			token, height, err := retriever.GetToken(cliCtx, symbol, "")
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)

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
			retriever := clienttypes.NewRetriever(cliCtx)

			tokens, height, err := retriever.GetTokens(cliCtx)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)
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
			retriever := clienttypes.NewRetriever(cliCtx)

			symbol := args[0]
			collection, height, err := retriever.GetCollection(cliCtx, symbol)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)
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
			retriever := clienttypes.NewRetriever(cliCtx)

			collections, height, err := retriever.GetCollections(cliCtx)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)
			return cliCtx.PrintOutput(collections)
		},
	}

	return client.GetCommands(cmd)[0]
}
func GetCollectionTokenCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-collection [symbol] [token-id]",
		Short: "Query collection token with collection symbol and token's token-id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			symbol := args[0]
			tokenID := args[1]
			token, height, err := retriever.GetToken(cliCtx, symbol, tokenID)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)

			return cliCtx.PrintOutput(token)
		},
	}

	return client.GetCommands(cmd)[0]
}

func GetCollectionTokensCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens-collection [symbol]",
		Short: "Query all collection tokens with collection symbol",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			symbol := args[0]

			collection, height, err := retriever.GetCollection(cliCtx, symbol)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)
			return cliCtx.PrintOutput(collection.Tokens)
		},
	}

	return client.GetCommands(cmd)[0]
}

func GetCollectionTokenSupplyCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "supply-token-collection [symbol] [token-id]",
		Short: "Query supply of collection token with collection symbol and tokens's token-id.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			symbol := args[0]
			tokenID := args[1]

			supply, height, err := retriever.GetSupply(cliCtx, symbol, tokenID)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)
			return cliCtx.PrintOutput(supply)
		},
	}

	return client.GetCommands(cmd)[0]
}
func GetCollectionTokenCountCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "count-token-collection [symbol] [token-type]",
		Short: "Query count of collection tokens with collection symbol and the base-id.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			symbol := args[0]
			baseID := args[1]

			supply, height, err := retriever.GetCollectionNFTCount(cliCtx, symbol, baseID)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)
			return cliCtx.PrintOutput(supply)
		},
	}

	return client.GetCommands(cmd)[0]
}

func GetSupplyCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "supply-token [symbol]",
		Short: "Query supply of collection token with collection symbol and token's token-id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			symbol := args[0]

			supply, height, err := retriever.GetSupply(cliCtx, symbol, "")
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)
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
			retriever := clienttypes.NewRetriever(cliCtx)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			pms, height, err := retriever.GetAccountPermission(cliCtx, addr)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)
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
