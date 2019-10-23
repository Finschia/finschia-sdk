package k8s

import (
	"crypto/sha256"
	"fmt"
	"hash"
	. "strings"
	"time"
)

func ErrCheckedStrParam(param string, err error) string {
	if err != nil {
		panic(err)
	}
	return param
}
func ErrCheckedIntParam(param int, err error) int {
	if err != nil {
		panic(err)
	}
	return param
}

func ErrCheckedBoolParam(param bool, err error) bool {
	if err != nil {
		panic(err)
	}
	return param

}

func ErrCheckedStrArrParam(param []string, err error) []string {
	if err != nil {
		panic(err)
	}
	return param

}

func DefIfLTEZero(confField *string, format string, defVal int, setVal int) int {
	if !(setVal <= 0) {
		*confField = fmt.Sprintf(format, setVal)
		return setVal
	}
	*confField = fmt.Sprintf(format, defVal)
	return defVal
}

func DefIfEmpty(confField *string, defVal string, setVal string) string {
	if TrimSpace(setVal) != "" {
		*confField = setVal
		return setVal
	}
	*confField = defVal
	return defVal
}

func RandomHash() (hash.Hash, error) {
	hashKey := fmt.Sprintf("%d", time.Now().Second())
	sha256Hash := sha256.New()
	_, err := sha256Hash.Write([]byte(hashKey))
	if err != nil {
		return nil, err
	}
	return sha256Hash, nil
}
