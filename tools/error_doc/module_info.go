package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func findModuleName(s string) string {
	startIndex := strings.Index(s, "/x/") + len("/x/")
	endIndex := strings.Index(s[startIndex:], "/")

	if startIndex != -1 && endIndex != -1 {
		return s[startIndex : startIndex+endIndex]
	}
	return ""
}

func getModuleNameValue(filePath string) (string, error) {
	possibleFileNames := []string{"keys.go", "key.go"}
	var keyFilePath string
	for _, fileName := range possibleFileNames {
		paramPath := strings.Replace(filePath, "errors.go", fileName, 1)
		if _, err := os.Stat(paramPath); err == nil {
			keyFilePath = paramPath
			break
		}
	}

	if keyFilePath != "" {
		file, err := os.Open(keyFilePath)
		if err != nil {
			return "", errors.New(keyFilePath + " cannot be opened")
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			// get module name
			if strings.Contains(line, "ModuleName = ") {
				_, val, err := getConst(line)
				if err != nil {
					return "", err
				}
				return val, nil
			}
		}
	}

	return "", nil
}
