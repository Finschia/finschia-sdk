package main

import (
	"os"

	"github.com/line/lbm-sdk/server"
	svrcmd "github.com/line/lbm-sdk/server/cmd"

	"github.com/line/lbm-sdk/x/ibc/testing/simapp"
	"github.com/line/lbm-sdk/x/ibc/testing/simapp/simd/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, simapp.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
