package generator

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func getCodeSpace(line string) (string, string, error) {
	line = strings.Replace(line, "const", "", 1)
	parts := strings.Split(line, "=")
	if len(parts) == 2 {
		i := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		return i, val, nil
	}
	return "", "", errors.New("failed to get the value in: " + line)
}

func (mi *moduleInfo) parseKeysFile() error {
	// find keys.go or key.go
	possibleFileNames := []string{"keys.go", "key.go"}
	var keyFilePath string
	for _, fileName := range possibleFileNames {
		paramPath := strings.Replace(mi.filepath, "errors.go", fileName, 1)
		if _, err := os.Stat(paramPath); err == nil {
			keyFilePath = paramPath
			break
		}
	}
	// if keys.go or key.go is exist
	if keyFilePath != "" {
		file, err := os.Open(keyFilePath)
		if err != nil {
			return errors.New(keyFilePath + " cannot be opened")
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			// get module name
			if strings.Contains(line, "ModuleName = ") {
				_, val, err := getCodeSpace(line)
				if err != nil {
					return err
				}
				mi.codespace = val
			}
		}
	}

	return nil
}
