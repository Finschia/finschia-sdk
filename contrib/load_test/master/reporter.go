package master

import (
	"fmt"
	"io/ioutil"
)

// TODO: functions that write result to stdout will be added later
type Reporter struct {
	results [][]byte
}

func NewReporter(results [][]byte) *Reporter {
	return &Reporter{results: results}
}

func (r *Reporter) WriteResult(outputDir string) error {
	for i, result := range r.results {
		if err := ioutil.WriteFile(fmt.Sprintf("%s/result%d.bin", outputDir, i), result, Permission); err != nil {
			return err
		}
	}
	return nil
}
