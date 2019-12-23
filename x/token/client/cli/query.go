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
	)

	return cmd
}

func GetTokenCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token [symbol]",
		Short: "Query token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			tokenGetter := clienttypes.NewTokenRetriever(cliCtx)

			symbol := args[0]
			if err := tokenGetter.EnsureExists(symbol); err != nil {
				return err
			}

			token, err := tokenGetter.GetToken(symbol)
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

			tokens, err := tokenGetter.GetAllTokens()
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
			if err := collectionGetter.EnsureExists(symbol); err != nil {
				return err
			}

			collection, err := collectionGetter.GetCollection(symbol)
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

			collections, err := collectionGetter.GetAllCollections()
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
			if err := supplyGetter.EnsureExists(symbol); err != nil {
				return err
			}

			supply, err := supplyGetter.GetSupply(symbol)
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
			pms, err := permGetter.GetAccountPermission(addr)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(pms)
		},
	}

	return client.GetCommands(cmd)[0]
}
