package types

type Slave struct {
	URL        string `json:"url"`
	Mnemonic   string `json:"mnemonic"`
	TargetType string `json:"target_type"`
}

func NewSlave(url, mnemonic, targetType string) Slave {
	return Slave{
		URL:        url,
		Mnemonic:   mnemonic,
		TargetType: targetType,
	}
}
