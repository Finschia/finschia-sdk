package cli

import (
	"errors"
	linktype "github.com/link-chain/link/types"
	"github.com/link-chain/link/x/token/internal/types"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/link-chain/link/client"
)

var (
	flagTotalSupply = "total-supply"
	flagDecimals    = "decimals"
	flagMintable    = "mintable"
	flagTokenURI    = "token-uri"
	flagTokenID     = "token-id"
	flagFungible    = "fungible"
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
		IssueTxCmd(cdc),
		MintTxCmd(cdc),
		BurnTxCmd(cdc),
		GrantPermTxCmd(cdc),
		RevokePermTxCmd(cdc),
	)
	return txCmd
}

func IssueTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue [from_key_or_address] [symbol] [name]",
		Short: "Create and sign an issue token tx",
		Args:  cobra.ExactArgs(3),
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
			tokenID := viper.GetString(flagTokenID)
			fungible := viper.GetBool(flagFungible)
			aas := viper.GetBool(flagAAS)

			if aas {
				symbol = symbol + to.String()[len(to.String())-linktype.AccAddrSuffixLen:]
			}

			if err := linktype.ValidateSymbolUserDefined(symbol); err != nil {
				return err
			}

			if decimals < 0 || decimals > 18 {
				return errors.New("invalid decimals. 0 <= decimals <= 18")
			}

			var msg sdk.Msg
			if len(tokenID) == 0 {
				if !fungible {
					msg = types.NewMsgIssueNFT(name, symbol, tokenURI, to)
				} else {
					msg = types.NewMsgIssue(name, symbol, tokenURI, to, sdk.NewInt(supply), sdk.NewInt(decimals), mintable)
				}
			} else {
				if !fungible {
					msg = types.NewMsgIssueNFTCollection(name, symbol, tokenURI, to, tokenID)
				} else {
					msg = types.NewMsgIssueCollection(name, symbol, tokenURI, to, sdk.NewInt(supply), sdk.NewInt(decimals), mintable, tokenID)
				}

			}
			// build and sign the transaction, then broadcast to Tendermint
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().Int64(flagTotalSupply, DefaultTotalSupply, "total supply")
	cmd.Flags().Int64(flagDecimals, DefaultDecimals, "set decimals")
	cmd.Flags().Bool(flagMintable, false, "set mintable")
	cmd.Flags().String(flagTokenURI, "", "set token-uri")
	cmd.Flags().String(flagTokenID, "", "token-id in the collection")
	cmd.Flags().Bool(flagFungible, true, "set fungible. it overwrite values of decimals, mintable to 0 when set false")
	cmd.Flags().Bool(flagAAS, true, "attach address suffix to symbol")

	return client.PostCommands(cmd)[0]
}

func MintTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [to_key_or_address] [amount_with_denom]",
		Short: "Create and sign a mint token tx",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			// parse coins trying to be sent
			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgMint(cliCtx.GetFromAddress(), coins)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func BurnTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [from_key_or_address] [amount_with_denom]",
		Short: "Create and sign a burn token tx",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			// parse coins trying to be sent
			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgBurn(cliCtx.GetFromAddress(), coins)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func GrantPermTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grant [from_key_or_address] [to] [token] [action]",
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
		Use:   "revoke [from_key_or_address] [token] [action]",
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
