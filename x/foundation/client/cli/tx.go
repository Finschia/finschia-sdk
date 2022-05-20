package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/flags"
	"github.com/line/lbm-sdk/client/tx"
	"github.com/line/lbm-sdk/codec"
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

	FlagExec = "exec"
	ExecTry  = "try"
)

func parseMembers(codec codec.Codec, membersJSON string) ([]foundation.Member, error) {
	var cliMembers []json.RawMessage
	if err := json.Unmarshal([]byte(membersJSON), &cliMembers); err != nil {
		return nil, err
	}

	members := make([]foundation.Member, len(cliMembers))
	for i, cliMember := range cliMembers {
		var member foundation.Member
		if err := codec.UnmarshalJSON(cliMember, &member); err != nil {
			return nil, err
		}
		members[i] = member
	}

	return members, nil
}

func parseAddresses(addressesJSON string) ([]string, error) {
	var addresses []string
	if err := json.Unmarshal([]byte(addressesJSON), &addresses); err != nil {
		return nil, err
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("you must provide one address at least")
	}

	return addresses, nil
}

func parseDecisionPolicy(codec codec.Codec, policyJSON string) (foundation.DecisionPolicy, error) {
	var policy foundation.DecisionPolicy
	if err := codec.UnmarshalInterfaceJSON([]byte(policyJSON), &policy); err != nil {
		return nil, err
	}

	return policy, nil
}

func execFromString(execStr string) foundation.Exec {
	exec := foundation.Exec_EXEC_UNSPECIFIED
	switch execStr {
	case ExecTry:
		exec = foundation.Exec_EXEC_TRY
	}
	return exec
}

// VoteOptionFromString returns a VoteOption from a string. It returns an error
// if the string is invalid.
func voteOptionFromString(str string) (foundation.VoteOption, error) {
	vo, ok := foundation.VoteOption_value[str]
	if !ok {
		return foundation.VOTE_OPTION_UNSPECIFIED, fmt.Errorf("'%s' is not a valid vote option", str)
	}
	return foundation.VoteOption(vo), nil
}

func parseMsgs(cdc codec.Codec, msgsJSON string) ([]sdk.Msg, error) {
	var cliMsgs []json.RawMessage
	if err := json.Unmarshal([]byte(msgsJSON), &cliMsgs); err != nil {
		return nil, err
	}

	msgs := make([]sdk.Msg, len(cliMsgs))
	for i, anyJSON := range cliMsgs {
		var msg sdk.Msg
		err := cdc.UnmarshalInterfaceJSON(anyJSON, &msg)
		if err != nil {
			return nil, err
		}

		msgs[i] = msg
	}

	return msgs, nil
}

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
		NewTxCmdUpdateMembers(),
		NewTxCmdUpdateDecisionPolicy(),
		NewTxCmdSubmitProposal(),
		NewTxCmdWithdrawProposal(),
		NewTxCmdVote(),
		NewTxCmdExec(),
		NewTxCmdLeaveFoundation(),
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

			createAuths := func(addings, deletings []string) []foundation.ValidatorAuth {
				var auths []foundation.ValidatorAuth
				for _, addr := range addings {
					auth := foundation.ValidatorAuth{
						OperatorAddress: addr,
						CreationAllowed: true,
					}
					auths = append(auths, auth)
				}
				for _, addr := range deletings {
					auth := foundation.ValidatorAuth{
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
		Short: "Fund the treasury",
		Long: `Fund the treasury
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			from := args[0]
			if err := cmd.Flags().Set(flags.FlagFrom, from); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}

			msg := foundation.MsgFundTreasury{
				From:   from,
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
		Use:   "withdraw-from-treasury [operator] [to] [amount]",
		Args:  cobra.ExactArgs(3),
		Short: "Withdraw coins from the treasury",
		Long: `Withdraw coins from the treasury
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			operator := args[0]
			if err := cmd.Flags().Set(flags.FlagFrom, operator); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return err
			}

			msg := foundation.MsgWithdrawFromTreasury{
				Operator: operator,
				To:       args[1],
				Amount:   amount,
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

func NewTxCmdUpdateMembers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-members [operator] [members-json]",
		Args:  cobra.ExactArgs(2),
		Short: "Update the foundation members",
		Long: `Update the foundation members

Example of the content of members-json:

[
  {
    "address": "addr1",
    "participating": true,
    "metadata": "some new metadata"
  },
  {
    "address": "addr2",
    "participating": false,
    "metadata": "some metadata"
  }
]

Set a member's participating to false to delete it.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			operator := args[0]
			if err := cmd.Flags().Set(flags.FlagFrom, operator); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			updates, err := parseMembers(clientCtx.Codec, args[1])
			if err != nil {
				return err
			}

			msg := foundation.MsgUpdateMembers{
				Operator:      operator,
				MemberUpdates: updates,
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

func NewTxCmdUpdateDecisionPolicy() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-decision-policy [operator] [policy-json]",
		Args:  cobra.ExactArgs(2),
		Short: "Update the foundation decision policy",
		Long: `Update the foundation decision policy

Example of the content of policy-json:

{
  "@type": "/lbm.foundation.v1.ThresholdDecisionPolicy",
  "threshold": "10",
  "windows": {
    "voting_period": "24h",
    "min_execution_period": "0s",
  }
}
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			operator := args[0]
			if err := cmd.Flags().Set(flags.FlagFrom, operator); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := foundation.MsgUpdateDecisionPolicy{
				Operator: operator,
			}
			policy, err := parseDecisionPolicy(clientCtx.Codec, args[1])
			if err != nil {
				return err
			}
			if err := msg.SetDecisionPolicy(policy); err != nil {
				return err
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

func NewTxCmdSubmitProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-proposal [metadata] [proposers-json] [messages-json]",
		Args:  cobra.ExactArgs(3),
		Short: "Submit a new proposal",
		Long: `Submit a new proposal

Parameters:
    metadata: metadata of the proposal.
    proposers-json: the addresses of the proposers in json format.
    messages-json: messages in json format that will be executed if the proposal is accepted.

Example of the content of proposers-json:

[
  "addr1",
  "addr2"
]

Example of the content of messages-json:

[
  {
    "@type": "/lbm.foundation.v1.MsgWithdrawFromTreasury",
    "operator": "addr1",
    "to": "addr2",
    "amount": "10000stake"
  }
]
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			proposers, err := parseAddresses(args[1])
			if err != nil {
				return err
			}

			signer := proposers[0]
			if err := cmd.Flags().Set(flags.FlagFrom, signer); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			messages, err := parseMsgs(clientCtx.Codec, args[2])
			if err != nil {
				return err
			}

			execStr, err := cmd.Flags().GetString(FlagExec)
			if err != nil {
				return err
			}
			exec := execFromString(execStr)

			msg := foundation.MsgSubmitProposal{
				Proposers: proposers,
				Metadata:  args[0],
				Exec:      exec,
			}
			if err := msg.SetMsgs(messages); err != nil {
				return err
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(FlagExec, "", "Set to 1 to try to execute proposal immediately after creation")

	return cmd
}

func NewTxCmdWithdrawProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-proposal [proposal-id] [address]",
		Args:  cobra.ExactArgs(2),
		Short: "Withdraw a submitted proposal",
		Long: `Withdraw a submitted proposal.

Parameters:
    proposal-id: unique ID of the proposal.
    address: one of the proposer of the proposal.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			address := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, address); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := foundation.MsgWithdrawProposal{
				ProposalId: proposalID,
				Address:    address,
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

func NewTxCmdVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote [proposal-id] [voter] [option] [metadata]",
		Args:  cobra.ExactArgs(4),
		Short: "Vote on a proposal",
		Long: `Vote on a proposal.

Parameters:
    proposal-id: unique ID of the proposal
    voter: voter account addresses.
    vote-option: choice of the voter(s)
        VOTE_OPTION_UNSPECIFIED: no-op
        VOTE_OPTION_NO: no
        VOTE_OPTION_YES: yes
        VOTE_OPTION_ABSTAIN: abstain
        VOTE_OPTION_NO_WITH_VETO: no-with-veto
    metadata: metadata for the vote
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			voter := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, voter); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			option, err := voteOptionFromString(args[2])
			if err != nil {
				return err
			}

			execStr, err := cmd.Flags().GetString(FlagExec)
			if err != nil {
				return err
			}
			exec := execFromString(execStr)

			msg := foundation.MsgVote{
				ProposalId: proposalID,
				Voter:      voter,
				Option:     option,
				Metadata:   args[3],
				Exec:       exec,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(FlagExec, "", "Set to 1 to try to execute proposal immediately after voting")

	return cmd
}

func NewTxCmdExec() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec [proposal-id] [signer]",
		Args:  cobra.ExactArgs(2),
		Short: "Execute a proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			signer := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, signer); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := foundation.MsgExec{
				ProposalId: proposalID,
				Signer:     signer,
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

func NewTxCmdLeaveFoundation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "leave-foundation [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Leave the foundation",
		RunE: func(cmd *cobra.Command, args []string) error {
			address := args[0]
			if err := cmd.Flags().Set(flags.FlagFrom, address); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := foundation.MsgLeaveFoundation{
				Address: address,
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
