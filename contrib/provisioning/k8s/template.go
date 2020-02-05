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
	sourceTemplateFile, err := template.ParseFiles(sourceTemplateFilePath)
	if err != nil {
		panic(err)
	}

	destinationTemplateFilePath := to.MakeFilebeatResultPath(metaData.ChainID, validatorOrder)
	destinationTemplateFile, err := os.Create(destinationTemplateFilePath)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := destinationTemplateFile.Close()
		fmt.Println(err)
	}()

	var configMap = make(map[string]string)
	configMap["logPath"] = fmt.Sprintf("/linkd/%s/node%d/linkd/linkd.log", metaData.ChainID, validatorOrder)
	configMap["validatorOrder"] = strconv.Itoa(validatorOrder)

	err = sourceTemplateFile.Execute(destinationTemplateFile, configMap)
	if err != nil {
		panic(err)
	}

}
