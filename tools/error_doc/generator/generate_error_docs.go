package generator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ErrorDocumentGenerator struct {
	targetPath    string
	errorsFiles   []string
	modules       []string
	errorDocument map[string][]*moduleInfo
}

type moduleInfo struct {
	filepath  string
	codespace string
	constDict map[string]string
	errorDict []errorInfo
}

func NewErrorDocumentGenerator(p string) *ErrorDocumentGenerator {
	return &ErrorDocumentGenerator{
		targetPath:    p,
		errorDocument: make(map[string][]*moduleInfo),
	}
}

func (edg *ErrorDocumentGenerator) extractModuleName() error {
	for _, filepath := range edg.errorsFiles {
		var moduleName string
		startIndex := strings.Index(filepath, "/x/") + len("/x/")
		endIndex := strings.Index(filepath[startIndex:], "/")
		if startIndex != -1 && endIndex != -1 {
			moduleName = filepath[startIndex : startIndex+endIndex]
		}
		if moduleName == "" {
			return errors.New("failed to get module name for " + filepath)
		}
		edg.errorDocument[moduleName] = append(edg.errorDocument[moduleName], &moduleInfo{
			filepath:  filepath,
			codespace: "",
			constDict: make(map[string]string),
			errorDict: []errorInfo{},
		})
	}
	// sort keys and filepaths
	for moduleName := range edg.errorDocument {
		edg.modules = append(edg.modules, moduleName)
		// sort.Strings(edg.errorDocument[moduleName])
	}
	sort.Strings(edg.modules)
	return nil
}

func (edg ErrorDocumentGenerator) outputCategory(file *os.File) {
	file.WriteString("<!-- TOC -->\n")
	file.WriteString("# Category\n")
	columnTemplate := "  * [%s](#%s)\n"
	for _, moduleName := range edg.modules {
		file.WriteString(fmt.Sprintf(columnTemplate, cases.Title(language.Und).String(moduleName), moduleName))
	}
	file.WriteString("<!-- TOC -->\n")
}

func (edg *ErrorDocumentGenerator) generateContent() error {
	// generate errors in each module
	for _, moduleName := range edg.modules {
		mods := edg.errorDocument[moduleName]
		for _, mod := range mods {
			if err := mod.errorsFileParse(); err != nil {
				return err
			}
			if err := mod.keysFileParse(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (edg ErrorDocumentGenerator) outputContent(file *os.File) error {
	extraInfoTemplate := "  * [%s](%s)\n"
	for _, moduleName := range edg.modules {
		// module name
		file.WriteString("\n")
		file.WriteString("## " + cases.Title(language.Und).String(moduleName) + "\n")
		// table header
		file.WriteString("\n")
		file.WriteString("|Error Name|Codespace|Code|Description|\n")
		file.WriteString("|:-|:-|:-|:-|\n")
		// table contents
		mods := edg.errorDocument[moduleName]
		for _, mod := range mods {
			for _, errInfo := range mod.errorDict {
				// assign value to field "codespace"
				if s, err := errInfo.toString(mod.codespace); err != nil {
					return err
				} else {
					file.WriteString(s)
				}
			}
		}
		// extract infomation
		file.WriteString("\n>You can also find detailed information in the following Errors.go files:\n")
		for _, mod := range mods {
			relPath, err := filepath.Rel(edg.targetPath, mod.filepath)
			if err != nil {
				return err
			}
			file.WriteString(fmt.Sprintf(extraInfoTemplate, relPath, relPath))
		}
	}
	return nil
}

func (edg ErrorDocumentGenerator) AutoGenerate() error {
	// get all errors.go in x folder
	errorsFileName := "errors.go"
	err := edg.listUpErrorsGoFiles(edg.targetPath, errorsFileName)
	if len(edg.errorsFiles) == 0 || err != nil {
		return errors.New("not find target files in x folder")
	}
	// get each module name and bind it to paths (one module may have multiple errors.go)
	if err := edg.extractModuleName(); err != nil {
		return err
	}
	// generate content
	if err := edg.generateContent(); err != nil {
		return err
	}
	// prepare the file for writing
	filepath := edg.targetPath + "/ERRORS.md"
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	// output category
	edg.outputCategory(file)
	// output content
	if err := edg.outputContent(file); err != nil {
		return err
	}
	return nil
}
