package main

import (
	"os"

	"github.com/line/lfb-sdk/server"
	svrcmd "github.com/line/lfb-sdk/server/cmd"

	app "github.com/line/lfb-sdk/x/wasm/linkwasmd/app"
	"github.com/line/lfb-sdk/x/wasm/linkwasmd/cmd/linkwasmd/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
