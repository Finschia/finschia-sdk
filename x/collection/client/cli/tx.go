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
	"github.com/line/lbm-sdk/x/collection"
)

const (
	// common flags for the entities
	FlagName = "name"
	FlagMeta = "meta"

	// flag for contracts
	FlagBaseImgURI = "base-img-uri"

	// flag for fungible token classes
	FlagDecimals = "decimals"
	FlagMintable = "mintable"
	FlagTo       = "to"
	FlagSupply   = "supply"

	DefaultDecimals = 8
	DefaultSupply   = "0"
)

// NewTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        collection.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", collection.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewTxCmdSendFT(),
		NewTxCmdOperatorSendFT(),
		NewTxCmdSendNFT(),
		NewTxCmdOperatorSendNFT(),
		NewTxCmdCreateContract(),
		NewTxCmdIssueFT(),
		NewTxCmdIssueNFT(),
		NewTxCmdMintFT(),
		NewTxCmdMintNFT(),
		NewTxCmdAttach(),
		NewTxCmdDetach(),
		NewTxCmdOperatorAttach(),
		NewTxCmdOperatorDetach(),
		NewTxCmdGrantPermission(),
		NewTxCmdRevokePermission(),
		NewTxCmdAuthorizeOperator(),
		NewTxCmdRevokeOperator(),
	)

	return txCmd
}

func NewTxCmdSendFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-ft [contract-id] [from] [to] [amount]",
		Args:  cobra.ExactArgs(4),
		Short: "send fungible tokens",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s transfer-ft [contract-id] [from] [to] [amount]`, version.AppName, collection.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			from := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, from); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amountStr := args[3]
			amount, err := collection.ParseCoins(amountStr)
			if err != nil {
				return err
			}

			msg := &collection.MsgSendFT{
				ContractId: args[0],
				From:       from,
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

func NewTxCmdOperatorSendFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-ft-from [contract-id] [operator] [from] [to] [amount]",
		Args:  cobra.ExactArgs(5),
		Short: "send tokens by operator",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s operator-send [contract-id] [operator] [from] [to] [amount]`, version.AppName, collection.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			operator := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, operator); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amountStr := args[4]
			amount, err := collection.ParseCoins(amountStr)
			if err != nil {
				return err
			}

			msg := collection.MsgOperatorSendFT{
				ContractId: args[0],
				Operator:   operator,
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

func NewTxCmdSendNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-nft [contract-id] [from] [to] [token-id]",
		Args:  cobra.ExactArgs(4),
		Short: "send non-fungible tokens",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s transfer-nft [contract-id] [from] [to] [token-id]`, version.AppName, collection.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			from := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, from); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &collection.MsgSendNFT{
				ContractId: args[0],
				From:       from,
				To:         args[2],
				TokenIds:   []string{args[3]},
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

func NewTxCmdOperatorSendNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-nft-from [contract-id] [operator] [from] [to] [amount]",
		Args:  cobra.ExactArgs(5),
		Short: "send tokens by operator",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s transfer-nft-from [contract-id] [operator] [from] [to] [amount]`, version.AppName, collection.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			operator := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, operator); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := collection.MsgOperatorSendNFT{
				ContractId: args[0],
				Operator:   operator,
				From:       args[2],
				To:         args[3],
				TokenIds:   []string{args[4]},
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

func NewTxCmdCreateContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-contract [creator]",
		Args:  cobra.ExactArgs(1),
		Short: "create a contract",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s create-contract [creator]`, version.AppName, collection.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			creator := args[0]
			if err := cmd.Flags().Set(flags.FlagFrom, creator); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return err
			}

			baseImgURI, err := cmd.Flags().GetString(FlagBaseImgURI)
			if err != nil {
				return err
			}

			meta, err := cmd.Flags().GetString(FlagMeta)
			if err != nil {
				return err
			}

			msg := collection.MsgCreateContract{
				Owner: creator,
				Name:  name,
				Uri:   baseImgURI,
				Meta:  meta,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(FlagName, "", "set name")
	cmd.Flags().String(FlagBaseImgURI, "", "set base-img-uri")
	cmd.Flags().String(FlagMeta, "", "set meta")

	return cmd
}

func NewTxCmdIssueFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-ft [contract-id] [operator]",
		Args:  cobra.ExactArgs(2),
		Short: "create a fungible token class",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s issue-ft [contract-id] [operator]`, version.AppName, collection.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			operator := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, operator); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return err
			}

			meta, err := cmd.Flags().GetString(FlagMeta)
			if err != nil {
				return err
			}

			decimals, err := cmd.Flags().GetInt32(FlagDecimals)
			if err != nil {
				return err
			}

			mintable, err := cmd.Flags().GetBool(FlagMintable)
			if err != nil {
				return err
			}

			supplyStr, err := cmd.Flags().GetString(FlagSupply)
			if err != nil {
				return err
			}
			supply, ok := sdk.NewIntFromString(supplyStr)
			if !ok {
				return sdkerrors.ErrInvalidType.Wrapf("failed to set supply: %s", supplyStr)
			}

			to, err := cmd.Flags().GetString(FlagTo)
			if err != nil {
				return err
			}

			msg := collection.MsgIssueFT{
				ContractId: args[0],
				Owner:      operator,
				Name:       name,
				Meta:       meta,
				Decimals:   decimals,
				Mintable:   mintable,
				To:         to,
				Amount:     supply,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(FlagName, "", "set name")
	cmd.MarkFlagRequired(FlagName)
	cmd.Flags().String(FlagMeta, "", "set meta")
	cmd.Flags().String(FlagTo, "", "address to send the initial supply")
	cmd.MarkFlagRequired(FlagTo)
	cmd.Flags().Bool(FlagMintable, false, "set mintable")
	cmd.Flags().String(FlagSupply, DefaultSupply, "initial supply")
	cmd.Flags().Int32(FlagDecimals, DefaultDecimals, "set decimals")

	return cmd
}

func NewTxCmdIssueNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-nft [contract-id] [operator]",
		Args:  cobra.ExactArgs(2),
		Short: "create a non-fungible token class",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s issue-nft [contract-id] [operator]`, version.AppName, collection.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			operator := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, operator); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return err
			}

			meta, err := cmd.Flags().GetString(FlagMeta)
			if err != nil {
				return err
			}

			msg := collection.MsgIssueNFT{
				ContractId: args[0],
				Owner:      operator,
				Name:       name,
				Meta:       meta,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(FlagName, "", "set name")
	cmd.Flags().String(FlagMeta, "", "set meta")

	return cmd
}

func NewTxCmdMintFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-ft [contract-id] [operator] [to] [class-id] [amount]",
		Args:  cobra.ExactArgs(5),
		Short: "mint fungible tokens",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s mint-ft [contract-id] [operator] [to] [class-id] [amount]`, version.AppName, collection.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			operator := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, operator); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amountStr := args[4]
			amount, ok := sdk.NewIntFromString(amountStr)
			if !ok {
				return sdkerrors.ErrInvalidType.Wrapf("failed to set amount: %s", amountStr)
			}

			coins := collection.NewCoins(collection.NewFTCoin(args[3], amount))
			msg := collection.MsgMintFT{
				ContractId: args[0],
				From:       args[1],
				To:         args[2],
				Amount:     coins,
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

func NewTxCmdMintNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-nft [contract-id] [operator] [to] [class-id]",
		Args:  cobra.ExactArgs(4),
		Short: "mint non-fungible tokens",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s mint-nft [contract-id] [operator] [to] [class-id]`, version.AppName, collection.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			operator := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, operator); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return err
			}

			meta, err := cmd.Flags().GetString(FlagMeta)
			if err != nil {
				return err
			}

			params := []collection.MintNFTParam{{
				TokenType: args[3],
				Name:      name,
				Meta:      meta,
			}}

			msg := collection.MsgMintNFT{
				ContractId: args[0],
				From:       args[1],
				To:         args[2],
				Params:     params,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(FlagName, "", "set name")
	cmd.Flags().String(FlagMeta, "", "set meta")
	cmd.MarkFlagRequired(FlagName)

	return cmd
}

func NewTxCmdBurnFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-ft [contract-id] [from] [amount]",
		Args:  cobra.ExactArgs(3),
		Short: "burn tokens",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s burn-ft [contract-id] [from] [amount]`, version.AppName, collection.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			from := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, from); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amountStr := args[2]
			amount, err := collection.ParseCoins(amountStr)
			if err != nil {
				return err
			}

			msg := collection.MsgBurnFT{
				ContractId: args[0],
				From:       from,
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

func NewTxCmdOperatorBurnFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-ft-from [contract-id] [operator] [from] [amount]",
		Args:  cobra.ExactArgs(4),
		Short: "burn tokens by a given operator",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s burn-ft-from [contract-id] [operator] [from] [amount]`, version.AppName, collection.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			operator := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, operator); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amountStr := args[3]
			amount, err := collection.ParseCoins(amountStr)
			if err != nil {
				return err
			}

			msg := collection.MsgOperatorBurnFT{
				ContractId: args[0],
				Operator:   operator,
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

func NewTxCmdBurnNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-nft [contract-id] [from] [token-id]",
		Args:  cobra.ExactArgs(3),
		Short: "burn tokens",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s burn-nft [contract-id] [from] [token-id]`, version.AppName, collection.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			from := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, from); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := collection.MsgBurnNFT{
				ContractId: args[0],
				From:       from,
				TokenIds:   []string{args[2]},
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

func NewTxCmdOperatorBurnNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-nft-from [contract-id] [operator] [from] [token-id]",
		Args:  cobra.ExactArgs(4),
		Short: "burn tokens by a given operator",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s burn-nft-from [contract-id] [operator] [from] [token-id]`, version.AppName, collection.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			operator := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, operator); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := collection.MsgOperatorBurnNFT{
				ContractId: args[0],
				Operator:   operator,
				From:       args[2],
				TokenIds:   []string{args[3]},
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
		Use:   "modify [contract-id] [operator] [token-type] [token-index] [key] [value]",
		Args:  cobra.ExactArgs(6),
		Short: "modify",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s modify [contract-id] [operator] [token-type] [token-index] [key] [value]`, version.AppName, collection.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			operator := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, operator); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			changes := []collection.Attribute{{
				Key:   args[4],
				Value: args[5],
			}}
			msg := collection.MsgModify{
				ContractId: args[0],
				Owner:      args[1],
				TokenType:  args[2],
				TokenIndex: args[3],
				Changes:    changes,
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

func NewTxCmdAttach() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attach [contract-id] [holder] [subject] [target]",
		Args:  cobra.ExactArgs(4),
		Short: "attach a token to another",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s attach [contract-id] [holder] [subject] [target]`, version.AppName, collection.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			holder := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, holder); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := collection.MsgAttach{
				ContractId: args[0],
				From:       holder,
				TokenId:    args[2],
				ToTokenId:  args[3],
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

func NewTxCmdDetach() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach [contract-id] [holder] [subject]",
		Args:  cobra.ExactArgs(3),
		Short: "detach a token from its parent",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s detach [contract-id] [holder] [subject]`, version.AppName, collection.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			holder := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, holder); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := collection.MsgDetach{
				ContractId: args[0],
				From:       holder,
				TokenId:    args[2],
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

func NewTxCmdOperatorAttach() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attach-from [contract-id] [operator] [holder] [subject] [target]",
		Args:  cobra.ExactArgs(5),
		Short: "attach a token to another by the operator",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s attach-from [contract-id] [operator] [holder] [subject] [target]`, version.AppName, collection.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			operator := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, operator); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := collection.MsgOperatorAttach{
				ContractId: args[0],
				Operator:   operator,
				From:       args[2],
				TokenId:    args[3],
				ToTokenId:  args[4],
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

func NewTxCmdOperatorDetach() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach-from [contract-id] [operator] [holder] [subject]",
		Args:  cobra.ExactArgs(4),
		Short: "detach a token from its parent by the operator",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s detach-from [contract-id] [operator] [holder] [subject]`, version.AppName, collection.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			operator := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, operator); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := collection.MsgOperatorDetach{
				ContractId: args[0],
				Operator:   operator,
				From:       args[2],
				TokenId:    args[3],
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

func NewTxCmdGrantPermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grant-permission [contract-id] [granter] [grantee] [permission]",
		Args:  cobra.ExactArgs(4),
		Short: "grant a permission for mint, burn, modify and issue",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s grant-permission [contract-id] [granter] [grantee] [permission]`, version.AppName, collection.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			granter := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, granter); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := collection.MsgGrantPermission{
				ContractId: args[0],
				From:       granter,
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

func NewTxCmdRevokePermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke-permission [contract-id] [grantee] [permission]",
		Args:  cobra.ExactArgs(3),
		Short: "revoke a permission by a given grantee",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s revoke-permission [contract-id] [grantee] [permission]`, version.AppName, collection.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			grantee := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, grantee); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := collection.MsgRevokePermission{
				ContractId: args[0],
				From:       grantee,
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

func NewTxCmdAuthorizeOperator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve [contract-id] [holder] [operator]",
		Args:  cobra.ExactArgs(3),
		Short: "authorize operator to manipulate tokens of holder",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s approve [contract-id] [holder] [operator]`, version.AppName, collection.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			holder := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, holder); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := collection.MsgAuthorizeOperator{
				ContractId: args[0],
				Holder:     holder,
				Operator:   args[2],
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
		Use:   "disapprove [contract-id] [holder] [operator]",
		Args:  cobra.ExactArgs(3),
		Short: "revoke operator",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s disapprove [contract-id] [holder] [operator]`, version.AppName, collection.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			holder := args[1]
			if err := cmd.Flags().Set(flags.FlagFrom, holder); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := collection.MsgRevokeOperator{
				ContractId: args[0],
				Holder:     holder,
				Operator:   args[2],
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
