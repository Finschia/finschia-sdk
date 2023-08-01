package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type errorInfo struct {
	errorName   string
	codeSpace   string
	code        string
	description string
}

func (err errorInfo) toString(moduleName string) (string, error) {
	errorInfoTemplate := "|%s|%s|%s|%s|\n"
	if err.codeSpace == "ModuleName" {
		if moduleName == "" {
			return "", errors.New("failed to find moduleName")
		} else {
			return fmt.Sprintf(errorInfoTemplate, err.errorName, moduleName, err.code, err.description), nil
		}
	} else {
		return fmt.Sprintf(errorInfoTemplate, err.errorName, err.codeSpace, err.code, err.description), nil
	}
}

func addError(line string, errorDict map[string]string) (errorInfo, error) {
	parts := strings.SplitN(line, "=", 2)
	errName := strings.TrimSpace(parts[0])
	errBody := strings.TrimSpace(parts[1])
	// error info is like as sdkerrors.Register(...)
	pattern := regexp.MustCompile(`sdkerrors\.Register\((.*)\)`)
	match := pattern.FindStringSubmatch(errBody)

	if len(match) == 2 {
		parts := strings.SplitN(match[1], ",", 3)

		if len(parts) == 3 {
			codeSpace := strings.TrimSpace(parts[0])
			code := strings.TrimSpace(parts[1])
			description := strings.Trim(strings.TrimSpace(parts[2]), `"`)

			if constValue, found := errorDict[codeSpace]; found {
				codeSpace = constValue
			}

			return errorInfo{
				errorName:   errName,
				codeSpace:   codeSpace,
				code:        code,
				description: description,
			}, nil
		} else {
			return errorInfo{}, errors.New("failed to get error info in: " + line)
		}
	} else {
		return errorInfo{}, errors.New("failed to parse error info in: " + line)
	}

}

func getErrors(p string) ([]errorInfo, error) {
	var errorDict []errorInfo
	constDict := make(map[string]string)

	file, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// get const
		if strings.Contains(line, "=") {
			if !strings.Contains(line, "sdkerrors.Register") {
				identifier, value, err := getConst(line)
				if err != nil {
					return nil, err
				}
				constDict[identifier] = value
			} else {
				errInfo, err := addError(line, constDict)
				if err != nil {
					return nil, err
				}
				errorDict = append(errorDict, errInfo)
			}
		}
	}
	return errorDict, nil
}
