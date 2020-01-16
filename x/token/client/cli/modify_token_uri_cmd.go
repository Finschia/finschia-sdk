package cli

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/line/link/x/token/internal/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/line/link/client"
)

func ModifyTokenURICmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify-token-uri [owner_address] [symbol] [token_uri] [token_id]",
		Short: "Modify token_uri of token ",
		Args:  cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)
			symbol := args[1]
			tokenURI := args[2]
			tokenID := ""
			if len(args) > 3 {
				tokenID = args[3]
			}

			msg := types.NewMsgModifyTokenURI(cliCtx.FromAddress, symbol, tokenURI, tokenID)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}
