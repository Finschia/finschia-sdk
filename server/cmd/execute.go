package cmd

import (
	"context"

	ostcfg "github.com/line/ostracon/config"
	ostcli "github.com/line/ostracon/libs/cli"
	"github.com/spf13/cobra"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/client/flags"
	"github.com/Finschia/finschia-sdk/server"
)

const (
	envPrefix = "LBM"
)

// Execute executes the root command of an application. It handles creating a
// server context object with the appropriate server and client objects injected
// into the underlying stdlib Context. It also handles adding core CLI flags,
// specifically the logging flags. It returns an error upon execution failure.
func Execute(rootCmd *cobra.Command, defaultHome string) error {
	// Create and set a client.Context on the command's Context. During the pre-run
	// of the root command, a default initialized client.Context is provided to
	// seed child command execution with values such as AccountRetriver, Keyring,
	// and a Tendermint RPC. This requires the use of a pointer reference when
	// getting and setting the client.Context. Ideally, we utilize
	// https://github.com/spf13/cobra/pull/1118.
	srvCtx := server.NewDefaultContext()
	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &client.Context{})
	ctx = context.WithValue(ctx, server.ServerContextKey, srvCtx)

	rootCmd.PersistentFlags().String(flags.FlagLogLevel, ostcfg.DefaultPackageLogLevels(), "The logging level by modules (debug|info|error|none)")
	rootCmd.PersistentFlags().String(flags.FlagLogFormat, ostcfg.LogFormatPlain, "The logging format (json|plain)")

	executor := ostcli.PrepareBaseCmd(rootCmd, envPrefix, defaultHome)
	return executor.ExecuteContext(ctx)
}
