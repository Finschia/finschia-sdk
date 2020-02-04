package block

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestCommand(t *testing.T) {
	rootCmd := &cobra.Command{
		Use: "linkcli",
	}
	queryCmd := &cobra.Command{
		Use: "query",
	}
	rootCmd.AddCommand(queryCmd)

	blockCmd := Command(codec.New())
	queryCmd.AddCommand(blockCmd)

	tests := []struct {
		reason  string
		args    []string
		wantErr bool
	}{
		{"misspelled command", []string{"block_"}, true},
		{"no command provided", []string{}, false},
		{"help flag", []string{"block", "--help"}, false},
		{"shorthand help flag", []string{"block", "-h"}, false},
	}

	for _, tt := range tests {
		err := client.ValidateCmd(blockCmd, tt.args)

		require.Equal(t, tt.wantErr, err != nil, tt.reason)
	}
}
