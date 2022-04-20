package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/flags"
	"github.com/line/lbm-sdk/client/tx"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/version"
	"github.com/line/lbm-sdk/x/foundation"
	"github.com/line/lbm-sdk/x/gov/client/cli"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

// Proposal flags
const (
	FlagAllowedValidatorAdd    = "add"
	FlagAllowedValidatorDelete = "delete"
)

func phraseIgnoreFromFlag(name string) string {
	return fmt.Sprintf(" note, the '--from' flag is ignored as it is implied from [%s].", name)
}

var phraseNotForUsers = " not intended to be used by users."

// NewTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        foundation.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", foundation.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewTxCmdFundTreasury(),
		NewTxCmdWithdrawFromTreasury(),
	)

	return txCmd
}

// NewProposalCmdUpdateFoundationParams implements the command to submit an update-foundation-params proposal
func NewProposalCmdUpdateFoundationParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-foundation-params",
		Args:  cobra.NoArgs,
		Short: "Submit an update foundation params proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit an update foundation params proposal.
For now, you have no other options, so we make the corresponding params json file for you.

Example:
$ %s tx gov submit-proposal update-foundation-params [flags]
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
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

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			params := &foundation.Params{
				Enabled: false,
			}
			content := foundation.NewUpdateFoundationParamsProposal(title, description, params)
			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")

	return cmd
}

// NewProposalCmdUpdateValidatorAuths implements the command to submit an update-validator-auths proposal
func NewProposalCmdUpdateValidatorAuths() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-validator-auths",
		Args:  cobra.NoArgs,
		Short: "Submit an update validator auths proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit an update validator auths proposal.

Example:
$ %s tx gov submit-proposal update-validator-auths [flags]
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
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

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			parseCommaSeparated := func(concat string) []string {
				if concat == "" {
					return []string{}
				}
				return strings.Split(concat, ",")
			}

			addingValidatorsStr, err := cmd.Flags().GetString(FlagAllowedValidatorAdd)
			if err != nil {
				return err
			}
			addingValidators := parseCommaSeparated(addingValidatorsStr)

			deletingValidatorsStr, err := cmd.Flags().GetString(FlagAllowedValidatorDelete)
			if err != nil {
				return err
			}
			deletingValidators := parseCommaSeparated(deletingValidatorsStr)

			createAuths := func(addings, deletings []string) []*foundation.ValidatorAuth {
				var auths []*foundation.ValidatorAuth
				for _, addr := range addings {
					auth := &foundation.ValidatorAuth{
						OperatorAddress: addr,
						CreationAllowed: true,
					}
					auths = append(auths, auth)
				}
				for _, addr := range deletings {
					auth := &foundation.ValidatorAuth{
						OperatorAddress: addr,
						CreationAllowed: false,
					}
					auths = append(auths, auth)
				}

				return auths
			}

			auths := createAuths(addingValidators, deletingValidators)
			content := foundation.NewUpdateValidatorAuthsProposal(title, description, auths)
			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")

	cmd.Flags().String(FlagAllowedValidatorAdd, "", "validator addresses to add")
	cmd.Flags().String(FlagAllowedValidatorDelete, "", "validator addresses to delete")

	return cmd
}

func NewTxCmdFundTreasury() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fund-treasury [from] [amount]",
		Args:  cobra.ExactArgs(2),
		Short: "fund the treasury." + phraseIgnoreFromFlag("from"),
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s fund-treasury <from> <amount>`, version.AppName, foundation.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			from := args[0]
			cmd.Flags().Set(flags.FlagFrom, from)

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}

			msg := foundation.MsgFundTreasury{
				From: from,
				Amount: amount,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewTxCmdWithdrawFromTreasury() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-from-treasury [to] [amount]",
		Args:  cobra.ExactArgs(2),
		Short: "withdraw coins from the treasury." + phraseNotForUsers,
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s withdraw-from-treasury <to> <amount>`, version.AppName, foundation.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}

			msg := foundation.MsgWithdrawFromTreasury{
				To: args[0],
				Amount: amount,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
