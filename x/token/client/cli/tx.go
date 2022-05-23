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
	FlagSupply   = "supply"
	FlagDecimals = "decimals"
	FlagMintable = "mintable"
	FlagMeta     = "meta"
	FlagImageURI = "image-uri"

	DefaultDecimals = 8
	DefaultSupply   = "1"
)

// NewTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        token.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", token.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewTxCmdSend(),
		NewTxCmdOperatorSend(),
		NewTxCmdAuthorizeOperator(),
		NewTxCmdRevokeOperator(),
		NewTxCmdIssue(),
		NewTxCmdGrant(),
		NewTxCmdAbandon(),
		NewTxCmdMint(),
		NewTxCmdBurn(),
		NewTxCmdOperatorBurn(),
		NewTxCmdModify(),
	)

	return txCmd
}

func NewTxCmdSend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [class-id] [from] [to] [amount]",
		Args:  cobra.ExactArgs(4),
		Short: "send tokens",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s send <class-id> <from> <to> <amount>`, version.AppName, token.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amountStr := args[3]
			amount, ok := sdk.NewIntFromString(amountStr)
			if !ok {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "failed to set amount: %s", amountStr)
			}
			msg := &token.MsgSend{
				ContractId: args[0],
				From:       args[1],
				To:         args[2],
				Amount:     amount,
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

func NewTxCmdOperatorSend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operator-send [class-id] [proxy] [from] [to] [amount]",
		Args:  cobra.ExactArgs(5),
		Short: "send tokens by proxy",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s operator-send <class-id> <proxy> <from> <to> <amount>`, version.AppName, token.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amountStr := args[4]
			amount, ok := sdk.NewIntFromString(amountStr)
			if !ok {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "failed to set amount: %s", amountStr)
			}
			msg := token.MsgOperatorSend{
				ContractId: args[0],
				Proxy:      args[1],
				From:       args[2],
				To:         args[3],
				Amount:     amount,
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

func NewTxCmdAuthorizeOperator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "authorize-operator [class-id] [approver] [proxy]",
		Args:  cobra.ExactArgs(3),
		Short: "authorize operator to send tokens to a given proxy",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s authorize-operator <class-id> <approver> <proxy>`, version.AppName, token.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := token.MsgAuthorizeOperator{
				ContractId: args[0],
				Approver:   args[1],
				Proxy:      args[2],
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

func NewTxCmdRevokeOperator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke-operator [class-id] [approver] [proxy]",
		Args:  cobra.ExactArgs(3),
		Short: "revoke operator to send tokens to a given proxy",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s revoke-operator <class-id> <approver> <proxy>`, version.AppName, token.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := token.MsgRevokeOperator{
				ContractId: args[0],
				Approver:   args[1],
				Proxy:      args[2],
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

			imageURI, err := cmd.Flags().GetString(FlagImageURI)
			if err != nil {
				return err
			}

			meta, err := cmd.Flags().GetString(FlagMeta)
			if err != nil {
				return err
			}

			supplyStr, err := cmd.Flags().GetString(FlagSupply)
			if err != nil {
				return err
			}
			supply, ok := sdk.NewIntFromString(supplyStr)
			if !ok {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "failed to set supply: %s", supplyStr)
			}

			mintable, err := cmd.Flags().GetBool(FlagMintable)
			if err != nil {
				return err
			}

			decimals, err := cmd.Flags().GetInt32(FlagDecimals)
			if err != nil {
				return err
			}

			msg := token.MsgIssue{
				Owner:    args[0],
				To:       args[1],
				Name:     args[2],
				Symbol:   args[3],
				ImageUri: imageURI,
				Meta:     meta,
				Amount:   supply,
				Mintable: mintable,
				Decimals: decimals,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(FlagImageURI, "", "set image-uri")
	cmd.Flags().String(FlagMeta, "", "set meta")
	cmd.Flags().String(FlagSupply, DefaultSupply, "initial supply")
	cmd.Flags().Bool(FlagMintable, false, "set mintable")
	cmd.Flags().Int32(FlagDecimals, DefaultDecimals, "set decimals")

	return cmd
}

func NewTxCmdGrant() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grant [class-id] [granter] [grantee] [permission]",
		Args:  cobra.ExactArgs(4),
		Short: "grant a permission for mint, burn and modify",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s grant <class-id> <granter> <grantee> <permission>`, version.AppName, token.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := token.MsgGrant{
				ContractId: args[0],
				From:       args[1],
				To:         args[2],
				Permission: args[3],
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

func NewTxCmdAbandon() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "abandon [class-id] [grantee] [permission]",
		Args:  cobra.ExactArgs(3),
		Short: "abandon a permission by a given grantee",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s revoke <class-id> <grantee> <permission>`, version.AppName, token.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := token.MsgAbandon{
				ContractId: args[0],
				Grantee:    args[1],
				Permission: args[2],
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

			amountStr := args[3]
			amount, ok := sdk.NewIntFromString(amountStr)
			if !ok {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "failed to set amount: %s", amountStr)
			}

			msg := token.MsgMint{
				ContractId: args[0],
				From:       args[1],
				To:         args[2],
				Amount:     amount,
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

			amountStr := args[2]
			amount, ok := sdk.NewIntFromString(amountStr)
			if !ok {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "failed to set amount: %s", amountStr)
			}

			msg := token.MsgBurn{
				ContractId: args[0],
				From:       args[1],
				Amount:     amount,
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

func NewTxCmdOperatorBurn() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operator-burn [class-id] [grantee] [from] [amount]",
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

			amountStr := args[3]
			amount, ok := sdk.NewIntFromString(amountStr)
			if !ok {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "failed to set amount: %s", amountStr)
			}

			msg := token.MsgOperatorBurn{
				ContractId: args[0],
				Proxy:      args[1],
				From:       args[2],
				Amount:     amount,
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

			change := token.Pair{Field: args[2], Value: args[3]}
			msg := token.MsgModify{
				ContractId: args[0],
				Owner:      args[1],
				Changes:    []token.Pair{change},
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
