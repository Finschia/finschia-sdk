package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/client/flags"
	"github.com/Finschia/finschia-sdk/client/tx"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/or/settlement/types"
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

	cmd.AddCommand(NewTxCmdStartChallenge())
	cmd.AddCommand(NewTxCmdNsectChallenge())
	cmd.AddCommand(NewTxCmdFinishChallenge())

	return cmd
}

func NewTxCmdStartChallenge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-challenge [from_address] [to_address] [rollup_name] [block_height] [step_count]",
		Short: "Start challenge.",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			blockHeight, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return err
			}

			stepCount, ok := sdk.NewIntFromString(args[4])
			if !ok {
				return types.ErrInvalidStepCount
			}

			msg := &types.MsgStartChallenge{
				From:        args[0],
				To:          args[1],
				RollupName:  args[2],
				BlockHeight: blockHeight,
				StepCount:   stepCount.Uint64(),
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewTxCmdNsectChallenge() *cobra.Command {
	return &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			panic("implement me")
		},
	}
}

func NewTxCmdFinishChallenge() *cobra.Command {
	return &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			panic("implement me")
		},
	}
}
