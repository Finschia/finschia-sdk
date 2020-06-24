package types

type Slave struct {
	URL      string `json:"url"`
	Mnemonic string `json:"mnemonic"`
	Scenario string `json:"scenario"`
}

func NewSlave(url, mnemonic, scenario string) Slave {
	return Slave{
		URL:      url,
		Mnemonic: mnemonic,
		Scenario: scenario,
	}
}
