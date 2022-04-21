package cli

import (
	"context"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/flags"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

// NewQueryCmd returns the parent command for all x/foundation CLi query commands.
func NewQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   foundation.ModuleName,
		Short: "Querying commands for the foundation module",
	}

	cmd.AddCommand(
		NewQueryCmdParams(),
		NewQueryCmdValidatorAuth(),
		NewQueryCmdValidatorAuths(),
		NewQueryCmdTreasury(),
		NewQueryCmdFoundationInfo(),
		NewQueryCmdMember(),
		NewQueryCmdMembers(),
		NewQueryCmdProposal(),
		NewQueryCmdProposals(),
		NewQueryCmdVote(),
		NewQueryCmdVotes(),
		NewQueryCmdTallyResult(),
	)

	return cmd
}

// NewQueryCmdParams returns the query foundation parameters command.
func NewQueryCmdParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query foundation params",
		Long:  "Gets the current parameters of foundation",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := foundation.NewQueryClient(clientCtx)

			params := foundation.QueryParamsRequest{}
			res, err := queryClient.Params(context.Background(), &params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryCmdValidatorAuth returns validator authorization by foundation
func NewQueryCmdValidatorAuth() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-auth [validator-address]",
		Short: "Query validator authorization",
		Long:  "Gets validator authorization by foundation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := foundation.NewQueryClient(clientCtx)

			valAddr := args[0]
			if err = sdk.ValidateValAddress(valAddr); err != nil {
				return err
			}

			params := foundation.QueryValidatorAuthRequest{ValidatorAddress: valAddr}
			res, err := queryClient.ValidatorAuth(context.Background(), &params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryCmdValidatorAuths returns validator authorizations by foundation
func NewQueryCmdValidatorAuths() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-auths",
		Short: "Query validator authorizations",
		Long:  "Gets validator authorizations by foundation",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := foundation.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := foundation.QueryValidatorAuthsRequest{Pagination: pageReq}
			res, err := queryClient.ValidatorAuths(context.Background(), &params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "validator auths")

	return cmd
}

// NewQueryCmdTreasury returns the amount of coins in the foundation treasury
func NewQueryCmdTreasury() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "treasury",
		Short: "Query foundation treasury",
		Long:  "Gets the amount of coins in the foundation treasury",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := foundation.NewQueryClient(clientCtx)

			req := foundation.QueryTreasuryRequest{}
			res, err := queryClient.Treasury(context.Background(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}


// NewQueryCmdFoundationInfo returns the information of the foundation.
func NewQueryCmdFoundationInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "foundation-info",
		Args:  cobra.NoArgs,
		Short: "Query the foundation information",
		Long: `Query the foundation information
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := foundation.NewQueryClient(clientCtx)

			req := foundation.QueryFoundationInfoRequest{}
			res, err := queryClient.FoundationInfo(context.Background(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// NewQueryCmdMember returns a member of the foundation.
func NewQueryCmdMember() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Query a foundation member",
		Long: `Query a foundation member
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := foundation.NewQueryClient(clientCtx)

			address := args[0]
			if err = sdk.ValidateValAddress(address); err != nil {
				return err
			}

			req := foundation.QueryMemberRequest{Address: address}
			res, err := queryClient.Member(context.Background(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// NewQueryCmdMembers returns the members of the foundation.
func NewQueryCmdMembers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "members",
		Args:  cobra.NoArgs,
		Short: "Query the foundation members",
		Long: `Query the foundation members
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := foundation.NewQueryClient(clientCtx)

			req := foundation.QueryMembersRequest{}
			res, err := queryClient.Members(context.Background(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// NewQueryCmdProposal returns a proposal baesd on proposal id.
func NewQueryCmdProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query a proposal",
		Long: `Query a proposal
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := foundation.NewQueryClient(clientCtx)

			proposalId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			req := foundation.QueryProposalRequest{ProposalId: proposalId}
			res, err := queryClient.Proposal(context.Background(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// NewQueryCmdProposals returns all proposals of the foundation.
func NewQueryCmdProposals() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposals",
		Args:  cobra.NoArgs,
		Short: "Query all proposals",
		Long: `Query all proposals
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := foundation.NewQueryClient(clientCtx)

			req := foundation.QueryProposalsRequest{}
			res, err := queryClient.Proposals(context.Background(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// NewQueryCmdVote returns the vote of a voter on a proposal.
func NewQueryCmdVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote [proposal-id] [voter]",
		Args:  cobra.ExactArgs(2),
		Short: "Query the vote of a voter on a proposal",
		Long: `Query the vote of a voter on a proposal
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := foundation.NewQueryClient(clientCtx)

			proposalId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			voter := args[1]
			if err := sdk.ValidateAccAddress(voter); err != nil {
				return err
			}

			req := foundation.QueryVoteRequest{ProposalId: proposalId, Voter: voter}
			res, err := queryClient.Vote(context.Background(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// NewQueryCmdVotes returns the votes on a proposal.
func NewQueryCmdVotes() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "votes [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the votes on a proposal",
		Long: `Query the votes on a proposal
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := foundation.NewQueryClient(clientCtx)

			proposalId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			req := foundation.QueryVotesRequest{ProposalId: proposalId}
			res, err := queryClient.Votes(context.Background(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// NewQueryCmdTallyResult returns the tally of proposal votes.
func NewQueryCmdTallyResult() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tally [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the tally of proposal votes",
		Long: `Query the tally of proposal votes
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := foundation.NewQueryClient(clientCtx)

			proposalId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			req := foundation.QueryTallyResultRequest{ProposalId: proposalId}
			res, err := queryClient.TallyResult(context.Background(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
