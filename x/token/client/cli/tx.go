package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/flags"
	"github.com/line/lbm-sdk/client/tx"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/version"
	"github.com/line/lbm-sdk/x/token"
)

const (
	flagSupply = "supply"
	flagDecimals = "decimals"
	flagMintable = "mintable"
	flagMeta = "meta"
	flagImageUri = "image-uri"

	DefaultDecimals = "8"
	DefaultSupply = "1"
)

// NewTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
	nftTxCmd := &cobra.Command{
		Use:                        token.ModuleName,
		Short:                      "token transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	nftTxCmd.AddCommand(
		NewTxCmdTransfer(),
		NewTxCmdTransferFrom(),
		NewTxCmdApprove(),
		NewTxCmdIssue(),
		NewTxCmdGrant(),
		NewTxCmdRevoke(),
		NewTxCmdMint(),
		NewTxCmdBurn(),
		NewTxCmdBurnFrom(),
		NewTxCmdModify(),
	)

	return nftTxCmd
}

func NewTxCmdTransfer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer [class-id] [from] [to] [amount]",
		Args:  cobra.ExactArgs(4),
		Short: "transfer tokens",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s transfer <class-id> <from> <to> <amount>`, version.AppName, token.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "failed to set amount: %s")
			}
			msg := token.MsgTransfer{
				ClassId:  args[0],
				From: args[1],
				To: args[2],
				Amount: amount,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewTxCmdTransferFrom() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-from [class-id] [proxy] [from] [to] [amount]",
		Args:  cobra.ExactArgs(5),
		Short: "transfer tokens by proxy",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s transfer-from <class-id> <proxy> <from> <to> <amount>`, version.AppName, token.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[4])
			if !ok {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "failed to set amount: %s")
			}
			msg := token.MsgTransferFrom{
				ClassId:  args[0],
				Proxy: args[1],
				From: args[2],
				To: args[3],
				Amount: amount,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewTxCmdApprove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve [class-id] [approver] [proxy]",
		Args:  cobra.ExactArgs(3),
		Short: "approve transfer of tokens to a given proxy",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s approve <class-id> <approver> <proxy>`, version.AppName, token.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := token.MsgApprove{
				ClassId:  args[0],
				Approver: args[1],
				Proxy: args[2],
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewTxCmdIssue() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue [owner] [to] [name] [symbol]",
		Args:  cobra.ExactArgs(4),
		Short: "issue token",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s issue <owner> <to> <name> <symbol>`, version.AppName, token.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			imageUri, err := cmd.Flags().GetString(flagImageUri)
			if err != nil {
				return err
			}

			meta, err := cmd.Flags().GetString(flagMeta)
			if err != nil {
				return err
			}

			supplyStr, err := cmd.Flags().GetString(flagSupply)
			if err != nil {
				return err
			}
			supply, ok := sdk.NewIntFromString(supplyStr)
			if !ok {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "failed to set supply: %s")
			}

			mintable, err := cmd.Flags().GetBool(flagMintable)
			if err != nil {
				return err
			}

			decimalsStr, err := cmd.Flags().GetString(flagDecimals)
			if err != nil {
				return err
			}
			decimals, ok := sdk.NewIntFromString(decimalsStr)
			if !ok {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "failed to set decimals: %s")
			}

			msg := token.MsgIssue{
				Owner: args[0],
				To: args[1],
				Name: args[2],
				Symbol: args[3],
				ImageUri: imageUri,
				Meta: meta,
				Amount: supply,
				Mintable: mintable,
				Decimals: decimals,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(flagImageUri, "", "set image-uri")
	cmd.Flags().String(flagMeta, "", "set meta")
	cmd.Flags().String(flagSupply, DefaultSupply, "initial supply")
	cmd.Flags().Bool(flagMintable, false, "set mintable")
	cmd.Flags().String(flagDecimals, DefaultDecimals, "set decimals")

	return cmd
}

func NewTxCmdGrant() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grant [class-id] [granter] [grantee] [action]",
		Args:  cobra.ExactArgs(4),
		Short: "grant an action for mint, burn and modify",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s grant <class-id> <granter> <grantee> <action>`, version.AppName, token.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := token.MsgGrant{
				ClassId:  args[0],
				Granter: args[1],
				Grantee: args[2],
				Action: args[3],
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewTxCmdRevoke() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke [class-id] [grantee] [action]",
		Args:  cobra.ExactArgs(3),
		Short: "revoke an action by a given grantee",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s revoke <class-id> <grantee> <action>`, version.AppName, token.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := token.MsgRevoke{
				ClassId:  args[0],
				Grantee: args[1],
				Action: args[2],
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewTxCmdMint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [class-id] [grantee] [to] [amount]",
		Args:  cobra.ExactArgs(4),
		Short: "mint tokens",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s mint <class-id> <grantee> <to> <amount>`, version.AppName, token.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "failed to set amount: %s")
			}

			msg := token.MsgMint{
				ClassId:  args[0],
				Grantee: args[1],
				To: args[2],
				Amount: amount,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewTxCmdBurn() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [class-id] [from] [amount]",
		Args:  cobra.ExactArgs(3),
		Short: "burn tokens",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s burn <class-id> <from> <amount>`, version.AppName, token.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "failed to set amount: %s")
			}

			msg := token.MsgBurn{
				ClassId:  args[0],
				From: args[1],
				Amount: amount,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewTxCmdBurnFrom() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-from [class-id] [grantee] [from] [amount]",
		Args:  cobra.ExactArgs(4),
		Short: "burn tokens by a given grantee",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s burn-from <class-id> <grantee> <from> <amount>`, version.AppName, token.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "failed to set amount: %s")
			}

			msg := token.MsgBurnFrom{
				ClassId:  args[0],
				Grantee: args[1],
				From: args[2],
				Amount: amount,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewTxCmdModify() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify [class-id] [grantee] [key] [value]",
		Args:  cobra.ExactArgs(4),
		Short: "modify token metadata",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s modify <class-id> <grantee> <key> <value>`, version.AppName, token.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			change := token.Pair{Key: args[2], Value: args[3]}
			msg := token.MsgModify{
				ClassId:  args[0],
				Grantee: args[1],
				Changes: []token.Pair{change},
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
