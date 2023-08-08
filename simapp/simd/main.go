package main

import (
	"os"

	"github.com/Finschia/finschia-rdk/server"
	svrcmd "github.com/Finschia/finschia-rdk/server/cmd"
	"github.com/Finschia/finschia-rdk/simapp"
	"github.com/Finschia/finschia-rdk/simapp/simd/cmd"
)

func main() {
	var _ = os.Args // To avoid linter "imported but not used" false positive
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
