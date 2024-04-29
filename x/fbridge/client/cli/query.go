package cli

import (
	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

// NewQueryCmd returns the query commands for fbridge module
func NewQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the fbridge module",
	}

	cmd.AddCommand(
		NewQueryNextSeqSendCmd(),
	)

	return cmd
}

func NewQueryNextSeqSendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nextseq-send",
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
