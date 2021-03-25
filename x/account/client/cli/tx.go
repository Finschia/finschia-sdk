package cli

import (
	"bufio"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/line/link-modules/x/account/client/utils"
	"github.com/line/link-modules/x/account/internal/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Account commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		CreateAccountCmd(cdc),
		EmptyCmd(cdc),
	)
	return txCmd
}

func CreateAccountCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-account [from_key_or_address] [target_address]",
		Short: "Create an account having target_address",
		Args:  cobra.ExactArgs(2),
		RunE:  makeCreateAccountCmd(cdc),
	}

	return flags.PostCommands(cmd)[0]
}

func EmptyCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "empty [from_key_or_address]",
		Short: "Do nothing",
		Args:  cobra.ExactArgs(1),
		RunE:  makeEmptyCmd(cdc),
	}

	return flags.PostCommands(cmd)[0]
}

func makeCreateAccountCmd(cdc *codec.Codec) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		inBuf := bufio.NewReader(cmd.InOrStdin())
		txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
		cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

		target, err := sdk.AccAddressFromBech32(args[1])
		if err != nil {
			return err
		}

		// build and sign the transaction, then broadcast to Tendermint
		msg := types.NewMsgCreateAccount(cliCtx.GetFromAddress(), target)
		return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
	}
}

func makeEmptyCmd(cdc *codec.Codec) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		inBuf := bufio.NewReader(cmd.InOrStdin())
		txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
		cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

		msg := types.NewMsgEmpty(cliCtx.GetFromAddress())
		return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
	}
}
