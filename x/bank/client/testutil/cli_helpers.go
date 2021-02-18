package testutil

import (
	"context"
	"fmt"

	"github.com/line/ostracon/libs/cli"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/tx"
	"github.com/line/lbm-sdk/testutil"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	"github.com/line/lbm-sdk/testutil/testdata"
	sdk "github.com/line/lbm-sdk/types"
	txtypes "github.com/line/lbm-sdk/types/tx"
	"github.com/line/lbm-sdk/types/tx/signing"
	authclient "github.com/line/lbm-sdk/x/auth/client"
	bankcli "github.com/line/lbm-sdk/x/bank/client/cli"
	"github.com/line/lbm-sdk/x/bank/types"
)

func MsgSendExec(clientCtx client.Context, from, to, amount fmt.Stringer, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{from.String(), to.String(), amount.String()}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, bankcli.NewSendTxCmd(), args)
}

func QueryBalancesExec(clientCtx client.Context, address fmt.Stringer, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{address.String(), fmt.Sprintf("--%s=json", cli.OutputFlag)}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, bankcli.GetBalancesCmd(), args)
}

// LegacyGRPCProtoMsgSend is a legacy method to broadcast a legacy proto MsgSend.
//
// Deprecated.
//nolint:interfacer
func LegacyGRPCProtoMsgSend(clientCtx client.Context, keyName string, from, to sdk.Address, fee, amount []sdk.Coin, extraArgs ...string) (*txtypes.BroadcastTxResponse, error) {
	// prepare txBuilder with msg
	txBuilder := clientCtx.TxConfig.NewTxBuilder()
	feeAmount := fee
	gasLimit := testdata.NewTestGasLimit()

	// This sets a legacy Proto MsgSend.
	err := txBuilder.SetMsgs(&types.MsgSend{
		FromAddress: from.String(),
		ToAddress:   to.String(),
		Amount:      amount,
	})
	if err != nil {
		return nil, err
	}

	txBuilder.SetFeeAmount(feeAmount)
	txBuilder.SetGasLimit(gasLimit)

	// setup txFactory
	txFactory := tx.Factory{}.
		WithChainID(clientCtx.ChainID).
		WithKeybase(clientCtx.Keyring).
		WithTxConfig(clientCtx.TxConfig).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)

	// Sign Tx.
	err = authclient.SignTx(txFactory, clientCtx, keyName, txBuilder, false, true)
	if err != nil {
		return nil, err
	}

	txBytes, err := clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	// Broadcast the tx via gRPC.
	queryClient := txtypes.NewServiceClient(clientCtx)

	return queryClient.BroadcastTx(
		context.Background(),
		&txtypes.BroadcastTxRequest{
			Mode:    txtypes.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: txBytes,
		},
	)
}
