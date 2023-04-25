package rpc

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Finschia/ostracon/libs/bytes"
	"github.com/Finschia/ostracon/p2p"
	ctypes "github.com/Finschia/ostracon/rpc/core/types"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/client/flags"
	cryptocodec "github.com/Finschia/finschia-sdk/crypto/codec"
	cryptotypes "github.com/Finschia/finschia-sdk/crypto/types"
)

// ValidatorInfo is info about the node's validator, same as Ostracon,
// except that we use our own PubKey.
type validatorInfo struct {
	Address     bytes.HexBytes
	PubKey      cryptotypes.PubKey
	VotingPower int64
}

// ResultStatus is node's info, same as Ostracon, except that we use our own
// PubKey.
type resultStatus struct {
	NodeInfo      p2p.DefaultNodeInfo
	SyncInfo      ctypes.SyncInfo
	ValidatorInfo validatorInfo
}

// StatusCommand returns the command to return the status of the network.
func StatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Query remote node for status",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			status, err := getNodeStatus(clientCtx)
			if err != nil {
				return err
			}

			// `status` has OC pubkeys, we need to convert them to our pubkeys.
			pk, err := cryptocodec.FromOcPubKeyInterface(status.ValidatorInfo.PubKey)
			if err != nil {
				return err
			}
			statusWithPk := resultStatus{
				NodeInfo: status.NodeInfo,
				SyncInfo: status.SyncInfo,
				ValidatorInfo: validatorInfo{
					Address:     status.ValidatorInfo.Address,
					PubKey:      pk,
					VotingPower: status.ValidatorInfo.VotingPower,
				},
			}

			output, err := clientCtx.LegacyAmino.MarshalJSON(statusWithPk)
			if err != nil {
				return err
			}

			cmd.Println(string(output))
			return nil
		},
	}

	cmd.Flags().StringP(flags.FlagNode, "n", "tcp://localhost:26657", "Node to connect to")

	return cmd
}

func getNodeStatus(clientCtx client.Context) (*ctypes.ResultStatus, error) {
	node, err := clientCtx.GetNode()
	if err != nil {
		return &ctypes.ResultStatus{}, err
	}

	return node.Status(context.Background())
}
