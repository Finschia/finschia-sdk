package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/client/flags"
	"github.com/Finschia/finschia-sdk/client/tx"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/version"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

// NewTxCmd returns the transaction commands for fbridge module
func NewTxCmd() *cobra.Command {
	TxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "fbridge transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	TxCmd.AddCommand(
		NewTransferTxCmd(),
		NewSuggestRoleTxCmd(),
		NewAddVoteForRoleTxCmd(),
		NewSetBridgeStatusTxCmd(),
	)

	return TxCmd
}

func NewTransferTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "transfer [to_address] [amount]",
		Short:   `Transfer token from current chain to counterparty chain`,
		Example: fmt.Sprintf("%s tx %s transfer link1... 1000cony --from mykey", version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fromAddr := clientCtx.GetFromAddress().String()
			if _, err := sdk.AccAddressFromBech32(fromAddr); err != nil {
				return sdkerrors.ErrInvalidAddress.Wrapf("invalid address: %s", fromAddr)
			}
			toAddr := args[0]
			coins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}
			if len(coins) != 1 {
				return sdkerrors.ErrInvalidRequest.Wrapf("only one native coin type is allowed")
			}

			msg := types.MsgTransfer{
				Sender:   fromAddr,
				Receiver: toAddr,
				Amount:   coins[0].Amount,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewSuggestRoleTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "suggest-role [target_address] [role]",
		Short:   `Suggest a role to a specific address (unspecified|guardian|operator|judge)`,
		Args:    cobra.ExactArgs(2),
		Example: fmt.Sprintf("%s tx %s suggest-role link1... guardian --from guardiankey", version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			from := clientCtx.GetFromAddress().String()
			if _, err := sdk.AccAddressFromBech32(from); err != nil {
				return sdkerrors.ErrInvalidAddress.Wrapf("invalid address: %s", from)
			}
			target := args[0]
			role, found := types.QueryParamToRole[args[1]]
			if !found {
				return sdkerrors.ErrInvalidRequest.Wrapf("invalid role: %s", args[1])
			}

			msg := types.MsgSuggestRole{
				From:   from,
				Target: target,
				Role:   role,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewAddVoteForRoleTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add-vote-for-role [proposal_id] [option]",
		Short:   `Vote for a role proposal (yes|no)`,
		Args:    cobra.ExactArgs(2),
		Example: fmt.Sprintf("%s tx %s add-vote-for-role 1 yes --from guardiankey", version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			from := clientCtx.GetFromAddress().String()
			if _, err := sdk.AccAddressFromBech32(from); err != nil {
				return sdkerrors.ErrInvalidAddress.Wrapf("invalid address: %s", from)
			}
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return sdkerrors.ErrInvalidRequest.Wrapf("invalid proposal ID: %s", args[0])
			}

			voteOpts := map[string]types.VoteOption{
				"yes": types.OptionYes,
				"no":  types.OptionNo,
			}
			option, found := voteOpts[args[1]]
			if !found {
				return sdkerrors.ErrInvalidRequest.Wrapf("invalid vote option: %s", args[1])
			}

			msg := types.MsgAddVoteForRole{
				From:       from,
				ProposalId: proposalID,
				Option:     option,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewSetBridgeStatusTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-bridge-status [status]",
		Short: `Set sender's bridge switch for halting/resuming the bridge module. Each guardian has their own switch. (halt|resume)`,
		Args:  cobra.ExactArgs(1),
		Example: fmt.Sprintf("%s tx %s set-bridge-status halt --from guardiankey\n"+
			"%s tx %s set-bridge-status resume --from guardiankey\n", version.AppName, types.ModuleName, version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			from := clientCtx.GetFromAddress().String()
			if _, err := sdk.AccAddressFromBech32(from); err != nil {
				return sdkerrors.ErrInvalidAddress.Wrapf("invalid address: %s", from)
			}

			conv := map[string]types.BridgeStatus{
				"halt":   types.StatusInactive,
				"resume": types.StatusActive,
			}

			msg := types.MsgSetBridgeStatus{
				Guardian: from,
				Status:   conv[args[0]],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
