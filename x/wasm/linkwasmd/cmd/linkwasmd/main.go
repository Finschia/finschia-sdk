package main

import (
	"os"

	"github.com/line/lbm-sdk/v2/server"
	svrcmd "github.com/line/lbm-sdk/v2/server/cmd"

	app "github.com/line/lbm-sdk/v2/x/wasm/linkwasmd/app"
	"github.com/line/lbm-sdk/v2/x/wasm/linkwasmd/cmd/linkwasmd/cmd"
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
