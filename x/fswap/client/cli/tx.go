package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/client/flags"
	"github.com/Finschia/finschia-sdk/client/tx"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
	"github.com/Finschia/finschia-sdk/x/gov/client/cli"
	gov "github.com/Finschia/finschia-sdk/x/gov/types"
)

const (
	FlagFromDenom   = "from-denom"
	FlagToDenom     = "to-denom"
	FlagAmountLimit = "amount-limit"
	FlagSwapRate    = "swap-rate"
)

var DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())

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
		CmdTxMsgSwap(),
		CmdTxMsgSwapAll(),
	)

	return cmd
}

func CmdTxMsgSwap() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap [address] [fnsa_amount]",
		Short: "swap amounts of old coin to new coin. Note, the'--from' flag is ignored as it is implied from [address].",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Flags().Set(flags.FlagFrom, args[0]); err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			coin, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			msg := &types.MsgSwapRequest{
				FromAddress: clientCtx.GetFromAddress().String(),
				Amount:      coin,
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

func CmdTxMsgSwapAll() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap-all [address]",
		Short: "swap all the old coins. Note, the'--from' flag is ignored as it is implied from [address].",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Flags().Set(flags.FlagFrom, args[0]); err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgSwapAllRequest{
				FromAddress: clientCtx.GetFromAddress().String(),
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

// NewCmdFswapInitProposal implements a command handler for submitting a fswap init proposal transaction.
func NewCmdFswapInitProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fswap-init [flags]",
		Args:  cobra.ExactArgs(0),
		Short: "todo",
		Long:  "todo",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			// need confirm parse with flags or parse with args
			contents, err := parseArgsToContent(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			msg, err := gov.NewMsgSubmitProposal(contents, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	cmd.Flags().String(FlagFromDenom, "", "cony")
	cmd.Flags().String(FlagToDenom, "", "PDT")
	cmd.Flags().Int64(FlagAmountLimit, 0, "tbd")
	cmd.Flags().Int64(FlagSwapRate, 0, "tbd")

	return cmd
}

func parseArgsToContent(cmd *cobra.Command) (gov.Content, error) {
	title, err := cmd.Flags().GetString(cli.FlagTitle)
	if err != nil {
		return nil, err
	}

	description, err := cmd.Flags().GetString(cli.FlagDescription)
	if err != nil {
		return nil, err
	}

	from_denom, err := cmd.Flags().GetString(FlagFromDenom)
	if err != nil {
		return nil, err
	}

	to_denom, err := cmd.Flags().GetString(FlagToDenom)
	if err != nil {
		return nil, err
	}

	amount_limit, err := cmd.Flags().GetInt64(FlagAmountLimit)
	if err != nil {
		return nil, err
	}

	swap_rate, err := cmd.Flags().GetInt64(FlagSwapRate)
	if err != nil {
		return nil, err
	}

	fswapInit := types.FswapInit{
		FromDenom:           from_denom,
		ToDenom:             to_denom,
		AmountCapForToDenom: sdk.NewInt(amount_limit),
		SwapMultiple:        sdk.NewInt(swap_rate),
	}

	content := types.NewFswapInitProposal(title, description, fswapInit)
	return content, nil
}
