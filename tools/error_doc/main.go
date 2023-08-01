package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func findFilesWithName(startPath, fileName string) ([]string, error) {
	var foundFiles []string

	err := filepath.Walk(startPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == fileName {
			foundFiles = append(foundFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return foundFiles, nil
}

func findModuleWithFiles(targetPath string) (map[string][]string, []string, error) {

	// get all errors.go in x folder
	errorFile := "errors.go"
	filePaths, err := findFilesWithName(targetPath, errorFile)
	if len(filePaths) == 0 || err != nil {
		return nil, nil, errors.New("not find target files in x folder")
	}

	// get each module name and bind it to paths (one module may have multiple errors.go)
	moduleWithPaths := make(map[string][]string)
	for _, filePath := range filePaths {
		moduleName := findModuleName(filePath)
		if moduleName == "" {
			return nil, nil, errors.New("failed to get module name for " + filePath)
		}
		moduleWithPaths[moduleName] = append(moduleWithPaths[moduleName], filePath)
	}

	// sort keys and filepaths
	n := len(moduleWithPaths)
	modules := make([]string, 0, n)

	for moduleName := range moduleWithPaths {
		modules = append(modules, moduleName)
		sort.Strings(moduleWithPaths[moduleName])
	}
	sort.Strings(modules)

	return moduleWithPaths, modules, nil
}

func autoGenerate(targetPath string, moduleWithPaths map[string][]string, modules []string) error {
	filePath := targetPath + "/ERRORS.md"
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// generate category
	file.WriteString("<!-- TOC -->\n")
	file.WriteString("# Category\n")
	columnTemplate := "  * [%s](#%s)\n"
	for _, moduleName := range modules {
		file.WriteString(fmt.Sprintf(columnTemplate, cases.Title(language.Und).String(moduleName), moduleName))
	}
	file.WriteString("<!-- TOC -->\n")

	extraInfoTemplate := "  * [%s](%s)\n"
	// errors in each module
	for _, moduleName := range modules {

		// table header
		file.WriteString("\n")
		file.WriteString("## " + cases.Title(language.Und).String(moduleName) + "\n")
		file.WriteString("\n")
		file.WriteString("|Error Name|Codespace|Code|Description|\n")
		file.WriteString("|:-|:-|:-|:-|\n")

		filePaths := moduleWithPaths[moduleName]
		for _, filePath := range filePaths {
			errDict, err := getErrors(filePath)
			if err != nil {
				return err
			}
			moduleName, err := getModuleNameValue(filePath)
			if err != nil {
				return err
			}
			for _, errInfo := range errDict {
				column, err := errInfo.toString(moduleName)
				if err != nil {
					return err
				}
				file.WriteString(column)
			}
		}

		file.WriteString("\n>You can also find detailed information in the following Errors.go files:\n")
		for _, filePath := range filePaths {
			relPath, err := filepath.Rel(targetPath, filePath)
			if err != nil {
				return err
			}
			file.WriteString(fmt.Sprintf(extraInfoTemplate, relPath, relPath))
		}

	}
	return nil
}

func main() {
	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		os.Exit(1)
	}
	targetPath := filepath.Join(currentPath, "..", "..", "x")

	moduleWithPaths, modules, err := findModuleWithFiles(targetPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := autoGenerate(targetPath, moduleWithPaths, modules); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
