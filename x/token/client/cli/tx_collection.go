package cli

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/client/context"
	linktype "github.com/line/link/types"
	"github.com/line/link/x/token/internal/types"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/line/link/client"
)

func CreateCollectionTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection-create [from_key_or_address] [symbol] [name]",
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
			msg := types.NewMsgCreateCollection(name, symbol, owner)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().Bool(flagAAS, true, "attach address suffix to symbol")

	return client.PostCommands(cmd)[0]
}

func IssueCollectionNFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection-issue-nft [from_key_or_address] [symbol]",
		Short: "Create and sign an collection-issue-nft tx",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			to := cliCtx.FromAddress
			symbol := args[1]

			msg := types.NewMsgIssueCNFT(symbol, to)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return client.PostCommands(cmd)[0]
}

func IssueCollectionFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection-issue-ft [from_key_or_address] [symbol] [name]",
		Short: "Create and sign an collection-issue-ft tx",
		Long: `
[Collective Fungible Token]
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

			msg := types.NewMsgIssueCFT(name, symbol, tokenURI, to, sdk.NewInt(supply), sdk.NewInt(decimals), mintable)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().Int64(flagTotalSupply, DefaultTotalSupply, "total supply")
	cmd.Flags().Int64(flagDecimals, DefaultDecimals, "set decimals")
	cmd.Flags().Bool(flagMintable, false, "set mintable")
	cmd.Flags().String(flagTokenURI, "", "set token-uri")

	return client.PostCommands(cmd)[0]
}

func MintCollectionNFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection-mint-nft [from_key_or_address] [to] [symbol] [token_type] [name]",
		Short: "Create and sign an collection-mint-nft tx",
		Long: `
[Collective NonFungible Token]
linkcli tx token collection-mint-nft [from_key_or_address] [symbol] [token_type] [name]
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

			msg := types.NewMsgMintCNFT(name, symbol, tokenURI, tokenType, from, to)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(flagTokenURI, "", "set token-uri")
	cmd.Flags().String(flagTokenType, "", "token-type for the nft")

	return client.PostCommands(cmd)[0]
}

func BurnCollectionNFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection-burn-nft [from_key_or_address] [symbol] [token_id]",
		Short: "Create and sign an collection-burn-nft tx",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			symbol := args[1]
			tokenID := args[2]

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgBurnCNFT(symbol, tokenID, cliCtx.GetFromAddress())
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func TransferCFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection-transfer-ft [from_key_or_address] [to_address] [symbol] [token_id] [amount]",
		Short: "Create and sign a tx transferring non-reserved collective fungible tokens",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			to, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[4])
			if !ok {
				return types.ErrInvalidAmount(types.DefaultCodespace, args[4])
			}

			msg := types.NewMsgTransferCFT(cliCtx.GetFromAddress(), to, args[2], args[3], amount)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

func TransferCNFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection-transfer-nft [from_key_or_address] [to_address] [symbol] [token_id]",
		Short: "Create and sign a tx transferring a collective non-fungible token",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			to, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgTransferCNFT(cliCtx.GetFromAddress(), to, args[2], args[3])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func TransferCFTFromTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection-transfer-ft-from [proxy_key_or_address] [approver_address] [to_address] [symbol] [token_id] [amount]",
		Short: "Create and sign a tx transferring non-reserved collective fungible tokens by approved proxy",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			approver, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			to, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[5])
			if !ok {
				return types.ErrInvalidAmount(types.DefaultCodespace, args[5])
			}

			msg := types.NewMsgTransferCFTFrom(cliCtx.GetFromAddress(), approver, to, args[3], args[4], amount)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

func TransferCNFTFromTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection-transfer-nft-from [proxy_key_or_address] [approver] [to_address] [symbol] [token_id]",
		Short: "Create and sign a tx transferring a collective non-fungible token by approved proxy",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			approver, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			to, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgTransferCNFTFrom(cliCtx.GetFromAddress(), approver, to, args[3], args[4])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func AttachTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection-attach [from_key_or_address] [symbol] [to_token_id] [token_id]",
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
		Use:   "collection-detach [from_key_or_address] [to_address] [symbol] [token_id]",
		Short: "Create and sign a tx detaching a token",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			to, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDetach(cliCtx.GetFromAddress(), to, args[2], args[3])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func MintCollectionFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection-mint-ft [from_key_or_address] [to] [symbol] [token-id] [amount]",
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

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgMintCFT(cliCtx.GetFromAddress(), to, linktype.NewCoinWithTokenIDs(linktype.NewCoinWithTokenID(symbol, tokenID, amount)))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func BurnCollectionFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection-burn-ft [from_key_or_address][symbol] [token-id] [amount]",
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

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgBurnCFT(cliCtx.GetFromAddress(), linktype.NewCoinWithTokenIDs(linktype.NewCoinWithTokenID(symbol, tokenID, amount)))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func AttachFromTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection-attach-from [proxy_key_or_address] [approver_address] [symbol] [to_token_id] [token_id]",
		Short: "Create and sign a tx attaching a token to other by approved proxy",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			approver, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgAttachFrom(cliCtx.GetFromAddress(), approver, args[2], args[3], args[4])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func DetachFromTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection-detach-from [proxy_key_or_address] [approver_address] [to_address] [symbol] [token_id]",
		Short: "Create and sign a tx detaching a token by approved proxy",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			approver, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			to, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgDetachFrom(cliCtx.GetFromAddress(), approver, to, args[3], args[4])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func ApproveCollectionTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection-approve [approver_key_or_address] [proxy_address] [symbol]",
		Short: "Create and sign a tx approve all token operations of a collection to a proxy",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			proxy, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgApproveCollection(cliCtx.GetFromAddress(), proxy, args[2])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func DisapproveCollectionTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection-disapprove [approver_key_or_address] [proxy_address] [symbol]",
		Short: "Create and sign a tx disapprove all token operations of a collection to a proxy",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			proxy, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDisapproveCollection(cliCtx.GetFromAddress(), proxy, args[2])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}
