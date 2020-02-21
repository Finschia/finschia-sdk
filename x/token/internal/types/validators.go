package types

import "unicode/utf8"

const TokenURIMaxLength = 1000

func ValidTokenURI(tokenURI string) bool {
	// must be shorter than 1000 UTF characters
	return utf8.RuneCountInString(tokenURI) < TokenURIMaxLength
}
