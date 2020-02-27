package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/client"
	clienttypes "github.com/line/link/x/collection/client/internal/types"
	"github.com/line/link/x/collection/internal/types"
	"github.com/spf13/cobra"
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
		GetTokenCmd(cdc),
		GetTokensCmd(cdc),
		GetCollectionCmd(cdc),
		GetCollectionsCmd(cdc),
		GetTokenTotalCmd(cdc),
		GetTokenCountCmd(cdc),
		GetPermsCmd(cdc),
		GetParentCmd(cdc),
		GetRootCmd(cdc),
		GetChildrenCmd(cdc),
		GetIsApproved(cdc),
	)

	return cmd
}

func GetBalanceCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance [symbol] [token_id] [addr]",
		Short: "Query balance of the account",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			symbol := args[0]
			tokenID := args[1]
			addr, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			supply, height, err := retriever.GetAccountBalance(cliCtx, symbol, tokenID, addr)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)
			return cliCtx.PrintOutput(supply)
		},
	}

	return client.GetCommands(cmd)[0]
}

func GetCollectionCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection [symbol]",
		Short: "Query collection",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			symbol := args[0]
			collection, height, err := retriever.GetCollection(cliCtx, symbol)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)
			return cliCtx.PrintOutput(collection)
		},
	}

	return client.GetCommands(cmd)[0]
}

func GetCollectionsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collections",
		Short: "Query all collections",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			collections, height, err := retriever.GetCollections(cliCtx)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)
			return cliCtx.PrintOutput(collections)
		},
	}

	return client.GetCommands(cmd)[0]
}
func GetTokenCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token [symbol] [token-id]",
		Short: "Query collection token with collection symbol and token's token-id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			symbol := args[0]
			tokenID := args[1]
			token, height, err := retriever.GetToken(cliCtx, symbol, tokenID)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)

			return cliCtx.PrintOutput(token)
		},
	}

	return client.GetCommands(cmd)[0]
}

func GetTokensCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens [symbol]",
		Short: "Query all collection tokens with collection symbol",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			symbol := args[0]

			tokens, height, err := retriever.GetTokens(cliCtx, symbol)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)
			return cliCtx.PrintOutput(tokens)
		},
	}

	return client.GetCommands(cmd)[0]
}

func GetTokenTotalCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total [supply|mint|burn] [symbol] [token-id]",
		Short: "Query supply/mint/burn of collection token with collection symbol and tokens's token-id.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			target := args[0]
			symbol := args[1]
			tokenID := args[2]

			supply, height, err := retriever.GetTotal(cliCtx, symbol, tokenID, target)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)
			return cliCtx.PrintOutput(supply)
		},
	}

	return client.GetCommands(cmd)[0]
}
func GetTokenCountCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "count [symbol] [token-type]",
		Short: "Query count of collection tokens with collection symbol and the base-id.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			symbol := args[0]
			baseID := args[1]

			supply, height, err := retriever.GetCollectionNFTCount(cliCtx, symbol, baseID)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)
			return cliCtx.PrintOutput(supply)
		},
	}

	return client.GetCommands(cmd)[0]
}

func GetPermsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "perm [addr]",
		Short: "Get Permission of the Account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			pms, height, err := retriever.GetAccountPermission(cliCtx, addr)
			if err != nil {
				return err
			}

			cliCtx = cliCtx.WithHeight(height)
			return cliCtx.PrintOutput(pms)
		},
	}

	return client.GetCommands(cmd)[0]
}

func GetParentCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parent [symbol] [token-id]",
		Short: "Query parent token with symbol and token-id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			tokenGetter := clienttypes.NewRetriever(cliCtx)

			symbol := args[0]
			tokenID := args[1]

			if err := tokenGetter.EnsureExists(cliCtx, symbol, tokenID); err != nil {
				return err
			}

			token, _, err := tokenGetter.GetParent(cliCtx, symbol, tokenID)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(token)
		},
	}

	return client.GetCommands(cmd)[0]
}

func GetRootCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "root [symbol] [token-id]",
		Short: "Query root token with symbol and token-id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			tokenGetter := clienttypes.NewRetriever(cliCtx)

			symbol := args[0]
			tokenID := args[1]

			if err := tokenGetter.EnsureExists(cliCtx, symbol, tokenID); err != nil {
				return err
			}

			token, _, err := tokenGetter.GetRoot(cliCtx, symbol, tokenID)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(token)
		},
	}

	return client.GetCommands(cmd)[0]
}

func GetChildrenCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "children [symbol] [token-id]",
		Short: "Query children tokens with symbol and token-id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			tokenGetter := clienttypes.NewRetriever(cliCtx)

			symbol := args[0]
			tokenID := args[1]

			if err := tokenGetter.EnsureExists(cliCtx, symbol, tokenID); err != nil {
				return err
			}

			tokens, _, err := tokenGetter.GetChildren(cliCtx, symbol, tokenID)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(tokens)
		},
	}

	return client.GetCommands(cmd)[0]
}

type Approved struct {
	Approved bool
}

func (a Approved) String() string {
	return string(codec.MustMarshalJSONIndent(types.ModuleCdc, a))
}

var _ fmt.Stringer = (*Approved)(nil)

func GetIsApproved(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approved [proxy] [approver] [symbol]",
		Short: "Query whether a proxy is approved by approver on a collection",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			retriever := clienttypes.NewRetriever(cliCtx)

			proxy, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			approver, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			approved, height, err := retriever.IsApproved(cliCtx, proxy, approver, args[2])
			if err != nil {
				return err
			}

			return cliCtx.WithHeight(height).PrintOutput(Approved{Approved: approved})
		},
	}

	return client.GetCommands(cmd)[0]
}
