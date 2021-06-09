package cli

import (
	"github.com/line/lfb-sdk/client"
	"github.com/line/lfb-sdk/client/flags"
	"github.com/line/lfb-sdk/client/tx"
	"github.com/line/lfb-sdk/x/auth/types"
	"github.com/spf13/cobra"
)

// NewTxCmd returns a root CLI command handler for all x/auth transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Auth transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(NewEmptyTxCmd())

	return txCmd
}

// NewEmptyTxCmd returns a CLI command handler for creating a MsgEmpty transaction.
func NewEmptyTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "empty [from_key_or_address]",
		Short: `Empty doesn't do anything. Used to measure performance.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.Flags().Set(flags.FlagFrom, args[0])
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewServiceMsgEmpty(clientCtx.GetFromAddress())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
