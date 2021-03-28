package cli

import (
	"bufio"
	"errors"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/line/lbm-sdk/v2/x/collection/internal/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
	flagTokenType   = "token-type"
	flagTokenIndex  = "token-index"
)

const (
	DefaultDecimals    = 8
	DefaultTotalSupply = 1
	DefaultTokenType   = ""
	DefaultTokenIndex  = ""
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
		ModifyCmd(cdc),
	)
	return txCmd
}

func ModifyCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify [owner_address] [contract_id] [field] [new_value]",
		Short: "Create and sign a modify tx",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			contractID := args[1]
			field := args[2]
			newValue := args[3]
			tokenType := viper.GetString(flagTokenType)
			tokenIndex := viper.GetString(flagTokenIndex)

			msg := types.NewMsgModify(cliCtx.FromAddress, contractID, tokenType, tokenIndex,
				types.NewChanges(types.NewChange(field, newValue)))
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(flagTokenType, DefaultTokenType, "token type")
	cmd.Flags().String(flagTokenIndex, DefaultTokenIndex, "token index")

	return flags.PostCommands(cmd)[0]
}

func CreateCollectionTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [from_key_or_address] [name] [meta] [base_img_uri]",
		Short: "Create and sign an create collection tx",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			owner := cliCtx.FromAddress
			name := args[1]
			meta := args[2]
			baseImgURI := args[3]

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgCreateCollection(owner, name, meta, baseImgURI)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return flags.PostCommands(cmd)[0]
}

func IssueNFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-nft [from_key_or_address] [contract_id] [name] [meta]",
		Short: "Create and sign an issue-nft tx",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			to := cliCtx.FromAddress
			contractID := args[1]
			name := args[2]
			meta := args[3]

			msg := types.NewMsgIssueNFT(to, contractID, name, meta)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return flags.PostCommands(cmd)[0]
}

func IssueFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-ft [from_key_or_address] [contract_id] [to] [name] [meta]",
		Short: "Create and sign an issue-ft tx",
		Long: `
[Fungible Token]
linkcli tx collection issue-ft [from_key_or_address] [contract_id] [to] [name] [meta]
--decimals=[decimals]
--mintable=[mintable]
--total-supply=[initial amount of the token]
`,
		Args: cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			owner := cliCtx.FromAddress
			contractID := args[1]
			to, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}
			name := args[3]
			meta := args[4]
			supply := viper.GetInt64(flagTotalSupply)
			decimals := viper.GetInt64(flagDecimals)
			mintable := viper.GetBool(flagMintable)

			if decimals < 0 || decimals > 18 {
				return errors.New("invalid decimals. 0 <= decimals <= 18")
			}

			msg := types.NewMsgIssueFT(owner, to, contractID, name, meta, sdk.NewInt(supply), sdk.NewInt(decimals), mintable)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().Int64(flagTotalSupply, DefaultTotalSupply, "total supply")
	cmd.Flags().Int64(flagDecimals, DefaultDecimals, "set decimals")
	cmd.Flags().Bool(flagMintable, false, "set mintable")

	return flags.PostCommands(cmd)[0]
}

func MintNFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-nft [from_key_or_address] [contract_id] [to] [token_type:name:meta][,[token_type:name:meta]]",
		Short: "Create and sign an mint-nft tx",
		Long: `
[NonFungible Token]
linkcli tx collection mint-nft [from_key_or_address] [contract_id] [to] [token_type] [name] [meta]
`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			contractID := args[1]
			from := cliCtx.FromAddress
			to, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			mintNFTParamStrs := strings.Split(args[3], ",")

			mintNFTParams := make([]types.MintNFTParam, len(mintNFTParamStrs))
			for i, mintNFTParamStr := range mintNFTParamStrs {
				strs := strings.Split(mintNFTParamStr, ":")
				if len(strs) != 3 {
					return errors.New("invalid format: <token_type:name:meta>")
				}

				mintNFTParams[i] = types.NewMintNFTParam(strs[1], strs[2], strs[0])
			}

			msg := types.NewMsgMintNFT(from, contractID, to, mintNFTParams...)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(flagTokenType, "", "token-type for the nft")

	return flags.PostCommands(cmd)[0]
}

func BurnNFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-nft [from_key_or_address] [contract_id] [token_id]",
		Short: "Create and sign an burn-nft tx",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			contractID := args[1]
			tokenID := args[2]

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgBurnNFT(cliCtx.GetFromAddress(), contractID, tokenID)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return flags.PostCommands(cmd)[0]
}

func BurnNFTFromTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-nft-from [proxy_key_or_address] [contract_id] [from_address] [token_id]",
		Short: "Create and sign an burn-nft-from tx",
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

			tokenID := args[3]

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgBurnNFTFrom(cliCtx.GetFromAddress(), contractID, from, tokenID)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return flags.PostCommands(cmd)[0]
}

// nolint:dupl
func TransferFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-ft [from_key_or_address] [contract_id] [to_address] [amount]",
		Short: "Create and sign a tx transferring non-reserved collective fungible tokens",
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

			amount, err := types.ParseCoins(args[3])
			if err != nil {
				return sdkerrors.Wrap(types.ErrInvalidAmount, args[3])
			}

			msg := types.NewMsgTransferFT(cliCtx.GetFromAddress(), contractID, to, amount...)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

// nolint:dupl
func TransferNFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-nft [from_key_or_address] [contract_id] [to_address] [token_id][,[token_id]]",
		Short: "Create and sign a tx transferring a collective non-fungible token",
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

			tokenIDs := strings.Split(args[3], ",")

			msg := types.NewMsgTransferNFT(cliCtx.GetFromAddress(), contractID, to, tokenIDs...)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return flags.PostCommands(cmd)[0]
}

func TransferFTFromTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-ft-from [proxy_key_or_address] [contract_id] [from_address] [to_address] [amount]",
		Short: "Create and sign a tx transferring non-reserved collective fungible tokens by approved proxy",
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

			amount, err := types.ParseCoins(args[4])
			if err != nil {
				return sdkerrors.Wrap(types.ErrInvalidAmount, args[4])
			}

			msg := types.NewMsgTransferFTFrom(cliCtx.GetFromAddress(), contractID, from, to, amount...)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

// nolint:dupl
func TransferNFTFromTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-nft-from [proxy_key_or_address] [contract_id] [from_address] [to_address] [token_id]",
		Short: "Create and sign a tx transferring a collective non-fungible token by approved proxy",
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

			tokenIDs := strings.Split(args[4], ",")

			msg := types.NewMsgTransferNFTFrom(cliCtx.GetFromAddress(), contractID, from, to, tokenIDs...)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return flags.PostCommands(cmd)[0]
}

func AttachTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attach [from_key_or_address] [contract_id] [to_token_id] [token_id]",
		Short: "Create and sign a tx attaching a token to other",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			msg := types.NewMsgAttach(cliCtx.GetFromAddress(), args[1], args[2], args[3])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return flags.PostCommands(cmd)[0]
}

func DetachTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach [from_key_or_address] [contract_id] [token_id]",
		Short: "Create and sign a tx detaching a token",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			msg := types.NewMsgDetach(cliCtx.GetFromAddress(), args[1], args[2])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return flags.PostCommands(cmd)[0]
}

// nolint:dupl
func MintFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-ft [from_key_or_address] [contract_id] [to] [amount]",
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

			amount, err := types.ParseCoins(args[3])
			if err != nil {
				return sdkerrors.Wrap(types.ErrInvalidAmount, args[3])
			}

			msg := types.NewMsgMintFT(cliCtx.GetFromAddress(), contractID, to, amount...)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return flags.PostCommands(cmd)[0]
}

func BurnFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-ft [from_key_or_address] [contract_id] [token-id] [amount]",
		Short: "Create and sign a mint token tx",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)
			contractID := args[1]
			tokenID := args[2]
			if err := types.ValidateDenom(tokenID); err != nil {
				return errors.New("invalid tokenID")
			}
			amount, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return errors.New("invalid amount")
			}

			msg := types.NewMsgBurnFT(cliCtx.GetFromAddress(), contractID, types.NewCoin(tokenID, amount))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return flags.PostCommands(cmd)[0]
}

// nolint:dupl
func BurnFTFromTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-ft-from [proxy_key_or_address] [contract_id] [from_address] [token-id] [amount]",
		Short: "Create and sign a mint token tx",
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
			tokenID := args[3]
			if err := types.ValidateDenom(tokenID); err != nil {
				return errors.New("invalid tokenID")
			}
			amount, ok := sdk.NewIntFromString(args[4])
			if !ok {
				return errors.New("invalid amount")
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgBurnFTFrom(cliCtx.GetFromAddress(), contractID, from, types.NewCoin(tokenID, amount))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return flags.PostCommands(cmd)[0]
}

func AttachFromTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attach-from [proxy_key_or_address] [contract_id] [from_address] [to_token_id] [token_id]",
		Short: "Create and sign a tx attaching a token to other by approved proxy",
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

			msg := types.NewMsgAttachFrom(cliCtx.GetFromAddress(), contractID, from, args[3], args[4])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return flags.PostCommands(cmd)[0]
}

// nolint:dupl
func DetachFromTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach-from [proxy_key_or_address] [contract_id] [from_address] [token_id]",
		Short: "Create and sign a tx detaching a token by approved proxy",
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

			msg := types.NewMsgDetachFrom(cliCtx.GetFromAddress(), contractID, from, args[3])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return flags.PostCommands(cmd)[0]
}

// nolint:dupl
func ApproveCollectionTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve [approver_key_or_address] [contract_id] [proxy_address]",
		Short: "Create and sign a tx approve all token operations of a collection to a proxy",
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

// nolint:dupl
func DisapproveCollectionTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disapprove [approver_key_or_address] [contract_id] [proxy_address]",
		Short: "Create and sign a tx disapprove all token operations of a collection to a proxy",
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

			msg := types.NewMsgDisapprove(cliCtx.GetFromAddress(), contractID, proxy)
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
			if !perm.Validate() {
				return errors.New("permission invalid")
			}

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
			if !perm.Validate() {
				return errors.New("permission invalid")
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgRevokePermission(cliCtx.GetFromAddress(), contractID, perm)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return flags.PostCommands(cmd)[0]
}
