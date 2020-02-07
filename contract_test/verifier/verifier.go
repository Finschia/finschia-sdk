package verifier

import (
	"log"
)

func CompareJSONFormat(expected interface{}, actual interface{}) bool {
	switch v := expected.(type) {
	case []interface{}:
		array2, ok := actual.([]interface{})
		if !ok {
			log.Printf("Expected is an array but actual is not (expected: %v, actual: %v)", v, actual)
			return false
		}
		if v == nil || array2 == nil {
			log.Println("Array can not be null")
			return false
		}

		minLen := min(len(v), len(array2))
		for i := 0; i < minLen; i++ {
			if !CompareJSONFormat(v[i], array2[i]) {
				return false
			}
		}
		return true

	case map[string]interface{}:
		map2, ok := actual.(map[string]interface{})
		if !ok {
			log.Printf("Expected is an object but actual is not (expected: %v, actual: %v)", v, actual)
			return false
		}
		if v == nil || map2 == nil {
			log.Println("Object can not be null")
			return false
		}
		if len(v) != len(map2) {
			log.Printf("Objects have different size: (expected: %d, actual: %d)", len(v), len(map2))
			return false
		}

		for key, val1 := range v {
			if val2, ok := map2[key]; ok {
				if !CompareJSONFormat(val1, val2) {
					return false
				}
			} else {
				log.Printf("Expected has %v property but actual is not", key)
				return false
			}
		}
		return true

	default:
		if v == "" {
			return true
		}
		switch actual.(type) {
		case []interface{}:
			log.Printf("Expected is element type but actual is array (expected: %v, actual: %v)", v, actual)
			return false
		case map[string]interface{}:
			log.Printf("Expected is element type but actual is object (expected: %v, actual: %v)", v, actual)
			return false
		}
		return true
	}
}

func min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}
