package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/client/flags"
	"github.com/Finschia/finschia-sdk/client/tx"
	"github.com/Finschia/finschia-sdk/x/or/da/types"
)

const (
	FlagCompression = "compression"
	FlagL2GasLimit  = "l2-gas-limit"
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

	cmd.AddCommand(
		CmdTxAppendCCBatch(),
		CmdTxEnqueue(),
	)

	return cmd
}

func CmdTxAppendCCBatch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "append-cc-batch [from_key_or_address] [rollup_name] [compressed_batch_bytes]",
		Short: "append a batch of rollup transactions",
		Long: `append a batch of rollup transactions.
Note, the '--from' flag is ignored as it is implied from [from_key_or_address].
When using '--dry-run' a key name cannot be used, only a bech32 address.`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.Flags().Set(flags.FlagFrom, args[0])
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			compressOpt, err := cmd.Flags().GetInt32(FlagCompression)
			if err != nil {
				return err
			}
			msg := &types.MsgAppendCCBatch{
				FromAddress: clientCtx.GetFromAddress().String(),
				RollupName:  args[1],
				Batch: types.CompressedCCBatch{
					Data:        []byte(args[2]),
					Compression: types.CompressionOption(compressOpt),
				},
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().Int32(FlagCompression, 0, "compression algorithm to use for the batch")
	return cmd
}

func CmdTxEnqueue() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enqueue [from_key_or_address] [rollup_name] [tx_bytes]",
		Short: "enqueue a rollup transaction on L1",
		Long: `enqueue a rollup transaction on L1.
Note, the '--from' flag is ignored as it is implied from [from_key_or_address].
When using '--dry-run' a key name cannot be used, only a bech32 address.
`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.Flags().Set(flags.FlagFrom, args[0])
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			l2gasLimit, err := cmd.Flags().GetUint64(FlagL2GasLimit)
			if err != nil {
				return err
			}

			msg := &types.MsgEnqueue{
				FromAddress: clientCtx.GetFromAddress().String(),
				RollupName:  args[1],
				GasLimit:    l2gasLimit,
				Txraw:       []byte(args[2]),
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().Uint64(FlagL2GasLimit, 0, "gas limit for the transaction on L2")
	err := cmd.MarkFlagRequired(FlagL2GasLimit)
	if err != nil {
		return nil
	}

	return cmd
}
