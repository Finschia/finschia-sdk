package cli

import (
	"github.com/line/link/contrib/load_test/types"
	"github.com/spf13/viper"
)

func getSlaves() []types.Slave {
	slavesMap := viper.GetStringMap(FlagSlaves)
	slaves := make([]types.Slave, len(slavesMap))
	i := 0
	// The index of slaves and slave numbers do not match.
	for _, slaveMap := range slavesMap {
		s, ok := slaveMap.(map[string]interface{})
		if !ok {
			return nil
		}

		url, ok := s["url"].(string)
		if !ok {
			panic("url in config.yaml is not string")
		}
		mnemonic, ok := s["mnemonic"].(string)
		if !ok {
			panic("mnemonic in config.yaml is not string")
		}
		scenario, ok := s["scenario"].(string)
		if !ok {
			panic("scenario in config.yaml is not string")
		}
		slaves[i] = types.NewSlave(url, mnemonic, scenario, convertToStrings(s["params"].([]interface{})))
		i++
	}
	return slaves
}

func convertToStrings(array []interface{}) []string {
	var ok bool
	strings := make([]string, len(array))
	for i, v := range array {
		strings[i], ok = v.(string)
		if !ok {
			panic("There are params, not string types.")
		}
	}
	return strings
}
