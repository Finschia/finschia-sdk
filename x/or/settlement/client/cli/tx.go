package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/Finschia/finschia-sdk/client"
	// "github.com/Finschia/finschia-sdk/client/flags"
	"github.com/Finschia/finschia-sdk/x/or/settlement/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(NewTxCmdStartChallenge())
	cmd.AddCommand(NewTxCmdNsectChallenge())
	cmd.AddCommand(NewTxCmdFinishChallenge())

	return cmd
}

func NewTxCmdStartChallenge() *cobra.Command {
	return &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			panic("implement me")
		},
	}
}

func NewTxCmdNsectChallenge() *cobra.Command {
	return &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			panic("implement me")
		},
	}
}

func NewTxCmdFinishChallenge() *cobra.Command {
	return &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			panic("implement me")
		},
	}
}
