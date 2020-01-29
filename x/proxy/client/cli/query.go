package cli

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/client"
	"github.com/line/link/x/proxy/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the proxy module",
		Aliases:                    []string{"px"},
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		GetAllowanceCmd(cdc),
	)

	return cmd
}

func GetAllowanceCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "allowance [proxy] [on_behalf_of] [denom]",
		Short: "Get allowance of [proxy] on behalf of [on_behalf_of] for [denom]",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			proxyAllowanceGetter := types.NewProxyAllowanceRetriever(cliCtx)

			proxy, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			onBehalfOf, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			denom := args[2]

			allowance, height, err := proxyAllowanceGetter.GetProxyAllowance(proxy, onBehalfOf, denom)
			if err != nil {
				return err
			}

			return cliCtx.WithHeight(height).PrintOutput(allowance)
		},
	}

	return client.GetCommands(cmd)[0]
}
