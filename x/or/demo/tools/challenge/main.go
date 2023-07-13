package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "challenge"
	app.Usage = "Challnge Structure tool"
	app.Description = "Challnge Structure tool"
	app.Commands = []*cli.Command{
		IDCommand,
		IDFromCommand,
		IsSearchingCommand,
		StepsCommand,
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			<-c
			cancel()
			fmt.Println("\r\nExiting...")
			os.Exit(0)
		}
	}()

	err := app.RunContext(ctx, os.Args)
	if err != nil {
		if errors.Is(err, ctx.Err()) {
			_, _ = fmt.Fprintf(os.Stderr, "command interrupted")
			os.Exit(130)
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "error: %v", err)
			os.Exit(1)
		}
	}
}
