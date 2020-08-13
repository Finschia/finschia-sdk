package rest

//TODO : Integrate http status mapping for every REST API
import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

//HTTPStatusMappingTable is map to mapping an error type and a http status
type HTTPStatusMappingTable map[string]map[uint32]int

var (
	table = HTTPStatusMappingTable{
		errors.RootCodespace: {
			errors.ErrUnknownAddress.ABCICode(): http.StatusNotFound,
		},
	}
)

func RegisterHTTPStatusMapping(rawErr error, httpStatus int) {
	err := parsingError(rawErr)
	table[err.Codespace()][err.ABCICode()] = httpStatus
}

func parsingError(rawErr error) *errors.Error {
	if rawErr == nil {
		return nil
	}
	if err, ok := rawErr.(*context.Error); ok {
		return errors.New(err.Codespace, err.Code, err.Message)
	}
	if err, ok := rawErr.(*errors.Error); ok {
		return err
	}
	return errors.New(errors.UndefinedCodespace, 1, "internal")
}

//GetHTTPStatus is method to get http status for given error
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
