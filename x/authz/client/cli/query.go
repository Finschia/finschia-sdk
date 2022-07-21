package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/flags"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/version"
	"github.com/line/lbm-sdk/x/authz"
	bank "github.com/line/lbm-sdk/x/bank/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	authorizationQueryCmd := &cobra.Command{
		Use:                        authz.ModuleName,
		Short:                      "Querying commands for the authz module",
		Long:                       "",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	authorizationQueryCmd.AddCommand(
		GetCmdQueryGrants(),
	)

	return authorizationQueryCmd
}

// GetCmdQueryGrants implements the query authorization command.
func GetCmdQueryGrants() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grants [granter-addr] [grantee-addr] [msg-type-url]?",
		Args:  cobra.RangeArgs(2, 3),
		Short: "query grants for a granter-grantee pair and optionally a msg-type-url",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query authorization grants for a granter-grantee pair. If msg-type-url
is set, it will select grants only for that msg type.
Examples:
$ %s query %s grants link1skj.. link1skjwj..
$ %s query %s grants link1skjw.. link1skjwj.. %s
`,
				version.AppName, authz.ModuleName,
				version.AppName, authz.ModuleName, bank.SendAuthorization{}.MsgTypeURL()),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := authz.NewQueryClient(clientCtx)

			granter, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			grantee, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			var msgAuthorized = ""
			if len(args) >= 3 {
				msgAuthorized = args[2]
			}
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Grants(
				cmd.Context(),
				&authz.QueryGrantsRequest{
					Granter:    granter.String(),
					Grantee:    grantee.String(),
					MsgTypeUrl: msgAuthorized,
					Pagination: pageReq},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "authorizations")
	return cmd
}
