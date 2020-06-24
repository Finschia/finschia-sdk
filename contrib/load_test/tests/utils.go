package tests

import "os"

func RemoveFile(path string) {
	if _, err := os.Stat(path); err != nil {
		return
	}
	if err := os.Remove(path); err != nil {
		panic(err)
	}
}
