package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/client/tx"
	authclient "github.com/Finschia/finschia-sdk/x/auth/client"
	"github.com/Finschia/finschia-sdk/x/zkauth/types"
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
		NewCmdExecution(),
	)

	return cmd
}

func NewCmdExecution() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execuion [tx_json_file] [zk_proof_json_file] [max_block_height]",
		Short: "execute zkauth tx",
		Long:  "execute zkauth tx",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			theTx, err := authclient.ReadTxFromFile(clientCtx, args[0])
			if err != nil {
				return err
			}

			var zkAuthInputs types.ZKAuthInputs
			bytes, err := os.ReadFile(args[1])
			if err != nil {
				return err
			}
			err = clientCtx.Codec.UnmarshalInterfaceJSON(bytes, &zkAuthInputs)
			if err != nil {
				return err
			}
			maxBlockHeight, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}

			zkAuthSign := types.ZKAuthSignature{
				ZkAuthInputs:   &zkAuthInputs,
				MaxBlockHeight: maxBlockHeight,
			}
			msg := types.NewMsgExecution(theTx.GetMsgs(), zkAuthSign)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}
