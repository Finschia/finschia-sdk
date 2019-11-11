package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/link-chain/link/x/lrc3/internal/types"

	nft "github.com/link-chain/link/x/nft"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	nftQueryCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the NFT module",
	}

	nftQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryLRC3(queryRoute, cdc),
		GetCmdQueryAllLRC3(queryRoute, cdc),
		GetCmdGetApproved(queryRoute, cdc),
		GetCmdIsApprovedForAll(queryRoute, cdc),
		GetCmdBalanceOf(queryRoute, cdc),
		GetCmdOwnerOf(queryRoute, cdc),
	)...)

	return nftQueryCmd
}

// GetCmdQueryLRC3 queries the lrc3
func GetCmdQueryLRC3(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get [denom]",
		Short: "Get LRC3 token information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get the LRC3 Token information by symbol.
Example:
$ %s query %s get sample_dapp
`, version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			denom := args[0]

			params := types.NewQueryLRC3Params(denom)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/lrc3", queryRoute), bz)
			if err != nil {
				return err
			}

			var out nft.Collection
			err = cdc.UnmarshalJSON(res, &out)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdQueryLRC3 queries the lrc3
func GetCmdQueryAllLRC3(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "getAll",
		Short: "Get ALL LRC3 token symbol",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get All LRC3 Token Contract Address information.
Example:
$ %s query %s getAll`, version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/allLRC3", queryRoute), nil)
			if err != nil {
				return err
			}

			var out nft.Collections
			err = cdc.UnmarshalJSON(res, &out)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdGetApproved queries a approve for token
func GetCmdGetApproved(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "getApproved [denom] [tokenId]",
		Short: "Query users approved by one NFT",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Queries whether the access permission of all tokens owned by owner is granted to operator.
Example:
$ %s query %s getApproved sample_dapp 0`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			denom := args[0]
			tokenId := args[1]

			params := types.NewQueryApproveParams(denom, tokenId)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/getApproved", queryRoute), bz)
			if err != nil {
				return err
			}

			var out types.Approval
			err = cdc.UnmarshalJSON(res, &out)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdIsApprovedForAll queries whether '_operator' has authority to control all NFTs owned by '_owner'.
func GetCmdIsApprovedForAll(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "isApprovedForAll [denom] [owner] [operator]",
		Short: "Query whether you have the authority to control all NFTs you own",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query whether 'operator' has the authority to control all NFTs owned by owner.
Example:
$ %s query %s isApprovedForAll sample_dapp cosmos12u008kswqd54tu3uupqsk2zzq4hg2t9nlf56rf cosmos19fw638zr5rrwj6cgl2r2cnwxxtm0azgza79vm3`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := args[0]
			owner, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			operator, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			params := types.NewQueryOperatorApproveParams(denom, owner, operator)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/isApprovedForAll", queryRoute), bz)
			if err != nil {
				return err
			}

			var out types.OperatorApprovals
			err = cdc.UnmarshalJSON(res, &out)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdBalanceOf queries a balance of account
func GetCmdBalanceOf(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "balanceOf [denom] [owner]",
		Short: "Query balance of user corresponding to address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query balance of user corresponding to address.
Example:
$ %s query %s balanceOf sample_dapp cosmos12u008kswqd54tu3uupqsk2zzq4hg2t9nlf56rf`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			denom := args[0]
			owner, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			params := types.NewQueryBalanceParams(owner, denom)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/balanceOf", queryRoute), bz)
			if err != nil {
				return err
			}

			var out types.TokenBalance
			err = cdc.UnmarshalJSON(res, &out)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdOwnerOf queries a single NFTs from a collection
func GetCmdOwnerOf(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "ownerOf [denom] [tokenId]",
		Short: "Query the token owner address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the token owner address for tokenId.
Example:
$ %s query %s ownerOf sample_dapp 0
`, version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			denom := args[0]
			tokenId := args[1]

			params := types.NewQueryOwnerOfParams(denom, tokenId)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/ownerOf", queryRoute), bz)
			if err != nil {
				return err
			}

			var out types.TokenOwner
			err = cdc.UnmarshalJSON(res, &out)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(out)
		},
	}
}
