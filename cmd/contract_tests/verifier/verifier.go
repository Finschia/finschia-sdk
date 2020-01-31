package verifier

import (
	"log"
)

func CompareJSONFormat(expected interface{}, actual interface{}) bool {
	switch v := expected.(type) {
	case []interface{}:
		array2, ok := actual.([]interface{})
		if !ok {
			go log.Println("One is an array but not the other")
			return false
		}
		if v == nil || array2 == nil {
			go log.Println("Array can not be null")
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
			go log.Println("One is an object but not the other")
			return false
		}
		if v == nil || map2 == nil {
			go log.Println("Object can not be null")
			return false
		}
		if len(v) != len(map2) {
			go log.Println("Objects have different size")
			return false
		}

		for key, val1 := range v {
			if val2, ok := map2[key]; ok {
				if !CompareJSONFormat(val1, val2) {
					return false
				}
			} else {
				go log.Println("Objects have different properties")
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
			go log.Println("One is element type but another is array")
			return false
		case map[string]interface{}:
			go log.Println("One is element type but another is object")
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
