package main

import (
	"errors"
	"strings"
)

func getConst(line string) (string, string, error) {
	line = strings.Replace(line, "const", "", 1)
	parts := strings.Split(line, "=")
	if len(parts) == 2 {
		i := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		return i, val, nil
	} else {
		return "", "", errors.New("failed to get the value in: " + line)
	}
}
