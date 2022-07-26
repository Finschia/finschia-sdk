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
		NewTxCmdSend(),
		NewTxCmdOperatorSend(),
		NewTxCmdCreateFTClass(),
		NewTxCmdCreateNFTClass(),
		NewTxCmdMintFT(),
		NewTxCmdMintNFT(),
		NewTxCmdBurn(),
		NewTxCmdOperatorBurn(),
		NewTxCmdModifyContract(),
		NewTxCmdModifyTokenClass(),
		NewTxCmdModifyNFT(),
		NewTxCmdAttach(),
		NewTxCmdDetach(),
		NewTxCmdOperatorAttach(),
		NewTxCmdOperatorDetach(),
		NewTxCmdGrant(),
		NewTxCmdAbandon(),
		NewTxCmdAuthorizeOperator(),
		NewTxCmdRevokeOperator(),
	)

	return txCmd
}

func NewTxCmdSend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [contract-id] [from] [to] [amount]",
		Args:  cobra.ExactArgs(4),
		Short: "send tokens",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s send [contract-id] [from] [to] [amount]`, version.AppName, collection.ModuleName),
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

			msg := &collection.MsgSend{
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

func NewTxCmdOperatorSend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operator-send [contract-id] [operator] [from] [to] [amount]",
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

			msg := collection.MsgOperatorSend{
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
				Owner:      creator,
				Name:       name,
				BaseImgUri: baseImgURI,
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
	cmd.Flags().String(FlagBaseImgURI, "", "set base-img-uri")
	cmd.Flags().String(FlagMeta, "", "set meta")

	return cmd
}

func NewTxCmdCreateFTClass() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-ft-class [contract-id] [operator]",
		Args:  cobra.ExactArgs(2),
		Short: "create a fungible token class",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s create-ft-class [contract-id] [operator]`, version.AppName, collection.ModuleName),
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

			msg := collection.MsgCreateFTClass{
				ContractId: args[0],
				Operator:   operator,
				Name:       name,
				Meta:       meta,
				Decimals:   decimals,
				To:         to,
				Supply:     supply,
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
	cmd.Flags().String(FlagTo, "", "address to send the initial supply")
	cmd.Flags().String(FlagSupply, DefaultSupply, "initial supply")
	cmd.Flags().Int32(FlagDecimals, DefaultDecimals, "set decimals")

	return cmd
}

func NewTxCmdCreateNFTClass() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-nft-class [contract-id] [operator]",
		Args:  cobra.ExactArgs(2),
		Short: "create a non-fungible token class",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s create-nft-class [contract-id] [operator]`, version.AppName, collection.ModuleName),
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

			msg := collection.MsgCreateNFTClass{
				ContractId: args[0],
				Operator:   operator,
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

	return cmd
}

func NewTxCmdBurn() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [contract-id] [from] [amount]",
		Args:  cobra.ExactArgs(3),
		Short: "burn tokens",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s burn [contract-id] [from] [amount]`, version.AppName, collection.ModuleName),
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

			msg := collection.MsgBurn{
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

func NewTxCmdOperatorBurn() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operator-burn [contract-id] [operator] [from] [amount]",
		Args:  cobra.ExactArgs(4),
		Short: "burn tokens by a given operator",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s operator-burn [contract-id] [operator] [from] [amount]`, version.AppName, collection.ModuleName),
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

			msg := collection.MsgOperatorBurn{
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

func NewTxCmdModifyContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify-contract [contract-id] [operator] [key] [value]",
		Args:  cobra.ExactArgs(4),
		Short: "modify contract metadata",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s modify-contract [contract-id] [operator] [key] [value]`, version.AppName, collection.ModuleName),
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
				Key:   args[2],
				Value: args[3],
			}}
			msg := collection.MsgModifyContract{
				ContractId: args[0],
				Operator:   args[1],
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

func NewTxCmdModifyTokenClass() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify-token-class [contract-id] [operator] [class-id] [key] [value]",
		Args:  cobra.ExactArgs(5),
		Short: "modify token class metadata",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s modify-token-class [contract-id] [operator] [class-id] [key] [value]`, version.AppName, collection.ModuleName),
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
				Key:   args[3],
				Value: args[4],
			}}
			msg := collection.MsgModifyTokenClass{
				ContractId: args[0],
				Operator:   args[1],
				ClassId:    args[2],
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

func NewTxCmdModifyNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify-nft [contract-id] [operator] [token-id] [key] [value]",
		Args:  cobra.ExactArgs(5),
		Short: "modify nft metadata",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s modify-nft [contract-id] [operator] [token-id] [key] [value]`, version.AppName, collection.ModuleName),
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
				Key:   args[3],
				Value: args[4],
			}}
			msg := collection.MsgModifyNFT{
				ContractId: args[0],
				Operator:   args[1],
				TokenId:    args[2],
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
		Use:   "operator-attach [contract-id] [operator] [holder] [subject] [target]",
		Args:  cobra.ExactArgs(5),
		Short: "attach a token to another by the operator",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s operator-attach [contract-id] [operator] [holder] [subject] [target]`, version.AppName, collection.ModuleName),
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
				Owner:      args[2],
				Subject:    args[3],
				Target:     args[4],
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
		Use:   "operator-detach [contract-id] [operator] [holder] [subject]",
		Args:  cobra.ExactArgs(4),
		Short: "detach a token from its parent by the operator",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s operator-detach [contract-id] [operator] [holder] [subject]`, version.AppName, collection.ModuleName),
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
				Owner:      args[2],
				Subject:    args[3],
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

func NewTxCmdGrant() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grant [contract-id] [granter] [grantee] [permission]",
		Args:  cobra.ExactArgs(4),
		Short: "grant a permission for mint, burn, modify and issue",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s grant [contract-id] [granter] [grantee] [permission]`, version.AppName, collection.ModuleName),
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

			permission := collection.Permission(collection.Permission_value[args[3]])

			msg := collection.MsgGrant{
				ContractId: args[0],
				Granter:    granter,
				Grantee:    args[2],
				Permission: permission,
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
		Use:   "abandon [contract-id] [grantee] [permission]",
		Args:  cobra.ExactArgs(3),
		Short: "abandon a permission by a given grantee",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s abandon [contract-id] [grantee] [permission]`, version.AppName, collection.ModuleName),
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

			permission := collection.Permission(collection.Permission_value[args[2]])

			msg := collection.MsgAbandon{
				ContractId: args[0],
				Grantee:    grantee,
				Permission: permission,
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
		Use:   "authorize-operator [contract-id] [holder] [operator]",
		Args:  cobra.ExactArgs(3),
		Short: "authorize operator to manipulate tokens of holder",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s authorize-operator [contract-id] [holder] [operator]`, version.AppName, collection.ModuleName),
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
		Use:   "revoke-operator [contract-id] [holder] [operator]",
		Args:  cobra.ExactArgs(3),
		Short: "revoke operator",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s revoke-operator [contract-id] [holder] [operator]`, version.AppName, collection.ModuleName),
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
