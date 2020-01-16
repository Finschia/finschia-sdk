package cli

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/line/link/client"
	"github.com/line/link/x/proxy/types"
	"github.com/spf13/cobra"
	"strconv"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Proxy transaction subcommands",
		Aliases:                    []string{"px"},
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		ProxyApproveCoinsCmd(cdc),
		ProxyDisapproveCoinsCmd(cdc),
		ProxySendCoinsFromCmd(cdc),
	)
	return txCmd
}

func ProxyApproveCoinsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve [proxy] [on_behalf_of] [denom] [amount]",
		Short: "Approve [proxy] to send [amount] [denom] coins on behalf of [on_behalf_of]",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[1]).WithCodec(cdc)

			proxy, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			onBehalfOf := cliCtx.FromAddress
			denom := args[2]
			amount, err := strconv.ParseInt(args[3], 10, 0)
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgProxyApproveCoins(proxy, onBehalfOf, denom, sdk.NewInt(amount))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func ProxyDisapproveCoinsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disapprove [proxy] [on_behalf_of] [denom] [amount]",
		Short: "Disapprove [proxy] to send [amount] [denom] coins on behalf of [on_behalf_of]",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[1]).WithCodec(cdc)

			proxy, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			onBehalfOf := cliCtx.FromAddress
			denom := args[2]
			amount, err := strconv.ParseInt(args[3], 10, 0)
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgProxyDisapproveCoins(proxy, onBehalfOf, denom, sdk.NewInt(amount))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func ProxySendCoinsFromCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-coins-from [proxy] [on_behalf_of] [to] [denom] [amount]",
		Short: "Send [amount] [denom] coins to [to] on behalf of [on_behalf_of] by [proxy]",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			proxy := cliCtx.FromAddress
			onBehalfOf, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			to, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}
			denom := args[3]
			amount, err := strconv.ParseInt(args[4], 10, 0)
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgProxySendCoinsFrom(proxy, onBehalfOf, to, denom, sdk.NewInt(amount))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}
