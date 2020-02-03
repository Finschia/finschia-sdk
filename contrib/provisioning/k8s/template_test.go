package k8s

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func MakeDummyData() BuildMetaData {
	dummyData := BuildMetaData{}
	dummyData.filebeatTemplateFilePath = "./filebeat-validator-template.yaml"
	dummyData.ChainID = "zero"

	return dummyData
}

func TestMakeTemplate(t *testing.T) {
	buildMetaData := MakeDummyData()
	validateOrder := 0

	expectedFileFormat := "./%s%d.yaml"
	expectedFileName := fmt.Sprintf(expectedFileFormat, buildMetaData.ChainID, validateOrder)

	to := NewTemplateObject(expectedFileFormat)
	to.MakeTemplate(&buildMetaData, validateOrder)

	defer func() {
		_ = os.Remove(expectedFileName)
	}()
	require.FileExists(t, expectedFileName)
}
