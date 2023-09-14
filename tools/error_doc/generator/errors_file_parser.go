package generator

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
	codespace   string
	code        string
	description string
}

func (ei errorInfo) toString(cs string) (string, error) {
	errorInfoTemplate := "|%s|%s|%s|%s|\n"
	if ei.codespace == "ModuleName" {
		if cs == "" {
			return "", errors.New("failed to find moduleName")
		}
		ei.codespace = cs
	}
	return fmt.Sprintf(errorInfoTemplate, ei.errorName, ei.codespace, ei.code, ei.description), nil
}

func (ei *errorInfo) getError(line string, constDict map[string]string) error {
	parts := strings.SplitN(line, "=", 2)
	ei.errorName = strings.TrimSpace(parts[0])
	errBody := strings.TrimSpace(parts[1])
	// error info is like as sdkerrors.Register(...)
	pattern := regexp.MustCompile(`sdkerrors\.Register\((.*)\)`)
	match := pattern.FindStringSubmatch(errBody)
	if len(match) == 2 {
		parts := strings.SplitN(match[1], ",", 3)
		if len(parts) == 3 {
			ei.codespace = strings.TrimSpace(parts[0])
			ei.code = strings.TrimSpace(parts[1])
			ei.description = strings.Trim(strings.TrimSpace(parts[2]), `"`)
			if constValue, found := constDict[ei.codespace]; found {
				ei.codespace = constValue
			}
			return nil
		}
		return errors.New("failed to get error info in: " + line)
	}
	return errors.New("failed to parse error info in: " + line)
}

func getConst(line string) (string, string, error) {
	line = strings.Replace(line, "const", "", 1)
	parts := strings.Split(line, "=")
	if len(parts) == 2 {
		i := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		return i, val, nil
	}
	return "", "", errors.New("failed to get the value in: " + line)
}

func (mi *moduleInfo) parseErrorsFile() error {
	// var errorDict []errorInfo
	// constDict := make(map[string]string)
	file, err := os.Open(mi.filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=") {
			// get const
			if !strings.Contains(line, "sdkerrors.Register") {
				identifier, value, err := getConst(line)
				if err != nil {
					return err
				}
				mi.constDict[identifier] = value
			} else {
				// get error
				var errInfo errorInfo
				if err := errInfo.getError(line, mi.constDict); err != nil {
					return err
				}
				mi.errorDict = append(mi.errorDict, errInfo)
			}
		}
	}
	return nil
}
