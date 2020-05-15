package main

import (
	"fmt"
	"sort"

	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
)

func errorCodesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "errorcodes",
		Short: "List errorcodes registered",
		RunE: func(cmd *cobra.Command, _ []string) error {
			registeredErrors := errors.RegisteredErrors()
			keys := make([]string, 0, len(registeredErrors))
			errorCodesMap := map[string]*errors.Error{}
			for _, e := range registeredErrors {
				key := fmt.Sprintf("%s:%10d", e.Codespace(), e.ABCICode())
				keys = append(keys, key)
				errorCodesMap[key] = e
			}

			sort.Strings(keys)

			fmt.Println("| codespace | error code | description | ")
			fmt.Println("| --------- | ---------- | ----------- | ")
			for _, k := range keys {
				e := errorCodesMap[k]
				fmt.Printf("| %v | %v | %s |\n", e.Codespace(), e.ABCICode(), e.Error())
			}

			return nil
		},
	}
	return cmd
}
