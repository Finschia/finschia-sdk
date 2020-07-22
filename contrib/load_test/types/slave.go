package types

type Slave struct {
	URL      string   `json:"url"`
	Mnemonic string   `json:"mnemonic"`
	Scenario string   `json:"scenario"`
	Params   []string `json:"params"`
}

func NewSlave(url, mnemonic, scenario string, params []string) Slave {
	return Slave{
		URL:      url,
		Mnemonic: mnemonic,
		Scenario: scenario,
		Params:   params,
	}
}
