package rest

// TODO : Intergrate http status mapping for every REST API
import (
	"net/http"

	"github.com/line/lfb-sdk/client"
	"github.com/line/lfb-sdk/types/errors"
)

// HTTPStatusMappingTable is map to mapping an error type and a http status
type HTTPStatusMappingTable map[string]map[uint32]int

var (
	table = HTTPStatusMappingTable{
		errors.RootCodespace: {
			9: http.StatusNotFound,
		},
	}
)

func parsingError(rawErr error) *errors.Error {
	if rawErr == nil {
		return nil
	}
	if err, ok := rawErr.(client.Error); ok {
		return errors.New(err.Codespace, err.Code, err.Message)
	}
	if err, ok := rawErr.(*errors.Error); ok {
		return err
	}
	return errors.New(errors.UndefinedCodespace, 1, "internal")
}

// GetHTTPStatus is method to get http status for given error
func GetHTTPStatusWithError(err error) int {
	abciErr := parsingError(err)
	if abciErr == nil {
		return http.StatusOK
	}
	result := http.StatusInternalServerError
	if codeTable, ok := table[abciErr.Codespace()]; ok {
		if status, ok := codeTable[abciErr.ABCICode()]; ok {
			result = status
		}
	}
	return result
}
