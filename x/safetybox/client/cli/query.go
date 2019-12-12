package cli

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/link-chain/link/client"
	"github.com/link-chain/link/x/safetybox/internal/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the safety box module",
		Aliases:                    []string{"sb"},
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		GetSafetyBoxCmd(cdc),
		GetRoleCmd(cdc),
	)

	return cmd
}

func GetSafetyBoxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [safety_box_id]",
		Short: "Query the safety box",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			safetyBoxGetter := types.NewSafetyBoxRetriever(cliCtx)

			id := args[0]
			sb, err := safetyBoxGetter.GetSafetyBox(id)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(sb)
		},
	}

	return client.GetCommands(cmd)[0]
}

func GetRoleCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role [safety_box_id] [owner|operator|allocator|issuer|returner] [addr]",
		Short: "Get if the account has the role on the safety box",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewCLIContext().WithCodec(cdc)
			permGetter := types.NewAccountPermissionRetriever(cliCtx)

			id := args[0]
			role := args[1]
			addr, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			pms, err := permGetter.GetAccountPermissions(id, role, addr)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(pms)
		},
	}

	return client.GetCommands(cmd)[0]
}
