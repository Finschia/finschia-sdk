package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/client/flags"
	"github.com/Finschia/finschia-sdk/client/tx"
	sdk "github.com/Finschia/finschia-sdk/types"
	bank "github.com/Finschia/finschia-sdk/x/bank/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
	govcli "github.com/Finschia/finschia-sdk/x/gov/client/cli"
	gov "github.com/Finschia/finschia-sdk/x/gov/types"
)

const (
	FlagFromDenom           = "from-denom"
	FlagToDenom             = "to-denom"
	FlagAmountCapForToDenom = "to-coin-amount-cap"
	FlagSwapRate            = "swap-rate"
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
		CmdTxMsgSwap(),
		CmdTxMsgSwapAll(),
	)

	return cmd
}

func CmdTxMsgSwap() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap [from] [from_coin_amount] [to_denom]",
		Short: "swap amount of from-coin to to-coin",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			from := args[0]
			if err := cmd.Flags().Set(flags.FlagFrom, from); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			toDenom := args[2]

			msg := &types.MsgSwap{
				FromAddress:    from,
				FromCoinAmount: amount,
				ToDenom:        toDenom,
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
		Use:   "swap-all [from_address] [from_denom] [to_denom]",
		Short: "swap all the from-coin to to-coin",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			from := args[0]
			if err := cmd.Flags().Set(flags.FlagFrom, from); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			fromDenom := args[1]
			toDenom := args[2]
			msg := &types.MsgSwapAll{
				FromAddress: clientCtx.GetFromAddress().String(),
				FromDenom:   fromDenom,
				ToDenom:     toDenom,
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

// NewCmdMakeSwapProposal implements a command handler for submitting a swap init proposal transaction.
func NewCmdMakeSwapProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "make-swap [messages-json]",
		Args:  cobra.ExactArgs(1),
		Short: "todo",
		Long: `
Parameters:
    messages-json: messages in json format that will be executed if the proposal is accepted.

Example of the content of messages-json:

{
  "metadata": {
    "description": "the base coin of Finschia mainnet",
    "denom_units": [
      {
        "denom": "cony",
        "exponent": 0,
        "aliases": [
          "microfinschia"
        ]
      },
      {
        "denom": "finschia",
        "exponent": 6,
        "aliases": []
      }
    ],
    "base": "cony",
    "display": "finschia",
    "name": "FINSCHIA",
    "symbol": "FNSA"
  }
}
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			title, err := cmd.Flags().GetString(govcli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(govcli.FlagDescription)
			if err != nil {
				return err
			}

			fromDenom, err := cmd.Flags().GetString(FlagFromDenom)
			if err != nil {
				return err
			}

			toDenom, err := cmd.Flags().GetString(FlagToDenom)
			if err != nil {
				return err
			}
			amountCapStr, err := cmd.Flags().GetString(FlagAmountCapForToDenom)
			if err != nil {
				return err
			}
			amountCap, err := sdk.NewDecFromStr(amountCapStr)
			if err != nil {
				return err
			}

			swapRate, err := cmd.Flags().GetString(FlagSwapRate)
			if err != nil {
				return err
			}
			swapRateDec, err := sdk.NewDecFromStr(swapRate)
			if err != nil {
				return err
			}
			swap := types.Swap{
				FromDenom:           fromDenom,
				ToDenom:             toDenom,
				AmountCapForToDenom: amountCap.TruncateInt(),
				SwapRate:            swapRateDec,
			}

			toDenomMetadata, err := parseToDenomMetadata(args[0])
			if err != nil {
				return err
			}

			content := types.NewMakeSwapProposal(title, description, swap, toDenomMetadata)

			depositStr, err := cmd.Flags().GetString(govcli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			msg, err := gov.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(govcli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(govcli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(govcli.FlagDeposit, "", "deposit of proposal")
	cmd.Flags().String(FlagFromDenom, "", "cony")
	cmd.Flags().String(FlagToDenom, "", "PDT")
	cmd.Flags().String(FlagAmountCapForToDenom, "0", "tbd")
	cmd.Flags().String(FlagSwapRate, "0", "tbd")

	return cmd
}

func parseToDenomMetadata(jsonDenomMetadata string) (bank.Metadata, error) {
	type toDenomMeta struct {
		Metadata bank.Metadata `json:"metadata"`
	}
	denomMeta := toDenomMeta{}
	if err := json.Unmarshal([]byte(jsonDenomMetadata), &denomMeta); err != nil {
		return bank.Metadata{}, err
	}

	if err := denomMeta.Metadata.Validate(); err != nil {
		return bank.Metadata{}, err
	}

	return denomMeta.Metadata, nil
}
