package k8s

import (
	"fmt"
	"os"
	"strconv"
	"text/template"
)

type TemplateObject struct {
	templateOutputPath string
}

func NewTemplateObject(templateOutputPath string) *TemplateObject {
	return &TemplateObject{templateOutputPath}
}

func (to *TemplateObject) MakeFilebeatResultPath(chainId string, validatorOrder int) string {
	return fmt.Sprintf(to.templateOutputPath, chainId, validatorOrder)
}

func (to *TemplateObject) MakeTemplate(metaData *BuildMetaData, validatorOrder int) {
	sourceTemplateFilePath := metaData.filebeatTemplateFilePath
	sourceTemplateFile, _ := template.ParseFiles(sourceTemplateFilePath)

	destinationTemplateFilePath := to.MakeFilebeatResultPath(metaData.ChainID, validatorOrder)
	destinationTemplateFile, _ := os.Create(destinationTemplateFilePath)

	var configMap = make(map[string]string)
	configMap["logPath"] = fmt.Sprintf("/linkd/%s/node%d/linkd/linkd.log", metaData.ChainID, validatorOrder)
	configMap["validatorOrder"] = strconv.Itoa(validatorOrder)

	_ = sourceTemplateFile.Execute(destinationTemplateFile, configMap)
	_ = destinationTemplateFile.Close()
}
