package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/client/flags"
	"github.com/Finschia/finschia-sdk/version"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

const (
	FlagSequences = "sequences"
)

// NewQueryCmd returns the query commands for fbridge module
func NewQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the fbridge module",
	}

	cmd.AddCommand(
		NewQueryParamsCmd(),
		NewQueryNextSeqSendCmd(),
		NewQuerySeqToBlocknumsCmd(),
		NewQueryMembersCmd(),
		NewQueryMemberCmd(),
		NewQueryProposalsCmd(),
		NewQueryProposalCmd(),
		NewQueryVotesCmd(),
		NewQueryVoteCmd(),
		NewQueryBridgeStatusCmd(),
	)

	return cmd
}

func NewQueryParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current fbridge module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := types.NewQueryClient(clientCtx)
			res, err := qc.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryNextSeqSendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sending-next-seq",
		Short: "Query the next sequence number for sending",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := types.NewQueryClient(clientCtx)
			res, err := qc.NextSeqSend(cmd.Context(), &types.QueryNextSeqSendRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQuerySeqToBlocknumsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "seq-to-blocknums",
		Short:   "Query the block number for given sequence numbers",
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf("%s query %s seq-to-blocknums --sequences=1,2,3", version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := types.NewQueryClient(clientCtx)

			seqSlice, err := cmd.Flags().GetInt64Slice(FlagSequences)
			if err != nil {
				return err
			}

			seqs := make([]uint64, len(seqSlice))
			for i, seq := range seqSlice {
				seqs[i] = uint64(seq)
			}

			res, err := qc.SeqToBlocknums(cmd.Context(), &types.QuerySeqToBlocknumsRequest{Seqs: seqs})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().Int64Slice(FlagSequences, []int64{}, "comma separated list of bridge sequnece numbers")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryMembersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "members [role]",
		Short:   "Query the members of spcific group registered on the bridge (guardian|operator|judge)",
		Args:    cobra.ExactArgs(1),
		Example: fmt.Sprintf("%s query %s members guardian", version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := types.NewQueryClient(clientCtx)

			res, err := qc.Members(cmd.Context(), &types.QueryMembersRequest{Role: args[0]})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryMemberCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "member [address]",
		Short:   "Query the roles of a specific member registered on the bridge",
		Args:    cobra.ExactArgs(1),
		Example: fmt.Sprintf("%s query %s member link1...", version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := types.NewQueryClient(clientCtx)

			res, err := qc.Member(cmd.Context(), &types.QueryMemberRequest{Address: args[0]})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryProposalsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "proposals",
		Short:   "Query all role proposals",
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf("%s query %s proposals", version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := qc.Proposals(cmd.Context(), &types.QueryProposalsRequest{Pagination: pageReq})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all proposals")
	return cmd
}

func NewQueryProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "proposal [proposal_id]",
		Short:   "Query a specific role proposal",
		Args:    cobra.ExactArgs(1),
		Example: fmt.Sprintf("%s query %s proposal 1", version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := types.NewQueryClient(clientCtx)

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			res, err := qc.Proposal(cmd.Context(), &types.QueryProposalRequest{ProposalId: id})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryVotesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "votes [proposal_id]",
		Short:   "Query all votes for a specific role proposal",
		Args:    cobra.ExactArgs(1),
		Example: fmt.Sprintf("%s query %s votes 1", version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := types.NewQueryClient(clientCtx)

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			res, err := qc.Votes(cmd.Context(), &types.QueryVotesRequest{ProposalId: id})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryVoteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vote [proposal_id] [voter]",
		Short:   "Query a specific vote for a role proposal",
		Args:    cobra.ExactArgs(2),
		Example: fmt.Sprintf("%s query %s vote 1 link1...", version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := types.NewQueryClient(clientCtx)

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			res, err := qc.Vote(cmd.Context(), &types.QueryVoteRequest{ProposalId: id, Voter: args[1]})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func NewQueryBridgeStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Query the current status of the bridge",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := types.NewQueryClient(clientCtx)
			res, err := qc.BridgeStatus(cmd.Context(), &types.QueryBridgeStatusRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
