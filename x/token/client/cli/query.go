package cli

import (
	"github.com/link-chain/link/x/token/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
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

	cmd.AddCommand(GetTokenCmd(cdc))

	return cmd
}

func GetTokenCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "symbol [symbol]",
		Short: "Query symbol",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			tokenGetter := types.NewTokenRetriever(cliCtx)

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

	return flags.GetCommands(cmd)[0]
}
