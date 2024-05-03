package cli

import (
	"fmt"
	"github.com/Finschia/finschia-sdk/version"
	"github.com/spf13/cobra"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/client/flags"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

const (
	flagSequences = "sequences"
)

// NewQueryCmd returns the query commands for fbridge module
func NewQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the fbridge module",
	}

	cmd.AddCommand(
		NewQueryNextSeqSendCmd(),
		NewQuerySeqToBlocknumsCmd(),
	)

	return cmd
}

func NewQueryNextSeqSendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sending-next-seq",
		Short: "Query the next sequence number for sending",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := types.NewQueryClient(clientCtx)
			res, err := qc.NextSeqSend(cmd.Context(), &types.QueryNextSeqSendRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQuerySeqToBlocknumsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "seq-to-blocknums",
		Short:   "Query the block number for given sequence numbers",
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf("%s query %s sending seq-to-blocknums --sequences=1,2,3", version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := types.NewQueryClient(clientCtx)

			seqSlice, err := cmd.Flags().GetInt64Slice(flagSequences)
			if err != nil {
				return err
			}

			seqs := make([]uint64, len(seqSlice))
			for i, seq := range seqSlice {
				seqs[i] = uint64(seq)
			}

			res, err := qc.SeqToBlocknums(cmd.Context(), &types.QuerySeqToBlocknumsRequest{Seqs: seqs})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().Int64Slice(flagSequences, []int64{}, "comma separated list of bridge sequnece numbers")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
