package cli

import (
	"github.com/Finschia/finschia-sdk/client/flags"
	"github.com/Finschia/finschia-sdk/client/tx"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/spf13/cobra"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

// NewTxCmd returns the transaction commands for fbridge module
func NewTxCmd() *cobra.Command {
	TxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "fbridge transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	TxCmd.AddCommand(
		NewTransferTxCmd(),
	)

	return TxCmd
}

func NewTransferTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer [to_address] [amount]",
		Short: `Transfer token from current chain to counterparty chain`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fromAddr := clientCtx.GetFromAddress().String()
			if _, err := sdk.AccAddressFromBech32(fromAddr); err != nil {
				return sdkerrors.ErrInvalidAddress.Wrapf("invalid address: %s", fromAddr)
			}
			toAddr := args[0]
			coins, err := sdk.ParseCoinsNormalized(args[1])
			if len(coins) != 1 {
				return sdkerrors.ErrInvalidRequest.Wrapf("only one native coin type is allowed")
			}

			msg := types.MsgTransfer{
				Sender:   fromAddr,
				Receiver: toAddr,
				Amount:   coins[0].Amount,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
