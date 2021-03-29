package cli

import (
	"bufio"
	"errors"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/line/lbm-sdk/v2/x/token/internal/types"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

var (
	flagTotalSupply = "total-supply"
	flagDecimals    = "decimals"
	flagMintable    = "mintable"
	flagMeta        = "meta"
	flagImageURI    = "image-uri"
)

const (
	DefaultDecimals    = 8
	DefaultTotalSupply = 1
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Token transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		IssueTxCmd(cdc),
		TransferTxCmd(cdc),
		MintTxCmd(cdc),
		BurnTxCmd(cdc),
		GrantPermTxCmd(cdc),
		RevokePermTxCmd(cdc),
		ModifyTokenCmd(cdc),
		TransferFromTxCmd(cdc),
		ApproveTokenTxCmd(cdc),
		BurnFromTxCmd(cdc),
	)
	return txCmd
}

func IssueTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue [from_key_or_address] [to] [name] [symbol]",
		Short: "Create and sign an issue token tx",
		Long: `
[Issue a token command]
To query or send the token, you should remember the contract id


[Fungible Token]
linkcli tx token issue [from_key_or_address] [to] [name] [symbol]
--decimals=[decimals]
--mintable=[mintable]
--total-supply=[initial amount of the token]
--meta=[metadata for the token]
--image-uri=[image uri for the token]
`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			from := cliCtx.FromAddress
			to, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			name := args[2]
			symbol := args[3]
			supply := viper.GetInt64(flagTotalSupply)
			decimals := viper.GetInt64(flagDecimals)
			mintable := viper.GetBool(flagMintable)
			meta := viper.GetString(flagMeta)
			imageURI := viper.GetString(flagImageURI)

			if decimals < 0 || decimals > 18 {
				return errors.New("invalid decimals. 0 <= decimals <= 18")
			}

			msg := types.NewMsgIssue(from, to, name, symbol, meta, imageURI, sdk.NewInt(supply), sdk.NewInt(decimals), mintable)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().Int64(flagTotalSupply, DefaultTotalSupply, "total supply")
	cmd.Flags().Int64(flagDecimals, DefaultDecimals, "set decimals")
	cmd.Flags().Bool(flagMintable, false, "set mintable")
	cmd.Flags().String(flagMeta, "", "set meta")
	cmd.Flags().String(flagImageURI, "", "set img-uri")

	return flags.PostCommands(cmd)[0]
}

func TransferTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer [from_key_or_address] [to_address] [contract_id] [amount]",
		Short: "Create and sign a tx transferring non-reserved fungible tokens",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			to, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return sdkerrors.Wrap(types.ErrInvalidAmount, args[3])
			}

			msg := types.NewMsgTransfer(cliCtx.GetFromAddress(), to, args[2], amount)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func MintTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [from_key_or_address] [contract_id] [to] [amount]",
		Short: "Create and sign a mint token tx",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			contractID := args[1]

			to, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			amount, err := strconv.Atoi(args[3])
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgMint(cliCtx.GetFromAddress(), contractID, to, sdk.NewInt(int64(amount)))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return flags.PostCommands(cmd)[0]
}

func BurnTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [from_key_or_address] [contract_id] [amount]",
		Short: "Create and sign a burn token tx",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			contractID := args[1]
			amount, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return sdkerrors.Wrap(types.ErrInvalidAmount, args[4])
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgBurn(cliCtx.GetFromAddress(), contractID, amount)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return flags.PostCommands(cmd)[0]
}

func GrantPermTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grant [from_key_or_address] [contract_id] [to] [action]",
		Short: "Create and sign a grant permission for token tx",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			contractID := args[1]

			to, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}
			perm := types.Permission(args[3])

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgGrantPermission(cliCtx.GetFromAddress(), contractID, to, perm)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return flags.PostCommands(cmd)[0]
}

func RevokePermTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke [from_key_or_address] [contract_id] [action]",
		Short: "Create and sign a revoke permission for token tx",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			contractID := args[1]
			perm := types.Permission(args[2])

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgRevokePermission(cliCtx.GetFromAddress(), contractID, perm)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return flags.PostCommands(cmd)[0]
}

func ModifyTokenCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify [owner_address] [contract_id] [field] [new_value]",
		Short: "Create and sign a modify token tx",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)
			contractID := args[1]
			field := args[2]
			newValue := args[3]

			msg := types.NewMsgModify(
				cliCtx.FromAddress,
				contractID,
				types.NewChanges(types.NewChange(field, newValue)),
			)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return flags.PostCommands(cmd)[0]
}

func TransferFromTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-from [proxy_key_or_address] [contract_id] [from_address] [to_address] [amount]",
		Short: "Create and sign a tx transferring tokens by approved proxy",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			contractID := args[1]

			from, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			to, err := sdk.AccAddressFromBech32(args[3])
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[4])
			if !ok {
				return sdkerrors.Wrap(types.ErrInvalidAmount, args[4])
			}

			msg := types.NewMsgTransferFrom(cliCtx.GetFromAddress(), contractID, from, to, amount)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func ApproveTokenTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve [approver_key_or_address] [contract_id] [proxy_address]",
		Short: "Create and sign a tx approve all token operations of a token to a proxy",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			contractID := args[1]

			proxy, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgApprove(cliCtx.GetFromAddress(), contractID, proxy)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return flags.PostCommands(cmd)[0]
}
func BurnFromTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-from [proxy_key_or_address] [contract_id] [from_address] [amount]",
		Short: "Create and sign a burn token tx by approved proxy",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)
			contractID := args[1]
			from, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}
			amount, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return errors.New("invalid amount")
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgBurnFrom(cliCtx.GetFromAddress(), contractID, from, amount)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return flags.PostCommands(cmd)[0]
}
