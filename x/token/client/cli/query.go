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
		GetSupplyCmd(cdc),
		GetPermsCmd(cdc),
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
			token, height, err := retriever.GetToken(cliCtx, symbol)
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

func GetSupplyCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "supply [symbol]",
		Short: "Query supply of token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			symbol := args[0]

			supply, height, err := retriever.GetSupply(cliCtx, symbol)
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
