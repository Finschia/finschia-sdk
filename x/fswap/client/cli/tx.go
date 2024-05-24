package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/client/flags"
	"github.com/Finschia/finschia-sdk/client/tx"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	bank "github.com/Finschia/finschia-sdk/x/bank/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
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
		CmdMsgSetSwap(),
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

// CmdMsgSetSwap implements a command handler for submitting a swap init proposal transaction.
func CmdMsgSetSwap() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-swap [authority] [metadata-json]",
		Args:  cobra.ExactArgs(2),
		Short: "Set a swap",
		Long: `
Parameters:
    metadata-json: messages in json format that will be executed if the proposal is accepted.

Example of the content of metadata-json:

{
  "metadata": {
    "description": "example of to-denom is finschia cony",
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
			if err := validateGenerateOnly(cmd); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
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

			amountCap, ok := sdk.NewIntFromString(amountCapStr)
			if !ok {
				return sdkerrors.ErrInvalidRequest.Wrapf("failed to parse %s %s", FlagAmountCapForToDenom, amountCap.String())
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
				AmountCapForToDenom: amountCap,
				SwapRate:            swapRateDec,
			}

			authority := args[0]
			toDenomMetadata, err := parseToDenomMetadata(args[1])
			if err != nil {
				return err
			}

			msg := types.MsgSetSwap{
				Authority:       authority,
				Swap:            swap,
				ToDenomMetadata: toDenomMetadata,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().String(FlagFromDenom, "", "set fromDenom string, ex) cony")
	cmd.Flags().String(FlagToDenom, "", "set toDenom string, ex) peb")
	cmd.Flags().String(FlagAmountCapForToDenom, "0", "set integer value for limit cap for the amount to swap to to-denom, ex 1000000000")
	cmd.Flags().String(FlagSwapRate, "0", "set swap rate for swap from fromDenom to toDenom, ex(rate for cony to peb)  148079656000000")

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func validateGenerateOnly(cmd *cobra.Command) error {
	generateOnly, err := cmd.Flags().GetBool(flags.FlagGenerateOnly)
	if err != nil {
		return err
	}
	if !generateOnly {
		return fmt.Errorf("you must use it with the flag --%s", flags.FlagGenerateOnly)
	}
	return nil
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
