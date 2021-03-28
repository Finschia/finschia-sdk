package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/line/lbm-sdk/v2/x/collection/client/internal/types"
	"github.com/line/lbm-sdk/v2/x/collection/internal/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the collection module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetBalanceCmd(cdc),
		GetBalancesCmd(cdc),
		GetTokenCmd(cdc),
		GetTokensCmd(cdc),
		GetTokenTypeCmd(cdc),
		GetTokenTypesCmd(cdc),
		GetCollectionCmd(cdc),
		GetTokenTotalCmd(cdc),
		GetTokenCountCmd(cdc),
		GetPermsCmd(cdc),
		GetParentCmd(cdc),
		GetRootCmd(cdc),
		GetChildrenCmd(cdc),
		GetApproversCmd(cdc),
		GetIsApprovedCmd(cdc),
	)

	return cmd
}

func GetBalanceCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance [contract_id] [token_id] [addr]",
		Short: "Query balance of the account",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			contractID := args[0]
			tokenID := args[1]
			addr, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			supply, height, err := retriever.GetAccountBalance(cliCtx, contractID, tokenID, addr)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)
			return cliCtx.PrintOutput(supply)
		},
	}

	return flags.GetCommands(cmd)[0]
}

func GetBalancesCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balances [contract_id] [addr]",
		Short: "Query balances of the account for each token_id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			contractID := args[0]
			addr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			coins, height, err := retriever.GetAccountBalances(cliCtx, contractID, addr)
			if err != nil {
				return err
			}
			cliCtx = cliCtx.WithHeight(height)
			return cliCtx.PrintOutput(coins)
		},
	}

	return flags.GetCommands(cmd)[0]
}

func GetCollectionCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection [contract_id]",
		Short: "Query collection",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			contractID := args[0]
			collection, height, err := retriever.GetCollection(cliCtx, contractID)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)
			return cliCtx.PrintOutput(collection)
		},
	}

	return flags.GetCommands(cmd)[0]
}

func GetTokenTypeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokentype [contract_id] [token-type]",
		Short: "Query collection token-type with collection contract_id and token-type",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			contractID := args[0]
			tokenTypeID := args[1]
			tokenType, height, err := retriever.GetTokenType(cliCtx, contractID, tokenTypeID)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)

			return cliCtx.PrintOutput(tokenType)
		},
	}

	return flags.GetCommands(cmd)[0]
}

func GetTokenTypesCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokentypes [contract_id]",
		Short: "Query all collection token-types with collection contract_id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			contractID := args[0]

			tokenTypes, height, err := retriever.GetTokenTypes(cliCtx, contractID)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)
			return cliCtx.PrintOutput(tokenTypes)
		},
	}

	return flags.GetCommands(cmd)[0]
}

func GetTokenCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token [contract_id] [token_id]",
		Short: "Query collection token with collection contractID and token's token_id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			contractID := args[0]
			tokenID := args[1]
			token, height, err := retriever.GetToken(cliCtx, contractID, tokenID)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)

			return cliCtx.PrintOutput(token)
		},
	}

	return flags.GetCommands(cmd)[0]
}

func GetTokensCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens [contract_id]",
		Short: "Query all collection tokens with collection contractID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			contractID := args[0]
			var tokens types.Tokens
			var height int64
			var err error

			tokenType := viper.GetString(flagTokenType)
			if len(tokenType) > 0 {
				tokens, height, err = retriever.GetTokensWithTokenType(cliCtx, contractID, tokenType)
			} else {
				tokens, height, err = retriever.GetTokens(cliCtx, contractID)
			}

			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)
			return cliCtx.PrintOutput(tokens)
		},
	}
	cmd.Flags().String(flagTokenType, DefaultTokenType, "get tokens belong to the token-type")

	return flags.GetCommands(cmd)[0]
}

func GetTokenTotalCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total [supply|mint|burn] [contract_id] [token_id]",
		Short: "Query supply/mint/burn of collection token with contract-id and tokens's token_id.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			target := args[0]
			contractID := args[1]
			tokenID := args[2]

			supply, height, err := retriever.GetTotal(cliCtx, contractID, tokenID, target)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)
			return cliCtx.PrintOutput(supply)
		},
	}

	return flags.GetCommands(cmd)[0]
}

func GetTokenCountCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "count [total|mint|burn] [contract_id] [token_type]",
		Short: "Query count of collection tokens with collection contractID and the type_type.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			target := args[0]
			switch target {
			case "total":
				target = types.QueryNFTCount
			case "mint":
				target = types.QueryNFTMint
			case "burn":
				target = types.QueryNFTBurn
			default:
				return fmt.Errorf("argument is not total, mint, or burn %s", target)
			}

			contractID := args[1]
			tokenType := args[2]

			supply, height, err := retriever.GetCollectionNFTCount(cliCtx, contractID, tokenType, target)
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

func GetParentCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parent [contract_id] [token-id]",
		Short: "Query parent token with contractID and token-id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			tokenGetter := clienttypes.NewRetriever(cliCtx)

			contractID := args[0]
			tokenID := args[1]

			if err := tokenGetter.EnsureExists(cliCtx, contractID, tokenID); err != nil {
				return err
			}

			token, _, err := tokenGetter.GetParent(cliCtx, contractID, tokenID)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(token)
		},
	}

	return flags.GetCommands(cmd)[0]
}

func GetRootCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "root [contract_id] [token-id]",
		Short: "Query root token with contractID and token-id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			tokenGetter := clienttypes.NewRetriever(cliCtx)

			contractID := args[0]
			tokenID := args[1]

			if err := tokenGetter.EnsureExists(cliCtx, contractID, tokenID); err != nil {
				return err
			}

			token, _, err := tokenGetter.GetRoot(cliCtx, contractID, tokenID)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(token)
		},
	}

	return flags.GetCommands(cmd)[0]
}

func GetChildrenCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "children [contract_id] [token-id]",
		Short: "Query children tokens with contractID and token-id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			tokenGetter := clienttypes.NewRetriever(cliCtx)

			contractID := args[0]
			tokenID := args[1]

			if err := tokenGetter.EnsureExists(cliCtx, contractID, tokenID); err != nil {
				return err
			}

			tokens, _, err := tokenGetter.GetChildren(cliCtx, contractID, tokenID)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(tokens)
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

func GetIsApprovedCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approved [contract_id] [proxy] [approver]",
		Short: "Query whether a proxy is approved by approver on a collection",
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
