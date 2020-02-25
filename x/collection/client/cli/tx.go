package cli

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/line/link/x/collection/internal/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/line/link/client"

	linktype "github.com/line/link/types"
)

var (
	flagTotalSupply = "total-supply"
	flagDecimals    = "decimals"
	flagMintable    = "mintable"
	flagTokenURI    = "token-uri"
	flagTokenType   = "token-type"
	flagAAS         = "address-suffix"
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
		CreateCollectionTxCmd(cdc),
		IssueFTTxCmd(cdc),
		IssueNFTTxCmd(cdc),
		MintFTTxCmd(cdc),
		MintNFTTxCmd(cdc),
		BurnFTTxCmd(cdc),
		BurnFTFromTxCmd(cdc),
		BurnNFTTxCmd(cdc),
		BurnNFTFromTxCmd(cdc),
		TransferFTTxCmd(cdc),
		TransferFTFromTxCmd(cdc),
		TransferNFTTxCmd(cdc),
		TransferNFTFromTxCmd(cdc),
		AttachTxCmd(cdc),
		AttachFromTxCmd(cdc),
		DetachTxCmd(cdc),
		DetachFromTxCmd(cdc),
		ApproveCollectionTxCmd(cdc),
		DisapproveCollectionTxCmd(cdc),
		GrantPermTxCmd(cdc),
		RevokePermTxCmd(cdc),
		ModifyTokenURICmd(cdc),
	)
	return txCmd
}

func ModifyTokenURICmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify-token-uri [owner_address] [symbol] [token_id] [token_uri]",
		Short: "Create and sign a modify token_uri of token tx",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)
			symbol := args[1]
			tokenID := args[2]
			tokenURI := args[3]

			msg := types.NewMsgModifyTokenURI(cliCtx.FromAddress, symbol, tokenURI, tokenID)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func CreateCollectionTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [from_key_or_address] [symbol] [name]",
		Short: "Create and sign an create collection tx",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			owner := cliCtx.FromAddress
			symbol := args[1]
			name := args[2]

			aas := viper.GetBool(flagAAS)

			if aas {
				symbol += owner.String()[len(owner.String())-linktype.AccAddrSuffixLen:]
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgCreateCollection(owner, name, symbol)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().Bool(flagAAS, true, "attach address suffix to symbol")

	return client.PostCommands(cmd)[0]
}

func IssueNFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-nft [from_key_or_address] [symbol]",
		Short: "Create and sign an issue-nft tx",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			to := cliCtx.FromAddress
			symbol := args[1]

			msg := types.NewMsgIssueNFT(to, symbol)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return client.PostCommands(cmd)[0]
}

func IssueFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-ft [from_key_or_address] [symbol] [name]",
		Short: "Create and sign an issue-ft tx",
		Long: `
[Fungible Token]
linkcli tx token issue [from_key_or_address] [symbol] [name]
--decimals=[decimals]
--mintable=[mintable]
--total-supply=[initial amount of the token]
--token-uri=[metadata for the token]
`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			to := cliCtx.FromAddress
			symbol := args[1]
			name := args[2]
			supply := viper.GetInt64(flagTotalSupply)
			decimals := viper.GetInt64(flagDecimals)
			mintable := viper.GetBool(flagMintable)
			tokenURI := viper.GetString(flagTokenURI)

			if err := linktype.ValidateSymbolUserDefined(symbol); err != nil {
				return err
			}

			if decimals < 0 || decimals > 18 {
				return errors.New("invalid decimals. 0 <= decimals <= 18")
			}

			msg := types.NewMsgIssueFT(to, name, symbol, tokenURI, sdk.NewInt(supply), sdk.NewInt(decimals), mintable)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().Int64(flagTotalSupply, DefaultTotalSupply, "total supply")
	cmd.Flags().Int64(flagDecimals, DefaultDecimals, "set decimals")
	cmd.Flags().Bool(flagMintable, false, "set mintable")
	cmd.Flags().String(flagTokenURI, "", "set token-uri")

	return client.PostCommands(cmd)[0]
}

func MintNFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-nft [from_key_or_address] [to] [symbol] [token_type] [name]",
		Short: "Create and sign an mint-nft tx",
		Long: `
[NonFungible Token]
linkcli tx token mint-nft [from_key_or_address] [symbol] [token_type] [name]
--token-uri=[metadata for the token]
`,
		Args: cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			from := cliCtx.FromAddress
			to, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			symbol := args[2]
			tokenType := args[3]
			name := args[4]
			tokenURI := viper.GetString(flagTokenURI)

			if err := linktype.ValidateSymbolUserDefined(symbol); err != nil {
				return err
			}

			msg := types.NewMsgMintNFT(from, to, name, symbol, tokenURI, tokenType)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(flagTokenURI, "", "set token-uri")
	cmd.Flags().String(flagTokenType, "", "token-type for the nft")

	return client.PostCommands(cmd)[0]
}

func BurnNFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-nft [from_key_or_address] [symbol] [token_id]",
		Short: "Create and sign an burn-nft tx",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			symbol := args[1]
			tokenID := args[2]

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgBurnNFT(cliCtx.GetFromAddress(), symbol, tokenID)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func BurnNFTFromTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-nft-from [proxy_key_or_address] [from_address] [symbol] [token_id]",
		Short: "Create and sign an burn-nft-from tx",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			from, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			symbol := args[2]
			tokenID := args[3]

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgBurnNFTFrom(cliCtx.GetFromAddress(), from, symbol, tokenID)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func TransferFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-ft [from_key_or_address] [to_address] [symbol] [amount]",
		Short: "Create and sign a tx transferring non-reserved collective fungible tokens",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			to, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			amount, err := types.ParseCoins(args[3])
			if err != nil {
				return types.ErrInvalidAmount(types.DefaultCodespace, args[3])
			}

			msg := types.NewMsgTransferFT(cliCtx.GetFromAddress(), to, args[2], amount...)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

func TransferNFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-nft [from_key_or_address] [to_address] [symbol] [token_id]",
		Short: "Create and sign a tx transferring a collective non-fungible token",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			to, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgTransferNFT(cliCtx.GetFromAddress(), to, args[2], args[3])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func TransferFTFromTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-ft-from [proxy_key_or_address] [from_address] [to_address] [symbol] [amount]",
		Short: "Create and sign a tx transferring non-reserved collective fungible tokens by approved proxy",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			from, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			to, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			amount, err := types.ParseCoins(args[4])
			if err != nil {
				return types.ErrInvalidAmount(types.DefaultCodespace, args[4])
			}

			msg := types.NewMsgTransferFTFrom(cliCtx.GetFromAddress(), from, to, args[3], amount...)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

//nolint:dupl
func TransferNFTFromTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-nft-from [proxy_key_or_address] [from_address] [to_address] [symbol] [token_id]",
		Short: "Create and sign a tx transferring a collective non-fungible token by approved proxy",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			from, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			to, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgTransferNFTFrom(cliCtx.GetFromAddress(), from, to, args[3], args[4])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func AttachTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attach [from_key_or_address] [symbol] [to_token_id] [token_id]",
		Short: "Create and sign a tx attaching a token to other",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			msg := types.NewMsgAttach(cliCtx.GetFromAddress(), args[1], args[2], args[3])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func DetachTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach [from_key_or_address] [symbol] [token_id]",
		Short: "Create and sign a tx detaching a token",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			msg := types.NewMsgDetach(cliCtx.GetFromAddress(), args[1], args[2])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

//nolint:dupl
func MintFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-ft [from_key_or_address] [to] [symbol] [token-id] [amount]",
		Short: "Create and sign a mint token tx",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)
			to, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			symbol := args[2]
			tokenID := args[3]
			amount, ok := sdk.NewIntFromString(args[4])
			if !ok {
				return errors.New("invalid amount")
			}

			msg := types.NewMsgMintFT(symbol, cliCtx.GetFromAddress(), to, types.NewCoin(tokenID, amount))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func BurnFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-ft [from_key_or_address] [symbol] [token-id] [amount]",
		Short: "Create and sign a mint token tx",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)
			symbol := args[1]
			tokenID := args[2]
			amount, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return errors.New("invalid amount")
			}

			msg := types.NewMsgBurnFT(symbol, cliCtx.GetFromAddress(), types.NewCoin(tokenID, amount))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

//nolint:dupl
func BurnFTFromTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-ft-from [proxy_key_or_address] [from_address] [symbol] [token-id] [amount]",
		Short: "Create and sign a mint token tx",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)
			from, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			symbol := args[2]
			tokenID := args[3]
			amount, ok := sdk.NewIntFromString(args[4])
			if !ok {
				return errors.New("invalid amount")
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgBurnFTFrom(symbol, cliCtx.GetFromAddress(), from, types.NewCoin(tokenID, amount))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func AttachFromTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attach-from [proxy_key_or_address] [from_address] [symbol] [to_token_id] [token_id]",
		Short: "Create and sign a tx attaching a token to other by approved proxy",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			from, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgAttachFrom(cliCtx.GetFromAddress(), from, args[2], args[3], args[4])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

//nolint:dupl
func DetachFromTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach-from [proxy_key_or_address] [from_address] [symbol] [token_id]",
		Short: "Create and sign a tx detaching a token by approved proxy",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			from, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDetachFrom(cliCtx.GetFromAddress(), from, args[2], args[3])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func ApproveCollectionTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve [approver_key_or_address] [proxy_address] [symbol]",
		Short: "Create and sign a tx approve all token operations of a collection to a proxy",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			proxy, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgApprove(cliCtx.GetFromAddress(), proxy, args[2])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func DisapproveCollectionTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disapprove [approver_key_or_address] [proxy_address] [symbol]",
		Short: "Create and sign a tx disapprove all token operations of a collection to a proxy",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			proxy, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDisapprove(cliCtx.GetFromAddress(), proxy, args[2])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func GrantPermTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grant [from_key_or_address] [to] [resource] [action]",
		Short: "Create and sign a grant permission for token tx",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			to, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			perm := types.Permission{Resource: args[2], Action: args[3]}
			if !perm.Validate() {
				return errors.New("permission invalid")
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgGrantPermission(cliCtx.GetFromAddress(), to, perm)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func RevokePermTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke [from_key_or_address] [resource] [action]",
		Short: "Create and sign a revoke permission for token tx",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			perm := types.Permission{Resource: args[1], Action: args[2]}
			if !perm.Validate() {
				return errors.New("permission invalid")
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgRevokePermission(cliCtx.GetFromAddress(), perm)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}
