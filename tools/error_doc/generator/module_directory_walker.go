package generator

import (
	"os"
	"path/filepath"
)

func (edg *ErrorDocumentGenerator) listUpErrorsGoFiles(startPath, errorsFileName string) error {
	err := filepath.Walk(startPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == errorsFileName {
			edg.errorsFiles = append(edg.errorsFiles, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
