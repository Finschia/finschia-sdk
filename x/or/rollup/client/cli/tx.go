package cli

import (
	"encoding/json"
	"fmt"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/client/flags"
	"github.com/Finschia/finschia-sdk/client/tx"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/or/rollup/types"
	"github.com/spf13/cobra"

	cryptotypes "github.com/Finschia/finschia-sdk/crypto/types"
)

const (
	FlagAmount = "amount"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(NewCreateRollupCmd())
	cmd.AddCommand(NewRegisterSequencerCmd())
	cmd.AddCommand(NewDepositCmd())
	cmd.AddCommand(NewWithdrawCmd())
	// cmd.AddCommand(NewRemoveSequencerCmd())

	return cmd
}

func NewCreateRollupCmd() *cobra.Command {
	cmd := &cobra.Command{
		// TODO: If rollup:sequencer=1:N, add [max-sequencers] parameter
		Use:   "create-rollup [rollup_name] [permissioned-addresses]",
		Short: "Create a new rollup",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argRollupName := args[0]
			argPermissionedAddresses := new(types.Sequencers)
			if len(args) == 2 {
				err = json.Unmarshal([]byte(args[1]), argPermissionedAddresses)
				if err != nil {
					return err
				}
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateRollup(
				argRollupName,
				clientCtx.GetFromAddress().String(),
				argPermissionedAddresses,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewRegisterSequencerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-sequencer [rollup_name] [pubkey] [amount]",
		Short: "Register a new sequencer for a rollup. Amount is the value to deposit",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argRollupName := args[0]
			argPubkey := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var pk cryptotypes.PubKey
			if err := clientCtx.Codec.UnmarshalInterfaceJSON([]byte(argPubkey), &pk); err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			msg, err := types.NewMsgRegisterSequencer(
				pk,
				argRollupName,
				clientCtx.GetFromAddress().String(),
				amount,
			)
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(FlagAmount)

	return cmd
}

func NewDepositCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-by-sequencer [rollup_name] [amount]",
		Short: "Deposit by sequencer",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argRollupName := args[0]
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDeposit(argRollupName, clientCtx.GetFromAddress().String(), amount)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(FlagAmount)

	return cmd
}

func NewWithdrawCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-by-sequencer [rollup_name] [amount]",
		Short: "Withdraw by sequencer",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argRollupName := args[0]
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdraw(argRollupName, clientCtx.GetFromAddress().String(), amount)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(FlagAmount)

	return cmd
}

// func NewRemoveSequencerCmd() *cobra.Command {
// 	panic("implement me")
// }
