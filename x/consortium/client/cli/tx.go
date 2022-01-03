package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/tx"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/version"
	"github.com/line/lbm-sdk/x/consortium/types"
	"github.com/line/lbm-sdk/x/gov/client/cli"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

// Proposal flags
const (
	FlagAllowedValidatorAdd    = "add"
	FlagAllowedValidatorDelete = "delete"
)

// NewProposalCmdUpdateConsortiumParams implements the command to submit an update-consortium-params proposal
func NewProposalCmdUpdateConsortiumParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-consortium-params",
		Args:  cobra.NoArgs,
		Short: "Submit an update consortium params proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit an update consortium params proposal.
For now, you have no other options, so we make the corresponding params json file for you.

Example:
$ %s tx gov submit-proposal update-consortium-params [flags]
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

			params := &types.Params{
				Enabled: false,
			}
			content := types.NewUpdateConsortiumParamsProposal(title, description, params)
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

			createAuths := func(addings, deletings []string) []*types.ValidatorAuth {
				var auths []*types.ValidatorAuth
				for _, addr := range addings {
					auth := &types.ValidatorAuth{
						OperatorAddress: addr,
						CreationAllowed: true,
					}
					auths = append(auths, auth)
				}
				for _, addr := range deletings {
					auth := &types.ValidatorAuth{
						OperatorAddress: addr,
						CreationAllowed: false,
					}
					auths = append(auths, auth)
				}

				return auths
			}

			auths := createAuths(addingValidators, deletingValidators)
			content := types.NewUpdateValidatorAuthsProposal(title, description, auths)
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
