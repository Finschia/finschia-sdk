package collection_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/x/collection"
)

func TestAttributeKeyStringer(t *testing.T) {
	for _, name := range collection.AttributeKey_name {
		t.Run(name, func(t *testing.T) {
			value := collection.AttributeKey(collection.AttributeKey_value[name])
			customName := value.String()
			require.EqualValues(t, value, collection.AttributeKeyFromString(customName))
		})
	}
}
