package cli

import (
	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/context"
	"github.com/line/lbm-sdk/client/flags"
	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
	clienttypes "github.com/line/lbm-sdk/x/token/client/internal/types"
	"github.com/line/lbm-sdk/x/token/internal/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the token module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		GetTokenCmd(cdc),
		GetBalanceCmd(cdc),
		GetTotalCmd(cdc),
		GetPermsCmd(cdc),
		GetIsApprovedCmd(cdc),
		GetApproversCmd(cdc),
	)

	return cmd
}

func GetTokenCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token [contract_id]",
		Short: "Query token with its contract_id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			contractID := args[0]
			token, height, err := retriever.GetToken(cliCtx, contractID)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)

			return cliCtx.PrintOutput(token)
		},
	}

	return flags.GetCommands(cmd)[0]
}

func GetBalanceCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance [contract_id] [addr]",
		Short: "Query balance of the account",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			contractID := args[0]
			addr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			supply, height, err := retriever.GetAccountBalance(cliCtx, contractID, addr)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)
			return cliCtx.PrintOutput(supply)
		},
	}

	return flags.GetCommands(cmd)[0]
}

func GetTotalCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total [supply|mint|burn] [contract_id] ",
		Short: "Query total supply/mint/burn of token",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			target := args[0]
			contractID := args[1]

			supply, height, err := retriever.GetTotal(cliCtx, contractID, target)

			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)
			return cliCtx.PrintOutput(supply)
		},
	}

	return flags.GetCommands(cmd)[0]
}

func GetPermsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "perm [addr] [contract_id]",
		Short: "Get Permission of the Account",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			contractID := args[1]
			pms, height, err := retriever.GetAccountPermission(cliCtx, contractID, addr)
			if err != nil {
				return err
			}
			cliCtx = cliCtx.WithHeight(height)
			return cliCtx.PrintOutput(pms)
		},
	}

	return flags.GetCommands(cmd)[0]
}

func GetIsApprovedCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approved [contract_id] [proxy] [approver]",
		Short: "Query whether a proxy is approved by approver on a token",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			contractID := args[0]

			proxy, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			approver, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			approved, height, err := retriever.IsApproved(cliCtx, contractID, proxy, approver)
			if err != nil {
				return err
			}

			return cliCtx.WithHeight(height).PrintOutput(approved)
		},
	}

	return flags.GetCommands(cmd)[0]
}

func GetApproversCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approvers [contract_id] [proxy]",
		Short: "Query approvers by the proxy",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			contractID := args[0]

			proxy, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			approvers, height, err := retriever.GetApprovers(cliCtx, contractID, proxy)
			if err != nil {
				return err
			}

			return cliCtx.WithHeight(height).PrintOutput(approvers)
		},
	}

	return flags.GetCommands(cmd)[0]
}
