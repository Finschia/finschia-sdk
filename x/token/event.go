package token

import (
	"strings"
)

const (
	prefixAttributeKey = "ATTRIBUTE_KEY_"
)

func (x AttributeKey) String() string {
	lenPrefix := len(prefixAttributeKey)
	return strings.ToLower(AttributeKey_name[int32(x)][lenPrefix:])
}

func AttributeKeyFromString(name string) AttributeKey {
	attributeKeyName := prefixAttributeKey + strings.ToUpper(name)
	return AttributeKey(AttributeKey_value[attributeKeyName])
}
