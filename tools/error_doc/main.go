package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Finschia/finschia-sdk/tools/error_doc/generator"
)

func main() {
	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		os.Exit(1)
	}
	targetPath := filepath.Join(currentPath, "..", "..", "x")

	errorDocumentGenerator := generator.NewErrorDocumentGenerator(targetPath)
	if err := errorDocumentGenerator.AutoGenerate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
