package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/tx"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/gov/client/cli"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
	"github.com/line/lbm-sdk/version"
	"github.com/line/lbm-sdk/x/consortium/types"
)

// Proposal flags
const (
	FlagAllowedValidatorAdd    = "add"
	FlagAllowedValidatorRemove = "remove"
)

// NewCmdSubmitDisableConsortiumProposal implements the command to submit a disable-consortium proposal
func NewCmdSubmitDisableConsortiumProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable-consortium",
		Args:  cobra.NoArgs,
		Short: "Submit a disable consortium proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a disable consortium proposal.

Example:
$ %s tx gov submit-proposal disable-consortium [flags]
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

			content := types.NewDisableConsortiumProposal(title, description)

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

// NewCmdSubmitEditAllowedValidatorsProposal implements the command to submit a edit-allowed-validators proposal
func NewCmdSubmitEditAllowedValidatorsProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-allowed-validators",
		Args:  cobra.NoArgs,
		Short: "Submit a edit allowed validators proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a edit allowed validators proposal.

Example:
$ %s tx gov submit-proposal edit-allowed-validators [flags]
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
				} else {
					return strings.Split(concat, ",")
				}
			}

			addingValidatorsStr, err := cmd.Flags().GetString(FlagAllowedValidatorAdd)
			if err != nil {
				return err
			}
			addingValidators := parseCommaSeparated(addingValidatorsStr)

			removingValidatorsStr, err := cmd.Flags().GetString(FlagAllowedValidatorRemove)
			if err != nil {
				return err
			}
			removingValidators := parseCommaSeparated(removingValidatorsStr)

			content := types.NewEditAllowedValidatorsProposal(title, description, addingValidators, removingValidators)

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
	cmd.Flags().String(FlagAllowedValidatorRemove, "", "validator addresses to remove")

	return cmd
}
