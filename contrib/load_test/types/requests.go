package types

type LoadRequest struct {
	Scenario string            `json:"scenario"`
	Config   Config            `json:"config"`
	Params   map[string]string `json:"params"`
}

func NewLoadRequest(scenario string, config Config, params map[string]string) *LoadRequest {
	return &LoadRequest{
		Scenario: scenario,
		Config:   config,
		Params:   params,
	}
}
