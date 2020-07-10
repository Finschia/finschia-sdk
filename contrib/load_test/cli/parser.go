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
		slaves[i] = types.NewSlave(s["url"].(string), s["mnemonic"].(string), s["scenario"].(string))
		i++
	}
	return slaves
}
