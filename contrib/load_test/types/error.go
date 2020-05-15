package types

import "fmt"

type InaccessibleFieldError string

func (e InaccessibleFieldError) Error() string {
	return fmt.Sprintf("Inaccessible Field Error: %s", string(e))
}

type InvalidLoadParameterError string

func (e InvalidLoadParameterError) Error() string {
	return fmt.Sprintf("Invalid Load Parameter Error: %s", string(e))
}

type InvalidTargetTypeError string

func (e InvalidTargetTypeError) Error() string {
	return fmt.Sprintf("Invalid target Type Error: %s", string(e))
}

type InvalidMasterMnemonic struct {
	Mnemonic string
}

func (e InvalidMasterMnemonic) Error() string {
	return fmt.Sprintf("Invalid master mnemonic: %s", e.Mnemonic)
}

type InvalidMnemonic struct {
	Mnemonic string
}

func (e InvalidMnemonic) Error() string {
	return fmt.Sprintf("Invalid mnemonic: %s", e.Mnemonic)
}

type RequestFailed struct {
	URL    string
	Status string
	Body   []byte
}

func (e RequestFailed) Error() string {
	return fmt.Sprintf("Request failed: %s %s\n%s", e.URL, e.Status, e.Body)
}

type InvalidPacerTypeError struct {
	PacerType string
}

func (e InvalidPacerTypeError) Error() string {
	return fmt.Sprintf("Invalid pacer type: %s", e.PacerType)
}
