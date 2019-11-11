package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/link-chain/link/x/lrc3/internal/types"
)

// Edit metadata flags
const (
	flagTokenURI = "tokenURI"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	nftTxCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "NFT transactions subcommands",
	}

	nftTxCmd.AddCommand(client.PostCommands(
		GetCmdInit(cdc),
		GetCmdMintNFT(cdc),
		GetCmdBurnNFT(cdc),
		GetCmdTransferNFT(cdc),
		GetCmdApprove(cdc),
		GetCmdSetApprovalForAll(cdc),
	)...)

	return nftTxCmd
}

func GetCmdInit(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "init [denom]",
		Short: "Initialize a Non Fungible Token",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Specify Name and Symbol to generate Non Fungible Token.
			The account that executed Generate is set to Service Operator by default.
Example:
$ %s tx %s init sample_dapp --from mykey`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authtypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			denom := args[0]

			msg := types.NewMsgInit(denom, cliCtx.GetFromAddress())
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdMintNFT is the CLI command for a MintNFT transaction
func GetCmdMintNFT(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "mint [denom] [tokenURI] [to]",
		Short: "Mint an NFT and set the owner to the recipient",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Mint an NFT from a given collection that has a
			specific tokenId and set the ownership to a specific address.
Example:
$ %s tx %s mint sample_dapp cosmos1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p --from mykey
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authtypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			denom := args[0]
			tokenURI := viper.GetString(args[1])

			recipient, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}
			msg := types.NewMsgMintNFT(cliCtx.GetFromAddress(), recipient, denom, tokenURI)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdBurnNFT is the CLI command for a MintNFT transaction
func GetCmdBurnNFT(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "burn [denom] [tokenID]",
		Short: "burn an NFT",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Burn (i.e permanently delete) an NFT from a given collection that has a 
			specific tokenId.

Example:
$ %s tx %s burn sample_dapp 0 --from mykey
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authtypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			denom := args[0]
			tokenID := args[1]

			msg := types.NewMsgBurn(cliCtx.GetFromAddress(), tokenID, denom, cliCtx.GetFromAddress())
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdTransferNFT is the CLI command for sending a TransferNFT transaction
func GetCmdTransferNFT(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "transfer [sender] [recipient] [denom] [tokenID]",
		Short: "transfer a NFT to a recipient",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Transfer a NFT from a given collection that has a 
			specific tokenId to a specific recipient.

Example:
$ %s tx %s transfer cosmos1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p cosmos1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm \
sample_dapp 0 --from mykey
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authtypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			sender, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			recipient, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			denom := args[2]
			tokenID := args[3]

			msg := types.NewMsgTransfer(sender, recipient, denom, tokenID, cliCtx.GetFromAddress())
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdEditNFTMetadata is the CLI command for sending an EditMetadata transaction
func GetCmdEditNFTMetadata(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-metadata [denom] [tokenID]",
		Short: "edit the metadata of an NFT",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Edit the metadata of an NFT from a given collection that has a 
			specific tokenId.

Example:
$ %s tx %s edit-metadata sample_dapp 0 --tokenURI path_to_token_URI_JSON --from mykey
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authtypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			denom := args[0]
			tokenID := args[1]
			tokenURI := viper.GetString(flagTokenURI)

			msg := types.NewMsgEditMetadata(cliCtx.GetFromAddress(), tokenID, denom, tokenURI, cliCtx.GetFromAddress())
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagTokenURI, "", "Extra properties available for querying")
	return cmd
}

func GetCmdApprove(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "approve [denom] [tokenId] [to]",
		Short: "Give authority to transfer NFT to others",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Gives 'to' address the right to transfer the NFT corresponding to tokenId to another person.
Example:
$ %s tx %s approve sample_dapp 0 cosmos12u008kswqd54tu3uupqsk2zzq4hg2t9nlf56rf --from mykey`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authtypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			denom := args[0]
			tokenId := args[1]
			to, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgApprove(cliCtx.GetFromAddress(), denom, tokenId, to)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdSetApprovalForAll(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "setApprovalForAll [denom] [operator] [approved]",
		Short: "Set the authority to control all NFT assets owned by the requested user",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Set so that the control of all NFT assets held by the requesting user cannot be set to 'operator'.
Example:
$ %s tx %s setApprovalForAll sample_dapp cosmos12u008kswqd54tu3uupqsk2zzq4hg2t9nlf56rf cosmos19fw638zr5rrwj6cgl2r2cnwxxtm0azgza79vm3 true --from mykey`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authtypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			symbol := args[0]

			operator, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			approved, err := strconv.ParseBool(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgSetApprovalForAll(symbol, cliCtx.GetFromAddress(), operator, approved)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
